using Insuretech.Beneficiary.Services.V1;

namespace InsuranceEngine.Beneficiary;

public interface IBeneficiaryDataGateway
{
    Task<CreateIndividualBeneficiaryResponse> CreateIndividualBeneficiaryAsync(CreateIndividualBeneficiaryRequest request, CancellationToken ct = default);
    Task<CreateBusinessBeneficiaryResponse> CreateBusinessBeneficiaryAsync(CreateBusinessBeneficiaryRequest request, CancellationToken ct = default);
    Task<GetBeneficiaryResponse> GetBeneficiaryAsync(GetBeneficiaryRequest request, CancellationToken ct = default);
    Task<UpdateBeneficiaryResponse> UpdateBeneficiaryAsync(UpdateBeneficiaryRequest request, CancellationToken ct = default);
    Task<CompleteKYCResponse> CompleteKYCAsync(CompleteKYCRequest request, CancellationToken ct = default);
    Task<UpdateRiskScoreResponse> UpdateRiskScoreAsync(UpdateRiskScoreRequest request, CancellationToken ct = default);
}
