using MediatR;
using PoliSync.Quotes.Infrastructure;
using PoliSync.SharedKernel.Auth;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Quotes.Application.Queries;

/// <summary>
/// Handler for listing quotations with filtering
/// </summary>
public sealed class ListQuotationsQueryHandler : IRequestHandler<ListQuotationsQuery, Result<ListQuotationsResponse>>
{
    private readonly IQuotationDataGateway _gateway;
    private readonly ICurrentUser _currentUser;

    public ListQuotationsQueryHandler(IQuotationDataGateway gateway, ICurrentUser currentUser)
    {
        _gateway = gateway;
        _currentUser = currentUser;
    }

    public async Task<Result<ListQuotationsResponse>> Handle(ListQuotationsQuery request, CancellationToken cancellationToken)
    {
        var tenantId = _currentUser.TenantId;

        var (quotations, totalCount) = await _gateway.ListAsync(
            tenantId,
            request.CustomerId,
            request.Status,
            request.PageNumber,
            request.PageSize,
            cancellationToken);

        var dtos = quotations.Select(q => new QuotationDto(
            q.Id,
            q.QuotationNumber,
            q.ProductId,
            q.PlanId,
            q.CustomerId,
            q.Status,
            q.ExpiryDate,
            q.BasePremium,
            q.RiderPremium,
            q.LoadingAmount,
            q.DiscountAmount,
            q.VatTax,
            q.ServiceFee,
            q.TotalPayable,
            q.RejectionReason,
            q.CreatedAt,
            q.UpdatedAt ?? q.CreatedAt
        )).ToList();

        var response = new ListQuotationsResponse(
            dtos,
            totalCount,
            request.PageNumber,
            request.PageSize
        );

        return Result.Ok(response);
    }
}
