using System.Threading;
using System.Threading.Tasks;
using MediatR;
using InsuranceEngine.Application.Interfaces;
using InsuranceEngine.Application.DTOs;

namespace InsuranceEngine.Application.Features.Products.Queries.GetProduct;

public class GetProductQueryHandler : IRequestHandler<GetProductQuery, ProductDto?>
{
    private readonly IProductRepository _productRepository;

    public GetProductQueryHandler(IProductRepository productRepository)
    {
        _productRepository = productRepository;
    }

    public async Task<ProductDto?> Handle(GetProductQuery request, CancellationToken cancellationToken)
    {
        var p = await _productRepository.GetByIdAsync(request.Id);
        if (p == null) return null;

        return new ProductDto(
            p.Id, p.ProductCode, p.ProductName, p.ProductNameBn, p.Description, 
            p.Category, p.Status, p.MinSumInsured, p.MaxSumInsured, p.MinAge, p.MaxAge, p.InsurerId
        );
    }
}
