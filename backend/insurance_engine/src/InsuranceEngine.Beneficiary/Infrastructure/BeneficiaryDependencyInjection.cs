using Microsoft.Extensions.DependencyInjection;
using MediatR;
using InsuranceEngine.Beneficiary.Application;
using InsuranceEngine.Beneficiary.Infrastructure;

namespace InsuranceEngine.Beneficiary;

/// <summary>
/// Registers all Beneficiary module services into the DI container.
/// Follows the professional pattern used in PoliSync.
/// </summary>
public static class BeneficiaryDependencyInjection
{
    public static IServiceCollection AddBeneficiaryModule(this IServiceCollection services)
    {
        // Register MediatR for this assembly using the AssemblyMarker
        services.AddMediatR(cfg => 
        {
            cfg.RegisterServicesFromAssembly(typeof(AssemblyMarker).Assembly);
        });

        // Add module-specific gateways
        services.AddScoped<IBeneficiaryDataGateway, GoBeneficiaryDataGateway>();
        
        return services;
    }
}
