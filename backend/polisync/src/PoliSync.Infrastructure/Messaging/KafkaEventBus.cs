using Confluent.Kafka;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using PoliSync.SharedKernel.Domain;
using PoliSync.SharedKernel.Messaging;
using System.Text.Json;

namespace PoliSync.Infrastructure.Messaging;

/// <summary>
/// Kafka event bus implementation for publishing domain events
/// </summary>
public sealed class KafkaEventBus : IEventBus, IDisposable
{
    private readonly IProducer<string, string> _producer;
    private readonly ILogger<KafkaEventBus> _logger;
    private readonly KafkaOptions _options;
    private readonly JsonSerializerOptions _jsonOptions;

    public KafkaEventBus(
        IOptions<KafkaOptions> options,
        ILogger<KafkaEventBus> logger)
    {
        _options = options.Value;
        _logger = logger;

        var config = new ProducerConfig
        {
            BootstrapServers = _options.BootstrapServers,
            Acks = Acks.All,
            EnableIdempotence = true,
            MaxInFlight = 5,
            MessageSendMaxRetries = 3,
            RetryBackoffMs = 100
        };

        _producer = new ProducerBuilder<string, string>(config).Build();

        _jsonOptions = new JsonSerializerOptions
        {
            PropertyNamingPolicy = JsonNamingPolicy.CamelCase,
            WriteIndented = false
        };
    }

    public async Task PublishAsync<TEvent>(
        TEvent @event, 
        string topic, 
        CancellationToken cancellationToken = default) 
        where TEvent : DomainEvent
    {
        try
        {
            var message = new Message<string, string>
            {
                Key = @event.EventId.ToString(),
                Value = JsonSerializer.Serialize(@event, _jsonOptions),
                Headers = new Headers
                {
                    { "event-type", System.Text.Encoding.UTF8.GetBytes(@event.EventType) },
                    { "event-id", System.Text.Encoding.UTF8.GetBytes(@event.EventId.ToString()) },
                    { "occurred-at", System.Text.Encoding.UTF8.GetBytes(@event.OccurredAt.ToString("O")) }
                }
            };

            var result = await _producer.ProduceAsync(topic, message, cancellationToken);

            _logger.LogInformation(
                "Published event {EventType} to topic {Topic} at offset {Offset}",
                @event.EventType, topic, result.Offset);
        }
        catch (ProduceException<string, string> ex)
        {
            _logger.LogError(ex, 
                "Failed to publish event {EventType} to topic {Topic}", 
                @event.EventType, topic);
            throw;
        }
    }

    public async Task PublishAsync<TEvent>(
        TEvent @event, 
        CancellationToken cancellationToken = default) 
        where TEvent : DomainEvent
    {
        var topic = GetTopicForEvent(@event);
        await PublishAsync(@event, topic, cancellationToken);
    }

    public async Task PublishBatchAsync<TEvent>(
        IEnumerable<TEvent> events, 
        string topic, 
        CancellationToken cancellationToken = default) 
        where TEvent : DomainEvent
    {
        var tasks = events.Select(e => PublishAsync(e, topic, cancellationToken));
        await Task.WhenAll(tasks);
    }

    private string GetTopicForEvent<TEvent>(TEvent @event) where TEvent : DomainEvent
    {
        // Map event type to Kafka topic based on naming convention
        // e.g., PolicyIssuedEvent -> insuretech.policy.issued.v1
        var eventTypeName = @event.EventType.Replace("Event", "");
        var parts = SplitCamelCase(eventTypeName);
        
        return $"insuretech.{string.Join(".", parts).ToLowerInvariant()}.v1";
    }

    private static string[] SplitCamelCase(string input)
    {
        return System.Text.RegularExpressions.Regex
            .Replace(input, "([A-Z])", " $1", System.Text.RegularExpressions.RegexOptions.Compiled)
            .Trim()
            .Split(' ');
    }

    public void Dispose()
    {
        _producer?.Flush(TimeSpan.FromSeconds(10));
        _producer?.Dispose();
    }
}

public sealed class KafkaOptions
{
    public string BootstrapServers { get; set; } = "localhost:9092";
    public Dictionary<string, string> Topics { get; set; } = new();
}
