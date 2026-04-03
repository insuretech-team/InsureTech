namespace InsuranceEngine.Claims.Domain.Events;

public sealed record ClaimSubmittedEvent(
    Guid ClaimId, 
    string ClaimNumber, 
    Guid PolicyId, 
    Guid CustomerId, 
    long ClaimedAmount, 
    Guid? PartnerId, 
    Guid? AgentId) : SharedKernel.Domain.DomainEvent;
