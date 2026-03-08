using System.ComponentModel.DataAnnotations;

namespace InsuranceEngine.Domain.Entities;

public class FieldViolation
{
    [Key]
    public Guid FieldViolationId { get; set; }

    [Required]
    [MaxLength(200)]
    public string Field { get; set; } = string.Empty;

    [Required]
    [MaxLength(100)]
    public string Code { get; set; } = string.Empty;

    [Required]
    [MaxLength(1000)]
    public string Description { get; set; } = string.Empty;

    [Required]
    [MaxLength(500)]
    public string RejectedValue { get; set; } = string.Empty;

    [Required]
    public Guid ErrorId { get; set; }

    public Error Error { get; set; } = null!;
}
