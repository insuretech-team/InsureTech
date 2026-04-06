using Microsoft.Extensions.DependencyInjection;
using MediatR;
using Microsoft.EntityFrameworkCore;
using InsuranceEngine.Underwriting.Application;
using InsuranceEngine.Underwriting.Infrastructure;

namespace InsuranceEngine.Underwriting;

public static class UnderwritingDependencyInjection
{
    public static IServiceCollection AddUnderwritingModule(this IServiceCollection services, string connectionString)
    {
        services.AddMediatR(cfg => 
        {
            cfg.RegisterServicesFromAssembly(typeof(AssemblyMarker).Assembly);
        });

        services.AddDbContext<UnderwritingDbContext>(options =>
            options.UseNpgsql(connectionString));

        services.AddScoped<IUnderwritingDataGateway, SqlUnderwritingDataGateway>();
        
        return services;
    }
}
