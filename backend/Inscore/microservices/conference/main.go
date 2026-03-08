package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	conferenceGrpc "github.com/newage-saint/insuretech/backend/inscore/microservices/conference/grpc"
)

const (
	defaultPort = "50052"
	serviceName = "conference"
)

func main() {
	// Initialize logger with default config
	_ = appLogger.Initialize(appLogger.Config{
		Level:  "info",
		Format: "text",
		Output: "console",
	})

	appLogger.Info("Starting Conference microservice...")

	// Get port from environment or use default
	port := os.Getenv("CONFERENCE_PORT")
	if port == "" {
		port = defaultPort
	}

	// Initialize database connection
	configPath := "inscore/configs/database.yaml"
	if err := db.InitializeManagerForService(configPath); err != nil {
		appLogger.Errorf("Failed to initialize database: %v", err)
		appLogger.Fatal("Database initialization failed")
	}
	defer func() {
		if db.Manager != nil {
			if err := db.Manager.Close(); err != nil {
				appLogger.Errorf("Failed to close database: %v", err)
			}
		}
	}()

	// Verify database connection
	database := db.GetDB()
	if database == nil {
		appLogger.Fatal("Database connection is nil")
	}
	appLogger.Info("Database connection established")

	// Create server configuration
	config := conferenceGrpc.DefaultServerConfig()
	config.Port = port

	// Create and configure gRPC server
	server, err := conferenceGrpc.NewServer(config)
	if err != nil {
		appLogger.Errorf("Failed to create gRPC server: %v", err)
		appLogger.Fatal("Server initialization failed")
	}

	// Perform health check
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.HealthCheck(ctx); err != nil {
		appLogger.Errorf("Health check failed: %v", err)
		appLogger.Fatal("Server health check failed")
	}
	appLogger.Info("Server health check passed")

	// Start server in goroutine
	go func() {
		if err := server.Start(); err != nil {
			appLogger.Errorf("Failed to serve: %v", err)
			appLogger.Fatal("Server crashed")
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	appLogger.Info("Shutting down Conference microservice...")
	server.Stop()
	appLogger.Info("Conference microservice stopped gracefully")
}
