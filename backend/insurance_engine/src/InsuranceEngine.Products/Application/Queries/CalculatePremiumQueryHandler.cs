using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Products.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
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
        // 0. Extract age from ApplicantData if available
        int age = 30; // Default
        if (request.ApplicantData.TryGetValue("age", out var ageStr) && int.TryParse(ageStr, out var parsedAge))
        {
            age = parsedAge;
        }

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

            // 4. Calculate Premium (FR-024)
            // Implementation detail: (BasePremium * Tenure) + (SumAssured * 0.01) adjusted by age
            decimal rateMultiplier = 1.0m;
            if (age > 50) rateMultiplier += 0.2m; // 20% increase for senior citizens
            if (request.TenureMonths > 12) rateMultiplier -= 0.05m; // 5% discount for long term

            long basePremiumAmount = (long)(product.BasePremium * request.TenureMonths * rateMultiplier);
            
            // For complex sum-assured based premium
            long coveragePremium = (long)(request.SumInsured.Amount * 0.01m * rateMultiplier); 
            
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
                Description = "Premium proportional to selected sum insured"
            });

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
