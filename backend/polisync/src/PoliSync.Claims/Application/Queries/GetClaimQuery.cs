using PoliSync.SharedKernel.CQRS;
using ClaimEntity = Insuretech.Claims.Entity.V1.Claim;

namespace PoliSync.Claims.Application.Queries;

public sealed record GetClaimQuery(string ClaimId) : IQuery<ClaimEntity>;
