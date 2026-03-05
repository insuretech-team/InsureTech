package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/service"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
)

// Config holds server configuration
type Config struct {
	Host string
	Port string
	DB   *gorm.DB
}

// DefaultServerConfig returns default configuration
func DefaultServerConfig() *Config {
	return &Config{
		Host: "0.0.0.0",
		Port: "50053",
	}
}

// Server represents the auth gRPC server
type Server struct {
	server      *grpc.Server
	config      *Config
	health      *health.Server
	authService *service.AuthService
}

// NewServer creates a new gRPC server
func NewServer(config *Config, authService *service.AuthService) (*Server, error) {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(defaultUnaryInterceptors()...),
		grpc.ChainStreamInterceptor(defaultStreamInterceptors()...),
	)
	healthServer := health.NewServer()

	s := &Server{
		server:      grpcServer,
		config:      config,
		health:      healthServer,
		authService: authService,
	}

	s.registerServices()
	return s, nil
}

func (s *Server) registerServices() {
	grpc_health_v1.RegisterHealthServer(s.server, s.health)

	// Use the handler factory/constructor
	authHandler := NewAuthServiceHandler(s.authService)
	authnservicev1.RegisterAuthServiceServer(s.server, authHandler)

	reflection.Register(s.server)
}

// Start starts the server
func (s *Server) Start() error {
	host := s.config.Host
	if host == "" {
		host = "0.0.0.0"
	}
	addr := net.JoinHostPort(host, s.config.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Errorf("failed to listen on %s: %v", addr, err)
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	appLogger.Infof("gRPC server listening on %s", addr)

	s.health.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	s.health.SetServingStatus("insuretech.authn.services.v1.AuthService", grpc_health_v1.HealthCheckResponse_SERVING)

	if err := s.server.Serve(lis); err != nil {
		logger.Errorf("failed to serve: %v", err)
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

// Stop stops the server
func (s *Server) Stop() {
	s.health.Shutdown()
	s.server.GracefulStop()
}

// HealthCheck performs a self-check
func (s *Server) HealthCheck(ctx context.Context) error {
	dbConn := s.config.DB
	if dbConn == nil {
		// Fallback to global manager-backed DB in case caller omitted static DB pointer.
		dbConn = db.GetDB()
	}
	if dbConn == nil {
		logger.Errorf("database connection is nil")
		return errors.New("database connection is nil")
	}

	sqlDB, err := dbConn.DB()
	if err != nil {
		logger.Errorf("failed to get sql db: %v", err)
		return fmt.Errorf("failed to get sql db: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		logger.Errorf("database ping failed: %v", err)
		return fmt.Errorf("database ping failed: %w", err)
	}
	return nil
}
