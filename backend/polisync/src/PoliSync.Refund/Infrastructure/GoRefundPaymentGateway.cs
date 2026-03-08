using Insuretech.Common.V1;
using Insuretech.Payment.Services.V1;
using PoliSync.Infrastructure.Clients;

namespace PoliSync.Refund.Infrastructure;

public sealed class GoRefundPaymentGateway : IRefundPaymentGateway
{
    private readonly PaymentServiceGrpcClient _paymentClient;

    public GoRefundPaymentGateway(PaymentServiceGrpcClient paymentClient)
    {
        _paymentClient = paymentClient;
    }

    public async Task<InitiateRefundResponse> InitiateRefundAsync(
        string paymentId,
        Money refundAmount,
        string reason,
        string initiatedBy,
        CancellationToken cancellationToken = default)
    {
        return await _paymentClient.Client.InitiateRefundAsync(new InitiateRefundRequest
        {
            PaymentId = paymentId,
            RefundAmount = refundAmount,
            Reason = reason,
            InitiatedBy = initiatedBy
        }, cancellationToken: cancellationToken);
    }
}
