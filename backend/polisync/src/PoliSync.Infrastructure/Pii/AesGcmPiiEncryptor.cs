using Microsoft.Extensions.Options;
using PoliSync.SharedKernel.Pii;
using System.Security.Cryptography;
using System.Text;

namespace PoliSync.Infrastructure.Pii;

/// <summary>
/// AES-256-GCM encryption for PII data (NID, phone numbers)
/// </summary>
public sealed class AesGcmPiiEncryptor : IPiiEncryptor
{
    private readonly byte[] _key;
    private const int NonceSize = 12; // 96 bits
    private const int TagSize = 16;   // 128 bits

    public AesGcmPiiEncryptor(IOptions<PiiEncryptionOptions> options)
    {
        var keyString = options.Value.EncryptionKey;
        
        if (string.IsNullOrEmpty(keyString))
            throw new ArgumentException("PII encryption key is not configured");

        // Key should be 32 bytes (256 bits) for AES-256
        _key = Convert.FromBase64String(keyString);
        
        if (_key.Length != 32)
            throw new ArgumentException("PII encryption key must be 32 bytes (256 bits)");
    }

    public string Encrypt(string plainText)
    {
        if (string.IsNullOrEmpty(plainText))
            return plainText;

        var plainBytes = Encoding.UTF8.GetBytes(plainText);
        var encryptedBytes = EncryptBytes(plainBytes);
        return Convert.ToBase64String(encryptedBytes);
    }

    public string Decrypt(string cipherText)
    {
        if (string.IsNullOrEmpty(cipherText))
            return cipherText;

        var cipherBytes = Convert.FromBase64String(cipherText);
        var decryptedBytes = DecryptBytes(cipherBytes);
        return Encoding.UTF8.GetString(decryptedBytes);
    }

    public byte[] EncryptBytes(byte[] plainBytes)
    {
        // Generate random nonce
        var nonce = new byte[NonceSize];
        RandomNumberGenerator.Fill(nonce);

        // Allocate buffer for nonce + ciphertext + tag
        var cipherBytes = new byte[NonceSize + plainBytes.Length + TagSize];
        
        // Copy nonce to beginning
        Buffer.BlockCopy(nonce, 0, cipherBytes, 0, NonceSize);

        using var aesGcm = new AesGcm(_key, TagSize);
        
        var cipherSpan = cipherBytes.AsSpan(NonceSize, plainBytes.Length);
        var tagSpan = cipherBytes.AsSpan(NonceSize + plainBytes.Length, TagSize);
        
        aesGcm.Encrypt(nonce, plainBytes, cipherSpan, tagSpan);

        return cipherBytes;
    }

    public byte[] DecryptBytes(byte[] cipherBytes)
    {
        if (cipherBytes.Length < NonceSize + TagSize)
            throw new ArgumentException("Invalid cipher text");

        // Extract nonce
        var nonce = cipherBytes.AsSpan(0, NonceSize);
        
        // Extract ciphertext
        var cipherTextLength = cipherBytes.Length - NonceSize - TagSize;
        var cipherText = cipherBytes.AsSpan(NonceSize, cipherTextLength);
        
        // Extract tag
        var tag = cipherBytes.AsSpan(NonceSize + cipherTextLength, TagSize);

        var plainBytes = new byte[cipherTextLength];

        using var aesGcm = new AesGcm(_key, TagSize);
        aesGcm.Decrypt(nonce, cipherText, tag, plainBytes);

        return plainBytes;
    }
}

public sealed class PiiEncryptionOptions
{
    public string EncryptionKey { get; set; } = string.Empty;
}
