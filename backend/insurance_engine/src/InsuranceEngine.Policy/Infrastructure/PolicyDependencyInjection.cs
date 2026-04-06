using Microsoft.Extensions.DependencyInjection;
using MediatR;
using InsuranceEngine.Policy.Application;
using InsuranceEngine.Policy.Infrastructure;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Policy;

/// <summary>
/// Registers all Policy module services into the DI container.
/// Follows the professional pattern used in PoliSync.
/// </summary>
public static class PolicyDependencyInjection
{
    public static IServiceCollection AddPolicyModule(this IServiceCollection services, string connectionString)
    {
        services.AddMediatR(cfg => 
        {
            cfg.RegisterServicesFromAssembly(typeof(AssemblyMarker).Assembly);
        });

        services.AddDbContext<PolicyDbContext>(options =>
            options.UseNpgsql(connectionString));

        services.AddScoped<IPolicyDataGateway, SqlPolicyDataGateway>();
        
        return services;
    }
}
