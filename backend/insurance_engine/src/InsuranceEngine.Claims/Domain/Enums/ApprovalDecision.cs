namespace InsuranceEngine.Claims.Domain.Enums;

public enum ApprovalDecision
{
    Unspecified = 0,
    Pending = 1,
    Approved = 2,
    Rejected = 3,
    NeedsMoreInfo = 4,
    Escalated = 5
}
