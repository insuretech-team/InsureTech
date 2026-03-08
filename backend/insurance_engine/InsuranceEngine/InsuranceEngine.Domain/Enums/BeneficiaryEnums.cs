using System.Runtime.Serialization;

namespace InsuranceEngine.Domain.Enums;

public enum BeneficiaryGender
{
    [EnumMember(Value = "GENDER_UNSPECIFIED")]
    GenderUnspecified,
    [EnumMember(Value = "GENDER_MALE")]
    GenderMale,
    [EnumMember(Value = "GENDER_FEMALE")]
    GenderFemale,
    [EnumMember(Value = "GENDER_OTHER")]
    GenderOther
}

public enum BeneficiaryStatus
{
    [EnumMember(Value = "BENEFICIARY_STATUS_UNSPECIFIED")]
    BeneficiaryStatusUnspecified,
    [EnumMember(Value = "BENEFICIARY_STATUS_PENDING_KYC")]
    BeneficiaryStatusPendingKyc,
    [EnumMember(Value = "BENEFICIARY_STATUS_ACTIVE")]
    BeneficiaryStatusActive,
    [EnumMember(Value = "BENEFICIARY_STATUS_SUSPENDED")]
    BeneficiaryStatusSuspended,
    [EnumMember(Value = "BENEFICIARY_STATUS_BLOCKED")]
    BeneficiaryStatusBlocked,
    [EnumMember(Value = "BENEFICIARY_STATUS_DEACTIVATED")]
    BeneficiaryStatusDeactivated
}

public enum BeneficiaryType
{
    [EnumMember(Value = "BENEFICIARY_TYPE_UNSPECIFIED")]
    BeneficiaryTypeUnspecified,
    [EnumMember(Value = "BENEFICIARY_TYPE_INDIVIDUAL")]
    BeneficiaryTypeIndividual,
    [EnumMember(Value = "BENEFICIARY_TYPE_BUSINESS")]
    BeneficiaryTypeBusiness
}

public enum BusinessType
{
    [EnumMember(Value = "BUSINESS_TYPE_UNSPECIFIED")]
    BusinessTypeUnspecified,
    [EnumMember(Value = "BUSINESS_TYPE_SOLE_PROPRIETORSHIP")]
    BusinessTypeSoleProprietorship,
    [EnumMember(Value = "BUSINESS_TYPE_PARTNERSHIP")]
    BusinessTypePartnership,
    [EnumMember(Value = "BUSINESS_TYPE_PRIVATE_LIMITED")]
    BusinessTypePrivateLimited,
    [EnumMember(Value = "BUSINESS_TYPE_PUBLIC_LIMITED")]
    BusinessTypePublicLimited,
    [EnumMember(Value = "BUSINESS_TYPE_NGO")]
    BusinessTypeNgo,
    [EnumMember(Value = "BUSINESS_TYPE_GOVERNMENT")]
    BusinessTypeGovernment
}

public enum MaritalStatus
{
    [EnumMember(Value = "MARITAL_STATUS_UNSPECIFIED")]
    MaritalStatusUnspecified,
    [EnumMember(Value = "MARITAL_STATUS_SINGLE")]
    MaritalStatusSingle,
    [EnumMember(Value = "MARITAL_STATUS_MARRIED")]
    MaritalStatusMarried,
    [EnumMember(Value = "MARITAL_STATUS_DIVORCED")]
    MaritalStatusDivorced,
    [EnumMember(Value = "MARITAL_STATUS_WIDOWED")]
    MaritalStatusWidowed
}
