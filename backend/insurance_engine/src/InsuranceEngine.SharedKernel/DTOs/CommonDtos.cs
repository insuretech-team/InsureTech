using System;
using System.Collections.Generic;

namespace InsuranceEngine.SharedKernel.DTOs;

public record MoneyDto(long Amount, string Currency = "BDT")
{
    public decimal DecimalAmount => Amount / 100m;
}

public record AddressDto(
    string? AddressLine1,
    string? AddressLine2,
    string? City,
    string? District,
    string? Division,
    string? PostalCode,
    string Country = "Bangladesh",
    double? Latitude = null,
    double? Longitude = null
);

public record ContactInfoDto(
    string? MobileNumber,
    string? AlternateMobile,
    string? Email,
    string? Landline
);

public record AuditInfoDto(
    DateTime CreatedAt,
    string? CreatedBy,
    DateTime UpdatedAt,
    string? UpdatedBy,
    DateTime? DeletedAt = null,
    string? DeletedBy = null
);

public record PrimaryContactDto(
    string? Name,
    string? Email,
    string? Phone,
    string? Department
);

public record PaginatedResponse<T>(
    List<T> Items,
    int TotalCount,
    int Page,
    int PageSize
);
