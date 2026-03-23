using System;
using System.Collections.Generic;
using InsuranceEngine.Products.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Domain.ValueObjects;

namespace InsuranceEngine.Products.Domain;

public class Product
{
    public Guid Id { get; set; }
    public string ProductCode { get; set; } = string.Empty;
    public string ProductName { get; set; } = string.Empty;
    public string? Description { get; set; }
    public ProductCategory Category { get; set; }
    public ProductStatus Status { get; set; }

    // Money fields — stored as bigint (paisa)
    public long BasePremium { get; set; }
    public string BasePremiumCurrency { get; set; } = "BDT";
    public long MinSumInsured { get; set; }
    public string MinSumInsuredCurrency { get; set; } = "BDT";
    public long MaxSumInsured { get; set; }
    public string MaxSumInsuredCurrency { get; set; } = "BDT";

    public int MinTenureMonths { get; set; }
    public int MaxTenureMonths { get; set; }

    public List<string> Exclusions { get; set; } = new();
    public string? ProductAttributes { get; set; } // JSONB

    // Navigation properties
    public List<Rider> AvailableRiders { get; set; } = new();
    public List<ProductPlan> Plans { get; set; } = new();
    public List<RiskAssessmentQuestion> RiskAssessmentQuestions { get; set; } = new();
    public PricingConfig? PricingConfig { get; set; }

    // Audit
    public Guid CreatedBy { get; set; }
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    public DateTime? DeletedAt { get; set; }

    // --- Money convenience accessors ---
    public Money BasePremiumMoney
    {
        get => new(BasePremium, BasePremiumCurrency);
        set { BasePremium = value.Amount; BasePremiumCurrency = value.CurrencyCode; }
    }

    public Money MinSumInsuredMoney
    {
        get => new(MinSumInsured, MinSumInsuredCurrency);
        set { MinSumInsured = value.Amount; MinSumInsuredCurrency = value.CurrencyCode; }
    }

    public Money MaxSumInsuredMoney
    {
        get => new(MaxSumInsured, MaxSumInsuredCurrency);
        set { MaxSumInsured = value.Amount; MaxSumInsuredCurrency = value.CurrencyCode; }
    }

    // --- Status transition methods ---

    public Result Activate()
    {
        if (Status != ProductStatus.Draft && Status != ProductStatus.Inactive)
            return Result.Fail(Error.InvalidStateTransition(
                $"Cannot activate product in '{Status}' status. Only DRAFT or INACTIVE products can be activated."));

        Status = ProductStatus.Active;
        UpdatedAt = DateTime.UtcNow;
        return Result.Ok();
    }

    public Result Deactivate()
    {
        if (Status != ProductStatus.Active)
            return Result.Fail(Error.InvalidStateTransition(
                $"Cannot deactivate product in '{Status}' status. Only ACTIVE products can be deactivated."));

        Status = ProductStatus.Inactive;
        UpdatedAt = DateTime.UtcNow;
        return Result.Ok();
    }

    public Result Discontinue()
    {
        if (Status == ProductStatus.Discontinued)
            return Result.Fail(Error.InvalidStateTransition(
                "Product is already discontinued."));

        Status = ProductStatus.Discontinued;
        UpdatedAt = DateTime.UtcNow;
        return Result.Ok();
    }
}
