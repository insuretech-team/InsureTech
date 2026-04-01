using System.Collections.Concurrent;
using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Vehicle.Entity.V1;
using Microsoft.Extensions.Logging;

namespace PoliSync.VehicleInsurance.Services;

public class InMemoryVehicleRegistrationRepository : IVehicleRegistrationRepository
{
    private readonly ConcurrentDictionary<string, VehicleRegistration> _registrations = new();
    private readonly ILogger<InMemoryVehicleRegistrationRepository> _logger;

    public InMemoryVehicleRegistrationRepository(ILogger<InMemoryVehicleRegistrationRepository> logger)
    {
        _logger = logger;
    }

    public Task<VehicleRegistration?> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        _registrations.TryGetValue(id, out var registration);
        return Task.FromResult(registration);
    }

    public Task<VehicleRegistration?> GetByRegistrationNumberAsync(string number, CancellationToken cancellationToken = default)
    {
        var registration = _registrations.Values.FirstOrDefault(r => 
            r.RegistrationNumber.Equals(number, StringComparison.OrdinalIgnoreCase));
        return Task.FromResult(registration);
    }

    public Task<IEnumerable<VehicleRegistration>> GetByOwnerIdAsync(string ownerId, CancellationToken cancellationToken = default)
    {
        var registrations = _registrations.Values
            .Where(r => r.OwnerId == ownerId)
            .AsEnumerable();
        return Task.FromResult(registrations);
    }

    public Task<VehicleRegistration> CreateAsync(VehicleRegistration registration, CancellationToken cancellationToken = default)
    {
        if (string.IsNullOrEmpty(registration.RegistrationId))
        {
            registration.RegistrationId = Guid.NewGuid().ToString();
        }
        
        // Calculate age and current value
        var currentYear = DateTime.UtcNow.Year;
        registration.VehicleAge = currentYear - registration.RegistrationYear;
        
        // Set state from registration number (first 2 chars)
        if (registration.RegistrationNumber.Length >= 2)
        {
            registration.RegistrationState = registration.RegistrationNumber.Substring(0, 2).ToUpper();
        }
        
        registration.CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        registration.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        
        _registrations[registration.RegistrationId] = registration;
        _logger.LogInformation("Created vehicle registration: {RegistrationId} - {RegistrationNumber}", 
            registration.RegistrationId, registration.RegistrationNumber);
        
        return Task.FromResult(registration);
    }

    public Task<VehicleRegistration?> UpdateAsync(VehicleRegistration registration, CancellationToken cancellationToken = default)
    {
        if (!_registrations.ContainsKey(registration.RegistrationId))
        {
            return Task.FromResult<VehicleRegistration?>(null);
        }

        registration.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        _registrations[registration.RegistrationId] = registration;
        
        _logger.LogInformation("Updated vehicle registration: {RegistrationId}", registration.RegistrationId);
        
        return Task.FromResult<VehicleRegistration?>(registration);
    }
}
