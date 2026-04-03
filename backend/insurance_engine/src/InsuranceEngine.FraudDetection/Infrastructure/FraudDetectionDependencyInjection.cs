using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Options;
using InsuranceEngine.FraudDetection.Infrastructure;

namespace InsuranceEngine.FraudDetection;

public static class FraudDetectionDependencyInjection
{
    public static IServiceCollection AddFraudDetectionModule(this IServiceCollection services)
    {
        services.AddScoped<IFraudDetectionDataGateway, GoFraudDetectionDataGateway>();
        services.AddScoped<IFraudDetectionService, FraudDetectionService>();

        services.Configure<FraudCheckSettings>(options =>
        {
            options.RapidClaimHoursThreshold = 48;
            options.ClaimFrequencyThreshold = 2;
            options.ClaimFrequencyWindowMonths = 12;
            options.FullCoverageClaimThreshold = 1.0m;
            options.DeviceAccountThreshold = 3;
            options.EnablePatternAnalysis = true;
            options.EnableProviderValidation = true;
        });

        return services;
    }
}
