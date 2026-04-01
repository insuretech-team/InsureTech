export type SurveyorDeskTab = "details" | "chat" | "call" | "documents";

export interface SurveyorChatEntry {
  id: string;
  author: string;
  text: string;
}

export interface SurveyorDocFormState {
  surveyReport: boolean;
  estimate: boolean;
  fir: boolean;
  photos: boolean;
  diagnostics: boolean;
  note: string;
}

type SurveyorDocToggleKey = Exclude<keyof SurveyorDocFormState, "note">;

export interface SurveyorDocRequestField {
  key: SurveyorDocToggleKey;
  label: string;
  fullWidth?: boolean;
}

export const surveyorDefaultTab: SurveyorDeskTab = "details";

export const surveyorDeskTabs: Array<{ id: SurveyorDeskTab; label: string }> = [
  { id: "details", label: "Details" },
  { id: "chat", label: "Chat" },
  { id: "call", label: "Web Call" },
  { id: "documents", label: "Request Docs" },
];

export const initialSurveyorDocForm: SurveyorDocFormState = {
  surveyReport: true,
  estimate: true,
  fir: false,
  photos: true,
  diagnostics: false,
  note: "",
};

export const surveyorDocRequestFields: SurveyorDocRequestField[] = [
  { key: "surveyReport", label: "Survey report" },
  { key: "estimate", label: "Estimate / quotation" },
  { key: "fir", label: "FIR / incident report" },
  { key: "photos", label: "Photos / videos" },
  { key: "diagnostics", label: "Diagnostic / vet report", fullWidth: true },
];

export const initialSurveyorChatLog: SurveyorChatEntry[] = [
  {
    id: "desk-seed-1",
    author: "Surveyor Desk",
    text: "Surveyor review queue initialized. Select a claim to begin assessment.",
  },
];

export const surveyorChatAck = "Claims desk acknowledged the surveyor note and linked it to the case file.";

export const surveyorCallActions = ["Start survey call", "Invite claims desk", "Add focal person"] as const;

export function buildSurveyorRequestList(form: SurveyorDocFormState) {
  return surveyorDocRequestFields.filter((field) => form[field.key]).map((field) => field.label);
}

export const surveyorDeskActors = {
  surveyor: "Surveyor",
  claimsDesk: "Claims Desk",
  surveyorDesk: "Surveyor Desk",
} as const;

export function buildSurveyorChatEntries(text: string): SurveyorChatEntry[] {
  return [
    { id: `chat-${Date.now()}`, author: surveyorDeskActors.surveyor, text },
    {
      id: `chat-reply-${Date.now() + 1}`,
      author: surveyorDeskActors.claimsDesk,
      text: surveyorChatAck,
    },
  ];
}

export function buildSurveyorRequestMessage(claimNumber: string, requested: string[], note: string) {
  const noteSuffix = note.trim() ? ` ${note.trim()}` : "";
  return `Additional documents requested for ${claimNumber}: ${requested.join(", ")}.${noteSuffix}`.trim();
}

export function getSurveyorRequestPreparedMessage(claimNumber: string) {
  return `${surveyorDeskCopy.queuePanel.queuePreparedPrefix} ${claimNumber}.`;
}

export const surveyorDeskCopy = {
  queuePanel: {
    title: "Surveyor queue",
    description: "Motor, fire/property, and pet claims that need surveyor assessment.",
    searchPlaceholder: "Search by claim number, insured, plan, or category",
    loadingLabel: "Loading surveyor queue...",
    emptyLabel: "No surveyor claims available right now.",
    errorLabel: "Unable to load surveyor claims.",
    queuePreparedPrefix: "Document request prepared for",
  },
  queueTable: ["Claim", "Category", "Surveyor", "Status"] as const,
  workspacePanel: {
    title: "Surveyor workspace",
    description: "Claim details, surveyor communication, call setup, and additional document requests.",
    emptyLabel: "Select a surveyor claim from the queue to begin review.",
  },
  fields: {
    category: "Category",
    plan: "Plan",
    incidentDate: "Incident date",
    filed: "Filed",
    requiredDocuments: "Required documents",
  },
  chat: {
    placeholder: "Add surveyor observations, inspection updates, or clarification notes...",
    sendButton: "Send message",
  },
  call: {
    roomLabel: "WebRTC verification room",
    roomTitleSuffix: "inspection call",
    roomDescription: "Use this room for video-based claim verification and inspection as described in the SRS.",
    videoPanelTitle: "Video panel",
    videoPlaceholder: "Camera preview area",
    controlsTitle: "Call controls",
  },
  documents: {
    placeholder: "Explain exactly what the surveyor needs and why the claim cannot proceed yet...",
    sendButton: "Send document request",
  },
} as const;
