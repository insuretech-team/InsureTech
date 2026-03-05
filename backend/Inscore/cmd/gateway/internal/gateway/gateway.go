package gateway

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/handlers"
	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/lifecycle"
	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/middleware"
	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/resilience"
	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/routes"
	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/security"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	opsconfig "github.com/newage-saint/insuretech/ops/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

// Config holds gateway configuration.
type Config struct {
	Port                string
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	IdleTimeout         time.Duration
	HealthCheckInterval time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		Port:                "8080",
		ReadTimeout:         10 * time.Second,
		WriteTimeout:        0,
		IdleTimeout:         60 * time.Second,
		HealthCheckInterval: 30 * time.Second,
	}
}

// Gateway orchestrates HTTP routing and resilient gRPC clients.
type Gateway struct {
	config *Config
	port   string

	clientManager *resilience.ResilientClientManager

	// handlers
	authnHandler *handlers.AuthnHandler
	dlrHandler   *handlers.DLRHandler

	// http server
	httpServer *http.Server
	shutdown   *lifecycle.GracefulShutdown

	// tls
	tlsConfig *tls.Config

	shutdownOnce sync.Once
}

func NewGateway(cfg *Config) (*Gateway, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	tlsCfg, err := security.TLSConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS env config: %w", err)
	}
	var tlsConfig *tls.Config
	if tlsCfg != nil {
		tlsConfig, err = security.LoadTLSConfig(tlsCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to build TLS config: %w", err)
		}
		logger.Info("TLS enabled",
			zap.String("cert", tlsCfg.CertFile),
			zap.Bool("mtls", tlsCfg.RequireClientCert),
		)
	}

	return &Gateway{
		config:        cfg,
		port:          cfg.Port,
		clientManager: resilience.NewResilientClientManager(),
		tlsConfig:     tlsConfig,
	}, nil
}

func (g *Gateway) IsTLSEnabled() bool { return g.tlsConfig != nil }

func (g *Gateway) Start(ctx context.Context) error {
	logger.Info("Starting gateway", zap.String("port", g.port))

	if err := g.registerAllServices(ctx); err != nil {
		return err
	}
	_ = g.clientManager.InitializeAll(ctx)
	if err := g.initHandlers(); err != nil {
		return err
	}

	router := g.createRouter()

	g.httpServer = &http.Server{
		Addr:           ":" + g.port,
		Handler:        router,
		ReadTimeout:    g.config.ReadTimeout,
		WriteTimeout:   g.config.WriteTimeout,
		IdleTimeout:    g.config.IdleTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	g.shutdown = lifecycle.NewGracefulShutdown(g.httpServer)
	g.httpServer.Handler = g.shutdown.Middleware(g.httpServer.Handler)

	go func() {
		if g.tlsConfig != nil {
			g.httpServer.TLSConfig = g.tlsConfig
			logger.Info("HTTPS gateway listening", zap.String("addr", g.httpServer.Addr))
			if err := g.httpServer.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				logger.Fatal("HTTPS gateway failed", zap.Error(err))
			}
			return
		}

		logger.Info("HTTP gateway listening", zap.String("addr", g.httpServer.Addr))
		if err := g.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("HTTP gateway failed", zap.Error(err))
		}
	}()

	return nil
}

func (g *Gateway) Shutdown(ctx context.Context) error {
	var err error
	g.shutdownOnce.Do(func() {
		logger.Info("Gateway shutdown started")

		if g.shutdown != nil {
			err = g.shutdown.Shutdown(ctx)
		}

		if closeErr := g.clientManager.CloseAll(); closeErr != nil {
			logger.Warn("Failed to close some clients", zap.Error(closeErr))
			if err == nil {
				err = closeErr
			}
		}
	})
	return err
}

func (g *Gateway) registerAllServices(ctx context.Context) error {
	_ = ctx
	services := []string{
		// InScore (Go) services
		"authn", "authz", "tenant", "audit", "kyc", "partner", "fraud", "b2b",
		"beneficiary", "notification", "storage", "media", "docgen", "webrtc", "workflow",
		// PoliSync (C# .NET 8) services
		"product", "quote", "order", "commission", "policy", "underwriting", "claim",
	}
	addresses := g.resolveServiceAddresses(services)

	registered := 0
	for _, name := range services {
		addr := strings.TrimSpace(addresses[name])
		if addr == "" {
			continue
		}

		cfg := resilience.DefaultServiceClientConfig(name, addr)
		switch name {
		case "authn":
			cfg.RetryPolicy = resilience.AggressiveRetryPolicy()
			cfg.PoolConfig.PoolSize = 10
		default:
			cfg.PoolConfig.PoolSize = 5
		}

		if _, err := g.clientManager.RegisterClient(cfg); err != nil {
			logger.Warn("Failed to register service client", zap.String("service", name), zap.Error(err))
			continue
		}

		logger.Info("Service registered", zap.String("service", name), zap.String("address", addr))
		registered++
	}

	if registered == 0 {
		return fmt.Errorf("no services registered: set *_GRPC_ADDR env vars or verify backend/inscore/configs/services.yaml")
	}
	return nil
}

type gatewayServicesConfig struct {
	Services map[string]struct {
		Ports struct {
			Grpc int `yaml:"grpc"`
		} `yaml:"ports"`
	} `yaml:"services"`
}

func (g *Gateway) resolveServiceAddresses(serviceNames []string) map[string]string {
	resolved := make(map[string]string, len(serviceNames))

	// 1) Explicit env wins.
	for _, name := range serviceNames {
		addrEnv := strings.ToUpper(name) + "_GRPC_ADDR"
		addr := sanitizeGRPCAddr(strings.TrimSpace(os.Getenv(addrEnv)))
		if addr == "" {
			continue
		}
		resolved[name] = addr
	}

	// 2) Fallback to services.yaml (localhost:<grpc_port>) for missing services.
	missing := make([]string, 0, len(serviceNames))
	for _, name := range serviceNames {
		if resolved[name] == "" {
			missing = append(missing, name)
		}
	}
	if len(missing) == 0 {
		return resolved
	}

	servicesConfigPath, err := opsconfig.ResolveConfigPath("services.yaml")
	if err != nil {
		logger.Warn("Failed to resolve services.yaml for gateway service discovery", zap.Error(err))
		return resolved
	}
	raw, err := os.ReadFile(servicesConfigPath)
	if err != nil {
		logger.Warn("Failed to read services.yaml for gateway service discovery",
			zap.String("path", servicesConfigPath),
			zap.Error(err),
		)
		return resolved
	}
	var cfg gatewayServicesConfig
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		logger.Warn("Failed to parse services.yaml for gateway service discovery",
			zap.String("path", servicesConfigPath),
			zap.Error(err),
		)
		return resolved
	}

	defaultHost := strings.TrimSpace(os.Getenv("GATEWAY_GRPC_HOST"))
	if defaultHost == "" {
		defaultHost = "localhost"
	}

	for _, name := range missing {
		svc, ok := cfg.Services[name]
		if !ok || svc.Ports.Grpc <= 0 {
			continue
		}
		addr := defaultHost + ":" + strconv.Itoa(svc.Ports.Grpc)
		resolved[name] = addr
		logger.Info("Service address resolved from services.yaml",
			zap.String("service", name),
			zap.String("address", addr),
		)
	}

	return resolved
}

func sanitizeGRPCAddr(raw string) string {
	if raw == "" {
		return ""
	}
	addr := strings.TrimSpace(raw)
	addr = strings.TrimPrefix(addr, "grpc://")
	addr = strings.TrimPrefix(addr, "http://")
	addr = strings.TrimPrefix(addr, "https://")
	return addr
}

func (g *Gateway) initHandlers() error {
	if authnClient, err := g.clientManager.GetClient("authn"); err == nil {
		conn, err := authnClient.GetConnection()
		if err == nil {
			g.authnHandler = handlers.NewAuthnHandler(conn)
			g.dlrHandler = handlers.NewDLRHandler(conn)
			logger.Info("Authn handler initialized")
			logger.Info("DLR handler initialized")
		} else {
			logger.Warn("Authn connection unavailable", zap.Error(err))
		}
	}
	return nil
}

func (g *Gateway) createRouter() http.Handler {
	var authnConn *grpc.ClientConn
	if authnClient, err := g.clientManager.GetClient("authn"); err == nil {
		conn, err := authnClient.GetConnection()
		if err == nil {
			authnConn = conn
		}
	}

	// Wire authzConn when the AuthZ microservice is available.
	// Falls back gracefully to portal-gate (user-type) only when authz is not yet deployed.
	var authzConn *grpc.ClientConn
	if authzClient, err := g.clientManager.GetClient("authz"); err == nil {
		conn, err := authzClient.GetConnection()
		if err == nil {
			authzConn = conn
			logger.Info("AuthZ middleware enabled — Casbin PERM enforcement active")
		} else {
			logger.Warn("AuthZ client registered but connection unavailable — falling back to portal-gate only", zap.Error(err))
		}
	} else {
		logger.Info("AuthZ service not configured (no resolved address) — portal-gate enforcement only")
	}

	return routes.NewRouter(g.authnHandler, authnConn, authzConn, g.clientManager, g.dlrHandler)
}

// MiddlewareStack returns the default middleware stack (primarily for tests/custom wiring).
func MiddlewareStack(h http.Handler) http.Handler {
	h = middleware.Recovery(h)
	h = middleware.RequestID(h)
	h = middleware.SecurityHeaders(h)
	h = middleware.Metrics(h)
	return h
}
