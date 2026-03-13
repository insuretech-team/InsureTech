/**
 * GET /api/employees/template?format=csv
 *
 * Single table CSV template for bulk employee upload.
 *
 * Columns:
 *   name, employee_id, department_name, email, mobile_number,
 *   date_of_birth, date_of_joining, gender, insurance_category,
 *   coverage_amount (BDT), number_of_dependent, assigned_plan_name,
 *   [gap], available_plan_name (ref), premium_amount_BDT (ref)
 *
 * The last two columns are read-only reference columns showing available
 * plans fetched live from the catalog — users copy a plan name into
 * assigned_plan_name. These columns are ignored by the upload parser.
 *
 * UTF-8 BOM prepended so Excel renders Bengali/Unicode correctly.
 * Dates: DD/MM/YYYY. coverage_amount: plain BDT integer, no currency symbol.
 */

import { NextRequest, NextResponse } from "next/server";

// ─── Employee column headers (uploaded / parsed by gateway) ───────────────────

// Column order matches the portal employee table:
// Name → Employee ID → Department → Email → Mobile → DOB → DOJ → Gender →
// Insurance Category → Assigned Plan → Coverage → Dependents
const EMP_HEADERS = [
  "name",
  "employee_id",
  "department_name",
  "email",
  "mobile_number",
  "date_of_birth",
  "date_of_joining",
  "gender",
  "insurance_category",
  "assigned_plan_name",
  "coverage_amount",
  "number_of_dependent",
];

// ─── Reference column headers removed ────────────────────────────────────────
// Previously we appended reference plan columns to the right of the template.
// This caused "premium_amount (BDT)" and numeric premium values to be parsed as
// department names by the bulk upload handler when the column map misaligned.
// Reference info is now shown only in the instructions, not as extra CSV columns.
const REF_HEADERS: string[] = [];

// ─── Example employee rows (4 rows, Bengali names) ────────────────────────────

// Rows match EMP_HEADERS column order:
// name, employee_id, department_name, email, mobile_number, date_of_birth,
// date_of_joining, gender, insurance_category, assigned_plan_name, coverage_amount, number_of_dependent
const EXAMPLE_ROWS: string[][] = [
  ["মোহাম্মদ রহিম উদ্দিন", "EMP001", "Engineering",     "rahim@company.com",  "+8801712345678", "15/06/1990", "01/01/2023", "MALE",   "HEALTH", "", "500000", "2"],
  ["ফাতেমা বেগম",           "EMP002", "Human Resources", "fatema@company.com", "+8801812345678", "22/03/1988", "15/03/2022", "FEMALE", "LIFE",   "", "300000", "1"],
  ["করিম হোসেন",            "EMP003", "Engineering",     "karim@company.com",  "+8801912345678", "10/11/1992", "01/06/2021", "MALE",   "HEALTH", "", "400000", "3"],
  ["সালমা আক্তার",          "EMP004", "Finance",         "salma@company.com",  "+8801612345678", "05/08/1995", "01/09/2023", "FEMALE", "LIFE",   "", "250000", "0"],
];

// ─── Helpers ──────────────────────────────────────────────────────────────────

function q(value: string): string {
  return `"${value.replace(/"/g, '""')}"`;
}

// ─── Types ────────────────────────────────────────────────────────────────────

interface CatalogPlan {
  planName?: string;
  insuranceCategory?: string;
  premiumAmount?: string;
}

interface CatalogResponse {
  items?: CatalogPlan[];
}

// ─── Route handler ────────────────────────────────────────────────────────────

export async function GET(req: NextRequest) {
  try {
    // ── Fetch plan catalog (best-effort — silently omit if it fails) ──────────
    let plans: CatalogPlan[] = [];
    try {
      const res = await fetch(`${req.nextUrl.origin}/api/purchase-orders/catalog`, {
        headers: { cookie: req.headers.get("cookie") ?? "" },
        cache: "no-store",
      });
      if (res.ok) {
        const data = (await res.json()) as CatalogResponse;
        plans = Array.isArray(data.items) ? data.items : [];
      }
    } catch {
      // silently skip
    }

    // ── Build rows ────────────────────────────────────────────────────────────
    const hasPlans = plans.length > 0;
    const numRows = Math.max(EXAMPLE_ROWS.length, plans.length);

    // Header row: employee cols + gap + ref cols (only if plans available)
    const headerRow = hasPlans
      ? [...EMP_HEADERS, ...REF_HEADERS].map(q).join(",")
      : EMP_HEADERS.map(q).join(",");

    const rows: string[] = [headerRow];

    // assigned_plan_name column index in EMP_HEADERS
    const planNameColIdx = EMP_HEADERS.indexOf("assigned_plan_name");

    for (let i = 0; i < numRows; i++) {
      // Copy example row (or blank row if beyond example rows)
      const emp = (EXAMPLE_ROWS[i] ?? Array(EMP_HEADERS.length).fill("")).map((v) => v);

      // Fill assigned_plan_name from catalog for example rows
      if (hasPlans && planNameColIdx >= 0 && plans[i % plans.length]) {
        emp[planNameColIdx] = plans[i % plans.length].planName ?? "";
      }

      const empCells = emp.map(q).join(",");

      if (hasPlans) {
        const plan = plans[i];
        // premiumAmount formatted as ৳X,XX,XXX — strip all non-numeric to get plain integer
        const premiumRaw = plan ? (plan.premiumAmount ?? "") : "";
        const premiumPlain = premiumRaw.replace(/\D/g, "");
        const planName = plan ? (plan.planName ?? "") : "";
        rows.push(`${empCells},${q(planName)},${q(premiumPlain)}`);
      } else {
        rows.push(empCells);
      }
    }

    // UTF-8 BOM so Excel renders Bengali correctly
    const csv = "\uFEFF" + rows.join("\r\n");

    return new Response(csv, {
      status: 200,
      headers: {
        "Content-Type": "text/csv; charset=utf-8",
        "Content-Disposition": 'attachment; filename="employees_template.csv"',
        "Cache-Control": "no-store",
      },
    });
  } catch (err) {
    return NextResponse.json(
      { ok: false, message: err instanceof Error ? err.message : "Failed to generate template" },
      { status: 500 }
    );
  }
}
