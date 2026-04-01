using Google.Protobuf.WellKnownTypes;
using Insuretech.Insurance.Services.V1;
using MediatR;
using Microsoft.AspNetCore.Mvc;
using PoliSync.Infrastructure.Clients;
using PoliSync.Policy.Application.Commands;
using PoliSync.Policy.Application.Queries;
using PoliSync.Policy.Domain;
using PoliSync.Policy.Infrastructure;
using PoliSync.SharedKernel.Auth;
using System.Text.Json;

namespace PoliSync.ApiHost.Controllers;

/// <summary>
/// HTTP companion for the Policy gRPC service.
/// The InScore gateway reverse-proxies /v1/policies/* to this controller on port 50161.
/// </summary>
[ApiController]
public sealed class PoliciesController : ControllerBase
{
    private readonly IMediator _mediator;
    private readonly IPolicyDataGateway _policyDataGateway;
    private readonly InsuranceServiceClient _insuranceClient;
    private readonly ICurrentUser _currentUser;
    private readonly ILogger<PoliciesController> _logger;

    public PoliciesController(
        IMediator mediator,
        IPolicyDataGateway policyDataGateway,
        InsuranceServiceClient insuranceClient,
        ICurrentUser currentUser,
        ILogger<PoliciesController> logger)
    {
        _mediator = mediator;
        _policyDataGateway = policyDataGateway;
        _insuranceClient = insuranceClient;
        _currentUser = currentUser;
        _logger = logger;
    }

    /// <summary>List policies for the authenticated user.</summary>
    [HttpGet("/v1/policies")]
    public async Task<IActionResult> ListPolicies(
        [FromQuery(Name = "user_id")] string? userId = null,
        [FromQuery] int page = 1,
        [FromQuery(Name = "page_size")] int pageSize = 20,
        CancellationToken cancellationToken = default)
    {
        // B2C: use authenticated user's ID; admins can pass explicit user_id
        var effectiveUserId = userId
            ?? (_currentUser.UserId != Guid.Empty ? _currentUser.UserId.ToString() : null)
            ?? string.Empty;

        var query = new ListUserPoliciesQuery(effectiveUserId, page, pageSize);
        var result = await _mediator.Send(query, cancellationToken);

        if (!result.IsSuccess)
            return StatusCode(500, new { success = false, error = new { message = result.Error } });

        var policies = result.Value?.Policies ?? [];
        var totalCount = result.Value?.TotalCount ?? 0;

        return Ok(new
        {
            success = true,
            data = new
            {
                policies,
                total_count = totalCount,
                page,
                page_size = pageSize
            }
        });
    }

    /// <summary>Get policy by ID.</summary>
    [HttpGet("/v1/policies/{policyId}")]
    public async Task<IActionResult> GetPolicy(string policyId, CancellationToken cancellationToken = default)
    {
        var query = new GetPolicyQuery(policyId);
        var result = await _mediator.Send(query, cancellationToken);

        if (!result.IsSuccess || result.Value is null)
            return NotFound(new { success = false, error = new { message = $"Policy not found: {policyId}" } });

        return Ok(new { success = true, data = result.Value });
    }

    /// <summary>Get policy by policy number.</summary>
    [HttpGet("/v1/policies/number/{policyNumber}/lookup")]
    public async Task<IActionResult> GetPolicyByNumber(string policyNumber, CancellationToken cancellationToken = default)
    {
        // List and filter — or query by number if the handler supports it
        var query = new GetPolicyQuery(policyNumber); // handler resolves by number too
        var result = await _mediator.Send(query, cancellationToken);

        if (!result.IsSuccess || result.Value is null)
            return NotFound(new { success = false, error = new { message = $"Policy not found: {policyNumber}" } });

        return Ok(new { success = true, data = result.Value });
    }

    /// <summary>Cancel a policy.</summary>
    [HttpPost("/v1/policies/{policyId}/cancel")]
    public async Task<IActionResult> CancelPolicy(
        string policyId,
        [FromBody] CancelPolicyRequest request,
        CancellationToken cancellationToken = default)
    {
        var requestedBy = _currentUser.UserId != Guid.Empty ? _currentUser.UserId.ToString() : request.RequestedBy;
        var command = new CancelPolicyCommand(policyId, request.Reason, requestedBy, "B2C");
        var result = await _mediator.Send(command, cancellationToken);

        if (!result.IsSuccess)
            return BadRequest(new { success = false, error = new { message = result.Error } });

        var updatedPolicy = await _policyDataGateway.GetPolicyAsync(policyId, cancellationToken);

        return Ok(new
        {
            success = true,
            data = new
            {
                policy_id = policyId,
                status = updatedPolicy?.Status.ToString(),
                policy = updatedPolicy,
                reason = request.Reason
            }
        });
    }

    /// <summary>Suspend a policy.</summary>
    [HttpPost("/v1/policies/{policyId}/suspend")]
    public async Task<IActionResult> SuspendPolicy(
        string policyId,
        [FromBody] PolicyActionRequest request,
        CancellationToken cancellationToken = default)
    {
        var policy = await _policyDataGateway.GetPolicyAsync(policyId, cancellationToken);
        if (policy is null)
            return NotFound(new { success = false, error = new { message = $"Policy not found: {policyId}" } });

        try
        {
            var aggregate = new PolicyAggregate(policy);
            aggregate.SuspendPolicy();
            var updated = await _policyDataGateway.UpdatePolicyAsync(policy, cancellationToken);

            _logger.LogInformation("Policy suspended for {PolicyId} by {User}", policyId, _currentUser.UserId);
            return Ok(new { success = true, data = new { message = "Policy suspended", policy = updated, reason = request.Reason } });
        }
        catch (InvalidOperationException ex)
        {
            return BadRequest(new { success = false, error = new { message = ex.Message } });
        }
    }

    /// <summary>Reinstate a suspended policy.</summary>
    [HttpPost("/v1/policies/{policyId}/reinstate")]
    public async Task<IActionResult> ReinstatePolicy(
        string policyId,
        [FromBody] PolicyActionRequest? request,
        CancellationToken cancellationToken = default)
    {
        var policy = await _policyDataGateway.GetPolicyAsync(policyId, cancellationToken);
        if (policy is null)
            return NotFound(new { success = false, error = new { message = $"Policy not found: {policyId}" } });

        try
        {
            var aggregate = new PolicyAggregate(policy);
            aggregate.ReinstatePolicy();
            var updated = await _policyDataGateway.UpdatePolicyAsync(policy, cancellationToken);

            _logger.LogInformation("Policy reinstated for {PolicyId} by {User}", policyId, _currentUser.UserId);
            return Ok(new { success = true, data = new { message = "Policy reinstated", policy = updated, reason = request?.Reason } });
        }
        catch (InvalidOperationException ex)
        {
            return BadRequest(new { success = false, error = new { message = ex.Message } });
        }
    }

    /// <summary>Get the document for a policy (e.g. policy certificate PDF).</summary>
    [HttpPost("/v1/policies/{policyId}/document")]
    public async Task<IActionResult> RequestPolicyDocument(
        string policyId,
        CancellationToken cancellationToken = default)
    {
        var policy = await _policyDataGateway.GetPolicyAsync(policyId, cancellationToken);
        if (policy is null)
            return NotFound(new { success = false, error = new { message = $"Policy not found: {policyId}" } });

        if (!string.IsNullOrWhiteSpace(policy.PolicyDocumentUrl))
        {
            return Ok(new
            {
                success = true,
                data = new
                {
                    policy_id = policyId,
                    document_url = policy.PolicyDocumentUrl,
                    status = "AVAILABLE"
                }
            });
        }

        var now = DateTime.UtcNow;
        var actorId = _currentUser.UserId != Guid.Empty ? _currentUser.UserId.ToString() : policy.CustomerId;
        var serviceRequest = new Insuretech.Policy.Entity.V1.PolicyServiceRequest
        {
            RequestId = Guid.NewGuid().ToString("N"),
            PolicyId = policyId,
            CustomerId = policy.CustomerId,
            RequestType = Insuretech.Policy.Entity.V1.ServiceRequestType.DownloadPolicy,
            RequestData = JsonSerializer.Serialize(new
            {
                requested_by = actorId,
                requested_at = now,
                policy_number = policy.PolicyNumber
            }),
            Status = Insuretech.Policy.Entity.V1.ServiceRequestStatus.Pending,
            CreatedAt = Timestamp.FromDateTime(now)
        };

        var created = await _insuranceClient.Client.CreatePolicyServiceRequestAsync(
            new CreatePolicyServiceRequestRequest { Request = serviceRequest },
            _insuranceClient.BuildCallOptions(cancellationToken));

        _logger.LogInformation("Policy document request persisted for {PolicyId} as {RequestId}", policyId, created.Request.RequestId);

        return Accepted(new
        {
            success = true,
            data = new
            {
                policy_id = policyId,
                request_id = created.Request.RequestId,
                request_type = created.Request.RequestType.ToString(),
                status = created.Request.Status.ToString()
            }
        });
    }

    /// <summary>Get nominees for a policy.</summary>
    [HttpGet("/v1/policies/{policyId}/nominees")]
    public async Task<IActionResult> GetNominees(string policyId, CancellationToken cancellationToken = default)
    {
        var policy = await _policyDataGateway.GetPolicyAsync(policyId, cancellationToken);
        if (policy is null)
            return NotFound(new { success = false, error = new { message = $"Policy not found: {policyId}" } });

        return Ok(new
        {
            success = true,
            data = new
            {
                nominees = policy.Nominees,
                total = policy.Nominees.Count,
                policy_id = policyId
            }
        });
    }

    /// <summary>Get renewal info for a policy.</summary>
    [HttpGet("/v1/policies/{policyId}/renewal")]
    public async Task<IActionResult> GetRenewal(string policyId, CancellationToken cancellationToken = default)
    {
        var policy = await _policyDataGateway.GetPolicyAsync(policyId, cancellationToken);
        if (policy is null)
            return NotFound(new { success = false, error = new { message = $"Policy not found: {policyId}" } });

        var schedulesResponse = await _insuranceClient.Client.ListRenewalSchedulesAsync(
            new ListRenewalSchedulesRequest { PolicyId = policyId },
            _insuranceClient.BuildCallOptions(cancellationToken));

        if (schedulesResponse.Schedules.Count == 0)
            return NotFound(new { success = false, error = new { message = $"No renewal schedule found for policy {policyId}" } });

        var latestSchedule = schedulesResponse.Schedules
            .OrderByDescending(schedule => schedule.RenewalDueDate?.Seconds ?? 0)
            .First();

        var remindersResponse = await _insuranceClient.Client.ListRenewalRemindersAsync(
            new ListRenewalRemindersRequest { ScheduleId = latestSchedule.Id },
            _insuranceClient.BuildCallOptions(cancellationToken));

        return Ok(new
        {
            success = true,
            data = new
            {
                policy_id = policyId,
                schedule = latestSchedule,
                reminders = remindersResponse.Reminders,
                reminder_count = remindersResponse.Reminders.Count
            }
        });
    }
}

// ── Request DTOs ──────────────────────────────────────────────────────────────

public sealed record CancelPolicyRequest(
    string Reason,
    string? RequestedBy = null);

public sealed record PolicyActionRequest(
    string? Reason = null);
