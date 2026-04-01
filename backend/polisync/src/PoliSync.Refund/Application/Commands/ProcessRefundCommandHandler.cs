using Insuretech.Common.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Refund.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Messaging;

namespace PoliSync.Refund.Application.Commands;

public sealed class ProcessRefundCommandHandler : IRequestHandler<ProcessRefundCommand, Result<ProcessRefundResult>>
{
    private readonly IRefundPaymentGateway _paymentGateway;
    private readonly IEventBus _eventBus;
    private readonly ILogger<ProcessRefundCommandHandler> _logger;

    public ProcessRefundCommandHandler(
        IRefundPaymentGateway paymentGateway,
        IEventBus eventBus,
        ILogger<ProcessRefundCommandHandler> logger)
    {
        _paymentGateway = paymentGateway;
        _eventBus = eventBus;
        _logger = logger;
    }

    public async Task<Result<ProcessRefundResult>> Handle(ProcessRefundCommand request, CancellationToken cancellationToken)
    {
        try
        {
            if (string.IsNullOrWhiteSpace(request.PaymentReference))
                return Result.Fail<ProcessRefundResult>("VALIDATION_ERROR", "PaymentReference is required");

            var refundAmount = new Money
            {
                Amount = request.RefundAmountPaisa,
                Currency = "BDT"
            };

            var response = await _paymentGateway.InitiateRefundAsync(
                request.PaymentReference,
                refundAmount,
                request.Reason,
                request.InitiatedBy,
                cancellationToken);

            if (response.Error is not null && !string.IsNullOrWhiteSpace(response.Error.Code))
                return Result.Fail<ProcessRefundResult>(response.Error.Code, response.Error.Message);

            var paymentRefundId = string.IsNullOrWhiteSpace(response.RefundId)
                ? $"PAYREF-{Guid.NewGuid():N}"[..18]
                : response.RefundId;

            var processedAt = DateTime.UtcNow.ToString("O");

            await _eventBus.PublishAsync(
                new RefundProcessedEvent(request.RefundId, paymentRefundId, request.RefundAmountPaisa),
                cancellationToken);

            _logger.LogInformation("Refund {RefundId} processed. PaymentRefundId: {PaymentRefundId}",
                request.RefundId, paymentRefundId);

            return Result.Ok(new ProcessRefundResult(paymentRefundId, processedAt));
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to process refund {RefundId}", request.RefundId);
            return Result.Fail<ProcessRefundResult>("PROCESS_REFUND_FAILED", ex.Message);
        }
    }
}

public sealed record RefundProcessedEvent(string RefundId, string PaymentRefundId, long AmountPaisa)
    : PoliSync.SharedKernel.Domain.DomainEvent;

