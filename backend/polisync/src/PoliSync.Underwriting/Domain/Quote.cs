using PoliSync.SharedKernel.Domain;

namespace PoliSync.Underwriting.Domain;

public class Quote : Entity
{
    public Guid QuoteId { get; private set; }
    public string QuoteNumber { get; private set; } = string.Empty;
    public Guid BeneficiaryId { get; private set; }
    public Guid InsurerProductId { get; private set; } // FK to Product
    public QuoteStatus Status { get; private set; } = QuoteStatus.Draft;
    
    // Coverage Details
    public long SumAssured { get; private set; }
    public string Currency { get; private set; } = "BDT";
    public int TermYears { get; private set; }
    public string PremiumPaymentMode { get; private set; } = "YEARLY"; // YEARLY, HALF_YEARLY, MONTHLY
    
    // Calculated Premiums
    public long BasePremiumAmount { get; private set; }
    public long RiderPremiumAmount { get; private set; }
    public long TaxAmount { get; private set; }
    public long TotalPremiumAmount { get; private set; }
    
    // Underwriting Inputs
    public int ApplicantAgeDays { get; private set; }
    public bool IsSmoker { get; private set; }
    
    // Timestamps
    public DateTime CreatedAt { get; private set; }
    public DateTime UpdatedAt { get; private set; }
    public DateTime ValidUntil { get; private set; }

    // Navigation properties
    public HealthDeclaration? HealthDeclaration { get; private set; }
    public UnderwritingDecision? Decision { get; private set; }

    private Quote() { }

    public static Quote Create(
        Guid beneficiaryId,
        Guid insurerProductId,
        long sumAssured,
        int termYears,
        string premiumPaymentMode,
        long basePremiumAmount,
        long riderPremiumAmount,
        long taxAmount,
        int applicantAgeDays,
        bool isSmoker)
    {
        var quote = new Quote
        {
            QuoteId = Guid.NewGuid(),
            QuoteNumber = $"QUT-{Guid.NewGuid().ToString()[..8].ToUpper()}",
            BeneficiaryId = beneficiaryId,
            InsurerProductId = insurerProductId,
            Status = QuoteStatus.PendingUnderwriting,
            SumAssured = sumAssured,
            TermYears = termYears,
            PremiumPaymentMode = premiumPaymentMode,
            BasePremiumAmount = basePremiumAmount,
            RiderPremiumAmount = riderPremiumAmount,
            TaxAmount = taxAmount,
            TotalPremiumAmount = basePremiumAmount + riderPremiumAmount + taxAmount,
            ApplicantAgeDays = applicantAgeDays,
            IsSmoker = isSmoker,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow,
            ValidUntil = DateTime.UtcNow.AddDays(30) // Quotes valid for 30 days
        };

        quote.RaiseDomainEvent(new QuoteRequestedEvent(quote.QuoteId, quote.QuoteNumber));
        return quote;
    }

    public void UpdateStatus(QuoteStatus status)
    {
        Status = status;
        UpdatedAt = DateTime.UtcNow;
    }

    public void Approve()
    {
        if (Status != QuoteStatus.PendingUnderwriting)
            throw new InvalidOperationException("Can only approve quotes in pending state");
            
        Status = QuoteStatus.Approved;
        UpdatedAt = DateTime.UtcNow;
        RaiseDomainEvent(new QuoteApprovedEvent(QuoteId));
    }

    public void Reject()
    {
        if (Status != QuoteStatus.PendingUnderwriting)
            throw new InvalidOperationException("Can only reject quotes in pending state");
            
        Status = QuoteStatus.Rejected;
        UpdatedAt = DateTime.UtcNow;
        RaiseDomainEvent(new QuoteRejectedEvent(QuoteId));
    }

    public void ConvertToPolicy()
    {
        if (Status != QuoteStatus.Approved)
            throw new InvalidOperationException("Only approved quotes can be converted to policies");
            
        Status = QuoteStatus.ConvertedToPolicy;
        UpdatedAt = DateTime.UtcNow;
    }
}
