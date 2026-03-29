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

    // RPC 1: CreatePolicy
    public override async Task<CreatePolicyResponse> CreatePolicy(
        CreatePolicyRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.ProductId) || string.IsNullOrEmpty(request.CustomerId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Product ID and Customer ID are required"));
        }

        var command = new CreatePolicyCommand(
            ProductId: request.ProductId,
            CustomerId: request.CustomerId,
            PartnerId: request.PartnerId,
            AgentId: request.AgentId,
            QuoteId: null,
            PremiumAmount: request.PremiumAmount != null ? (decimal)request.PremiumAmount.Amount / 100m : 0m,
            SumInsured: request.SumInsured != null ? (decimal)request.SumInsured.Amount / 100m : 0m,
            TenureMonths: request.TenureMonths,
            StartDate: DateTime.UtcNow,
            ProposerDetails: null,
            Nominees: request.Nominees.ToList()
        );

        return await _mediator.Send(command, context.CancellationToken);
    }

    // RPC 2: GetPolicy
    public override async Task<GetPolicyResponse> GetPolicy(
        GetPolicyRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.PolicyId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Policy ID is required"));
        }

        return await _mediator.Send(new GetPolicyQuery(request.PolicyId), context.CancellationToken);
    }

    // RPC 3: ListUserPolicies
    public override async Task<ListUserPoliciesResponse> ListUserPolicies(
        ListUserPoliciesRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.CustomerId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Customer ID is required"));
        }

        var query = new ListUserPoliciesQuery(
            CustomerId: request.CustomerId,
            Status: request.Status != Insuretech.Policy.Entity.V1.PolicyStatus.Unspecified ? request.Status.ToString() : null,
            ProductId: null,
            Page: request.Page <= 0 ? 1 : request.Page,
            PageSize: request.PageSize <= 0 ? 10 : request.PageSize
        );

        return await _mediator.Send(query, context.CancellationToken);
    }

    // RPC 4: UpdatePolicy
    public override async Task<UpdatePolicyResponse> UpdatePolicy(
        UpdatePolicyRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.PolicyId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Policy ID is required"));
        }

        var command = new UpdatePolicyCommand(
            PolicyId: request.PolicyId,
            Nominees: request.Nominees.Count > 0 ? request.Nominees.ToList() : null,
            Address: request.Address
        );

        return await _mediator.Send(command, context.CancellationToken);
    }

    // RPC 5: CancelPolicy
    public override async Task<CancelPolicyResponse> CancelPolicy(
        CancelPolicyRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.PolicyId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Policy ID is required"));
        }
        if (string.IsNullOrEmpty(request.Reason))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Cancellation reason is required"));
        }

        return await _mediator.Send(new CancelPolicyCommand(request.PolicyId, request.Reason), context.CancellationToken);
    }

    // RPC 6: RenewPolicy
    public override async Task<RenewPolicyResponse> RenewPolicy(
        RenewPolicyRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.PolicyId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Policy ID is required"));
        }
        if (request.TenureMonths <= 0)
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Tenure months must be positive"));
        }

        var command = new RenewPolicyCommand(
            PolicyId: request.PolicyId,
            TenureMonths: request.TenureMonths,
            UpdateNominees: request.UpdateNominees,
            Nominees: request.Nominees.Count > 0 ? request.Nominees.ToList() : null
        );

        return await _mediator.Send(command, context.CancellationToken);
    }

    // RPC 7: GeneratePolicyDocument
    public override async Task<GeneratePolicyDocumentResponse> GeneratePolicyDocument(
        GeneratePolicyDocumentRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.PolicyId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Policy ID is required"));
        }

        return await _mediator.Send(new GeneratePolicyDocumentCommand(request.PolicyId), context.CancellationToken);
    }

    // RPC 8: IssuePolicy
    public override async Task<IssuePolicyResponse> IssuePolicy(
        IssuePolicyRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.PolicyId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Policy ID is required"));
        }

        return await _mediator.Send(
            new IssuePolicyCommand(request.PolicyId, request.QuoteId, request.PaymentId), 
            context.CancellationToken);
    }
}
