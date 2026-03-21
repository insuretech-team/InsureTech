using Microsoft.AspNetCore.Mvc;
using MediatR;
using InsuranceEngine.Policy.Application.Features.Commands.Beneficiaries;
using InsuranceEngine.Policy.Application.Features.Queries.Beneficiaries;
using InsuranceEngine.Policy.Application.DTOs;

namespace InsuranceEngine.Policy.Controllers;

[ApiController]
[Route("api/v1/beneficiaries")]
public class BeneficiariesController : ControllerBase
{
    private readonly IMediator _mediator;

    public BeneficiariesController(IMediator mediator)
    {
        _mediator = mediator;
    }

    [HttpPost("individual")]
    public async Task<IActionResult> CreateIndividual([FromBody] CreateIndividualBeneficiaryCommand command)
    {
        var result = await _mediator.Send(command);
        return result.IsSuccess ? CreatedAtAction(nameof(Get), new { id = result.Value!.Id }, result.Value) : BadRequest(result.Error);
    }

    [HttpPost("business")]
    public async Task<IActionResult> CreateBusiness([FromBody] CreateBusinessBeneficiaryCommand command)
    {
        var result = await _mediator.Send(command);
        return result.IsSuccess ? CreatedAtAction(nameof(Get), new { id = result.Value!.Id }, result.Value) : BadRequest(result.Error);
    }

    [HttpGet("{id}")]
    public async Task<IActionResult> Get(Guid id)
    {
        var result = await _mediator.Send(new GetBeneficiaryQuery(id));
        return result.IsSuccess ? Ok(result.Value) : NotFound(result.Error);
    }

    [HttpPatch("{id}")]
    public async Task<IActionResult> Update(Guid id, [FromBody] UpdateBeneficiaryCommand command)
    {
        if (id != command.BeneficiaryId) return BadRequest();
        var result = await _mediator.Send(command);
        return result.IsSuccess ? NoContent() : BadRequest(result.Error);
    }

    [HttpPost("{id}/kyc")]
    public async Task<IActionResult> CompleteKyc(Guid id, [FromBody] CompleteKYCCommand command)
    {
        if (id != command.BeneficiaryId) return BadRequest();
        var result = await _mediator.Send(command);
        return result.IsSuccess ? Ok(new { message = "KYC completed successfully" }) : BadRequest(result.Error);
    }

    [HttpPost("{id}/risk-score")]
    public async Task<IActionResult> UpdateRiskScore(Guid id, [FromBody] UpdateRiskScoreCommand command)
    {
        if (id != command.BeneficiaryId) return BadRequest();
        var result = await _mediator.Send(command);
        return result.IsSuccess ? Ok(new { message = "Risk score updated successfully" }) : BadRequest(result.Error);
    }

    [HttpGet]
    public async Task<IActionResult> List([FromQuery] string? type, [FromQuery] string? status, [FromQuery] int page = 1, [FromQuery] int pageSize = 10)
    {
        var result = await _mediator.Send(new ListBeneficiariesQuery(type, status, pageSize, page));
        return Ok(result.Value);
    }



    [HttpGet("{id}/audit-trail")]
    public async Task<IActionResult> GetAuditTrail(Guid id)
    {
        return await Task.FromResult(Ok(new { beneficiaryId = id, audits = new List<object>() }));
    }

    [HttpGet("{id}/documents")]
    public async Task<IActionResult> GetDocuments(Guid id)
    {
        return await Task.FromResult(Ok(new { beneficiaryId = id, documents = new List<object>() }));
    }

    [HttpGet("{id}/media")]
    public async Task<IActionResult> GetMedia(Guid id)
    {
        return await Task.FromResult(Ok(new { beneficiaryId = id, media = new List<object>() }));
    }

    [HttpGet("{id}/workflow-history")]
    public async Task<IActionResult> GetWorkflowHistory(Guid id)
    {
        return await Task.FromResult(Ok(new { beneficiaryId = id, history = new List<object>() }));
    }

    [HttpGet("{id}/commission-statement")]
    public async Task<IActionResult> GetCommissionStatement(Guid id)
    {
        return await Task.FromResult(Ok(new { beneficiaryId = id, statement = new { totalCommission = 0 } }));
    }
}
