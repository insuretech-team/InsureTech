using System;
using MediatR;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Policy.Application.Features.Commands.Claims;

public record SubmitClaimCommand(
    Guid PolicyId,
    Guid CustomerId,
    ClaimType Type,
    long ClaimedAmount,
    DateTime IncidentDate,
    string IncidentDescription,
    string? PlaceOfIncident
) : IRequest<Result<Guid>>;
