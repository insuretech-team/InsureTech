using MediatR;
using Insuretech.Products.Services.V1;
using Insuretech.Products.Entity.V1;

namespace InsuranceEngine.Products.Application.Queries;

public sealed record ListProductsQuery(ProductCategory? Category = null, string? SearchTerm = null, int Page = 1, int PageSize = 20) : IRequest<ListProductsResponse>;

public sealed record GetProductQuery(string ProductId) : IRequest<GetProductResponse>;

public sealed record GetProductByCodeQuery(string ProductCode) : IRequest<GetProductResponse>;
