using System;
using System.Collections.Generic;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace InsuranceEngine.Service.Migrations
{
    /// <inheritdoc />
    public partial class init : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.CreateTable(
                name: "beneficiaries",
                columns: table => new
                {
                    beneficiary_id = table.Column<Guid>(type: "uuid", nullable: false),
                    user_id = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: false),
                    partner_id = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    beneficiary_code = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: false),
                    policy_id = table.Column<Guid>(type: "uuid", nullable: true),
                    type = table.Column<string>(type: "character varying(32)", maxLength: 32, nullable: false),
                    status = table.Column<string>(type: "character varying(32)", maxLength: 32, nullable: false),
                    kyc_status = table.Column<string>(type: "character varying(32)", maxLength: 32, nullable: false),
                    kyc_completed_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    risk_score = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: true),
                    referral_code = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    referred_by = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    created_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    updated_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    created_by = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: false),
                    updated_by = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: false),
                    deleted_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    deleted_by = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_beneficiaries", x => x.beneficiary_id);
                });

            migrationBuilder.CreateTable(
                name: "errors",
                columns: table => new
                {
                    error_id = table.Column<Guid>(type: "uuid", nullable: false),
                    code = table.Column<string>(type: "character varying(64)", maxLength: 64, nullable: false),
                    message = table.Column<string>(type: "character varying(500)", maxLength: 500, nullable: false),
                    details = table.Column<Dictionary<string, string>>(type: "jsonb", nullable: true),
                    retryable = table.Column<bool>(type: "boolean", nullable: false),
                    retry_after_seconds = table.Column<int>(type: "integer", nullable: true),
                    http_status_code = table.Column<int>(type: "integer", nullable: false),
                    documentation_url = table.Column<string>(type: "character varying(2048)", maxLength: 2048, nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_errors", x => x.error_id);
                });

            migrationBuilder.CreateTable(
                name: "beneficiary_businesses",
                columns: table => new
                {
                    id = table.Column<Guid>(type: "uuid", nullable: false),
                    beneficiary_id = table.Column<Guid>(type: "uuid", nullable: false),
                    business_name = table.Column<string>(type: "character varying(255)", maxLength: 255, nullable: false),
                    business_name_bn = table.Column<string>(type: "character varying(255)", maxLength: 255, nullable: true),
                    trade_license_number = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: false),
                    trade_license_issue_date = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    trade_license_expiry_date = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    tin_number = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: false),
                    bin_number = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    business_type = table.Column<string>(type: "character varying(32)", maxLength: 32, nullable: false),
                    industry_sector = table.Column<string>(type: "character varying(150)", maxLength: 150, nullable: true),
                    employee_count = table.Column<int>(type: "integer", nullable: true),
                    incorporation_date = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    contact_mobile_number = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: true),
                    contact_email = table.Column<string>(type: "character varying(255)", maxLength: 255, nullable: true),
                    contact_alternate_mobile = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: true),
                    contact_landline = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: true),
                    registered_address_line1 = table.Column<string>(type: "character varying(255)", maxLength: 255, nullable: true),
                    registered_address_line2 = table.Column<string>(type: "character varying(255)", maxLength: 255, nullable: true),
                    registered_city = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    registered_district = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    registered_division = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    registered_postal_code = table.Column<string>(type: "character varying(20)", maxLength: 20, nullable: true),
                    registered_country = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    registered_latitude = table.Column<decimal>(type: "numeric", nullable: true),
                    registered_longitude = table.Column<decimal>(type: "numeric", nullable: true),
                    business_address_line1 = table.Column<string>(type: "character varying(255)", maxLength: 255, nullable: true),
                    business_address_line2 = table.Column<string>(type: "character varying(255)", maxLength: 255, nullable: true),
                    business_city = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    business_district = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    business_division = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    business_postal_code = table.Column<string>(type: "character varying(20)", maxLength: 20, nullable: true),
                    business_country = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    business_latitude = table.Column<decimal>(type: "numeric", nullable: true),
                    business_longitude = table.Column<decimal>(type: "numeric", nullable: true),
                    focal_person_name = table.Column<string>(type: "character varying(255)", maxLength: 255, nullable: false),
                    focal_person_designation = table.Column<string>(type: "character varying(150)", maxLength: 150, nullable: true),
                    focal_person_nid = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    focal_person_mobile_number = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: true),
                    focal_person_email = table.Column<string>(type: "character varying(255)", maxLength: 255, nullable: true),
                    focal_person_alternate_mobile = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: true),
                    focal_person_landline = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: true),
                    created_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    updated_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    created_by = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: false),
                    updated_by = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: false),
                    deleted_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    deleted_by = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_beneficiary_businesses", x => x.id);
                    table.ForeignKey(
                        name: "FK_beneficiary_businesses_beneficiaries_beneficiary_id",
                        column: x => x.beneficiary_id,
                        principalTable: "beneficiaries",
                        principalColumn: "beneficiary_id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateTable(
                name: "beneficiary_individuals",
                columns: table => new
                {
                    id = table.Column<Guid>(type: "uuid", nullable: false),
                    beneficiary_id = table.Column<Guid>(type: "uuid", nullable: false),
                    full_name = table.Column<string>(type: "character varying(200)", maxLength: 200, nullable: false),
                    full_name_bn = table.Column<string>(type: "character varying(200)", maxLength: 200, nullable: true),
                    date_of_birth = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    gender = table.Column<string>(type: "character varying(32)", maxLength: 32, nullable: false),
                    nid_number = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    passport_number = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    birth_certificate_number = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    tin_number = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    marital_status = table.Column<string>(type: "character varying(32)", maxLength: 32, nullable: true),
                    occupation = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    mobile_number = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: true),
                    email = table.Column<string>(type: "character varying(255)", maxLength: 255, nullable: true),
                    alternate_mobile = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: true),
                    landline = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: true),
                    permanent_address_line1 = table.Column<string>(type: "character varying(255)", maxLength: 255, nullable: true),
                    permanent_address_line2 = table.Column<string>(type: "character varying(255)", maxLength: 255, nullable: true),
                    permanent_city = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    permanent_district = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    permanent_division = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    permanent_postal_code = table.Column<string>(type: "character varying(20)", maxLength: 20, nullable: true),
                    permanent_country = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    permanent_latitude = table.Column<decimal>(type: "numeric", nullable: true),
                    permanent_longitude = table.Column<decimal>(type: "numeric", nullable: true),
                    present_address_line1 = table.Column<string>(type: "character varying(255)", maxLength: 255, nullable: true),
                    present_address_line2 = table.Column<string>(type: "character varying(255)", maxLength: 255, nullable: true),
                    present_city = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    present_district = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    present_division = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    present_postal_code = table.Column<string>(type: "character varying(20)", maxLength: 20, nullable: true),
                    present_country = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    present_latitude = table.Column<decimal>(type: "numeric", nullable: true),
                    present_longitude = table.Column<decimal>(type: "numeric", nullable: true),
                    nominee_name = table.Column<string>(type: "character varying(200)", maxLength: 200, nullable: true),
                    nominee_relationship = table.Column<string>(type: "character varying(200)", maxLength: 200, nullable: true),
                    created_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    updated_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    created_by = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: false),
                    updated_by = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: false),
                    deleted_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    deleted_by = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_beneficiary_individuals", x => x.id);
                    table.ForeignKey(
                        name: "FK_beneficiary_individuals_beneficiaries_beneficiary_id",
                        column: x => x.beneficiary_id,
                        principalTable: "beneficiaries",
                        principalColumn: "beneficiary_id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateTable(
                name: "field_violations",
                columns: table => new
                {
                    field_violation_id = table.Column<Guid>(type: "uuid", nullable: false),
                    field = table.Column<string>(type: "character varying(200)", maxLength: 200, nullable: false),
                    code = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: false),
                    description = table.Column<string>(type: "character varying(1000)", maxLength: 1000, nullable: false),
                    rejected_value = table.Column<string>(type: "character varying(500)", maxLength: 500, nullable: false),
                    error_id = table.Column<Guid>(type: "uuid", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_field_violations", x => x.field_violation_id);
                    table.ForeignKey(
                        name: "FK_field_violations_errors_error_id",
                        column: x => x.error_id,
                        principalTable: "errors",
                        principalColumn: "error_id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateIndex(
                name: "IX_beneficiary_businesses_beneficiary_id",
                table: "beneficiary_businesses",
                column: "beneficiary_id",
                unique: true);

            migrationBuilder.CreateIndex(
                name: "IX_beneficiary_individuals_beneficiary_id",
                table: "beneficiary_individuals",
                column: "beneficiary_id",
                unique: true);

            migrationBuilder.CreateIndex(
                name: "IX_field_violations_error_id",
                table: "field_violations",
                column: "error_id");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "beneficiary_businesses");

            migrationBuilder.DropTable(
                name: "beneficiary_individuals");

            migrationBuilder.DropTable(
                name: "field_violations");

            migrationBuilder.DropTable(
                name: "beneficiaries");

            migrationBuilder.DropTable(
                name: "errors");
        }
    }
}
