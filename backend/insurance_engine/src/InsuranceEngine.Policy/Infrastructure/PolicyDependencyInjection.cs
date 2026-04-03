using Microsoft.Extensions.DependencyInjection;
using MediatR;
using InsuranceEngine.Policy.Application;
using InsuranceEngine.Policy.Infrastructure;

namespace InsuranceEngine.Policy;

/// <summary>
/// Registers all Policy module services into the DI container.
/// Follows the professional pattern used in PoliSync.
/// </summary>
public static class PolicyDependencyInjection
{
    public static IServiceCollection AddPolicyModule(this IServiceCollection services)
    {
        // Register MediatR for this assembly using the AssemblyMarker
        services.AddMediatR(cfg => 
        {
            cfg.RegisterServicesFromAssembly(typeof(AssemblyMarker).Assembly);
        });

        // Add module-specific gateways
        services.AddScoped<IPolicyDataGateway, GoPolicyDataGateway>();
        
        return services;
    }
}
