using System;
using System.Collections.Generic;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace InsuranceEngine.Products.Infrastructure.Persistence.Migrations
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
                name: "products",
                schema: "insurance_schema",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uuid", nullable: false),
                    ProductCode = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: false),
                    ProductName = table.Column<string>(type: "character varying(255)", maxLength: 255, nullable: false),
                    ProductNameBn = table.Column<string>(type: "text", nullable: true),
                    Description = table.Column<string>(type: "text", nullable: true),
                    DescriptionBn = table.Column<string>(type: "text", nullable: true),
                    Category = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: false),
                    Status = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: false),
                    base_premium = table.Column<long>(type: "bigint", nullable: false),
                    base_premium_currency = table.Column<string>(type: "character varying(3)", maxLength: 3, nullable: false, defaultValue: "BDT"),
                    min_sum_insured = table.Column<long>(type: "bigint", nullable: false),
                    min_sum_insured_currency = table.Column<string>(type: "character varying(3)", maxLength: 3, nullable: false, defaultValue: "BDT"),
                    max_sum_insured = table.Column<long>(type: "bigint", nullable: false),
                    max_sum_insured_currency = table.Column<string>(type: "character varying(3)", maxLength: 3, nullable: false, defaultValue: "BDT"),
                    MinAge = table.Column<int>(type: "integer", nullable: false),
                    MaxAge = table.Column<int>(type: "integer", nullable: false),
                    MinTenureMonths = table.Column<int>(type: "integer", nullable: false),
                    MaxTenureMonths = table.Column<int>(type: "integer", nullable: false),
                    Exclusions = table.Column<List<string>>(type: "text[]", nullable: false),
                    ProductAttributes = table.Column<string>(type: "jsonb", nullable: true),
                    DeductiblePercentage = table.Column<double>(type: "double precision", nullable: false),
                    CoPayPercentage = table.Column<double>(type: "double precision", nullable: false),
                    MaxDeductibleAmount = table.Column<long>(type: "bigint", nullable: false),
                    CreatedBy = table.Column<Guid>(type: "uuid", nullable: false),
                    CreatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    UpdatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    DeletedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    IsDeleted = table.Column<bool>(type: "boolean", nullable: false),
                    TenantId = table.Column<Guid>(type: "uuid", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_products", x => x.Id);
                });

            migrationBuilder.CreateTable(
                name: "pricing_configs",
                schema: "insurance_schema",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uuid", nullable: false),
                    ProductId = table.Column<Guid>(type: "uuid", nullable: false),
                    Rules = table.Column<string>(type: "jsonb", nullable: false),
                    EffectiveFrom = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    EffectiveTo = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    CreatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    UpdatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_pricing_configs", x => x.Id);
                    table.ForeignKey(
                        name: "FK_pricing_configs_products_ProductId",
                        column: x => x.ProductId,
                        principalSchema: "insurance_schema",
                        principalTable: "products",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateTable(
                name: "product_plans",
                schema: "insurance_schema",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uuid", nullable: false),
                    ProductId = table.Column<Guid>(type: "uuid", nullable: false),
                    PlanName = table.Column<string>(type: "character varying(255)", maxLength: 255, nullable: false),
                    PlanNameBn = table.Column<string>(type: "text", nullable: true),
                    PlanDescription = table.Column<string>(type: "text", nullable: true),
                    DescriptionBn = table.Column<string>(type: "text", nullable: true),
                    premium_amount = table.Column<long>(type: "bigint", nullable: false),
                    premium_currency = table.Column<string>(type: "character varying(3)", maxLength: 3, nullable: false, defaultValue: "BDT"),
                    min_sum_insured = table.Column<long>(type: "bigint", nullable: false),
                    min_sum_insured_currency = table.Column<string>(type: "character varying(3)", maxLength: 3, nullable: false, defaultValue: "BDT"),
                    max_sum_insured = table.Column<long>(type: "bigint", nullable: false),
                    max_sum_insured_currency = table.Column<string>(type: "character varying(3)", maxLength: 3, nullable: false, defaultValue: "BDT"),
                    Attributes = table.Column<string>(type: "jsonb", nullable: true),
                    CreatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    UpdatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_product_plans", x => x.Id);
                    table.ForeignKey(
                        name: "FK_product_plans_products_ProductId",
                        column: x => x.ProductId,
                        principalSchema: "insurance_schema",
                        principalTable: "products",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateTable(
                name: "product_riders",
                schema: "insurance_schema",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uuid", nullable: false),
                    ProductId = table.Column<Guid>(type: "uuid", nullable: false),
                    RiderName = table.Column<string>(type: "character varying(255)", maxLength: 255, nullable: false),
                    Description = table.Column<string>(type: "text", nullable: true),
                    premium_amount = table.Column<long>(type: "bigint", nullable: false),
                    premium_currency = table.Column<string>(type: "character varying(3)", maxLength: 3, nullable: false, defaultValue: "BDT"),
                    coverage_amount = table.Column<long>(type: "bigint", nullable: false),
                    coverage_currency = table.Column<string>(type: "character varying(3)", maxLength: 3, nullable: false, defaultValue: "BDT"),
                    IsMandatory = table.Column<bool>(type: "boolean", nullable: false),
                    CreatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    UpdatedAt = table.Column<DateTime>(type: "timestamp with time zone", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_product_riders", x => x.Id);
                    table.ForeignKey(
                        name: "FK_product_riders_products_ProductId",
                        column: x => x.ProductId,
                        principalSchema: "insurance_schema",
                        principalTable: "products",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateTable(
                name: "risk_assessment_questions",
                schema: "insurance_schema",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uuid", nullable: false),
                    ProductId = table.Column<Guid>(type: "uuid", nullable: false),
                    ProductId1 = table.Column<Guid>(type: "uuid", nullable: true),
                    QuestionText = table.Column<string>(type: "text", nullable: false),
                    QuestionTextBn = table.Column<string>(type: "text", nullable: true),
                    OptionsJson = table.Column<string>(type: "text", nullable: false),
                    Weight = table.Column<int>(type: "integer", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_risk_assessment_questions", x => x.Id);
                    table.ForeignKey(
                        name: "FK_risk_assessment_questions_products_ProductId",
                        column: x => x.ProductId,
                        principalSchema: "insurance_schema",
                        principalTable: "products",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                    table.ForeignKey(
                        name: "FK_risk_assessment_questions_products_ProductId1",
                        column: x => x.ProductId1,
                        principalSchema: "insurance_schema",
                        principalTable: "products",
                        principalColumn: "Id");
                });

            migrationBuilder.CreateIndex(
                name: "IX_pricing_configs_ProductId",
                schema: "insurance_schema",
                table: "pricing_configs",
                column: "ProductId",
                unique: true);

            migrationBuilder.CreateIndex(
                name: "IX_product_plans_ProductId",
                schema: "insurance_schema",
                table: "product_plans",
                column: "ProductId");

            migrationBuilder.CreateIndex(
                name: "IX_product_riders_ProductId",
                schema: "insurance_schema",
                table: "product_riders",
                column: "ProductId");

            migrationBuilder.CreateIndex(
                name: "IX_products_Category",
                schema: "insurance_schema",
                table: "products",
                column: "Category");

            migrationBuilder.CreateIndex(
                name: "IX_products_ProductCode",
                schema: "insurance_schema",
                table: "products",
                column: "ProductCode",
                unique: true);

            migrationBuilder.CreateIndex(
                name: "IX_products_Status",
                schema: "insurance_schema",
                table: "products",
                column: "Status");

            migrationBuilder.CreateIndex(
                name: "IX_risk_assessment_questions_ProductId",
                schema: "insurance_schema",
                table: "risk_assessment_questions",
                column: "ProductId");

            migrationBuilder.CreateIndex(
                name: "IX_risk_assessment_questions_ProductId1",
                schema: "insurance_schema",
                table: "risk_assessment_questions",
                column: "ProductId1");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "pricing_configs",
                schema: "insurance_schema");

            migrationBuilder.DropTable(
                name: "product_plans",
                schema: "insurance_schema");

            migrationBuilder.DropTable(
                name: "product_riders",
                schema: "insurance_schema");

            migrationBuilder.DropTable(
                name: "risk_assessment_questions",
                schema: "insurance_schema");

            migrationBuilder.DropTable(
                name: "products",
                schema: "insurance_schema");
        }
    }
}
