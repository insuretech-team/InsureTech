using Xunit;
using Moq;
using FluentAssertions;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Logging.Abstractions;
using InsuranceEngine.Infrastructure.Renewals;
using InsuranceEngine.Infrastructure.Notifications;
using InsuranceEngine.Infrastructure.Messaging;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using InsuranceEngine.Renewals;
using Insuretech.Policy.Services.V1;

namespace InsuranceEngine.Infrastructure.Tests;

public class GracePeriodServiceTests
{
    private readonly Mock<IRepository<PolicyEntity>> _policyRepoMock;
    private readonly Mock<IRenewalDataGateway> _renewalGatewayMock;
    private readonly Mock<INotificationService> _notificationServiceMock;
    private readonly Mock<IEventPublisher> _eventPublisherMock;
    private readonly GracePeriodService _service;

    public GracePeriodServiceTests()
    {
        _policyRepoMock = new Mock<IRepository<PolicyEntity>>();
        _renewalGatewayMock = new Mock<IRenewalDataGateway>();
        _notificationServiceMock = new Mock<INotificationService>();
        _eventPublisherMock = new Mock<IEventPublisher>();
        
        var settings = Microsoft.Extensions.Options.Options.Create(new GracePeriodSettings
        {
            GracePeriodDays = 30,
            ReinstatementWindowDays = 90,
            ReinstatementPenaltyPercent = 10.0m,
            EnableDailyReminders = true
        });

        _service = new GracePeriodService(
            _policyRepoMock.Object,
            _renewalGatewayMock.Object,
            _notificationServiceMock.Object,
            _eventPublisherMock.Object,
            settings,
            NullLogger<GracePeriodService>.Instance);
    }

    [Fact]
    public async Task GetGracePeriodInfoAsync_ShouldReturnNull_WhenPolicyNotFound()
    {
        _policyRepoMock.Setup(r => r.FindAsync(
            It.IsAny<System.Linq.Expressions.Expression<Func<PolicyEntity, bool>>>(),
            It.IsAny<CancellationToken>()))
            .ReturnsAsync(new List<PolicyEntity>());

        var result = await _service.GetGracePeriodInfoAsync(Guid.NewGuid().ToString());

        result.Should().BeNull();
    }

    [Fact]
    public async Task GetGracePeriodInfoAsync_ShouldReturnGracePeriodInfo_WhenPolicyInGracePeriod()
    {
        var policyId = Guid.NewGuid();
        var policy = new PolicyEntity
        {
            PolicyId = policyId,
            PolicyNumber = "POL-001",
            Status = "GRACE_PERIOD",
            CustomerId = Guid.NewGuid(),
            PremiumAmount = 10000,
            UnderwritingData = "{\"GracePeriodEndDate\":\"" + DateTime.UtcNow.AddDays(15).ToString("O") + "\"}"
        };

        _policyRepoMock.Setup(r => r.FindAsync(
            It.IsAny<System.Linq.Expressions.Expression<Func<PolicyEntity, bool>>>(),
            It.IsAny<CancellationToken>()))
            .ReturnsAsync(new List<PolicyEntity> { policy });

        var result = await _service.GetGracePeriodInfoAsync(policyId.ToString());

        result.Should().NotBeNull();
        result!.PolicyNumber.Should().Be("POL-001");
        result.Status.Should().Be("GRACE_PERIOD");
        result.DaysRemaining.Should().BeGreaterThan(0);
        result.CanRenew.Should().BeTrue();
    }

    [Fact]
    public async Task GetGracePeriodInfoAsync_ShouldReturnNull_WhenPolicyIsActive()
    {
        var policyId = Guid.NewGuid();
        var policy = new PolicyEntity
        {
            PolicyId = policyId,
            PolicyNumber = "POL-001",
            Status = "ACTIVE",
            CustomerId = Guid.NewGuid(),
            PremiumAmount = 10000
        };

        _policyRepoMock.Setup(r => r.FindAsync(
            It.IsAny<System.Linq.Expressions.Expression<Func<PolicyEntity, bool>>>(),
            It.IsAny<CancellationToken>()))
            .ReturnsAsync(new List<PolicyEntity> { policy });

        var result = await _service.GetGracePeriodInfoAsync(policyId.ToString());

        result.Should().BeNull();
    }

    [Fact]
    public async Task CanPolicyBeReinstatedAsync_ShouldReturnFalse_WhenPolicyNotFound()
    {
        _policyRepoMock.Setup(r => r.FindAsync(
            It.IsAny<System.Linq.Expressions.Expression<Func<PolicyEntity, bool>>>(),
            It.IsAny<CancellationToken>()))
            .ReturnsAsync(new List<PolicyEntity>());

        var result = await _service.CanPolicyBeReinstatedAsync(Guid.NewGuid().ToString());

        result.Should().BeFalse();
    }

    [Fact]
    public async Task CanPolicyBeReinstatedAsync_ShouldReturnTrue_WhenPolicyIsLapsedAndWithinWindow()
    {
        var policyId = Guid.NewGuid();
        var policy = new PolicyEntity
        {
            PolicyId = policyId,
            PolicyNumber = "POL-001",
            Status = "LAPSED",
            CustomerId = Guid.NewGuid(),
            UnderwritingData = "{\"ReinstatementWindowEndDate\":\"" + DateTime.UtcNow.AddDays(30).ToString("O") + "\"}"
        };

        _policyRepoMock.Setup(r => r.FindAsync(
            It.IsAny<System.Linq.Expressions.Expression<Func<PolicyEntity, bool>>>(),
            It.IsAny<CancellationToken>()))
            .ReturnsAsync(new List<PolicyEntity> { policy });

        var result = await _service.CanPolicyBeReinstatedAsync(policyId.ToString());

        result.Should().BeTrue();
    }

    [Fact]
    public async Task CanPolicyBeReinstatedAsync_ShouldReturnFalse_WhenPolicyIsActive()
    {
        var policyId = Guid.NewGuid();
        var policy = new PolicyEntity
        {
            PolicyId = policyId,
            PolicyNumber = "POL-001",
            Status = "ACTIVE",
            CustomerId = Guid.NewGuid()
        };

        _policyRepoMock.Setup(r => r.FindAsync(
            It.IsAny<System.Linq.Expressions.Expression<Func<PolicyEntity, bool>>>(),
            It.IsAny<CancellationToken>()))
            .ReturnsAsync(new List<PolicyEntity> { policy });

        var result = await _service.CanPolicyBeReinstatedAsync(policyId.ToString());

        result.Should().BeFalse();
    }

    [Fact]
    public async Task ReinstatePolicyAsync_ShouldReturnError_WhenPolicyNotFound()
    {
        _policyRepoMock.Setup(r => r.FindAsync(
            It.IsAny<System.Linq.Expressions.Expression<Func<PolicyEntity, bool>>>(),
            It.IsAny<CancellationToken>()))
            .ReturnsAsync(new List<PolicyEntity>());

        var request = new ReinstatementRequest
        {
            PolicyId = Guid.NewGuid().ToString(),
            TenureMonths = 12
        };

        var result = await _service.ReinstatePolicyAsync(request);

        result.Success.Should().BeFalse();
        result.ErrorMessage.Should().Contain("not found");
    }

    [Fact]
    public async Task ReinstatePolicyAsync_ShouldReturnError_WhenInvalidPolicyId()
    {
        var request = new ReinstatementRequest
        {
            PolicyId = "not-a-guid",
            TenureMonths = 12
        };

        var result = await _service.ReinstatePolicyAsync(request);

        result.Success.Should().BeFalse();
        result.ErrorMessage.Should().Contain("Invalid policy ID");
    }

    [Fact]
    public async Task ReinstatePolicyAsync_ShouldReturnError_WhenOutsideReinstatementWindow()
    {
        var policyId = Guid.NewGuid();
        var policy = new PolicyEntity
        {
            PolicyId = policyId,
            PolicyNumber = "POL-001",
            Status = "LAPSED",
            CustomerId = Guid.NewGuid(),
            UnderwritingData = "{\"ReinstatementWindowEndDate\":\"" + DateTime.UtcNow.AddDays(-1).ToString("O") + "\"}"
        };

        _policyRepoMock.Setup(r => r.FindAsync(
            It.IsAny<System.Linq.Expressions.Expression<Func<PolicyEntity, bool>>>(),
            It.IsAny<CancellationToken>()))
            .ReturnsAsync(new List<PolicyEntity> { policy });

        var request = new ReinstatementRequest
        {
            PolicyId = policyId.ToString(),
            TenureMonths = 12
        };

        var result = await _service.ReinstatePolicyAsync(request);

        result.Success.Should().BeFalse();
        result.ErrorMessage.Should().Contain("reinstatement window");
    }

    [Fact]
    public async Task ReinstatePolicyAsync_ShouldSucceed_WhenWithinWindow()
    {
        var policyId = Guid.NewGuid();
        var policy = new PolicyEntity
        {
            PolicyId = policyId,
            PolicyNumber = "POL-001",
            Status = "LAPSED",
            CustomerId = Guid.NewGuid(),
            PremiumAmount = 10000,
            UnderwritingData = "{\"ReinstatementWindowEndDate\":\"" + DateTime.UtcNow.AddDays(30).ToString("O") + "\"}"
        };

        _policyRepoMock.Setup(r => r.FindAsync(
            It.IsAny<System.Linq.Expressions.Expression<Func<PolicyEntity, bool>>>(),
            It.IsAny<CancellationToken>()))
            .ReturnsAsync(new List<PolicyEntity> { policy });

        _renewalGatewayMock.Setup(g => g.RenewPolicyAsync(
            It.IsAny<RenewPolicyTenureRequest>(),
            It.IsAny<CancellationToken>()))
            .ReturnsAsync(new RenewPolicyTenureResponse
            {
                NewPolicyId = Guid.NewGuid().ToString(),
                NewPolicyNumber = "POL-002"
            });

        var request = new ReinstatementRequest
        {
            PolicyId = policyId.ToString(),
            TenureMonths = 12,
            ApplyReinstatementPenalty = true
        };

        var result = await _service.ReinstatePolicyAsync(request);

        result.Success.Should().BeTrue();
        result.NewPolicyNumber.Should().Be("POL-002");
        result.ReinstatementPenalty.Should().BeGreaterThan(0);
    }

    [Fact]
    public async Task ReinstatePolicyAsync_ShouldCalculatePenalty_WhenApplyReinstatementPenaltyIsTrue()
    {
        var policyId = Guid.NewGuid();
        var policy = new PolicyEntity
        {
            PolicyId = policyId,
            PolicyNumber = "POL-001",
            Status = "LAPSED",
            CustomerId = Guid.NewGuid(),
            PremiumAmount = 1000000,
            UnderwritingData = "{\"ReinstatementWindowEndDate\":\"" + DateTime.UtcNow.AddDays(30).ToString("O") + "\"}"
        };

        _policyRepoMock.Setup(r => r.FindAsync(
            It.IsAny<System.Linq.Expressions.Expression<Func<PolicyEntity, bool>>>(),
            It.IsAny<CancellationToken>()))
            .ReturnsAsync(new List<PolicyEntity> { policy });

        _renewalGatewayMock.Setup(g => g.RenewPolicyAsync(
            It.IsAny<RenewPolicyTenureRequest>(),
            It.IsAny<CancellationToken>()))
            .ReturnsAsync(new RenewPolicyTenureResponse
            {
                NewPolicyId = Guid.NewGuid().ToString(),
                NewPolicyNumber = "POL-002"
            });

        var request = new ReinstatementRequest
        {
            PolicyId = policyId.ToString(),
            TenureMonths = 12,
            ApplyReinstatementPenalty = true
        };

        var result = await _service.ReinstatePolicyAsync(request);

        result.Success.Should().BeTrue();
        result.ReinstatementPenalty.Should().BeGreaterThan(0);
    }

    [Fact]
    public async Task ReinstatePolicyAsync_ShouldNotCalculatePenalty_WhenApplyReinstatementPenaltyIsFalse()
    {
        var policyId = Guid.NewGuid();
        var policy = new PolicyEntity
        {
            PolicyId = policyId,
            PolicyNumber = "POL-001",
            Status = "LAPSED",
            CustomerId = Guid.NewGuid(),
            PremiumAmount = 1000000,
            UnderwritingData = "{\"ReinstatementWindowEndDate\":\"" + DateTime.UtcNow.AddDays(30).ToString("O") + "\"}"
        };

        _policyRepoMock.Setup(r => r.FindAsync(
            It.IsAny<System.Linq.Expressions.Expression<Func<PolicyEntity, bool>>>(),
            It.IsAny<CancellationToken>()))
            .ReturnsAsync(new List<PolicyEntity> { policy });

        _renewalGatewayMock.Setup(g => g.RenewPolicyAsync(
            It.IsAny<RenewPolicyTenureRequest>(),
            It.IsAny<CancellationToken>()))
            .ReturnsAsync(new RenewPolicyTenureResponse
            {
                NewPolicyId = Guid.NewGuid().ToString(),
                NewPolicyNumber = "POL-002"
            });

        var request = new ReinstatementRequest
        {
            PolicyId = policyId.ToString(),
            TenureMonths = 12,
            ApplyReinstatementPenalty = false
        };

        var result = await _service.ReinstatePolicyAsync(request);

        result.Success.Should().BeTrue();
        result.ReinstatementPenalty.Should().Be(0);
    }
}
