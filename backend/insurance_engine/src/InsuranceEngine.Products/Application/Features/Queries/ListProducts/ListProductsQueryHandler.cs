using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using InsuranceEngine.Products.Application.DTOs;
using InsuranceEngine.Products.Application.Interfaces;

namespace InsuranceEngine.Products.Application.Features.Queries.ListProducts;

public class ListProductsQueryHandler : IRequestHandler<ListProductsQuery, PaginatedResponse<ProductListDto>>
{
    private readonly IProductRepository _productRepository;

    public ListProductsQueryHandler(IProductRepository productRepository)
    {
        _productRepository = productRepository;
    }

    public async Task<PaginatedResponse<ProductListDto>> Handle(ListProductsQuery request, CancellationToken cancellationToken)
    {
        var (items, totalCount) = await _productRepository.ListActiveAsync(request.Category, request.Page, request.PageSize);

        return new PaginatedResponse<ProductListDto>(
            Items: items.Select(p => p.ToListDto()).ToList(),
            TotalCount: totalCount,
            Page: request.Page,
            PageSize: request.PageSize
        );
    }
}
