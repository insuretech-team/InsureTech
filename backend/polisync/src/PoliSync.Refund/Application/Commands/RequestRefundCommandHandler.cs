using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Refund.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Messaging;
using PoliSync.Workflow.Application.Commands;
using PoliSync.Workflow.Domain;

namespace PoliSync.Refund.Application.Commands;

public sealed class RequestRefundCommandHandler : IRequestHandler<RequestRefundCommand, Result<RequestRefundResult>>
{
    private readonly IRefundPaymentGateway _paymentGateway;
    private readonly IEventBus _eventBus;
    private readonly IMediator _mediator;
    private readonly ILogger<RequestRefundCommandHandler> _logger;

    public RequestRefundCommandHandler(
        IRefundPaymentGateway paymentGateway,
        IEventBus eventBus,
        IMediator mediator,
        ILogger<RequestRefundCommandHandler> logger)
    {
        _paymentGateway = paymentGateway;
        _eventBus = eventBus;
        _mediator = mediator;
        _logger = logger;
    }

    public async Task<Result<RequestRefundResult>> Handle(RequestRefundCommand request, CancellationToken cancellationToken)
    {
        try
        {
            if (string.IsNullOrWhiteSpace(request.PolicyId))
                return Result.Fail<RequestRefundResult>("VALIDATION_ERROR", "PolicyId is required");

            // Refund request is recorded in the in-memory store via the GrpcService
            // This handler captures the domain intent and publishes the event
            var refundId = Guid.NewGuid().ToString("N");
            var refundNumber = $"RFD-{DateTime.UtcNow:yyyyMMdd}-{Random.Shared.Next(100000, 999999)}";

            await _eventBus.PublishAsync(
                new RefundRequestedEvent(refundId, request.PolicyId, request.Reason),
                cancellationToken);

            _logger.LogInformation("Refund requested {RefundId} for policy {PolicyId}", refundId, request.PolicyId);

            // Trigger approval workflow — high vs standard resolved by amount
            var workflowResult = await _mediator.Send(new TriggerWorkflowCommand(
                new WorkflowTriggerContext
                {
                    EntityType  = "REFUND",
                    EntityId    = refundId,
                    InitiatedBy = "SYSTEM",
                    AmountPaisa = request.AmountPaisa,
                    Portal      = "B2C",
                    Metadata    = new Dictionary<string, string>
                    {
                        ["policy_id"]     = request.PolicyId,
                        ["refund_number"] = refundNumber,
                        ["reason"]        = request.Reason
                    }
                }), cancellationToken);

            if (workflowResult.IsSuccess && workflowResult.Value!.WasTriggered)
                _logger.LogInformation(
                    "Refund approval workflow started: instance={InstanceId} template='{Template}'",
                    workflowResult.Value.WorkflowInstanceId, workflowResult.Value.TemplateName);

            return Result.Ok(new RequestRefundResult(refundId, refundNumber));
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to request refund for policy {PolicyId}", request.PolicyId);
            return Result.Fail<RequestRefundResult>("REQUEST_REFUND_FAILED", ex.Message);
        }
    }
}

public sealed record RefundRequestedEvent(string RefundId, string PolicyId, string Reason)
    : PoliSync.SharedKernel.Domain.DomainEvent;

