using MediatR;
using PoliSync.Quotes.Domain;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Quotes.Application.Queries;

/// <summary>
/// Query to list quotations with optional filtering
/// </summary>
public sealed record ListQuotationsQuery(
    Guid? CustomerId = null,
    QuotationStatus? Status = null,
    int PageNumber = 1,
    int PageSize = 20
) : IRequest<Result<ListQuotationsResponse>>;

public sealed record ListQuotationsResponse(
    List<QuotationDto> Quotations,
    int TotalCount,
    int PageNumber,
    int PageSize
);

public sealed record QuotationDto(
    Guid Id,
    string QuotationNumber,
    Guid ProductId,
    Guid PlanId,
    Guid CustomerId,
    QuotationStatus Status,
    DateTime ExpiryDate,
    long BasePremium,
    long RiderPremium,
    long LoadingAmount,
    long DiscountAmount,
    long VatTax,
    long ServiceFee,
    long TotalPayable,
    string? RejectionReason,
    DateTime CreatedAt,
    DateTime UpdatedAt
);
