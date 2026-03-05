using Microsoft.Extensions.DependencyInjection;
using PoliSync.Beneficiaries.Domain;
using PoliSync.Beneficiaries.Persistence;

namespace PoliSync.Beneficiaries;

public static class BeneficiaryModule
{
    public static IServiceCollection AddBeneficiaryModule(this IServiceCollection services)
    {
        // Persistence
        services.AddScoped<IBeneficiaryRepository, BeneficiaryRepository>();

        // MediatR is already registered in ApiHost, it will pick up handlers 
        // if this assembly is scanned.

        return services;
    }
}
