using PoliSync.SharedKernel.Domain;

namespace PoliSync.Beneficiaries.Domain;

public class Beneficiary : Entity
{
    public Guid BeneficiaryId { get; private set; }
    public Guid UserId { get; private set; } // Link to authn_schema.users
    public BeneficiaryType Type { get; private set; }
    public string Code { get; private set; } = string.Empty; // BEN-XXXXXX
    public BeneficiaryStatus Status { get; private set; } = BeneficiaryStatus.PendingKyc;
    public KycStatus KycStatus { get; private set; } = KycStatus.NotStarted;
    public DateTime? KycCompletedAt { get; private set; }
    public string? RiskScore { get; private set; } // LOW, MEDIUM, HIGH
    public string? ReferralCode { get; private set; }
    public Guid? ReferredBy { get; private set; }
    public Guid? PartnerId { get; private set; }
    public string? AuditInfo { get; private set; } // JSONB

    // Navigation properties for details
    public IndividualBeneficiary? IndividualDetails { get; private set; }
    public BusinessBeneficiary? BusinessDetails { get; private set; }

    private Beneficiary() { }

    public static Beneficiary Create(
        Guid userId,
        BeneficiaryType type,
        string code,
        Guid? partnerId = null,
        string? referralCode = null,
        Guid? referredBy = null)
    {
        return new Beneficiary
        {
            BeneficiaryId = Guid.NewGuid(),
            UserId = userId,
            Type = type,
            Code = code,
            Status = BeneficiaryStatus.PendingKyc,
            KycStatus = KycStatus.NotStarted,
            PartnerId = partnerId,
            ReferralCode = referralCode,
            ReferredBy = referredBy,
            AuditInfo = "{}"
        };
    }

    public void UpdateStatus(BeneficiaryStatus status)
    {
        Status = status;
    }

    public void CompleteKyc(KycStatus status)
    {
        KycStatus = status;
        if (status == KycStatus.Completed)
        {
            KycCompletedAt = DateTime.UtcNow;
            Status = BeneficiaryStatus.Active;
        }
    }

    public void UpdateRiskScore(string riskScore)
    {
        RiskScore = riskScore;
    }
}
