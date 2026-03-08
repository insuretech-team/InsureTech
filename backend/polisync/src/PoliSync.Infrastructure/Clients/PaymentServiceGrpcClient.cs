using Grpc.Net.Client;
using Insuretech.Payment.Services.V1;
using Microsoft.Extensions.Configuration;

namespace PoliSync.Infrastructure.Clients;

/// <summary>
/// Client for calling the Go Payment Service.
/// </summary>
public sealed class PaymentServiceGrpcClient : IDisposable
{
    private readonly GrpcChannel _channel;
    private readonly PaymentService.PaymentServiceClient _client;

    public PaymentServiceGrpcClient(IConfiguration configuration)
    {
        var paymentServiceUrl = configuration["GrpcClients:PaymentService"] ?? "http://localhost:50190";
        _channel = GrpcChannel.ForAddress(paymentServiceUrl);
        _client = new PaymentService.PaymentServiceClient(_channel);
    }

    public PaymentService.PaymentServiceClient Client => _client;

    public void Dispose()
    {
        _channel.Dispose();
    }
}
