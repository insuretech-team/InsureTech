using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Products.Application.Interfaces;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Products.Application.Features.Queries.GetProductCode;

public class GetProductCodeQueryHandler : IRequestHandler<GetProductCodeQuery, Result<string>>
{
    private readonly IProductRepository _repository;

    public GetProductCodeQueryHandler(IProductRepository repository)
    {
        _repository = repository;
    }

    public async Task<Result<string>> Handle(GetProductCodeQuery request, CancellationToken cancellationToken)
    {
        var product = await _repository.GetByIdAsync(request.ProductId);
        if (product == null)
            return Result.Fail<string>(Error.NotFound("Product", request.ProductId.ToString()));

        return Result.Ok(product.ProductCode);
    }
}
