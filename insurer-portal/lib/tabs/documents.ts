import {
  FileText,
  Image as ImageIcon,
  Presentation,
  Sheet,
  type LucideIcon,
} from "lucide-react";

import {
  getCategoryOptions,
  type DigitalDocumentBlock,
  type DocumentKind,
  type DocumentStage,
  type InsurerManagedDocument,
  type ReferenceDocument,
} from "@/lib/pragati-documents";

export type SheetDraft = {
  fields: Record<string, string>;
  tables: Record<string, string[][]>;
};

export type DocumentEditorState = {
  title: string;
  category: string;
  stage: DocumentStage;
  kind: DocumentKind;
  summary: string;
  owner: string;
  uploadStatus: string;
  suggestedUse: string;
  fileName: string;
};

export type DocumentOverrides = Record<string, Partial<DocumentEditorState>>;
export type ManagedDocument = InsurerManagedDocument | ReferenceDocument;
export type DocumentEditorMode = "none" | "create" | "edit";
export type DocumentViewMode = "preview" | "source";

export const documentDraftStorageKey = "insurer-portal-documents-gallery-drafts";
export const documentOverrideStorageKey = "insurer-portal-documents-gallery-overrides";
export const documentCustomStorageKey = "insurer-portal-documents-gallery-custom";
export const documentStageFilterOptions = ["All", "Proposal", "Pricing", "Claim Reimbursement", "Reference"] as const;
export const documentKindOptions: Array<{ value: DocumentKind; label: string }> = [
  { value: "proposal-form", label: "Proposal form" },
  { value: "medical-questionnaire", label: "Medical questionnaire" },
  { value: "declaration", label: "Declaration" },
  { value: "rate-table", label: "Rate table" },
  { value: "claim-form", label: "Claim form" },
  { value: "schedule", label: "Schedule" },
  { value: "reference-file", label: "Reference file" },
  { value: "process-note", label: "Process note" },
  { value: "presentation", label: "Presentation" },
  { value: "image-reference", label: "Image reference" },
];
export const documentChoiceOptions = [
  { value: "", label: "Select response" },
  { value: "Yes", label: "Yes" },
  { value: "No", label: "No" },
  { value: "Not applicable", label: "Not applicable" },
] as const;

export const documentsTabCopy = {
  hero: {
    eyebrow: "Insurer document library",
    title: "Source-aware cards, multi-page previews, and print-ready forms",
    description:
      "Each document is shown as a compact card. Opening it launches a modal with page-by-page layout, header, footer, and source-aware structure based on the original insurer files and workbook sheets.",
    addButton: "Add document",
  },
  filters: {
    searchPlaceholder: "Search document title, category, stage, owner, or source",
    emptyLabel: "No documents matched the current filters.",
  },
  libraryCard: {
    slideCountLabel: "slides",
    pageCountLabel: "pages",
  },
  digital: {
    criticalYes: "Yes",
    criticalNo: "No",
  },
  modal: {
    createEyebrow: "Library management",
    createTitle: "Add document",
    editButton: "Edit",
    savePdfButton: "Save as PDF",
    saveButton: "Save document",
    cancelButton: "Cancel edit",
    closeButton: "Close",
    infoLabel: "Document info",
    pagesLabel: "Pages",
    viewLabel: "View",
    previewButton: "Preview",
    sourceButton: "Source",
    sourceFallbackTitle: "Document source",
    officeDownloadLabel: "This source file opens as a downloadable office document.",
    openSourceButton: "Open source file",
  },
  editor: {
    titleLabel: "Title",
    categoryLabel: "Category",
    stageLabel: "Stage",
    kindLabel: "Kind",
    ownerLabel: "Owner",
    statusLabel: "Status",
    sourceFileLabel: "Source file name",
    summaryLabel: "Summary",
    suggestedUseLabel: "Suggested use",
  },
} as const;

export function createInitialDocumentDraft(blocks: DigitalDocumentBlock[]): SheetDraft {
  const fields: Record<string, string> = {};
  const tables: Record<string, string[][]> = {};

  blocks.forEach((block) => {
    if (block.type === "field-group") {
      block.fields.forEach((field) => {
        fields[field.id] = field.defaultValue ?? "";
      });
    }

    if (block.type === "table") {
      tables[block.id] = Array.from({ length: block.editableRows }, (_, rowIndex) => {
        const source = block.rows[rowIndex] ?? [];
        return block.headers.map((_, columnIndex) => source[columnIndex] ?? "");
      });
    }
  });

  return { fields, tables };
}

export function createDocumentEditorState(document?: InsurerManagedDocument | ReferenceDocument): DocumentEditorState {
  if (!document) {
    return {
      title: "",
      category: "Reference",
      stage: "Reference",
      kind: "reference-file",
      summary: "",
      owner: "Labaid Insuretech",
      uploadStatus: "Draft",
      suggestedUse: "",
      fileName: "",
    };
  }

  return {
    title: document.title,
    category: document.category,
    stage: document.stage,
    kind: document.kind,
    summary: document.summary,
    owner: document.owner,
    uploadStatus: document.uploadStatus,
    suggestedUse: document.suggestedUse,
    fileName: document.fileName,
  };
}

export function toCustomDocument(state: DocumentEditorState): ReferenceDocument {
  const id = `reference-library-${Date.now()}`;

  return {
    id,
    sourceType: "reference-file",
    title: state.title.trim() || "Untitled document",
    category: state.category.trim() || "Reference",
    stage: state.stage,
    kind: state.kind,
    summary: state.summary.trim() || "Custom document created in the insurer portal library.",
    owner: state.owner.trim() || "Labaid Insuretech",
    uploadStatus: state.uploadStatus.trim() || "Draft",
    suggestedUse: state.suggestedUse.trim() || "Use this document entry inside the insurer library.",
    packId: null,
    fileName: state.fileName.trim() || "manual-document.txt",
    sourceLabel: "Manual library entry",
    format: "manual",
    rows: [],
  };
}

export function mergeManagedDocument(
  document: InsurerManagedDocument | ReferenceDocument,
  overrides: DocumentOverrides,
) {
  const override = overrides[document.id];
  return override ? { ...document, ...override } : document;
}

export function getDocumentCategoryFilterOptions() {
  return ["All", ...getCategoryOptions()];
}

export function getManagedDocuments(documents: ManagedDocument[], overrides: DocumentOverrides) {
  return documents.map((document) => mergeManagedDocument(document, overrides));
}

export function filterManagedDocuments(
  documents: ManagedDocument[],
  query: string,
  stageFilter: (typeof documentStageFilterOptions)[number],
  categoryFilter: string,
) {
  const lowered = query.trim().toLowerCase();

  return documents.filter((document) => {
    if (stageFilter !== "All" && document.stage !== stageFilter) return false;
    if (categoryFilter !== "All" && document.category !== categoryFilter) return false;
    if (!lowered) return true;

    return [
      document.title,
      document.category,
      document.summary,
      document.owner,
      document.sourceLabel,
      document.suggestedUse,
    ].some((value) => value.toLowerCase().includes(lowered));
  });
}

export function getDocumentStageTone(stage: DocumentStage | "All") {
  if (stage === "Claim Reimbursement") return "pill-danger";
  if (stage === "Pricing") return "pill-warn";
  if (stage === "Reference") return "pill-neutral";
  return "pill-live";
}

export function getDocumentFormatIcon(document: ManagedDocument): LucideIcon {
  if (document.format === "pdf" || document.kind === "proposal-form" || document.kind === "claim-form") return FileText;
  if (document.format === "pptx" || document.kind === "presentation") return Presentation;
  if (document.format === "jpeg" || document.kind === "image-reference") return ImageIcon;
  return Sheet;
}

export function getDocumentPageCountLabel(variant: "paper" | "slides") {
  return variant === "slides"
    ? documentsTabCopy.libraryCard.slideCountLabel
    : documentsTabCopy.libraryCard.pageCountLabel;
}

export function getDocumentSourceTitle(document?: ManagedDocument | null) {
  return document?.title ?? documentsTabCopy.modal.sourceFallbackTitle;
}
