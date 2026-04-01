using Microsoft.AspNetCore.Mvc;
using MediatR;
using InsuranceEngine.Beneficiary.Application.Features.Commands;
using InsuranceEngine.Beneficiary.Application.Features.Queries;
using InsuranceEngine.Beneficiary.Application.DTOs;
using InsuranceEngine.SharedKernel.DTOs;

namespace InsuranceEngine.Beneficiary.Controllers;

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

    private ErrorDto MapError(SharedKernel.CQRS.Error error)
    {
        return new ErrorDto(error.Code, error.Message);
    }
}
