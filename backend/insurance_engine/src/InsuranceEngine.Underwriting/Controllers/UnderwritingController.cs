using System;
using System.Threading.Tasks;
using MediatR;
using Microsoft.AspNetCore.Mvc;
using InsuranceEngine.Underwriting.Application.Features.Commands.ApplyForQuote;
using InsuranceEngine.Underwriting.Application.Features.Commands.RecordUnderwritingDecision;
using InsuranceEngine.Underwriting.Application.Features.Queries.GetQuote;
using InsuranceEngine.Underwriting.Application.Features.Queries.ListQuotes;
using InsuranceEngine.Underwriting.Application.Features.Queries.GetUnderwritingHistory;
// using InsuranceEngine.Underwriting.Application.Features.Queries.GetUnderwritingDecision; 
// using InsuranceEngine.Underwriting.Application.Features.Queries.GetHealthDeclaration; 
using InsuranceEngine.Underwriting.Application.DTOs;
using InsuranceEngine.Underwriting.Domain.Enums;
using InsuranceEngine.SharedKernel.DTOs;

namespace InsuranceEngine.Underwriting.Controllers;

[ApiController]
[Route("v1/quotes")]
public class UnderwritingController : ControllerBase
{
    private readonly IMediator _mediator;

    public UnderwritingController(IMediator mediator)
    {
        _mediator = mediator;
    }

    /// <summary>
    /// Request premium quote
    /// </summary>
    [HttpPost]
    public async Task<IActionResult> ApplyForQuote([FromBody] ApplyForQuoteCommand command)
    {
        var result = await _mediator.Send(command);
        if (result.IsSuccess)
        {
            var r = result.Value!;
            // documentation shows: quote_id, quote_number, base_premium, total_premium, valid_until, message
            return Ok(new RequestQuoteResponse(
                r.Id,
                r.QuoteNumber,
                r.BasePremium,
                r.TotalPremium,
                r.ValidUntil,
                "Quote generated successfully."));
        }
        return BadRequest(MapError(result.Error!));
    }

    /// <summary>
    /// Get quote by ID
    /// </summary>
    [HttpGet("{id}")]
    public async Task<IActionResult> GetQuote(Guid id)
    {
        var result = await _mediator.Send(new GetQuoteQuery(id));
        if (result.IsSuccess) return Ok(new QuoteRetrievalResponse(result.Value!));
        return NotFound(MapError(result.Error!));
    }

    /// <summary>
    /// List quotes
    /// </summary>
    [HttpGet]
    public async Task<IActionResult> ListQuotes([FromQuery] Guid? beneficiaryId, [FromQuery] QuoteStatus? status, [FromQuery] int page = 1, [FromQuery] int pageSize = 20)
    {
        var result = await _mediator.Send(new ListQuotesQuery(beneficiaryId, status, page, pageSize));
        return Ok(new QuotesListingResponse(result.Items, result.TotalCount));
    }

    /// <summary>
    /// Approve underwriting (manual)
    /// </summary>
    [HttpPost("{id}")]
    public async Task<IActionResult> Approve(Guid id, [FromBody] RecordUnderwritingDecisionCommand command)
    {
        if (id != command.QuoteId) return BadRequest(new ErrorDto("VALIDATION_ERROR", "Quote ID mismatch."));

        var result = await _mediator.Send(command);
        if (result.IsSuccess)
        {
            return Ok(new UnderwritingApprovalResponse(result.Value!.Id, "Underwriting decision recorded."));
        }
        return HandleErrorResult(result.Error!);
    }

    /// <summary>
    /// Get underwriting decision
    /// </summary>
    [HttpGet("{id}/decision")]
    public async Task<IActionResult> GetDecision(Guid id)
    {
        // Placeholder until application layer is implemented
        return StatusCode(501, new ErrorDto("NOT_IMPLEMENTED", "Underwriting decision query not yet implemented."));
        /*
        var result = await _mediator.Send(new GetUnderwritingDecisionQuery(id));
        if (result.IsSuccess) return Ok(new UnderwritingDecisionRetrievalResponse(result.Value!));
        return NotFound(MapError(result.Error!));
        */
    }

    /// <summary>
    /// Get health declaration
    /// </summary>
    [HttpGet("{id}/health-declaration")]
    public async Task<IActionResult> GetHealthDeclaration(Guid id)
    {
        // Placeholder until application layer is implemented
        return StatusCode(501, new ErrorDto("NOT_IMPLEMENTED", "Health declaration query not yet implemented."));
        /*
        var result = await _mediator.Send(new GetHealthDeclarationQuery(id));
        if (result.IsSuccess) return Ok(new HealthDeclarationRetrievalResponse(result.Value!));
        return NotFound(MapError(result.Error!));
        */
    }

    /// <summary>
    /// Get underwriting history
    /// </summary>
    [HttpGet("{id}/history")]
    public async Task<IActionResult> GetHistory(Guid id)
    {
        var result = await _mediator.Send(new GetUnderwritingHistoryQuery(id));
        if (result.IsSuccess) return Ok(result.Value); // Needs a response wrapper if documented
        return BadRequest(MapError(result.Error!));
    }

    // ===================== Helpers =====================

    private ErrorDto MapError(InsuranceEngine.SharedKernel.CQRS.Error error)
    {
        return new ErrorDto(error.Code, error.Message);
    }

    private IActionResult HandleErrorResult(InsuranceEngine.SharedKernel.CQRS.Error error)
    {
        var errorDto = MapError(error);
        return error.Code switch
        {
            "NOT_FOUND" => NotFound(errorDto),
            "CONFLICT" => Conflict(errorDto),
            "VALIDATION_ERROR" => BadRequest(errorDto),
            _ => BadRequest(errorDto)
        };
    }
}
