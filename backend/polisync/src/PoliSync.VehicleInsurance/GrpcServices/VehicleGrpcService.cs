using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Insuretech.Common.V1;
using Insuretech.Vehicle.Entity.V1;
using Insuretech.Vehicle.Services.V1;
using Microsoft.Extensions.Logging;
using PoliSync.VehicleInsurance.Services;

namespace PoliSync.VehicleInsurance.GrpcServices;

public class VehicleGrpcService : VehicleService.VehicleServiceBase
{
    private readonly IVehicleService _vehicleService;
    private readonly ILogger<VehicleGrpcService> _logger;

    public VehicleGrpcService(
        IVehicleService vehicleService,
        ILogger<VehicleGrpcService> logger)
    {
        _vehicleService = vehicleService;
        _logger = logger;
    }

    public override async Task<GetVehicleResponse> GetVehicle(
        GetVehicleRequest request,
        ServerCallContext context)
    {
        var vehicle = await _vehicleService.GetVehicleAsync(request.VehicleId, context.CancellationToken);
        
        if (vehicle == null)
        {
            return new GetVehicleResponse
            {
                Error = new Error
                {
                    Code = "VEHICLE_NOT_FOUND",
                    Message = $"Vehicle {request.VehicleId} not found",
                    HttpStatusCode = 404
                }
            };
        }

        return new GetVehicleResponse { Vehicle = vehicle };
    }

    public override async Task<GetVehicleResponse> GetVehicleByModel(
        GetVehicleByModelRequest request,
        ServerCallContext context)
    {
        var vehicle = await _vehicleService.GetVehicleByModelAsync(request.Model, context.CancellationToken);
        
        if (vehicle == null)
        {
            return new GetVehicleResponse
            {
                Error = new Error
                {
                    Code = "VEHICLE_NOT_FOUND",
                    Message = $"Vehicle model '{request.Model}' not found",
                    HttpStatusCode = 404
                }
            };
        }

        return new GetVehicleResponse { Vehicle = vehicle };
    }

    public override async Task<ListVehiclesResponse> ListVehicles(
        ListVehiclesRequest request,
        ServerCallContext context)
    {
        var vehicles = await _vehicleService.ListVehiclesAsync(
            request.VehicleType != VehicleType.Unspecified ? request.VehicleType : null,
            string.IsNullOrEmpty(request.Manufacturer) ? null : request.Manufacturer,
            request.OnlyActive,
            context.CancellationToken);

        var vehicleList = vehicles.ToList();

        return new ListVehiclesResponse
        {
            Vehicles = { vehicleList },
            TotalCount = vehicleList.Count
        };
    }

    public override async Task<CreateVehicleResponse> CreateVehicle(
        CreateVehicleRequest request,
        ServerCallContext context)
    {
        try
        {
            var vehicle = new Vehicle
            {
                Model = request.Model,
                Type = request.Type,
                Price = request.Price,
                Manufacturer = request.Manufacturer,
                Year = request.Year,
                ImageUri = request.ImageUri,
                Metadata = { request.Metadata },
                IsActive = true
            };

            var created = await _vehicleService.CreateVehicleAsync(vehicle, context.CancellationToken);

            return new CreateVehicleResponse { Vehicle = created };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error creating vehicle");
            return new CreateVehicleResponse
            {
                Error = new Error
                {
                    Code = "CREATE_ERROR",
                    Message = ex.Message,
                    HttpStatusCode = 500
                }
            };
        }
    }

    public override async Task<UpdateVehicleResponse> UpdateVehicle(
        UpdateVehicleRequest request,
        ServerCallContext context)
    {
        try
        {
            var vehicle = new Vehicle
            {
                VehicleId = request.VehicleId,
                Model = request.Model,
                Type = request.Type,
                Price = request.Price,
                Manufacturer = request.Manufacturer,
                Year = request.Year,
                ImageUri = request.ImageUri,
                IsActive = request.IsActive,
                Metadata = { request.Metadata }
            };

            var updated = await _vehicleService.UpdateVehicleAsync(
                request.VehicleId,
                vehicle,
                context.CancellationToken);

            if (updated == null)
            {
                return new UpdateVehicleResponse
                {
                    Error = new Error
                    {
                        Code = "VEHICLE_NOT_FOUND",
                        Message = $"Vehicle {request.VehicleId} not found",
                        HttpStatusCode = 404
                    }
                };
            }

            return new UpdateVehicleResponse { Vehicle = updated };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error updating vehicle");
            return new UpdateVehicleResponse
            {
                Error = new Error
                {
                    Code = "UPDATE_ERROR",
                    Message = ex.Message,
                    HttpStatusCode = 500
                }
            };
        }
    }

    public override async Task<DeleteVehicleResponse> DeleteVehicle(
        DeleteVehicleRequest request,
        ServerCallContext context)
    {
        var result = await _vehicleService.DeleteVehicleAsync(
            request.VehicleId,
            request.Permanent,
            context.CancellationToken);

        return new DeleteVehicleResponse { Success = result };
    }

    public override async Task<GetVehicleModelsResponse> GetVehicleModels(
        GetVehicleModelsRequest request,
        ServerCallContext context)
    {
        var models = await _vehicleService.GetVehicleModelsAsync(
            request.VehicleType,
            context.CancellationToken);

        return new GetVehicleModelsResponse
        {
            Models = { models }
        };
    }

    public override async Task<RegisterVehicleResponse> RegisterVehicle(
        RegisterVehicleRequest request,
        ServerCallContext context)
    {
        try
        {
            var registration = new VehicleRegistration
            {
                VehicleId = request.VehicleId,
                OwnerId = request.OwnerId,
                RegistrationNumber = request.RegistrationNumber,
                RegistrationYear = request.RegistrationYear,
                Status = RegistrationStatus.Active,
                AdditionalInfo = { request.AdditionalInfo }
            };

            var created = await _vehicleService.RegisterVehicleAsync(registration, context.CancellationToken);

            return new RegisterVehicleResponse { Registration = created };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error registering vehicle");
            return new RegisterVehicleResponse
            {
                Error = new Error
                {
                    Code = "REGISTER_ERROR",
                    Message = ex.Message,
                    HttpStatusCode = 500
                }
            };
        }
    }

    public override async Task<GetVehicleRegistrationResponse> GetVehicleRegistration(
        GetVehicleRegistrationRequest request,
        ServerCallContext context)
    {
        var registration = await _vehicleService.GetVehicleRegistrationAsync(
            request.RegistrationId,
            context.CancellationToken);

        if (registration == null)
        {
            return new GetVehicleRegistrationResponse
            {
                Error = new Error
                {
                    Code = "REGISTRATION_NOT_FOUND",
                    Message = $"Registration {request.RegistrationId} not found",
                    HttpStatusCode = 404
                }
            };
        }

        var vehicle = await _vehicleService.GetVehicleAsync(registration.VehicleId, context.CancellationToken);

        return new GetVehicleRegistrationResponse
        {
            Registration = registration,
            Vehicle = vehicle
        };
    }

    public override async Task<CalculatePremiumResponse> CalculatePremium(
        CalculatePremiumRequest request,
        ServerCallContext context)
    {
        try
        {
            var calculation = await _vehicleService.CalculatePremiumAsync(
                request.RegistrationId,
                request.AccidentalCover,
                context.CancellationToken);

            return new CalculatePremiumResponse { Calculation = calculation };
        }
        catch (InvalidOperationException ex)
        {
            return new CalculatePremiumResponse
            {
                Error = new Error
                {
                    Code = "NOT_FOUND",
                    Message = ex.Message,
                    HttpStatusCode = 404
                }
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error calculating premium");
            return new CalculatePremiumResponse
            {
                Error = new Error
                {
                    Code = "CALCULATION_ERROR",
                    Message = ex.Message,
                    HttpStatusCode = 500
                }
            };
        }
    }

    public override async Task<GetPremiumCalculationsResponse> GetPremiumCalculations(
        GetPremiumCalculationsRequest request,
        ServerCallContext context)
    {
        var calculations = await _vehicleService.GetPremiumCalculationsAsync(
            request.RegistrationId,
            context.CancellationToken);

        var calculationList = calculations.ToList();

        return new GetPremiumCalculationsResponse
        {
            Calculations = { calculationList },
            TotalCount = calculationList.Count
        };
    }

    public override async Task<GetVehicleImagesResponse> GetVehicleImages(
        GetVehicleImagesRequest request,
        ServerCallContext context)
    {
        var images = await _vehicleService.GetVehicleImagesAsync(context.CancellationToken);

        return new GetVehicleImagesResponse
        {
            Images = { images }
        };
    }
}
