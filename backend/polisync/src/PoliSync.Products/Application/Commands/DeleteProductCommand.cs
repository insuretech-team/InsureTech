using MediatR;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Products.Application.Commands;

public record DeleteProductCommand : IRequest<Result<bool>>
{
    public string ProductId { get; init; } = string.Empty;
}
