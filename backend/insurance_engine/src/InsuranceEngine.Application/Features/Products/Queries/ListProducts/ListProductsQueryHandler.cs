using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using InsuranceEngine.Application.Interfaces;
using InsuranceEngine.Application.DTOs;

namespace InsuranceEngine.Application.Features.Products.Queries.ListProducts;

public class ListProductsQueryHandler : IRequestHandler<ListProductsQuery, List<ProductDto>>
{
    private readonly IProductRepository _productRepository;

    public ListProductsQueryHandler(IProductRepository productRepository)
    {
        _productRepository = productRepository;
    }

    public async Task<List<ProductDto>> Handle(ListProductsQuery request, CancellationToken cancellationToken)
    {
        var products = await _productRepository.ListAsync();
        return products.Select(p => new ProductDto(
            p.Id, p.ProductCode, p.ProductName, p.ProductNameBn, p.Description, 
            p.Category, p.Status, p.MinSumInsured, p.MaxSumInsured, p.MinAge, p.MaxAge, p.InsurerId
        )).ToList();
    }
}
