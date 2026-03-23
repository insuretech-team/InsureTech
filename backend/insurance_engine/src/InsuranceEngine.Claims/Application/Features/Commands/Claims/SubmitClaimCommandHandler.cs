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

        // 1.1 Fetch Product Info
        var productQuery = new InsuranceEngine.Products.Application.Features.Queries.GetProduct.GetProductQuery(policy.ProductId);
        var product = await _mediator.Send(productQuery, cancellationToken);
        if (product == null)
            return Result<Guid>.Fail(Error.NotFound("Product", policy.ProductId.ToString()));

        // 2. Document Validation (FR-099)
        var docRequests = request.Documents?.Select(d => new ValidateDocumentRequest(d.FileName, d.FileSize)) ?? Enumerable.Empty<ValidateDocumentRequest>();
        var docValidation = _documentValidator.Validate(docRequests);
        if (!docValidation.IsSuccess)
        {
            return Result<Guid>.Fail(docValidation.Error!);
        }

        // 3. Eligibility Validation (FR-042)
        var eligibility = await _eligibilityValidator.ValidateAsync(
            policy, product, request.Type, request.IncidentDate, cancellationToken);
        
        if (!eligibility.IsSuccess)
        {
            _logger.LogWarning("Claim eligibility failed for policy {PolicyId}: {Error}",
                request.PolicyId, eligibility.Error?.Message);
            return Result<Guid>.Fail(eligibility.Error!);
        }

        var claimNumber = await _claimsRepository.GetNextClaimNumberAsync(cancellationToken);
        
        // 4. Perform Synchronous Fraud Check
        var fraudCommand = new InsuranceEngine.Fraud.Application.Features.Commands.CheckFraud.CheckClaimForFraudCommand(
            Guid.Empty,
            request.PolicyId,
            request.CustomerId,
            request.ClaimedAmount.Amount,
            policy.SumInsured.Amount,
            request.Type.ToString(),
            request.PlaceOfIncident,
            request.IncidentDate,
            policy.IssuedAt ?? policy.CreatedAt);

        var fraudResult = await _mediator.Send(fraudCommand, cancellationToken);
        var fraudScore = fraudResult.IsSuccess ? (double)fraudResult.Value.RiskScore / 100.0 : 1.0;

        // 5. Create Aggregate using Factory
        var claim = Claim.File(
            claimNumber,
            request.PolicyId,
            request.CustomerId,
            request.Type,
            request.ClaimedAmount.Amount,
            request.IncidentDate,
            request.IncidentDescription,
            request.PlaceOfIncident);

        claim.BankDetailsForPayout = request.BankDetailsForPayout;

        // 5.1 Calculate Financials (FR-100)
        // Deductible/CoPay was removed from Product to align with proto. 
        // Defaulting to 0 for now as it's not in the canonical schema.
        long deductible = 0; 
        claim.CalculateFinancials(deductible, 0.0);

        // 7. Apply Fraud Check (ZHTC Auto-Approve if < 10k and low risk)
        if (fraudResult.IsSuccess)
        {
            var riskFactors = fraudResult.Value.Findings ?? new List<string>();
            claim.ApplyFraudCheck(
                fraudResult.Value.CheckId,
                fraudResult.Value.RiskScore,
                riskFactors);
        }

        // 7. Add Documents via Aggregate
        if (request.Documents != null)
        {
            foreach (var d in request.Documents)
            {
                claim.AddDocument(
                    d.DocumentType,
                    d.FileUrl,
                    d.FileHash
                );
            }
        }

        // Check if flagged for review (FR-097)
        var isFlagged = fraudResult.IsSuccess && (
            fraudResult.Value.RiskLevel >= InsuranceEngine.Fraud.Domain.Enums.FraudRiskLevel.High || 
            fraudResult.Value.Status == InsuranceEngine.Fraud.Domain.Enums.FraudCheckStatus.Flagged);
        
        if (isFlagged && claim.Status != ClaimStatus.Approved) // Don't override ZHTC if it somehow reached here
        {
            claim.UpdateStatus(ClaimStatus.UnderReview);
        }

        await _claimsRepository.CreateAsync(claim, cancellationToken);
        
        _logger.LogInformation("Claim {ClaimNumber} created with ID {ClaimId}. Status: {Status}. Processing: {Processing}",
            claim.ClaimNumber, claim.Id, claim.Status, claim.ProcessingType);

        // Publish Submission Event
        await _eventBus.PublishAsync("insurance.claims.v1", new ClaimSubmittedEvent(
            ClaimId: claim.Id,
            ClaimNumber: claim.ClaimNumber,
            PolicyId: claim.PolicyId,
            CustomerId: claim.CustomerId,
            Amount: claim.ClaimedAmount,
            Currency: claim.ClaimedCurrency,
            IncidentDate: claim.IncidentDate
        ));

        // Publish Approval Event for ZHTC (Instant Settlement)
        if (claim.Status == ClaimStatus.Approved)
        {
            await _eventBus.PublishAsync("insurance.claims.v1", new ClaimProcessedEvent(
                ClaimId: claim.Id,
                ClaimNumber: claim.ClaimNumber,
                NewStatus: claim.Status,
                ApprovedAmount: claim.ApprovedAmount,
                Notes: $"Auto-Approved via ZHTC. Fraud Score: {fraudScore:P2}"
            ));
        }
        
        return Result<Guid>.Success(claim.Id);
    }
}

