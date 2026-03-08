using Insuretech.Common.V1;
using Insuretech.Payment.Services.V1;

namespace PoliSync.Refund.Infrastructure;

public interface IRefundPaymentGateway
{
    Task<InitiateRefundResponse> InitiateRefundAsync(
        string paymentId,
        Money refundAmount,
        string reason,
        string initiatedBy,
        CancellationToken cancellationToken = default);
}
