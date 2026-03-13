using System;

namespace InsuranceEngine.Claims.Domain.Entities;

public class ClaimDocument
{
    public Guid Id { get; set; }
    public Guid ClaimId { get; set; }
    public string DocumentType { get; set; } = string.Empty;
    public string FileUrl { get; set; } = string.Empty;
    public string FileHash { get; set; } = string.Empty;
    public bool Verified { get; set; }
    public Guid? VerifiedBy { get; set; }
    public DateTime UploadedAt { get; set; }
    public DateTime CreatedAt { get; set; }
}
