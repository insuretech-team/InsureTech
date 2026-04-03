using Grpc.Net.Client;
using Microsoft.Extensions.Configuration;
using System.Collections.Concurrent;

namespace InsuranceEngine.Grpc.Clients;

/// <summary>
/// A factory for creating gRPC channels and service clients.
/// Caches channels to avoid overhead and resource leaks.
/// </summary>
public sealed class GrpcClientFactory : IDisposable
{
    private readonly string _address;
    private readonly ConcurrentDictionary<string, GrpcChannel> _channels = new();

    public GrpcClientFactory(IConfiguration configuration)
    {
        _address = configuration["InsuranceService:Url"] ?? "http://localhost:50115";
    }

    public GrpcChannel GetChannel()
    {
        return _channels.GetOrAdd("default", _ => GrpcChannel.ForAddress(_address));
    }

    public void Dispose()
    {
        foreach (var channel in _channels.Values)
        {
            channel.Dispose();
        }
    }
}
