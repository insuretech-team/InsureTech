package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/insurance/service"
	insurancev1 "github.com/newage-saint/insuretech/gen/go/insuretech/insurance/services/v1"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/ops/config"
)

func main() {
	// Initialize logger
	_ = logger.Initialize(logger.Config{
		Level:  "info",
		Format: "text",
		Output: "console",
	})
	logger.Info("Starting Insurance Service...")

	// Load database config
	dbConfigPath, err := config.ResolveConfigPath("database.yaml")
	if err != nil {
		logger.Fatalf("Failed to resolve database config path: %v", err)
	}
	if err := db.InitializeManagerForService(dbConfigPath); err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Manager.Close()

	// Get database connection
	database := db.GetDB()

	// Create insurance service
	insuranceService := service.NewInsuranceService(database)

	// Setup gRPC server
	grpcPort := os.Getenv("INSURANCE_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50115"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		logger.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer()
	insurancev1.RegisterInsuranceServiceServer(grpcServer, insuranceService)
	
	// Enable reflection for grpcurl
	reflection.Register(grpcServer)

	logger.Infof("Insurance Service starting on gRPC port %s", grpcPort)

	// Start server in goroutine
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Insurance Service...")
	grpcServer.GracefulStop()
}
