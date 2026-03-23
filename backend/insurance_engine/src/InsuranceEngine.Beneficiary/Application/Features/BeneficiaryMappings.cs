using System;
using InsuranceEngine.Beneficiary.Application.DTOs;
using InsuranceEngine.Beneficiary.Domain.Entities;

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
            b.Status.ToString(),
            b.KycStatus.ToString(),
            b.KycCompletedAt,
            b.RiskScore,
            null, // ReferralCode — can be added to entity later if needed
            b.IndividualDetails != null ? new IndividualBeneficiaryDto(
                b.IndividualDetails.FullName,
                b.IndividualDetails.FullNameBn,
                b.IndividualDetails.DateOfBirth,
                b.IndividualDetails.Gender.ToString(),
                b.IndividualDetails.NidNumber,
                b.IndividualDetails.PassportNumber,
                b.IndividualDetails.BirthCertificateNumber,
                b.IndividualDetails.TinNumber,
                b.IndividualDetails.MaritalStatus.ToString(),
                b.IndividualDetails.Occupation,
                b.IndividualDetails.ContactInfoJson,
                b.IndividualDetails.PermanentAddressJson,
                b.IndividualDetails.PresentAddressJson,
                null,
                null
            ) : null,
            b.BusinessDetails != null ? new BusinessBeneficiaryDto(
                b.BusinessDetails.BusinessName,
                b.BusinessDetails.BusinessNameBn,
                b.BusinessDetails.TradeLicenseNumber,
                b.BusinessDetails.TinNumber,
                null,
                b.BusinessDetails.BusinessType.ToString(),
                b.BusinessDetails.IndustrySector,
                b.BusinessDetails.FocalPersonName,
                b.BusinessDetails.FocalPersonMobile,
                null,
                null,
                null
            ) : null
        );
    }
}
