using MediatR;
using Insuretech.Products.Services.V1;
using Insuretech.Products.Entity.V1;
using Insuretech.Common.V1;

namespace InsuranceEngine.Products.Application.Queries;

public sealed record ListProductsQuery(ProductCategory Category = ProductCategory.Unspecified, string? SearchTerm = null, int Page = 1, int PageSize = 20) : IRequest<ListProductsResponse>;

public sealed record GetProductQuery(string ProductId) : IRequest<GetProductResponse>;

public sealed record GetProductByCodeQuery(string ProductCode) : IRequest<GetProductResponse>;

public sealed record SearchProductsQuery(
    string? Query, 
    ProductCategory Category, 
    Money? MinPremium, 
    Money? MaxPremium) : IRequest<SearchProductsResponse>;

public sealed record CalculatePremiumQuery(
    string ProductId, 
    Money SumInsured, 
    int TenureMonths, 
    IEnumerable<string>? RiderIds, 
    IDictionary<string, string> ApplicantData) : IRequest<CalculatePremiumResponse>;
