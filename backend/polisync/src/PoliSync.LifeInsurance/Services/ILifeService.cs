using Insuretech.Life.Entity.V1;
using Insuretech.Common.V1;

namespace PoliSync.LifeInsurance.Services;

public interface ILifeProductService
{
    Task<LifeProduct?> GetProductAsync(string productId, CancellationToken cancellationToken = default);
    Task<IEnumerable<LifeProduct>> ListProductsAsync(LifeProductType? productType, bool onlyActive, CancellationToken cancellationToken = default);
    Task<LifeProduct> CreateProductAsync(LifeProduct product, CancellationToken cancellationToken = default);
    Task<LifeProduct?> UpdateProductAsync(string productId, LifeProduct product, CancellationToken cancellationToken = default);
    Task<bool> DeleteProductAsync(string productId, bool permanent = false, CancellationToken cancellationToken = default);
    Task<IEnumerable<ConditionMultiplier>> GetHealthConditionsAsync(string productId, CancellationToken cancellationToken = default);
}

public interface ILifeQuoteService
{
    Task<(long BasePremium, long AgeAddition, float ConditionMultiplier, long ConditionAddition, long BonusDiscount, long TotalPremium, 
          List<PremiumBreakdown> Breakdown, List<string> AppliedConditions, List<string> AppliedBonuses)> 
        CalculatePremiumAsync(
            string productId,
            InsuredPerson insuredPerson,
            int ageAtEntry,
            int policyTermYears,
            long sumAssured,
            List<string> bonusCodes,
            CancellationToken cancellationToken = default);

    Task<LifeQuote> GenerateQuoteAsync(
        string productId,
        string customerId,
        InsuredPerson insuredPerson,
        int ageAtEntry,
        int policyTermYears,
        long sumAssured,
        List<string> bonusCodes,
        string? agentId = null,
        int validityDays = 30,
        CancellationToken cancellationToken = default);

    Task<LifeQuote?> GetQuoteAsync(string quoteId, CancellationToken cancellationToken = default);
    Task<LifeQuote?> GetQuoteByNumberAsync(string quoteNumber, CancellationToken cancellationToken = default);
    Task<IEnumerable<LifeQuote>> ListQuotesAsync(string? customerId, string? productId, LifeQuoteStatus? status, CancellationToken cancellationToken = default);
    Task<bool> ConvertQuoteToPolicyAsync(string quoteId, string policyId, string? convertedBy = null, CancellationToken cancellationToken = default);
}

public interface ILifeProductRepository
{
    Task<LifeProduct?> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<IEnumerable<LifeProduct>> GetAllAsync(CancellationToken cancellationToken = default);
    Task<IEnumerable<LifeProduct>> GetByTypeAsync(LifeProductType type, CancellationToken cancellationToken = default);
    Task<IEnumerable<LifeProduct>> GetByFilterAsync(LifeProductType? type, bool onlyActive, CancellationToken cancellationToken = default);
    Task<LifeProduct> CreateAsync(LifeProduct product, CancellationToken cancellationToken = default);
    Task<LifeProduct?> UpdateAsync(LifeProduct product, CancellationToken cancellationToken = default);
    Task<bool> DeleteAsync(string id, bool permanent = false, CancellationToken cancellationToken = default);
}

public interface ILifeQuoteRepository
{
    Task<LifeQuote?> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<LifeQuote?> GetByNumberAsync(string number, CancellationToken cancellationToken = default);
    Task<IEnumerable<LifeQuote>> GetByCustomerIdAsync(string customerId, CancellationToken cancellationToken = default);
    Task<IEnumerable<LifeQuote>> GetByFilterAsync(string? customerId, string? productId, LifeQuoteStatus? status, CancellationToken cancellationToken = default);
    Task<LifeQuote> CreateAsync(LifeQuote quote, CancellationToken cancellationToken = default);
    Task<LifeQuote?> UpdateAsync(LifeQuote quote, CancellationToken cancellationToken = default);
}

public interface ILifePremiumCalculator
{
    Task<(long BasePremium, long AgeAddition, float ConditionMultiplier, long ConditionAddition, long BonusDiscount, long TotalPremium,
          List<PremiumBreakdown> Breakdown, List<string> AppliedConditions, List<string> AppliedBonuses)>
        CalculatePremiumAsync(
            LifeProduct product,
            InsuredPerson insuredPerson,
            int ageAtEntry,
            int policyTermYears,
            long sumAssured,
            List<string> bonusCodes,
            CancellationToken cancellationToken = default);
}

public interface IQuoteNumberGenerator
{
    string GenerateQuoteNumber();
}
