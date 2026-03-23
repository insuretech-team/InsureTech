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
            b.Status, // BeneficiaryStatusInfo
            b.KycStatus, // KYCStatusInfo
            b.KycCompletedAt,
            b.RiskScore,
            b.ReferralCode,
            b.AuditInfo.ToDto(),
            b.Individual?.ToDto(),
            b.Business?.ToDto()
        );
    }

    public static IndividualBeneficiaryDto ToDto(this IndividualBeneficiary i)
    {
        return new IndividualBeneficiaryDto(
            i.FullName,
            i.FullNameBn,
            i.DateOfBirth,
            i.Gender.ToString(),
            i.NidNumber,
            i.PassportNumber,
            i.BirthCertificateNumber,
            i.TinNumber,
            i.MaritalStatus.ToString(),
            i.Occupation,
            i.ContactInfo.ToDto(),
            i.PermanentAddress.ToDto(),
            i.PresentAddress.ToDto(),
            i.NomineeName,
            i.NomineeRelationship,
            i.AuditInfo.ToDto()
        );
    }

    public static BusinessBeneficiaryDto ToDto(this BusinessBeneficiary b)
    {
        return new BusinessBeneficiaryDto(
            b.BusinessName,
            b.BusinessNameBn,
            b.TradeLicenseNumber,
            b.TinNumber,
            b.BinNumber,
            b.BusinessType.ToString(),
            b.IndustrySector,
            b.FocalPersonName,
            b.FocalPersonContact.ToDto(),
            b.FocalPersonDesignation,
            b.FocalPersonNid,
            b.ContactInfo.ToDto(),
            b.RegisteredAddress.ToDto(),
            b.BusinessAddress.ToDto(),
            b.ActivePoliciesCount,
            b.PendingActionsCount,
            b.TotalEmployeesCovered,
            new MoneyDto(b.TotalPremiumAmount.Amount, b.TotalPremiumAmount.CurrencyCode),
            b.AuditInfo.ToDto()
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
