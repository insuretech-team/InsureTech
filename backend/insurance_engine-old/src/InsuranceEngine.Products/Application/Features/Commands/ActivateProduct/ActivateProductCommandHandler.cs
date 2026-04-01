using System;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using InsuranceEngine.Products.Application.Interfaces;
using InsuranceEngine.Products.Domain.Events;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Interfaces;

namespace InsuranceEngine.Products.Application.Features.Commands.ActivateProduct;

public class ActivateProductCommandHandler : IRequestHandler<ActivateProductCommand, Result>
{
    private readonly IProductRepository _productRepository;
    private readonly IEventBus _eventBus;

    public ActivateProductCommandHandler(IProductRepository productRepository, IEventBus eventBus)
    {
        _productRepository = productRepository;
        _eventBus = eventBus;
    }

    public async Task<Result> Handle(ActivateProductCommand request, CancellationToken cancellationToken)
    {
        var product = await _productRepository.GetByIdAsync(request.ProductId);
        if (product == null)
            return Result.Fail(Error.NotFound("Product", request.ProductId.ToString()));

        var result = product.Activate();
        if (result.IsFailure)
            return result;

        await _productRepository.UpdateAsync(product);

        await _eventBus.PublishAsync("product.events", new ProductActivatedDomainEvent(
            ProductId: product.Id,
            ProductCode: product.ProductCode,
            ActivatedAt: DateTime.UtcNow
        ));

        return Result.Ok();
    }
}
