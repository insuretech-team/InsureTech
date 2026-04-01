import { enrollmentBatches } from "@/lib/docs-forms-operations";

export function getEnrollmentCensusTabData() {
  const totalMembers = enrollmentBatches.reduce((sum, item) => sum + item.memberCount, 0);
  const totalDependents = enrollmentBatches.reduce((sum, item) => sum + item.dependentCount, 0);
  const outstandingFlags = enrollmentBatches.reduce((sum, item) => sum + item.validationFlags.length, 0);

  return {
    batches: enrollmentBatches,
    metrics: { totalMembers, totalDependents, outstandingFlags },
    copy: {
      hero: {
        eyebrow: "Member Enrollment Desk",
        title: "Member census and dependent roster management",
        description:
          "The Pragati census sheet and the Alpha Force / Prime Shine enrollment files imply a dedicated desk for member rosters, dependent cleanup, and insurer-ready dispatch.",
        primaryAction: "Open source templates",
        secondaryAction: "Proposal queue",
        membersLabel: "Covered lives in staging",
        membersDescription: "Employee rows currently modeled across the docs-driven batches.",
        dependentsLabel: "Dependents + blockers",
        dependentsDescriptionSuffix: "active cleanup flags before dispatch.",
      },
      focusPanel: {
        title: "Operational focus",
        description: "What this workspace should control before insurer dispatch.",
      },
      batchesPanel: {
        title: "Enrollment batches",
        description: "Source-linked census packets waiting for cleanup or dispatch.",
        membersLabel: "Members",
        dependentsLabel: "Dependents",
      },
      detailPanel: {
        description: "Batch overview, validation flags, and source linkage.",
        actionLabel: "Open in Documents",
        coverageLabel: "Coverage window",
        postureLabel: "Dispatch posture",
        notesLabel: "Operational notes",
        flagsLabel: "Validation flags",
        filesLabel: "Source files",
      },
      rosterPanel: {
        title: "Roster preview",
        description: "Sample rows from the selected census batch.",
        headers: ["Member ID", "Name", "Relation", "Designation", "Sum assured", "Coverage start", "Nominee", "Phone"] as const,
      },
      checklistPanel: {
        title: "Dispatch checklist",
        description: "What Labaid should complete before sending to Pragati.",
      },
    },
    operationalFocus: [
      "Roster validation should happen before the insurer receives a proposal packet.",
      "Dependent and nominee issues need a visible queue, not hidden workbook edits.",
      "Commercial summaries should stay tied to the same census batch that drives pricing.",
    ],
    dispatchChecklist: [
      "Member and dependent counts reconciled against the quote.",
      "Nominee names, relation codes, and contact numbers normalized.",
      "Census packet attached to the correct proposal and document set.",
    ],
  };
}
