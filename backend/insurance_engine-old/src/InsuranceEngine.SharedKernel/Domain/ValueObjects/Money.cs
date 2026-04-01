namespace InsuranceEngine.SharedKernel.Domain.ValueObjects;

/// <summary>
/// Represents a monetary amount in the smallest currency unit (paisa for BDT).
/// 1 BDT = 100 paisa. All amounts are stored as long (bigint) to avoid floating-point issues.
/// </summary>
public record Money
{
    public long Amount { get; init; }
    public string CurrencyCode { get; init; } = "BDT";

    public Money() { }

    public Money(long amount, string currencyCode = "BDT")
    {
        if (amount < 0)
            throw new ArgumentException("Money amount cannot be negative", nameof(amount));

        Amount = amount;
        CurrencyCode = currencyCode;
    }

    /// <summary>
    /// Create a Money instance in BDT paisa
    /// </summary>
    public static Money Bdt(long paisa) => new(paisa, "BDT");

    /// <summary>
    /// Zero amount in BDT
    /// </summary>
    public static Money Zero => new(0, "BDT");

    /// <summary>
    /// Get decimal amount for display (amount / 100)
    /// </summary>
    public decimal DecimalAmount => Amount / 100m;

    public static Money operator +(Money a, Money b)
    {
        if (a.CurrencyCode != b.CurrencyCode)
            throw new InvalidOperationException("Cannot add money with different currencies");
        return new Money(a.Amount + b.Amount, a.CurrencyCode);
    }

    public static Money operator -(Money a, Money b)
    {
        if (a.CurrencyCode != b.CurrencyCode)
            throw new InvalidOperationException("Cannot subtract money with different currencies");
        return new Money(a.Amount - b.Amount, a.CurrencyCode);
    }

    public static Money operator *(Money a, decimal factor)
    {
        return new Money((long)Math.Round(a.Amount * factor, MidpointRounding.AwayFromZero), a.CurrencyCode);
    }

    public static Money operator *(Money a, double factor)
    {
        return new Money((long)Math.Round(a.Amount * factor, MidpointRounding.AwayFromZero), a.CurrencyCode);
    }
}
