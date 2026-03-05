package main

import (
	"database/sql"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/docgen"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"google.golang.org/grpc"
)

func main() {
	// Initialize database connection (placeholder - implement based on actual DB setup)
	db, err := setupDatabase()
	if err != nil {
		logger.Fatalf("failed to setup database: %v", err)
	}
	defer db.Close()

	// Initialize storage service connection (optional)
	var storageConn *grpc.ClientConn
	// Storage connection setup would go here if needed

	// Create DocGen server with Kafka publisher and async worker
	docgenServer, err := docgen.NewDocumentServer(db, storageConn)
	if err != nil {
		logger.Fatalf("failed to create docgen server: %v", err)
	}

	// Setup gRPC server
	grpcServer := grpc.NewServer()
	docgenServer.Handler() // Handler is available but not registered here - registration would happen in caller

	// Start listening on gRPC port
	listener, err := net.Listen("tcp", ":50051") // Use appropriate port
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Run server in a goroutine
	go func() {
		logger.Info("Starting DocGen gRPC server")
		if err := grpcServer.Serve(listener); err != nil {
			logger.Errorf("gRPC server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	sig := <-sigChan
	logger.Infof("Received signal: %v, shutting down gracefully", sig)

	// Graceful shutdown: close DocGen server resources (Kafka, async worker)
	if err := docgenServer.Close(); err != nil {
		logger.Errorf("error closing docgen server: %v", err)
	}

	// Stop gRPC server gracefully
	grpcServer.GracefulStop()
	logger.Info("DocGen server stopped")
}

// setupDatabase initializes the database connection
// This is a placeholder - implement based on your actual database setup
func setupDatabase() (*sql.DB, error) {
	// Example: return sql.Open("postgres", "connection_string")
	// For now, this is a stub that should be implemented
	logger.Warnf("setupDatabase is not fully implemented - using stub")
	return nil, nil
}
