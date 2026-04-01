using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using InsuranceEngine.Fraud.Application.DTOs;
using InsuranceEngine.SharedKernel.DTOs;
using MediatR;
using Microsoft.AspNetCore.Mvc;

namespace InsuranceEngine.Fraud.Controllers;

[ApiController]
[Route("v1")]
public class FraudController : ControllerBase
{
    private readonly IMediator _mediator;

    public FraudController(IMediator mediator)
    {
        _mediator = mediator;
    }

    /// <summary>
    /// Check for fraud
    /// </summary>
    [HttpPost("fraud-checks")]
    public async Task<IActionResult> CheckFraud([FromBody] CheckFraudRequest request)
    {
        // Placeholder for MediatR send
        // var result = await _mediator.Send(new CheckFraudCommand(request.EntityType, request.EntityId, request.Data));
        // if (result.IsSuccess) return Ok(new CheckFraudResponse(...));
        return Ok(new CheckFraudResponse(false, 10, "LOW", new List<string>(), Guid.NewGuid().ToString()));
    }

    /// <summary>
    /// List fraud alerts
    /// </summary>
    [HttpGet("fraud-alerts")]
    public async Task<IActionResult> ListAlerts([FromQuery] int page = 1, [FromQuery] int pageSize = 20)
    {
        return Ok(new FraudAlertsListingResponse(new List<FraudAlertDto>(), 0));
    }

    /// <summary>
    /// Get fraud alert details
    /// </summary>
    [HttpGet("fraud-alerts/{id}")]
    public async Task<IActionResult> GetAlert(string id)
    {
        return Ok(new FraudAlertRetrievalResponse(new FraudAlertDto(id, "Claim", "CL-123", 80, "HIGH", new List<string> { "RULE_001" }, "OPEN", DateTime.UtcNow)));
    }

    /// <summary>
    /// Create fraud case
    /// </summary>
    [HttpPost("fraud-cases")]
    public async Task<IActionResult> CreateCase([FromBody] object request) // Placeholder request type
    {
        return Created("", new FraudCaseRetrievalResponse(new FraudCaseDto(Guid.NewGuid().ToString(), Guid.NewGuid().ToString(), null, "OPEN", "HIGH", null, DateTime.UtcNow, DateTime.UtcNow)));
    }

    /// <summary>
    /// List fraud cases
    /// </summary>
    [HttpGet("fraud-cases")]
    public async Task<IActionResult> ListCases([FromQuery] int page = 1, [FromQuery] int pageSize = 20)
    {
        return Ok(new FraudCasesListingResponse(new List<FraudCaseDto>(), 0));
    }

    /// <summary>
    /// Get fraud case details
    /// </summary>
    [HttpGet("fraud-cases/{id}")]
    public async Task<IActionResult> GetCase(string id)
    {
        return Ok(new FraudCaseRetrievalResponse(new FraudCaseDto(id, Guid.NewGuid().ToString(), "INV-001", "IN_PROGRESS", "MEDIUM", "Investigating...", DateTime.UtcNow.AddDays(-1), DateTime.UtcNow)));
    }

    /// <summary>
    /// Update fraud case
    /// </summary>
    [HttpPatch("fraud-cases/{id}")]
    public async Task<IActionResult> UpdateCase(string id, [FromBody] object request)
    {
        return Ok(new FraudCaseRetrievalResponse(new FraudCaseDto(id, Guid.NewGuid().ToString(), "INV-001", "RESOLVED", "LOW", "No fraud found.", DateTime.UtcNow.AddDays(-2), DateTime.UtcNow)));
    }

    /// <summary>
    /// List fraud rules
    /// </summary>
    [HttpGet("fraud-rules")]
    public async Task<IActionResult> ListRules([FromQuery] int page = 1, [FromQuery] int pageSize = 20)
    {
        return Ok(new FraudRulesListingResponse(new List<FraudRuleDto>(), 0));
    }

    /// <summary>
    /// Create fraud rule
    /// </summary>
    [HttpPost("fraud-rules")]
    public async Task<IActionResult> CreateRule([FromBody] object request)
    {
        return Created("", new FraudRuleRetrievalResponse(new FraudRuleDto(Guid.NewGuid().ToString(), "Large Claim Rule", "Flag claims > $1M", "Claim", "amount > 1000000", true, DateTime.UtcNow)));
    }

    /// <summary>
    /// Update fraud rule
    /// </summary>
    [HttpPatch("fraud-rules/{id}")]
    public async Task<IActionResult> UpdateRule(string id, [FromBody] object request)
    {
        return Ok(new FraudRuleOperationResponse("Rule updated successfully."));
    }

    /// <summary>
    /// Activate fraud rule
    /// </summary>
    [HttpPost("fraud-rules/{id}/activate")]
    public async Task<IActionResult> ActivateRule(string id)
    {
        return Ok(new FraudRuleOperationResponse("Rule activated successfully."));
    }

    /// <summary>
    /// Deactivate fraud rule
    /// </summary>
    [HttpPost("fraud-rules/{id}/deactivate")]
    public async Task<IActionResult> DeactivateRule(string id)
    {
        return Ok(new FraudRuleOperationResponse("Rule deactivated successfully."));
    }

    // ===================== Helpers =====================

    private ErrorDto MapError(InsuranceEngine.SharedKernel.CQRS.Error error)
    {
        return new ErrorDto(error.Code, error.Message);
    }
}
