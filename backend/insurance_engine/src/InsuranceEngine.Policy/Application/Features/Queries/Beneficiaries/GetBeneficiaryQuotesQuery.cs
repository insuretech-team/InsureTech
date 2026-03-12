using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.Policy.Infrastructure.Persistence;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Policy.Application.Features.Queries.Beneficiaries;

public record GetBeneficiaryQuotesQuery(Guid BeneficiaryId) : IRequest<Result<List<QuoteDto>>>;

public class GetBeneficiaryQuotesQueryHandler : IRequestHandler<GetBeneficiaryQuotesQuery, Result<List<QuoteDto>>>
{
    private readonly PolicyDbContext _context;

    public GetBeneficiaryQuotesQueryHandler(PolicyDbContext context)
    {
        _context = context;
    }

    public async Task<Result<List<QuoteDto>>> Handle(GetBeneficiaryQuotesQuery request, CancellationToken cancellationToken)
    {
        // Actually, quotes are usually linked to a customer/beneficiary
        // We need to see if the Quote entity has a BeneficiaryId or CustomerId.
        // Let's assume it links via CustomerId for now, or we add BeneficiaryId to Quote.
        
        var quotes = await _context.Quotes
            .Where(q => q.BeneficiaryId == request.BeneficiaryId)
            .OrderByDescending(q => q.CreatedAt)
            .ToListAsync(cancellationToken);

        return Result.Ok(quotes.Select(q => new QuoteDto(
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
        )).ToList());
    }
}
