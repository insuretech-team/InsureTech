using Microsoft.Extensions.DependencyInjection;
using MediatR;
using InsuranceEngine.Claims.Application;
using InsuranceEngine.Claims.Infrastructure;

namespace InsuranceEngine.Claims;

/// <summary>
/// Registers all Claims module services into the DI container.
/// Follows the professional pattern used in PoliSync.
/// </summary>
public static class ClaimsDependencyInjection
{
    public static IServiceCollection AddClaimsModule(this IServiceCollection services)
    {
        services.AddMediatR(cfg => 
        {
            cfg.RegisterServicesFromAssembly(typeof(AssemblyMarker).Assembly);
        });

        services.AddScoped<IClaimDataGateway, GoClaimsDataGateway>();
        
        return services;
    }
}
