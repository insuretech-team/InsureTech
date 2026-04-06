using Microsoft.Extensions.DependencyInjection;
using MediatR;
using InsuranceEngine.Renewals.Application;
using InsuranceEngine.Renewals.Infrastructure;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Renewals;

public static class RenewalsDependencyInjection
{
    public static IServiceCollection AddRenewalsModule(this IServiceCollection services, string connectionString)
    {
        services.AddMediatR(cfg => 
        {
            cfg.RegisterServicesFromAssembly(typeof(AssemblyMarker).Assembly);
        });

        services.AddDbContext<RenewalsDbContext>(options =>
            options.UseNpgsql(connectionString));

        services.AddScoped<IRenewalDataGateway, SqlRenewalDataGateway>();

        return services;
    }
}
