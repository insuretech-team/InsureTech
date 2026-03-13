using System;
using InsuranceEngine.Fraud.Infrastructure.Persistence;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;

namespace InsuranceEngine.Fraud;

public static class DependencyInjection
{
    public static IServiceCollection AddFraudModule(this IServiceCollection services, IConfiguration configuration)
    {
        var connectionString = configuration.GetConnectionString("DefaultConnection");

        services.AddDbContext<FraudDbContext>(options =>
            options.UseNpgsql(connectionString,
                b => b.MigrationsAssembly(typeof(FraudDbContext).Assembly.FullName)));


        return services;
    }
}
