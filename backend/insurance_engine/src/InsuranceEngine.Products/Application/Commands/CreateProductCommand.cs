using Insuretech.Products.Services.V1;
using MediatR;

namespace InsuranceEngine.Products.Application.Commands;

public sealed record CreateProductCommand(
    string ProductCode,
    string ProductName,
    string Category,
    string? Description,
    decimal BasePremium,
    decimal MinSumInsured,
    decimal MaxSumInsured,
    int MinTenureMonths,
    int MaxTenureMonths,
    int MinAge,
    int MaxAge,
    string? CreatedBy) : IRequest<CreateProductResponse>;
