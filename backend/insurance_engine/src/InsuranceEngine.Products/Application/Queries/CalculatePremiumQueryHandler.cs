using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Products.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.Products.Domain.Entities;
using Microsoft.Extensions.Caching.Distributed;
using System.Text.Json;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Products.Application.Queries;

public sealed class CalculatePremiumQueryHandler : IRequestHandler<CalculatePremiumQuery, CalculatePremiumResponse>
{
    private readonly IRepository<ProductEntity> _productRepository;
    private readonly IRepository<ProductRiderEntity> _riderRepository;
    private readonly ILogger<CalculatePremiumQueryHandler> _logger;
    private readonly IDistributedCache _cache;

    public CalculatePremiumQueryHandler(
        IRepository<ProductEntity> productRepository,
        IRepository<ProductRiderEntity> riderRepository,
        ILogger<CalculatePremiumQueryHandler> logger,
        IDistributedCache cache)
    {
        _productRepository = productRepository;
        _riderRepository = riderRepository;
        _logger = logger;
        _cache = cache;
    }

    public async Task<CalculatePremiumResponse> Handle(CalculatePremiumQuery request, CancellationToken cancellationToken)
    {
        // 0. Extract loading factors from ApplicantData (FR-062)
        int age = 30; 
        if (request.ApplicantData.TryGetValue("age", out var ageStr) && int.TryParse(ageStr, out var parsedAge))
            age = parsedAge;

        string occupation = "Normal";
        if (request.ApplicantData.TryGetValue("occupation", out var occ))
            occupation = occ;

        bool hasPreExistingCondition = false;
        if (request.ApplicantData.TryGetValue("has_pre_existing_condition", out var pecStr) && bool.TryParse(pecStr, out var pec))
            hasPreExistingCondition = pec;

        int units = 1;
        if (request.ApplicantData.TryGetValue("units", out var unitsStr) && int.TryParse(unitsStr, out var parsedUnits))
            units = parsedUnits;

        // 1. Check Cache (FR-025 strategy)
        string riderKey = request.RiderIds != null ? string.Join("-", request.RiderIds.OrderBy(id => id)) : "none";
        string cacheKey = $"product:premium:{request.ProductId}:{age}:{request.SumInsured.Amount}:{request.TenureMonths}:{riderKey}";
        
        var cachedData = await _cache.GetStringAsync(cacheKey, cancellationToken);
        if (!string.IsNullOrEmpty(cachedData))
        {
            try
            {
                return JsonSerializer.Deserialize<CalculatePremiumResponse>(cachedData)!;
            }
            catch (JsonException) { }
        }

        try
        {
            // 2. Fetch Product (Include riders)
            var product = await _productRepository.GetByIdAsync(Guid.Parse(request.ProductId), cancellationToken);
            if (product == null)
            {
                return new CalculatePremiumResponse
                {
                    Error = new Error { Code = "PRODUCT_NOT_FOUND", Message = "Product not found" }
                };
            }

            // 3. Basic Validation
            if (request.SumInsured.Amount < product.MinSumInsured || request.SumInsured.Amount > product.MaxSumInsured)
            {
                 return new CalculatePremiumResponse
                {
                    Error = new Error { Code = "INVALID_SUM_INSURED", Message = $"Sum insured must be between {product.MinSumInsured/100m} and {product.MaxSumInsured/100m}" }
                };
            }

            // 4. Calculate Premium (FR-024, FR-062)
            decimal rateMultiplier = 1.0m;
            if (age > 50) rateMultiplier += 0.2m; // 20% increase for senior citizens
            if (request.TenureMonths > 12) rateMultiplier -= 0.05m; // 5% discount for long term
            
            // Occupation Loading (FR-062)
            decimal occupationLoading = 0;
            if (occupation.Contains("High Risk", StringComparison.OrdinalIgnoreCase) || occupation == "Driver" || occupation == "Construction")
            {
                occupationLoading = 0.25m;
                rateMultiplier += occupationLoading;
            }

            // Health Condition Loading (FR-062)
            decimal healthLoading = 0;
            if (hasPreExistingCondition)
            {
                healthLoading = 0.30m;
                rateMultiplier += healthLoading;
            }

            // Adjust Sum Insured based on units if provided (FR-023-A)
            long targetSumInsured = request.SumInsured.Amount;
            if (units > 1) 
            {
                targetSumInsured = (long)product.UnitAmount * units;
            }

            long basePremiumAmount = (long)(product.BasePremium * request.TenureMonths * rateMultiplier);
            
            // For complex sum-assured based premium
            long coveragePremium = (long)(targetSumInsured * 0.01m * rateMultiplier); 
            
            long totalBasePremium = basePremiumAmount + coveragePremium;

            var response = new CalculatePremiumResponse
            {
                BasePremium = new Money { Amount = totalBasePremium, Currency = "BDT" }
            };

            response.Breakdown.Add(new PremiumBreakdown 
            { 
                Item = "Base Policy Premium", 
                Amount = new Money { Amount = basePremiumAmount, Currency = "BDT" },
                Description = $"Base rate adjusted by age ({age}) and tenure ({request.TenureMonths}m)"
            });

            response.Breakdown.Add(new PremiumBreakdown 
            { 
                Item = "Coverage Adjustment", 
                Amount = new Money { Amount = coveragePremium, Currency = "BDT" },
                Description = units > 1 
                    ? $"Premium for {units} units of coverage (Unit Size: {product.UnitAmount/100m} BDT)"
                    : "Premium proportional to selected sum insured"
            });

            if (occupationLoading > 0)
            {
                response.Breakdown.Add(new PremiumBreakdown 
                { 
                    Item = "Occupation Risk Loading", 
                    Amount = new Money { Amount = (long)(totalBasePremium * occupationLoading / rateMultiplier), Currency = "BDT" },
                    Description = $"25% loading for high-risk occupation: {occupation}"
                });
            }

            if (healthLoading > 0)
            {
                response.Breakdown.Add(new PremiumBreakdown 
                { 
                    Item = "Health Risk Loading", 
                    Amount = new Money { Amount = (long)(totalBasePremium * healthLoading / rateMultiplier), Currency = "BDT" },
                    Description = "30% loading for pre-existing medical conditions"
                });
            }

            // 5. Riders Calculation
            long totalRiderPremium = 0;
            if (request.RiderIds != null && request.RiderIds.Any())
            {
                foreach (var riderId in request.RiderIds)
                {
                    var rider = await _riderRepository.GetByIdAsync(Guid.Parse(riderId), cancellationToken);
                    if (rider != null && rider.ProductId == product.ProductId)
                    {
                        totalRiderPremium += rider.AdditionalPremium;
                        response.Breakdown.Add(new PremiumBreakdown
                        {
                            Item = $"Rider: {rider.NameEn}",
                            Amount = new Money { Amount = rider.AdditionalPremium, Currency = "BDT" },
                            Description = "Added rider coverage"
                        });
                    }
                }
            }

            response.RiderPremium = new Money { Amount = totalRiderPremium, Currency = "BDT" };
            response.TotalPremium = new Money { Amount = totalBasePremium + totalRiderPremium, Currency = "BDT" };

            // 6. Cache Result for 1 hour
            var cacheOptions = new DistributedCacheEntryOptions { AbsoluteExpirationRelativeToNow = TimeSpan.FromHours(1) };
            await _cache.SetStringAsync(cacheKey, JsonSerializer.Serialize(response), cacheOptions, cancellationToken);

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to calculate premium for product {ProductId}", request.ProductId);
            throw;
        }
    }
}
