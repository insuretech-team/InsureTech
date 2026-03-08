using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Domain;

namespace PoliSync.Products.Domain;

/// <summary>
/// Rider entity owned by Product. Represents an optional or mandatory add-on to a product.
/// </summary>
public class Rider : Entity
{
    public Guid Id { get; private set; }
    public Guid ProductId { get; private set; }
    public string RiderCode { get; private set; } = string.Empty;
    public string RiderName { get; private set; } = string.Empty;
    public string Description { get; private set; } = string.Empty;
    public long PremiumAmountPaisa { get; private set; }
    public long SumInsuredPaisa { get; private set; }
    public string Currency { get; private set; } = "BDT";
    public string Category { get; private set; } = string.Empty;
    public bool IsMandatory { get; private set; }
    public bool IsActive { get; private set; }
    public DateTimeOffset CreatedAt { get; private set; }
    public DateTimeOffset UpdatedAt { get; private set; }

    // For EF Core
    protected Rider() { }

    private Rider(
        Guid id,
        Guid productId,
        string riderCode,
        string riderName,
        string description,
        long premiumAmountPaisa,
        long sumInsuredPaisa,
        string category,
        bool isMandatory,
        string currency)
    {
        Id = id;
        ProductId = productId;
        RiderCode = riderCode;
        RiderName = riderName;
        Description = description;
        PremiumAmountPaisa = premiumAmountPaisa;
        SumInsuredPaisa = sumInsuredPaisa;
        Category = category;
        IsMandatory = isMandatory;
        Currency = currency;
        IsActive = true;
        CreatedAt = DateTimeOffset.UtcNow;
        UpdatedAt = DateTimeOffset.UtcNow;
    }

    /// <summary>
    /// Factory method to create a new Rider with validation.
    /// </summary>
    public static Result<Rider> Create(
        Guid productId,
        string riderCode,
        string riderName,
        string description,
        long premiumAmountPaisa,
        long sumInsuredPaisa,
        string category,
        bool isMandatory,
        string currency)
    {
        if (string.IsNullOrWhiteSpace(riderCode))
            return Result<Rider>.Fail("INVALID_CODE", "Rider code cannot be empty.");

        if (string.IsNullOrWhiteSpace(riderName))
            return Result<Rider>.Fail("INVALID_NAME", "Rider name cannot be empty.");

        if (premiumAmountPaisa < 0)
            return Result<Rider>.Fail("INVALID_PREMIUM", "Rider premium cannot be negative.");

        if (sumInsuredPaisa <= 0)
            return Result<Rider>.Fail("INVALID_SUM_INSURED", "Rider sum insured must be greater than zero.");

        var rider = new Rider(
            Guid.NewGuid(),
            productId,
            riderCode,
            riderName,
            description,
            premiumAmountPaisa,
            sumInsuredPaisa,
            category,
            isMandatory,
            currency);

        return Result<Rider>.Ok(rider);
    }
}
