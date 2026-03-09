using System;
using MediatR;
using InsuranceEngine.Products.Domain.Enums;

namespace InsuranceEngine.Products.Application.Features.Commands.CreateProduct;

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
