using System;
using System.Collections.Generic;
using InsuranceEngine.SharedKernel.DTOs;

namespace InsuranceEngine.Beneficiary.Application.DTOs;

public record BeneficiaryDto(
    Guid Id,
    Guid UserId,
    string Type,
    string Code,
    object Status,
    object KycStatus,
    DateTime? KycCompletedAt,
    string? RiskScore,
    string? ReferralCode,
    AuditInfoDto AuditInfo,
    IndividualBeneficiaryDto? Individual = null,
    BusinessBeneficiaryDto? Business = null
);

public record IndividualBeneficiaryDto(
    string FullName,
    string? FullNameBn,
    DateTime DateOfBirth,
    string Gender,
    string? NidNumber,
    string? PassportNumber,
    string? BirthCertificateNumber,
    string? TinNumber,
    string MaritalStatus,
    string? Occupation,
    ContactInfoDto ContactInfo,
    AddressDto PermanentAddress,
    AddressDto PresentAddress,
    string? NomineeName,
    string? NomineeRelationship,
    AuditInfoDto AuditInfo
);

public record BusinessBeneficiaryDto(
    string BusinessName,
    string? BusinessNameBn,
    string TradeLicenseNumber,
    string TinNumber,
    string? BinNumber,
    string BusinessType,
    string? IndustrySector,
    string FocalPersonName,
    ContactInfoDto FocalPersonContact,
    string? FocalPersonDesignation,
    string? FocalPersonNid,
    ContactInfoDto ContactInfo,
    AddressDto RegisteredAddress,
    AddressDto BusinessAddress,
    int ActivePoliciesCount,
    int PendingActionsCount,
    int TotalEmployeesCovered,
    MoneyDto TotalPremiumAmount,
    AuditInfoDto AuditInfo
);
