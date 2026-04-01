using MediatR;
using Insuretech.Orders.Entity.V1;
using Insuretech.Orders.Services.V1;
using PoliSync.Orders.Infrastructure;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Orders.Application.Queries;

/// <summary>
/// Handler for listing orders with filtering
/// </summary>
public sealed class ListOrdersQueryHandler : IRequestHandler<ListOrdersQuery, Result<ListOrdersResponse>>
{
    private readonly IOrderDataGateway _gateway;

    public ListOrdersQueryHandler(IOrderDataGateway gateway)
    {
        _gateway = gateway;
    }

    public async Task<Result<ListOrdersResponse>> Handle(ListOrdersQuery request, CancellationToken cancellationToken)
    {
        var response = await _gateway.ListOrdersAsync(
            new ListOrdersRequest
            {
                CustomerId = request.CustomerId?.ToString() ?? string.Empty,
                Status = request.Status is null
                    ? OrderStatus.Unspecified
                    : OrderViewMapper.ToProtoStatus(request.Status.Value),
                PageSize = request.PageSize,
                PageToken = request.PageNumber <= 1 ? string.Empty : request.PageNumber.ToString()
            },
            cancellationToken);

        if (response.Error is not null && !string.IsNullOrWhiteSpace(response.Error.Code))
        {
            return Result.Fail<ListOrdersResponse>(response.Error.Code, response.Error.Message);
        }

        var mappedResponse = new ListOrdersResponse(
            response.Orders.Select(OrderViewMapper.ToDto).ToList(),
            response.TotalCount,
            request.PageNumber,
            request.PageSize
        );

        return Result.Ok(mappedResponse);
    }
}
