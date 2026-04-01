package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	notificationservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/notification/services/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

type Config struct {
	Host string
	Port int
	DB   *gorm.DB
}

func DefaultServerConfig() *Config {
	return &Config{
		Host: "0.0.0.0",
		Port: 50230,
	}
}

type Server struct {
	server *grpc.Server
	config *Config
	health *health.Server
}

func NewServer(config *Config, service NotificationServiceIface) (*Server, error) {
	if service == nil {
		return nil, errors.New("notification service is required")
	}

	grpcServer := grpc.NewServer()
	healthServer := health.NewServer()
	s := &Server{
		server: grpcServer,
		config: config,
		health: healthServer,
	}

	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	notificationservicev1.RegisterNotificationServiceServer(grpcServer, NewNotificationServiceHandler(service))
	reflection.Register(grpcServer)
	return s, nil
}

func (s *Server) Start() error {
	host := "0.0.0.0"
	if s.config != nil && strings.TrimSpace(s.config.Host) != "" {
		host = s.config.Host
	}
	port := 50230
	if s.config != nil && s.config.Port > 0 {
		port = s.config.Port
	}

	addr := net.JoinHostPort(host, strconv.Itoa(port))
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}

	appLogger.Infof("notification gRPC server listening on %s", addr)
	s.health.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	s.health.SetServingStatus("insuretech.notification.services.v1.NotificationService", grpc_health_v1.HealthCheckResponse_SERVING)
	return s.server.Serve(lis)
}

func (s *Server) Stop() {
	s.health.Shutdown()
	s.server.GracefulStop()
}

func (s *Server) HealthCheck(ctx context.Context) error {
	dbConn := s.config.DB
	if dbConn == nil {
		dbConn = db.GetDB()
	}
	if dbConn == nil {
		return errors.New("database connection is nil")
	}
	sqlDB, err := dbConn.DB()
	if err != nil {
		return fmt.Errorf("get sql db: %w", err)
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}
	return nil
}
