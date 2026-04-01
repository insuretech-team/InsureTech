using Grpc.Core;
using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Products.Services.V1;
using InsuranceEngine.Products.Application.Queries;
using InsuranceEngine.Products.Application.Commands;
using Insuretech.Common.V1;

namespace InsuranceEngine.Products.GrpcServices;

public sealed class ProductGrpcService : ProductService.ProductServiceBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<ProductGrpcService> _logger;

    public ProductGrpcService(IMediator mediator, ILogger<ProductGrpcService> logger)
    {
        _mediator = mediator;
        _logger = logger;
    }

    public override async Task<ListProductsResponse> ListProducts(
        ListProductsRequest request, ServerCallContext context)
    {
        var query = new ListProductsQuery(
            Category: request.Category,
            Page: request.Page <= 0 ? 1 : request.Page,
            PageSize: request.PageSize <= 0 ? 10 : request.PageSize
        );

        return await _mediator.Send(query, context.CancellationToken);
    }

    public override async Task<GetProductResponse> GetProduct(
        GetProductRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.ProductId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Product ID is required"));
        }

        return await _mediator.Send(new GetProductQuery(request.ProductId), context.CancellationToken);
    }

    public override async Task<CreateProductResponse> CreateProduct(
        CreateProductRequest request, ServerCallContext context)
    {
        if (request.Product == null)
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Product data is required"));
        }

        var command = new CreateProductCommand(
            ProductCode: request.Product.ProductCode,
            ProductName: request.Product.ProductName,
            Category: request.Product.Category.ToString(),
            Description: request.Product.Description,
            BasePremium: (decimal)request.Product.BasePremium.Amount / 100m,
            MinSumInsured: (decimal)request.Product.MinSumInsured.Amount / 100m,
            MaxSumInsured: (decimal)request.Product.MaxSumInsured.Amount / 100m,
            MinTenureMonths: request.Product.MinTenureMonths,
            MaxTenureMonths: request.Product.MaxTenureMonths,
            CreatedBy: "System"
        );

        return await _mediator.Send(command, context.CancellationToken);
    }

    public override async Task<UpdateProductResponse> UpdateProduct(
        UpdateProductRequest request, ServerCallContext context)
    {
        if (request.Product == null || string.IsNullOrEmpty(request.Product.ProductId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Product data and ID are required"));
        }

        var command = new UpdateProductCommand(
            request.Product.ProductId,
            request.Product?.ProductName,
            request.Product?.Description,
            request.Product != null && request.Product.BasePremium != null ? (decimal)request.Product.BasePremium.Amount / 100m : null,
            request.Product != null && request.Product.MinSumInsured != null ? (decimal)request.Product.MinSumInsured.Amount / 100m : null,
            request.Product != null && request.Product.MaxSumInsured != null ? (decimal)request.Product.MaxSumInsured.Amount / 100m : null,
            request.Product?.MinTenureMonths,
            request.Product?.MaxTenureMonths);

        var result = await _mediator.Send(command, context.CancellationToken);

        if (result.IsFailure)
        {
            return new UpdateProductResponse
            {
                Error = new Error
                {
                    Code = result.Error?.Code ?? "UPDATE_FAILED",
                    Message = result.Error?.Message ?? "Failed to update product"
                }
            };
        }

        return new UpdateProductResponse { Message = "Product updated successfully" };
    }

    public override async Task<SearchProductsResponse> SearchProducts(
        SearchProductsRequest request, ServerCallContext context)
    {
        var query = new SearchProductsQuery(
            Query: request.Query,
            Category: request.Category,
            MinPremium: request.MinPremium,
            MaxPremium: request.MaxPremium
        );

        return await _mediator.Send(query, context.CancellationToken);
    }

    public override async Task<CalculatePremiumResponse> CalculatePremium(
        CalculatePremiumRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.ProductId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Product ID is required"));
        }

        var query = new CalculatePremiumQuery(
            ProductId: request.ProductId,
            SumInsured: request.SumInsured,
            TenureMonths: request.TenureMonths,
            RiderIds: request.RiderIds,
            ApplicantData: request.ApplicantData
        );

        return await _mediator.Send(query, context.CancellationToken);
    }

    public override async Task<ActivateProductResponse> ActivateProduct(
        ActivateProductRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.ProductId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Product ID is required"));
        }

        var result = await _mediator.Send(new ActivateProductCommand(request.ProductId), context.CancellationToken);

        if (result.IsFailure)
        {
            return new ActivateProductResponse
            {
                Error = new Error
                {
                    Code = result.Error?.Code ?? "ACTIVATION_FAILED",
                    Message = result.Error?.Message ?? "Failed to activate product"
                }
            };
        }

        return new ActivateProductResponse { Message = "Product activated successfully" };
    }

    public override async Task<DeactivateProductResponse> DeactivateProduct(
        DeactivateProductRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.ProductId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Product ID is required"));
        }

        var result = await _mediator.Send(new DeactivateProductCommand(request.ProductId), context.CancellationToken);

        if (result.IsFailure)
        {
            return new DeactivateProductResponse
            {
                Error = new Error
                {
                    Code = result.Error?.Code ?? "DEACTIVATION_FAILED",
                    Message = result.Error?.Message ?? "Failed to deactivate product"
                }
            };
        }

        return new DeactivateProductResponse { Message = "Product deactivated successfully" };
    }

    public override async Task<DiscontinueProductResponse> DiscontinueProduct(
        DiscontinueProductRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.ProductId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Product ID is required"));
        }

        var result = await _mediator.Send(new DiscontinueProductCommand(request.ProductId), context.CancellationToken);

        if (result.IsFailure)
        {
            return new DiscontinueProductResponse
            {
                Error = new Error
                {
                    Code = result.Error?.Code ?? "DISCONTINUE_FAILED",
                    Message = result.Error?.Message ?? "Failed to discontinue product"
                }
            };
        }

        return new DiscontinueProductResponse { Message = "Product discontinued successfully" };
    }
}
