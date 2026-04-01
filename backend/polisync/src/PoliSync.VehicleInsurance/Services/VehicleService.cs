using Insuretech.Vehicle.Entity.V1;
using Microsoft.Extensions.Logging;

namespace PoliSync.VehicleInsurance.Services;

public class VehicleInsuranceService : IVehicleService
{
    private readonly IVehicleRepository _vehicleRepository;
    private readonly IVehicleRegistrationRepository _registrationRepository;
    private readonly IVehiclePremiumCalculator _premiumCalculator;
    private readonly ILogger<VehicleInsuranceService> _logger;

    public VehicleInsuranceService(
        IVehicleRepository vehicleRepository,
        IVehicleRegistrationRepository registrationRepository,
        IVehiclePremiumCalculator premiumCalculator,
        ILogger<VehicleInsuranceService> logger)
    {
        _vehicleRepository = vehicleRepository;
        _registrationRepository = registrationRepository;
        _premiumCalculator = premiumCalculator;
        _logger = logger;
    }

    public Task<Vehicle?> GetVehicleAsync(string vehicleId, CancellationToken cancellationToken = default)
    {
        return _vehicleRepository.GetByIdAsync(vehicleId, cancellationToken);
    }

    public Task<Vehicle?> GetVehicleByModelAsync(string model, CancellationToken cancellationToken = default)
    {
        return _vehicleRepository.GetByModelAsync(model, cancellationToken);
    }

    public Task<IEnumerable<Vehicle>> ListVehiclesAsync(
        VehicleType? vehicleType, 
        string? manufacturer, 
        bool onlyActive, 
        CancellationToken cancellationToken = default)
    {
        return _vehicleRepository.GetByFilterAsync(vehicleType, manufacturer, onlyActive, cancellationToken);
    }

    public Task<Vehicle> CreateVehicleAsync(Vehicle vehicle, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Creating vehicle: {Model}", vehicle.Model);
        return _vehicleRepository.CreateAsync(vehicle, cancellationToken);
    }

    public Task<Vehicle?> UpdateVehicleAsync(string vehicleId, Vehicle vehicle, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Updating vehicle: {VehicleId}", vehicleId);
        vehicle.VehicleId = vehicleId;
        return _vehicleRepository.UpdateAsync(vehicle, cancellationToken);
    }

    public Task<bool> DeleteVehicleAsync(string vehicleId, bool permanent = false, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Deleting vehicle: {VehicleId} (permanent: {Permanent})", vehicleId, permanent);
        return _vehicleRepository.DeleteAsync(vehicleId, permanent, cancellationToken);
    }

    public async Task<IEnumerable<string>> GetVehicleModelsAsync(VehicleType vehicleType, CancellationToken cancellationToken = default)
    {
        var vehicles = await _vehicleRepository.GetByTypeAsync(vehicleType, cancellationToken);
        return vehicles.Select(v => v.Model).Distinct().OrderBy(m => m);
    }

    public Task<VehicleRegistration> RegisterVehicleAsync(VehicleRegistration registration, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation(
            "Registering vehicle: {RegistrationNumber} for owner {OwnerId}",
            registration.RegistrationNumber, registration.OwnerId);
        return _registrationRepository.CreateAsync(registration, cancellationToken);
    }

    public Task<VehicleRegistration?> GetVehicleRegistrationAsync(string registrationId, CancellationToken cancellationToken = default)
    {
        return _registrationRepository.GetByIdAsync(registrationId, cancellationToken);
    }

    public async Task<VehiclePremiumCalculation> CalculatePremiumAsync(
        string registrationId, 
        bool accidentalCover, 
        CancellationToken cancellationToken = default)
    {
        var registration = await _registrationRepository.GetByIdAsync(registrationId, cancellationToken);
        if (registration == null)
        {
            throw new InvalidOperationException($"Vehicle registration {registrationId} not found");
        }

        var vehicle = await _vehicleRepository.GetByIdAsync(registration.VehicleId, cancellationToken);
        if (vehicle == null)
        {
            throw new InvalidOperationException($"Vehicle {registration.VehicleId} not found");
        }

        return await _premiumCalculator.CalculatePremiumAsync(vehicle, registration, accidentalCover, cancellationToken);
    }

    public Task<IEnumerable<VehiclePremiumCalculation>> GetPremiumCalculationsAsync(
        string registrationId, 
        CancellationToken cancellationToken = default)
    {
        // For now, return empty as we don't store calculation history in memory
        return Task.FromResult(Enumerable.Empty<VehiclePremiumCalculation>());
    }

    public Task<Dictionary<string, string>> GetVehicleImagesAsync(CancellationToken cancellationToken = default)
    {
        // Return default images
        var images = new Dictionary<string, string>
        {
            { "bike", "/assets/images/bike.jpg" },
            { "car", "/assets/images/car.jpg" },
            { "commercial", "/assets/images/commercial.jpg" }
        };
        return Task.FromResult(images);
    }
}
