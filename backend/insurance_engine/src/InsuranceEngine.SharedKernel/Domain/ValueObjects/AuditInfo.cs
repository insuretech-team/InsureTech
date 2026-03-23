using System;

namespace InsuranceEngine.SharedKernel.Domain.ValueObjects;

/// <summary>
/// Value object for audit information.
/// </summary>
public record AuditInfo
{
    public DateTime CreatedAt { get; init; } = DateTime.UtcNow;
    public string? CreatedBy { get; init; }
    public DateTime UpdatedAt { get; init; } = DateTime.UtcNow;
    public string? UpdatedBy { get; init; }
    public DateTime? DeletedAt { get; init; }
    public string? DeletedBy { get; init; }

    public AuditInfo() { }

    public AuditInfo(string? createdBy)
    {
        CreatedAt = DateTime.UtcNow;
        CreatedBy = createdBy;
        UpdatedAt = DateTime.UtcNow;
        UpdatedBy = createdBy;
    }
}
