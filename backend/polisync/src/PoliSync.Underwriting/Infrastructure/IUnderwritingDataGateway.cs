using HealthDeclarationEntity = Insuretech.Underwriting.Entity.V1.HealthDeclaration;
using QuotationEntity = Insuretech.Policy.Entity.V1.Quotation;
using QuoteEntity = Insuretech.Underwriting.Entity.V1.Quote;
using UnderwritingDecisionEntity = Insuretech.Underwriting.Entity.V1.UnderwritingDecision;

namespace PoliSync.Underwriting.Infrastructure;

public interface IUnderwritingDataGateway
{
    Task<QuoteEntity> CreateQuoteAsync(QuoteEntity quote, CancellationToken cancellationToken = default);
    Task<QuoteEntity?> GetQuoteAsync(string quoteId, CancellationToken cancellationToken = default);
    Task<QuoteEntity> UpdateQuoteAsync(QuoteEntity quote, CancellationToken cancellationToken = default);
    Task<IReadOnlyList<QuoteEntity>> ListQuotesAsync(string beneficiaryId, int page, int pageSize, CancellationToken cancellationToken = default);

    Task<HealthDeclarationEntity?> GetHealthDeclarationByQuoteAsync(string quoteId, CancellationToken cancellationToken = default);
    Task<HealthDeclarationEntity> UpsertHealthDeclarationAsync(HealthDeclarationEntity declaration, CancellationToken cancellationToken = default);

    Task<UnderwritingDecisionEntity?> GetLatestDecisionByQuoteAsync(string quoteId, CancellationToken cancellationToken = default);
    Task<UnderwritingDecisionEntity> UpsertUnderwritingDecisionAsync(UnderwritingDecisionEntity decision, CancellationToken cancellationToken = default);

    Task<QuotationEntity?> GetQuotationAsync(string quotationId, CancellationToken cancellationToken = default);
    Task<QuotationEntity> UpdateQuotationAsync(QuotationEntity quotation, CancellationToken cancellationToken = default);
}
