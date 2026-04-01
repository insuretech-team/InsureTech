using Microsoft.Extensions.Options;
using PoliSync.Quotes.Domain;
using PoliSync.Quotes.Infrastructure;

namespace PoliSync.ApiHost.BackgroundServices;

/// <summary>
/// Background service that expires quotations past their expiry date
/// Runs hourly to check for expired quotations
/// </summary>
public sealed class QuotationExpiryService : BackgroundService
{
    private readonly IServiceProvider _serviceProvider;
    private readonly ILogger<QuotationExpiryService> _logger;
    private readonly TimeSpan _interval = TimeSpan.FromHours(1);

    public QuotationExpiryService(
        IServiceProvider serviceProvider,
        ILogger<QuotationExpiryService> logger)
    {
        _serviceProvider = serviceProvider;
        _logger = logger;
    }

    protected override async Task ExecuteAsync(CancellationToken stoppingToken)
    {
        _logger.LogInformation("Quotation Expiry Service started");

        while (!stoppingToken.IsCancellationRequested)
        {
            try
            {
                await ProcessExpiredQuotationsAsync(stoppingToken);
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Error processing expired quotations");
            }

            await Task.Delay(_interval, stoppingToken);
        }

        _logger.LogInformation("Quotation Expiry Service stopped");
    }

    private async Task ProcessExpiredQuotationsAsync(CancellationToken cancellationToken)
    {
        using var scope = _serviceProvider.CreateScope();
        var gateway = scope.ServiceProvider.GetRequiredService<IQuotationDataGateway>();

        _logger.LogInformation("Checking for expired quotations...");

        // Get all quotations that are past expiry and not in terminal states
        var expiredQuotations = await gateway.GetExpiredQuotationsAsync(cancellationToken);

        if (!expiredQuotations.Any())
        {
            _logger.LogInformation("No expired quotations found");
            return;
        }

        _logger.LogInformation("Found {Count} expired quotations", expiredQuotations.Count);

        var expiredCount = 0;
        foreach (var quotation in expiredQuotations)
        {
            try
            {
                var result = quotation.Expire();
                if (result.IsSuccess)
                {
                    await gateway.UpdateAsync(quotation, cancellationToken);
                    expiredCount++;
                    _logger.LogInformation(
                        "Expired quotation {QuotationNumber} (ID: {QuotationId})",
                        quotation.QuotationNumber,
                        quotation.Id);
                }
                else
                {
                    _logger.LogWarning(
                        "Failed to expire quotation {QuotationNumber}: {Error}",
                        quotation.QuotationNumber,
                        result.Error?.Message);
                }
            }
            catch (Exception ex)
            {
                _logger.LogError(
                    ex,
                    "Error expiring quotation {QuotationNumber}",
                    quotation.QuotationNumber);
            }
        }

        _logger.LogInformation("Successfully expired {Count} quotations", expiredCount);
    }
}
