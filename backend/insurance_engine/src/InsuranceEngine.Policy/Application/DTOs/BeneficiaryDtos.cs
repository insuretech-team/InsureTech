using System;

namespace InsuranceEngine.Policy.Application.DTOs;

public record BeneficiaryDto(
    Guid Id,
    Guid UserId,
    string Type,
    string Code,
    string Status,
    string KycStatus,
    DateTime? KycCompletedAt,
    string? RiskScore,
    string? ReferralCode,
    IndividualBeneficiaryDto? IndividualDetails = null,
    BusinessBeneficiaryDto? BusinessDetails = null
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
    string? ContactInfoJson,
    string? PermanentAddressJson,
    string? PresentAddressJson,
    string? NomineeName,
    string? NomineeRelationship
);

public record BusinessBeneficiaryDto(
    string BusinessName,
    string? BusinessNameBn,
    string TradeLicenseNumber,
    DateTime? TradeLicenseIssueDate,
    DateTime? TradeLicenseExpiryDate,
    string TinNumber,
    string? BinNumber,
    string BusinessType,
    string? IndustrySector,
    int EmployeeCount,
    DateTime? IncorporationDate,
    string? ContactInfoJson,
    string? RegisteredAddressJson,
    string? BusinessAddressJson,
    string FocalPersonName,
    string? FocalPersonDesignation,
    string? FocalPersonNid,
    string? FocalPersonContactJson,
    string? RegistrationNumber,
    string? TaxId,
    int TotalEmployeesCovered,
    int ActivePoliciesCount,
    long TotalPremiumAmount,
    int PendingActionsCount
);

public record CreateIndividualBeneficiaryRequest(
    Guid UserId,
    string FullName,
    DateTime DateOfBirth,
    string Gender,
    string NidNumber,
    string MobileNumber,
    string? Email = null,
    Guid? PartnerId = null
);

public record CreateBusinessBeneficiaryRequest(
    Guid UserId,
    string BusinessName,
    string TradeLicenseNumber,
    string TinNumber,
    string FocalPersonName,
    string FocalPersonMobile,
    Guid? PartnerId = null
);
