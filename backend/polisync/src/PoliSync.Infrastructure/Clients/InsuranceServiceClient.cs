using Grpc.Core;
using Grpc.Net.Client;
using Insuretech.Actuarial.Services.V1;
using Insuretech.Crm.Services.V1;
using Insuretech.Insurance.Services.V1;
using Insuretech.Life.Services.V1;
using Insuretech.Quoting.Services.V1;
using Insuretech.Vehicle.Services.V1;
using Microsoft.AspNetCore.Http;
using Microsoft.Extensions.Configuration;

namespace PoliSync.Infrastructure.Clients;

/// <summary>
/// Client for calling the Go Insurance Service for database CRUD operations.
/// Includes X-* identity header forwarding for all gRPC calls (required by Go service auth).
/// </summary>
public class InsuranceServiceClient : IDisposable
{
    private readonly GrpcChannel _channel;
    private readonly InsuranceService.InsuranceServiceClient _client;
    private readonly QuotingService.QuotingServiceClient _quotingClient;
    private readonly LifeInsuranceService.LifeInsuranceServiceClient _lifeClient;
    private readonly VehicleService.VehicleServiceClient _vehicleClient;
    private readonly CrmService.CrmServiceClient _crmClient;
    private readonly ActuarialService.ActuarialServiceClient _actuarialClient;
    private readonly IHttpContextAccessor _httpContextAccessor;

    public InsuranceServiceClient(IConfiguration configuration, IHttpContextAccessor httpContextAccessor)
    {
        var insuranceServiceUrl = configuration["InsuranceService:Url"] ?? "http://localhost:50115";
        _channel = GrpcChannel.ForAddress(insuranceServiceUrl);
        _client = new InsuranceService.InsuranceServiceClient(_channel);
        _quotingClient = new QuotingService.QuotingServiceClient(_channel);
        _lifeClient = new LifeInsuranceService.LifeInsuranceServiceClient(_channel);
        _vehicleClient = new VehicleService.VehicleServiceClient(_channel);
        _crmClient = new CrmService.CrmServiceClient(_channel);
        _actuarialClient = new ActuarialService.ActuarialServiceClient(_channel);
        _httpContextAccessor = httpContextAccessor;
    }

    public InsuranceService.InsuranceServiceClient Client => _client;
    public QuotingService.QuotingServiceClient QuotingClient => _quotingClient;
    public LifeInsuranceService.LifeInsuranceServiceClient LifeClient => _lifeClient;
    public VehicleService.VehicleServiceClient VehicleClient => _vehicleClient;
    public CrmService.CrmServiceClient CrmClient => _crmClient;
    public ActuarialService.ActuarialServiceClient ActuarialClient => _actuarialClient;

    /// <summary>
    /// Builds gRPC CallOptions with X-* identity headers forwarded from the active HTTP request.
    /// The Go insurance service requires x-user-id, x-tenant-id etc. in gRPC metadata.
    /// Call this on every gRPC invocation instead of passing cancellationToken directly.
    /// </summary>
    public CallOptions BuildCallOptions(CancellationToken ct = default)
    {
        var headers = new Metadata();
        var httpCtx = _httpContextAccessor.HttpContext;
        if (httpCtx != null)
        {
            foreach (var key in new[] {
                "x-user-id", "x-tenant-id", "x-partner-id", "x-token-id",
                "x-user-type", "x-portal", "x-roles", "x-request-id", "x-session-id" })
            {
                var val = httpCtx.Request.Headers[key].FirstOrDefault();
                if (!string.IsNullOrEmpty(val)) headers.Add(key, val);
            }
        }
        return new CallOptions(headers: headers, cancellationToken: ct);
    }

    public void Dispose()
    {
        _channel?.Dispose();
    }
}
