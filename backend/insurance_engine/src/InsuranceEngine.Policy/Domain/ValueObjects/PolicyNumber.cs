using InsuranceEngine.SharedKernel.Domain;

namespace InsuranceEngine.Policy.Domain.ValueObjects;

public class PolicyNumber : ValueObject
{
    public string Value { get; }

    private PolicyNumber(string value) => Value = value;

    public static PolicyNumber Generate(string insuranceType, string productCode)
    {
        var year = DateTime.UtcNow.Year;
        var sequence = new Random().Next(100000, 999999).ToString();
        
        // Format: LBT-YYYY-XXXX-NNNNNN (e.g., LBT-2026-001-123456)
        return new PolicyNumber($"LBT-{year}-{productCode}-{sequence}");
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
