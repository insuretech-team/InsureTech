using PoliSync.Products.Domain.Events;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Domain;

namespace PoliSync.Products.Domain;

/// <summary>
/// PricingConfig entity owned by Product. Manages pricing rules and configurations.
/// </summary>
public class PricingConfig : Entity
{
    public Guid Id { get; private set; }
    public Guid ProductId { get; private set; }
    public List<PricingRule> Rules { get; private set; } = [];
    public DateTimeOffset EffectiveFrom { get; private set; }
    public DateTimeOffset? EffectiveTo { get; private set; }
    public int Version { get; private set; }
    public DateTimeOffset CreatedAt { get; private set; }
    public DateTimeOffset UpdatedAt { get; private set; }

    // For EF Core
    protected PricingConfig() { }

    private PricingConfig(
        Guid id,
        Guid productId,
        List<PricingRule> rules,
        DateTimeOffset effectiveFrom,
        DateTimeOffset? effectiveTo)
    {
        Id = id;
        ProductId = productId;
        Rules = rules;
        EffectiveFrom = effectiveFrom;
        EffectiveTo = effectiveTo;
        Version = 1;
        CreatedAt = DateTimeOffset.UtcNow;
        UpdatedAt = DateTimeOffset.UtcNow;
    }

    /// <summary>
    /// Factory method to create a new PricingConfig with validation.
    /// </summary>
    public static Result<PricingConfig> Create(
        Guid productId,
        List<PricingRule> rules,
        DateTimeOffset effectiveFrom,
        DateTimeOffset? effectiveTo)
    {
        if (!rules.Any())
            return Result<PricingConfig>.Fail("EMPTY_RULES", "Pricing config must contain at least one rule.");

        if (effectiveTo.HasValue && effectiveTo <= effectiveFrom)
            return Result<PricingConfig>.Fail("INVALID_DATES", "Effective To must be after Effective From.");

        var config = new PricingConfig(
            Guid.NewGuid(),
            productId,
            rules,
            effectiveFrom,
            effectiveTo);

        return Result<PricingConfig>.Ok(config);
    }

    /// <summary>
    /// Update the pricing configuration rules and effective dates.
    /// </summary>
    public Result Update(
        List<PricingRule> rules,
        DateTimeOffset effectiveFrom,
        DateTimeOffset? effectiveTo,
        string updatedBy)
    {
        if (!rules.Any())
            return Result.Fail("EMPTY_RULES", "Pricing config must contain at least one rule.");

        if (effectiveTo.HasValue && effectiveTo <= effectiveFrom)
            return Result.Fail("INVALID_DATES", "Effective To must be after Effective From.");

        Rules = rules;
        EffectiveFrom = effectiveFrom;
        EffectiveTo = effectiveTo;
        Version++;
        UpdatedAt = DateTimeOffset.UtcNow;

        RaiseDomainEvent(new PricingUpdatedDomainEvent(
            PricingConfigId: Id,
            ProductId: ProductId,
            TenantId: Guid.Empty, // Will be set by application layer
            Version: Version,
            UpdatedBy: updatedBy));

        return Result.Ok();
    }
}
