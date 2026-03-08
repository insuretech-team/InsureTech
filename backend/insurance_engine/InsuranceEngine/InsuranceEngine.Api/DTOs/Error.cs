using System.Text.Json.Serialization;
using Microsoft.AspNetCore.Http;

namespace InsuranceEngine.Api.DTOs;

public class Error
{
    [JsonPropertyName("code")]
    public string Code { get; set; } = string.Empty;

    [JsonPropertyName("message")]
    public string Message { get; set; } = string.Empty;

    [JsonPropertyName("details")]
    public IDictionary<string, object?>? Details { get; set; }

    [JsonPropertyName("field_violations")]
    public IReadOnlyList<FieldViolation>? FieldViolations { get; set; }

    [JsonPropertyName("retryable")]
    public bool Retryable { get; set; }

    [JsonPropertyName("retry_after_seconds")]
    public int? RetryAfterSeconds { get; set; }

    [JsonPropertyName("http_status_code")]
    public int HttpStatusCode { get; set; }

    [JsonPropertyName("error_id")]
    public string ErrorId { get; set; } = string.Empty;

    [JsonPropertyName("documentation_url")]
    public string? DocumentationUrl { get; set; }

    public static Error Create(
        string code,
        string message,
        int httpStatusCode,
        bool retryable = false,
        int? retryAfterSeconds = null,
        IDictionary<string, object?>? details = null,
        IReadOnlyList<FieldViolation>? fieldViolations = null,
        string? documentationUrl = null)
    {
        return new Error
        {
            Code = code,
            Message = message,
            Details = details,
            FieldViolations = fieldViolations,
            Retryable = retryable,
            RetryAfterSeconds = retryAfterSeconds,
            HttpStatusCode = httpStatusCode,
            ErrorId = Guid.NewGuid().ToString(),
            DocumentationUrl = documentationUrl
        };
    }

    public static Error None(int httpStatusCode = StatusCodes.Status200OK)
    {
        return Create("NONE", "No error.", httpStatusCode);
    }
}

public class FieldViolation
{
    [JsonPropertyName("field")]
    public string Field { get; set; } = string.Empty;

    [JsonPropertyName("code")]
    public string Code { get; set; } = string.Empty;

    [JsonPropertyName("description")]
    public string Description { get; set; } = string.Empty;

    [JsonPropertyName("rejected_value")]
    public string RejectedValue { get; set; } = string.Empty;
}
