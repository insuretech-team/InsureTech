using System;

namespace InsuranceEngine.Beneficiary.Application.DTOs;

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
    string TinNumber,
    string? BinNumber,
    string BusinessType,
    string? IndustrySector,
    string FocalPersonName,
    string? FocalPersonMobile,
    string? ContactInfoJson,
    string? RegisteredAddressJson,
    string? BusinessAddressJson
);

public record PaginatedResponse<T>(
    List<T> Items,
    int TotalCount,
    int Page,
    int PageSize
);
