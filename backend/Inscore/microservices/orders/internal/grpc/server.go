package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/orders/internal/domain"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	orderservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/orders/services/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

type Config struct {
	Port string
	DB   *gorm.DB
}

func DefaultServerConfig() *Config {
	return &Config{Port: "50142"}
}

type Server struct {
	server  *grpc.Server
	config  *Config
	health  *health.Server
	handler *OrderHandler
}

func NewServer(config *Config, svc domain.OrderService, interceptors ...grpc.ServerOption) (*Server, error) {
	grpcServer := grpc.NewServer(interceptors...)
	healthServer := health.NewServer()
	handler := NewOrderHandler(svc)

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
	orderservicev1.RegisterOrderServiceServer(s.server, s.handler)
	reflection.Register(s.server)
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.config.Port))
	if err != nil {
		appLogger.Errorf("failed to listen on port %s: %v", s.config.Port, err)
		return fmt.Errorf("failed to listen on port %s: %w", s.config.Port, err)
	}

	appLogger.Infof("orders gRPC server listening on port %s", s.config.Port)
	s.health.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	s.health.SetServingStatus("insuretech.orders.services.v1.OrderService", grpc_health_v1.HealthCheckResponse_SERVING)

	if err := s.server.Serve(lis); err != nil {
		appLogger.Errorf("failed to serve: %v", err)
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}

func (s *Server) Stop() {
	s.health.Shutdown()
	s.server.GracefulStop()
}

func (s *Server) HealthCheck(ctx context.Context) error {
	if s.config.DB == nil {
		return errors.New("database connection is nil")
	}
	sqlDB, err := s.config.DB.DB()
	if err != nil {
		return errors.New("failed to get sql db")
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return errors.New("database ping failed")
	}
	return nil
}
