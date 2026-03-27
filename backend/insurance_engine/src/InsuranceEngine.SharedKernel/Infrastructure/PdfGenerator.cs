using System.Threading.Tasks;

namespace InsuranceEngine.SharedKernel.Infrastructure;

public interface IPdfGenerator
{
    Task<byte[]> GeneratePolicyDocumentAsync(string policyNumber, string customerName, string productName, decimal premium);
}

public class MockPdfGenerator : IPdfGenerator
{
    public Task<byte[]> GeneratePolicyDocumentAsync(string policyNumber, string customerName, string productName, decimal premium)
    {
        // Simulate PDF generation logic with a placeholder byte array
        var docInfo = $"Policy: {policyNumber}\nCustomer: {customerName}\nProduct: {productName}\nPremium: {premium}";
        return Task.FromResult(System.Text.Encoding.UTF8.GetBytes(docInfo));
    }
}
