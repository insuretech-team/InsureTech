using System;

namespace InsuranceEngine.Beneficiary.Domain.Enums;

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
    Inactive = 3,
    Suspended = 4,
    Blacklisted = 5
}

public enum KYCStatus
{
    Unspecified = 0,
    NotStarted = 1,
    Pending = 2,
    InReview = 3,
    Verified = 4,
    Rejected = 5
}

public enum BeneficiaryGender
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
    Unspecified = 0,
    Proprietorship = 1,
    Partnership = 2,
    PrivateLimited = 3,
    PublicLimited = 4,
    Society = 5,
    Trust = 6
}
