using System.Collections.Generic;
using MediatR;
using InsuranceEngine.Application.DTOs;

namespace InsuranceEngine.Application.Features.Products.Queries.SearchProducts;

public record SearchProductsQuery(
    string? Query, 
    decimal? MinPremium, 
    decimal? MaxPremium
) : IRequest<List<ProductDto>>;
