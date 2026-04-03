using Microsoft.Extensions.DependencyInjection;
using MediatR;
using InsuranceEngine.Cancellations.Application;

namespace InsuranceEngine.Cancellations;

/// <summary>
/// Registers all Cancellations module services into the DI container.
/// Follows the professional pattern used in PoliSync.
/// </summary>
public static class CancellationsDependencyInjection
{
    public static IServiceCollection AddCancellationsModule(this IServiceCollection services)
    {
        // Register MediatR for this assembly using the AssemblyMarker
        services.AddMediatR(cfg => 
        {
            cfg.RegisterServicesFromAssembly(typeof(AssemblyMarker).Assembly);
        });
        // Add module-specific gateways
        services.AddScoped<Infrastructure.ICancellationDataGateway, Infrastructure.GoCancellationDataGateway>();

        return services;
    }
}
