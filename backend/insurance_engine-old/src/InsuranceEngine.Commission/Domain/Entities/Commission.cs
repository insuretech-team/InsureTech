using System;
using InsuranceEngine.SharedKernel.Domain;
using InsuranceEngine.Commission.Domain.Enums;

namespace InsuranceEngine.Commission.Domain.Entities;

public class Commission : AggregateRoot<Guid>
{
    public Guid PolicyId { get; private set; }
    public Guid? PartnerId { get; private set; }
    public Guid? AgentId { get; private set; }
    public CommissionType Type { get; private set; }
    public long Amount { get; private set; }
    public string Currency { get; private set; } = "BDT";
    public CommissionStatus Status { get; private set; }
    public Guid? PayoutId { get; private set; }
    public DateTime CreatedAt { get; private set; }
    public DateTime UpdatedAt { get; private set; }

    private Commission() { }

    public static Commission Create(
        Guid policyId, 
        Guid? partnerId, 
        Guid? agentId, 
        CommissionType type, 
        long amount)
    {
        return new Commission
        {
            Id = Guid.NewGuid(),
            PolicyId = policyId,
            PartnerId = partnerId,
            AgentId = agentId,
            Type = type,
            Amount = amount,
            Status = CommissionStatus.Pending,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };
    }

    public void MarkAsProcessing(Guid payoutId)
    {
        Status = CommissionStatus.Processing;
        PayoutId = payoutId;
        UpdatedAt = DateTime.UtcNow;
    }

    public void MarkAsPaid()
    {
        Status = CommissionStatus.Paid;
        UpdatedAt = DateTime.UtcNow;
    }
}

public class Payout : AggregateRoot<Guid>
{
    public Guid RecipientId { get; private set; }
    public long TotalAmount { get; private set; }
    public string Currency { get; private set; } = "BDT";
    public PayoutStatus Status { get; private set; }
    public string? PaymentReference { get; private set; }
    public DateTime CreatedAt { get; private set; }
    public DateTime? PaidAt { get; private set; }

    private Payout() { }

    public static Payout Create(Guid recipientId, long totalAmount)
    {
        return new Payout
        {
            Id = Guid.NewGuid(),
            RecipientId = recipientId,
            TotalAmount = totalAmount,
            Status = PayoutStatus.Pending,
            CreatedAt = DateTime.UtcNow
        };
    }

    public void MarkAsPaid(string reference)
    {
        Status = PayoutStatus.Paid;
        PaymentReference = reference;
        PaidAt = DateTime.UtcNow;
    }
}
