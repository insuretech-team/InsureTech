export interface EnrollmentMember {
  employeeId: string;
  memberName: string;
  relation: string;
  designation: string;
  sumAssuredText: string;
  coverageStart: string;
  nominee: string;
  phone: string;
}

export interface EnrollmentBatch {
  id: string;
  title: string;
  clientName: string;
  proposalNumber: string;
  status: string;
  memberCount: number;
  dependentCount: number;
  coverageWindow: string;
  sourceFiles: string[];
  notes: string[];
  validationFlags: string[];
  rows: EnrollmentMember[];
}

export interface PricingLineItem {
  label: string;
  members: string;
  benefit: string;
  rateText: string;
  annualPremiumText: string;
}

export interface PricingScenario {
  id: string;
  title: string;
  clientName: string;
  status: string;
  segment: string;
  headlinePremiumText: string;
  sourceFiles: string[];
  assumptions: string[];
  openItems: string[];
  deliverables: string[];
  lineItems: PricingLineItem[];
}

export interface ClaimChecklistDocument {
  name: string;
  owner: string;
  purpose: string;
  critical: boolean;
}

export interface ClaimChecklist {
  id: string;
  title: string;
  category: string;
  surveyorRequired: boolean;
  responseWindow: string;
  sourceFiles: string[];
  intakeSteps: string[];
  deficiencySignals: string[];
  documents: ClaimChecklistDocument[];
}

export interface AssistanceProvider {
  id: string;
  name: string;
  role: string;
  supportWindow: string;
  phones: string[];
  emails: string[];
  address: string;
  useCases: string[];
  handoffPacket: string[];
  sourceFiles: string[];
}

export interface AssistanceRunbook {
  title: string;
  description: string;
  steps: string[];
}

export interface KnowledgeAsset {
  id: string;
  title: string;
  category: string;
  audience: string;
  summary: string;
  sourceFile: string;
  linkedTabs: string[];
  keyPoints: string[];
}

export const enrollmentBatches: EnrollmentBatch[] = [
  {
    id: "batch-pragati-member-census",
    title: "Pragati Member Census Schedule",
    clientName: "Labaid SME Family Float",
    proposalNumber: "PRG-ENR-0326-001",
    status: "Ready for insurer dispatch",
    memberCount: 142,
    dependentCount: 117,
    coverageWindow: "1 Apr 2026 to 31 Mar 2027",
    sourceFiles: ["pragati.xlsx / Sheet13", "Enrollment Format  Alpha Force.xlsx"],
    notes: [
      "Structured around the insurer-ready member census sheet used for employee, spouse, and child enrollment.",
      "Labaid operations should validate family composition and nominee rows before forwarding to Pragati.",
    ],
    validationFlags: [
      "7 members missing nominee relation",
      "3 children missing date of birth normalization",
      "2 records need coverage effective-date correction",
    ],
    rows: [
      {
        employeeId: "AF-1001",
        memberName: "Sadia Rahman",
        relation: "Employee",
        designation: "Accounts Executive",
        sumAssuredText: "BDT 250,000",
        coverageStart: "1 Apr 2026",
        nominee: "Tariq Hasan",
        phone: "01711-220011",
      },
      {
        employeeId: "AF-1001-S",
        memberName: "Tariq Hasan",
        relation: "Spouse",
        designation: "Dependent",
        sumAssuredText: "BDT 250,000",
        coverageStart: "1 Apr 2026",
        nominee: "Sadia Rahman",
        phone: "01711-220011",
      },
      {
        employeeId: "AF-1088",
        memberName: "Nusrat Jahan",
        relation: "Employee",
        designation: "HR Lead",
        sumAssuredText: "BDT 300,000",
        coverageStart: "1 Apr 2026",
        nominee: "Ayan Kabir",
        phone: "01819-550042",
      },
      {
        employeeId: "AF-1088-C1",
        memberName: "Ayan Kabir",
        relation: "Child",
        designation: "Dependent",
        sumAssuredText: "BDT 150,000",
        coverageStart: "1 Apr 2026",
        nominee: "Nusrat Jahan",
        phone: "01819-550042",
      },
    ],
  },
  {
    id: "batch-alpha-force",
    title: "Alpha Force Enrollment Register",
    clientName: "Alpha Force",
    proposalNumber: "PRG-ENR-0326-002",
    status: "Needs cleanup",
    memberCount: 300,
    dependentCount: 184,
    coverageWindow: "15 Apr 2026 to 14 Apr 2027",
    sourceFiles: ["Enrollment Format  Alpha Force.xlsx", "Insurance Coverage & Income.xlsx / Alpha Force"],
    notes: [
      "This workbook behaves like a true batch register and should sit in a dedicated enrollment workspace.",
      "The insurer packet should only move once dependent dates and nominee information are complete.",
    ],
    validationFlags: [
      "11 rows missing mobile numbers",
      "5 spouse records missing relationship code",
      "Bank-facing roster summary not yet attached",
    ],
    rows: [
      {
        employeeId: "ALP-2101",
        memberName: "Md. Sohel Rana",
        relation: "Employee",
        designation: "Supervisor",
        sumAssuredText: "BDT 100,000",
        coverageStart: "15 Apr 2026",
        nominee: "Mst. Salma",
        phone: "01977-110210",
      },
      {
        employeeId: "ALP-2102",
        memberName: "Farhana Yeasmin",
        relation: "Employee",
        designation: "Officer",
        sumAssuredText: "BDT 100,000",
        coverageStart: "15 Apr 2026",
        nominee: "Sabbir Ahmed",
        phone: "01688-772311",
      },
      {
        employeeId: "ALP-2102-C1",
        memberName: "Safa Ahmed",
        relation: "Child",
        designation: "Dependent",
        sumAssuredText: "BDT 50,000",
        coverageStart: "15 Apr 2026",
        nominee: "Farhana Yeasmin",
        phone: "01688-772311",
      },
    ],
  },
  {
    id: "batch-prime-shine",
    title: "Prime Shine Enrollment Register",
    clientName: "Prime Shine",
    proposalNumber: "PRG-ENR-0326-003",
    status: "Draft packet",
    memberCount: 218,
    dependentCount: 96,
    coverageWindow: "1 May 2026 to 30 Apr 2027",
    sourceFiles: ["Enrollment format Prime Shine.xlsx", "Insurance Coverage & Income.xlsx / Prime Shine"],
    notes: [
      "Prime Shine has cleaner census structure but still needs insurer packet assembly and sign-off sequencing.",
      "A dispatch checklist should sit beside the roster, not inside a generic document card.",
    ],
    validationFlags: [
      "Client HR sign-off pending",
      "Coverage start date aligned, but 4 nominee names need English transliteration",
    ],
    rows: [
      {
        employeeId: "PS-7001",
        memberName: "Arefin Islam",
        relation: "Employee",
        designation: "Sales Manager",
        sumAssuredText: "BDT 150,000",
        coverageStart: "1 May 2026",
        nominee: "Sharmin Akter",
        phone: "01555-440701",
      },
      {
        employeeId: "PS-7002",
        memberName: "Sharmin Akter",
        relation: "Employee",
        designation: "Accounts Officer",
        sumAssuredText: "BDT 150,000",
        coverageStart: "1 May 2026",
        nominee: "Arefin Islam",
        phone: "01555-440702",
      },
      {
        employeeId: "PS-7002-S",
        memberName: "Rafid Islam",
        relation: "Spouse",
        designation: "Dependent",
        sumAssuredText: "BDT 100,000",
        coverageStart: "1 May 2026",
        nominee: "Sharmin Akter",
        phone: "01555-440702",
      },
    ],
  },
];

export const pricingScenarios: PricingScenario[] = [
  {
    id: "pricing-example-b2b",
    title: "Example B2B Financial Proposal",
    clientName: "Life Plus Bangladesh",
    status: "Ready for negotiation",
    segment: "Group Life & Medical",
    headlinePremiumText: "BDT 33,92,000 annual premium",
    sourceFiles: ["Financial Proposal LifePlus_Shanta-2026 (Group Insurance).pdf"],
    assumptions: [
      "3,000 members assumed for the life and medical quotation set.",
      "Life section combines group life and accidental death benefit pricing.",
      "Medical section separates hospitalization and OPC benefit slabs.",
    ],
    openItems: [
      "VAT bearing party needs final confirmation before release.",
      "Offer validity window should be tracked inside the commercial workspace.",
      "Dispatch should include a formal quotation cover note.",
    ],
    deliverables: [
      "Insurer-facing quotation PDF",
      "Pricing summary for partner review",
      "Negotiation note with scenario deltas",
    ],
    lineItems: [
      {
        label: "Group Life (GL) 100K",
        members: "3,000",
        benefit: "BDT 30,00,00,000 total sum assured",
        rateText: "2.44 per 1,000 SA",
        annualPremiumText: "BDT 7,32,000",
      },
      {
        label: "ADB 200K",
        members: "3,000",
        benefit: "Accidental death rider",
        rateText: "1.10 per 1,000 SA",
        annualPremiumText: "BDT 3,30,000",
      },
      {
        label: "IPC Hospitalization 20K",
        members: "3,000",
        benefit: "BDT 250,000 per member per year",
        rateText: "BDT 410 per member",
        annualPremiumText: "BDT 12,30,000",
      },
      {
        label: "OPC General 2K",
        members: "3,000",
        benefit: "BDT 50,000 per member per year",
        rateText: "BDT 550 per member",
        annualPremiumText: "BDT 11,00,000",
      },
    ],
  },
  {
    id: "pricing-alpha-force",
    title: "Alpha Force Pricing Summary",
    clientName: "Alpha Force",
    status: "Working draft",
    segment: "Group life commercials",
    headlinePremiumText: "Premium assumptions under review",
    sourceFiles: ["Insurance Coverage & Income.xlsx / Alpha Force", "Enrollment Format  Alpha Force.xlsx"],
    assumptions: [
      "Commercials are driven by policy count, coverage assumptions, and gross premium totals.",
      "Enrollment volume should stay synchronized with the pricing sheet before release.",
    ],
    openItems: [
      "Gross premium values need insurer sign-off.",
      "Coverage summary should be packaged with member census before submission.",
    ],
    deliverables: [
      "Commercial summary sheet",
      "Census-linked pricing note",
      "Partner-facing premium justification",
    ],
    lineItems: [
      {
        label: "Flat life cover",
        members: "300",
        benefit: "Employee block coverage",
        rateText: "Workbook basis",
        annualPremiumText: "Imported from pricing summary",
      },
      {
        label: "Dependent extension",
        members: "184",
        benefit: "Spouse and child inclusion",
        rateText: "Workbook basis",
        annualPremiumText: "Pending insurer confirmation",
      },
    ],
  },
  {
    id: "pricing-prime-shine",
    title: "Prime Shine Pricing Summary",
    clientName: "Prime Shine",
    status: "Commercial review",
    segment: "Group life commercials",
    headlinePremiumText: "Coverage and premium assumptions staged",
    sourceFiles: ["Insurance Coverage & Income.xlsx / Prime Shine", "Enrollment format Prime Shine.xlsx"],
    assumptions: [
      "Prime Shine pricing should stay aligned with the cleanest available enrollment roster.",
      "Commercial packet needs a compact insurer-ready summary rather than raw sheet export.",
    ],
    openItems: [
      "Confirm final sum assured slabs with insurer operations.",
      "Translate workbook assumptions into a shareable proposal note.",
    ],
    deliverables: [
      "Draft insurer commercial packet",
      "Final pricing comparison sheet",
      "Dispatch-ready premium summary",
    ],
    lineItems: [
      {
        label: "Base member block",
        members: "218",
        benefit: "Annual group cover",
        rateText: "Workbook basis",
        annualPremiumText: "Imported from pricing summary",
      },
      {
        label: "Family extension block",
        members: "96",
        benefit: "Dependent coverage",
        rateText: "Workbook basis",
        annualPremiumText: "Pending packaging",
      },
    ],
  },
];

export const claimChecklists: ClaimChecklist[] = [
  {
    id: "claims-fire",
    title: "Fire Claim Checklist",
    category: "Fire / Property",
    surveyorRequired: true,
    responseWindow: "Immediate notice, full pack before settlement review",
    sourceFiles: ["Documents are normally required for Claims.docx", "Fire Insurance Proposal Form_20230622_0001.pdf"],
    intakeSteps: [
      "Register claim notice and assign fire/property surveyor immediately.",
      "Collect bank-signed claim form and origin reports before financial loss validation.",
      "Validate stock and power-source evidence before insurer review begins.",
    ],
    deficiencySignals: [
      "Missing fire brigade report",
      "Stock records not countersigned by the concerned bank",
      "No generator/power authority documentation",
    ],
    documents: [
      {
        name: "Claim form signed by office and bank",
        owner: "Labaid + bank",
        purpose: "Formal claim initiation and banker confirmation",
        critical: true,
      },
      {
        name: "Fire brigade report in original",
        owner: "Claimant",
        purpose: "Official incident confirmation",
        critical: true,
      },
      {
        name: "Daily stock report for 90 days",
        owner: "Client + bank",
        purpose: "Pre-loss stock validation",
        critical: true,
      },
      {
        name: "Monthly stock statement",
        owner: "Client + bank",
        purpose: "Insurer inventory reconciliation",
        critical: true,
      },
      {
        name: "Trade license, fire license, power-source papers",
        owner: "Client",
        purpose: "Regulatory and risk-context validation",
        critical: true,
      },
      {
        name: "Layout plan and witness statements",
        owner: "Client",
        purpose: "Surveyor context and loss mapping",
        critical: false,
      },
    ],
  },
  {
    id: "claims-motor",
    title: "Motor Claim Checklist",
    category: "Motor",
    surveyorRequired: true,
    responseWindow: "Claim notice first, survey pack before repair approval",
    sourceFiles: ["Documents are normally required for Claims.docx", "Motor Insurance Proposal Form.pdf"],
    intakeSteps: [
      "Log claim intimation and assign motor surveyor.",
      "Collect repair estimates, driver statement, and police documents before approval.",
      "Validate registration, permit, and license before settlement movement.",
    ],
    deficiencySignals: [
      "Repair estimates from fewer than three garages",
      "Driver license not attested by BRTA",
      "No MVI report or GD copy attached",
    ],
    documents: [
      {
        name: "Claim intimation letter",
        owner: "Labaid / claimant",
        purpose: "First notice of loss",
        critical: true,
      },
      {
        name: "Signed claim form with bank seal",
        owner: "Labaid + bank",
        purpose: "Formal insurer claim filing",
        critical: true,
      },
      {
        name: "Three repair estimates",
        owner: "Garage / claimant",
        purpose: "Damage cost benchmarking",
        critical: true,
      },
      {
        name: "GD copy and MVI report",
        owner: "Claimant",
        purpose: "Police and inspection support",
        critical: true,
      },
      {
        name: "Registration, tax token, fitness, route permit",
        owner: "Claimant",
        purpose: "Vehicle compliance verification",
        critical: true,
      },
      {
        name: "Driver statement and BRTA-attested license",
        owner: "Driver / claimant",
        purpose: "Liability and driver eligibility review",
        critical: true,
      },
    ],
  },
  {
    id: "claims-health",
    title: "Health Reimbursement Checklist",
    category: "Health",
    surveyorRequired: false,
    responseWindow: "Document validation before reimbursement dispatch",
    sourceFiles: ["pragati.xlsx / Sheet14"],
    intakeSteps: [
      "Capture hospitalization and expense details in the insurer form.",
      "Validate all discharge papers, bills, prescriptions, and diagnostics.",
      "Complete coordinator forwarding block before insurer submission.",
    ],
    deficiencySignals: [
      "Expense breakup not aligned with the attached bills",
      "Missing discharge summary",
      "Coordinator forwarding section left blank",
    ],
    documents: [
      {
        name: "Completed health insurance claim form",
        owner: "Labaid operations",
        purpose: "Primary reimbursement submission",
        critical: true,
      },
      {
        name: "Original bills and money receipts",
        owner: "Claimant",
        purpose: "Expense proof",
        critical: true,
      },
      {
        name: "Discharge certificate and prescriptions",
        owner: "Hospital / claimant",
        purpose: "Medical evidence",
        critical: true,
      },
      {
        name: "Diagnostics and investigation reports",
        owner: "Hospital / claimant",
        purpose: "Clinical validation",
        critical: false,
      },
    ],
  },
  {
    id: "claims-travel",
    title: "Travel Assistance Claim Checklist",
    category: "Travel",
    surveyorRequired: false,
    responseWindow: "Emergency contact at incident time, final packet to TPA/assistance partner",
    sourceFiles: ["OMP New  Claim Process.docx", "OMP Proposal Form (New).pdf"],
    intakeSteps: [
      "Direct claimant or helper to the assistance provider during the medical emergency.",
      "Collect insurance certificate and claim form from the TPA flow.",
      "Forward full packet with supporting documents to Crisis24 or Van Ameyde.",
    ],
    deficiencySignals: [
      "No evidence of assistance-provider contact",
      "Insurance certificate not attached",
      "Trip dates not aligned with policy window",
    ],
    documents: [
      {
        name: "Travel claim form from assistance provider",
        owner: "Claimant / TPA",
        purpose: "Formal overseas mediclaim claim pack",
        critical: true,
      },
      {
        name: "Insurance certificate",
        owner: "Claimant",
        purpose: "Policy verification",
        critical: true,
      },
      {
        name: "Emergency medical records and bills",
        owner: "Hospital / claimant",
        purpose: "Treatment evidence",
        critical: true,
      },
      {
        name: "Trip itinerary and travel dates",
        owner: "Claimant",
        purpose: "Coverage window validation",
        critical: false,
      },
    ],
  },
];

export const assistanceProviders: AssistanceProvider[] = [
  {
    id: "provider-crisis24",
    name: "Crisis24",
    role: "Medical emergency assistance provider",
    supportWindow: "Emergency response / overseas assistance",
    phones: ["+44 207 902 7131"],
    emails: ["opsassist@crisis24.com", "corporateteam@crisis24.com"],
    address: "2 London Bridge, London, SE1 9RA, UK",
    useCases: [
      "Illness or accident abroad leading to hospital treatment",
      "Trip curtailment linked to medical emergency",
      "Travel-claim first response and case coordination",
    ],
    handoffPacket: [
      "Insurance certificate",
      "Traveller identity and callback number",
      "Hospital / treating doctor details",
      "Incident summary and travel dates",
    ],
    sourceFiles: ["OMP New  Claim Process.docx"],
  },
  {
    id: "provider-van-ameyde",
    name: "Van Ameyde",
    role: "Third-party administrator / assistance partner",
    supportWindow: "Claims handling and form issuance",
    phones: ["+44 208 315 0732"],
    emails: [],
    address: "Office G18, Bromley Old Town Hall, 30 Tweedy Road, Bromley, BR1 3FE, UK",
    useCases: [
      "Claim form issuance for overseas mediclaim",
      "Direct claim submission with relevant documentation",
      "Support escalation where a callback is required",
    ],
    handoffPacket: [
      "Completed claim form",
      "Insurance certificate",
      "Relevant treatment or incident documentation",
      "Claimant contact number for callback",
    ],
    sourceFiles: ["OMP New  Claim Process.docx"],
  },
];

export const assistanceRunbooks: AssistanceRunbook[] = [
  {
    title: "Medical emergency abroad",
    description: "Use the travel assistance path before trying to process a reimbursement claim internally.",
    steps: [
      "Contact Crisis24 immediately when hospitalization or emergency treatment begins.",
      "Provide claimant callback number to avoid reverse-charge failures.",
      "Collect case instructions from the assistance provider and log them in the portal.",
      "Move the full claim packet to the TPA / assistance desk once the incident stabilizes.",
    ],
  },
  {
    title: "Travel-claim handoff",
    description: "Labaid should package the claim to the external assistance desk, not only to the insurer queue.",
    steps: [
      "Request claim form from Crisis24 or Van Ameyde.",
      "Attach insurance certificate and relevant documentation.",
      "Track handoff status until acknowledgement is confirmed.",
      "Push the final record into claim settlement only after assistance review is complete.",
    ],
  },
];

export const knowledgeAssets: KnowledgeAsset[] = [
  {
    id: "knowledge-motor-deck",
    title: "Motor Underwriting & Claims Playbook",
    category: "Motor",
    audience: "Claims officer, surveyor, underwriting desk",
    summary: "A working insurer-side playbook derived from the motor underwriting and claims presentation deck.",
    sourceFile: "Motor insurance policy-LifePlus Bangladesh(Final)(Underwrite & Claims).pptx",
    linkedTabs: ["Claim Settlement", "Surveyor Desk", "Claims Checklists"],
    keyPoints: [
      "Motor claims require proposal, underwriting, exclusion, and claims-procedure knowledge in one place.",
      "Document requirements and settlement logic should be accessible during review, not hidden in offline slides.",
      "The deck implies a reusable training and SOP layer for insurer and ops teams.",
    ],
  },
  {
    id: "knowledge-travel-proposal",
    title: "Overseas Mediclaim Proposal Anatomy",
    category: "Travel",
    audience: "Proposal desk, medical reviewer",
    summary: "Maps the three-page overseas mediclaim flow into an insurer-ready digital form sequence.",
    sourceFile: "OMP Proposal Form (New).pdf",
    linkedTabs: ["Documents", "Travel Assistance", "Proposals"],
    keyPoints: [
      "Page 1 collects proposer identity, itinerary, and plan selection.",
      "Page 2 drives medical history and physician disclosure.",
      "Page 3 closes declaration, Schengen list, and benefit schedule acknowledgement.",
    ],
  },
  {
    id: "knowledge-fire-proposal",
    title: "Fire Risk Disclosure Guide",
    category: "Fire",
    audience: "Proposal desk, surveyor, insurer reviewer",
    summary: "Explains how fire proposal details, stock values, and site context should be captured before review.",
    sourceFile: "Fire Insurance Proposal Form_20230622_0001.pdf",
    linkedTabs: ["Documents", "Claims Checklists", "Surveyor Desk"],
    keyPoints: [
      "Risk location, construction type, occupation, and stock values are core insurer inputs.",
      "Proposal and claim review should stay connected through the same risk vocabulary.",
      "Poor scan quality makes a structured digital rendering and layout governance essential.",
    ],
  },
  {
    id: "knowledge-claims-note",
    title: "Claims Required Documents Reference",
    category: "Claims Ops",
    audience: "Claims coordinator, bank liaison, ops lead",
    summary: "Turns the raw claims-required-documents note into an operational checklist and deficiency playbook.",
    sourceFile: "Documents are normally required for Claims.docx",
    linkedTabs: ["Claims Checklists", "Claim Settlement", "Documents"],
    keyPoints: [
      "Fire and motor claims have different evidence packs and ownership paths.",
      "Bank countersignature and official reports are central, not optional attachments.",
      "Checklist completion should gate settlement review readiness.",
    ],
  },
  {
    id: "knowledge-commercials",
    title: "Group Pricing & Commercials Guide",
    category: "Commercials",
    audience: "Proposal desk, commercial lead, insurer liaison",
    summary: "Connects financial proposal layout, coverage summaries, and census-linked pricing into a reusable workspace.",
    sourceFile: "Financial Proposal LifePlus_Shanta-2026 (Group Insurance).pdf",
    linkedTabs: ["Pricing & Commercials", "Enrollment & Census", "Proposals"],
    keyPoints: [
      "Pricing should be scenario-based and linked to actual member counts.",
      "Commercial packet generation deserves its own tab instead of sitting inside documents.",
      "Offer validity, VAT notes, and dispatch deliverables should be tracked explicitly.",
    ],
  },
];
