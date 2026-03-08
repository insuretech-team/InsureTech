using QuotationEntity = Insuretech.Policy.Entity.V1.Quotation;

namespace PoliSync.Quotes.Infrastructure;

public interface IQuotationDataGateway
{
    Task<QuotationEntity> CreateQuotationAsync(QuotationEntity quotation, CancellationToken cancellationToken = default);
    Task<QuotationEntity?> GetQuotationAsync(string quotationId, CancellationToken cancellationToken = default);
    Task<QuotationEntity> UpdateQuotationAsync(QuotationEntity quotation, CancellationToken cancellationToken = default);
    Task<IReadOnlyList<QuotationEntity>> ListQuotationsAsync(string businessId, int page, int pageSize, CancellationToken cancellationToken = default);
    Task DeleteQuotationAsync(string quotationId, CancellationToken cancellationToken = default);
}
