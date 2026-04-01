package grpc

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestDefaultServerConfigAndNewServer(t *testing.T) {
	cfg := DefaultServerConfig()
	if cfg.Port != "50142" {
		t.Fatalf("expected default port 50142, got %q", cfg.Port)
	}

	srv, err := NewServer(cfg, &fakeOrderService{})
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}
	if srv == nil || srv.handler == nil || srv.health == nil {
		t.Fatalf("expected initialized server components")
	}
}

func TestServerHealthCheck(t *testing.T) {
	srv := &Server{config: &Config{}}
	if err := srv.HealthCheck(context.Background()); err == nil {
		t.Fatalf("expected nil db error")
	}

	dbConn, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	srv.config.DB = dbConn
	if err := srv.HealthCheck(context.Background()); err != nil {
		t.Fatalf("HealthCheck() error = %v", err)
	}
}
