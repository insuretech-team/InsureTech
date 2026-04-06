using Microsoft.Extensions.DependencyInjection;
using MediatR;
using Microsoft.EntityFrameworkCore;
using InsuranceEngine.Endorsements.Application;

namespace InsuranceEngine.Endorsements;

public static class EndorsementsDependencyInjection
{
    public static IServiceCollection AddEndorsementsModule(this IServiceCollection services, string connectionString)
    {
        services.AddMediatR(cfg => 
        {
            cfg.RegisterServicesFromAssembly(typeof(AssemblyMarker).Assembly);
        });

        services.AddDbContext<Infrastructure.EndorsementsDbContext>(options =>
            options.UseNpgsql(connectionString));

        services.AddScoped<Infrastructure.ISqlEndorsementDataGateway, Infrastructure.SqlEndorsementDataGateway>();
        services.AddScoped<IEndorsementProcessingService, EndorsementProcessingService>();

        return services;
    }
}
