using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using PoliSync.Products.Infrastructure;

namespace PoliSync.Products;

/// <summary>
/// Registers all Products module services into the DI container.
/// Called from ApiHost Program.cs.
/// </summary>
public static class DependencyInjection
{
    public static IServiceCollection AddProductsModule(this IServiceCollection services)
    {
        services.AddScoped<IProductRepository, ProductRepository>();
        return services;
    }
}
