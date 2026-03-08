using Insuretech.Common.V1;
using Insuretech.Products.Entity.V1;
using MediatR;
using PoliSync.SharedKernel.CQRS;
using System.Collections.Generic;

namespace PoliSync.Products.Application.Commands;

public record CreateProductCommand : IRequest<Result<Product>>
{
    public string ProductCode { get; init; } = string.Empty;
    public string ProductName { get; init; } = string.Empty;
    public ProductCategory Category { get; init; }
    public string Description { get; init; } = string.Empty;
    public long BasePremiumAmount { get; init; } // In paisa
    public long MinSumInsuredAmount { get; init; }
    public long MaxSumInsuredAmount { get; init; }
    public int MinTenureMonths { get; init; }
    public int MaxTenureMonths { get; init; }
    public List<string> Exclusions { get; init; } = new();
    public string CreatedBy { get; init; } = string.Empty;
}
