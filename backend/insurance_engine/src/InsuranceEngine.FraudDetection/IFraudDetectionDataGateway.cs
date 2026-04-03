using Insuretech.Fraud.Services.V1;

namespace InsuranceEngine.FraudDetection;

public interface IFraudDetectionDataGateway
{
    Task<CheckFraudResponse> CheckFraudAsync(CheckFraudRequest request, CancellationToken ct = default);
}
