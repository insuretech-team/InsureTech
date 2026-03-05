using Microsoft.Extensions.DependencyInjection;
using PoliSync.Underwriting.Domain;
using PoliSync.Underwriting.Persistence;

namespace PoliSync.Underwriting;

public static class UnderwritingModule
{
    public static IServiceCollection AddUnderwritingModule(this IServiceCollection services)
    {
        // Persistence
        services.AddScoped<IQuoteRepository, QuoteRepository>();

        // MediatR handlers are automatically picked up by ApiHost's AddMediatR scanner
        
        return services;
    }
}
