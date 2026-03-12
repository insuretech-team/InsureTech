using System.Threading;
using System.Threading.Tasks;
using MediatR;
using InsuranceEngine.Products.Application.DTOs;
using InsuranceEngine.Products.Application.Interfaces;

namespace InsuranceEngine.Products.Application.Features.Queries.GetProduct;

public class GetProductQueryHandler : IRequestHandler<GetProductQuery, ProductDto?>
{
    private readonly IProductRepository _productRepository;

    public GetProductQueryHandler(IProductRepository productRepository)
    {
        _productRepository = productRepository;
    }

    public async Task<ProductDto?> Handle(GetProductQuery request, CancellationToken cancellationToken)
    {
        var product = await _productRepository.GetByIdWithRidersAsync(request.Id);
        return product?.ToDto();
    }
}
