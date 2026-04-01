using PoliSync.SharedKernel.Domain;

namespace PoliSync.Quotes.Domain;

/// <summary>
/// Raised when a new quotation is created
/// </summary>
public sealed record QuotationCreatedEvent(Guid QuotationId, string QuotationNumber) : DomainEvent;

/// <summary>
/// Raised when a quotation is submitted for underwriting
/// </summary>
public sealed record QuotationSubmittedEvent(
    Guid QuotationId,
    string QuotationNumber,
    Guid ProductId,
    Guid CustomerId) : DomainEvent;

/// <summary>
/// Raised when a quotation is approved
/// </summary>
public sealed record QuotationApprovedEvent(
    Guid QuotationId,
    string QuotationNumber,
    long TotalPayable) : DomainEvent;

/// <summary>
/// Raised when a quotation is rejected
/// </summary>
public sealed record QuotationRejectedEvent(
    Guid QuotationId,
    string QuotationNumber,
    string Reason) : DomainEvent;

/// <summary>
/// Raised when a quotation expires
/// </summary>
public sealed record QuotationExpiredEvent(
    Guid QuotationId,
    string QuotationNumber) : DomainEvent;
