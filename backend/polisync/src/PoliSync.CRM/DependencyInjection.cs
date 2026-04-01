using Microsoft.Extensions.DependencyInjection;
using PoliSync.CRM.Services;

namespace PoliSync.CRM;

public static class DependencyInjection
{
    public static IServiceCollection AddCrmServices(this IServiceCollection services)
    {
        services.AddSingleton<ICrmService, GoCrmService>();
        
        return services;
    }
}
