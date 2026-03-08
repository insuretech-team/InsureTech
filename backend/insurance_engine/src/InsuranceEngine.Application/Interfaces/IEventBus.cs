using System.Threading.Tasks;

namespace InsuranceEngine.Application.Interfaces;

public interface IEventBus
{
    Task PublishAsync<T>(string topic, T @event);
}
