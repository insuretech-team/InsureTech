using DomainQuotation = PoliSync.Quotes.Domain.Quotation;
using DomainQuotationStatus = PoliSync.Quotes.Domain.QuotationStatus;
using QuotationEntity = Insuretech.Policy.Entity.V1.Quotation;

namespace PoliSync.Quotes.Infrastructure;

public interface IQuotationDataGateway
{
    Task CreateAsync(DomainQuotation quotation, CancellationToken cancellationToken = default);
    Task<DomainQuotation?> GetByIdAsync(Guid quotationId, CancellationToken cancellationToken = default);
    Task UpdateAsync(DomainQuotation quotation, CancellationToken cancellationToken = default);
    Task<(IReadOnlyList<DomainQuotation> Quotations, int TotalCount)> ListAsync(
        Guid tenantId,
        Guid? customerId,
        DomainQuotationStatus? status,
        int pageNumber,
        int pageSize,
        CancellationToken cancellationToken = default);
    Task<IReadOnlyList<DomainQuotation>> GetExpiredQuotationsAsync(CancellationToken cancellationToken = default);

    Task<QuotationEntity> CreateQuotationAsync(QuotationEntity quotation, CancellationToken cancellationToken = default);
    Task<QuotationEntity?> GetQuotationAsync(string quotationId, CancellationToken cancellationToken = default);
    Task<QuotationEntity> UpdateQuotationAsync(QuotationEntity quotation, CancellationToken cancellationToken = default);
    Task<IReadOnlyList<QuotationEntity>> ListQuotationsAsync(string businessId, int page, int pageSize, CancellationToken cancellationToken = default);
    Task DeleteQuotationAsync(string quotationId, CancellationToken cancellationToken = default);
}
