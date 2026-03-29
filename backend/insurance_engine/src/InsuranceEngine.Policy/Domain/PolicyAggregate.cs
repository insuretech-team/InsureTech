using InsuranceEngine.SharedKernel.Domain;
using Insuretech.Policy.Services.V1;
using Insuretech.Policy.Entity.V1;

namespace InsuranceEngine.Policy.Domain;

public sealed class PolicyAggregate : AggregateRoot<Guid>
{
    public string PolicyNumber { get; private set; } = string.Empty;
    public Guid ProductId { get; private set; }
    public Guid CustomerId { get; private set; }
    public Guid? PartnerId { get; private set; }
    public Guid? AgentId { get; private set; }
    public string InsuranceType { get; private set; } = string.Empty;
    public string ProductCode { get; private set; } = string.Empty;
    public string Status { get; private set; } = "DRAFT";
    public Money PremiumAmount { get; private set; } = default!;
    public Money SumInsured { get; private set; } = default!;
    public int TenureMonths { get; private set; }
    public DateTime StartDate { get; private set; }
    public DateTime EndDate { get; private set; }
    public DateTime CreatedAt { get; private set; }
    
    private readonly List<Insuretech.Policy.Entity.V1.Nominee> _nominees = new();
    public IReadOnlyCollection<Insuretech.Policy.Entity.V1.Nominee> Nominees => _nominees.AsReadOnly();

    private PolicyAggregate(Guid id, string policyNumber, Guid productId, string productCode, string insuranceType, Guid customerId, 
                         Money premium, Money sumInsured, int tenure, DateTime startDate)
    {
        Id = id;
        PolicyNumber = policyNumber;
        ProductId = productId;
        ProductCode = productCode;
        InsuranceType = insuranceType;
        CustomerId = customerId;
        PremiumAmount = premium;
        SumInsured = sumInsured;
        TenureMonths = tenure;
        StartDate = startDate;
        EndDate = startDate.AddMonths(tenure);
        Status = "DRAFT";
        CreatedAt = DateTime.UtcNow;
    }

    public static PolicyAggregate Create(Guid productId, string productCode, string insuranceType, Guid customerId, decimal premium, decimal sumInsured, int tenure, DateTime startDate, long sequenceNumber)
    {
        return new PolicyAggregate(
            Guid.NewGuid(),
            ValueObjects.PolicyNumber.Generate(productCode, sequenceNumber).Value,
            productId,
            productCode,
            insuranceType,
            customerId,
            Money.FromDecimal(premium),
            Money.FromDecimal(sumInsured),
            tenure,
            startDate
        );
    }

    public void AddNominees(IEnumerable<Insuretech.Policy.Entity.V1.Nominee>? nominees)
    {
        if (nominees != null)
        {
            _nominees.AddRange(nominees);
        }
    }

    public void Activate()
    {
        Status = "ACTIVE";
    }
}
