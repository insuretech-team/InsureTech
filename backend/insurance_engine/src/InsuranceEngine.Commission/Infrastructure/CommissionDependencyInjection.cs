using Microsoft.Extensions.DependencyInjection;
using Microsoft.EntityFrameworkCore;
using InsuranceEngine.Commission.Infrastructure;

namespace InsuranceEngine.Commission;

public static class CommissionDependencyInjection
{
    public static IServiceCollection AddCommissionModule(this IServiceCollection services, string connectionString)
    {
        services.AddDbContext<CommissionDbContext>(options =>
            options.UseNpgsql(connectionString));

        services.AddScoped<ICommissionDataGateway, SqlCommissionDataGateway>();
        
        return services;
    }
}
