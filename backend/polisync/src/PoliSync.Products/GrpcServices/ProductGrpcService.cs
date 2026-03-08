using Grpc.Core;
using Insuretech.Products.Entity.V1;
using Insuretech.Products.Services.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Products.Application.Commands;
using PoliSync.Products.Application.Queries;

namespace PoliSync.Products.GrpcServices;

public sealed class ProductGrpcService : ProductService.ProductServiceBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<ProductGrpcService> _logger;

    public ProductGrpcService(IMediator mediator, ILogger<ProductGrpcService> logger)
    {
        _mediator = mediator;
        _logger = logger;
    }

    public override async Task<CreateProductResponse> CreateProduct(CreateProductRequest request, ServerCallContext context)
    {
        if (request.Product == null)
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Product payload is required"));
        }

        var command = new CreateProductCommand
        {
            ProductCode = request.Product.ProductCode,
            ProductName = request.Product.ProductName,
            Category = request.Product.Category,
            Description = request.Product.Description,
            BasePremiumAmount = request.Product.BasePremium?.Amount ?? 0,
            MinSumInsuredAmount = request.Product.MinSumInsured?.Amount ?? 0,
            MaxSumInsuredAmount = request.Product.MaxSumInsured?.Amount ?? 0,
            MinTenureMonths = request.Product.MinTenureMonths,
            MaxTenureMonths = request.Product.MaxTenureMonths,
            Exclusions = new List<string>(request.Product.Exclusions),
            CreatedBy = request.Product.CreatedBy
        };

        var result = await _mediator.Send(command, context.CancellationToken);
        if (result.IsFailure || result.Value == null)
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, result.Error?.Message ?? "Failed to create product"));
        }

        return new CreateProductResponse
        {
            ProductId = result.Value.ProductId,
            Message = "Product created"
        };
    }

    public override async Task<GetProductResponse> GetProduct(GetProductRequest request, ServerCallContext context)
    {
        var result = await _mediator.Send(new GetProductQuery { ProductId = request.ProductId }, context.CancellationToken);
        if (result.IsFailure || result.Value == null)
        {
            throw new RpcException(new Status(StatusCode.NotFound, result.Error?.Message ?? "Product not found"));
        }

        return new GetProductResponse { Product = result.Value };
    }

    public override async Task<ListProductsResponse> ListProducts(ListProductsRequest request, ServerCallContext context)
    {
        var result = await _mediator.Send(new ListProductsQuery
        {
            Page = request.Page <= 0 ? 1 : request.Page,
            PageSize = request.PageSize <= 0 ? 50 : request.PageSize,
            Category = request.Category == ProductCategory.Unspecified ? null : request.Category
        }, context.CancellationToken);

        if (result.IsFailure || result.Value == null)
        {
            throw new RpcException(new Status(StatusCode.Internal, result.Error?.Message ?? "Failed to list products"));
        }

        var response = new ListProductsResponse
        {
            TotalCount = result.Value.TotalCount
        };
        response.Products.AddRange(result.Value.Products);
        return response;
    }

    public override async Task<UpdateProductResponse> UpdateProduct(UpdateProductRequest request, ServerCallContext context)
    {
        if (request.Product == null)
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Product payload is required"));
        }

        var command = new UpdateProductCommand
        {
            ProductId = request.Product.ProductId,
            ProductName = string.IsNullOrWhiteSpace(request.Product.ProductName) ? null : request.Product.ProductName,
            Description = string.IsNullOrWhiteSpace(request.Product.Description) ? null : request.Product.Description,
            BasePremiumAmount = request.Product.BasePremium?.Amount,
            MinSumInsuredAmount = request.Product.MinSumInsured?.Amount,
            MaxSumInsuredAmount = request.Product.MaxSumInsured?.Amount,
            MinTenureMonths = request.Product.MinTenureMonths,
            MaxTenureMonths = request.Product.MaxTenureMonths,
            Exclusions = request.Product.Exclusions.Count == 0 ? null : new List<string>(request.Product.Exclusions),
            Status = request.Product.Status == ProductStatus.Unspecified ? null : request.Product.Status
        };

        var result = await _mediator.Send(command, context.CancellationToken);
        if (result.IsFailure)
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, result.Error?.Message ?? "Failed to update product"));
        }

        return new UpdateProductResponse
        {
            Message = "Product updated"
        };
    }

    public override async Task<ActivateProductResponse> ActivateProduct(ActivateProductRequest request, ServerCallContext context)
    {
        var result = await _mediator.Send(new UpdateProductCommand
        {
            ProductId = request.ProductId,
            Status = ProductStatus.Active
        }, context.CancellationToken);

        if (result.IsFailure)
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, result.Error?.Message ?? "Failed to activate product"));
        }

        return new ActivateProductResponse { Message = "Product activated" };
    }

    public override async Task<DeactivateProductResponse> DeactivateProduct(DeactivateProductRequest request, ServerCallContext context)
    {
        var result = await _mediator.Send(new UpdateProductCommand
        {
            ProductId = request.ProductId,
            Status = ProductStatus.Inactive
        }, context.CancellationToken);

        if (result.IsFailure)
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, result.Error?.Message ?? "Failed to deactivate product"));
        }

        return new DeactivateProductResponse { Message = "Product deactivated" };
    }
}
