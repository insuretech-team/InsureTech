using Insuretech.Payment.Services.V1;

namespace PoliSync.Infrastructure.GrpcClients;

/// <summary>
/// Typed wrapper for the Go Payment service gRPC client.
/// </summary>
public sealed class PaymentGrpcClient
{
    private readonly GrpcClientFactory _factory;

    public PaymentGrpcClient(GrpcClientFactory factory) => _factory = factory;

    private PaymentService.PaymentServiceClient Client =>
        _factory.GetClient("PaymentService", ch => new PaymentService.PaymentServiceClient(ch));

    public async Task<InitiatePaymentResponse> InitiateAsync(
        string orderId, long amountPaisa, string currency,
        string customerId, string method,
        CancellationToken ct = default)
    {
        return await Client.InitiatePaymentAsync(new InitiatePaymentRequest
        {
            OrderId    = orderId,
            Amount     = amountPaisa,
            Currency   = currency,
            CustomerId = customerId,
            Method     = method,
        }, cancellationToken: ct);
    }

    public async Task<VerifyPaymentResponse> VerifyAsync(
        string paymentId, CancellationToken ct = default)
    {
        return await Client.VerifyPaymentAsync(
            new VerifyPaymentRequest { PaymentId = paymentId },
            cancellationToken: ct);
    }
}
