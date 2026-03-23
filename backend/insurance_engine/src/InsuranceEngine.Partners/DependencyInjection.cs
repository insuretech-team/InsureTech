using InsuranceEngine.Partners.Domain.Interfaces;
using InsuranceEngine.Partners.Infrastructure.Persistence;
using InsuranceEngine.Partners.Infrastructure.Repositories;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;

namespace InsuranceEngine.Partners;

public static class DependencyInjection
{
    public static IServiceCollection AddPartnersModule(this IServiceCollection services, IConfiguration configuration)
    {
        services.AddDbContext<PartnerDbContext>(options =>
            options.UseNpgsql(configuration.GetConnectionString("DefaultConnection"))
                .UseSnakeCaseNamingConvention());

        services.AddScoped<IPartnerRepository, PartnerRepository>();

        services.AddMediatR(cfg => cfg.RegisterServicesFromAssembly(typeof(DependencyInjection).Assembly));

        return services;
    }
}
