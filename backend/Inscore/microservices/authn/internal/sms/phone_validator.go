package sms

import (
	"errors"
	"regexp"
	"strings"
)

// Bangladesh mobile number validation
// Valid operators:
// - Grameenphone (GP): 017, 013
// - Robi: 018
// - Banglalink: 019, 014
// - Teletalk: 015, 016

var (
	// ErrInvalidPhoneNumber indicates the phone number format is invalid
	ErrInvalidPhoneNumber = errors.New("invalid phone number format")

	// ErrUnsupportedOperator indicates the operator is not supported
	ErrUnsupportedOperator = errors.New("unsupported mobile operator")

	// validBDPhonePattern matches valid Bangladesh numbers in 880XXXXXXXXXX format
	validBDPhonePattern = regexp.MustCompile(`^880(13|14|15|16|17|18|19)\d{8}$`)
)

// Operator represents a Bangladesh mobile operator
type Operator string

const (
	OperatorGrameenphone Operator = "GRAMEENPHONE"
	OperatorRobi         Operator = "ROBI"
	OperatorBanglalink   Operator = "BANGLALINK"
	OperatorTeletalk     Operator = "TELETALK"
	OperatorUnknown      Operator = "UNKNOWN"
)

// NormalizePhoneNumber converts various Bangladesh phone formats to international format
// Input formats handled:
//   - 01712345678    → 8801712345678
//   - 8801712345678  → 8801712345678
//   - +8801712345678 → 8801712345678
//   - 00880171234... → 8801712345678
func NormalizePhoneNumber(phone string) (string, error) {
	if phone == "" {
		return "", ErrInvalidPhoneNumber
	}

	// Remove all non-numeric characters (spaces, dashes, plus, parentheses)
	cleaned := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")

	// Handle different input formats
	switch {
	case strings.HasPrefix(cleaned, "00880"):
		// Remove 00 prefix (international dialing)
		cleaned = cleaned[2:]
	case strings.HasPrefix(cleaned, "880"):
		// Already in international format (without +)
		// No change needed
	case strings.HasPrefix(cleaned, "0"):
		// Local format: 01XXXXXXXXX
		cleaned = "880" + cleaned[1:]
	case len(cleaned) == 10:
		// 10 digits without leading 0: 1XXXXXXXXX
		cleaned = "880" + cleaned
	default:
		// Unknown format
		return "", ErrInvalidPhoneNumber
	}

	// Validate the normalized number
	if !validBDPhonePattern.MatchString(cleaned) {
		return "", ErrInvalidPhoneNumber
	}

	return cleaned, nil
}

// ValidatePhoneNumber checks if a phone number is valid Bangladesh number
func ValidatePhoneNumber(phone string) bool {
	normalized, err := NormalizePhoneNumber(phone)
	if err != nil {
		return false
	}
	return validBDPhonePattern.MatchString(normalized)
}

// GetOperator returns the mobile operator for a given phone number
func GetOperator(phone string) (Operator, error) {
	normalized, err := NormalizePhoneNumber(phone)
	if err != nil {
		return OperatorUnknown, err
	}

	// Extract operator prefix (characters 3-5 of 880XXYYYYYYYY)
	if len(normalized) < 5 {
		return OperatorUnknown, ErrInvalidPhoneNumber
	}

	prefix := normalized[3:5]

	switch prefix {
	case "17", "13":
		return OperatorGrameenphone, nil
	case "18":
		return OperatorRobi, nil
	case "19", "14":
		return OperatorBanglalink, nil
	case "15", "16":
		return OperatorTeletalk, nil
	default:
		return OperatorUnknown, ErrUnsupportedOperator
	}
}

// MaskPhoneNumber masks a phone number for logging (privacy)
// Example: 8801712345678 → 880171***5678
func MaskPhoneNumber(phone string) string {
	if len(phone) < 10 {
		return phone
	}
	// Show first 6 and last 4 characters
	return phone[:6] + "***" + phone[len(phone)-4:]
}

// FormatForDisplay formats phone number for user display
// Example: 8801712345678 → +880 171 234 5678
func FormatForDisplay(phone string) string {
	normalized, err := NormalizePhoneNumber(phone)
	if err != nil {
		return phone
	}
	if len(normalized) != 13 {
		return phone
	}
	return "+" + normalized[:3] + " " + normalized[3:6] + " " + normalized[6:9] + " " + normalized[9:]
}

// Note: ValidateMSISDN, normalizeMSISDN, detectCarrier, DLRStatusPending
// are defined in sslwireless.go - do not duplicate here
