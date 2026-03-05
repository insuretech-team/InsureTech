package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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
			h.writeJSONError(w, http.StatusBadGateway, "UNAVAILABLE",
				"PoliSync service address not configured: "+h.serviceName)
			return
		}

		targetURL, err := url.Parse(target)
		if err != nil {
			h.writeJSONError(w, http.StatusInternalServerError, "INTERNAL", "invalid upstream URL")
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
				h.writeJSONError(w, http.StatusBadGateway, "UNAVAILABLE",
					"PoliSync "+h.serviceName+" unreachable: "+err.Error())
			},
		}
		proxy.ServeHTTP(w, r)
	})
}

// writeJSONError writes a JSON error response.
func (h *PoliSyncHandler) writeJSONError(w http.ResponseWriter, httpCode int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"code":    code,
		"message": message,
		"service": h.serviceName,
	})
}

// grpcStatusToHTTP maps gRPC status codes to HTTP status codes.
// Used when PoliSync returns gRPC errors (future direct gRPC call support).
func grpcStatusToHTTP(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	default:
		return http.StatusInternalServerError
	}
}

// writeGRPCError writes a JSON error from a gRPC error.
func (h *PoliSyncHandler) writeGRPCError(w http.ResponseWriter, err error) {
	st, _ := status.FromError(err)
	h.writeJSONError(w, grpcStatusToHTTP(st.Code()), st.Code().String(), st.Message())
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
// In production, Docker DNS resolves the service names.
// In dev, override with PRODUCT_HTTP_ADDR etc. env vars or update appsettings.Development.json.
func poliSyncServiceURL(serviceName string) string {
	m := map[string]string{
		"product-service":      "http://product-service:50121",
		"quote-service":        "http://quote-service:50131",
		"order-service":        "http://order-service:50141",
		"commission-service":   "http://commission-service:50151",
		"policy-service":       "http://policy-service:50161",
		"underwriting-service": "http://underwriting-service:50171",
		"claim-service":        "http://claim-service:50211",
	}
	return m[serviceName]
}

// Ensure exported methods satisfy interfaces (compile-time checks).
var _ http.Handler = (*PoliSyncHandler)(nil).Proxy()

// suppress unused warning on future-use helpers during scaffold phase
var (
	_ = buildOutgoingMD
	_ = (*PoliSyncHandler).writeGRPCError
)
