using PoliSync.SharedKernel.Domain;

namespace PoliSync.Products.Domain;

/// <summary>
/// Product aggregate root — represents an insurance product (Motor, Health, Travel, etc.).
/// Maps to 'products' table in insurance_schema.
/// </summary>
public class Product : Entity
{
    public Guid ProductId { get; private set; }
    public string ProductCode { get; private set; } = string.Empty;  // e.g., MOT-001, HLT-001
    public string ProductName { get; private set; } = string.Empty;
    public ProductCategory Category { get; private set; }
    public string? Description { get; private set; }
    public long BasePremium { get; private set; }           // Amount in paisa
    public string BasePremiumCurrency { get; private set; } = "BDT";
    public long MinSumInsured { get; private set; }         // Minimum coverage in paisa
    public string MinSumInsuredCurrency { get; private set; } = "BDT";
    public long MaxSumInsured { get; private set; }         // Maximum coverage in paisa
    public string MaxSumInsuredCurrency { get; private set; } = "BDT";
    public int MinTenureMonths { get; private set; }
    public int MaxTenureMonths { get; private set; }
    public List<string> Exclusions { get; private set; } = [];
    public ProductStatus Status { get; private set; } = ProductStatus.Draft;
    public string? ProductAttributes { get; private set; }  // JSONB flexible attributes
    public string CreatedBy { get; private set; } = string.Empty;
    public DateTime CreatedAt { get; private set; }
    public DateTime UpdatedAt { get; private set; }
    public DateTime? DeletedAt { get; private set; }

    // Navigation properties
    public List<ProductPlan> Plans { get; private set; } = [];
    public List<Rider> AvailableRiders { get; private set; } = [];
    public PricingConfig? PricingConfig { get; private set; }

    // Private constructor for EF Core
    private Product() { }

    // Factory method for creating new products
    public static Product Create(
        string productCode,
        string productName,
        ProductCategory category,
        long basePremium,
        long minSumInsured,
        long maxSumInsured,
        int minTenureMonths,
        int maxTenureMonths,
        string createdBy,
        string? description = null,
        List<string>? exclusions = null,
        string? productAttributes = null)
    {
        var product = new Product
        {
            ProductId = Guid.NewGuid(),
            ProductCode = productCode,
            ProductName = productName,
            Category = category,
            Description = description,
            BasePremium = basePremium,
            MinSumInsured = minSumInsured,
            MaxSumInsured = maxSumInsured,
            MinTenureMonths = minTenureMonths,
            MaxTenureMonths = maxTenureMonths,
            Exclusions = exclusions ?? [],
            Status = ProductStatus.Draft,
            ProductAttributes = productAttributes,
            CreatedBy = createdBy,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        product.RaiseDomainEvent(new ProductCreatedEvent(product.ProductId, productCode, productName, category));
        return product;
    }

    public void Update(
        string? productName = null,
        string? description = null,
        long? basePremium = null,
        long? minSumInsured = null,
        long? maxSumInsured = null,
        int? minTenureMonths = null,
        int? maxTenureMonths = null,
        List<string>? exclusions = null,
        string? productAttributes = null)
    {
        if (productName is not null) ProductName = productName;
        if (description is not null) Description = description;
        if (basePremium.HasValue) BasePremium = basePremium.Value;
        if (minSumInsured.HasValue) MinSumInsured = minSumInsured.Value;
        if (maxSumInsured.HasValue) MaxSumInsured = maxSumInsured.Value;
        if (minTenureMonths.HasValue) MinTenureMonths = minTenureMonths.Value;
        if (maxTenureMonths.HasValue) MaxTenureMonths = maxTenureMonths.Value;
        if (exclusions is not null) Exclusions = exclusions;
        if (productAttributes is not null) ProductAttributes = productAttributes;

        UpdatedAt = DateTime.UtcNow;
        RaiseDomainEvent(new ProductUpdatedEvent(ProductId, ProductCode));
    }

    public void Activate()
    {
        if (Status is not (ProductStatus.Draft or ProductStatus.Inactive))
            throw new InvalidOperationException($"Cannot activate product in {Status} status");

        Status = ProductStatus.Active;
        UpdatedAt = DateTime.UtcNow;
        RaiseDomainEvent(new ProductActivatedEvent(ProductId, ProductCode, ProductName));
    }

    public void Deactivate(string? reason = null)
    {
        if (Status != ProductStatus.Active)
            throw new InvalidOperationException($"Cannot deactivate product in {Status} status");

        Status = ProductStatus.Inactive;
        UpdatedAt = DateTime.UtcNow;
        RaiseDomainEvent(new ProductDeactivatedEvent(ProductId, ProductCode, reason));
    }

    public void Discontinue(string? reason = null)
    {
        if (Status == ProductStatus.Discontinued)
            throw new InvalidOperationException("Product is already discontinued");

        Status = ProductStatus.Discontinued;
        UpdatedAt = DateTime.UtcNow;
        DeletedAt = DateTime.UtcNow; // soft delete
        RaiseDomainEvent(new ProductDiscontinuedEvent(ProductId, ProductCode, reason));
    }

    public void AddPlan(ProductPlan plan)
    {
        Plans.Add(plan);
        UpdatedAt = DateTime.UtcNow;
    }

    public void AddRider(Rider rider)
    {
        AvailableRiders.Add(rider);
        UpdatedAt = DateTime.UtcNow;
    }

    public void SetPricingConfig(PricingConfig config)
    {
        PricingConfig = config;
        UpdatedAt = DateTime.UtcNow;
    }
}
