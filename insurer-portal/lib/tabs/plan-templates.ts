import { productPlaybooks } from "@/lib/product-playbooks";

export const planTemplatesStorageKey = "insurer-portal.plan-template-studio";

export type TemplateBlock = {
  id: string;
  kind: "coverage" | "pricing" | "eligibility" | "documents" | "terms";
  title: string;
  body: string;
};

export type PlanTemplate = {
  id: string;
  productId: string;
  title: string;
  summary: string;
  blocks: TemplateBlock[];
};

export function createTemplateBlock(kind: TemplateBlock["kind"] = "terms"): TemplateBlock {
  return {
    id: `blk-${Date.now()}`,
    kind,
    title: "New section",
    body: "Add content for this section.",
  };
}

export function createEmptyPlanTemplate(productId = ""): PlanTemplate {
  return {
    id: `tpl-${Date.now()}`,
    productId,
    title: "New plan template",
    summary: "Draft a concise storefront and operations template for this plan.",
    blocks: [
      {
        ...createTemplateBlock("coverage"),
        title: "Coverage overview",
        body: "Summarize the most important benefits and coverage boundaries here.",
      },
    ],
  };
}

export function createPlanTemplateFromPlaybook(playbook: (typeof productPlaybooks)[number], productId = ""): PlanTemplate {
  return {
    id: `tpl-${playbook.code.toLowerCase()}-${Date.now()}`,
    productId,
    title: `${playbook.name} Template`,
    summary: playbook.summary,
    blocks: [
      {
        id: `blk-${playbook.code}-overview`,
        kind: "coverage",
        title: "Coverage overview",
        body: `${playbook.coverageLimitText}. ${playbook.summary}`,
      },
      {
        id: `blk-${playbook.code}-eligibility`,
        kind: "eligibility",
        title: "Eligibility",
        body: `Audience: ${playbook.audience}. Age range: ${playbook.ageRange}. Policy term: ${playbook.policyTerm}.`,
      },
      {
        id: `blk-${playbook.code}-documents`,
        kind: "documents",
        title: "Required documents",
        body: playbook.requiredDocuments.join("; "),
      },
      {
        id: `blk-${playbook.code}-terms`,
        kind: "terms",
        title: "Key exclusions and flags",
        body: [...playbook.exclusions.slice(0, 3), ...playbook.operationalFlags].join("; "),
      },
    ],
  };
}

export const planTemplateKindOptions = [
  { value: "coverage", label: "Coverage" },
  { value: "pricing", label: "Pricing" },
  { value: "eligibility", label: "Eligibility" },
  { value: "documents", label: "Documents" },
  { value: "terms", label: "Terms" },
] as const;

export const planTemplatesTabCopy = {
  libraryPanel: {
    title: "Template library",
    description: "Draft reusable plan content for insurer products.",
    newButton: "New template",
    emptyProduct: "Unassigned product",
  },
  studioPanel: {
    title: "Template studio",
    description: "Shape plan structure, sections, and positioning text before backend persistence is introduced.",
    titleLabel: "Template title",
    productLabel: "Product",
    emptyProductOption: "Select a product",
    summaryLabel: "Summary",
    sectionTypeLabel: "Section type",
    sectionTitleLabel: "Section title",
    sectionBodyLabel: "Section body",
    removeBlockButton: "Remove",
    addSectionButton: "Add section",
    saveButton: "Save local templates",
    deleteButton: "Delete template",
    emptyLabel: "Create or select a template to begin editing.",
    saveMessage: "Template drafts saved locally in this browser.",
  },
} as const;
