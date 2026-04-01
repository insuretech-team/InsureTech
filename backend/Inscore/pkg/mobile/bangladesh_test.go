package mobile

import "testing"

func TestNormalizeBangladeshMobileDigits(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "local", input: "01712345678", want: "8801712345678"},
		{name: "plus", input: "+8801712345678", want: "8801712345678"},
		{name: "double zero country code", input: "008801712345678", want: "8801712345678"},
		{name: "ten digits", input: "1712345678", want: "8801712345678"},
		{name: "with separators", input: "017-1234-5678", want: "8801712345678"},
		{name: "invalid", input: "01112345678", wantErr: true},
		{name: "empty", input: "", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NormalizeBangladeshMobileDigits(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestNormalizeBangladeshMobileE164(t *testing.T) {
	got, err := NormalizeBangladeshMobileE164("01712345678")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "+8801712345678" {
		t.Fatalf("got %q", got)
	}
}
