using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Configuration;
using Microsoft.EntityFrameworkCore;
using InsuranceEngine.Underwriting.Infrastructure.Persistence;
using InsuranceEngine.Underwriting.Domain.Services;
using InsuranceEngine.Underwriting.Application.Interfaces;
using InsuranceEngine.Underwriting.Infrastructure.Repositories;

namespace InsuranceEngine.Underwriting;

public static class UnderwritingModule
{
    public static IServiceCollection AddUnderwritingModule(this IServiceCollection services, IConfiguration configuration)
    {
        services.AddDbContext<UnderwritingDbContext>(options =>
            options.UseNpgsql(configuration.GetConnectionString("DefaultConnection"),
                b => b.MigrationsHistoryTable("__EFMigrationsHistory", "insurance_schema"))
                .UseSnakeCaseNamingConvention());

        services.AddScoped<IUnderwritingRepository, UnderwritingRepository>();
        services.AddSingleton<QuoteNumberGenerator>();


        return services;
    }
}
