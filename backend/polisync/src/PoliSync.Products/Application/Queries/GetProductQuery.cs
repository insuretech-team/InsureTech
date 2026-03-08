using Insuretech.Products.Entity.V1;
using MediatR;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Products.Application.Queries;

public record GetProductQuery : IRequest<Result<Product?>>
{
    public string ProductId { get; init; } = string.Empty;
}
