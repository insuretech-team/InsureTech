using MediatR;
using Microsoft.Extensions.Logging;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using InsuranceEngine.Products.Domain.Entities;
using Microsoft.EntityFrameworkCore;
using System;
using System.Linq;
using System.Threading.Tasks;

namespace InsuranceEngine.Policy.Application.BackgroundJobs;

public class PolicyBackgroundJobs
{
    private readonly IRepository<PolicyEntity> _repository;
    private readonly IMediator _mediator;
    private readonly ILogger<PolicyBackgroundJobs> _logger;

    public PolicyBackgroundJobs(
        IRepository<PolicyEntity> repository,
        IMediator mediator,
        ILogger<PolicyBackgroundJobs> logger)
    {
        _repository = repository;
        _mediator = mediator;
        _logger = logger;
    }

    /// <summary>
    /// FR-048: Automated policy status check for expiration/lapsing.
    /// Runs daily at midnight.
    /// </summary>
    public async Task ProcessAutoLapseAsync()
    {
        _logger.LogInformation("Starting Auto-Lapse background job at {Time}", DateTime.UtcNow);

        // Find ACTIVE policies where EndDate < Now
        var query = await _repository.FindAsync(p => p.Status == "ACTIVE" && p.EndDate < DateTime.UtcNow);
        var policiesToLapse = query.ToList();

        _logger.LogInformation("Found {Count} policies to lapse.", policiesToLapse.Count);

        foreach (var policy in policiesToLapse)
        {
            try
            {
                policy.Status = "LAPSED";
                policy.UpdatedAt = DateTime.UtcNow;
                await _repository.UpdateAsync(policy, default);
                
                _logger.LogInformation("Policy {PolicyNumber} auto-lapsed (Expired on {EndDate}).", policy.PolicyNumber, policy.EndDate);
                
                // Note: In Phase 4, we will add a Kafka event for PolicyLapsed.
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Failed to auto-lapse policy {PolicyId}", policy.PolicyId);
            }
        }
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
            // Logic to trigger a notification event (Phase 4)
            _logger.LogInformation("Renewal reminder triggered for policy {PolicyNumber} (Expires {EndDate}).", policy.PolicyNumber, policy.EndDate);
        }
    }
}
