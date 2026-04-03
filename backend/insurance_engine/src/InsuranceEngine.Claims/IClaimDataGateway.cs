using Insuretech.Claims.Entity.V1;
using Insuretech.Claims.Services.V1;

namespace InsuranceEngine.Claims;

public interface IClaimDataGateway
{
    Task<Claim?> GetClaimAsync(string claimId, CancellationToken ct = default);
    Task<ListClaimsResponse> ListClaimsAsync(ListClaimsRequest request, CancellationToken ct = default);
    Task<string> CreateClaimAsync(Claim claim, CancellationToken ct = default);
    Task UpdateClaimAsync(Claim claim, CancellationToken ct = default);
    Task ApproveClaimAsync(string claimId, string notes, CancellationToken ct = default);
    Task RejectClaimAsync(string claimId, string reason, CancellationToken ct = default);
}
