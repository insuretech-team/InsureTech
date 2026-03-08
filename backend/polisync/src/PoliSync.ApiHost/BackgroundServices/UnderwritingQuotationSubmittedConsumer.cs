using System.Text.Json;
using Confluent.Kafka;
using Grpc.Core;
using Insuretech.Common.V1;
using Insuretech.Insurance.Services.V1;
using Insuretech.Underwriting.Entity.V1;
using PoliSync.Infrastructure.Clients;

namespace PoliSync.ApiHost.BackgroundServices;

public sealed class UnderwritingQuotationSubmittedConsumer : BackgroundService
{
    private readonly ILogger<UnderwritingQuotationSubmittedConsumer> _logger;
    private readonly InsuranceServiceClient _insuranceClient;
    private readonly IConsumer<Ignore, string> _consumer;
    private readonly string _topic;

    public UnderwritingQuotationSubmittedConsumer(
        IConfiguration configuration,
        ILogger<UnderwritingQuotationSubmittedConsumer> logger,
        InsuranceServiceClient insuranceClient)
    {
        _logger = logger;
        _insuranceClient = insuranceClient;

        _topic = configuration["Kafka:Topics:QuotationSubmitted"] ?? "insuretech.quotation.submitted.v1";
        var bootstrapServers = configuration["Kafka:BootstrapServers"] ?? "localhost:9092";
        var groupId = configuration["Kafka:Consumer:UnderwritingQuotationSubmitted:GroupId"] ?? "polisync-underwriting-quotation-submitted";

        var consumerConfig = new ConsumerConfig
        {
            BootstrapServers = bootstrapServers,
            GroupId = groupId,
            AutoOffsetReset = AutoOffsetReset.Earliest,
            EnableAutoCommit = false
        };

        _consumer = new ConsumerBuilder<Ignore, string>(consumerConfig).Build();
    }

    protected override Task ExecuteAsync(CancellationToken stoppingToken)
    {
        _consumer.Subscribe(_topic);
        _logger.LogInformation("Subscribed to Kafka topic {Topic} for underwriting quote seeding", _topic);

        return Task.Run(async () =>
        {
            while (!stoppingToken.IsCancellationRequested)
            {
                ConsumeResult<Ignore, string>? result = null;
                try
                {
                    result = _consumer.Consume(stoppingToken);
                }
                catch (OperationCanceledException)
                {
                    break;
                }
                catch (ConsumeException ex)
                {
                    _logger.LogError(ex, "Kafka consume error on topic {Topic}", _topic);
                    continue;
                }

                if (result is null)
                {
                    continue;
                }

                var processed = await TryProcessMessageAsync(result.Message.Value, stoppingToken);
                if (processed)
                {
                    _consumer.Commit(result);
                }
            }
        }, stoppingToken);
    }

    public override void Dispose()
    {
        _consumer.Close();
        _consumer.Dispose();
        base.Dispose();
    }

    private async Task<bool> TryProcessMessageAsync(string payload, CancellationToken cancellationToken)
    {
        if (string.IsNullOrWhiteSpace(payload))
        {
            return true;
        }

        try
        {
            using var doc = JsonDocument.Parse(payload);
            var root = doc.RootElement;

            var quotationId = FirstString(root, "quotation_id", "quotationId", "quote_id", "quoteId");
            if (string.IsNullOrWhiteSpace(quotationId))
            {
                _logger.LogWarning("QuotationSubmitted event missing quotation/quote id. Payload ignored.");
                return true;
            }

            var existing = await _insuranceClient.Client.GetQuoteAsync(
                new GetQuoteRequest { QuoteId = quotationId },
                cancellationToken: cancellationToken);

            if (existing?.Quote is not null && !string.IsNullOrWhiteSpace(existing.Quote.Id))
            {
                return true;
            }
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            // Expected for first create path.
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to parse quotation submitted event");
            return false;
        }

        try
        {
            using var doc = JsonDocument.Parse(payload);
            var root = doc.RootElement;

            var quotationId = FirstString(root, "quotation_id", "quotationId", "quote_id", "quoteId")!;
            var beneficiaryId = FirstString(root, "beneficiary_id", "beneficiaryId", "customer_id", "customerId", "business_id", "businessId")
                ?? quotationId;
            var insurerProductId = FirstString(root, "insurer_product_id", "insurerProductId", "product_id", "productId", "plan_id", "planId")
                ?? "unknown-product";

            var currency = FirstString(root, "currency") ?? "BDT";
            var sumAssuredAmount = FirstMoneyAmount(root, "sum_assured", "sumAssured") ?? 500_000L;
            var quotedAmount = FirstMoneyAmount(root, "quoted_amount", "quotedAmount", "estimated_premium", "estimatedPremium") ?? 0L;
            if (quotedAmount <= 0)
            {
                quotedAmount = (long)Math.Round(sumAssuredAmount * 0.04m, MidpointRounding.AwayFromZero);
            }

            var termYears = FirstInt(root, "term_years", "termYears", "tenure_years", "tenureYears");
            if (termYears <= 0)
            {
                termYears = 1;
            }

            var now = DateTime.UtcNow;
            var quote = new Quote
            {
                Id = quotationId,
                QuoteNumber = FirstString(root, "quotation_number", "quotationNumber", "quote_number", "quoteNumber")
                    ?? $"QTE-{now:yyyyMMdd}-{Random.Shared.Next(100000, 999999)}",
                BeneficiaryId = beneficiaryId,
                InsurerProductId = insurerProductId,
                Status = QuoteStatus.PendingUnderwriting,
                SumAssured = new Money { Amount = sumAssuredAmount, Currency = currency },
                TermYears = termYears,
                PremiumPaymentMode = FirstString(root, "premium_payment_mode", "premiumPaymentMode") ?? "YEARLY",
                BasePremium = new Money { Amount = quotedAmount, Currency = currency },
                RiderPremium = new Money { Amount = 0, Currency = currency },
                TaxAmount = new Money { Amount = 0, Currency = currency },
                TotalPremium = new Money { Amount = quotedAmount, Currency = currency },
                PremiumCalculation = "seeded_from_quotation_submitted_event=true",
                SelectedRiders = string.Empty,
                ApplicantAge = FirstInt(root, "applicant_age", "applicantAge"),
                ApplicantOccupation = FirstString(root, "applicant_occupation", "applicantOccupation") ?? string.Empty,
                Smoker = FirstBool(root, "smoker"),
                ValidUntil = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(now.AddDays(30))
            };

            await _insuranceClient.Client.CreateQuoteAsync(
                new CreateQuoteRequest { Quote = quote },
                cancellationToken: cancellationToken);

            _logger.LogInformation("Seeded underwriting quote from quotation submitted event. QuoteId={QuoteId}", quote.Id);
            return true;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to seed underwriting quote from quotation submitted event");
            return false;
        }
    }

    private static string? FirstString(JsonElement root, params string[] names)
    {
        foreach (var name in names)
        {
            if (root.TryGetProperty(name, out var element) && element.ValueKind == JsonValueKind.String)
            {
                return element.GetString();
            }
        }

        return null;
    }

    private static int FirstInt(JsonElement root, params string[] names)
    {
        foreach (var name in names)
        {
            if (root.TryGetProperty(name, out var element))
            {
                if (element.ValueKind == JsonValueKind.Number && element.TryGetInt32(out var number))
                {
                    return number;
                }

                if (element.ValueKind == JsonValueKind.String && int.TryParse(element.GetString(), out var parsed))
                {
                    return parsed;
                }
            }
        }

        return 0;
    }

    private static bool FirstBool(JsonElement root, params string[] names)
    {
        foreach (var name in names)
        {
            if (root.TryGetProperty(name, out var element))
            {
                if (element.ValueKind == JsonValueKind.True)
                {
                    return true;
                }

                if (element.ValueKind == JsonValueKind.False)
                {
                    return false;
                }

                if (element.ValueKind == JsonValueKind.String && bool.TryParse(element.GetString(), out var parsed))
                {
                    return parsed;
                }
            }
        }

        return false;
    }

    private static long? FirstMoneyAmount(JsonElement root, params string[] fieldNames)
    {
        foreach (var fieldName in fieldNames)
        {
            if (!root.TryGetProperty(fieldName, out var field))
            {
                continue;
            }

            if (field.ValueKind == JsonValueKind.Number && field.TryGetInt64(out var amountNumber))
            {
                return amountNumber;
            }

            if (field.ValueKind == JsonValueKind.Object && field.TryGetProperty("amount", out var amount))
            {
                if (amount.ValueKind == JsonValueKind.Number && amount.TryGetInt64(out var nestedAmountNumber))
                {
                    return nestedAmountNumber;
                }

                if (amount.ValueKind == JsonValueKind.String && long.TryParse(amount.GetString(), out var nestedAmountString))
                {
                    return nestedAmountString;
                }
            }
        }

        return null;
    }
}
