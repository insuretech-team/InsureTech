package events

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestPublisherCoversAllEventTypes verifies that Publisher exposes publish methods
// corresponding to all event messages defined in proto/auth_events.proto.
//
// It is intentionally defensive: it validates method existence + that each method
// is callable with zero-value arguments and returns nil (publisher is best-effort).
func TestPublisherCoversAllEventTypes(t *testing.T) {
	p := NewPublisher(nil)
	ctx := context.Background()

	// Keep this list in sync with proto message names in:
	// e:/Projects/InsureTech/proto/insuretech/authn/events/v1/auth_events.proto
	//
	// Naming convention: Message FooBarEvent => method PublishFooBar
	calls := []struct {
		method string
		args   []any
	}{
		{"PublishUserRegistered", []any{ctx, "u", "+8801", "a@b.com", "1.1.1.1", "WEB"}},
		{"PublishUserLoggedIn", []any{ctx, "u", "s", "JWT", "1.1.1.1", "MOBILE", "ua"}},
		{"PublishUserLoggedOut", []any{ctx, "u", "s", "JWT", "user_initiated", "1.1.1.1", "MOBILE"}},
		{"PublishTokenRefreshed", []any{ctx, "u", "s", "old", "new", "newr", "1.1.1.1", "MOBILE", "ua"}},
		{"PublishSessionRevoked", []any{ctx, "u", "s", "JWT", "system", "security"}},
		{"PublishPasswordChanged", []any{ctx, "u", "1.1.1.1", "u"}},
		{"PublishAccountLocked", []any{ctx, "u", "failed_login", time.Now().Add(time.Minute)}},
		{"PublishOTPSent", []any{ctx, "o", "+8801***", "login", "sms", "sslwireless", "", "", true}},
		{"PublishOTPVerified", []any{ctx, "o", "u", int32(1)}},
		{"PublishSMSDeliveryReport", []any{ctx, "o", "pmid", "8801***", "DELIVERED", "", "GP", time.Now()}},
		{"PublishCSRFValidationFailed", []any{ctx, "u", "s", "eh", "rh", "1.1.1.1", "ua", "/x", "POST"}},
		{"PublishSessionExpired", []any{ctx, "u", "s", "JWT", time.Now(), int32(10)}},
		{"PublishLoginFailed", []any{ctx, "u", "+8801", "invalid_password", "1.1.1.1", "MOBILE", "ua", int32(3)}},
		{"PublishPasswordResetRequested", []any{ctx, "u", "+8801", "1.1.1.1", "MOBILE"}},
		{"PublishEmailVerified", []any{ctx, "u", "user@domain.com"}},
		{"PublishEmailVerificationSent", []any{ctx, "u", "user@domain.com", "o2", "email_verification", "1.1.1.1"}},
		{"PublishEmailLoginSucceeded", []any{ctx, "u", "s", "user@domain.com", "SYSTEM_USER", "1.1.1.1", "ua", "dev"}},
		{"PublishEmailLoginFailed", []any{ctx, "user@domain.com", "invalid_otp", int32(2), "1.1.1.1", "ua"}},
		{"PublishPasswordResetByEmailRequested", []any{ctx, "u", "user@domain.com", "o3", "1.1.1.1"}},
	}

	pv := reflect.ValueOf(p)
	for _, tc := range calls {
		m := pv.MethodByName(tc.method)
		require.True(t, m.IsValid(), "missing method %s", tc.method)

		// Convert args to reflect.Value.
		in := make([]reflect.Value, 0, len(tc.args))
		for _, a := range tc.args {
			in = append(in, reflect.ValueOf(a))
		}

		out := m.Call(in)
		require.Len(t, out, 1, "expected single return (error) from %s", tc.method)
		if !out[0].IsNil() {
			err, ok := out[0].Interface().(error)
			require.True(t, ok)
			require.NoError(t, err, "method %s returned non-nil error", tc.method)
		}
	}
}
