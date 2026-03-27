using InsuranceEngine.SharedKernel.CQRS;
using Insuretech.Claims.Services.V1;
using MediatR;

namespace InsuranceEngine.Claims.Application.Queries;

public sealed record ListUserClaimsQuery(
    string? PolicyId,
    string? Status,
    int Page = 1,
    int PageSize = 20) : IRequest<ListUserClaimsResponse>;

public sealed record ListClaimsResult(
    IReadOnlyList<ClaimDto> Items,
    int TotalCount);

public sealed record ClaimDto(
    string ClaimId,
    string ClaimNumber,
    string PolicyId,
    string ClaimType,
    decimal ClaimAmount,
    decimal? ApprovedAmount,
    decimal? SettlementAmount,
    string? Description,
    string Status,
    string? RejectionReason,
    DateTime? SettledAt,
    DateTime? CreatedAt)
{
    public ClaimDto() : this("", "", "", "", 0, null, null, null, "", null, null, null) { }
}

public sealed record GetClaimQuery(string ClaimId) : IRequest<GetClaimResponse>;
