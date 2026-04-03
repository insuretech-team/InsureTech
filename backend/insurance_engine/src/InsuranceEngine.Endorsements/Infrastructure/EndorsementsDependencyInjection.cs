using Microsoft.Extensions.DependencyInjection;
using MediatR;
using InsuranceEngine.Endorsements.Application;

namespace InsuranceEngine.Endorsements;

public static class EndorsementsDependencyInjection
{
    public static IServiceCollection AddEndorsementsModule(this IServiceCollection services)
    {
        services.AddMediatR(cfg => 
        {
            cfg.RegisterServicesFromAssembly(typeof(AssemblyMarker).Assembly);
        });

        services.AddScoped<IEndorsementDataGateway, Infrastructure.GoEndorsementDataGateway>();
        services.AddScoped<IEndorsementProcessingService, EndorsementProcessingService>();

        return services;
    }
}
