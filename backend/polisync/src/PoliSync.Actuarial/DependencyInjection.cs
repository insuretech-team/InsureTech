using Microsoft.Extensions.DependencyInjection;
using PoliSync.Actuarial.Services;

namespace PoliSync.Actuarial;

public static class DependencyInjection
{
    public static IServiceCollection AddActuarialServices(this IServiceCollection services)
    {
        services.AddSingleton<GoActuarialService>();
        services.AddSingleton<IActuarialService>(sp => sp.GetRequiredService<GoActuarialService>());
        services.AddSingleton<IRatingFormulaService>(sp => sp.GetRequiredService<GoActuarialService>());
        services.AddSingleton<IReserveCalculationService>(sp => sp.GetRequiredService<GoActuarialService>());
        services.AddSingleton<ILossRatioService>(sp => sp.GetRequiredService<GoActuarialService>());
        services.AddSingleton<IFormulaEvaluator, FormulaEvaluator>();
        
        return services;
    }
}
