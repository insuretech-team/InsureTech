using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using InsuranceEngine.Infrastructure.Notifications;
using InsuranceEngine.Infrastructure.Messaging;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using System.Text.Json;

namespace InsuranceEngine.Infrastructure.Renewals;

public interface IGracePeriodService
{
    Task ProcessExpiredPoliciesAsync(CancellationToken ct = default);
    Task ProcessGracePeriodRemindersAsync(CancellationToken ct = default);
    Task ProcessGracePeriodExpiryAsync(CancellationToken ct = default);
    Task<GracePeriodInfo?> GetGracePeriodInfoAsync(string policyId, CancellationToken ct = default);
    Task<ReinstatementResult> ReinstatePolicyAsync(ReinstatementRequest request, CancellationToken ct = default);
    Task<bool> CanPolicyBeReinstatedAsync(string policyId, CancellationToken ct = default);
}

public class GracePeriodService : IGracePeriodService
{
    private readonly IRepository<PolicyEntity> _policyRepository;
    private readonly InsuranceEngine.Renewals.IRenewalDataGateway _renewalGateway;
    private readonly INotificationService _notificationService;
    private readonly IEventPublisher _eventPublisher;
    private readonly ILogger<GracePeriodService> _logger;
    private readonly GracePeriodSettings _settings;

    public GracePeriodService(
        IRepository<PolicyEntity> policyRepository,
        InsuranceEngine.Renewals.IRenewalDataGateway renewalGateway,
        INotificationService notificationService,
        IEventPublisher eventPublisher,
        IOptions<GracePeriodSettings> settings,
        ILogger<GracePeriodService> logger)
    {
        _policyRepository = policyRepository;
        _renewalGateway = renewalGateway;
        _notificationService = notificationService;
        _eventPublisher = eventPublisher;
        _settings = settings.Value;
        _logger = logger;
    }

    public async Task ProcessExpiredPoliciesAsync(CancellationToken ct = default)
    {
        _logger.LogInformation("Processing expired policies for grace period entry at {Time}", DateTime.UtcNow);

        var today = DateTime.UtcNow.Date;
        var gracePeriodEnd = today.AddDays(_settings.GracePeriodDays);

        var query = await _policyRepository.FindAsync(
            p => p.Status == "ACTIVE" && p.EndDate < today,
            ct);

        var expiredPolicies = query.ToList();
        _logger.LogInformation("Found {Count} expired policies to move to grace period", expiredPolicies.Count);

        foreach (var policy in expiredPolicies)
        {
            try
            {
                policy.Status = "GRACE_PERIOD";
                policy.UpdatedAt = DateTime.UtcNow;

                var metadata = ParseMetadata(policy.UnderwritingData);
                metadata["GracePeriodStartDate"] = today.ToString("O");
                metadata["GracePeriodEndDate"] = gracePeriodEnd.ToString("O");
                policy.UnderwritingData = JsonSerializer.Serialize(metadata);

                await _policyRepository.UpdateAsync(policy, ct);

                _logger.LogInformation(
                    "Policy {PolicyNumber} moved to GRACE_PERIOD. EndDate: {EndDate}, GracePeriodEnds: {GraceEnd}",
                    policy.PolicyNumber, policy.EndDate, gracePeriodEnd);

                await _eventPublisher.PublishAsync("policy.lifecycle.events", new
                {
                    EventType = "PolicyEnteredGracePeriod",
                    PolicyId = policy.PolicyId.ToString(),
                    PolicyNumber = policy.PolicyNumber,
                    CustomerId = policy.CustomerId.ToString(),
                    GracePeriodStartDate = today,
                    GracePeriodEndDate = gracePeriodEnd,
                    Timestamp = DateTime.UtcNow
                });

                await _notificationService.NotifyGracePeriodAsync(
                    policy.CustomerId.ToString(),
                    policy.PolicyNumber,
                    _settings.GracePeriodDays,
                    ct);

                _logger.LogInformation("Grace period notification sent for policy {PolicyNumber}", policy.PolicyNumber);
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Failed to process grace period for policy {PolicyId}", policy.PolicyId);
            }
        }
    }

    public async Task ProcessGracePeriodRemindersAsync(CancellationToken ct = default)
    {
        if (!_settings.EnableDailyReminders)
        {
            _logger.LogDebug("Grace period reminders disabled");
            return;
        }

        _logger.LogInformation("Processing grace period daily reminders at {Time}", DateTime.UtcNow);

        var today = DateTime.UtcNow.Date;
        var gracePeriodPolicies = await _policyRepository.FindAsync(
            p => p.Status == "GRACE_PERIOD",
            ct);

        foreach (var policy in gracePeriodPolicies)
        {
            try
            {
                var metadata = ParseMetadata(policy.UnderwritingData);
                var graceEndDate = ParseGraceEndDate(metadata);
                var daysRemaining = (int)(graceEndDate - today).TotalDays;

                if (daysRemaining < 0)
                {
                    _logger.LogWarning(
                        "Policy {PolicyNumber} grace period expired but status still GRACE_PERIOD",
                        policy.PolicyNumber);
                    continue;
                }

                var shouldNotify = daysRemaining switch
                {
                    30 or 15 or 7 or 3 or 1 => true,
                    _ => daysRemaining <= 7
                };

                if (shouldNotify)
                {
                    await _notificationService.NotifyGracePeriodAsync(
                        policy.CustomerId.ToString(),
                        policy.PolicyNumber,
                        daysRemaining,
                        ct);

                    _logger.LogInformation(
                        "Grace period reminder sent for policy {PolicyNumber}: {DaysRemaining} days remaining",
                        policy.PolicyNumber, daysRemaining);
                }
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Failed to send grace period reminder for policy {PolicyId}", policy.PolicyId);
            }
        }
    }

    public async Task ProcessGracePeriodExpiryAsync(CancellationToken ct = default)
    {
        _logger.LogInformation("Processing grace period expiry (auto-lapse) at {Time}", DateTime.UtcNow);

        var today = DateTime.UtcNow.Date;
        var policies = await _policyRepository.FindAsync(
            p => p.Status == "GRACE_PERIOD",
            ct);

        var expiredPolicies = policies
            .Where(p =>
            {
                var metadata = ParseMetadata(p.UnderwritingData);
                var graceEndDate = ParseGraceEndDate(metadata);
                return graceEndDate < today;
            })
            .ToList();

        _logger.LogInformation("Found {Count} policies with expired grace period to lapse", expiredPolicies.Count);

        foreach (var policy in expiredPolicies)
        {
            try
            {
                policy.Status = "LAPSED";
                policy.UpdatedAt = DateTime.UtcNow;

                var metadata = ParseMetadata(policy.UnderwritingData);
                metadata["GracePeriodExpiredDate"] = DateTime.UtcNow.ToString("O");
                metadata["ReinstatementWindowEndDate"] = today.AddDays(_settings.ReinstatementWindowDays).ToString("O");
                policy.UnderwritingData = JsonSerializer.Serialize(metadata);

                await _policyRepository.UpdateAsync(policy, ct);

                _logger.LogInformation(
                    "Policy {PolicyNumber} auto-lapsed after grace period expired. Reinstatement window: {WindowEnd}",
                    policy.PolicyNumber,
                    today.AddDays(_settings.ReinstatementWindowDays));

                await _eventPublisher.PublishAsync("policy.lifecycle.events", new
                {
                    EventType = "PolicyLapsed",
                    PolicyId = policy.PolicyId.ToString(),
                    PolicyNumber = policy.PolicyNumber,
                    CustomerId = policy.CustomerId.ToString(),
                    Reason = "GracePeriodExpired",
                    ReinstatementWindowEndDate = today.AddDays(_settings.ReinstatementWindowDays),
                    Timestamp = DateTime.UtcNow
                });

                await _notificationService.NotifyPolicyLapsedAsync(
                    policy.CustomerId.ToString(),
                    policy.PolicyNumber,
                    ct);

                _logger.LogInformation("Lapsed notification sent for policy {PolicyNumber}", policy.PolicyNumber);
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Failed to lapse policy {PolicyId}", policy.PolicyId);
            }
        }
    }

    public async Task<GracePeriodInfo?> GetGracePeriodInfoAsync(string policyId, CancellationToken ct = default)
    {
        if (!Guid.TryParse(policyId, out var policyGuid))
            return null;

        var policies = await _policyRepository.FindAsync(p => p.PolicyId == policyGuid, ct);
        var policy = policies.FirstOrDefault();

        if (policy == null)
            return null;

        if (policy.Status != "GRACE_PERIOD" && policy.Status != "LAPSED")
            return null;

        var today = DateTime.UtcNow.Date;
        var metadata = ParseMetadata(policy.UnderwritingData);
        DateTime graceEndDate = ParseGraceEndDate(metadata);

        var daysRemaining = (int)(graceEndDate - today).TotalDays;
        if (daysRemaining < 0) daysRemaining = 0;

        var canReinstate = policy.Status == "LAPSED" && IsWithinReinstatementWindow(metadata);
        var canRenew = policy.Status == "ACTIVE" || policy.Status == "GRACE_PERIOD";

        var reinstatementPenalty = canReinstate
            ? CalculateReinstatementPenalty(policy, _settings.ReinstatementPenaltyPercent)
            : 0;

        var nextAction = DetermineNextAction(policy.Status, daysRemaining, canReinstate);

        return new GracePeriodInfo
        {
            PolicyId = policy.PolicyId.ToString(),
            PolicyNumber = policy.PolicyNumber,
            Status = policy.Status,
            GracePeriodStartDate = ParseGraceStartDate(metadata),
            GracePeriodEndDate = graceEndDate,
            DaysRemaining = daysRemaining,
            ReinstatementPenaltyAmount = reinstatementPenalty,
            CanReinstate = canReinstate,
            CanRenew = canRenew,
            NextAction = nextAction
        };
    }

    public async Task<bool> CanPolicyBeReinstatedAsync(string policyId, CancellationToken ct = default)
    {
        if (!Guid.TryParse(policyId, out var policyGuid))
            return false;

        var policies = await _policyRepository.FindAsync(p => p.PolicyId == policyGuid, ct);
        var policy = policies.FirstOrDefault();

        if (policy == null)
            return false;

        var metadata = ParseMetadata(policy.UnderwritingData);
        return policy.Status == "LAPSED" && IsWithinReinstatementWindow(metadata);
    }

    public async Task<ReinstatementResult> ReinstatePolicyAsync(ReinstatementRequest request, CancellationToken ct = default)
    {
        _logger.LogInformation("Reinstatement requested for policy {PolicyId}", request.PolicyId);

        if (!Guid.TryParse(request.PolicyId, out var policyGuid))
        {
            return new ReinstatementResult
            {
                Success = false,
                ErrorMessage = "Invalid policy ID format"
            };
        }

        var policies = await _policyRepository.FindAsync(p => p.PolicyId == policyGuid, ct);
        var policy = policies.FirstOrDefault();

        if (policy == null)
        {
            return new ReinstatementResult
            {
                Success = false,
                ErrorMessage = "Policy not found"
            };
        }

        var metadata = ParseMetadata(policy.UnderwritingData);
        if (!IsWithinReinstatementWindow(metadata))
        {
            return new ReinstatementResult
            {
                Success = false,
                ErrorMessage = $"Policy {policy.PolicyNumber} is outside the reinstatement window (90 days from lapse)"
            };
        }

        var reinstatementPenalty = request.ApplyReinstatementPenalty
            ? CalculateReinstatementPenalty(policy, _settings.ReinstatementPenaltyPercent)
            : 0;

        try
        {
            var grpcRequest = new Insuretech.Policy.Services.V1.RenewPolicyTenureRequest
            {
                PolicyId = request.PolicyId,
                TenureMonths = request.TenureMonths
            };

            var response = await _renewalGateway.RenewPolicyAsync(grpcRequest, ct);

            if (response.Error != null)
            {
                return new ReinstatementResult
                {
                    Success = false,
                    ErrorMessage = response.Error.Message
                };
            }

            await _eventPublisher.PublishAsync("policy.lifecycle.events", new
            {
                EventType = "PolicyReinstated",
                OriginalPolicyId = request.PolicyId,
                NewPolicyId = response.NewPolicyId,
                ReinstatementPenalty = reinstatementPenalty,
                TenureMonths = request.TenureMonths,
                Notes = request.Notes,
                Timestamp = DateTime.UtcNow
            });

            _logger.LogInformation(
                "Policy {OriginalPolicyId} reinstated successfully. New Policy: {NewPolicyId}",
                request.PolicyId, response.NewPolicyId);

            return new ReinstatementResult
            {
                Success = true,
                NewPolicyId = response.NewPolicyId,
                NewPolicyNumber = response.NewPolicyNumber,
                ReinstatementPenalty = reinstatementPenalty,
                TotalAmountDue = reinstatementPenalty
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to reinstate policy {PolicyId}", request.PolicyId);
            return new ReinstatementResult
            {
                Success = false,
                ErrorMessage = ex.Message
            };
        }
    }

    private Dictionary<string, string> ParseMetadata(string? json)
    {
        if (string.IsNullOrEmpty(json))
            return new Dictionary<string, string>();

        try
        {
            return JsonSerializer.Deserialize<Dictionary<string, string>>(json) ?? new Dictionary<string, string>();
        }
        catch
        {
            return new Dictionary<string, string>();
        }
    }

    private DateTime ParseGraceStartDate(Dictionary<string, string> metadata)
    {
        if (metadata.TryGetValue("GracePeriodStartDate", out var startStr) &&
            DateTime.TryParse(startStr, out var startDate))
        {
            return startDate;
        }
        return DateTime.UtcNow;
    }

    private DateTime ParseGraceEndDate(Dictionary<string, string> metadata)
    {
        if (metadata.TryGetValue("GracePeriodEndDate", out var endStr) &&
            DateTime.TryParse(endStr, out var endDate))
        {
            return endDate;
        }
        return DateTime.UtcNow.AddDays(-1);
    }

    private bool IsWithinReinstatementWindow(Dictionary<string, string> metadata)
    {
        if (metadata.TryGetValue("ReinstatementWindowEndDate", out var windowEndStr) &&
            DateTime.TryParse(windowEndStr, out var windowEnd))
        {
            return DateTime.UtcNow <= windowEnd;
        }

        if (metadata.TryGetValue("GracePeriodEndDate", out var graceEndStr) &&
            DateTime.TryParse(graceEndStr, out var graceEnd))
        {
            return DateTime.UtcNow <= graceEnd.AddDays(_settings.ReinstatementWindowDays - _settings.GracePeriodDays);
        }
        return false;
    }

    private decimal CalculateReinstatementPenalty(PolicyEntity policy, decimal penaltyPercent)
    {
        if (policy.PremiumAmount <= 0)
            return 0;

        return (policy.PremiumAmount / 100.0m) * (penaltyPercent / 100);
    }

    private string DetermineNextAction(string status, int daysRemaining, bool canReinstate)
    {
        return status switch
        {
            "GRACE_PERIOD" when daysRemaining > 0 => $"Pay renewal premium within {daysRemaining} days to maintain coverage",
            "GRACE_PERIOD" when daysRemaining <= 0 => "Grace period expired - policy will lapse",
            "LAPSED" when canReinstate => $"Reinstate within {_settings.ReinstatementWindowDays} days (penalty may apply)",
            "LAPSED" => "Reinstatement window expired - purchase new policy",
            _ => "Renew your policy"
        };
    }
}

public class GracePeriodEventPublisher : IEventPublisher
{
    private readonly IEventPublisher _kafkaPublisher;
    private readonly ILogger<GracePeriodEventPublisher> _logger;

    public GracePeriodEventPublisher(
        IEventPublisher kafkaPublisher,
        ILogger<GracePeriodEventPublisher> logger)
    {
        _kafkaPublisher = kafkaPublisher;
        _logger = logger;
    }

    public async Task PublishAsync<T>(string topic, T @event, string? key = null) where T : class
    {
        try
        {
            await _kafkaPublisher.PublishAsync(topic, @event, key);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to publish event to {Topic}", topic);
        }
    }
}
