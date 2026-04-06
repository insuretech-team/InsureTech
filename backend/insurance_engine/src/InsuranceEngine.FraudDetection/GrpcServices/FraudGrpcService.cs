using Grpc.Core;
using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Fraud.Services.V1;
using InsuranceEngine.FraudDetection.Application.Commands;
using InsuranceEngine.FraudDetection;
using System.Threading.Tasks;

namespace InsuranceEngine.FraudDetection.GrpcServices;

public sealed class FraudGrpcService : FraudService.FraudServiceBase
{
    private readonly IMediator _mediator;
    private readonly IFraudDetectionService _fraudService;
    private readonly ILogger<FraudGrpcService> _logger;

    public FraudGrpcService(
        IMediator mediator,
        IFraudDetectionService fraudService,
        ILogger<FraudGrpcService> logger)
    {
        _mediator = mediator;
        _fraudService = fraudService;
        _logger = logger;
    }

    public override async Task<CheckFraudResponse> CheckFraud(
        CheckFraudRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.EntityType) || string.IsNullOrEmpty(request.EntityId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Entity type and ID are required"));
        }

        var command = new CheckFraudCommand(
            request.EntityType,
            request.EntityId
        );

        return await _mediator.Send(command, context.CancellationToken);
    }
}
