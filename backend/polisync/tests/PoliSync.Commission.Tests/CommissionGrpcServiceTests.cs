using FluentAssertions;
using Insuretech.Commission.Services.V1;
using Insuretech.Partner.Entity.V1;
using MediatR;
using Microsoft.Extensions.Logging.Abstractions;
using Moq;
using PoliSync.Commission.GrpcServices;
using Xunit;

namespace PoliSync.Commission.Tests;

public class CommissionGrpcServiceTests
{
    private static CommissionGrpcService CreateService()
        => new(Mock.Of<IMediator>(), NullLogger<CommissionGrpcService>.Instance);

    [Fact]
    public async Task CalculateAndGetCommission_ReturnsPendingCommission()
    {
        var service = CreateService();

        var calculated = await service.CalculateCommission(new CalculateCommissionRequest
        {
            PolicyId = $"pol-{Guid.NewGuid():N}",
            CommissionType = "COMMISSION_TYPE_ACQUISITION",
            RecipientType = "agent",
            RecipientId = $"agent-{Guid.NewGuid():N}"
        }, null!);

        var get = await service.GetCommission(new GetCommissionRequest
        {
            CommissionId = calculated.CommissionId
        }, null!);

        calculated.CommissionId.Should().NotBeNullOrWhiteSpace();
        get.Commission.Status.Should().Be(CommissionStatus.Pending);
        get.Commission.CommissionAmount.Amount.Should().BeGreaterThan(0);
    }

    [Fact]
    public async Task CreateAndProcessPayout_MarksCommissionAsPaid()
    {
        var service = CreateService();
        var recipientId = $"agent-{Guid.NewGuid():N}";

        var commission = await service.CalculateCommission(new CalculateCommissionRequest
        {
            PolicyId = $"pol-{Guid.NewGuid():N}",
            CommissionType = "COMMISSION_TYPE_RENEWAL",
            RecipientType = "agent",
            RecipientId = recipientId
        }, null!);

        var payout = await service.CreatePayout(new CreatePayoutRequest
        {
            RecipientType = "agent",
            RecipientId = recipientId,
            PeriodStart = DateTime.UtcNow.AddDays(-30).ToString("yyyy-MM-dd"),
            PeriodEnd = DateTime.UtcNow.ToString("yyyy-MM-dd")
        }, null!);

        await service.ProcessPayout(new ProcessPayoutRequest
        {
            PayoutId = payout.PayoutId,
            PaymentMethod = "bank_transfer",
            PaymentReference = "payout-ref-1"
        }, null!);

        var get = await service.GetCommission(new GetCommissionRequest
        {
            CommissionId = commission.CommissionId
        }, null!);

        payout.PayoutId.Should().NotBeNullOrWhiteSpace();
        get.Commission.Status.Should().Be(CommissionStatus.Paid);
        get.Commission.PaymentId.Should().Be(payout.PayoutId);
    }
}
