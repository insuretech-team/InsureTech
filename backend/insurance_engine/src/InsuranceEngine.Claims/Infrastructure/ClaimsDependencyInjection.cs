using InsuranceEngine.Grpc.Gateways;
using Microsoft.Extensions.DependencyInjection;
using MediatR;
using InsuranceEngine.Claims.Application;

namespace InsuranceEngine.Claims;

/// <summary>
/// Registers all Claims module services into the DI container.
/// Follows the professional pattern used in PoliSync.
/// </summary>
public static class ClaimsDependencyInjection
{
    public static IServiceCollection AddClaimsModule(this IServiceCollection services)
    {
        // Register MediatR for this assembly using the AssemblyMarker
        services.AddMediatR(cfg => 
        {
            cfg.RegisterServicesFromAssembly(typeof(AssemblyMarker).Assembly);
        });

        // Add module-specific gateways
        services.AddScoped<IClaimDataGateway, GoClaimDataGateway>();
        
        return services;
    }
}
