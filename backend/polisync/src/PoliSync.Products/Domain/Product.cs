using PoliSync.Products.Domain.Events;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Domain;

namespace PoliSync.Products.Domain;

/// <summary>
/// Product aggregate root representing an insurance product.
/// </summary>
public class Product : Entity
{
    public Guid Id { get; private set; }
    public Guid TenantId { get; private set; }
    public Guid PartnerId { get; private set; }
    public string ProductCode { get; private set; } = string.Empty;
    public string ProductName { get; private set; } = string.Empty;
    public string Description { get; private set; } = string.Empty;
    public string Category { get; private set; } = string.Empty; // ProductCategory as string
    public string Status { get; private set; } = string.Empty; // ProductStatus as string
    public long BasePremiumPaisa { get; private set; }
    public string Currency { get; private set; } = "BDT";
    public long SumInsuredMinPaisa { get; private set; }
    public long SumInsuredMaxPaisa { get; private set; }
    public int MinTenureMonths { get; private set; }
    public int MaxTenureMonths { get; private set; }
    public List<string> Exclusions { get; private set; } = [];
    public int Version { get; private set; }
    public DateTimeOffset CreatedAt { get; private set; }
    public DateTimeOffset UpdatedAt { get; private set; }
    public string CreatedBy { get; private set; } = string.Empty;

    // Navigation properties
    private readonly List<ProductPlan> _plans = [];
    private readonly List<Rider> _riders = [];
    public IReadOnlyList<ProductPlan> Plans => _plans.AsReadOnly();
    public IReadOnlyList<Rider> Riders => _riders.AsReadOnly();
    public PricingConfig? PricingConfig { get; private set; }

    // For EF Core
    protected Product() { }

    private Product(
        Guid id,
        Guid tenantId,
        Guid partnerId,
        string productCode,
        string productName,
        string description,
        string category,
        long basePremiumPaisa,
        long sumInsuredMinPaisa,
        long sumInsuredMaxPaisa,
        int minTenureMonths,
        int maxTenureMonths,
        List<string> exclusions,
        string currency,
        string createdBy)
    {
        Id = id;
        TenantId = tenantId;
        PartnerId = partnerId;
        ProductCode = productCode;
        ProductName = productName;
        Description = description;
        Category = category;
        Status = "Draft";
        BasePremiumPaisa = basePremiumPaisa;
        SumInsuredMinPaisa = sumInsuredMinPaisa;
        SumInsuredMaxPaisa = sumInsuredMaxPaisa;
        MinTenureMonths = minTenureMonths;
        MaxTenureMonths = maxTenureMonths;
        Exclusions = exclusions;
        Currency = currency;
        Version = 1;
        CreatedAt = DateTimeOffset.UtcNow;
        UpdatedAt = DateTimeOffset.UtcNow;
        CreatedBy = createdBy;
    }

    /// <summary>
    /// Factory method to create a new Product with business rule validation.
    /// </summary>
    public static Result<Product> Create(
        Guid tenantId,
        Guid partnerId,
        string productCode,
        string productName,
        string description,
        string category,
        long basePremiumPaisa,
        long sumInsuredMinPaisa,
        long sumInsuredMaxPaisa,
        int minTenureMonths,
        int maxTenureMonths,
        List<string> exclusions,
        string currency,
        string createdBy)
    {
        if (string.IsNullOrWhiteSpace(productCode))
            return Result<Product>.Fail("INVALID_CODE", "Product code cannot be empty.");

        if (string.IsNullOrWhiteSpace(productName))
            return Result<Product>.Fail("INVALID_NAME", "Product name cannot be empty.");

        if (basePremiumPaisa <= 0)
            return Result<Product>.Fail("INVALID_PREMIUM", "Base premium must be greater than zero.");

        if (sumInsuredMinPaisa > sumInsuredMaxPaisa)
            return Result<Product>.Fail("INVALID_SUM_INSURED", "Minimum sum insured cannot exceed maximum.");

        if (minTenureMonths < 1)
            return Result<Product>.Fail("INVALID_TENURE", "Minimum tenure must be at least 1 month.");

        if (maxTenureMonths < minTenureMonths)
            return Result<Product>.Fail("INVALID_TENURE", "Maximum tenure must be >= minimum tenure.");

        // TRAVEL category specific validation
        if (category == "Travel" && (minTenureMonths < 1 || minTenureMonths > 12 || maxTenureMonths > 12))
            return Result<Product>.Fail("INVALID_TENURE", "Travel products must have tenure between 1-12 months.");

        var product = new Product(
            Guid.NewGuid(),
            tenantId,
            partnerId,
            productCode,
            productName,
            description,
            category,
            basePremiumPaisa,
            sumInsuredMinPaisa,
            sumInsuredMaxPaisa,
            minTenureMonths,
            maxTenureMonths,
            exclusions,
            currency,
            createdBy);

        product.RaiseDomainEvent(new ProductCreatedDomainEvent(
            ProductId: product.Id,
            TenantId: product.TenantId,
            PartnerId: product.PartnerId,
            ProductCode: product.ProductCode,
            ProductName: product.ProductName,
            Category: product.Category,
            BasePremiumPaisa: product.BasePremiumPaisa,
            CreatedBy: product.CreatedBy));

        return Result<Product>.Ok(product);
    }

    /// <summary>
    /// Activate product. Only allowed from DRAFT or INACTIVE status.
    /// </summary>
    public Result Activate()
    {
        if (Status != "Draft" && Status != "Inactive")
            return Result.Fail("INVALID_STATE", $"Cannot activate product in {Status} status.");

        Status = "Active";
        UpdatedAt = DateTimeOffset.UtcNow;

        RaiseDomainEvent(new ProductActivatedDomainEvent(
            ProductId: Id,
            TenantId: TenantId,
            ProductCode: ProductCode,
            ActivatedBy: CreatedBy));

        return Result.Ok();
    }

    /// <summary>
    /// Deactivate product. Only allowed from ACTIVE status.
    /// </summary>
    public Result Deactivate(string reason)
    {
        if (Status != "Active")
            return Result.Fail("INVALID_STATE", $"Cannot deactivate product in {Status} status.");

        if (string.IsNullOrWhiteSpace(reason))
            return Result.Fail("INVALID_REASON", "Deactivation reason cannot be empty.");

        Status = "Inactive";
        UpdatedAt = DateTimeOffset.UtcNow;

        RaiseDomainEvent(new ProductDeactivatedDomainEvent(
            ProductId: Id,
            TenantId: TenantId,
            ProductCode: ProductCode,
            Reason: reason,
            DeactivatedBy: CreatedBy));

        return Result.Ok();
    }

    /// <summary>
    /// Discontinue product. Allowed from ACTIVE or INACTIVE status.
    /// </summary>
    public Result Discontinue(string reason)
    {
        if (Status != "Active" && Status != "Inactive")
            return Result.Fail("INVALID_STATE", $"Cannot discontinue product in {Status} status.");

        if (string.IsNullOrWhiteSpace(reason))
            return Result.Fail("INVALID_REASON", "Discontinuation reason cannot be empty.");

        Status = "Discontinued";
        UpdatedAt = DateTimeOffset.UtcNow;

        RaiseDomainEvent(new ProductDiscontinuedDomainEvent(
            ProductId: Id,
            TenantId: TenantId,
            ProductCode: ProductCode,
            Reason: reason,
            DiscontinuedBy: CreatedBy));

        return Result.Ok();
    }

    /// <summary>
    /// Update product. Only allowed in DRAFT status.
    /// </summary>
    public Result Update(
        string name,
        string description,
        long basePremiumPaisa,
        long sumInsuredMinPaisa,
        long sumInsuredMaxPaisa,
        List<string> exclusions,
        string updatedBy)
    {
        if (Status != "Draft")
            return Result.Fail("INVALID_STATE", "Product can only be updated in DRAFT status.");

        if (string.IsNullOrWhiteSpace(name))
            return Result.Fail("INVALID_NAME", "Product name cannot be empty.");

        if (basePremiumPaisa <= 0)
            return Result.Fail("INVALID_PREMIUM", "Base premium must be greater than zero.");

        if (sumInsuredMinPaisa > sumInsuredMaxPaisa)
            return Result.Fail("INVALID_SUM_INSURED", "Minimum sum insured cannot exceed maximum.");

        ProductName = name;
        Description = description;
        BasePremiumPaisa = basePremiumPaisa;
        SumInsuredMinPaisa = sumInsuredMinPaisa;
        SumInsuredMaxPaisa = sumInsuredMaxPaisa;
        Exclusions = exclusions;
        Version++;
        UpdatedAt = DateTimeOffset.UtcNow;

        RaiseDomainEvent(new ProductUpdatedDomainEvent(
            ProductId: Id,
            TenantId: TenantId,
            ProductCode: ProductCode,
            ProductName: ProductName,
            Version: Version,
            UpdatedBy: updatedBy));

        return Result.Ok();
    }

    /// <summary>
    /// Validate if amount is within sum insured range.
    /// </summary>
    public bool ValidateSumInsured(long amount)
        => amount >= SumInsuredMinPaisa && amount <= SumInsuredMaxPaisa;
}
