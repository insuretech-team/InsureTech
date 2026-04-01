using System;

namespace InsuranceEngine.Underwriting.Domain.Services;

public class QuoteNumberGenerator
{
    public string Generate(string productCode, long sequenceNumber)
    {
        var year = DateTime.UtcNow.Year;
        var suffix = ExtractSuffix(productCode);
        return $"QTE-{year}-{suffix}-{sequenceNumber:D6}";
    }

    private static string ExtractSuffix(string productCode)
    {
        if (string.IsNullOrEmpty(productCode) || productCode.Length < 7)
            return "XXXX";

        var parts = productCode.Split('-');
        if (parts.Length == 2)
        {
            var prefix = parts[0].Length > 0 ? parts[0][0].ToString() : "Q";
            var number = parts[1].PadLeft(3, '0');
            return $"{prefix}{number}";
        }

        return productCode.Replace("-", "").PadRight(4).Substring(0, 4).ToUpperInvariant();
    }
}
