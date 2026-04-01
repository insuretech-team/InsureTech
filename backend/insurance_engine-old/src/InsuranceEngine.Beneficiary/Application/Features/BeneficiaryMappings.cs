using System;
using InsuranceEngine.Beneficiary.Application.DTOs;
using InsuranceEngine.Beneficiary.Domain.Entities;
using InsuranceEngine.SharedKernel.DTOs;

namespace InsuranceEngine.Beneficiary.Application.Features;

public static class BeneficiaryMappings
{
    public static BeneficiaryDto ToDto(this Domain.Entities.Beneficiary b)
    {
        return new BeneficiaryDto(
            b.Id,
            b.UserId,
            b.Type.ToString(),
            b.Code,
            b.Status,
            b.KycStatus,
            b.AuditInfo.ToDto(),
            b.KycCompletedAt,
            b.RiskScore,
            b.ReferralCode,
            b.ReferredBy,
            b.PartnerId,
            b.Individual?.ToDto(),
            b.Business?.ToDto()
        );
    }

    public static IndividualBeneficiaryDto ToDto(this IndividualBeneficiary i)
    {
        return new IndividualBeneficiaryDto(
            i.BeneficiaryId,
            i.FullName,
            i.DateOfBirth,
            i.Gender.ToString(),
            i.MaritalStatus.ToString(),
            i.ContactInfo.ToDto(),
            i.PermanentAddress.ToDto(),
            i.AuditInfo.ToDto(),
            i.FullNameBn,
            i.NidNumber,
            i.PassportNumber,
            i.BirthCertificateNumber,
            i.TinNumber,
            i.Occupation,
            i.PresentAddress.ToDto(),
            i.NomineeName,
            i.NomineeRelationship
        );
    }

    public static BusinessBeneficiaryDto ToDto(this BusinessBeneficiary b)
    {
        return new BusinessBeneficiaryDto(
            b.Id,
            b.BeneficiaryId,
            b.BusinessName,
            b.TradeLicenseNumber,
            b.TinNumber,
            b.BusinessType.ToString(),
            b.ContactInfo.ToDto(),
            b.RegisteredAddress.ToDto(),
            b.BusinessAddress.ToDto(),
            b.FocalPersonName,
            b.FocalPersonContact.ToDto(),
            b.AuditInfo.ToDto(),
            b.BusinessNameBn,
            b.TradeLicenseIssueDate,
            b.TradeLicenseExpiryDate,
            b.BinNumber,
            b.IndustrySector,
            b.EmployeeCount,
            b.IncorporationDate,
            b.FocalPersonDesignation,
            b.FocalPersonNid,
            b.RegistrationNumber,
            b.TaxId,
            b.TotalEmployeesCovered,
            b.ActivePoliciesCount,
            new MoneyDto(b.TotalPremiumAmount.Amount, b.TotalPremiumAmount.CurrencyCode),
            b.PendingActionsCount
        );
    }

    public static AuditInfoDto ToDto(this InsuranceEngine.SharedKernel.Domain.ValueObjects.AuditInfo a)
    {
        return new AuditInfoDto(
            a.CreatedAt,
            a.CreatedBy,
            a.UpdatedAt,
            a.UpdatedBy,
            a.DeletedAt,
            a.DeletedBy
        );
    }

    public static AddressDto ToDto(this InsuranceEngine.SharedKernel.Domain.ValueObjects.Address a)
    {
        return new AddressDto(
            a.AddressLine1,
            a.AddressLine2,
            a.City,
            a.District,
            a.Division,
            a.PostalCode,
            a.Country,
            a.Latitude,
            a.Longitude
        );
    }

    public static ContactInfoDto ToDto(this InsuranceEngine.SharedKernel.Domain.ValueObjects.ContactInfo c)
    {
        return new ContactInfoDto(
            c.MobileNumber,
            c.AlternateMobile,
            c.Email,
            c.Landline
        );
    }
}
