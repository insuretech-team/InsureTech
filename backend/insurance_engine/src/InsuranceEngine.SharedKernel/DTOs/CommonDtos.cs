using System.Collections.Generic;

namespace InsuranceEngine.SharedKernel.DTOs;

public record MoneyDto(long Amount, string CurrencyCode = "BDT");

public record PaginatedResponse<T>(
    List<T> Items,
    int TotalCount,
    int Page,
    int PageSize
);
