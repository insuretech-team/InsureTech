namespace PoliSync.SharedKernel.Domain;

/// <summary>
/// Money value object - maps to common.v1.Money proto (int64 paisa + currency)
/// BDT subunit: 1 BDT = 100 paisa
/// </summary>
public sealed class Money : ValueObject
{
    public long AmountInPaisa { get; }
    public string Currency { get; }

    private Money(long amountInPaisa, string currency)
    {
        AmountInPaisa = amountInPaisa;
        Currency = currency;
    }

    public static Money FromPaisa(long paisa, string currency = "BDT")
    {
        return new Money(paisa, currency);
    }

    public static Money FromBdt(decimal bdt)
    {
        return new Money((long)(bdt * 100), "BDT");
    }

    public static Money Zero(string currency = "BDT")
    {
        return new Money(0, currency);
    }

    public decimal ToBdt()
    {
        if (Currency != "BDT")
            throw new InvalidOperationException($"Cannot convert {Currency} to BDT");
        
        return AmountInPaisa / 100m;
    }

    public Money Add(Money other)
    {
        if (Currency != other.Currency)
            throw new InvalidOperationException($"Cannot add {Currency} and {other.Currency}");

        return new Money(AmountInPaisa + other.AmountInPaisa, Currency);
    }

    public Money Subtract(Money other)
    {
        if (Currency != other.Currency)
            throw new InvalidOperationException($"Cannot subtract {other.Currency} from {Currency}");

        return new Money(AmountInPaisa - other.AmountInPaisa, Currency);
    }

    public Money Multiply(decimal factor)
    {
        return new Money((long)(AmountInPaisa * factor), Currency);
    }

    public Money MultiplyByPercentage(decimal percentage)
    {
        return new Money((long)(AmountInPaisa * percentage / 100), Currency);
    }

    public bool IsPositive => AmountInPaisa > 0;
    public bool IsNegative => AmountInPaisa < 0;
    public bool IsZero => AmountInPaisa == 0;

    protected override IEnumerable<object?> GetEqualityComponents()
    {
        yield return AmountInPaisa;
        yield return Currency;
    }

    public override string ToString()
    {
        return $"{ToBdt():N2} {Currency}";
    }

    public static bool operator >(Money left, Money right)
    {
        if (left.Currency != right.Currency)
            throw new InvalidOperationException($"Cannot compare {left.Currency} and {right.Currency}");
        
        return left.AmountInPaisa > right.AmountInPaisa;
    }

    public static bool operator <(Money left, Money right)
    {
        if (left.Currency != right.Currency)
            throw new InvalidOperationException($"Cannot compare {left.Currency} and {right.Currency}");
        
        return left.AmountInPaisa < right.AmountInPaisa;
    }

    public static bool operator >=(Money left, Money right)
    {
        return left > right || left.Equals(right);
    }

    public static bool operator <=(Money left, Money right)
    {
        return left < right || left.Equals(right);
    }

    public static Money operator +(Money left, Money right)
    {
        return left.Add(right);
    }

    public static Money operator -(Money left, Money right)
    {
        return left.Subtract(right);
    }

    public static Money operator *(Money money, decimal factor)
    {
        return money.Multiply(factor);
    }
}
