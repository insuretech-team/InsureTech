package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/gateway"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/ops/env"
	"go.uber.org/zap"
)

func main() {
	if err := logger.Initialize(logger.NoFileConfig()); err != nil {
		logger.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.GetLogger().Sync()

	logger.Info("Starting InScore API Gateway...")
	_ = env.Load()

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	cfg := &gateway.Config{
		Port:                port,
		ReadTimeout:         5 * time.Second,
		WriteTimeout:        0,
		IdleTimeout:         30 * time.Second,
		HealthCheckInterval: 10 * time.Second,
	}

	gw, err := gateway.NewGateway(cfg)
	if err != nil {
		logger.Fatal("Failed to create gateway", zap.Error(err))
	}

	if err := gw.Start(context.Background()); err != nil {
		logger.Fatal("Failed to start gateway", zap.Error(err))
	}

	protocol := "http"
	if gw.IsTLSEnabled() {
		protocol = "https"
	}

	logger.Info("Gateway started",
		zap.String("protocol", protocol),
		zap.String("port", port),
		zap.String("health", protocol+"://localhost:"+port+"/healthz"),
		zap.String("ready", protocol+"://localhost:"+port+"/readyz"),
	)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down gateway...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := gw.Shutdown(shutdownCtx); err != nil {
		logger.Error("Gateway shutdown error", zap.Error(err))
	}
	logger.Info("Gateway shutdown complete")
}
