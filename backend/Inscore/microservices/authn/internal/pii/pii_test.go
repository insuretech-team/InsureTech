package pii

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// testEncryptor returns a deterministic Encryptor suitable for unit tests.
func testEncryptor(t *testing.T) *Encryptor {
	t.Helper()
	aesKey := []byte("01234567890123456789012345678901")  // 32 bytes
	hmacKey := []byte("abcdefghijklmnopqrstuvwxyz012345") // 32 bytes
	enc, err := NewEncryptor(aesKey, hmacKey)
	require.NoError(t, err)
	return enc
}

func TestNewEncryptor_BadKeyLengths(t *testing.T) {
	_, err := NewEncryptor([]byte("short"), []byte("01234567890123456789012345678901"))
	require.Error(t, err)

	_, err = NewEncryptor([]byte("01234567890123456789012345678901"), []byte("short"))
	require.Error(t, err)
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	enc := testEncryptor(t)

	plaintexts := []string{
		"01711234567",
		"user@example.com",
		"1234567890123456789", // NID
		"a",
	}

	for _, pt := range plaintexts {
		ct, err := enc.Encrypt(pt)
		require.NoError(t, err, "Encrypt(%q)", pt)
		require.NotEmpty(t, ct)
		require.NotEqual(t, pt, ct, "ciphertext should differ from plaintext")

		got, err := enc.Decrypt(ct)
		require.NoError(t, err, "Decrypt for %q", pt)
		require.Equal(t, pt, got)
	}
}

func TestEncrypt_ProducesUniqueCiphertexts(t *testing.T) {
	// Each Encrypt call must produce a different ciphertext (random nonce).
	enc := testEncryptor(t)
	ct1, err := enc.Encrypt("01711234567")
	require.NoError(t, err)
	ct2, err := enc.Encrypt("01711234567")
	require.NoError(t, err)
	require.NotEqual(t, ct1, ct2, "two encryptions of same plaintext must differ")
}

func TestEncrypt_EmptyReturnsError(t *testing.T) {
	enc := testEncryptor(t)
	_, err := enc.Encrypt("")
	require.ErrorIs(t, err, ErrEmptyPlaintext)
}

func TestDecrypt_EmptyReturnsEmpty(t *testing.T) {
	enc := testEncryptor(t)
	got, err := enc.Decrypt("")
	require.NoError(t, err)
	require.Empty(t, got)
}

func TestDecrypt_CorruptCiphertext(t *testing.T) {
	enc := testEncryptor(t)
	_, err := enc.Decrypt("not-valid-base64!!!")
	require.Error(t, err)
}

func TestDecrypt_WrongKey(t *testing.T) {
	enc := testEncryptor(t)
	ct, err := enc.Encrypt("secret-mobile")
	require.NoError(t, err)

	// Different key
	wrongKey := []byte("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX") // 32 bytes
	enc2, err := NewEncryptor(wrongKey, wrongKey)
	require.NoError(t, err)
	_, err = enc2.Decrypt(ct)
	require.Error(t, err, "decryption with wrong key must fail")
}

func TestBlindIndex_Deterministic(t *testing.T) {
	enc := testEncryptor(t)
	idx1 := enc.BlindIndex("01711234567")
	idx2 := enc.BlindIndex("01711234567")
	require.Equal(t, idx1, idx2, "blind index must be deterministic")
	require.Len(t, idx1, 64, "HMAC-SHA256 hex must be 64 chars")
}

func TestBlindIndex_DifferentValues(t *testing.T) {
	enc := testEncryptor(t)
	idx1 := enc.BlindIndex("01711234567")
	idx2 := enc.BlindIndex("01719999999")
	require.NotEqual(t, idx1, idx2)
}

func TestBlindIndex_Empty(t *testing.T) {
	enc := testEncryptor(t)
	require.Empty(t, enc.BlindIndex(""))
}

func TestEncryptIfNonEmpty(t *testing.T) {
	enc := testEncryptor(t)

	// Non-empty should encrypt
	ct, err := enc.EncryptIfNonEmpty("nid-12345")
	require.NoError(t, err)
	require.NotEmpty(t, ct)

	// Empty should pass through
	ct2, err := enc.EncryptIfNonEmpty("")
	require.NoError(t, err)
	require.Empty(t, ct2)
}

func TestDecryptIfNonEmpty(t *testing.T) {
	enc := testEncryptor(t)

	ct, _ := enc.Encrypt("mobile-number")
	got, err := enc.DecryptIfNonEmpty(ct)
	require.NoError(t, err)
	require.Equal(t, "mobile-number", got)

	// Empty input
	got2, err := enc.DecryptIfNonEmpty("")
	require.NoError(t, err)
	require.Empty(t, got2)
}

func TestGenerateBiometricToken(t *testing.T) {
	enc := testEncryptor(t)

	raw, encrypted, lookup, err := enc.GenerateBiometricToken()
	require.NoError(t, err)

	// raw is 32-byte hex = 64 chars
	require.Len(t, raw, 64)
	require.True(t, isHex(raw), "raw token should be hex")

	// encrypted is non-empty and different from raw
	require.NotEmpty(t, encrypted)
	require.NotEqual(t, raw, encrypted)

	// lookup is 64-char hex (HMAC-SHA256)
	require.Len(t, lookup, 64)
	require.True(t, isHex(lookup))

	// Round-trip: decrypt should recover raw
	decrypted, err := enc.Decrypt(encrypted)
	require.NoError(t, err)
	require.Equal(t, raw, decrypted)

	// lookup matches BlindIndex of raw
	require.Equal(t, enc.BlindIndex(raw), lookup)

	// Two calls produce different tokens
	raw2, _, _, _ := enc.GenerateBiometricToken()
	require.NotEqual(t, raw, raw2)
}

func isHex(s string) bool {
	if len(s) == 0 {
		return false
	}
	s = strings.ToLower(s)
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}
