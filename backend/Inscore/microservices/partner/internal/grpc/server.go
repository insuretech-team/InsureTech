package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/partner/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/interceptors"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	partnerservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/services/v1"
)

// Config holds server configuration
type Config struct {
	Port string
	DB   *gorm.DB
}

// DefaultServerConfig returns default configuration
func DefaultServerConfig() *Config {
	return &Config{
		Port: "50058",
	}
}

// Server represents the partner gRPC server
type Server struct {
	server  *grpc.Server
	config  *Config
	health  *health.Server
	handler *PartnerHandler
}

// NewServer creates a new gRPC server
func NewServer(config *Config, partnerService domain.PartnerService, authzClient authzservicev1.AuthZServiceClient) (*Server, error) {
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptors.NewAuthZInterceptor(authzClient, interceptors.DefaultSkipMethods),
		),
	}
	grpcServer := grpc.NewServer(opts...)
	healthServer := health.NewServer()
	handler := NewPartnerHandler(partnerService)

	s := &Server{
		server:  grpcServer,
		config:  config,
		health:  healthServer,
		handler: handler,
	}

	s.registerServices()
	return s, nil
}

func (s *Server) registerServices() {
	grpc_health_v1.RegisterHealthServer(s.server, s.health)
	partnerservicev1.RegisterPartnerServiceServer(s.server, s.handler)
	reflection.Register(s.server)
}

// Start starts the server
func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.config.Port))
	if err != nil {
		appLogger.Errorf("failed to listen on port %s: %v", s.config.Port, err)
		return fmt.Errorf("failed to listen on port %s: %w", s.config.Port, err)
	}

	appLogger.Infof("gRPC server listening on port %s", s.config.Port)

	s.health.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	s.health.SetServingStatus("insuretech.partner.services.v1.PartnerService", grpc_health_v1.HealthCheckResponse_SERVING)

	if err := s.server.Serve(lis); err != nil {
		appLogger.Errorf("failed to serve: %v", err)
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

// Stop stops the server
func (s *Server) Stop() {
	s.health.Shutdown()
	s.server.GracefulStop()
}

// HealthCheck performs a database ping
func (s *Server) HealthCheck(ctx context.Context) error {
	if s.config.DB == nil {
		return errors.New("database connection is nil")
	}
	db, err := s.config.DB.DB()
	if err != nil {
		return errors.New("failed to get sql db")
	}
	if err := db.PingContext(ctx); err != nil {
		return errors.New("database ping failed")
	}
	return nil
}
