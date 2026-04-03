using Microsoft.Extensions.DependencyInjection;
using InsuranceEngine.Grpc.Gateways;

namespace InsuranceEngine.Commission;

public static class CommissionDependencyInjection
{
    public static IServiceCollection AddCommissionModule(this IServiceCollection services)
    {
        // Data Gateways (PoliSync Standard)
        services.AddScoped<ICommissionDataGateway, GoCommissionDataGateway>();
        
        return services;
    }
}
