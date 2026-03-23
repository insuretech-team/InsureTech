using InsuranceEngine.Policy.Domain.Entities;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.SharedKernel.Interfaces;
using Microsoft.EntityFrameworkCore;
using System;

namespace InsuranceEngine.Policy.Infrastructure.Persistence;

public class PolicyDbContext : DbContext
{
    private readonly Guid _tenantId;

    public PolicyDbContext(DbContextOptions<PolicyDbContext> options, ITenantService tenantService) : base(options)
    {
        _tenantId = tenantService.GetTenantId();
    }

    public DbSet<PolicyAggregate> Policies { get; set; } = null!;
    public DbSet<Nominee> Nominees { get; set; } = null!;
    public DbSet<PolicyRider> Riders { get; set; } = null!; // Renamed from PolicyRiders
    public DbSet<Endorsement> Endorsements { get; set; } = null!;



    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("insurance"); // Changed schema name

        // --- Policy ---
        modelBuilder.Entity<PolicyAggregate>(entity => // Changed from PolicyEntity to PolicyAggregate
        {
            entity.ToTable("policies");
            entity.HasKey(e => e.Id);
            entity.HasQueryFilter(e => !e.IsDeleted); // Kept from PolicyEntity

            entity.Property(e => e.PolicyNumber).HasMaxLength(50).IsRequired();
            entity.HasIndex(e => e.PolicyNumber).IsUnique();

            entity.Property(e => e.Status).HasConversion<string>().HasMaxLength(50).IsRequired();

            // Money columns as bigint
            entity.Property(e => e.PremiumAmount).HasColumnName("premium_amount").IsRequired();
            entity.Property(e => e.PremiumCurrency).HasColumnName("premium_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.SumInsuredAmount).HasColumnName("sum_insured_amount").IsRequired();
            entity.Property(e => e.SumInsuredCurrency).HasColumnName("sum_insured_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.VatTaxAmount).HasColumnName("vat_tax_amount");
            entity.Property(e => e.ServiceFeeAmount).HasColumnName("service_fee_amount");
            entity.Property(e => e.TotalPayableAmount).HasColumnName("total_payable_amount");

            entity.Property(e => e.ProposerDetailsJson).HasColumnType("jsonb").HasColumnName("proposer_details");
            entity.Property(e => e.UnderwritingData).HasColumnType("jsonb").HasColumnName("underwriting_data");

            entity.Ignore(e => e.PremiumMoney);
            entity.Ignore(e => e.SumInsuredMoney);

            entity.HasIndex(e => e.CustomerId);
            entity.HasIndex(e => e.ProductId);
            entity.HasIndex(e => e.Status);

            entity.HasMany(e => e.Nominees)
                  .WithOne()
                  .HasForeignKey(n => n.PolicyId)
                  .OnDelete(DeleteBehavior.Cascade);

            entity.HasMany(e => e.Riders)
                  .WithOne()
                  .HasForeignKey(r => r.PolicyId)
                  .OnDelete(DeleteBehavior.Cascade);
        });


        // --- PolicyRider ---
        modelBuilder.Entity<PolicyRider>(entity =>
        {
            entity.ToTable("policy_riders");
            entity.HasKey(e => e.Id);

            entity.Property(e => e.RiderName).HasMaxLength(255).IsRequired();
            entity.Property(e => e.PremiumAmount).HasColumnName("premium_amount").IsRequired();
            entity.Property(e => e.PremiumCurrency).HasColumnName("premium_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.CoverageAmount).HasColumnName("coverage_amount").IsRequired();
            entity.Property(e => e.CoverageCurrency).HasColumnName("coverage_currency").HasMaxLength(3).HasDefaultValue("BDT");

            entity.Ignore(e => e.Premium);
            entity.Ignore(e => e.Coverage);

            entity.HasIndex(e => e.PolicyId);
        });

        modelBuilder.Entity<Nominee>(entity =>
        {
            entity.ToTable("policy_nominees", "insurance_schema");
            entity.HasKey(e => e.Id);
            entity.HasQueryFilter(e => !e.IsDeleted);

            entity.Property(e => e.FullName).HasColumnName("full_name").HasMaxLength(200).IsRequired();
            entity.Property(e => e.Relationship).HasMaxLength(50).IsRequired();
            entity.Property(e => e.DateOfBirth).HasColumnName("date_of_birth");
            entity.Property(e => e.NomineeDobText).HasColumnName("nominee_dob_text").HasMaxLength(50);
            entity.Property(e => e.NidNumber).HasColumnName("nid_number").HasMaxLength(20);
            entity.Property(e => e.PhoneNumber).HasColumnName("phone_number").HasMaxLength(20);

            // Navigation removed for decoupling

            entity.HasIndex(e => e.PolicyId);
            entity.HasIndex(e => e.NidNumber);
        });

        // --- Endorsement ---
        modelBuilder.Entity<Endorsement>(entity =>
        {
            entity.ToTable("endorsements");
            entity.HasKey(e => e.Id);
            entity.Property(e => e.EndorsementNumber).IsRequired().HasMaxLength(50);
            entity.Property(e => e.Type).HasConversion<string>().HasMaxLength(50).IsRequired();
            entity.Property(e => e.Status).HasConversion<string>().HasMaxLength(50).IsRequired();
            entity.Property(e => e.Reason).HasMaxLength(500);
            entity.Property(e => e.ChangesJson).HasColumnType("jsonb").HasColumnName("changes");
            entity.Property(e => e.AuditInfoJson).HasColumnType("jsonb").HasColumnName("audit_info");
            
            entity.Property(e => e.PremiumAdjustmentAmount).HasColumnName("premium_adjustment_amount");
            entity.Property(e => e.PremiumAdjustmentCurrency).HasColumnName("premium_adjustment_currency").HasMaxLength(3).HasDefaultValue("BDT");

            entity.HasIndex(e => e.EndorsementNumber).IsUnique();
            entity.HasIndex(e => e.PolicyId);
            entity.HasIndex(e => e.Status);
        });

        // --- DB Sequences ---

        modelBuilder.HasSequence<long>("policy_number_seq", "insurance_schema")
            .StartsAt(1)
            .IncrementsBy(1);

        modelBuilder.HasSequence<long>("endorsement_number_seq", "insurance_schema")
            .StartsAt(1)
            .IncrementsBy(1);

    }
}
