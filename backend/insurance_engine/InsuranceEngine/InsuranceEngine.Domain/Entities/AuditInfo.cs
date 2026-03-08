using System.ComponentModel.DataAnnotations;

namespace InsuranceEngine.Domain.Entities;

public class AuditInfo
{
    [Required]
    public DateTime CreatedAt { get; set; }

    [Required]
    public DateTime UpdatedAt { get; set; }

    [MaxLength(100)]
    [Required]
    public string CreatedBy { get; set; } = string.Empty;

    [MaxLength(100)]
    [Required]
    public string UpdatedBy { get; set; } = string.Empty;

    public DateTime? DeletedAt { get; set; }

    [MaxLength(100)]
    public string? DeletedBy { get; set; }
}
