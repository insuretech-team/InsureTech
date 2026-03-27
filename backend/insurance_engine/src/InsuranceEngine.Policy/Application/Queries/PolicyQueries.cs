using MediatR;
using Insuretech.Policy.Services.V1;

namespace InsuranceEngine.Policy.Application.Queries;

public sealed record ListUserPoliciesQuery(
    string? CustomerId,
    string? Status,
    string? ProductId,
    int Page = 1,
    int PageSize = 20) : IRequest<ListUserPoliciesResponse>;

public sealed record GetPolicyQuery(string PolicyId) : IRequest<GetPolicyResponse>;
