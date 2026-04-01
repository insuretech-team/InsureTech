import { knowledgeAssets } from "@/lib/docs-forms-operations";

export function getKnowledgeCenterTabData() {
  return {
    assets: knowledgeAssets,
    copy: {
      hero: {
        eyebrow: "Playbooks & Guidance",
        title: "Operational playbooks for underwriting, claims, and proposals",
        description:
          "The motor deck, OMP forms, fire proposal source, and claims note imply a searchable playbook layer so the team can work from shared operational guidance instead of offline files.",
        primaryAction: "Open document source archive",
        secondaryAction: "Open claims checklists",
        modulesLabel: "Knowledge modules",
        modulesDescription: "Playbooks and references grounded in the source files.",
        portalImpactLabel: "Portal impact",
      },
      searchPanel: {
        title: "Search the playbooks",
        description: "Filter by category, audience, source file, or linked workflow.",
        placeholder: "Search knowledge center",
      },
      listPanel: {
        title: "Knowledge modules",
        description: "Select a module to inspect the playbook content.",
      },
      detailPanel: {
        audienceLabel: "Audience",
        linkedTabsLabel: "Linked tabs",
        sourceFileLabel: "Source file",
      },
    },
    impactSummary:
      "These assets inform proposals, documents, claims, surveyor review, and commercial packaging.",
  };
}
