using Insuretech.Vehicle.Entity.V1;
using Insuretech.Common.V1;

namespace PoliSync.VehicleInsurance.Services;

public interface IVehicleService
{
    Task<Vehicle?> GetVehicleAsync(string vehicleId, CancellationToken cancellationToken = default);
    Task<Vehicle?> GetVehicleByModelAsync(string model, CancellationToken cancellationToken = default);
    Task<IEnumerable<Vehicle>> ListVehiclesAsync(VehicleType? vehicleType, string? manufacturer, bool onlyActive, CancellationToken cancellationToken = default);
    Task<Vehicle> CreateVehicleAsync(Vehicle vehicle, CancellationToken cancellationToken = default);
    Task<Vehicle?> UpdateVehicleAsync(string vehicleId, Vehicle vehicle, CancellationToken cancellationToken = default);
    Task<bool> DeleteVehicleAsync(string vehicleId, bool permanent = false, CancellationToken cancellationToken = default);
    Task<IEnumerable<string>> GetVehicleModelsAsync(VehicleType vehicleType, CancellationToken cancellationToken = default);
    
    Task<VehicleRegistration> RegisterVehicleAsync(VehicleRegistration registration, CancellationToken cancellationToken = default);
    Task<VehicleRegistration?> GetVehicleRegistrationAsync(string registrationId, CancellationToken cancellationToken = default);
    Task<VehiclePremiumCalculation> CalculatePremiumAsync(string registrationId, bool accidentalCover, CancellationToken cancellationToken = default);
    Task<IEnumerable<VehiclePremiumCalculation>> GetPremiumCalculationsAsync(string registrationId, CancellationToken cancellationToken = default);
    
    Task<Dictionary<string, string>> GetVehicleImagesAsync(CancellationToken cancellationToken = default);
}

public interface IVehicleRepository
{
    Task<Vehicle?> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<Vehicle?> GetByModelAsync(string model, CancellationToken cancellationToken = default);
    Task<IEnumerable<Vehicle>> GetAllAsync(CancellationToken cancellationToken = default);
    Task<IEnumerable<Vehicle>> GetByTypeAsync(VehicleType type, CancellationToken cancellationToken = default);
    Task<IEnumerable<Vehicle>> GetByFilterAsync(VehicleType? type, string? manufacturer, bool onlyActive, CancellationToken cancellationToken = default);
    Task<Vehicle> CreateAsync(Vehicle vehicle, CancellationToken cancellationToken = default);
    Task<Vehicle?> UpdateAsync(Vehicle vehicle, CancellationToken cancellationToken = default);
    Task<bool> DeleteAsync(string id, bool permanent = false, CancellationToken cancellationToken = default);
}

public interface IVehicleRegistrationRepository
{
    Task<VehicleRegistration?> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<VehicleRegistration?> GetByRegistrationNumberAsync(string number, CancellationToken cancellationToken = default);
    Task<IEnumerable<VehicleRegistration>> GetByOwnerIdAsync(string ownerId, CancellationToken cancellationToken = default);
    Task<VehicleRegistration> CreateAsync(VehicleRegistration registration, CancellationToken cancellationToken = default);
    Task<VehicleRegistration?> UpdateAsync(VehicleRegistration registration, CancellationToken cancellationToken = default);
}

public interface IVehiclePremiumCalculator
{
    Task<VehiclePremiumCalculation> CalculatePremiumAsync(
        Vehicle vehicle,
        VehicleRegistration registration,
        bool accidentalCover,
        CancellationToken cancellationToken = default);
}
