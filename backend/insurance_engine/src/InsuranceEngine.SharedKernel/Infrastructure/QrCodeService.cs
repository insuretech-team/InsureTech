using QRCoder;
using System;
using System.Drawing;
using System.Drawing.Imaging;
using System.IO;

namespace InsuranceEngine.SharedKernel.Infrastructure;

public interface IQrCodeService
{
    byte[] GenerateQrCode(string text);
    string GenerateQrCodeBase64(string text);
}

public class QrCodeService : IQrCodeService
{
    public byte[] GenerateQrCode(string text)
    {
        using var qrGenerator = new QRCodeGenerator();
        using var qrCodeData = qrGenerator.CreateQrCode(text, QRCodeGenerator.ECCLevel.Q);
        using var qrCode = new PngByteQRCode(qrCodeData);
        return qrCode.GetGraphic(20);
    }

    public string GenerateQrCodeBase64(string text)
    {
        var bytes = GenerateQrCode(text);
        return Convert.ToBase64String(bytes);
    }
}
