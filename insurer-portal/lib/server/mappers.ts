import {
  fallbackClaims,
  fallbackConfig,
  fallbackInsurers,
  fallbackProducts,
  fallbackProposals,
} from "@/lib/mock-data";
import type {
  InsurerConfigForm,
  PortalClaim,
  PortalInsurer,
  PortalProduct,
  PortalProposal,
} from "@/lib/types";
import { formatMoney, inferCategory, titleFromEnum } from "@/lib/utils";

type JsonRecord = Record<string, unknown>;

function record(value: unknown): JsonRecord {
  return value && typeof value === "object" && !Array.isArray(value) ? (value as JsonRecord) : {};
}

function text(value: unknown) {
  return typeof value === "string" ? value : "";
}

function moneyToNumber(value: unknown) {
  const bag = record(value);
  if (typeof bag.decimal_amount === "number") return bag.decimal_amount;
  if (typeof bag.decimalAmount === "number") return bag.decimalAmount;
  if (typeof bag.amount === "number") return bag.amount / 100;
  if (typeof bag.units === "number") return bag.units;
  if (typeof value === "number") return value;
  return 0;
}

function moneyText(minValue?: unknown, maxValue?: unknown) {
  const min = moneyToNumber(minValue);
  const max = moneyToNumber(maxValue);
  if (!min && !max) return "Not configured";
  if (min && max) return `${formatMoney(min)} - ${formatMoney(max)}`;
  return formatMoney(min || max);
}

function parseSubmissionPayload(value: unknown) {
  const raw = text(value);
  if (!raw) return {};
  try {
    return record(JSON.parse(raw));
  } catch {
    return {};
  }
}

function userDisplayName(source: JsonRecord, fallbackId: string) {
  const payload = parseSubmissionPayload(source.submission_payload);
  const nestedCustomer = record(payload.customer);
  return (
    text(payload.customer_name) ||
    text(payload.full_name) ||
    text(payload.name) ||
    text(nestedCustomer.full_name) ||
    text(nestedCustomer.name) ||
    (fallbackId ? `Customer ${fallbackId.slice(0, 8)}` : "Customer")
  );
}

export function mapInsurer(value: unknown): PortalInsurer {
  const source = record(value);
  const contact = record(source.contact_info);
  return {
    id: text(source.id),
    code: text(source.code),
    name: text(source.name) || fallbackInsurers[0].name,
    type: titleFromEnum(text(source.type)),
    status: titleFromEnum(text(source.status)),
    email: text(contact.email) || text(source.email),
    phone: text(contact.phone_number) || text(source.phone),
    websiteUrl: text(source.website_url),
    businessModel: text(source.business_model) || fallbackInsurers[0].businessModel,
  };
}

export function mapProduct(value: unknown): PortalProduct {
  const source = record(value);
  const name = text(source.name) || "Untitled product";
  return {
    id: text(source.id),
    insurerId: text(source.insurer_id),
    code: text(source.code),
    name,
    category: inferCategory(name, text(source.code), text(source.features)),
    status: titleFromEnum(text(source.status)),
    premiumRangeText: moneyText(source.min_premium, source.max_premium),
    coverageRangeText: moneyText(source.min_sum_assured, source.max_sum_assured),
    features: text(source.features)
      .split(/[,;\n]/)
      .map((item) => item.trim())
      .filter(Boolean)
      .slice(0, 5),
  };
}

export function mapProposal(value: unknown, sourceLabel: "live" | "fallback" = "live"): PortalProposal {
  const source = record(value);
  const productName = text(source.plan_name) || text(source.product_name) || text(source.product_id) || "Insurance plan";
  return {
    id: text(source.proposal_id) || text(source.id),
    proposalNumber: text(source.proposal_number) || "Proposal",
    customerName: userDisplayName(source, text(source.customer_id)),
    customerId: text(source.customer_id),
    orderId: text(source.order_id),
    category: inferCategory(productName, text(source.product_id), text(source.plan_id)),
    planName: productName,
    coverageText: formatMoney(moneyToNumber(source.proposed_sum_insured)),
    premiumText: formatMoney(moneyToNumber(source.proposed_premium)),
    status: titleFromEnum(text(source.status)),
    submittedAt: text(source.submitted_at) || text(source.created_at),
    decisionReason: text(source.decision_reason),
    source: sourceLabel,
  };
}

export function mapClaim(value: unknown, sourceLabel: "live" | "fallback" = "live"): PortalClaim {
  const source = record(value);
  const reason = text(source.incident_description);
  const inferredPlan = inferCategory(text(source.type), reason);
  return {
    id: text(source.claim_id) || text(source.id),
    claimNumber: text(source.claim_number) || "Claim",
    insuredName: text(source.customer_name) || (text(source.customer_id) ? `Customer ${text(source.customer_id).slice(0, 8)}` : "Insured"),
    customerId: text(source.customer_id),
    category: inferCategory(text(source.type), reason),
    planName: inferredPlan,
    claimedAmountText: formatMoney(moneyToNumber(source.claimed_amount)),
    approvedAmountText: formatMoney(moneyToNumber(source.approved_amount)),
    settledAmountText: formatMoney(moneyToNumber(source.settled_amount)),
    status: titleFromEnum(text(source.status)),
    submittedAt: text(source.submitted_at) || text(source.created_at),
    incidentDate: text(source.incident_date),
    reason,
    source: sourceLabel,
  };
}

export function mapConfig(insurerId: string, value: unknown): InsurerConfigForm {
  const source = record(value);
  return {
    insurerId,
    apiBaseUrl: text(source.api_base_url) || fallbackConfig.apiBaseUrl,
    authType: titleFromEnum(text(source.auth_type) || fallbackConfig.authType),
    authCredentials: text(source.auth_credentials) || fallbackConfig.authCredentials,
    businessModel: text(source.business_model) || fallbackConfig.businessModel,
    autoUnderwritingEnabled: Boolean(source.auto_underwriting_enabled ?? fallbackConfig.autoUnderwritingEnabled),
    realTimeClaimNotification: Boolean(source.real_time_claim_notification ?? fallbackConfig.realTimeClaimNotification),
    paymentTerms: text(source.payment_terms) || fallbackConfig.paymentTerms,
    claimSettlementDays:
      typeof source.claim_settlement_days === "number"
        ? source.claim_settlement_days
        : fallbackConfig.claimSettlementDays,
  };
}

export function mapLiveDocument(value: unknown): import("@/lib/types").LiveDocument {
  const source = record(value);
  const entityType = text(source.entity_type).toUpperCase();
  const kind =
    entityType.includes("PROPOSAL") ? "proposal" : entityType.includes("CLAIM") ? "claim" : "other";
  const rawStatus = titleFromEnum(text(source.status));
  const status = (["Completed", "Failed", "Cancelled", "Pending"].includes(rawStatus)
    ? rawStatus : "Pending") as import("@/lib/types").LiveDocument["status"];
  return {
    id: text(source.id),
    templateId: text(source.document_template_id) || text(source.template_id),
    templateName: text(source.template_name) || text(source.name) || "Document",
    entityType,
    entityId: text(source.entity_id),
    status,
    fileUrl: text(source.file_url),
    fileSizeBytes: typeof source.file_size_bytes === "number" ? source.file_size_bytes : 0,
    contentType: text(source.content_type),
    filename: text(source.filename),
    generatedBy: text(source.generated_by),
    generatedAt: text(source.generated_at),
    kind,
  };
}

export function fallbackContext() {
  return {
    insurers: fallbackInsurers,
    currentInsurer: fallbackInsurers[0],
    products: fallbackProducts,
    config: fallbackConfig,
    proposals: fallbackProposals,
    claims: fallbackClaims,
  };
}
