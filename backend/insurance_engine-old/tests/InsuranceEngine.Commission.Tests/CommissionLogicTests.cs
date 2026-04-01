using System;
using System.Threading;
using System.Threading.Tasks;
using System.Reflection;
using FluentAssertions;
using InsuranceEngine.Commission.Application.Features.Events;
using InsuranceEngine.Commission.Domain.Enums;
using InsuranceEngine.Commission.Domain.Interfaces;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.Policy.Domain.Entities;
using InsuranceEngine.Policy.Domain.Events;
using Microsoft.Extensions.Logging;
using Moq;
using Xunit;

namespace InsuranceEngine.Commission.Tests;

public class CommissionLogicTests
{
    private readonly Mock<ICommissionRepository> _commissionRepoMock;
    private readonly Mock<IPolicyRepository> _policyRepoMock;
    private readonly Mock<ILogger<PolicyIssuedEventHandler>> _issuedLoggerMock;
    private readonly Mock<ILogger<PolicyRenewedEventHandler>> _renewedLoggerMock;

    public CommissionLogicTests()
    {
        _commissionRepoMock = new Mock<ICommissionRepository>();
        _policyRepoMock = new Mock<IPolicyRepository>();
        _issuedLoggerMock = new Mock<ILogger<PolicyIssuedEventHandler>>();
        _renewedLoggerMock = new Mock<ILogger<PolicyRenewedEventHandler>>();
    }

    [Fact]
    public async Task PolicyIssued_WithPartner_ShouldCreate15PercentCommission()
    {
        // Arrange
        var policyId = Guid.NewGuid();
        var partnerId = Guid.NewGuid();
        var premium = 100_000L; // 1000 BDT
        
        var policy = PolicyAggregate.Create("POL-001", Guid.NewGuid(), Guid.NewGuid(), partnerId, 1000000L, premium, 12, DateTime.UtcNow);
        
        _policyRepoMock.Setup(r => r.GetByIdAsync(policyId, It.IsAny<CancellationToken>()))
            .ReturnsAsync(policy);

        var handler = new PolicyIssuedEventHandler(_commissionRepoMock.Object, _policyRepoMock.Object, _issuedLoggerMock.Object);
        var notification = new PolicyIssuedEvent(policyId, "POL-001", premium, DateTime.UtcNow);

        // Act
        await handler.Handle(notification, CancellationToken.None);

        // Assert
        _commissionRepoMock.Verify(r => r.CreateAsync(
            It.Is<Domain.Entities.Commission>(c => 
                c.PolicyId == policyId && 
                c.PartnerId == partnerId && 
                c.Type == CommissionType.Acquisition &&
                c.Amount == 15_000L), // 15% of 100,000
            It.IsAny<CancellationToken>()), Times.Once);
    }

    [Fact]
    public async Task PolicyRenewed_WithAgent_ShouldCreate5PercentCommission()
    {
        // Arrange
        var newPolicyId = Guid.NewGuid();
        var agentId = Guid.NewGuid();
        var premium = 100_000L;

        // Using reflection to set AgentId as it's private set and Renew might not set it directly in this test setup
        var policy = PolicyAggregate.Create("POL-Renew", Guid.NewGuid(), Guid.NewGuid(), null, 1000000L, premium, 12, DateTime.UtcNow);
        var agentField = typeof(PolicyAggregate).GetProperty("AgentId");
        agentField?.SetValue(policy, agentId);

        _policyRepoMock.Setup(r => r.GetByIdAsync(newPolicyId, It.IsAny<CancellationToken>()))
            .ReturnsAsync(policy);

        var handler = new PolicyRenewedEventHandler(_commissionRepoMock.Object, _policyRepoMock.Object, _renewedLoggerMock.Object);
        var notification = new PolicyRenewedEvent(Guid.NewGuid(), newPolicyId, "POL-Renew", premium, DateTime.UtcNow);

        // Act
        await handler.Handle(notification, CancellationToken.None);

        // Assert
        _commissionRepoMock.Verify(r => r.CreateAsync(
            It.Is<Domain.Entities.Commission>(c => 
                c.PolicyId == newPolicyId && 
                c.AgentId == agentId && 
                c.Type == CommissionType.Renewal &&
                c.Amount == 5_000L), // 5% of 100,000
            It.IsAny<CancellationToken>()), Times.Once);
    }
}
