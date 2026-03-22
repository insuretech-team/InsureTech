using System;
using System.Collections.Generic;
using MediatR;
using InsuranceEngine.Claims.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.Claims.Application.DTOs;

namespace InsuranceEngine.Claims.Application.Features.Commands.Claims;

public record SubmitClaimCommand(
    Guid PolicyId,
    Guid CustomerId,
    ClaimType Type,
    long ClaimedAmount,
    DateTime IncidentDate,
    string IncidentDescription,
    string? PlaceOfIncident,
    string? BankDetailsForPayout,
    List<ClaimDocumentDto> Documents
) : IRequest<Result<Guid>>;
