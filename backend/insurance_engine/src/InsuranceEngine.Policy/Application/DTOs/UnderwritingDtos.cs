using System;
using System.Collections.Generic;
using InsuranceEngine.Policy.Domain.Enums;

namespace InsuranceEngine.Policy.Application.DTOs;

public record QuoteDto(
    Guid Id,
    string QuoteNumber,
    Guid BeneficiaryId,
    Guid InsurerProductId,
    QuoteStatus Status,
    MoneyDto SumAssured,
    int TermYears,
    string PremiumPaymentMode,
    MoneyDto BasePremium,
    MoneyDto RiderPremium,
    MoneyDto TotalPremium,
    int ApplicantAge,
    string? ApplicantOccupation,
    bool IsSmoker,
    DateTime ValidUntil,
    DateTime CreatedAt
);

public record UnderwritingHealthDeclarationDto(
    int HeightCm,
    decimal WeightKg,
    decimal Bmi,
    bool HasPreExistingConditions,
    List<string>? PreExistingConditions,
    bool IsCurrentlyHospitalized,
    bool HasFamilyHistory,
    List<string>? FamilyHistory,
    bool IsSmoker,
    bool IsAlcoholConsumer,
    string? OccupationRiskLevel,
    bool IsMedicalExamRequired,
    bool IsMedicalExamCompleted,
    DateTime? MedicalExamDate,
    List<string>? MedicalDocuments
);

public record UnderwritingDecisionDto(
    Guid Id,
    Guid QuoteId,
    DecisionType Decision,
    DecisionMethod Method,
    decimal RiskScore,
    RiskLevel RiskLevel,
    List<string>? RiskFactors,
    string? Reason,
    List<string>? Conditions,
    bool IsPremiumAdjusted,
    MoneyDto? AdjustedPremium,
    string? AdjustmentReason,
    Guid? UnderwriterId,
    string? UnderwriterComments,
    DateTime DecidedAt,
    DateTime? ValidUntil
);

public record UnderwritingHealthDeclarationResponseDto(
    Guid Id,
    Guid QuoteId,
    int HeightCm,
    decimal WeightKg,
    decimal Bmi,
    bool HasPreExistingConditions,
    List<string>? PreExistingConditions,
    bool IsCurrentlyHospitalized,
    bool HasFamilyHistory,
    List<string>? FamilyHistory,
    bool IsSmoker,
    bool IsAlcoholConsumer,
    string? OccupationRiskLevel,
    bool IsMedicalExamRequired,
    bool IsMedicalExamCompleted,
    DateTime? MedicalExamDate,
    DateTime CreatedAt,
    DateTime UpdatedAt
);
