using System.Threading.Tasks;
using Microsoft.Extensions.Logging;

namespace InsuranceEngine.SharedKernel.Infrastructure;

public interface IKafkaPublisher
{
    Task PublishAsync<T>(string topic, T message) where T : class;
}

public class MockKafkaPublisher : IKafkaPublisher
{
    private readonly ILogger<MockKafkaPublisher> _logger;

    public MockKafkaPublisher(ILogger<MockKafkaPublisher> logger)
    {
        _logger = logger;
    }

    public Task PublishAsync<T>(string topic, T message) where T : class
    {
        _logger.LogInformation("Streaming event to Kafka Topic [{Topic}]: {Message}", topic, System.Text.Json.JsonSerializer.Serialize(message));
        return Task.CompletedTask;
    }
}
