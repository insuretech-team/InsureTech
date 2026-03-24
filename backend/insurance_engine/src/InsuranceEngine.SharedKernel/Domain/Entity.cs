namespace InsuranceEngine.SharedKernel.Domain;

public abstract class Entity
{
    public string Id { get; protected set; } = string.Empty;

    protected Entity(string id) => Id = id;
    protected Entity() { }
}

public abstract class AggregateRoot : Entity
{
    private readonly List<object> _domainEvents = new();
    public IReadOnlyCollection<object> DomainEvents => _domainEvents.AsReadOnly();

    protected void AddDomainEvent(object domainEvent) => _domainEvents.Add(domainEvent);
    public void ClearDomainEvents() => _domainEvents.Clear();
}
