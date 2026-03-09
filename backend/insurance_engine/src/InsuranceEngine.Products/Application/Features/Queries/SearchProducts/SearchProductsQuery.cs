using System.Collections.Generic;
using MediatR;
using InsuranceEngine.Products.Application.DTOs;

namespace InsuranceEngine.Products.Application.Features.Queries.SearchProducts;

public record SearchProductsQuery(
    string? Query, 
    decimal? MinPremium, 
    decimal? MaxPremium
) : IRequest<List<ProductDto>>;
