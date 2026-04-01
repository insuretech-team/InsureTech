using System.Threading;

namespace PoliSync.LifeInsurance.Services;

public class QuoteNumberGenerator : IQuoteNumberGenerator
{
    private static int _counter = 0;
    private readonly string _prefix = "LF";

    public string GenerateQuoteNumber()
    {
        var year = DateTime.UtcNow.Year;
        var sequence = Interlocked.Increment(ref _counter);
        return $"{_prefix}-{year}-{sequence:D6}";
    }
}
