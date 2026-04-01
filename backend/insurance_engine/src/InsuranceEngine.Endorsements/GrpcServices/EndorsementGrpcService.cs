using Grpc.Core;
using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Endorsement.Services.V1;
using InsuranceEngine.Endorsements.Application.Commands;
using System.Threading.Tasks;

namespace InsuranceEngine.Endorsements.GrpcServices;

public sealed class EndorsementGrpcService : EndorsementService.EndorsementServiceBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<EndorsementGrpcService> _logger;

    public EndorsementGrpcService(IMediator mediator, ILogger<EndorsementGrpcService> logger)
    {
        _mediator = mediator;
        _logger = logger;
    }

    // RPC: RequestEndorsement
    public override async Task<RequestEndorsementResponse> RequestEndorsement(
        RequestEndorsementRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.PolicyId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Policy ID is required"));
        }

        // Mapping to UpdatePolicyCommand for current extraction task
        // We might want to expand the command to handle more endorsement-specific fields (type, reason, changes)
        var command = new UpdatePolicyCommand(
            PolicyId: request.PolicyId,
            Nominees: null, // Endorsement-specific nominees logic can be added later
            Address: null
        );

        var result = await _mediator.Send(command, context.CancellationToken);

        return new RequestEndorsementResponse
        {
            EndorsementId = request.PolicyId, // Reusing ID for now
            EndorsementNumber = "PENDING",
            Message = result.Message
        };
    }

    // Other RPCs (GetEndorsement, Approval, Rejection, etc.)
}
