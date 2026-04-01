import { ArrowRightLeft, BadgeDollarSign, ReceiptText } from "lucide-react";

import { pricingScenarios } from "@/lib/docs-forms-operations";

export function getPricingCommercialsTabData() {
  return {
    scenarios: pricingScenarios,
    copy: {
      hero: {
        eyebrow: "Commercials Desk",
        title: "Pricing scenarios and insurer commercial proposals",
        description:
          "The financial proposal PDF and coverage-income sheets imply structured quote preparation, scenario comparison, and dispatch deliverables instead of a raw workbook download.",
        primaryAction: "Open supporting documents",
        secondaryAction: "View census linkage",
        leadQuoteLabel: "Lead quote",
        leadQuoteDescription: "Current flagship commercial packet from the source set.",
        activeCasesLabel: "Active pricing cases",
        activeCasesDescription: "Each should track assumptions, open items, and deliverables.",
      },
      outputsPanel: {
        title: "Commercial outputs",
        description: "Artifacts the pricing desk should own.",
      },
      casesPanel: {
        title: "Pricing cases",
        description: "Commercial packets derived from the available source files.",
        premiumLabel: "Headline premium",
      },
      detailPanel: {
        description: "Scenario detail, pricing lines, and release blockers.",
        actionLabel: "Open source packet",
        clientLabel: "Client",
        segmentLabel: "Segment",
        headers: ["Benefit line", "Members", "Benefit", "Rate", "Annual premium"] as const,
        assumptionsLabel: "Assumptions",
        openItemsLabel: "Open items",
      },
      deliverablesPanel: {
        title: "Release deliverables",
        description: "What the commercial desk should generate from this scenario.",
      },
      sourceFilesPanel: {
        title: "Source files",
        description: "Files currently backing the commercial packet.",
      },
      valuePanel: {
        title: "Commercial desk value",
        description: "Why these pricing tasks need a dedicated workspace.",
      },
    },
    commercialOutputs: [
      "Scenario comparison sheet",
      "Formal insurer quotation packet",
      "Coverage and premium summary",
      "Dispatch note with validity and VAT assumptions",
    ],
    valueCards: [
      { icon: BadgeDollarSign, text: "Pricing is scenario-driven and not just a file attachment." },
      { icon: ArrowRightLeft, text: "Enrollment counts and quote assumptions need visible reconciliation." },
      { icon: ReceiptText, text: "Commercial dispatch output should be tracked like a workflow artifact." },
    ],
  };
}
