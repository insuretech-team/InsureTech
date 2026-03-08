package main

import (
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/insurance/service"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/ops/config"
	"gopkg.in/yaml.v3"
	"google.golang.org/grpc"
	insurancev1 "github.com/newage-saint/insuretech/gen/go/insuretech/insurance/services/v1"
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
	appLogger.Info("Starting Insurance microservice...")

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

	insuranceConfig, exists := svcConfig.Services["insurance"]
	if !exists {
		appLogger.Fatal("Configuration for 'insurance' service not found in services.yaml")
	}
	port := strconv.Itoa(insuranceConfig.Ports.Grpc)
	appLogger.Infof("Service '%s' configured on port %s", insuranceConfig.Name, port)

	// 3. Initialize Database
	dbConfigPath, err := config.ResolveConfigPath("database.yaml")
	if err != nil {
		appLogger.Fatalf("Failed to resolve database config path: %v", err)
	}
	appLogger.Infof("Using database config: %s", dbConfigPath)
	
	if err := db.InitializeManagerForService(dbConfigPath); err != nil {
		appLogger.Errorf("Failed to initialize database: %v", err)
		appLogger.Fatal("Database initialization failed")
	}
	defer db.Manager.Close()

	database := db.GetDB()

	// 4. Initialize Service (repositories are created internally)
	insuranceService := service.NewInsuranceService(database)

	// 5. Initialize gRPC Server
	grpcServer := grpc.NewServer()
	insurancev1.RegisterInsuranceServiceServer(grpcServer, insuranceService)

	// 6. Start Server
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		appLogger.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	appLogger.Infof("Insurance Service listening on port %s", port)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			appLogger.Fatalf("Failed to serve: %v", err)
		}
	}()

	// 7. Graceful Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	appLogger.Info("Shutting down...")
	grpcServer.GracefulStop()
	appLogger.Info("Stopped.")
}
