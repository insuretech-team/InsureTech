using Grpc.Core;
using Insuretech.Insurance.Services.V1;
using PoliSync.Infrastructure.Clients;
using DomainQuotation = PoliSync.Quotes.Domain.Quotation;
using DomainQuotationStatus = PoliSync.Quotes.Domain.QuotationStatus;
using QuotationEntity = Insuretech.Policy.Entity.V1.Quotation;

namespace PoliSync.Quotes.Infrastructure;

public sealed class GoQuotationDataGateway : IQuotationDataGateway
{
    private readonly InsuranceServiceClient _insuranceClient;

    public GoQuotationDataGateway(InsuranceServiceClient insuranceClient) =>
        _insuranceClient = insuranceClient;

    public async Task CreateAsync(DomainQuotation quotation, CancellationToken cancellationToken = default)
    {
        await CreateQuotationAsync(QuotationProtoMapper.ToProto(quotation), cancellationToken);
    }

    public async Task<DomainQuotation?> GetByIdAsync(Guid quotationId, CancellationToken cancellationToken = default)
    {
        var quotation = await GetQuotationAsync(quotationId.ToString(), cancellationToken);
        return quotation is null ? null : QuotationProtoMapper.ToDomain(quotation);
    }

    public async Task UpdateAsync(DomainQuotation quotation, CancellationToken cancellationToken = default)
    {
        await UpdateQuotationAsync(QuotationProtoMapper.ToProto(quotation), cancellationToken);
    }

    public async Task<(IReadOnlyList<DomainQuotation> Quotations, int TotalCount)> ListAsync(
        Guid tenantId,
        Guid? customerId,
        DomainQuotationStatus? status,
        int pageNumber,
        int pageSize,
        CancellationToken cancellationToken = default)
    {
        var quotations = await ListQuotationsAsync(
            tenantId == Guid.Empty ? string.Empty : tenantId.ToString(),
            pageNumber,
            pageSize,
            cancellationToken);

        var mapped = quotations
            .Select(QuotationProtoMapper.ToDomain)
            .Where(q => customerId is null || q.CustomerId == customerId.Value)
            .Where(q => status is null || q.Status == status.Value)
            .ToList();

        return (mapped, mapped.Count);
    }

    public async Task<IReadOnlyList<DomainQuotation>> GetExpiredQuotationsAsync(CancellationToken cancellationToken = default)
    {
        const int pageSize = 200;
        var page = 1;
        var expired = new List<DomainQuotation>();

        while (true)
        {
            var batch = await ListQuotationsAsync(string.Empty, page, pageSize, cancellationToken);
            if (batch.Count == 0)
            {
                break;
            }

            expired.AddRange(batch
                .Select(QuotationProtoMapper.ToDomain)
                .Where(q => q.ExpiryDate < DateTime.UtcNow)
                .Where(q => q.Status is not (DomainQuotationStatus.Approved or DomainQuotationStatus.Rejected or DomainQuotationStatus.Expired)));

            if (batch.Count < pageSize)
            {
                break;
            }

            page++;
        }

        return expired;
    }

    public async Task<QuotationEntity> CreateQuotationAsync(QuotationEntity quotation, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.CreateQuotationAsync(
            new CreateQuotationRequest { Quotation = quotation },
            _insuranceClient.BuildCallOptions(cancellationToken));
        return response.Quotation;
    }

    public async Task<QuotationEntity?> GetQuotationAsync(string quotationId, CancellationToken cancellationToken = default)
    {
        try
        {
            var response = await _insuranceClient.Client.GetQuotationAsync(
                new GetQuotationRequest { QuotationId = quotationId },
                _insuranceClient.BuildCallOptions(cancellationToken));
            return response.Quotation;
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            return null;
        }
    }

    public async Task<QuotationEntity> UpdateQuotationAsync(QuotationEntity quotation, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.UpdateQuotationAsync(
            new UpdateQuotationRequest { Quotation = quotation },
            _insuranceClient.BuildCallOptions(cancellationToken));
        return response.Quotation;
    }

    public async Task<IReadOnlyList<QuotationEntity>> ListQuotationsAsync(string businessId, int page, int pageSize, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.ListQuotationsAsync(
            new ListQuotationsRequest { BusinessId = businessId, Page = page, PageSize = pageSize },
            _insuranceClient.BuildCallOptions(cancellationToken));
        return response.Quotations;
    }

    public Task DeleteQuotationAsync(string quotationId, CancellationToken cancellationToken = default)
    {
        return _insuranceClient.Client.DeleteQuotationAsync(
            new DeleteQuotationRequest { QuotationId = quotationId },
            _insuranceClient.BuildCallOptions(cancellationToken)).ResponseAsync;
    }
}

