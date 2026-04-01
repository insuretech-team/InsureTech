import { TimerReset } from "lucide-react";

import { assistanceProviders, assistanceRunbooks } from "@/lib/docs-forms-operations";

export function getTravelAssistanceTabData() {
  return {
    providers: assistanceProviders,
    runbooks: assistanceRunbooks,
    copy: {
      hero: {
        eyebrow: "Overseas mediclaim assistance desk",
        title: "Travel assistance contacts and overseas claim handoff",
        description:
          "The OMP claim-process note clearly shows Crisis24 and Van Ameyde as active assistance and TPA partners. This workspace makes those contacts, handoff packets, and emergency runbooks visible.",
        primaryAction: "Open travel claims",
        secondaryAction: "View OMP source forms",
        providersLabel: "Assistance providers",
        providersDescription: "Named travel-support partners from the source doc.",
        runbooksLabel: "Runbooks",
        runbooksDescription: "Emergency and handoff flows modeled from OMP instructions.",
      },
      focusPanel: {
        title: "Assistance routing focus",
        description: "What the OMP process note requires the team to do.",
      },
      providersPanel: {
        title: "Provider directory",
        description: "Select an assistance partner to inspect contact and handoff details.",
      },
      detailPanel: {
        description: "Emergency contact posture and packet expectations.",
        phonesLabel: "Phones",
        emailsLabel: "Emails",
        emptyEmail: "No email listed in the source note.",
        useCasesLabel: "Use cases",
        addressLabel: "Address",
        handoffLabel: "Handoff packet",
        sourceFilesLabel: "Source file",
      },
      runbooksPanel: {
        title: "Travel assistance runbooks",
        description: "Operational flows derived from the OMP note.",
      },
      rulesPanel: {
        title: "Time-sensitive rules",
        description: "Travel support expectations to surface in the dashboard.",
      },
    },
    routingFocus: [
      "Emergency overseas cases should route to the assistance provider first.",
      "Callback numbers and contact instructions matter operationally and should be visible in-system.",
      "Claim packets need a TPA handoff log, not just a reimbursement note.",
    ],
    timeSensitiveRules: [
      {
        icon: TimerReset,
        text: "Emergency contact should happen when treatment begins, not after return.",
      },
      {
        icon: TimerReset,
        text: "Callback number should always be captured for international assistance flow.",
      },
      {
        icon: TimerReset,
        text: "Claim settlement should start after assistance handoff is acknowledged.",
      },
    ],
  };
}
