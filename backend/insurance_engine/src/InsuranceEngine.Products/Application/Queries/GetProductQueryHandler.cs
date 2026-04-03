using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Products.Services.V1;

namespace InsuranceEngine.Products.Application.Queries;

public sealed class GetProductQueryHandler : IRequestHandler<GetProductQuery, GetProductResponse>
{
    private readonly IProductDataGateway _gateway;
    private readonly ILogger<GetProductQueryHandler> _logger;

    public GetProductQueryHandler(IProductDataGateway gateway, ILogger<GetProductQueryHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<GetProductResponse> Handle(GetProductQuery request, CancellationToken cancellationToken)
    {
        try
        {
            var product = await _gateway.GetProductAsync(request.ProductId, cancellationToken);

            if (product == null)
            {
                return new GetProductResponse
                {
                    Error = new Insuretech.Common.V1.Error
                    {
                        Code = "PRODUCT_NOT_FOUND",
                        Message = "Product not found"
                    }
                };
            }

            return new GetProductResponse { Product = product };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to get product {ProductId}", request.ProductId);
            throw;
        }
    }
}
