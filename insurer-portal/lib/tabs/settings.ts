export const settingsAuthTypeOptions = [
  { value: "AUTH_TYPE_BEARER", label: "Bearer token" },
  { value: "AUTH_TYPE_BASIC", label: "Basic auth" },
  { value: "AUTH_TYPE_API_KEY", label: "API key" },
] as const;

export const settingsFeatureToggles = [
  {
    key: "autoUnderwritingEnabled",
    title: "Auto underwriting",
    description: "Let straightforward proposals move with reduced manual intervention.",
  },
  {
    key: "realTimeClaimNotification",
    title: "Real-time claim notifications",
    description: "Keep insurer staff informed as new claims and updates hit the queue.",
  },
] as const;

export const settingsSummaryMetrics = [
  { key: "productCount", label: "Products" },
  { key: "proposalCount", label: "Open proposals" },
  { key: "claimCount", label: "Claims" },
] as const;

export const settingsTabCopy = {
  formPanel: {
    title: "Insurer configuration",
    description: "Adjust operational settings used by the insurer integration.",
  },
  formFields: {
    apiBaseUrl: "API base URL",
    authType: "Auth type",
    businessModel: "Business model",
    authCredentials: "Auth credentials",
    paymentTerms: "Payment terms",
    claimSettlementDays: "Claim settlement days",
  },
  messages: {
    saved: "Insurer configuration saved.",
    saveFailed: "Unable to save configuration.",
    serviceDown: "The configuration service could not be reached.",
  },
  saveButton: {
    idle: "Save configuration",
    saving: "Saving...",
  },
  summaryPanel: {
    title: "Workspace summary",
    description: "A quick snapshot of what this insurer configuration supports.",
    currentInsurerLabel: "Current insurer",
    fallbackInsurer: "Selected insurer",
    fallbackModel: "Embedded partnership",
    alignmentTitle: "Portal alignment",
    alignmentBody:
      "This settings screen uses the same BFF pattern as the rest of the React portal, so config changes can stay inside the authenticated workspace.",
  },
} as const;
