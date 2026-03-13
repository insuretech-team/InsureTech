using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Underwriting.Application.DTOs;
using InsuranceEngine.Underwriting.Application.Interfaces;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.DTOs;
using MediatR;

namespace InsuranceEngine.Underwriting.Application.Features.Queries.GetQuote;

public class GetQuoteQueryHandler : IRequestHandler<GetQuoteQuery, Result<QuoteDto>>
{
    private readonly IUnderwritingRepository _repository;

    public GetQuoteQueryHandler(IUnderwritingRepository repository)
    {
        _repository = repository;
    }

    public async Task<Result<QuoteDto>> Handle(GetQuoteQuery request, CancellationToken cancellationToken)
    {
        var quote = await _repository.GetQuoteByIdAsync(request.QuoteId);
        if (quote == null)
            return Result.Fail<QuoteDto>(Error.NotFound("Quote", request.QuoteId.ToString()));

        return Result.Ok(new QuoteDto(
            quote.Id,
            quote.QuoteNumber,
            quote.BeneficiaryId,
            quote.InsurerProductId,
            quote.Status,
            new MoneyDto(quote.SumAssuredAmount, quote.SumAssuredCurrency),
            quote.TermYears,
            quote.PremiumPaymentMode,
            new MoneyDto(quote.BasePremiumAmount, quote.Currency),
            new MoneyDto(quote.RiderPremiumAmount, quote.Currency),
            new MoneyDto(quote.TotalPremiumAmount, quote.Currency),
            quote.ApplicantAge,
            quote.ApplicantOccupation,
            quote.IsSmoker,
            quote.ValidUntil,
            quote.CreatedAt
        ));
    }
}
