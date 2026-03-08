using ClaimEntity = Insuretech.Claims.Entity.V1.Claim;

namespace PoliSync.Claims.Infrastructure;

public interface IClaimDataGateway
{
    Task<ClaimEntity> CreateClaimAsync(ClaimEntity claim, CancellationToken cancellationToken = default);
    Task<ClaimEntity?> GetClaimAsync(string claimId, CancellationToken cancellationToken = default);
    Task<ClaimEntity> UpdateClaimAsync(ClaimEntity claim, CancellationToken cancellationToken = default);
    Task<IReadOnlyList<ClaimEntity>> ListClaimsAsync(string customerId, string policyId, int page, int pageSize, CancellationToken cancellationToken = default);
}
