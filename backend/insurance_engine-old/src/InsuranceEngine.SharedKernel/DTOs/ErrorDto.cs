using System.Collections.Generic;

namespace InsuranceEngine.SharedKernel.DTOs;

public record ErrorDto(
    string Code,
    string Message,
    Dictionary<string, string>? Details = null,
    List<FieldViolationDto>? FieldViolations = null,
    bool Retryable = false,
    int? RetryAfterSeconds = null,
    int? HttpStatusCode = null,
    string? ErrorId = null,
    string? DocumentationUrl = null
);

public record FieldViolationDto(
    string Field,
    string Code,
    string Description,
    string? RejectedValue = null
);
