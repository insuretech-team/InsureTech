/**
 * Document Library API
 *
 * Serves the full structured document library from two sources (merged):
 *   1. The template definitions saved in backend/inscore/templates/insurance/definitions/
 *   2. A seed file that stores human-readable metadata (title, category, stage, kind, pack, summary)
 *
 * This replaces the hardcoded pragati-documents.ts / docs-forms-workbooks.json approach.
 * Any new template saved via Template Creator instantly appears here.
 */
import { NextResponse } from "next/server";
import path from "path";
import fs from "fs";

export interface LibraryDocument {
  id: string;
  title: string;
  category: string;
  stage: string;
  kind: string;
  summary: string;
  owner: string;
  uploadStatus: string;
  suggestedUse: string;
  packId: string;
  templateDefinitionId: string;  // maps to definitions/<id>.json
  format: string;
  isGenerated: boolean;          // true = backed by DocxBuilder template
}

export interface LibraryPack {
  id: string;
  title: string;
  category: string;
  stage: string;
  description: string;
  requiredFor: string[];
  notes: string[];
  documentIds: string[];
}

export interface LibraryResponse {
  documents: LibraryDocument[];
  packs: LibraryPack[];
  source: "db" | "seed";
}

const PROJECT_ROOT = path.resolve(process.cwd(), "..");
const DEFS_DIR   = path.join(PROJECT_ROOT, "backend", "inscore", "templates", "insurance", "definitions");
const SEED_FILE  = path.join(PROJECT_ROOT, "backend", "inscore", "templates", "insurance", "library_seed.json");

// ── Default seed — matches all 21 original cards with proper metadata ─────────
const DEFAULT_SEED: { documents: LibraryDocument[]; packs: LibraryPack[] } = {
  documents: [
    // ── Travel / Overseas Mediclaim ──────────────────────────────────────────
    { id: "doc-omp-proposal", title: "Overseas Mediclaim Proposal Form", category: "Travel", stage: "Proposal", kind: "proposal-form", summary: "Full proposal form for overseas mediclaim policy covering business and holiday travel. Available for ages 6 months to 79 years.", owner: "Pragati Insurance PLC", uploadStatus: "Active", suggestedUse: "Complete this form for all outbound overseas travel insurance proposals.", packId: "omp-pack", templateDefinitionId: "overseas_mediclaim_proposal", format: "docx", isGenerated: true },
    { id: "doc-omp-medical", title: "Mediclaim Medical History Questionnaire", category: "Travel", stage: "Underwriting", kind: "medical-questionnaire", summary: "Medical history form to be completed by proposer and spouse. Required for all overseas mediclaim proposals.", owner: "Pragati Insurance PLC", uploadStatus: "Active", suggestedUse: "Attach to the proposal form during underwriting. All questions must be answered Yes or No.", packId: "omp-pack", templateDefinitionId: "mediclaim_medical_history", format: "docx", isGenerated: true },
    { id: "doc-omp-declaration", title: "Mediclaim Declaration & Benefits Schedule", category: "Travel", stage: "Issuance", kind: "declaration", summary: "Proposer declaration, policy acceptance, Schengen country list, and full product benefits & limitations schedule.", owner: "Pragati Insurance PLC", uploadStatus: "Active", suggestedUse: "Issue alongside the policy certificate. Client must sign the declaration before travel.", packId: "omp-pack", templateDefinitionId: "mediclaim_declaration", format: "docx", isGenerated: true },
    { id: "doc-omp-rate-non-schengen", title: "Non-Schengen Premium Rate Matrix", category: "Travel", stage: "Pricing", kind: "rate-table", summary: "Premium rate schedule for worldwide travel excluding USA & Canada and including USA & Canada (Non-Schengen plans).", owner: "Pragati Insurance PLC", uploadStatus: "Active", suggestedUse: "Use for quoting and pricing non-Schengen travel plans by age band and duration.", packId: "omp-pack", templateDefinitionId: "travel_rate_table", format: "docx", isGenerated: true },
    { id: "doc-omp-rate-addendum", title: "Travel Addendum — Employment & Student Rates", category: "Travel", stage: "Pricing", kind: "rate-table", summary: "Special rate addendum for employment and student visa travel categories.", owner: "Pragati Insurance PLC", uploadStatus: "Active", suggestedUse: "Apply these rates for employment or study visa applicants travelling abroad.", packId: "omp-pack", templateDefinitionId: "travel_rate_table", format: "docx", isGenerated: true },
    { id: "doc-omp-rate-schengen", title: "Schengen Premium Rate Matrix", category: "Travel", stage: "Pricing", kind: "rate-table", summary: "Premium rate schedule for Schengen country travel (Euro 30,000 cover, nil deductible).", owner: "Pragati Insurance PLC", uploadStatus: "Active", suggestedUse: "Use for all Schengen visa applications requiring mandatory travel insurance.", packId: "omp-pack", templateDefinitionId: "travel_rate_table", format: "docx", isGenerated: true },
    { id: "doc-omp-rate-schengen-frequent", title: "Schengen Frequent Travel Rate Addendum", category: "Travel", stage: "Pricing", kind: "rate-table", summary: "Rate addendum for frequent Schengen travellers — multi-trip annual plan pricing.", owner: "Pragati Insurance PLC", uploadStatus: "Active", suggestedUse: "Apply for clients who travel to Schengen countries multiple times per year.", packId: "omp-pack", templateDefinitionId: "travel_rate_table", format: "docx", isGenerated: true },
    // ── Auto / Motor ─────────────────────────────────────────────────────────
    { id: "doc-pvehicle-proposal", title: "Private Vehicle Insurance Proposal Form", category: "Auto", stage: "Proposal", kind: "proposal-form", summary: "Motor insurance proposal form for private vehicles (cars, jeeps, SUVs). Covers comprehensive and third-party liability.", owner: "Pragati Insurance PLC", uploadStatus: "Active", suggestedUse: "Complete this form for all private vehicle motor insurance proposals submitted via the bank channel.", packId: "motor-pack", templateDefinitionId: "private_vehicle_proposal", format: "docx", isGenerated: true },
    { id: "doc-pvehicle-declaration", title: "Private Vehicle — Declaration Continuation", category: "Auto", stage: "Proposal", kind: "declaration", summary: "Continuation sheet for private vehicle proposal covering additional declarations and bank use section.", owner: "Pragati Insurance PLC", uploadStatus: "Active", suggestedUse: "Attach to the private vehicle proposal when additional declaration space is required.", packId: "motor-pack", templateDefinitionId: "private_vehicle_proposal", format: "docx", isGenerated: true },
    // ── Fire ─────────────────────────────────────────────────────────────────
    { id: "doc-fire-proposal", title: "Fire & Allied Perils Insurance Proposal Form", category: "Fire", stage: "Proposal", kind: "proposal-form", summary: "Proposal form for fire and allied perils insurance covering buildings, contents, machinery, stock, and electronic equipment.", owner: "Pragati Insurance PLC", uploadStatus: "Active", suggestedUse: "Complete for all fire insurance proposals for residential, commercial, or industrial properties.", packId: "fire-pack", templateDefinitionId: "fire_proposal", format: "docx", isGenerated: true },
    // ── Commercial Vehicle ────────────────────────────────────────────────────
    { id: "doc-cvehicle-proposal", title: "Commercial Vehicle Insurance Proposal Form", category: "Commercial Vehicle", stage: "Proposal", kind: "proposal-form", summary: "Motor insurance proposal form for commercial vehicles — trucks, buses, covered vans, and goods carriers.", owner: "Pragati Insurance PLC", uploadStatus: "Active", suggestedUse: "Complete for all commercial vehicle insurance proposals including route permits and goods types.", packId: "motor-pack", templateDefinitionId: "commercial_vehicle_proposal", format: "docx", isGenerated: true },
    // ── Livestock ─────────────────────────────────────────────────────────────
    { id: "doc-livestock-proposal", title: "Livestock Insurance Proposal Form", category: "Livestock", stage: "Proposal", kind: "proposal-form", summary: "Proposal form for livestock insurance covering cattle, poultry, and other farm animals against mortality and disease.", owner: "Pragati Insurance PLC", uploadStatus: "Active", suggestedUse: "Complete for all livestock insurance proposals. Attach veterinary certificates for each animal.", packId: "livestock-pack", templateDefinitionId: "livestock_proposal", format: "docx", isGenerated: true },
    // ── Health / Group ────────────────────────────────────────────────────────
    { id: "doc-member-census", title: "Group Health — Member Census Schedule", category: "Health", stage: "Enrollment", kind: "schedule", summary: "Employee and dependent census register for group health insurance enrollment. Lists all principal members and their dependents.", owner: "Pragati Insurance PLC", uploadStatus: "Active", suggestedUse: "Submit with the group proposal. Update within 30 days of any member addition, deletion, or amendment.", packId: "group-health-pack", templateDefinitionId: "member_census", format: "docx", isGenerated: true },
    { id: "doc-health-claim", title: "Health Insurance Reimbursement Claim Form", category: "Health", stage: "Claims", kind: "claim-form", summary: "Reimbursement claim form for hospitalisation and medical expenses under group or individual health policies.", owner: "Pragati Insurance PLC", uploadStatus: "Active", suggestedUse: "Submit within 90 days of discharge with all original bills, discharge summary, and investigation reports.", packId: "group-health-pack", templateDefinitionId: "health_claim", format: "docx", isGenerated: true },
    // ── Reference Documents ───────────────────────────────────────────────────
    { id: "ref-claims-docs-required", title: "Documents Required for Claims — All Lines", category: "Claims", stage: "Claims", kind: "process-note", summary: "Master checklist of documents required for health, motor, and fire insurance claims. Includes submission process and contact details.", owner: "Pragati Insurance PLC", uploadStatus: "Active", suggestedUse: "Share with claimants at the time of intimation. Available in the insurer portal document library.", packId: "claims-pack", templateDefinitionId: "claims_required_docs", format: "docx", isGenerated: true },
    { id: "ref-omp-claim-process", title: "Overseas Mediclaim — Claim Process Note", category: "Travel", stage: "Claims", kind: "process-note", summary: "Step-by-step claim submission guide for overseas mediclaim policies including emergency contact numbers and TPA details.", owner: "Pragati Insurance PLC", uploadStatus: "Active", suggestedUse: "Distribute to all OMP policyholders before travel. Include in policy issuance kit.", packId: "omp-pack", templateDefinitionId: "mediclaim_declaration", format: "docx", isGenerated: true },
    { id: "ref-omp-proposal-pdf", title: "Overseas Mediclaim Proposal — Source Form", category: "Travel", stage: "Reference", kind: "reference-file", summary: "Original OMP proposal form PDF from Pragati Insurance for reference and comparison with generated DOCX.", owner: "Pragati Insurance PLC", uploadStatus: "Reference", suggestedUse: "Use as reference to verify the generated DOCX matches the official form layout.", packId: "omp-pack", templateDefinitionId: "overseas_mediclaim_proposal", format: "docx", isGenerated: true },
    { id: "ref-motor-proposal-pdf", title: "Motor Insurance Proposal — Source Form", category: "Auto", stage: "Reference", kind: "reference-file", summary: "Original private vehicle insurance proposal PDF from Pragati Insurance for reference.", owner: "Pragati Insurance PLC", uploadStatus: "Reference", suggestedUse: "Use as reference to verify the generated DOCX matches the official form layout.", packId: "motor-pack", templateDefinitionId: "private_vehicle_proposal", format: "docx", isGenerated: true },
    { id: "ref-fire-proposal-pdf", title: "Fire Insurance Proposal — Source Form", category: "Fire", stage: "Reference", kind: "reference-file", summary: "Original fire insurance proposal PDF from Pragati Insurance for reference.", owner: "Pragati Insurance PLC", uploadStatus: "Reference", suggestedUse: "Use as reference to verify the generated DOCX matches the official form layout.", packId: "fire-pack", templateDefinitionId: "fire_proposal", format: "docx", isGenerated: true },
    { id: "ref-group-financial-proposal", title: "Group Life & Health — Financial Proposal", category: "Group Life", stage: "Pricing", kind: "rate-table", summary: "Financial proposal template for group life and health insurance covering coverage structure, premiums, and terms.", owner: "Pragati Insurance PLC", uploadStatus: "Active", suggestedUse: "Use when preparing commercial proposals for corporate group insurance clients.", packId: "group-health-pack", templateDefinitionId: "group_life_proposal", format: "docx", isGenerated: true },
    { id: "ref-motor-policy-deck", title: "Motor Policy — Underwriting & Claims Reference", category: "Auto", stage: "Reference", kind: "reference-file", summary: "Reference deck covering motor insurance policy terms, underwriting guidelines, and claims settlement procedures.", owner: "Pragati Insurance PLC", uploadStatus: "Reference", suggestedUse: "Use during underwriting review and claims assessment for motor policies.", packId: "motor-pack", templateDefinitionId: "private_vehicle_proposal", format: "docx", isGenerated: true },
  ],
  packs: [
    { id: "omp-pack", title: "Overseas Mediclaim Pack", category: "Travel", stage: "Proposal", description: "Complete document set for overseas mediclaim insurance — proposal, medical history, declaration, rate tables, and claim process.", requiredFor: ["Outbound travel insurance", "Schengen visa applications", "Corporate travel programmes"], notes: ["Submit proposal + medical history together for underwriting.", "Declaration must be signed before policy issuance."], documentIds: ["doc-omp-proposal", "doc-omp-medical", "doc-omp-declaration", "doc-omp-rate-non-schengen", "doc-omp-rate-schengen", "doc-omp-rate-addendum", "doc-omp-rate-schengen-frequent", "ref-omp-claim-process", "ref-omp-proposal-pdf"] },
    { id: "motor-pack", title: "Motor Insurance Pack", category: "Auto", stage: "Proposal", description: "Complete document set for private and commercial vehicle insurance — proposals, declarations, and reference materials.", requiredFor: ["Private vehicle insurance", "Commercial vehicle insurance", "Fleet insurance"], notes: ["Use private vehicle form for cars/SUVs, commercial form for trucks/buses."], documentIds: ["doc-pvehicle-proposal", "doc-pvehicle-declaration", "doc-cvehicle-proposal", "ref-motor-proposal-pdf", "ref-motor-policy-deck"] },
    { id: "fire-pack", title: "Fire Insurance Pack", category: "Fire", stage: "Proposal", description: "Proposal and reference documents for fire and allied perils insurance.", requiredFor: ["Building insurance", "Contents insurance", "Business fire insurance"], notes: ["Include allied perils schedule with the proposal."], documentIds: ["doc-fire-proposal", "ref-fire-proposal-pdf"] },
    { id: "livestock-pack", title: "Livestock Insurance Pack", category: "Livestock", stage: "Proposal", description: "Proposal form for livestock and farm animal insurance.", requiredFor: ["Cattle insurance", "Poultry insurance", "Agri-finance collateral"], notes: ["Veterinary certificate required for each animal at proposal stage."], documentIds: ["doc-livestock-proposal"] },
    { id: "group-health-pack", title: "Group Health Pack", category: "Health", stage: "Enrollment", description: "Complete document set for group health insurance — financial proposal, member census, and claim forms.", requiredFor: ["Corporate group health insurance", "SME health plans", "Employee benefit schemes"], notes: ["Submit census within 7 days of proposal acceptance."], documentIds: ["doc-member-census", "doc-health-claim", "ref-group-financial-proposal"] },
    { id: "claims-pack", title: "Claims Documents Pack", category: "Claims", stage: "Claims", description: "Master claims document set covering all insurance lines.", requiredFor: ["Health claims", "Motor claims", "Fire claims"], notes: ["Original documents required. Certified copies accepted at claims officer discretion."], documentIds: ["ref-claims-docs-required"] },
  ],
};

export async function GET() {
  // Ensure seed file exists
  if (!fs.existsSync(SEED_FILE)) {
    fs.mkdirSync(path.dirname(SEED_FILE), { recursive: true });
    fs.writeFileSync(SEED_FILE, JSON.stringify(DEFAULT_SEED, null, 2), "utf8");
  }

  let seed = DEFAULT_SEED;
  try {
    seed = JSON.parse(fs.readFileSync(SEED_FILE, "utf8")) as typeof DEFAULT_SEED;
  } catch { /* fallback to default */ }

  // Enrich with any NEW template definitions that aren't in the seed yet
  if (fs.existsSync(DEFS_DIR)) {
    const defFiles = fs.readdirSync(DEFS_DIR).filter((f) => f.endsWith(".json"));
    const existingDefIds = new Set(seed.documents.map((d) => d.templateDefinitionId));
    for (const file of defFiles) {
      const defId = file.replace(".json", "");
      if (!existingDefIds.has(defId)) {
        // Auto-add user-created template as a new document
        try {
          const def = JSON.parse(fs.readFileSync(path.join(DEFS_DIR, file), "utf8")) as { id?: string; company?: { name?: string }; sections?: unknown[] };
          const label = defId.replace(/_/g, " ").replace(/\b\w/g, (c) => c.toUpperCase());
          seed.documents.push({
            id: `doc-custom-${defId}`,
            title: label,
            category: "Custom",
            stage: "Custom",
            kind: "custom",
            summary: `User-created template with ${(def.sections ?? []).length} sections.`,
            owner: def.company?.name ?? "Custom",
            uploadStatus: "Active",
            suggestedUse: "User-defined document template.",
            packId: "custom-pack",
            templateDefinitionId: defId,
            format: "docx",
            isGenerated: true,
          });
        } catch { /* skip invalid files */ }
      }
    }
  }

  return NextResponse.json({ ok: true, data: seed, source: "seed" });
}

export async function POST(request: Request) {
  // Allow saving updated seed (for future admin use)
  let body: typeof DEFAULT_SEED;
  try {
    body = (await request.json()) as typeof DEFAULT_SEED;
  } catch {
    return NextResponse.json({ ok: false, message: "Invalid JSON" }, { status: 400 });
  }
  fs.mkdirSync(path.dirname(SEED_FILE), { recursive: true });
  fs.writeFileSync(SEED_FILE, JSON.stringify(body, null, 2), "utf8");
  return NextResponse.json({ ok: true });
}
