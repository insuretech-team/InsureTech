import { Activity, HeartPulse, Landmark, ShieldCheck, Stethoscope, Waypoints } from "lucide-react";

import {
  approvalMatrix,
  claimCategoryMatrix,
  matrixDeskSignals,
  tpaSnapshot,
} from "@/lib/claims-intelligence";

export const tpaMatrixTabCopy = {
  operatingModel: {
    title: "TPA operating model",
    description:
      "Health claims are handled as a TPA-style workflow with network providers, cashless lanes, and manual fallback safeguards.",
    modelLabel: "Model",
    integrationLabel: "Integration",
    cards: [
      { title: "Claim modes", icon: HeartPulse },
      { title: "Fallback rules", icon: Stethoscope },
      { title: "Operating rules", icon: ShieldCheck },
    ],
  },
  priorities: {
    title: "Matrix desk priorities",
    description: "Lane-level controls that operators should watch before a claim reaches settlement.",
    surveyorLabel: "Surveyor queue",
    surveyorDescription:
      "Motor, fire/property, and pet lanes should not jump straight into settlement review without a survey outcome.",
    ownershipTitle: "Lane ownership",
    ownershipIcon: Waypoints,
  },
  matrixPanel: {
    title: "Claim matrix by plan category",
    description:
      "Category-by-category claim handling matrix shaped from the SRS, KBank data design notes, and current insurer product assumptions.",
    actionLabel: "Download matrix",
    exportFileName: "tpa-claim-matrix.csv",
    headers: ["Category", "Plan Type", "Intake Gate", "Claim Mode", "Settlement Rail", "TAT"] as const,
  },
  controlsPanel: {
    title: "Document and fraud controls",
    description: "What each category typically needs before approval can move forward.",
    documentsLabel: "Primary documents",
    settlementLabel: "Settlement rail",
    escalationLabel: "Escalation owner",
    fraudLabel: "Fraud and review checks",
  },
  approvalsPanel: {
    title: "Approval tiers",
    description: "SRS-aligned routing by claim amount.",
    aboveLabel: "Above",
  },
  touchpointsPanel: {
    title: "Operational partner touchpoints",
    description: "Who usually participates in each claim lane.",
    intakeGatePrefix: "Intake gate:",
  },
} as const;

export function getTouchpointIcon(category: string) {
  if (category === "Health Care & Hospitalization") return HeartPulse;
  if (category === "Motor") return Landmark;
  if (category === "Fire") return ShieldCheck;
  return Activity;
}

const tpaMatrixExportHeaders = [
  "category",
  "plan_type",
  "intake_gate",
  "claim_mode",
  "settlement_rail",
  "tpa_model",
  "escalation_owner",
  "typical_tat",
  "coverage_signal",
  "primary_documents",
  "fraud_checks",
  "partner_touchpoints",
] as const;

function buildCsvDownload(fileName: string, rows: string[][]) {
  const csv = `\uFEFF${rows
    .map((row) => row.map((cell) => `"${cell.replace(/"/g, '""')}"`).join(","))
    .join("\r\n")}`;
  const blob = new Blob([csv], { type: "text/csv;charset=utf-8" });
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = fileName;
  link.click();
  URL.revokeObjectURL(url);
}

export function getTpaClaimMatrixWorkspace() {
  return {
    snapshot: tpaSnapshot,
    deskSignals: matrixDeskSignals,
    claimMatrix: claimCategoryMatrix,
    approvals: approvalMatrix,
    surveyorLanes: claimCategoryMatrix.filter((row) => row.surveyorRequired),
  };
}

export function downloadTpaClaimMatrixCsv() {
  const rows = [
    [...tpaMatrixExportHeaders],
    ...claimCategoryMatrix.map((row) => [
      row.category,
      row.planType,
      row.intakeGate,
      row.claimMode,
      row.settlementRail,
      row.tpaModel,
      row.escalationOwner,
      row.typicalTat,
      row.coverageSignal,
      row.primaryDocuments.join(" | "),
      row.fraudChecks.join(" | "),
      row.partnerTouchpoints.join(" | "),
    ]),
  ];

  buildCsvDownload(tpaMatrixTabCopy.matrixPanel.exportFileName, rows);
}

export function getApprovalRangeLabel(min: number, max: number | null) {
  return `BDT ${min.toLocaleString()} - ${max ? max.toLocaleString() : tpaMatrixTabCopy.approvalsPanel.aboveLabel}`;
}
