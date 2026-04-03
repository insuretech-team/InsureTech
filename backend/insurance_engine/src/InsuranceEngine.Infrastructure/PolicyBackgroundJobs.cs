using MediatR;
using Microsoft.Extensions.Logging;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using InsuranceEngine.Infrastructure.Notifications;
using InsuranceEngine.Infrastructure.Renewals;
using System;
using System.Linq;
using System.Threading.Tasks;

namespace InsuranceEngine.Infrastructure;

public class PolicyBackgroundJobs
{
    private readonly IRepository<PolicyEntity> _repository;
    private readonly IMediator _mediator;
    private readonly IGracePeriodService _gracePeriodService;
    private readonly INotificationService _notificationService;
    private readonly ILogger<PolicyBackgroundJobs> _logger;

    public PolicyBackgroundJobs(
        IRepository<PolicyEntity> repository,
        IMediator mediator,
        IGracePeriodService gracePeriodService,
        INotificationService notificationService,
        ILogger<PolicyBackgroundJobs> logger)
    {
        _repository = repository;
        _mediator = mediator;
        _gracePeriodService = gracePeriodService;
        _notificationService = notificationService;
        _logger = logger;
    }

    /// <summary>
    /// FR-047 + FR-048: Process grace period workflow.
    /// 1. Move expired ACTIVE policies to GRACE_PERIOD (30 days)
    /// 2. Send daily reminders during grace period
    /// 3. Auto-lapse after grace period expires
    /// </summary>
    public async Task ProcessGracePeriodWorkflowAsync()
    {
        _logger.LogInformation("Starting Grace Period workflow at {Time}", DateTime.UtcNow);

        try
        {
            await _gracePeriodService.ProcessExpiredPoliciesAsync();
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error in ProcessExpiredPoliciesAsync");
        }

        try
        {
            await _gracePeriodService.ProcessGracePeriodRemindersAsync();
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error in ProcessGracePeriodRemindersAsync");
        }

        try
        {
            await _gracePeriodService.ProcessGracePeriodExpiryAsync();
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error in ProcessGracePeriodExpiryAsync");
        }

        _logger.LogInformation("Grace Period workflow completed at {Time}", DateTime.UtcNow);
    }

    /// <summary>
    /// FR-048: Legacy auto-lapse for backward compatibility.
    /// Now delegated to GracePeriodService.
    /// </summary>
    [Obsolete("Use ProcessGracePeriodWorkflowAsync instead")]
    public async Task ProcessAutoLapseAsync()
    {
        _logger.LogInformation("Starting Legacy Auto-Lapse background job at {Time}", DateTime.UtcNow);
        await ProcessGracePeriodWorkflowAsync();
    }

    /// <summary>
    /// FR-045: Send renewal reminders for policies expiring in 30 days.
    /// Runs daily.
    /// </summary>
    public async Task ProcessRenewalRemindersAsync()
    {
        _logger.LogInformation("Starting Renewal Reminders background job at {Time}", DateTime.UtcNow);

        var thirtyDaysFromNow = DateTime.UtcNow.AddDays(30);
        var query = await _repository.FindAsync(p => p.Status == "ACTIVE" && p.EndDate <= thirtyDaysFromNow && p.EndDate > DateTime.UtcNow);
        var policiesToRemind = query.ToList();

        _logger.LogInformation("Found {Count} policies requiring renewal reminders.", policiesToRemind.Count);

        foreach (var policy in policiesToRemind)
        {
            try
            {
                var daysUntilExpiry = (int)(policy.EndDate - DateTime.UtcNow).TotalDays;
                
                var shouldNotify = daysUntilExpiry switch
                {
                    30 or 15 or 7 or 3 or 1 => true,
                    _ => daysUntilExpiry <= 3
                };

                if (shouldNotify)
                {
                    await _notificationService.NotifyRenewalReminderAsync(
                        policy.CustomerId.ToString(),
                        policy.PolicyNumber,
                        policy.EndDate,
                        CancellationToken.None);

                    _logger.LogInformation(
                        "Renewal reminder sent for policy {PolicyNumber} (Expires in {Days} days)",
                        policy.PolicyNumber, daysUntilExpiry);
                }
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Failed to send renewal reminder for policy {PolicyId}", policy.PolicyId);
            }
        }
    }
}
