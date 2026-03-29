using Insuretech.Claims.Services.V1;
using MediatR;

namespace InsuranceEngine.Claims.Application.Queries;

public sealed record ListUserClaimsQuery(
    string? CustomerId,
    string? PolicyId,
    string? Status,
    int Page = 1,
    int PageSize = 20) : IRequest<ListUserClaimsResponse>;

public sealed record GetClaimQuery(string ClaimId) : IRequest<GetClaimResponse>;
