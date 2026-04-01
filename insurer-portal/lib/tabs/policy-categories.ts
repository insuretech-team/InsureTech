export const policyCategoriesStorageKey = "insurer-portal.policy-category-drafts";

export const policyCategoryOptions = ["Travel", "Auto", "Health", "Life", "Fire", "Device", "General"] as const;

export const policyCategoriesTabCopy = {
  footprintPanel: {
    title: "Category footprint",
    description: "Use live products and local notes to refine insurer categorization.",
    productsDescription: "Products currently mapped here.",
  },
  mapPanel: {
    title: "Product map",
    description: "Select a product to review and adjust its category notes.",
  },
  editorPanel: {
    title: "Category editor",
    description:
      "This editor is local-only for now, so you can shape the React flow before backend persistence is finalized.",
    coveragePrefix: "Coverage:",
    premiumPrefix: "Premium:",
    categoryLabel: "Category",
    tagsLabel: "Tags",
    tagsPlaceholder: "travel, baggage, assistance",
    notesLabel: "Notes",
    notesPlaceholder: "Add underwriting or storefront notes for this category...",
    playbookLabel: "Playbook insight",
    docsLabel: "Required docs",
    flagsLabel: "Operational flags",
    saveMessage: "Local category notes saved for this browser.",
    saveButton: "Save local category draft",
    emptyLabel: "Select a product to edit its category draft.",
  },
} as const;
