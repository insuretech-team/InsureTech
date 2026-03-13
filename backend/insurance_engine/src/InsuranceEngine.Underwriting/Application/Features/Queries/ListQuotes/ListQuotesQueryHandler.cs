using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Underwriting.Application.DTOs;
using InsuranceEngine.Underwriting.Application.Interfaces;
using InsuranceEngine.SharedKernel.DTOs;
using MediatR;

namespace InsuranceEngine.Underwriting.Application.Features.Queries.ListQuotes;

public class ListQuotesQueryHandler : IRequestHandler<ListQuotesQuery, PaginatedResponse<QuoteDto>>
{
    private readonly IUnderwritingRepository _repository;

    public ListQuotesQueryHandler(IUnderwritingRepository repository)
    {
        _repository = repository;
    }

    public async Task<PaginatedResponse<QuoteDto>> Handle(ListQuotesQuery request, CancellationToken cancellationToken)
    {
        var (items, totalCount) = await _repository.ListQuotesAsync(
            request.BeneficiaryId, request.Status, request.Page, request.PageSize);

        var dtos = items.Select(q => new QuoteDto(
            q.Id,
            q.QuoteNumber,
            q.BeneficiaryId,
            q.InsurerProductId,
            q.Status,
            new MoneyDto(q.SumAssuredAmount, q.SumAssuredCurrency),
            q.TermYears,
            q.PremiumPaymentMode,
            new MoneyDto(q.BasePremiumAmount, q.Currency),
            new MoneyDto(q.RiderPremiumAmount, q.Currency),
            new MoneyDto(q.TotalPremiumAmount, q.Currency),
            q.ApplicantAge,
            q.ApplicantOccupation,
            q.IsSmoker,
            q.ValidUntil,
            q.CreatedAt
        )).ToList();

        return new PaginatedResponse<QuoteDto>(dtos, totalCount, request.Page, request.PageSize);
    }
}
