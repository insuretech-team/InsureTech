using Google.Protobuf.WellKnownTypes;
using Insuretech.Common.V1;
using Insuretech.Products.Entity.V1;

namespace PoliSync.Products.Persistence;

/// <summary>
/// Maps between EF persistence records and proto-generated types.
/// Proto types are the canonical domain representation (proto-first).
/// EF records are the DB storage representation.
/// </summary>
public static class ProductMapper
{
    // ── Record → Proto ─────────────────────────────────────────────────────

    public static Product ToProto(this ProductRecord r)
    {
        var p = new Product
        {
            ProductId       = r.ProductId.ToString(),
            ProductCode     = r.ProductCode,
            ProductName     = r.ProductName,
            Category        = Enum.TryParse<ProductCategory>(r.Category, out var cat) ? cat : ProductCategory.Unspecified,
            Description     = r.Description ?? string.Empty,
            BasePremium     = new Money { Amount = r.BasePremium, Currency = r.BasePremiumCurrency },
            MinSumInsured   = new Money { Amount = r.MinSumInsured, Currency = r.MinSumInsuredCurrency },
            MaxSumInsured   = new Money { Amount = r.MaxSumInsured, Currency = r.MaxSumInsuredCurrency },
            MinTenureMonths = r.MinTenureMonths,
            MaxTenureMonths = r.MaxTenureMonths,
            Status          = Enum.TryParse<ProductStatus>(r.Status, out var st) ? st : ProductStatus.Unspecified,
            CreatedAt       = Timestamp.FromDateTime(r.CreatedAt),
            UpdatedAt       = Timestamp.FromDateTime(r.UpdatedAt),
            CreatedBy       = r.CreatedBy.ToString(),
        };
        p.Exclusions.AddRange(r.Exclusions);
        p.Plans.AddRange(r.Plans.Select(pl => pl.ToProto()));
        p.AvailableRiders.AddRange(r.Riders.Select(rd => rd.ToProto()));
        return p;
    }

    public static ProductPlan ToProto(this ProductPlanRecord r) => new()
    {
        PlanId          = r.PlanId.ToString(),
        ProductId       = r.ProductId.ToString(),
        PlanName        = r.PlanName,
        PlanDescription = r.PlanDescription ?? string.Empty,
        PremiumAmount   = new Money { Amount = r.PremiumAmount, Currency = r.PremiumCurrency },
        MinSumInsured   = new Money { Amount = r.MinSumInsured, Currency = r.MinSumInsuredCurrency },
        MaxSumInsured   = new Money { Amount = r.MaxSumInsured, Currency = r.MaxSumInsuredCurrency },
        CreatedAt       = Timestamp.FromDateTime(r.CreatedAt),
        UpdatedAt       = Timestamp.FromDateTime(r.UpdatedAt),
    };

    public static Rider ToProto(this RiderRecord r) => new()
    {
        RiderId         = r.RiderId.ToString(),
        ProductId       = r.ProductId.ToString(),
        RiderName       = r.RiderName,
        Description     = r.Description ?? string.Empty,
        PremiumAmount   = new Money { Amount = r.PremiumAmount, Currency = r.PremiumCurrency },
        CoverageAmount  = new Money { Amount = r.CoverageAmount, Currency = r.CoverageCurrency },
        IsMandatory     = r.IsMandatory,
        CreatedAt       = Timestamp.FromDateTime(r.CreatedAt),
        UpdatedAt       = Timestamp.FromDateTime(r.UpdatedAt),
    };

    // ── Proto → Record ─────────────────────────────────────────────────────

    public static ProductRecord ToRecord(this Product p, Guid createdBy) => new()
    {
        ProductId              = string.IsNullOrEmpty(p.ProductId) ? Guid.NewGuid() : Guid.Parse(p.ProductId),
        ProductCode            = p.ProductCode,
        ProductName            = p.ProductName,
        Category               = p.Category.ToString(),
        Description            = p.Description,
        BasePremium            = p.BasePremium?.Amount ?? 0,
        BasePremiumCurrency    = p.BasePremium?.Currency ?? "BDT",
        MinSumInsured          = p.MinSumInsured?.Amount ?? 0,
        MinSumInsuredCurrency  = p.MinSumInsured?.Currency ?? "BDT",
        MaxSumInsured          = p.MaxSumInsured?.Amount ?? 0,
        MaxSumInsuredCurrency  = p.MaxSumInsured?.Currency ?? "BDT",
        MinTenureMonths        = p.MinTenureMonths,
        MaxTenureMonths        = p.MaxTenureMonths,
        Exclusions             = [.. p.Exclusions],
        Status                 = ProductStatus.Draft.ToString(),
        CreatedBy              = createdBy,
    };
}
