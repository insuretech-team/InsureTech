using Microsoft.Extensions.DependencyInjection;
using PoliSync.VehicleInsurance.Services;

namespace PoliSync.VehicleInsurance;

public static class DependencyInjection
{
    public static IServiceCollection AddVehicleInsuranceServices(this IServiceCollection services)
    {
        services.AddSingleton<IVehicleService, GoVehicleService>();
        
        return services;
    }
}
