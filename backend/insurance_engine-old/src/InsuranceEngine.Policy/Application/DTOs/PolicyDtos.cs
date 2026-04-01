using System;
using System.Collections.Generic;

namespace InsuranceEngine.Policy.Application.DTOs;

public record MoneyDto(long Amount, string Currency = "BDT")
{
    public decimal DecimalAmount => Amount / 100m;
}

public record PolicyDto(
    Guid Id,
    string PolicyNumber,
    Guid ProductId,
    Guid CustomerId,
    Guid? PartnerId,
    Guid? AgentId,
    InsuranceEngine.Policy.Domain.Enums.PolicyStatus Status,
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
    InsuranceEngine.Policy.Domain.Enums.PolicyStatus Status,
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
    MoneyDto AnnualIncome,
    string? Address,
    string? PhoneNumber,
    HealthDeclarationDto? HealthDeclaration
);

public record HealthDeclarationDto(
    bool HasPreExistingConditions,
    List<string>? Conditions,
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
    Guid? Id,
    string RiderName,
    MoneyDto PremiumAmount,
    MoneyDto CoverageAmount
);

public record GracePeriodDto(
    Guid PolicyId,
    InsuranceEngine.Policy.Domain.Enums.PolicyStatus Status,
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

public record CreatePolicyResponse(Guid PolicyId, string PolicyNumber);
