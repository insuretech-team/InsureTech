using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Insuretech.Commission.Services.V1;
using Insuretech.Common.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Commission.Application.Commands;
using PoliSync.Commission.Infrastructure;

namespace PoliSync.Commission.GrpcServices;

public sealed class CommissionGrpcService : CommissionService.CommissionServiceBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<CommissionGrpcService> _logger;
    private readonly ICommissionDataGateway _dataGateway;

    public CommissionGrpcService(
        IMediator mediator,
        ILogger<CommissionGrpcService> logger,
        ICommissionDataGateway dataGateway)
    {
        _mediator = mediator;
        _logger = logger;
        _dataGateway = dataGateway;
    }

    public override async Task<CalculateCommissionResponse> CalculateCommission(
        CalculateCommissionRequest request, ServerCallContext context)
    {
        if (string.IsNullOrWhiteSpace(request.PolicyId) || string.IsNullOrWhiteSpace(request.RecipientId))
        {
            return new CalculateCommissionResponse
            {
                Error = BuildError("VALIDATION_ERROR", "PolicyId and RecipientId are required")
            };
        }

        try
        {
            var result = await _mediator.Send(new CalculateCommissionCommand(
                request.PolicyId,
                request.CommissionType,
                request.RecipientType,
                request.RecipientId), context.CancellationToken);

            if (result.IsFailure)
                return new CalculateCommissionResponse { Error = BuildError(result.Error!.Code, result.Error.Message) };

            return new CalculateCommissionResponse
            {
                CommissionId = result.Value!.CommissionId,
                CommissionNumber = result.Value.CommissionNumber,
                Amount = result.Value.Amount,
                CalculationBreakdown = result.Value.CalculationBreakdown
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to calculate commission for policy {PolicyId}", request.PolicyId);
            return new CalculateCommissionResponse { Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail) };
        }
    }

    public override async Task<GetCommissionResponse> GetCommission(
        GetCommissionRequest request, ServerCallContext context)
    {
        try
        {
            var commission = await _dataGateway.GetCommissionAsync(request.CommissionId, context.CancellationToken);
            if (commission is null)
                return new GetCommissionResponse { Error = BuildError("NOT_FOUND", "Commission not found") };

            return new GetCommissionResponse { Commission = commission };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to get commission {CommissionId}", request.CommissionId);
            return new GetCommissionResponse { Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail) };
        }
    }

    public override async Task<ListCommissionsResponse> ListCommissions(
        ListCommissionsRequest request, ServerCallContext context)
    {
        try
        {
            var page = request.Page <= 0 ? 1 : request.Page;
            var pageSize = request.PageSize <= 0 ? 20 : request.PageSize;

            var (items, totalCount, totalAmount) = await _dataGateway.ListCommissionsAsync(
                request.RecipientType, request.RecipientId, request.Status,
                request.StartDate, request.EndDate, page, pageSize,
                context.CancellationToken);

            var response = new ListCommissionsResponse
            {
                TotalCount = totalCount,
                TotalAmount = new Money { Amount = totalAmount, Currency = "BDT" }
            };
            response.Commissions.AddRange(items);
            return response;
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to list commissions for recipient {RecipientId}", request.RecipientId);
            return new ListCommissionsResponse { Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail) };
        }
    }

    public override async Task<CreatePayoutResponse> CreatePayout(
        CreatePayoutRequest request, ServerCallContext context)
    {
        try
        {
            var result = await _mediator.Send(new CreateCommissionPayoutCommand(
                request.RecipientType,
                request.RecipientId,
                request.PeriodStart,
                request.PeriodEnd,
                request.CommissionIds.ToList()), context.CancellationToken);

            if (result.IsFailure)
                return new CreatePayoutResponse { Error = BuildError(result.Error!.Code, result.Error.Message) };

            return new CreatePayoutResponse
            {
                PayoutId = result.Value!.PayoutId,
                PayoutNumber = result.Value.PayoutNumber,
                TotalAmount = result.Value.TotalAmount,
                CommissionCount = result.Value.CommissionCount
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to create payout for recipient {RecipientId}", request.RecipientId);
            return new CreatePayoutResponse { Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail) };
        }
    }

    public override async Task<ProcessPayoutResponse> ProcessPayout(
        ProcessPayoutRequest request, ServerCallContext context)
    {
        try
        {
            var result = await _mediator.Send(new ProcessCommissionPayoutCommand(
                request.PayoutId,
                request.PaymentMethod,
                request.PaymentReference), context.CancellationToken);

            if (result.IsFailure)
                return new ProcessPayoutResponse { Error = BuildError(result.Error!.Code, result.Error.Message) };

            return new ProcessPayoutResponse
            {
                Message = "Payout processed",
                PaidAt = result.Value!
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to process payout {PayoutId}", request.PayoutId);
            return new ProcessPayoutResponse { Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail) };
        }
    }

    public override async Task<GetCommissionStatementResponse> GetCommissionStatement(
        GetCommissionStatementRequest request, ServerCallContext context)
    {
        try
        {
            return await _dataGateway.GetCommissionStatementAsync(
                request.RecipientId, request.PeriodStart, request.PeriodEnd,
                context.CancellationToken);
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to get commission statement for {RecipientId}", request.RecipientId);
            return new GetCommissionStatementResponse { Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail) };
        }
    }

    public override async Task<GetRevenueShareReportResponse> GetRevenueShareReport(
        GetRevenueShareReportRequest request, ServerCallContext context)
    {
        try
        {
            return await _dataGateway.GetRevenueShareReportAsync(
                request.InsurerId, request.StartDate, request.EndDate,
                context.CancellationToken);
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to get revenue share report for insurer {InsurerId}", request.InsurerId);
            return new GetRevenueShareReportResponse { Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail) };
        }
    }

    private static Error BuildError(string code, string message) => new() { Code = code, Message = message };
}
