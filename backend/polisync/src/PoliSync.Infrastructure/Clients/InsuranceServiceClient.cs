using Grpc.Net.Client;
using Insuretech.Insurance.Services.V1;
using Microsoft.Extensions.Configuration;

namespace PoliSync.Infrastructure.Clients;

/// <summary>
/// Client for calling the Go Insurance Service for database CRUD operations
/// </summary>
public class InsuranceServiceClient : IDisposable
{
    private readonly GrpcChannel _channel;
    private readonly InsuranceService.InsuranceServiceClient _client;

    public InsuranceServiceClient(IConfiguration configuration)
    {
        var insuranceServiceUrl = configuration["InsuranceService:Url"] ?? "http://localhost:50115";
        
        _channel = GrpcChannel.ForAddress(insuranceServiceUrl);
        _client = new InsuranceService.InsuranceServiceClient(_channel);
    }

    public InsuranceService.InsuranceServiceClient Client => _client;

    public void Dispose()
    {
        _channel?.Dispose();
    }
}
