using System;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Policy.Application.Features.Commands.HealthDeclarations;

/// <summary>
/// Submit a standalone health declaration for an existing quote.
/// </summary>
public record SubmitHealthDeclarationCommand(
    Guid QuoteId,
    UnderwritingHealthDeclarationDto HealthDeclaration
) : IRequest<Result<Guid>>;
