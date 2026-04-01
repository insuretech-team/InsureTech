using Insuretech.Payment.Services.V1;
using PaymentService = Insuretech.Payment.Services.V1.PaymentService;
using Money = Insuretech.Common.V1.Money;

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
            UserId = customerId,
            CustomerId = customerId,
            OrderId = orderId,
            Amount = new Money
            {
                Amount = amountPaisa,
                Currency = currency,
                DecimalAmount = amountPaisa / 100d
            },
            Currency = currency,
            PaymentMethod = method,
            IdempotencyKey = $"polisync:{orderId}:{method}",
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
