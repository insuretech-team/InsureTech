using FluentAssertions;
using Insuretech.Commission.Services.V1;
using Insuretech.Common.V1;
using MediatR;
using Microsoft.Extensions.Logging.Abstractions;
using Moq;
using PoliSync.Commission.Application.Commands;
using PoliSync.Commission.GrpcServices;
using PoliSync.Commission.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using Xunit;
using Insuretech.Partner.Entity.V1;
using PartnerCommission = Insuretech.Partner.Entity.V1.Commission;

namespace PoliSync.Commission.Tests;

public class CommissionGrpcServiceTests
{
    private static CommissionGrpcService CreateService(IMediator? mediator = null, ICommissionDataGateway? gateway = null)
    {
        mediator ??= Mock.Of<IMediator>();
        gateway ??= Mock.Of<ICommissionDataGateway>();
        return new(mediator, NullLogger<CommissionGrpcService>.Instance, gateway);
    }

    [Fact]
    public async Task CalculateCommission_ReturnsCommissionId_OnSuccess()
    {
        var commissionId = $"COM-{Guid.NewGuid():N}"[..16];
        var mediatorMock = new Mock<IMediator>();
        mediatorMock
            .Setup(m => m.Send(It.IsAny<CalculateCommissionCommand>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(Result.Ok(new CalculateCommissionResult(
                commissionId,
                $"COM-{DateTime.UtcNow:yyyyMMdd}-{commissionId[..8].ToUpperInvariant()}",
                new Money { Amount = 150_000, Currency = "BDT" },
                "base=1000000; rate=15.00%; amount=150000")));

        var service = CreateService(mediatorMock.Object);

        var result = await service.CalculateCommission(new CalculateCommissionRequest
        {
            PolicyId = $"pol-{Guid.NewGuid():N}",
            CommissionType = "ACQUISITION",
            RecipientType = "agent",
            RecipientId = $"agent-{Guid.NewGuid():N}"
        }, null!);

        result.CommissionId.Should().Be(commissionId);
        result.Amount.Amount.Should().Be(150_000);
        result.Error.Should().BeNull();
    }

    [Fact]
    public async Task CalculateCommission_ReturnsError_WhenValidationFails()
    {
        var service = CreateService();

        var result = await service.CalculateCommission(new CalculateCommissionRequest
        {
            PolicyId = "", // Missing
            RecipientId = ""
        }, null!);

        result.Error.Should().NotBeNull();
        result.Error!.Code.Should().Be("VALIDATION_ERROR");
    }

    [Fact]
    public async Task GetCommission_ReturnsCommission_WhenFound()
    {
        var commissionId = Guid.NewGuid().ToString("N");
        var gatewayMock = new Mock<ICommissionDataGateway>();
        gatewayMock
            .Setup(g => g.GetCommissionAsync(commissionId, It.IsAny<CancellationToken>()))
            .ReturnsAsync(new PartnerCommission
            {
                CommissionId = commissionId,
                PolicyId = "pol-001",
                Status = CommissionStatus.Pending,
                CommissionAmount = new Money { Amount = 150_000, Currency = "BDT" }
            });

        var service = CreateService(gateway: gatewayMock.Object);

        var result = await service.GetCommission(
            new GetCommissionRequest { CommissionId = commissionId }, null!);

        result.Commission.Should().NotBeNull();
        result.Commission!.CommissionId.Should().Be(commissionId);
        result.Commission.Status.Should().Be(CommissionStatus.Pending);
    }

    [Fact]
    public async Task GetCommission_ReturnsNotFound_WhenMissing()
    {
        var gatewayMock = new Mock<ICommissionDataGateway>();
        gatewayMock
            .Setup(g => g.GetCommissionAsync(It.IsAny<string>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync((PartnerCommission?)null);

        var service = CreateService(gateway: gatewayMock.Object);

        var result = await service.GetCommission(
            new GetCommissionRequest { CommissionId = "missing" }, null!);

        result.Error.Should().NotBeNull();
        result.Error!.Code.Should().Be("NOT_FOUND");
    }

    [Fact]
    public async Task CreatePayout_ReturnsPayout_OnSuccess()
    {
        var payoutId = Guid.NewGuid().ToString("N");
        var mediatorMock = new Mock<IMediator>();
        mediatorMock
            .Setup(m => m.Send(It.IsAny<CreateCommissionPayoutCommand>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(Result.Ok(new CreateCommissionPayoutResult(
                payoutId,
                $"PYO-{DateTime.UtcNow:yyyyMMdd}-{payoutId[..8].ToUpperInvariant()}",
                new Money { Amount = 300_000, Currency = "BDT" },
                2)));

        var service = CreateService(mediatorMock.Object);

        var result = await service.CreatePayout(new CreatePayoutRequest
        {
            RecipientType = "agent",
            RecipientId = $"agent-{Guid.NewGuid():N}",
            PeriodStart = DateTime.UtcNow.AddDays(-30).ToString("yyyy-MM-dd"),
            PeriodEnd = DateTime.UtcNow.ToString("yyyy-MM-dd")
        }, null!);

        result.PayoutId.Should().Be(payoutId);
        result.CommissionCount.Should().Be(2);
        result.TotalAmount.Amount.Should().Be(300_000);
    }

    [Fact]
    public async Task ProcessPayout_ReturnsPaidAt_OnSuccess()
    {
        var mediatorMock = new Mock<IMediator>();
        mediatorMock
            .Setup(m => m.Send(It.IsAny<ProcessCommissionPayoutCommand>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(Result.Ok(DateTime.UtcNow.ToString("O")));

        var service = CreateService(mediatorMock.Object);

        var result = await service.ProcessPayout(new ProcessPayoutRequest
        {
            PayoutId = Guid.NewGuid().ToString("N"),
            PaymentMethod = "bank_transfer",
            PaymentReference = "payout-ref-1"
        }, null!);

        result.Message.Should().Be("Payout processed");
        result.PaidAt.Should().NotBeNullOrWhiteSpace();
        result.Error.Should().BeNull();
    }
}
