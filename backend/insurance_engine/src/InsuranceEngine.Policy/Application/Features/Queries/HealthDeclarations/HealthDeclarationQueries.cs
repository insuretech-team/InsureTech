using System;
using InsuranceEngine.Policy.Application.DTOs;
using MediatR;

namespace InsuranceEngine.Policy.Application.Features.Queries.HealthDeclarations;

/// <summary>
/// Get health declaration by its ID.
/// </summary>
public record GetHealthDeclarationQuery(Guid DeclarationId) : IRequest<UnderwritingHealthDeclarationResponseDto?>;

/// <summary>
/// Get health declaration by the associated quote ID.
/// </summary>
public record GetHealthDeclarationByQuoteQuery(Guid QuoteId) : IRequest<UnderwritingHealthDeclarationResponseDto?>;
