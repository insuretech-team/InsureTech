import { NextResponse } from "next/server";
import { execFile } from "child_process";
import { promisify } from "util";
import path from "path";
import fs from "fs";

const execFileAsync = promisify(execFile);

/**
 * Maps every document card (by exact sheet ID, then category+kind fallback)
 * to a template definition filename in:
 *   backend/inscore/templates/insurance/definitions/<name>.json
 *
 * To add a new document type:
 * 1. Create the template definition JSON in the definitions folder
 * 2. Add the mapping here
 * No Python code changes needed.
 */
const TEMPLATE_MAP: Record<string, string> = {
  // ── 14 Workbook Sheets — each maps to its own unique template ────────────
  "sheet-pragati-sheet1":  "overseas_mediclaim_proposal",   // Overseas Mediclaim Proposal Form
  "sheet-pragati-sheet2":  "mediclaim_medical_history",     // Mediclaim Medical History Questionnaire
  "sheet-pragati-sheet3":  "mediclaim_declaration",         // Mediclaim Declaration & Benefits Schedule
  "sheet-pragati-sheet4":  "travel_rate_table",             // Non-Schengen Premium Matrix
  "sheet-pragati-sheet5":  "travel_rate_table",             // Travel Addendum & Employment / Study Rates
  "sheet-pragati-sheet6":  "travel_rate_table",             // Schengen Premium Matrix
  "sheet-pragati-sheet7":  "travel_rate_table",             // Schengen Frequent Travel Addendum
  "sheet-pragati-sheet8":  "private_vehicle_proposal",      // Private Vehicle Proposal
  "sheet-pragati-sheet9":  "private_vehicle_proposal",      // Private Vehicle Declaration Continuation
  "sheet-pragati-sheet10": "fire_proposal",                 // Fire Insurance Proposal
  "sheet-pragati-sheet11": "commercial_vehicle_proposal",   // Commercial Vehicle Proposal
  "sheet-pragati-sheet12": "livestock_proposal",            // Livestock Proposal
  "sheet-pragati-sheet13": "member_census",                 // Member Census Schedule
  "sheet-pragati-sheet14": "health_claim",                  // Health Insurance Claim Form
  // ── 7 Reference Documents ────────────────────────────────────────────────
  "reference-documents-are-normally-required-for-claims-docx":                          "claims_required_docs",
  "reference-omp-new-claim-process-docx":                                               "mediclaim_declaration",
  "reference-omp-proposal-form-new-pdf":                                                "overseas_mediclaim_proposal",
  "reference-motor-insurance-proposal-form-pdf":                                        "private_vehicle_proposal",
  "reference-fire-insurance-proposal-form-20230622-0001-pdf":                           "fire_proposal",
  "reference-financial-proposal-lifeplus-shanta-2026-group-insurance-pdf":              "group_life_proposal",
  "reference-motor-insurance-policy-lifeplus-bangladesh-final-underwrite-claims-pptx":  "private_vehicle_proposal",
};

function resolveTemplate(docId: string, category: string, title = "", kind = ""): string | null {
  if (TEMPLATE_MAP[docId]) return TEMPLATE_MAP[docId];

  const cat = category.toLowerCase();
  const ttl = title.toLowerCase();
  const knd = kind.toLowerCase();

  if (cat === "travel" || ttl.includes("mediclaim") || ttl.includes("overseas")) return "overseas_mediclaim_proposal";
  if (cat === "auto" || cat === "commercial vehicle" || ttl.includes("vehicle") || ttl.includes("motor")) return "private_vehicle_proposal";
  if (cat === "fire" || ttl.includes("fire")) return "fire_proposal";
  if (cat === "health" || knd === "claim-form" || ttl.includes("claim")) return "health_claim";

  return null;
}

export async function POST(request: Request) {
  let body: Record<string, unknown> = {};
  try {
    body = (await request.json()) as Record<string, unknown>;
  } catch {
    return NextResponse.json({ ok: false, message: "Invalid request body" }, { status: 400 });
  }

  const { documentId, category, title, kind, data, filename } = body as {
    documentId?: string;
    category?: string;
    title?: string;
    kind?: string;
    data?: Record<string, unknown>;
    filename?: string;
  };

  const templateName = resolveTemplate(documentId ?? "", category ?? "", title, kind);

  if (!templateName) {
    return NextResponse.json({
      ok: false,
      message: `No template definition found for "${title ?? category}". Add a definition JSON to templates/insurance/definitions/ and map it here.`,
    }, { status: 404 });
  }

  const projectRoot = path.resolve(process.cwd(), "..");
  const generatedDir    = path.join(projectRoot, "backend", "inscore", "generated");
  const venvPython      = path.join(projectRoot, ".venv", "Scripts", "python.exe");
  const sidecarDir      = path.join(projectRoot, "backend", "inscore", "microservices", "docgen", "sidecar");
  const definitionsDir  = path.join(projectRoot, "backend", "inscore", "templates", "insurance", "definitions");
  const logoPath        = path.join(projectRoot, "web_shared", "insurers", "pragati_logo.png");

  fs.mkdirSync(generatedDir, { recursive: true });

  const safeFile = ((filename ?? `${templateName}_${Date.now()}`).replace(/[^a-zA-Z0-9_\-]/g, "_")) + ".docx";
  const outPath  = path.join(generatedDir, safeFile);

  // Merge logo path into data (generator reads logo_path from data)
  const fullData: Record<string, unknown> = {
    logo_path: logoPath,
    proposal_id: documentId ?? "",
    generated_at: new Date().toLocaleString("en-GB", {
      day: "2-digit", month: "short", year: "numeric",
      hour: "2-digit", minute: "2-digit",
    }),
    ...(data ?? {}),
  };

  const tmpJson    = path.join(generatedDir, `_tmp_${Date.now()}.json`);
  const runnerPath = path.join(generatedDir, `_runner_${Date.now()}.py`);

  fs.writeFileSync(tmpJson, JSON.stringify(fullData), "utf8");

  const tplPath = path.join(definitionsDir, `${templateName}.json`).replace(/\\/g, "\\\\");
  const runner = `
import sys, json, io
sys.path.insert(0, r"${sidecarDir.replace(/\\/g, "\\\\")}")
from docx_builder import DocxBuilder

with open(r"${tplPath}", encoding="utf-8") as f:
    template = json.load(f)

with open(r"${tmpJson.replace(/\\/g, "\\\\")}", encoding="utf-8") as f:
    data = json.load(f)

buf = DocxBuilder(template, data).build()
with open(r"${outPath.replace(/\\/g, "\\\\")}", "wb") as f:
    f.write(buf.read())

print("OK")
`.trim();

  fs.writeFileSync(runnerPath, runner, "utf8");

  try {
    const { stdout, stderr } = await execFileAsync(venvPython, [runnerPath], { timeout: 30000 });
    fs.unlinkSync(tmpJson);
    fs.unlinkSync(runnerPath);

    if (!stdout.includes("OK")) {
      return NextResponse.json({ ok: false, message: stderr || stdout || "Generator failed" });
    }

    return NextResponse.json({
      ok: true,
      data: {
        filename: safeFile,
        downloadUrl: `/api/documents/file/${encodeURIComponent(safeFile)}`,
        message: `${title ?? templateName} generated successfully.`,
      },
    });
  } catch (err: unknown) {
    try { fs.unlinkSync(tmpJson); } catch { /* ignore */ }
    try { fs.unlinkSync(runnerPath); } catch { /* ignore */ }
    const message = err instanceof Error ? err.message : String(err);
    return NextResponse.json({ ok: false, message: `Generator error: ${message}` });
  }
}
