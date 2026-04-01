using System.Diagnostics;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Vehicle.Entity.V1;
using Insuretech.Common.V1;
using Microsoft.Extensions.Logging;

namespace PoliSync.VehicleInsurance.Services;

public class VehiclePremiumCalculator : IVehiclePremiumCalculator
{
    private readonly ILogger<VehiclePremiumCalculator> _logger;

    public VehiclePremiumCalculator(ILogger<VehiclePremiumCalculator> logger)
    {
        _logger = logger;
    }

    public Task<VehiclePremiumCalculation> CalculatePremiumAsync(
        Vehicle vehicle,
        VehicleRegistration registration,
        bool accidentalCover,
        CancellationToken cancellationToken = default)
    {
        var stopwatch = Stopwatch.StartNew();
        var calculationId = Guid.NewGuid().ToString();

        _logger.LogInformation(
            "Calculating premium for vehicle {Model}, registration {RegNumber}",
            vehicle.Model, registration.RegistrationNumber);

        // Calculate multipliers based on reference VehicleInsuranceAPI logic
        
        // 1. Type Multiplier
        var typeMultiplier = vehicle.Type switch
        {
            VehicleType.Bike => 1.0f,
            VehicleType.Car => 1.2f,
            _ => 1.0f
        };

        // 2. Age Multiplier
        var age = registration.VehicleAge;
        var ageMultiplier = age switch
        {
            0 or 1 or 2 => 1.0f,
            3 or 4 or 5 => 1.2f,
            _ => 1.5f
        };

        // 3. Value Multiplier (based on depreciation)
        // Calculate depreciated value: Price * (1 - (AgeMultiplier - 1))
        var depreciationFactor = (ageMultiplier - 1.0f);
        var currentValue = (long)(vehicle.Price * (1 - depreciationFactor));
        var minValue = (long)(vehicle.Price * 0.5f);
        
        // Ensure current value doesn't go below 50%
        if (currentValue < minValue)
        {
            currentValue = minValue;
        }
        
        var valueMultiplier = currentValue switch
        {
            _ when currentValue == vehicle.Price => 1.0f,
            _ when currentValue >= minValue => 1.2f,
            _ => 1.5f
        };

        // 4. Location Multiplier
        var nearbyStates = new[] { "AP", "TL", "KL", "KA" }; // Nearby states to TN
        var location = registration.RegistrationState.ToUpper();
        
        var locationMultiplier = location switch
        {
            "TN" => 1.0f,
            _ when nearbyStates.Contains(location) => 1.1f,
            _ => 1.2f
        };

        // Calculate base premium: 10% of original price
        var basePremiumAmount = vehicle.Price * 0.10m;

        // Apply multipliers
        var adjustedPremium = basePremiumAmount * (decimal)typeMultiplier;
        adjustedPremium = adjustedPremium * (decimal)ageMultiplier;
        adjustedPremium = adjustedPremium * (decimal)valueMultiplier;
        adjustedPremium = adjustedPremium * (decimal)locationMultiplier;

        var totalPremium = (long)Math.Round(adjustedPremium);

        // Calculate Third Party premiums for 1, 2, 3 years
        var tpPremium1Year = totalPremium;
        var tpPremium2Year = totalPremium * 2;
        var tpPremium3Year = totalPremium * 3;

        // Calculate Comprehensive premiums (1.2x Third Party)
        var compPremium1Year = (long)Math.Round(tpPremium1Year * 1.2m);
        var compPremium2Year = (long)Math.Round(tpPremium2Year * 1.2m);
        var compPremium3Year = (long)Math.Round(tpPremium3Year * 1.2m);

        stopwatch.Stop();

        _logger.LogInformation(
            "Premium calculated: Base={Base}, Total={Total}, Multipliers: Type={TypeM}, Age={AgeM}, Value={ValueM}, Location={LocM}, Time={Ms}ms",
            basePremiumAmount, totalPremium, typeMultiplier, ageMultiplier, valueMultiplier, locationMultiplier, stopwatch.ElapsedMilliseconds);

        return Task.FromResult(new VehiclePremiumCalculation
        {
            CalculationId = calculationId,
            RegistrationId = registration.RegistrationId,
            BasePremium = new Money { Amount = (long)(basePremiumAmount), Currency = "BDT" },
            TypeMultiplier = typeMultiplier,
            AgeMultiplier = ageMultiplier,
            ValueMultiplier = valueMultiplier,
            LocationMultiplier = locationMultiplier,
            TpPremium1Year = new Money { Amount = tpPremium1Year, Currency = "BDT" },
            TpPremium2Year = new Money { Amount = tpPremium2Year, Currency = "BDT" },
            TpPremium3Year = new Money { Amount = tpPremium3Year, Currency = "BDT" },
            CompPremium1Year = new Money { Amount = compPremium1Year, Currency = "BDT" },
            CompPremium2Year = new Money { Amount = compPremium2Year, Currency = "BDT" },
            CompPremium3Year = new Money { Amount = compPremium3Year, Currency = "BDT" },
            AccidentalCover = accidentalCover,
            CalculatedAt = Timestamp.FromDateTime(DateTime.UtcNow),
            CalculationDurationMs = (int)stopwatch.ElapsedMilliseconds
        });
    }
}
