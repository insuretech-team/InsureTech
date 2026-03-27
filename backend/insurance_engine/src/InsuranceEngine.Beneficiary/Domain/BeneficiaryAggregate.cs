using InsuranceEngine.SharedKernel.Domain;

namespace InsuranceEngine.Beneficiary.Domain;

public sealed class BeneficiaryAggregate : AggregateRoot<Guid>
{
    public Guid? UserId { get; private set; }
    public string Type { get; private set; } = string.Empty;
    public string Code { get; private set; } = string.Empty;
    public string Status { get; private set; } = string.Empty;
    public string KycStatus { get; private set; } = string.Empty;
    public Guid? PartnerId { get; private set; }
    public DateTime CreatedAt { get; private set; }

    private BeneficiaryAggregate(Guid id, Guid? userId, string type, string code, Guid? partnerId)
    {
        Id = id;
        UserId = userId;
        Type = type;
        Code = code;
        PartnerId = partnerId;
        Status = "PENDING_KYC";
        KycStatus = "NOT_STARTED";
        CreatedAt = DateTime.UtcNow;
    }

    public static BeneficiaryAggregate CreateIndividual(Guid? userId, string fullName, Guid? partnerId = null)
    {
        var code = $"BEN-I-{Guid.NewGuid().ToString()[..8].ToUpper()}";
        return new BeneficiaryAggregate(Guid.NewGuid(), userId, "INDIVIDUAL", code, partnerId);
    }

    public static BeneficiaryAggregate CreateBusiness(Guid userId, string businessName, Guid? partnerId)
    {
        var code = $"BEN-B-{Guid.NewGuid().ToString()[..8].ToUpper()}";
        return new BeneficiaryAggregate(Guid.NewGuid(), userId, "BUSINESS", code, partnerId);
    }
}
