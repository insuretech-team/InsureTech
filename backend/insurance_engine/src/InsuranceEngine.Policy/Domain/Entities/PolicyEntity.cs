using System;
using System.Collections.Generic;
using System.Linq;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.Policy.Domain.ValueObjects;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Domain.ValueObjects;

namespace InsuranceEngine.Policy.Domain.Entities;

/// <summary>
/// Policy aggregate root. Maps to 'policies' table in insurance_schema.
/// Enforces lifecycle state machine and nominee share invariant.
/// </summary>
public class PolicyEntity
{
    public Guid Id { get; set; }
    public string PolicyNumber { get; set; } = string.Empty;
    public Guid ProductId { get; set; }
    public Guid CustomerId { get; set; }
    public Guid? PartnerId { get; set; }
    public Guid? AgentId { get; set; }
    public Guid? QuoteId { get; set; }
    public Guid? UnderwritingDecisionId { get; set; }

    public PolicyStatus Status { get; set; }

    // Money fields — stored as bigint (paisa)
    public long PremiumAmount { get; set; }
    public string PremiumCurrency { get; set; } = "BDT";
    public long SumInsuredAmount { get; set; }
    public string SumInsuredCurrency { get; set; } = "BDT";
    public long VatTaxAmount { get; set; }
    public long ServiceFeeAmount { get; set; }
    public long TotalPayableAmount { get; set; }

    public int TenureMonths { get; set; }
    public DateTime StartDate { get; set; }
    public DateTime EndDate { get; set; }
    public DateTime? IssuedAt { get; set; }

    public string? PaymentFrequency { get; set; }
    public string? PaymentGatewayReference { get; set; }
    public string? ReceiptNumber { get; set; }
    public string? PolicyDocumentUrl { get; set; }

    // Applicant stored as JSONB
    public string? ProposerDetailsJson { get; set; }

    public string? OccupationRiskClass { get; set; }
    public bool HasExistingPolicies { get; set; }
    public string? ClaimsHistorySummary { get; set; }
    public string? ProviderName { get; set; }
    public DateTime? EnrollmentStartDate { get; set; }
    public DateTime? EnrollmentEndDate { get; set; }
    public string? UnderwritingData { get; set; }

    // Collections
    public List<Nominee> Nominees { get; set; } = new();
    public List<PolicyRider> Riders { get; set; } = new();

    // Audit
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    public DateTime? DeletedAt { get; set; }
    public bool IsDeleted { get; set; }

    // --- Money convenience accessors ---
    public Money PremiumMoney => new(PremiumAmount, PremiumCurrency);
    public Money SumInsuredMoney => new(SumInsuredAmount, SumInsuredCurrency);

    // --- Lifecycle State Machine ---

    public Result Issue(DateTime issuedAt)
    {
        if (Status != PolicyStatus.PendingPayment)
            return Result.Fail(Error.InvalidStateTransition(
                $"Cannot issue policy in '{Status}' status. Only PENDING_PAYMENT policies can be issued."));

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
            return Result.Fail(Error.InvalidStateTransition(
                "Only active policies can be suspended."));

        Status = PolicyStatus.Suspended;
        UpdatedAt = DateTime.UtcNow;
        return Result.Ok();
    }

    public Result EnterGracePeriod()
    {
        if (Status != PolicyStatus.Active)
            return Result.Fail(Error.InvalidStateTransition(
                "Only active policies can enter grace period."));

        Status = PolicyStatus.GracePeriod;
        UpdatedAt = DateTime.UtcNow;
        return Result.Ok();
    }

    public Result Lapse()
    {
        if (Status != PolicyStatus.GracePeriod)
            return Result.Fail(Error.InvalidStateTransition(
                "Only policies in grace period can lapse."));

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

    public bool CanEndorse()
    {
        return Status == PolicyStatus.Active || Status == PolicyStatus.GracePeriod;
    }

    // --- Nominee share invariant ---

    public Result AddNominee(Guid? beneficiaryId, string fullName, string relationship, double sharePercentage,
        DateTime? dateOfBirth = null, string? nidNumber = null, string? phoneNumber = null, string? nomineeDobText = null)
    {
        var nominee = new Nominee
        {
            Id = Guid.NewGuid(),
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

        // If there are remaining active nominees, validate shares
        var activeNominees = Nominees.Where(n => !n.IsDeleted).ToList();
        if (activeNominees.Count > 0)
            return ValidateNomineeShares();

        return Result.Ok();
    }

    private Result ValidateNomineeShares()
    {
        var activeNominees = Nominees.Where(n => !n.IsDeleted).ToList();
        if (activeNominees.Count == 0)
            return Result.Ok();

        var totalShare = activeNominees.Sum(n => n.SharePercentage);
        if (Math.Abs(totalShare - 100.0) > 0.001)
            return Result.Fail(Error.Validation(
                $"Nominee share percentages must sum to 100. Current sum: {totalShare:F2}"));

        return Result.Ok();
    }
}
