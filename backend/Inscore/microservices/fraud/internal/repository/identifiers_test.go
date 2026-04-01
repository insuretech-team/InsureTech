package repository

import (
	"strings"
	"testing"
	"time"
)

func TestNewSequenceNumber(t *testing.T) {
	now := time.Date(2026, time.March, 13, 11, 22, 33, 0, time.UTC)

	first := newSequenceNumber("FAL", now)
	second := newSequenceNumber("FAL", now)

	if first == second {
		t.Fatalf("expected unique numbers, got %q", first)
	}
	if !strings.HasPrefix(first, "FAL-20260313-112233-") {
		t.Fatalf("unexpected prefix format: %q", first)
	}
	if len(first) <= len("FAL-20260313-112233-") {
		t.Fatalf("expected suffix in %q", first)
	}
}

func TestNormalizeAlertStatus(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: ""},
		{name: "already enum", in: "ALERT_STATUS_OPEN", want: "ALERT_STATUS_OPEN"},
		{name: "short token", in: "open", want: "ALERT_STATUS_OPEN"},
		{name: "unknown preserved", in: "something_else", want: "something_else"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeAlertStatus(tt.in); got != tt.want {
				t.Fatalf("normalizeAlertStatus(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
