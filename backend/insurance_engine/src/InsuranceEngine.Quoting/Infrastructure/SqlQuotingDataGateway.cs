using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Insuretech.Quoting.Entity.V1;
using Insuretech.Quoting.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence.Entities;

namespace InsuranceEngine.Quoting.Infrastructure;

public class SqlQuotingDataGateway : IQuotingDataGateway
{
    private readonly QuotingDbContext _context;
    private readonly ILogger<SqlQuotingDataGateway> _logger;

    public SqlQuotingDataGateway(QuotingDbContext context, ILogger<SqlQuotingDataGateway> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<Quote> GenerateQuoteAsync(GenerateQuoteRequest request, CancellationToken cancellationToken = default)
    {
        var quoteId = Guid.NewGuid();
        var quoteNumber = $"QT-{DateTime.UtcNow.Year}-{DateTime.UtcNow:MMdd}-{Guid.NewGuid().ToString()[..8].ToUpper()}";
        var now = DateTime.UtcNow;

        var sumAssured = 1000000L;
        var basePremium = CalculateBasePremium(sumAssured, 1, 30);
        var totalPremium = basePremium;

        var quote = new QuoteEntity
        {
            QuoteId = quoteId,
            QuoteNumber = quoteNumber,
            BeneficiaryId = Guid.Empty,
            InsurerProductId = Guid.Empty,
            Status = "DRAFT",
            SumAssured = sumAssured,
            SumAssuredCurrency = "BDT",
            TermYears = 1,
            BasePremium = basePremium,
            BasePremiumCurrency = "BDT",
            TotalPremium = totalPremium,
            TotalPremiumCurrency = "BDT",
            ApplicantAge = 30,
            ValidUntil = now.AddDays(30),
            CreatedAt = now,
            UpdatedAt = now
        };

        _context.Quotes.Add(quote);
        await _context.SaveChangesAsync(cancellationToken);

        _logger.LogInformation("SQL: Generated quote {QuoteNumber}", quoteNumber);

        return MapToProto(quote);
    }

    public async Task<Quote> GetQuoteAsync(string quoteId, CancellationToken cancellationToken = default)
    {
        var id = Guid.TryParse(quoteId, out var qid) ? qid : Guid.Empty;
        var quote = await _context.Quotes.FindAsync([id], cancellationToken);

        if (quote == null)
        {
            throw new Exception("Quote not found");
        }

        return MapToProto(quote);
    }

    public async Task<Quote> GetQuoteByNumberAsync(string quoteNumber, CancellationToken cancellationToken = default)
    {
        var quote = await _context.Quotes.FirstOrDefaultAsync(q => q.QuoteNumber == quoteNumber, cancellationToken);

        if (quote == null)
        {
            throw new Exception("Quote not found");
        }

        return MapToProto(quote);
    }

    public async Task<ListQuotesResponse> ListQuotesAsync(ListQuotesRequest request, CancellationToken cancellationToken = default)
    {
        var query = _context.Quotes.AsQueryable();

        var page = 1;
        var pageSize = 10;

        var quotes = await query
            .OrderByDescending(q => q.CreatedAt)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);

        var response = new ListQuotesResponse();
        response.Quotes.AddRange(quotes.Select(MapToProto));

        return response;
    }

    public async Task<Quote> ReviseQuoteAsync(ReviseQuoteRequest request, CancellationToken cancellationToken = default)
    {
        var id = Guid.TryParse(request.QuoteId, out var qid) ? qid : Guid.Empty;
        var quote = await _context.Quotes.FindAsync([id], cancellationToken);

        if (quote == null)
        {
            throw new Exception("Quote not found");
        }

        quote.BasePremium = CalculateBasePremium(quote.SumAssured, quote.TermYears, quote.ApplicantAge);
        quote.TotalPremium = quote.BasePremium;
        quote.UpdatedAt = DateTime.UtcNow;

        await _context.SaveChangesAsync(cancellationToken);

        _logger.LogInformation("SQL: Revised quote {QuoteId}", request.QuoteId);

        return MapToProto(quote);
    }

    public async Task<ConvertQuoteToPolicyResponse> ConvertQuoteToPolicyAsync(ConvertQuoteToPolicyRequest request, CancellationToken cancellationToken = default)
    {
        var id = Guid.TryParse(request.QuoteId, out var qid) ? qid : Guid.Empty;
        var quote = await _context.Quotes.FindAsync([id], cancellationToken);

        if (quote == null)
        {
            return new ConvertQuoteToPolicyResponse
            {
                Error = new Error { Code = "NOT_FOUND", Message = "Quote not found" }
            };
        }

        var policyId = Guid.NewGuid();

        quote.Status = "CONVERTED";
        quote.ConvertedPolicyId = policyId;
        quote.ConvertedAt = DateTime.UtcNow;
        quote.UpdatedAt = DateTime.UtcNow;

        await _context.SaveChangesAsync(cancellationToken);

        _logger.LogInformation("SQL: Converted quote {QuoteId} to policy {PolicyId}", request.QuoteId, policyId);

        return new ConvertQuoteToPolicyResponse
        {
            PolicyId = policyId.ToString()
        };
    }

    private static long CalculateBasePremium(long sumAssured, int termYears, int age)
    {
        var baseRate = 0.02;
        var ageMultiplier = 1 + (age - 30) * 0.01;
        var termMultiplier = 1 + (termYears - 1) * 0.05;
        return (long)(sumAssured * baseRate * ageMultiplier * termMultiplier / 100);
    }

    private static Quote MapToProto(QuoteEntity entity)
    {
        return new Quote
        {
            QuoteId = entity.QuoteId.ToString(),
            QuoteNumber = entity.QuoteNumber,
            Status = Enum.TryParse<QuoteStatus>(entity.Status, true, out var s) ? s : QuoteStatus.Draft,
            TotalPremium = new Money { Amount = entity.TotalPremium, Currency = entity.TotalPremiumCurrency },
            ValidUntil = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(entity.ValidUntil),
            ConvertedPolicyId = entity.ConvertedPolicyId?.ToString() ?? "",
            CreatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(entity.CreatedAt)
        };
    }
}
