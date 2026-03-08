using Google.Protobuf.WellKnownTypes;
using PoliSync.Products.Domain;

namespace PoliSync.Products.Application.Mappers;

/// <summary>
/// Maps domain models to protobuf messages and vice versa.
/// </summary>
public static class ProductMapper
{
    /// <summary>
    /// Maps a domain Product to protobuf Product message.
    /// </summary>
    public static Insuretech.Products.Entity.V1.Product ToProto(this Product domain)
    {
        var protoProduct = new Insuretech.Products.Entity.V1.Product
        {
            Id = domain.Id.ToString(),
            TenantId = domain.TenantId.ToString(),
            ProductCode = domain.ProductCode,
            ProductName = domain.ProductName,
            Description = domain.Description ?? string.Empty,
            Category = ParseProductCategory(domain.Category),
            BasePremium = ToProtoMoney(domain.BasePremiumPaisa, domain.Currency),
            SumInsuredMin = ToProtoMoney(domain.SumInsuredMinPaisa, domain.Currency),
            SumInsuredMax = ToProtoMoney(domain.SumInsuredMaxPaisa, domain.Currency),
            MinTenureMonths = domain.MinTenureMonths,
            MaxTenureMonths = domain.MaxTenureMonths,
            Status = ParseProductStatus(domain.Status),
            Currency = domain.Currency,
            CreatedAt = domain.CreatedAt.ToTimestamp(),
            UpdatedAt = domain.UpdatedAt.ToTimestamp(),
        };

        if (domain.PartnerId.HasValue)
        {
            protoProduct.PartnerId = domain.PartnerId.Value.ToString();
        }

        if (domain.Exclusions != null && domain.Exclusions.Count > 0)
        {
            protoProduct.Exclusions.AddRange(domain.Exclusions);
        }

        return protoProduct;
    }

    /// <summary>
    /// Maps a domain ProductPlan to protobuf ProductPlan message.
    /// </summary>
    public static Insuretech.Products.Entity.V1.ProductPlan ToProto(this ProductPlan domain)
    {
        var protoPlan = new Insuretech.Products.Entity.V1.ProductPlan
        {
            Id = domain.Id.ToString(),
            ProductId = domain.ProductId.ToString(),
            PlanCode = domain.PlanCode,
            PlanName = domain.PlanName,
            Description = domain.Description ?? string.Empty,
            BasePremium = ToProtoMoney(domain.BasePremiumPaisa, domain.Currency),
            SumInsured = ToProtoMoney(domain.SumInsuredPaisa, domain.Currency),
            Currency = domain.Currency,
            CreatedAt = domain.CreatedAt.ToTimestamp(),
        };

        if (domain.Features != null && domain.Features.Count > 0)
        {
            protoPlan.Features.AddRange(domain.Features);
        }

        return protoPlan;
    }

    /// <summary>
    /// Maps a domain Rider to protobuf Rider message.
    /// </summary>
    public static Insuretech.Products.Entity.V1.Rider ToProto(this Rider domain)
    {
        var protoRider = new Insuretech.Products.Entity.V1.Rider
        {
            Id = domain.Id.ToString(),
            ProductId = domain.ProductId.ToString(),
            RiderCode = domain.RiderCode,
            RiderName = domain.RiderName,
            Description = domain.Description ?? string.Empty,
            PremiumAmount = ToProtoMoney(domain.PremiumAmountPaisa, domain.Currency),
            SumInsured = ToProtoMoney(domain.SumInsuredPaisa, domain.Currency),
            Category = domain.Category,
            IsMandatory = domain.IsMandatory,
            Currency = domain.Currency,
            CreatedAt = domain.CreatedAt.ToTimestamp(),
        };

        return protoRider;
    }

    /// <summary>
    /// Maps a domain PricingConfig to protobuf PricingConfig message.
    /// </summary>
    public static Insuretech.Products.Entity.V1.PricingConfig ToProto(this PricingConfig domain)
    {
        var protoConfig = new Insuretech.Products.Entity.V1.PricingConfig
        {
            Id = domain.Id.ToString(),
            ProductId = domain.ProductId.ToString(),
            EffectiveFrom = domain.EffectiveFrom.ToTimestamp(),
            CreatedAt = domain.CreatedAt.ToTimestamp(),
        };

        if (domain.EffectiveTo.HasValue)
        {
            protoConfig.EffectiveTo = domain.EffectiveTo.Value.ToTimestamp();
        }

        foreach (var rule in domain.Rules)
        {
            protoConfig.Rules.Add(rule.ToProto());
        }

        return protoConfig;
    }

    /// <summary>
    /// Maps a domain PricingRule to protobuf PricingRule message.
    /// </summary>
    public static Insuretech.Products.Entity.V1.PricingRule ToProto(this PricingRule domain)
    {
        var protoRule = new Insuretech.Products.Entity.V1.PricingRule
        {
            Id = domain.Id.ToString(),
            Name = domain.Name,
            Type = domain.Type,
            Priority = domain.Priority,
            ApplyAll = domain.ApplyAll,
        };

        foreach (var condition in domain.Conditions)
        {
            protoRule.Conditions.Add(condition.ToProto());
        }

        protoRule.Action = domain.Action.ToProto();

        return protoRule;
    }

    /// <summary>
    /// Maps a domain RuleCondition to protobuf RuleCondition message.
    /// </summary>
    public static Insuretech.Products.Entity.V1.RuleCondition ToProto(this RuleCondition domain)
    {
        return new Insuretech.Products.Entity.V1.RuleCondition
        {
            Field = domain.Field,
            Operator = domain.Operator,
            Value = domain.Value
        };
    }

    /// <summary>
    /// Maps a domain RuleAction to protobuf RuleAction message.
    /// </summary>
    public static Insuretech.Products.Entity.V1.RuleAction ToProto(this RuleAction domain)
    {
        return new Insuretech.Products.Entity.V1.RuleAction
        {
            ActionType = domain.ActionType,
            Value = (double)domain.Value
        };
    }

    /// <summary>
    /// Maps a PricingRuleDto to domain PricingRule.
    /// </summary>
    public static PricingRule MapPricingRuleDtoToDomain(Commands.PricingRuleDto dto)
    {
        var conditions = dto.Conditions
            .Select(c => new RuleCondition(c.Field, c.Operator, c.Value))
            .ToList();

        var action = new RuleAction(dto.Action.ActionType, dto.Action.Value);

        return new PricingRule(
            Guid.Parse(dto.RuleId),
            dto.RuleName,
            dto.RuleType,
            dto.Priority,
            dto.ApplyAll,
            conditions,
            action
        );
    }

    /// <summary>
    /// Converts paisa amount to protobuf Money message.
    /// </summary>
    private static Insuretech.Common.V1.Money ToProtoMoney(long paisaAmount, string currency)
    {
        return new Insuretech.Common.V1.Money
        {
            Amount = paisaAmount,
            Currency = currency
        };
    }

    /// <summary>
    /// Parses string category to protobuf enum.
    /// </summary>
    private static Insuretech.Products.Entity.V1.ProductCategory ParseProductCategory(string category)
    {
        return category.ToUpper() switch
        {
            "HEALTH" => Insuretech.Products.Entity.V1.ProductCategory.Health,
            "LIFE" => Insuretech.Products.Entity.V1.ProductCategory.Life,
            "MOTOR" => Insuretech.Products.Entity.V1.ProductCategory.Motor,
            "HOME" => Insuretech.Products.Entity.V1.ProductCategory.Home,
            "TRAVEL" => Insuretech.Products.Entity.V1.ProductCategory.Travel,
            "GENERAL" => Insuretech.Products.Entity.V1.ProductCategory.General,
            _ => Insuretech.Products.Entity.V1.ProductCategory.Unspecified
        };
    }

    /// <summary>
    /// Parses string status to protobuf enum.
    /// </summary>
    private static Insuretech.Products.Entity.V1.ProductStatus ParseProductStatus(string status)
    {
        return status.ToUpper() switch
        {
            "DRAFT" => Insuretech.Products.Entity.V1.ProductStatus.Draft,
            "ACTIVE" => Insuretech.Products.Entity.V1.ProductStatus.Active,
            "INACTIVE" => Insuretech.Products.Entity.V1.ProductStatus.Inactive,
            "DISCONTINUED" => Insuretech.Products.Entity.V1.ProductStatus.Discontinued,
            _ => Insuretech.Products.Entity.V1.ProductStatus.Unspecified
        };
    }
}
