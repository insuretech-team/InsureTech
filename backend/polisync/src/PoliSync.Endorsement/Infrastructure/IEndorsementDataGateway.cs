using EndorsementEntity = Insuretech.Endorsement.Entity.V1.Endorsement;

namespace PoliSync.Endorsement.Infrastructure;

public interface IEndorsementDataGateway
{
    Task<EndorsementEntity> CreateEndorsementAsync(EndorsementEntity endorsement, CancellationToken cancellationToken = default);
    Task<EndorsementEntity?> GetEndorsementAsync(string endorsementId, CancellationToken cancellationToken = default);
    Task<EndorsementEntity> UpdateEndorsementAsync(EndorsementEntity endorsement, CancellationToken cancellationToken = default);
    Task<IReadOnlyList<EndorsementEntity>> ListEndorsementsByPolicyAsync(string policyId, CancellationToken cancellationToken = default);
}
