using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using InsuranceEngine.Commission.Application.Features.Queries;
using InsuranceEngine.Commission.Application.Features.Queries.Commissions;
// using InsuranceEngine.Commission.Application.Features.Commands;
using InsuranceEngine.Commission.Application.DTOs;
using InsuranceEngine.SharedKernel.DTOs;
using MediatR;
using Microsoft.AspNetCore.Mvc;

namespace InsuranceEngine.Commission.Controllers;

[ApiController]
[Route("v1/commissions")]
public class CommissionsController : ControllerBase
{
    private readonly IMediator _mediator;

    public CommissionsController(IMediator mediator)
    {
        _mediator = mediator;
    }

    /// <summary>
    /// List commissions for recipient
    /// </summary>
    [HttpGet]
    public async Task<IActionResult> ListCommissions([FromQuery] Guid recipientId, [FromQuery] int page = 1, [FromQuery] int pageSize = 20)
    {
        var result = await _mediator.Send(new GetRecipientCommissionsQuery(recipientId));
        if (result.IsSuccess)
        {
            var commissions = result.Value!.Select(MapToDto).ToList();
            var totalAmount = new MoneyDto(0, "USD");
            return Ok(new CommissionsListingResponse(commissions, commissions.Count, totalAmount));
        }
        return BadRequest(MapError(result.Error!));
    }

    /// <summary>
    /// Get commission details
    /// </summary>
    [HttpGet("{id}")]
    public async Task<IActionResult> GetCommission(Guid id)
    {
        return StatusCode(501, new ErrorDto("NOT_IMPLEMENTED", "Get commission details query not yet implemented."));
        /*
        var result = await _mediator.Send(new GetCommissionQuery(id)); // Assuming this query exists
        if (result.IsSuccess) return Ok(new CommissionRetrievalResponse(MapToDto(result.Value!)));
        return NotFound(MapError(result.Error!));
        */
    }

    /// <summary>
    /// Calculate and record commission for policy
    /// </summary>
    [HttpPost]
    public async Task<IActionResult> CalculateCommission([FromQuery] string action, [FromBody] object command) // Using object placeholder
    {
        return StatusCode(501, new ErrorDto("NOT_IMPLEMENTED", "Calculate commission command not yet implemented."));
    }

    // ===================== Helpers =====================

    private InsuranceEngine.Commission.Application.DTOs.CommissionDto MapToDto(InsuranceEngine.Commission.Application.Features.Queries.CommissionDto c)
    {
        return new InsuranceEngine.Commission.Application.DTOs.CommissionDto(c.Id, c.PolicyId, c.PartnerId, c.AgentId, c.Type, c.Amount, c.Currency, c.Status, c.CreatedAt);
    }

    private ErrorDto MapError(InsuranceEngine.SharedKernel.CQRS.Error error)
    {
        return new ErrorDto(error.Code, error.Message);
    }
}
