using Google.Protobuf.WellKnownTypes;
using Insuretech.Common.V1;
using Microsoft.AspNetCore.Mvc;
using PoliSync.ApiHost.Services;
using PoliSync.Policy.Infrastructure;
using PoliSync.SharedKernel.Auth;
using ProposalEntity = Insuretech.Policy.Entity.V1.InsuranceProposal;
using ProposalStatus = Insuretech.Policy.Entity.V1.ProposalStatus;

namespace PoliSync.ApiHost.Controllers;

[ApiController]
public sealed class InsuranceProposalsController : ControllerBase
{
    private readonly IPolicyDataGateway _policyGateway;
    private readonly InsuranceProposalWorkflowService _workflowService;
    private readonly ICurrentUser _currentUser;

    public InsuranceProposalsController(
        IPolicyDataGateway policyGateway,
        InsuranceProposalWorkflowService workflowService,
        ICurrentUser currentUser)
    {
        _policyGateway = policyGateway;
        _workflowService = workflowService;
        _currentUser = currentUser;
    }

    [HttpGet("/v1/insurance-proposals")]
    public async Task<IActionResult> List(
        [FromQuery(Name = "order_id")] string? orderId = null,
        [FromQuery(Name = "insurer_id")] string? insurerId = null,
        [FromQuery(Name = "customer_id")] string? customerId = null,
        [FromQuery] string? status = null,
        [FromQuery] int page = 1,
        [FromQuery(Name = "page_size")] int pageSize = 20,
        CancellationToken cancellationToken = default)
    {
        if (!TryParseStatus(status, out var proposalStatus))
        {
            return BadRequest(new { error = $"Unknown proposal status '{status}'." });
        }

        var proposals = await _policyGateway.ListInsuranceProposalsAsync(
            orderId,
            insurerId,
            customerId,
            proposalStatus,
            page,
            pageSize,
            cancellationToken);

        return Ok(proposals);
    }

    [HttpGet("/v1/insurance-proposals/{proposalId}")]
    public async Task<IActionResult> GetById(string proposalId, CancellationToken cancellationToken = default)
    {
        var proposal = await _policyGateway.GetInsuranceProposalAsync(proposalId, cancellationToken);
        return proposal is null
            ? NotFound(new { error = "Insurance proposal not found." })
            : Ok(proposal);
    }

    [HttpPost("/v1/insurance-proposals")]
    public async Task<IActionResult> Create(
        [FromBody] CreateInsuranceProposalRequest request,
        CancellationToken cancellationToken = default)
    {
        var now = DateTime.UtcNow;
        var proposal = new ProposalEntity
        {
            ProposalId = Guid.NewGuid().ToString(),
            ProposalNumber = string.IsNullOrWhiteSpace(request.ProposalNumber)
                ? $"PRP-{DateTime.UtcNow:yyyyMMdd}-{Guid.NewGuid():N}"[..24]
                : request.ProposalNumber,
            TenantId = FirstNonEmpty(request.TenantId, _currentUser.TenantId == Guid.Empty ? null : _currentUser.TenantId.ToString()) ?? string.Empty,
            OrderId = request.OrderId,
            QuotationId = request.QuotationId,
            CustomerId = request.CustomerId,
            InsurerId = request.InsurerId,
            ProductId = request.ProductId,
            PlanId = request.PlanId,
            ProposedPremium = NewMoney(request.ProposedPremiumAmount, request.ProposedPremiumCurrency),
            ProposedSumInsured = NewMoney(request.ProposedSumInsuredAmount, request.ProposedSumInsuredCurrency),
            Status = request.Status ?? ProposalStatus.Submitted,
            SubmissionPayload = request.SubmissionPayload ?? string.Empty,
            CorrelationId = request.CorrelationId ?? string.Empty,
            SubmittedAt = Timestamp.FromDateTime((request.SubmittedAt ?? now).ToUniversalTime()),
            CreatedAt = Timestamp.FromDateTime(now),
            UpdatedAt = Timestamp.FromDateTime(now)
        };

        var created = await _policyGateway.CreateInsuranceProposalAsync(proposal, cancellationToken);
        return Created($"/v1/insurance-proposals/{created.ProposalId}", created);
    }

    [HttpPatch("/v1/insurance-proposals/{proposalId}")]
    public async Task<IActionResult> Update(
        string proposalId,
        [FromBody] UpdateInsuranceProposalRequest request,
        CancellationToken cancellationToken = default)
    {
        var proposal = await _policyGateway.GetInsuranceProposalAsync(proposalId, cancellationToken);
        if (proposal is null)
        {
            return NotFound(new { error = "Insurance proposal not found." });
        }

        if (request.Status.HasValue)
        {
            proposal.Status = request.Status.Value;
        }

        if (request.ProposedPremiumAmount.HasValue)
        {
            proposal.ProposedPremium = NewMoney(
                request.ProposedPremiumAmount.Value,
                request.ProposedPremiumCurrency ?? proposal.ProposedPremium?.Currency ?? "BDT");
        }

        if (request.ProposedSumInsuredAmount.HasValue)
        {
            proposal.ProposedSumInsured = NewMoney(
                request.ProposedSumInsuredAmount.Value,
                request.ProposedSumInsuredCurrency ?? proposal.ProposedSumInsured?.Currency ?? "BDT");
        }

        if (request.SubmissionPayload is not null)
        {
            proposal.SubmissionPayload = request.SubmissionPayload;
        }

        if (request.InsurerResponsePayload is not null)
        {
            proposal.InsurerResponsePayload = request.InsurerResponsePayload;
        }

        if (request.DecisionReason is not null)
        {
            proposal.DecisionReason = request.DecisionReason;
        }

        if (request.ApprovedPolicyId is not null)
        {
            proposal.ApprovedPolicyId = request.ApprovedPolicyId;
        }

        if (request.RefundId is not null)
        {
            proposal.RefundId = request.RefundId;
        }

        if (request.CorrelationId is not null)
        {
            proposal.CorrelationId = request.CorrelationId;
        }

        if (request.ReviewedAt.HasValue)
        {
            proposal.ReviewedAt = Timestamp.FromDateTime(request.ReviewedAt.Value.ToUniversalTime());
        }

        if (request.ReviewedByUserId is not null)
        {
            proposal.ReviewedByUserId = request.ReviewedByUserId;
        }

        proposal.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);

        var updated = await _policyGateway.UpdateInsuranceProposalAsync(proposal, cancellationToken);
        return Ok(updated);
    }

    [HttpDelete("/v1/insurance-proposals/{proposalId}")]
    public async Task<IActionResult> Delete(string proposalId, CancellationToken cancellationToken = default)
    {
        await _policyGateway.DeleteInsuranceProposalAsync(proposalId, cancellationToken);
        return NoContent();
    }

    [HttpGet("/v1/orders/{orderId}/proposal")]
    public async Task<IActionResult> GetForOrder(string orderId, CancellationToken cancellationToken = default)
    {
        var proposal = await _workflowService.TryGetProposalByOrderAsync(orderId, cancellationToken);
        return proposal is null
            ? NotFound(new { error = "No insurance proposal exists for this order." })
            : Ok(proposal);
    }

    [HttpPost("/v1/orders/{orderId}/proposal")]
    public async Task<IActionResult> CreateForOrder(
        string orderId,
        [FromBody] SubmitOrderProposalRequest? request,
        CancellationToken cancellationToken = default)
    {
        try
        {
            var proposal = await _workflowService.SubmitProposalForOrderAsync(
                orderId,
                request?.InsurerId,
                request?.CorrelationId,
                request?.SubmissionPayload,
                request?.TotalPayableAmount,
                request?.TotalPayableCurrency,
                cancellationToken);

            return Created($"/v1/insurance-proposals/{proposal.ProposalId}", proposal);
        }
        catch (KeyNotFoundException ex)
        {
            return NotFound(new { error = ex.Message });
        }
        catch (InvalidOperationException ex)
        {
            return Conflict(new { error = ex.Message });
        }
    }

    [HttpPost("/v1/insurance-proposals/{proposalId}/approve")]
    public async Task<IActionResult> Approve(
        string proposalId,
        [FromBody] ProposalDecisionRequest? request,
        CancellationToken cancellationToken = default)
    {
        try
        {
            var reviewedByUserId = ResolveReviewedByUserId(request?.ReviewedByUserId);
            var result = await _workflowService.ApproveProposalAsync(
                proposalId,
                reviewedByUserId,
                request?.InsurerResponsePayload,
                request?.DecisionReason,
                cancellationToken);

            return Ok(result);
        }
        catch (KeyNotFoundException ex)
        {
            return NotFound(new { error = ex.Message });
        }
        catch (InvalidOperationException ex)
        {
            return Conflict(new { error = ex.Message });
        }
    }

    [HttpPost("/v1/insurance-proposals/{proposalId}/reject")]
    public async Task<IActionResult> Reject(
        string proposalId,
        [FromBody] ProposalDecisionRequest? request,
        CancellationToken cancellationToken = default)
    {
        try
        {
            var reviewedByUserId = ResolveReviewedByUserId(request?.ReviewedByUserId);
            var result = await _workflowService.RejectProposalAsync(
                proposalId,
                reviewedByUserId,
                request?.InsurerResponsePayload,
                request?.DecisionReason,
                cancellationToken);

            return Ok(result);
        }
        catch (KeyNotFoundException ex)
        {
            return NotFound(new { error = ex.Message });
        }
        catch (InvalidOperationException ex)
        {
            return Conflict(new { error = ex.Message });
        }
    }

    private string ResolveReviewedByUserId(string? explicitUserId)
        => FirstNonEmpty(
            explicitUserId,
            _currentUser.UserId == Guid.Empty ? null : _currentUser.UserId.ToString(),
            "polisync-system")!;

    private static bool TryParseStatus(string? rawStatus, out ProposalStatus? status)
    {
        status = null;
        if (string.IsNullOrWhiteSpace(rawStatus))
        {
            return true;
        }

        if (System.Enum.TryParse<ProposalStatus>(rawStatus, ignoreCase: true, out var parsed))
        {
            status = parsed;
            return true;
        }

        return false;
    }

    private static string? FirstNonEmpty(params string?[] values)
        => values.FirstOrDefault(value => !string.IsNullOrWhiteSpace(value));

    private static Money NewMoney(long amount, string currency)
        => new() { Amount = amount, Currency = string.IsNullOrWhiteSpace(currency) ? "BDT" : currency };
}

public sealed record CreateInsuranceProposalRequest(
    string TenantId,
    string OrderId,
    string QuotationId,
    string CustomerId,
    string InsurerId,
    string ProductId,
    string PlanId,
    long ProposedPremiumAmount,
    long ProposedSumInsuredAmount,
    string ProposedPremiumCurrency = "BDT",
    string ProposedSumInsuredCurrency = "BDT",
    string? ProposalNumber = null,
    ProposalStatus? Status = null,
    string? SubmissionPayload = null,
    string? CorrelationId = null,
    DateTime? SubmittedAt = null);

public sealed record UpdateInsuranceProposalRequest(
    ProposalStatus? Status = null,
    long? ProposedPremiumAmount = null,
    long? ProposedSumInsuredAmount = null,
    string? ProposedPremiumCurrency = null,
    string? ProposedSumInsuredCurrency = null,
    string? SubmissionPayload = null,
    string? InsurerResponsePayload = null,
    string? DecisionReason = null,
    string? ApprovedPolicyId = null,
    string? RefundId = null,
    string? CorrelationId = null,
    DateTime? ReviewedAt = null,
    string? ReviewedByUserId = null);

public sealed record SubmitOrderProposalRequest(
    string? InsurerId = null,
    string? CorrelationId = null,
    string? SubmissionPayload = null,
    long? TotalPayableAmount = null,
    string? TotalPayableCurrency = null);

public sealed record ProposalDecisionRequest(
    string? ReviewedByUserId = null,
    string? InsurerResponsePayload = null,
    string? DecisionReason = null);
