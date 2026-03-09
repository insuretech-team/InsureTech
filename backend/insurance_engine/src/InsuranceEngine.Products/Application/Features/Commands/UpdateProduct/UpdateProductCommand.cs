using System;
using MediatR;
using InsuranceEngine.Products.Domain.Enums;

namespace InsuranceEngine.Products.Application.Features.Commands.UpdateProduct;

public record UpdateProductCommand(
    Guid Id,
    string ProductName,
    string? ProductNameBn,
    string? Description,
    ProductCategory Category,
    decimal MinSumInsured,
    decimal MaxSumInsured,
    int MinAge,
    int MaxAge
) : IRequest<bool>;
