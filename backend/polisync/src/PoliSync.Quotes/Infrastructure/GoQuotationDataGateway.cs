using Grpc.Core;
using Insuretech.Insurance.Services.V1;
using PoliSync.Infrastructure.Clients;
using QuotationEntity = Insuretech.Policy.Entity.V1.Quotation;

namespace PoliSync.Quotes.Infrastructure;

public sealed class GoQuotationDataGateway : IQuotationDataGateway
{
    private readonly InsuranceServiceClient _insuranceClient;

    public GoQuotationDataGateway(InsuranceServiceClient insuranceClient)
    {
        _insuranceClient = insuranceClient;
    }

    public async Task<QuotationEntity> CreateQuotationAsync(QuotationEntity quotation, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.CreateQuotationAsync(
            new CreateQuotationRequest { Quotation = quotation },
            cancellationToken: cancellationToken);

        return response.Quotation;
    }

    public async Task<QuotationEntity?> GetQuotationAsync(string quotationId, CancellationToken cancellationToken = default)
    {
        try
        {
            var response = await _insuranceClient.Client.GetQuotationAsync(
                new GetQuotationRequest { QuotationId = quotationId },
                cancellationToken: cancellationToken);

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
            cancellationToken: cancellationToken);

        return response.Quotation;
    }

    public async Task<IReadOnlyList<QuotationEntity>> ListQuotationsAsync(string businessId, int page, int pageSize, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.ListQuotationsAsync(
            new ListQuotationsRequest
            {
                BusinessId = businessId,
                Page = page,
                PageSize = pageSize
            },
            cancellationToken: cancellationToken);

        return response.Quotations;
    }

    public Task DeleteQuotationAsync(string quotationId, CancellationToken cancellationToken = default)
    {
        return _insuranceClient.Client.DeleteQuotationAsync(
            new DeleteQuotationRequest { QuotationId = quotationId },
            cancellationToken: cancellationToken).ResponseAsync;
    }
}
