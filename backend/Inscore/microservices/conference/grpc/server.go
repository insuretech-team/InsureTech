package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/service"
	"github.com/newage-saint/insuretech/backend/inscore/db"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/webrtc"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/webrtc/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// ServerConfig holds the configuration for the gRPC server
type ServerConfig struct {
	Port                string
	MaxConnectionIdle   time.Duration
	MaxConnectionAge    time.Duration
	MaxRecvMsgSize      int
	MaxSendMsgSize      int
	EnableReflection    bool
	EnableKeepalive     bool
	KeepaliveInterval   time.Duration
	KeepaliveTimeout    time.Duration
}

// DefaultServerConfig returns a default server configuration
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:                "50052",
		MaxConnectionIdle:   5 * time.Minute,
		MaxConnectionAge:    30 * time.Minute,
		MaxRecvMsgSize:      10 * 1024 * 1024, // 10MB
		MaxSendMsgSize:      10 * 1024 * 1024, // 10MB
		EnableReflection:    true,
		EnableKeepalive:     true,
		KeepaliveInterval:   30 * time.Second,
		KeepaliveTimeout:    10 * time.Second,
	}
}

// Server wraps the gRPC server with WebRTC services
type Server struct {
	config           *ServerConfig
	grpcServer       *grpc.Server
	listener         net.Listener
	
	// WebRTC services
	roomService      *services.RoomService
	peerService      *services.PeerService
	trackService     *services.TrackService
	signalingService *services.SignalingService
	statsService     *services.StatsService
}

// NewServer creates a new gRPC server with WebRTC services
func NewServer(config *ServerConfig) (*Server, error) {
	if config == nil {
		config = DefaultServerConfig()
	}

	// Create gRPC server options
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(config.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(config.MaxSendMsgSize),
		grpc.Creds(insecure.NewCredentials()),
	}

	// Add keepalive parameters
	if config.EnableKeepalive {
		opts = append(opts,
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionIdle: config.MaxConnectionIdle,
				MaxConnectionAge:  config.MaxConnectionAge,
				Time:              config.KeepaliveInterval,
				Timeout:           config.KeepaliveTimeout,
			}),
			grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
				MinTime:             10 * time.Second,
				PermitWithoutStream: true,
			}),
		)
	}

	grpcServer := grpc.NewServer(opts...)

	// Get database connection
	database := db.GetDB()
	if database == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	// Create repository and services using ServiceManager
	manager := webrtc.NewServiceManager(database)
	roomService := manager.GetRoomService()
	peerService := manager.GetPeerService()
	trackService := manager.GetTrackService()
	signalingService := manager.GetSignalingService()
	statsService := manager.GetStatsService()

	server := &Server{
		config:           config,
		grpcServer:       grpcServer,
		roomService:      roomService,
		peerService:      peerService,
		trackService:     trackService,
		signalingService: signalingService,
		statsService:     statsService,
	}

	// Register services
	server.registerServices()

	return server, nil
}

// registerServices registers all WebRTC services with the gRPC server
func (s *Server) registerServices() {
	service.RegisterRoomServiceServer(s.grpcServer, s.roomService)
	service.RegisterPeerServiceServer(s.grpcServer, s.peerService)
	service.RegisterTrackServiceServer(s.grpcServer, s.trackService)
	service.RegisterSignalingServiceServer(s.grpcServer, s.signalingService)
	service.RegisterStatsServiceServer(s.grpcServer, s.statsService)

	if s.config.EnableReflection {
		reflection.Register(s.grpcServer)
	}

	appLogger.Info("Registered WebRTC gRPC services: RoomService, PeerService, TrackService, SignalingService, StatsService")
}

// Start starts the gRPC server
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", s.config.Port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", s.config.Port, err)
	}
	s.listener = listener

	appLogger.Infof("Conference gRPC server listening on port %s", s.config.Port)
	appLogger.Infof("Server configuration: MaxRecvMsgSize=%d, MaxSendMsgSize=%d, Keepalive=%v",
		s.config.MaxRecvMsgSize, s.config.MaxSendMsgSize, s.config.EnableKeepalive)

	return s.grpcServer.Serve(listener)
}

// Stop gracefully stops the gRPC server
func (s *Server) Stop() {
	appLogger.Info("Shutting down gRPC server...")
	s.grpcServer.GracefulStop()
	if s.listener != nil {
		s.listener.Close()
	}
	appLogger.Info("gRPC server stopped")
}

// GetGRPCServer returns the underlying gRPC server
func (s *Server) GetGRPCServer() *grpc.Server {
	return s.grpcServer
}

// Health check implementation
func (s *Server) HealthCheck(ctx context.Context) error {
	// Check database connection
	database := db.GetDB()
	if database == nil {
		return fmt.Errorf("database connection is nil")
	}

	// Ping database
	sqlDB, err := database.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}
