using Insuretech.Commission.Services.V1;
using MediatR;

namespace InsuranceEngine.Commission.Application.Queries;

public sealed record GetCommissionQuery(string CommissionId) : IRequest<GetCommissionResponse>;

public sealed record ListCommissionsQuery(
    string RecipientType,
    string RecipientId,
    string? Status,
    string? StartDate,
    string? EndDate,
    int Page = 1,
    int PageSize = 20) : IRequest<ListCommissionsResponse>;

public sealed record GetCommissionStatementQuery(
    string RecipientId,
    string PeriodStart,
    string PeriodEnd) : IRequest<GetCommissionStatementResponse>;

public sealed record GetRevenueShareReportQuery(
    string InsurerId,
    string StartDate,
    string EndDate) : IRequest<GetRevenueShareReportResponse>;
