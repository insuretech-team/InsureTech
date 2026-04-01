using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using InsuranceEngine.Products.Application.DTOs;
using InsuranceEngine.Products.Application.Interfaces;

namespace InsuranceEngine.Products.Application.Features.Queries.SearchProducts;

public class SearchProductsQueryHandler : IRequestHandler<SearchProductsQuery, List<ProductListDto>>
{
    private readonly IProductRepository _productRepository;

    public SearchProductsQueryHandler(IProductRepository productRepository)
    {
        _productRepository = productRepository;
    }

    public async Task<List<ProductListDto>> Handle(SearchProductsQuery request, CancellationToken cancellationToken)
    {
        var products = await _productRepository.SearchAsync(request.Query, request.MinPremium, request.MaxPremium);
        return products.Select(p => p.ToListDto()).ToList();
    }
}
