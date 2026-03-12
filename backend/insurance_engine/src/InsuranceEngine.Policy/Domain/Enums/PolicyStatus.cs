namespace InsuranceEngine.Policy.Domain.Enums;

public enum PolicyStatus
{
    Unspecified = 0,
    PendingPayment = 1,
    Active = 2,
    GracePeriod = 3,
    Lapsed = 4,
    Suspended = 5,
    Cancelled = 6,
    Expired = 7
}
