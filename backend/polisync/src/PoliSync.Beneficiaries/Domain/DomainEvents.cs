using PoliSync.SharedKernel.Domain;

namespace PoliSync.Beneficiaries.Domain;

public record BeneficiaryCreatedEvent(Guid BeneficiaryId, string Code, BeneficiaryType Type) : DomainEvent
{
    public override string EventType => "Beneficiary.Created";
}

public record KycCompletedEvent(Guid BeneficiaryId, KycStatus Status) : DomainEvent
{
    public override string EventType => "Beneficiary.KycCompleted";
}

public record BeneficiaryStatusUpdatedEvent(Guid BeneficiaryId, BeneficiaryStatus Status) : DomainEvent
{
    public override string EventType => "Beneficiary.StatusUpdated";
}
