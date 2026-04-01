export const proposalStatusOptions = ["All", "In Review", "Approved", "Rejected"] as const;

export type ProposalActionMode = "approve" | "reject";

export interface ProposalActionState {
  mode: ProposalActionMode | null;
  proposalId: string;
  proposalNumber: string;
  note: string;
}

export const proposalsTabCopy = {
  searchPlaceholder: "Search by proposal number, customer, plan, or category",
  queueTitle: "Proposal queue",
  queueDescription: "Search and move proposals through underwriting decisions.",
  detailTitle: "Proposal detail",
  detailDescription: "Review the selected proposal before you take action.",
  loadingLabel: "Loading proposals...",
  emptyLabel: "No proposals matched the current filter.",
  tableHeaders: ["Proposal", "Plan", "Status", "Submitted"] as const,
  fields: {
    plan: "Plan",
    category: "Category",
    coverage: "Coverage",
    premium: "Premium",
    orderId: "Order ID",
    orderFallback: "Unavailable",
    submitted: "Submitted",
    decisionNote: "Decision note",
    emptyDecisionNote: "No decision note has been added yet.",
    reviewGuide: "Product review guide",
    audience: "Audience",
    eligibility: "Eligibility",
    requiredDocs: "Required docs",
    flags: "Operational flags",
    documentsPack: "Documents pack",
    documentsAction: "Open documents",
  },
  actionButtons: {
    approve: "Approve proposal",
    reject: "Reject proposal",
    emptyLabel: "Select a proposal to review its details.",
  },
  modalButtons: {
    close: "Close",
    cancel: "Cancel",
  },
  approveModal: {
    title: "Approve proposal",
    description: "Capture the approval note in a proper form instead of a browser prompt.",
    submitLabel: "Approve proposal",
    noteLabel: "Approval note",
  },
  rejectModal: {
    title: "Reject proposal",
    description: "Add a clear rejection reason for audit and downstream support teams.",
    submitLabel: "Reject proposal",
    noteLabel: "Rejection reason",
  },
};

export const proposalsValidationCopy = {
  rejectRequired: "A rejection reason is required before the proposal can be rejected.",
  saveFailed: "The decision could not be saved.",
} as const;

export const initialProposalActionState: ProposalActionState = {
  mode: null,
  proposalId: "",
  proposalNumber: "",
  note: "",
};

export function createProposalActionState(
  mode: ProposalActionMode,
  proposal: {
    id: string;
    proposalNumber: string;
    decisionReason: string;
  },
): ProposalActionState {
  return {
    mode,
    proposalId: proposal.id,
    proposalNumber: proposal.proposalNumber,
    note: proposal.decisionReason || "",
  };
}

export function getProposalActionModalConfig(mode: ProposalActionMode | null) {
  if (mode === "approve") {
    return {
      title: proposalsTabCopy.approveModal.title,
      description: proposalsTabCopy.approveModal.description,
      submitLabel: proposalsTabCopy.approveModal.submitLabel,
      fields: [{ key: "note", label: proposalsTabCopy.approveModal.noteLabel, type: "textarea" as const }],
    };
  }

  return {
    title: proposalsTabCopy.rejectModal.title,
    description: proposalsTabCopy.rejectModal.description,
    submitLabel: proposalsTabCopy.rejectModal.submitLabel,
    fields: [{ key: "note", label: proposalsTabCopy.rejectModal.noteLabel, type: "textarea" as const }],
  };
}

export function getProposalActionFailureMessage(action: "approve" | "reject") {
  return action === "approve" ? "Unable to approve this proposal." : "Unable to reject this proposal.";
}

export function getProposalActionSuccessMessage(
  proposalNumber: string,
  action: "approve" | "reject",
) {
  const status = action === "approve" ? "approved" : "rejected";
  return `Proposal ${proposalNumber} ${status} successfully.`;
}
