using FluentAssertions;
using InsuranceEngine.FraudDetection;
using Insuretech.Fraud.Services.V1;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using Moq;
using Xunit;

namespace InsuranceEngine.Infrastructure.Tests;

public class FraudDetectionServiceTests
{
    private readonly Mock<IFraudDetectionDataGateway> _gatewayMock;
    private readonly Mock<ILogger<FraudDetectionService>> _loggerMock;
    private readonly IOptions<FraudCheckSettings> _settings;
    private readonly FraudDetectionService _service;

    public FraudDetectionServiceTests()
    {
        _gatewayMock = new Mock<IFraudDetectionDataGateway>();
        _loggerMock = new Mock<ILogger<FraudDetectionService>>();
        _settings = Options.Create(new FraudCheckSettings
        {
            RapidClaimHoursThreshold = 48,
            ClaimFrequencyThreshold = 2,
            ClaimFrequencyWindowMonths = 12,
            FullCoverageClaimThreshold = 1.0m,
            DeviceAccountThreshold = 3,
            EnablePatternAnalysis = true,
            EnableProviderValidation = true
        });

        _service = new FraudDetectionService(_gatewayMock.Object, _settings, _loggerMock.Object);
    }

    [Fact]
    public async Task CheckRapidClaimFlagAsync_ShouldFlag_WhenClaimWithin48Hours()
    {
        var customerId = "cust-123";
        var purchaseDate = DateTime.UtcNow.AddHours(-24);
        var claimDate = DateTime.UtcNow;

        var result = await _service.CheckRapidClaimFlagAsync(customerId, purchaseDate, claimDate);

        result.Code.Should().Be("FR-175");
        result.ScoreContribution.Should().BeGreaterThan(0);
        result.Severity.Should().Be("HIGH");
    }

    [Fact]
    public async Task CheckRapidClaimFlagAsync_ShouldNotFlag_WhenClaimAfter48Hours()
    {
        var customerId = "cust-123";
        var purchaseDate = DateTime.UtcNow.AddHours(-72);
        var claimDate = DateTime.UtcNow;

        var result = await _service.CheckRapidClaimFlagAsync(customerId, purchaseDate, claimDate);

        result.Code.Should().Be("FR-175");
        result.ScoreContribution.Should().Be(0);
        result.Severity.Should().Be("LOW");
    }

    [Fact]
    public async Task CheckRapidClaimFlagAsync_ShouldSkip_WhenPurchaseDateNull()
    {
        var customerId = "cust-123";

        var result = await _service.CheckRapidClaimFlagAsync(customerId, null, DateTime.UtcNow);

        result.Code.Should().Be("FR-175_CHECK_SKIP");
        result.ScoreContribution.Should().Be(0);
    }

    [Fact]
    public async Task CheckFullCoverageClaimFlagAsync_ShouldFlag_WhenClaimEqualsCoverage()
    {
        var claimAmount = 100000m;
        var policyCoverage = 100000m;

        var result = await _service.CheckFullCoverageClaimFlagAsync(claimAmount, policyCoverage);

        result.Code.Should().Be("FR-177");
        result.ScoreContribution.Should().Be(35);
        result.Severity.Should().Be("HIGH");
    }

    [Fact]
    public async Task CheckFullCoverageClaimFlagAsync_ShouldFlag_WhenClaimExceedsCoverage()
    {
        var claimAmount = 120000m;
        var policyCoverage = 100000m;

        var result = await _service.CheckFullCoverageClaimFlagAsync(claimAmount, policyCoverage);

        result.Code.Should().Be("FR-177");
        result.ScoreContribution.Should().Be(35);
        result.Severity.Should().Be("HIGH");
    }

    [Fact]
    public async Task CheckFullCoverageClaimFlagAsync_ShouldNotFlag_WhenClaimBelowCoverage()
    {
        var claimAmount = 50000m;
        var policyCoverage = 100000m;

        var result = await _service.CheckFullCoverageClaimFlagAsync(claimAmount, policyCoverage);

        result.Code.Should().Be("FR-177");
        result.ScoreContribution.Should().Be(0);
    }

    [Fact]
    public async Task CheckFullCoverageClaimFlagAsync_ShouldSkip_WhenCoverageNull()
    {
        var result = await _service.CheckFullCoverageClaimFlagAsync(100000m, null);

        result.Code.Should().Be("FR-177_CHECK_SKIP");
    }

    [Fact]
    public async Task ValidateProviderAsync_ShouldApprove_KnownProviders()
    {
        var providerId = "apollo";

        var result = await _service.ValidateProviderAsync(providerId);

        result.IsApproved.Should().BeTrue();
        result.ProviderId.Should().Be(providerId);
    }

    [Fact]
    public async Task ValidateProviderAsync_ShouldNotApprove_UnknownProviders()
    {
        var providerId = "unknown-provider";

        var result = await _service.ValidateProviderAsync(providerId);

        result.IsApproved.Should().BeFalse();
        result.ProviderId.Should().Be(providerId);
    }

    [Fact]
    public async Task GetDashboardSummaryAsync_ShouldReturnEmptySummary_WhenNoAlerts()
    {
        var result = await _service.GetDashboardSummaryAsync();

        result.TotalFlagsToday.Should().Be(0);
        result.HighRiskFlags.Should().Be(0);
        result.PendingReviewCount.Should().Be(0);
    }

    [Fact]
    public async Task UpdateAlertStatusAsync_ShouldReturnFalse_WhenAlertNotFound()
    {
        var alertId = "non-existent-alert";
        var status = "RESOLVED";

        var result = await _service.UpdateAlertStatusAsync(alertId, status);

        result.Should().BeFalse();
    }

    [Fact]
    public async Task CheckForFraudAsync_ShouldCombineIndicators()
    {
        var request = new FraudCheckRequest
        {
            EntityId = "claim-123",
            EntityType = "claim",
            CustomerId = "cust-123",
            ClaimId = "claim-123",
            PolicyPurchaseDate = DateTime.UtcNow.AddHours(-24),
            ClaimSubmissionDate = DateTime.UtcNow,
            ClaimAmount = 100000m,
            PolicyCoverageAmount = 100000m,
            ProviderId = "apollo"
        };

        _gatewayMock.Setup(g => g.CheckFraudAsync(It.IsAny<Insuretech.Fraud.Services.V1.CheckFraudRequest>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new CheckFraudResponse { IsFraudDetected = false });

        var result = await _service.CheckForFraudAsync(request);

        result.Indicators.Should().NotBeEmpty();
        result.FraudScore.Should().BeGreaterThan(0);
    }

    [Fact]
    public async Task CheckForFraudAsync_ShouldFlagRapidClaim()
    {
        var request = new FraudCheckRequest
        {
            EntityId = "claim-123",
            EntityType = "claim",
            CustomerId = "cust-123",
            ClaimId = "claim-123",
            PolicyPurchaseDate = DateTime.UtcNow.AddHours(-12),
            ClaimSubmissionDate = DateTime.UtcNow
        };

        _gatewayMock.Setup(g => g.CheckFraudAsync(It.IsAny<Insuretech.Fraud.Services.V1.CheckFraudRequest>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new CheckFraudResponse { IsFraudDetected = false });

        var result = await _service.CheckForFraudAsync(request);

        result.Indicators.Should().Contain(i => i.Code == "FR-175");
    }

    [Fact]
    public async Task CheckForFraudAsync_ShouldFlagFullCoverageClaim()
    {
        var request = new FraudCheckRequest
        {
            EntityId = "claim-123",
            EntityType = "claim",
            CustomerId = "cust-123",
            ClaimId = "claim-123",
            ClaimAmount = 100000m,
            PolicyCoverageAmount = 100000m
        };

        _gatewayMock.Setup(g => g.CheckFraudAsync(It.IsAny<Insuretech.Fraud.Services.V1.CheckFraudRequest>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new CheckFraudResponse { IsFraudDetected = false });

        var result = await _service.CheckForFraudAsync(request);

        result.Indicators.Should().Contain(i => i.Code == "FR-177");
    }

    [Fact]
    public async Task CheckForFraudAsync_ShouldFlagNonApprovedProvider()
    {
        var request = new FraudCheckRequest
        {
            EntityId = "claim-123",
            EntityType = "claim",
            CustomerId = "cust-123",
            ClaimId = "claim-123",
            ProviderId = "unknown-clinic"
        };

        _gatewayMock.Setup(g => g.CheckFraudAsync(It.IsAny<Insuretech.Fraud.Services.V1.CheckFraudRequest>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new CheckFraudResponse { IsFraudDetected = false });

        var result = await _service.CheckForFraudAsync(request);

        result.Indicators.Should().Contain(i => i.Code == "FR-178");
    }
}
