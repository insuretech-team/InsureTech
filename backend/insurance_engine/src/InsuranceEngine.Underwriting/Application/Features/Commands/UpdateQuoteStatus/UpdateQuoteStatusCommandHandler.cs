using System;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Underwriting.Application.Interfaces;
using InsuranceEngine.Underwriting.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Underwriting.Application.Features.Commands.UpdateQuoteStatus;

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
            return Result.Failure("Quote not found.");

        quote.Status = request.Status;
        quote.UpdatedAt = DateTime.UtcNow;

        await _repository.UpdateQuoteAsync(quote);
        return Result.Ok();
    }
}
