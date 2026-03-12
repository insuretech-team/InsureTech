using System;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using InsuranceEngine.Products.Application.Interfaces;
using InsuranceEngine.Products.Domain.Events;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Interfaces;

namespace InsuranceEngine.Products.Application.Features.Commands.DiscontinueProduct;

public class DiscontinueProductCommandHandler : IRequestHandler<DiscontinueProductCommand, Result>
{
    private readonly IProductRepository _productRepository;
    private readonly IEventBus _eventBus;

    public DiscontinueProductCommandHandler(IProductRepository productRepository, IEventBus eventBus)
    {
        _productRepository = productRepository;
        _eventBus = eventBus;
    }

    public async Task<Result> Handle(DiscontinueProductCommand request, CancellationToken cancellationToken)
    {
        var product = await _productRepository.GetByIdAsync(request.ProductId);
        if (product == null)
            return Result.Fail(Error.NotFound("Product", request.ProductId.ToString()));

        var result = product.Discontinue();
        if (result.IsFailure)
            return result;

        await _productRepository.UpdateAsync(product);

        await _eventBus.PublishAsync("product.events", new ProductDiscontinuedDomainEvent(
            ProductId: product.Id,
            ProductCode: product.ProductCode,
            Reason: request.Reason,
            DiscontinuedAt: DateTime.UtcNow
        ));

        return Result.Ok();
    }
}
