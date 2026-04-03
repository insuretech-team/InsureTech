using Microsoft.Extensions.DependencyInjection;
using MediatR;
using InsuranceEngine.Products.Application;
using InsuranceEngine.Products.Infrastructure;

namespace InsuranceEngine.Products;

/// <summary>
/// Registers all Products module services into the DI container.
/// Follows the professional pattern used in PoliSync.
/// </summary>
public static class ProductsDependencyInjection
{
    public static IServiceCollection AddProductsModule(this IServiceCollection services)
    {
        // Register MediatR for this assembly using the AssemblyMarker
        services.AddMediatR(cfg => 
        {
            cfg.RegisterServicesFromAssembly(typeof(AssemblyMarker).Assembly);
        });

        // Add module-specific gateways
        services.AddScoped<IProductDataGateway, GoProductDataGateway>();
        
        return services;
    }
}
