using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;
using InsuranceEngine.SharedKernel.DTOs;

namespace InsuranceEngine.Policy.Application.DTOs;

public record PolicyCreationResponse(
    [property: JsonPropertyName("policy_id")] Guid PolicyId,
    [property: JsonPropertyName("policy_number")] string PolicyNumber,
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record PolicyRetrievalResponse(
    [property: JsonPropertyName("policy")] PolicyDto Policy,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record PolicyListingResponse(
    [property: JsonPropertyName("policies")] List<PolicyListDto> Policies,
    [property: JsonPropertyName("total_count")] int TotalCount,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record PolicyIssueResponse(
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record PolicyCancelResponse(
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record PolicyRenewResponse(
    [property: JsonPropertyName("new_policy_id")] Guid NewPolicyId,
    [property: JsonPropertyName("new_policy_number")] string NewPolicyNumber,
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record GracePeriodResponse(
    [property: JsonPropertyName("grace_period")] GracePeriodDto GracePeriod,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record RenewalScheduleResponse(
    [property: JsonPropertyName("renewal_schedule")] RenewalScheduleDto RenewalSchedule,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record NomineeResponse(
    [property: JsonPropertyName("nominee_id")] Guid NomineeId,
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record NomineeListingResponse(
    [property: JsonPropertyName("nominees")] List<NomineeDto> Nominees,
    [property: JsonPropertyName("total_count")] int TotalCount,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);
