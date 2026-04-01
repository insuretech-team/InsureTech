import {
  buildDigitalBlocks,
  isDigitalTemplate,
  type DigitalDocumentBlock,
  type InsurerManagedDocument,
} from "@/lib/pragati-documents";

export type DocumentPreviewSection =
  | { type: "paragraph"; text: string }
  | { type: "bullet-list"; title?: string; items: string[] }
  | { type: "field-grid"; title?: string; columns?: number; items: Array<{ label: string; value?: string }> }
  | { type: "table"; title?: string; headers: string[]; rows: string[][] }
  | { type: "notice"; title?: string; text: string };

export type DocumentPreviewPage = {
  id: string;
  label: string;
  title: string;
  subtitle?: string;
  footer?: string;
  sections?: DocumentPreviewSection[];
  blocks?: DigitalDocumentBlock[];
  sourceUrl?: string;
};

export type DocumentPreviewBundle = {
  variant: "paper" | "slides";
  sourceUrl?: string;
  sourceLabel?: string;
  pages: DocumentPreviewPage[];
};

const sourceUrlByFileName: Record<string, string> = {
  "OMP Proposal Form (New).pdf": "/insurer-docs/omp-proposal-form-new.pdf",
  "Motor Insurance Proposal Form.pdf": "/insurer-docs/motor-insurance-proposal-form.pdf",
  "Financial Proposal LifePlus_Shanta-2026 (Group Insurance).pdf": "/insurer-docs/financial-proposal-lifeplus-shanta-2026.pdf",
  "Fire Insurance Proposal Form_20230622_0001.pdf": "/insurer-docs/fire-insurance-proposal-form-20230622-0001.pdf",
  "Documents are normally required for Claims.docx": "/insurer-docs/claims-required-documents.docx",
  "OMP New  Claim Process.docx": "/insurer-docs/omp-new-claim-process.docx",
  "Motor insurance policy-LifePlus Bangladesh(Final)(Underwrite & Claims).pptx": "/insurer-docs/motor-underwrite-claims-deck.pptx",
};

const referencePreviewById: Record<string, DocumentPreviewBundle> = {
  "reference-documents-are-normally-required-for-claims-docx": {
    variant: "paper",
    sourceUrl: "/insurer-docs/claims-required-documents.docx",
    sourceLabel: "Original claim checklist note",
    pages: [
      {
        id: "claims-checklist-fire",
        label: "Page 1",
        title: "Fire Claim Checklist",
        subtitle: "Documents normally required before insurer review",
        footer: "Claims reference / Fire",
        sections: [
          {
            type: "notice",
            title: "Context",
            text: "This note is a working operational checklist for fire-loss claims where Labaid Insuretech must assemble the insurer-ready pack before submission.",
          },
          {
            type: "bullet-list",
            title: "Required documents",
            items: [
              "Claim form signed by the insured office and concerned bank",
              "Original fire brigade report",
              "Daily stock report for the last 90 days countersigned by the bank",
              "Monthly stock statement countersigned by the bank",
              "Stock register copies for the last six months",
              "Tally book copies for the affected godown",
              "Purchase invoices, challans, and local bill copies for damaged stock",
              "Fire license and trade license",
              "Power-source permissions and electricity test reports where generators are used",
              "Generator logbook and GD entry copy",
              "Layout plan of the mill and godown",
            ],
          },
          {
            type: "bullet-list",
            title: "Statements to collect",
            items: [
              "Factory in-charge statement",
              "Electrical engineer statement",
              "Godown keeper statement",
              "Guard or duty personnel statement attested by the appropriate authority",
            ],
          },
        ],
      },
      {
        id: "claims-checklist-motor",
        label: "Page 2",
        title: "Motor Claim Checklist",
        subtitle: "Documents required before motor survey and settlement",
        footer: "Claims reference / Motor",
        sections: [
          {
            type: "bullet-list",
            title: "Required documents",
            items: [
              "Claim intimation letter",
              "Claim form signed and sealed by the insured and concerned bank",
              "GD entry copy signed by the concerned police station",
              "Three repair estimates from reputed workshops",
              "Driver statement countersigned by the insured",
              "Registration certificate copy",
              "Tax token, fitness certificate, and route permit copy",
              "Motor vehicle inspector report",
              "Driver license copy attested by BRTA",
              "Original challan of the carrying goods at the time of accident",
            ],
          },
          {
            type: "notice",
            title: "Operational note",
            text: "This checklist directly supports the surveyor and claims desk workflow. The pack should be complete before the motor claim advances to full insurer assessment.",
          },
        ],
      },
    ],
  },
  "reference-omp-new-claim-process-docx": {
    variant: "paper",
    sourceUrl: "/insurer-docs/omp-new-claim-process.docx",
    sourceLabel: "Original travel claim process note",
    pages: [
      {
        id: "omp-claim-process",
        label: "Page 1",
        title: "Travel Medical Emergency Process",
        subtitle: "Overseas mediclaim assistance and claim submission flow",
        footer: "Travel claims / Assistance workflow",
        sections: [
          {
            type: "notice",
            title: "Emergency contact instruction",
            text: "In illness or accident abroad that may lead to hospitalization or trip curtailment, the claimant should contact the TPA or assistance provider immediately.",
          },
          {
            type: "field-grid",
            title: "Primary assistance contacts",
            columns: 2,
            items: [
              { label: "Provider", value: "Crisis24" },
              { label: "Postal address", value: "2 London Bridge, London, SE1 9RA, UK" },
              { label: "Telephone", value: "+44 207 902 7131" },
              { label: "Email", value: "opsassist@crisis24.com / corporateteam@crisis24.com" },
              { label: "Alternate provider", value: "Van Ameyd" },
              { label: "Alternate contact", value: "Office G18, Bromley Old Town Hall, 30 Tweedy Road, Bromley" },
            ],
          },
          {
            type: "bullet-list",
            title: "Process reminders",
            items: [
              "Ask the claimant to provide a callback number to avoid reverse-charge call issues",
              "Crisis24 or Van Ameyd act as third-party administrators and assistance providers",
              "The claimant applies for a claim form through the assistance provider",
              "The completed claim form, insurance certificate, and relevant documentation are submitted back to the assistance provider",
            ],
          },
        ],
      },
    ],
  },
  "reference-financial-proposal-lifeplus-shanta-2026-group-insurance-pdf": {
    variant: "paper",
    sourceUrl: "/insurer-docs/financial-proposal-lifeplus-shanta-2026.pdf",
    sourceLabel: "Original financial proposal PDF",
    pages: [
      {
        id: "financial-proposal",
        label: "Page 1",
        title: "Example B2B Financial Proposal",
        subtitle: "Commercial quotation letter example for group life and medical insurance",
        footer: "Commercial proposal / Page 1 of 1",
        sections: [
          {
            type: "field-grid",
            title: "Letter header",
            columns: 2,
            items: [
              { label: "Reference", value: "SLI/HO/GID/2025-031" },
              { label: "Date", value: "4 September 2026" },
              { label: "Recipient", value: "CEO, Life Plus Bangladesh" },
              { label: "Address", value: "House-66, Mirpur Road-67, Kalabagan, 2nd Lane, Dhaka-1205" },
            ],
          },
          {
            type: "table",
            title: "Life coverage",
            headers: ["Benefit", "Members", "Sum Assured", "Rate", "Annual Premium"],
            rows: [
              ["Group Life (100K)", "3000", "30,00,00,000", "2.44 / 1000 SA", "BDT 7,32,000"],
              ["ADB (200K)", "3000", "Included", "1.10 / 1000 SA", "BDT 3,30,000"],
              ["Life & disability total", "-", "-", "3.54", "BDT 10,62,000"],
            ],
          },
          {
            type: "table",
            title: "Medical coverage",
            headers: ["Benefit", "Members", "Coverage per member", "Premium per member", "Annual Premium"],
            rows: [
              ["IPC Hospitalization 20K", "3000", "250,000", "410", "BDT 12,30,000"],
              ["OPC General 2K", "3000", "50,000", "550", "BDT 11,00,000"],
              ["Health total", "-", "-", "-", "BDT 23,30,000"],
              ["Combined annual premium", "-", "-", "-", "BDT 33,92,000"],
            ],
          },
          {
            type: "bullet-list",
            title: "Proposal notes",
            items: [
              "PPD and PTD schedule compliant with Bangladesh labor law first schedule",
              "Applicable VAT borne by the insurer side",
              "Premium payable yearly in advance",
              "Offer valid for two months",
            ],
          },
        ],
      },
    ],
  },
  "reference-motor-insurance-proposal-form-pdf": {
    variant: "paper",
    sourceUrl: "/insurer-docs/motor-insurance-proposal-form.pdf",
    sourceLabel: "Original motor proposal PDF",
    pages: [
      {
        id: "motor-proposal-page-1",
        label: "Page 1",
        title: "Private Vehicle Insurance Proposal Form",
        subtitle: "Pragati insurer proposal sheet for private motor underwriting",
        footer: "Motor proposal / Page 1 of 2",
        sections: [
          {
            type: "field-grid",
            title: "Header and proposer details",
            columns: 2,
            items: [
              { label: "Certificate no.", value: "Blank in source" },
              { label: "Policy no.", value: "Blank in source" },
              { label: "Issuing office", value: "Pragati Insurance, Head Office / issuing office block" },
              { label: "Proposer fields", value: "Full name, address, mobile, email, business or profession" },
            ],
          },
          {
            type: "field-grid",
            title: "Vehicle schedule",
            columns: 2,
            items: [
              { label: "Registration mark and no.", value: "Vehicle identity block" },
              { label: "Make / engine / chassis", value: "Vehicle particulars block" },
              { label: "Type of body / CC / HP", value: "Technical vehicle details" },
              { label: "Year of manufacture", value: "Model year field" },
              { label: "Carrying or seating capacity", value: "Capacity field" },
              { label: "Estimated value", value: "Insured value segregation" },
            ],
          },
          {
            type: "bullet-list",
            title: "Page-one declaration questions",
            items: [
              "Will the car be used solely for social, domestic, and pleasure purposes?",
              "Is the proposer the owner and is the vehicle registered in the proposer name?",
              "Date of purchase, whether new, price paid, and present market value",
              "Any driver physical infirmity affecting driving",
              "Any driver conviction or pending prosecution in the last five years",
              "Driving experience duration",
            ],
          },
        ],
      },
      {
        id: "motor-proposal-page-2",
        label: "Page 2",
        title: "Private Vehicle Insurance Proposal Form",
        subtitle: "Underwriter history, premium choice, and bank-use block",
        footer: "Motor proposal / Page 2 of 2",
        sections: [
          {
            type: "bullet-list",
            title: "Underwriter and claim history",
            items: [
              "Previous insurer name",
              "No claim bonus entitlement",
              "Prior proposal decline, cancellation, or renewal refusal",
              "Special conditions or increased premium history",
              "Accident and loss history for the past three years",
            ],
          },
          {
            type: "bullet-list",
            title: "Coverage and premium choices",
            items: [
              "Comprehensive policy / Act only / statutory minimum option",
              "First-loss deductible amount if the proposer wants to bear part of the loss",
              "Rugs, coats, and luggage extension",
              "Other additional benefits",
              "Policy commencement and proposer signature block",
            ],
          },
          {
            type: "field-grid",
            title: "Bank use only block",
            columns: 2,
            items: [
              { label: "Branch name", value: "Bank branch confirmation field" },
              { label: "Account number", value: "Lien or finance-linked account block" },
              { label: "Account name", value: "Borrower or vehicle owner block" },
              { label: "Bank officer code", value: "Officer verification field" },
            ],
          },
        ],
      },
    ],
  },
  "reference-fire-insurance-proposal-form-20230622-0001-pdf": {
    variant: "paper",
    sourceUrl: "/insurer-docs/fire-insurance-proposal-form-20230622-0001.pdf",
    sourceLabel: "Original fire proposal PDF",
    pages: [
      {
        id: "fire-proposal-page-1",
        label: "Page 1",
        title: "Proposal for Fire Insurance",
        subtitle: "Property risk schedule and sum insured declaration",
        footer: "Fire proposal / Page 1 of 2",
        sections: [
          {
            type: "field-grid",
            title: "Header",
            columns: 2,
            items: [
              { label: "Insurer", value: "Pragati Insurance Limited" },
              { label: "Office block", value: "Head office details and contact lines" },
              { label: "Proposer details", value: "Full name, address, trade or profession" },
              { label: "Term of insurance", value: "From / To" },
            ],
          },
          {
            type: "table",
            title: "Amount to be insured",
            headers: ["Diagram no.", "Building", "Machinery", "Furniture & effects", "Merchandise / stock", "Total"],
            rows: [
              ["No. 1", "", "", "", "", ""],
              ["No. 2", "", "", "", "", ""],
              ["No. 3", "", "", "", "", ""],
              ["Total", "", "", "", "", ""],
            ],
          },
          {
            type: "notice",
            title: "Source context",
            text: "The scanned PDF is low-text, but the matching Pragati workbook confirms that this page is the property schedule and amount-to-be-insured section.",
          },
        ],
      },
      {
        id: "fire-proposal-page-2",
        label: "Page 2",
        title: "Proposal for Fire Insurance",
        subtitle: "Location, construction, occupation, and declaration",
        footer: "Fire proposal / Page 2 of 2",
        sections: [
          {
            type: "field-grid",
            title: "Location and construction",
            columns: 2,
            items: [
              { label: "Building name / owner / plot / holding / street / town", value: "Location block" },
              { label: "Number of storeys", value: "Construction block" },
              { label: "Walls / roof / floors", value: "Construction materials" },
              { label: "Adjoining building / building within 50 feet", value: "Exposure details" },
            ],
          },
          {
            type: "bullet-list",
            title: "Operational sections",
            items: [
              "Occupation declaration",
              "Hazard and storage characteristics",
              "Previous insurance or loss history",
              "Declaration and proposer signature",
            ],
          },
        ],
      },
    ],
  },
  "reference-motor-insurance-policy-lifeplus-bangladesh-final-underwrite-claims-pptx": {
    variant: "slides",
    sourceUrl: "/insurer-docs/motor-underwrite-claims-deck.pptx",
    sourceLabel: "Original training deck",
    pages: [
      {
        id: "motor-slide-1",
        label: "Slide 1",
        title: "Training on Motor Insurance",
        subtitle: "Deck overview",
        sections: [{ type: "paragraph", text: "The presentation introduces motor insurance concepts, coverage types, underwriting logic, and claims handling." }],
      },
      {
        id: "motor-slide-2",
        label: "Slide 5",
        title: "Document Required for Comprehensive Motor Insurance",
        sections: [
          {
            type: "bullet-list",
            items: [
              "Motor insurance proposal form",
              "Vehicle purchase invoice or present market value support",
              "Registration copy and owner NID",
              "CC, seating capacity, manufacturer, year of manufacturing",
              "KYC form",
            ],
          },
        ],
      },
      {
        id: "motor-slide-3",
        label: "Slide 19",
        title: "Claim Loading for Private and Commercial Vehicle",
        sections: [{ type: "paragraph", text: "The deck contains underwriting and renewal loading guidance linked to own-damage history and no-claim bonus logic." }],
      },
      {
        id: "motor-slide-4",
        label: "Slide 32",
        title: "Claims Procedure (Documents)",
        sections: [
          {
            type: "bullet-list",
            items: [
              "Claim form signed by insured and bank",
              "Registration certificate, fitness certificate, tax token",
              "Driver statement and driver license",
              "GD entry or FIR where applicable",
              "Police final report in theft cases",
              "Original survey report and three repair estimates",
            ],
          },
        ],
      },
      {
        id: "motor-slide-5",
        label: "Slide 33",
        title: "Claims Settlement Way",
        sections: [{ type: "paragraph", text: "The closing section positions survey, documentation, and settlement decisioning as the final operational path for motor claims." }],
      },
    ],
  },
};

function chunkBlocks(blocks: DigitalDocumentBlock[]) {
  const pages: DigitalDocumentBlock[][] = [];
  let current: DigitalDocumentBlock[] = [];
  let score = 0;

  const blockWeight = (block: DigitalDocumentBlock) => {
    if (block.type === "table") return 3;
    if (block.type === "field-group") return Math.max(1, Math.ceil(block.fields.length / 2));
    return 1;
  };

  blocks.forEach((block) => {
    const weight = blockWeight(block);
    if (current.length && score + weight > 6) {
      pages.push(current);
      current = [];
      score = 0;
    }

    current.push(block);
    score += weight;
  });

  if (current.length) pages.push(current);
  return pages;
}

export function getDocumentSourceUrl(document: InsurerManagedDocument) {
  return sourceUrlByFileName[document.fileName];
}

export function buildPreviewBundle(document: InsurerManagedDocument): DocumentPreviewBundle {
  if (!isDigitalTemplate(document)) {
    return (
      referencePreviewById[document.id] ?? {
        variant: "paper" as const,
        sourceUrl: getDocumentSourceUrl(document),
        sourceLabel: "Original source file",
        pages: [
          {
            id: `${document.id}-profile`,
            label: "Page 1",
            title: document.title,
            subtitle: document.category,
            footer: `${document.owner} / ${document.format.toUpperCase()}`,
            sections: [
              { type: "paragraph", text: document.summary },
              { type: "notice", title: "Suggested use", text: document.suggestedUse },
            ],
          },
        ],
      }
    );
  }

  const blocks = buildDigitalBlocks(document);
  const blockPages = chunkBlocks(blocks);
  const pages: DocumentPreviewPage[] = blockPages.map((pageBlocks, index) => ({
    id: `${document.id}-page-${index + 1}`,
    label: `Page ${index + 1}`,
    title: document.title,
    subtitle: index === 0 ? document.sourceLabel : `${document.category} continuation`,
    footer: `${document.owner} / ${index + 1} of ${Math.max(1, blockPages.length)}`,
    blocks: pageBlocks,
  }));

  return {
    variant: "paper" as const,
    sourceLabel: "Workbook source",
    pages,
  };
}
