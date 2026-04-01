using System.Collections.Concurrent;
using System.Text.Json;
using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Life.Entity.V1;

namespace PoliSync.LifeInsurance.Services;

public class InMemoryLifeProductRepository : ILifeProductRepository
{
    private readonly ConcurrentDictionary<string, LifeProduct> _products = new();
    private readonly ILogger<InMemoryLifeProductRepository> _logger;

    public InMemoryLifeProductRepository(ILogger<InMemoryLifeProductRepository> logger)
    {
        _logger = logger;
        SeedDefaultProducts();
    }

    public Task<LifeProduct?> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        _products.TryGetValue(id, out var product);
        return Task.FromResult(product);
    }

    public Task<IEnumerable<LifeProduct>> GetAllAsync(CancellationToken cancellationToken = default)
    {
        return Task.FromResult(_products.Values.AsEnumerable());
    }

    public Task<IEnumerable<LifeProduct>> GetByTypeAsync(LifeProductType type, CancellationToken cancellationToken = default)
    {
        var products = _products.Values
            .Where(p => p.ProductType == type)
            .AsEnumerable();
        return Task.FromResult(products);
    }

    public Task<IEnumerable<LifeProduct>> GetByFilterAsync(LifeProductType? type, bool onlyActive, CancellationToken cancellationToken = default)
    {
        var query = _products.Values.AsEnumerable();
        
        if (type.HasValue)
        {
            query = query.Where(p => p.ProductType == type.Value);
        }
        
        if (onlyActive)
        {
            query = query.Where(p => p.IsActive);
        }
        
        return Task.FromResult(query);
    }

    public Task<LifeProduct> CreateAsync(LifeProduct product, CancellationToken cancellationToken = default)
    {
        if (string.IsNullOrEmpty(product.ProductId))
        {
            product.ProductId = Guid.NewGuid().ToString();
        }
        
        product.CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        product.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        
        _products[product.ProductId] = product;
        _logger.LogInformation("Created life product: {ProductId} - {ProductCode}", 
            product.ProductId, product.ProductCode);
        
        return Task.FromResult(product);
    }

    public Task<LifeProduct?> UpdateAsync(LifeProduct product, CancellationToken cancellationToken = default)
    {
        if (!_products.ContainsKey(product.ProductId))
        {
            return Task.FromResult<LifeProduct?>(null);
        }

        product.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        _products[product.ProductId] = product;
        
        _logger.LogInformation("Updated life product: {ProductId}", product.ProductId);
        
        return Task.FromResult<LifeProduct?>(product);
    }

    public Task<bool> DeleteAsync(string id, bool permanent = false, CancellationToken cancellationToken = default)
    {
        if (permanent)
        {
            var result = _products.TryRemove(id, out _);
            if (result)
            {
                _logger.LogInformation("Permanently deleted life product: {ProductId}", id);
            }
            return Task.FromResult(result);
        }
        else
        {
            if (_products.TryGetValue(id, out var product))
            {
                product.DeletedAt = Timestamp.FromDateTime(DateTime.UtcNow);
                product.IsActive = false;
                _logger.LogInformation("Soft deleted life product: {ProductId}", id);
                return Task.FromResult(true);
            }
            return Task.FromResult(false);
        }
    }

    private void SeedDefaultProducts()
    {
        // Term Life Insurance
        var termProduct = new LifeProduct
        {
            ProductId = Guid.NewGuid().ToString(),
            ProductCode = "LIFE-TERM-001",
            ProductName = "Secure Future Term Life",
            ProductType = LifeProductType.Term,
            Description = "Affordable term life insurance with high coverage",
            BaseRate = 50000, // 500 BDT base rate
            AgeAdditionConfig = new AgeAdditionConfig
            {
                StartAge = 18,
                AgeIncrement = 5,
                PriceToAdd = 10000 // 100 BDT per 5 years
            },
            MinSumAssured = 10000000, // 100,000 BDT
            MaxSumAssured = 500000000, // 5,000,000 BDT
            MinEntryAge = 18,
            MaxEntryAge = 60,
            MinPolicyTerm = 5,
            MaxPolicyTerm = 30,
            IsActive = true
        };

        // Set condition multipliers
        var conditionMultipliers = new List<ConditionMultiplier>
        {
            new ConditionMultiplier { ConditionCode = "DIABETES", ConditionName = "Diabetes", Multiplier = 0.2f, Severity = "MEDIUM", Description = "Type 1 or Type 2 diabetes" },
            new ConditionMultiplier { ConditionCode = "HYPERTENSION", ConditionName = "Hypertension", Multiplier = 0.15f, Severity = "MEDIUM", Description = "High blood pressure" },
            new ConditionMultiplier { ConditionCode = "SMOKING", ConditionName = "Smoking", Multiplier = 0.3f, Severity = "HIGH", Description = "Current smoker" },
            new ConditionMultiplier { ConditionCode = "HEART_DISEASE", ConditionName = "Heart Disease", Multiplier = 0.5f, Severity = "HIGH", Description = "Heart conditions" },
            new ConditionMultiplier { ConditionCode = "CANCER", ConditionName = "Cancer History", Multiplier = 0.4f, Severity = "HIGH", Description = "Previous cancer diagnosis" },
            new ConditionMultiplier { ConditionCode = "OBESITY", ConditionName = "Obesity", Multiplier = 0.15f, Severity = "MEDIUM", Description = "BMI > 30" }
        };
        termProduct.ConditionMultipliersJson = JsonSerializer.Serialize(conditionMultipliers);

        // Set bonus config
        var bonusConfigs = new List<BonusConfig>
        {
            new BonusConfig { BonusCode = "NON_SMOKER", BonusName = "Non-Smoker Discount", BonusType = "PERCENTAGE", Percentage = 0.1f, Description = "10% discount for non-smokers" },
            new BonusConfig { BonusCode = "HEALTHY", BonusName = "Healthy Lifestyle", BonusType = "PERCENTAGE", Percentage = 0.05f, Description = "5% discount for healthy lifestyle" },
            new BonusConfig { BonusCode = "EARLY_BIRD", BonusName = "Early Bird", BonusType = "PERCENTAGE", Percentage = 0.03f, Description = "3% discount for applying before age 30" }
        };
        termProduct.BonusConfigJson = JsonSerializer.Serialize(bonusConfigs);

        // Whole Life Insurance
        var wholeLifeProduct = new LifeProduct
        {
            ProductId = Guid.NewGuid().ToString(),
            ProductCode = "LIFE-WHOLE-001",
            ProductName = "Lifetime Protection Plan",
            ProductType = LifeProductType.WholeLife,
            Description = "Whole life insurance with cash value accumulation",
            BaseRate = 100000, // 1000 BDT base rate
            AgeAdditionConfig = new AgeAdditionConfig
            {
                StartAge = 18,
                AgeIncrement = 5,
                PriceToAdd = 20000 // 200 BDT per 5 years
            },
            MinSumAssured = 20000000, // 200,000 BDT
            MaxSumAssured = 1000000000, // 10,000,000 BDT
            MinEntryAge = 18,
            MaxEntryAge = 55,
            MinPolicyTerm = 10,
            MaxPolicyTerm = 40,
            IsActive = true,
            ConditionMultipliersJson = JsonSerializer.Serialize(conditionMultipliers),
            BonusConfigJson = JsonSerializer.Serialize(bonusConfigs)
        };

        // Endowment Plan
        var endowmentProduct = new LifeProduct
        {
            ProductId = Guid.NewGuid().ToString(),
            ProductCode = "LIFE-ENDOW-001",
            ProductName = "Smart Savings Endowment",
            ProductType = LifeProductType.Endowment,
            Description = "Savings plan with life coverage and maturity benefits",
            BaseRate = 75000, // 750 BDT base rate
            AgeAdditionConfig = new AgeAdditionConfig
            {
                StartAge = 18,
                AgeIncrement = 5,
                PriceToAdd = 15000 // 150 BDT per 5 years
            },
            MinSumAssured = 15000000, // 150,000 BDT
            MaxSumAssured = 750000000, // 7,500,000 BDT
            MinEntryAge = 18,
            MaxEntryAge = 50,
            MinPolicyTerm = 10,
            MaxPolicyTerm = 25,
            IsActive = true,
            ConditionMultipliersJson = JsonSerializer.Serialize(conditionMultipliers),
            BonusConfigJson = JsonSerializer.Serialize(bonusConfigs)
        };

        foreach (var product in new[] { termProduct, wholeLifeProduct, endowmentProduct })
        {
            product.CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            product.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            _products[product.ProductId] = product;
        }

        _logger.LogInformation("Seeded {Count} default life insurance products", _products.Count);
    }
}
