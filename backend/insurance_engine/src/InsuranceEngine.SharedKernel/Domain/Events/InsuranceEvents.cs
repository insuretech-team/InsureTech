using System;
using InsuranceEngine.SharedKernel.Domain;

namespace InsuranceEngine.SharedKernel.Domain.Events;

public record PolicyIssuedEvent(Guid PolicyId, string PolicyNumber, string CustomerId, long PremiumAmount) : DomainEvent;

public record ClaimSubmittedEvent(Guid ClaimId, string ClaimNumber, Guid PolicyId, long ClaimAmount) : DomainEvent;
