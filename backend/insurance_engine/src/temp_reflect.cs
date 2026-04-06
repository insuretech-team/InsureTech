using System.Security.Cryptography;

try {
    Console.WriteLine($"AES-GCM Nonce Sizes: {AesGcm.NonceByteSizes.Min} - {AesGcm.NonceByteSizes.Max}");
} catch (Exception ex) {
    Console.WriteLine(ex.Message);
}
