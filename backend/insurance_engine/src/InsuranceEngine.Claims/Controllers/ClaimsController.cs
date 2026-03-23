using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using MediatR;
using Microsoft.AspNetCore.Mvc;
using InsuranceEngine.Claims.Application.Features.Commands.Claims;
using InsuranceEngine.Claims.Application.Features.Queries.Claims;
using InsuranceEngine.Claims.Application.DTOs;
using InsuranceEngine.SharedKernel.DTOs;

namespace InsuranceEngine.Claims.Controllers;

[ApiController]
[Route("v1/claims")]
public class ClaimsController : ControllerBase
{
    private readonly IMediator _mediator;

    public ClaimsController(IMediator mediator)
    {
        _mediator = mediator;
    }

    [HttpPost]
    public async Task<IActionResult> SubmitClaim([FromBody] SubmitClaimRestRequest request)
    {
        var command = new SubmitClaimCommand(
            request.PolicyId,
            request.CustomerId,
            request.Type,
            request.ClaimedAmount,
            request.IncidentDate,
            request.IncidentDescription,
            request.PlaceOfIncident,
            request.BankDetailsForPayout,
            request.Documents ?? new List<ClaimDocumentDto>()
        );

        var result = await _mediator.Send(command);

        if (result.IsSuccess)
        {
            // Value is expected to be Guid from SubmitClaimCommandHandler
            // Need to fetch the claim to get the claim number if required by the response DTO
            // For now, returning message and claim id
            var response = new ClaimSubmissionResponse(
                result.Value,
                "CLAIM-" + result.Value.ToString().Substring(0, 8).ToUpper(), // Placeholder claim number
                "Claim submitted successfully."
            );
            return Ok(response);
        }

        return BadRequest(MapError(result.Error!));
    }

    [HttpPost("{id}/documents")]
    public async Task<IActionResult> UploadDocuments(Guid id, [FromBody] List<ClaimDocumentDto> documents)
    {
        var command = new UploadClaimDocumentCommand(id, documents);
        var result = await _mediator.Send(command);

        if (result.IsSuccess)
        {
            // Documentation expects document_id, document_url, file_hash
            // For now, returning first document if available
            var doc = documents.Count > 0 ? documents[0] : new ClaimDocumentDto();
            var response = new ClaimsDocumentUploadResponse(
                doc.Id,
                doc.FileUrl,
                doc.FileHash
            );
            return Ok(response);
        }

        return BadRequest(MapError(result.Error!));
    }

    [HttpGet("{id}")]
    public async Task<IActionResult> GetClaim(Guid id)
    {
        var result = await _mediator.Send(new GetClaimByIdQuery(id));

        if (result.IsSuccess)
        {
            return Ok(new ClaimRetrievalResponse(result.Value));
        }

        return NotFound(MapError(result.Error!));
    }

    [HttpGet("/v1/users/{customerId}/claims")]
    public async Task<IActionResult> ListByCustomer(Guid customerId, [FromQuery] int page = 1, [FromQuery] int pageSize = 10)
    {
        var result = await _mediator.Send(new ListClaimsByCustomerQuery(customerId, page, pageSize));

        if (result.IsSuccess)
        {
            var claims = result.Value.Items.Select(c => new ClaimListDto
            {
                Id = c.Id,
                ClaimNumber = c.ClaimNumber,
                PolicyId = c.PolicyId,
                Type = c.Type,
                Status = c.Status,
                ClaimedAmount = c.ClaimedAmount,
                ApprovedAmount = c.ApprovedAmount,
                SubmittedAt = c.SubmittedAt
            }).ToList();
            
            return Ok(new UserClaimsListingResponse(claims, result.Value.TotalCount));
        }

        return BadRequest(MapError(result.Error!));
    }

    [HttpPost("{id}/approve")]
    public async Task<IActionResult> ApproveClaim(Guid id, [FromBody] ApproveClaimRestRequest request)
    {
        var command = new ApproveClaimCommand(
            id,
            request.ApproverId,
            request.ApproverRole,
            request.ApprovalLevel,
            request.Decision,
            request.ApprovedAmount,
            request.Notes
        );

        var result = await _mediator.Send(command);
        if (result.IsSuccess)
        {
            return Ok(new ClaimApprovalResponse("Claim approved successfully."));
        }
        return BadRequest(MapError(result.Error!));
    }

    // Helper to map Domain/Application errors to documented ErrorDto
    private ErrorDto MapError(SharedKernel.CQRS.Error error)
    {
        return new ErrorDto(error.Code, error.Message);
    }
}

public record ApproveClaimRestRequest(
    Guid ApproverId,
    string ApproverRole,
    int ApprovalLevel,
    Domain.Enums.ApprovalDecision Decision,
    long ApprovedAmount,
    string Notes
);
