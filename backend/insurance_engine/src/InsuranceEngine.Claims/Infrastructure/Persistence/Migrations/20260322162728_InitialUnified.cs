using System;
using System.Collections.Generic;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace InsuranceEngine.Claims.Infrastructure.Persistence.Migrations
{
    /// <inheritdoc />
    public partial class InitialUnified : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.EnsureSchema(
                name: "insurance_schema");

            migrationBuilder.CreateTable(
                name: "claims",
                schema: "insurance_schema",
                columns: table => new
                {
                    id = table.Column<Guid>(type: "uuid", nullable: false),
                    claim_number = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: false),
                    policy_id = table.Column<Guid>(type: "uuid", nullable: false),
                    customer_id = table.Column<Guid>(type: "uuid", nullable: false),
                    status = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: false),
                    type = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: false),
                    claimed_amount = table.Column<long>(type: "bigint", nullable: false),
                    claimed_currency = table.Column<string>(type: "character varying(3)", maxLength: 3, nullable: false, defaultValue: "BDT"),
                    approved_amount = table.Column<long>(type: "bigint", nullable: false),
                    approved_currency = table.Column<string>(type: "character varying(3)", maxLength: 3, nullable: false, defaultValue: "BDT"),
                    settled_amount = table.Column<long>(type: "bigint", nullable: false),
                    settled_currency = table.Column<string>(type: "character varying(3)", maxLength: 3, nullable: false, defaultValue: "BDT"),
                    incident_date = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    incident_description = table.Column<string>(type: "text", nullable: false),
                    place_of_incident = table.Column<string>(type: "text", nullable: true),
                    submitted_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    approved_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    settled_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    rejection_reason = table.Column<string>(type: "text", nullable: true),
                    processing_type = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: false),
                    deductible_amount = table.Column<long>(type: "bigint", nullable: false),
                    deductible_currency = table.Column<string>(type: "character varying(3)", maxLength: 3, nullable: false, defaultValue: "BDT"),
                    co_pay_amount = table.Column<long>(type: "bigint", nullable: false),
                    co_pay_currency = table.Column<string>(type: "character varying(3)", maxLength: 3, nullable: false, defaultValue: "BDT"),
                    bank_details_for_payout = table.Column<string>(type: "text", nullable: true),
                    appeal_option_available = table.Column<bool>(type: "boolean", nullable: false, defaultValue: false),
                    in_app_messages = table.Column<string>(type: "jsonb", nullable: true),
                    processor_notes = table.Column<string>(type: "text", nullable: true),
                    fraud_check_id = table.Column<Guid>(type: "uuid", nullable: true),
                    created_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    updated_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    deleted_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    is_deleted = table.Column<bool>(type: "boolean", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("pk_claims", x => x.id);
                });

            migrationBuilder.CreateTable(
                name: "claim_approvals",
                schema: "insurance_schema",
                columns: table => new
                {
                    id = table.Column<Guid>(type: "uuid", nullable: false),
                    claim_id = table.Column<Guid>(type: "uuid", nullable: false),
                    approver_id = table.Column<Guid>(type: "uuid", nullable: false),
                    approver_role = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: false),
                    approval_level = table.Column<int>(type: "integer", nullable: false),
                    decision = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: false),
                    approved_amount = table.Column<long>(type: "bigint", nullable: false),
                    approved_currency = table.Column<string>(type: "character varying(3)", maxLength: 3, nullable: false, defaultValue: "BDT"),
                    notes = table.Column<string>(type: "text", nullable: true),
                    approved_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    created_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("pk_claim_approvals", x => x.id);
                    table.ForeignKey(
                        name: "fk_claim_approvals_claims_claim_id",
                        column: x => x.claim_id,
                        principalSchema: "insurance_schema",
                        principalTable: "claims",
                        principalColumn: "id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateTable(
                name: "claim_documents",
                schema: "insurance_schema",
                columns: table => new
                {
                    id = table.Column<Guid>(type: "uuid", nullable: false),
                    claim_id = table.Column<Guid>(type: "uuid", nullable: false),
                    document_type = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: false),
                    file_url = table.Column<string>(type: "text", nullable: false),
                    file_hash = table.Column<string>(type: "character varying(64)", maxLength: 64, nullable: false),
                    verified = table.Column<bool>(type: "boolean", nullable: false),
                    verified_by = table.Column<Guid>(type: "uuid", nullable: true),
                    uploaded_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    created_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    updated_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("pk_claim_documents", x => x.id);
                    table.ForeignKey(
                        name: "fk_claim_documents_claims_claim_id",
                        column: x => x.claim_id,
                        principalSchema: "insurance_schema",
                        principalTable: "claims",
                        principalColumn: "id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateTable(
                name: "fraud_checks",
                schema: "insurance_schema",
                columns: table => new
                {
                    id = table.Column<Guid>(type: "uuid", nullable: false),
                    claim_id = table.Column<Guid>(type: "uuid", nullable: false),
                    fraud_score = table.Column<double>(type: "numeric(5,2)", nullable: false),
                    risk_factors = table.Column<List<string>>(type: "text[]", nullable: false),
                    flagged = table.Column<bool>(type: "boolean", nullable: false, defaultValue: false),
                    reviewed_by = table.Column<Guid>(type: "uuid", nullable: true),
                    reviewed_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    created_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("pk_fraud_checks", x => x.id);
                    table.ForeignKey(
                        name: "fk_fraud_checks_claims_claim_id",
                        column: x => x.claim_id,
                        principalSchema: "insurance_schema",
                        principalTable: "claims",
                        principalColumn: "id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateIndex(
                name: "ix_claim_approvals_approver_id",
                schema: "insurance_schema",
                table: "claim_approvals",
                column: "approver_id");

            migrationBuilder.CreateIndex(
                name: "ix_claim_approvals_claim_id",
                schema: "insurance_schema",
                table: "claim_approvals",
                column: "claim_id");

            migrationBuilder.CreateIndex(
                name: "ix_claim_documents_claim_id",
                schema: "insurance_schema",
                table: "claim_documents",
                column: "claim_id");

            migrationBuilder.CreateIndex(
                name: "ix_claim_documents_file_hash",
                schema: "insurance_schema",
                table: "claim_documents",
                column: "file_hash");

            migrationBuilder.CreateIndex(
                name: "ix_claims_claim_number",
                schema: "insurance_schema",
                table: "claims",
                column: "claim_number",
                unique: true);

            migrationBuilder.CreateIndex(
                name: "ix_claims_customer_id",
                schema: "insurance_schema",
                table: "claims",
                column: "customer_id");

            migrationBuilder.CreateIndex(
                name: "ix_claims_incident_date",
                schema: "insurance_schema",
                table: "claims",
                column: "incident_date");

            migrationBuilder.CreateIndex(
                name: "ix_claims_policy_id",
                schema: "insurance_schema",
                table: "claims",
                column: "policy_id");

            migrationBuilder.CreateIndex(
                name: "ix_claims_status",
                schema: "insurance_schema",
                table: "claims",
                column: "status");

            migrationBuilder.CreateIndex(
                name: "ix_fraud_checks_claim_id",
                schema: "insurance_schema",
                table: "fraud_checks",
                column: "claim_id",
                unique: true);

            migrationBuilder.CreateIndex(
                name: "ix_fraud_checks_flagged",
                schema: "insurance_schema",
                table: "fraud_checks",
                column: "flagged");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "claim_approvals",
                schema: "insurance_schema");

            migrationBuilder.DropTable(
                name: "claim_documents",
                schema: "insurance_schema");

            migrationBuilder.DropTable(
                name: "fraud_checks",
                schema: "insurance_schema");

            migrationBuilder.DropTable(
                name: "claims",
                schema: "insurance_schema");
        }
    }
}
