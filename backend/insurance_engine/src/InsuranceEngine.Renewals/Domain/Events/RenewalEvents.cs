namespace InsuranceEngine.Renewals.Domain.Events;

public sealed record PolicyRenewedEvent(
    Guid PolicyId, 
    string PolicyNumber, 
    Guid CustomerId, 
    long PremiumAmount, 
    string NewEndDate, 
    Guid? PartnerId, 
    Guid? AgentId) : SharedKernel.Domain.DomainEvent;
