using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace InsuranceEngine.Products.Migrations
{
    /// <inheritdoc />
    public partial class InitialProductsMigration : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.CreateTable(
                name: "products",
                columns: table => new
                {
                    id = table.Column<Guid>(type: "uuid", nullable: false),
                    product_code = table.Column<string>(type: "text", nullable: false),
                    product_name = table.Column<string>(type: "text", nullable: false),
                    product_name_bn = table.Column<string>(type: "text", nullable: true),
                    description = table.Column<string>(type: "text", nullable: true),
                    description_bn = table.Column<string>(type: "text", nullable: true),
                    category = table.Column<string>(type: "text", nullable: false),
                    status = table.Column<string>(type: "text", nullable: false),
                    min_sum_insured = table.Column<decimal>(type: "numeric", nullable: false),
                    max_sum_insured = table.Column<decimal>(type: "numeric", nullable: false),
                    min_age = table.Column<int>(type: "integer", nullable: false),
                    max_age = table.Column<int>(type: "integer", nullable: false),
                    min_tenure_months = table.Column<int>(type: "integer", nullable: false),
                    max_tenure_months = table.Column<int>(type: "integer", nullable: false),
                    created_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    updated_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    tenant_id = table.Column<Guid>(type: "uuid", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("pk_products", x => x.id);
                });

            migrationBuilder.CreateTable(
                name: "pricing_rules",
                columns: table => new
                {
                    id = table.Column<Guid>(type: "uuid", nullable: false),
                    product_id = table.Column<Guid>(type: "uuid", nullable: false),
                    product_id1 = table.Column<Guid>(type: "uuid", nullable: true),
                    rule_name = table.Column<string>(type: "text", nullable: false),
                    rule_expression = table.Column<string>(type: "text", nullable: false),
                    adjustment_amount = table.Column<decimal>(type: "numeric", nullable: false),
                    is_percentage = table.Column<bool>(type: "boolean", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("pk_pricing_rules", x => x.id);
                    table.ForeignKey(
                        name: "fk_pricing_rules_products_product_id",
                        column: x => x.product_id,
                        principalTable: "products",
                        principalColumn: "id",
                        onDelete: ReferentialAction.Cascade);
                    table.ForeignKey(
                        name: "fk_pricing_rules_products_product_id1",
                        column: x => x.product_id1,
                        principalTable: "products",
                        principalColumn: "id");
                });

            migrationBuilder.CreateTable(
                name: "product_plans",
                columns: table => new
                {
                    id = table.Column<Guid>(type: "uuid", nullable: false),
                    product_id = table.Column<Guid>(type: "uuid", nullable: false),
                    product_id1 = table.Column<Guid>(type: "uuid", nullable: true),
                    plan_name = table.Column<string>(type: "text", nullable: false),
                    plan_name_bn = table.Column<string>(type: "text", nullable: true),
                    description = table.Column<string>(type: "text", nullable: true),
                    description_bn = table.Column<string>(type: "text", nullable: true),
                    premium_amount = table.Column<decimal>(type: "numeric", nullable: false),
                    sum_insured = table.Column<decimal>(type: "numeric", nullable: false),
                    is_unit_wise = table.Column<bool>(type: "boolean", nullable: false),
                    unit_price = table.Column<decimal>(type: "numeric", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("pk_product_plans", x => x.id);
                    table.ForeignKey(
                        name: "fk_product_plans_products_product_id",
                        column: x => x.product_id,
                        principalTable: "products",
                        principalColumn: "id",
                        onDelete: ReferentialAction.Cascade);
                    table.ForeignKey(
                        name: "fk_product_plans_products_product_id1",
                        column: x => x.product_id1,
                        principalTable: "products",
                        principalColumn: "id");
                });

            migrationBuilder.CreateTable(
                name: "risk_assessment_questions",
                columns: table => new
                {
                    id = table.Column<Guid>(type: "uuid", nullable: false),
                    product_id = table.Column<Guid>(type: "uuid", nullable: false),
                    product_id1 = table.Column<Guid>(type: "uuid", nullable: true),
                    question_text = table.Column<string>(type: "text", nullable: false),
                    question_text_bn = table.Column<string>(type: "text", nullable: true),
                    options_json = table.Column<string>(type: "text", nullable: false),
                    weight = table.Column<int>(type: "integer", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("pk_risk_assessment_questions", x => x.id);
                    table.ForeignKey(
                        name: "fk_risk_assessment_questions_products_product_id",
                        column: x => x.product_id,
                        principalTable: "products",
                        principalColumn: "id",
                        onDelete: ReferentialAction.Cascade);
                    table.ForeignKey(
                        name: "fk_risk_assessment_questions_products_product_id1",
                        column: x => x.product_id1,
                        principalTable: "products",
                        principalColumn: "id");
                });

            migrationBuilder.CreateIndex(
                name: "ix_pricing_rules_product_id",
                table: "pricing_rules",
                column: "product_id");

            migrationBuilder.CreateIndex(
                name: "ix_pricing_rules_product_id1",
                table: "pricing_rules",
                column: "product_id1");

            migrationBuilder.CreateIndex(
                name: "ix_product_plans_product_id",
                table: "product_plans",
                column: "product_id");

            migrationBuilder.CreateIndex(
                name: "ix_product_plans_product_id1",
                table: "product_plans",
                column: "product_id1");

            migrationBuilder.CreateIndex(
                name: "ix_risk_assessment_questions_product_id",
                table: "risk_assessment_questions",
                column: "product_id");

            migrationBuilder.CreateIndex(
                name: "ix_risk_assessment_questions_product_id1",
                table: "risk_assessment_questions",
                column: "product_id1");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "pricing_rules");

            migrationBuilder.DropTable(
                name: "product_plans");

            migrationBuilder.DropTable(
                name: "risk_assessment_questions");

            migrationBuilder.DropTable(
                name: "products");
        }
    }
}
