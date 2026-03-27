using Microsoft.AspNetCore.Mvc;
using MediatR;
using InsuranceEngine.Underwriting.Application.Commands;
using InsuranceEngine.Underwriting.Application.Queries;
using InsuranceEngine.SharedKernel.DTOs;

namespace InsuranceEngine.ApiHost.Controllers;

[ApiController]
[Route("v1/quotes")]
public class QuotesController : ControllerBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<QuotesController> _logger;

    public QuotesController(IMediator mediator, ILogger<QuotesController> logger)
    {
        _mediator = mediator;
        _logger = logger;
    }

    [HttpPost]
    public async Task<IActionResult> Create([FromBody] CreateQuoteRequest request)
    {
        var command = new CreateQuoteCommand(
            request.BeneficiaryId,
            request.ProductId,
            request.SumAssured,
            request.TermYears,
            request.PremiumPaymentMode,
            request.ApplicantAge,
            request.ApplicantOccupation,
            request.Smoker,
            request.SelectedRiders);

        var result = await _mediator.Send(command);

        if (result.IsSuccess)
            return CreatedAtAction(nameof(Get), new { id = result.Value }, new { id = result.Value, message = "Quote created successfully" });

        return BadRequest(MapError(result.Error!));
    }

    [HttpGet("{id}")]
    public async Task<IActionResult> Get(string id)
    {
        var result = await _mediator.Send(new GetQuoteQuery(id));

        if (result.IsNotFound)
            return NotFound(MapError(result.Error!));

        if (result.IsFailure)
            return BadRequest(MapError(result.Error!));

        return Ok(MapToResponse(result.Value!));
    }

    [HttpGet]
    public async Task<IActionResult> List(
        [FromQuery] string? beneficiaryId,
        [FromQuery] string? status,
        [FromQuery] int page = 1,
        [FromQuery] int pageSize = 20)
    {
        var result = await _mediator.Send(new ListQuotesQuery(beneficiaryId, status, page, pageSize));

        if (result.IsFailure)
            return BadRequest(MapError(result.Error!));

        var items = result.Value.Items.Select(MapToResponse).ToList();
        return Ok(new QuoteListResponse(items, result.Value.TotalCount));
    }

    [HttpPost("{id}/submit")]
    public async Task<IActionResult> Submit(string id)
    {
        var result = await _mediator.Send(new SubmitQuoteForUnderwritingCommand(id));

        if (result.IsFailure)
            return BadRequest(MapError(result.Error!));

        return Ok(new { message = "Quote submitted for underwriting" });
    }

    [HttpPost("{id}/approve")]
    public async Task<IActionResult> Approve(string id)
    {
        var result = await _mediator.Send(new ApproveQuoteCommand(id));

        if (result.IsFailure)
            return BadRequest(MapError(result.Error!));

        return Ok(new { message = "Quote approved" });
    }

    [HttpPost("{id}/reject")]
    public async Task<IActionResult> Reject(string id, [FromBody] RejectQuoteRequest request)
    {
        var result = await _mediator.Send(new RejectQuoteCommand(id, request.Reason));

        if (result.IsFailure)
            return BadRequest(MapError(result.Error!));

        return Ok(new { message = "Quote rejected" });
    }

    private static ErrorDto MapError(SharedKernel.CQRS.Error error) => new(error.Code, error.Message);

    private static QuoteResponse MapToResponse(QuoteDto dto) => new(
        dto.QuoteId,
        dto.QuoteNumber,
        dto.BeneficiaryId,
        dto.ProductId,
        dto.Status,
        dto.SumAssured,
        dto.TermYears,
        dto.BasePremium,
        dto.TotalPremium,
        dto.ApplicantAge,
        dto.Smoker,
        dto.ValidUntil,
        dto.CreatedAt);
}

public record CreateQuoteRequest(
    string? BeneficiaryId,
    string ProductId,
    decimal SumAssured,
    int TermYears,
    string PremiumPaymentMode,
    int ApplicantAge,
    string? ApplicantOccupation,
    bool Smoker,
    string? SelectedRiders);

public record RejectQuoteRequest(string Reason);

public record QuoteResponse(
    string QuoteId,
    string QuoteNumber,
    string? BeneficiaryId,
    string ProductId,
    string Status,
    decimal SumAssured,
    int TermYears,
    decimal BasePremium,
    decimal TotalPremium,
    int ApplicantAge,
    bool Smoker,
    DateTime? ValidUntil,
    DateTime? CreatedAt);

public record QuoteListResponse(IReadOnlyList<QuoteResponse> Items, int TotalCount);
