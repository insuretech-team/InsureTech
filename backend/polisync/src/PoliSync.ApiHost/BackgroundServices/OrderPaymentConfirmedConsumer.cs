using System.Text.Json;
using Confluent.Kafka;
using PoliSync.ApiHost.Services;
using PoliSync.Orders.Infrastructure;

namespace PoliSync.ApiHost.BackgroundServices;

public sealed class OrderPaymentConfirmedConsumer : BackgroundService
{
    private readonly ILogger<OrderPaymentConfirmedConsumer> _logger;
    private readonly IServiceScopeFactory _scopeFactory;
    private readonly IConsumer<Ignore, string> _consumer;
    private readonly string _consumeTopic;

    public OrderPaymentConfirmedConsumer(
        IConfiguration configuration,
        ILogger<OrderPaymentConfirmedConsumer> logger,
        IServiceScopeFactory scopeFactory)
    {
        _logger = logger;
        _scopeFactory = scopeFactory;

        _consumeTopic = configuration["Kafka:Topics:OrderPaymentConfirmed"] ?? "insuretech.orders.v1.payment_confirmed";
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
        _logger.LogInformation("Subscribed to Kafka topic {Topic} for proposal submission", _consumeTopic);

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
            var workflowService = scope.ServiceProvider.GetRequiredService<InsuranceProposalWorkflowService>();

            var order = await orderGateway.GetOrderAsync(evt.OrderId, cancellationToken);
            if (order?.Order is null)
            {
                _logger.LogWarning("Order {OrderId} not found while processing payment confirmed event", evt.OrderId);
                return true;
            }

            if (!string.IsNullOrWhiteSpace(order.Order.PolicyId))
            {
                return true;
            }

            var insurerId = FirstNonEmpty(evt.InsurerId, order.Order.InsurerId);

            if (string.IsNullOrWhiteSpace(insurerId))
            {
                _logger.LogWarning(
                    "Skipping proposal submission for order {OrderId} because insurer_id is missing",
                    evt.OrderId);
                return true;
            }

            var created = await workflowService.SubmitProposalForOrderAsync(
                evt.OrderId,
                insurerId,
                correlationId: FirstNonEmpty(evt.CorrelationId, order.Order.CorrelationId),
                submissionPayload: payload,
                totalPayableAmount: evt.TotalPayableAmount > 0 ? evt.TotalPayableAmount : null,
                totalPayableCurrency: evt.TotalPayableCurrency,
                cancellationToken);

            _logger.LogInformation(
                "Submitted proposal {ProposalId} for order {OrderId} from payment confirmed event",
                created.ProposalId,
                evt.OrderId);

            return true;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to submit proposal from order payment confirmed event for order {OrderId}", evt.OrderId);
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
            InsurerId: FirstString(root, "insurer_id", "insurerId"),
            CorrelationId: FirstString(root, "correlation_id", "correlationId"),
            TotalPayableAmount: money?.Amount ?? 0,
            TotalPayableCurrency: money?.Currency ?? "BDT");
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

    private sealed record OrderPaymentConfirmedPayload(
        string? OrderId,
        string? PaymentId,
        string? QuotationId,
        string? CustomerId,
        string? ProductId,
        string? InsurerId,
        string? CorrelationId,
        long TotalPayableAmount,
        string TotalPayableCurrency);
}
