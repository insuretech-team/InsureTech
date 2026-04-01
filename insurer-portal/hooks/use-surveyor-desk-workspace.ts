"use client";

import { useMemo, useState } from "react";

import { useClaimsWorkspace } from "@/hooks/use-claims-workspace";
import { buildSurveyorReview, findClaimCategoryMatrix } from "@/lib/claims-intelligence";
import {
  buildSurveyorChatEntries,
  buildSurveyorRequestList,
  buildSurveyorRequestMessage,
  getSurveyorRequestPreparedMessage,
  initialSurveyorChatLog,
  initialSurveyorDocForm,
  surveyorDeskActors,
  surveyorDefaultTab,
  type SurveyorDeskTab,
} from "@/lib/tabs/surveyor-desk";

export function useSurveyorDeskWorkspace(insurerId?: string) {
  const claims = useClaimsWorkspace(insurerId, "All", "surveyor-only");
  const [activeTab, setActiveTab] = useState<SurveyorDeskTab>(surveyorDefaultTab);
  const [chatDraft, setChatDraft] = useState("");
  const [docForm, setDocForm] = useState(initialSurveyorDocForm);
  const [chatLog, setChatLog] = useState(initialSurveyorChatLog);

  const selectedMatrix = useMemo(
    () => (claims.selected ? findClaimCategoryMatrix(claims.selected.category, claims.selected.planName) : undefined),
    [claims.selected],
  );
  const surveyor = useMemo(
    () => (claims.selected ? buildSurveyorReview(claims.selected) : null),
    [claims.selected],
  );
  const queueItems = useMemo(
    () =>
      claims.visibleItems.map((claim) => ({
        claim,
        review: buildSurveyorReview(claim),
      })),
    [claims.visibleItems],
  );

  function openTab(tab: SurveyorDeskTab) {
    setActiveTab(tab);
  }

  function sendChat() {
    if (!chatDraft.trim() || !claims.selected) return;
    setChatLog((current) => [...current, ...buildSurveyorChatEntries(chatDraft.trim())]);
    setChatDraft("");
  }

  function submitDocRequest() {
    if (!claims.selected) return;

    const requested = buildSurveyorRequestList(docForm);
    setChatLog((current) => [
      ...current,
      {
        id: `docs-${Date.now()}`,
        author: surveyorDeskActors.surveyorDesk,
        text: buildSurveyorRequestMessage(claims.selected.claimNumber, requested, docForm.note),
      },
    ]);
    claims.setMessage(getSurveyorRequestPreparedMessage(claims.selected.claimNumber));
    setActiveTab("chat");
  }

  return {
    ...claims,
    activeTab,
    openTab,
    chatDraft,
    setChatDraft,
    docForm,
    setDocForm,
    chatLog,
    queueItems,
    selectedMatrix,
    surveyor,
    sendChat,
    submitDocRequest,
  };
}
