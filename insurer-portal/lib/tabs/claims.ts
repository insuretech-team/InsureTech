export const claimStatusOptions = [
  "All",
  "Under Review",
  "Pending Documents",
  "Approved",
  "Settled",
  "Rejected",
] as const;

export type ClaimActionMode = "approve" | "reject" | "settle";

export interface ClaimActionState {
  mode: ClaimActionMode | null;
  claimNumber: string;
  claimId: string;
  reason: string;
  amount: string;
  paymentReference: string;
}

export const claimsTabCopy = {
  searchPlaceholder: "Search by claim number, insured, plan, or category",
  queueTitle: "Claim queue",
  queueDescription: "Move claims from review through approval and settlement.",
  detailTitle: "Claim detail",
  detailDescription: "Review documents, amounts, and settlement actions for the selected claim.",
  loadingLabel: "Loading claims...",
  emptyLabel: "No claims matched the current filter.",
  tableHeaders: ["Claim", "Plan", "Status", "Amount"] as const,
  validation: {
    rejectRequired: "A rejection reason is required before the claim can be rejected.",
    settleAmount: "Settlement amount must be a valid positive number.",
    saveFailed: "The claim update could not be saved.",
    successSuffix: "updated successfully.",
  },
  fields: {
    filed: "Filed",
    incidentDate: "Incident date",
    claimedAmount: "Claimed amount",
    approvedAmount: "Approved amount",
    settledAmount: "Settled amount",
    category: "Category",
    claimReason: "Claim reason",
    emptyReason: "No claim note available.",
    guideTitle: "Claim handling guide",
    requiredDocs: "Required docs",
    exclusions: "Key exclusions",
    flow: "Processing flow",
    claimMode: "Claim mode",
    approvalTier: "Approval tier",
    touchpoints: "Partner touchpoints",
    surveyorTitle: "Surveyor workflow",
    surveyorBody:
      "This claim category requires surveyor assessment before final decisioning. Use the dedicated surveyor workspace for inspection review, chat, web call, and additional document requests.",
    surveyorAction: "Open Surveyor Desk",
    settlementNoteTitle: "Settlement note",
    settlementNoteBody:
      "Claim updates now use structured action forms so amount, reason, and payment reference are ready to map directly to backend payloads.",
    emptyLabel: "Select a claim to review its detail and actions.",
  },
  buttons: {
    approve: "Approve claim",
    settle: "Settle claim",
    reject: "Reject claim",
  },
  modalButtons: {
    close: "Close",
    cancel: "Cancel",
  },
  approveModal: {
    title: "Approve claim",
    description: "Record the approval note and amount in a structured claim action form.",
    submitLabel: "Approve claim",
    noteLabel: "Approval note",
  },
  rejectModal: {
    title: "Reject claim",
    description: "Provide a reason so the rejection is traceable and reusable downstream.",
    submitLabel: "Reject claim",
    noteLabel: "Rejection reason",
  },
  settleModal: {
    title: "Settle claim",
    description: "Capture the settlement amount and payment reference before posting the update.",
    submitLabel: "Settle claim",
    amountLabel: "Settlement amount",
    referenceLabel: "Payment reference",
  },
};

export const initialClaimActionState: ClaimActionState = {
  mode: null,
  claimNumber: "",
  claimId: "",
  reason: "",
  amount: "",
  paymentReference: "",
};

export function extractClaimAmount(value: string) {
  const numeric = Number(value.replace(/[^\d.]/g, ""));
  return Number.isFinite(numeric) ? numeric : 0;
}

export function createClaimActionState(
  mode: ClaimActionMode,
  claim: {
    id: string;
    claimNumber: string;
    reason: string;
    approvedAmountText: string;
    claimedAmountText: string;
  },
): ClaimActionState {
  const defaultAmount = String(
    extractClaimAmount(claim.approvedAmountText) || extractClaimAmount(claim.claimedAmountText),
  );

  return {
    mode,
    claimId: claim.id,
    claimNumber: claim.claimNumber,
    reason: claim.reason || "",
    amount: mode === "reject" ? "" : defaultAmount,
    paymentReference: claim.claimNumber,
  };
}

export function getClaimActionModalConfig(mode: ClaimActionMode | null) {
  if (mode === "approve") {
    return {
      title: claimsTabCopy.approveModal.title,
      description: claimsTabCopy.approveModal.description,
      submitLabel: claimsTabCopy.approveModal.submitLabel,
      fields: [{ key: "reason", label: claimsTabCopy.approveModal.noteLabel, type: "textarea" as const }],
    };
  }

  if (mode === "reject") {
    return {
      title: claimsTabCopy.rejectModal.title,
      description: claimsTabCopy.rejectModal.description,
      submitLabel: claimsTabCopy.rejectModal.submitLabel,
      fields: [{ key: "reason", label: claimsTabCopy.rejectModal.noteLabel, type: "textarea" as const }],
    };
  }

  return {
    title: claimsTabCopy.settleModal.title,
    description: claimsTabCopy.settleModal.description,
    submitLabel: claimsTabCopy.settleModal.submitLabel,
    fields: [
      { key: "amount", label: claimsTabCopy.settleModal.amountLabel, type: "number" as const, min: 1 },
      { key: "paymentReference", label: claimsTabCopy.settleModal.referenceLabel, type: "text" as const },
    ],
  };
}

export function getClaimActionFailureMessage(action: "approve" | "reject" | "settle") {
  if (action === "approve") return "Unable to approve this claim.";
  if (action === "reject") return "Unable to reject this claim.";
  return "Unable to settle this claim.";
}

export function getClaimActionSuccessMessage(claimNumber: string) {
  return `Claim ${claimNumber} ${claimsTabCopy.validation.successSuffix}`;
}
