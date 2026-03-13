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
    private readonly IMediator _mediator;
    private readonly ILogger<SubmitClaimCommandHandler> _logger;

    public SubmitClaimCommandHandler(
        IClaimsRepository claimsRepository, 
        IEventBus eventBus,
        IMediator mediator,
        ILogger<SubmitClaimCommandHandler> logger)
    {
        _claimsRepository = claimsRepository;
        _eventBus = eventBus;
        _mediator = mediator;
        _logger = logger;
    }

    public async Task<Result<Guid>> Handle(SubmitClaimCommand request, CancellationToken cancellationToken)
    {
        _logger.LogInformation($"Submitting claim for policy {request.PolicyId}");
        
        // 1. Fetch Policy Info for Fraud Check (FD-001)
        var policyQuery = new InsuranceEngine.Policy.Application.Features.Queries.GetPolicyQuery(request.PolicyId);
        var policy = await _mediator.Send(policyQuery, cancellationToken);
        
        if (policy == null)
            return Result<Guid>.Fail(Error.NotFound("Policy", request.PolicyId.ToString()));

        var claimNumber = await _claimsRepository.GetNextClaimNumberAsync(cancellationToken);
        
        // 2. Perform Synchronous Fraud Check
        var fraudCommand = new InsuranceEngine.Fraud.Application.Features.Commands.CheckFraud.CheckClaimForFraudCommand(
            Guid.Empty, // Temporary ID since claim isn't persisted yet
            request.PolicyId,
            request.ClaimedAmount,
            request.IncidentDate,
            policy.IssuedAt ?? policy.CreatedAt);

        var fraudResult = await _mediator.Send(fraudCommand, cancellationToken);
        
        var claim = new Claim
        {
            Id = Guid.NewGuid(),
            ClaimNumber = claimNumber,
            PolicyId = request.PolicyId,
            CustomerId = request.CustomerId,
            Type = request.Type,
            Status = (fraudResult.IsSuccess && fraudResult.Value.Status == InsuranceEngine.Fraud.Domain.Enums.FraudCheckStatus.Flagged) 
                     ? ClaimStatus.UnderReview 
                     : ClaimStatus.Submitted,
            ClaimedAmount = request.ClaimedAmount,
            ClaimedCurrency = "BDT",
            IncidentDate = request.IncidentDate,
            IncidentDescription = request.IncidentDescription,
            PlaceOfIncident = request.PlaceOfIncident,
            SubmittedAt = DateTime.UtcNow,
            ProcessingType = ClaimProcessingType.Manual,
            FraudScore = fraudResult.IsSuccess ? fraudResult.Value.RiskScore : 0,
            FraudCheckData = fraudResult.IsSuccess ? System.Text.Json.JsonSerializer.Serialize(fraudResult.Value.Findings) : null,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        if (fraudResult.IsSuccess)
        {
            // Update the fraud check with the actual claim ID
            // In a more robust system, we would do this via a domain event or after persistence.
        }

        await _claimsRepository.CreateAsync(claim, cancellationToken);
        
        _logger.LogInformation($"Claim {claim.ClaimNumber} created with ID {claim.Id}. Fraud Score: {claim.FraudScore}");

        // Publish to Kafka
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
