namespace PoliSync.Products.Domain;

/// <summary>
/// Product rider/add-on (optional coverage enhancement).
/// Maps to 'product_riders' table in insurance_schema.
/// </summary>
public class Rider
{
    public Guid RiderId { get; private set; }
    public Guid ProductId { get; private set; }
    public string RiderName { get; private set; } = string.Empty;
    public string? Description { get; private set; }
    public long PremiumAmount { get; private set; }     // Additional premium in paisa
    public string PremiumCurrency { get; private set; } = "BDT";
    public long CoverageAmount { get; private set; }    // Additional coverage in paisa
    public string CoverageCurrency { get; private set; } = "BDT";
    public bool IsMandatory { get; private set; }
    public DateTime CreatedAt { get; private set; }
    public DateTime UpdatedAt { get; private set; }

    // Navigation
    public Product? Product { get; private set; }

    private Rider() { }

    public static Rider Create(
        Guid productId,
        string riderName,
        long premiumAmount,
        long coverageAmount,
        bool isMandatory = false,
        string? description = null)
    {
        return new Rider
        {
            RiderId = Guid.NewGuid(),
            ProductId = productId,
            RiderName = riderName,
            Description = description,
            PremiumAmount = premiumAmount,
            CoverageAmount = coverageAmount,
            IsMandatory = isMandatory,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };
    }
}
