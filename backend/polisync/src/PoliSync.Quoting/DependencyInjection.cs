using Microsoft.Extensions.DependencyInjection;
using PoliSync.Quoting.Services;

namespace PoliSync.Quoting;

public static class DependencyInjection
{
    public static IServiceCollection AddQuotingServices(this IServiceCollection services)
    {
        services.AddSingleton<IQuoteService, GoQuoteService>();
        
        return services;
    }
}
