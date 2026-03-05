package service

import (
	"testing"
)

func TestPasswordValidator_ValidPasswords(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{"strong alphanumeric+special", "MyP@ssw0rd!"},
		{"long password", "ThisIsALongPassword1!"},
		{"mixed case + number + symbol", "Secure#99Pass"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := validatePasswordStrength(tc.password); err != nil {
				t.Errorf("expected valid password %q to pass, got: %v", tc.password, err)
			}
		})
	}
}

func TestPasswordValidator_InvalidPasswords(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  string
	}{
		{"too short", "Ab1!", "at least 8 characters"},
		{"no uppercase", "myp@ssw0rd!", "uppercase"},
		{"no lowercase", "MYP@SSW0RD!", "lowercase"},
		{"no digit", "MyP@ssword!", "digit"},
		{"no special char", "MyPassw0rd", "special character"},
		{"empty", "", "at least 8 characters"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePasswordStrength(tc.password)
			if err == nil {
				t.Errorf("expected password %q to fail, but it passed", tc.password)
				return
			}
			if tc.wantErr != "" {
				found := false
				msg := err.Error()
				for i := 0; i < len(msg)-len(tc.wantErr)+1; i++ {
					if msg[i:i+len(tc.wantErr)] == tc.wantErr {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error containing %q, got: %v", tc.wantErr, err)
				}
			}
		})
	}
}
