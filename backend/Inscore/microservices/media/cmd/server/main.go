package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/media"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/ops/config"
	mediaservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/media/services/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"
)

// ServicesConfig structure matches services.yaml
type ServicesConfig struct {
	Services map[string]struct {
		Name  string `yaml:"name"`
		Ports struct {
			Grpc int `yaml:"grpc"`
			Http int `yaml:"http"`
		} `yaml:"ports"`
	} `yaml:"services"`
}

func main() {
	// 1. Initialize Logger
	_ = appLogger.Initialize(appLogger.Config{
		Level:  "info",
		Format: "text",
		Output: "console",
	})
	appLogger.Info("Starting Media microservice...")

	// 2. Load Configuration (Port)
	servicesConfigPath, err := config.ResolveConfigPath("services.yaml")
	if err != nil {
		appLogger.Fatalf("Failed to resolve services.yaml path: %v", err)
	}

	servicesData, err := os.ReadFile(servicesConfigPath)
	if err != nil {
		appLogger.Fatalf("Failed to read services.yaml: %v", err)
	}

	var svcConfig ServicesConfig
	if err := yaml.Unmarshal(servicesData, &svcConfig); err != nil {
		appLogger.Fatalf("Failed to parse services.yaml: %v", err)
	}

	mediaConfig, exists := svcConfig.Services["media"]
	if !exists {
		appLogger.Fatal("Configuration for 'media' service not found in services.yaml")
	}
	port := strconv.Itoa(mediaConfig.Ports.Grpc)
	if os.Getenv("MEDIA_PORT") != "" || os.Getenv("MEDIA_GRPC_PORT") != "" || os.Getenv("MEDIA_HTTP_PORT") != "" {
		appLogger.Warn("MEDIA_PORT/MEDIA_GRPC_PORT/MEDIA_HTTP_PORT env values are ignored; using backend/inscore/configs/services.yaml")
	}

	appLogger.Infof("Service '%s' configured on port %s", mediaConfig.Name, port)

	// 3. Initialize Infrastructure (DB)
	dbConfigPath, err := config.ResolveConfigPath("database.yaml")
	if err != nil {
		appLogger.Fatalf("Failed to resolve database config path: %v", err)
	}
	if err := db.InitializeManagerForService(dbConfigPath); err != nil {
		appLogger.Errorf("Failed to initialize database: %v", err)
		appLogger.Fatal("Database initialization failed")
	}
	defer db.Manager.Close()

	database := db.GetDB()

	// 4. Initialize optional Storage gRPC client (if configured)
	var storageConn *grpc.ClientConn
	storageAddr := os.Getenv("STORAGE_SERVICE_ADDR")
	if storageAddr != "" {
		dialCtx, dialCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer dialCancel()
		conn, err := grpc.DialContext(
			dialCtx,
			storageAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		if err != nil {
			appLogger.Warnf("Storage service dial failed (%s): %v — running without storage integration", storageAddr, err)
		} else {
			storageConn = conn
			defer func() { _ = storageConn.Close() }()
			appLogger.Infof("Storage service client initialized: %s", storageAddr)
		}
	}

	// 5. Initialize Media Server (includes Kafka publisher and processing worker)
	mediaServer, err := media.NewMediaServer(database, storageConn)
	if err != nil {
		appLogger.Fatalf("Failed to create media server: %v", err)
	}

	// 6. Register gRPC service and start listening
	grpcServer := grpc.NewServer()
	mediaservicev1.RegisterMediaServiceServer(grpcServer, mediaServer.GetGRPCServer())

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		appLogger.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// 7. Start gRPC server in background
	go func() {
		appLogger.Infof("Media gRPC server listening on :%s", port)
		if err := grpcServer.Serve(lis); err != nil {
			appLogger.Fatalf("gRPC server error: %v", err)
		}
	}()

	// 8. Graceful Shutdown on signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	appLogger.Info("Shutting down...")
	grpcServer.GracefulStop()
	_ = mediaServer.Close()
	appLogger.Info("Media service stopped.")
}
