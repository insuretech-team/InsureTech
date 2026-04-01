using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Refund.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Messaging;

namespace PoliSync.Refund.Application.Commands;

public sealed class ApproveRefundCommandHandler : IRequestHandler<ApproveRefundCommand, Result>
{
    private readonly IEventBus _eventBus;
    private readonly ILogger<ApproveRefundCommandHandler> _logger;

    public ApproveRefundCommandHandler(
        IEventBus eventBus,
        ILogger<ApproveRefundCommandHandler> logger)
    {
        _eventBus = eventBus;
        _logger = logger;
    }

    public async Task<Result> Handle(ApproveRefundCommand request, CancellationToken cancellationToken)
    {
        try
        {
            if (string.IsNullOrWhiteSpace(request.RefundId))
                return Result.Fail("VALIDATION_ERROR", "RefundId is required");

            await _eventBus.PublishAsync(
                new RefundApprovedEvent(request.RefundId, request.ApprovedBy),
                cancellationToken);

            _logger.LogInformation("Refund {RefundId} approved by {ApprovedBy}", request.RefundId, request.ApprovedBy);
            return Result.Ok();
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to approve refund {RefundId}", request.RefundId);
            return Result.Fail("APPROVE_REFUND_FAILED", ex.Message);
        }
    }
}

public sealed record RefundApprovedEvent(string RefundId, string ApprovedBy)
    : PoliSync.SharedKernel.Domain.DomainEvent;

