using Microsoft.EntityFrameworkCore;
using InsuranceEngine.SharedKernel.Persistence.Entities;

namespace InsuranceEngine.FraudDetection.Infrastructure;

public class FraudDetectionDbContext : DbContext
{
    public FraudDetectionDbContext(DbContextOptions<FraudDetectionDbContext> options) : base(options) { }

    public DbSet<FraudCheckEntity> FraudChecks => Set<FraudCheckEntity>();
    public DbSet<FraudAlertEntity> FraudAlerts => Set<FraudAlertEntity>();
    public DbSet<FraudDashboardSummaryEntity> FraudDashboardSummaries => Set<FraudDashboardSummaryEntity>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("insurance_schema");

        modelBuilder.Entity<FraudCheckEntity>(entity =>
        {
            entity.ToTable("fraud_checks");
            entity.HasKey(e => e.FraudCheckId);
            entity.Property(e => e.FraudCheckId).HasColumnName("fraud_check_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.EntityId).HasColumnName("entity_id").IsRequired();
            entity.Property(e => e.EntityType).HasColumnName("entity_type").HasMaxLength(50).IsRequired();
            entity.Property(e => e.CheckType).HasColumnName("check_type").HasMaxLength(100);
            entity.Property(e => e.FraudScore).HasColumnName("fraud_score").IsRequired();
            entity.Property(e => e.RiskLevel).HasColumnName("risk_level").HasMaxLength(20).IsRequired();
            entity.Property(e => e.Flagged).HasColumnName("flagged").IsRequired();
            entity.Property(e => e.ClaimId).HasColumnName("claim_id");
            entity.Property(e => e.CustomerId).HasColumnName("customer_id");
            entity.Property(e => e.ClaimType).HasColumnName("claim_type").HasMaxLength(50);
            entity.Property(e => e.ClaimAmount).HasColumnName("claim_amount");
            entity.Property(e => e.Recommendation).HasColumnName("recommendation");
            entity.Property(e => e.CheckedAt).HasColumnName("checked_at").IsRequired();
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();

            if (Database.ProviderName == "Npgsql.EntityFrameworkCore.PostgreSQL")
            {
                entity.Property(e => e.RiskFactors).HasColumnName("risk_factors").HasColumnType("TEXT[]");
            }

            entity.HasIndex(e => e.EntityId).HasDatabaseName("idx_fraud_checks_entity_id");
            entity.HasIndex(e => e.ClaimId).HasDatabaseName("idx_fraud_checks_claim_id");
            entity.HasIndex(e => e.Flagged).HasDatabaseName("idx_fraud_checks_flagged");
        });

        modelBuilder.Entity<FraudAlertEntity>(entity =>
        {
            entity.ToTable("fraud_alerts");
            entity.HasKey(e => e.AlertId);
            entity.Property(e => e.AlertId).HasColumnName("alert_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.AlertNumber).HasColumnName("alert_number").HasMaxLength(50).IsRequired();
            entity.Property(e => e.EntityType).HasColumnName("entity_type").HasMaxLength(50).IsRequired();
            entity.Property(e => e.EntityId).HasColumnName("entity_id").IsRequired();
            entity.Property(e => e.AlertType).HasColumnName("alert_type").HasMaxLength(100).IsRequired();
            entity.Property(e => e.Severity).HasColumnName("severity").HasMaxLength(20).IsRequired();
            entity.Property(e => e.Status).HasColumnName("status").HasMaxLength(50).HasDefaultValue("OPEN").IsRequired();
            entity.Property(e => e.Description).HasColumnName("description");
            entity.Property(e => e.FraudScore).HasColumnName("fraud_score").IsRequired();
            entity.Property(e => e.RecommendedAction).HasColumnName("recommended_action");
            entity.Property(e => e.ResolvedBy).HasColumnName("resolved_by");
            entity.Property(e => e.ResolutionNotes).HasColumnName("resolution_notes");
            entity.Property(e => e.ResolvedAt).HasColumnName("resolved_at");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();

            entity.HasIndex(e => e.AlertNumber).IsUnique().HasDatabaseName("idx_fraud_alerts_number");
            entity.HasIndex(e => e.EntityId).HasDatabaseName("idx_fraud_alerts_entity_id");
            entity.HasIndex(e => e.Status).HasDatabaseName("idx_fraud_alerts_status");
            entity.HasIndex(e => e.Severity).HasDatabaseName("idx_fraud_alerts_severity");
        });

        modelBuilder.Entity<FraudDashboardSummaryEntity>(entity =>
        {
            entity.ToTable("fraud_dashboard_summaries");
            entity.HasKey(e => e.SummaryId);
            entity.Property(e => e.SummaryId).HasColumnName("summary_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.SummaryDate).HasColumnName("summary_date").IsRequired();
            entity.Property(e => e.TotalFlagsToday).HasColumnName("total_flags_today").IsRequired();
            entity.Property(e => e.HighRiskFlags).HasColumnName("high_risk_flags").IsRequired();
            entity.Property(e => e.MediumRiskFlags).HasColumnName("medium_risk_flags").IsRequired();
            entity.Property(e => e.LowRiskFlags).HasColumnName("low_risk_flags").IsRequired();
            entity.Property(e => e.PendingReviewCount).HasColumnName("pending_review_count").IsRequired();
            entity.Property(e => e.ResolvedCount).HasColumnName("resolved_count").IsRequired();
            entity.Property(e => e.AverageFraudScore).HasColumnName("average_fraud_score").IsRequired();
            entity.Property(e => e.TopFraudTypes).HasColumnName("top_fraud_types");
            entity.Property(e => e.GeneratedAt).HasColumnName("generated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();

            entity.HasIndex(e => e.SummaryDate).IsUnique().HasDatabaseName("idx_fraud_dashboard_summaries_date");
        });
    }
}
