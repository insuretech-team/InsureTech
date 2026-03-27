using Grpc.Core;
using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Policy.Services.V1;
using InsuranceEngine.Policy.Application.Queries;
using InsuranceEngine.Policy.Application.Commands;

namespace InsuranceEngine.Policy.GrpcServices;

public sealed class PolicyGrpcService : PolicyService.PolicyServiceBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<PolicyGrpcService> _logger;

    public PolicyGrpcService(IMediator mediator, ILogger<PolicyGrpcService> logger)
    {
        _mediator = mediator;
        _logger = logger;
    }

    public override async Task<GetPolicyResponse> GetPolicy(
        GetPolicyRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.PolicyId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Policy ID is required"));
        }

        return await _mediator.Send(new GetPolicyQuery(request.PolicyId), context.CancellationToken);
    }

    public override async Task<ListUserPoliciesResponse> ListUserPolicies(
        ListUserPoliciesRequest request, ServerCallContext context)
    {
        var query = new ListUserPoliciesQuery(
            CustomerId: request.CustomerId,
            Status: request.Status != Insuretech.Policy.Entity.V1.PolicyStatus.Unspecified ? request.Status.ToString() : null,
            ProductId: null,
            Page: request.Page <= 0 ? 1 : request.Page,
            PageSize: request.PageSize <= 0 ? 10 : request.PageSize
        );

        return await _mediator.Send(query, context.CancellationToken);
    }

    public override async Task<CreatePolicyResponse> CreatePolicy(
        CreatePolicyRequest request, ServerCallContext context)
    {
        var command = new CreatePolicyCommand(
            ProductId: request.ProductId,
            CustomerId: request.CustomerId,
            PartnerId: request.PartnerId,
            AgentId: request.AgentId,
            QuoteId: null,
            PremiumAmount: (decimal)(request.PremiumAmount?.Amount ?? 0) / 100m,
            SumInsured: (decimal)(request.SumInsured?.Amount ?? 0) / 100m,
            TenureMonths: request.TenureMonths,
            StartDate: DateTime.UtcNow,
            ProposerDetails: null,
            Nominees: request.Nominees.ToList()
        );

        return await _mediator.Send(command, context.CancellationToken);
    }

    public override async Task<IssuePolicyResponse> IssuePolicy(
        IssuePolicyRequest request, ServerCallContext context)
    {
        return await _mediator.Send(new IssuePolicyCommand(request.PolicyId), context.CancellationToken);
    }
}
