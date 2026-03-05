package service

import (
	"errors"
	"strconv"
	"unicode"
)

// passwordPolicy defines the minimum requirements for a valid password.
const (
	minPasswordLength = 8
	maxPasswordLength = 128
)

// validatePasswordStrength checks that a password meets the minimum security policy:
//   - At least 8 characters, no more than 128
//   - At least one uppercase letter
//   - At least one lowercase letter
//   - At least one digit
//   - At least one special character
//
// It returns a descriptive error if the policy is not met, nil otherwise.
func validatePasswordStrength(password string) error {
	if len(password) < minPasswordLength {
		return errors.New("password must be at least " + strconv.Itoa(minPasswordLength) + " characters long")
	}
	if len(password) > maxPasswordLength {
		return errors.New("password must be no more than " + strconv.Itoa(maxPasswordLength) + " characters long")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return errors.New("password must contain at least one digit")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}
