using Microsoft.Extensions.DependencyInjection;
using PoliSync.LifeInsurance.Services;

namespace PoliSync.LifeInsurance;

public static class DependencyInjection
{
    public static IServiceCollection AddLifeInsuranceServices(this IServiceCollection services)
    {
        services.AddSingleton<ILifeProductService, GoLifeProductService>();
        services.AddSingleton<ILifeQuoteService, GoLifeQuoteService>();
        
        return services;
    }
}
