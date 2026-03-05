namespace PoliSync.Beneficiaries.Domain;

public enum BeneficiaryType
{
    Unspecified = 0,
    Individual = 1,
    Business = 2
}

public enum BeneficiaryStatus
{
    Unspecified = 0,
    PendingKyc = 1,
    Active = 2,
    Suspended = 3,
    Blocked = 4,
    Closed = 5
}

public enum KycStatus
{
    Unspecified = 0,
    NotStarted = 1,
    InProgress = 2,
    Completed = 3,
    Failed = 4,
    Expired = 5,
    Rejected = 6
}

public enum Gender
{
    Unspecified = 0,
    Male = 1,
    Female = 2,
    Other = 3
}

public enum MaritalStatus
{
    Unspecified = 0,
    Single = 1,
    Married = 2,
    Divorced = 3,
    Widowed = 4
}
