package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/respond"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// PoliSyncHandler is a generic HTTP reverse-proxy for all PoliSync C# services.
// The Go gateway validates JWT and injects X-* identity headers; this handler
// forwards the full HTTP request (including those headers) to the PoliSync
// REST endpoint (each gRPC service exposes a companion HTTP/1.1 port).
//
// Architecture:
//
//	HTTP client → gateway (auth + authz + X-* headers injected)
//	           → PoliSyncHandler.Proxy()
//	           → PoliSync C# HTTP companion port
//	           → PoliSync AuthInterceptor reads X-* headers → ICurrentUser
//
// Identity propagation (gateway injects, PoliSync reads):
//
//	X-User-ID, X-Tenant-ID, X-Partner-ID, X-Token-ID,
//	X-User-Type, X-Portal, X-Roles, X-Request-ID
type PoliSyncHandler struct {
	conn        *grpc.ClientConn // reserved for future direct gRPC calls
	serviceName string
	overrideURL string // if set, used instead of poliSyncServiceURL map (for testing)
}

// NewPoliSyncHandler creates a handler that reverse-proxies to a PoliSync HTTP companion port.
// The service URL is resolved from the built-in poliSyncServiceURL map (Docker DNS).
func NewPoliSyncHandler(conn *grpc.ClientConn, serviceName string) *PoliSyncHandler {
	return &PoliSyncHandler{
		conn:        conn,
		serviceName: serviceName,
		overrideURL: "",
	}
}

// NewPoliSyncHandlerWithURL creates a handler with an explicit target URL.
// Used in tests to point at mock backends, and optionally in dev for local overrides.
func NewPoliSyncHandlerWithURL(conn *grpc.ClientConn, serviceName, targetURL string) *PoliSyncHandler {
	return &PoliSyncHandler{
		conn:        conn,
		serviceName: serviceName,
		overrideURL: targetURL,
	}
}

// Proxy returns an http.Handler that reverse-proxies the current request to the
// PoliSync service HTTP companion port. All X-* identity headers that the gateway
// auth middleware injected are preserved and forwarded.
func (h *PoliSyncHandler) Proxy() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		target := h.overrideURL
		if target == "" {
			target = poliSyncServiceURL(h.serviceName)
		}
		if target == "" {
			respond.Error(w, r, http.StatusBadGateway, "UNAVAILABLE",
				"PoliSync service address not configured: "+h.serviceName)
			return
		}

		targetURL, err := url.Parse(target)
		if err != nil {
			respond.Error(w, r, http.StatusInternalServerError, "INTERNAL", "invalid upstream URL")
			return
		}

		proxy := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = targetURL.Scheme
				req.URL.Host = targetURL.Host
				req.Host = targetURL.Host
				// Remove hop-by-hop headers
				req.Header.Del("Te")
				req.Header.Del("Trailers")
			},
			ErrorHandler: func(w http.ResponseWriter, req *http.Request, err error) {
				respond.Error(w, req, http.StatusBadGateway, "UNAVAILABLE",
					"PoliSync "+h.serviceName+" unreachable: "+err.Error())
			},
		}
		proxy.ServeHTTP(w, r)
	})
}

// writeJSONError was replaced by the unified respond package.
// See: github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/respond

// grpcStatusToHTTP was replaced by respond.GRPCCodeToHTTP from the unified respond package.
// See: github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/respond

// writeGRPCError writes a JSON error from a gRPC error.
func (h *PoliSyncHandler) writeGRPCError(w http.ResponseWriter, r *http.Request, err error) {
	respond.GRPCError(w, r, err)
}

// buildOutgoingMD extracts X-* identity headers as gRPC metadata.
// Reserved for future direct gRPC call path.
func buildOutgoingMD(r *http.Request) metadata.MD {
	md := metadata.New(nil)
	pairs := map[string]string{
		"X-User-ID":      "x-user-id",
		"X-Tenant-ID":    "x-tenant-id",
		"X-Partner-ID":   "x-partner-id",
		"X-Token-ID":     "x-token-id",
		"X-User-Type":    "x-user-type",
		"X-Portal":       "x-portal",
		"X-Roles":        "x-roles",
		"X-Request-ID":   "x-request-id",
		"X-Session-ID":   "x-session-id",
		"X-Session-Type": "x-session-type",
	}
	for header, mdKey := range pairs {
		if v := r.Header.Get(header); v != "" {
			md.Set(mdKey, v)
		}
	}
	return md
}

// poliSyncServiceURL maps service names to their HTTP companion port base URLs.
//
// Two categories of services use this handler:
//
//  1. PoliSync C# services — all run inside the single "polisync" Docker container,
//     each exposing a Kestrel HTTP/1.1 companion port alongside their gRPC port.
//     Host defaults to POLISYNC_HOST env var (falls back to "polisync").
//
//  2. InScore Go services — each runs in its own container with a dedicated HTTP
//     companion port. Host is read from {SERVICE}_HOST env var (e.g. PAYMENT_HOST).
//
// Port layout (must match services.yaml and appsettings.json Kestrel config):
//
//	PoliSync C# (host: "polisync"):
//	  insurance-service    → :50116
//	  product-service      → :50121
//	  quote-service        → :50131
//	  order-service        → :50141   (NOTE: not sync_order Go service at :50142)
//	  commission-service   → :50151
//	  policy-service       → :50161
//	  underwriting-service → :50171
//	  claim-service        → :50211
//
//	InScore Go (host from {SVC}_HOST env, e.g. PAYMENT_HOST=payment):
//	  payment-service      → :50191
//	  notification-service → :50231
//	  kyc-service          → :50091
//	  beneficiary-service  → :50111
//	  tenant-service       → :50051
//	  audit-service        → :50081
//	  billing-service      → :50196
//	  b2b-service          → :50113
//	  media-service        → :50261
//	  storage-service      → :50291
//	  fraud-service        → :50221
//	  partner-service      → :50101
func poliSyncServiceURL(serviceName string) string {
	// Allow full URL override per service: {SERVICE}_HTTP_ADDR takes highest priority.
	// e.g. PAYMENT_HTTP_ADDR=http://payment:50191
	envKey := strings.ToUpper(strings.ReplaceAll(strings.TrimSuffix(serviceName, "-service"), "-", "_")) + "_HTTP_ADDR"
	if override := os.Getenv(envKey); override != "" {
		return override
	}

	type serviceConfig struct {
		defaultHost string // Docker DNS name (overridable via {SVC}_HOST env)
		hostEnvKey  string // env var for host override
		port        string
	}

	configs := map[string]serviceConfig{
		// PoliSync C# — single container "polisync"
		"insurance-service":    {"polisync", "POLISYNC_HOST", "50116"},
		"product-service":      {"polisync", "POLISYNC_HOST", "50121"},
		"quote-service":        {"polisync", "POLISYNC_HOST", "50131"},
		"order-service":        {"polisync", "POLISYNC_HOST", "50141"},
		"commission-service":   {"polisync", "POLISYNC_HOST", "50151"},
		"policy-service":       {"polisync", "POLISYNC_HOST", "50161"},
		"underwriting-service": {"polisync", "POLISYNC_HOST", "50171"},
		"claim-service":        {"polisync", "POLISYNC_HOST", "50211"},
		// InScore Go — separate containers
		"payment-service":      {"payment", "PAYMENT_HOST", "50191"},
		"notification-service": {"notification", "NOTIFICATION_HOST", "50231"},
		"kyc-service":          {"kyc", "KYC_HOST", "50091"},
		"beneficiary-service":  {"beneficiary", "BENEFICIARY_HOST", "50111"},
		"tenant-service":       {"tenant", "TENANT_HOST", "50051"},
		"audit-service":        {"audit", "AUDIT_HOST", "50081"},
		"billing-service":      {"billing", "BILLING_HOST", "50196"},
		"b2b-service":          {"b2b", "B2B_HOST", "50113"},
		"media-service":        {"media", "MEDIA_HOST", "50261"},
		"storage-service":      {"storage", "STORAGE_HOST", "50291"},
		"fraud-service":        {"fraud", "FRAUD_HOST", "50221"},
		"partner-service":      {"partner", "PARTNER_HOST", "50101"},
	}

	cfg, ok := configs[serviceName]
	if !ok {
		return ""
	}

	host := os.Getenv(cfg.hostEnvKey)
	if host == "" {
		host = cfg.defaultHost
	}
	return "http://" + host + ":" + cfg.port
}

// Ensure exported methods satisfy interfaces (compile-time checks).
var _ http.Handler = (*PoliSyncHandler)(nil).Proxy()

// suppress unused warning on future-use helpers during scaffold phase
var (
	_ = buildOutgoingMD
	_ = (*PoliSyncHandler).writeGRPCError
)
