using System;
using System.Threading.Tasks;
using MediatR;
using Microsoft.AspNetCore.Mvc;
using InsuranceEngine.Underwriting.Application.Features.Commands.ApplyForQuote;
using InsuranceEngine.Underwriting.Application.Features.Commands.RecordUnderwritingDecision;
using InsuranceEngine.Underwriting.Application.Features.Queries.GetQuote;
using InsuranceEngine.Underwriting.Application.Features.Queries.ListQuotes;
using InsuranceEngine.Underwriting.Application.Features.Queries.GetUnderwritingHistory;
using InsuranceEngine.Underwriting.Domain.Enums;

namespace InsuranceEngine.Underwriting.Controllers;

[ApiController]
[Route("api/[controller]")]
public class UnderwritingController : ControllerBase
{
    private readonly IMediator _mediator;

    public UnderwritingController(IMediator mediator)
    {
        _mediator = mediator;
    }

    [HttpPost("quotes")]
    public async Task<IActionResult> ApplyForQuote([FromBody] ApplyForQuoteCommand command)
    {
        var result = await _mediator.Send(command);
        if (result.IsSuccess)
            return CreatedAtAction(nameof(GetQuote), new { id = result.Value.Id }, result.Value);
        return BadRequest(result.Error);
    }

    [HttpGet("quotes/{id}")]
    public async Task<IActionResult> GetQuote(Guid id)
    {
        var result = await _mediator.Send(new GetQuoteQuery(id));
        if (result.IsSuccess) return Ok(result.Value);
        return NotFound(result.Error);
    }

    [HttpGet("quotes")]
    public async Task<IActionResult> ListQuotes([FromQuery] Guid? beneficiaryId, [FromQuery] QuoteStatus? status, [FromQuery] int page = 1, [FromQuery] int pageSize = 20)
    {
        var result = await _mediator.Send(new ListQuotesQuery(beneficiaryId, status, page, pageSize));
        return Ok(result);
    }

    [HttpPost("decisions")]
    public async Task<IActionResult> RecordDecision([FromBody] RecordUnderwritingDecisionCommand command)
    {
        var result = await _mediator.Send(command);
        if (result.IsSuccess) return Ok(result.Value);
        return BadRequest(result.Error);
    }

    [HttpGet("quotes/{id}/history")]
    public async Task<IActionResult> GetHistory(Guid id)
    {
        var result = await _mediator.Send(new GetUnderwritingHistoryQuery(id));
        if (result.IsSuccess) return Ok(result.Value);
        return BadRequest(result.Error);
    }
}
