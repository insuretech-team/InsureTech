using Google.Protobuf.WellKnownTypes;
using Insuretech.Common.V1;
using Insuretech.Insurance.Services.V1;
using MediatR;
using Microsoft.AspNetCore.Mvc;
using PoliSync.Claims.Infrastructure;
using PoliSync.Infrastructure.Clients;
using PoliSync.Infrastructure.GrpcClients;
using PoliSync.Claims.Application.Commands;
using PoliSync.Claims.Application.Queries;
using PoliSync.SharedKernel.Auth;
using Insuretech.Claims.Entity.V1;
using System.Security.Cryptography;
using System.Text;

namespace PoliSync.ApiHost.Controllers;

/// <summary>
/// HTTP companion for the Claims gRPC service.
/// The InScore gateway reverse-proxies /v1/claims/* to this controller on port 50211.
/// BUG-005 FIX: Created new controller with explicit /v1/claims routes.
/// </summary>
[ApiController]
public sealed class ClaimsController : ControllerBase
{
    private readonly IMediator _mediator;
    private readonly IClaimDataGateway _claimDataGateway;
    private readonly InsuranceServiceClient _insuranceClient;
    private readonly DocgenGrpcClient _docgenClient;
    private readonly ICurrentUser _currentUser;
    private readonly ILogger<ClaimsController> _logger;

    public ClaimsController(
        IMediator mediator,
        IClaimDataGateway claimDataGateway,
        InsuranceServiceClient insuranceClient,
        DocgenGrpcClient docgenClient,
        ICurrentUser currentUser,
        ILogger<ClaimsController> logger)
    {
        _mediator = mediator;
        _claimDataGateway = claimDataGateway;
        _insuranceClient = insuranceClient;
        _docgenClient = docgenClient;
        _currentUser = currentUser;
        _logger = logger;
    }

    /// <summary>List claims — B2C users see only their own claims.</summary>
    [HttpGet("/v1/claims")]
    public async Task<IActionResult> ListClaims(
        [FromQuery(Name = "policy_id")] string? policyId = null,
        [FromQuery(Name = "customer_id")] string? customerId = null,
        [FromQuery] int page = 1,
        [FromQuery(Name = "page_size")] int pageSize = 20,
        CancellationToken cancellationToken = default)
    {
        var effectiveCustomerId = customerId
            ?? (_currentUser.UserId != Guid.Empty ? _currentUser.UserId.ToString() : string.Empty);

        var query = new ListClaimsQuery(effectiveCustomerId, policyId ?? string.Empty, page, pageSize);
        var result = await _mediator.Send(query, cancellationToken);

        if (!result.IsSuccess)
            return StatusCode(500, new { success = false, error = new { message = result.Error } });

        var claims = result.Value?.Claims ?? [];
        var totalCount = result.Value?.TotalCount ?? 0;

        return Ok(new
        {
            success = true,
            data = new
            {
                claims,
                total_count = totalCount,
                page,
                page_size = pageSize
            }
        });
    }

    /// <summary>File a new claim.</summary>
    [HttpPost("/v1/claims")]
    public async Task<IActionResult> FileClaim(
        [FromBody] FileClaimRequest request,
        CancellationToken cancellationToken = default)
    {
        if (!System.Enum.TryParse<ClaimType>(request.ClaimType, ignoreCase: true, out var claimType))
            claimType = ClaimType.Death;

        var customerId = request.CustomerId
            ?? (_currentUser.UserId != Guid.Empty ? _currentUser.UserId.ToString() : string.Empty);

        var command = new FileClaimCommand(
            request.PolicyId,
            customerId,
            claimType,
            request.ClaimedAmountPaisa,
            request.IncidentDate,
            request.IncidentDescription,
            request.PlaceOfIncident ?? string.Empty);

        var result = await _mediator.Send(command, cancellationToken);

        if (!result.IsSuccess)
            return BadRequest(new { success = false, error = new { message = result.Error } });

        var claimId = result.Value!;
        var createdClaim = await _claimDataGateway.GetClaimAsync(claimId, cancellationToken);

        return Created($"/v1/claims/{claimId}",
            new
            {
                success = true,
                data = new
                {
                    claim_id = claimId,
                    claim = createdClaim
                }
            });
    }

    /// <summary>Get claim by ID.</summary>
    [HttpGet("/v1/claims/{claimId}")]
    public async Task<IActionResult> GetClaim(string claimId, CancellationToken cancellationToken = default)
    {
        var query = new GetClaimQuery(claimId);
        var result = await _mediator.Send(query, cancellationToken);

        if (!result.IsSuccess || result.Value is null)
            return NotFound(new { success = false, error = new { message = $"Claim not found: {claimId}" } });

        return Ok(new { success = true, data = result.Value });
    }

    /// <summary>Update a claim (add notes, docs).</summary>
    [HttpPatch("/v1/claims/{claimId}")]
    public async Task<IActionResult> UpdateClaim(
        string claimId,
        [FromBody] UpdateClaimRequest request,
        CancellationToken cancellationToken = default)
    {
        var claim = await _claimDataGateway.GetClaimAsync(claimId, cancellationToken);
        if (claim is null)
            return NotFound(new { success = false, error = new { message = $"Claim not found: {claimId}" } });

        var changed = false;
        var createdDocuments = new List<ClaimDocument>();

        if (!string.IsNullOrWhiteSpace(request.Notes))
        {
            claim.ProcessorNotes = request.Notes;
            changed = true;
        }

        foreach (var documentId in request.DocumentIds?.Where(id => !string.IsNullOrWhiteSpace(id)).Distinct() ?? [])
        {
            var createdDocument = await CreateClaimDocumentAsync(
                claimId,
                new ClaimDocumentRequest(documentId, "SUPPORTING"),
                cancellationToken);
            createdDocuments.Add(createdDocument);
            changed = true;
        }

        if (!changed)
            return BadRequest(new { success = false, error = new { message = "No claim updates were provided" } });

        claim.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        var updatedClaim = await _claimDataGateway.UpdateClaimAsync(claim, cancellationToken);

        return Ok(new
        {
            success = true,
            data = new
            {
                message = "Claim updated",
                claim = updatedClaim,
                linked_documents = createdDocuments
            }
        });
    }

    /// <summary>Approve a claim (admin/underwriter only — AuthZ enforced at gateway).</summary>
    [HttpPost("/v1/claims/{claimId}/approve")]
    public async Task<IActionResult> ApproveClaim(
        string claimId,
        [FromBody] ClaimApproveRequest request,
        CancellationToken cancellationToken = default)
    {
        var approverId = _currentUser.UserId != Guid.Empty ? _currentUser.UserId.ToString() : "system";
        // ApproveClaimCommand(string ClaimId, string ApproverId, long ApprovedAmountPaisa, string Notes)
        var command = new ApproveClaimCommand(
            claimId,
            approverId,
            request.ApprovedAmountPaisa,
            request.Notes ?? string.Empty);

        var result = await _mediator.Send(command, cancellationToken);

        if (!result.IsSuccess)
            return BadRequest(new { success = false, error = new { message = result.Error } });

        var updatedClaim = await _claimDataGateway.GetClaimAsync(claimId, cancellationToken);

        return Ok(new
        {
            success = true,
            data = new
            {
                claim_id = claimId,
                status = updatedClaim?.Status.ToString(),
                claim = updatedClaim
            }
        });
    }

    /// <summary>Reject a claim (admin/underwriter only).</summary>
    [HttpPost("/v1/claims/{claimId}/reject")]
    public async Task<IActionResult> RejectClaim(
        string claimId,
        [FromBody] ClaimRejectRequest request,
        CancellationToken cancellationToken = default)
    {
        var approverId = _currentUser.UserId != Guid.Empty ? _currentUser.UserId.ToString() : "system";
        // RejectClaimCommand(string ClaimId, string ApproverId, string Reason)
        var command = new RejectClaimCommand(
            claimId,
            approverId,
            request.Reason ?? "No reason provided");

        var result = await _mediator.Send(command, cancellationToken);

        if (!result.IsSuccess)
            return BadRequest(new { success = false, error = new { message = result.Error } });

        var updatedClaim = await _claimDataGateway.GetClaimAsync(claimId, cancellationToken);

        return Ok(new
        {
            success = true,
            data = new
            {
                claim_id = claimId,
                status = updatedClaim?.Status.ToString(),
                claim = updatedClaim
            }
        });
    }

    /// <summary>Settle a claim (admin only).</summary>
    [HttpPost("/v1/claims/{claimId}/settle")]
    public async Task<IActionResult> SettleClaim(
        string claimId,
        [FromBody] ClaimSettleRequest request,
        CancellationToken cancellationToken = default)
    {
        var settledBy = _currentUser.UserId != Guid.Empty ? _currentUser.UserId.ToString() : "system";
        // SettleClaimCommand(string ClaimId, string SettledBy, string PaymentReference)
        var command = new SettleClaimCommand(
            claimId,
            settledBy,
            request.PaymentReference ?? $"PAY-{Guid.NewGuid():N}"[..16]);

        var result = await _mediator.Send(command, cancellationToken);

        if (!result.IsSuccess)
            return BadRequest(new { success = false, error = new { message = result.Error } });

        var updatedClaim = await _claimDataGateway.GetClaimAsync(claimId, cancellationToken);

        return Ok(new
        {
            success = true,
            data = new
            {
                claim_id = claimId,
                payment_id = result.Value,
                status = updatedClaim?.Status.ToString(),
                claim = updatedClaim
            }
        });
    }

    /// <summary>Get claim settlement details.</summary>
    [HttpGet("/v1/claims/{claimId}/settlement")]
    public async Task<IActionResult> GetSettlement(string claimId, CancellationToken cancellationToken = default)
    {
        var query = new GetClaimQuery(claimId);
        var result = await _mediator.Send(query, cancellationToken);

        if (!result.IsSuccess || result.Value is null)
            return NotFound(new { success = false, error = new { message = "Claim not found" } });

        return Ok(new { success = true, data = new { claim_id = claimId, claim = result.Value } });
    }

    /// <summary>Add document to a claim.</summary>
    [HttpPost("/v1/claims/{claimId}/documents")]
    public async Task<IActionResult> AddClaimDocument(
        string claimId,
        [FromBody] ClaimDocumentRequest request,
        CancellationToken cancellationToken = default)
    {
        var claim = await _claimDataGateway.GetClaimAsync(claimId, cancellationToken);
        if (claim is null)
            return NotFound(new { success = false, error = new { message = $"Claim not found: {claimId}" } });

        var createdDocument = await CreateClaimDocumentAsync(claimId, request, cancellationToken);

        claim.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        await _claimDataGateway.UpdateClaimAsync(claim, cancellationToken);

        _logger.LogInformation("Document {DocId} linked to claim {ClaimId}", createdDocument.DocumentId, claimId);

        return Created($"/v1/claims/{claimId}/documents/{createdDocument.DocumentId}",
            new { success = true, data = createdDocument });
    }

    /// <summary>Review a claim (L1/L2/L3 underwriting).</summary>
    [HttpPost("/v1/claims/{claimId}/review")]
    public async Task<IActionResult> ReviewClaim(
        string claimId,
        [FromBody] ClaimReviewRequest request,
        CancellationToken cancellationToken = default)
    {
        var claim = await _claimDataGateway.GetClaimAsync(claimId, cancellationToken);
        if (claim is null)
            return NotFound(new { success = false, error = new { message = $"Claim not found: {claimId}" } });

        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        var approverId = _currentUser.UserId != Guid.Empty ? _currentUser.UserId.ToString() : "system";
        var approvalLevel = ParseApprovalLevel(request.ReviewLevel);

        claim.Status = ClaimStatus.UnderReview;
        claim.ProcessorNotes = AppendNote(claim.ProcessorNotes, $"[{request.ReviewLevel}] {request.Notes}".Trim());
        claim.UpdatedAt = now;
        var updatedClaim = await _claimDataGateway.UpdateClaimAsync(claim, cancellationToken);

        var approval = new ClaimApproval
        {
            ApprovalId = Guid.NewGuid().ToString("N"),
            ClaimId = claimId,
            ApproverId = approverId,
            ApproverRole = string.IsNullOrWhiteSpace(request.ReviewLevel) ? "REVIEW" : request.ReviewLevel,
            ApprovalLevel = approvalLevel,
            Decision = ApprovalDecision.Pending,
            ApprovedAmount = new Money { Amount = 0, Currency = "BDT" },
            Notes = request.Notes ?? string.Empty,
            CreatedAt = now,
            ApprovedCurrency = "BDT"
        };

        var createdApproval = await _insuranceClient.Client.CreateClaimApprovalAsync(
            new CreateClaimApprovalRequest { Approval = approval },
            _insuranceClient.BuildCallOptions(cancellationToken));

        _logger.LogInformation("Claim review persisted for {ClaimId} level={Level}", claimId, request.ReviewLevel);

        return Ok(new
        {
            success = true,
            data = new
            {
                message = "Claim review submitted",
                claim = updatedClaim,
                approval = createdApproval.Approval
            }
        });
    }

    private async Task<ClaimDocument> CreateClaimDocumentAsync(
        string claimId,
        ClaimDocumentRequest request,
        CancellationToken cancellationToken)
    {
        if (string.IsNullOrWhiteSpace(request.DocumentId) && string.IsNullOrWhiteSpace(request.FileUrl))
            throw new InvalidOperationException("Either document_id or file_url is required");

        var fileUrl = request.FileUrl;
        var resolvedDocumentId = request.DocumentId;

        if (string.IsNullOrWhiteSpace(fileUrl) && !string.IsNullOrWhiteSpace(request.DocumentId))
        {
            var document = await _docgenClient.GetDocumentAsync(request.DocumentId, cancellationToken);
            fileUrl = document?.FileUrl;
            resolvedDocumentId = string.IsNullOrWhiteSpace(document?.Id) ? request.DocumentId : document.Id;
        }

        if (string.IsNullOrWhiteSpace(fileUrl))
            throw new InvalidOperationException($"Unable to resolve a file URL for document {request.DocumentId}");

        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        var documentToCreate = new ClaimDocument
        {
            DocumentId = string.IsNullOrWhiteSpace(resolvedDocumentId) ? Guid.NewGuid().ToString("N") : resolvedDocumentId,
            ClaimId = claimId,
            DocumentType = string.IsNullOrWhiteSpace(request.DocumentType) ? "SUPPORTING" : request.DocumentType,
            FileUrl = fileUrl,
            FileHash = string.IsNullOrWhiteSpace(request.FileHash) ? ComputeSha256(fileUrl) : request.FileHash,
            UploadedAt = now,
            CreatedAt = now,
            UpdatedAt = now
        };

        var createdDocument = await _insuranceClient.Client.CreateClaimDocumentAsync(
            new CreateClaimDocumentRequest { Document = documentToCreate },
            _insuranceClient.BuildCallOptions(cancellationToken));

        return createdDocument.Document;
    }

    private static int ParseApprovalLevel(string? reviewLevel)
    {
        if (string.IsNullOrWhiteSpace(reviewLevel))
            return 1;

        var normalized = reviewLevel.Trim().ToUpperInvariant();
        if (normalized.StartsWith('L') && int.TryParse(normalized[1..], out var levelFromLabel) && levelFromLabel > 0)
            return levelFromLabel;

        return int.TryParse(normalized, out var level) && level > 0 ? level : 1;
    }

    private static string AppendNote(string? existingNotes, string newNote)
    {
        if (string.IsNullOrWhiteSpace(newNote))
            return existingNotes ?? string.Empty;
        if (string.IsNullOrWhiteSpace(existingNotes))
            return newNote;

        return $"{existingNotes}{Environment.NewLine}{newNote}";
    }

    private static string ComputeSha256(string text)
    {
        var bytes = Encoding.UTF8.GetBytes(text);
        using var sha = SHA256.Create();
        return Convert.ToHexString(sha.ComputeHash(bytes)).ToLowerInvariant();
    }
}

// ── Request DTOs ──────────────────────────────────────────────────────────────

public sealed record FileClaimRequest(
    string PolicyId,
    string ClaimType,
    long ClaimedAmountPaisa,
    DateTime IncidentDate,
    string IncidentDescription,
    string? PlaceOfIncident = null,
    string? CustomerId = null);

public sealed record UpdateClaimRequest(
    string? Notes = null,
    List<string>? DocumentIds = null);

public sealed record ClaimApproveRequest(
    long ApprovedAmountPaisa,
    string? Notes = null);

public sealed record ClaimRejectRequest(
    string? Reason = null);

public sealed record ClaimSettleRequest(
    string? PaymentReference = null,
    string? Notes = null);

public sealed record ClaimDocumentRequest(
    string? DocumentId = null,
    string? DocumentType = null,
    string? FileUrl = null,
    string? FileHash = null);

public sealed record ClaimReviewRequest(
    string ReviewLevel,
    string? Notes = null);
