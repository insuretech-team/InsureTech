using Microsoft.EntityFrameworkCore;
using InsuranceEngine.SharedKernel.Persistence.Entities;

namespace InsuranceEngine.Endorsements.Infrastructure;

public class EndorsementsDbContext : DbContext
{
    public EndorsementsDbContext(DbContextOptions<EndorsementsDbContext> options) : base(options) { }

    public DbSet<EndorsementEntity> Endorsements => Set<EndorsementEntity>();
    public DbSet<EndorsementDocumentEntity> EndorsementDocuments => Set<EndorsementDocumentEntity>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("insurance_schema");

        modelBuilder.Entity<EndorsementEntity>(entity =>
        {
            entity.ToTable("endorsements");
            entity.HasKey(e => e.EndorsementId);
            entity.Property(e => e.EndorsementId).HasColumnName("endorsement_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.EndorsementNumber).HasColumnName("endorsement_number").HasMaxLength(50).IsRequired();
            entity.Property(e => e.PolicyId).HasColumnName("policy_id").IsRequired();
            entity.Property(e => e.Type).HasColumnName("type").HasMaxLength(50).IsRequired();
            entity.Property(e => e.Status).HasColumnName("status").HasMaxLength(50).HasDefaultValue("PENDING").IsRequired();
            entity.Property(e => e.Reason).HasColumnName("reason");
            entity.Property(e => e.Changes).HasColumnName("changes");
            entity.Property(e => e.OldSumAssured).HasColumnName("old_sum_assured");
            entity.Property(e => e.NewSumAssured).HasColumnName("new_sum_assured");
            entity.Property(e => e.RefundAmount).HasColumnName("refund_amount");
            entity.Property(e => e.AdditionalPremium).HasColumnName("additional_premium");
            entity.Property(e => e.EffectiveDate).HasColumnName("effective_date");
            entity.Property(e => e.ApprovedBy).HasColumnName("approved_by");
            entity.Property(e => e.RejectedBy).HasColumnName("rejected_by");
            entity.Property(e => e.RejectionReason).HasColumnName("rejection_reason");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();

            entity.HasIndex(e => e.EndorsementNumber).IsUnique().HasDatabaseName("idx_endorsements_number");
            entity.HasIndex(e => e.PolicyId).HasDatabaseName("idx_endorsements_policy_id");
            entity.HasIndex(e => e.Status).HasDatabaseName("idx_endorsements_status");
        });

        modelBuilder.Entity<EndorsementDocumentEntity>(entity =>
        {
            entity.ToTable("endorsement_documents");
            entity.HasKey(e => e.DocumentId);
            entity.Property(e => e.DocumentId).HasColumnName("document_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.EndorsementId).HasColumnName("endorsement_id").IsRequired();
            entity.Property(e => e.DocumentType).HasColumnName("document_type").HasMaxLength(50).IsRequired();
            entity.Property(e => e.DocumentNumber).HasColumnName("document_number").HasMaxLength(50);
            entity.Property(e => e.FileUrl).HasColumnName("file_url").HasMaxLength(500);
            entity.Property(e => e.Status).HasColumnName("status").HasMaxLength(50).HasDefaultValue("PENDING").IsRequired();
            entity.Property(e => e.GeneratedAt).HasColumnName("generated_at");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();

            entity.HasIndex(e => e.EndorsementId).HasDatabaseName("idx_endorsement_documents_endorsement_id");
        });
    }
}
