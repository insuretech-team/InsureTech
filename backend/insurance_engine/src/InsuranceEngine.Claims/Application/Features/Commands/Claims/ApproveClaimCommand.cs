using System;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Claims.Application.Interfaces;
using InsuranceEngine.Claims.Domain.Enums;
using InsuranceEngine.Claims.Domain.Services;
using InsuranceEngine.Products.Application.Interfaces;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;
using Microsoft.Extensions.Logging;

namespace InsuranceEngine.Claims.Application.Features.Commands.Claims;

/// <summary>
/// Approve a claim with automatic deductible/co-pay calculation (FR-100/FR-104).
/// </summary>
public record ApproveClaimCommand(
    Guid ClaimId,
    Guid ApproverId,
    string ApproverRole,
    int ApprovalLevel,
    ApprovalDecision Decision,
    long ApprovedAmount,
    string Notes
) : IRequest<Result>;

public class ApproveClaimCommandHandler : IRequestHandler<ApproveClaimCommand, Result>
{
    private readonly IClaimsRepository _claimsRepository;
    private readonly IProductRepository _productRepository;
    private readonly ClaimSettlementCalculator _calculator;
    private readonly IMediator _mediator;
    private readonly ILogger<ApproveClaimCommandHandler> _logger;

    public ApproveClaimCommandHandler(
        IClaimsRepository claimsRepository,
        IProductRepository productRepository,
        ClaimSettlementCalculator calculator,
        IMediator mediator,
        ILogger<ApproveClaimCommandHandler> logger)
    {
        _claimsRepository = claimsRepository;
        _productRepository = productRepository;
        _calculator = calculator;
        _mediator = mediator;
        _logger = logger;
    }

    public async Task<Result> Handle(ApproveClaimCommand request, CancellationToken cancellationToken)
    {
        var claim = await _claimsRepository.GetByIdAsync(request.ClaimId, cancellationToken);
        if (claim == null) return Result.Failure("Claim not found");

        if (claim.Status != ClaimStatus.Submitted && claim.Status != ClaimStatus.UnderReview)
            return Result.Failure($"Claim cannot be approved in '{claim.Status}' status");

        // --- Fetch product for deductible/co-pay config (cross-module via MediatR) ---
        var policyQuery = new InsuranceEngine.Policy.Application.Features.Queries.GetPolicyQuery(claim.PolicyId);
        var policy = await _mediator.Send(policyQuery, cancellationToken);

        double deductiblePct = 0;
        double coPayPct = 0;
        long maxDeductible = 0;

        if (policy != null)
        {
            var product = await _productRepository.GetByIdAsync(policy.ProductId);
            if (product != null)
            {
                deductiblePct = product.DeductiblePercentage;
                coPayPct = product.CoPayPercentage;
                maxDeductible = product.MaxDeductibleAmount;
            }
        }

        // --- Calculate settlement ---
        if (request.Decision == ApprovalDecision.Approved && request.ApprovedAmount > 0)
        {
            var settlement = _calculator.Calculate(
                claim.ClaimedAmount,
                request.ApprovedAmount,
                deductiblePct,
                coPayPct,
                maxDeductible);

            claim.DeductibleAmount = settlement.DeductibleAmount;
            claim.CoPayAmount = settlement.CoPayAmount;
            claim.SettledAmount = settlement.NetSettlementAmount;

            _logger.LogInformation(
                "Claim {ClaimId}: Approved={Approved}, Deductible={Deductible}, CoPay={CoPay}, Net={Net}",
                claim.Id, settlement.ApprovedAmount, settlement.DeductibleAmount,
                settlement.CoPayAmount, settlement.NetSettlementAmount);
        }

        // --- Standard approval flow ---
        var result = claim.AddApproval(
            request.ApproverId, request.ApproverRole, request.ApprovalLevel,
            request.Decision, request.ApprovedAmount, request.Notes);

        if (!result.IsSuccess) return result;

        await _claimsRepository.UpdateAsync(claim, cancellationToken);
        return Result.Ok();
    }
}
