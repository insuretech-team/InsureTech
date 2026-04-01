using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'claim_documents' table in insurance_schema.
/// Aligned with the database schema managed by Go migrations.
/// </summary>
[Table("claim_documents", Schema = "insurance_schema")]
public class ClaimDocumentEntity
{
    [Key]
    [Column("document_id")]
    public Guid DocumentId { get; set; }

    [Column("claim_id")]
    public Guid ClaimId { get; set; }

    [Column("document_type")]
    public string DocumentType { get; set; } = string.Empty;

    [Column("file_url")]
    public string FileUrl { get; set; } = string.Empty;

    [Column("file_hash")]
    public string FileHash { get; set; } = string.Empty; // SHA-256

    [Column("uploaded_at")]
    public DateTime UploadedAt { get; set; }

    [Column("verified")]
    public bool Verified { get; set; }

    [Column("verified_by")]
    public Guid? VerifiedBy { get; set; }

    [Column("created_at")]
    public DateTime CreatedAt { get; set; }

    [Column("updated_at")]
    public DateTime UpdatedAt { get; set; }

    // Navigation
    public ClaimEntity Claim { get; set; } = null!;
}
