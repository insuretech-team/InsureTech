using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using InsuranceEngine.Application.Interfaces;
using InsuranceEngine.Application.DTOs;

namespace InsuranceEngine.Application.Features.Products.Queries.SearchProducts;

public class SearchProductsQueryHandler : IRequestHandler<SearchProductsQuery, List<ProductDto>>
{
    private readonly IProductRepository _productRepository;

    public SearchProductsQueryHandler(IProductRepository productRepository)
    {
        _productRepository = productRepository;
    }

    public async Task<List<ProductDto>> Handle(SearchProductsQuery request, CancellationToken cancellationToken)
    {
        var products = await _productRepository.SearchAsync(request.Query, request.MinPremium, request.MaxPremium);
        return products.Select(p => new ProductDto(
            p.Id, p.ProductCode, p.ProductName, p.ProductNameBn, p.Description, 
            p.Category, p.Status, p.MinSumInsured, p.MaxSumInsured, p.MinAge, p.MaxAge, p.InsurerId
        )).ToList();
    }
}
