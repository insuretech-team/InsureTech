using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;
using InsuranceEngine.SharedKernel.DTOs;

namespace InsuranceEngine.Products.Application.DTOs;

public record ProductCreationResponse(
    [property: JsonPropertyName("product_id")] Guid ProductId,
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record ProductRetrievalResponse(
    [property: JsonPropertyName("product")] ProductDto Product,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record ProductsListingResponse(
    [property: JsonPropertyName("products")] List<ProductListDto> Products,
    [property: JsonPropertyName("total_count")] int TotalCount,
    [property: JsonPropertyName("page")] int Page,
    [property: JsonPropertyName("page_size")] int PageSize,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record ProductUpdateResponse(
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record PremiumCalculationResponse(
    [property: JsonPropertyName("base_premium")] MoneyDto BasePremium,
    [property: JsonPropertyName("rider_premium")] MoneyDto RiderPremium,
    [property: JsonPropertyName("total_premium")] MoneyDto TotalPremium,
    [property: JsonPropertyName("breakdown")] List<PremiumBreakdownDto> Breakdown,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);
