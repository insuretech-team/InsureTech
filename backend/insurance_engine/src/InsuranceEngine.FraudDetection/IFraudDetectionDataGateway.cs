using Insuretech.Fraud.Services.V1;

namespace InsuranceEngine.FraudDetection;

public interface IFraudDetectionDataGateway
{
    Task<CheckFraudResponse> CheckFraudAsync(CheckFraudRequest request, CancellationToken ct = default);
    Task<List<ClaimRecord>> GetRecentClaimsAsync(string customerId, int months, CancellationToken ct = default);
}

public class ClaimRecord
{
    public string ClaimId { get; set; } = string.Empty;
    public string ClaimType { get; set; } = string.Empty;
    public decimal ClaimAmount { get; set; }
    public DateTime SubmittedAt { get; set; }
}
