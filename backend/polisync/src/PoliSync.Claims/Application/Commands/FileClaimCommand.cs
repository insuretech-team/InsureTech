using MediatR;
using PoliSync.SharedKernel.CQRS;
using Insuretech.Claims.Entity.V1;

namespace PoliSync.Claims.Application.Commands;

public sealed record FileClaimCommand(
    string PolicyId,
    string CustomerId,
    ClaimType ClaimType,
    long ClaimedAmountPaisa,
    DateTime IncidentDate,
    string IncidentDescription,
    string PlaceOfIncident
) : ICommand<string>; // Returns claim_id
