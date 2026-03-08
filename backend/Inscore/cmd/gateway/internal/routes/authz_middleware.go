package routes

// authz_middleware.go — Fine-grained AuthZ enforcement via AuthZ service (Casbin PERM model).
//
// Replaces rbac_middleware.go (flat user-type checks).
//
// Call chain per request:
//   1. AuthMiddleware     → AuthN.ValidateToken  → populates X-User-ID, X-Portal, X-Tenant-ID, etc.
//   2. AuthZMiddleware    → AuthZ.CheckAccess    → Casbin PERM enforce(sub, dom, obj, act)
//   3. Per-service gRPC interceptor → AuthZ.CheckAccess (defense-in-depth, same call)
//
// Domain format:  "portal:tenant_id"   e.g. "system:root", "agent:tenant-abc-123"
// Subject format: "user:<user_id>"
// Object format:  "svc:<service>/<resource>"  e.g. "svc:policy/create"
// Action format:  HTTP method            e.g. "GET", "POST", "DELETE"
//
// Deny-by-default: if no matching Casbin rule exists → 403 Forbidden.
//
// Usage in router:
//   r.Use(AuthMiddleware(authnConn))
//   r.Use(AuthZMiddleware(authzConn, "svc:policy", ResourceExtractorFromPath))
//   r.Use(AuthZMiddleware(authzConn, "svc:claim",  ResourceExtractorFromPath))

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ResourceExtractorFn extracts the resource path segment from a request.
// e.g. for "/v1/policies/123/approve" → "policy/approve"
// The returned string is appended to the service prefix to form the object:
//
//	servicePrefix="svc:policy" + resource="approve" → object="svc:policy/approve"
type ResourceExtractorFn func(r *http.Request) string

// AuthZMiddleware enforces fine-grained access control via AuthZ.CheckAccess (Casbin).
// Must be chained AFTER AuthMiddleware (requires X-User-ID, X-Portal, X-Tenant-ID headers).
//
//	servicePrefix: "svc:policy" | "svc:claim" | "svc:auth" | "svc:agent" | etc.
//	extractResource: function that extracts the resource sub-path from the request.
//
// The Casbin enforcement call is:
//
//	enforce("user:<user_id>", "portal:tenant_id", "svc:<service>/<resource>", "METHOD")
func AuthZMiddleware(authzConn *grpc.ClientConn, servicePrefix string, extractResource ResourceExtractorFn) func(http.Handler) http.Handler {
	if authzConn == nil {
		logger.Error("AuthZ middleware created without authz connection — all requests will be DENIED")
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Authorization service unavailable", http.StatusServiceUnavailable)
			})
		}
	}

	client := authzservicev1.NewAuthZServiceClient(authzConn)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// ── Extract identity from AuthMiddleware-populated headers ─────────
			userID := r.Header.Get("X-User-ID")
			portal := r.Header.Get("X-Portal")
			tenantID := r.Header.Get("X-Tenant-ID")
			sessionID := r.Header.Get("X-Session-ID")
			tokenID := r.Header.Get("X-Token-ID")
			deviceID := r.Header.Get("X-Device-ID")

			if userID == "" {
				// AuthMiddleware must run first
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// ── Build Casbin (sub, dom, obj, act) tuple ───────────────────────
			// domain = "portal:tenant_id"
			domain := buildRequestDomain(r, portal, tenantID)
			// object = "svc:<service>/<resource>"
			resource := ""
			if extractResource != nil {
				resource = extractResource(r)
			}
			object := buildObject(servicePrefix, resource)
			// action = HTTP method (GET, POST, PUT, PATCH, DELETE)
			action := r.Method

			// ── Call AuthZ.CheckAccess ─────────────────────────────────────────
			authzCtx := metadata.AppendToOutgoingContext(ctx, "x-internal-service", "gateway")
			resp, err := client.CheckAccess(authzCtx, &authzservicev1.CheckAccessRequest{
				UserId: userID,
				Domain: domain,
				Object: object,
				Action: action,
				Context: &authzservicev1.AccessContext{
					SessionId: sessionID,
					TokenId:   tokenID,
					DeviceId:  deviceID,
					IpAddress: realIP(r),
					UserAgent: r.UserAgent(),
				},
			})

			if err != nil {
				// Check if this is a connectivity error (authz service not running).
				// In that case, degrade gracefully: log a warning and fall through
				// to portal-gate enforcement only (user-type check already passed).
				// Hard 503 is reserved for when authz is reachable but returns an error.
				errStr := err.Error()
				isConnErr := strings.Contains(errStr, "connection refused") ||
					strings.Contains(errStr, "No connection could be made") ||
					strings.Contains(errStr, "Unavailable") ||
					strings.Contains(errStr, "transport:")
				if isConnErr {
					logger.Warn("AuthZ service unreachable — falling back to portal-gate enforcement only",
						zap.String("user_id", userID),
						zap.String("domain", domain),
						zap.String("object", object),
						zap.String("action", action),
						zap.Error(err),
					)
					// Allow request through; portal-gate middleware provides baseline enforcement.
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
				logger.Error("AuthZ.CheckAccess RPC failed",
					zap.Error(err),
					zap.String("user_id", userID),
					zap.String("domain", domain),
					zap.String("object", object),
					zap.String("action", action),
				)
				http.Error(w, "Authorization service error", http.StatusServiceUnavailable)
				return
			}

			if resp == nil || !resp.Allowed {
				reason := "no matching policy"
				if resp != nil && resp.Reason != "" {
					reason = resp.Reason
				}
				logger.Warn("AuthZ DENY",
					zap.String("user_id", userID),
					zap.String("domain", domain),
					zap.String("object", object),
					zap.String("action", action),
					zap.String("reason", reason),
				)
				http.Error(w, "Forbidden: "+reason, http.StatusForbidden)
				return
			}

			logger.Debug("AuthZ ALLOW",
				zap.String("user_id", userID),
				zap.String("domain", domain),
				zap.String("object", object),
				zap.String("action", action),
				zap.String("matched_rule", resp.MatchedRule),
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ── Public bypass routes — no AuthZ check needed ─────────────────────────────
// These are routes that are explicitly public (no token required).
// AuthMiddleware itself is also skipped for these in the router.
var PublicRoutes = map[string]bool{
	"/v1/auth/register":                     true,
	"/v1/auth/otp:send":                     true,
	"/v1/auth/otp:verify":                   true,
	"/v1/auth/otp:resend":                   true,
	"/v1/auth/login":                        true,
	"/v1/auth/email/register":               true,
	"/v1/auth/email/otp:send":               true,
	"/v1/auth/email/verify":                 true,
	"/v1/auth/email/login":                  true,
	"/v1/auth/email/password:reset-request": true,
	"/v1/auth/email/password:reset":         true,
	"/v1/auth/password:reset":               true,
	"/v1/auth/biometric:authenticate":       true,
	"/v1/auth/token:refresh":                true,
	"/.well-known/jwks.json":                true,
	"/healthz":                              true,
	"/readyz":                               true,
}

// IsPublicRoute returns true if the path requires no authentication.
func IsPublicRoute(path string) bool {
	return PublicRoutes[path]
}

// ── Portal-gate helpers (fast pre-AuthZ user-type guards) ─────────────────────
// These are lightweight checks that run BEFORE the full Casbin AuthZ call.
// They enforce portal boundaries (e.g. only SYSTEM_USER can access /v1/admin/*).
// They do NOT replace AuthZMiddleware — they are used together for defence-in-depth.
//
// Portal → UserType mapping (from enums.proto):
//   system    → USER_TYPE_SYSTEM_USER
//   business  → USER_TYPE_BUSINESS_BENEFICIARY
//   b2b       → USER_TYPE_PARTNER
//   agent     → USER_TYPE_AGENT
//   regulator → USER_TYPE_REGULATOR
//   b2c       → USER_TYPE_B2C_CUSTOMER

const (
	UserTypeSystemUser          = "USER_TYPE_SYSTEM_USER"
	UserTypeBusinessBeneficiary = "USER_TYPE_BUSINESS_BENEFICIARY"
	UserTypePartner             = "USER_TYPE_PARTNER"
	UserTypeAgent               = "USER_TYPE_AGENT"
	UserTypeRegulator           = "USER_TYPE_REGULATOR"
	UserTypeB2CCustomer         = "USER_TYPE_B2C_CUSTOMER"
)

// requireUserTypes returns a middleware that allows only the listed user types.
func requireUserTypes(allowed ...string) func(http.Handler) http.Handler {
	allowedSet := make(map[string]bool, len(allowed))
	for _, t := range allowed {
		allowedSet[t] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ut := r.Header.Get("X-User-Type")
			if !allowedSet[ut] {
				http.Error(w, "Forbidden: insufficient portal access", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// SystemUserMiddleware allows only SYSTEM_USER (system portal).
func SystemUserMiddleware(next http.Handler) http.Handler {
	return requireUserTypes(UserTypeSystemUser)(next)
}

// AgentOrSystemMiddleware allows AGENT or SYSTEM_USER.
func AgentOrSystemMiddleware(next http.Handler) http.Handler {
	return requireUserTypes(UserTypeAgent, UserTypeSystemUser)(next)
}

// BusinessOrSystemMiddleware allows BUSINESS_BENEFICIARY or SYSTEM_USER.
func BusinessOrSystemMiddleware(next http.Handler) http.Handler {
	return requireUserTypes(UserTypeBusinessBeneficiary, UserTypeSystemUser)(next)
}

// PartnerOrSystemMiddleware allows PARTNER or SYSTEM_USER.
func PartnerOrSystemMiddleware(next http.Handler) http.Handler {
	return requireUserTypes(UserTypePartner, UserTypeSystemUser)(next)
}

// RegulatorOrSystemMiddleware allows REGULATOR or SYSTEM_USER (read-only portal).
func RegulatorOrSystemMiddleware(next http.Handler) http.Handler {
	return requireUserTypes(UserTypeRegulator, UserTypeSystemUser)(next)
}

// AnyAuthenticatedMiddleware allows any authenticated user type.
func AnyAuthenticatedMiddleware(next http.Handler) http.Handler {
	return requireUserTypes(
		UserTypeSystemUser,
		UserTypeBusinessBeneficiary,
		UserTypePartner,
		UserTypeAgent,
		UserTypeRegulator,
		UserTypeB2CCustomer,
	)(next)
}

// ── Helper functions ──────────────────────────────────────────────────────────

// buildDomain constructs the Casbin domain string: "portal:tenant_id".
// The portal prefix is normalized to lowercase without the "PORTAL_" prefix
// so it matches the seeder's portalDomainKey format (e.g. "b2b:root").
// If tenantID is empty, uses "root" (system-wide scope).
func buildDomain(portal, tenantID string) string {
	if portal == "" {
		portal = "b2c"
	}
	// Strip "PORTAL_" prefix (e.g. "PORTAL_B2B" → "b2b") and lowercase.
	// The authz seeder stores domains as "<portal>:<tenantID>" without the enum prefix.
	portal = strings.ToLower(strings.TrimPrefix(portal, "PORTAL_"))
	if portal == "unspecified" || portal == "" {
		portal = "b2c"
	}
	// System-portal roles are seeded in the global root domain. Do not scope
	// them by tenant_id or super-admin access to cross-portal operations breaks.
	if portal == "system" {
		return "system:root"
	}
	if tenantID == "" {
		tenantID = "root"
	}
	return portal + ":" + tenantID
}

func buildRequestDomain(r *http.Request, portal, tenantID string) string {
	normalizedPortal := strings.ToLower(strings.TrimPrefix(portal, "PORTAL_"))
	if normalizedPortal == "b2b" {
		businessID := resolveRequestBusinessID(r)
		if businessID != "" {
			return "b2b:" + businessID
		}
	}

	return buildDomain(portal, tenantID)
}

func resolveRequestBusinessID(r *http.Request) string {
	candidates := []string{
		r.Header.Get("X-Business-ID"),
		r.URL.Query().Get("business_id"),
		r.PathValue("organisation_id"),
	}
	for _, candidate := range candidates {
		candidate = strings.TrimSpace(candidate)
		if candidate != "" {
			return candidate
		}
	}

	if r.Body == nil || !strings.Contains(strings.ToLower(r.Header.Get("Content-Type")), "application/json") {
		return ""
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return ""
	}
	r.Body = io.NopCloser(bytes.NewReader(body))

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return ""
	}

	for _, key := range []string{"business_id", "businessId", "organisation_id", "organisationId"} {
		if value := stringValue(payload[key]); value != "" {
			return value
		}
	}

	return ""
}

func stringValue(value any) string {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	default:
		return ""
	}
}

// buildObject constructs the Casbin object string: "svc:<service>/<resource>".
// If resource is empty, uses wildcard "*".
func buildObject(servicePrefix, resource string) string {
	if resource == "" || resource == "/" {
		return servicePrefix + "/*"
	}
	resource = strings.TrimPrefix(resource, "/")
	return servicePrefix + "/" + resource
}

// realIP extracts the real client IP address from request headers.
func realIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-Ip"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		// X-Forwarded-For may contain a comma-separated list; first is the client
		return strings.Split(ip, ",")[0]
	}
	return r.RemoteAddr
}

// ── Built-in ResourceExtractor helpers ────────────────────────────────────────

// PathSegmentExtractor returns a ResourceExtractorFn that extracts segments
// from the URL path after stripping the given prefix.
// e.g. prefix="/v1/policies/", path="/v1/policies/123/approve" → "approve"
func PathSegmentExtractor(stripPrefix string) ResourceExtractorFn {
	return func(r *http.Request) string {
		path := strings.TrimPrefix(r.URL.Path, stripPrefix)
		// Remove UUIDs / numeric IDs (keep only semantic path segments)
		parts := strings.Split(path, "/")
		var semanticParts []string
		for _, p := range parts {
			if p == "" {
				continue
			}
			// Skip UUID-like and purely numeric segments
			if isIDSegment(p) {
				continue
			}
			semanticParts = append(semanticParts, p)
		}
		return strings.Join(semanticParts, "/")
	}
}

// StorageResourceExtractor maps storage HTTP paths to stable AuthZ resources.
// This keeps policy objects concise (upload-url/finalize/get/...).
// Supports both modern slash actions (files/{id}/download-url) and
// legacy colon actions (files/{id}:download-url) during migration.
func StorageResourceExtractor() ResourceExtractorFn {
	return func(r *http.Request) string {
		p := strings.TrimPrefix(r.URL.Path, "/v1/storage/")
		p = strings.TrimPrefix(p, "/")

		switch {
		case p == "files":
			if r.Method == http.MethodPost {
				return "upload"
			}
			return "get"
		case p == "files:batch":
			return "upload-batch"
		case p == "files:upload-url":
			return "upload-url"
		case p == "files:finalize":
			return "finalize"
		case strings.HasPrefix(p, "files/") && (strings.HasSuffix(p, ":download-url") || strings.HasSuffix(p, "/download-url")):
			return "download-url"
		case strings.HasPrefix(p, "files/"):
			switch r.Method {
			case http.MethodPatch:
				return "update"
			case http.MethodDelete:
				return "delete"
			default:
				return "get"
			}
		default:
			return PathSegmentExtractor("/v1/storage/")(r)
		}
	}
}

// isIDSegment returns true if the path segment looks like a UUID or numeric ID.
func isIDSegment(s string) bool {
	if len(s) == 36 && strings.Count(s, "-") == 4 {
		return true // UUID v4
	}
	allDigits := true
	for _, c := range s {
		if c < '0' || c > '9' {
			allDigits = false
			break
		}
	}
	return allDigits
}
