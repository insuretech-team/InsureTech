using System;
using System.Collections.Generic;
using System.Linq;
using InsuranceEngine.Products.Domain;
using InsuranceEngine.Products.Domain.Enums;

namespace InsuranceEngine.Products.Application.DTOs;

/// <summary>
/// Maps between domain entities and DTOs
/// </summary>
public static class ProductMappings
{
    public static ProductDto ToDto(this Product product)
    {
        return new ProductDto(
            Id: product.Id,
            ProductCode: product.ProductCode,
            ProductName: product.ProductName,
            ProductNameBn: product.ProductNameBn,
            Description: product.Description,
            Category: product.Category,
            Status: product.Status,
            BasePremium: new MoneyDto(product.BasePremiumAmount, product.BasePremiumCurrency),
            MinSumInsured: new MoneyDto(product.MinSumInsuredAmount, product.MinSumInsuredCurrency),
            MaxSumInsured: new MoneyDto(product.MaxSumInsuredAmount, product.MaxSumInsuredCurrency),
            MinAge: product.MinAge,
            MaxAge: product.MaxAge,
            MinTenureMonths: product.MinTenureMonths,
            MaxTenureMonths: product.MaxTenureMonths,
            Exclusions: product.Exclusions,
            AvailableRiders: product.AvailableRiders?.Select(r => r.ToDto()).ToList(),
            Plans: product.Plans?.Select(p => p.ToDto()).ToList(),
            PricingConfig: product.PricingConfig?.ToDto(),
            CreatedBy: product.CreatedBy,
            CreatedAt: product.CreatedAt,
            UpdatedAt: product.UpdatedAt
        );
    }

    public static ProductListDto ToListDto(this Product product)
    {
        return new ProductListDto(
            Id: product.Id,
            ProductCode: product.ProductCode,
            ProductName: product.ProductName,
            Category: product.Category,
            Status: product.Status,
            BasePremium: new MoneyDto(product.BasePremiumAmount, product.BasePremiumCurrency),
            MinSumInsured: new MoneyDto(product.MinSumInsuredAmount, product.MinSumInsuredCurrency),
            MaxSumInsured: new MoneyDto(product.MaxSumInsuredAmount, product.MaxSumInsuredCurrency)
        );
    }

    public static RiderDto ToDto(this Rider rider)
    {
        return new RiderDto(
            Id: rider.Id,
            RiderName: rider.RiderName,
            Description: rider.Description,
            PremiumAmount: new MoneyDto(rider.PremiumAmount, rider.PremiumCurrency),
            CoverageAmount: new MoneyDto(rider.CoverageAmount, rider.CoverageCurrency),
            IsMandatory: rider.IsMandatory
        );
    }

    public static ProductPlanDto ToDto(this ProductPlan plan)
    {
        return new ProductPlanDto(
            Id: plan.Id,
            PlanName: plan.PlanName,
            PlanDescription: plan.PlanDescription,
            PremiumAmount: new MoneyDto(plan.PremiumAmount, plan.PremiumCurrency),
            MinSumInsured: new MoneyDto(plan.MinSumInsuredAmount, plan.MinSumInsuredCurrency),
            MaxSumInsured: new MoneyDto(plan.MaxSumInsuredAmount, plan.MaxSumInsuredCurrency),
            Attributes: plan.Attributes
        );
    }

    public static PricingConfigDto ToDto(this PricingConfig config)
    {
        return new PricingConfigDto(
            Id: config.Id,
            Rules: config.Rules?.Select(r => r.ToDto()).ToList() ?? new(),
            EffectiveFrom: config.EffectiveFrom,
            EffectiveTo: config.EffectiveTo
        );
    }

    public static PricingRuleDto ToDto(this PricingRule rule)
    {
        return new PricingRuleDto(
            Id: rule.Id,
            RuleName: rule.RuleName,
            Type: rule.Type,
            Conditions: rule.Conditions?.Select(c => new RuleConditionDto(c.Field, c.Operator, c.Value)).ToList() ?? new(),
            Action: new RuleActionDto(rule.Action.Type, rule.Action.Value)
        );
    }
}
