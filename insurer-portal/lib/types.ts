export type PortalRole =
  | "SYSTEM_ADMIN"
  | "PARTNER_ADMIN"
  | "PARTNER_USER"
  | "BUSINESS_ADMIN"
  | "B2B_ORG_ADMIN"
  | "HR_MANAGER"
  | "VIEWER";

export interface PortalSessionSnapshot {
  userId: string;
  sessionId: string;
  role: PortalRole;
  portal: string;
  businessId: string;
  email: string;
  mobile: string;
  expiresAt: string;
}

export interface PortalInsurer {
  id: string;
  code: string;
  name: string;
  type: string;
  status: string;
  email: string;
  phone: string;
  websiteUrl: string;
  businessModel: string;
}

export interface PortalProduct {
  id: string;
  insurerId: string;
  code: string;
  name: string;
  category: string;
  status: string;
  premiumRangeText: string;
  coverageRangeText: string;
  features: string[];
}

export interface PortalProposal {
  id: string;
  proposalNumber: string;
  customerName: string;
  customerId: string;
  orderId: string;
  category: string;
  planName: string;
  coverageText: string;
  premiumText: string;
  status: string;
  submittedAt: string;
  decisionReason: string;
  source: "live" | "fallback";
}

export interface PortalClaim {
  id: string;
  claimNumber: string;
  insuredName: string;
  customerId: string;
  category: string;
  planName: string;
  claimedAmountText: string;
  approvedAmountText: string;
  settledAmountText: string;
  status: string;
  submittedAt: string;
  incidentDate: string;
  reason: string;
  source: "live" | "fallback";
}

export interface InsurerConfigForm {
  insurerId: string;
  apiBaseUrl: string;
  authType: string;
  authCredentials: string;
  businessModel: string;
  autoUnderwritingEnabled: boolean;
  realTimeClaimNotification: boolean;
  paymentTerms: string;
  claimSettlementDays: number;
}

export interface PortalOverview {
  currentInsurer: PortalInsurer | null;
  insurers: PortalInsurer[];
  products: PortalProduct[];
  proposals: PortalProposal[];
  claims: PortalClaim[];
  config: InsurerConfigForm | null;
  metrics: {
    insurerCount: number;
    productCount: number;
    proposalCount: number;
    approvedProposalCount: number;
    claimCount: number;
    settledClaimCount: number;
    requestedClaimCount: number;
    underReviewClaimCount: number;
  };
  source: "live" | "fallback" | "mixed";
}

export interface ApiEnvelope<T> {
  ok: boolean;
  message?: string;
  data?: T;
}

export interface LiveDocument {
  id: string;
  templateId: string;
  templateName: string;
  entityType: string;
  entityId: string;
  status: "Pending" | "Completed" | "Failed" | "Cancelled";
  fileUrl: string;
  fileSizeBytes: number;
  contentType: string;
  filename: string;
  generatedBy: string;
  generatedAt: string;
  kind: "proposal" | "claim" | "other";
}

export interface GenerateDocumentPayload {
  templateId: string;
  entityType: string;
  entityId: string;
  data?: Record<string, unknown>;
  outputFormat?: "pdf" | "docx" | "xlsx" | "";
}

export interface GenerateDocumentResult {
  documentId: string;
  fileUrl: string;
  message: string;
}

export interface LiveDocumentTemplate {
  id: string;
  name: string;
  type: string;
  description: string;
  outputFormat: string;
}

// ─── Document Library (from API, replaces hardcoded pragati-documents.ts) ─────

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
  templateDefinitionId: string;
  format: string;
  isGenerated: boolean;
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
