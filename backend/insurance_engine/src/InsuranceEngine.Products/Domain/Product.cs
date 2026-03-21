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
    public string? ProductNameBn { get; set; }
    public string? Description { get; set; }
    public string? DescriptionBn { get; set; }
    public ProductCategory Category { get; set; }
    public ProductStatus Status { get; set; }

    // Money fields — stored as bigint (paisa)
    public long BasePremiumAmount { get; set; }
    public string BasePremiumCurrency { get; set; } = "BDT";
    public long MinSumInsuredAmount { get; set; }
    public string MinSumInsuredCurrency { get; set; } = "BDT";
    public long MaxSumInsuredAmount { get; set; }
    public string MaxSumInsuredCurrency { get; set; } = "BDT";

    public int MinAge { get; set; }
    public int MaxAge { get; set; }
    public int MinTenureMonths { get; set; }
    public int MaxTenureMonths { get; set; }

    public List<string> Exclusions { get; set; } = new();
    public string? ProductAttributes { get; set; } // JSONB

    // Co-pay / Deductible configuration (FR-100/FR-104)
    public double DeductiblePercentage { get; set; }  // 0-100
    public double CoPayPercentage { get; set; }       // 0-100
    public long MaxDeductibleAmount { get; set; }     // paisa, 0 = no cap

    // Navigation properties
    public List<Rider> AvailableRiders { get; set; } = new();
    public List<ProductPlan> Plans { get; set; } = new();
    public List<RiskAssessmentQuestion> Questions { get; set; } = new();
    public PricingConfig? PricingConfig { get; set; }

    // Audit
    public Guid CreatedBy { get; set; }
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    public DateTime? DeletedAt { get; set; }
    public bool IsDeleted { get; set; }
    public Guid TenantId { get; set; }

    // --- Money convenience accessors ---
    public Money BasePremium
    {
        get => new(BasePremiumAmount, BasePremiumCurrency);
        set { BasePremiumAmount = value.Amount; BasePremiumCurrency = value.CurrencyCode; }
    }

    public Money MinSumInsured
    {
        get => new(MinSumInsuredAmount, MinSumInsuredCurrency);
        set { MinSumInsuredAmount = value.Amount; MinSumInsuredCurrency = value.CurrencyCode; }
    }

    public Money MaxSumInsured
    {
        get => new(MaxSumInsuredAmount, MaxSumInsuredCurrency);
        set { MaxSumInsuredAmount = value.Amount; MaxSumInsuredCurrency = value.CurrencyCode; }
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
