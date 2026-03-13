using System;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Claims.Application.Interfaces;
using InsuranceEngine.Claims.Domain.Entities;
using InsuranceEngine.Claims.Domain.Enums;
using InsuranceEngine.Claims.Domain.Events;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Interfaces;
using MediatR;
using Microsoft.Extensions.Logging;

namespace InsuranceEngine.Claims.Application.Features.Commands.Claims;

public class SubmitClaimCommandHandler : IRequestHandler<SubmitClaimCommand, Result<Guid>>
{
    private readonly IClaimsRepository _claimsRepository;
    private readonly IEventBus _eventBus;
    private readonly ILogger<SubmitClaimCommandHandler> _logger;

    public SubmitClaimCommandHandler(
        IClaimsRepository claimsRepository, 
        IEventBus eventBus,
        ILogger<SubmitClaimCommandHandler> logger)
    {
        _claimsRepository = claimsRepository;
        _eventBus = eventBus;
        _logger = logger;
    }

    public async Task<Result<Guid>> Handle(SubmitClaimCommand request, CancellationToken cancellationToken)
    {
        _logger.LogInformation($"Submitting claim for policy {request.PolicyId}");
        
        var claimNumber = await _claimsRepository.GetNextClaimNumberAsync(cancellationToken);
        
        var claim = new Claim
        {
            Id = Guid.NewGuid(),
            ClaimNumber = claimNumber,
            PolicyId = request.PolicyId,
            CustomerId = request.CustomerId,
            Type = request.Type,
            Status = ClaimStatus.Submitted,
            ClaimedAmount = request.ClaimedAmount,
            ClaimedCurrency = "BDT",
            IncidentDate = request.IncidentDate,
            IncidentDescription = request.IncidentDescription,
            PlaceOfIncident = request.PlaceOfIncident,
            SubmittedAt = DateTime.UtcNow,
            ProcessingType = ClaimProcessingType.Manual,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        await _claimsRepository.CreateAsync(claim, cancellationToken);
        
        _logger.LogInformation($"Claim {claim.ClaimNumber} created with ID {claim.Id}");

        // Publish to Kafka (fixed property access)
        await _eventBus.PublishAsync("insurance.claims.v1", new ClaimSubmittedEvent(
            ClaimId: claim.Id,
            ClaimNumber: claim.ClaimNumber,
            PolicyId: claim.PolicyId,
            CustomerId: claim.CustomerId,
            Amount: claim.ClaimedAmount,
            Currency: claim.ClaimedCurrency,
            IncidentDate: claim.IncidentDate
        ));
        
        return Result<Guid>.Success(claim.Id);
    }
}
