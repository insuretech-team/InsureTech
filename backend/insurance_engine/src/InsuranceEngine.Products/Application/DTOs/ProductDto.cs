using System;
using System.Collections.Generic;
using InsuranceEngine.Products.Domain.Enums;

namespace InsuranceEngine.Products.Application.DTOs;

// --- Money DTO ---
public record MoneyDto(long Amount, string CurrencyCode = "BDT");

// --- Product DTOs ---
public record ProductDto(
    Guid Id,
    string ProductCode,
    string ProductName,
    string? ProductNameBn,
    string? Description,
    ProductCategory Category,
    ProductStatus Status,
    MoneyDto BasePremium,
    MoneyDto MinSumInsured,
    MoneyDto MaxSumInsured,
    int MinAge,
    int MaxAge,
    int MinTenureMonths,
    int MaxTenureMonths,
    List<string> Exclusions,
    List<RiderDto>? AvailableRiders,
    List<ProductPlanDto>? Plans,
    PricingConfigDto? PricingConfig,
    Guid CreatedBy,
    DateTime CreatedAt,
    DateTime UpdatedAt
);

public record ProductListDto(
    Guid Id,
    string ProductCode,
    string ProductName,
    ProductCategory Category,
    ProductStatus Status,
    MoneyDto BasePremium,
    MoneyDto MinSumInsured,
    MoneyDto MaxSumInsured
);

// --- Rider DTO ---
public record RiderDto(
    Guid Id,
    string RiderName,
    string? Description,
    MoneyDto PremiumAmount,
    MoneyDto CoverageAmount,
    bool IsMandatory
);

// --- ProductPlan DTO ---
public record ProductPlanDto(
    Guid Id,
    string PlanName,
    string? PlanDescription,
    MoneyDto PremiumAmount,
    MoneyDto MinSumInsured,
    MoneyDto MaxSumInsured,
    string? Attributes
);

// --- PricingConfig DTO ---
public record PricingConfigDto(
    Guid Id,
    List<PricingRuleDto> Rules,
    DateTime EffectiveFrom,
    DateTime? EffectiveTo
);

public record PricingRuleDto(
    Guid Id,
    string RuleName,
    RuleType Type,
    List<RuleConditionDto> Conditions,
    RuleActionDto Action
);

public record RuleConditionDto(string Field, string Operator, string Value);
public record RuleActionDto(ActionType Type, double Value);

// --- Premium Calculation DTOs ---
public record CalculatePremiumRequest(
    long SumInsuredAmount,
    int TenureMonths,
    List<Guid>? RiderIds,
    Dictionary<string, string> ApplicantData
);

public record CalculatePremiumResponse(
    MoneyDto BasePremium,
    MoneyDto RiderPremium,
    MoneyDto Vat,
    MoneyDto ServiceFee,
    MoneyDto TotalPremium,
    List<PremiumBreakdownDto> Breakdown
);

public record PremiumBreakdownDto(string Item, MoneyDto Amount, string Description);

// --- Paginated Response ---
public record PaginatedResponse<T>(
    List<T> Items,
    int TotalCount,
    int Page,
    int PageSize
);
