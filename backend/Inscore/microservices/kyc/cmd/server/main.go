package main

// KYC microservice entry point.
//
// This binary exposes the insuretech.kyc.services.v1.KYCService gRPC API.
// It is a thin wrapper around the kyc package, backed by the shared InsureTech
// PostgreSQL database (authn_schema.kyc_verifications).
//
// Configuration is via environment variables and shared ops/config YAML files.

import (
	"context"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/kyc"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/serviceaddr"
	kycservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/kyc/services/v1"
	"github.com/newage-saint/insuretech/ops/config"
	"github.com/newage-saint/insuretech/ops/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
)

// ServicesConfig mirrors the relevant parts of services.yaml.
type ServicesConfig = serviceaddr.ServicesConfig

func main() {
	// 1. Logger
	_ = appLogger.Initialize(appLogger.Config{
		Level:  "info",
		Format: "text",
		Output: "console",
	})
	appLogger.Info("Starting KYC microservice...")

	// Load and normalize repo-root env vars before resolving config or opening DB
	// connections. This keeps KYC aligned with AuthN startup behavior.
	if err := env.Load(); err != nil {
		appLogger.Warnf("Failed to load .env from repo root: %v", err)
	}

	// 2. Resolve port from services.yaml
	servicesConfigPath, err := config.ResolveConfigPath("services.yaml")
	if err != nil {
		appLogger.Fatalf("Failed to resolve services.yaml: %v", err)
	}
	servicesData, err := os.ReadFile(servicesConfigPath)
	if err != nil {
		appLogger.Fatalf("Failed to read services.yaml: %v", err)
	}
	var svcConfig ServicesConfig
	if err := yaml.Unmarshal(servicesData, &svcConfig); err != nil {
		appLogger.Fatalf("Failed to parse services.yaml: %v", err)
	}
	kycSvc, exists := svcConfig.Services["kyc"]
	if !exists {
		appLogger.Fatal("'kyc' service not found in services.yaml")
	}
	port := strconv.Itoa(kycSvc.Ports.Grpc)
	if os.Getenv("KYC_PORT") != "" || os.Getenv("KYC_GRPC_PORT") != "" || os.Getenv("KYC_HTTP_PORT") != "" {
		appLogger.Warn("KYC_PORT/KYC_GRPC_PORT/KYC_HTTP_PORT env values are ignored; using backend/inscore/configs/services.yaml")
	}
	appLogger.Infof("KYC service '%s' starting on gRPC port %s", kycSvc.Name, port)

	// 3. Database
	dbConfigPath, err := config.ResolveConfigPath("database.yaml")
	if err != nil {
		appLogger.Fatalf("Failed to resolve database.yaml: %v", err)
	}
	if err := db.InitializeManagerForService(dbConfigPath); err != nil {
		appLogger.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Manager.Close()
	database := db.GetDB()

	// 4. KYC Service
	kycService := kyc.NewKYCService(database)

	// Wire FLVE adapter so KYC service can start eKYC sessions.
	// Read endpoint and token from flve.yaml (same source as authn service) so
	// godotenv not-overwrite behaviour for existing OS env vars doesn't cause
	// empty token issues when FLVE_HF_TOKEN was not set in the parent shell.
	var flveYAML struct {
		FLVE struct {
			HFEndpoint string `yaml:"hf_endpoint"`
			HFToken    string `yaml:"hf_token"`
		} `yaml:"flve"`
	}
	if flveConfigPath, err := config.ResolveConfigPath("flve.yaml"); err == nil {
		if data, err := os.ReadFile(flveConfigPath); err == nil {
			_ = yaml.Unmarshal(data, &flveYAML)
		}
	}
	flveEndpoint := os.Getenv("FLVE_HF_ENDPOINT")
	if flveEndpoint == "" {
		flveEndpoint = flveYAML.FLVE.HFEndpoint
	}
	if flveEndpoint == "" {
		flveEndpoint = "https://farukhannan-flve.hf.space"
	}
	flveToken := os.Getenv("FLVE_HF_TOKEN")
	if flveToken == "" {
		flveToken = flveYAML.FLVE.HFToken
	}
	kycService.SetFLVEAdapter(kyc.NewFLVEAdapter(flveEndpoint, flveToken, 30*time.Second))
	appLogger.Infof("FLVE adapter configured: endpoint=%s token_len=%d", flveEndpoint, len(flveToken))

	// 5. gRPC Server
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		appLogger.Fatalf("Failed to listen on :%s: %v", port, err)
	}

	grpcServer := grpc.NewServer(
		grpc.ConnectionTimeout(30 * time.Second),
	)
	kycservicev1.RegisterKYCServiceServer(grpcServer, kycService)

	// Health check
	healthSvc := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthSvc)
	healthSvc.SetServingStatus("insuretech.kyc.services.v1.KYCService", grpc_health_v1.HealthCheckResponse_SERVING)

	// Reflection (dev/debug)
	reflection.Register(grpcServer)

	// 6. Start
	go func() {
		appLogger.Infof("KYC gRPC server listening on :%s", port)
		if err := grpcServer.Serve(lis); err != nil {
			appLogger.Fatalf("KYC gRPC server crashed: %v", err)
		}
	}()

	// 7. Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	appLogger.Info("Shutting down KYC service...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()
	select {
	case <-stopped:
		appLogger.Info("KYC service stopped gracefully.")
	case <-ctx.Done():
		appLogger.Warn("Graceful shutdown timed out — forcing stop.")
		grpcServer.Stop()
	}
}
