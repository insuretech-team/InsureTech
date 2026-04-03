using System.Text.Json;
using Confluent.Kafka;
using InsuranceEngine.SharedKernel.Infrastructure;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;

namespace InsuranceEngine.Infrastructure.Messaging;

public class KafkaSettings
{
    public string BootstrapServers { get; set; } = "localhost:9092";
    public string ClientId { get; set; } = "insurance-engine";
    public int Acks { get; set; } = 2; // All = 2
    public int Retries { get; set; } = 3;
    public int RetryBackoffMs { get; set; } = 100;
    public bool EnableIdempotence { get; set; } = true;
}

public interface IEventPublisher
{
    Task PublishAsync<T>(string topic, T @event, string? key = null) where T : class;
}

public class KafkaEventPublisher : IEventPublisher, IDisposable
{
    private readonly IProducer<string, string> _producer;
    private readonly ILogger<KafkaEventPublisher> _logger;
    private readonly KafkaSettings _settings;
    private readonly JsonSerializerOptions _jsonOptions;
    private bool _disposed;

    public KafkaEventPublisher(
        ILogger<KafkaEventPublisher> logger,
        IOptions<KafkaSettings> settings)
    {
        _logger = logger;
        _settings = settings.Value;
        
        _jsonOptions = new JsonSerializerOptions
        {
            PropertyNamingPolicy = JsonNamingPolicy.CamelCase,
            WriteIndented = false
        };

        var config = new ProducerConfig
        {
            BootstrapServers = _settings.BootstrapServers,
            ClientId = _settings.ClientId,
            Acks = (Acks)_settings.Acks,
            MessageSendMaxRetries = _settings.Retries,
            RetryBackoffMs = _settings.RetryBackoffMs,
            EnableIdempotence = _settings.EnableIdempotence,
            LingerMs = 5,
            BatchSize = 16384
        };

        _producer = new ProducerBuilder<string, string>(config)
            .SetErrorHandler((_, e) => _logger.LogError("Kafka error: {Reason}", e.Reason))
            .SetLogHandler((_, log) => _logger.LogDebug("Kafka: {Message}", log.Message))
            .Build();

        _logger.LogInformation("Kafka producer initialized: {BootstrapServers}", _settings.BootstrapServers);
    }

    public async Task PublishAsync<T>(string topic, T @event, string? key = null) where T : class
    {
        try
        {
            var message = new Message<string, string>
            {
                Key = key ?? Guid.NewGuid().ToString(),
                Value = JsonSerializer.Serialize(@event, _jsonOptions),
                Headers = new Headers
                {
                    { "event-type", System.Text.Encoding.UTF8.GetBytes(typeof(T).Name) },
                    { "timestamp", System.Text.Encoding.UTF8.GetBytes(DateTime.UtcNow.ToString("O")) },
                    { "source", System.Text.Encoding.UTF8.GetBytes("insurance-engine") }
                }
            };

            var deliveryResult = await _producer.ProduceAsync(topic, message);

            _logger.LogInformation(
                "Event published to Kafka: Topic={Topic}, Key={Key}, Partition={Partition}, Offset={Offset}",
                topic,
                message.Key,
                deliveryResult.Partition.Value,
                deliveryResult.Offset.Value);
        }
        catch (ProduceException<string, string> ex)
        {
            _logger.LogError(ex, "Failed to publish event to Kafka: Topic={Topic}, Error={Error}", topic, ex.Error.Reason);
            throw;
        }
    }

    public void Dispose()
    {
        if (_disposed) return;
        
        _producer?.Flush(TimeSpan.FromSeconds(10));
        _producer?.Dispose();
        _disposed = true;
        
        _logger.LogInformation("Kafka producer disposed");
    }
}

public class KafkaEventPublisherFactory
{
    private readonly ILoggerFactory _loggerFactory;
    private readonly KafkaSettings _settings;

    public KafkaEventPublisherFactory(ILoggerFactory loggerFactory, KafkaSettings settings)
    {
        _loggerFactory = loggerFactory;
        _settings = settings;
    }

    public IEventPublisher CreatePublisher()
    {
        return new KafkaEventPublisher(
            _loggerFactory.CreateLogger<KafkaEventPublisher>(),
            Options.Create(_settings));
    }
}
