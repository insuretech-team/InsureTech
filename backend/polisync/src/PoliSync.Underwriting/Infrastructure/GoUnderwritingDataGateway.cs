using Grpc.Core;
using Insuretech.Insurance.Services.V1;
using PoliSync.Infrastructure.Clients;
using HealthDeclarationEntity = Insuretech.Underwriting.Entity.V1.HealthDeclaration;
using QuotationEntity = Insuretech.Policy.Entity.V1.Quotation;
using QuoteEntity = Insuretech.Underwriting.Entity.V1.Quote;
using UnderwritingDecisionEntity = Insuretech.Underwriting.Entity.V1.UnderwritingDecision;

namespace PoliSync.Underwriting.Infrastructure;

public sealed class GoUnderwritingDataGateway : IUnderwritingDataGateway
{
    private readonly InsuranceServiceClient _insuranceClient;

    public GoUnderwritingDataGateway(InsuranceServiceClient insuranceClient)
    {
        _insuranceClient = insuranceClient;
    }

    public async Task<QuoteEntity> CreateQuoteAsync(QuoteEntity quote, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.CreateQuoteAsync(
            new CreateQuoteRequest { Quote = quote },
            cancellationToken: cancellationToken);

        return response.Quote;
    }

    public async Task<QuoteEntity?> GetQuoteAsync(string quoteId, CancellationToken cancellationToken = default)
    {
        try
        {
            var response = await _insuranceClient.Client.GetQuoteAsync(
                new GetQuoteRequest { QuoteId = quoteId },
                cancellationToken: cancellationToken);

            return response.Quote;
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            return null;
        }
    }

    public async Task<QuoteEntity> UpdateQuoteAsync(QuoteEntity quote, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.UpdateQuoteAsync(
            new UpdateQuoteRequest { Quote = quote },
            cancellationToken: cancellationToken);

        return response.Quote;
    }

    public async Task<IReadOnlyList<QuoteEntity>> ListQuotesAsync(string beneficiaryId, int page, int pageSize, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.ListQuotesAsync(
            new ListQuotesRequest
            {
                BeneficiaryId = beneficiaryId ?? string.Empty,
                Page = page,
                PageSize = pageSize
            },
            cancellationToken: cancellationToken);

        return response.Quotes;
    }

    public async Task<HealthDeclarationEntity?> GetHealthDeclarationByQuoteAsync(string quoteId, CancellationToken cancellationToken = default)
    {
        try
        {
            var response = await _insuranceClient.Client.GetHealthDeclarationByQuoteAsync(
                new GetHealthDeclarationByQuoteRequest { QuoteId = quoteId },
                cancellationToken: cancellationToken);

            if (response.Declaration is null || string.IsNullOrWhiteSpace(response.Declaration.Id))
            {
                return null;
            }

            return response.Declaration;
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            return null;
        }
    }

    public async Task<HealthDeclarationEntity> UpsertHealthDeclarationAsync(HealthDeclarationEntity declaration, CancellationToken cancellationToken = default)
    {
        var existing = await GetHealthDeclarationByQuoteAsync(declaration.QuoteId, cancellationToken);
        if (existing is null)
        {
            var created = await _insuranceClient.Client.CreateHealthDeclarationAsync(
                new CreateHealthDeclarationRequest { Declaration = declaration },
                cancellationToken: cancellationToken);

            return created.Declaration;
        }

        declaration.Id = existing.Id;
        var updated = await _insuranceClient.Client.UpdateHealthDeclarationAsync(
            new UpdateHealthDeclarationRequest { Declaration = declaration },
            cancellationToken: cancellationToken);

        return updated.Declaration;
    }

    public async Task<UnderwritingDecisionEntity?> GetLatestDecisionByQuoteAsync(string quoteId, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.ListUnderwritingDecisionsAsync(
            new ListUnderwritingDecisionsRequest { QuoteId = quoteId },
            cancellationToken: cancellationToken);

        return response.Decisions
            .OrderByDescending(x => x.DecidedAt?.Seconds ?? 0)
            .FirstOrDefault();
    }

    public async Task<UnderwritingDecisionEntity> UpsertUnderwritingDecisionAsync(UnderwritingDecisionEntity decision, CancellationToken cancellationToken = default)
    {
        var existing = await GetLatestDecisionByQuoteAsync(decision.QuoteId, cancellationToken);
        if (existing is null)
        {
            var created = await _insuranceClient.Client.CreateUnderwritingDecisionAsync(
                new CreateUnderwritingDecisionRequest { Decision = decision },
                cancellationToken: cancellationToken);

            return created.Decision;
        }

        decision.Id = existing.Id;
        var updated = await _insuranceClient.Client.UpdateUnderwritingDecisionAsync(
            new UpdateUnderwritingDecisionRequest { Decision = decision },
            cancellationToken: cancellationToken);

        return updated.Decision;
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
}
