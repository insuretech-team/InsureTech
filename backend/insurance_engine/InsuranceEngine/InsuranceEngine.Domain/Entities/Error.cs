using System.ComponentModel.DataAnnotations;
using InsuranceEngine.Domain.Enums;

namespace InsuranceEngine.Domain.Entities;

public class Error
{
    [Key]
    public Guid ErrorId { get; set; }

    [Required]
    public ErrorCode Code { get; set; }

    [Required]
    [MaxLength(500)]
    public string Message { get; set; } = string.Empty;

    public Dictionary<string, string?>? Details { get; set; }

    public bool Retryable { get; set; }

    public int? RetryAfterSeconds { get; set; }

    public int HttpStatusCode { get; set; }

    [MaxLength(2048)]
    public string? DocumentationUrl { get; set; }

    public ICollection<FieldViolation> FieldViolations { get; set; } = new List<FieldViolation>();
}
