using System.Security.Cryptography;
using System.Text;
using InsuranceEngine.SharedKernel.Domain.Security;
using Microsoft.Extensions.Configuration;
using Org.BouncyCastle.Crypto;
using Org.BouncyCastle.Crypto.Engines;
using Org.BouncyCastle.Crypto.Modes;
using Org.BouncyCastle.Crypto.Parameters;
using Org.BouncyCastle.Security;

namespace InsuranceEngine.SharedKernel.Infrastructure.Security;

/// <summary>
/// Implementation of ICryptoService using BouncyCastle to support 8-byte nonces.
/// Mirrors the Go backend's AES-256-GCM implementation (compactGCMNonceSize = 8).
/// </summary>
public class BouncyGcmCryptoService : ICryptoService
{
    private readonly byte[] _key;
    private const int NonceSize = 8;
    private const int TagSizeBits = 128;

    public BouncyGcmCryptoService(IConfiguration configuration)
    {
        var keyStr = configuration["Security:PIIEncryptionKey"] 
                     ?? "0123456789abcdef0123456789abcdef"; // Default 32-char key for dev
        
        _key = Encoding.UTF8.GetBytes(keyStr);
        
        if (_key.Length != 32)
        {
            throw new ArgumentException("PII Encryption Key must be exactly 32 bytes (AES-256).");
        }
    }

    public string Encrypt(string plaintext)
    {
        if (string.IsNullOrEmpty(plaintext)) return string.Empty;

        var nonce = new byte[NonceSize];
        using (var rng = RandomNumberGenerator.Create())
        {
            rng.GetBytes(nonce);
        }

        var plaintextBytes = Encoding.UTF8.GetBytes(plaintext);
        var cipher = CreateCipher(true, nonce);
        
        var ciphertext = new byte[cipher.GetOutputSize(plaintextBytes.Length)];
        var len = cipher.ProcessBytes(plaintextBytes, 0, plaintextBytes.Length, ciphertext, 0);
        cipher.DoFinal(ciphertext, len);

        // Prepend nonce to ciphertext (Go parity)
        var combined = new byte[nonce.Length + ciphertext.Length];
        Buffer.BlockCopy(nonce, 0, combined, 0, nonce.Length);
        Buffer.BlockCopy(ciphertext, 0, combined, nonce.Length, ciphertext.Length);

        return Convert.ToBase64String(combined);
    }

    public string Decrypt(string ciphertext)
    {
        if (string.IsNullOrEmpty(ciphertext)) return string.Empty;

        var combined = Convert.FromBase64String(ciphertext);
        if (combined.Length < NonceSize)
        {
            throw new InvalidCiphertextException("Invalid ciphertext length.");
        }

        var nonce = new byte[NonceSize];
        var ciphertextBytes = new byte[combined.Length - NonceSize];
        
        Buffer.BlockCopy(combined, 0, nonce, 0, NonceSize);
        Buffer.BlockCopy(combined, NonceSize, ciphertextBytes, 0, ciphertextBytes.Length);

        var cipher = CreateCipher(false, nonce);
        
        var plaintextBytes = new byte[cipher.GetOutputSize(ciphertextBytes.Length)];
        var len = cipher.ProcessBytes(ciphertextBytes, 0, ciphertextBytes.Length, plaintextBytes, 0);
        var finalLen = cipher.DoFinal(plaintextBytes, len);

        return Encoding.UTF8.GetString(plaintextBytes, 0, finalLen);
    }

    private GcmBlockCipher CreateCipher(bool forEncryption, byte[] nonce)
    {
        var cipher = new GcmBlockCipher(new AesEngine());
        var parameters = new AeadParameters(new KeyParameter(_key), TagSizeBits, nonce);
        cipher.Init(forEncryption, parameters);
        return cipher;
    }
}
