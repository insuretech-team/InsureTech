using System.Text.Json;
using Insuretech.Life.Entity.V1;
using Insuretech.Life.Services.V1;
using PoliSync.Infrastructure.Clients;

namespace PoliSync.LifeInsurance.Services;

public sealed class GoLifeProductService : ILifeProductService
{
    private readonly InsuranceServiceClient _client;

    public GoLifeProductService(InsuranceServiceClient client) => _client = client;

    public async Task<LifeProduct?> GetProductAsync(string productId, CancellationToken cancellationToken = default)
    {
        var response = await _client.LifeClient.GetLifeProductAsync(
            new GetLifeProductRequest { ProductId = productId },
            _client.BuildCallOptions(cancellationToken));
        return response.Error is null ? response.Product : null;
    }

    public async Task<IEnumerable<LifeProduct>> ListProductsAsync(LifeProductType? productType, bool onlyActive, CancellationToken cancellationToken = default)
    {
        var response = await _client.LifeClient.ListLifeProductsAsync(
            new ListLifeProductsRequest
            {
                ProductType = productType ?? LifeProductType.Unspecified,
                OnlyActive = onlyActive,
                PageSize = 200
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Products;
    }

    public async Task<LifeProduct> CreateProductAsync(LifeProduct product, CancellationToken cancellationToken = default)
    {
        var response = await _client.LifeClient.CreateLifeProductAsync(
            new CreateLifeProductRequest
            {
                ProductCode = product.ProductCode,
                ProductName = product.ProductName,
                ProductType = product.ProductType,
                Description = product.Description,
                BaseRate = product.BaseRate,
                AgeAdditionConfig = product.AgeAdditionConfig,
                ConditionMultipliers = { Deserialize<List<ConditionMultiplier>>(product.ConditionMultipliersJson) ?? [] },
                Bonuses = { Deserialize<List<BonusConfig>>(product.BonusConfigJson) ?? [] },
                MinSumAssured = product.MinSumAssured,
                MaxSumAssured = product.MaxSumAssured,
                MinEntryAge = product.MinEntryAge,
                MaxEntryAge = product.MaxEntryAge,
                MinPolicyTerm = product.MinPolicyTerm,
                MaxPolicyTerm = product.MaxPolicyTerm,
                Metadata = { product.Metadata }
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Product;
    }

    public async Task<LifeProduct?> UpdateProductAsync(string productId, LifeProduct product, CancellationToken cancellationToken = default)
    {
        var response = await _client.LifeClient.UpdateLifeProductAsync(
            new UpdateLifeProductRequest
            {
                ProductId = productId,
                ProductName = product.ProductName,
                Description = product.Description,
                BaseRate = product.BaseRate,
                AgeAdditionConfig = product.AgeAdditionConfig,
                ConditionMultipliers = { Deserialize<List<ConditionMultiplier>>(product.ConditionMultipliersJson) ?? [] },
                Bonuses = { Deserialize<List<BonusConfig>>(product.BonusConfigJson) ?? [] },
                MinSumAssured = product.MinSumAssured,
                MaxSumAssured = product.MaxSumAssured,
                MinEntryAge = product.MinEntryAge,
                MaxEntryAge = product.MaxEntryAge,
                MinPolicyTerm = product.MinPolicyTerm,
                MaxPolicyTerm = product.MaxPolicyTerm,
                IsActive = product.IsActive,
                Metadata = { product.Metadata }
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Error is null ? response.Product : null;
    }

    public async Task<bool> DeleteProductAsync(string productId, bool permanent = false, CancellationToken cancellationToken = default)
    {
        var response = await _client.LifeClient.DeleteLifeProductAsync(
            new DeleteLifeProductRequest { ProductId = productId, Permanent = permanent },
            _client.BuildCallOptions(cancellationToken));
        return response.Success && response.Error is null;
    }

    public async Task<IEnumerable<ConditionMultiplier>> GetHealthConditionsAsync(string productId, CancellationToken cancellationToken = default)
    {
        var response = await _client.LifeClient.GetHealthConditionsAsync(
            new GetHealthConditionsRequest { ProductId = productId },
            _client.BuildCallOptions(cancellationToken));
        return response.Conditions;
    }

    private static T? Deserialize<T>(string json)
    {
        return string.IsNullOrWhiteSpace(json) ? default : JsonSerializer.Deserialize<T>(json);
    }
}

public sealed class GoLifeQuoteService : ILifeQuoteService
{
    private readonly InsuranceServiceClient _client;

    public GoLifeQuoteService(InsuranceServiceClient client) => _client = client;

    public async Task<(long BasePremium, long AgeAddition, float ConditionMultiplier, long ConditionAddition, long BonusDiscount, long TotalPremium, List<PremiumBreakdown> Breakdown, List<string> AppliedConditions, List<string> AppliedBonuses)> CalculatePremiumAsync(
        string productId,
        InsuredPerson insuredPerson,
        int ageAtEntry,
        int policyTermYears,
        long sumAssured,
        List<string> bonusCodes,
        CancellationToken cancellationToken = default)
    {
        var response = await _client.LifeClient.CalculatePremiumAsync(
            new CalculatePremiumRequest
            {
                ProductId = productId,
                InsuredPerson = insuredPerson,
                AgeAtEntry = ageAtEntry,
                PolicyTermYears = policyTermYears,
                SumAssured = sumAssured,
                BonusCodes = { bonusCodes }
            },
            _client.BuildCallOptions(cancellationToken));

        return (
            response.BasePremium,
            response.AgeAddition,
            response.ConditionMultiplier,
            response.ConditionAddition,
            response.BonusDiscount,
            response.TotalPremium,
            response.Breakdown.ToList(),
            response.AppliedConditions.ToList(),
            response.AppliedBonuses.ToList());
    }

    public async Task<LifeQuote> GenerateQuoteAsync(
        string productId,
        string customerId,
        InsuredPerson insuredPerson,
        int ageAtEntry,
        int policyTermYears,
        long sumAssured,
        List<string> bonusCodes,
        string? agentId = null,
        int validityDays = 30,
        CancellationToken cancellationToken = default)
    {
        var response = await _client.LifeClient.GenerateQuoteAsync(
            new GenerateQuoteRequest
            {
                ProductId = productId,
                CustomerId = customerId,
                AgentId = agentId ?? string.Empty,
                InsuredPerson = insuredPerson,
                AgeAtEntry = ageAtEntry,
                PolicyTermYears = policyTermYears,
                SumAssured = sumAssured,
                BonusCodes = { bonusCodes },
                ValidityDays = validityDays
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Quote;
    }

    public async Task<LifeQuote?> GetQuoteAsync(string quoteId, CancellationToken cancellationToken = default)
    {
        var response = await _client.LifeClient.GetQuoteAsync(
            new GetQuoteRequest { QuoteId = quoteId },
            _client.BuildCallOptions(cancellationToken));
        return response.Error is null ? response.Quote : null;
    }

    public async Task<LifeQuote?> GetQuoteByNumberAsync(string quoteNumber, CancellationToken cancellationToken = default)
    {
        var response = await _client.LifeClient.GetQuoteByNumberAsync(
            new GetQuoteByNumberRequest { QuoteNumber = quoteNumber },
            _client.BuildCallOptions(cancellationToken));
        return response.Error is null ? response.Quote : null;
    }

    public async Task<IEnumerable<LifeQuote>> ListQuotesAsync(string? customerId, string? productId, LifeQuoteStatus? status, CancellationToken cancellationToken = default)
    {
        var response = await _client.LifeClient.ListQuotesAsync(
            new ListQuotesRequest
            {
                CustomerId = customerId ?? string.Empty,
                ProductId = productId ?? string.Empty,
                Status = status ?? LifeQuoteStatus.Unspecified,
                PageSize = 200
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Quotes;
    }

    public async Task<bool> ConvertQuoteToPolicyAsync(string quoteId, string policyId, string? convertedBy = null, CancellationToken cancellationToken = default)
    {
        var response = await _client.LifeClient.ConvertQuoteToPolicyAsync(
            new ConvertQuoteToPolicyRequest
            {
                QuoteId = quoteId,
                PolicyId = policyId,
                ConvertedBy = convertedBy ?? string.Empty
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Error is null;
    }
}
