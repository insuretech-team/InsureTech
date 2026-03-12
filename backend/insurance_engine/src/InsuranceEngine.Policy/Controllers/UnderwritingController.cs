using System;
using System.Threading.Tasks;
using MediatR;
using Microsoft.AspNetCore.Mvc;
using InsuranceEngine.Policy.Application.Features.Commands.ApplyForQuote;
using InsuranceEngine.Policy.Application.Features.Commands.RecordUnderwritingDecision;
using InsuranceEngine.Policy.Application.Features.Commands.HealthDeclarations;
using InsuranceEngine.Policy.Application.Features.Commands.UpdateQuoteStatus;
using InsuranceEngine.Policy.Application.Features.Queries.GetQuote;
using InsuranceEngine.Policy.Application.Features.Queries.ListQuotes;
using InsuranceEngine.Policy.Application.Features.Queries.GetUnderwritingHistory;
using InsuranceEngine.Policy.Application.Features.Queries.HealthDeclarations;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.Policy.Domain.Enums;

namespace InsuranceEngine.Policy.Controllers;

[ApiController]
[Route("api/[controller]")]
public class UnderwritingController : ControllerBase
{
    private readonly IMediator _mediator;

    public UnderwritingController(IMediator mediator) => _mediator = mediator;

    // ===================== Quotes =====================

    /// <summary>
    /// Apply for a new insurance quote (initiates underwriting).
    /// </summary>
    [HttpPost("quotes")]
    public async Task<IActionResult> ApplyForQuote([FromBody] ApplyForQuoteCommand command)
    {
        var result = await _mediator.Send(command);
        return result.Match<IActionResult>(
            onSuccess: quote => CreatedAtAction(nameof(GetQuote), new { id = quote.Id },
                new { quote.Id, quote.QuoteNumber, message = "Quote created successfully." }),
            onFailure: error => HandleErrorResult(error)
        );
    }

    /// <summary>
    /// Get quote by ID.
    /// </summary>
    [HttpGet("quotes/{id}")]
    public async Task<IActionResult> GetQuote(Guid id)
    {
        var result = await _mediator.Send(new GetQuoteQuery(id));
        return result != null ? Ok(result) : NotFound();
    }

    /// <summary>
    /// List quotes with optional filters.
    /// </summary>
    [HttpGet("quotes")]
    public async Task<ActionResult<PaginatedResponse<QuoteDto>>> ListQuotes(
        [FromQuery] Guid? beneficiaryId,
        [FromQuery] QuoteStatus? status,
        [FromQuery] int page = 1,
        [FromQuery] int pageSize = 20)
    {
        return Ok(await _mediator.Send(new ListQuotesQuery(beneficiaryId, status, page, pageSize)));
    }

    /// <summary>
    /// Update quote status (expire, cancel, convert).
    /// </summary>
    [HttpPatch("quotes/{id}/status")]
    public async Task<IActionResult> UpdateQuoteStatus(Guid id, [FromBody] UpdateQuoteStatusRequest request)
    {
        var result = await _mediator.Send(new UpdateQuoteStatusCommand(id, request.NewStatus, request.ConvertedPolicyId));
        return HandleResult(result, "Quote status updated successfully.");
    }

    // ===================== Health Declarations =====================

    /// <summary>
    /// Submit a health declaration for an existing quote.
    /// </summary>
    [HttpPost("quotes/{quoteId}/health-declarations")]
    public async Task<IActionResult> SubmitHealthDeclaration(Guid quoteId, [FromBody] UnderwritingHealthDeclarationDto dto)
    {
        var result = await _mediator.Send(new SubmitHealthDeclarationCommand(quoteId, dto));
        return result.Match<IActionResult>(
            onSuccess: id => Created($"api/underwriting/health-declarations/{id}",
                new { declarationId = id, message = "Health declaration submitted successfully." }),
            onFailure: error => HandleErrorResult(error)
        );
    }

    /// <summary>
    /// Get health declaration by ID.
    /// </summary>
    [HttpGet("health-declarations/{id}")]
    public async Task<IActionResult> GetHealthDeclaration(Guid id)
    {
        var result = await _mediator.Send(new GetHealthDeclarationQuery(id));
        return result != null ? Ok(result) : NotFound();
    }

    /// <summary>
    /// Get health declaration by quote ID.
    /// </summary>
    [HttpGet("quotes/{quoteId}/health-declarations")]
    public async Task<IActionResult> GetHealthDeclarationByQuote(Guid quoteId)
    {
        var result = await _mediator.Send(new GetHealthDeclarationByQuoteQuery(quoteId));
        return result != null ? Ok(result) : NotFound();
    }

    // ===================== Underwriting Decisions =====================

    /// <summary>
    /// Record an underwriting decision for a quote.
    /// </summary>
    [HttpPost("quotes/{quoteId}/decisions")]
    public async Task<IActionResult> RecordDecision(Guid quoteId, [FromBody] RecordDecisionRequest request)
    {
        var command = new RecordUnderwritingDecisionCommand(
            quoteId,
            request.Decision,
            request.Method,
            request.RiskScore,
            request.RiskLevel,
            request.RiskFactors,
            request.Reason,
            request.Conditions,
            request.IsPremiumAdjusted,
            request.AdjustedPremiumAmount,
            request.UnderwriterId,
            request.UnderwriterComments
        );
        var result = await _mediator.Send(command);
        return result.Match<IActionResult>(
            onSuccess: id => Created($"api/underwriting/decisions/{id}",
                new { decisionId = id, message = "Underwriting decision recorded successfully." }),
            onFailure: error => HandleErrorResult(error)
        );
    }

    /// <summary>
    /// Get underwriting decision history for a quote.
    /// </summary>
    [HttpGet("quotes/{quoteId}/decisions")]
    public async Task<IActionResult> GetDecisionHistory(Guid quoteId)
    {
        var result = await _mediator.Send(new GetUnderwritingHistoryQuery(quoteId));
        return Ok(result);
    }

    // ===================== Helpers =====================

    private IActionResult HandleResult(SharedKernel.CQRS.Result result, string successMessage)
    {
        if (result.IsSuccess) return Ok(new { message = successMessage });
        return HandleErrorResult(result.Error!);
    }

    private IActionResult HandleErrorResult(SharedKernel.CQRS.Error error)
    {
        return error.Code switch
        {
            "NOT_FOUND" => NotFound(new { error.Code, error.Message }),
            "INVALID_STATE_TRANSITION" => Conflict(new { error.Code, error.Message }),
            "DUPLICATE" => Conflict(new { error.Code, error.Message }),
            "VALIDATION_ERROR" => BadRequest(new { error.Code, error.Message }),
            _ => BadRequest(new { error.Code, error.Message })
        };
    }
}

// --- Request DTOs ---
public record UpdateQuoteStatusRequest(string NewStatus, Guid? ConvertedPolicyId = null);

public record RecordDecisionRequest(
    DecisionType Decision,
    DecisionMethod Method,
    decimal RiskScore,
    RiskLevel RiskLevel,
    System.Collections.Generic.List<string>? RiskFactors,
    string? Reason,
    System.Collections.Generic.List<string>? Conditions,
    bool IsPremiumAdjusted,
    long AdjustedPremiumAmount,
    Guid? UnderwriterId,
    string? UnderwriterComments
);
