using System;
using InsuranceEngine.Domain.Enums;

namespace InsuranceEngine.Application.DTOs;

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
    int MaxAge,
    Guid InsurerId
);
