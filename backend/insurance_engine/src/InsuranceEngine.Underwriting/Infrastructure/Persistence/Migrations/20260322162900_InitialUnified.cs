using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace InsuranceEngine.Underwriting.Infrastructure.Persistence.Migrations
{
    /// <inheritdoc />
    public partial class InitialUnified : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.EnsureSchema(
                name: "insurance_schema");

            migrationBuilder.CreateSequence(
                name: "quote_number_seq",
                schema: "insurance_schema");

            migrationBuilder.CreateTable(
                name: "beneficiaries",
                schema: "insurance_schema",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uuid", nullable: false),
                    UserId = table.Column<Guid>(type: "uuid", nullable: false),
                    Type = table.Column<string>(type: "text", nullable: false),
                    Code = table.Column<string>(type: "character varying(20)", maxLength: 20, nullable: false),
                    Status = table.Column<string>(type: "text", nullable: false),
                    KycStatus = table.Column<string>(type: "text", nullable: false),
                    KycCompletedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    RiskScore = table.Column<string>(type: "text", nullable: true),
                    ReferralCode = table.Column<string>(type: "text", nullable: true),
                    ReferredBy = table.Column<Guid>(type: "uuid", nullable: true),
                    PartnerId = table.Column<Guid>(type: "uuid", nullable: true),
                    Name = table.Column<string>(type: "text", nullable: false),
                    ContactNumber = table.Column<string>(type: "text", nullable: false),
                    Email = table.Column<string>(type: "text", nullable: false),
                    Address = table.Column<string>(type: "text", nullable: false),
                    AuditInfo = table.Column<string>(type: "jsonb", nullable: true),
                    CreatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    UpdatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    IsDeleted = table.Column<bool>(type: "boolean", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_beneficiaries", x => x.Id);
                });

            migrationBuilder.CreateTable(
                name: "health_declarations",
                schema: "insurance_schema",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uuid", nullable: false),
                    QuoteId = table.Column<Guid>(type: "uuid", nullable: false),
                    HeightCm = table.Column<int>(type: "integer", nullable: false),
                    WeightKg = table.Column<decimal>(type: "numeric(5,2)", precision: 5, scale: 2, nullable: false),
                    Bmi = table.Column<decimal>(type: "numeric(5,2)", precision: 5, scale: 2, nullable: false),
                    HasPreExistingConditions = table.Column<bool>(type: "boolean", nullable: false),
                    pre_existing_conditions = table.Column<string>(type: "jsonb", nullable: true),
                    IsCurrentlyHospitalized = table.Column<bool>(type: "boolean", nullable: false),
                    HasFamilyHistory = table.Column<bool>(type: "boolean", nullable: false),
                    family_history = table.Column<string>(type: "jsonb", nullable: true),
                    IsSmoker = table.Column<bool>(type: "boolean", nullable: false),
                    IsAlcoholConsumer = table.Column<bool>(type: "boolean", nullable: false),
                    OccupationRiskLevel = table.Column<string>(type: "text", nullable: true),
                    IsMedicalExamRequired = table.Column<bool>(type: "boolean", nullable: false),
                    IsMedicalExamCompleted = table.Column<bool>(type: "boolean", nullable: false),
                    medical_exam_results = table.Column<string>(type: "jsonb", nullable: true),
                    MedicalExamDate = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    medical_documents = table.Column<string>(type: "jsonb", nullable: true),
                    audit_info = table.Column<string>(type: "jsonb", nullable: true),
                    CreatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    UpdatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_health_declarations", x => x.Id);
                });

            migrationBuilder.CreateTable(
                name: "quotes",
                schema: "insurance_schema",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uuid", nullable: false),
                    QuoteNumber = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: false),
                    BeneficiaryId = table.Column<Guid>(type: "uuid", nullable: false),
                    InsurerProductId = table.Column<Guid>(type: "uuid", nullable: false),
                    Status = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: false),
                    sum_assured_amount = table.Column<long>(type: "bigint", nullable: false),
                    sum_assured_currency = table.Column<string>(type: "character varying(3)", maxLength: 3, nullable: false, defaultValue: "BDT"),
                    TermYears = table.Column<int>(type: "integer", nullable: false),
                    PremiumPaymentMode = table.Column<string>(type: "text", nullable: false),
                    base_premium_amount = table.Column<long>(type: "bigint", nullable: false),
                    rider_premium_amount = table.Column<long>(type: "bigint", nullable: false),
                    tax_amount = table.Column<long>(type: "bigint", nullable: false),
                    total_premium_amount = table.Column<long>(type: "bigint", nullable: false),
                    currency = table.Column<string>(type: "character varying(3)", maxLength: 3, nullable: false, defaultValue: "BDT"),
                    premium_calculation = table.Column<string>(type: "jsonb", nullable: true),
                    selected_riders = table.Column<string>(type: "jsonb", nullable: true),
                    ApplicantAge = table.Column<int>(type: "integer", nullable: false),
                    ApplicantOccupation = table.Column<string>(type: "text", nullable: true),
                    IsSmoker = table.Column<bool>(type: "boolean", nullable: false),
                    ValidUntil = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    ConvertedPolicyId = table.Column<Guid>(type: "uuid", nullable: true),
                    ConvertedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    audit_info = table.Column<string>(type: "jsonb", nullable: true),
                    CreatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    UpdatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    DeletedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    IsDeleted = table.Column<bool>(type: "boolean", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_quotes", x => x.Id);
                });

            migrationBuilder.CreateTable(
                name: "underwriting_decisions",
                schema: "insurance_schema",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uuid", nullable: false),
                    QuoteId = table.Column<Guid>(type: "uuid", nullable: false),
                    Decision = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: false),
                    Method = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: false),
                    RiskScore = table.Column<decimal>(type: "numeric(5,2)", precision: 5, scale: 2, nullable: false),
                    RiskLevel = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: false),
                    risk_factors = table.Column<string>(type: "jsonb", nullable: true),
                    Reason = table.Column<string>(type: "text", nullable: true),
                    conditions = table.Column<string>(type: "jsonb", nullable: true),
                    IsPremiumAdjusted = table.Column<bool>(type: "boolean", nullable: false),
                    adjusted_premium_amount = table.Column<long>(type: "bigint", nullable: false),
                    adjusted_premium_currency = table.Column<string>(type: "character varying(3)", maxLength: 3, nullable: false, defaultValue: "BDT"),
                    AdjustmentReason = table.Column<string>(type: "text", nullable: true),
                    UnderwriterId = table.Column<Guid>(type: "uuid", nullable: true),
                    UnderwriterComments = table.Column<string>(type: "text", nullable: true),
                    DecidedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    ValidUntil = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    audit_info = table.Column<string>(type: "jsonb", nullable: true),
                    CreatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    UpdatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_underwriting_decisions", x => x.Id);
                });

            migrationBuilder.CreateTable(
                name: "business_beneficiaries",
                schema: "insurance_schema",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uuid", nullable: false),
                    BeneficiaryId = table.Column<Guid>(type: "uuid", nullable: false),
                    BusinessName = table.Column<string>(type: "text", nullable: false),
                    BusinessNameBn = table.Column<string>(type: "text", nullable: true),
                    TradeLicenseNumber = table.Column<string>(type: "text", nullable: false),
                    TradeLicenseIssueDate = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    TradeLicenseExpiryDate = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    TinNumber = table.Column<string>(type: "text", nullable: false),
                    BinNumber = table.Column<string>(type: "text", nullable: true),
                    BusinessType = table.Column<string>(type: "text", nullable: false),
                    IndustrySector = table.Column<string>(type: "text", nullable: true),
                    EmployeeCount = table.Column<int>(type: "integer", nullable: false),
                    IncorporationDate = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    ContactInfoJson = table.Column<string>(type: "jsonb", nullable: true),
                    RegisteredAddressJson = table.Column<string>(type: "jsonb", nullable: true),
                    BusinessAddressJson = table.Column<string>(type: "jsonb", nullable: true),
                    FocalPersonName = table.Column<string>(type: "text", nullable: false),
                    FocalPersonDesignation = table.Column<string>(type: "text", nullable: true),
                    FocalPersonNid = table.Column<string>(type: "text", nullable: true),
                    FocalPersonContactJson = table.Column<string>(type: "jsonb", nullable: true),
                    Industry = table.Column<string>(type: "text", nullable: true),
                    FocalPersonContact = table.Column<string>(type: "text", nullable: true),
                    AuditInfo = table.Column<string>(type: "jsonb", nullable: true),
                    RegistrationNumber = table.Column<string>(type: "text", nullable: true),
                    TaxId = table.Column<string>(type: "text", nullable: true),
                    PrimaryContactJson = table.Column<string>(type: "jsonb", nullable: true),
                    TotalEmployeesCovered = table.Column<int>(type: "integer", nullable: false),
                    ActivePoliciesCount = table.Column<int>(type: "integer", nullable: false),
                    TotalPremiumAmount = table.Column<long>(type: "bigint", nullable: false),
                    PendingActionsCount = table.Column<int>(type: "integer", nullable: false),
                    CreatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    UpdatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_business_beneficiaries", x => x.Id);
                    table.ForeignKey(
                        name: "FK_business_beneficiaries_beneficiaries_BeneficiaryId",
                        column: x => x.BeneficiaryId,
                        principalSchema: "insurance_schema",
                        principalTable: "beneficiaries",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateTable(
                name: "individual_beneficiaries",
                schema: "insurance_schema",
                columns: table => new
                {
                    BeneficiaryId = table.Column<Guid>(type: "uuid", nullable: false),
                    FullName = table.Column<string>(type: "text", nullable: false),
                    FullNameBn = table.Column<string>(type: "text", nullable: true),
                    DateOfBirth = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    Gender = table.Column<string>(type: "text", nullable: false),
                    NidNumber = table.Column<string>(type: "text", nullable: true),
                    PassportNumber = table.Column<string>(type: "text", nullable: true),
                    BirthCertificateNumber = table.Column<string>(type: "text", nullable: true),
                    TinNumber = table.Column<string>(type: "text", nullable: true),
                    MaritalStatus = table.Column<string>(type: "text", nullable: false),
                    Occupation = table.Column<string>(type: "text", nullable: true),
                    FatherName = table.Column<string>(type: "text", nullable: true),
                    MotherName = table.Column<string>(type: "text", nullable: true),
                    MonthlyIncome = table.Column<decimal>(type: "numeric", nullable: false),
                    ContactInfoJson = table.Column<string>(type: "jsonb", nullable: true),
                    PermanentAddressJson = table.Column<string>(type: "jsonb", nullable: true),
                    PresentAddressJson = table.Column<string>(type: "jsonb", nullable: true),
                    NomineeName = table.Column<string>(type: "text", nullable: true),
                    NomineeRelationship = table.Column<string>(type: "text", nullable: true),
                    AuditInfo = table.Column<string>(type: "jsonb", nullable: true),
                    CreatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    UpdatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_individual_beneficiaries", x => x.BeneficiaryId);
                    table.ForeignKey(
                        name: "FK_individual_beneficiaries_beneficiaries_BeneficiaryId",
                        column: x => x.BeneficiaryId,
                        principalSchema: "insurance_schema",
                        principalTable: "beneficiaries",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateIndex(
                name: "IX_business_beneficiaries_BeneficiaryId",
                schema: "insurance_schema",
                table: "business_beneficiaries",
                column: "BeneficiaryId",
                unique: true);

            migrationBuilder.CreateIndex(
                name: "IX_health_declarations_QuoteId",
                schema: "insurance_schema",
                table: "health_declarations",
                column: "QuoteId",
                unique: true);

            migrationBuilder.CreateIndex(
                name: "IX_quotes_BeneficiaryId",
                schema: "insurance_schema",
                table: "quotes",
                column: "BeneficiaryId");

            migrationBuilder.CreateIndex(
                name: "IX_quotes_InsurerProductId",
                schema: "insurance_schema",
                table: "quotes",
                column: "InsurerProductId");

            migrationBuilder.CreateIndex(
                name: "IX_quotes_QuoteNumber",
                schema: "insurance_schema",
                table: "quotes",
                column: "QuoteNumber",
                unique: true);

            migrationBuilder.CreateIndex(
                name: "IX_quotes_Status",
                schema: "insurance_schema",
                table: "quotes",
                column: "Status");

            migrationBuilder.CreateIndex(
                name: "IX_underwriting_decisions_Decision",
                schema: "insurance_schema",
                table: "underwriting_decisions",
                column: "Decision");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "business_beneficiaries",
                schema: "insurance_schema");

            migrationBuilder.DropTable(
                name: "health_declarations",
                schema: "insurance_schema");

            migrationBuilder.DropTable(
                name: "individual_beneficiaries",
                schema: "insurance_schema");

            migrationBuilder.DropTable(
                name: "quotes",
                schema: "insurance_schema");

            migrationBuilder.DropTable(
                name: "underwriting_decisions",
                schema: "insurance_schema");

            migrationBuilder.DropTable(
                name: "beneficiaries",
                schema: "insurance_schema");

            migrationBuilder.DropSequence(
                name: "quote_number_seq",
                schema: "insurance_schema");
        }
    }
}
