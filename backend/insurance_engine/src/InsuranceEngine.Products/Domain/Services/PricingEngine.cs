using System;
using System.Collections.Generic;
using System.Linq;
using InsuranceEngine.SharedKernel.Domain.ValueObjects;

namespace InsuranceEngine.Products.Domain.Services;

/// <summary>
/// Domain service for premium calculation.
/// Applies pricing rules, rider surcharges, VAT (15%), and service fees.
/// All output values are in paisa (long).
/// </summary>
public class PricingEngine
{
    private const double VatRate = 0.15; // 15% VAT, Bangladesh standard
    private const long DefaultServiceFee = 0; // Configurable per product

    /// <summary>
    /// Calculate premium for a product given applicant data and selected riders.
    /// </summary>
    public PremiumCalculationResult Calculate(
        Product product,
        long sumInsuredAmount,
        int tenureMonths,
        List<Rider> selectedRiders,
        Dictionary<string, string> applicantData,
        double riskLoadingFactor = 1.0)
    {
        var breakdown = new List<PremiumBreakdownItem>();

        // Base premium
        long basePremium = product.BasePremiumAmount;
        breakdown.Add(new PremiumBreakdownItem("Base Premium", basePremium, "Base premium amount"));

        // Apply pricing rules
        double loadingFactor = 1.0;
        if (product.PricingConfig?.Rules != null)
        {
            foreach (var rule in product.PricingConfig.Rules)
            {
                var ruleResult = EvaluateRule(rule, applicantData);
                if (ruleResult.HasValue)
                {
                    loadingFactor *= ruleResult.Value;
                    breakdown.Add(new PremiumBreakdownItem(
                        $"Rule: {rule.RuleName}",
                        (long)Math.Round(basePremium * (ruleResult.Value - 1.0), MidpointRounding.AwayFromZero),
                        $"Adjustment factor: {ruleResult.Value:F4}"));
                }
            }
        }

        // Risk Loading factor
        if (Math.Abs(riskLoadingFactor - 1.0) > 0.0001)
        {
            loadingFactor *= riskLoadingFactor;
            breakdown.Add(new PremiumBreakdownItem(
                "Risk Loading",
                (long)Math.Round(basePremium * (riskLoadingFactor - 1.0), MidpointRounding.AwayFromZero),
                $"Risk loading factor: {riskLoadingFactor:F4}"));
        }

        // Age loading (from applicant_data if provided)
        if (applicantData.TryGetValue("age", out var ageStr) && int.TryParse(ageStr, out var age))
        {
            var ageFactor = GetAgeLoadingFactor(age);
            if (Math.Abs(ageFactor - 1.0) > 0.0001)
            {
                loadingFactor *= ageFactor;
                breakdown.Add(new PremiumBreakdownItem(
                    "Age Loading",
                    (long)Math.Round(basePremium * (ageFactor - 1.0), MidpointRounding.AwayFromZero),
                    $"Age {age}, factor: {ageFactor:F4}"));
            }
        }

        // Pre-Existing Conditions load
        if (applicantData.TryGetValue("pre_existing_conditions", out var hasConditionsStr) && bool.TryParse(hasConditionsStr, out var hasConditions) && hasConditions)
        {
            var pecFactor = 1.25; // 25% loading for pre-existing conditions
            loadingFactor *= pecFactor;
            breakdown.Add(new PremiumBreakdownItem(
                "Pre-Existing Conditions Load",
                (long)Math.Round(basePremium * (pecFactor - 1.0), MidpointRounding.AwayFromZero),
                $"Condition flag true, factor: {pecFactor:F4}"));
        }

        // Occupational Hazards load
        if (applicantData.TryGetValue("occupation_category", out var occupationStr))
        {
            var occFactor = GetOccupationalHazardsFactor(occupationStr);
            if (Math.Abs(occFactor - 1.0) > 0.0001)
            {
                loadingFactor *= occFactor;
                breakdown.Add(new PremiumBreakdownItem(
                    "Occupational Hazards Load",
                    (long)Math.Round(basePremium * (occFactor - 1.0), MidpointRounding.AwayFromZero),
                    $"Occupation: {occupationStr}, factor: {occFactor:F4}"));
            }
        }

        // Family Discount
        if (applicantData.TryGetValue("family_discount_eligible", out var familyDiscountStr) && bool.TryParse(familyDiscountStr, out var familyDiscountEligible) && familyDiscountEligible)
        {
            var familyFactor = 0.90; // 10% discount for family
            loadingFactor *= familyFactor;
            breakdown.Add(new PremiumBreakdownItem(
                "Family Discount",
                (long)Math.Round(basePremium * (familyFactor - 1.0), MidpointRounding.AwayFromZero),
                $"Family discount applied, factor: {familyFactor:F4}"));
        }

        // Tenure discount
        var tenureFactor = GetTenureDiscountFactor(tenureMonths);
        if (Math.Abs(tenureFactor - 1.0) > 0.0001)
        {
            loadingFactor *= tenureFactor;
            breakdown.Add(new PremiumBreakdownItem(
                "Tenure Discount",
                (long)Math.Round(basePremium * (tenureFactor - 1.0), MidpointRounding.AwayFromZero),
                $"Tenure {tenureMonths} months, factor: {tenureFactor:F4}"));
        }

        // Adjusted base premium
        long adjustedPremium = (long)Math.Round(basePremium * loadingFactor, MidpointRounding.AwayFromZero);

        // Rider surcharge
        long riderSurcharge = selectedRiders.Sum(r => r.PremiumAmount);
        if (riderSurcharge > 0)
        {
            breakdown.Add(new PremiumBreakdownItem(
                "Rider Surcharge",
                riderSurcharge,
                $"{selectedRiders.Count} rider(s) selected"));
        }

        long premium = adjustedPremium + riderSurcharge;

        // VAT (15%)
        long vat = (long)Math.Round(premium * VatRate, MidpointRounding.AwayFromZero);
        breakdown.Add(new PremiumBreakdownItem("VAT (15%)", vat, "Bangladesh standard VAT"));

        // Service fee
        long serviceFee = DefaultServiceFee;
        if (serviceFee > 0)
        {
            breakdown.Add(new PremiumBreakdownItem("Service Fee", serviceFee, "Platform service fee"));
        }

        long totalPayable = premium + vat + serviceFee;

        return new PremiumCalculationResult
        {
            BasePremium = Money.Bdt(adjustedPremium),
            RiderPremium = Money.Bdt(riderSurcharge),
            Vat = Money.Bdt(vat),
            ServiceFee = Money.Bdt(serviceFee),
            TotalPremium = Money.Bdt(totalPayable),
            Breakdown = breakdown
        };
    }

    private double? EvaluateRule(PricingRule rule, Dictionary<string, string> applicantData)
    {
        if (rule.Conditions == null || !rule.Conditions.Any())
            return null;

        foreach (var condition in rule.Conditions)
        {
            if (!applicantData.TryGetValue(condition.Field, out var fieldValue))
                return null;

            if (!EvaluateCondition(condition, fieldValue))
                return null;
        }

        // All conditions met — apply the action
        return rule.Action.Type switch
        {
            Enums.ActionType.IncreasePercentage => 1.0 + (rule.Action.Value / 100.0),
            Enums.ActionType.DecreasePercentage => 1.0 - (rule.Action.Value / 100.0),
            _ => null
        };
    }

    private static bool EvaluateCondition(RuleCondition condition, string fieldValue)
    {
        if (double.TryParse(fieldValue, out var numValue) && double.TryParse(condition.Value, out var condValue))
        {
            return condition.Operator switch
            {
                ">=" => numValue >= condValue,
                "<=" => numValue <= condValue,
                ">" => numValue > condValue,
                "<" => numValue < condValue,
                "==" or "=" => Math.Abs(numValue - condValue) < 0.001,
                _ => false
            };
        }

        return condition.Operator switch
        {
            "==" or "=" => fieldValue.Equals(condition.Value, StringComparison.OrdinalIgnoreCase),
            "in" => condition.Value.Split(',').Contains(fieldValue, StringComparer.OrdinalIgnoreCase),
            "not_in" => !condition.Value.Split(',').Contains(fieldValue, StringComparer.OrdinalIgnoreCase),
            _ => false
        };
    }

    private static double GetAgeLoadingFactor(int age)
    {
        return age switch
        {
            <= 25 => 0.90,  // Young — discount
            <= 35 => 1.00,  // Standard
            <= 45 => 1.10,  // Slight increase
            <= 55 => 1.25,  // Moderate increase
            <= 65 => 1.50,  // Significant increase
            _ => 1.75       // Senior — highest loading
        };
    }

    private static double GetTenureDiscountFactor(int tenureMonths)
    {
        return tenureMonths switch
        {
            >= 36 => 0.85,  // 3+ years — 15% discount
            >= 24 => 0.90,  // 2+ years — 10% discount
            >= 12 => 0.95,  // 1+ year — 5% discount
            _ => 1.00       // Less than 1 year — no discount
        };
    }

    private static double GetOccupationalHazardsFactor(string occupationCategory)
    {
        if (string.IsNullOrWhiteSpace(occupationCategory)) return 1.00;

        return occupationCategory.Trim().ToLowerInvariant() switch
        {
            "hazardous" or "high_risk" or "manual_labor_heavy" => 1.30,   // 30% loading for hazardous
            "moderate_risk" or "manual_labor_light" => 1.15,              // 15% loading for moderate
            "low_risk" or "desk_job" or "academic" => 1.00,               // Standard
            "safe" => 0.95,                                               // 5% discount for ultra-safe occupations
            _ => 1.00
        };
    }
}

public class PremiumCalculationResult
{
    public Money BasePremium { get; set; } = Money.Zero;
    public Money RiderPremium { get; set; } = Money.Zero;
    public Money Vat { get; set; } = Money.Zero;
    public Money ServiceFee { get; set; } = Money.Zero;
    public Money TotalPremium { get; set; } = Money.Zero;
    public List<PremiumBreakdownItem> Breakdown { get; set; } = new();
}

public record PremiumBreakdownItem(string Item, long Amount, string Description);
