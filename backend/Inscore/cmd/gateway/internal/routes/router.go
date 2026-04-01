package routes

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/handlers"
	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/middleware"
	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/resilience"
	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/respond"
	authnv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	authzv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	"google.golang.org/grpc"
)

// NewRouter wires the main HTTP routing for the gateway.
// authnConn: gRPC connection to authn service (required)
// authzConn: gRPC connection to authz service (required for AuthZ enforcement; nil = portal-gate only)
func NewRouter(authnHandler *handlers.AuthnHandler, authnConn *grpc.ClientConn, authzConn *grpc.ClientConn, clientManager *resilience.ResilientClientManager, dlrHandler *handlers.DLRHandler) http.Handler {
	mux := http.NewServeMux()
	paymentConn := getServiceConn(clientManager, "payment")

	authMW := AuthMiddleware(authnConn)
	csrfMW := CSRFMiddleware(authnConn)
	// OTP rate limit: read from env OTP_RATE_LIMIT (default 10/hour per IP+mobile).
	// Keep this generous enough for testing but strict enough for production abuse prevention.
	otpRateMax := 100
	if v := os.Getenv("OTP_RATE_LIMIT_MAX"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			otpRateMax = n
		}
	}
	otpRateWin := time.Hour
	if v := os.Getenv("OTP_RATE_LIMIT_WINDOW_MINUTES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			otpRateWin = time.Duration(n) * time.Minute
		}
	}
	otpRL := middleware.OTPRateLimit(otpRateMax, otpRateWin)
	// Login rate limit: read from env (default 30 per 15 minutes per IP).
	// Previous 10/hour was too aggressive — locked out legitimate users
	// who retry OTP flows or have device-binding re-logins.
	loginRateMax := 30
	if v := os.Getenv("RATE_LIMIT_LOGIN_MAX"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			loginRateMax = n
		}
	}
	loginRateWin := 15 * time.Minute
	if v := os.Getenv("RATE_LIMIT_LOGIN_WINDOW_MINUTES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			loginRateWin = time.Duration(n) * time.Minute
		}
	}
	loginRL := middleware.IPWindowRateLimit(loginRateMax, loginRateWin)
	registerRL := middleware.IPWindowRateLimit(5, 24*time.Hour)

	// authzMW builds a per-service AuthZ middleware when authzConn is available.
	// Falls back to portal-gate (user-type) only when authzConn is nil.
	authzMW := func(svcPrefix string, extractor ResourceExtractorFn) func(http.Handler) http.Handler {
		if authzConn != nil {
			return AuthZMiddleware(authzConn, svcPrefix, extractor)
		}
		return func(next http.Handler) http.Handler { return next }
	}

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		health := map[string]any{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
			"svc":    "gateway",
		}
		if clientManager != nil {
			health["services"] = clientManager.HealthCheck()
		}
		_ = json.NewEncoder(w).Encode(health)
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		// Same behavior for now; can become stricter later.
		w.Header().Set("Content-Type", "application/json")
		health := map[string]any{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
			"svc":    "gateway",
		}
		if clientManager != nil {
			health["services"] = clientManager.HealthCheck()
		}
		_ = json.NewEncoder(w).Encode(health)
	})

	// ── JWKS endpoint (public — RS256 public key for JWT verification) ──────
	// Served by authn handler; consumed by all services + external verifiers.
	if authnHandler != nil {
		mux.HandleFunc("GET /.well-known/jwks.json", authnHandler.JWKS)
	}

	if authnHandler != nil {
		// ── PUBLIC (no auth required) ───────────────────────────────────────
		mux.HandleFunc("POST /v1/auth/register", registerRL(http.HandlerFunc(authnHandler.Register)).ServeHTTP)
		mux.HandleFunc("POST /v1/auth/otp:send", otpRL(http.HandlerFunc(authnHandler.SendOTP)).ServeHTTP)
		mux.HandleFunc("POST /v1/auth/otp:verify", authnHandler.VerifyOTP)
		mux.HandleFunc("POST /v1/auth/otp:resend", otpRL(http.HandlerFunc(authnHandler.ResendOTP)).ServeHTTP)
		mux.HandleFunc("POST /v1/auth/login", loginRL(http.HandlerFunc(authnHandler.Login)).ServeHTTP)
		mux.HandleFunc("POST /v1/auth/token:refresh", authnHandler.RefreshToken)
		mux.HandleFunc("POST /v1/auth/token/refresh", authnHandler.RefreshToken)
		mux.HandleFunc("POST /v1/auth/token:validate", authnHandler.ValidateToken)
		mux.HandleFunc("POST /v1/auth/token/validate", authnHandler.ValidateToken)
		mux.HandleFunc("GET /v1/auth/csrf-token", authnHandler.GetCSRFToken)
		// Mobile password reset: send OTP to mobile then reset using OTP verify flow.
		// Email password reset: separate endpoint under /v1/auth/email/
		mux.HandleFunc("POST /v1/auth/password:reset-request", authnHandler.RequestPasswordResetByEmail)
		mux.HandleFunc("POST /v1/auth/password:reset", authnHandler.ResetPassword)
		mux.HandleFunc("POST /v1/auth/biometric:authenticate", authnHandler.BiometricAuthenticate)

		// Email auth public routes (business/system/partner/regulator portals)
		mux.HandleFunc("POST /v1/auth/email/register", registerRL(http.HandlerFunc(authnHandler.RegisterEmailUser)).ServeHTTP)
		mux.HandleFunc("POST /v1/auth/email/otp:send", authnHandler.SendEmailOTP)
		mux.HandleFunc("POST /v1/auth/email/verify", authnHandler.VerifyEmail)
		mux.HandleFunc("POST /v1/auth/email/login", loginRL(http.HandlerFunc(authnHandler.EmailLogin)).ServeHTTP)
		mux.HandleFunc("POST /v1/auth/email-password:login", loginRL(http.HandlerFunc(authnHandler.EmailPasswordLogin)).ServeHTTP)
		mux.HandleFunc("POST /v1/auth/email/password:reset-request", authnHandler.RequestPasswordResetByEmail)
		mux.HandleFunc("POST /v1/auth/email/password:reset", authnHandler.ResetPasswordByEmail)

		// ── PROTECTED — any authenticated user ─────────────────────────────
		mux.Handle("POST /v1/auth/logout", authMW(csrfMW(http.HandlerFunc(authnHandler.Logout))))
		mux.Handle("POST /v1/auth/csrf:validate", authMW(csrfMW(http.HandlerFunc(authnHandler.ValidateCSRF))))
		mux.Handle("POST /v1/auth/password:change", authMW(csrfMW(http.HandlerFunc(authnHandler.ChangePassword))))
		// Current session — support both singular (legacy) and plural (canonical) paths.
		mux.Handle("GET /v1/auth/session/current", authMW(http.HandlerFunc(authnHandler.GetCurrentSession)))
		mux.Handle("GET /v1/auth/sessions/current", authMW(http.HandlerFunc(authnHandler.GetCurrentSession)))
		mux.Handle("POST /v1/auth/users:find", authMW(SystemUserMiddleware(http.HandlerFunc(authnHandler.FindPortalUser))))
		mux.Handle("POST /v1/auth/users/{user_id}/password:temporary", authMW(SystemUserMiddleware(csrfMW(http.HandlerFunc(authnHandler.SetTemporaryPassword)))))
		mux.Handle("GET /v1/auth/sessions/{session_id}", authMW(http.HandlerFunc(authnHandler.GetSession)))
		mux.Handle("DELETE /v1/auth/sessions/{session_id}", authMW(csrfMW(http.HandlerFunc(authnHandler.RevokeSession))))
		// Session refresh — PUBLIC (no auth middleware):
		// The caller's access token may be expired, so authMW would reject with 401
		// before the refresh handler is reached.  Each branch validates independently:
		// - JWT sessions: RefreshToken validates the refresh_token itself.
		// - Server-side sessions: GetCurrentSession reads session_token from cookie metadata.
		mux.HandleFunc("POST /v1/auth/sessions/refresh", func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			_ = r.Body.Close()
			if len(body) > 0 && strings.Contains(string(body), "refresh_token") {
				r.Body = io.NopCloser(strings.NewReader(string(body)))
				authnHandler.RefreshToken(w, r)
				return
			}
			r.Body = io.NopCloser(strings.NewReader(string(body)))
			authnHandler.GetCurrentSession(w, r)
		})
		mux.Handle("GET /v1/auth/users/{user_id}/sessions", authMW(http.HandlerFunc(authnHandler.ListSessions)))
		mux.Handle("POST /v1/auth/users/{user_id}/sessions/revoke-all", authMW(csrfMW(http.HandlerFunc(authnHandler.RevokeAllSessions))))
		// User info endpoint (get own profile/basic info by user_id)
		mux.Handle("GET /v1/auth/users/{user_id}", authMW(http.HandlerFunc(authnHandler.GetUserProfile)))

		// Profile (any authenticated user)
		mux.Handle("POST /v1/auth/users/{user_id}/profile", authMW(http.HandlerFunc(authnHandler.CreateUserProfile)))
		mux.Handle("GET /v1/auth/users/{user_id}/profile", authMW(http.HandlerFunc(authnHandler.GetUserProfile)))
		mux.Handle("PATCH /v1/auth/users/{user_id}/profile", authMW(csrfMW(http.HandlerFunc(authnHandler.UpdateUserProfile))))
		mux.Handle("POST /v1/auth/users/{user_id}/profile/photo:upload-url", authMW(http.HandlerFunc(authnHandler.GetProfilePhotoUploadURL)))
		// BUG-010 FIX: Added GET endpoint for notification preferences (was missing, only PATCH existed)
		mux.Handle("GET /v1/auth/users/{user_id}/notification-preferences", authMW(http.HandlerFunc(authnHandler.GetNotificationPreferences)))
		mux.Handle("PATCH /v1/auth/users/{user_id}/notification-preferences", authMW(csrfMW(http.HandlerFunc(authnHandler.UpdateNotificationPreferences))))

		// TOTP / 2FA (any authenticated user)
		mux.Handle("POST /v1/auth/users/{user_id}/totp:enable", authMW(csrfMW(http.HandlerFunc(authnHandler.EnableTOTP))))
		// VerifyTOTP must be callable without AuthMiddleware for MFA step-up flow
		// (Login returns mfa_session_token before any auth token exists).
		mux.HandleFunc("POST /v1/auth/users/{user_id}/totp:verify", authnHandler.VerifyTOTP)
		mux.Handle("POST /v1/auth/users/{user_id}/totp:disable", authMW(csrfMW(http.HandlerFunc(authnHandler.DisableTOTP))))

		// Documents (any authenticated user — authz enforces finer rules)
		authzDoc := authzMW("svc:document", PathSegmentExtractor("/v1/auth/"))
		mux.Handle("POST /v1/auth/users/{user_id}/documents", authMW(authzDoc(http.HandlerFunc(authnHandler.UploadUserDocument))))
		mux.Handle("GET /v1/auth/users/{user_id}/documents", authMW(authzDoc(http.HandlerFunc(authnHandler.ListUserDocuments))))
		mux.Handle("GET /v1/auth/documents/{user_document_id}", authMW(authzDoc(http.HandlerFunc(authnHandler.GetUserDocument))))
		mux.Handle("PATCH /v1/auth/documents/{user_document_id}", authMW(csrfMW(authzDoc(http.HandlerFunc(authnHandler.UpdateUserDocument)))))
		mux.Handle("DELETE /v1/auth/documents/{user_document_id}", authMW(csrfMW(authzDoc(http.HandlerFunc(authnHandler.DeleteUserDocument)))))
		mux.Handle("GET /v1/auth/document-types", authMW(authzDoc(http.HandlerFunc(authnHandler.ListDocumentTypes))))

		// ── PROTECTED — system/agent portal only (portal-gate + AuthZ) ─────
		// KYC admin actions: agent or system user + Casbin policy check
		authzKYC := authzMW("svc:kyc", PathSegmentExtractor("/v1/auth/"))
		// KYC self-service routes: any authenticated user can initiate/submit/complete
		// their own KYC. B2B org admins need these to complete eKYC on first login.
		// No Casbin authz check — AnyAuthenticatedMiddleware is the only gate needed
		// since the user can only act on their own KYC (user_id enforced by authn).
		mux.Handle("POST /v1/auth/users/{user_id}/kyc", authMW(AnyAuthenticatedMiddleware(http.HandlerFunc(authnHandler.InitiateKYC))))
		mux.Handle("GET /v1/auth/users/{user_id}/kyc", authMW(AnyAuthenticatedMiddleware(http.HandlerFunc(authnHandler.GetKYCStatus))))
		mux.Handle("POST /v1/auth/users/{user_id}/kyc:submit-frame", authMW(AnyAuthenticatedMiddleware(http.HandlerFunc(authnHandler.SubmitKYCFrame))))
		mux.Handle("POST /v1/auth/users/{user_id}/kyc:complete", authMW(AnyAuthenticatedMiddleware(http.HandlerFunc(authnHandler.CompleteKYCSession))))
		mux.Handle("POST /v1/auth/kyc/{kyc_id}/approve", authMW(SystemUserMiddleware(authzKYC(csrfMW(http.HandlerFunc(authnHandler.ApproveKYC))))))
		mux.Handle("POST /v1/auth/kyc/{kyc_id}/reject", authMW(SystemUserMiddleware(authzKYC(csrfMW(http.HandlerFunc(authnHandler.RejectKYC))))))
		mux.Handle("POST /v1/auth/documents/{user_document_id}/verify", authMW(AgentOrSystemMiddleware(authzKYC(csrfMW(http.HandlerFunc(authnHandler.VerifyDocument))))))

		// Voice sessions (agent or system only)
		authzVoice := authzMW("svc:voice", PathSegmentExtractor("/v1/auth/"))
		mux.Handle("POST /v1/auth/voice-sessions", authMW(AgentOrSystemMiddleware(authzVoice(http.HandlerFunc(authnHandler.CreateVoiceSession)))))
		mux.Handle("GET /v1/auth/voice-sessions/{voice_session_id}", authMW(AgentOrSystemMiddleware(authzVoice(http.HandlerFunc(authnHandler.GetVoiceSession)))))
		mux.Handle("POST /v1/auth/voice-sessions/{voice_session_id}/end", authMW(AgentOrSystemMiddleware(authzVoice(csrfMW(http.HandlerFunc(authnHandler.EndVoiceSession))))))

		// API key management (system or partner)
		authzAPIKey := authzMW("svc:apikey", PathSegmentExtractor("/v1/auth/"))
		mux.Handle("POST /v1/auth/api-keys", authMW(csrfMW(authzAPIKey(http.HandlerFunc(authnHandler.CreateAPIKey)))))
		mux.Handle("GET /v1/auth/api-keys", authMW(authzAPIKey(http.HandlerFunc(authnHandler.ListAPIKeys))))
		mux.Handle("POST /v1/auth/api-keys/{key_id}/revoke", authMW(csrfMW(authzAPIKey(http.HandlerFunc(authnHandler.RevokeAPIKey)))))
	}

	if paymentConn != nil {
		paymentCallbackHandler := handlers.NewPaymentCallbackHandler(paymentConn)
		mux.HandleFunc("POST /v1/payments/webhook/sslcommerz", paymentCallbackHandler.Webhook)
		mux.HandleFunc("POST /v1/payments/sslcommerz/success", paymentCallbackHandler.Success)
		mux.HandleFunc("POST /v1/payments/sslcommerz/fail", paymentCallbackHandler.Fail)
		mux.HandleFunc("POST /v1/payments/sslcommerz/cancel", paymentCallbackHandler.Cancel)
		mux.HandleFunc("GET /v1/payments/sslcommerz/success", paymentCallbackHandler.Success)
		mux.HandleFunc("GET /v1/payments/sslcommerz/fail", paymentCallbackHandler.Fail)
		mux.HandleFunc("GET /v1/payments/sslcommerz/cancel", paymentCallbackHandler.Cancel)
	}

	// ── B2B APIs (auth + authz) ───────────────────────────────────────────────
	b2bConn := getServiceConn(clientManager, "b2b")
	authnConn = getServiceConn(clientManager, "authn")
	authzConn = getServiceConn(clientManager, "authz")

	if b2bConn != nil && authnConn != nil && authzConn != nil {
		authnClient := authnv1.NewAuthServiceClient(authnConn)
		authzClient := authzv1.NewAuthZServiceClient(authzConn)
		b2bHandler := handlers.NewB2BServiceHandler(b2bConn, authnClient, authzClient)
		b2bContext := middleware.NewB2BContextMiddleware(b2bConn)
		authzB2B := authzMW("svc:b2b", PathSegmentExtractor("/v1/b2b/"))

		// Employee public self-service bootstrap
		mux.HandleFunc("GET /v1/b2b/organisations:employee-login", b2bHandler.ListEmployeeLoginOrganisations)
		mux.HandleFunc("POST /v1/b2b/employees:activate", b2bHandler.ActivateEmployee)

		// Purchase Orders
		mux.Handle("GET /v1/b2b/purchase-orders/catalog", authMW(authzB2B(http.HandlerFunc(b2bHandler.ListPurchaseOrderCatalog))))
		mux.Handle("GET /v1/b2b/purchase-orders", authMW(b2bContext.InjectOrganisationContext(authzB2B(http.HandlerFunc(b2bHandler.ListPurchaseOrders)))))
		mux.Handle("GET /v1/b2b/purchase-orders/{purchase_order_id}", authMW(authzB2B(http.HandlerFunc(b2bHandler.GetPurchaseOrder))))
		mux.Handle("POST /v1/b2b/purchase-orders", authMW(csrfMW(authzB2B(http.HandlerFunc(b2bHandler.CreatePurchaseOrder)))))
		// UpdatePurchaseOrder / DeletePurchaseOrder RPCs are not yet defined in the proto.
		// Return 501 so clients get a clear "not implemented" instead of a catch-all 404.
		mux.Handle("PATCH /v1/b2b/purchase-orders/{purchase_order_id}", authMW(http.HandlerFunc(notImplementedHandler)))
		mux.Handle("DELETE /v1/b2b/purchase-orders/{purchase_order_id}", authMW(http.HandlerFunc(notImplementedHandler)))

		// Departments (full CRUD)
		mux.Handle("GET /v1/b2b/departments", authMW(b2bContext.InjectOrganisationContext(authzB2B(http.HandlerFunc(b2bHandler.ListDepartments)))))
		mux.Handle("GET /v1/b2b/departments/{department_id}", authMW(authzB2B(http.HandlerFunc(b2bHandler.GetDepartment))))
		mux.Handle("POST /v1/b2b/departments", authMW(csrfMW(authzB2B(http.HandlerFunc(b2bHandler.CreateDepartment)))))
		mux.Handle("PATCH /v1/b2b/departments/{department_id}", authMW(csrfMW(authzB2B(http.HandlerFunc(b2bHandler.UpdateDepartment)))))
		mux.Handle("DELETE /v1/b2b/departments/{department_id}", authMW(csrfMW(authzB2B(http.HandlerFunc(b2bHandler.DeleteDepartment)))))

		// Employees (full CRUD + bulk upload)
		mux.Handle("GET /v1/b2b/employees", authMW(b2bContext.InjectOrganisationContext(authzB2B(http.HandlerFunc(b2bHandler.ListEmployees)))))
		mux.Handle("GET /v1/b2b/employees/{employee_uuid}", authMW(authzB2B(http.HandlerFunc(b2bHandler.GetEmployee))))
		mux.Handle("POST /v1/b2b/employees", authMW(csrfMW(authzB2B(http.HandlerFunc(b2bHandler.CreateEmployee)))))
		mux.Handle("PATCH /v1/b2b/employees/{employee_uuid}", authMW(csrfMW(authzB2B(http.HandlerFunc(b2bHandler.UpdateEmployee)))))
		mux.Handle("DELETE /v1/b2b/employees/{employee_uuid}", authMW(csrfMW(authzB2B(http.HandlerFunc(b2bHandler.DeleteEmployee)))))
		// Bulk upload: POST /v1/b2b/employees:bulkUpload  (multipart/form-data)
		mux.Handle("POST /v1/b2b/employees/bulk-upload", authMW(csrfMW(b2bContext.InjectOrganisationContext(authzB2B(http.HandlerFunc(b2bHandler.BulkUploadEmployees))))))

		// Employee self-service
		// These routes are already self-scoped by the B2B service using x-user-id,
		// so gateway-level Casbin adds brittle org-domain coupling without improving safety.
		// Keeping authMW only lets beneficiary sessions bootstrap before portal_biz_id exists.
		mux.Handle("GET /v1/b2b-self/profile", authMW(http.HandlerFunc(b2bHandler.GetMyEmployeeProfile)))
		mux.Handle("GET /v1/b2b-self/coverage", authMW(http.HandlerFunc(b2bHandler.GetMyEmployeeCoverage)))

		// Organisations (full CRUD + members)
		mux.Handle("GET /v1/b2b/organisations", authMW(authzB2B(http.HandlerFunc(b2bHandler.ListOrganisations))))
		mux.Handle("GET /v1/b2b/organisations/me", authMW(http.HandlerFunc(b2bHandler.ResolveMyOrganisation)))
		mux.Handle("GET /v1/b2b/organisations/{organisation_id}", authMW(authzB2B(http.HandlerFunc(b2bHandler.GetOrganisation))))
		mux.Handle("POST /v1/b2b/organisations", authMW(csrfMW(authzB2B(http.HandlerFunc(b2bHandler.CreateOrganisation)))))
		mux.Handle("PATCH /v1/b2b/organisations/{organisation_id}", authMW(csrfMW(authzB2B(http.HandlerFunc(b2bHandler.UpdateOrganisation)))))
		mux.Handle("DELETE /v1/b2b/organisations/{organisation_id}", authMW(csrfMW(authzB2B(http.HandlerFunc(b2bHandler.DeleteOrganisation)))))
		mux.Handle("GET /v1/b2b/organisations/{organisation_id}/members", authMW(authzB2B(http.HandlerFunc(b2bHandler.ListOrgMembers))))
		mux.Handle("POST /v1/b2b/organisations/{organisation_id}/members", authMW(csrfMW(authzB2B(http.HandlerFunc(b2bHandler.AddOrgMember)))))
		mux.Handle("POST /v1/b2b/organisations/{organisation_id}/admins", authMW(csrfMW(authzB2B(http.HandlerFunc(b2bHandler.AssignOrgAdmin)))))
		mux.Handle("POST /v1/b2b/organisations/{organisation_id}/admins:assign", authMW(csrfMW(authzB2B(http.HandlerFunc(b2bHandler.AssignOrgAdmin)))))
		mux.Handle("DELETE /v1/b2b/organisations/{organisation_id}/members/{member_id}", authMW(csrfMW(authzB2B(http.HandlerFunc(b2bHandler.RemoveOrgMember)))))
	}

	// ── Media APIs (auth + authz) ─────────────────────────────────────────────
	if mediaConn := getServiceConn(clientManager, "media"); mediaConn != nil {
		mediaHandler := handlers.NewMediaHandler(mediaConn)
		authzMedia := authzMW("svc:media", PathSegmentExtractor("/v1/"))

		mux.Handle("POST /v1/media", authMW(authzMedia(http.HandlerFunc(mediaHandler.Upload))))
		mux.Handle("GET /v1/media/{media_id}", authMW(authzMedia(http.HandlerFunc(mediaHandler.Get))))
		mux.Handle("GET /v1/entities/{entity_type}/{entity_id}/media", authMW(authzMedia(http.HandlerFunc(mediaHandler.List))))
		mux.Handle("GET /v1/media/{media_id}/download", authMW(authzMedia(http.HandlerFunc(mediaHandler.Download))))
		mux.Handle("GET /v1/media/{media_id}/optimized", authMW(authzMedia(http.HandlerFunc(mediaHandler.DownloadOptimized))))
		mux.Handle("GET /v1/media/{media_id}/thumbnail", authMW(authzMedia(http.HandlerFunc(mediaHandler.DownloadThumbnail))))
		mux.Handle("DELETE /v1/media/{media_id}", authMW(csrfMW(authzMedia(http.HandlerFunc(mediaHandler.Delete)))))
		mux.Handle("POST /v1/media/{media_id}/validate", authMW(csrfMW(authzMedia(http.HandlerFunc(mediaHandler.Validate)))))
		mux.Handle("POST /v1/media/{media_id}/process", authMW(csrfMW(authzMedia(http.HandlerFunc(mediaHandler.RequestProcessing)))))
		mux.Handle("GET /v1/processing-jobs/{job_id}", authMW(authzMedia(http.HandlerFunc(mediaHandler.GetProcessingJob))))
		mux.Handle("GET /v1/processing-jobs", authMW(authzMedia(http.HandlerFunc(mediaHandler.ListProcessingJobs))))
	}

	// ── Storage APIs (auth + authz) ───────────────────────────────────────────
	if storageConn := getServiceConn(clientManager, "storage"); storageConn != nil {
		documentHandler := handlers.NewDocumentHandler(storageConn)
		authzStorage := authzMW("svc:storage", StorageResourceExtractor())

		mux.Handle("POST /v1/storage/files", authMW(authzStorage(http.HandlerFunc(documentHandler.Upload))))
		mux.Handle("POST /v1/storage/files:batch", authMW(authzStorage(http.HandlerFunc(documentHandler.UploadBatch))))
		mux.Handle("POST /v1/storage/files:upload-url", authMW(authzStorage(http.HandlerFunc(documentHandler.GetUploadURL))))
		mux.Handle("POST /v1/storage/files:finalize", authMW(csrfMW(authzStorage(http.HandlerFunc(documentHandler.FinalizeUpload)))))
		mux.Handle("GET /v1/storage/files", authMW(authzStorage(http.HandlerFunc(documentHandler.List))))
		mux.Handle("GET /v1/storage/files/{id}", authMW(authzStorage(http.HandlerFunc(documentHandler.Get))))
		mux.Handle("PATCH /v1/storage/files/{id}", authMW(csrfMW(authzStorage(http.HandlerFunc(documentHandler.Update)))))
		mux.Handle("GET /v1/storage/files/{id}/download-url", authMW(authzStorage(http.HandlerFunc(documentHandler.GetDownloadURL))))
		mux.Handle("POST /v1/storage/files/{id}/download-url", authMW(authzStorage(http.HandlerFunc(documentHandler.GetDownloadURL))))
		mux.Handle("DELETE /v1/storage/files/{id}", authMW(csrfMW(authzStorage(http.HandlerFunc(documentHandler.Delete)))))
	}

	// ── Partner APIs (auth + authz) ───────────────────────────────────────────
	if partnerConn := getServiceConn(clientManager, "partner"); partnerConn != nil {
		partnerHandler := handlers.NewPartnerHandler(partnerConn)
		authzPartner := authzMW("svc:partner", PathSegmentExtractor("/v1/"))

		mux.Handle("POST /v1/partners", authMW(csrfMW(authzPartner(http.HandlerFunc(partnerHandler.Create)))))
		mux.Handle("GET /v1/partners", authMW(authzPartner(http.HandlerFunc(partnerHandler.List))))
		mux.Handle("GET /v1/partners/{partner_id}", authMW(authzPartner(http.HandlerFunc(partnerHandler.Get))))
		mux.Handle("PATCH /v1/partners/{partner_id}", authMW(csrfMW(authzPartner(http.HandlerFunc(partnerHandler.Update)))))
		mux.Handle("DELETE /v1/partners/{partner_id}", authMW(csrfMW(authzPartner(http.HandlerFunc(partnerHandler.Delete)))))
		mux.Handle("POST /v1/partners/{partner_id}/verify", authMW(csrfMW(authzPartner(http.HandlerFunc(partnerHandler.Verify)))))
		mux.Handle("POST /v1/partners/{partner_id}/updateStatus", authMW(csrfMW(authzPartner(http.HandlerFunc(partnerHandler.UpdateStatus)))))
		mux.Handle("GET /v1/partners/{partner_id}/commission", authMW(authzPartner(http.HandlerFunc(partnerHandler.GetCommission))))
		mux.Handle("PUT /v1/partners/{partner_id}/commission", authMW(csrfMW(authzPartner(http.HandlerFunc(partnerHandler.UpdateCommission)))))
		mux.Handle("GET /v1/partners/{partner_id}/credentials", authMW(authzPartner(http.HandlerFunc(partnerHandler.GetCredentials))))
		mux.Handle("POST /v1/partners/{partner_id}/credentials:rotate", authMW(csrfMW(authzPartner(http.HandlerFunc(partnerHandler.RotateAPIKey)))))
	}

	// ── Fraud APIs (auth + authz) ─────────────────────────────────────────────
	if fraudConn := getServiceConn(clientManager, "fraud"); fraudConn != nil {
		fraudHandler := handlers.NewFraudHandler(fraudConn)
		authzFraud := authzMW("svc:fraud", PathSegmentExtractor("/v1/"))

		mux.Handle("POST /v1/fraud-checks", authMW(csrfMW(authzFraud(http.HandlerFunc(fraudHandler.Check)))))
		mux.Handle("GET /v1/fraud-alerts", authMW(authzFraud(http.HandlerFunc(fraudHandler.ListAlerts))))
		mux.Handle("GET /v1/fraud-alerts/{fraud_alert_id}", authMW(authzFraud(http.HandlerFunc(fraudHandler.GetAlert))))
		mux.Handle("POST /v1/fraud-cases", authMW(csrfMW(authzFraud(http.HandlerFunc(fraudHandler.CreateCase)))))
		mux.Handle("GET /v1/fraud-cases/{fraud_case_id}", authMW(authzFraud(http.HandlerFunc(fraudHandler.GetCase))))
		mux.Handle("PATCH /v1/fraud-cases/{fraud_case_id}", authMW(csrfMW(authzFraud(http.HandlerFunc(fraudHandler.UpdateCase)))))
		mux.Handle("GET /v1/fraud-rules", authMW(authzFraud(http.HandlerFunc(fraudHandler.ListRules))))
		mux.Handle("POST /v1/fraud-rules", authMW(csrfMW(authzFraud(http.HandlerFunc(fraudHandler.CreateRule)))))
		mux.Handle("PATCH /v1/fraud-rules/{rule_id}", authMW(csrfMW(authzFraud(http.HandlerFunc(fraudHandler.UpdateRule)))))
		mux.Handle("POST /v1/fraud-rules/{rule_id}/activate", authMW(csrfMW(authzFraud(http.HandlerFunc(fraudHandler.ActivateRule)))))
		mux.Handle("POST /v1/fraud-rules/{rule_id}/deactivate", authMW(csrfMW(authzFraud(http.HandlerFunc(fraudHandler.DeactivateRule)))))
	}

	if dlrHandler != nil {
		mux.Handle("POST /v1/internal/sms/dlr", dlrHandler)
	}

	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	// 🛡️  PoliSync (C# .NET 8) — Insurance Commerce & Policy Engine
	// JWT validated by this gateway; identity forwarded as X-* headers → C# services
	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

	// 📦 Product APIs (product-service :50121)
	if productConn := getServiceConn(clientManager, "product"); productConn != nil {
		productHandler := handlers.NewPoliSyncHandler(productConn, "product-service")
		authzProduct := authzMW("svc:product", PathSegmentExtractor("/v1/"))

		mux.Handle("GET /v1/products", authMW(authzProduct(productHandler.Proxy())))
		mux.Handle("POST /v1/products", authMW(csrfMW(authzProduct(productHandler.Proxy()))))
		mux.Handle("GET /v1/products/{product_id}", authMW(authzProduct(productHandler.Proxy())))
		mux.Handle("PATCH /v1/products/{product_id}", authMW(csrfMW(authzProduct(productHandler.Proxy()))))
		mux.Handle("POST /v1/products/{product_id}/activate", authMW(csrfMW(authzProduct(productHandler.Proxy()))))
		mux.Handle("POST /v1/products/{product_id}/deactivate", authMW(csrfMW(authzProduct(productHandler.Proxy()))))
		mux.Handle("GET /v1/products/{product_id}/plans", authMW(authzProduct(productHandler.Proxy())))
		mux.Handle("POST /v1/products/{product_id}/plans", authMW(csrfMW(authzProduct(productHandler.Proxy()))))
		mux.Handle("GET /v1/products/{product_id}/plans/{plan_id}", authMW(authzProduct(productHandler.Proxy())))
		mux.Handle("GET /v1/products/{product_id}/riders", authMW(authzProduct(productHandler.Proxy())))
		mux.Handle("POST /v1/products/{product_id}/riders", authMW(csrfMW(authzProduct(productHandler.Proxy()))))
		mux.Handle("POST /v1/products/{product_id}/pricing", authMW(csrfMW(authzProduct(productHandler.Proxy()))))
		mux.Handle("GET /v1/products/{product_id}/pricing", authMW(authzProduct(productHandler.Proxy())))
		mux.Handle("POST /v1/premium:calculate", authMW(productHandler.Proxy()))
	}

	// 💬 Quotation APIs (quote-service :50131)
	if quoteConn := getServiceConn(clientManager, "quote"); quoteConn != nil {
		quoteHandler := handlers.NewPoliSyncHandler(quoteConn, "quote-service")
		authzQuote := authzMW("svc:quote", PathSegmentExtractor("/v1/"))

		mux.Handle("POST /v1/quotations", authMW(authzQuote(quoteHandler.Proxy())))
		mux.Handle("GET /v1/quotations", authMW(authzQuote(quoteHandler.Proxy())))
		mux.Handle("GET /v1/quotations/{quotation_id}", authMW(authzQuote(quoteHandler.Proxy())))
		mux.Handle("PATCH /v1/quotations/{quotation_id}", authMW(csrfMW(authzQuote(quoteHandler.Proxy()))))
		mux.Handle("POST /v1/quotations/{quotation_id}/submit", authMW(csrfMW(authzQuote(quoteHandler.Proxy()))))
		mux.Handle("POST /v1/quotations/{quotation_id}/approve", authMW(csrfMW(authzQuote(quoteHandler.Proxy()))))
		mux.Handle("POST /v1/quotations/{quotation_id}/reject", authMW(csrfMW(authzQuote(quoteHandler.Proxy()))))
	}

	// 🛒 Order APIs (order-service :50141)
	if orderConn := getServiceConn(clientManager, "order"); orderConn != nil {
		orderHandler := handlers.NewPoliSyncHandler(orderConn, "order-service")
		authzOrder := authzMW("svc:order", PathSegmentExtractor("/v1/"))

		mux.Handle("POST /v1/orders", authMW(authzOrder(orderHandler.Proxy())))
		mux.Handle("GET /v1/orders", authMW(authzOrder(orderHandler.Proxy())))
		mux.Handle("GET /v1/orders/{order_id}", authMW(authzOrder(orderHandler.Proxy())))
		mux.Handle("POST /v1/orders/{order_id}/initiate-payment", authMW(csrfMW(authzOrder(orderHandler.Proxy()))))
		mux.Handle("POST /v1/orders/{order_id}/confirm", authMW(csrfMW(authzOrder(orderHandler.Proxy()))))
		mux.Handle("POST /v1/orders/{order_id}/cancel", authMW(csrfMW(authzOrder(orderHandler.Proxy()))))
	}

	// Payment APIs (payment-service :50190) — gRPC handler (BUG-006 FIX)
	// Previously used PoliSyncHandler (HTTP proxy) which fails because payment is gRPC-only.
	if paymentConn := getServiceConn(clientManager, "payment"); paymentConn != nil {
		paymentHandler := handlers.NewPaymentHandler(paymentConn)
		authzPayment := authzMW("svc:payment", PathSegmentExtractor("/v1/"))

		// Core payment CRUD
		mux.Handle("POST /v1/payments", authMW(csrfMW(authzPayment(http.HandlerFunc(paymentHandler.Initiate)))))
		mux.Handle("GET /v1/payments", authMW(authzPayment(http.HandlerFunc(paymentHandler.List))))
		mux.Handle("GET /v1/payments/{payment_id}", authMW(authzPayment(http.HandlerFunc(paymentHandler.Get))))
		mux.Handle("POST /v1/payments/{payment_id}/verify", authMW(csrfMW(authzPayment(http.HandlerFunc(paymentHandler.Verify)))))
		mux.Handle("POST /v1/payments/{payment_id}/refunds", authMW(csrfMW(authzPayment(http.HandlerFunc(paymentHandler.InitiateRefund)))))
		mux.Handle("GET /v1/refunds/{refund_id}", authMW(authzPayment(http.HandlerFunc(paymentHandler.GetRefund))))
		mux.Handle("GET /v1/users/{user_id}/payment-methods", authMW(authzPayment(http.HandlerFunc(paymentHandler.ListMethods))))
		mux.Handle("POST /v1/users/{user_id}/payment-methods", authMW(csrfMW(authzPayment(http.HandlerFunc(paymentHandler.AddMethod)))))
		mux.Handle("POST /v1/payments/reconcile", authMW(csrfMW(authzPayment(http.HandlerFunc(paymentHandler.Reconcile)))))

		// Provider reference lookup (admin/agent)
		mux.Handle("GET /v1/payments/provider/{provider}/references/{provider_reference}", authMW(authzPayment(http.HandlerFunc(paymentHandler.GetByProviderRef))))

		// Manual payment proof — customer submits bank transfer proof, admin reviews
		mux.Handle("POST /v1/payments/{payment_id}/submit-proof", authMW(csrfMW(authzPayment(http.HandlerFunc(paymentHandler.SubmitProof)))))
		mux.Handle("POST /v1/payments/{payment_id}/review", authMW(csrfMW(authzPayment(http.HandlerFunc(paymentHandler.Review)))))

		// Receipt generation and retrieval
		mux.Handle("POST /v1/payments/{payment_id}/generate-receipt", authMW(csrfMW(authzPayment(http.HandlerFunc(paymentHandler.GenerateReceipt)))))
		mux.Handle("GET /v1/payments/{payment_id}/receipt", authMW(authzPayment(http.HandlerFunc(paymentHandler.GetReceipt))))
	}

	// 🧾 Billing APIs (billing-service :50195)
	if billingConn := getServiceConn(clientManager, "billing"); billingConn != nil {
		billingHandler := handlers.NewPoliSyncHandler(billingConn, "billing-service")
		authzBilling := authzMW("svc:billing", PathSegmentExtractor("/v1/"))

		mux.Handle("POST /v1/invoices", authMW(csrfMW(authzBilling(billingHandler.Proxy()))))
		mux.Handle("GET /v1/invoices", authMW(authzBilling(billingHandler.Proxy())))
		mux.Handle("GET /v1/invoices/{invoice_id}", authMW(authzBilling(billingHandler.Proxy())))
		mux.Handle("POST /v1/invoices/{invoice_id}/mark-paid", authMW(csrfMW(authzBilling(billingHandler.Proxy()))))
		mux.Handle("POST /v1/invoices/{invoice_id}/cancel", authMW(csrfMW(authzBilling(billingHandler.Proxy()))))
		mux.Handle("POST /v1/invoices/{invoice_id}/issue", authMW(csrfMW(authzBilling(billingHandler.Proxy()))))
		mux.Handle("GET /v1/invoices/{invoice_id}/pdf", authMW(authzBilling(billingHandler.Proxy())))
		mux.Handle("POST /v1/invoices/{invoice_id}/generate-pdf", authMW(csrfMW(authzBilling(billingHandler.Proxy()))))
		mux.Handle("GET /v1/orders/{order_id}/invoice", authMW(authzBilling(billingHandler.Proxy())))
	}

	// 📋 Policy APIs (policy-service :50161) — incl. Endorsement + Renewal
	if policyConn := getServiceConn(clientManager, "policy"); policyConn != nil {
		policyHandler := handlers.NewPoliSyncHandler(policyConn, "policy-service")
		authzPolicy := authzMW("svc:policy", PathSegmentExtractor("/v1/"))
		policyProxy := policyHandler.Proxy()

		// Keep canonical upstream path as /v1/policies/number/{policy_number}.
		// We register a non-conflicting mux path (with /lookup suffix) and strip it before proxying.
		policyNumberLookupProxy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/lookup") {
				r2 := r.Clone(r.Context())
				r2.URL.Path = strings.TrimSuffix(r.URL.Path, "/lookup")
				r2.URL.RawPath = r2.URL.Path
				policyProxy.ServeHTTP(w, r2)
				return
			}
			policyProxy.ServeHTTP(w, r)
		})

		mux.Handle("POST /v1/policies", authMW(csrfMW(authzPolicy(policyProxy))))
		mux.Handle("GET /v1/policies", authMW(authzPolicy(policyProxy)))
		mux.Handle("GET /v1/policies/{policy_id}", authMW(authzPolicy(policyProxy)))
		mux.Handle("GET /v1/policies/number/{policy_number}/lookup", authMW(authzPolicy(policyNumberLookupProxy)))
		mux.Handle("POST /v1/policies/{policy_id}/cancel", authMW(csrfMW(authzPolicy(policyProxy))))
		mux.Handle("POST /v1/policies/{policy_id}/suspend", authMW(csrfMW(authzPolicy(policyProxy))))
		mux.Handle("POST /v1/policies/{policy_id}/reinstate", authMW(csrfMW(authzPolicy(policyProxy))))
		mux.Handle("GET /v1/policies/{policy_id}/document", authMW(authzPolicy(policyProxy)))
		mux.Handle("GET /v1/policies/{policy_id}/nominees", authMW(authzPolicy(policyProxy)))
		mux.Handle("POST /v1/policies/{policy_id}/nominees", authMW(csrfMW(authzPolicy(policyProxy))))
		mux.Handle("PATCH /v1/policies/{policy_id}/nominees/{nominee_id}", authMW(csrfMW(authzPolicy(policyProxy))))
		mux.Handle("DELETE /v1/policies/{policy_id}/nominees/{nominee_id}", authMW(csrfMW(authzPolicy(policyProxy))))
		// Endorsements (co-hosted on policy-service)
		mux.Handle("POST /v1/policies/{policy_id}/endorsements", authMW(csrfMW(authzPolicy(policyProxy))))
		mux.Handle("GET /v1/policies/{policy_id}/endorsements", authMW(authzPolicy(policyProxy)))
		mux.Handle("GET /v1/endorsements/{endorsement_id}", authMW(authzPolicy(policyProxy)))
		mux.Handle("POST /v1/endorsements/{endorsement_id}/approve", authMW(csrfMW(authzPolicy(policyProxy))))
		mux.Handle("POST /v1/endorsements/{endorsement_id}/reject", authMW(csrfMW(authzPolicy(policyProxy))))
		// Renewal (co-hosted on policy-service)
		mux.Handle("GET /v1/policies/{policy_id}/renewal", authMW(authzPolicy(policyProxy)))
		mux.Handle("POST /v1/policies/{policy_id}/renewal/process", authMW(csrfMW(authzPolicy(policyProxy))))
	}

	// 🏥 Underwriting APIs (underwriting-service :50171)
	if uwConn := getServiceConn(clientManager, "underwriting"); uwConn != nil {
		uwHandler := handlers.NewPoliSyncHandler(uwConn, "underwriting-service")
		authzUW := authzMW("svc:underwriting", PathSegmentExtractor("/v1/"))

		mux.Handle("POST /v1/health-declarations", authMW(csrfMW(authzUW(uwHandler.Proxy()))))
		mux.Handle("GET /v1/health-declarations/{declaration_id}", authMW(authzUW(uwHandler.Proxy())))
		mux.Handle("GET /v1/quotations/{quotation_id}/underwriting-decision", authMW(authzUW(uwHandler.Proxy())))
		mux.Handle("POST /v1/underwriting/risk-score", authMW(authzUW(uwHandler.Proxy())))
	}

	// 🏛️  Claims APIs (claim-service :50211)
	if claimConn := getServiceConn(clientManager, "claim"); claimConn != nil {
		claimHandler := handlers.NewPoliSyncHandler(claimConn, "claim-service")
		authzClaim := authzMW("svc:claim", PathSegmentExtractor("/v1/"))

		mux.Handle("POST /v1/claims", authMW(csrfMW(authzClaim(claimHandler.Proxy()))))
		mux.Handle("GET /v1/claims", authMW(authzClaim(claimHandler.Proxy())))
		mux.Handle("GET /v1/claims/{claim_id}", authMW(authzClaim(claimHandler.Proxy())))
		mux.Handle("POST /v1/claims/{claim_id}/documents", authMW(csrfMW(authzClaim(claimHandler.Proxy()))))
		mux.Handle("POST /v1/claims/{claim_id}/review", authMW(csrfMW(authzClaim(claimHandler.Proxy()))))
		mux.Handle("POST /v1/claims/{claim_id}/approve", authMW(csrfMW(authzClaim(claimHandler.Proxy()))))
		mux.Handle("POST /v1/claims/{claim_id}/reject", authMW(csrfMW(authzClaim(claimHandler.Proxy()))))
		mux.Handle("POST /v1/claims/{claim_id}/settle", authMW(csrfMW(authzClaim(claimHandler.Proxy()))))
		mux.Handle("GET /v1/claims/{claim_id}/settlement", authMW(authzClaim(claimHandler.Proxy())))
	}

	// 💰 Commission APIs (commission-service :50151)
	if commConn := getServiceConn(clientManager, "commission"); commConn != nil {
		commHandler := handlers.NewPoliSyncHandler(commConn, "commission-service")
		authzComm := authzMW("svc:commission", PathSegmentExtractor("/v1/"))

		mux.Handle("GET /v1/commission/configs", authMW(authzComm(commHandler.Proxy())))
		mux.Handle("POST /v1/commission/configs", authMW(csrfMW(authzComm(commHandler.Proxy()))))
		mux.Handle("GET /v1/commission/payouts", authMW(authzComm(commHandler.Proxy())))
		mux.Handle("GET /v1/commission/payouts/{payout_id}", authMW(authzComm(commHandler.Proxy())))
		mux.Handle("GET /v1/commission/revenue-shares", authMW(authzComm(commHandler.Proxy())))
		mux.Handle("GET /v1/commission/summary", authMW(authzComm(commHandler.Proxy())))

		// Canonical commission paths from ENDPOINT_MAP
		mux.Handle("GET /v1/commissions", authMW(authzComm(commHandler.Proxy())))
		mux.Handle("GET /v1/commissions/{commission_id}", authMW(authzComm(commHandler.Proxy())))
		mux.Handle("POST /v1/commissions/calculate", authMW(csrfMW(authzComm(commHandler.Proxy()))))
		mux.Handle("POST /v1/commission-payouts", authMW(csrfMW(authzComm(commHandler.Proxy()))))
		mux.Handle("POST /v1/commission-payouts/{payout_id}/process", authMW(csrfMW(authzComm(commHandler.Proxy()))))
		mux.Handle("GET /v1/recipients/{recipient_id}/commission-statement", authMW(authzComm(commHandler.Proxy())))
	}

	// ── /v1/me — identity introspection (any authenticated user) ─────────────
	mux.Handle("GET /v1/me", authMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"user_id":      r.Header.Get("X-User-ID"),
			"session_id":   r.Header.Get("X-Session-ID"),
			"session_type": r.Header.Get("X-Session-Type"),
			"user_type":    r.Header.Get("X-User-Type"),
			"portal":       r.Header.Get("X-Portal"),
			"tenant_id":    r.Header.Get("X-Tenant-ID"),
		})
	})))

	// ── AuthZ APIs (authz-service — standalone gRPC, NOT a PoliSync HTTP proxy) ──
	// AuthZ is isolated and accessible by ALL portals (B2B, B2C, system, agent, regulator).
	// Uses AuthZHandler which calls authz gRPC methods directly.
	if authzSvcConn := getServiceConn(clientManager, "authz"); authzSvcConn != nil {
		authzSvcHandler := handlers.NewAuthZHandler(authzSvcConn)
		authzAuthz := authzMW("svc:authz", PathSegmentExtractor("/v1/authz/"))

		// Roles
		mux.Handle("GET /v1/authz/roles", authMW(authzAuthz(http.HandlerFunc(authzSvcHandler.ListRoles))))
		mux.Handle("POST /v1/authz/roles", authMW(csrfMW(authzAuthz(http.HandlerFunc(authzSvcHandler.CreateRole)))))
		mux.Handle("GET /v1/authz/roles/{role_id}", authMW(authzAuthz(http.HandlerFunc(authzSvcHandler.GetRole))))
		mux.Handle("PATCH /v1/authz/roles/{role_id}", authMW(csrfMW(authzAuthz(http.HandlerFunc(authzSvcHandler.UpdateRole)))))
		mux.Handle("DELETE /v1/authz/roles/{role_id}", authMW(csrfMW(authzAuthz(http.HandlerFunc(authzSvcHandler.DeleteRole)))))
		// User role assignment
		mux.Handle("GET /v1/authz/users/{user_id}/roles", authMW(authzAuthz(http.HandlerFunc(authzSvcHandler.ListUserRoles))))
		mux.Handle("POST /v1/authz/users/{user_id}/roles", authMW(csrfMW(authzAuthz(http.HandlerFunc(authzSvcHandler.AssignRole)))))
		mux.Handle("DELETE /v1/authz/users/{user_id}/roles/{role_id}", authMW(csrfMW(authzAuthz(http.HandlerFunc(authzSvcHandler.RemoveRole)))))
		// Policies
		mux.Handle("GET /v1/authz/policies", authMW(authzAuthz(http.HandlerFunc(authzSvcHandler.ListPolicies))))
		mux.Handle("POST /v1/authz/policies", authMW(csrfMW(authzAuthz(http.HandlerFunc(authzSvcHandler.CreatePolicy)))))
		mux.Handle("PATCH /v1/authz/policies/{policy_id}", authMW(csrfMW(authzAuthz(http.HandlerFunc(authzSvcHandler.UpdatePolicy)))))
		mux.Handle("DELETE /v1/authz/policies/{policy_id}", authMW(csrfMW(authzAuthz(http.HandlerFunc(authzSvcHandler.DeletePolicy)))))
		// Access check — authMW only (no authzAuthz self-check — circular bootstrap problem).
		// /v1/authz/check is the enforcement endpoint itself; wrapping it in authzAuthz
		// creates a self-referential deny loop. JWT auth is sufficient to call it.
		mux.Handle("POST /v1/authz/check", authMW(http.HandlerFunc(authzSvcHandler.CheckAccess)))
		mux.Handle("POST /v1/authz/check/batch", authMW(http.HandlerFunc(authzSvcHandler.BatchCheckAccess)))
		mux.Handle("POST /v1/authz/check:batch", authMW(http.HandlerFunc(authzSvcHandler.BatchCheckAccess)))
		// Audits & permissions
		mux.Handle("GET /v1/authz/audits", authMW(authzAuthz(http.HandlerFunc(authzSvcHandler.ListAudits))))
		mux.Handle("GET /v1/authz/users/{user_id}/permissions", authMW(authzAuthz(http.HandlerFunc(authzSvcHandler.GetUserPermissions))))
		// Portal config
		mux.Handle("GET /v1/authz/portals/configs", authMW(authzAuthz(http.HandlerFunc(authzSvcHandler.ListPortalConfigs))))
		mux.Handle("GET /v1/authz/portals/{portal}/config", authMW(authzAuthz(http.HandlerFunc(authzSvcHandler.GetPortalConfig))))
		mux.Handle("PATCH /v1/authz/portals/{portal}/config", authMW(csrfMW(authzAuthz(http.HandlerFunc(authzSvcHandler.UpdatePortalConfig)))))
	}
	// ── Tenant APIs (tenant-service) ─────────────────────────────────────────
	if tenantConn := getServiceConn(clientManager, "tenant"); tenantConn != nil {
		tenantHandler := handlers.NewPoliSyncHandler(tenantConn, "tenant-service")
		authzTenant := authzMW("svc:tenant", PathSegmentExtractor("/v1/"))
		mux.Handle("GET /v1/tenants", authMW(authzTenant(tenantHandler.Proxy())))
		mux.Handle("POST /v1/tenants", authMW(csrfMW(authzTenant(tenantHandler.Proxy()))))
		mux.Handle("GET /v1/tenants/{tenant_id}", authMW(authzTenant(tenantHandler.Proxy())))
		mux.Handle("PATCH /v1/tenants/{tenant_id}", authMW(csrfMW(authzTenant(tenantHandler.Proxy()))))
		mux.Handle("GET /v1/tenants/{tenant_id}/config", authMW(authzTenant(tenantHandler.Proxy())))
		mux.Handle("PUT /v1/tenants/{tenant_id}/config", authMW(csrfMW(authzTenant(tenantHandler.Proxy()))))
	}

	// ── Audit APIs (audit-service) ────────────────────────────────────────────
	if auditConn := getServiceConn(clientManager, "audit"); auditConn != nil {
		auditHandler := handlers.NewPoliSyncHandler(auditConn, "audit-service")
		authzAudit := authzMW("svc:audit", PathSegmentExtractor("/v1/"))
		mux.Handle("GET /v1/audit-events", authMW(authzAudit(auditHandler.Proxy())))
		mux.Handle("POST /v1/audit-events", authMW(csrfMW(authzAudit(auditHandler.Proxy()))))
		mux.Handle("GET /v1/audit-logs", authMW(authzAudit(auditHandler.Proxy())))
		mux.Handle("POST /v1/audit-logs", authMW(csrfMW(authzAudit(auditHandler.Proxy()))))
		mux.Handle("GET /v1/entities/{entity_type}/{entity_id}/audit-trail", authMW(authzAudit(auditHandler.Proxy())))
		mux.Handle("GET /v1/compliance-logs", authMW(authzAudit(auditHandler.Proxy())))
		mux.Handle("POST /v1/compliance-logs", authMW(csrfMW(authzAudit(auditHandler.Proxy()))))
		mux.Handle("POST /v1/compliance-reports/generate", authMW(csrfMW(authzAudit(auditHandler.Proxy()))))
	}

	// ── KYC Verifications (kyc-service) ───────────────────────────────────────
	if kycConn := getServiceConn(clientManager, "kyc"); kycConn != nil {
		kycHandler := handlers.NewPoliSyncHandler(kycConn, "kyc-service")
		authzKYCSvc := authzMW("svc:kyc", PathSegmentExtractor("/v1/"))
		mux.Handle("POST /v1/kyc-verifications", authMW(csrfMW(authzKYCSvc(kycHandler.Proxy()))))
		mux.Handle("GET /v1/kyc-verifications/{kyc_verification_id}", authMW(authzKYCSvc(kycHandler.Proxy())))
		mux.Handle("POST /v1/kyc-verifications/{kyc_verification_id}/verify", authMW(csrfMW(authzKYCSvc(kycHandler.Proxy()))))
		mux.Handle("POST /v1/kyc-verifications/{kyc_verification_id}/reject", authMW(csrfMW(authzKYCSvc(kycHandler.Proxy()))))
		mux.Handle("GET /v1/kyc-verifications/pending", authMW(authzKYCSvc(kycHandler.Proxy())))
		mux.Handle("POST /v1/kyc-verifications/{kyc_verification_id}/documents", authMW(csrfMW(authzKYCSvc(kycHandler.Proxy()))))
	}
	// ── Beneficiary APIs (beneficiary-service) ───────────────────────────────
	if beneficiaryConn := getServiceConn(clientManager, "beneficiary"); beneficiaryConn != nil {
		beneficiaryHandler := handlers.NewPoliSyncHandler(beneficiaryConn, "beneficiary-service")
		authzBeneficiary := authzMW("svc:beneficiary", PathSegmentExtractor("/v1/"))
		mux.Handle("GET /v1/beneficiaries", authMW(authzBeneficiary(beneficiaryHandler.Proxy())))
		mux.Handle("POST /v1/beneficiaries/individual", authMW(csrfMW(authzBeneficiary(beneficiaryHandler.Proxy()))))
		mux.Handle("POST /v1/beneficiaries/business", authMW(csrfMW(authzBeneficiary(beneficiaryHandler.Proxy()))))
		mux.Handle("GET /v1/beneficiaries/{beneficiary_id}", authMW(authzBeneficiary(beneficiaryHandler.Proxy())))
		mux.Handle("PATCH /v1/beneficiaries/{beneficiary_id}", authMW(csrfMW(authzBeneficiary(beneficiaryHandler.Proxy()))))
		mux.Handle("POST /v1/beneficiaries/{beneficiary_id}/kyc", authMW(csrfMW(authzBeneficiary(beneficiaryHandler.Proxy()))))
		mux.Handle("GET /v1/beneficiaries/{beneficiary_id}/quotes", authMW(authzBeneficiary(beneficiaryHandler.Proxy())))
		mux.Handle("POST /v1/beneficiaries/{beneficiary_id}/risk-score", authMW(csrfMW(authzBeneficiary(beneficiaryHandler.Proxy()))))
	}

	// ── Insurance / Insurer APIs (insurance-service) ──────────────────────────
	if insuranceConn := getServiceConn(clientManager, "insurance"); insuranceConn != nil {
		insuranceHandler := handlers.NewPoliSyncHandler(insuranceConn, "insurance-service")
		authzInsurance := authzMW("svc:insurance", PathSegmentExtractor("/v1/"))
		mux.Handle("GET /v1/insurers", authMW(authzInsurance(insuranceHandler.Proxy())))
		mux.Handle("POST /v1/insurers", authMW(csrfMW(authzInsurance(insuranceHandler.Proxy()))))
		mux.Handle("GET /v1/insurers/{insurer_id}", authMW(authzInsurance(insuranceHandler.Proxy())))
		mux.Handle("PATCH /v1/insurers/{insurer_id}", authMW(csrfMW(authzInsurance(insuranceHandler.Proxy()))))
		mux.Handle("PUT /v1/insurers/{insurer_id}/config", authMW(csrfMW(authzInsurance(insuranceHandler.Proxy()))))
		mux.Handle("GET /v1/insurers/{insurer_id}/products", authMW(authzInsurance(insuranceHandler.Proxy())))
		mux.Handle("POST /v1/insurers/{insurer_id}/products", authMW(csrfMW(authzInsurance(insuranceHandler.Proxy()))))
		mux.Handle("GET /v1/insurers/{insurer_id}/revenue-share", authMW(authzInsurance(insuranceHandler.Proxy())))
		mux.Handle("GET /v1/insurer-products/{insurer_product_id}", authMW(authzInsurance(insuranceHandler.Proxy())))
		mux.Handle("PATCH /v1/insurer-products/{insurer_product_id}", authMW(csrfMW(authzInsurance(insuranceHandler.Proxy()))))
	}

	// ── Notification APIs (notification-service) — gRPC handler (BUG-006 FIX) ─
	// Previously used PoliSyncHandler (HTTP proxy) which fails because notification is gRPC-only.
	if notifConn := getServiceConn(clientManager, "notification"); notifConn != nil {
		notifHandler := handlers.NewNotificationHandler(notifConn)
		authzNotif := authzMW("svc:notification", PathSegmentExtractor("/v1/"))
		mux.Handle("POST /v1/notifications", authMW(csrfMW(authzNotif(http.HandlerFunc(notifHandler.Send)))))
		mux.Handle("GET /v1/notifications/{notification_id}", authMW(authzNotif(http.HandlerFunc(notifHandler.GetStatus))))
		mux.Handle("POST /v1/notifications/mark-as-read", authMW(csrfMW(authzNotif(http.HandlerFunc(notifHandler.MarkAsRead)))))
		mux.Handle("POST /v1/notifications/send-bulk", authMW(csrfMW(authzNotif(http.HandlerFunc(notifHandler.SendBulk)))))
		mux.Handle("GET /v1/users/{user_id}/notifications", authMW(authzNotif(http.HandlerFunc(notifHandler.GetUserNotifications))))
		mux.Handle("POST /v1/notification-templates", authMW(csrfMW(authzNotif(http.HandlerFunc(notifHandler.CreateTemplate)))))
		mux.Handle("PATCH /v1/notification-templates/{template_id}", authMW(csrfMW(authzNotif(http.HandlerFunc(notifHandler.UpdateTemplate)))))
		mux.Handle("POST /v1/notification-templates/{template_id}/deactivate", authMW(csrfMW(authzNotif(http.HandlerFunc(notifHandler.DeactivateTemplate)))))
	}

	// ── API Keys (standalone v1/api-keys — separate from /v1/auth/api-keys) ──
	if apiKeyConn := getServiceConn(clientManager, "authn"); apiKeyConn != nil {
		authzAPIKeySvc := authzMW("svc:apikey", PathSegmentExtractor("/v1/"))
		// Re-use authn handler for now; the authn service owns API key management.
		mux.Handle("GET /v1/api-keys", authMW(authzAPIKeySvc(http.HandlerFunc(notImplementedHandler))))
		mux.Handle("POST /v1/api-keys", authMW(csrfMW(authzAPIKeySvc(http.HandlerFunc(notImplementedHandler)))))
		mux.Handle("GET /v1/api-keys/{api_key_id}", authMW(authzAPIKeySvc(http.HandlerFunc(notImplementedHandler))))
		mux.Handle("DELETE /v1/api-keys/{api_key_id}", authMW(csrfMW(authzAPIKeySvc(http.HandlerFunc(notImplementedHandler)))))
		mux.Handle("POST /v1/api-keys/{api_key_id}/rotate", authMW(csrfMW(authzAPIKeySvc(http.HandlerFunc(notImplementedHandler)))))
		mux.Handle("POST /v1/api-keys/{api_key_id}/usage", authMW(authzAPIKeySvc(http.HandlerFunc(notImplementedHandler))))
		mux.Handle("GET /v1/api-keys/{api_key_id}/usage", authMW(authzAPIKeySvc(http.HandlerFunc(notImplementedHandler))))
		mux.Handle("POST /v1/api-keys/validate", authMW(authzAPIKeySvc(http.HandlerFunc(notImplementedHandler))))
		mux.Handle("POST /v1/auth/api-keys/{key_id}/rotate", authMW(csrfMW(authzAPIKeySvc(http.HandlerFunc(notImplementedHandler)))))
	}
	// ── Workflow APIs (workflow-service) ─────────────────────────────────────
	if workflowConn := getServiceConn(clientManager, "workflow"); workflowConn != nil {
		workflowHandler := handlers.NewPoliSyncHandler(workflowConn, "workflow-service")
		authzWorkflow := authzMW("svc:workflow", PathSegmentExtractor("/v1/"))
		mux.Handle("POST /v1/workflow-definitions", authMW(csrfMW(authzWorkflow(workflowHandler.Proxy()))))
		mux.Handle("GET /v1/workflow-definitions/{workflow_definition_id}", authMW(authzWorkflow(workflowHandler.Proxy())))
		mux.Handle("POST /v1/workflow-instances", authMW(csrfMW(authzWorkflow(workflowHandler.Proxy()))))
		mux.Handle("GET /v1/workflow-instances/{workflow_instance_id}", authMW(authzWorkflow(workflowHandler.Proxy())))
		mux.Handle("POST /v1/workflow-tasks/{task_id}/complete", authMW(csrfMW(authzWorkflow(workflowHandler.Proxy()))))
		mux.Handle("GET /v1/workflow-tasks/my-tasks", authMW(authzWorkflow(workflowHandler.Proxy())))
		mux.Handle("GET /v1/entities/{entity_type}/{entity_id}/workflow-history", authMW(authzWorkflow(workflowHandler.Proxy())))
		// Tasks
		mux.Handle("POST /v1/tasks", authMW(csrfMW(authzWorkflow(workflowHandler.Proxy()))))
		mux.Handle("GET /v1/tasks/{task_id}", authMW(authzWorkflow(workflowHandler.Proxy())))
		mux.Handle("PATCH /v1/tasks/{task_id}", authMW(csrfMW(authzWorkflow(workflowHandler.Proxy()))))
		mux.Handle("POST /v1/tasks/{task_id}/assign", authMW(csrfMW(authzWorkflow(workflowHandler.Proxy()))))
		mux.Handle("POST /v1/tasks/{task_id}/complete", authMW(csrfMW(authzWorkflow(workflowHandler.Proxy()))))
		mux.Handle("GET /v1/tasks/my-tasks", authMW(authzWorkflow(workflowHandler.Proxy())))
	}

	// ── Analytics APIs (analytics-service) ───────────────────────────────────
	if analyticsConn := getServiceConn(clientManager, "analytics"); analyticsConn != nil {
		analyticsHandler := handlers.NewPoliSyncHandler(analyticsConn, "analytics-service")
		authzAnalytics := authzMW("svc:analytics", PathSegmentExtractor("/v1/"))
		mux.Handle("POST /v1/analytics/dashboards", authMW(csrfMW(authzAnalytics(analyticsHandler.Proxy()))))
		mux.Handle("GET /v1/analytics/dashboards/{dashboard_id}", authMW(authzAnalytics(analyticsHandler.Proxy())))
		mux.Handle("POST /v1/analytics/metrics", authMW(authzAnalytics(analyticsHandler.Proxy())))
		mux.Handle("POST /v1/analytics/queries/run", authMW(authzAnalytics(analyticsHandler.Proxy())))
		mux.Handle("POST /v1/analytics/reports/{report_id}/generate", authMW(csrfMW(authzAnalytics(analyticsHandler.Proxy()))))
		mux.Handle("POST /v1/analytics/reports/schedule", authMW(csrfMW(authzAnalytics(analyticsHandler.Proxy()))))
		// Reporting
		mux.Handle("GET /v1/report-definitions", authMW(authzAnalytics(analyticsHandler.Proxy())))
		mux.Handle("GET /v1/report-executions", authMW(authzAnalytics(analyticsHandler.Proxy())))
		mux.Handle("GET /v1/report-executions/{report_execution_id}", authMW(authzAnalytics(analyticsHandler.Proxy())))
		mux.Handle("GET /v1/report-executions/{report_execution_id}/download", authMW(authzAnalytics(analyticsHandler.Proxy())))
		mux.Handle("GET /v1/report-schedules", authMW(authzAnalytics(analyticsHandler.Proxy())))
		mux.Handle("POST /v1/report-schedules", authMW(csrfMW(authzAnalytics(analyticsHandler.Proxy()))))
		mux.Handle("POST /v1/reports/{report_definition_id}/execute", authMW(csrfMW(authzAnalytics(analyticsHandler.Proxy()))))
	}

	// ── AI APIs (ai-service) ──────────────────────────────────────────────────
	if aiConn := getServiceConn(clientManager, "ai"); aiConn != nil {
		aiHandler := handlers.NewPoliSyncHandler(aiConn, "ai-service")
		authzAI := authzMW("svc:ai", PathSegmentExtractor("/v1/"))
		mux.Handle("POST /v1/ai/chat", authMW(authzAI(aiHandler.Proxy())))
		mux.Handle("POST /v1/ai/claims/evaluate", authMW(csrfMW(authzAI(aiHandler.Proxy()))))
		mux.Handle("POST /v1/ai/documents/analyze", authMW(csrfMW(authzAI(aiHandler.Proxy()))))
		mux.Handle("POST /v1/ai/fraud/detect", authMW(csrfMW(authzAI(aiHandler.Proxy()))))
		mux.Handle("POST /v1/ai/risk/assess", authMW(csrfMW(authzAI(aiHandler.Proxy()))))
	}

	// ── Support / Tickets (support-service) ───────────────────────────────────
	if supportConn := getServiceConn(clientManager, "support"); supportConn != nil {
		supportHandler := handlers.NewPoliSyncHandler(supportConn, "support-service")
		authzSupport := authzMW("svc:support", PathSegmentExtractor("/v1/"))
		mux.Handle("GET /v1/tickets", authMW(authzSupport(supportHandler.Proxy())))
		mux.Handle("POST /v1/tickets", authMW(csrfMW(authzSupport(supportHandler.Proxy()))))
		mux.Handle("GET /v1/tickets/{ticket_id}", authMW(authzSupport(supportHandler.Proxy())))
		mux.Handle("POST /v1/tickets/{ticket_id}/assign", authMW(csrfMW(authzSupport(supportHandler.Proxy()))))
		mux.Handle("PATCH /v1/tickets/{ticket_id}/status", authMW(csrfMW(authzSupport(supportHandler.Proxy()))))
		mux.Handle("POST /v1/tickets/{ticket_id}/messages", authMW(csrfMW(authzSupport(supportHandler.Proxy()))))
		// Knowledge base
		mux.Handle("POST /v1/knowledge-base", authMW(csrfMW(authzSupport(supportHandler.Proxy()))))
		mux.Handle("GET /v1/knowledge-base/{slug}", authMW(authzSupport(supportHandler.Proxy())))
		mux.Handle("PATCH /v1/knowledge-base/{article_id}", authMW(csrfMW(authzSupport(supportHandler.Proxy()))))
		mux.Handle("DELETE /v1/knowledge-base/{article_id}", authMW(csrfMW(authzSupport(supportHandler.Proxy()))))
		mux.Handle("GET /v1/knowledge-base/search", authMW(authzSupport(supportHandler.Proxy())))
		// FAQs
		mux.Handle("GET /v1/faqs", authMW(authzSupport(supportHandler.Proxy())))
		mux.Handle("POST /v1/faqs", authMW(csrfMW(authzSupport(supportHandler.Proxy()))))
		mux.Handle("PATCH /v1/faqs/{faq_id}", authMW(csrfMW(authzSupport(supportHandler.Proxy()))))
		mux.Handle("DELETE /v1/faqs/{faq_id}", authMW(csrfMW(authzSupport(supportHandler.Proxy()))))
	}

	// ── IoT APIs (iot-service) ────────────────────────────────────────────────
	if iotConn := getServiceConn(clientManager, "iot"); iotConn != nil {
		iotHandler := handlers.NewPoliSyncHandler(iotConn, "iot-service")
		authzIoT := authzMW("svc:iot", PathSegmentExtractor("/v1/"))
		mux.Handle("POST /v1/iot/devices", authMW(csrfMW(authzIoT(iotHandler.Proxy()))))
		mux.Handle("GET /v1/iot/devices/{device_id}", authMW(authzIoT(iotHandler.Proxy())))
		mux.Handle("POST /v1/iot/devices/{device_id}/deactivate", authMW(csrfMW(authzIoT(iotHandler.Proxy()))))
		mux.Handle("GET /v1/iot/devices/{device_id}/risk", authMW(authzIoT(iotHandler.Proxy())))
		mux.Handle("POST /v1/iot/telemetry", authMW(csrfMW(authzIoT(iotHandler.Proxy()))))
	}

	// ── Voice Sessions (standalone, non-authn) ────────────────────────────────
	if webrtcConn := getServiceConn(clientManager, "webrtc"); webrtcConn != nil {
		voiceHandler := handlers.NewPoliSyncHandler(webrtcConn, "webrtc-service")
		authzVoiceSvc := authzMW("svc:voice", PathSegmentExtractor("/v1/"))
		mux.Handle("POST /v1/voice-sessions", authMW(authzVoiceSvc(voiceHandler.Proxy())))
		mux.Handle("GET /v1/voice-sessions/{voice_session_id}", authMW(authzVoiceSvc(voiceHandler.Proxy())))
		mux.Handle("POST /v1/voice-sessions/{voice_session_id}/end", authMW(csrfMW(authzVoiceSvc(voiceHandler.Proxy()))))
		mux.Handle("POST /v1/voice-sessions/{voice_session_id}/commands", authMW(csrfMW(authzVoiceSvc(voiceHandler.Proxy()))))
		mux.Handle("GET /v1/voice-sessions/{voice_session_id}/transcript", authMW(authzVoiceSvc(voiceHandler.Proxy())))
		// Voice biometric (under auth namespace)
		mux.Handle("POST /v1/auth/voice-biometric/initiate", authMW(authzVoiceSvc(voiceHandler.Proxy())))
		mux.Handle("POST /v1/auth/voice-biometric/submit", authMW(csrfMW(authzVoiceSvc(voiceHandler.Proxy()))))
		mux.Handle("POST /v1/auth/voice-biometric/verify", authMW(authzVoiceSvc(voiceHandler.Proxy())))
	}

	// ── Quotes — canonical paths (quote-service) ──────────────────────────────
	// These are the ENDPOINT_MAP canonical /v1/quotes/* routes; proxy to quote-service.
	if quoteCanonConn := getServiceConn(clientManager, "quote"); quoteCanonConn != nil {
		quoteCanonHandler := handlers.NewPoliSyncHandler(quoteCanonConn, "quote-service")
		authzQuoteCanon := authzMW("svc:quote", PathSegmentExtractor("/v1/"))
		mux.Handle("POST /v1/quotes", authMW(csrfMW(authzQuoteCanon(quoteCanonHandler.Proxy()))))
		mux.Handle("GET /v1/quotes/{quote_id}", authMW(authzQuoteCanon(quoteCanonHandler.Proxy())))
		mux.Handle("POST /v1/quotes/{quote_id}/approve", authMW(csrfMW(authzQuoteCanon(quoteCanonHandler.Proxy()))))
		mux.Handle("POST /v1/quotes/{quote_id}/reject", authMW(csrfMW(authzQuoteCanon(quoteCanonHandler.Proxy()))))
		mux.Handle("POST /v1/quotes/{quote_id}/convert", authMW(csrfMW(authzQuoteCanon(quoteCanonHandler.Proxy()))))
		mux.Handle("GET /v1/quotes/{quote_id}/decision", authMW(authzQuoteCanon(quoteCanonHandler.Proxy())))
		mux.Handle("GET /v1/quotes/{quote_id}/health-declaration", authMW(authzQuoteCanon(quoteCanonHandler.Proxy())))
		mux.Handle("POST /v1/quotes/{quote_id}/health-declaration", authMW(csrfMW(authzQuoteCanon(quoteCanonHandler.Proxy()))))
	}

	// ── Policy — missing canonical actions (policy-service) ───────────────────
	if policyCanonConn := getServiceConn(clientManager, "policy"); policyCanonConn != nil {
		policyCanonHandler := handlers.NewPoliSyncHandler(policyCanonConn, "policy-service")
		authzPolicyCanon := authzMW("svc:policy", PathSegmentExtractor("/v1/"))
		policyCanonProxy := policyCanonHandler.Proxy()
		mux.Handle("GET /v1/insurance-proposals", authMW(authzPolicyCanon(policyCanonProxy)))
		mux.Handle("POST /v1/insurance-proposals", authMW(csrfMW(authzPolicyCanon(policyCanonProxy))))
		mux.Handle("GET /v1/insurance-proposals/{proposal_id}", authMW(authzPolicyCanon(policyCanonProxy)))
		mux.Handle("PATCH /v1/insurance-proposals/{proposal_id}", authMW(csrfMW(authzPolicyCanon(policyCanonProxy))))
		mux.Handle("DELETE /v1/insurance-proposals/{proposal_id}", authMW(csrfMW(authzPolicyCanon(policyCanonProxy))))
		mux.Handle("POST /v1/insurance-proposals/{proposal_id}/approve", authMW(csrfMW(authzPolicyCanon(policyCanonProxy))))
		mux.Handle("POST /v1/insurance-proposals/{proposal_id}/reject", authMW(csrfMW(authzPolicyCanon(policyCanonProxy))))
		mux.Handle("GET /v1/orders/{order_id}/proposal", authMW(authzPolicyCanon(policyCanonProxy)))
		mux.Handle("POST /v1/orders/{order_id}/proposal", authMW(csrfMW(authzPolicyCanon(policyCanonProxy))))
		mux.Handle("PATCH /v1/policies/{policy_id}", authMW(csrfMW(authzPolicyCanon(policyCanonProxy))))
		mux.Handle("POST /v1/policies/{policy_id}/issue", authMW(csrfMW(authzPolicyCanon(policyCanonProxy))))
		mux.Handle("POST /v1/policies/{policy_id}/revive", authMW(csrfMW(authzPolicyCanon(policyCanonProxy))))
		mux.Handle("POST /v1/policies/{policy_id}/renew", authMW(csrfMW(authzPolicyCanon(policyCanonProxy))))
		mux.Handle("POST /v1/policies/{policy_id}/renew-tenure", authMW(csrfMW(authzPolicyCanon(policyCanonProxy))))
		mux.Handle("POST /v1/policies/{policy_id}/generate-document", authMW(csrfMW(authzPolicyCanon(policyCanonProxy))))
		mux.Handle("GET /v1/policies/{policy_id}/grace-period", authMW(authzPolicyCanon(policyCanonProxy)))
		mux.Handle("GET /v1/policies/{policy_id}/renewal-schedule", authMW(authzPolicyCanon(policyCanonProxy)))
		mux.Handle("POST /v1/policies/{policy_id}/refund", authMW(csrfMW(authzPolicyCanon(policyCanonProxy))))
		mux.Handle("POST /v1/policies/{policy_id}/refunds/calculate", authMW(csrfMW(authzPolicyCanon(policyCanonProxy))))
		mux.Handle("GET /v1/users/{customer_id}/policies", authMW(authzPolicyCanon(policyCanonProxy)))
		mux.Handle("GET /v1/users/{customer_id}/claims", authMW(authzPolicyCanon(policyCanonProxy)))
	}

	// ── Refunds — canonical (payment/policy-service) ──────────────────────────
	if refundPayConn := getServiceConn(clientManager, "payment"); refundPayConn != nil {
		refundHandler := handlers.NewPoliSyncHandler(refundPayConn, "payment-service")
		authzRefund := authzMW("svc:payment", PathSegmentExtractor("/v1/"))
		mux.Handle("GET /v1/refunds", authMW(authzRefund(refundHandler.Proxy())))
		mux.Handle("GET /v1/refunds/{refund_id}/status", authMW(authzRefund(refundHandler.Proxy())))
		mux.Handle("POST /v1/refunds/{refund_id}/approve", authMW(csrfMW(authzRefund(refundHandler.Proxy()))))
		mux.Handle("POST /v1/refunds/{refund_id}/process", authMW(csrfMW(authzRefund(refundHandler.Proxy()))))
	}

	// ── Renewal APIs (renewal via policy-service) ─────────────────────────────
	if renewalConn := getServiceConn(clientManager, "policy"); renewalConn != nil {
		renewalHandler := handlers.NewPoliSyncHandler(renewalConn, "policy-service")
		authzRenewal := authzMW("svc:policy", PathSegmentExtractor("/v1/"))
		mux.Handle("GET /v1/renewals/upcoming", authMW(authzRenewal(renewalHandler.Proxy())))
		mux.Handle("POST /v1/renewal-schedules/{renewal_schedule_id}/reminders", authMW(csrfMW(authzRenewal(renewalHandler.Proxy()))))
	}

	// ── MFS (Mobile Financial Services) ──────────────────────────────────────
	if mfsConn := getServiceConn(clientManager, "payment"); mfsConn != nil {
		mfsHandler := handlers.NewPoliSyncHandler(mfsConn, "payment-service")
		authzMFS := authzMW("svc:payment", PathSegmentExtractor("/v1/"))
		mux.Handle("POST /v1/mfs/payments", authMW(csrfMW(authzMFS(mfsHandler.Proxy()))))
		mux.Handle("POST /v1/mfs/refunds", authMW(csrfMW(authzMFS(mfsHandler.Proxy()))))
		mux.Handle("GET /v1/mfs/transactions", authMW(authzMFS(mfsHandler.Proxy())))
		mux.Handle("GET /v1/mfs/transactions/{mfs_transaction_id}", authMW(authzMFS(mfsHandler.Proxy())))
		mux.Handle("POST /v1/mfs/webhooks/{provider}", authMW(csrfMW(authzMFS(mfsHandler.Proxy()))))
	}

	// ── Orders — missing canonical actions (order-service) ────────────────────
	if orderCanonConn := getServiceConn(clientManager, "order"); orderCanonConn != nil {
		orderCanonHandler := handlers.NewPoliSyncHandler(orderCanonConn, "order-service")
		authzOrderCanon := authzMW("svc:order", PathSegmentExtractor("/v1/"))
		mux.Handle("POST /v1/orders/{order_id}/pay", authMW(csrfMW(authzOrderCanon(orderCanonHandler.Proxy()))))
		mux.Handle("POST /v1/orders/{order_id}/confirm-payment", authMW(csrfMW(authzOrderCanon(orderCanonHandler.Proxy()))))
		mux.Handle("GET /v1/orders/{order_id}/status", authMW(authzOrderCanon(orderCanonHandler.Proxy())))
	}

	// ── Payment webhook for specific providers (literal paths to avoid wildcard conflicts) ─
	// The SSLCommerz-specific routes are already registered above (sslcommerz, bkash).
	// Register explicit provider webhook paths for each known provider.
	if payWebhookConn := getServiceConn(clientManager, "payment"); payWebhookConn != nil {
		pwHandler := handlers.NewPoliSyncHandler(payWebhookConn, "payment-service")
		mux.Handle("POST /v1/payments/webhooks/bkash", http.HandlerFunc(pwHandler.Proxy().ServeHTTP))
		mux.Handle("POST /v1/payments/webhooks/nagad", http.HandlerFunc(pwHandler.Proxy().ServeHTTP))
		mux.Handle("POST /v1/payments/webhooks/rocket", http.HandlerFunc(pwHandler.Proxy().ServeHTTP))
		mux.Handle("POST /v1/payments/webhooks/upay", http.HandlerFunc(pwHandler.Proxy().ServeHTTP))
	}

	// ── DocGen — Document generation & template management ───────────────────
	// Google API style:
	//   POST /v1/documents:generate              → generate (JSON response with file_url)
	//   POST /v1/documents:generate?$alt=media   → generate + stream raw bytes immediately
	//   GET  /v1/documents/{document_id}         → get document metadata
	//   GET  /v1/documents/{document_id}/download → download raw file bytes
	//   GET  /v1/entities/{type}/{id}/documents  → list documents for an entity
	//   DELETE /v1/documents/{document_id}       → delete document
	//   POST /v1/document-templates              → create template
	//   GET  /v1/document-templates              → list templates
	//   GET  /v1/document-templates/{template_id} → get template
	//   PATCH /v1/document-templates/{template_id} → update template
	//   POST /v1/document-templates/{template_id}:deactivate → deactivate
	//   DELETE /v1/document-templates/{template_id} → delete template
	//
	// Format selection (Google API style, in order of precedence):
	//   1. ?format=xlsx   or  ?output_format=xlsx  (query param)
	//   2. Body field:  {"output_format": "xlsx"}
	//   Supported values: pdf, html, docx, xlsx
	if docgenConn := getServiceConn(clientManager, "docgen"); docgenConn != nil {
		dg := handlers.NewDocGenHandler(docgenConn)
		authzDoc := authzMW("svc:document", PathSegmentExtractor("/v1/"))

		// ── Document operations ────────────────────────────────────────────
		// POST /v1/documents  (legacy alias)
		mux.Handle("POST /v1/documents", authMW(csrfMW(authzDoc(http.HandlerFunc(dg.Generate)))))
		// POST /v1/documents:generate  (custom verb — rewritten by customVerbCompatMiddleware)
		mux.Handle("POST /v1/documents/generate", authMW(csrfMW(authzDoc(http.HandlerFunc(dg.Generate)))))
		// GET  /v1/documents/{document_id}
		mux.Handle("GET /v1/documents/{document_id}", authMW(authzDoc(http.HandlerFunc(dg.GetDocument))))
		// GET  /v1/documents/{document_id}/download  — streams raw bytes
		mux.Handle("GET /v1/documents/{document_id}/download", authMW(authzDoc(http.HandlerFunc(dg.DownloadRaw))))
		// GET  /v1/entities/{entity_type}/{entity_id}/documents
		mux.Handle("GET /v1/entities/{entity_type}/{entity_id}/documents", authMW(authzDoc(http.HandlerFunc(dg.ListDocuments))))
		// DELETE /v1/documents/{document_id}
		mux.Handle("DELETE /v1/documents/{document_id}", authMW(csrfMW(authzDoc(http.HandlerFunc(dg.DeleteDocument)))))

		// ── Template CRUD ──────────────────────────────────────────────────
		// POST /v1/document-templates
		mux.Handle("POST /v1/document-templates", authMW(csrfMW(authzDoc(http.HandlerFunc(dg.CreateTemplate)))))
		// GET  /v1/document-templates
		mux.Handle("GET /v1/document-templates", authMW(authzDoc(http.HandlerFunc(dg.ListTemplates))))
		// GET  /v1/document-templates/{template_id}
		mux.Handle("GET /v1/document-templates/{template_id}", authMW(authzDoc(http.HandlerFunc(dg.GetTemplate))))
		// PATCH /v1/document-templates/{template_id}
		mux.Handle("PATCH /v1/document-templates/{template_id}", authMW(csrfMW(authzDoc(http.HandlerFunc(dg.UpdateTemplate)))))
		// POST /v1/document-templates/{template_id}/deactivate  (rewritten from :deactivate)
		mux.Handle("POST /v1/document-templates/{template_id}/deactivate", authMW(csrfMW(authzDoc(http.HandlerFunc(dg.DeactivateTemplate)))))
		// DELETE /v1/document-templates/{template_id}
		mux.Handle("DELETE /v1/document-templates/{template_id}", authMW(csrfMW(authzDoc(http.HandlerFunc(dg.DeleteTemplate)))))
	}
	// Phase E route consolidation:
	// authn topology routes remain hidden under /v1/auth/* while public service
	// APIs are exposed directly under their own namespaces (for example /v1/media/*).

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.NotFound(w, r) })

	var h http.Handler = mux
	h = customVerbCompatMiddleware(h)
	h = middleware.Recovery(h)
	h = middleware.RequestID(h)
	h = middleware.SecurityHeaders(h)
	h = corsMiddleware(h)
	h = middleware.Metrics(h)
	h = middleware.MaxBodySize(10 * 1024 * 1024)(h)
	h = middleware.Compression(middleware.CompressionDefault)(h)
	h = middleware.Timeout(30 * time.Second)(h)
	return h
}

func notImplementedHandler(w http.ResponseWriter, r *http.Request) {
	respond.Error(w, r, http.StatusNotImplemented, "NOT_IMPLEMENTED", "not implemented")
}

// getServiceConn retrieves a gRPC connection for a named service from the client manager.
// Returns nil if the service is not registered or the connection is unavailable.
func getServiceConn(cm *resilience.ResilientClientManager, name string) *grpc.ClientConn {
	if cm == nil {
		return nil
	}
	client, err := cm.GetClient(name)
	if err != nil {
		return nil
	}
	conn, err := client.GetConnection()
	if err != nil {
		return nil
	}
	return conn
}

func corsMiddleware(next http.Handler) http.Handler {
	allowedOriginsEnv := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOriginsEnv == "" {
		allowedOriginsEnv = "http://localhost:3000,http://localhost:5173,http://b2b.labaidinsuretech.com,https://b2b.labaidinsuretech.com,http://system.labaidinsuretech.com,https://system.labaidinsuretech.com,http://146.190.97.242,http://146.190.97.242:3000"
	}
	allowedOrigins := strings.Split(allowedOriginsEnv, ",")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowed := false
		for _, o := range allowedOrigins {
			if strings.TrimSpace(o) == origin {
				allowed = true
				break
			}
		}
		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Vary", "Origin")
		}

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-CSRF-Token, X-Device-Id")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// customVerbCompatMiddleware keeps compatibility with legacy Google-style custom verb
// routes that use "{id}:action" segments, which Go's ServeMux pattern parser rejects
// when combined with path wildcards.
func customVerbCompatMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// normalize rewrites "{prefix}{id}{suffix}" → "{prefix}{id}{replacement}".
		// Used for routes with a variable ID segment before the custom verb colon.
		normalize := func(prefix, suffix, replacement string) bool {
			if !strings.HasPrefix(r.URL.Path, prefix) || !strings.HasSuffix(r.URL.Path, suffix) {
				return false
			}
			id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, prefix), suffix)
			if id == "" || strings.Contains(id, "/") {
				return false
			}
			r.URL.Path = prefix + id + replacement
			r.URL.RawPath = r.URL.Path
			return true
		}

		// normalizeFixed rewrites an exact fixed path (no variable ID segment).
		// Used for custom-verb routes where the entire resource is in the prefix.
		normalizeFixed := func(from, to string) bool {
			if r.URL.Path != from {
				return false
			}
			r.URL.Path = to
			r.URL.RawPath = to
			return true
		}

		// Document generation custom verb
		_ = normalizeFixed("/v1/documents:generate", "/v1/documents/generate") ||
			normalizeFixed("/v1/auth/voice-biometric:initiate", "/v1/auth/voice-biometric/initiate") ||
			normalizeFixed("/v1/auth/voice-biometric:submit", "/v1/auth/voice-biometric/submit") ||
			normalizeFixed("/v1/auth/voice-biometric:verify", "/v1/auth/voice-biometric/verify") ||
			normalize("/v1/auth/documents/", ":verify", "/verify") ||
			normalize("/v1/auth/kyc/", ":approve", "/approve") ||
			normalize("/v1/auth/kyc/", ":reject", "/reject") ||
			normalize("/v1/partners/", ":verify", "/verify") ||
			normalize("/v1/partners/", ":updateStatus", "/updateStatus") ||
			normalize("/v1/partners/", ":update-status", "/update-status") ||
			normalize("/v1/fraud-rules/", ":activate", "/activate") ||
			normalize("/v1/fraud-rules/", ":deactivate", "/deactivate") ||
			normalize("/v1/payments/", ":submit-proof", "/submit-proof") ||
			normalize("/v1/payments/", ":review", "/review") ||
			normalize("/v1/payments/", ":generate-receipt", "/generate-receipt") ||
			normalize("/v1/payments/", ":verify", "/verify") ||
			normalize("/v1/payments/", ":reconcile", "/reconcile") ||
			normalize("/v1/invoices/", ":mark-paid", "/mark-paid") ||
			normalize("/v1/invoices/", ":cancel", "/cancel") ||
			normalize("/v1/invoices/", ":issue", "/issue") ||
			normalize("/v1/invoices/", ":generate-pdf", "/generate-pdf") ||
			normalize("/v1/quotes/", ":approve", "/approve") ||
			normalize("/v1/quotes/", ":reject", "/reject") ||
			normalize("/v1/quotes/", ":convert", "/convert") ||
			normalize("/v1/insurance-proposals/", ":approve", "/approve") ||
			normalize("/v1/insurance-proposals/", ":reject", "/reject") ||
			normalize("/v1/policies/", ":cancel", "/cancel") ||
			normalize("/v1/policies/", ":issue", "/issue") ||
			normalize("/v1/policies/", ":revive", "/revive") ||
			normalize("/v1/policies/", ":renew", "/renew") ||
			normalize("/v1/policies/", ":renew-tenure", "/renew-tenure") ||
			normalize("/v1/policies/", ":generate-document", "/generate-document") ||
			normalize("/v1/policies/", ":suspend", "/suspend") ||
			normalize("/v1/policies/", ":reinstate", "/reinstate") ||
			normalize("/v1/claims/", ":approve", "/approve") ||
			normalize("/v1/claims/", ":reject", "/reject") ||
			normalize("/v1/claims/", ":settle", "/settle") ||
			normalize("/v1/claims/", ":dispute", "/dispute") ||
			normalize("/v1/claims/", ":request-documents", "/request-documents") ||
			normalize("/v1/endorsements/", ":approve", "/approve") ||
			normalize("/v1/endorsements/", ":reject", "/reject") ||
			normalize("/v1/orders/", ":cancel", "/cancel") ||
			normalize("/v1/orders/", ":pay", "/pay") ||
			normalize("/v1/orders/", ":confirm-payment", "/confirm-payment") ||
			normalize("/v1/refunds/", ":approve", "/approve") ||
			normalize("/v1/refunds/", ":process", "/process") ||
			normalize("/v1/commission-payouts/", ":process", "/process") ||
			normalize("/v1/products/", ":activate", "/activate") ||
			normalize("/v1/products/", ":deactivate", "/deactivate") ||
			normalize("/v1/products/", ":discontinue", "/discontinue") ||
			normalize("/v1/products/", ":calculate-premium", "/calculate-premium") ||
			normalize("/v1/kyc-verifications/", ":verify", "/verify") ||
			normalize("/v1/kyc-verifications/", ":reject", "/reject") ||
			normalize("/v1/document-templates/", ":deactivate", "/deactivate") ||
			normalize("/v1/notification-templates/", ":deactivate", "/deactivate") ||
			normalize("/v1/workflow-tasks/", ":complete", "/complete") ||
			normalize("/v1/tasks/", ":assign", "/assign") ||
			normalize("/v1/tasks/", ":complete", "/complete") ||
			normalize("/v1/tickets/", ":assign", "/assign") ||
			normalize("/v1/iot/devices/", ":deactivate", "/deactivate") ||
			normalize("/v1/voice-sessions/", ":end", "/end") ||
			normalize("/v1/auth/api-keys/", ":revoke", "/revoke") ||
			normalize("/v1/auth/api-keys/", ":rotate", "/rotate") ||
			normalize("/v1/auth/users/", ":sessions:revoke-all", "/sessions/revoke-all") ||
			normalize("/v1/auth/users/", ":totp:enable", "/totp/enable") ||
			normalize("/v1/auth/users/", ":totp:disable", "/totp/disable") ||
			normalize("/v1/auth/users/", ":totp:verify", "/totp/verify") ||
			normalize("/v1/auth/users/", ":kyc:complete", "/kyc/complete") ||
			normalize("/v1/auth/users/", ":kyc:submit-frame", "/kyc/submit-frame") ||
			normalize("/v1/auth/users/", ":profile:photo:upload-url", "/profile/photo/upload-url") ||
			normalize("/v1/auth/csrf", ":validate", "/csrf/validate") ||
			normalize("/v1/auth/biometric", ":authenticate", "/biometric/authenticate") ||
			normalize("/v1/partners/", ":credentials:rotate", "/credentials/rotate") ||
			normalize("/v1/kyc-verifications", ":pending", "/pending") ||
			normalize("/v1/media/", ":validate", "/validate") ||
			normalize("/v1/policies/", ":renew", "/renew") ||
			normalize("/v1/payments", ":reconcile", "/reconcile") ||
			normalize("/v1/analytics/reports/", ":generate", "/generate") ||
			normalize("/v1/analytics/reports", ":schedule", "/schedule") ||
			normalize("/v1/analytics/queries", ":run", "/run")

		// Avoid Go 1.22+ ServeMux route ambiguity while preserving public API path.
		if strings.HasPrefix(r.URL.Path, "/v1/policies/number/") &&
			!strings.HasSuffix(r.URL.Path, "/lookup") {
			id := strings.TrimPrefix(r.URL.Path, "/v1/policies/number/")
			if id != "" && !strings.Contains(id, "/") {
				r.URL.Path = "/v1/policies/number/" + id + "/lookup"
				r.URL.RawPath = r.URL.Path
			}
		}

		next.ServeHTTP(w, r)
	})
}
