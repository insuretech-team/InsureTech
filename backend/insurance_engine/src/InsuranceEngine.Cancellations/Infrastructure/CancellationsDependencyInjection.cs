using Microsoft.Extensions.DependencyInjection;
using MediatR;
using Microsoft.EntityFrameworkCore;
using InsuranceEngine.Cancellations.Application;
using InsuranceEngine.Cancellations.Infrastructure;

namespace InsuranceEngine.Cancellations;

public static class CancellationsDependencyInjection
{
    public static IServiceCollection AddCancellationsModule(this IServiceCollection services, string connectionString)
    {
        services.AddMediatR(cfg => 
        {
            cfg.RegisterServicesFromAssembly(typeof(AssemblyMarker).Assembly);
        });

        services.AddDbContext<CancellationsDbContext>(options =>
            options.UseNpgsql(connectionString));

        services.AddScoped<ICancellationDataGateway, SqlCancellationDataGateway>();

        return services;
    }
}
