using Grpc.Core;
using InsuranceEngine.Grpc.Clients;
using Insuretech.Beneficiary.Services.V1;

namespace InsuranceEngine.Beneficiary.Infrastructure;

/// <summary>
/// Implementation of IBeneficiaryDataGateway using gRPC calls to the Go backend.
/// </summary>
public sealed class GoBeneficiaryDataGateway : IBeneficiaryDataGateway
{
    private readonly InsuranceServiceClient _serviceClient;

    public GoBeneficiaryDataGateway(InsuranceServiceClient serviceClient)
    {
        _serviceClient = serviceClient;
    }

    public async Task<CreateIndividualBeneficiaryResponse> CreateIndividualBeneficiaryAsync(CreateIndividualBeneficiaryRequest request, CancellationToken ct = default)
    {
        return await _serviceClient.Beneficiaries.CreateIndividualBeneficiaryAsync(request, _serviceClient.BuildCallOptions(ct));
    }

    public async Task<CreateBusinessBeneficiaryResponse> CreateBusinessBeneficiaryAsync(CreateBusinessBeneficiaryRequest request, CancellationToken ct = default)
    {
        return await _serviceClient.Beneficiaries.CreateBusinessBeneficiaryAsync(request, _serviceClient.BuildCallOptions(ct));
    }

    public async Task<GetBeneficiaryResponse> GetBeneficiaryAsync(GetBeneficiaryRequest request, CancellationToken ct = default)
    {
        return await _serviceClient.Beneficiaries.GetBeneficiaryAsync(request, _serviceClient.BuildCallOptions(ct));
    }

    public async Task<UpdateBeneficiaryResponse> UpdateBeneficiaryAsync(UpdateBeneficiaryRequest request, CancellationToken ct = default)
    {
        return await _serviceClient.Beneficiaries.UpdateBeneficiaryAsync(request, _serviceClient.BuildCallOptions(ct));
    }

    public async Task<CompleteKYCResponse> CompleteKYCAsync(CompleteKYCRequest request, CancellationToken ct = default)
    {
        return await _serviceClient.Beneficiaries.CompleteKYCAsync(request, _serviceClient.BuildCallOptions(ct));
    }

    public async Task<UpdateRiskScoreResponse> UpdateRiskScoreAsync(UpdateRiskScoreRequest request, CancellationToken ct = default)
    {
        return await _serviceClient.Beneficiaries.UpdateRiskScoreAsync(request, _serviceClient.BuildCallOptions(ct));
    }
}
