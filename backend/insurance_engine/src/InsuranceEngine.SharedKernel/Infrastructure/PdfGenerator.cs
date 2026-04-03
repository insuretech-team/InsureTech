using System.Threading.Tasks;

namespace InsuranceEngine.SharedKernel.Infrastructure;

public interface IPdfGenerator
{
    Task<byte[]> GeneratePolicyDocumentAsync(string policyNumber, string customerName, string productName, decimal premium);
}

public class MockPdfGenerator : IPdfGenerator
{
    private readonly IQrCodeService _qrCodeService;

    public MockPdfGenerator(IQrCodeService qrCodeService)
    {
        _qrCodeService = qrCodeService;
    }

    public Task<byte[]> GeneratePolicyDocumentAsync(string policyNumber, string customerName, string productName, decimal premium)
    {
        // FR-035: Generate QR Code for policy verification
        var qrUrl = $"https://insuretech.labaid.com/verify/{policyNumber}";
        var qrBase64 = _qrCodeService.GenerateQrCodeBase64(qrUrl);

        // Simulate PDF generation with embedded QR metadata
        var docInfo = $"--- LABAID INSURETECH POLICY ---\n" +
                      $"Policy Number: {policyNumber}\n" +
                      $"Customer: {customerName}\n" +
                      $"Product: {productName}\n" +
                      $"Premium: {premium} BDT\n" +
                      $"QR Verification: [Embedded Base64: {qrBase64[..20]}...]";
        
        return Task.FromResult(System.Text.Encoding.UTF8.GetBytes(docInfo));
    }
}
