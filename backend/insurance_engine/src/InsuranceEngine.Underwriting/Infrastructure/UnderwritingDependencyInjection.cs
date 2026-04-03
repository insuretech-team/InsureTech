using InsuranceEngine.Grpc.Gateways;
using Microsoft.Extensions.DependencyInjection;
using MediatR;
using InsuranceEngine.Underwriting.Application;

namespace InsuranceEngine.Underwriting;

/// <summary>
/// Registers all Underwriting module services into the DI container.
/// Follows the professional pattern used in PoliSync.
/// </summary>
public static class UnderwritingDependencyInjection
{
    public static IServiceCollection AddUnderwritingModule(this IServiceCollection services)
    {
        // Register MediatR for this assembly using the AssemblyMarker
        services.AddMediatR(cfg => 
        {
            cfg.RegisterServicesFromAssembly(typeof(AssemblyMarker).Assembly);
        });

        // Add module-specific gateways
        services.AddScoped<IUnderwritingDataGateway, GoUnderwritingDataGateway>();
        
        return services;
    }
}
