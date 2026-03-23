using System;
using System.Collections.Generic;
using InsuranceEngine.Beneficiary.Domain.Enums;
using InsuranceEngine.SharedKernel.Domain;

namespace InsuranceEngine.Beneficiary.Domain.Entities;

public class Beneficiary : AggregateRoot<Guid>
{
    public Guid UserId { get; private set; }
    public BeneficiaryType Type { get; private set; }
    public string Code { get; private set; } = string.Empty; // BEN-XXXXXX
    public BeneficiaryStatus Status { get; private set; }
    public KYCStatus KycStatus { get; private set; }
    public DateTime? KycCompletedAt { get; private set; }
    public string? RiskScore { get; private set; } // LOW, MEDIUM, HIGH
    
    public Guid? PartnerId { get; private set; }
    public string? ReferralCode { get; private set; }
    public Guid? ReferredBy { get; private set; }
    public string? AuditInfoJson { get; set; }

    public IndividualBeneficiary? IndividualDetails { get; private set; }
    public BusinessBeneficiary? BusinessDetails { get; private set; }

    public DateTime CreatedAt { get; private set; }
    public DateTime UpdatedAt { get; set; }

    // EF Core constructor
    public Beneficiary() { }

    private Beneficiary(
        Guid id,
        Guid userId,
        BeneficiaryType type) : base(id)
    {
        UserId = userId;
        Type = type;
        Status = BeneficiaryStatus.PendingKyc;
        KycStatus = KYCStatus.NotStarted;
        CreatedAt = DateTime.UtcNow;
        UpdatedAt = DateTime.UtcNow;
    }

    public static Beneficiary CreateIndividual(
        Guid userId,
        string fullName,
        DateTime dob,
        Gender gender,
        string mobile,
        string? email)
    {
        var beneficiary = new Beneficiary(Guid.NewGuid(), userId, BeneficiaryType.Individual);
        beneficiary.IndividualDetails = new IndividualBeneficiary(Guid.NewGuid(), beneficiary.Id, fullName, dob, gender);
        // Contact info will be set in detail entity
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
        beneficiary.BusinessDetails = new BusinessBeneficiary(Guid.NewGuid(), beneficiary.Id, businessName, tradeLicense, tin);
        beneficiary.BusinessDetails.UpdateFocalPerson(focalPersonName, focalPersonMobile);
        beneficiary.UpdateCode();
        return beneficiary;
    }

    public void CompleteKYC(KYCStatus status)
    {
        KycStatus = status;
        if (status == KYCStatus.Verified)
        {
            Status = BeneficiaryStatus.Active;
            KycCompletedAt = DateTime.UtcNow;
        }
        UpdatedAt = DateTime.UtcNow;
    }

    public void UpdateRiskScore(string score)
    {
        RiskScore = score;
        UpdatedAt = DateTime.UtcNow;
    }

    private void UpdateCode()
    {
        Code = $"BEN-{Id.ToString().Substring(0, 8).ToUpper()}";
    }
}
