using System;
using MediatR;
using InsuranceEngine.Domain.Enums;

namespace InsuranceEngine.Application.Features.Products.Commands.CreateProduct;

public record CreateProductCommand(
    string ProductCode,
    string ProductName,
    string? ProductNameBn,
    string? Description,
    ProductCategory Category,
    decimal MinSumInsured,
    decimal MaxSumInsured,
    int MinAge,
    int MaxAge,
    Guid InsurerId
) : IRequest<Guid>;
