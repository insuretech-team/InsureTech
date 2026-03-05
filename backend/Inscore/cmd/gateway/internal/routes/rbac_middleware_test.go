package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// nextOK is a simple next handler that records it was called.
func nextOK(t *testing.T, called *bool) http.Handler {
	t.Helper()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*called = true
		w.WriteHeader(http.StatusOK)
	})
}

func requestWithUserType(userType string) *http.Request {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	if userType != "" {
		r.Header.Set("X-User-Type", userType)
	}
	return r
}

// ---------------------------------------------------------------------------
// SystemUserMiddleware
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// SystemUserMiddleware
// ---------------------------------------------------------------------------

func TestSystemUserMiddleware_AllowsSystemUser(t *testing.T) {
	called := false
	h := SystemUserMiddleware(nextOK(t, &called))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, requestWithUserType(UserTypeSystemUser))
	require.True(t, called)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestSystemUserMiddleware_RejectsAgent(t *testing.T) {
	called := false
	h := SystemUserMiddleware(nextOK(t, &called))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, requestWithUserType(UserTypeAgent))
	require.False(t, called)
	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestSystemUserMiddleware_RejectsB2C(t *testing.T) {
	called := false
	h := SystemUserMiddleware(nextOK(t, &called))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, requestWithUserType(UserTypeB2CCustomer))
	require.False(t, called)
	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestSystemUserMiddleware_RejectsPartner(t *testing.T) {
	called := false
	h := SystemUserMiddleware(nextOK(t, &called))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, requestWithUserType(UserTypePartner))
	require.False(t, called)
	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestSystemUserMiddleware_RejectsRegulator(t *testing.T) {
	called := false
	h := SystemUserMiddleware(nextOK(t, &called))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, requestWithUserType(UserTypeRegulator))
	require.False(t, called)
	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestSystemUserMiddleware_RejectsEmpty(t *testing.T) {
	called := false
	h := SystemUserMiddleware(nextOK(t, &called))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, requestWithUserType(""))
	require.False(t, called)
	require.Equal(t, http.StatusForbidden, w.Code)
}

// ---------------------------------------------------------------------------
// AgentOrSystemMiddleware
// ---------------------------------------------------------------------------

func TestAgentOrSystemMiddleware_AllowsAgent(t *testing.T) {
	called := false
	h := AgentOrSystemMiddleware(nextOK(t, &called))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, requestWithUserType(UserTypeAgent))
	require.True(t, called)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestAgentOrSystemMiddleware_AllowsSystemUser(t *testing.T) {
	called := false
	h := AgentOrSystemMiddleware(nextOK(t, &called))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, requestWithUserType(UserTypeSystemUser))
	require.True(t, called)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestAgentOrSystemMiddleware_RejectsB2C(t *testing.T) {
	called := false
	h := AgentOrSystemMiddleware(nextOK(t, &called))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, requestWithUserType(UserTypeB2CCustomer))
	require.False(t, called)
	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestAgentOrSystemMiddleware_RejectsPartner(t *testing.T) {
	called := false
	h := AgentOrSystemMiddleware(nextOK(t, &called))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, requestWithUserType(UserTypePartner))
	require.False(t, called)
	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestAgentOrSystemMiddleware_RejectsRegulator(t *testing.T) {
	called := false
	h := AgentOrSystemMiddleware(nextOK(t, &called))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, requestWithUserType(UserTypeRegulator))
	require.False(t, called)
	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestAgentOrSystemMiddleware_RejectsEmpty(t *testing.T) {
	called := false
	h := AgentOrSystemMiddleware(nextOK(t, &called))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, requestWithUserType(""))
	require.False(t, called)
	require.Equal(t, http.StatusForbidden, w.Code)
}

// ---------------------------------------------------------------------------
// BusinessOrSystemMiddleware
// ---------------------------------------------------------------------------

func TestBusinessOrSystemMiddleware_AllowsBusiness(t *testing.T) {
	called := false
	h := BusinessOrSystemMiddleware(nextOK(t, &called))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, requestWithUserType(UserTypeBusinessBeneficiary))
	require.True(t, called)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestBusinessOrSystemMiddleware_AllowsSystem(t *testing.T) {
	called := false
	h := BusinessOrSystemMiddleware(nextOK(t, &called))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, requestWithUserType(UserTypeSystemUser))
	require.True(t, called)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestBusinessOrSystemMiddleware_RejectsAgent(t *testing.T) {
	called := false
	h := BusinessOrSystemMiddleware(nextOK(t, &called))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, requestWithUserType(UserTypeAgent))
	require.False(t, called)
	require.Equal(t, http.StatusForbidden, w.Code)
}

// ---------------------------------------------------------------------------
// PartnerOrSystemMiddleware
// ---------------------------------------------------------------------------

func TestPartnerOrSystemMiddleware_AllowsPartner(t *testing.T) {
	called := false
	h := PartnerOrSystemMiddleware(nextOK(t, &called))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, requestWithUserType(UserTypePartner))
	require.True(t, called)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestPartnerOrSystemMiddleware_AllowsSystem(t *testing.T) {
	called := false
	h := PartnerOrSystemMiddleware(nextOK(t, &called))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, requestWithUserType(UserTypeSystemUser))
	require.True(t, called)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestPartnerOrSystemMiddleware_RejectsAgent(t *testing.T) {
	called := false
	h := PartnerOrSystemMiddleware(nextOK(t, &called))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, requestWithUserType(UserTypeAgent))
	require.False(t, called)
	require.Equal(t, http.StatusForbidden, w.Code)
}

// ---------------------------------------------------------------------------
// RegulatorOrSystemMiddleware
// ---------------------------------------------------------------------------

func TestRegulatorOrSystemMiddleware_AllowsRegulator(t *testing.T) {
	called := false
	h := RegulatorOrSystemMiddleware(nextOK(t, &called))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, requestWithUserType(UserTypeRegulator))
	require.True(t, called)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestRegulatorOrSystemMiddleware_AllowsSystem(t *testing.T) {
	called := false
	h := RegulatorOrSystemMiddleware(nextOK(t, &called))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, requestWithUserType(UserTypeSystemUser))
	require.True(t, called)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestRegulatorOrSystemMiddleware_RejectsB2C(t *testing.T) {
	called := false
	h := RegulatorOrSystemMiddleware(nextOK(t, &called))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, requestWithUserType(UserTypeB2CCustomer))
	require.False(t, called)
	require.Equal(t, http.StatusForbidden, w.Code)
}

// ---------------------------------------------------------------------------
// buildDomain helper
// ---------------------------------------------------------------------------

func TestBuildDomain_WithPortalAndTenant(t *testing.T) {
	require.Equal(t, "system:root", buildDomain("system", "root"))
	require.Equal(t, "agent:tenant-abc", buildDomain("agent", "tenant-abc"))
}

func TestBuildDomain_EmptyPortalDefaultsToB2C(t *testing.T) {
	require.Equal(t, "b2c:root", buildDomain("", "root"))
}

func TestBuildDomain_EmptyTenantDefaultsToRoot(t *testing.T) {
	require.Equal(t, "system:root", buildDomain("system", ""))
}

// ---------------------------------------------------------------------------
// buildObject helper
// ---------------------------------------------------------------------------

func TestBuildObject_WithResource(t *testing.T) {
	require.Equal(t, "svc:policy/create", buildObject("svc:policy", "create"))
	require.Equal(t, "svc:claim/approve", buildObject("svc:claim", "approve"))
}

func TestBuildObject_EmptyResourceBecomesWildcard(t *testing.T) {
	require.Equal(t, "svc:policy/*", buildObject("svc:policy", ""))
}

func TestBuildObject_SlashResourceBecomesWildcard(t *testing.T) {
	require.Equal(t, "svc:policy/*", buildObject("svc:policy", "/"))
}
