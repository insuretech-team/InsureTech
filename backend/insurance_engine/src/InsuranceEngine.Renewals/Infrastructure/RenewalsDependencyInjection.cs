using Microsoft.Extensions.DependencyInjection;
using MediatR;
using InsuranceEngine.Renewals.Application;

namespace InsuranceEngine.Renewals;

public static class RenewalsDependencyInjection
{
    public static IServiceCollection AddRenewalsModule(this IServiceCollection services)
    {
        services.AddMediatR(cfg => 
        {
            cfg.RegisterServicesFromAssembly(typeof(AssemblyMarker).Assembly);
        });

        services.AddScoped<IRenewalDataGateway, Infrastructure.GoRenewalDataGateway>();

        return services;
    }
}
