import type { PortalClaim } from "@/lib/types";

export interface ApprovalMatrixRow {
  min: number;
  max: number | null;
  approvalLevel: string;
  approvers: string;
  maxTat: string;
  mode: string;
}

export interface ClaimCategoryMatrix {
  category: string;
  planType: string;
  claimMode: "Cashless" | "Reimbursement" | "Hybrid";
  tpaModel: string;
  typicalTat: string;
  coverageSignal: string;
  settlementRail: string;
  escalationOwner: string;
  intakeGate: string;
  surveyorRequired: boolean;
  surveyorType: string;
  primaryDocuments: string[];
  fraudChecks: string[];
  partnerTouchpoints: string[];
}

export const approvalMatrix: ApprovalMatrixRow[] = [
  {
    min: 0,
    max: 10000,
    approvalLevel: "L1 Auto / Officer",
    approvers: "System auto-approval or Claims Officer",
    maxTat: "24 hours",
    mode: "Fast-track / ZHTC candidate",
  },
  {
    min: 10001,
    max: 50000,
    approvalLevel: "L2 Manager",
    approvers: "Claims Manager",
    maxTat: "3 days",
    mode: "Manager review",
  },
  {
    min: 50001,
    max: 200000,
    approvalLevel: "L3 Head",
    approvers: "Business Admin + Focal Person (joint)",
    maxTat: "7 days",
    mode: "Joint approval",
  },
  {
    min: 200001,
    max: null,
    approvalLevel: "Board",
    approvers: "Board + Insurer approval",
    maxTat: "15 days",
    mode: "Escalated governance",
  },
];

export const tpaSnapshot = {
  model: "Health as TPA with network management and fee-per-claim operating model",
  network: "LabAid Hospitals at Phase M1 (5 locations), expanding to 20+ partner hospitals in later phases",
  integration: "HL7 FHIR R4 with OAuth 2.0 + JWT, patient consent and digital signature before EHR access",
  claimModes: ["Cashless hospital-network claims", "Reimbursement with bank/MFS payout"],
  fallbacks: [
    "Queue for manual verification if EHR integration times out",
    "Notify hospital focal point by SMS if systems are unavailable",
    "Notify support team through Slack/email escalation",
  ],
  channels: ["Hospital EHR", "bKash", "Nagad", "Bank transfer", "SMS", "Partner portal"],
  operatingRules: [
    "300 DPI minimum claim document quality",
    "10 MB per file and 50 MB total document cap per claim",
    "Co-pay and deductible formula at product level",
    "Provider validation against approved network list",
  ],
};

export const claimCategoryMatrix: ClaimCategoryMatrix[] = [
  {
    category: "Health Care & Hospitalization",
    planType: "Network hospitalization and SME health",
    claimMode: "Hybrid",
    tpaModel: "Primary TPA lane with LabAid hospital network and manual fallback",
    typicalTat: "Cashless pre-auth immediate to same day; reimbursement 7-15 working days",
    coverageSignal: "Hospitalization, surgery, room rent, IPD/OPD, co-pay and deductible controls",
    settlementRail: "Cashless hospital settlement or reimbursement to member bank/MFS account",
    escalationOwner: "TPA desk and insurer claims manager",
    intakeGate: "Provider validation, admission notice, and benefit schedule match",
    surveyorRequired: false,
    surveyorType: "Medical desk / hospital focal review",
    primaryDocuments: [
      "Claim form",
      "Hospital bill and discharge summary",
      "Diagnostic report",
      "Network Hospital ID",
      "Bank/MFS details for reimbursement",
    ],
    fraudChecks: [
      "Non-network provider flag",
      "Rapid policy-to-claim movement",
      "Claim amount near coverage limit",
    ],
    partnerTouchpoints: ["Hospital focal person", "TPA desk", "Insurer claims manager"],
  },
  {
    category: "Personal Accident",
    planType: "Nibedita / People's Personal Accident / Sorbojonin Surokkha Bima",
    claimMode: "Reimbursement",
    tpaModel: "Insurer-led review with document and benefit-schedule validation",
    typicalTat: "24 hours to 7 days based on amount tier",
    coverageSignal: "Accidental death, disability, trauma allowance, named benefit schedule",
    settlementRail: "Direct payout after entitlement and beneficiary verification",
    escalationOwner: "Claims officer with amount-tier approval chain",
    intakeGate: "Incident narrative, identity proof, and policy schedule confirmation",
    surveyorRequired: false,
    surveyorType: "Claims officer review",
    primaryDocuments: [
      "Claim form",
      "Medical or disability certificate",
      "Death certificate where applicable",
      "Incident or police report",
      "Identity / citizenship evidence",
    ],
    fraudChecks: [
      "Duplicate claim type frequency",
      "Geographic anomaly",
      "Amount exactly matching benefit cap",
    ],
    partnerTouchpoints: ["Partner agent", "Focal person", "Claims officer"],
  },
  {
    category: "Pet",
    planType: "Cat & Dog Insurance",
    claimMode: "Reimbursement",
    tpaModel: "Clinic-assisted but document-led reimbursement workflow",
    typicalTat: "3 to 10 working days after vet document verification",
    coverageSignal: "Accidents, surgery, hospitalization, diagnostics, critical illness sections",
    settlementRail: "Reimbursement after surveyor and veterinary record validation",
    escalationOwner: "Veterinary surveyor and insurer claims desk",
    intakeGate: "Vaccination history, treatment chronology, and diagnosis confirmation",
    surveyorRequired: true,
    surveyorType: "Veterinary surveyor / field verifier",
    primaryDocuments: [
      "Completed claim form",
      "Vaccination certificate",
      "Vet medical papers and bills",
      "Hospitalization bill if admitted",
      "Diagnostic report",
    ],
    fraudChecks: [
      "Missing vaccination history",
      "Non-covered illness list mismatch",
      "Duplicate or repeated surgery expense patterns",
    ],
    partnerTouchpoints: ["Vet clinic", "Insurer claims desk"],
  },
  {
    category: "Motor",
    planType: "Motor Insurance",
    claimMode: "Hybrid",
    tpaModel: "Garage-network operational model with digital survey and optional cashless repair",
    typicalTat: "Target digital settlement within 48 hours for straightforward claims",
    coverageSignal: "Accident, theft, fire, natural calamity, third-party liability",
    settlementRail: "Garage settlement for approved repairs or reimbursement against final bills",
    escalationOwner: "Motor surveyor, garage coordinator, and claims manager",
    intakeGate: "Incident notice, vehicle identity match, and initial damage evidence",
    surveyorRequired: true,
    surveyorType: "Motor surveyor",
    primaryDocuments: [
      "Claim form",
      "Repair estimate or bills",
      "Survey report",
      "FIR if required",
      "Vehicle registration and driver license",
      "Photos / videos",
    ],
    fraudChecks: [
      "Same-device or same-vehicle repeat incidents",
      "Location anomaly",
      "Provider / garage validation",
    ],
    partnerTouchpoints: ["Partner garage", "Surveyor", "Claims manager"],
  },
  {
    category: "Travel",
    planType: "Travel Insurance",
    claimMode: "Reimbursement",
    tpaModel: "Assistance-first, reimbursement-backed travel workflow",
    typicalTat: "3 to 7 working days after document validation",
    coverageSignal: "Medical emergency, cancellation, baggage loss, delay support",
    settlementRail: "Assistance partner coordination followed by reimbursement settlement",
    escalationOwner: "Travel assistance provider and insurer claims desk",
    intakeGate: "Trip eligibility, incident timing, and overseas evidence pack",
    surveyorRequired: false,
    surveyorType: "Travel assistance review",
    primaryDocuments: [
      "Claim form",
      "Travel ticket / itinerary",
      "Medical papers or airline certificate",
      "Loss / delay confirmation",
      "Payment receipts",
    ],
    fraudChecks: [
      "Duplicate travel incident claims",
      "Coverage limit matching",
      "Date mismatch with trip window",
    ],
    partnerTouchpoints: ["Travel assistance provider", "Insurer claims desk"],
  },
  {
    category: "Fire / Property",
    planType: "Property, fire, burglary, casualty",
    claimMode: "Reimbursement",
    tpaModel: "Survey-driven insurer workflow with optional smart risk assessment",
    typicalTat: "Up to 14 days for standard fire/property loss assessment",
    coverageSignal: "Fire, theft, burglary, natural disaster, high-value item extensions",
    settlementRail: "Assessment-led reimbursement after survey and loss quantification",
    escalationOwner: "Property surveyor and insurer approver",
    intakeGate: "Peril confirmation, site photos, and insured asset schedule mapping",
    surveyorRequired: true,
    surveyorType: "Fire / property surveyor",
    primaryDocuments: [
      "Claim form",
      "Survey report",
      "Loss photos / videos",
      "Police / fire service report where applicable",
      "Repair or replacement estimate",
    ],
    fraudChecks: [
      "Coverage-cap matching",
      "Geographic anomaly",
      "Suspicious repeat peril pattern",
    ],
    partnerTouchpoints: ["Surveyor", "Property assessor", "Insurer approver"],
  },
];

export const pragatiBusinessModelPillars = [
  "Dual-Entity Market Coverage",
  "Bancassurance",
  "B2B / Corporate Group Solutions",
  "Financial Inclusion & Micro-Insurance",
  "Health-Tech & TPA",
  "Reinsurance Backing",
] as const;

export const pragatiBusinessModelNotes = [
  "Pragati sits well inside the portal as both a retail-facing insurer and a partner-facing carrier for institutional schemes.",
  "The insurer’s role in bank-linked distribution, partner programs, and group enrollment makes proposal, census, and claims coordination equally important.",
  "Health-tech workflows need TPA-grade controls, while non-life categories still rely on surveyor, branch, and insurer approval lanes.",
] as const;

export const matrixDeskSignals = [
  {
    title: "Survey-driven categories",
    value: "Motor, fire/property, and pet",
    detail: "These lanes should move to Surveyor Desk before financial settlement review.",
  },
  {
    title: "Fastest operating lane",
    value: "Health cashless pre-auth",
    detail: "Provider validation and hospital network routing drive same-day decisions for eligible admissions.",
  },
  {
    title: "Primary escalation model",
    value: "Amount tier plus category owner",
    detail: "Large losses escalate by financial authority, while category ownership keeps lane-specific evidence review intact.",
  },
] as const;

export const pragatiOfficialSnapshot = {
  established: "1986",
  hotline: "09613115511",
  supportPhone: "+88-02-55012680-2",
  headquarters: "Pragati Insurance Bhaban, 20-21 Kawran Bazar, Dhaka-1215",
  website: "www.pragatiinsurance.com",
  products: ["Fire", "Motor", "Health Care & Hospitalization", "Marine", "All Risks", "Aviation", "Miscellaneous"],
  recentSignals: [
    "IDRA Insurance Excellence Award 2025 - first position in non-life",
    "39th AGM reported 2024 gross premium Tk. 2503.65 million",
    "2024 net claim settled Tk. 289.79 million",
    "Business cooperation and medical support agreement signed with LabAid",
  ],
};

export function routeApprovalTier(claimAmount: number) {
  return (
    approvalMatrix.find((row) => claimAmount >= row.min && (row.max === null || claimAmount <= row.max)) ??
    approvalMatrix[approvalMatrix.length - 1]
  );
}

export function findClaimCategoryMatrix(category?: string, planName?: string) {
  const query = `${category ?? ""} ${planName ?? ""}`.toLowerCase();

  return (
    claimCategoryMatrix.find((item) => item.category.toLowerCase().includes(query)) ??
    claimCategoryMatrix.find((item) =>
      [item.category, item.planType].some((value) => query.includes(value.toLowerCase())),
    ) ??
    claimCategoryMatrix.find((item) => {
      if (query.includes("nibedita") || query.includes("personal accident") || query.includes("sorbojonin")) {
        return item.category === "Personal Accident";
      }
      if (query.includes("pet") || query.includes("cat") || query.includes("dog")) return item.category === "Pet";
      if (query.includes("motor") || query.includes("auto") || query.includes("vehicle")) return item.category === "Motor";
      if (query.includes("fire") || query.includes("property")) return item.category === "Fire / Property";
      return false;
    })
  );
}

export function buildSurveyorReview(claim: PortalClaim) {
  const query = `${claim.category} ${claim.planName}`.toLowerCase();

  const notes =
    query.includes("pet")
      ? [
          "Vaccination record and treatment chronology should be cross-checked with the attending clinic.",
          "Final diagnostic report and hospitalization bill must be attached before reimbursement decision.",
        ]
      : query.includes("motor") || query.includes("auto") || query.includes("vehicle")
        ? [
            "Damage photos and repair estimate should be matched against the incident description.",
            "Surveyor should validate garage quotation, affected parts, and any third-party liability indicators.",
          ]
        : [
            "Site inspection and loss photos should be reviewed against the submitted incident narrative.",
            "Surveyor should confirm extent of damage, probable cause, and adequacy of supporting estimates.",
          ];

  return {
    name:
      query.includes("pet")
        ? "Dr. Samia Rahman"
        : query.includes("motor") || query.includes("auto") || query.includes("vehicle")
          ? "Md. Rezaul Karim"
          : "Engr. Farhan Islam",
    title:
      query.includes("pet")
        ? "Veterinary Surveyor"
        : query.includes("motor") || query.includes("auto") || query.includes("vehicle")
          ? "Motor Surveyor"
          : "Fire and Property Surveyor",
    status: claim.status === "Pending Documents" ? "Awaiting documents" : "Review in progress",
    notes,
  };
}

export function isSurveyorRequiredClaim(claim: Pick<PortalClaim, "category" | "planName">) {
  return Boolean(findClaimCategoryMatrix(claim.category, claim.planName)?.surveyorRequired);
}
