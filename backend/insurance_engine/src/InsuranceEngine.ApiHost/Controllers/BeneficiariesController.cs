using MediatR;
using Microsoft.AspNetCore.Mvc;
using InsuranceEngine.Beneficiary.Application.Commands;
using InsuranceEngine.Beneficiary.Application.Queries;
using Insuretech.Beneficiary.Services.V1;

namespace InsuranceEngine.ApiHost.Controllers;

[ApiController]
[Route("v1/beneficiaries")]
public sealed class BeneficiariesController : ControllerBase
{
    private readonly IMediator _mediator;

    public BeneficiariesController(IMediator mediator)
    {
        _mediator = mediator;
    }

    [HttpGet]
    public async Task<IActionResult> List([FromQuery] string? type, [FromQuery] string? status, [FromQuery] int page = 1, [FromQuery] int pageSize = 10)
    {
        var query = new ListBeneficiariesQuery(type, status, page, pageSize);
        var response = await _mediator.Send(query);
        return Ok(response);
    }

    [HttpGet("{id}")]
    public async Task<IActionResult> Get(string id)
    {
        var query = new GetBeneficiaryQuery(id);
        var response = await _mediator.Send(query);
        
        if (response.Error != null)
        {
            return response.Error.HttpStatusCode == 404 ? NotFound(response) : BadRequest(response);
        }
        
        return Ok(response);
    }

    [HttpPost("individual")]
    public async Task<IActionResult> CreateIndividual([FromBody] CreateIndividualBeneficiaryCommand command)
    {
        var result = await _mediator.Send(command);

        return result.IsSuccess 
            ? Ok(new { beneficiary_id = result.Value, message = "Beneficiary created successfully" }) 
            : BadRequest(new { error = result.Error });
    }

    [HttpPost("business")]
    public async Task<IActionResult> CreateBusiness([FromBody] CreateBusinessBeneficiaryCommand command)
    {
        var result = await _mediator.Send(command);

        return result.IsSuccess 
            ? Ok(new { beneficiary_id = result.Value, message = "Business beneficiary created successfully" }) 
            : BadRequest(new { error = result.Error });
    }

    [HttpPut("{id}")]
    public async Task<IActionResult> Update(string id, [FromBody] UpdateBeneficiaryRequest request)
    {
        // Ensure ID matches
        request.BeneficiaryId = id;
        var response = await _mediator.Send(new UpdateBeneficiaryCommand(request));
        
        if (response.Error != null)
        {
            return response.Error.HttpStatusCode == 404 ? NotFound(response) : BadRequest(response);
        }
        
        return Ok(response);
    }

    [HttpPost("{id}/kyc")]
    public async Task<IActionResult> CompleteKYC(string id, [FromBody] CompleteKYCRequest request)
    {
        request.BeneficiaryId = id;
        var response = await _mediator.Send(new CompleteKYCCommand(request));
        
        if (response.Error != null)
        {
            return response.Error.HttpStatusCode == 404 ? NotFound(response) : BadRequest(response);
        }
        
        return Ok(response);
    }

    [HttpPost("{id}/risk-score")]
    public async Task<IActionResult> UpdateRiskScore(string id, [FromBody] UpdateRiskScoreRequest request)
    {
        request.BeneficiaryId = id;
        var response = await _mediator.Send(new UpdateRiskScoreCommand(request));
        
        if (response.Error != null)
        {
            return response.Error.HttpStatusCode == 404 ? NotFound(response) : BadRequest(response);
        }
        
        return Ok(response);
    }
}
