using Grpc.Core;
using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Renewal.Services.V1;
using InsuranceEngine.Renewals.Application.Commands;
using System.Threading.Tasks;

namespace InsuranceEngine.Renewals.GrpcServices;

public sealed class RenewalGrpcService : RenewalService.RenewalServiceBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<RenewalGrpcService> _logger;

    public RenewalGrpcService(IMediator mediator, ILogger<RenewalGrpcService> logger)
    {
        _mediator = mediator;
        _logger = logger;
    }

    // RPC: RenewPolicy
    public override async Task<RenewPolicyResponse> RenewPolicy(
        RenewPolicyRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.PolicyId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Policy ID is required"));
        }

        // Note: The logic in RenewPolicyCommandHandler currently expects TenureMonths.
        // The RenewalService.proto RenewPolicyRequest has PolicyId, PaymentMethod, PaymentReference.
        // We might need to adjust the command or the handler to match the specialized service.
        // For now, mapping to the existing handler logic (assuming 12 months as default if not in proto).
        
        var command = new RenewPolicyCommand(
            PolicyId: request.PolicyId,
            TenureMonths: 12, // Defaulting for now as the proto didn't have it
            UpdateNominees: false,
            Nominees: null
        );

        var result = await _mediator.Send(command, context.CancellationToken);

        return new RenewPolicyResponse
        {
            NewPolicyId = result.NewPolicyId,
            Message = result.Message,
            Error = result.Error
        };
    }

    // Other RPCs (GetRenewalSchedule, ListUpcomingRenewals, etc.) would go here
    // For P0, we are just isolating the existing logic.
}
