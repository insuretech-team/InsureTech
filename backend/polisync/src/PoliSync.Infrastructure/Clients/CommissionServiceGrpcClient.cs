using Grpc.Net.Client;
using Insuretech.Commission.Services.V1;
using Microsoft.Extensions.Configuration;

namespace PoliSync.Infrastructure.Clients;

/// <summary>
/// Client for calling the Go Commission Service gRPC API.
/// </summary>
public sealed class CommissionServiceGrpcClient : IDisposable
{
    private readonly GrpcChannel _channel;
    private readonly CommissionService.CommissionServiceClient _client;

    public CommissionServiceGrpcClient(IConfiguration configuration)
    {
        var url = configuration["GrpcClients:CommissionService"] ?? "http://localhost:50160";
        _channel = GrpcChannel.ForAddress(url);
        _client = new CommissionService.CommissionServiceClient(_channel);
    }

    public CommissionService.CommissionServiceClient Client => _client;

    public void Dispose() => _channel.Dispose();
}
