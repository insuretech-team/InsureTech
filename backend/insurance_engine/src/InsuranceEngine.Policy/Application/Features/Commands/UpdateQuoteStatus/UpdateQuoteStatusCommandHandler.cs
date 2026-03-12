using System;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Policy.Application.Features.Commands.UpdateQuoteStatus;

public class UpdateQuoteStatusCommandHandler : IRequestHandler<UpdateQuoteStatusCommand, Result>
{
    private readonly IUnderwritingRepository _repository;

    public UpdateQuoteStatusCommandHandler(IUnderwritingRepository repository)
    {
        _repository = repository;
    }

    public async Task<Result> Handle(UpdateQuoteStatusCommand request, CancellationToken cancellationToken)
    {
        var quote = await _repository.GetQuoteByIdAsync(request.QuoteId);
        if (quote == null)
            return Result.Fail(new Error("NOT_FOUND", $"Quote with ID {request.QuoteId} not found."));

        // Parse and validate status transition
        if (!Enum.TryParse<QuoteStatus>(request.NewStatus, ignoreCase: true, out var newStatus))
            return Result.Fail(new Error("VALIDATION_ERROR", $"Invalid status '{request.NewStatus}'. Valid values: Draft, PendingUnderwriting, Approved, Rejected, Expired, Converted."));

        // Validate allowed transitions
        var isValid = (quote.Status, newStatus) switch
        {
            (QuoteStatus.Draft, QuoteStatus.PendingUnderwriting) => true,
            (QuoteStatus.Draft, QuoteStatus.Expired) => true,
            (QuoteStatus.PendingUnderwriting, QuoteStatus.Approved) => true,
            (QuoteStatus.PendingUnderwriting, QuoteStatus.Rejected) => true,
            (QuoteStatus.PendingUnderwriting, QuoteStatus.Expired) => true,
            (QuoteStatus.Approved, QuoteStatus.Converted) => true,
            (QuoteStatus.Approved, QuoteStatus.Expired) => true,
            _ => false
        };

        if (!isValid)
            return Result.Fail(new Error("INVALID_STATE_TRANSITION", $"Cannot transition from {quote.Status} to {newStatus}."));

        quote.Status = newStatus;
        quote.UpdatedAt = DateTime.UtcNow;

        if (newStatus == QuoteStatus.Converted && request.ConvertedPolicyId.HasValue)
        {
            quote.ConvertedPolicyId = request.ConvertedPolicyId.Value;
            quote.ConvertedAt = DateTime.UtcNow;
        }

        await _repository.UpdateQuoteAsync(quote);
        return Result.Ok();
    }
}
