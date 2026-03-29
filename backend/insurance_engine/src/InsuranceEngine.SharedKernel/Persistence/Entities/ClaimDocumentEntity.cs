namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'claim_documents' table in insurance_schema.
/// Aligned with insuretech.claims.entity.v1.ClaimDocument proto definition.
/// </summary>
public class ClaimDocumentEntity
{
    public Guid DocumentId { get; set; }
    public Guid ClaimId { get; set; }
    public string DocumentType { get; set; } = string.Empty;
    public string FileUrl { get; set; } = string.Empty;
    public string FileHash { get; set; } = string.Empty; // SHA-256
    public DateTime UploadedAt { get; set; }
    public bool Verified { get; set; }
    public Guid? VerifiedBy { get; set; }
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }

    // Navigation
    public ClaimEntity Claim { get; set; } = null!;
}
