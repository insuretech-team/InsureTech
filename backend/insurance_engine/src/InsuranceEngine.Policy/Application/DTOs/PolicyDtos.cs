using System;
using System.Collections.Generic;
using InsuranceEngine.Policy.Domain.Enums;

namespace InsuranceEngine.Policy.Application.DTOs;

public record MoneyDto(long Amount, string CurrencyCode = "BDT");

public record PolicyDto(
    Guid Id,
    string PolicyNumber,
    Guid ProductId,
    Guid CustomerId,
    Guid? PartnerId,
    Guid? AgentId,
    PolicyStatus Status,
    MoneyDto PremiumAmount,
    MoneyDto SumInsured,
    MoneyDto? VatTax,
    MoneyDto? ServiceFee,
    MoneyDto? TotalPayable,
    int TenureMonths,
    DateTime StartDate,
    DateTime EndDate,
    DateTime? IssuedAt,
    string? PaymentFrequency,
    string? ProviderName,
    ApplicantDto? ProposerDetails,
    List<NomineeDto>? Nominees,
    List<PolicyRiderDto>? Riders,
    DateTime CreatedAt,
    DateTime UpdatedAt
);

public record PolicyListDto(
    Guid Id,
    string PolicyNumber,
    Guid ProductId,
    Guid CustomerId,
    PolicyStatus Status,
    MoneyDto PremiumAmount,
    MoneyDto SumInsured,
    DateTime StartDate,
    DateTime EndDate,
    DateTime? IssuedAt
);

public record ApplicantDto(
    string FullName,
    DateTime? DateOfBirth,
    string? NidNumber,
    string? Occupation,
    long AnnualIncome,
    string? Address,
    string? PhoneNumber,
    HealthDeclarationDto? HealthDeclaration
);

public record HealthDeclarationDto(
    bool HasPreExistingConditions,
    List<string> Conditions,
    bool IsSmoker,
    string? BloodGroup
);

public record NomineeDto(
    Guid? Id,
    Guid? BeneficiaryId,
    string FullName,
    string Relationship,
    double SharePercentage,
    DateTime? DateOfBirth,
    string? NomineeDobText,
    string? NidNumber,
    string? PhoneNumber
);

public record PolicyRiderDto(
    Guid Id,
    string RiderName,
    MoneyDto PremiumAmount,
    MoneyDto CoverageAmount
);

public record GracePeriodDto(
    Guid PolicyId,
    PolicyStatus Status,
    DateTime EndDate,
    DateTime GracePeriodEndDate,
    int DaysRemaining,
    bool IsInGracePeriod
);

public record RenewalScheduleDto(
    Guid PolicyId,
    string PolicyNumber,
    DateTime CurrentEndDate,
    DateTime NextRenewalDate,
    MoneyDto EstimatedPremium,
    bool IsEligibleForRenewal
);

public record PaginatedResponse<T>(
    List<T> Items,
    int TotalCount,
    int Page,
    int PageSize
);
