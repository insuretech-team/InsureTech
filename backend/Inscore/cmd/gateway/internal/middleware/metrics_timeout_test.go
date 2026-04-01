package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestTimeout_WritesBufferedResponseOnSuccess(t *testing.T) {
	handler := Timeout(200 * time.Millisecond)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test", "ok")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
	if rec.Header().Get("X-Test") != "ok" {
		t.Fatalf("expected buffered header to be copied")
	}
	if body := rec.Body.String(); body != `{"ok":true}` {
		t.Fatalf("expected buffered body, got %q", body)
	}
}

func TestTimeout_ReturnsGatewayTimeoutWithoutRacingLateWrites(t *testing.T) {
	handler := Timeout(20 * time.Millisecond)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(60 * time.Millisecond)
		w.Header().Set("X-Test", "late")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusGatewayTimeout {
		t.Fatalf("expected status %d, got %d", http.StatusGatewayTimeout, rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `"code":"DEADLINE_EXCEEDED"`) {
		t.Fatalf("expected timeout envelope, got %q", body)
	}
	if rec.Header().Get("X-Test") != "" {
		t.Fatalf("late handler headers should not leak into timeout response")
	}

	time.Sleep(80 * time.Millisecond)
	if rec.Code != http.StatusGatewayTimeout {
		t.Fatalf("late writes changed response code to %d", rec.Code)
	}
	if rec.Header().Get("X-Test") != "" {
		t.Fatalf("late header write leaked after timeout")
	}
}
