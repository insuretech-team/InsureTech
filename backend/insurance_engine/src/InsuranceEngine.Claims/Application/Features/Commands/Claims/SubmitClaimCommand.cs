using System;
using MediatR;
using InsuranceEngine.Claims.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Claims.Application.Features.Commands.Claims;

public record SubmitClaimCommand(
    Guid PolicyId,
    Guid CustomerId,
    ClaimType Type,
    long ClaimedAmount,
    DateTime IncidentDate,
    string IncidentDescription,
    string? PlaceOfIncident,
    string? BankDetailsForPayout
) : IRequest<Result<Guid>>;
