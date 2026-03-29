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
        var command = new RequestQuoteCommand(
            BeneficiaryId: request.BeneficiaryId ?? "",
            InsurerProductId: request.ProductId,
            SumAssured: (long)(request.SumAssured * 100),
            TermYears: request.TermYears,
            PremiumPaymentMode: request.PremiumPaymentMode,
            RiderCodes: request.SelectedRiders?.Split(',').ToList(),
            ApplicantAge: request.ApplicantAge,
            Smoker: request.Smoker);

        var result = await _mediator.Send(command);

        if (string.IsNullOrEmpty(result.Error?.Code))
            return CreatedAtAction(nameof(Get), new { id = result.QuoteId }, new { id = result.QuoteId, message = "Quote created successfully" });

        return BadRequest(new { code = result.Error.Code, message = result.Error.Message });
    }

    [HttpGet("{id}")]
    public async Task<IActionResult> Get(string id)
    {
        var response = await _mediator.Send(new GetQuoteQuery(id));

        if (!string.IsNullOrEmpty(response.Error?.Code))
        {
            if (response.Error.Code == "NOT_FOUND") return NotFound(new { code = response.Error.Code, message = response.Error.Message });
            return BadRequest(new { code = response.Error.Code, message = response.Error.Message });
        }

        return Ok(MapToResponse(response.Quote));
    }

    [HttpGet]
    public async Task<IActionResult> List(
        [FromQuery] string? beneficiaryId,
        [FromQuery] string? status,
        [FromQuery] int page = 1,
        [FromQuery] int pageSize = 20)
    {
        var response = await _mediator.Send(new ListQuotesQuery(beneficiaryId, status, page, pageSize));

        var items = response.Quotes.Select(MapToResponse).ToList();
        return Ok(new QuoteListResponse(items, response.TotalCount));
    }

    [HttpPost("{id}/submit")]
    public async Task<IActionResult> Submit(string id)
    {
        var command = new SubmitHealthDeclarationCommand(
            QuoteId: id,
            HeightCm: 170, // Default placeholders for simplified REST API
            WeightKg: "70",
            HasPreExistingConditions: false,
            PreExistingConditions: null,
            Smoker: false,
            AlcoholConsumer: false,
            OccupationRiskLevel: "LOW");

        var result = await _mediator.Send(command);

        if (string.IsNullOrEmpty(result.Error?.Code))
            return Ok(new { message = "Quote submitted for underwriting" });

        return BadRequest(new { code = result.Error.Code, message = result.Error.Message });
    }

    [HttpPost("{id}/approve")]
    public async Task<IActionResult> Approve(string id)
    {
        var command = new ApproveUnderwritingCommand(
            QuoteId: id,
            UnderwriterId: Guid.NewGuid().ToString(), // Placeholder
            RiskLevel: "STANDARD",
            PremiumAdjusted: false,
            AdjustedPremium: null,
            Comments: "Approved via REST API");

        var result = await _mediator.Send(command);

        if (string.IsNullOrEmpty(result.Error?.Code))
            return Ok(new { message = "Quote approved" });

        return BadRequest(new { code = result.Error.Code, message = result.Error.Message });
    }

    [HttpPost("{id}/reject")]
    public async Task<IActionResult> Reject(string id, [FromBody] RejectQuoteRequest request)
    {
        var command = new RejectUnderwritingCommand(
            QuoteId: id,
            UnderwriterId: Guid.NewGuid().ToString(), // Placeholder
            Reason: request.Reason,
            RiskLevel: null,
            Comments: "Rejected via REST API");

        var result = await _mediator.Send(command);

        if (string.IsNullOrEmpty(result.Error?.Code))
            return Ok(new { message = "Quote rejected" });

        return BadRequest(new { code = result.Error.Code, message = result.Error.Message });
    }

    private static ErrorDto MapError(SharedKernel.CQRS.Error error) => new(error.Code, error.Message);

    private static QuoteResponse MapToResponse(Insuretech.Underwriting.Entity.V1.Quote q) => new(
        q.Id,
        q.QuoteNumber,
        q.BeneficiaryId,
        q.InsurerProductId,
        q.Status.ToString(),
        (decimal)q.SumAssured.Amount / 100m,
        q.TermYears,
        (decimal)q.BasePremium.Amount / 100m,
        (decimal)q.TotalPremium.Amount / 100m,
        q.ApplicantAge,
        q.Smoker,
        q.ValidUntil?.ToDateTime(),
        q.AuditInfo?.CreatedAt?.ToDateTime());
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
