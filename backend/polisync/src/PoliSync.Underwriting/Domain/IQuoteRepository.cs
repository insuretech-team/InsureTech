namespace PoliSync.Underwriting.Domain;

public interface IQuoteRepository
{
    Task<Quote?> GetByIdAsync(Guid id, CancellationToken ct = default);
    Task<Quote?> GetByNumberAsync(string quoteNumber, CancellationToken ct = default);
    Task AddAsync(Quote quote, CancellationToken ct = default);
    void Update(Quote quote);
    
    Task<(IEnumerable<Quote> Items, int TotalCount)> ListAsync(
        Guid? beneficiaryId,
        QuoteStatus? status,
        int page,
        int pageSize,
        CancellationToken ct = default);

    Task AddHealthDeclarationAsync(HealthDeclaration declaration, CancellationToken ct = default);
    Task AddDecisionAsync(UnderwritingDecision decision, CancellationToken ct = default);
    Task<HealthDeclaration?> GetHealthDeclarationByQuoteIdAsync(Guid quoteId, CancellationToken ct = default);
    Task<UnderwritingDecision?> GetDecisionByQuoteIdAsync(Guid quoteId, CancellationToken ct = default);
}
