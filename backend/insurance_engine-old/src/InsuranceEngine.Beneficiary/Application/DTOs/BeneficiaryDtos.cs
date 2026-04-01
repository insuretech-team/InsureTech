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
    AuditInfoDto AuditInfo,
    DateTime? KycCompletedAt = null,
    string? RiskScore = null,
    string? ReferralCode = null,
    Guid? ReferredBy = null,
    Guid? PartnerId = null,
    IndividualBeneficiaryDto? Individual = null,
    BusinessBeneficiaryDto? Business = null
);

public record IndividualBeneficiaryDto(
    Guid BeneficiaryId,
    string FullName,
    DateTime DateOfBirth,
    string Gender,
    string MaritalStatus,
    ContactInfoDto ContactInfo,
    AddressDto PermanentAddress,
    AuditInfoDto AuditInfo,
    string? FullNameBn = null,
    string? NidNumber = null,
    string? PassportNumber = null,
    string? BirthCertificateNumber = null,
    string? TinNumber = null,
    string? Occupation = null,
    AddressDto? PresentAddress = null,
    string? NomineeName = null,
    string? NomineeRelationship = null
);

public record BusinessBeneficiaryDto(
    Guid Id,
    Guid BeneficiaryId,
    string BusinessName,
    string TradeLicenseNumber,
    string TinNumber,
    string BusinessType,
    ContactInfoDto ContactInfo,
    AddressDto RegisteredAddress,
    AddressDto BusinessAddress,
    string FocalPersonName,
    ContactInfoDto FocalPersonContact,
    AuditInfoDto AuditInfo,
    string? BusinessNameBn = null,
    DateTime? TradeLicenseIssueDate = null,
    DateTime? TradeLicenseExpiryDate = null,
    string? BinNumber = null,
    string? IndustrySector = null,
    int? EmployeeCount = null,
    DateTime? IncorporationDate = null,
    string? FocalPersonDesignation = null,
    string? FocalPersonNid = null,
    string? RegistrationNumber = null,
    string? TaxId = null,
    int? TotalEmployeesCovered = null,
    int? ActivePoliciesCount = null,
    MoneyDto? TotalPremiumAmount = null,
    int? PendingActionsCount = null
);
