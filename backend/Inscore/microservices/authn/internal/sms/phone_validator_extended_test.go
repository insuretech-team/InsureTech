package sms

import (
	"testing"
)

func TestNormalizePhoneNumberExtended(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		// ── Standard local format ─────────────────────────────────────────────
		{name: "local GP 017",          input: "01712345678",    want: "8801712345678", wantErr: false},
		{name: "local GP 013",          input: "01312345678",    want: "8801312345678", wantErr: false},
		{name: "local Robi 018",        input: "01812345678",    want: "8801812345678", wantErr: false},
		{name: "local Banglalink 019",  input: "01912345678",    want: "8801912345678", wantErr: false},
		{name: "local Banglalink 014",  input: "01412345678",    want: "8801412345678", wantErr: false},
		{name: "local Teletalk 015",    input: "01512345678",    want: "8801512345678", wantErr: false},
		{name: "local Teletalk 016",    input: "01612345678",    want: "8801612345678", wantErr: false},

		// ── With + prefix ────────────────────────────────────────────────────
		{name: "+880 format",           input: "+8801712345678", want: "8801712345678", wantErr: false},
		{name: "+880 Banglalink",       input: "+8801912345678", want: "8801912345678", wantErr: false},

		// ── 00880 international dialing prefix ───────────────────────────────
		{name: "00880 prefix GP",       input: "008801712345678", want: "8801712345678", wantErr: false},
		{name: "00880 prefix Robi",     input: "008801812345678", want: "8801812345678", wantErr: false},
		{name: "00880 Banglalink",      input: "008801912345678", want: "8801912345678", wantErr: false},

		// ── 880 without + ────────────────────────────────────────────────────
		{name: "880 bare GP",           input: "8801712345678",  want: "8801712345678", wantErr: false},
		{name: "880 bare Robi",         input: "8801812345678",  want: "8801812345678", wantErr: false},

		// ── 10-digit (no leading 0) ───────────────────────────────────────────
		{name: "10-digit 1712345678",   input: "1712345678",     want: "8801712345678", wantErr: false},
		{name: "10-digit Robi",         input: "1812345678",     want: "8801812345678", wantErr: false},

		// ── With separators (spaces, dashes, dots) ────────────────────────────
		{name: "spaces in local",       input: "017 1234 5678",  want: "8801712345678", wantErr: false},
		{name: "dashes in local",       input: "017-1234-5678",  want: "8801712345678", wantErr: false},
		{name: "dots in local",         input: "017.1234.5678",  want: "8801712345678", wantErr: false},
		{name: "+880 with dashes",      input: "+880-171-234-5678", want: "8801712345678", wantErr: false},
		{name: "00880 with spaces",     input: "00880 171 234 5678", want: "8801712345678", wantErr: false},
		{name: "local with parens",     input: "(017)12345678",  want: "8801712345678", wantErr: false},

		// ── Bad formats — should error ─────────────────────────────────────────
		{name: "empty string",          input: "",               want: "", wantErr: true},
		{name: "whitespace only",       input: "   ",            want: "", wantErr: true},
		{name: "letters only",          input: "abcdefgh",       want: "", wantErr: true},
		{name: "too short",             input: "0171234",        want: "", wantErr: true},
		{name: "too long local",        input: "017123456789",   want: "", wantErr: true},
		{name: "too long +880",         input: "+88017123456789", want: "", wantErr: true},
		{name: "wrong country code",    input: "+91712345678",   want: "", wantErr: true},
		{name: "invalid operator 012",  input: "01212345678",    want: "", wantErr: true},
		{name: "invalid operator 011",  input: "01112345678",    want: "", wantErr: true},
		{name: "invalid operator 020",  input: "02012345678",    want: "", wantErr: true},
		{name: "all zeros",             input: "00000000000",    want: "", wantErr: true},
		{name: "random garbage",        input: "xyz!@#$%",       want: "", wantErr: true},
		{name: "US number",             input: "+12025551234",   want: "", wantErr: true},
		{name: "UK number",             input: "+447911123456",  want: "", wantErr: true},
		{name: "9 digits only",         input: "171234567",      want: "", wantErr: true},
		{name: "11 digits no prefix",   input: "17123456789",    want: "", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NormalizePhoneNumber(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("NormalizePhoneNumber(%q) expected error, got %q", tc.input, got)
				}
			} else {
				if err != nil {
					t.Errorf("NormalizePhoneNumber(%q) unexpected error: %v", tc.input, err)
					return
				}
				if got != tc.want {
					t.Errorf("NormalizePhoneNumber(%q)\n  got:  %q\n  want: %q", tc.input, got, tc.want)
				}
			}
		})
	}
}

func TestValidatePhoneNumberExtended(t *testing.T) {
	valid := []string{
		"01712345678",
		"+8801712345678",
		"008801712345678",
		"8801812345678",
		"1712345678",
		"017-123-45678",
		"+880 171 234 5678",
	}
	for _, v := range valid {
		if !ValidatePhoneNumber(v) {
			t.Errorf("ValidatePhoneNumber(%q) = false, want true", v)
		}
	}

	invalid := []string{
		"",
		"abc",
		"01212345678",
		"+9101712345678",
		"1234",
		"00000000000",
	}
	for _, v := range invalid {
		if ValidatePhoneNumber(v) {
			t.Errorf("ValidatePhoneNumber(%q) = true, want false", v)
		}
	}
}

func TestGetOperatorExtended(t *testing.T) {
	tests := []struct {
		input    string
		wantOp   Operator
		wantErr  bool
	}{
		{"01712345678",    OperatorGrameenphone, false},
		{"01312345678",    OperatorGrameenphone, false},
		{"01812345678",    OperatorRobi,         false},
		{"01912345678",    OperatorBanglalink,   false},
		{"01412345678",    OperatorBanglalink,   false},
		{"01512345678",    OperatorTeletalk,     false},
		{"01612345678",    OperatorTeletalk,     false},
		{"+8801712345678", OperatorGrameenphone, false},
		{"invalid",        OperatorUnknown,      true},
		{"01212345678",    OperatorUnknown,      true},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			op, err := GetOperator(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("GetOperator(%q) expected error, got operator %q", tc.input, op)
				}
			} else {
				if err != nil {
					t.Errorf("GetOperator(%q) unexpected error: %v", tc.input, err)
					return
				}
				if op != tc.wantOp {
					t.Errorf("GetOperator(%q) = %q, want %q", tc.input, op, tc.wantOp)
				}
			}
		})
	}
}

func TestMaskPhoneNumber(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// len >= 10: show first 6 + *** + last 4
		{"8801712345678", "880171***5678"},
		{"8801812345678", "880181***5678"},
		{"1234567890",    "123456***7890"}, // exactly 10 chars
		// len < 10: returned as-is
		{"short",      "short"},
		{"123456789",  "123456789"}, // 9 chars — returned unchanged
		{"",           ""},
	}
	for _, tc := range tests {
		got := MaskPhoneNumber(tc.input)
		if got != tc.want {
			t.Errorf("MaskPhoneNumber(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestFormatForDisplay(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"01712345678",    "+880 171 234 5678"},
		{"+8801712345678", "+880 171 234 5678"},
		{"8801712345678",  "+880 171 234 5678"},
		{"invalid",        "invalid"},
	}
	for _, tc := range tests {
		got := FormatForDisplay(tc.input)
		if got != tc.want {
			t.Errorf("FormatForDisplay(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
