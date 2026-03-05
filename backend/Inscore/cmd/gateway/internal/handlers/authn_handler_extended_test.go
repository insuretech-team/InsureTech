package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// ---------------------------------------------------------------------------
// T5a: ResendOTP calls the gRPC client
// ---------------------------------------------------------------------------

func TestAuthnHandler_ResendOTP_CallsClient(t *testing.T) {
	called := false
	m := &mockAuthServiceClient{}
	m.resendOTPFn = func(ctx context.Context, in *authnservicev1.ResendOTPRequest, opts ...grpc.CallOption) (*authnservicev1.ResendOTPResponse, error) {
		called = true
		require.Equal(t, "otp-original-123", in.OriginalOtpId)
		return &authnservicev1.ResendOTPResponse{OtpId: "otp-new-456"}, nil
	}

	h := &AuthnHandler{client: m}
	body := bytes.NewBufferString(`{"original_otp_id":"otp-original-123","reason":"not_received"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/otp:resend", body)
	w := httptest.NewRecorder()

	h.ResendOTP(w, req)

	require.True(t, called)
	require.Equal(t, http.StatusOK, w.Code)
}

// ---------------------------------------------------------------------------
// T5b: JWT login should NOT set a session cookie
// ---------------------------------------------------------------------------

func TestAuthnHandler_Login_JWTResponse_NoSessionCookie(t *testing.T) {
	m := &mockAuthServiceClient{}
	m.loginFn = func(ctx context.Context, in *authnservicev1.LoginRequest, _ ...grpc.CallOption) (*authnservicev1.LoginResponse, error) {
		return &authnservicev1.LoginResponse{
			UserId:       "u1",
			SessionId:    "s1",
			SessionType:  "JWT", // <-- JWT, not SERVER_SIDE
			AccessToken:  "access.token.here",
			RefreshToken: "refresh.token.here",
		}, nil
	}

	h := &AuthnHandler{client: m}
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewBufferString(`{"mobile_number":"+8801700000001","password":"pass"}`))
	w := httptest.NewRecorder()

	h.Login(w, req)

	res := w.Result()
	require.Equal(t, http.StatusOK, res.StatusCode)
	// No session cookie should be set for JWT login
	for _, c := range res.Cookies() {
		require.NotEqual(t, sessionCookieName, c.Name, "JWT login must not set session cookie")
	}
}

// ---------------------------------------------------------------------------
// T5c: Register with invalid JSON body → 400
// ---------------------------------------------------------------------------

func TestAuthnHandler_Register_InvalidBody_Returns400(t *testing.T) {
	h := &AuthnHandler{client: &mockAuthServiceClient{}}
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(`{invalid json`))
	w := httptest.NewRecorder()

	h.Register(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

// ---------------------------------------------------------------------------
// T5d: VerifyOTP basic wiring — calls client and returns 200
// ---------------------------------------------------------------------------

func TestAuthnHandler_VerifyOTP_CallsClient(t *testing.T) {
	called := false
	m := &mockAuthServiceClient{}
	m.verifyOTPFn = func(ctx context.Context, in *authnservicev1.VerifyOTPRequest, opts ...grpc.CallOption) (*authnservicev1.VerifyOTPResponse, error) {
		called = true
		return &authnservicev1.VerifyOTPResponse{}, nil
	}

	h := &AuthnHandler{client: m}
	body := bytes.NewBufferString(`{"otp_id":"otp-123","code":"456789"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/otp:verify", body)
	w := httptest.NewRecorder()

	h.VerifyOTP(w, req)

	require.True(t, called)
	require.Equal(t, http.StatusOK, w.Code)
}

// ---------------------------------------------------------------------------
// T5e: EmailLogin sets session cookie
// ---------------------------------------------------------------------------

func TestAuthnHandler_EmailLogin_SetsSessionCookie(t *testing.T) {
	m := &mockAuthServiceClient{}
	m.emailLoginFn = func(ctx context.Context, in *authnservicev1.EmailLoginRequest, _ ...grpc.CallOption) (*authnservicev1.EmailLoginResponse, error) {
		return &authnservicev1.EmailLoginResponse{
			UserId:       "u2",
			SessionToken: "email-session-token",
			CsrfToken:    "email-csrf",
		}, nil
	}

	h := &AuthnHandler{client: m}
	body := bytes.NewBufferString(`{"email":"test@example.com","otp_id":"otp-abc","code":"123456"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/email/login", body)
	w := httptest.NewRecorder()

	h.EmailLogin(w, req)

	res := w.Result()
	require.Equal(t, http.StatusOK, res.StatusCode)

	var sessionCookie *http.Cookie
	for _, c := range res.Cookies() {
		if c.Name == sessionCookieName {
			sessionCookie = c
		}
	}
	require.NotNil(t, sessionCookie, "email login must set session cookie")
	require.Equal(t, "email-session-token", sessionCookie.Value)
	require.True(t, sessionCookie.HttpOnly)
	require.Equal(t, "email-csrf", res.Header.Get("X-CSRF-Token"))
}

// ---------------------------------------------------------------------------
// T5f: Logout clears the session cookie
// ---------------------------------------------------------------------------

func TestAuthnHandler_Logout_ClearsCookie(t *testing.T) {
	h := &AuthnHandler{client: &mockAuthServiceClient{}}
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/logout", bytes.NewBufferString(`{}`))
	w := httptest.NewRecorder()

	h.Logout(w, req)

	res := w.Result()
	require.Equal(t, http.StatusOK, res.StatusCode)

	var sessionCookie *http.Cookie
	for _, c := range res.Cookies() {
		if c.Name == sessionCookieName {
			sessionCookie = c
		}
	}
	require.NotNil(t, sessionCookie, "logout must set a cleared cookie")
	require.Equal(t, "", sessionCookie.Value)
	require.Equal(t, -1, sessionCookie.MaxAge)
}
