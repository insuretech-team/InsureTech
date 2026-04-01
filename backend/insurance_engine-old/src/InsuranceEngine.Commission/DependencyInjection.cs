using InsuranceEngine.Commission.Domain.Interfaces;
using InsuranceEngine.Commission.Infrastructure.Persistence;
using InsuranceEngine.Commission.Infrastructure.Repositories;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;

namespace InsuranceEngine.Commission;

public static class DependencyInjection
{
    public static IServiceCollection AddCommissionModule(this IServiceCollection services, IConfiguration configuration)
    {
        services.AddDbContext<CommissionDbContext>(options =>
            options.UseNpgsql(configuration.GetConnectionString("DefaultConnection")).UseSnakeCaseNamingConvention());

        services.AddScoped<ICommissionRepository, CommissionRepository>();

        services.AddMediatR(cfg => cfg.RegisterServicesFromAssembly(typeof(DependencyInjection).Assembly));

        return services;
    }
}
