using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Configuration;
using Microsoft.EntityFrameworkCore;
using InsuranceEngine.Claims.Infrastructure.Persistence;
using InsuranceEngine.Claims.Application.Interfaces;
using InsuranceEngine.Claims.Infrastructure.Repositories;
using InsuranceEngine.Claims.Domain.Services;

namespace InsuranceEngine.Claims;

public static class ClaimsModule
{
    public static IServiceCollection AddClaimsModule(this IServiceCollection services, IConfiguration configuration)
    {
        services.AddDbContext<ClaimsDbContext>(options =>
            options.UseNpgsql(configuration.GetConnectionString("DefaultConnection"),
                b => b.MigrationsHistoryTable("__EFMigrationsHistory", "insurance_schema"))
                .UseSnakeCaseNamingConvention());

        services.AddScoped<IClaimsRepository, ClaimsRepository>();
        services.AddScoped<ClaimEligibilityValidator>();
        services.AddSingleton<ClaimSettlementCalculator>();

        return services;
    }
}
