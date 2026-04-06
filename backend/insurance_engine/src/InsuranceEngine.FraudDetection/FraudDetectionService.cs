using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;

namespace InsuranceEngine.FraudDetection;

public interface IFraudDetectionService
{
    Task<FraudCheckResult> CheckForFraudAsync(FraudCheckRequest request, CancellationToken ct = default);
    Task<ClaimPatternAnalysis> AnalyzeClaimPatternsAsync(string customerId, CancellationToken ct = default);
    Task<ProviderValidationResult> ValidateProviderAsync(string providerId, CancellationToken ct = default);
    Task<FraudDashboardSummary> GetDashboardSummaryAsync(string? assignedTo = null, CancellationToken ct = default);
    Task<List<FraudAlert>> GetPendingAlertsAsync(string? assignedTo = null, int limit = 100, CancellationToken ct = default);
    Task<bool> UpdateAlertStatusAsync(string alertId, string status, string? assignedTo = null, CancellationToken ct = default);
    Task<FraudIndicator> CheckRapidClaimFlagAsync(string customerId, DateTime? policyPurchaseDate, DateTime claimDate, CancellationToken ct = default);
    Task<FraudIndicator> CheckClaimFrequencyFlagAsync(string customerId, string claimType, int months = 12, CancellationToken ct = default);
    Task<FraudIndicator> CheckFullCoverageClaimFlagAsync(decimal claimAmount, decimal? policyCoverage, CancellationToken ct = default);
}

public class FraudDetectionService : IFraudDetectionService
{
    private readonly IFraudDetectionDataGateway _gateway;
    private readonly ILogger<FraudDetectionService> _logger;
    private readonly FraudCheckSettings _settings;
    private readonly List<FraudAlert> _alerts = new();
    private readonly object _alertsLock = new();

    private static readonly HashSet<string> ApprovedProviders = new(StringComparer.OrdinalIgnoreCase)
    {
        "provider-001", "provider-002", "provider-003",
        "apollo", "square", "united", "national", "bay", "labaid"
    };

    public FraudDetectionService(
        IFraudDetectionDataGateway gateway,
        IOptions<FraudCheckSettings> settings,
        ILogger<FraudDetectionService> logger)
    {
        _gateway = gateway;
        _settings = settings.Value;
        _logger = logger;
    }

    public async Task<FraudCheckResult> CheckForFraudAsync(FraudCheckRequest request, CancellationToken ct = default)
    {
        var result = new FraudCheckResult();
        var indicators = new List<FraudIndicator>();

        _logger.LogInformation("Performing fraud check for {EntityType}: {EntityId}", 
            request.EntityType, request.EntityId);

        if (request.ClaimId != null && request.CustomerId != null)
        {
            if (request.PolicyPurchaseDate.HasValue && request.ClaimSubmissionDate.HasValue)
            {
                var rapidClaimIndicator = await CheckRapidClaimFlagAsync(
                    request.CustomerId,
                    request.PolicyPurchaseDate,
                    request.ClaimSubmissionDate.Value,
                    ct);
                indicators.Add(rapidClaimIndicator);
            }

            if (!string.IsNullOrEmpty(request.ClaimType))
            {
                var frequencyIndicator = await CheckClaimFrequencyFlagAsync(
                    request.CustomerId,
                    request.ClaimType,
                    _settings.ClaimFrequencyWindowMonths,
                    ct);
                indicators.Add(frequencyIndicator);
            }

            if (request.ClaimAmount.HasValue)
            {
                var coverageIndicator = await CheckFullCoverageClaimFlagAsync(
                    request.ClaimAmount.Value,
                    request.PolicyCoverageAmount,
                    ct);
                indicators.Add(coverageIndicator);
            }
        }

        if (!string.IsNullOrEmpty(request.ProviderId) && _settings.EnableProviderValidation)
        {
            var providerResult = await ValidateProviderAsync(request.ProviderId, ct);
            if (!providerResult.IsApproved)
            {
                indicators.Add(new FraudIndicator
                {
                    Code = "FR-178",
                    Description = $"Non-approved network provider: {providerResult.ProviderId}",
                    ScoreContribution = 30,
                    Severity = "MEDIUM",
                    FrId = "FR-178"
                });
            }
        }

        if (!string.IsNullOrEmpty(request.DeviceFingerprint) && _settings.EnablePatternAnalysis)
        {
            var deviceIndicator = await CheckDeviceFingerprintAsync(request.DeviceFingerprint, ct);
            if (deviceIndicator != null)
            {
                indicators.Add(deviceIndicator);
            }
        }

        var fraudCheck = new InsuranceEngine.SharedKernel.Persistence.Entities.FraudCheckEntity
        {
            FraudCheckId = Guid.NewGuid().ToString(),
            EntityId = request.EntityId,
            EntityType = request.EntityType,
            CheckType = "COMPREHENSIVE",
            FraudScore = 0,
            RiskLevel = "LOW",
            Flagged = false,
            ClaimId = request.ClaimId,
            CustomerId = request.CustomerId,
            CheckedAt = DateTime.UtcNow,
            CreatedAt = DateTime.UtcNow
        };

        await _gateway.CreateFraudCheckAsync(fraudCheck, ct);

        result.Indicators = indicators;
        result.FraudScore = indicators.Sum(i => i.ScoreContribution);
        result.IsFraudDetected = result.FraudScore >= 50;
        result.RiskLevel = DetermineRiskLevel(result.FraudScore);
        result.Recommendation = GenerateRecommendation(result);

        if (result.IsFraudDetected)
        {
            CreateAlerts(request, result);
        }

        _logger.LogInformation(
            "Fraud check completed for {EntityId}: Score={Score}, Risk={Risk}, Detected={Detected}",
            request.EntityId, result.FraudScore, result.RiskLevel, result.IsFraudDetected);

        return result;
    }

    public async Task<ClaimPatternAnalysis> AnalyzeClaimPatternsAsync(string customerId, CancellationToken ct = default)
    {
        _logger.LogInformation("Analyzing claim patterns for customer: {CustomerId}", customerId);

        var recentClaims = await _gateway.GetRecentClaimsAsync(customerId, _settings.ClaimFrequencyWindowMonths, ct);

        var analysis = new ClaimPatternAnalysis
        {
            CustomerId = customerId,
            TotalClaimsInPeriod = recentClaims.Count,
            TotalClaimAmount = recentClaims.Sum(c => c.ClaimAmount ?? 0)
        };

        var claimTypeGroups = recentClaims
            .Where(c => !string.IsNullOrEmpty(c.ClaimType))
            .GroupBy(c => c.ClaimType!)
            .ToList();

        foreach (var group in claimTypeGroups)
        {
            var pattern = new ClaimPattern
            {
                ClaimType = group.Key,
                Count = group.Count(),
                TotalAmount = group.Sum(c => c.ClaimAmount ?? 0),
                ClaimIds = group.Where(c => !string.IsNullOrEmpty(c.ClaimId)).Select(c => c.ClaimId!).ToList()
            };

            analysis.Patterns.Add(pattern);

            if (pattern.Count > _settings.ClaimFrequencyThreshold)
            {
                analysis.HasSuspiciousPattern = true;
                analysis.SameTypeClaimsInPeriod += pattern.Count;
            }
        }

        return analysis;
    }

    public async Task<ProviderValidationResult> ValidateProviderAsync(string providerId, CancellationToken ct = default)
    {
        _logger.LogInformation("Validating provider: {ProviderId}", providerId);

        var normalizedId = providerId.ToLowerInvariant().Replace(" ", "");
        var isApproved = ApprovedProviders.Any(p => normalizedId.Contains(p.ToLowerInvariant()));

        var result = new ProviderValidationResult
        {
            ProviderId = providerId,
            IsApproved = isApproved,
            RequiresVerification = !isApproved,
            NetworkName = isApproved ? "InsureTech Network" : null,
            VerificationStatus = isApproved ? "APPROVED" : "PENDING_REVIEW"
        };

        return await Task.FromResult(result);
    }

    public Task<FraudDashboardSummary> GetDashboardSummaryAsync(string? assignedTo = null, CancellationToken ct = default)
    {
        List<FraudAlert> filteredAlerts;
        lock (_alertsLock)
        {
            filteredAlerts = string.IsNullOrEmpty(assignedTo)
                ? _alerts.ToList()
                : _alerts.Where(a => a.AssignedTo == assignedTo || a.AssignedTo == null).ToList();
        }

        var today = DateTime.UtcNow.Date;
        var todayAlerts = filteredAlerts.Where(a => a.CreatedAt.Date == today).ToList();

        var summary = new FraudDashboardSummary
        {
            TotalFlagsToday = todayAlerts.Count,
            HighRiskFlags = todayAlerts.Count(a => a.Severity == "HIGH"),
            MediumRiskFlags = todayAlerts.Count(a => a.Severity == "MEDIUM"),
            LowRiskFlags = todayAlerts.Count(a => a.Severity == "LOW"),
            PendingReviewCount = filteredAlerts.Count(a => a.Status == "PENDING"),
            RecentAlerts = filteredAlerts.OrderByDescending(a => a.CreatedAt).Take(20).ToList()
        };

        summary.FraudByType = todayAlerts
            .GroupBy(a => a.AlertType)
            .ToDictionary(g => g.Key, g => g.Count());

        return Task.FromResult(summary);
    }

    public Task<List<FraudAlert>> GetPendingAlertsAsync(string? assignedTo = null, int limit = 100, CancellationToken ct = default)
    {
        List<FraudAlert> result;
        lock (_alertsLock)
        {
            var query = _alerts.Where(a => a.Status == "PENDING");
            
            if (!string.IsNullOrEmpty(assignedTo))
            {
                query = query.Where(a => a.AssignedTo == assignedTo || a.AssignedTo == null);
            }

            result = query
                .OrderByDescending(a => a.CreatedAt)
                .Take(limit)
                .ToList();
        }

        return Task.FromResult(result);
    }

    public Task<bool> UpdateAlertStatusAsync(string alertId, string status, string? assignedTo = null, CancellationToken ct = default)
    {
        lock (_alertsLock)
        {
            var alert = _alerts.FirstOrDefault(a => a.AlertId == alertId);
            if (alert == null) return Task.FromResult(false);

            alert.Status = status;
            if (!string.IsNullOrEmpty(assignedTo))
            {
                alert.AssignedTo = assignedTo;
            }
        }

        _logger.LogInformation("Alert {AlertId} updated to status: {Status}", alertId, status);
        return Task.FromResult(true);
    }

    public async Task<FraudIndicator> CheckRapidClaimFlagAsync(
        string customerId,
        DateTime? policyPurchaseDate,
        DateTime claimDate,
        CancellationToken ct = default)
    {
        if (!policyPurchaseDate.HasValue)
        {
            return new FraudIndicator
            {
                Code = "FR-175_CHECK_SKIP",
                Description = "Policy purchase date not available",
                ScoreContribution = 0,
                Severity = "NONE"
            };
        }

        var hoursSincePurchase = (claimDate - policyPurchaseDate.Value).TotalHours;

        if (hoursSincePurchase < _settings.RapidClaimHoursThreshold)
        {
            _logger.LogWarning(
                "FR-175: Rapid claim detected for customer {CustomerId}: {Hours:F1} hours since purchase",
                customerId, hoursSincePurchase);

            return new FraudIndicator
            {
                Code = "FR-175",
                Description = $"Claim submitted {hoursSincePurchase:F1} hours after policy purchase (threshold: {_settings.RapidClaimHoursThreshold} hours)",
                ScoreContribution = 40,
                Severity = "HIGH",
                FrId = "FR-175"
            };
        }

        return new FraudIndicator
        {
            Code = "FR-175",
            Description = $"Claim submitted {hoursSincePurchase:F1} hours after policy purchase - no flag",
            ScoreContribution = 0,
            Severity = "LOW"
        };
    }

    public async Task<FraudIndicator> CheckClaimFrequencyFlagAsync(
        string customerId,
        string claimType,
        int months = 12,
        CancellationToken ct = default)
    {
        var analysis = await AnalyzeClaimPatternsAsync(customerId, ct);
        var pattern = analysis.Patterns.FirstOrDefault(p => 
            p.ClaimType.Equals(claimType, StringComparison.OrdinalIgnoreCase));

        if (pattern == null || pattern.Count <= _settings.ClaimFrequencyThreshold)
        {
            return new FraudIndicator
            {
                Code = "FR-176",
                Description = $"No suspicious frequency pattern for claim type '{claimType}'",
                ScoreContribution = 0,
                Severity = "LOW"
            };
        }

        _logger.LogWarning(
            "FR-176: Frequent claim pattern detected for customer {CustomerId}: {Count}x '{ClaimType}' in {Months} months",
            customerId, pattern.Count, claimType, months);

        var score = Math.Min(50, 10 + (pattern.Count - _settings.ClaimFrequencyThreshold) * 15);

        return new FraudIndicator
        {
            Code = "FR-176",
            Description = $"Same claim type '{claimType}' submitted {pattern.Count} times in {months} months (threshold: {_settings.ClaimFrequencyThreshold})",
            ScoreContribution = score,
            Severity = pattern.Count > 4 ? "HIGH" : "MEDIUM",
            FrId = "FR-176"
        };
    }

    public Task<FraudIndicator> CheckFullCoverageClaimFlagAsync(
        decimal claimAmount,
        decimal? policyCoverage,
        CancellationToken ct = default)
    {
        if (!policyCoverage.HasValue || policyCoverage.Value <= 0)
        {
            return Task.FromResult(new FraudIndicator
            {
                Code = "FR-177_CHECK_SKIP",
                Description = "Policy coverage amount not available",
                ScoreContribution = 0,
                Severity = "NONE"
            });
        }

        var coverageRatio = claimAmount / policyCoverage.Value;

        if (coverageRatio >= _settings.FullCoverageClaimThreshold)
        {
            _logger.LogWarning(
                "FR-177: Full coverage claim detected: {Amount}/{Coverage} = {Ratio:P0}",
                claimAmount, policyCoverage.Value, coverageRatio);

            return Task.FromResult(new FraudIndicator
            {
                Code = "FR-177",
                Description = $"Claim amount ({claimAmount:N0}) equals or exceeds policy coverage ({policyCoverage.Value:N0})",
                ScoreContribution = 35,
                Severity = "HIGH",
                FrId = "FR-177"
            });
        }

        return Task.FromResult(new FraudIndicator
        {
            Code = "FR-177",
            Description = $"Claim is {coverageRatio:P0} of coverage - no flag",
            ScoreContribution = 0,
            Severity = "LOW"
        });
    }

    private async Task<FraudIndicator?> CheckDeviceFingerprintAsync(string deviceFingerprint, CancellationToken ct)
    {
        if (string.IsNullOrEmpty(deviceFingerprint))
            return null;

        return await Task.FromResult<FraudIndicator?>(null);
    }

    private void CreateAlerts(FraudCheckRequest request, FraudCheckResult result)
    {
        lock (_alertsLock)
        {
            foreach (var indicator in result.Indicators.Where(i => i.ScoreContribution > 0))
            {
                var alert = new FraudAlert
                {
                    AlertId = $"alert_{Guid.NewGuid():N}",
                    EntityId = request.EntityId,
                    EntityType = request.EntityType,
                    AlertType = indicator.Code,
                    Severity = indicator.Severity,
                    Status = "PENDING",
                    Description = indicator.Description,
                    CreatedAt = DateTime.UtcNow,
                    AssignedTo = DetermineAssignedRole(indicator)
                };

                _alerts.Add(alert);

                _logger.LogInformation(
                    "Created fraud alert: {AlertId} for {EntityType}:{EntityId} - {Description}",
                    alert.AlertId, alert.EntityType, alert.EntityId, alert.Description);
            }
        }
    }

    private string DetermineAssignedRole(FraudIndicator indicator)
    {
        return indicator.Severity switch
        {
            "HIGH" => "ClaimsOfficer",
            "MEDIUM" => "FocalPerson",
            _ => "System"
        };
    }

    private string DetermineRiskLevel(int score)
    {
        return score switch
        {
            >= 70 => "HIGH",
            >= 40 => "MEDIUM",
            _ => "LOW"
        };
    }

    private string GenerateRecommendation(FraudCheckResult result)
    {
        if (!result.IsFraudDetected)
            return "Proceed with normal processing";

        var highSeverityIndicators = result.Indicators
            .Where(i => i.Severity == "HIGH")
            .Select(i => i.Code)
            .ToList();

        if (highSeverityIndicators.Contains("FR-175"))
            return "MANUAL_REVIEW: Claim within 48 hours of purchase - requires supervisor approval";
        
        if (highSeverityIndicators.Contains("FR-176"))
            return "MANUAL_REVIEW: Frequent claim pattern detected - verify claim legitimacy";
        
        if (highSeverityIndicators.Contains("FR-177"))
            return "MANUAL_REVIEW: Full coverage claim - additional documentation required";

        return "ELEVATED_MONITORING: Multiple fraud indicators detected";
    }
}

public class MockFraudDetectionService : IFraudDetectionService
{
    private readonly ILogger<MockFraudDetectionService> _logger;

    public MockFraudDetectionService(ILogger<MockFraudDetectionService> logger)
    {
        _logger = logger;
    }

    public Task<FraudCheckResult> CheckForFraudAsync(FraudCheckRequest request, CancellationToken ct = default)
    {
        _logger.LogInformation("[MOCK] Fraud check for {EntityType}: {EntityId}", 
            request.EntityType, request.EntityId);

        return Task.FromResult(new FraudCheckResult
        {
            IsFraudDetected = false,
            FraudScore = 0,
            RiskLevel = "LOW",
            Recommendation = "Proceed with normal processing"
        });
    }

    public Task<ClaimPatternAnalysis> AnalyzeClaimPatternsAsync(string customerId, CancellationToken ct = default)
    {
        _logger.LogInformation("[MOCK] Claim pattern analysis for {CustomerId}", customerId);
        return Task.FromResult(new ClaimPatternAnalysis { CustomerId = customerId });
    }

    public Task<ProviderValidationResult> ValidateProviderAsync(string providerId, CancellationToken ct = default)
    {
        _logger.LogInformation("[MOCK] Provider validation for {ProviderId}", providerId);
        return Task.FromResult(new ProviderValidationResult
        {
            ProviderId = providerId,
            IsApproved = true,
            NetworkName = "Mock Network"
        });
    }

    public Task<FraudDashboardSummary> GetDashboardSummaryAsync(string? assignedTo = null, CancellationToken ct = default)
    {
        return Task.FromResult(new FraudDashboardSummary());
    }

    public Task<List<FraudAlert>> GetPendingAlertsAsync(string? assignedTo = null, int limit = 100, CancellationToken ct = default)
    {
        return Task.FromResult(new List<FraudAlert>());
    }

    public Task<bool> UpdateAlertStatusAsync(string alertId, string status, string? assignedTo = null, CancellationToken ct = default)
    {
        return Task.FromResult(true);
    }

    public Task<FraudIndicator> CheckRapidClaimFlagAsync(string customerId, DateTime? policyPurchaseDate, DateTime claimDate, CancellationToken ct = default)
    {
        return Task.FromResult(new FraudIndicator { Code = "FR-175", ScoreContribution = 0, Severity = "LOW" });
    }

    public Task<FraudIndicator> CheckClaimFrequencyFlagAsync(string customerId, string claimType, int months = 12, CancellationToken ct = default)
    {
        return Task.FromResult(new FraudIndicator { Code = "FR-176", ScoreContribution = 0, Severity = "LOW" });
    }

    public Task<FraudIndicator> CheckFullCoverageClaimFlagAsync(decimal claimAmount, decimal? policyCoverage, CancellationToken ct = default)
    {
        return Task.FromResult(new FraudIndicator { Code = "FR-177", ScoreContribution = 0, Severity = "LOW" });
    }
}

