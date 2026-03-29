using Microsoft.EntityFrameworkCore;
using InsuranceEngine.SharedKernel.Persistence.Entities;

namespace InsuranceEngine.SharedKernel.Persistence;

/// <summary>
/// Typed DbContext for the Insurance Engine.
/// Maps to 'insurance_schema' in PostgreSQL.
/// Replaces the bare DbContext used previously.
/// </summary>
public class InsuranceDbContext : DbContext
{
    public InsuranceDbContext(DbContextOptions<InsuranceDbContext> options) : base(options) { }

    // ===== Core Tables =====
    public DbSet<ProductEntity> Products => Set<ProductEntity>();
    public DbSet<PolicyEntity> Policies => Set<PolicyEntity>();
    public DbSet<PolicyNomineeEntity> PolicyNominees => Set<PolicyNomineeEntity>();
    public DbSet<PolicyRiderEntity> PolicyRiders => Set<PolicyRiderEntity>();
    public DbSet<ClaimEntity> Claims => Set<ClaimEntity>();
    public DbSet<ClaimDocumentEntity> ClaimDocuments => Set<ClaimDocumentEntity>();
    public DbSet<ClaimApprovalEntity> ClaimApprovals => Set<ClaimApprovalEntity>();
    public DbSet<FraudCheckEntity> FraudChecks => Set<FraudCheckEntity>();
    public DbSet<ProductRiderEntity> ProductRiders => Set<ProductRiderEntity>();
    public DbSet<QuoteEntity> Quotes => Set<QuoteEntity>();
    public DbSet<HealthDeclarationEntity> HealthDeclarations => Set<HealthDeclarationEntity>();
    public DbSet<UnderwritingDecisionEntity> UnderwritingDecisions => Set<UnderwritingDecisionEntity>();
    public DbSet<CommissionEntity> Commissions => Set<CommissionEntity>();
    public DbSet<CommissionPayoutEntity> CommissionPayouts => Set<CommissionPayoutEntity>();
    public DbSet<BeneficiaryEntity> Beneficiaries => Set<BeneficiaryEntity>();
    public DbSet<IndividualBeneficiaryEntity> IndividualBeneficiaries => Set<IndividualBeneficiaryEntity>();
    public DbSet<BusinessBeneficiaryEntity> BusinessBeneficiaries => Set<BusinessBeneficiaryEntity>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);

        // Use insurance_schema as default schema
        modelBuilder.HasDefaultSchema("insurance_schema");

        // ===== Product Configuration =====
        modelBuilder.Entity<ProductEntity>(entity =>
        {
            entity.ToTable("products");
            entity.HasKey(e => e.ProductId);
            entity.Property(e => e.ProductId).HasColumnName("product_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.ProductCode).HasColumnName("product_code").HasMaxLength(50).IsRequired();
            entity.Property(e => e.ProductName).HasColumnName("product_name").HasMaxLength(255).IsRequired();
            entity.Property(e => e.Category).HasColumnName("category").HasMaxLength(50).IsRequired();
            entity.Property(e => e.Description).HasColumnName("description");
            entity.Property(e => e.BasePremium).HasColumnName("base_premium").IsRequired();
            entity.Property(e => e.BasePremiumCurrency).HasColumnName("base_premium_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
            entity.Property(e => e.MinSumInsured).HasColumnName("min_sum_insured").IsRequired();
            entity.Property(e => e.MinSumInsuredCurrency).HasColumnName("min_sum_insured_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
            entity.Property(e => e.MaxSumInsured).HasColumnName("max_sum_insured").IsRequired();
            entity.Property(e => e.MaxSumInsuredCurrency).HasColumnName("max_sum_insured_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
            entity.Property(e => e.MaxTenureMonths).HasColumnName("max_tenure_months").IsRequired();
            entity.Property(e => e.MinAge).HasColumnName("min_age").IsRequired();
            entity.Property(e => e.MaxAge).HasColumnName("max_age").IsRequired();
            entity.Property(e => e.TermsUrl).HasColumnName("terms_url").HasMaxLength(500);
            entity.Property(e => e.Questions).HasColumnName("questions").HasColumnType("JSONB");
            entity.Property(e => e.Version).HasColumnName("version").HasDefaultValue(1).IsRequired();
            entity.Property(e => e.Exclusions).HasColumnName("exclusions").HasColumnType("TEXT[]");
            entity.Property(e => e.Status).HasColumnName("status").HasMaxLength(50).HasDefaultValue("DRAFT").IsRequired();
            entity.Property(e => e.ProductAttributes).HasColumnName("product_attributes").HasColumnType("JSONB");
            entity.Property(e => e.CreatedBy).HasColumnName("created_by").IsRequired();
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.DeletedAt).HasColumnName("deleted_at");

            // Indexes
            entity.HasIndex(e => e.ProductCode).IsUnique().HasDatabaseName("idx_products_product_code");
            entity.HasIndex(e => e.Category).HasDatabaseName("idx_products_category");
            entity.HasIndex(e => e.Status).HasDatabaseName("idx_products_status");

            // Soft delete filter
            entity.HasQueryFilter(e => e.DeletedAt == null);
        });

        // ===== Policy Configuration =====
        modelBuilder.Entity<PolicyEntity>(entity =>
        {
            entity.ToTable("policies");
            entity.HasKey(e => e.PolicyId);
            entity.Property(e => e.PolicyId).HasColumnName("policy_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.PolicyNumber).HasColumnName("policy_number").HasMaxLength(50).IsRequired();
            entity.Property(e => e.ProductId).HasColumnName("product_id").IsRequired();
            entity.Property(e => e.CustomerId).HasColumnName("customer_id").IsRequired();
            entity.Property(e => e.PartnerId).HasColumnName("partner_id");
            entity.Property(e => e.AgentId).HasColumnName("agent_id");
            entity.Property(e => e.QuoteId).HasColumnName("quote_id");
            entity.Property(e => e.UnderwritingDecisionId).HasColumnName("underwriting_decision_id");
            entity.Property(e => e.Status).HasColumnName("status").HasMaxLength(50).HasDefaultValue("PENDING_PAYMENT").IsRequired();
            entity.Property(e => e.PremiumAmount).HasColumnName("premium_amount").IsRequired();
            entity.Property(e => e.PremiumCurrency).HasColumnName("premium_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
            entity.Property(e => e.SumInsured).HasColumnName("sum_insured").IsRequired();
            entity.Property(e => e.SumInsuredCurrency).HasColumnName("sum_insured_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
            entity.Property(e => e.TenureMonths).HasColumnName("tenure_months").IsRequired();
            entity.Property(e => e.StartDate).HasColumnName("start_date").HasColumnType("DATE").IsRequired();
            entity.Property(e => e.EndDate).HasColumnName("end_date").HasColumnType("DATE").IsRequired();
            entity.Property(e => e.IssuedAt).HasColumnName("issued_at");
            entity.Property(e => e.PolicyDocumentUrl).HasColumnName("policy_document_url");
            entity.Property(e => e.PaymentFrequency).HasColumnName("payment_frequency").HasMaxLength(50);
            entity.Property(e => e.VatTax).HasColumnName("vat_tax");
            entity.Property(e => e.ServiceFee).HasColumnName("service_fee");
            entity.Property(e => e.TotalPayable).HasColumnName("total_payable");
            entity.Property(e => e.PaymentGatewayReference).HasColumnName("payment_gateway_reference").HasMaxLength(255);
            entity.Property(e => e.ReceiptNumber).HasColumnName("receipt_number").HasMaxLength(100);
            entity.Property(e => e.OccupationRiskClass).HasColumnName("occupation_risk_class").HasMaxLength(50);
            entity.Property(e => e.HasExistingPolicies).HasColumnName("has_existing_policies").HasDefaultValue(false);
            entity.Property(e => e.ClaimsHistorySummary).HasColumnName("claims_history_summary");
            entity.Property(e => e.ProviderName).HasColumnName("provider_name").HasMaxLength(255);
            entity.Property(e => e.EnrollmentStartDate).HasColumnName("enrollment_start_date").HasColumnType("DATE");
            entity.Property(e => e.EnrollmentEndDate).HasColumnName("enrollment_end_date").HasColumnType("DATE");
            entity.Property(e => e.UnderwritingData).HasColumnName("underwriting_data").HasColumnType("JSONB");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.DeletedAt).HasColumnName("deleted_at");

            // Indexes
            entity.HasIndex(e => e.PolicyNumber).IsUnique().HasDatabaseName("idx_policies_policy_number");
            entity.HasIndex(e => e.ProductId).HasDatabaseName("idx_policies_product_id");
            entity.HasIndex(e => e.CustomerId).HasDatabaseName("idx_policies_customer_id");
            entity.HasIndex(e => e.PartnerId).HasDatabaseName("idx_policies_partner_id");
            entity.HasIndex(e => e.AgentId).HasDatabaseName("idx_policies_agent_id");
            entity.HasIndex(e => e.QuoteId).HasDatabaseName("idx_policies_quote_id");
            entity.HasIndex(e => e.Status).HasDatabaseName("idx_policies_status");
            entity.HasIndex(e => e.StartDate).HasDatabaseName("idx_policies_start_date");
            entity.HasIndex(e => e.EndDate).HasDatabaseName("idx_policies_end_date");

            // Relationships
            entity.HasOne(e => e.Product).WithMany(p => p.Policies).HasForeignKey(e => e.ProductId).OnDelete(DeleteBehavior.Restrict);

            // Soft delete filter
            entity.HasQueryFilter(e => e.DeletedAt == null);
        });

        // ===== Policy Nominee Configuration =====
        modelBuilder.Entity<PolicyNomineeEntity>(entity =>
        {
            entity.ToTable("policy_nominees");
            entity.HasKey(e => e.NomineeId);
            entity.Property(e => e.NomineeId).HasColumnName("nominee_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.PolicyId).HasColumnName("policy_id").IsRequired();
            entity.Property(e => e.FullName).HasColumnName("full_name").HasMaxLength(255).IsRequired();
            entity.Property(e => e.Relationship).HasColumnName("relationship").HasMaxLength(50).IsRequired();
            entity.Property(e => e.SharePercentage).HasColumnName("share_percentage").HasColumnType("DECIMAL(5,2)").IsRequired();
            entity.Property(e => e.DateOfBirth).HasColumnName("date_of_birth").HasColumnType("DATE").IsRequired();
            entity.Property(e => e.NidNumber).HasColumnName("nid_number").HasMaxLength(17);
            entity.Property(e => e.PhoneNumber).HasColumnName("phone_number").HasMaxLength(20);
            entity.Property(e => e.NomineeDobText).HasColumnName("nominee_dob_text").HasMaxLength(15);
            entity.Property(e => e.NomineeSharePercent).HasColumnName("nominee_share_percent").HasColumnType("DECIMAL(5,2)");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();

            entity.HasIndex(e => e.PolicyId).HasDatabaseName("idx_policy_nominees_policy_id");
            entity.HasOne(e => e.Policy).WithMany(p => p.Nominees).HasForeignKey(e => e.PolicyId).OnDelete(DeleteBehavior.Cascade);
        });

        // ===== Policy Rider Configuration =====
        modelBuilder.Entity<PolicyRiderEntity>(entity =>
        {
            entity.ToTable("policy_riders");
            entity.HasKey(e => e.RiderId);
            entity.Property(e => e.RiderId).HasColumnName("rider_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.PolicyId).HasColumnName("policy_id").IsRequired();
            entity.Property(e => e.RiderName).HasColumnName("rider_name").HasMaxLength(255).IsRequired();
            entity.Property(e => e.PremiumAmount).HasColumnName("premium_amount").IsRequired();
            entity.Property(e => e.PremiumCurrency).HasColumnName("premium_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
            entity.Property(e => e.CoverageAmount).HasColumnName("coverage_amount").IsRequired();
            entity.Property(e => e.CoverageCurrency).HasColumnName("coverage_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();

            entity.HasIndex(e => e.PolicyId).HasDatabaseName("idx_policy_riders_policy_id");
            entity.HasOne(e => e.Policy).WithMany(p => p.Riders).HasForeignKey(e => e.PolicyId).OnDelete(DeleteBehavior.Cascade);
        });

        // ===== Claim Configuration =====
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
            entity.Property(e => e.ClaimedCurrency).HasColumnName("claimed_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
            entity.Property(e => e.ApprovedAmount).HasColumnName("approved_amount");
            entity.Property(e => e.ApprovedCurrency).HasColumnName("approved_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
            entity.Property(e => e.SettledAmount).HasColumnName("settled_amount");
            entity.Property(e => e.SettledCurrency).HasColumnName("settled_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
            entity.Property(e => e.IncidentDate).HasColumnName("incident_date").HasColumnType("DATE").IsRequired();
            entity.Property(e => e.IncidentDescription).HasColumnName("incident_description").IsRequired();
            entity.Property(e => e.SubmittedAt).HasColumnName("submitted_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.ApprovedAt).HasColumnName("approved_at");
            entity.Property(e => e.SettledAt).HasColumnName("settled_at");
            entity.Property(e => e.RejectionReason).HasColumnName("rejection_reason");
            entity.Property(e => e.PlaceOfIncident).HasColumnName("place_of_incident");
            entity.Property(e => e.BankDetailsForPayout).HasColumnName("bank_details_for_payout");
            entity.Property(e => e.AppealOptionAvailable).HasColumnName("appeal_option_available").HasDefaultValue(false);
            entity.Property(e => e.InAppMessages).HasColumnName("in_app_messages").HasColumnType("JSONB");
            entity.Property(e => e.ProcessingType).HasColumnName("processing_type").HasMaxLength(50).HasDefaultValue("MANUAL").IsRequired();
            entity.Property(e => e.DeductibleAmount).HasColumnName("deductible_amount");
            entity.Property(e => e.CoPayAmount).HasColumnName("co_pay_amount");
            entity.Property(e => e.ProcessorNotes).HasColumnName("processor_notes");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.DeletedAt).HasColumnName("deleted_at");

            // Indexes
            entity.HasIndex(e => e.ClaimNumber).IsUnique().HasDatabaseName("idx_claims_claim_number");
            entity.HasIndex(e => e.PolicyId).HasDatabaseName("idx_claims_policy_id");
            entity.HasIndex(e => e.CustomerId).HasDatabaseName("idx_claims_customer_id");
            entity.HasIndex(e => e.Status).HasDatabaseName("idx_claims_status");
            entity.HasIndex(e => e.Type).HasDatabaseName("idx_claims_type");
            entity.HasIndex(e => e.IncidentDate).HasDatabaseName("idx_claims_incident_date");

            // Relationships
            entity.HasOne(e => e.Policy).WithMany(p => p.Claims).HasForeignKey(e => e.PolicyId).OnDelete(DeleteBehavior.Restrict);

            // Soft delete filter
            entity.HasQueryFilter(e => e.DeletedAt == null);
        });

        // ===== Claim Document Configuration =====
        modelBuilder.Entity<ClaimDocumentEntity>(entity =>
        {
            entity.ToTable("claim_documents");
            entity.HasKey(e => e.DocumentId);
            entity.Property(e => e.DocumentId).HasColumnName("document_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.ClaimId).HasColumnName("claim_id").IsRequired();
            entity.Property(e => e.DocumentType).HasColumnName("document_type").HasMaxLength(100).IsRequired();
            entity.Property(e => e.FileUrl).HasColumnName("file_url").IsRequired();
            entity.Property(e => e.FileHash).HasColumnName("file_hash").HasMaxLength(64).IsRequired();
            entity.Property(e => e.UploadedAt).HasColumnName("uploaded_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.Verified).HasColumnName("verified").HasDefaultValue(false).IsRequired();
            entity.Property(e => e.VerifiedBy).HasColumnName("verified_by");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();

            entity.HasIndex(e => e.ClaimId).HasDatabaseName("idx_claim_documents_claim_id");
            entity.HasIndex(e => e.FileHash).HasDatabaseName("idx_claim_documents_file_hash");
            entity.HasOne(e => e.Claim).WithMany(c => c.Documents).HasForeignKey(e => e.ClaimId).OnDelete(DeleteBehavior.Cascade);
        });

        // ===== Claim Approval Configuration =====
        modelBuilder.Entity<ClaimApprovalEntity>(entity =>
        {
            entity.ToTable("claim_approvals");
            entity.HasKey(e => e.ApprovalId);
            entity.Property(e => e.ApprovalId).HasColumnName("approval_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.ClaimId).HasColumnName("claim_id").IsRequired();
            entity.Property(e => e.ApproverId).HasColumnName("approver_id").IsRequired();
            entity.Property(e => e.ApproverRole).HasColumnName("approver_role").HasMaxLength(100).IsRequired();
            entity.Property(e => e.ApprovalLevel).HasColumnName("approval_level").IsRequired();
            entity.Property(e => e.Decision).HasColumnName("decision").HasMaxLength(50).HasDefaultValue("PENDING").IsRequired();
            entity.Property(e => e.ApprovedAmount).HasColumnName("approved_amount");
            entity.Property(e => e.ApprovedCurrency).HasColumnName("approved_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
            entity.Property(e => e.Notes).HasColumnName("notes");
            entity.Property(e => e.ApprovedAt).HasColumnName("approved_at");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP");

            entity.HasIndex(e => e.ClaimId).HasDatabaseName("idx_claim_approvals_claim_id");
            entity.HasIndex(e => e.ApproverId).HasDatabaseName("idx_claim_approvals_approver_id");
            entity.HasOne(e => e.Claim).WithMany(c => c.Approvals).HasForeignKey(e => e.ClaimId).OnDelete(DeleteBehavior.Cascade);
        });

        // ===== Fraud Check Configuration =====
        modelBuilder.Entity<FraudCheckEntity>(entity =>
        {
            entity.ToTable("fraud_checks");
            entity.HasKey(e => e.FraudCheckId);
            entity.Property(e => e.FraudCheckId).HasColumnName("fraud_check_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.ClaimId).HasColumnName("claim_id").IsRequired();
            entity.Property(e => e.FraudScore).HasColumnName("fraud_score").HasColumnType("DECIMAL(5,2)").IsRequired();
            entity.Property(e => e.RiskFactors).HasColumnName("risk_factors").HasColumnType("TEXT[]");
            entity.Property(e => e.Flagged).HasColumnName("flagged").HasDefaultValue(false).IsRequired();
            entity.Property(e => e.ReviewedBy).HasColumnName("reviewed_by");
            entity.Property(e => e.ReviewedAt).HasColumnName("reviewed_at");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();

            entity.HasIndex(e => e.ClaimId).IsUnique().HasDatabaseName("idx_fraud_checks_claim_id");
            entity.HasIndex(e => e.Flagged).HasDatabaseName("idx_fraud_checks_flagged");
            entity.HasOne(e => e.Claim).WithOne(c => c.FraudCheck).HasForeignKey<FraudCheckEntity>(e => e.ClaimId).OnDelete(DeleteBehavior.Cascade);
        });

        // ===== Product Rider Configuration =====
        modelBuilder.Entity<ProductRiderEntity>(entity =>
        {
            entity.ToTable("product_riders");
            entity.HasKey(e => e.RiderId);
            entity.Property(e => e.RiderId).HasColumnName("rider_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.ProductId).HasColumnName("product_id").IsRequired();
            entity.Property(e => e.NameEn).HasColumnName("name_en").HasMaxLength(255).IsRequired();
            entity.Property(e => e.NameBn).HasColumnName("name_bn").HasMaxLength(255).IsRequired();
            entity.Property(e => e.AdditionalPremium).HasColumnName("additional_premium").IsRequired();
            entity.Property(e => e.AdditionalPremiumCurrency).HasColumnName("additional_premium_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
            entity.Property(e => e.AdditionalCoverage).HasColumnName("additional_coverage").IsRequired();
            entity.Property(e => e.AdditionalCoverageCurrency).HasColumnName("additional_coverage_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();

            entity.HasIndex(e => e.ProductId).HasDatabaseName("idx_product_riders_product_id");
            entity.HasOne(e => e.Product).WithMany(p => p.Riders).HasForeignKey(e => e.ProductId).OnDelete(DeleteBehavior.Cascade);
        });

        // ===== Quote Configuration =====
        modelBuilder.Entity<QuoteEntity>(entity =>
        {
            entity.ToTable("quotes");
            entity.HasKey(e => e.QuoteId);
            entity.Property(e => e.QuoteId).HasColumnName("quote_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.QuoteNumber).HasColumnName("quote_number").HasMaxLength(50).IsRequired();
            entity.Property(e => e.BeneficiaryId).HasColumnName("beneficiary_id").IsRequired();
            entity.Property(e => e.InsurerProductId).HasColumnName("insurer_product_id").IsRequired();
            entity.Property(e => e.Status).HasColumnName("status").HasMaxLength(50).HasDefaultValue("DRAFT").IsRequired();
            entity.Property(e => e.SumAssured).HasColumnName("sum_assured").IsRequired();
            entity.Property(e => e.SumAssuredCurrency).HasColumnName("sum_assured_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
            entity.Property(e => e.TermYears).HasColumnName("term_years").IsRequired();
            entity.Property(e => e.PremiumPaymentMode).HasColumnName("premium_payment_mode").HasMaxLength(50).IsRequired();
            entity.Property(e => e.BasePremium).HasColumnName("base_premium").IsRequired();
            entity.Property(e => e.BasePremiumCurrency).HasColumnName("base_premium_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
            entity.Property(e => e.RiderPremium).HasColumnName("rider_premium");
            entity.Property(e => e.TaxAmount).HasColumnName("tax_amount");
            entity.Property(e => e.TotalPremium).HasColumnName("total_premium").IsRequired();
            entity.Property(e => e.TotalPremiumCurrency).HasColumnName("total_premium_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
            entity.Property(e => e.PremiumCalculation).HasColumnName("premium_calculation").HasColumnType("JSONB");
            entity.Property(e => e.SelectedRiders).HasColumnName("selected_riders").HasColumnType("JSONB");
            entity.Property(e => e.ApplicantAge).HasColumnName("applicant_age").IsRequired();
            entity.Property(e => e.ApplicantOccupation).HasColumnName("applicant_occupation");
            entity.Property(e => e.Smoker).HasColumnName("smoker").HasDefaultValue(false).IsRequired();
            entity.Property(e => e.ValidUntil).HasColumnName("valid_until").IsRequired();
            entity.Property(e => e.ConvertedPolicyId).HasColumnName("converted_policy_id");
            entity.Property(e => e.ConvertedAt).HasColumnName("converted_at");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.DeletedAt).HasColumnName("deleted_at");

            entity.HasIndex(e => e.QuoteNumber).IsUnique().HasDatabaseName("idx_quotes_quote_number");
            entity.HasIndex(e => e.BeneficiaryId).HasDatabaseName("idx_quotes_beneficiary_id");
            entity.HasIndex(e => e.InsurerProductId).HasDatabaseName("idx_quotes_insurer_product_id");
            entity.HasQueryFilter(e => e.DeletedAt == null);
        });

        // ===== Health Declaration Configuration =====
        modelBuilder.Entity<HealthDeclarationEntity>(entity =>
        {
            entity.ToTable("health_declarations");
            entity.HasKey(e => e.DeclarationId);
            entity.Property(e => e.DeclarationId).HasColumnName("declaration_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.QuoteId).HasColumnName("quote_id").IsRequired();
            entity.Property(e => e.HeightCm).HasColumnName("height_cm").IsRequired();
            entity.Property(e => e.WeightKg).HasColumnName("weight_kg").HasMaxLength(20).IsRequired();
            entity.Property(e => e.HasPreExistingConditions).HasColumnName("has_pre_existing_conditions").IsRequired();
            entity.Property(e => e.PreExistingConditions).HasColumnName("pre_existing_conditions").HasColumnType("JSONB");
            entity.Property(e => e.Smoker).HasColumnName("smoker").IsRequired();
            entity.Property(e => e.AlcoholConsumer).HasColumnName("alcohol_consumer").IsRequired();
            entity.Property(e => e.OccupationRiskLevel).HasColumnName("occupation_risk_level").HasMaxLength(50);
            entity.Property(e => e.MedicalExamRequired).HasColumnName("medical_exam_required").HasDefaultValue(false).IsRequired();
            entity.Property(e => e.AutoApprovalPossible).HasColumnName("auto_approval_possible").HasDefaultValue(false).IsRequired();
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();

            entity.HasIndex(e => e.QuoteId).IsUnique().HasDatabaseName("idx_health_declarations_quote_id");
        });

        // ===== Underwriting Decision Configuration =====
        modelBuilder.Entity<UnderwritingDecisionEntity>(entity =>
        {
            entity.ToTable("underwriting_decisions");
            entity.HasKey(e => e.DecisionId);
            entity.Property(e => e.DecisionId).HasColumnName("decision_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.QuoteId).HasColumnName("quote_id").IsRequired();
            entity.Property(e => e.UnderwriterId).HasColumnName("underwriter_id").IsRequired();
            entity.Property(e => e.Decision).HasColumnName("decision").HasMaxLength(50).HasDefaultValue("PENDING").IsRequired();
            entity.Property(e => e.RiskLevel).HasColumnName("risk_level").HasMaxLength(50);
            entity.Property(e => e.PremiumAdjusted).HasColumnName("premium_adjusted").HasDefaultValue(false).IsRequired();
            entity.Property(e => e.AdjustedPremium).HasColumnName("adjusted_premium");
            entity.Property(e => e.AdjustedPremiumCurrency).HasColumnName("adjusted_premium_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
            entity.Property(e => e.Conditions).HasColumnName("conditions").HasColumnType("JSONB");
            entity.Property(e => e.Comments).HasColumnName("comments");
            entity.Property(e => e.RejectionReason).HasColumnName("rejection_reason");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();

            entity.HasIndex(e => e.QuoteId).HasDatabaseName("idx_underwriting_decisions_quote_id");
        });

        // ===== Commission Configuration =====
        modelBuilder.Entity<CommissionEntity>(entity =>
        {
            entity.ToTable("commissions");
            entity.HasKey(e => e.CommissionId);
            entity.Property(e => e.CommissionId).HasColumnName("commission_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.CommissionNumber).HasColumnName("commission_number").HasMaxLength(50).IsRequired();
            entity.Property(e => e.PolicyId).HasColumnName("policy_id").IsRequired();
            entity.Property(e => e.CommissionType).HasColumnName("commission_type").HasMaxLength(50).HasDefaultValue("ACQUISITION").IsRequired();
            entity.Property(e => e.RecipientType).HasColumnName("recipient_type").HasMaxLength(50).IsRequired();
            entity.Property(e => e.RecipientId).HasColumnName("recipient_id").IsRequired();
            entity.Property(e => e.PremiumAmount).HasColumnName("premium_amount").IsRequired();
            entity.Property(e => e.PremiumCurrency).HasColumnName("premium_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
            entity.Property(e => e.CommissionRate).HasColumnName("commission_rate").IsRequired();
            entity.Property(e => e.CommissionAmount).HasColumnName("commission_amount").IsRequired();
            entity.Property(e => e.CommissionCurrency).HasColumnName("commission_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
            entity.Property(e => e.CalculationBreakdown).HasColumnName("calculation_breakdown").HasColumnType("JSONB");
            entity.Property(e => e.Status).HasColumnName("status").HasMaxLength(50).HasDefaultValue("PENDING").IsRequired();
            entity.Property(e => e.PayoutId).HasColumnName("payout_id");
            entity.Property(e => e.PaidAt).HasColumnName("paid_at");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.DeletedAt).HasColumnName("deleted_at");

            entity.HasIndex(e => e.CommissionNumber).IsUnique().HasDatabaseName("idx_commissions_number");
            entity.HasIndex(e => e.PolicyId).HasDatabaseName("idx_commissions_policy_id");
            entity.HasIndex(e => e.RecipientId).HasDatabaseName("idx_commissions_recipient_id");
            entity.HasQueryFilter(e => e.DeletedAt == null);
        });

        // ===== Commission Payout Configuration =====
        modelBuilder.Entity<CommissionPayoutEntity>(entity =>
        {
            entity.ToTable("commission_payouts");
            entity.HasKey(e => e.PayoutId);
            entity.Property(e => e.PayoutId).HasColumnName("payout_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.PayoutNumber).HasColumnName("payout_number").HasMaxLength(50).IsRequired();
            entity.Property(e => e.RecipientType).HasColumnName("recipient_type").HasMaxLength(50).IsRequired();
            entity.Property(e => e.RecipientId).HasColumnName("recipient_id").IsRequired();
            entity.Property(e => e.TotalAmount).HasColumnName("total_amount").IsRequired();
            entity.Property(e => e.TotalCurrency).HasColumnName("total_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
            entity.Property(e => e.CommissionCount).HasColumnName("commission_count").IsRequired();
            entity.Property(e => e.PeriodStart).HasColumnName("period_start").HasMaxLength(20).IsRequired();
            entity.Property(e => e.PeriodEnd).HasColumnName("period_end").HasMaxLength(20).IsRequired();
            entity.Property(e => e.Status).HasColumnName("status").HasMaxLength(50).HasDefaultValue("PENDING").IsRequired();
            entity.Property(e => e.PaymentMethod).HasColumnName("payment_method").HasMaxLength(50);
            entity.Property(e => e.PaymentReference).HasColumnName("payment_reference").HasMaxLength(255);
            entity.Property(e => e.PaidAt).HasColumnName("paid_at");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();

            entity.HasIndex(e => e.PayoutNumber).IsUnique().HasDatabaseName("idx_commission_payouts_number");
            entity.HasIndex(e => e.RecipientId).HasDatabaseName("idx_commission_payouts_recipient_id");
        });

        // ===== Beneficiary Configuration =====
        modelBuilder.Entity<BeneficiaryEntity>(entity =>
        {
            entity.ToTable("beneficiaries");
            entity.HasKey(e => e.BeneficiaryId);
            entity.Property(e => e.BeneficiaryId).HasColumnName("beneficiary_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.UserId).HasColumnName("user_id").IsRequired();
            entity.Property(e => e.Type).HasColumnName("type").HasMaxLength(20).IsRequired();
            entity.Property(e => e.Code).HasColumnName("code").HasMaxLength(20).IsRequired();
            entity.Property(e => e.Status).HasColumnName("status").HasMaxLength(20).HasDefaultValue("PENDING_KYC").IsRequired();
            entity.Property(e => e.KycStatus).HasColumnName("kyc_status").HasMaxLength(20).HasDefaultValue("NOT_STARTED").IsRequired();
            entity.Property(e => e.KycCompletedAt).HasColumnName("kyc_completed_at");
            entity.Property(e => e.RiskScore).HasColumnName("risk_score").HasMaxLength(10);
            entity.Property(e => e.ReferralCode).HasColumnName("referral_code").HasMaxLength(20);
            entity.Property(e => e.ReferredBy).HasColumnName("referred_by");
            entity.Property(e => e.PartnerId).HasColumnName("partner_id");
            entity.Property(e => e.AuditInfo).HasColumnName("audit_info").HasColumnType("JSONB");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.DeletedAt).HasColumnName("deleted_at");

            entity.HasIndex(e => e.Code).IsUnique().HasDatabaseName("idx_beneficiaries_code_unique");
            entity.HasQueryFilter(e => e.DeletedAt == null);

            entity.HasOne(e => e.IndividualDetails).WithOne(i => i.Beneficiary).HasForeignKey<IndividualBeneficiaryEntity>(i => i.BeneficiaryId).OnDelete(DeleteBehavior.Cascade);
            entity.HasOne(e => e.BusinessDetails).WithOne(b => b.Beneficiary).HasForeignKey<BusinessBeneficiaryEntity>(b => b.ParentBeneficiaryId).OnDelete(DeleteBehavior.Cascade);
        });

        modelBuilder.Entity<IndividualBeneficiaryEntity>(entity =>
        {
            entity.ToTable("individual_beneficiaries");
            entity.HasKey(e => e.BeneficiaryId);
            entity.Property(e => e.BeneficiaryId).HasColumnName("beneficiary_id");
            entity.Property(e => e.FullName).HasColumnName("full_name").HasMaxLength(255).IsRequired();
            entity.Property(e => e.FullNameBn).HasColumnName("full_name_bn").HasMaxLength(255);
            entity.Property(e => e.DateOfBirth).HasColumnName("date_of_birth").HasColumnType("DATE").IsRequired();
            entity.Property(e => e.Gender).HasColumnName("gender").HasMaxLength(10).IsRequired();
            entity.Property(e => e.NidNumber).HasColumnName("nid_number").HasMaxLength(17);
            entity.Property(e => e.PassportNumber).HasColumnName("passport_number").HasMaxLength(20);
            entity.Property(e => e.BirthCertificateNumber).HasColumnName("birth_certificate_number").HasMaxLength(20);
            entity.Property(e => e.TinNumber).HasColumnName("tin_number").HasMaxLength(12);
            entity.Property(e => e.MaritalStatus).HasColumnName("marital_status").HasMaxLength(20);
            entity.Property(e => e.Occupation).HasColumnName("occupation").HasMaxLength(100);
            entity.Property(e => e.ContactInfo).HasColumnName("contact_info").HasColumnType("JSONB");
            entity.Property(e => e.PermanentAddress).HasColumnName("permanent_address").HasColumnType("JSONB");
            entity.Property(e => e.PresentAddress).HasColumnName("present_address").HasColumnType("JSONB");
            entity.Property(e => e.NomineeName).HasColumnName("nominee_name").HasMaxLength(255);
            entity.Property(e => e.NomineeRelationship).HasColumnName("nominee_relationship").HasMaxLength(50);
            entity.Property(e => e.AuditInfo).HasColumnName("audit_info").HasColumnType("JSONB");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();

            entity.HasIndex(e => e.NidNumber).IsUnique().HasDatabaseName("idx_individual_beneficiaries_nid");
        });

        modelBuilder.Entity<BusinessBeneficiaryEntity>(entity =>
        {
            entity.ToTable("business_beneficiaries");
            entity.HasKey(e => e.BeneficiaryId);
            entity.Property(e => e.BeneficiaryId).HasColumnName("beneficiary_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.ParentBeneficiaryId).HasColumnName("parent_beneficiary_id").IsRequired();
            entity.Property(e => e.BusinessName).HasColumnName("business_name").HasMaxLength(255).IsRequired();
            entity.Property(e => e.BusinessNameBn).HasColumnName("business_name_bn").HasMaxLength(255);
            entity.Property(e => e.TradeLicenseNumber).HasColumnName("trade_license_number").HasMaxLength(50).IsRequired();
            entity.Property(e => e.TradeLicenseIssueDate).HasColumnName("trade_license_issue_date").HasColumnType("DATE");
            entity.Property(e => e.TradeLicenseExpiryDate).HasColumnName("trade_license_expiry_date").HasColumnType("DATE");
            entity.Property(e => e.TinNumber).HasColumnName("tin_number").HasMaxLength(12).IsRequired();
            entity.Property(e => e.BinNumber).HasColumnName("bin_number").HasMaxLength(15);
            entity.Property(e => e.BusinessType).HasColumnName("business_type").HasMaxLength(50).IsRequired();
            entity.Property(e => e.IndustrySector).HasColumnName("industry_sector").HasMaxLength(100);
            entity.Property(e => e.EmployeeCount).HasColumnName("employee_count");
            entity.Property(e => e.IncorporationDate).HasColumnName("incorporation_date").HasColumnType("DATE");
            entity.Property(e => e.ContactInfo).HasColumnName("contact_info").HasColumnType("JSONB");
            entity.Property(e => e.RegisteredAddress).HasColumnName("registered_address").HasColumnType("JSONB");
            entity.Property(e => e.BusinessAddress).HasColumnName("business_address").HasColumnType("JSONB");
            entity.Property(e => e.FocalPersonName).HasColumnName("focal_person_name").HasMaxLength(255).IsRequired();
            entity.Property(e => e.FocalPersonDesignation).HasColumnName("focal_person_designation").HasMaxLength(100);
            entity.Property(e => e.FocalPersonNid).HasColumnName("focal_person_nid").HasMaxLength(17);
            entity.Property(e => e.FocalPersonContact).HasColumnName("focal_person_contact").HasColumnType("JSONB");
            entity.Property(e => e.RegistrationNumber).HasColumnName("registration_number").HasMaxLength(100);
            entity.Property(e => e.TaxId).HasColumnName("tax_id").HasMaxLength(50);
            entity.Property(e => e.PrimaryContact).HasColumnName("primary_contact").HasColumnType("JSONB");
            entity.Property(e => e.TotalEmployeesCovered).HasColumnName("total_employees_covered").HasDefaultValue(0);
            entity.Property(e => e.ActivePoliciesCount).HasColumnName("active_policies_count").HasDefaultValue(0);
            entity.Property(e => e.TotalPremiumAmount).HasColumnName("total_premium_amount").HasDefaultValue(0);
            entity.Property(e => e.PendingActionsCount).HasColumnName("pending_actions_count").HasDefaultValue(0);
            entity.Property(e => e.AuditInfo).HasColumnName("audit_info").HasColumnType("JSONB");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();

            entity.HasIndex(e => e.ParentBeneficiaryId).IsUnique().HasDatabaseName("idx_business_beneficiaries_parent_beneficiary_id");
            entity.HasIndex(e => e.TradeLicenseNumber).IsUnique().HasDatabaseName("idx_business_beneficiaries_trade_license");
            entity.HasIndex(e => e.TinNumber).IsUnique().HasDatabaseName("idx_business_beneficiaries_tin");
        });
    }
}
