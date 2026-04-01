using Insuretech.Vehicle.Entity.V1;
using Insuretech.Vehicle.Services.V1;
using PoliSync.Infrastructure.Clients;

namespace PoliSync.VehicleInsurance.Services;

public sealed class GoVehicleService : IVehicleService
{
    private readonly InsuranceServiceClient _client;

    public GoVehicleService(InsuranceServiceClient client) => _client = client;

    public async Task<Vehicle?> GetVehicleAsync(string vehicleId, CancellationToken cancellationToken = default)
    {
        var response = await _client.VehicleClient.GetVehicleAsync(
            new GetVehicleRequest { VehicleId = vehicleId },
            _client.BuildCallOptions(cancellationToken));
        return response.Error is null ? response.Vehicle : null;
    }

    public async Task<Vehicle?> GetVehicleByModelAsync(string model, CancellationToken cancellationToken = default)
    {
        var response = await _client.VehicleClient.GetVehicleByModelAsync(
            new GetVehicleByModelRequest { Model = model },
            _client.BuildCallOptions(cancellationToken));
        return response.Error is null ? response.Vehicle : null;
    }

    public async Task<IEnumerable<Vehicle>> ListVehiclesAsync(VehicleType? vehicleType, string? manufacturer, bool onlyActive, CancellationToken cancellationToken = default)
    {
        var response = await _client.VehicleClient.ListVehiclesAsync(
            new ListVehiclesRequest
            {
                VehicleType = vehicleType ?? VehicleType.Unspecified,
                Manufacturer = manufacturer ?? string.Empty,
                OnlyActive = onlyActive,
                PageSize = 200
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Vehicles;
    }

    public async Task<Vehicle> CreateVehicleAsync(Vehicle vehicle, CancellationToken cancellationToken = default)
    {
        var response = await _client.VehicleClient.CreateVehicleAsync(
            new CreateVehicleRequest
            {
                Model = vehicle.Model,
                Type = vehicle.Type,
                Price = vehicle.Price,
                Manufacturer = vehicle.Manufacturer,
                Year = vehicle.Year,
                ImageUri = vehicle.ImageUri,
                Metadata = { vehicle.Metadata }
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Vehicle;
    }

    public async Task<Vehicle?> UpdateVehicleAsync(string vehicleId, Vehicle vehicle, CancellationToken cancellationToken = default)
    {
        var response = await _client.VehicleClient.UpdateVehicleAsync(
            new UpdateVehicleRequest
            {
                VehicleId = vehicleId,
                Model = vehicle.Model,
                Type = vehicle.Type,
                Price = vehicle.Price,
                Manufacturer = vehicle.Manufacturer,
                Year = vehicle.Year,
                ImageUri = vehicle.ImageUri,
                IsActive = vehicle.IsActive,
                Metadata = { vehicle.Metadata }
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Error is null ? response.Vehicle : null;
    }

    public async Task<bool> DeleteVehicleAsync(string vehicleId, bool permanent = false, CancellationToken cancellationToken = default)
    {
        var response = await _client.VehicleClient.DeleteVehicleAsync(
            new DeleteVehicleRequest { VehicleId = vehicleId, Permanent = permanent },
            _client.BuildCallOptions(cancellationToken));
        return response.Success && response.Error is null;
    }

    public async Task<IEnumerable<string>> GetVehicleModelsAsync(VehicleType vehicleType, CancellationToken cancellationToken = default)
    {
        var response = await _client.VehicleClient.GetVehicleModelsAsync(
            new GetVehicleModelsRequest { VehicleType = vehicleType },
            _client.BuildCallOptions(cancellationToken));
        return response.Models;
    }

    public async Task<VehicleRegistration> RegisterVehicleAsync(VehicleRegistration registration, CancellationToken cancellationToken = default)
    {
        var response = await _client.VehicleClient.RegisterVehicleAsync(
            new RegisterVehicleRequest
            {
                VehicleId = registration.VehicleId,
                OwnerId = registration.OwnerId,
                RegistrationNumber = registration.RegistrationNumber,
                RegistrationYear = registration.RegistrationYear,
                AdditionalInfo = { registration.AdditionalInfo }
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Registration;
    }

    public async Task<VehicleRegistration?> GetVehicleRegistrationAsync(string registrationId, CancellationToken cancellationToken = default)
    {
        var response = await _client.VehicleClient.GetVehicleRegistrationAsync(
            new GetVehicleRegistrationRequest { RegistrationId = registrationId },
            _client.BuildCallOptions(cancellationToken));
        return response.Error is null ? response.Registration : null;
    }

    public async Task<VehiclePremiumCalculation> CalculatePremiumAsync(string registrationId, bool accidentalCover, CancellationToken cancellationToken = default)
    {
        var response = await _client.VehicleClient.CalculatePremiumAsync(
            new CalculatePremiumRequest
            {
                RegistrationId = registrationId,
                AccidentalCover = accidentalCover
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Calculation;
    }

    public async Task<IEnumerable<VehiclePremiumCalculation>> GetPremiumCalculationsAsync(string registrationId, CancellationToken cancellationToken = default)
    {
        var response = await _client.VehicleClient.GetPremiumCalculationsAsync(
            new GetPremiumCalculationsRequest
            {
                RegistrationId = registrationId,
                PageSize = 200
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Calculations;
    }

    public async Task<Dictionary<string, string>> GetVehicleImagesAsync(CancellationToken cancellationToken = default)
    {
        var response = await _client.VehicleClient.GetVehicleImagesAsync(
            new GetVehicleImagesRequest(),
            _client.BuildCallOptions(cancellationToken));
        return response.Images.ToDictionary(kvp => kvp.Key, kvp => kvp.Value);
    }
}
