using Microsoft.EntityFrameworkCore;
using PoliSync.SharedKernel.Domain;
using Insuretech.Policy.Entity.V1;

namespace PoliSync.Policy.Domain;

/// <summary>
/// Policy aggregate root - uses proto entity as the domain model
/// Maps directly to insurance_schema.policies table
/// </summary>
[Index(nameof(PolicyNumber), IsUnique = true)]
[Index(nameof(CustomerId))]
[Index(nameof(ProductId))]
[Index(nameof(Status))]
public class PolicyAggregate
{
    // Proto entity is the domain model
    private readonly Insuretech.Policy.Entity.V1.Policy _policy;
    
    private readonly List<DomainEvent> _domainEvents = new();
    
    public PolicyAggregate(Insuretech.Policy.Entity.V1.Policy policy)
    {
        _policy = policy ?? throw new ArgumentNullException(nameof(policy));
    }
    
    // Expose proto entity
    public Insuretech.Policy.Entity.V1.Policy Policy => _policy;
    
    // Domain properties
    public string PolicyId => _policy.PolicyId;
    public string PolicyNumber => _policy.PolicyNumber;
    public string CustomerId => _policy.CustomerId;
    public string ProductId => _policy.ProductId;
    public PolicyStatus Status => _policy.Status;
    
    public IReadOnlyCollection<DomainEvent> DomainEvents => _domainEvents.AsReadOnly();
    
    public void ClearDomainEvents() => _domainEvents.Clear();
    
    // Factory method
    public static PolicyAggregate Create(
        string customerId,
        string productId,
        string quoteId,
        long premiumAmountPaisa,
        long sumInsuredPaisa,
        int tenureMonths,
        DateTime startDate,
        DateTime endDate)
    {
        var policy = new Insuretech.Policy.Entity.V1.Policy
        {
            PolicyId = Guid.NewGuid().ToString(),
            PolicyNumber = GeneratePolicyNumber(),
            CustomerId = customerId,
            ProductId = productId,
            QuoteId = quoteId,
            Status = PolicyStatus.PendingPayment,
            PremiumAmount = new Insuretech.Common.V1.Money 
            { 
                Amount = premiumAmountPaisa,
                Currency = "BDT"
            },
            SumInsured = new Insuretech.Common.V1.Money 
            { 
                Amount = sumInsuredPaisa,
                Currency = "BDT"
            },
            TenureMonths = tenureMonths,
            StartDate = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(startDate.ToUniversalTime()),
            EndDate = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(endDate.ToUniversalTime()),
            CreatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow)
        };
        
        var aggregate = new PolicyAggregate(policy);
        aggregate._domainEvents.Add(new PolicyCreatedEvent(policy.PolicyId));
        
        return aggregate;
    }
    
    // Business methods
    public void IssuePolicy()
    {
        if (Status != PolicyStatus.PendingPayment)
            throw new InvalidOperationException($"Cannot issue policy in status {Status}");
        
        _policy.Status = PolicyStatus.Active;
        _policy.IssuedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow);
        _policy.UpdatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow);
        
        _domainEvents.Add(new PolicyIssuedEvent(PolicyId, CustomerId, ProductId));
    }
    
    public void CancelPolicy(string reason)
    {
        if (Status != PolicyStatus.Active)
            throw new InvalidOperationException($"Cannot cancel policy in status {Status}");
        
        _policy.Status = PolicyStatus.Cancelled;
        _policy.UpdatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow);
        
        _domainEvents.Add(new PolicyCancelledEvent(PolicyId, reason));
    }
    
    public void SuspendPolicy()
    {
        if (Status != PolicyStatus.Active)
            throw new InvalidOperationException($"Cannot suspend policy in status {Status}");
        
        _policy.Status = PolicyStatus.Suspended;
        _policy.UpdatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow);
        
        _domainEvents.Add(new PolicySuspendedEvent(PolicyId));
    }
    
    public void ReinstatePolicy()
    {
        if (Status != PolicyStatus.Suspended && Status != PolicyStatus.Lapsed)
            throw new InvalidOperationException($"Cannot reinstate policy in status {Status}");
        
        _policy.Status = PolicyStatus.Active;
        _policy.UpdatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow);
        
        _domainEvents.Add(new PolicyReinstatedEvent(PolicyId));
    }
    
    public void MarkAsLapsed()
    {
        if (Status != PolicyStatus.Active)
            throw new InvalidOperationException($"Cannot lapse policy in status {Status}");
        
        _policy.Status = PolicyStatus.Lapsed;
        _policy.UpdatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow);
        
        _domainEvents.Add(new PolicyLapsedEvent(PolicyId));
    }
    
    public void SetDocumentUrl(string documentUrl)
    {
        _policy.PolicyDocumentUrl = documentUrl;
        _policy.UpdatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow);
    }
    
    // Policy number generation: LBT-YYYY-XXXX-NNNNNN
    private static string GeneratePolicyNumber()
    {
        var year = DateTime.UtcNow.Year;
        var random = new Random().Next(1000, 9999);
        var sequence = new Random().Next(100000, 999999);
        return $"LBT-{year}-{random:D4}-{sequence:D6}";
    }
}

// Domain Events
public sealed record PolicyCreatedEvent(string PolicyId) : DomainEvent;
public sealed record PolicyIssuedEvent(string PolicyId, string CustomerId, string ProductId) : DomainEvent;
public sealed record PolicyCancelledEvent(string PolicyId, string Reason) : DomainEvent;
public sealed record PolicySuspendedEvent(string PolicyId) : DomainEvent;
public sealed record PolicyReinstatedEvent(string PolicyId) : DomainEvent;
public sealed record PolicyLapsedEvent(string PolicyId) : DomainEvent;
