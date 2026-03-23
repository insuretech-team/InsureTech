using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using MediatR;
using Microsoft.AspNetCore.Mvc;
using InsuranceEngine.Policy.Application.Features.Commands.Endorsements;
using InsuranceEngine.Policy.Application.Features.Queries.Endorsements;
using InsuranceEngine.Policy.Domain.Entities;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.SharedKernel.DTOs;

namespace InsuranceEngine.Policy.Controllers;

[ApiController]
[Route("v1/endorsements")]
public class EndorsementsController : ControllerBase
{
    private readonly IMediator _mediator;

    public EndorsementsController(IMediator mediator)
    {
        _mediator = mediator;
    }

    [HttpPost]
    public async Task<IActionResult> Submit([FromBody] SubmitEndorsementCommand command)
    {
        var result = await _mediator.Send(command);
        if (result.IsSuccess)
        {
            return Ok(new EndorsementSubmissionResponse(
                result.Value!.Id,
                result.Value.EndorsementNumber,
                "Endorsement submitted successfully."
            ));
        }
        return BadRequest(MapError(result.Error!));
    }

    [HttpPost("{id}/approve")]
    public async Task<IActionResult> Approve(Guid id, [FromBody] Guid approvedBy)
    {
        var result = await _mediator.Send(new ApproveEndorsementCommand(id, approvedBy));
        if (result.IsSuccess) return Ok(new EndorsementApprovalResponse("Endorsement approved successfully."));
        return BadRequest(MapError(result.Error!));
    }

    [HttpGet("{id}")]
    public async Task<IActionResult> Get(Guid id)
    {
        var result = await _mediator.Send(new GetEndorsementQuery(id));
        if (result.IsSuccess)
        {
            return Ok(new EndorsementRetrievalResponse(MapDto(result.Value!)));
        }
        return NotFound(MapError(result.Error!));
    }

    [HttpGet]
    public async Task<IActionResult> List([FromQuery] Guid? policyId, [FromQuery] EndorsementStatus? status, [FromQuery] int page = 1, [FromQuery] int pageSize = 10)
    {
        var result = await _mediator.Send(new ListEndorsementsQuery(policyId, status, page, pageSize));
        if (result.IsSuccess)
        {
            var items = result.Value!.Items.Select(MapDto).ToList();
            return Ok(new EndorsementsListingResponse(items, result.Value.TotalCount));
        }
        return BadRequest(MapError(result.Error!));
    }

    private EndorsementDto MapDto(Endorsement e)
    {
        return new EndorsementDto(
            e.Id,
            e.EndorsementNumber,
            e.PolicyId,
            e.Type,
            e.Reason,
            e.Status,
            e.PremiumAdjustmentAmount,
            e.PremiumAdjustmentCurrency,
            e.EffectiveDate,
            e.CreatedAt
        );
    }

    private ErrorDto MapError(SharedKernel.CQRS.Error error)
    {
        return new ErrorDto(error.Code, error.Message);
    }
}
