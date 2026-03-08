using Grpc.Core;
using Insuretech.Insurance.Services.V1;
using PoliSync.Infrastructure.Clients;
using EndorsementEntity = Insuretech.Endorsement.Entity.V1.Endorsement;

namespace PoliSync.Endorsement.Infrastructure;

public sealed class GoEndorsementDataGateway : IEndorsementDataGateway
{
    private readonly InsuranceServiceClient _insuranceClient;

    public GoEndorsementDataGateway(InsuranceServiceClient insuranceClient)
    {
        _insuranceClient = insuranceClient;
    }

    public async Task<EndorsementEntity> CreateEndorsementAsync(EndorsementEntity endorsement, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.CreateEndorsementAsync(
            new CreateEndorsementRequest { Endorsement = endorsement },
            cancellationToken: cancellationToken);

        return response.Endorsement;
    }

    public async Task<EndorsementEntity?> GetEndorsementAsync(string endorsementId, CancellationToken cancellationToken = default)
    {
        try
        {
            var response = await _insuranceClient.Client.GetEndorsementAsync(
                new GetEndorsementRequest { EndorsementId = endorsementId },
                cancellationToken: cancellationToken);

            return response.Endorsement;
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            return null;
        }
    }

    public async Task<EndorsementEntity> UpdateEndorsementAsync(EndorsementEntity endorsement, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.UpdateEndorsementAsync(
            new UpdateEndorsementRequest { Endorsement = endorsement },
            cancellationToken: cancellationToken);

        return response.Endorsement;
    }

    public async Task<IReadOnlyList<EndorsementEntity>> ListEndorsementsByPolicyAsync(string policyId, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.ListEndorsementsByPolicyAsync(
            new ListEndorsementsByPolicyRequest { PolicyId = policyId },
            cancellationToken: cancellationToken);

        return response.Endorsements;
    }
}
