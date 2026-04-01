using ClaimEntity = Insuretech.Claims.Entity.V1.Claim;
using ClaimApprovalEntity = Insuretech.Claims.Entity.V1.ClaimApproval;
using ClaimDocumentEntity = Insuretech.Claims.Entity.V1.ClaimDocument;

namespace PoliSync.Claims.Infrastructure;

public interface IClaimDataGateway
{
    Task<ClaimEntity> CreateClaimAsync(ClaimEntity claim, CancellationToken cancellationToken = default);
    Task<ClaimDocumentEntity> CreateClaimDocumentAsync(ClaimDocumentEntity document, CancellationToken cancellationToken = default);
    Task<ClaimApprovalEntity> CreateClaimApprovalAsync(ClaimApprovalEntity approval, CancellationToken cancellationToken = default);
    Task<ClaimEntity?> GetClaimAsync(string claimId, CancellationToken cancellationToken = default);
    Task<ClaimEntity> UpdateClaimAsync(ClaimEntity claim, CancellationToken cancellationToken = default);
    Task<IReadOnlyList<ClaimEntity>> ListClaimsAsync(string customerId, string policyId, int page, int pageSize, CancellationToken cancellationToken = default);
}
