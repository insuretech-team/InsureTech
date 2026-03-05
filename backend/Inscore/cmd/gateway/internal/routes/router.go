package routes

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/handlers"
	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/middleware"
	"github.com/newage-saint/insuretech/backend/inscore/cmd/gateway/internal/resilience"
	"google.golang.org/grpc"
)

// NewRouter wires the main HTTP routing for the gateway.
// authnConn: gRPC connection to authn service (required)
// authzConn: gRPC connection to authz service (required for AuthZ enforcement; nil = portal-gate only)
func NewRouter(authnHandler *handlers.AuthnHandler, authnConn *grpc.ClientConn, authzConn *grpc.ClientConn, clientManager *resilience.ResilientClientManager, dlrHandler *handlers.DLRHandler) http.Handler {
	mux := http.NewServeMux()

	authMW := AuthMiddleware(authnConn)
	csrfMW := CSRFMiddleware(authnConn)
	otpRL := middleware.OTPRateLimit(3, time.Hour)
	loginRL := middleware.IPWindowRateLimit(10, time.Hour)
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
		mux.HandleFunc("POST /v1/auth/token:validate", authnHandler.ValidateToken)
		mux.HandleFunc("POST /v1/auth/password:reset", authnHandler.ResetPassword)
		mux.HandleFunc("POST /v1/auth/biometric:authenticate", authnHandler.BiometricAuthenticate)

		// Email auth public routes (business/system/partner/regulator portals)
		mux.HandleFunc("POST /v1/auth/email/register", registerRL(http.HandlerFunc(authnHandler.RegisterEmailUser)).ServeHTTP)
		mux.HandleFunc("POST /v1/auth/email/otp:send", authnHandler.SendEmailOTP)
		mux.HandleFunc("POST /v1/auth/email/verify", authnHandler.VerifyEmail)
		mux.HandleFunc("POST /v1/auth/email/login", loginRL(http.HandlerFunc(authnHandler.EmailLogin)).ServeHTTP)
		mux.HandleFunc("POST /v1/auth/email/password:reset-request", authnHandler.RequestPasswordResetByEmail)
		mux.HandleFunc("POST /v1/auth/email/password:reset", authnHandler.ResetPasswordByEmail)

		// ── PROTECTED — any authenticated user ─────────────────────────────
		mux.Handle("POST /v1/auth/logout", authMW(csrfMW(http.HandlerFunc(authnHandler.Logout))))
		mux.Handle("POST /v1/auth/csrf:validate", authMW(csrfMW(http.HandlerFunc(authnHandler.ValidateCSRF))))
		mux.Handle("POST /v1/auth/password:change", authMW(csrfMW(http.HandlerFunc(authnHandler.ChangePassword))))
		mux.Handle("GET /v1/auth/session/current", authMW(http.HandlerFunc(authnHandler.GetCurrentSession)))
		mux.Handle("GET /v1/auth/sessions/{session_id}", authMW(http.HandlerFunc(authnHandler.GetSession)))
		mux.Handle("DELETE /v1/auth/sessions/{session_id}", authMW(csrfMW(http.HandlerFunc(authnHandler.RevokeSession))))
		mux.Handle("GET /v1/auth/users/{user_id}/sessions", authMW(http.HandlerFunc(authnHandler.ListSessions)))
		mux.Handle("POST /v1/auth/users/{user_id}/sessions:revokeAll", authMW(csrfMW(http.HandlerFunc(authnHandler.RevokeAllSessions))))

		// Profile (any authenticated user)
		mux.Handle("POST /v1/auth/users/{user_id}/profile", authMW(http.HandlerFunc(authnHandler.CreateUserProfile)))
		mux.Handle("GET /v1/auth/users/{user_id}/profile", authMW(http.HandlerFunc(authnHandler.GetUserProfile)))
		mux.Handle("PATCH /v1/auth/users/{user_id}/profile", authMW(csrfMW(http.HandlerFunc(authnHandler.UpdateUserProfile))))
		mux.Handle("POST /v1/auth/users/{user_id}/profile/photo:upload-url", authMW(http.HandlerFunc(authnHandler.GetProfilePhotoUploadURL)))
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
		mux.Handle("POST /v1/auth/users/{user_id}/kyc", authMW(AgentOrSystemMiddleware(authzKYC(http.HandlerFunc(authnHandler.InitiateKYC)))))
		mux.Handle("GET /v1/auth/users/{user_id}/kyc", authMW(AgentOrSystemMiddleware(authzKYC(http.HandlerFunc(authnHandler.GetKYCStatus)))))
		mux.Handle("POST /v1/auth/users/{user_id}/kyc:submit-frame", authMW(AgentOrSystemMiddleware(authzKYC(http.HandlerFunc(authnHandler.SubmitKYCFrame)))))
		mux.Handle("POST /v1/auth/users/{user_id}/kyc:complete", authMW(AgentOrSystemMiddleware(authzKYC(http.HandlerFunc(authnHandler.CompleteKYCSession)))))
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

	// ── B2B APIs (auth + authz) ───────────────────────────────────────────────
	if b2bConn := getServiceConn(clientManager, "b2b"); b2bConn != nil {
		b2bHandler := handlers.NewB2BServiceHandler(b2bConn)
		authzB2B := authzMW("svc:b2b", PathSegmentExtractor("/v1/b2b/"))

		mux.Handle("GET /v1/b2b/purchase-orders/catalog", authMW(authzB2B(http.HandlerFunc(b2bHandler.ListPurchaseOrderCatalog))))
		mux.Handle("GET /v1/b2b/purchase-orders", authMW(authzB2B(http.HandlerFunc(b2bHandler.ListPurchaseOrders))))
		mux.Handle("GET /v1/b2b/purchase-orders/{purchase_order_id}", authMW(authzB2B(http.HandlerFunc(b2bHandler.GetPurchaseOrder))))
		mux.Handle("POST /v1/b2b/purchase-orders", authMW(csrfMW(authzB2B(http.HandlerFunc(b2bHandler.CreatePurchaseOrder)))))
		mux.Handle("GET /v1/b2b/departments", authMW(authzB2B(http.HandlerFunc(b2bHandler.ListDepartments))))
		mux.Handle("GET /v1/b2b/employees", authMW(authzB2B(http.HandlerFunc(b2bHandler.ListEmployees))))
		mux.Handle("GET /v1/b2b/employees/{employee_uuid}", authMW(authzB2B(http.HandlerFunc(b2bHandler.GetEmployee))))
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

	// ── Document Generation APIs (auth + authz) ──────────────────────────────
	if docgenConn := getServiceConn(clientManager, "docgen"); docgenConn != nil {
		docgenHandler := handlers.NewDocGenHandler(docgenConn)
		authzDocGen := authzMW("svc:document", PathSegmentExtractor("/v1/"))

		mux.Handle("POST /v1/documents", authMW(csrfMW(authzDocGen(http.HandlerFunc(docgenHandler.Generate)))))
		mux.Handle("GET /v1/documents/{document_id}", authMW(authzDocGen(http.HandlerFunc(docgenHandler.GetDocument))))
		mux.Handle("GET /v1/entities/{entity_type}/{entity_id}/documents", authMW(authzDocGen(http.HandlerFunc(docgenHandler.ListDocuments))))
		mux.Handle("GET /v1/documents/{document_id}/download", authMW(authzDocGen(http.HandlerFunc(docgenHandler.Download))))
		mux.Handle("DELETE /v1/documents/{document_id}", authMW(csrfMW(authzDocGen(http.HandlerFunc(docgenHandler.DeleteDocument)))))
		mux.Handle("POST /v1/document-templates", authMW(csrfMW(authzDocGen(http.HandlerFunc(docgenHandler.CreateTemplate)))))
		mux.Handle("GET /v1/document-templates/{template_id}", authMW(authzDocGen(http.HandlerFunc(docgenHandler.GetTemplate))))
		mux.Handle("GET /v1/document-templates", authMW(authzDocGen(http.HandlerFunc(docgenHandler.ListTemplates))))
		mux.Handle("PATCH /v1/document-templates/{template_id}", authMW(csrfMW(authzDocGen(http.HandlerFunc(docgenHandler.UpdateTemplate)))))
		mux.Handle("POST /v1/document-templates/{template_id}/deactivate", authMW(csrfMW(authzDocGen(http.HandlerFunc(docgenHandler.DeactivateTemplate)))))
		mux.Handle("DELETE /v1/document-templates/{template_id}", authMW(csrfMW(authzDocGen(http.HandlerFunc(docgenHandler.DeleteTemplate)))))
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
		allowedOriginsEnv = "http://localhost:3000,http://localhost:5173"
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

		_ = normalize("/v1/partners/", ":verify", "/verify") ||
			normalize("/v1/partners/", ":updateStatus", "/updateStatus") ||
			normalize("/v1/fraud-rules/", ":activate", "/activate") ||
			normalize("/v1/fraud-rules/", ":deactivate", "/deactivate")

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
