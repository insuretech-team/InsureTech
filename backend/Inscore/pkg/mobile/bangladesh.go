package mobile

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidBangladeshMobile = errors.New("invalid Bangladesh mobile number")

	nonDigitPattern = regexp.MustCompile(`\D`)
	bdPhonePattern  = regexp.MustCompile(`^880(13|14|15|16|17|18|19)\d{8}$`)
)

// NormalizeBangladeshMobileDigits accepts local or international BD mobile
// formats and returns the canonical digits-only value: 8801XXXXXXXXX.
func NormalizeBangladeshMobileDigits(phone string) (string, error) {
	if phone == "" {
		return "", ErrInvalidBangladeshMobile
	}

	cleaned := nonDigitPattern.ReplaceAllString(phone, "")
	if cleaned == "" {
		return "", ErrInvalidBangladeshMobile
	}

	switch {
	case hasPrefix(cleaned, "00880"):
		cleaned = cleaned[2:]
	case hasPrefix(cleaned, "0088"):
		cleaned = "88" + cleaned[4:]
	case hasPrefix(cleaned, "880"):
		// already normalized
	case hasPrefix(cleaned, "88") && len(cleaned) == 13:
		cleaned = "88" + cleaned[2:]
	case hasPrefix(cleaned, "0"):
		cleaned = "880" + cleaned[1:]
	case len(cleaned) == 10:
		cleaned = "880" + cleaned
	default:
		return "", ErrInvalidBangladeshMobile
	}

	if !bdPhonePattern.MatchString(cleaned) {
		return "", ErrInvalidBangladeshMobile
	}

	return cleaned, nil
}

func NormalizeBangladeshMobileE164(phone string) (string, error) {
	normalized, err := NormalizeBangladeshMobileDigits(phone)
	if err != nil {
		return "", err
	}
	return "+" + normalized, nil
}

func ValidateBangladeshMobile(phone string) bool {
	_, err := NormalizeBangladeshMobileDigits(phone)
	return err == nil
}

func hasPrefix(value, prefix string) bool {
	return len(value) >= len(prefix) && value[:len(prefix)] == prefix
}
