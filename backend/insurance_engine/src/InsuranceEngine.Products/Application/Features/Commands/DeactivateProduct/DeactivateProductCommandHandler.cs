using System;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using InsuranceEngine.Products.Application.Interfaces;
using InsuranceEngine.Products.Domain.Events;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Interfaces;

namespace InsuranceEngine.Products.Application.Features.Commands.DeactivateProduct;

public class DeactivateProductCommandHandler : IRequestHandler<DeactivateProductCommand, Result>
{
    private readonly IProductRepository _productRepository;
    private readonly IEventBus _eventBus;

    public DeactivateProductCommandHandler(IProductRepository productRepository, IEventBus eventBus)
    {
        _productRepository = productRepository;
        _eventBus = eventBus;
    }

    public async Task<Result> Handle(DeactivateProductCommand request, CancellationToken cancellationToken)
    {
        var product = await _productRepository.GetByIdAsync(request.ProductId);
        if (product == null)
            return Result.Fail(Error.NotFound("Product", request.ProductId.ToString()));

        var result = product.Deactivate();
        if (result.IsFailure)
            return result;

        await _productRepository.UpdateAsync(product);

        await _eventBus.PublishAsync("product.events", new ProductDeactivatedDomainEvent(
            ProductId: product.Id,
            ProductCode: product.ProductCode,
            Reason: request.Reason,
            DeactivatedAt: DateTime.UtcNow
        ));

        return Result.Ok();
    }
}
