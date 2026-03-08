using InsuranceEngine.Api.RequestModels;
using InsuranceEngine.Api.ResponseModels;

namespace InsuranceEngine.Application.Interfaces;

public interface IBeneficiaryService
{
    Task<ListBeneficiariesResponseV1> GetBeneficiariesAsync(int page, int pageSize, CancellationToken cancellationToken);

    Task<(GetBeneficiaryResponseV1 Response, bool NotFound)> GetBeneficiaryByIdAsync(
        string beneficiaryId,
        CancellationToken cancellationToken);

    Task<CreateBeneficiaryResponseV1> CreateIndividualBeneficiaryAsync(
        CreateIndividualBeneficiaryRequestV1 request,
        CancellationToken cancellationToken);

    Task<CreateBeneficiaryResponseV1> CreateBusinessBeneficiaryAsync(
        CreateBusinessBeneficiaryRequestV1 request,
        CancellationToken cancellationToken);

    Task<(UpdateBeneficiaryResponseV1 Response, bool NotFound)> UpdateBeneficiaryAsync(
        string beneficiaryId,
        UpdateBeneficiaryRequestV1 request,
        CancellationToken cancellationToken);
}
