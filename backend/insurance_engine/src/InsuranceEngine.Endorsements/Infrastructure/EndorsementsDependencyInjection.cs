using Microsoft.Extensions.DependencyInjection;
using MediatR;
using InsuranceEngine.Endorsements.Application;

namespace InsuranceEngine.Endorsements;

/// <summary>
/// Registers all Endorsements module services into the DI container.
/// Follows the professional pattern used in PoliSync.
/// </summary>
public static class EndorsementsDependencyInjection
{
    public static IServiceCollection AddEndorsementsModule(this IServiceCollection services)
    {
        // Register MediatR for this assembly using the AssemblyMarker
        services.AddMediatR(cfg => 
        {
            cfg.RegisterServicesFromAssembly(typeof(AssemblyMarker).Assembly);
        });
        // Add module-specific gateways
        services.AddScoped<Infrastructure.IEndorsementDataGateway, Infrastructure.GoEndorsementDataGateway>();

        return services;
    }
}
