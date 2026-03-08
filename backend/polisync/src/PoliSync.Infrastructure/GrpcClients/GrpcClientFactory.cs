using Grpc.Net.Client;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging;

namespace PoliSync.Infrastructure.GrpcClients;

/// <summary>
/// Factory for creating gRPC channels to upstream Go services.
/// Channels are cached (thread-safe) for connection reuse.
/// </summary>
public sealed class GrpcClientFactory : IDisposable
{
    private readonly IConfiguration _config;
    private readonly ILogger<GrpcClientFactory> _logger;
    private readonly Dictionary<string, GrpcChannel> _channels = [];
    private readonly object _lock = new();

    public GrpcClientFactory(IConfiguration config, ILogger<GrpcClientFactory> logger)
    {
        _config = config;
        _logger = logger;
    }

    public GrpcChannel GetChannel(string serviceName)
    {
        lock (_lock)
        {
            if (_channels.TryGetValue(serviceName, out var existing)) return existing;
            var address = _config[$"GrpcClients:{serviceName}"]
                ?? throw new InvalidOperationException($"gRPC address for '{serviceName}' not configured under GrpcClients:{serviceName}");
            _logger.LogInformation("Creating gRPC channel to {Service} at {Address}", serviceName, address);
            var channel = GrpcChannel.ForAddress(address, new GrpcChannelOptions
            {
                MaxReceiveMessageSize = 16 * 1024 * 1024, // 16MB
                MaxSendMessageSize = 16 * 1024 * 1024
            });
            _channels[serviceName] = channel;
            return channel;
        }
    }

    public T GetClient<T>(string serviceName, Func<GrpcChannel, T> factory)
        => factory(GetChannel(serviceName));

    public void Dispose()
    {
        foreach (var ch in _channels.Values) ch.Dispose();
        _channels.Clear();
    }
}
