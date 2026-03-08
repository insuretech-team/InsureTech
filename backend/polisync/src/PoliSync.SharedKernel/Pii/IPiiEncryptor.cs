namespace PoliSync.SharedKernel.Pii;

/// <summary>
/// PII encryption service for sensitive data (NID, phone numbers)
/// Uses AES-256-GCM
/// </summary>
public interface IPiiEncryptor
{
    string Encrypt(string plainText);
    
    string Decrypt(string cipherText);
    
    byte[] EncryptBytes(byte[] plainBytes);
    
    byte[] DecryptBytes(byte[] cipherBytes);
}
