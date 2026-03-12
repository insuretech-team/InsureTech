using System.Security.Cryptography;
using System.Text;
using InsuranceEngine.SharedKernel.Interfaces;
using Microsoft.Extensions.Configuration;

namespace InsuranceEngine.SharedKernel.Services;

/// <summary>
/// AES-256-CBC encryption service for PII fields.
/// Reads the encryption key from IConfiguration["Encryption:Key"].
/// </summary>
public class AesEncryptionService : IEncryptionService
{
    private readonly byte[] _key;

    public AesEncryptionService(IConfiguration configuration)
    {
        var keyString = configuration["Encryption:Key"]
            ?? "LabAidInsureTechDefaultDevKey32B"; // 32 bytes for AES-256, dev only
        _key = Encoding.UTF8.GetBytes(keyString.PadRight(32).Substring(0, 32));
    }

    public string Encrypt(string plaintext)
    {
        if (string.IsNullOrEmpty(plaintext))
            return plaintext;

        using var aes = Aes.Create();
        aes.Key = _key;
        aes.GenerateIV();

        using var encryptor = aes.CreateEncryptor();
        var plaintextBytes = Encoding.UTF8.GetBytes(plaintext);
        var cipherBytes = encryptor.TransformFinalBlock(plaintextBytes, 0, plaintextBytes.Length);

        // Prepend IV to ciphertext for storage
        var result = new byte[aes.IV.Length + cipherBytes.Length];
        aes.IV.CopyTo(result, 0);
        cipherBytes.CopyTo(result, aes.IV.Length);

        return Convert.ToBase64String(result);
    }

    public string Decrypt(string ciphertext)
    {
        if (string.IsNullOrEmpty(ciphertext))
            return ciphertext;

        var fullCipher = Convert.FromBase64String(ciphertext);

        using var aes = Aes.Create();
        aes.Key = _key;

        // Extract IV from first 16 bytes
        var iv = new byte[16];
        Array.Copy(fullCipher, 0, iv, 0, 16);
        aes.IV = iv;

        var cipherBytes = new byte[fullCipher.Length - 16];
        Array.Copy(fullCipher, 16, cipherBytes, 0, cipherBytes.Length);

        using var decryptor = aes.CreateDecryptor();
        var plaintextBytes = decryptor.TransformFinalBlock(cipherBytes, 0, cipherBytes.Length);

        return Encoding.UTF8.GetString(plaintextBytes);
    }

    public Task<string> EncryptAsync(string plaintext) => Task.FromResult(Encrypt(plaintext));
    public Task<string> DecryptAsync(string ciphertext) => Task.FromResult(Decrypt(ciphertext));
}
