using Microsoft.Extensions.DependencyInjection;
using MediatR;
using InsuranceEngine.Renewals.Application;

namespace InsuranceEngine.Renewals;

/// <summary>
/// Registers all Renewals module services into the DI container.
/// Follows the professional pattern used in PoliSync.
/// </summary>
public static class RenewalsDependencyInjection
{
    public static IServiceCollection AddRenewalsModule(this IServiceCollection services)
    {
        // Register MediatR for this assembly using the AssemblyMarker
        services.AddMediatR(cfg => 
        {
            cfg.RegisterServicesFromAssembly(typeof(AssemblyMarker).Assembly);
        });
        // Add module-specific gateways
        services.AddScoped<Infrastructure.IRenewalDataGateway, Infrastructure.GoRenewalDataGateway>();

        return services;
    }
}
