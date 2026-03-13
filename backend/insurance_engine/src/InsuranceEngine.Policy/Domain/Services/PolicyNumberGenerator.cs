using System;

namespace InsuranceEngine.Policy.Domain.Services;

/// <summary>
/// Generates policy numbers in format: LBT-{YYYY}-{XXXX}-{NNNNNN}
/// where YYYY = current year, XXXX = product code suffix, NNNNNN = zero-padded sequential number.
/// Uses a database sequence for thread-safe number generation.
/// </summary>
public class PolicyNumberGenerator
{
    private static long _counter = 0;
    private static readonly object _lock = new();

    /// <summary>
    /// Generate next policy number. In production, this should use a DB sequence.
    /// </summary>
    public string Generate(string productCode)
    {
        var year = DateTime.UtcNow.Year;

        // Extract suffix from product code (e.g., "HLT-001" -> "H001")
        var suffix = ExtractSuffix(productCode);

        long sequenceNumber;
        lock (_lock)
        {
            sequenceNumber = ++_counter;
        }

        return $"LBT-{year}-{suffix}-{sequenceNumber:D6}";
    }

    /// <summary>
    /// Generate using a specific sequence number (from DB sequence).
    /// </summary>
    public string Generate(string productCode, long sequenceNumber)
    {
        var year = DateTime.UtcNow.Year;
        var suffix = ExtractSuffix(productCode);
        return $"LBT-{year}-{suffix}-{sequenceNumber:D6}";
    }

    private static string ExtractSuffix(string productCode)
    {
        // Product code format: "HLT-001" → suffix "H001"
        // Take first char + last 3 chars
        if (string.IsNullOrEmpty(productCode) || productCode.Length < 7)
            return "XXXX";

        var parts = productCode.Split('-');
        if (parts.Length == 2)
        {
            var prefix = parts[0].Length > 0 ? parts[0][0].ToString() : "X";
            var number = parts[1].PadLeft(3, '0');
            return $"{prefix}{number}";
        }

        return productCode.Replace("-", "").PadRight(4).Substring(0, 4).ToUpperInvariant();
    }
}
