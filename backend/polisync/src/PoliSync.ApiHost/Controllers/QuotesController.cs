using MediatR;
using Microsoft.AspNetCore.Mvc;
using PoliSync.Underwriting.Application.Commands;
using PoliSync.Underwriting.Application.Queries;
using PoliSync.Underwriting.Domain;

namespace PoliSync.ApiHost.Controllers;

[ApiController]
[Route("v1/quotes")]
public class QuotesController : ControllerBase
{
    private readonly IMediator _mediator;

    public QuotesController(IMediator mediator) => _mediator = mediator;

    [HttpPost]
    public async Task<IActionResult> RequestQuote(RequestQuoteDto dto)
    {
        var cmd = new RequestQuoteCommand(
            dto.BeneficiaryId,
            dto.InsurerProductId,
            dto.SumAssured,
            dto.TermYears,
            dto.PremiumPaymentMode,
            dto.ApplicantAgeDays,
            dto.IsSmoker,
            dto.BasePremiumAmount,
            dto.RiderPremiumAmount,
            dto.TaxAmount
        );

        var result = await _mediator.Send(cmd);
        return result.IsSuccess ? Ok(new { quote_id = result.Value }) : BadRequest(result.Error);
    }

    [HttpGet("{id}")]
    public async Task<IActionResult> Get(Guid id)
    {
        var result = await _mediator.Send(new GetQuoteQuery(id));
        if (!result.IsSuccess) return BadRequest(result.Error);
        if (result.Value is null) return NotFound();

        return Ok(MapQuoteDto(result.Value));
    }

    [HttpGet]
    public async Task<IActionResult> List([FromQuery] Guid? beneficiaryId, [FromQuery] QuoteStatus? status, [FromQuery] int page = 1, [FromQuery] int pageSize = 20)
    {
        var result = await _mediator.Send(new ListQuotesQuery(beneficiaryId, status, page, pageSize));
        return Ok(new
        {
            items = result.Value.Items.Select(MapQuoteDto),
            total_count = result.Value.TotalCount,
            page = result.Value.Page,
            page_size = result.Value.PageSize
        });
    }

    [HttpPost("{id}/health-declaration")]
    public async Task<IActionResult> SubmitHealthDeclaration(Guid id, SubmitHealthDeclarationDto dto)
    {
        var cmd = new SubmitHealthDeclarationCommand(
            id,
            dto.HeightCm,
            dto.WeightKg,
            dto.IsSmoker,
            dto.ConsumesAlcohol,
            dto.HasPreExistingConditions,
            dto.ConditionDetails,
            dto.HasFamilyHistoryOfCriticalIllness,
            dto.OccupationRiskLevel
        );

        var result = await _mediator.Send(cmd);
        return result.IsSuccess ? Ok(new { declaration_id = result.Value }) : BadRequest(result.Error);
    }

    [HttpGet("{id}/health-declaration")]
    public async Task<IActionResult> GetHealthDeclaration(Guid id)
    {
        var result = await _mediator.Send(new GetHealthDeclarationQuery(id));
        if (!result.IsSuccess) return BadRequest(result.Error);
        if (result.Value is null) return NotFound();

        return Ok(result.Value);
    }

    [HttpPost("{id}:approve")]
    public async Task<IActionResult> Approve(Guid id, ApproveDto dto)
    {
        var cmd = new ApproveUnderwritingCommand(
            id,
            dto.RiskScore,
            dto.RiskLevel,
            dto.Reason,
            dto.Conditions,
            dto.RiskFactors,
            dto.AdjustedPremiumAmount,
            dto.UnderwriterId
        );

        var result = await _mediator.Send(cmd);
        return result.IsSuccess ? Ok() : BadRequest(result.Error);
    }

    [HttpPost("{id}:reject")]
    public async Task<IActionResult> Reject(Guid id, RejectDto dto)
    {
        var cmd = new RejectUnderwritingCommand(
            id,
            dto.RiskScore,
            dto.RiskLevel,
            dto.Reason,
            dto.RiskFactors,
            dto.UnderwriterId
        );

        var result = await _mediator.Send(cmd);
        return result.IsSuccess ? Ok() : BadRequest(result.Error);
    }

    [HttpPost("{id}:convert")]
    public async Task<IActionResult> ConvertToPolicy(Guid id)
    {
        var cmd = new ConvertQuoteToPolicyCommand(id);
        var result = await _mediator.Send(cmd);
        return result.IsSuccess ? Ok() : BadRequest(result.Error);
    }

    [HttpGet("{id}/decision")]
    public async Task<IActionResult> GetDecision(Guid id)
    {
        var result = await _mediator.Send(new GetUnderwritingDecisionQuery(id));
        if (!result.IsSuccess) return BadRequest(result.Error);
        if (result.Value is null) return NotFound();

        return Ok(result.Value);
    }

    private static object MapQuoteDto(Quote q) => new
    {
        quote_id = q.QuoteId,
        quote_number = q.QuoteNumber,
        beneficiary_id = q.BeneficiaryId,
        insurer_product_id = q.InsurerProductId,
        status = q.Status.ToString(),
        sum_assured = q.SumAssured,
        currency = q.Currency,
        term_years = q.TermYears,
        premium_payment_mode = q.PremiumPaymentMode,
        base_premium_amount = q.BasePremiumAmount,
        rider_premium_amount = q.RiderPremiumAmount,
        tax_amount = q.TaxAmount,
        total_premium_amount = q.TotalPremiumAmount,
        applicant_age_days = q.ApplicantAgeDays,
        is_smoker = q.IsSmoker,
        valid_until = q.ValidUntil,
        has_health_declaration = q.HealthDeclaration != null,
        has_decision = q.Decision != null
    };
}

public record RequestQuoteDto(
    Guid BeneficiaryId,
    Guid InsurerProductId,
    long SumAssured,
    int TermYears,
    string PremiumPaymentMode,
    int ApplicantAgeDays,
    bool IsSmoker,
    long BasePremiumAmount,
    long RiderPremiumAmount,
    long TaxAmount
);

public record SubmitHealthDeclarationDto(
    float HeightCm,
    float WeightKg,
    bool IsSmoker,
    bool ConsumesAlcohol,
    bool HasPreExistingConditions,
    string? ConditionDetails,
    bool HasFamilyHistoryOfCriticalIllness,
    string OccupationRiskLevel
);

public record ApproveDto(
    float RiskScore,
    RiskLevel RiskLevel,
    string? Reason,
    string? Conditions,
    string? RiskFactors,
    long? AdjustedPremiumAmount,
    Guid? UnderwriterId
);

public record RejectDto(
    float RiskScore,
    RiskLevel RiskLevel,
    string Reason,
    string? RiskFactors,
    Guid? UnderwriterId
);
