package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/webrtc/core"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/webrtc/mediaserver"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/webrtc/repository"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/webrtc/services"
	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/service"
	"github.com/newage-saint/insuretech/ops/env"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	webrtc "github.com/pion/webrtc/v4"
)

// Server represents the WebRTC gRPC server
type Server struct {
	config       *core.ProductionConfig
	grpcServer   *grpc.Server
	db           *gorm.DB
	sfu          *mediaserver.SFU
	logger       *zap.Logger
	healthServer *health.Server

	// Services
	roomService      *services.RoomService
	peerService      *services.PeerService
	trackService     *services.TrackService
	signalingService *services.SignalingServiceWithSFU
	statsService     *services.StatsService
}

func main() {
	// Load environment variables
	_ = env.Load()
	
	// Load configuration
	config := core.LoadFromEnv()

	// Initialize logger
	logger, err := initLogger(config.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Startinginsuretech WebRTC Conference Server",
		zap.String("version", "1.0.0"),
		zap.Bool("production", config.IsProduction()),
		zap.String("host", config.Host),
		zap.Int("port", config.Port),
	)

	// Validate configuration
	if err := config.Validate(); err != nil {
		logger.Fatal("Invalid configuration", zap.Error(err))
	}

	// Initialize server
	srv, err := NewServer(config, logger)
	if err != nil {
		logger.Fatal("Failed to create server", zap.Error(err))
	}

	// Start server
	if err := srv.Start(); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}

	// Wait for shutdown signal
	srv.WaitForShutdown()
}

// NewServer creates a new WebRTC server instance
func NewServer(config *core.ProductionConfig, logger *zap.Logger) (*Server, error) {
	srv := &Server{
		config: config,
		logger: logger,
	}

	// Initialize database if enabled
	if config.DatabaseEnabled {
		db, err := initDatabase(logger)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize database: %w", err)
		}
		srv.db = db
		logger.Info("Database connected successfully")
	}

	// Initialize SFU media server
	sfuConfig := &mediaserver.Config{
		ICEServers:      convertICEServersToWebRTC(config.ICEServers),
		MaxBitrate:      uint64(config.Media.MaxVideoBitrate) * 1000,
		EnableSimulcast: true,
		Debug:           config.LogLevel == "debug",
	}

	sfu, err := mediaserver.NewSFU(sfuConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SFU: %w", err)
	}
	srv.sfu = sfu
	logger.Info("SFU media server initialized")

	// Initialize repository
	repo := repository.NewRepository(srv.db)

	// Initialize services
	srv.roomService = services.NewRoomService(repo)
	srv.peerService = services.NewPeerService(repo)
	srv.trackService = services.NewTrackService(repo)
	srv.statsService = services.NewStatsService(repo)

	// Use SFU-integrated signaling service
	signalingServiceSFU := services.NewSignalingServiceWithSFU(repo, sfu)
	srv.signalingService = signalingServiceSFU

	logger.Info("All services initialized")

	// Initialize gRPC server
	if err := srv.initGRPCServer(); err != nil {
		return nil, fmt.Errorf("failed to initialize gRPC server: %w", err)
	}

	return srv, nil
}

// initGRPCServer initializes the gRPC server with all services
func (s *Server) initGRPCServer() error {
	var opts []grpc.ServerOption

	// Add TLS credentials if enabled
	if s.config.TLS.Enabled {
		creds, err := credentials.NewServerTLSFromFile(
			s.config.TLS.CertFile,
			s.config.TLS.KeyFile,
		)
		if err != nil {
			return fmt.Errorf("failed to load TLS credentials: %w", err)
		}
		opts = append(opts, grpc.Creds(creds))
		s.logger.Info("TLS enabled",
			zap.String("cert", s.config.TLS.CertFile),
			zap.String("key", s.config.TLS.KeyFile),
		)
	}

	// Add interceptors for logging, auth, rate limiting
	opts = append(opts,
		grpc.ChainUnaryInterceptor(
			s.loggingInterceptor(),
			s.authInterceptor(),
			s.rateLimitInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			s.streamLoggingInterceptor(),
			s.streamAuthInterceptor(),
		),
		grpc.MaxConcurrentStreams(1000),
		grpc.MaxRecvMsgSize(10*1024*1024), // 10MB
		grpc.MaxSendMsgSize(10*1024*1024),
	)

	// Create gRPC server
	s.grpcServer = grpc.NewServer(opts...)

	// Register services
	service.RegisterRoomServiceServer(s.grpcServer, s.roomService)
	service.RegisterPeerServiceServer(s.grpcServer, s.peerService)
	service.RegisterTrackServiceServer(s.grpcServer, s.trackService)
	service.RegisterSignalingServiceServer(s.grpcServer, s.signalingService)
	service.RegisterStatsServiceServer(s.grpcServer, s.statsService)

	// Register health check service
	s.healthServer = health.NewServer()
	grpc_health_v1.RegisterHealthServer(s.grpcServer, s.healthServer)
	s.healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Register reflection for grpcurl/grpc-ui
	reflection.Register(s.grpcServer)

	s.logger.Info("gRPC services registered",
		zap.Strings("services", []string{
			"RoomService",
			"PeerService",
			"TrackService",
			"SignalingService",
			"StatsService",
			"Health",
		}),
	)

	return nil
}

// Start starts the gRPC server
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.config.Address())
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.logger.Info("Starting gRPC server",
		zap.String("address", s.config.Address()),
		zap.Bool("tls", s.config.TLS.Enabled),
	)

	// Start server in goroutine
	go func() {
		if err := s.grpcServer.Serve(listener); err != nil {
			s.logger.Fatal("gRPC server failed", zap.Error(err))
		}
	}()

	s.logger.Info("✓ WebRTC Conference Server is running",
		zap.String("address", s.config.Address()),
		zap.Bool("production", s.config.IsProduction()),
	)

	return nil
}

// WaitForShutdown waits for shutdown signal and performs graceful shutdown
func (s *Server) WaitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	s.logger.Info("Received shutdown signal", zap.String("signal", sig.String()))

	// Graceful shutdown
	s.Shutdown()
}

// Shutdown performs graceful shutdown
func (s *Server) Shutdown() {
	s.logger.Info("Shutting down server...")

	// Set health check to NOT_SERVING
	s.healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	// Stop accepting new connections
	stopped := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(stopped)
	}()

	// Wait for graceful stop with timeout
	select {
	case <-stopped:
		s.logger.Info("Server stopped gracefully")
	case <-time.After(30 * time.Second):
		s.logger.Warn("Graceful shutdown timeout, forcing stop")
		s.grpcServer.Stop()
	}

	// Close database connection
	if s.db != nil {
		sqlDB, err := s.db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	s.logger.Info("Shutdown complete")
}

// Helper functions

func initLogger(level string) (*zap.Logger, error) {
	var config zap.Config

	if level == "debug" {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	return config.Build()
}

func initDatabase(logger *zap.Logger) (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=lifepluscore port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func convertICEServersToWebRTC(configs []core.ICEServerConfig) []webrtc.ICEServer {
	var servers []webrtc.ICEServer
	for _, cfg := range configs {
		server := webrtc.ICEServer{
			URLs: cfg.URLs,
		}
		if cfg.Username != "" {
			server.Username = cfg.Username
		}
		if cfg.Credential != "" {
			server.Credential = cfg.Credential
		}
		servers = append(servers, server)
	}
	return servers
}

// Interceptors

func (s *Server) loggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			s.logger.Error("gRPC call failed",
				zap.String("method", info.FullMethod),
				zap.Duration("duration", duration),
				zap.Error(err),
			)
		} else {
			s.logger.Debug("gRPC call",
				zap.String("method", info.FullMethod),
				zap.Duration("duration", duration),
			)
		}

		return resp, err
	}
}

func (s *Server) authInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip auth for health check
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(ctx, req)
		}

		// TODO: Implement JWT validation if enabled
		if s.config.Security.JWTEnabled {
			// Validate JWT token from metadata
			// Return unauthenticated error if invalid
		}

		return handler(ctx, req)
	}
}

func (s *Server) rateLimitInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// TODO: Implement rate limiting if enabled
		if s.config.Security.RateLimitEnabled {
			// Check rate limit
			// Return resource exhausted error if limit exceeded
		}

		return handler(ctx, req)
	}
}

func (s *Server) streamLoggingInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		s.logger.Info("Stream opened", zap.String("method", info.FullMethod))
		err := handler(srv, ss)
		if err != nil {
			s.logger.Error("Stream error", zap.String("method", info.FullMethod), zap.Error(err))
		}
		return err
	}
}

func (s *Server) streamAuthInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// TODO: Implement JWT validation for streams
		return handler(srv, ss)
	}
}
