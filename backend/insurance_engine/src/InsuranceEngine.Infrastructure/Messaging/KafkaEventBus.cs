using System.Threading.Tasks;
using Confluent.Kafka;
using InsuranceEngine.Application.Interfaces;
using Microsoft.Extensions.Options;

namespace InsuranceEngine.Infrastructure.Messaging;

public class KafkaEventBus : IEventBus
{
    private readonly InsuranceKafkaOptions _options;

    public KafkaEventBus(IOptions<InsuranceKafkaOptions> options)
    {
        _options = options.Value;
    }

    public async Task PublishAsync<T>(string topic, T @event)
    {
        var config = new ProducerConfig 
        { 
            BootstrapServers = _options.BootstrapServers 
        };
        
        using var producer = new ProducerBuilder<string, string>(config).Build();
        
        var message = Newtonsoft.Json.JsonConvert.SerializeObject(@event);
        await producer.ProduceAsync(topic, new Message<string, string> { Value = message });
    }
}
