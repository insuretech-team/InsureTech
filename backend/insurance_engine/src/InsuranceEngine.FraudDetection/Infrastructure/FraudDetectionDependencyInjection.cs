using Microsoft.Extensions.DependencyInjection;
using InsuranceEngine.Grpc.Gateways;

namespace InsuranceEngine.FraudDetection;

public static class FraudDetectionDependencyInjection
{
    public static IServiceCollection AddFraudDetectionModule(this IServiceCollection services)
    {
        // Data Gateways (PoliSync Standard)
        services.AddScoped<IFraudDetectionDataGateway, GoFraudDetectionDataGateway>();
        
        return services;
    }
}
