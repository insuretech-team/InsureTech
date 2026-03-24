using Microsoft.AspNetCore.Mvc;
using MediatR;
using InsuranceEngine.Beneficiary.Application.Features.Commands;
using InsuranceEngine.Beneficiary.Application.Features.Queries;
using InsuranceEngine.Beneficiary.Application.DTOs;
using InsuranceEngine.Underwriting.Application.Features.Queries.ListQuotes;
using InsuranceEngine.SharedKernel.DTOs;
using System.Linq;

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
        if (result.IsSuccess)
        {
            var response = new IndividualBeneficiaryCreationResponse(
                result.Value!.Id,
                result.Value.Code,
                "Individual beneficiary created successfully."
            );
            return CreatedAtAction(nameof(Get), new { id = result.Value.Id }, response);
        }
        return BadRequest(MapError(result.Error!));
    }

    [HttpPost("business")]
    public async Task<IActionResult> CreateBusiness([FromBody] CreateBusinessBeneficiaryCommand command)
    {
        var result = await _mediator.Send(command);
        if (result.IsSuccess)
        {
            var response = new BusinessBeneficiaryCreationResponse(
                result.Value!.Id,
                result.Value.Code,
                "Business beneficiary created successfully."
            );
            return CreatedAtAction(nameof(Get), new { id = result.Value.Id }, response);
        }
        return BadRequest(MapError(result.Error!));
    }

    [HttpGet("{id}")]
    public async Task<IActionResult> Get(Guid id)
    {
        var result = await _mediator.Send(new GetBeneficiaryQuery(id));
        if (result.IsSuccess)
        {
            var response = new BeneficiaryRetrievalResponse(
                result.Value!,
                result.Value!.Individual,
                result.Value!.Business
            );
            return Ok(response);
        }
        return NotFound(MapError(result.Error!));
    }

    [HttpPatch("{id}")]
    public async Task<IActionResult> Update(Guid id, [FromBody] UpdateBeneficiaryCommand command)
    {
        if (id != command.BeneficiaryId) return BadRequest(new ErrorDto("INVALID_ID", "ID in path does not match ID in body"));
        var result = await _mediator.Send(command);
        if (result.IsSuccess)
        {
            return Ok(new BeneficiaryUpdateResponse("Beneficiary updated successfully."));
        }
        return BadRequest(MapError(result.Error!));
    }

    //[HttpPost("{id}/kyc")]
    //public async Task<IActionResult> CompleteKyc(Guid id, [FromBody] CompleteKYCCommand command)
    //{
    //    if (id != command.BeneficiaryId) return BadRequest(new ErrorDto("INVALID_ID", "ID in path does not match ID in body"));
    //    var result = await _mediator.Send(command);
    //    if (result.IsSuccess)
    //    {
    //        return Ok(new KYCCompletionResponse(command.Status, "KYC completed successfully"));
    //    }
    //    return BadRequest(MapError(result.Error!));
    //}

    //[HttpPost("{id}/risk-score")]
    //public async Task<IActionResult> UpdateRiskScore(Guid id, [FromBody] UpdateRiskScoreCommand command)
    //{
    //    if (id != command.BeneficiaryId) return BadRequest(new ErrorDto("INVALID_ID", "ID in path does not match ID in body"));
    //    var result = await _mediator.Send(command);
    //    if (result.IsSuccess)
    //    {
    //        return Ok(new RiskScoreUpdateResponse("Risk score updated successfully"));
    //    }
    //    return BadRequest(MapError(result.Error!));
    //}

    //[HttpGet]
    //public async Task<IActionResult> List([FromQuery] string? type, [FromQuery] string? status, [FromQuery] int page = 1, [FromQuery] int pageSize = 10)
    //{
    //    var result = await _mediator.Send(new ListBeneficiariesQuery(type, status, pageSize, page));
    //    if (result.IsSuccess)
    //    {
    //        return Ok(new BeneficiariesListingResponse(result.Value.Items, result.Value.TotalCount));
    //    }
    //    return BadRequest(MapError(result.Error!));
    //}

    //[HttpGet("{id}/quotes")]
    //public async Task<IActionResult> GetQuotes(Guid id, [FromQuery] int page = 1, [FromQuery] int pageSize = 10)
    //{
    //    var result = await _mediator.Send(new ListQuotesQuery(id, null, page, pageSize));
        
    //    var quotes = result.Items.Select(q => new BeneficiaryQuoteDto(
    //        q.Id,
    //        q.QuoteNumber,
    //        q.Status.ToString(),
    //        q.TotalPremium,
    //        q.ValidUntil
    //    )).ToList();
        
    //    return Ok(new QuotesListingResponse(quotes, result.TotalCount));
    //}

    // Helper to map Domain/Application errors to documented ErrorDto
    private ErrorDto MapError(SharedKernel.CQRS.Error error)
    {
        return new ErrorDto(error.Code, error.Message);
    }
}
