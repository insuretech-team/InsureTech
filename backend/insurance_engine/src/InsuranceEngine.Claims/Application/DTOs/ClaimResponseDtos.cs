using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;
using InsuranceEngine.SharedKernel.DTOs;

namespace InsuranceEngine.Claims.Application.DTOs;

public record ClaimSubmissionResponse(
    [property: JsonPropertyName("claim_id")] Guid ClaimId,
    [property: JsonPropertyName("claim_number")] string ClaimNumber,
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record ClaimsDocumentUploadResponse(
    [property: JsonPropertyName("document_id")] Guid DocumentId,
    [property: JsonPropertyName("document_url")] string DocumentUrl,
    [property: JsonPropertyName("file_hash")] string FileHash,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record ClaimRetrievalResponse(
    [property: JsonPropertyName("claim")] ClaimResponseDto Claim,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record UserClaimsListingResponse(
    [property: JsonPropertyName("claims")] List<ClaimListDto> Claims,
    [property: JsonPropertyName("total_count")] int TotalCount,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record ClaimApprovalResponse(
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);
