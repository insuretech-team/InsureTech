using System.Runtime.Serialization;

namespace InsuranceEngine.Domain.Enums;

public enum ErrorCode
{
    [EnumMember(Value = "ERROR_CODE_UNSPECIFIED")]
    ErrorCodeUnspecified,

    [EnumMember(Value = "INVALID_REQUEST")]
    InvalidRequest,

    [EnumMember(Value = "INVALID_PAGINATION")]
    InvalidPagination,

    [EnumMember(Value = "INVALID_BENEFICIARY_ID")]
    InvalidBeneficiaryId,

    [EnumMember(Value = "BENEFICIARY_ID_MISMATCH")]
    BeneficiaryIdMismatch,

    [EnumMember(Value = "BENEFICIARY_NOT_FOUND")]
    BeneficiaryNotFound
}
