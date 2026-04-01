using Grpc.Core;
using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Fraud.Services.V1;
using InsuranceEngine.FraudDetection.Application.Commands;
using System.Threading.Tasks;

namespace InsuranceEngine.FraudDetection.GrpcServices;

public sealed class FraudGrpcService : FraudService.FraudServiceBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<FraudGrpcService> _logger;

    public FraudGrpcService(IMediator mediator, ILogger<FraudGrpcService> logger)
    {
        _mediator = mediator;
        _logger = logger;
    }

    // RPC: CheckFraud
    public override async Task<CheckFraudResponse> CheckFraud(
        CheckFraudRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.EntityType) || string.IsNullOrEmpty(request.EntityId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Entity type and ID are required"));
        }

        var command = new CheckFraudCommand(
            request.EntityType,
            request.EntityId,
            request.Data
        );

        return await _mediator.Send(command, context.CancellationToken);
    }
}
