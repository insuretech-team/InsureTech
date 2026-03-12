namespace InsuranceEngine.SharedKernel.Interfaces;

/// <summary>
/// Service for encrypting/decrypting PII fields (NID, phone numbers) at rest.
/// Uses AES-256 encryption.
/// </summary>
public interface IEncryptionService
{
    /// <summary>
    /// Encrypt a plaintext string. Returns base64-encoded ciphertext.
    /// </summary>
    string Encrypt(string plaintext);

    /// <summary>
    /// Decrypt a base64-encoded ciphertext back to plaintext.
    /// </summary>
    string Decrypt(string ciphertext);

    Task<string> EncryptAsync(string plaintext);
    Task<string> DecryptAsync(string ciphertext);
}
