using System.Threading.Tasks;

namespace InsuranceEngine.SharedKernel.Interfaces;

public interface IEventBus
{
    Task PublishAsync<T>(string topic, T @event);
}
