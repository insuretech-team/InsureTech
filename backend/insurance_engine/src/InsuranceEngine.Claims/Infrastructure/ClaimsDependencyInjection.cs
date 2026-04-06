using Microsoft.Extensions.DependencyInjection;
using MediatR;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Configuration;
using InsuranceEngine.Claims.Application;
using InsuranceEngine.Claims.Infrastructure;

namespace InsuranceEngine.Claims;

public static class ClaimsDependencyInjection
{
    public static IServiceCollection AddClaimsModule(this IServiceCollection services, IConfiguration configuration, string connectionString)
    {
        services.AddMediatR(cfg => 
        {
            cfg.RegisterServicesFromAssembly(typeof(AssemblyMarker).Assembly);
        });

        services.AddDbContext<ClaimsDbContext>(options =>
            options.UseNpgsql(connectionString));

        var useGoProxy = configuration.GetValue<bool>("Claims:UseGoProxy");
        if (useGoProxy)
        {
            services.AddScoped<IClaimDataGateway, GoClaimsDataGateway>();
        }
        else
        {
            services.AddScoped<IClaimDataGateway, SqlClaimsDataGateway>();
        }
        
        return services;
    }
}
