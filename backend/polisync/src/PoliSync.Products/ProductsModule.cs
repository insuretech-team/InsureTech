using Microsoft.Extensions.DependencyInjection;
using PoliSync.Products.Domain;
using PoliSync.Products.Persistence;

namespace PoliSync.Products;

/// <summary>
/// Products module DI registration.
/// </summary>
public static class ProductsModule
{
    public static IServiceCollection AddProductsModule(this IServiceCollection services)
    {
        // Repository
        services.AddScoped<IProductRepository, ProductRepository>();

        return services;
    }
}
