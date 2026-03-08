using Grpc.Net.Client;
using Insuretech.Orders.Services.V1;
using Microsoft.Extensions.Configuration;

namespace PoliSync.Infrastructure.Clients;

/// <summary>
/// Client for calling the Go Orders Service.
/// </summary>
public sealed class OrderServiceGrpcClient : IDisposable
{
    private readonly GrpcChannel _channel;
    private readonly OrderService.OrderServiceClient _client;

    public OrderServiceGrpcClient(IConfiguration configuration)
    {
        var orderServiceUrl = configuration["GrpcClients:OrdersService"] ?? "http://localhost:50142";
        _channel = GrpcChannel.ForAddress(orderServiceUrl);
        _client = new OrderService.OrderServiceClient(_channel);
    }

    public OrderService.OrderServiceClient Client => _client;

    public void Dispose()
    {
        _channel.Dispose();
    }
}
