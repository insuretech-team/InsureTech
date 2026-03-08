using System;
using MediatR;
using InsuranceEngine.Domain.Enums;

namespace InsuranceEngine.Application.Features.Products.Commands.UpdateProduct;

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
