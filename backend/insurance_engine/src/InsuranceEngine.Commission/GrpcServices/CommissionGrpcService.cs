using Grpc.Core;
using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Commission.Services.V1;
using InsuranceEngine.Commission.Application.Commands;
using InsuranceEngine.Commission.Application.Queries;

namespace InsuranceEngine.Commission.GrpcServices;

public sealed class CommissionGrpcService : CommissionService.CommissionServiceBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<CommissionGrpcService> _logger;

    public CommissionGrpcService(IMediator mediator, ILogger<CommissionGrpcService> logger)
    {
        _mediator = mediator;
        _logger = logger;
    }

    public override async Task<CalculateCommissionResponse> CalculateCommission(CalculateCommissionRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.PolicyId) || string.IsNullOrEmpty(request.RecipientId))
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Policy ID and Recipient ID are required"));
        return await _mediator.Send(new CalculateCommissionCommand(
            request.PolicyId, request.CommissionType, request.RecipientType, request.RecipientId), context.CancellationToken);
    }

    public override async Task<GetCommissionResponse> GetCommission(GetCommissionRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.CommissionId))
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Commission ID is required"));
        return await _mediator.Send(new GetCommissionQuery(request.CommissionId), context.CancellationToken);
    }

    public override async Task<ListCommissionsResponse> ListCommissions(ListCommissionsRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.RecipientId))
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Recipient ID is required"));
        return await _mediator.Send(new ListCommissionsQuery(
            request.RecipientType, request.RecipientId, request.Status,
            request.StartDate, request.EndDate,
            request.Page <= 0 ? 1 : request.Page,
            request.PageSize <= 0 ? 10 : request.PageSize), context.CancellationToken);
    }

    public override async Task<CreatePayoutResponse> CreatePayout(CreatePayoutRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.RecipientId))
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Recipient ID is required"));
        return await _mediator.Send(new CreatePayoutCommand(
            request.RecipientType, request.RecipientId,
            request.PeriodStart, request.PeriodEnd,
            request.CommissionIds.Count > 0 ? request.CommissionIds.ToList() : null), context.CancellationToken);
    }

    public override async Task<ProcessPayoutResponse> ProcessPayout(ProcessPayoutRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.PayoutId))
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Payout ID is required"));
        return await _mediator.Send(new ProcessPayoutCommand(
            request.PayoutId, request.PaymentMethod, request.PaymentReference), context.CancellationToken);
    }

    public override async Task<GetCommissionStatementResponse> GetCommissionStatement(GetCommissionStatementRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.RecipientId))
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Recipient ID is required"));
        return await _mediator.Send(new GetCommissionStatementQuery(
            request.RecipientId, request.PeriodStart, request.PeriodEnd), context.CancellationToken);
    }

    public override async Task<GetRevenueShareReportResponse> GetRevenueShareReport(GetRevenueShareReportRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.InsurerId))
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Insurer ID is required"));
        return await _mediator.Send(new GetRevenueShareReportQuery(
            request.InsurerId, request.StartDate, request.EndDate), context.CancellationToken);
    }
}
