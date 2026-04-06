namespace InsuranceEngine.SharedKernel.Domain.Security;

/// <summary>
/// Service for encrypting and decrypting Personally Identifiable Information (PII).
/// Must be compatible with the Go backend's AES-GCM (8-byte nonce) implementation.
/// </summary>
public interface ICryptoService
{
    /// <summary>
    /// Encrypts the plaintext using AES-256-GCM and returns a base64 encoded string.
    /// Nonce is prepended to the ciphertext.
    /// </summary>
    string Encrypt(string plaintext);

    /// <summary>
    /// Decrypts the base64 encoded ciphertext using AES-256-GCM.
    /// </summary>
    string Decrypt(string ciphertext);
}
