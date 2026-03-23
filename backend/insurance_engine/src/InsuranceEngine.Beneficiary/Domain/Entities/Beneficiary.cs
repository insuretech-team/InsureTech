using System;
using System.Collections.Generic;
using InsuranceEngine.Beneficiary.Domain.Enums;
using InsuranceEngine.SharedKernel.Domain;
using InsuranceEngine.SharedKernel.Domain.ValueObjects;

namespace InsuranceEngine.Beneficiary.Domain.Entities;

public class Beneficiary : AggregateRoot<Guid>
{
    public Guid UserId { get; private set; }
    public BeneficiaryType Type { get; private set; }
    public string Code { get; private set; } = string.Empty; // BEN-XXXXXX
    
    // Status as object (JSONB)
    public BeneficiaryStatusInfo Status { get; private set; } = new();
    
    // KYC status as object (JSONB)
    public KYCStatusInfo KycStatus { get; private set; } = new();
    
    public DateTime? KycCompletedAt { get; private set; }
    public string? RiskScore { get; private set; } // LOW, MEDIUM, HIGH
    
    public Guid? PartnerId { get; set; }
    public string? ReferralCode { get; set; }
    public Guid? ReferredBy { get; set; }
    
    public AuditInfo AuditInfo { get; private set; } = new();
    
    public IndividualBeneficiary? Individual { get; private set; }
    public BusinessBeneficiary? Business { get; private set; }

    // EF Core constructor
    public Beneficiary() { }

    public Beneficiary(
        Guid id,
        Guid userId,
        BeneficiaryType type) : base(id)
    {
        UserId = userId;
        Type = type;
        Status = new BeneficiaryStatusInfo { Value = BeneficiaryStatus.PendingKyc.ToString() };
        KycStatus = new KYCStatusInfo { Status = KYCStatus.NotStarted.ToString() };
        AuditInfo = new AuditInfo();
    }

    public static Beneficiary CreateIndividual(
        Guid userId,
        string fullName,
        DateTime dob,
        BeneficiaryGender gender,
        string mobile,
        string? email)
    {
        var beneficiary = new Beneficiary(Guid.NewGuid(), userId, BeneficiaryType.Individual);
        beneficiary.Individual = new IndividualBeneficiary(Guid.NewGuid(), beneficiary.Id, fullName, dob, gender);
        beneficiary.UpdateCode();
        return beneficiary;
    }

    public static Beneficiary CreateBusiness(
        Guid userId,
        string businessName,
        string tradeLicense,
        string tin,
        string focalPersonName,
        string focalPersonMobile)
    {
        var beneficiary = new Beneficiary(Guid.NewGuid(), userId, BeneficiaryType.Business);
        beneficiary.Business = new BusinessBeneficiary(Guid.NewGuid(), beneficiary.Id, businessName, tradeLicense, tin);
        beneficiary.Business.UpdateFocalPerson(focalPersonName, focalPersonMobile);
        beneficiary.UpdateCode();
        return beneficiary;
    }

    public static Beneficiary CreateIndividualEmpty(Guid userId, Guid partnerId)
    {
        var beneficiary = new Beneficiary(Guid.NewGuid(), userId, BeneficiaryType.Individual);
        beneficiary.Individual = IndividualBeneficiary.Create(Guid.NewGuid(), beneficiary.Id);
        beneficiary.Individual.ContactInfo = new SharedKernel.Domain.ValueObjects.ContactInfo();
        beneficiary.Individual.PermanentAddress = new SharedKernel.Domain.ValueObjects.Address();
        beneficiary.PartnerId = partnerId;
        beneficiary.UpdateCode();
        return beneficiary;
    }

    public static Beneficiary CreateBusinessEmpty(Guid userId, Guid partnerId)
    {
        var beneficiary = new Beneficiary(Guid.NewGuid(), userId, BeneficiaryType.Business);
        beneficiary.Business = BusinessBeneficiary.Create(Guid.NewGuid(), beneficiary.Id);
        beneficiary.Business.ContactInfo = new SharedKernel.Domain.ValueObjects.ContactInfo();
        beneficiary.Business.RegisteredAddress = new SharedKernel.Domain.ValueObjects.Address();
        beneficiary.Business.BusinessAddress = new SharedKernel.Domain.ValueObjects.Address();
        beneficiary.Business.FocalPersonContact = new SharedKernel.Domain.ValueObjects.ContactInfo();
        beneficiary.Business.PrimaryContact = new SharedKernel.Domain.ValueObjects.PrimaryContact();
        beneficiary.PartnerId = partnerId;
        beneficiary.UpdateCode();
        return beneficiary;
    }

    public void CompleteKYC(KYCStatus status)
    {
        KycStatus = new KYCStatusInfo { Status = status.ToString() };
        if (status == KYCStatus.Verified)
        {
            Status = new BeneficiaryStatusInfo { Value = BeneficiaryStatus.Active.ToString() };
            KycCompletedAt = DateTime.UtcNow;
        }
        AuditInfo = AuditInfo with { UpdatedAt = DateTime.UtcNow };
    }

    public void UpdateRiskScore(string score)
    {
        RiskScore = score;
        AuditInfo = AuditInfo with { UpdatedAt = DateTime.UtcNow };
    }

    public void UpdateCode()
    {
        Code = $"BEN-{Id.ToString().Substring(0, 8).ToUpper()}";
    }
}

public class BeneficiaryStatusInfo
{
    public string Value { get; set; } = string.Empty;
    public string? Description { get; set; }
}

public class KYCStatusInfo
{
    public string Status { get; set; } = string.Empty;
    public string? Remarks { get; set; }
}
