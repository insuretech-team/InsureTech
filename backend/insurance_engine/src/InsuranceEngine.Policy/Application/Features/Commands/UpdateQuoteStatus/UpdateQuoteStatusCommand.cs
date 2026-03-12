using System;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Policy.Application.Features.Commands.UpdateQuoteStatus;

/// <summary>
/// Update a quote's status (e.g., expire, cancel, convert to policy).
/// </summary>
public record UpdateQuoteStatusCommand(
    Guid QuoteId,
    string NewStatus, // "EXPIRED", "CANCELLED", "CONVERTED"
    Guid? ConvertedPolicyId = null
) : IRequest<Result>;
