using Google.Protobuf.WellKnownTypes;
using Insuretech.Quoting.Entity.V1;
using Insuretech.Common.V1;
using PoliSync.RulesEngine.Services;

namespace PoliSync.Quoting.Services;

public class PricingEngine : IPricingEngine
{
    private readonly IBusinessWorkflowService _rulesEngine;
    private readonly ILogger<PricingEngine> _logger;

    public PricingEngine(
        IBusinessWorkflowService rulesEngine,
        ILogger<PricingEngine> logger)
    {
        _rulesEngine = rulesEngine;
        _logger = logger;
    }

    public async Task<(PremiumCalculation Calculation, List<Coverage> Coverages, List<Discount> Discounts)> CalculatePremiumAsync(
        string productId,
        QuoteParameters parameters,
        CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Calculating premium for product {ProductId}", productId);

        var coverages = new List<Coverage>();
        var discounts = new List<Discount>();
        var basePremium = 0m;
        var riskAdjustment = 0m;
        var optionalCoveragesTotal = 0m;
        var discountsTotal = 0m;

        // Create base coverage
        coverages.Add(new Coverage
        {
            CoverageId = Guid.NewGuid().ToString(),
            Name = "Base Coverage",
            Description = "Standard coverage for the policy",
            IsIncluded = true,
            IsOptional = false
        });

        // Calculate base premium based on asset value and coverage type
        basePremium = CalculateBasePremium(parameters);

        // Apply risk adjustments based on parameters
        riskAdjustment = CalculateRiskAdjustment(parameters, basePremium);

        // Calculate optional coverages
        foreach (var optional in parameters.OptionalCoverages)
        {
            var coveragePremium = CalculateOptionalCoveragePremium(optional);
            optionalCoveragesTotal += coveragePremium.Amount;
            
            coverages.Add(new Coverage
            {
                CoverageId = optional.CoverageId,
                Name = optional.Name,
                Description = $"Optional coverage with limit {optional.SelectedLimit:C}",
                Limit = new Money { Amount = (long)(optional.SelectedLimit * 100), Currency = "USD" },
                Deductible = new Money { Amount = (long)(optional.SelectedDeductible * 100), Currency = "USD" },
                Premium = coveragePremium,
                IsIncluded = true,
                IsOptional = true
            });
        }

        // Try to get discounts from rules engine
        try
        {
            var discountInputs = new Dictionary<string, object>
            {
                { "assetValue", parameters.AssetValue },
                { "coverageType", parameters.CoverageType },
                { "coverageDurationMonths", parameters.CoverageDurationMonths }
            };

            var discountResult = await _rulesEngine.EvaluateWorkflowAsync(
                "DiscountCalculation",
                discountInputs,
                "QUOTE",
                productId,
                null,
                cancellationToken);

            if (discountResult.IsSuccess)
            {
                foreach (var result in discountResult.Results.Where(r => r.IsSuccess))
                {
                    if (decimal.TryParse(result.SuccessEvent, out var discountPercent))
                    {
                        var discountAmount = (basePremium + riskAdjustment + optionalCoveragesTotal) * (discountPercent / 100m);
                        discountsTotal += discountAmount;
                        
                        discounts.Add(new Discount
                        {
                            DiscountId = Guid.NewGuid().ToString(),
                            Name = result.RuleName,
                            Description = $"{discountPercent}% discount",
                            Amount = new Money { Amount = (long)(discountAmount * 100), Currency = "USD" },
                            Percentage = (double)discountPercent,
                            DiscountType = "Percentage"
                        });
                    }
                }
            }
        }
        catch (Exception ex)
        {
            _logger.LogWarning(ex, "Could not calculate discounts using rules engine");
        }

        // Calculate taxes (assuming 10%)
        var subtotal = basePremium + riskAdjustment + optionalCoveragesTotal - discountsTotal;
        var taxes = subtotal * 0.10m;

        // Calculate fees (assuming flat $50)
        var fees = 50m;

        var totalPremium = subtotal + taxes + fees;

        var calculation = new PremiumCalculation
        {
            BasePremium = new Money { Amount = (long)(basePremium * 100), Currency = "USD" },
            RiskAdjustment = new Money { Amount = (long)(riskAdjustment * 100), Currency = "USD" },
            OptionalCoveragesTotal = new Money { Amount = (long)(optionalCoveragesTotal * 100), Currency = "USD" },
            DiscountsTotal = new Money { Amount = (long)(discountsTotal * 100), Currency = "USD" },
            Taxes = new Money { Amount = (long)(taxes * 100), Currency = "USD" },
            Fees = new Money { Amount = (long)(fees * 100), Currency = "USD" },
            TotalPremium = new Money { Amount = (long)(totalPremium * 100), Currency = "USD" },
            Currency = "USD",
            Breakdown =
            {
                new PremiumBreakdown
                {
                    Category = "Base Premium",
                    Description = "Base coverage premium",
                    Amount = new Money { Amount = (long)(basePremium * 100), Currency = "USD" },
                    IsDiscount = false
                },
                new PremiumBreakdown
                {
                    Category = "Risk Adjustment",
                    Description = "Risk-based adjustment",
                    Amount = new Money { Amount = (long)(riskAdjustment * 100), Currency = "USD" },
                    IsDiscount = false
                },
                new PremiumBreakdown
                {
                    Category = "Optional Coverages",
                    Description = "Additional coverage options",
                    Amount = new Money { Amount = (long)(optionalCoveragesTotal * 100), Currency = "USD" },
                    IsDiscount = false
                },
                new PremiumBreakdown
                {
                    Category = "Discounts",
                    Description = "Applied discounts",
                    Amount = new Money { Amount = (long)(discountsTotal * 100), Currency = "USD" },
                    IsDiscount = true
                },
                new PremiumBreakdown
                {
                    Category = "Taxes",
                    Description = "Applicable taxes",
                    Amount = new Money { Amount = (long)(taxes * 100), Currency = "USD" },
                    IsDiscount = false
                },
                new PremiumBreakdown
                {
                    Category = "Fees",
                    Description = "Processing fees",
                    Amount = new Money { Amount = (long)(fees * 100), Currency = "USD" },
                    IsDiscount = false
                }
            }
        };

        _logger.LogInformation(
            "Premium calculated: Base={Base}, Total={Total}",
            basePremium, totalPremium);

        return (calculation, coverages, discounts);
    }

    private decimal CalculateBasePremium(QuoteParameters parameters)
    {
        // Simple calculation: 2% of asset value per year
        var assetValue = Convert.ToDecimal(parameters.AssetValue);
        var annualRate = 0.02m;
        var monthlyPremium = (assetValue * annualRate) / 12;
        return monthlyPremium * parameters.CoverageDurationMonths;
    }

    private decimal CalculateRiskAdjustment(QuoteParameters parameters, decimal basePremium)
    {
        var adjustment = 0m;

        // Adjust based on coverage type
        adjustment += parameters.CoverageType.ToLower() switch
        {
            "comprehensive" => basePremium * 0.15m,
            "collision" => basePremium * 0.10m,
            _ => 0m
        };

        return adjustment;
    }

    private Money CalculateOptionalCoveragePremium(OptionalCoverage optional)
    {
        // Simple calculation: 0.5% of selected limit
        var premium = Convert.ToDecimal(optional.SelectedLimit) * 0.005m;
        return new Money { Amount = (long)(premium * 100), Currency = "USD" };
    }
}
