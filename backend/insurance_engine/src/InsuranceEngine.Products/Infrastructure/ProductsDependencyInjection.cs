using Microsoft.Extensions.DependencyInjection;
using MediatR;
using Microsoft.EntityFrameworkCore;
using InsuranceEngine.Products.Application;
using InsuranceEngine.Products.Infrastructure;

namespace InsuranceEngine.Products;

public static class ProductsDependencyInjection
{
    public static IServiceCollection AddProductsModule(this IServiceCollection services, string connectionString)
    {
        services.AddMediatR(cfg => 
        {
            cfg.RegisterServicesFromAssembly(typeof(AssemblyMarker).Assembly);
        });

        services.AddDbContext<ProductsDbContext>(options =>
            options.UseNpgsql(connectionString));

        services.AddScoped<IProductDataGateway, SqlProductDataGateway>();
        
        return services;
    }
}
