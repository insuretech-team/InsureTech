using InsuranceEngine.Beneficiary.Application.Interfaces;
using InsuranceEngine.Beneficiary.Infrastructure.Persistence;
using InsuranceEngine.Beneficiary.Infrastructure.Repositories;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;

namespace InsuranceEngine.Beneficiary;

public static class DependencyInjection
{
    public static IServiceCollection AddBeneficiaryModule(this IServiceCollection services, IConfiguration configuration)
    {
        services.AddDbContext<BeneficiaryDbContext>(options =>
            options.UseNpgsql(configuration.GetConnectionString("DefaultConnection"),
                b => b.MigrationsAssembly(typeof(BeneficiaryDbContext).Assembly.FullName))
                .UseSnakeCaseNamingConvention());

        services.AddScoped<IBeneficiaryRepository, BeneficiaryRepository>();

        return services;
    }
}
