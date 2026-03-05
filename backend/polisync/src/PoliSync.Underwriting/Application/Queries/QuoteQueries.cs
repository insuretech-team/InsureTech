using PoliSync.SharedKernel.CQRS;
using PoliSync.Underwriting.Domain;

namespace PoliSync.Underwriting.Application.Queries;

// ── Get Quote ───────────────────────────────────────────────────────

public record GetQuoteQuery(Guid QuoteId) : IQuery<Quote?>;

public class GetQuoteHandler : IQueryHandler<GetQuoteQuery, Quote?>
{
    private readonly IQuoteRepository _repo;

    public GetQuoteHandler(IQuoteRepository repo) => _repo = repo;

    public async Task<Result<Quote?>> Handle(GetQuoteQuery q, CancellationToken ct)
    {
        var quote = await _repo.GetByIdAsync(q.QuoteId, ct);
        return Result<Quote?>.Ok(quote);
    }
}

// ── List Quotes ─────────────────────────────────────────────────────

public record ListQuotesQuery(
    Guid? BeneficiaryId = null,
    QuoteStatus? Status = null,
    int Page = 1,
    int PageSize = 20
) : IQuery<QuoteListResult>;

public record QuoteListResult(IEnumerable<Quote> Items, int TotalCount, int Page, int PageSize);

public class ListQuotesHandler : IQueryHandler<ListQuotesQuery, QuoteListResult>
{
    private readonly IQuoteRepository _repo;

    public ListQuotesHandler(IQuoteRepository repo) => _repo = repo;

    public async Task<Result<QuoteListResult>> Handle(ListQuotesQuery q, CancellationToken ct)
    {
        var (items, total) = await _repo.ListAsync(q.BeneficiaryId, q.Status, q.Page, q.PageSize, ct);
        return Result<QuoteListResult>.Ok(new QuoteListResult(items, total, q.Page, q.PageSize));
    }
}

// ── Get Health Declaration ──────────────────────────────────────────

public record GetHealthDeclarationQuery(Guid QuoteId) : IQuery<HealthDeclaration?>;

public class GetHealthDeclarationHandler : IQueryHandler<GetHealthDeclarationQuery, HealthDeclaration?>
{
    private readonly IQuoteRepository _repo;

    public GetHealthDeclarationHandler(IQuoteRepository repo) => _repo = repo;

    public async Task<Result<HealthDeclaration?>> Handle(GetHealthDeclarationQuery q, CancellationToken ct)
    {
        var hd = await _repo.GetHealthDeclarationByQuoteIdAsync(q.QuoteId, ct);
        return Result<HealthDeclaration?>.Ok(hd);
    }
}

// ── Get Underwriting Decision ───────────────────────────────────────

public record GetUnderwritingDecisionQuery(Guid QuoteId) : IQuery<UnderwritingDecision?>;

public class GetUnderwritingDecisionHandler : IQueryHandler<GetUnderwritingDecisionQuery, UnderwritingDecision?>
{
    private readonly IQuoteRepository _repo;

    public GetUnderwritingDecisionHandler(IQuoteRepository repo) => _repo = repo;

    public async Task<Result<UnderwritingDecision?>> Handle(GetUnderwritingDecisionQuery q, CancellationToken ct)
    {
        var decision = await _repo.GetDecisionByQuoteIdAsync(q.QuoteId, ct);
        return Result<UnderwritingDecision?>.Ok(decision);
    }
}
