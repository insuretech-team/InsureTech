using System;

namespace InsuranceEngine.Beneficiary.Domain.Enums;

public enum BeneficiaryType
{
    Individual = 0,
    Business = 1
}

public enum BeneficiaryStatus
{
    PendingKyc = 0,
    Active = 1,
    Inactive = 2,
    Suspended = 3,
    Blacklisted = 4
}

public enum KYCStatus
{
    NotStarted = 0,
    Pending = 1,
    InReview = 2,
    Verified = 3,
    Rejected = 4
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

public enum BusinessType
{
    Proprietorship = 0,
    Partnership = 1,
    PrivateLimited = 2,
    PublicLimited = 3,
    Society = 4,
    Trust = 5
}
