using InsuranceEngine.SharedKernel.Domain;

namespace InsuranceEngine.Policy.Domain.ValueObjects;

public class PolicyNumber : ValueObject
{
    public string Value { get; }

    private PolicyNumber(string value) => Value = value;

    /// <summary>
    /// Generates a policy number using DB-provided sequence.
    /// Format: LBT-YYYY-XXXX-NNNNNN (e.g., LBT-2026-0001-000042)
    /// </summary>
    /// <param name="productCode">4-char product code segment</param>
    /// <param name="sequenceNumber">Sequential number from PostgreSQL sequence</param>
    public static PolicyNumber Generate(string productCode, long sequenceNumber)
    {
        var year = DateTime.UtcNow.Year;
        var seq = sequenceNumber.ToString().PadLeft(6, '0');
        
        return new PolicyNumber($"LBT-{year}-{productCode}-{seq}");
    }

    public static PolicyNumber From(string value)
    {
        if (string.IsNullOrWhiteSpace(value)) throw new ArgumentException("Policy number cannot be empty");
        return new PolicyNumber(value);
    }

    protected override IEnumerable<object> GetEqualityComponents()
    {
        yield return Value;
    }

    public override string ToString() => Value;
}

