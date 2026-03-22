using System;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Claims.Application.Interfaces;
using InsuranceEngine.Claims.Domain.Entities;
using InsuranceEngine.Claims.Domain.Enums;
using InsuranceEngine.Claims.Domain.Events;
using InsuranceEngine.Claims.Domain.Services;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Interfaces;
using MediatR;
using Microsoft.Extensions.Logging;

namespace InsuranceEngine.Claims.Application.Features.Commands.Claims;

public class SubmitClaimCommandHandler : IRequestHandler<SubmitClaimCommand, Result<Guid>>
{
    private readonly IClaimsRepository _claimsRepository;
    private readonly ClaimEligibilityValidator _eligibilityValidator;
    private readonly ClaimDocumentValidator _documentValidator;
    private readonly IEventBus _eventBus;
    private readonly IMediator _mediator;
    private readonly ILogger<SubmitClaimCommandHandler> _logger;

    public SubmitClaimCommandHandler(
        IClaimsRepository claimsRepository,
        ClaimEligibilityValidator eligibilityValidator,
        ClaimDocumentValidator documentValidator,
        IEventBus eventBus,
        IMediator mediator,
        ILogger<SubmitClaimCommandHandler> logger)
    {
        _claimsRepository = claimsRepository;
        _eligibilityValidator = eligibilityValidator;
        _documentValidator = documentValidator;
        _eventBus = eventBus;
        _mediator = mediator;
        _logger = logger;
    }

    public async Task<Result<Guid>> Handle(SubmitClaimCommand request, CancellationToken cancellationToken)
    {
        _logger.LogInformation("Submitting claim for policy {PolicyId}", request.PolicyId);
        
        // 1. Fetch Policy Info
        var policyQuery = new InsuranceEngine.Policy.Application.Features.Queries.GetPolicyQuery(request.PolicyId);
        var policy = await _mediator.Send(policyQuery, cancellationToken);
        
        if (policy == null)
            return Result<Guid>.Fail(Error.NotFound("Policy", request.PolicyId.ToString()));

        // 2. Document Validation (FR-099)
        var docRequests = request.Documents.Select(d => new ValidateDocumentRequest(d.FileName, d.FileSize));
        var docValidation = _documentValidator.Validate(docRequests);
        if (!docValidation.IsSuccess)
        {
            return Result<Guid>.Fail(docValidation.Error!);
        }

        // 3. Eligibility Validation (FR-042)
        var eligibility = await _eligibilityValidator.ValidateAsync(
            policy, request.Type, request.IncidentDate, cancellationToken);
        
        if (!eligibility.IsSuccess)
        {
            _logger.LogWarning("Claim eligibility failed for policy {PolicyId}: {Error}",
                request.PolicyId, eligibility.Error?.Message);
            return Result<Guid>.Fail(eligibility.Error!);
        }

        var claimNumber = await _claimsRepository.GetNextClaimNumberAsync(cancellationToken);
        
        // 4. Perform Synchronous Fraud Check (FD-001 to FD-007)
        var fraudCommand = new InsuranceEngine.Fraud.Application.Features.Commands.CheckFraud.CheckClaimForFraudCommand(
            Guid.Empty, // Temporary ID since claim isn't persisted yet
            request.PolicyId,
            request.CustomerId,
            request.ClaimedAmount,
            policy.SumInsuredAmount,
            request.Type.ToString(),
            request.PlaceOfIncident,
            request.IncidentDate,
            policy.IssuedAt ?? policy.CreatedAt);

        var fraudResult = await _mediator.Send(fraudCommand, cancellationToken);

        var isFlagged = fraudResult.IsSuccess && fraudResult.Value.Status == InsuranceEngine.Fraud.Domain.Enums.FraudCheckStatus.Flagged;
        var fraudScore = fraudResult.IsSuccess ? fraudResult.Value.RiskScore : 100;

        // FR-093: Zero Human Touch Claims (ZHTC) auto-approval logic
        bool isZhtcEligible = !isFlagged 
            && fraudScore <= 15 
            && request.ClaimedAmount <= 5_000_000 // 50,000 BDT
            && (DateTime.UtcNow - (policy.IssuedAt ?? policy.CreatedAt)).TotalDays > 365;

        var claim = new Claim
        {
            Id = Guid.NewGuid(),
            ClaimNumber = claimNumber,
            PolicyId = request.PolicyId,
            CustomerId = request.CustomerId,
            Type = request.Type,
            Status = isZhtcEligible ? ClaimStatus.Approved : (isFlagged ? ClaimStatus.UnderReview : ClaimStatus.Submitted),
            ClaimedAmount = request.ClaimedAmount,
            ClaimedCurrency = "BDT",
            ApprovedAmount = isZhtcEligible ? request.ClaimedAmount : 0,
            ApprovedCurrency = "BDT",
            IncidentDate = request.IncidentDate,
            IncidentDescription = request.IncidentDescription,
            PlaceOfIncident = request.PlaceOfIncident,
            BankDetailsForPayout = request.BankDetailsForPayout,
            SubmittedAt = DateTime.UtcNow,
            ApprovedAt = isZhtcEligible ? DateTime.UtcNow : null,
            ProcessingType = isZhtcEligible ? ClaimProcessingType.Auto : ClaimProcessingType.Manual,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        if (isZhtcEligible)
        {
            claim.ProcessorNotes = "Auto-approved via ZHTC rule (FR-093).";
        }

        // Map Documents
        if (request.Documents != null)
        {
            foreach (var d in request.Documents)
            {
                claim.Documents.Add(new ClaimDocument
                {
                    Id = Guid.NewGuid(),
                    ClaimId = claim.Id,
                    DocumentType = d.DocumentType,
                    FileUrl = d.FileUrl,
                    FileHash = d.FileHash,
                    UploadedAt = DateTime.UtcNow,
                    CreatedAt = DateTime.UtcNow,
                    UpdatedAt = DateTime.UtcNow
                });
            }
        }

        // Create FraudCheckResult entity (proto-aligned) instead of inline fields
        if (fraudResult.IsSuccess)
        {
            claim.FraudCheck = new FraudCheckResult
            {
                Id = Guid.NewGuid(),
                ClaimId = claim.Id,
                FraudScore = fraudResult.Value.RiskScore,
                RiskFactors = fraudResult.Value.Findings ?? new(),
                Flagged = isFlagged,
                CreatedAt = DateTime.UtcNow
            };
        }

        await _claimsRepository.CreateAsync(claim, cancellationToken);
        
        _logger.LogInformation("Claim {ClaimNumber} created with ID {ClaimId}. Fraud flagged: {IsFlagged}",
            claim.ClaimNumber, claim.Id, isFlagged);

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

