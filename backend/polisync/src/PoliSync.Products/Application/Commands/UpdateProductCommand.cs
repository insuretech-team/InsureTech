using Insuretech.Products.Entity.V1;
using MediatR;
using PoliSync.SharedKernel.CQRS;
using System.Collections.Generic;

namespace PoliSync.Products.Application.Commands;

public record UpdateProductCommand : IRequest<Result<Product>>
{
    public string ProductId { get; init; } = string.Empty;
    public string? ProductName { get; init; }
    public string? Description { get; init; }
    public long? BasePremiumAmount { get; init; }
    public long? MinSumInsuredAmount { get; init; }
    public long? MaxSumInsuredAmount { get; init; }
    public int? MinTenureMonths { get; init; }
    public int? MaxTenureMonths { get; init; }
    public List<string>? Exclusions { get; init; }
    public ProductStatus? Status { get; init; }
}
