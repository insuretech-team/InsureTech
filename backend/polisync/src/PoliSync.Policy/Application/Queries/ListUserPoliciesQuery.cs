using MediatR;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Policy.Application.Queries;

public sealed record ListUserPoliciesQuery(
    string UserId,
    int PageNumber,
    int PageSize
) : IQuery<PolicyListResult>;

public sealed record PolicyListResult(
    List<Insuretech.Policy.Entity.V1.Policy> Policies,
    int TotalCount
);
