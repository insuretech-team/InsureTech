using Moq;
using Xunit;
using FluentAssertions;
using Microsoft.Extensions.Logging;
using InsuranceEngine.Policy.Application.Commands;
using InsuranceEngine.Policy.Domain;
using InsuranceEngine.SharedKernel.Infrastructure;
using Insuretech.Policy.Services.V1;
using Insuretech.Policy.Entity.V1;

namespace InsuranceEngine.Policy.Tests;

public class CreatePolicyCommandHandlerTests
{
    private readonly Mock<IPolicyDataGateway> _gatewayMock;
    private readonly Mock<ILogger<CreatePolicyCommandHandler>> _loggerMock;
    private readonly CreatePolicyCommandHandler _handler;

    public CreatePolicyCommandHandlerTests()
    {
        _gatewayMock = new Mock<IPolicyDataGateway>();
        _loggerMock = new Mock<ILogger<CreatePolicyCommandHandler>>();
        _handler = new CreatePolicyCommandHandler(
            _gatewayMock.Object,
            _loggerMock.Object);
    }

    [Fact]
    public async Task Handle_ShouldReturnError_WhenGatewayFails()
    {
        var command = new CreatePolicyCommand(
            Guid.NewGuid().ToString(),
            Guid.NewGuid().ToString(),
            null, null, null, 1000, 50000, 12, DateTime.UtcNow, null, null);

        _gatewayMock.Setup(g => g.CreatePolicyAsync(It.IsAny<CreatePolicyRequest>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new CreatePolicyResponse
            {
                Error = new Insuretech.Common.V1.Error { Code = "FAILED", Message = "Test error" }
            });

        var response = await _handler.Handle(command, CancellationToken.None);

        response.Error.Should().NotBeNull();
    }

    [Fact]
    public async Task Handle_ShouldCreatePolicy_WhenGatewaySucceeds()
    {
        var productId = Guid.NewGuid().ToString();
        var customerId = Guid.NewGuid().ToString();
        
        var command = new CreatePolicyCommand(
            productId,
            customerId,
            null, null, null, 1000, 50000, 12, DateTime.UtcNow, null, null);

        _gatewayMock.Setup(g => g.CreatePolicyAsync(It.IsAny<CreatePolicyRequest>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new CreatePolicyResponse
            {
                PolicyId = Guid.NewGuid().ToString(),
                PolicyNumber = "POL-123",
                Message = "Success"
            });

        var response = await _handler.Handle(command, CancellationToken.None);

        response.Error.Should().BeNull();
        response.PolicyNumber.Should().Be("POL-123");
    }

    [Theory]
    [InlineData(1239, 50000, "14-1 days, Age 0-40, Plan A - Non Schengen")]
    [InlineData(2499, 50000, "14-1 days, Age 51-55, Plan A - Non Schengen")]
    [InlineData(8131, 50000, "14-1 days, Age 60-65, Plan A - Non Schengen")]
    public async Task Handle_ShouldCreatePolicy_WithPremiumFromExcel(decimal expectedPremium, decimal sumInsured, string scenario)
    {
        var productId = Guid.NewGuid().ToString();
        var customerId = Guid.NewGuid().ToString();
        
        var command = new CreatePolicyCommand(
            productId,
            customerId,
            null, null, null, expectedPremium, sumInsured, 1, DateTime.UtcNow, null, null);

        _gatewayMock.Setup(g => g.CreatePolicyAsync(It.IsAny<CreatePolicyRequest>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new CreatePolicyResponse
            {
                PolicyId = Guid.NewGuid().ToString(),
                PolicyNumber = "POL-2026-0001",
                Message = "Success"
            });

        var response = await _handler.Handle(command, CancellationToken.None);

        response.Error.Should().BeNull();
        response.PolicyNumber.Should().NotBeEmpty();
    }

    [Fact]
    public async Task Handle_ShouldCreatePolicy_WithNomineesFromExcelData()
    {
        var productId = Guid.NewGuid().ToString();
        var customerId = Guid.NewGuid().ToString();
        
        var nominees = new List<Nominee>
        {
            new Nominee
            {
                NomineeId = Guid.NewGuid().ToString(),
                FullName = "John Doe",
                Relationship = "Spouse",
                SharePercentage = 50,
                DateOfBirth = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(new DateTime(1990, 1, 1).ToUniversalTime()),
                NidNumber = "1234567890",
                PhoneNumber = "+8801912345678"
            },
            new Nominee
            {
                NomineeId = Guid.NewGuid().ToString(),
                FullName = "Jane Doe",
                Relationship = "Child",
                SharePercentage = 50,
                DateOfBirth = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(new DateTime(2015, 5, 15).ToUniversalTime()),
                NidNumber = "",
                PhoneNumber = ""
            }
        };

        var command = new CreatePolicyCommand(
            productId,
            customerId,
            null, null, null, 1000, 50000, 12, DateTime.UtcNow, 
            "Proposer: Md. Zubayed Ur Rahman", nominees);

        _gatewayMock.Setup(g => g.CreatePolicyAsync(It.IsAny<CreatePolicyRequest>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new CreatePolicyResponse
            {
                PolicyId = Guid.NewGuid().ToString(),
                PolicyNumber = "POL-2026-0002",
                Message = "Success"
            });

        var response = await _handler.Handle(command, CancellationToken.None);

        response.Error.Should().BeNull();
        response.PolicyNumber.Should().NotBeEmpty();
    }
}

public class PolicyAggregateTests
{
    [Fact]
    public void Create_ShouldGenerateValidPolicyNumber()
    {
        var productCode = "OMC";
        var sequenceNumber = 1L;
        
        var policy = PolicyAggregate.Create(
            Guid.NewGuid(),
            productCode,
            "OVERSEAS_MEDICLAIM",
            Guid.NewGuid(),
            1239m,
            50000m,
            1,
            DateTime.UtcNow,
            sequenceNumber);

        policy.PolicyNumber.Should().NotBeEmpty();
        policy.PolicyNumber.Should().Contain(productCode);
    }

    [Fact]
    public void Create_ShouldSetCorrectStatus()
    {
        var policy = PolicyAggregate.Create(
            Guid.NewGuid(),
            "TEST",
            "OVERSEAS_MEDICLAIM",
            Guid.NewGuid(),
            1000m,
            50000m,
            12,
            DateTime.UtcNow,
            1);

        policy.Status.Should().Be("DRAFT");
    }

    [Fact]
    public void Activate_ShouldChangeStatusToActive()
    {
        var policy = PolicyAggregate.Create(
            Guid.NewGuid(),
            "TEST",
            "OVERSEAS_MEDICLAIM",
            Guid.NewGuid(),
            1000m,
            50000m,
            12,
            DateTime.UtcNow,
            1);

        policy.Activate();

        policy.Status.Should().Be("ACTIVE");
    }

    [Fact]
    public void AddNominees_ShouldAddMultipleNominees()
    {
        var policy = PolicyAggregate.Create(
            Guid.NewGuid(),
            "TEST",
            "OVERSEAS_MEDICLAIM",
            Guid.NewGuid(),
            1000m,
            50000m,
            12,
            DateTime.UtcNow,
            1);

        var nominees = new List<Nominee>
        {
            new Nominee { NomineeId = "1", FullName = "Nominee1", Relationship = "Spouse", SharePercentage = 50 },
            new Nominee { NomineeId = "2", FullName = "Nominee2", Relationship = "Child", SharePercentage = 50 }
        };

        policy.AddNominees(nominees);

        policy.Nominees.Should().HaveCount(2);
    }

    [Theory]
    [InlineData(14, 1, 1239, "Plan A - Non Schengen, Age 0-40")]
    [InlineData(15, 21, 1291, "Plan A - Non Schengen, Age 0-40")]
    [InlineData(29, 35, 1783, "Plan A - Non Schengen, Age 0-40")]
    public void Create_ShouldCalculateEndDateCorrectly(int days, int expectedMonths, decimal premium, string scenario)
    {
        var startDate = new DateTime(2026, 4, 1);
        var policy = PolicyAggregate.Create(
            Guid.NewGuid(),
            "OMC",
            "OVERSEAS_MEDICLAIM",
            Guid.NewGuid(),
            premium,
            50000m,
            expectedMonths,
            startDate,
            1);

        policy.StartDate.Should().Be(startDate);
        policy.EndDate.Should().Be(startDate.AddMonths(expectedMonths));
    }
}
