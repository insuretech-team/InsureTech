import { ClipboardList, FileCheck2, Shield } from "lucide-react";

import { claimChecklists } from "@/lib/docs-forms-operations";

export function getClaimChecklistsTabData() {
  return {
    checklists: claimChecklists,
    copy: {
      hero: {
        eyebrow: "Claim Evidence Desk",
        title: "Claim document checklists and readiness review",
        description:
          "The fire and motor claim note is effectively an operations playbook. This desk turns those source requirements into visible intake steps, mandatory evidence, and deficiency triggers.",
        primaryAction: "Open claim queue",
        secondaryAction: "Open source documents",
        checklistLabel: "Checklist lanes",
        checklistDescription: "Claim categories structured directly from the source set.",
        surveyorLabel: "Surveyor-routed lanes",
        surveyorDescription: "Fire and motor routes should gate on survey-led evidence.",
      },
      overviewPanel: {
        title: "Checklist overview",
        description: "How this workspace supports claim intake and review.",
      },
      lanesPanel: {
        title: "Checklist lanes",
        description: "Each checklist opens in a modal so the full evidence pack is readable without a cramped side panel.",
        docsLabel: "Required docs",
        blockersLabel: "Blockers",
        sourceSuffix: "source files",
        openLabel: "Open checklist",
        surveyorStatus: "Surveyor required",
        deskStatus: "Desk review",
      },
      stancePanel: {
        title: "Operational stance",
        description: "How this checklist workspace should influence downstream review.",
      },
      modal: {
        surveyorAction: "Open Surveyor Desk",
        settlementAction: "Open Claim Settlement",
        closeButton: "Close",
        docsLabel: "Required docs",
        blockersLabel: "Readiness blockers",
        reviewModeLabel: "Review mode",
        surveyorMode: "Survey-led lane",
        deskMode: "Desk-led lane",
        headers: ["Required document", "Owner", "Purpose", "Critical"] as const,
        criticalYes: "Yes",
        criticalNo: "No",
        blockersSidebarLabel: "Readiness blockers",
        sourceFilesLabel: "Source files",
        lanePostureTitle: "Lane posture",
        surveyorLaneBody: "Surveyor evidence should be complete before final approval or settlement movement.",
        deskLaneBody: "Coordinator review and checklist completion should drive claim readiness for this lane.",
      },
    },
    overviewItems: [
      "Settlement review sees a ready/not-ready claim state instead of document chaos.",
      "Surveyor-required categories show evidence expectations before handoff.",
      "Document ownership becomes explicit across bank, client, Labaid, and claimant.",
    ],
    stanceCards: [
      {
        icon: FileCheck2,
        text: "Claims should only enter settlement review after the checklist is materially complete.",
      },
      {
        icon: Shield,
        text: "Critical evidence should be visible before fraud, survey, or approval actions begin.",
      },
      {
        icon: ClipboardList,
        text: "Checklist ownership should sit with the coordinator desk, not with ad hoc notes.",
      },
    ],
  };
}
