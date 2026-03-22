using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace InsuranceEngine.Policy.Infrastructure.Persistence.Migrations
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
                name: "endorsement_number_seq",
                schema: "insurance_schema");

            migrationBuilder.CreateSequence(
                name: "policy_number_seq",
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
                name: "endorsements",
                schema: "insurance_schema",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uuid", nullable: false),
                    EndorsementNumber = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: false),
                    PolicyId = table.Column<Guid>(type: "uuid", nullable: false),
                    Type = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: false),
                    Reason = table.Column<string>(type: "character varying(500)", maxLength: 500, nullable: false),
                    changes = table.Column<string>(type: "jsonb", nullable: false),
                    premium_adjustment_amount = table.Column<decimal>(type: "numeric", nullable: false),
                    premium_adjustment_currency = table.Column<string>(type: "character varying(3)", maxLength: 3, nullable: false, defaultValue: "BDT"),
                    PremiumRefundRequired = table.Column<bool>(type: "boolean", nullable: false),
                    Status = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: false),
                    RequestedBy = table.Column<Guid>(type: "uuid", nullable: false),
                    ApprovedBy = table.Column<Guid>(type: "uuid", nullable: true),
                    EffectiveDate = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    At = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    audit_info = table.Column<string>(type: "jsonb", nullable: false),
                    CreatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    UpdatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_endorsements", x => x.Id);
                });

            migrationBuilder.CreateTable(
                name: "policies",
                schema: "insurance_schema",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uuid", nullable: false),
                    PolicyNumber = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: false),
                    ProductId = table.Column<Guid>(type: "uuid", nullable: false),
                    CustomerId = table.Column<Guid>(type: "uuid", nullable: false),
                    PartnerId = table.Column<Guid>(type: "uuid", nullable: true),
                    AgentId = table.Column<Guid>(type: "uuid", nullable: true),
                    QuoteId = table.Column<Guid>(type: "uuid", nullable: true),
                    UnderwritingDecisionId = table.Column<Guid>(type: "uuid", nullable: true),
                    Status = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: false),
                    premium_amount = table.Column<long>(type: "bigint", nullable: false),
                    premium_currency = table.Column<string>(type: "character varying(3)", maxLength: 3, nullable: false, defaultValue: "BDT"),
                    sum_insured_amount = table.Column<long>(type: "bigint", nullable: false),
                    sum_insured_currency = table.Column<string>(type: "character varying(3)", maxLength: 3, nullable: false, defaultValue: "BDT"),
                    vat_tax_amount = table.Column<long>(type: "bigint", nullable: false),
                    service_fee_amount = table.Column<long>(type: "bigint", nullable: false),
                    total_payable_amount = table.Column<long>(type: "bigint", nullable: false),
                    TenureMonths = table.Column<int>(type: "integer", nullable: false),
                    StartDate = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    EndDate = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    IssuedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    PaymentFrequency = table.Column<string>(type: "text", nullable: true),
                    PaymentGatewayReference = table.Column<string>(type: "text", nullable: true),
                    ReceiptNumber = table.Column<string>(type: "text", nullable: true),
                    PolicyDocumentUrl = table.Column<string>(type: "text", nullable: true),
                    proposer_details = table.Column<string>(type: "jsonb", nullable: true),
                    OccupationRiskClass = table.Column<string>(type: "text", nullable: true),
                    HasExistingPolicies = table.Column<bool>(type: "boolean", nullable: false),
                    ClaimsHistorySummary = table.Column<string>(type: "text", nullable: true),
                    ProviderName = table.Column<string>(type: "text", nullable: true),
                    EnrollmentStartDate = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    EnrollmentEndDate = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    underwriting_data = table.Column<string>(type: "jsonb", nullable: true),
                    CreatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    UpdatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    DeletedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    IsDeleted = table.Column<bool>(type: "boolean", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_policies", x => x.Id);
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

            migrationBuilder.CreateTable(
                name: "policy_nominees",
                schema: "insurance_schema",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uuid", nullable: false),
                    PolicyId = table.Column<Guid>(type: "uuid", nullable: false),
                    BeneficiaryId = table.Column<Guid>(type: "uuid", nullable: true),
                    full_name = table.Column<string>(type: "character varying(200)", maxLength: 200, nullable: false),
                    Relationship = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: false),
                    SharePercentage = table.Column<double>(type: "double precision", nullable: false),
                    date_of_birth = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    nominee_dob_text = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: true),
                    nid_number = table.Column<string>(type: "character varying(20)", maxLength: 20, nullable: true),
                    phone_number = table.Column<string>(type: "character varying(20)", maxLength: 20, nullable: true),
                    CreatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    UpdatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    DeletedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    IsDeleted = table.Column<bool>(type: "boolean", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_policy_nominees", x => x.Id);
                    table.ForeignKey(
                        name: "FK_policy_nominees_beneficiaries_BeneficiaryId",
                        column: x => x.BeneficiaryId,
                        principalSchema: "insurance_schema",
                        principalTable: "beneficiaries",
                        principalColumn: "Id");
                    table.ForeignKey(
                        name: "FK_policy_nominees_policies_PolicyId",
                        column: x => x.PolicyId,
                        principalSchema: "insurance_schema",
                        principalTable: "policies",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateTable(
                name: "policy_riders",
                schema: "insurance_schema",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uuid", nullable: false),
                    PolicyId = table.Column<Guid>(type: "uuid", nullable: false),
                    RiderName = table.Column<string>(type: "character varying(255)", maxLength: 255, nullable: false),
                    premium_amount = table.Column<long>(type: "bigint", nullable: false),
                    premium_currency = table.Column<string>(type: "character varying(3)", maxLength: 3, nullable: false, defaultValue: "BDT"),
                    coverage_amount = table.Column<long>(type: "bigint", nullable: false),
                    coverage_currency = table.Column<string>(type: "character varying(3)", maxLength: 3, nullable: false, defaultValue: "BDT"),
                    CreatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    UpdatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_policy_riders", x => x.Id);
                    table.ForeignKey(
                        name: "FK_policy_riders_policies_PolicyId",
                        column: x => x.PolicyId,
                        principalSchema: "insurance_schema",
                        principalTable: "policies",
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
                name: "IX_endorsements_EndorsementNumber",
                schema: "insurance_schema",
                table: "endorsements",
                column: "EndorsementNumber",
                unique: true);

            migrationBuilder.CreateIndex(
                name: "IX_endorsements_PolicyId",
                schema: "insurance_schema",
                table: "endorsements",
                column: "PolicyId");

            migrationBuilder.CreateIndex(
                name: "IX_endorsements_Status",
                schema: "insurance_schema",
                table: "endorsements",
                column: "Status");

            migrationBuilder.CreateIndex(
                name: "IX_policies_CustomerId",
                schema: "insurance_schema",
                table: "policies",
                column: "CustomerId");

            migrationBuilder.CreateIndex(
                name: "IX_policies_PolicyNumber",
                schema: "insurance_schema",
                table: "policies",
                column: "PolicyNumber",
                unique: true);

            migrationBuilder.CreateIndex(
                name: "IX_policies_ProductId",
                schema: "insurance_schema",
                table: "policies",
                column: "ProductId");

            migrationBuilder.CreateIndex(
                name: "IX_policies_Status",
                schema: "insurance_schema",
                table: "policies",
                column: "Status");

            migrationBuilder.CreateIndex(
                name: "IX_policy_nominees_BeneficiaryId",
                schema: "insurance_schema",
                table: "policy_nominees",
                column: "BeneficiaryId");

            migrationBuilder.CreateIndex(
                name: "IX_policy_nominees_nid_number",
                schema: "insurance_schema",
                table: "policy_nominees",
                column: "nid_number");

            migrationBuilder.CreateIndex(
                name: "IX_policy_nominees_PolicyId",
                schema: "insurance_schema",
                table: "policy_nominees",
                column: "PolicyId");

            migrationBuilder.CreateIndex(
                name: "IX_policy_riders_PolicyId",
                schema: "insurance_schema",
                table: "policy_riders",
                column: "PolicyId");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "business_beneficiaries",
                schema: "insurance_schema");

            migrationBuilder.DropTable(
                name: "endorsements",
                schema: "insurance_schema");

            migrationBuilder.DropTable(
                name: "individual_beneficiaries",
                schema: "insurance_schema");

            migrationBuilder.DropTable(
                name: "policy_nominees",
                schema: "insurance_schema");

            migrationBuilder.DropTable(
                name: "policy_riders",
                schema: "insurance_schema");

            migrationBuilder.DropTable(
                name: "beneficiaries",
                schema: "insurance_schema");

            migrationBuilder.DropTable(
                name: "policies",
                schema: "insurance_schema");

            migrationBuilder.DropSequence(
                name: "endorsement_number_seq",
                schema: "insurance_schema");

            migrationBuilder.DropSequence(
                name: "policy_number_seq",
                schema: "insurance_schema");
        }
    }
}
