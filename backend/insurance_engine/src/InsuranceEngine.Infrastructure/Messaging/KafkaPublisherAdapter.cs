using InsuranceEngine.SharedKernel.Infrastructure;
using Microsoft.Extensions.Logging;

namespace InsuranceEngine.Infrastructure.Messaging;

public class KafkaPublisherAdapter : IKafkaPublisher
{
    private readonly IEventPublisher _eventPublisher;
    private readonly ILogger<KafkaPublisherAdapter> _logger;

    public KafkaPublisherAdapter(IEventPublisher eventPublisher, ILogger<KafkaPublisherAdapter> logger)
    {
        _eventPublisher = eventPublisher;
        _logger = logger;
    }

    public async Task PublishAsync<T>(string topic, T message) where T : class
    {
        try
        {
            _logger.LogInformation("Publishing to Kafka topic {Topic}: {MessageType}", topic, typeof(T).Name);
            await _eventPublisher.PublishAsync(topic, message);
            _logger.LogInformation("Successfully published to {Topic}", topic);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to publish to Kafka topic {Topic}", topic);
            throw;
        }
    }
}
