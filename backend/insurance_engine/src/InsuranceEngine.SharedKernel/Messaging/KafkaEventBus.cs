using System;
using System.Threading.Tasks;
using Confluent.Kafka;
using InsuranceEngine.SharedKernel.Interfaces;
using Microsoft.Extensions.Options;

namespace InsuranceEngine.SharedKernel.Messaging;

public class KafkaEventBus : IEventBus, IDisposable
{
    private readonly InsuranceKafkaOptions _options;
    private readonly IProducer<string, string> _producer;

    public KafkaEventBus(IOptions<InsuranceKafkaOptions> options)
    {
        _options = options.Value;
        
        var config = new ProducerConfig 
        { 
            BootstrapServers = _options.BootstrapServers 
        };
        
        _producer = new ProducerBuilder<string, string>(config).Build();
    }

    public async Task PublishAsync<T>(string topic, T @event)
    {
        var message = Newtonsoft.Json.JsonConvert.SerializeObject(@event);
        await _producer.ProduceAsync(topic, new Message<string, string> { Value = message });
    }

    public void Dispose()
    {
        _producer?.Flush(TimeSpan.FromSeconds(10));
        _producer?.Dispose();
    }
}
