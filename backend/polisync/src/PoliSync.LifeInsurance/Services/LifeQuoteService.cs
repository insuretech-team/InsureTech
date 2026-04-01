using System.Text.Json;
using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Life.Entity.V1;

namespace PoliSync.LifeInsurance.Services;

public class LifeQuoteService : ILifeQuoteService
{
    private readonly ILifeProductRepository _productRepository;
    private readonly ILifeQuoteRepository _quoteRepository;
    private readonly ILifePremiumCalculator _premiumCalculator;
    private readonly IQuoteNumberGenerator _quoteNumberGenerator;
    private readonly ILogger<LifeQuoteService> _logger;

    public LifeQuoteService(
        ILifeProductRepository productRepository,
        ILifeQuoteRepository quoteRepository,
        ILifePremiumCalculator premiumCalculator,
        IQuoteNumberGenerator quoteNumberGenerator,
        ILogger<LifeQuoteService> logger)
    {
        _productRepository = productRepository;
        _quoteRepository = quoteRepository;
        _premiumCalculator = premiumCalculator;
        _quoteNumberGenerator = quoteNumberGenerator;
        _logger = logger;
    }

    public async Task<(long BasePremium, long AgeAddition, float ConditionMultiplier, long ConditionAddition, long BonusDiscount, long TotalPremium,
          List<PremiumBreakdown> Breakdown, List<string> AppliedConditions, List<string> AppliedBonuses)> 
        CalculatePremiumAsync(
            string productId,
            InsuredPerson insuredPerson,
            int ageAtEntry,
            int policyTermYears,
            long sumAssured,
            List<string> bonusCodes,
            CancellationToken cancellationToken = default)
    {
        var product = await _productRepository.GetByIdAsync(productId, cancellationToken);
        if (product == null)
        {
            throw new InvalidOperationException($"Life product {productId} not found");
        }

        return await _premiumCalculator.CalculatePremiumAsync(
            product, insuredPerson, ageAtEntry, policyTermYears, sumAssured, bonusCodes, cancellationToken);
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
        _logger.LogInformation(
            "Generating life quote for product {ProductId}, customer {CustomerId}",
            productId, customerId);

        var product = await _productRepository.GetByIdAsync(productId, cancellationToken);
        if (product == null)
        {
            throw new InvalidOperationException($"Life product {productId} not found");
        }

        // Calculate premium
        var (basePremium, ageAddition, conditionMultiplier, conditionAddition, bonusDiscount, totalPremium,
             breakdown, appliedConditions, appliedBonuses) = await _premiumCalculator.CalculatePremiumAsync(
            product, insuredPerson, ageAtEntry, policyTermYears, sumAssured, bonusCodes, cancellationToken);

        // Create quote
        var quote = new LifeQuote
        {
            QuoteId = Guid.NewGuid().ToString(),
            QuoteNumber = _quoteNumberGenerator.GenerateQuoteNumber(),
            ProductId = productId,
            CustomerId = customerId,
            AgentId = agentId ?? string.Empty,
            Status = LifeQuoteStatus.Generated,
            InsuredPersonJson = JsonSerializer.Serialize(insuredPerson),
            AgeAtEntry = ageAtEntry,
            PolicyTermYears = policyTermYears,
            SumAssured = sumAssured,
            BasePremium = basePremium,
            AgeAddition = ageAddition,
            ConditionMultiplier = conditionMultiplier,
            ConditionAddition = conditionAddition,
            BonusDiscount = bonusDiscount,
            TotalPremium = totalPremium,
            HealthConditionsJson = JsonSerializer.Serialize(insuredPerson.HealthConditions),
            BonusesAppliedJson = JsonSerializer.Serialize(appliedBonuses),
            ValidUntil = Timestamp.FromDateTime(DateTime.UtcNow.AddDays(validityDays)),
            CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow),
            UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow)
        };

        var created = await _quoteRepository.CreateAsync(quote, cancellationToken);
        
        _logger.LogInformation(
            "Life quote generated: {QuoteNumber} with premium {Premium}",
            created.QuoteNumber, created.TotalPremium);

        return created;
    }

    public Task<LifeQuote?> GetQuoteAsync(string quoteId, CancellationToken cancellationToken = default)
    {
        return _quoteRepository.GetByIdAsync(quoteId, cancellationToken);
    }

    public Task<LifeQuote?> GetQuoteByNumberAsync(string quoteNumber, CancellationToken cancellationToken = default)
    {
        return _quoteRepository.GetByNumberAsync(quoteNumber, cancellationToken);
    }

    public Task<IEnumerable<LifeQuote>> ListQuotesAsync(
        string? customerId, 
        string? productId, 
        LifeQuoteStatus? status, 
        CancellationToken cancellationToken = default)
    {
        return _quoteRepository.GetByFilterAsync(customerId, productId, status, cancellationToken);
    }

    public async Task<bool> ConvertQuoteToPolicyAsync(
        string quoteId, 
        string policyId, 
        string? convertedBy = null,
        CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Converting life quote {QuoteId} to policy {PolicyId}", quoteId, policyId);

        var quote = await _quoteRepository.GetByIdAsync(quoteId, cancellationToken);
        if (quote == null)
        {
            return false;
        }

        if (quote.Status != LifeQuoteStatus.Generated && 
            quote.Status != LifeQuoteStatus.Sent && 
            quote.Status != LifeQuoteStatus.Viewed)
        {
            _logger.LogWarning("Quote {QuoteId} cannot be converted. Status: {Status}", quoteId, quote.Status);
            return false;
        }

        quote.Status = LifeQuoteStatus.Converted;
        quote.ConvertedPolicyId = policyId;
        quote.ConvertedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        quote.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);

        await _quoteRepository.UpdateAsync(quote, cancellationToken);
        
        _logger.LogInformation("Life quote {QuoteNumber} converted to policy {PolicyId}", 
            quote.QuoteNumber, policyId);

        return true;
    }
}
