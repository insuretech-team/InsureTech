namespace InsuranceEngine.Underwriting.Domain.Enums;

public enum BeneficiaryStatus
{
    Unspecified,
    Active,
    Inactive,
    Suspended,
    Blacklisted
}

public enum BeneficiaryType
{
    Unspecified,
    Individual,
    Business
}

public enum KYCStatus
{
    Unspecified,
    Pending,
    Verified,
    Rejected,
    Expired
}

public enum Gender
{
    Unspecified,
    Male,
    Female,
    Other
}

public enum MaritalStatus
{
    Unspecified,
    Single,
    Married,
    Divorced,
    Widowed
}

public enum BusinessType
{
    Unspecified,
    Proprietorship,
    Partnership,
    PrivateLimited,
    PublicLimited,
    Ngo
}
