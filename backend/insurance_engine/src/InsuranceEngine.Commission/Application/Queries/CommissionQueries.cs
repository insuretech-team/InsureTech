using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Commission.Application.Queries;

public sealed record ListCommissionsQuery(
    string? AgentId,
    string? Status,
    int Page = 1,
    int PageSize = 20) : IQuery<ListCommissionsResult>;

public sealed record ListCommissionsResult(
    IReadOnlyList<CommissionDto> Items,
    int TotalCount);

public sealed record CommissionDto(
    string CommissionId,
    string PolicyId,
    string AgentId,
    decimal PremiumAmount,
    decimal CommissionRate,
    decimal CommissionAmount,
    string Status,
    DateTime? PaidAt,
    DateTime? CreatedAt)
{
    public CommissionDto() : this("", "", "", 0, 0, 0, "", null, null) { }
}

public sealed record GetCommissionQuery(string CommissionId) : IQuery<CommissionDto?>;
