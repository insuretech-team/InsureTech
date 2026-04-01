using System;
using InsuranceEngine.SharedKernel.Domain;

namespace InsuranceEngine.SharedKernel.Domain.Events;

// Global Event Base for Consistent Metadata (PolicyId, CustomerId, Amount, Partner/Agent)
public abstract record InsuranceEngineEvent(
    Guid PolicyId, 
    Guid CustomerId, 
    long Amount, 
    string Currency = "BDT", 
    string Module = "POLICY",
    Guid? PartnerId = null,
    Guid? AgentId = null,
    DateTime? OccurredAt = null) : DomainEvent;

public record PolicyIssuedEvent(
    Guid PolicyId, 
    string PolicyNumber, 
    Guid CustomerId, 
    long PremiumAmount,
    Guid? PartnerId = null,
    Guid? AgentId = null) 
    : InsuranceEngineEvent(PolicyId, CustomerId, PremiumAmount, "BDT", "POLICY", PartnerId, AgentId, DateTime.UtcNow);

public record PolicyRenewedEvent(
    Guid PolicyId, 
    string PolicyNumber, 
    Guid CustomerId, 
    long PremiumAmount, 
    string NewExpiryDate,
    Guid? PartnerId = null,
    Guid? AgentId = null) 
    : InsuranceEngineEvent(PolicyId, CustomerId, PremiumAmount, "BDT", "RENEWAL", PartnerId, AgentId, DateTime.UtcNow);

public record PolicyCancelledEvent(
    Guid PolicyId, 
    string PolicyNumber, 
    Guid CustomerId, 
    long RefundAmount, 
    string Reason,
    Guid? PartnerId = null,
    Guid? AgentId = null) 
    : InsuranceEngineEvent(PolicyId, CustomerId, RefundAmount, "BDT", "CANCELLATIONS", PartnerId, AgentId, DateTime.UtcNow);

public record PolicyEndorsedEvent(
    Guid PolicyId, 
    string PolicyNumber, 
    Guid CustomerId, 
    string ChangeType,
    Guid? PartnerId = null,
    Guid? AgentId = null) 
    : InsuranceEngineEvent(PolicyId, CustomerId, 0, "BDT", "ENDORSEMENTS", PartnerId, AgentId, DateTime.UtcNow);

public record ClaimSubmittedEvent(
    Guid ClaimId, 
    string ClaimNumber, 
    Guid PolicyId, 
    Guid CustomerId, 
    long ClaimAmount,
    Guid? PartnerId = null,
    Guid? AgentId = null) 
    : InsuranceEngineEvent(PolicyId, CustomerId, ClaimAmount, "BDT", "CLAIMS", PartnerId, AgentId, DateTime.UtcNow);
