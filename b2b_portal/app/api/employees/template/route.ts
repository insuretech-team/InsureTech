/**
 * GET /api/employees/template
 *
 * Serves the canonical employee bulk-upload CSV template directly from the
 * backend template file at:
 *   backend/inscore/templates/b2b/employees_template.csv
 *
 * Both the "Download Template" button in the employee data-table toolbar
 * and the bulk-upload modal use this endpoint.
 */

import { readFileSync } from "fs";
import { join } from "path";
import { NextResponse } from "next/server";

// Path to the canonical backend template (relative to the b2b_portal project root).
// process.cwd() == b2b_portal root in Next.js server context.
const TEMPLATE_PATH = join(
  process.cwd(),
  "..",
  "backend",
  "inscore",
  "templates",
  "b2b",
  "employees_template.csv"
);

export async function GET() {
  try {
    const csv = readFileSync(TEMPLATE_PATH);
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
      { ok: false, message: err instanceof Error ? err.message : "Template file not found" },
      { status: 500 }
    );
  }
}
