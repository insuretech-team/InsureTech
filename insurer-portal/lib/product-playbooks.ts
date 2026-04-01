export interface ProductPlaybook {
  code: string;
  name: string;
  insurerName: string;
  category: string;
  audience: string;
  coverageLimitText: string;
  premiumText: string;
  policyTerm: string;
  ageRange: string;
  summary: string;
  requiredDocuments: string[];
  claimSteps: string[];
  exclusions: string[];
  operationalFlags: string[];
  keywords: string[];
}

export const productPlaybooks: ProductPlaybook[] = [
  {
    code: "PRAGATI-PET",
    name: "Cat & Dog Insurance",
    insurerName: "PRAGATI INSURANCE",
    category: "Pet",
    audience: "Individual pet owners",
    coverageLimitText: "BDT 10,000 / 20,000 / 30,000",
    premiumText: "BDT 1,035 - 3,105 yearly",
    policyTerm: "1 year",
    ageRange: "8 weeks to 10 years",
    summary:
      "Pet accident and critical illness cover for cats and dogs with hospitalization, surgery, diagnostics, and reimbursement-based claims.",
    requiredDocuments: [
      "Completed claim form",
      "Vaccination certificate",
      "Vet medical papers and bills",
      "Hospital bill if hospitalized",
      "Diagnostic report for critical illness",
    ],
    claimSteps: [
      "Receive incident notification by email",
      "Validate pet identity and vaccination records",
      "Review vet papers, bills, diagnostics, and hospitalization proof",
      "Confirm covered illness or accidental injury before reimbursement",
    ],
    exclusions: [
      "Pre-existing conditions",
      "Vaccination and regular checkups",
      "Breeding, pregnancy, neuter",
      "Cosmetic procedures",
      "Food and nutritional supplements",
      "Fight wounds and behavioral conditions",
    ],
    operationalFlags: [
      "Reimbursement flow",
      "Critical illness validation",
      "Vet bill audit required",
    ],
    keywords: ["pet", "cat", "dog", "veterinary", "critical illness"],
  },
  {
    code: "PRAGATI-NIBEDITA",
    name: "Nibedita",
    insurerName: "PRAGATI INSURANCE",
    category: "Personal Accident",
    audience: "Women aged 18-65",
    coverageLimitText: "Capital sum insured based benefits",
    premiumText: "Configured per distribution agreement",
    policyTerm: "1 year",
    ageRange: "18 to 65 years",
    summary:
      "Women-focused personal accident protection covering accidental death, disability, childbirth-related death, trauma allowance, and household goods damage under specific events.",
    requiredDocuments: [
      "Claim form",
      "Identity and age proof",
      "Medical report or disability schedule",
      "Event evidence for natural calamity or riot-related claims",
      "Police or hospital record for trauma allowance cases",
    ],
    claimSteps: [
      "Confirm claimant is within the women-only eligibility window",
      "Map claim to accidental death, disability, childbirth, or trauma allowance section",
      "Validate supporting medical and event records",
      "Apply schedule of compensation for partial or permanent disability",
    ],
    exclusions: [
      "Out-of-scope events outside declared benefit sections",
      "Unsupported trauma evidence",
      "Non-covered property damage scenarios",
    ],
    operationalFlags: [
      "Benefit-schedule adjudication",
      "Gender-specific eligibility",
      "High documentation sensitivity",
    ],
    keywords: ["nibedita", "women", "accident", "trauma", "childbirth"],
  },
  {
    code: "PRAGATI-PPA",
    name: "People's Personal Accident",
    insurerName: "PRAGATI INSURANCE",
    category: "Personal Accident",
    audience: "Mass-market accident cover",
    coverageLimitText: "BDT 100,000 per person",
    premiumText: "BDT 74 yearly",
    policyTerm: "1 year",
    ageRange: "As per group scheme rules",
    summary:
      "Annual accidental death and permanent disability cover designed for simple large-scale enrollment with a fixed capital sum insured.",
    requiredDocuments: [
      "Claim form",
      "Medical evidence of injury or disability",
      "Death certificate where applicable",
      "Accident report or incident evidence",
    ],
    claimSteps: [
      "Validate accident occurred within the policy period",
      "Assess death or disablement outcome within policy timing rules",
      "Apply 100 percent or 50 percent benefit table as appropriate",
      "Reject excluded causes such as intoxication, crime, war, or adventure sports",
    ],
    exclusions: [
      "Intentional injury or criminal activity",
      "Drug or alcohol influence",
      "War or terrorism",
      "Adventure sports",
      "Past illness",
    ],
    operationalFlags: [
      "Simple annual pricing",
      "Benefit-table based review",
      "Mass enrollment friendly",
    ],
    keywords: ["personal accident", "ppa", "death", "disability"],
  },
  {
    code: "PRAGATI-SSB",
    name: "Sorbojonin Surokkha Bima",
    insurerName: "PRAGATI INSURANCE",
    category: "Personal Accident",
    audience: "Bangladesh citizens age 16-75",
    coverageLimitText: "Up to BDT 200,000",
    premiumText: "BDT 115 yearly",
    policyTerm: "1 year",
    ageRange: "16 to 75 years",
    summary:
      "Broad Bangladesh-only accident protection with higher benefits for death, total sight loss, limb loss, and permanent total disablement.",
    requiredDocuments: [
      "Claim form",
      "Bangladesh citizenship proof",
      "Medical evidence of injury outcome",
      "Death certificate or disability certificate if applicable",
      "Accident evidence",
    ],
    claimSteps: [
      "Verify Bangladesh-only territorial scope",
      "Check entry and exit age eligibility",
      "Map claim to death, total loss, partial loss, or total disablement benefit",
      "Confirm no refund or remaining benefit assumptions are being applied after expiry",
    ],
    exclusions: [
      "Out-of-territory incidents",
      "Outside age eligibility",
      "Unsupported accident evidence",
    ],
    operationalFlags: [
      "Citizenship check required",
      "Territorial scope control",
      "Higher benefit cap",
    ],
    keywords: ["sorbojonin", "surokkha", "bangladesh", "accident", "disability"],
  },
];

export function findPlaybook(planName?: string, category?: string) {
  const query = `${planName ?? ""} ${category ?? ""}`.toLowerCase();
  return productPlaybooks.find((item) =>
    [item.name, item.category, ...item.keywords].some((token) => query.includes(token.toLowerCase())),
  );
}

export function getInsurerPlaybooks(insurerName?: string) {
  const name = (insurerName ?? "").toLowerCase();
  if (!name) return productPlaybooks;
  return productPlaybooks.filter((item) => item.insurerName.toLowerCase().includes(name) || name.includes("pragati"));
}
