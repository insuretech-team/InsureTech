using Insuretech.Products.Entity.V1;
using MediatR;
using PoliSync.SharedKernel.CQRS;
using System.Collections.Generic;

namespace PoliSync.Products.Application.Queries;

public record ListProductsQuery : IRequest<Result<ListProductsResult>>
{
    public int Page { get; init; } = 1;
    public int PageSize { get; init; } = 50;
    public ProductCategory? Category { get; init; }
    public ProductStatus? Status { get; init; }
}

public record ListProductsResult
{
    public List<Product> Products { get; init; } = new();
    public int TotalCount { get; init; }
}
