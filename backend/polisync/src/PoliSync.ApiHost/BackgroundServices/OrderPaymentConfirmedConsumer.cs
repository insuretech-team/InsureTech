using System.Text.Json;
using System.Text.Json.Serialization;
using Confluent.Kafka;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Common.V1;
using Insuretech.Orders.Entity.V1;
using PoliSync.Orders.Infrastructure;
using PoliSync.Policy.Domain;
using PoliSync.Policy.Infrastructure;
using PoliSync.Quotes.Infrastructure;
using PoliSync.SharedKernel.Domain;
using PoliSync.SharedKernel.Messaging;

namespace PoliSync.ApiHost.BackgroundServices;

public sealed class OrderPaymentConfirmedConsumer : BackgroundService
{
    private readonly ILogger<OrderPaymentConfirmedConsumer> _logger;
    private readonly IServiceScopeFactory _scopeFactory;
    private readonly IEventBus _eventBus;
    private readonly IConsumer<Ignore, string> _consumer;
    private readonly string _consumeTopic;
    private readonly string _projectionTopic;

    public OrderPaymentConfirmedConsumer(
        IConfiguration configuration,
        ILogger<OrderPaymentConfirmedConsumer> logger,
        IServiceScopeFactory scopeFactory,
        IEventBus eventBus)
    {
        _logger = logger;
        _scopeFactory = scopeFactory;
        _eventBus = eventBus;

        _consumeTopic = configuration["Kafka:Topics:OrderPaymentConfirmed"] ?? "orders.order.payment_confirmed";
        _projectionTopic = configuration["Kafka:Topics:OrderPolicyProjectionIssued"] ?? "policy.issued";
        var bootstrapServers = configuration["Kafka:BootstrapServers"] ?? "localhost:9092";
        var groupId = configuration["Kafka:Consumer:OrderPaymentConfirmed:GroupId"] ?? "polisync-order-payment-confirmed";

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
        _consumer.Subscribe(_consumeTopic);
        _logger.LogInformation("Subscribed to Kafka topic {Topic} for policy issuance", _consumeTopic);

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
                    _logger.LogError(ex, "Kafka consume error on topic {Topic}", _consumeTopic);
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

        OrderPaymentConfirmedPayload? evt;
        try
        {
            using var doc = JsonDocument.Parse(payload);
            evt = ParsePayload(doc.RootElement);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to parse order payment confirmed event payload");
            return false;
        }

        if (evt is null || string.IsNullOrWhiteSpace(evt.OrderId))
        {
            _logger.LogWarning("OrderPaymentConfirmed event missing order_id/orderId. Payload ignored.");
            return true;
        }

        try
        {
            using var scope = _scopeFactory.CreateScope();
            var orderGateway = scope.ServiceProvider.GetRequiredService<IOrderDataGateway>();
            var quotationGateway = scope.ServiceProvider.GetRequiredService<IQuotationDataGateway>();
            var policyGateway = scope.ServiceProvider.GetRequiredService<IPolicyDataGateway>();

            var order = await orderGateway.GetOrderAsync(evt.OrderId, cancellationToken);
            if (order?.Order is null)
            {
                _logger.LogWarning("Order {OrderId} not found while processing payment confirmed event", evt.OrderId);
                return true;
            }

            if (!string.IsNullOrWhiteSpace(order.Order.PolicyId) || order.Order.Status == OrderStatus.PolicyIssued)
            {
                return true;
            }

            var quotationId = FirstNonEmpty(evt.QuotationId, order.Order.QuotationId);
            var customerId = FirstNonEmpty(evt.CustomerId, order.Order.CustomerId);
            var productId = FirstNonEmpty(evt.ProductId, order.Order.ProductId);
            var paymentId = FirstNonEmpty(evt.PaymentId, order.Order.PaymentId);

            if (string.IsNullOrWhiteSpace(quotationId) || string.IsNullOrWhiteSpace(customerId) || string.IsNullOrWhiteSpace(productId))
            {
                _logger.LogWarning(
                    "Skipping policy issuance for order {OrderId} because quotation/customer/product data is incomplete",
                    evt.OrderId);
                return true;
            }

            var quotation = await quotationGateway.GetQuotationAsync(quotationId, cancellationToken);
            if (quotation is null)
            {
                _logger.LogWarning("Quotation {QuotationId} not found for order {OrderId}", quotationId, evt.OrderId);
                return true;
            }

            if (quotation.Status != Insuretech.Policy.Entity.V1.QuotationStatus.Approved)
            {
                _logger.LogWarning(
                    "Skipping policy issuance for order {OrderId} because quotation {QuotationId} is in status {Status}",
                    evt.OrderId,
                    quotationId,
                    quotation.Status);
                return true;
            }

            var premium = ResolveMoney(evt.TotalPayableAmount, evt.TotalPayableCurrency, order.Order.TotalPayable, quotation);
            var sumInsured = ResolveSumInsured(quotation, premium.Currency);
            var now = DateTime.UtcNow;

            var aggregate = PolicyAggregate.Create(
                customerId,
                productId,
                quotationId,
                premium.Amount,
                sumInsured.Amount,
                12,
                now.Date,
                now.Date.AddMonths(12));

            aggregate.IssuePolicy();
            aggregate.Policy.PartnerId = quotation.BusinessId;
            aggregate.Policy.PaymentGatewayReference = paymentId;
            aggregate.Policy.ReceiptNumber = $"RCPT-{DateTime.UtcNow:yyyyMMdd}-{Random.Shared.Next(100000, 999999)}";
            aggregate.Policy.PaymentFrequency = "ANNUAL";
            aggregate.Policy.PremiumCurrency = premium.Currency;
            aggregate.Policy.SumInsuredCurrency = sumInsured.Currency;
            aggregate.Policy.PremiumAmount = premium;
            aggregate.Policy.SumInsured = sumInsured;
            aggregate.Policy.TotalPayable = premium;
            aggregate.Policy.VatTax = NewMoney(0, premium.Currency);
            aggregate.Policy.ServiceFee = NewMoney(0, premium.Currency);
            aggregate.Policy.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);

            var created = await policyGateway.CreatePolicyAsync(aggregate.Policy, cancellationToken);

            await _eventBus.PublishAsync(
                new PolicyIssuedProjectionEvent(created.PolicyId, evt.OrderId),
                _projectionTopic,
                cancellationToken);

            _logger.LogInformation(
                "Issued policy {PolicyId} for order {OrderId} from payment confirmed event",
                created.PolicyId,
                evt.OrderId);

            return true;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to issue policy from order payment confirmed event for order {OrderId}", evt.OrderId);
            return false;
        }
    }

    private static OrderPaymentConfirmedPayload ParsePayload(JsonElement root)
    {
        var money = FirstMoney(root, "totalPayable", "total_payable");

        return new OrderPaymentConfirmedPayload(
            OrderId: FirstString(root, "order_id", "orderId"),
            PaymentId: FirstString(root, "payment_id", "paymentId"),
            QuotationId: FirstString(root, "quotation_id", "quotationId"),
            CustomerId: FirstString(root, "customer_id", "customerId"),
            ProductId: FirstString(root, "product_id", "productId"),
            TotalPayableAmount: money?.Amount ?? 0,
            TotalPayableCurrency: money?.Currency ?? "BDT");
    }

    private static Money ResolveMoney(
        long eventAmount,
        string eventCurrency,
        Money? orderAmount,
        Insuretech.Policy.Entity.V1.Quotation quotation)
    {
        if (eventAmount > 0)
        {
            return NewMoney(eventAmount, string.IsNullOrWhiteSpace(eventCurrency) ? "BDT" : eventCurrency);
        }

        if (orderAmount is not null && orderAmount.Amount > 0)
        {
            return NewMoney(orderAmount.Amount, string.IsNullOrWhiteSpace(orderAmount.Currency) ? "BDT" : orderAmount.Currency);
        }

        if (quotation.QuotedAmount is not null && quotation.QuotedAmount.Amount > 0)
        {
            return NewMoney(quotation.QuotedAmount.Amount, string.IsNullOrWhiteSpace(quotation.QuotedAmount.Currency) ? "BDT" : quotation.QuotedAmount.Currency);
        }

        if (quotation.EstimatedPremium is not null && quotation.EstimatedPremium.Amount > 0)
        {
            return NewMoney(quotation.EstimatedPremium.Amount, string.IsNullOrWhiteSpace(quotation.EstimatedPremium.Currency) ? "BDT" : quotation.EstimatedPremium.Currency);
        }

        return NewMoney(120_000);
    }

    private static Money ResolveSumInsured(Insuretech.Policy.Entity.V1.Quotation quotation, string currency)
    {
        if (quotation.QuotedAmount is not null && quotation.QuotedAmount.Amount > 0)
        {
            return NewMoney(Math.Max(quotation.QuotedAmount.Amount * 10, 1_000_000), currency);
        }

        if (quotation.EstimatedPremium is not null && quotation.EstimatedPremium.Amount > 0)
        {
            return NewMoney(Math.Max(quotation.EstimatedPremium.Amount * 10, 1_000_000), currency);
        }

        return NewMoney(1_000_000, currency);
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

    private static (long Amount, string Currency)? FirstMoney(JsonElement root, params string[] names)
    {
        foreach (var name in names)
        {
            if (!root.TryGetProperty(name, out var element))
            {
                continue;
            }

            if (element.ValueKind == JsonValueKind.Object)
            {
                var amount = 0L;
                var currency = "BDT";

                if (element.TryGetProperty("amount", out var amountElement))
                {
                    if (amountElement.ValueKind == JsonValueKind.Number && amountElement.TryGetInt64(out var numericAmount))
                    {
                        amount = numericAmount;
                    }
                    else if (amountElement.ValueKind == JsonValueKind.String && long.TryParse(amountElement.GetString(), out var parsedAmount))
                    {
                        amount = parsedAmount;
                    }
                }

                if (element.TryGetProperty("currency", out var currencyElement) && currencyElement.ValueKind == JsonValueKind.String)
                {
                    currency = currencyElement.GetString() ?? "BDT";
                }

                return (amount, currency);
            }
        }

        return null;
    }

    private static string? FirstNonEmpty(params string?[] values)
        => values.FirstOrDefault(value => !string.IsNullOrWhiteSpace(value));

    private static Money NewMoney(long amount, string currency = "BDT")
        => new() { Amount = amount, Currency = currency };

    private sealed record OrderPaymentConfirmedPayload(
        string? OrderId,
        string? PaymentId,
        string? QuotationId,
        string? CustomerId,
        string? ProductId,
        long TotalPayableAmount,
        string TotalPayableCurrency);

    private sealed record PolicyIssuedProjectionEvent(
        [property: JsonPropertyName("policy_id")] string PolicyId,
        [property: JsonPropertyName("order_id")] string OrderId) : DomainEvent;
}
