using Grpc.Core;
using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Claims.Services.V1;
using InsuranceEngine.Claims.Application.Queries;
using InsuranceEngine.Claims.Application.Commands;

namespace InsuranceEngine.Claims.GrpcServices;

public sealed class ClaimGrpcService : ClaimService.ClaimServiceBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<ClaimGrpcService> _logger;

    public ClaimGrpcService(IMediator mediator, ILogger<ClaimGrpcService> logger)
    {
        _mediator = mediator;
        _logger = logger;
    }

    public override async Task<GetClaimResponse> GetClaim(
        GetClaimRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.ClaimId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Claim ID is required"));
        }

        return await _mediator.Send(new GetClaimQuery(request.ClaimId), context.CancellationToken);
    }

    public override async Task<ListUserClaimsResponse> ListUserClaims(
        ListUserClaimsRequest request, ServerCallContext context)
    {
        var query = new ListUserClaimsQuery(
            PolicyId: null, // Depending on requirements
            Status: request.Status != Insuretech.Claims.Entity.V1.ClaimStatus.Unspecified ? request.Status.ToString() : null,
            Page: request.Page <= 0 ? 1 : request.Page,
            PageSize: request.PageSize <= 0 ? 10 : request.PageSize
        );

        return await _mediator.Send(query, context.CancellationToken);
    }

    public override async Task<SubmitClaimResponse> SubmitClaim(
        SubmitClaimRequest request, ServerCallContext context)
    {
        var command = new SubmitClaimCommand(
            PolicyId: request.PolicyId,
            ClaimType: request.Type.ToString(),
            ClaimAmount: (decimal)request.ClaimedAmount.Amount / 100m,
            Description: request.IncidentDescription,
            BeneficiaryId: null
        );

        return await _mediator.Send(command, context.CancellationToken);
    }

    public override async Task<ApproveClaimResponse> ApproveClaim(
        ApproveClaimRequest request, ServerCallContext context)
    {
        var command = new ApproveClaimCommand(
            ClaimId: request.ClaimId,
            ApprovedAmount: (decimal)request.ApprovedAmount.Amount / 100m
        );

        return await _mediator.Send(command, context.CancellationToken);
    }
}
