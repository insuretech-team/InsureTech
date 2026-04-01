import workbookData from "@/lib/docs-forms-workbooks.json";

type RawWorkbookRow = {
  row: number;
  cells: string[];
};

type RawWorkbookSheet = {
  id: string;
  sheetName: string;
  rows: RawWorkbookRow[];
};

type RawWorkbook = {
  id: string;
  fileName: string;
  displayName: string;
  sheets: RawWorkbookSheet[];
};

export type DocumentStage =
  | "Proposal Intake"
  | "Medical Review"
  | "Pricing"
  | "Claim Reimbursement"
  | "Enrollment"
  | "Reference";

export type DocumentKind =
  | "proposal-form"
  | "medical-questionnaire"
  | "declaration"
  | "rate-table"
  | "claim-form"
  | "schedule"
  | "reference-file"
  | "process-note"
  | "presentation"
  | "image-reference";

export interface PragatiWorkbookRow {
  row: number;
  cells: string[];
}

export interface ManagedDocumentBase {
  id: string;
  title: string;
  category: string;
  stage: DocumentStage;
  kind: DocumentKind;
  summary: string;
  owner: string;
  uploadStatus: string;
  suggestedUse: string;
  packId: string | null;
  sourceLabel: string;
  format: string;
}

export interface PragatiDocumentTemplate extends ManagedDocumentBase {
  sourceType: "digital-form";
  workbookId: string;
  workbookName: string;
  fileName: string;
  sheetName: string;
  rows: PragatiWorkbookRow[];
}

export interface ReferenceDocument extends ManagedDocumentBase {
  sourceType: "reference-file";
  fileName: string;
  rows: [];
}

export type InsurerManagedDocument = PragatiDocumentTemplate | ReferenceDocument;

export interface PragatiDocumentPack {
  id: string;
  title: string;
  category: string;
  stage: DocumentStage;
  description: string;
  requiredFor: string[];
  templateIds: string[];
  notes: string[];
}

export type DigitalDocumentBlock =
  | { type: "heading"; id: string; text: string }
  | { type: "note"; id: string; text: string }
  | { type: "field-group"; id: string; fields: DigitalField[] }
  | { type: "table"; id: string; headers: string[]; rows: string[][]; editableRows: number };

export interface DigitalField {
  id: string;
  label: string;
  control: "text" | "textarea" | "number" | "tel" | "email" | "date" | "choice";
  defaultValue?: string;
}

type DocumentTemplateMeta = Omit<
  PragatiDocumentTemplate,
  "id" | "rows" | "workbookId" | "workbookName" | "fileName" | "sheetName" | "sourceType" | "sourceLabel" | "format"
>;

const stageByCategory: Record<string, DocumentStage> = {
  Travel: "Proposal Intake",
  Auto: "Proposal Intake",
  Fire: "Proposal Intake",
  "Commercial Vehicle": "Proposal Intake",
  Livestock: "Proposal Intake",
  "Group Health": "Enrollment",
  Health: "Claim Reimbursement",
  "Group Life": "Pricing",
};

const workbookMetadata: Record<string, DocumentTemplateMeta> = {
  "pragati.xlsx::Sheet1": {
    title: "Overseas Mediclaim Proposal",
    category: "Travel",
    stage: "Proposal Intake",
    kind: "proposal-form",
    summary: "Primary overseas mediclaim proposal form for business and holiday travel cases.",
    owner: "Pragati Insurance",
    uploadStatus: "Uploaded by Pragati",
    suggestedUse: "Use when Labaid Insuretech prepares outbound travel mediclaim submissions for proposal intake.",
    packId: "travel-proposal-pack",
  },
  "pragati.xlsx::Sheet2": {
    title: "Mediclaim Medical History",
    category: "Travel",
    stage: "Medical Review",
    kind: "medical-questionnaire",
    summary: "Health declaration and physician details that must accompany the travel mediclaim proposal.",
    owner: "Pragati Insurance",
    uploadStatus: "Uploaded by Pragati",
    suggestedUse: "Capture proposer and spouse medical history before sending the travel proposal to the insurer.",
    packId: "travel-proposal-pack",
  },
  "pragati.xlsx::Sheet3": {
    title: "Mediclaim Declaration & Benefits",
    category: "Travel",
    stage: "Proposal Intake",
    kind: "declaration",
    summary: "Traveller declaration, acknowledgement, and benefits schedule for overseas mediclaim.",
    owner: "Pragati Insurance",
    uploadStatus: "Uploaded by Pragati",
    suggestedUse: "Use as the declaration page and benefit acknowledgement in the travel proposal packet.",
    packId: "travel-proposal-pack",
  },
  "pragati.xlsx::Sheet4": {
    title: "Non-Schengen Premium Matrix",
    category: "Travel",
    stage: "Pricing",
    kind: "rate-table",
    summary: "Plan A and Plan B rating table for non-Schengen overseas mediclaim journeys.",
    owner: "Pragati Insurance",
    uploadStatus: "Uploaded by Pragati",
    suggestedUse: "Use the rate sheet to price non-Schengen travel proposals before dispatching them.",
    packId: "travel-pricing-pack",
  },
  "pragati.xlsx::Sheet5": {
    title: "Travel Addendum & Employment/Study Rates",
    category: "Travel",
    stage: "Pricing",
    kind: "rate-table",
    summary: "Children exclusions, corporate frequent travel pricing, and employment or study premium references.",
    owner: "Pragati Insurance",
    uploadStatus: "Uploaded by Pragati",
    suggestedUse: "Reference during frequent travel, employment, and study travel proposal preparation.",
    packId: "travel-pricing-pack",
  },
  "pragati.xlsx::Sheet6": {
    title: "Schengen Premium Matrix",
    category: "Travel",
    stage: "Pricing",
    kind: "rate-table",
    summary: "Plan A and Plan B rating matrix for Schengen-country overseas mediclaim cases.",
    owner: "Pragati Insurance",
    uploadStatus: "Uploaded by Pragati",
    suggestedUse: "Use for Schengen travel premium calculation and document preparation.",
    packId: "travel-pricing-pack",
  },
  "pragati.xlsx::Sheet7": {
    title: "Schengen Frequent Travel Addendum",
    category: "Travel",
    stage: "Pricing",
    kind: "rate-table",
    summary: "Schengen-specific frequent travel, employment, and studies addendum with premium guidance.",
    owner: "Pragati Insurance",
    uploadStatus: "Uploaded by Pragati",
    suggestedUse: "Keep attached for annual, repeated, or longer-term Schengen travel submissions.",
    packId: "travel-pricing-pack",
  },
  "pragati.xlsx::Sheet8": {
    title: "Private Vehicle Proposal",
    category: "Auto",
    stage: "Proposal Intake",
    kind: "proposal-form",
    summary: "Private vehicle proposal with proposer identity, vehicle particulars, valuation, and usage declarations.",
    owner: "Pragati Insurance",
    uploadStatus: "Uploaded by Pragati",
    suggestedUse: "Complete this digital form when sending private motor proposals to Pragati.",
    packId: "private-vehicle-pack",
  },
  "pragati.xlsx::Sheet9": {
    title: "Private Vehicle Declaration Continuation",
    category: "Auto",
    stage: "Proposal Intake",
    kind: "declaration",
    summary: "Continuation page for underwriter history, no-claim bonus, prior losses, and bank-use details.",
    owner: "Pragati Insurance",
    uploadStatus: "Uploaded by Pragati",
    suggestedUse: "Attach to the private vehicle proposal whenever declarations and bank metadata are needed.",
    packId: "private-vehicle-pack",
  },
  "pragati.xlsx::Sheet10": {
    title: "Fire Insurance Proposal",
    category: "Fire",
    stage: "Proposal Intake",
    kind: "proposal-form",
    summary: "Risk location, construction, stock values, and fire protection disclosures for fire proposals.",
    owner: "Pragati Insurance",
    uploadStatus: "Uploaded by Pragati",
    suggestedUse: "Use for property and fire-risk proposal submissions from Labaid Insuretech to Pragati.",
    packId: "fire-risk-pack",
  },
  "pragati.xlsx::Sheet11": {
    title: "Commercial Vehicle Proposal",
    category: "Commercial Vehicle",
    stage: "Proposal Intake",
    kind: "proposal-form",
    summary: "Commercial vehicle proposal with permit details, carrying capacity, valuation, and claims history.",
    owner: "Pragati Insurance",
    uploadStatus: "Uploaded by Pragati",
    suggestedUse: "Use when commercial fleet or logistics vehicle proposals are sent for underwriting.",
    packId: "commercial-vehicle-pack",
  },
  "pragati.xlsx::Sheet12": {
    title: "Livestock Proposal",
    category: "Livestock",
    stage: "Proposal Intake",
    kind: "proposal-form",
    summary: "Primary livestock proposal with insured details, farm location, and insured animal schedule.",
    owner: "Pragati Insurance",
    uploadStatus: "Uploaded by Pragati",
    suggestedUse: "Use for cattle and livestock submissions requiring farm and animal-level schedule capture.",
    packId: "livestock-pack",
  },
  "pragati.xlsx::Sheet13": {
    title: "Member Census Schedule",
    category: "Group Health",
    stage: "Enrollment",
    kind: "schedule",
    summary: "Family and dependent census schedule for employee, spouse, and child enrollment details.",
    owner: "Pragati Insurance",
    uploadStatus: "Uploaded by Pragati",
    suggestedUse: "Use as the member census sheet when Pragati requests grouped enrolment details.",
    packId: "member-census-pack",
  },
  "pragati.xlsx::Sheet14": {
    title: "Health Insurance Claim Form",
    category: "Health",
    stage: "Claim Reimbursement",
    kind: "claim-form",
    summary: "Reimbursement claim form with hospitalization expense breakup and required document checklist.",
    owner: "Pragati Insurance",
    uploadStatus: "Uploaded by Pragati",
    suggestedUse: "Use for reimbursement claims, claim pack validation, and coordinator forwarding to the insurer.",
    packId: "health-claim-pack",
  },
};

const referenceDocumentsMeta: Omit<ReferenceDocument, "id" | "rows" | "sourceType" | "sourceLabel">[] = [
  {
    title: "Claims Required Documents Note",
    category: "Claims",
    stage: "Claim Reimbursement",
    kind: "process-note",
    summary: "Reference note describing the documents normally required during claims submission and validation.",
    owner: "Labaid Insuretech",
    uploadStatus: "Imported from docs_forms",
    suggestedUse: "Use as an operational checklist while requesting or validating claim support documents.",
    packId: "claims-reference-pack",
    fileName: "Documents are normally required for Claims.docx",
    format: "docx",
  },
  {
    title: "OMP New Claim Process",
    category: "Travel",
    stage: "Claim Reimbursement",
    kind: "process-note",
    summary: "Operational process note for overseas mediclaim claims handling and forwarding.",
    owner: "Labaid Insuretech",
    uploadStatus: "Imported from docs_forms",
    suggestedUse: "Use as a process reference when coordinating travel-claim reimbursement cases.",
    packId: "travel-claims-pack",
    fileName: "OMP New  Claim Process.docx",
    format: "docx",
  },
  {
    title: "OMP Proposal Form Source PDF",
    category: "Travel",
    stage: "Proposal Intake",
    kind: "reference-file",
    summary: "Source PDF copy of the overseas mediclaim proposal form kept as a reference against the digital version.",
    owner: "Pragati Insurance",
    uploadStatus: "Imported from docs_forms",
    suggestedUse: "Use as the original insurer-supplied source when validating the digital travel proposal layout.",
    packId: "travel-proposal-pack",
    fileName: "OMP Proposal Form (New).pdf",
    format: "pdf",
  },
  {
    title: "Motor Insurance Proposal Source PDF",
    category: "Auto",
    stage: "Proposal Intake",
    kind: "reference-file",
    summary: "Source PDF of the motor insurance proposal form used as the original insurer reference.",
    owner: "Pragati Insurance",
    uploadStatus: "Imported from docs_forms",
    suggestedUse: "Use while validating the digital motor proposal against the original source document.",
    packId: "private-vehicle-pack",
    fileName: "Motor Insurance Proposal Form.pdf",
    format: "pdf",
  },
  {
    title: "Fire Insurance Proposal Source PDF",
    category: "Fire",
    stage: "Proposal Intake",
    kind: "reference-file",
    summary: "Source PDF copy of the fire insurance proposal form for property and stock risks.",
    owner: "Pragati Insurance",
    uploadStatus: "Imported from docs_forms",
    suggestedUse: "Use as the original source document while preparing or checking the digital fire form.",
    packId: "fire-risk-pack",
    fileName: "Fire Insurance Proposal Form_20230622_0001.pdf",
    format: "pdf",
  },
  {
    title: "Example B2B Financial Proposal",
    category: "Group Life",
    stage: "Pricing",
    kind: "reference-file",
    summary: "Financial proposal reference for group insurance commercial structuring and pricing comparison.",
    owner: "Labaid Insuretech",
    uploadStatus: "Imported from docs_forms",
    suggestedUse: "Use as a commercial reference when working on group-life or group-health insurer pricing discussions.",
    packId: "coverage-income-pack",
    fileName: "Financial Proposal LifePlus_Shanta-2026 (Group Insurance).pdf",
    format: "pdf",
  },
  {
    title: "Motor Underwrite & Claims Deck",
    category: "Auto",
    stage: "Reference",
    kind: "presentation",
    summary: "Working presentation covering motor underwriting and claims logic for LifePlus Bangladesh.",
    owner: "Labaid Insuretech",
    uploadStatus: "Imported from docs_forms",
    suggestedUse: "Use as an underwriting and claims reference for non-life motor operations and insurer discussion.",
    packId: "motor-claims-guide-pack",
    fileName: "Motor insurance policy-LifePlus Bangladesh(Final)(Underwrite & Claims).pptx",
    format: "pptx",
  },
];

const rawWorkbooks = workbookData as RawWorkbook[];
const excludedWorkbookSheets = new Set(["Insurance Coverage & Income.xlsx::Sheet2"]);

function defaultTemplateMeta(workbook: RawWorkbook, sheet: RawWorkbookSheet): DocumentTemplateMeta {
  return {
    title: `${workbook.displayName} - ${sheet.sheetName}`,
    category: "Reference",
    stage: "Reference",
    kind: "reference-file",
    summary: "Imported workbook sheet from docs_forms.",
    owner: "Labaid Insuretech",
    uploadStatus: "Imported from docs_forms",
    suggestedUse: "Review and map this sheet before using it in insurer operations.",
    packId: null,
  };
}

function compactLabel(value: string) {
  return normalizeText(value)
    .replace(/[._]{2,}/g, " ")
    .replace(/\s+/g, " ")
    .trim();
}

function slugify(value: string) {
  return compactLabel(value)
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");
}

export const pragatiDocumentTemplates: PragatiDocumentTemplate[] = rawWorkbooks.flatMap((workbook) =>
  workbook.sheets.flatMap((sheet) => {
    const key = `${workbook.fileName}::${sheet.sheetName}`;
    if (excludedWorkbookSheets.has(key)) return [];
    const meta = workbookMetadata[key] ?? defaultTemplateMeta(workbook, sheet);

    return [{
      id: sheet.id,
      workbookId: workbook.id,
      workbookName: workbook.displayName,
      fileName: workbook.fileName,
      sheetName: sheet.sheetName,
      rows: sheet.rows,
      sourceType: "digital-form",
      sourceLabel: `${workbook.displayName} / ${sheet.sheetName}`,
      format: "xlsx",
      ...meta,
    }];
  }),
);

export const referenceDocuments: ReferenceDocument[] = referenceDocumentsMeta.map((item) => ({
  id: `reference-${slugify(item.fileName)}`,
  rows: [],
  sourceType: "reference-file",
  sourceLabel: item.fileName,
  ...item,
}));

export const insurerManagedDocuments: InsurerManagedDocument[] = [...pragatiDocumentTemplates, ...referenceDocuments];

export const pragatiDocumentPacks: PragatiDocumentPack[] = [
  {
    id: "travel-proposal-pack",
    title: "Travel Proposal Pack",
    category: "Travel",
    stage: "Proposal Intake",
    description: "Proposal, medical history, declaration, and source reference packet for overseas mediclaim submissions.",
    requiredFor: ["Overseas mediclaim", "Business travel", "Holiday travel", "Schengen travel"],
    templateIds: [
      "sheet-pragati-sheet1",
      "sheet-pragati-sheet2",
      "sheet-pragati-sheet3",
      "reference-omp-proposal-form-new-pdf",
    ],
    notes: [
      "Passport and itinerary should be captured with the proposal before dispatch.",
      "Medical history and declaration pages must be completed by proposer and spouse where relevant.",
    ],
  },
  {
    id: "travel-pricing-pack",
    title: "Travel Pricing Reference Pack",
    category: "Travel",
    stage: "Pricing",
    description: "Premium matrices and addenda used to calculate travel proposal rates before insurer submission.",
    requiredFor: ["Travel pricing", "Schengen pricing", "Corporate frequent travel", "Employment or study cover"],
    templateIds: ["sheet-pragati-sheet4", "sheet-pragati-sheet5", "sheet-pragati-sheet6", "sheet-pragati-sheet7"],
    notes: [
      "Use as a reference pack while the proposal team prepares the primary travel proposal sheets.",
      "Schengen and non-Schengen ratings are kept separate to avoid pricing mismatch.",
    ],
  },
  {
    id: "private-vehicle-pack",
    title: "Private Vehicle Pack",
    category: "Auto",
    stage: "Proposal Intake",
    description: "Proposal, declaration continuation, and source PDF for private motor submissions.",
    requiredFor: ["Private vehicle", "Retail auto", "Personal car insurance"],
    templateIds: ["sheet-pragati-sheet8", "sheet-pragati-sheet9", "reference-motor-insurance-proposal-form-pdf"],
    notes: [
      "Collect valuation, usage, prior insurance, and no-claim disclosures before underwriter review.",
      "Bank-use metadata is captured on the continuation sheet and should travel with the proposal packet.",
    ],
  },
  {
    id: "fire-risk-pack",
    title: "Fire Risk Pack",
    category: "Fire",
    stage: "Proposal Intake",
    description: "Property and fire risk disclosure form backed by the original PDF source document.",
    requiredFor: ["Fire insurance", "Property cover", "Stock and contents cover"],
    templateIds: ["sheet-pragati-sheet10", "reference-fire-insurance-proposal-form-20230622-0001-pdf"],
    notes: ["The proposer must disclose building construction, stock values, and fire-fighting equipment readiness."],
  },
  {
    id: "commercial-vehicle-pack",
    title: "Commercial Vehicle Pack",
    category: "Commercial Vehicle",
    stage: "Proposal Intake",
    description: "Commercial vehicle proposal pack covering valuation, permit, carrying capacity, and claims history.",
    requiredFor: ["Commercial vehicle", "Fleet", "Passenger carrier", "Goods carrier"],
    templateIds: ["sheet-pragati-sheet11"],
    notes: ["Permit and carrying-capacity details need to be validated before the proposal moves to underwriting."],
  },
  {
    id: "livestock-pack",
    title: "Livestock Proposal Pack",
    category: "Livestock",
    stage: "Proposal Intake",
    description: "Farm-level proposal and animal schedule for cattle and livestock insurance submissions.",
    requiredFor: ["Livestock", "Cattle", "Farm insurance"],
    templateIds: ["sheet-pragati-sheet12"],
    notes: ["Each insured animal should be identified with tag, species, value, and purchase or birth date."],
  },
  {
    id: "member-census-pack",
    title: "Member Census Pack",
    category: "Group Health",
    stage: "Enrollment",
    description: "Employee and dependent roster set for group and family enrollment operations.",
    requiredFor: ["Group health", "Employee census", "Family enrollment"],
    templateIds: ["sheet-pragati-sheet13"],
    notes: ["Use this pack when an insurer asks for employee, spouse, child, or nominee roster data."],
  },
  {
    id: "health-claim-pack",
    title: "Health Reimbursement Claim Pack",
    category: "Health",
    stage: "Claim Reimbursement",
    description: "Claim reimbursement form and checklist for hospitalization expense submissions.",
    requiredFor: ["Health reimbursement", "Hospital claim", "Document validation"],
    templateIds: ["sheet-pragati-sheet14"],
    notes: [
      "Original bills, prescriptions, discharge papers, and diagnostics should be checked before forwarding.",
      "The plan coordinator forwarding block should be completed by Labaid Insuretech operations.",
    ],
  },
  {
    id: "coverage-income-pack",
    title: "Coverage & Income Pack",
    category: "Group Life",
    stage: "Pricing",
    description: "Coverage, income, and group commercial worksheets used during proposal pricing and negotiation.",
    requiredFor: ["Group insurance pricing", "Commercial review", "Income modeling"],
    templateIds: [
      "reference-financial-proposal-lifeplus-shanta-2026-group-insurance-pdf",
    ],
    notes: ["Use these files together when validating group premium assumptions and portfolio economics."],
  },
  {
    id: "travel-claims-pack",
    title: "Travel Claims Reference Pack",
    category: "Travel",
    stage: "Claim Reimbursement",
    description: "Travel-claim process notes and insurer-facing references for overseas mediclaim cases.",
    requiredFor: ["Travel reimbursement", "OMP claim process", "Coordinator handoff"],
    templateIds: ["reference-omp-new-claim-process-docx"],
    notes: ["Use this process pack when the claim team needs the original travel-claims handling note."],
  },
  {
    id: "claims-reference-pack",
    title: "Claims Reference Pack",
    category: "Claims",
    stage: "Claim Reimbursement",
    description: "General claims-supporting references, checklist notes, and field capture material.",
    requiredFor: ["Claim checklist", "Document deficiency review", "Operational reference"],
    templateIds: ["reference-documents-are-normally-required-for-claims-docx"],
    notes: ["Use this pack to guide document requests and manual claim completeness checks."],
  },
  {
    id: "motor-claims-guide-pack",
    title: "Motor Underwriting & Claims Guide Pack",
    category: "Auto",
    stage: "Reference",
    description: "Reference material for motor underwriting and claims coordination.",
    requiredFor: ["Motor claims review", "Underwriting discussion", "Training reference"],
    templateIds: ["reference-motor-insurance-policy-lifeplus-bangladesh-final-underwrite-claims-pptx"],
    notes: ["This pack is reference-only and should support insurer and operations discussion."],
  },
];

export function isDigitalTemplate(document: InsurerManagedDocument): document is PragatiDocumentTemplate {
  return document.sourceType === "digital-form";
}

export function normalizeText(value: string) {
  return value.replace(/\s+/g, " ").trim();
}

function makeId(prefix: string, value: string, index: number) {
  return `${prefix}-${slugify(value).slice(0, 56) || "field"}-${index}`;
}

function looksLikeHeading(value: string) {
  const cleaned = normalizeText(value);
  const letters = cleaned.replace(/[^A-Za-z]/g, "");
  if (!letters) return false;
  if (cleaned.length < 96 && cleaned === cleaned.toUpperCase()) return true;
  return /(proposal form|premium rating structure|medical history|product benefits|claim form|bank use only|head office)/i.test(
    cleaned,
  );
}

function looksLikeQuestion(value: string) {
  return /(\?|^are\b|^do\b|^will\b|^has\b|^have\b|^is\b|^how\b|^what\b|^were\b|^please give|^give particulars|^if)/i.test(
    normalizeText(value),
  );
}

function looksLikeField(value: string) {
  const cleaned = normalizeText(value);
  return (
    looksLikeQuestion(cleaned) ||
    cleaned.endsWith(":") ||
    /signature|policy no|certificate no|address|name of|occupation|membership no|passport number|date of|period of|itinerary|premium|amount|value|remarks|telephone|mobile|email|coverage|sum insured|designation|gender|age|relation/i.test(
      cleaned,
    ) ||
    /\d+[.)]/.test(cleaned.slice(0, 6))
  );
}

function looksLikeValue(value: string) {
  const cleaned = normalizeText(value);
  if (!cleaned) return false;
  if (looksLikeHeading(cleaned) || looksLikeField(cleaned) || looksLikeQuestion(cleaned)) return false;
  return cleaned.length < 96;
}

function controlFromLabel(label: string): DigitalField["control"] {
  const cleaned = label.toLowerCase();
  if (/yes or no|will|are you|have you|do you|is the/i.test(cleaned)) return "choice";
  if (/email/i.test(cleaned)) return "email";
  if (/telephone|mobile|phone/i.test(cleaned)) return "tel";
  if (/date|period of insurance|departure|discharge|admission|coverage start|coverage end/i.test(cleaned)) return "date";
  if (/amount|premium|value|weight|number of days|age|sum insured|horse power|capacity|policy|gross premium/i.test(cleaned))
    return "number";
  if (/address|particulars|history|itinerary|other information|remarks|description/i.test(cleaned) || label.length > 90)
    return "textarea";
  return "text";
}

function createField(label: string, index: number, initialValue?: string): DigitalField {
  return {
    id: makeId("field", label, index),
    label: compactLabel(label),
    control: controlFromLabel(label),
    defaultValue: initialValue,
  };
}

function rowStartsNewFieldGroup(cells: string[]) {
  return cells.length < 4;
}

export function buildDigitalBlocks(template: PragatiDocumentTemplate): DigitalDocumentBlock[] {
  const blocks: DigitalDocumentBlock[] = [];
  const rows = template.rows;
  let index = 0;
  let pendingFields: DigitalField[] = [];

  function flushPendingFields() {
    if (!pendingFields.length) return;

    blocks.push({
      type: "field-group",
      id: makeId("field-group", `${template.id}-${blocks.length}`, blocks.length),
      fields: pendingFields,
    });
    pendingFields = [];
  }

  while (index < rows.length) {
    const current = rows[index];
    const cells = current.cells.map(compactLabel).filter(Boolean);

    if (!cells.length) {
      index += 1;
      continue;
    }

    if (cells.length >= 4) {
      flushPendingFields();

      const headerRow = cells;
      const tableRows: string[][] = [];
      let pointer = index + 1;

      while (pointer < rows.length) {
        const nextCells = rows[pointer].cells.map(compactLabel).filter(Boolean);
        if (!nextCells.length) {
          pointer += 1;
          continue;
        }
        if (nextCells.length === 1 && (looksLikeHeading(nextCells[0]) || looksLikeField(nextCells[0]))) break;
        if (rowStartsNewFieldGroup(nextCells) && looksLikeField(nextCells[0])) break;

        tableRows.push(nextCells);
        pointer += 1;
      }

      blocks.push({
        type: "table",
        id: makeId("table", `${template.id}-${current.row}`, blocks.length),
        headers: headerRow,
        rows: tableRows,
        editableRows: Math.max(3, tableRows.length || 0),
      });
      index = pointer;
      continue;
    }

    if (cells.length === 1) {
      const text = cells[0];

      if (looksLikeHeading(text)) {
        flushPendingFields();
        blocks.push({ type: "heading", id: makeId("heading", text, blocks.length), text });
      } else if (looksLikeField(text)) {
        if (/signature/i.test(text) && /date/i.test(text)) {
          pendingFields.push(createField("Signature", blocks.length), createField("Date", blocks.length + 1));
        } else {
          pendingFields.push(createField(text, blocks.length));
        }
        if (pendingFields.length >= 3) flushPendingFields();
      } else {
        flushPendingFields();
        blocks.push({
          type: "note",
          id: makeId("note", text, blocks.length),
          text,
        });
      }

      index += 1;
      continue;
    }

    if (cells.length === 2 && looksLikeField(cells[0]) && looksLikeValue(cells[1])) {
      pendingFields.push(createField(cells[0], blocks.length, cells[1]));
      if (pendingFields.length >= 3) flushPendingFields();
      index += 1;
      continue;
    }

    flushPendingFields();
    blocks.push({
      type: "field-group",
      id: makeId("field-group", `${template.id}-${current.row}`, blocks.length),
      fields: cells.map((cell, cellIndex) => createField(cell, blocks.length + cellIndex)),
    });
    index += 1;
  }

  flushPendingFields();

  return blocks;
}

export function findDocumentPack(packId?: string) {
  return pragatiDocumentPacks.find((pack) => pack.id === packId) ?? null;
}

export function findDocumentPackForCategory(category?: string, planName?: string) {
  const haystack = `${category ?? ""} ${planName ?? ""}`.toLowerCase();

  if (/(travel|mediclaim|overseas|schengen)/i.test(haystack)) return findDocumentPack("travel-proposal-pack");
  if (/(private|motor|auto|car)/i.test(haystack)) return findDocumentPack("private-vehicle-pack");
  if (/(commercial vehicle|fleet|truck|bus|carrier)/i.test(haystack)) return findDocumentPack("commercial-vehicle-pack");
  if (/(fire|property)/i.test(haystack)) return findDocumentPack("fire-risk-pack");
  if (/(livestock|cattle|farm|pet)/i.test(haystack)) return findDocumentPack("livestock-pack");
  if (/(group|family|employee|member|enrollment)/i.test(haystack)) return findDocumentPack("member-census-pack");
  if (/(health|hospital|claim|medical)/i.test(haystack)) return findDocumentPack("health-claim-pack");
  if (/(life|income|coverage|financial proposal)/i.test(haystack)) return findDocumentPack("coverage-income-pack");

  return null;
}

export function getTemplatesForPack(packId?: string) {
  if (!packId) return [];
  const pack = findDocumentPack(packId);
  if (!pack) return [];
  return insurerManagedDocuments.filter((document) => pack.templateIds.includes(document.id));
}

export function getDocumentStageOptions() {
  return Array.from(new Set(insurerManagedDocuments.map((document) => document.stage))).sort();
}

export function getCategoryOptions() {
  return Array.from(new Set(insurerManagedDocuments.map((document) => document.category))).sort();
}

export function getDefaultStageForCategory(category: string): DocumentStage {
  return stageByCategory[category] ?? "Reference";
}
