using System.Diagnostics;
using System.Text.Json;
using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Life.Entity.V1;

namespace PoliSync.LifeInsurance.Services;

public class LifePremiumCalculator : ILifePremiumCalculator
{
    private readonly ILogger<LifePremiumCalculator> _logger;

    public LifePremiumCalculator(ILogger<LifePremiumCalculator> logger)
    {
        _logger = logger;
    }

    public Task<(long BasePremium, long AgeAddition, float ConditionMultiplier, long ConditionAddition, long BonusDiscount, long TotalPremium,
          List<PremiumBreakdown> Breakdown, List<string> AppliedConditions, List<string> AppliedBonuses)>
        CalculatePremiumAsync(
            LifeProduct product,
            InsuredPerson insuredPerson,
            int ageAtEntry,
            int policyTermYears,
            long sumAssured,
            List<string> bonusCodes,
            CancellationToken cancellationToken = default)
    {
        var stopwatch = Stopwatch.StartNew();

        _logger.LogInformation(
            "Calculating life premium for product {ProductCode}, age {Age}, term {Term} years, sum assured {SumAssured}",
            product.ProductCode, ageAtEntry, policyTermYears, sumAssured);

        var breakdown = new List<PremiumBreakdown>();
        var appliedConditions = new List<string>();
        var appliedBonuses = new List<string>();

        // 1. Calculate Base Premium (base_rate * policy_term)
        var basePremium = product.BaseRate * policyTermYears;
        breakdown.Add(new PremiumBreakdown
        {
            Component = "Base Premium",
            Amount = basePremium,
            Description = $"Base rate {product.BaseRate} × {policyTermYears} years",
            IsDiscount = false
        });

        // 2. Calculate Age Addition
        // Formula: ((age - startAge) / ageIncrement) * priceToAdd
        var ageConfig = product.AgeAdditionConfig;
        long ageAddition = 0;
        
        if (ageAtEntry > ageConfig.StartAge)
        {
            var ageDiff = ageAtEntry - ageConfig.StartAge;
            var incrementCount = (int)Math.Floor((double)ageDiff / ageConfig.AgeIncrement);
            ageAddition = incrementCount * ageConfig.PriceToAdd * policyTermYears;
        }
        
        breakdown.Add(new PremiumBreakdown
        {
            Component = "Age Addition",
            Amount = ageAddition,
            Description = $"Age {ageAtEntry} - addition based on age brackets",
            IsDiscount = false
        });

        // 3. Calculate Condition Multiplier
        var conditionMultipliers = JsonSerializer.Deserialize<List<ConditionMultiplier>>(product.ConditionMultipliersJson) ?? new List<ConditionMultiplier>();
        float totalConditionMultiplier = 1.0f;
        
        foreach (var condition in insuredPerson.HealthConditions)
        {
            var matchingMultiplier = conditionMultipliers.FirstOrDefault(cm => 
                cm.ConditionCode.Equals(condition.ConditionCode, StringComparison.OrdinalIgnoreCase));
            
            if (matchingMultiplier != null)
            {
                totalConditionMultiplier += matchingMultiplier.Multiplier;
                appliedConditions.Add(condition.ConditionName);
            }
        }

        // Calculate after age addition
        var afterAge = basePremium + ageAddition;
        var conditionAddition = (long)((totalConditionMultiplier - 1.0f) * afterAge);
        
        breakdown.Add(new PremiumBreakdown
        {
            Component = "Condition Addition",
            Amount = conditionAddition,
            Description = $"Health conditions multiplier: {totalConditionMultiplier:F2}",
            IsDiscount = false
        });

        // 4. Calculate Bonus/Discount
        var bonusConfigs = JsonSerializer.Deserialize<List<BonusConfig>>(product.BonusConfigJson) ?? new List<BonusConfig>();
        long bonusDiscount = 0;
        var subtotal = afterAge + conditionAddition;

        foreach (var bonusCode in bonusCodes)
        {
            var bonus = bonusConfigs.FirstOrDefault(b => 
                b.BonusCode.Equals(bonusCode, StringComparison.OrdinalIgnoreCase));
            
            if (bonus != null)
            {
                long bonusAmount = 0;
                
                if (bonus.BonusType == "PERCENTAGE")
                {
                    bonusAmount = (long)(subtotal * bonus.Percentage);
                }
                else if (bonus.BonusType == "FIXED_AMOUNT")
                {
                    bonusAmount = bonus.FixedAmount;
                }
                
                bonusDiscount += bonusAmount;
                appliedBonuses.Add(bonus.BonusName);
                
                breakdown.Add(new PremiumBreakdown
                {
                    Component = bonus.BonusName,
                    Amount = bonusAmount,
                    Description = bonus.Description,
                    IsDiscount = true
                });
            }
        }

        // 5. Calculate Total
        var totalPremium = subtotal + conditionAddition - bonusDiscount;

        stopwatch.Stop();

        _logger.LogInformation(
            "Life premium calculated: Base={Base}, AgeAdd={AgeAdd}, ConditionMult={CondMult}, ConditionAdd={CondAdd}, Bonus={Bonus}, Total={Total}, Time={Ms}ms",
            basePremium, ageAddition, totalConditionMultiplier, conditionAddition, bonusDiscount, totalPremium, stopwatch.ElapsedMilliseconds);

        return Task.FromResult((basePremium, ageAddition, totalConditionMultiplier, conditionAddition, bonusDiscount, totalPremium,
            breakdown, appliedConditions, appliedBonuses));
    }
}
