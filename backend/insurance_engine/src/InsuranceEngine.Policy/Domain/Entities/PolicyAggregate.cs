using System;
using System.Collections.Generic;
using System.Linq;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.Policy.Domain.Events;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Domain;
using InsuranceEngine.SharedKernel.Domain.ValueObjects;

namespace InsuranceEngine.Policy.Domain.Entities;

/// <summary>
/// Policy aggregate root. Maps to 'policies' table in insurance_schema.
/// Enforces lifecycle state machine and nominee share invariant.
/// </summary>
public class PolicyAggregate : AggregateRoot<Guid>
{
    public string PolicyNumber { get; private set; } = string.Empty;
    public Guid ProductId { get; private set; }
    public Guid CustomerId { get; private set; }
    public Guid? PartnerId { get; private set; }
    public Guid? AgentId { get; private set; }
    public Guid? QuoteId { get; private set; }
    public Guid? UnderwritingDecisionId { get; private set; }

    public PolicyStatus Status { get; private set; }

    // Money fields — stored as bigint (paisa)
    public long PremiumAmount { get; private set; }
    public string PremiumCurrency { get; private set; } = "BDT";
    public long SumInsuredAmount { get; private set; }
    public string SumInsuredCurrency { get; private set; } = "BDT";
    public long VatTaxAmount { get; private set; }
    public long ServiceFeeAmount { get; private set; }
    public long TotalPayableAmount { get; private set; }

    public int TenureMonths { get; private set; }
    public DateTime StartDate { get; private set; }
    public DateTime EndDate { get; private set; }
    public DateTime? IssuedAt { get; private set; }

    public string? PaymentFrequency { get; private set; }
    public string? PaymentGatewayReference { get; set; }
    public string? ReceiptNumber { get; set; }
    public string? PolicyDocumentUrl { get; set; }

    // Applicant stored as JSONB
    public string? ProposerDetailsJson { get; private set; }

    public string? OccupationRiskClass { get; private set; }
    public bool HasExistingPolicies { get; private set; }
    public string? ClaimsHistorySummary { get; private set; }
    public string? ProviderName { get; private set; }
    public DateTime? EnrollmentStartDate { get; private set; }
    public DateTime? EnrollmentEndDate { get; private set; }
    public string? UnderwritingData { get; private set; }

    // Collections
    public List<Nominee> Nominees { get; private set; } = new();
    public List<PolicyRider> Riders { get; private set; } = new();

    // Audit
    public DateTime CreatedAt { get; private set; }
    public DateTime UpdatedAt { get; private set; }
    public DateTime? DeletedAt { get; private set; }
    public bool IsDeleted { get; private set; }

    // --- Money convenience accessors ---
    public Money PremiumMoney => new(PremiumAmount, PremiumCurrency);
    public Money SumInsuredMoney => new(SumInsuredAmount, SumInsuredCurrency);

    // EF Core constructor
    public PolicyAggregate() { }

    public static PolicyAggregate Create(
        string policyNumber, Guid productId, Guid customerId, Guid? partnerId, 
        long sumInsured, long premium, int tenureMonths, DateTime startDate)
    {
        return new PolicyAggregate
        {
            Id = Guid.NewGuid(),
            PolicyNumber = policyNumber,
            ProductId = productId,
            CustomerId = customerId,
            PartnerId = partnerId,
            SumInsuredAmount = sumInsured,
            PremiumAmount = premium,
            TenureMonths = tenureMonths,
            StartDate = startDate,
            EndDate = startDate.AddMonths(tenureMonths),
            Status = PolicyStatus.PendingPayment,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };
    }

    public static PolicyAggregate Renew(
        PolicyAggregate oldPolicy, 
        string newPolicyNumber, 
        int tenureMonths)
    {
        var newStartDate = oldPolicy.EndDate;
        long adjustedPremium = oldPolicy.PremiumAmount;
        
        // FR-067: Gamified Renewals (No Claim Bonus vs Penalty)
        if (string.IsNullOrWhiteSpace(oldPolicy.ClaimsHistorySummary))
        {
            // NCB: 10% discount if no claims
            adjustedPremium = (long)Math.Round(adjustedPremium * 0.90, MidpointRounding.AwayFromZero);
        }
        else
        {
            // Penalty: 15% increase if claims exist (M2 requirement)
            adjustedPremium = (long)Math.Round(adjustedPremium * 1.15, MidpointRounding.AwayFromZero);
        }

        // FR-068: Grace Period tracking & late fee (5% penalty)
        if (oldPolicy.Status == PolicyStatus.GracePeriod)
        {
            adjustedPremium = (long)Math.Round(adjustedPremium * 1.05, MidpointRounding.AwayFromZero);
        }

        var newPolicy = new PolicyAggregate
        {
            Id = Guid.NewGuid(),
            PolicyNumber = newPolicyNumber,
            ProductId = oldPolicy.ProductId,
            CustomerId = oldPolicy.CustomerId,
            PartnerId = oldPolicy.PartnerId,
            AgentId = oldPolicy.AgentId,
            Status = PolicyStatus.PendingPayment,
            PremiumAmount = adjustedPremium,
            PremiumCurrency = oldPolicy.PremiumCurrency,
            SumInsuredAmount = oldPolicy.SumInsuredAmount,
            SumInsuredCurrency = oldPolicy.SumInsuredCurrency,
            TenureMonths = tenureMonths,
            StartDate = newStartDate,
            EndDate = newStartDate.AddMonths(tenureMonths),
            ProposerDetailsJson = oldPolicy.ProposerDetailsJson,
            ProviderName = oldPolicy.ProviderName,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        return newPolicy;
    }

    // --- Lifecycle State Machine ---

    public Result Issue(DateTime issuedAt)
    {
        if (Status != PolicyStatus.PendingPayment)
            return Result.Fail(Error.InvalidStateTransition(
                $"Cannot issue policy in '{Status}' status. Only PENDING_PAYMENT policies can be issued."));

        if (!UnderwritingDecisionId.HasValue)
            return Result.Fail(Error.Validation("Cannot issue policy without a valid Underwriting Decision."));

        Status = PolicyStatus.Active;
        IssuedAt = issuedAt;
        UpdatedAt = DateTime.UtcNow;
        return Result.Ok();
    }

    public Result Cancel(string reason)
    {
        if (Status == PolicyStatus.Cancelled)
            return Result.Fail(Error.InvalidStateTransition("Policy is already cancelled."));
        if (Status == PolicyStatus.Expired)
            return Result.Fail(Error.InvalidStateTransition("Cannot cancel an expired policy."));

        Status = PolicyStatus.Cancelled;
        UpdatedAt = DateTime.UtcNow;
        return Result.Ok();
    }

    public Result Suspend()
    {
        if (Status != PolicyStatus.Active)
            return Result.Fail(Error.InvalidStateTransition("Only active policies can be suspended."));

        Status = PolicyStatus.Suspended;
        UpdatedAt = DateTime.UtcNow;
        return Result.Ok();
    }

    public Result EnterGracePeriod()
    {
        if (Status != PolicyStatus.Active)
            return Result.Fail(Error.InvalidStateTransition("Only active policies can enter grace period."));

        Status = PolicyStatus.GracePeriod;
        UpdatedAt = DateTime.UtcNow;
        return Result.Ok();
    }

    public Result Lapse()
    {
        if (Status != PolicyStatus.GracePeriod)
            return Result.Fail(Error.InvalidStateTransition("Only policies in grace period can lapse."));

        Status = PolicyStatus.Lapsed;
        UpdatedAt = DateTime.UtcNow;
        return Result.Ok();
    }

    public Result Expire()
    {
        Status = PolicyStatus.Expired;
        UpdatedAt = DateTime.UtcNow;
        return Result.Ok();
    }

    public void SetUnderwritingDecision(Guid decisionId)
    {
        UnderwritingDecisionId = decisionId;
        UpdatedAt = DateTime.UtcNow;
    }

    // --- Metadata and Riders ---

    public void SetAgent(Guid? agentId)
    {
        AgentId = agentId;
        UpdatedAt = DateTime.UtcNow;
    }

    public void SetProposerDetails(string json)
    {
        ProposerDetailsJson = json;
        UpdatedAt = DateTime.UtcNow;
    }

    public void AddRider(string name, long premium, string premiumCurrency, long coverage, string coverageCurrency)
    {
        var rider = new PolicyRider(Guid.NewGuid())
        {
            PolicyId = Id,
            RiderName = name,
            PremiumAmount = premium,
            PremiumCurrency = premiumCurrency,
            CoverageAmount = coverage,
            CoverageCurrency = coverageCurrency,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };
        Riders.Add(rider);
        UpdatedAt = DateTime.UtcNow;
    }

    public void RemoveRider(string riderName)
    {
        var rider = Riders.FirstOrDefault(r => r.RiderName.Equals(riderName, StringComparison.OrdinalIgnoreCase));
        if (rider != null)
        {
            Riders.Remove(rider);
            UpdatedAt = DateTime.UtcNow;
        }
    }

    public bool CanEndorse()
    {
        return Status == PolicyStatus.Active || Status == PolicyStatus.GracePeriod;
    }

    // --- Nominee Management ---

    public Result AddNominee(Guid? beneficiaryId, string fullName, string relationship, double sharePercentage,
        DateTime? dateOfBirth = null, string? nidNumber = null, string? phoneNumber = null, string? nomineeDobText = null)
    {
        var nominee = new Nominee(Guid.NewGuid())
        {
            PolicyId = Id,
            BeneficiaryId = beneficiaryId,
            FullName = fullName,
            Relationship = relationship,
            SharePercentage = sharePercentage,
            DateOfBirth = dateOfBirth,
            NidNumber = nidNumber,
            PhoneNumber = phoneNumber,
            NomineeDobText = nomineeDobText,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };
        Nominees.Add(nominee);
        return ValidateNomineeShares();
    }

    public Result UpdateNominee(Guid nomineeId, string? fullName, string? relationship, double? sharePercentage,
        DateTime? dateOfBirth = null, string? nidNumber = null, string? phoneNumber = null, string? nomineeDobText = null)
    {
        var nominee = Nominees.FirstOrDefault(n => n.Id == nomineeId && !n.IsDeleted);
        if (nominee == null)
            return Result.Fail(Error.NotFound("Nominee", nomineeId.ToString()));

        if (fullName != null) nominee.FullName = fullName;
        if (relationship != null) nominee.Relationship = relationship;
        if (sharePercentage != null) nominee.SharePercentage = sharePercentage.Value;
        if (dateOfBirth != null) nominee.DateOfBirth = dateOfBirth;
        if (nidNumber != null) nominee.NidNumber = nidNumber;
        if (phoneNumber != null) nominee.PhoneNumber = phoneNumber;
        if (nomineeDobText != null) nominee.NomineeDobText = nomineeDobText;
        nominee.UpdatedAt = DateTime.UtcNow;

        return ValidateNomineeShares();
    }

    public Result RemoveNominee(Guid nomineeId)
    {
        var nominee = Nominees.FirstOrDefault(n => n.Id == nomineeId && !n.IsDeleted);
        if (nominee == null)
            return Result.Fail(Error.NotFound("Nominee", nomineeId.ToString()));

        nominee.IsDeleted = true;
        nominee.UpdatedAt = DateTime.UtcNow;

        return ValidateNomineeShares();
    }

    // --- Sum Insured and Premium Adjustments ---

    public void UpdateSumInsured(long amount)
    {
        SumInsuredAmount = amount;
        UpdatedAt = DateTime.UtcNow;
    }

    public void ApplyPremiumAdjustment(long basePremium, long vat, long serviceFee, long totalPayable)
    {
        PremiumAmount = basePremium;
        VatTaxAmount = vat;
        ServiceFeeAmount = serviceFee;
        TotalPayableAmount = totalPayable;
        UpdatedAt = DateTime.UtcNow;
    }

    private Result ValidateNomineeShares()
    {
        var activeNominees = Nominees.Where(n => !n.IsDeleted).ToList();
        if (activeNominees.Count == 0) return Result.Ok();

        var totalShare = activeNominees.Sum(n => n.SharePercentage);
        if (Math.Abs(totalShare - 100.0) > 0.001)
            return Result.Fail(Error.Validation($"Nominee share percentages must sum to 100. Current: {totalShare:F2}"));

        return Result.Ok();
    }
}
