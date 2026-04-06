using Microsoft.EntityFrameworkCore;
using InsuranceEngine.SharedKernel.Persistence.Entities;

namespace InsuranceEngine.Claims.Infrastructure;

public class ClaimsDbContext : DbContext
{
    public ClaimsDbContext(DbContextOptions<ClaimsDbContext> options) : base(options) { }

    public DbSet<ClaimEntity> Claims => Set<ClaimEntity>();
    public DbSet<ClaimDocumentEntity> ClaimDocuments => Set<ClaimDocumentEntity>();
    public DbSet<ClaimApprovalEntity> ClaimApprovals => Set<ClaimApprovalEntity>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("insurance_schema");

        modelBuilder.Entity<ClaimEntity>(entity =>
        {
            entity.ToTable("claims");
            entity.HasKey(e => e.ClaimId);
            entity.Property(e => e.ClaimId).HasColumnName("claim_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.ClaimNumber).HasColumnName("claim_number").HasMaxLength(50).IsRequired();
            entity.Property(e => e.PolicyId).HasColumnName("policy_id").IsRequired();
            entity.Property(e => e.CustomerId).HasColumnName("customer_id").IsRequired();
            entity.Property(e => e.Status).HasColumnName("status").HasMaxLength(50).HasDefaultValue("SUBMITTED").IsRequired();
            entity.Property(e => e.Type).HasColumnName("type").HasMaxLength(50).IsRequired();
            entity.Property(e => e.ClaimedAmount).HasColumnName("claimed_amount").IsRequired();
            entity.Property(e => e.ClaimedCurrency).HasColumnName("claimed_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.ApprovedAmount).HasColumnName("approved_amount");
            entity.Property(e => e.ApprovedCurrency).HasColumnName("approved_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.SettledAmount).HasColumnName("settled_amount");
            entity.Property(e => e.SettledCurrency).HasColumnName("settled_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.IncidentDate).HasColumnName("incident_date").IsRequired();
            entity.Property(e => e.IncidentDescription).HasColumnName("incident_description");
            entity.Property(e => e.SubmittedAt).HasColumnName("submitted_at").IsRequired();
            entity.Property(e => e.ApprovedAt).HasColumnName("approved_at");
            entity.Property(e => e.SettledAt).HasColumnName("settled_at");
            entity.Property(e => e.RejectionReason).HasColumnName("rejection_reason");
            entity.Property(e => e.PlaceOfIncident).HasColumnName("place_of_incident");
            entity.Property(e => e.BankDetailsForPayout).HasColumnName("bank_details_for_payout");
            entity.Property(e => e.AppealOptionAvailable).HasColumnName("appeal_option_available").HasDefaultValue(true);
            entity.Property(e => e.InAppMessages).HasColumnName("in_app_messages");
            entity.Property(e => e.ProcessingType).HasColumnName("processing_type").HasMaxLength(50).HasDefaultValue("MANUAL");
            entity.Property(e => e.DeductibleAmount).HasColumnName("deductible_amount");
            entity.Property(e => e.DeductibleCurrency).HasColumnName("deductible_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.CoPayAmount).HasColumnName("co_pay_amount");
            entity.Property(e => e.CoPayCurrency).HasColumnName("co_pay_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.ProcessorNotes).HasColumnName("processor_notes");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.DeletedAt).HasColumnName("deleted_at");

            entity.HasIndex(e => e.ClaimNumber).IsUnique().HasDatabaseName("idx_claims_number");
            entity.HasIndex(e => e.PolicyId).HasDatabaseName("idx_claims_policy_id");
            entity.HasIndex(e => e.CustomerId).HasDatabaseName("idx_claims_customer_id");
            entity.HasIndex(e => e.Status).HasDatabaseName("idx_claims_status");

            entity.HasOne(e => e.Policy).WithMany().HasForeignKey(e => e.PolicyId).OnDelete(DeleteBehavior.Restrict);
            entity.HasMany(e => e.Documents).WithOne(d => d.Claim).HasForeignKey(d => d.ClaimId).OnDelete(DeleteBehavior.Cascade);
            entity.HasMany(e => e.Approvals).WithOne(a => a.Claim).HasForeignKey(a => a.ClaimId).OnDelete(DeleteBehavior.Cascade);
        });

        modelBuilder.Entity<ClaimDocumentEntity>(entity =>
        {
            entity.ToTable("claim_documents");
            entity.HasKey(e => e.DocumentId);
            entity.Property(e => e.DocumentId).HasColumnName("document_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.ClaimId).HasColumnName("claim_id").IsRequired();
            entity.Property(e => e.DocumentType).HasColumnName("document_type").HasMaxLength(50).IsRequired();
            entity.Property(e => e.FileUrl).HasColumnName("file_url").HasMaxLength(500).IsRequired();
            entity.Property(e => e.FileHash).HasColumnName("file_hash").HasMaxLength(64);
            entity.Property(e => e.UploadedAt).HasColumnName("uploaded_at").IsRequired();
            entity.Property(e => e.Verified).HasColumnName("verified").HasDefaultValue(false);
            entity.Property(e => e.VerifiedBy).HasColumnName("verified_by");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();

            entity.HasIndex(e => e.ClaimId).HasDatabaseName("idx_claim_documents_claim_id");
        });

        modelBuilder.Entity<ClaimApprovalEntity>(entity =>
        {
            entity.ToTable("claim_approvals");
            entity.HasKey(e => e.ApprovalId);
            entity.Property(e => e.ApprovalId).HasColumnName("approval_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.ClaimId).HasColumnName("claim_id").IsRequired();
            entity.Property(e => e.ApproverId).HasColumnName("approver_id").IsRequired();
            entity.Property(e => e.ApproverRole).HasColumnName("approver_role").HasMaxLength(50).IsRequired();
            entity.Property(e => e.ApprovalLevel).HasColumnName("approval_level").IsRequired();
            entity.Property(e => e.Decision).HasColumnName("decision").HasMaxLength(50).HasDefaultValue("PENDING").IsRequired();
            entity.Property(e => e.ApprovedAmount).HasColumnName("approved_amount");
            entity.Property(e => e.ApprovedCurrency).HasColumnName("approved_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.Notes).HasColumnName("notes");
            entity.Property(e => e.ApprovedAt).HasColumnName("approved_at");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();

            entity.HasIndex(e => e.ClaimId).HasDatabaseName("idx_claim_approvals_claim_id");
        });
    }
}
