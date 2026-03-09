using System;
using InsuranceEngine.Products.Domain.Enums;

namespace InsuranceEngine.Products.Application.DTOs;

public record ProductDto(
    Guid Id,
    string ProductCode,
    string ProductName,
    string? ProductNameBn,
    string? Description,
    ProductCategory Category,
    ProductStatus Status,
    decimal MinSumInsured,
    decimal MaxSumInsured,
    int MinAge,
    int MaxAge
);
