namespace InsuranceEngine.Commission.Domain.Enums;

public enum CommissionType
{
    Acquisition,
    Renewal
}

public enum CommissionStatus
{
    Pending,
    Processing,
    Paid
}

public enum PayoutStatus
{
    Pending,
    Paid,
    Cancelled
}
