using System.Threading;

namespace PoliSync.Quoting.Services;

public class QuoteNumberGenerator : IQuoteNumberGenerator
{
    private static int _counter = 0;
    private readonly string _prefix = "QT";

    public string GenerateQuoteNumber()
    {
        var year = DateTime.UtcNow.Year;
        var sequence = Interlocked.Increment(ref _counter);
        return $"{_prefix}-{year}-{sequence:D6}";
    }
}
