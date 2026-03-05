using PoliSync.SharedKernel.Domain;

namespace PoliSync.Products.Domain;

public record ProductCreatedEvent(
    Guid ProductId, string ProductCode, string ProductName, ProductCategory Category
) : DomainEvent
{
    public override string EventType => "product.created.v1";
}

public record ProductUpdatedEvent(
    Guid ProductId, string ProductCode
) : DomainEvent
{
    public override string EventType => "product.updated.v1";
}

public record ProductActivatedEvent(
    Guid ProductId, string ProductCode, string ProductName
) : DomainEvent
{
    public override string EventType => "product.activated.v1";
}

public record ProductDeactivatedEvent(
    Guid ProductId, string ProductCode, string? Reason
) : DomainEvent
{
    public override string EventType => "product.deactivated.v1";
}

public record ProductDiscontinuedEvent(
    Guid ProductId, string ProductCode, string? Reason
) : DomainEvent
{
    public override string EventType => "product.discontinued.v1";
}
