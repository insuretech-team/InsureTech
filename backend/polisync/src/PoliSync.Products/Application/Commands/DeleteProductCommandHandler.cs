using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Products.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using System;
using System.Threading;
using System.Threading.Tasks;

namespace PoliSync.Products.Application.Commands;

public class DeleteProductCommandHandler : IRequestHandler<DeleteProductCommand, Result<bool>>
{
    private readonly IProductRepository _repository;
    private readonly ILogger<DeleteProductCommandHandler> _logger;

    public DeleteProductCommandHandler(
        IProductRepository repository,
        ILogger<DeleteProductCommandHandler> logger)
    {
        _repository = repository;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(DeleteProductCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var product = await _repository.GetByIdAsync(request.ProductId, cancellationToken);
            if (product == null)
            {
                return Result<bool>.Failure($"Product not found: {request.ProductId}");
            }

            await _repository.DeleteAsync(request.ProductId, cancellationToken);

            _logger.LogInformation("Product deleted: {ProductId}", request.ProductId);

            return Result<bool>.Success(true);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to delete product: {ProductId}", request.ProductId);
            return Result<bool>.Failure($"Failed to delete product: {ex.Message}");
        }
    }
}
