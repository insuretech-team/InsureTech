using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using InsuranceEngine.Products.Application.Interfaces;
using InsuranceEngine.Products.Domain;

namespace InsuranceEngine.Products.Application.Features.Queries.ListInsurers;

public class ListInsurersQueryHandler : IRequestHandler<ListInsurersQuery, List<Insurer>>
{
    private readonly IProductRepository _productRepository;

    public ListInsurersQueryHandler(IProductRepository insurerRepository)
    {
        _productRepository = insurerRepository;
    }

    public async Task<List<Insurer>> Handle(ListInsurersQuery request, CancellationToken cancellationToken)
    {
        return await _productRepository.ListInsurersAsync();
    }
}
