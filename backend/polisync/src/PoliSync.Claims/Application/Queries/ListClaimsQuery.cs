using PoliSync.SharedKernel.CQRS;
using ClaimEntity = Insuretech.Claims.Entity.V1.Claim;

namespace PoliSync.Claims.Application.Queries;

public sealed record ListClaimsQuery(
    string CustomerId,
    string PolicyId,
    int PageNumber,
    int PageSize
) : IQuery<ClaimListResult>;

public sealed record ClaimListResult(List<ClaimEntity> Claims, int TotalCount);
