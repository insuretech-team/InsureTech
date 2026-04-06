using Moq;
using Xunit;
using FluentAssertions;
using Microsoft.Extensions.Logging;
using InsuranceEngine.Policy.Application.Commands;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Infrastructure;
using Insuretech.Policy.Services.V1;
using Insuretech.Policy.Entity.V1;
using Google.Protobuf.WellKnownTypes;

namespace InsuranceEngine.Policy.Tests;

public class AddNomineeCommandHandlerTests
{
    private readonly Mock<IPolicyDataGateway> _gatewayMock;
    private readonly Mock<ILogger<AddNomineeCommandHandler>> _loggerMock;
    private readonly AddNomineeCommandHandler _handler;

    public AddNomineeCommandHandlerTests()
    {
        _gatewayMock = new Mock<IPolicyDataGateway>();
        _loggerMock = new Mock<ILogger<AddNomineeCommandHandler>>();
        _handler = new AddNomineeCommandHandler(_gatewayMock.Object, _loggerMock.Object);
    }

    [Fact]
    public async Task Handle_ShouldReturnNotFound_WhenPolicyNotFound()
    {
        var command = new AddNomineeCommand(
            Guid.NewGuid().ToString(),
            "John Doe",
            "Spouse",
            100,
            DateTime.UtcNow,
            "1234567890",
            "+8801912345678",
            null);

        _gatewayMock.Setup(g => g.GetPolicyAsync(It.IsAny<string>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new GetPolicyResponse { Policy = null });

        var result = await _handler.Handle(command, CancellationToken.None);

        result.IsFailure.Should().BeTrue();
        result.Error?.Code.Should().Be("POLICY_NOT_FOUND");
    }

    [Fact]
    public async Task Handle_ShouldAddNomineeSuccessfully()
    {
        var policyId = Guid.NewGuid().ToString();
        var command = new AddNomineeCommand(
            policyId,
            "John Doe",
            "Spouse",
            100,
            DateTime.UtcNow,
            "1234567890",
            "+8801912345678",
            null);

        var existingPolicy = CreatePolicyWithNominees(policyId, new List<Nominee>());

        _gatewayMock.Setup(g => g.GetPolicyAsync(It.IsAny<string>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new GetPolicyResponse { Policy = existingPolicy });

        _gatewayMock.Setup(g => g.UpdatePolicyAsync(It.IsAny<string>(), It.IsAny<List<Nominee>?>(), It.IsAny<string?>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new UpdatePolicyResponse { Error = null, Message = "Success" });

        var result = await _handler.Handle(command, CancellationToken.None);

        result.IsSuccess.Should().BeTrue();
        result.Value.Should().NotBeEmpty();
    }

    [Fact]
    public async Task Handle_ShouldAddNomineeWithFullDetails_FromExcelData()
    {
        var policyId = Guid.NewGuid().ToString();
        var command = new AddNomineeCommand(
            policyId,
            "Fatema Begum",
            "Wife",
            50,
            new DateTime(1985, 6, 15),
            "198515678900001",
            "+88017111234567",
            "15-06-1985");

        var existingPolicy = CreatePolicyWithNominees(policyId, new List<Nominee>());

        _gatewayMock.Setup(g => g.GetPolicyAsync(It.IsAny<string>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new GetPolicyResponse { Policy = existingPolicy });

        _gatewayMock.Setup(g => g.UpdatePolicyAsync(It.IsAny<string>(), It.IsAny<List<Nominee>?>(), It.IsAny<string?>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new UpdatePolicyResponse { Error = null, Message = "Success" });

        var result = await _handler.Handle(command, CancellationToken.None);

        result.IsSuccess.Should().BeTrue();
    }

    private static Insuretech.Policy.Entity.V1.Policy CreatePolicyWithNominees(string policyId, List<Nominee> nominees)
    {
        var policy = new Insuretech.Policy.Entity.V1.Policy
        {
            PolicyId = policyId,
            PolicyNumber = "POL-2026-0001",
            ProductId = "PROD-001",
            CustomerId = "CUST-001",
            Status = Insuretech.Policy.Entity.V1.PolicyStatus.Active,
            PremiumAmount = new Insuretech.Common.V1.Money { Amount = 1239, Currency = "BDT" },
            SumInsured = new Insuretech.Common.V1.Money { Amount = 50000, Currency = "BDT" },
            TenureMonths = 1,
            StartDate = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow),
            EndDate = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow.AddMonths(1))
        };
        foreach (var n in nominees)
        {
            policy.Nominees.Add(n);
        }
        return policy;
    }
}

public class UpdateNomineeCommandHandlerTests
{
    private readonly Mock<IPolicyDataGateway> _gatewayMock;
    private readonly Mock<ILogger<UpdateNomineeCommandHandler>> _loggerMock;
    private readonly UpdateNomineeCommandHandler _handler;

    public UpdateNomineeCommandHandlerTests()
    {
        _gatewayMock = new Mock<IPolicyDataGateway>();
        _loggerMock = new Mock<ILogger<UpdateNomineeCommandHandler>>();
        _handler = new UpdateNomineeCommandHandler(_gatewayMock.Object, _loggerMock.Object);
    }

    [Fact]
    public async Task Handle_ShouldReturnNotFound_WhenNomineeNotFound()
    {
        var policyId = Guid.NewGuid().ToString();
        var nomineeId = Guid.NewGuid().ToString();
        
        var command = new UpdateNomineeCommand(
            policyId,
            nomineeId,
            "Updated Name",
            "Spouse",
            100,
            null,
            null,
            null);

        var existingPolicy = CreatePolicyWithNominees(policyId, new List<Nominee>
        {
            new Nominee { NomineeId = "different-id", FullName = "Other", Relationship = "Child", SharePercentage = 100 }
        });

        _gatewayMock.Setup(g => g.GetPolicyAsync(It.IsAny<string>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new GetPolicyResponse { Policy = existingPolicy });

        var result = await _handler.Handle(command, CancellationToken.None);

        result.IsFailure.Should().BeTrue();
        result.Error?.Code.Should().Be("NOMINEE_NOT_FOUND");
    }

    [Fact]
    public async Task Handle_ShouldUpdateNomineeSuccessfully()
    {
        var policyId = Guid.NewGuid().ToString();
        var nomineeId = Guid.NewGuid().ToString();
        
        var command = new UpdateNomineeCommand(
            policyId,
            nomineeId,
            "Updated Name",
            "Spouse",
            50,
            null,
            null,
            null);

        var existingPolicy = CreatePolicyWithNominees(policyId, new List<Nominee>
        {
            new Nominee { NomineeId = nomineeId, FullName = "Original Name", Relationship = "Child", SharePercentage = 100 }
        });

        _gatewayMock.Setup(g => g.GetPolicyAsync(It.IsAny<string>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new GetPolicyResponse { Policy = existingPolicy });

        _gatewayMock.Setup(g => g.UpdatePolicyAsync(It.IsAny<string>(), It.IsAny<List<Nominee>?>(), It.IsAny<string?>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new UpdatePolicyResponse { Error = null, Message = "Success" });

        var result = await _handler.Handle(command, CancellationToken.None);

        result.IsSuccess.Should().BeTrue();
    }

    private static Insuretech.Policy.Entity.V1.Policy CreatePolicyWithNominees(string policyId, List<Nominee> nominees)
    {
        var policy = new Insuretech.Policy.Entity.V1.Policy
        {
            PolicyId = policyId,
            PolicyNumber = "POL-2026-0001"
        };
        foreach (var n in nominees)
        {
            policy.Nominees.Add(n);
        }
        return policy;
    }
}

public class DeleteNomineeCommandHandlerTests
{
    private readonly Mock<IPolicyDataGateway> _gatewayMock;
    private readonly Mock<ILogger<DeleteNomineeCommandHandler>> _loggerMock;
    private readonly DeleteNomineeCommandHandler _handler;

    public DeleteNomineeCommandHandlerTests()
    {
        _gatewayMock = new Mock<IPolicyDataGateway>();
        _loggerMock = new Mock<ILogger<DeleteNomineeCommandHandler>>();
        _handler = new DeleteNomineeCommandHandler(_gatewayMock.Object, _loggerMock.Object);
    }

    [Fact]
    public async Task Handle_ShouldDeleteNomineeSuccessfully()
    {
        var policyId = Guid.NewGuid().ToString();
        var nomineeId = Guid.NewGuid().ToString();
        
        var command = new DeleteNomineeCommand(policyId, nomineeId);

        var existingPolicy = CreatePolicyWithNominees(policyId, new List<Nominee>
        {
            new Nominee { NomineeId = nomineeId, FullName = "To Delete", Relationship = "Spouse", SharePercentage = 100 }
        });

        _gatewayMock.Setup(g => g.GetPolicyAsync(It.IsAny<string>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new GetPolicyResponse { Policy = existingPolicy });

        _gatewayMock.Setup(g => g.UpdatePolicyAsync(It.IsAny<string>(), It.IsAny<List<Nominee>?>(), It.IsAny<string?>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(new UpdatePolicyResponse { Error = null, Message = "Success" });

        var result = await _handler.Handle(command, CancellationToken.None);

        result.IsSuccess.Should().BeTrue();
    }

    private static Insuretech.Policy.Entity.V1.Policy CreatePolicyWithNominees(string policyId, List<Nominee> nominees)
    {
        var policy = new Insuretech.Policy.Entity.V1.Policy
        {
            PolicyId = policyId,
            PolicyNumber = "POL-2026-0001"
        };
        foreach (var n in nominees)
        {
            policy.Nominees.Add(n);
        }
        return policy;
    }
}
