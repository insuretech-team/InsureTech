using System.Collections.Concurrent;
using System.Text.Json;
using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Vehicle.Entity.V1;
using Microsoft.Extensions.Logging;

namespace PoliSync.VehicleInsurance.Services;

public class InMemoryVehicleRepository : IVehicleRepository
{
    private readonly ConcurrentDictionary<string, Vehicle> _vehicles = new();
    private readonly ILogger<InMemoryVehicleRepository> _logger;

    public InMemoryVehicleRepository(ILogger<InMemoryVehicleRepository> logger)
    {
        _logger = logger;
        SeedDefaultVehicles();
    }

    public Task<Vehicle?> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        _vehicles.TryGetValue(id, out var vehicle);
        return Task.FromResult(vehicle);
    }

    public Task<Vehicle?> GetByModelAsync(string model, CancellationToken cancellationToken = default)
    {
        var vehicle = _vehicles.Values.FirstOrDefault(v => 
            v.Model.Equals(model, StringComparison.OrdinalIgnoreCase));
        return Task.FromResult(vehicle);
    }

    public Task<IEnumerable<Vehicle>> GetAllAsync(CancellationToken cancellationToken = default)
    {
        return Task.FromResult(_vehicles.Values.AsEnumerable());
    }

    public Task<IEnumerable<Vehicle>> GetByTypeAsync(VehicleType type, CancellationToken cancellationToken = default)
    {
        var vehicles = _vehicles.Values
            .Where(v => v.Type == type)
            .AsEnumerable();
        return Task.FromResult(vehicles);
    }

    public Task<IEnumerable<Vehicle>> GetByFilterAsync(
        VehicleType? type, 
        string? manufacturer, 
        bool onlyActive, 
        CancellationToken cancellationToken = default)
    {
        var query = _vehicles.Values.AsEnumerable();
        
        if (type.HasValue)
        {
            query = query.Where(v => v.Type == type.Value);
        }
        
        if (!string.IsNullOrEmpty(manufacturer))
        {
            query = query.Where(v => v.Manufacturer.Contains(manufacturer, StringComparison.OrdinalIgnoreCase));
        }
        
        if (onlyActive)
        {
            query = query.Where(v => v.IsActive);
        }
        
        return Task.FromResult(query);
    }

    public Task<Vehicle> CreateAsync(Vehicle vehicle, CancellationToken cancellationToken = default)
    {
        if (string.IsNullOrEmpty(vehicle.VehicleId))
        {
            vehicle.VehicleId = Guid.NewGuid().ToString();
        }
        
        vehicle.CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        vehicle.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        
        _vehicles[vehicle.VehicleId] = vehicle;
        _logger.LogInformation("Created vehicle: {VehicleId} - {Model}", 
            vehicle.VehicleId, vehicle.Model);
        
        return Task.FromResult(vehicle);
    }

    public Task<Vehicle?> UpdateAsync(Vehicle vehicle, CancellationToken cancellationToken = default)
    {
        if (!_vehicles.ContainsKey(vehicle.VehicleId))
        {
            return Task.FromResult<Vehicle?>(null);
        }

        vehicle.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        _vehicles[vehicle.VehicleId] = vehicle;
        
        _logger.LogInformation("Updated vehicle: {VehicleId}", vehicle.VehicleId);
        
        return Task.FromResult<Vehicle?>(vehicle);
    }

    public Task<bool> DeleteAsync(string id, bool permanent = false, CancellationToken cancellationToken = default)
    {
        if (permanent)
        {
            var result = _vehicles.TryRemove(id, out _);
            if (result)
            {
                _logger.LogInformation("Permanently deleted vehicle: {VehicleId}", id);
            }
            return Task.FromResult(result);
        }
        else
        {
            if (_vehicles.TryGetValue(id, out var vehicle))
            {
                vehicle.DeletedAt = Timestamp.FromDateTime(DateTime.UtcNow);
                vehicle.IsActive = false;
                _logger.LogInformation("Soft deleted vehicle: {VehicleId}", id);
                return Task.FromResult(true);
            }
            return Task.FromResult(false);
        }
    }

    private void SeedDefaultVehicles()
    {
        // Bike models
        var bikes = new[]
        {
            new Vehicle { VehicleId = Guid.NewGuid().ToString(), Model = "Honda Activa", Type = VehicleType.Bike, Price = 7500000, Manufacturer = "Honda", Year = 2024, IsActive = true },
            new Vehicle { VehicleId = Guid.NewGuid().ToString(), Model = "Bajaj Pulsar", Type = VehicleType.Bike, Price = 12000000, Manufacturer = "Bajaj", Year = 2024, IsActive = true },
            new Vehicle { VehicleId = Guid.NewGuid().ToString(), Model = "Royal Enfield Classic", Type = VehicleType.Bike, Price = 20000000, Manufacturer = "Royal Enfield", Year = 2024, IsActive = true },
            new Vehicle { VehicleId = Guid.NewGuid().ToString(), Model = "TVS Apache", Type = VehicleType.Bike, Price = 14000000, Manufacturer = "TVS", Year = 2024, IsActive = true },
            new Vehicle { VehicleId = Guid.NewGuid().ToString(), Model = "Hero Splendor", Type = VehicleType.Bike, Price = 8000000, Manufacturer = "Hero", Year = 2024, IsActive = true }
        };

        // Car models
        var cars = new[]
        {
            new Vehicle { VehicleId = Guid.NewGuid().ToString(), Model = "Maruti Swift", Type = VehicleType.Car, Price = 60000000, Manufacturer = "Maruti", Year = 2024, IsActive = true },
            new Vehicle { VehicleId = Guid.NewGuid().ToString(), Model = "Hyundai i20", Type = VehicleType.Car, Price = 80000000, Manufacturer = "Hyundai", Year = 2024, IsActive = true },
            new Vehicle { VehicleId = Guid.NewGuid().ToString(), Model = "Honda City", Type = VehicleType.Car, Price = 120000000, Manufacturer = "Honda", Year = 2024, IsActive = true },
            new Vehicle { VehicleId = Guid.NewGuid().ToString(), Model = "Toyota Innova", Type = VehicleType.Car, Price = 200000000, Manufacturer = "Toyota", Year = 2024, IsActive = true },
            new Vehicle { VehicleId = Guid.NewGuid().ToString(), Model = "Kia Seltos", Type = VehicleType.Car, Price = 150000000, Manufacturer = "Kia", Year = 2024, IsActive = true }
        };

        foreach (var bike in bikes)
        {
            bike.CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            bike.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            _vehicles[bike.VehicleId] = bike;
        }

        foreach (var car in cars)
        {
            car.CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            car.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            _vehicles[car.VehicleId] = car;
        }

        _logger.LogInformation("Seeded {Count} default vehicles", _vehicles.Count);
    }
}
