using FluentAssertions;
using Insuretech.Refund.Entity.V1;
using Insuretech.Refund.Services.V1;
using Insuretech.Common.V1;
using Insuretech.Payment.Services.V1;
using MediatR;
using Microsoft.Extensions.Logging.Abstractions;
using Moq;
using PoliSync.Refund.GrpcServices;
using PoliSync.Refund.Infrastructure;
using Xunit;

namespace PoliSync.Refund.Tests;

public class RefundGrpcServiceTests
{
    private sealed class FakeRefundPaymentGateway : IRefundPaymentGateway
    {
        public Task<InitiateRefundResponse> InitiateRefundAsync(
            string paymentId,
            Money refundAmount,
            string reason,
            string initiatedBy,
            CancellationToken cancellationToken = default)
        {
            return Task.FromResult(new InitiateRefundResponse
            {
                RefundId = $"PAYREF-{Guid.NewGuid():N}"[..18],
                Status = "initiated"
            });
        }
    }

    private static RefundGrpcService CreateService()
        => new(Mock.Of<IMediator>(), NullLogger<RefundGrpcService>.Instance, new FakeRefundPaymentGateway());

    [Fact]
    public async Task RefundLifecycle_RequestCalculateApproveProcess_CompletesRefund()
    {
        var service = CreateService();
        var policyId = $"pol-{Guid.NewGuid():N}";

        var requested = await service.RequestRefund(new RequestRefundRequest
        {
            PolicyId = policyId,
            Reason = "REFUND_REASON_CUSTOMER_REQUEST",
            ReasonDetails = "Customer requested cancellation"
        }, null!);

        await service.CalculateRefund(new CalculateRefundRequest
        {
            PolicyId = policyId,
            Reason = "REFUND_REASON_CUSTOMER_REQUEST"
        }, null!);

        await service.ApproveRefund(new ApproveRefundRequest
        {
            RefundId = requested.RefundId,
            ApprovedBy = "refund-ops",
            Comments = "Approved"
        }, null!);

        await service.ProcessRefund(new ProcessRefundRequest
        {
            RefundId = requested.RefundId,
            PaymentMethod = "bank_transfer",
            PaymentReference = "refund-pay-1"
        }, null!);

        var get = await service.GetRefund(new GetRefundRequest
        {
            RefundId = requested.RefundId
        }, null!);

        get.Refund.Status.Should().Be(RefundStatus.Completed);
        get.Refund.PaymentMethod.Should().Be("bank_transfer");
        get.Refund.RefundableAmount.Amount.Should().BeGreaterThanOrEqualTo(0);
    }

    [Fact]
    public async Task ListRefunds_ByCompletedStatus_IncludesProcessedRefund()
    {
        var service = CreateService();
        var policyId = $"pol-{Guid.NewGuid():N}";

        var requested = await service.RequestRefund(new RequestRefundRequest
        {
            PolicyId = policyId,
            Reason = "REFUND_REASON_FREE_LOOK_CANCELLATION",
            ReasonDetails = "Within free look"
        }, null!);

        await service.CalculateRefund(new CalculateRefundRequest
        {
            PolicyId = policyId,
            Reason = "REFUND_REASON_FREE_LOOK_CANCELLATION"
        }, null!);

        await service.ApproveRefund(new ApproveRefundRequest
        {
            RefundId = requested.RefundId,
            ApprovedBy = "refund-ops",
            Comments = "Approved"
        }, null!);

        await service.ProcessRefund(new ProcessRefundRequest
        {
            RefundId = requested.RefundId,
            PaymentMethod = "wallet",
            PaymentReference = "refund-pay-2"
        }, null!);

        var list = await service.ListRefunds(new ListRefundsRequest
        {
            Status = "REFUND_STATUS_COMPLETED",
            Page = 1,
            PageSize = 20
        }, null!);

        list.Refunds.Should().Contain(x => x.Id == requested.RefundId);
    }
}
