using Microsoft.EntityFrameworkCore;
using PoliSync.Infrastructure.Persistence;
using PoliSync.Underwriting.Domain;

namespace PoliSync.Underwriting.Persistence;

public class QuoteRepository : IQuoteRepository
{
    private readonly PoliSyncDbContext _db;

    public QuoteRepository(PoliSyncDbContext db) => _db = db;

    public async Task<Quote?> GetByIdAsync(Guid id, CancellationToken ct = default)
    {
        return await _db.Set<Quote>()
            .Include(q => q.HealthDeclaration)
            .Include(q => q.Decision)
            .FirstOrDefaultAsync(q => q.QuoteId == id, ct);
    }

    public async Task<Quote?> GetByNumberAsync(string quoteNumber, CancellationToken ct = default)
    {
        return await _db.Set<Quote>()
            .Include(q => q.HealthDeclaration)
            .Include(q => q.Decision)
            .FirstOrDefaultAsync(q => q.QuoteNumber == quoteNumber, ct);
    }

    public async Task AddAsync(Quote quote, CancellationToken ct = default)
    {
        await _db.Set<Quote>().AddAsync(quote, ct);
    }

    public void Update(Quote quote)
    {
        _db.Set<Quote>().Update(quote);
    }

    public async Task<(IEnumerable<Quote> Items, int TotalCount)> ListAsync(
        Guid? beneficiaryId,
        QuoteStatus? status,
        int page,
        int pageSize,
        CancellationToken ct = default)
    {
        var query = _db.Set<Quote>().AsQueryable();

        if (beneficiaryId.HasValue)
            query = query.Where(q => q.BeneficiaryId == beneficiaryId.Value);

        if (status.HasValue && status.Value != QuoteStatus.Unspecified)
            query = query.Where(q => q.Status == status.Value);

        var total = await query.CountAsync(ct);
        var items = await query
            .OrderByDescending(q => q.CreatedAt)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(ct);

        return (items, total);
    }

    public async Task AddHealthDeclarationAsync(HealthDeclaration declaration, CancellationToken ct = default)
    {
        await _db.Set<HealthDeclaration>().AddAsync(declaration, ct);
    }

    public async Task AddDecisionAsync(UnderwritingDecision decision, CancellationToken ct = default)
    {
        await _db.Set<UnderwritingDecision>().AddAsync(decision, ct);
    }

    public async Task<HealthDeclaration?> GetHealthDeclarationByQuoteIdAsync(Guid quoteId, CancellationToken ct = default)
    {
        return await _db.Set<HealthDeclaration>()
            .FirstOrDefaultAsync(hd => hd.QuoteId == quoteId, ct);
    }

    public async Task<UnderwritingDecision?> GetDecisionByQuoteIdAsync(Guid quoteId, CancellationToken ct = default)
    {
        return await _db.Set<UnderwritingDecision>()
            .FirstOrDefaultAsync(d => d.QuoteId == quoteId, ct);
    }
}
