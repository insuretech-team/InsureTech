namespace InsuranceEngine.SharedKernel.Domain;

public class Money : ValueObject
{
    public long Amount { get; } // Paisa
    public string Currency { get; }

    public Money(long amount, string currency = "BDT")
    {
        Amount = amount;
        Currency = currency;
    }

    public static Money FromDecimal(decimal amount, string currency = "BDT") => new((long)(amount * 100), currency);
    public decimal ToDecimal() => Amount / 100m;

    protected override IEnumerable<object> GetEqualityComponents()
    {
        yield return Amount;
        yield return Currency;
    }

    public override string ToString() => $"{ToDecimal():F2} {Currency}";
}
