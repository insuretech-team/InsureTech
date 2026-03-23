using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using InsuranceEngine.Partners.Application.Features.Commands.Agents;
using InsuranceEngine.Partners.Application.Features.Commands.Partners;
using InsuranceEngine.Partners.Application.Features.Queries;
using InsuranceEngine.Partners.Application.Features.Queries.Partners;
using InsuranceEngine.Partners.Application.DTOs;
using InsuranceEngine.SharedKernel.DTOs;
using MediatR;
using Microsoft.AspNetCore.Mvc;

namespace InsuranceEngine.Partners.Controllers;

[ApiController]
[Route("v1/partners")]
public class PartnersController : ControllerBase
{
    private readonly IMediator _mediator;

    public PartnersController(IMediator mediator)
    {
        _mediator = mediator;
    }

    /// <summary>
    /// Partner Management - Create Partner
    /// </summary>
    [HttpPost]
    public async Task<IActionResult> CreatePartner([FromBody] CreatePartnerCommand command)
    {
        var result = await _mediator.Send(command);
        if (result.IsSuccess)
        {
            var partnerId = result.Value;
            var partnerResult = await _mediator.Send(new GetPartnerQuery(partnerId));
            var partner = MapToDto(partnerResult.Value!);
            return CreatedAtAction(nameof(GetPartner), new { id = partnerId }, new PartnerCreationResponse(partnerId, partner));
        }
        return BadRequest(MapError(result.Error!));
    }

    /// <summary>
    /// List partners
    /// </summary>
    [HttpGet]
    public async Task<IActionResult> ListPartners([FromQuery] int page = 1, [FromQuery] int pageSize = 20)
    {
        var result = await _mediator.Send(new ListPartnersQuery());
        if (result.IsSuccess)
        {
            var dtos = result.Value!.Select(MapToDto).ToList();
            return Ok(new PartnersListingResponse(dtos, dtos.Count));
        }
        return BadRequest(MapError(result.Error!));
    }

    /// <summary>
    /// Get partner by ID
    /// </summary>
    [HttpGet("{id}")]
    public async Task<IActionResult> GetPartner(Guid id)
    {
        var result = await _mediator.Send(new GetPartnerQuery(id));
        if (result.IsSuccess) return Ok(new PartnerRetrievalResponse(MapToDto(result.Value!)));
        return NotFound(MapError(result.Error!));
    }

    /// <summary>
    /// Update partner
    /// </summary>
    [HttpPatch("{id}")]
    public async Task<IActionResult> UpdatePartner(Guid id, [FromBody] object command) // Using object placeholder
    {
        // Placeholder until application layer is implemented
        return StatusCode(501, new ErrorDto("NOT_IMPLEMENTED", "Update partner command not yet implemented."));
        /*
        if (id != command.PartnerId) return BadRequest(new ErrorDto("VALIDATION_ERROR", "Partner ID mismatch."));
        var result = await _mediator.Send(command);
        if (result.IsSuccess) return Ok(new PartnerUpdateResponse("Partner information updated."));
        return HandleErrorResult(result.Error!);
        */
    }

    /// <summary>
    /// Delete partner
    /// </summary>
    [HttpDelete("{id}")]
    public async Task<IActionResult> DeletePartner(Guid id)
    {
        return StatusCode(501, new ErrorDto("NOT_IMPLEMENTED", "Delete partner command not yet implemented."));
    }

    /// <summary>
    /// Update partner status
    /// </summary>
    [HttpPost("{id}/status")]
    public async Task<IActionResult> UpdateStatus(Guid id, [FromBody] object command) // Using object placeholder
    {
        // Placeholder until application layer is implemented
        return StatusCode(501, new ErrorDto("NOT_IMPLEMENTED", "Update partner status command not yet implemented."));
        /*
        if (id != command.PartnerId) return BadRequest(new ErrorDto("VALIDATION_ERROR", "Partner ID mismatch."));
        var result = await _mediator.Send(command);
        if (result.IsSuccess) return Ok(new PartnerStatusUpdateResponse("Partner status updated."));
        return HandleErrorResult(result.Error!);
        */
    }

    /// <summary>
    /// Verify partner
    /// </summary>
    [HttpPost("{id}:verify")]
    public async Task<IActionResult> VerifyPartner(Guid id)
    {
        return StatusCode(501, new ErrorDto("NOT_IMPLEMENTED", "Verify partner command not yet implemented."));
    }

    /// <summary>
    /// Get partner credentials
    /// </summary>
    [HttpGet("{id}/credentials")]
    public async Task<IActionResult> GetCredentials(Guid id)
    {
        return StatusCode(501, new ErrorDto("NOT_IMPLEMENTED", "Get partner credentials query not yet implemented."));
    }

    /// <summary>
    /// Rotate partner credentials
    /// </summary>
    [HttpPost("{id}/credentials:rotate")]
    public async Task<IActionResult> RotateCredentials(Guid id)
    {
        return StatusCode(501, new ErrorDto("NOT_IMPLEMENTED", "Rotate partner credentials command not yet implemented."));
    }

    // ===================== Helpers =====================

    private InsuranceEngine.Partners.Application.DTOs.PartnerDto MapToDto(InsuranceEngine.Partners.Application.Features.Queries.PartnerDto p)
    {
        return new InsuranceEngine.Partners.Application.DTOs.PartnerDto(p.Id, p.Name, p.Code, p.Email, p.Phone, p.Address, p.Status);
    }

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
            "VALIDATION_ERROR" => BadRequest(errorDto),
            "CONFLICT" => Conflict(errorDto),
            _ => BadRequest(errorDto)
        };
    }
}
