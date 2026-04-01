"use client";

import { CheckCircle2, LoaderCircle, Search, ShieldX } from "lucide-react";
import Link from "next/link";
import { useState } from "react";

import { Panel } from "@/components/panel";
import { StatusPill } from "@/components/status-pill";
import { WorkspaceActionModal } from "@/components/workspace-action-modal";
import { useCurrentInsurerId } from "@/hooks/use-current-insurer-id";
import { useProposalsWorkspace } from "@/hooks/use-proposals-workspace";
import { api } from "@/lib/browser-client";
import { useLibraryDocuments } from "@/hooks/use-library-documents";
import { findPlaybook } from "@/lib/product-playbooks";
import {
  createProposalActionState,
  getProposalActionFailureMessage,
  getProposalActionModalConfig,
  getProposalActionSuccessMessage,
  initialProposalActionState,
  proposalStatusOptions as statusOptions,
  proposalsTabCopy,
  proposalsValidationCopy,
} from "@/lib/tabs/proposals";
import { formatDateTime } from "@/lib/utils";

export function ProposalsBoard() {
  const [status, setStatus] = useState("All");
  const [actionState, setActionState] = useState(initialProposalActionState);
  const { insurerId } = useCurrentInsurerId();
  const {
    selected,
    setSelectedId,
    query,
    setQuery,
    loading,
    pendingId,
    setPendingId,
    message,
    setMessage,
    refresh,
    visibleItems,
  } = useProposalsWorkspace(insurerId || undefined, status);
  const selectedPlaybook = selected ? findPlaybook(selected.planName, selected.category) : undefined;
  const documentLibrary = useLibraryDocuments();
  const selectedDocumentTemplates = selected
    ? documentLibrary.documents.filter(
        (d) =>
          d.category.toLowerCase() === (selected.category ?? "").toLowerCase() ||
          d.packId.toLowerCase().includes((selected.category ?? "").toLowerCase()),
      ).slice(0, 5)
    : [];
  const actionModal = getProposalActionModalConfig(actionState.mode);

  async function handleAction(proposalId: string, proposalNumber: string, action: "approve" | "reject") {
    const reason = actionState.note.trim();

    if (action === "reject" && !reason) {
      setMessage(proposalsValidationCopy.rejectRequired);
      return;
    }

    setPendingId(proposalId);
    setMessage("");

    try {
      const response = await api.proposals.updateStatus(proposalId, {
        action,
        reason: reason || undefined,
      });

      if (!response.ok) {
        setMessage(response.message ?? getProposalActionFailureMessage(action));
        return;
      }

      await refresh();
      setMessage(getProposalActionSuccessMessage(proposalNumber, action));
      setActionState(initialProposalActionState);
    } catch {
      setMessage(proposalsValidationCopy.saveFailed);
    } finally {
      setPendingId("");
    }
  }

  return (
    <div className="page-proposals grid gap-6 xl:grid-cols-[minmax(0,1.15fr)_380px]">
      <Panel
        title={proposalsTabCopy.queueTitle}
        description={proposalsTabCopy.queueDescription}
        action={
          <div className="flex flex-wrap gap-2">
            {statusOptions.map((option) => (
              <button
                key={option}
                className={status === option ? "portal-btn portal-btn-primary" : "portal-btn portal-btn-secondary"}
                onClick={() => setStatus(option)}
                type="button"
              >
                {option}
              </button>
            ))}
          </div>
        }
      >
        <div className="mb-5 flex items-center gap-3 rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 px-4 py-3">
          <Search className="h-4 w-4 text-[var(--muted)]" />
          <input
            className="w-full bg-transparent text-sm outline-none placeholder:text-[var(--muted)]"
            placeholder={proposalsTabCopy.searchPlaceholder}
            value={query}
            onChange={(event) => setQuery(event.target.value)}
          />
        </div>

        {message ? (
          <div className="mb-5 rounded-[22px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 px-4 py-3 text-sm text-[var(--muted)]">
            {message}
          </div>
        ) : null}

        <div className="table-wrap">
          <table className="table-base">
            <thead>
              <tr>
                {proposalsTabCopy.tableHeaders.map((header) => (
                  <th key={header}>{header}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr>
                  <td colSpan={4}>
                    <div className="flex items-center justify-center gap-2 py-8 text-sm text-[var(--muted)]">
                      <LoaderCircle className="h-4 w-4 animate-spin" />
                      {proposalsTabCopy.loadingLabel}
                    </div>
                  </td>
                </tr>
              ) : visibleItems.length ? (
                visibleItems.map((proposal) => (
                  <tr
                    key={proposal.id}
                    className={selected?.id === proposal.id ? "bg-[rgb(12_91_65_/_0.05)]" : ""}
                    onClick={() => setSelectedId(proposal.id)}
                  >
                    <td>
                      <p className="font-medium text-[var(--text)]">{proposal.customerName}</p>
                      <p className="mt-1 text-sm text-[var(--muted)]">{proposal.proposalNumber}</p>
                    </td>
                    <td>
                      <p className="font-medium text-[var(--text)]">{proposal.planName}</p>
                      <p className="mt-1 text-sm text-[var(--muted)]">{proposal.category}</p>
                    </td>
                    <td>
                      <StatusPill status={proposal.status} />
                    </td>
                    <td className="text-sm text-[var(--muted)]">{formatDateTime(proposal.submittedAt)}</td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={4}>
                    <div className="py-8 text-center text-sm text-[var(--muted)]">{proposalsTabCopy.emptyLabel}</div>
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </Panel>

      <Panel title={proposalsTabCopy.detailTitle} description={proposalsTabCopy.detailDescription}>
        {selected ? (
          <div className="space-y-5">
            <div className="rounded-[24px] bg-[rgb(12_91_65_/_0.05)] p-4">
              <div className="flex items-start justify-between gap-4">
                <div>
                  <p className="text-sm uppercase tracking-[0.14em] text-[var(--muted)]">{selected.proposalNumber}</p>
                  <h3 className="mt-2 font-[family:var(--font-heading)] text-2xl font-semibold text-[var(--text)]">
                    {selected.customerName}
                  </h3>
                </div>
                <StatusPill status={selected.status} />
              </div>
              <div className="mt-4 grid gap-3 sm:grid-cols-2">
                <div>
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{proposalsTabCopy.fields.plan}</p>
                  <p className="mt-2 text-sm text-[var(--text)]">{selected.planName}</p>
                </div>
                <div>
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{proposalsTabCopy.fields.category}</p>
                  <p className="mt-2 text-sm text-[var(--text)]">{selected.category}</p>
                </div>
                <div>
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{proposalsTabCopy.fields.coverage}</p>
                  <p className="mt-2 text-sm text-[var(--text)]">{selected.coverageText}</p>
                </div>
                <div>
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{proposalsTabCopy.fields.premium}</p>
                  <p className="mt-2 text-sm text-[var(--text)]">{selected.premiumText}</p>
                </div>
                <div>
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{proposalsTabCopy.fields.orderId}</p>
                  <p className="mt-2 text-sm text-[var(--text)]">{selected.orderId || proposalsTabCopy.fields.orderFallback}</p>
                </div>
                <div>
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{proposalsTabCopy.fields.submitted}</p>
                  <p className="mt-2 text-sm text-[var(--text)]">{formatDateTime(selected.submittedAt)}</p>
                </div>
              </div>
            </div>

            <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4">
              <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{proposalsTabCopy.fields.decisionNote}</p>
              <p className="mt-3 text-sm leading-6 text-[var(--text)]">
                {selected.decisionReason || proposalsTabCopy.fields.emptyDecisionNote}
              </p>
            </div>

            {selectedPlaybook ? (
              <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-[rgb(255_252_247_/_0.82)] p-4">
                <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{proposalsTabCopy.fields.reviewGuide}</p>
                <h4 className="mt-2 font-[family:var(--font-heading)] text-xl font-semibold text-[var(--text)]">
                  {selectedPlaybook.name}
                </h4>
                <p className="mt-2 text-sm leading-6 text-[var(--muted)]">{selectedPlaybook.summary}</p>
                <div className="mt-4 grid gap-3 sm:grid-cols-2">
                  <div className="rounded-[20px] bg-white p-3">
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{proposalsTabCopy.fields.audience}</p>
                    <p className="mt-2 text-sm text-[var(--text)]">{selectedPlaybook.audience}</p>
                  </div>
                  <div className="rounded-[20px] bg-white p-3">
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{proposalsTabCopy.fields.eligibility}</p>
                    <p className="mt-2 text-sm text-[var(--text)]">{selectedPlaybook.ageRange}</p>
                  </div>
                </div>
                <div className="mt-4">
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{proposalsTabCopy.fields.requiredDocs}</p>
                  <div className="mt-2 flex flex-wrap gap-2">
                    {selectedPlaybook.requiredDocuments.map((doc) => (
                      <span
                        key={doc}
                        className="rounded-full border border-[rgb(12_91_65_/_0.08)] bg-white px-3 py-1 text-xs font-medium text-[var(--muted)]"
                      >
                        {doc}
                      </span>
                    ))}
                  </div>
                </div>
                <div className="mt-4">
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{proposalsTabCopy.fields.flags}</p>
                  <div className="mt-2 flex flex-wrap gap-2">
                    {selectedPlaybook.operationalFlags.map((flag) => (
                      <span
                        key={flag}
                        className="rounded-full bg-[rgb(245_158_11_/_0.12)] px-3 py-1 text-xs font-semibold text-[#8a5200]"
                      >
                        {flag}
                      </span>
                    ))}
                  </div>
                </div>
              </div>
            ) : null}

            {selectedDocumentTemplates.length > 0 && (
              <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4">
                <div className="flex items-start justify-between gap-3">
                  <div>
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{proposalsTabCopy.fields.documentsPack}</p>
                    <h4 className="mt-2 font-[family:var(--font-heading)] text-xl font-semibold text-[var(--text)]">
                      Related Documents
                    </h4>
                  </div>
                  <Link className="portal-btn portal-btn-secondary" href="/documents">
                    {proposalsTabCopy.fields.documentsAction}
                  </Link>
                </div>
                <div className="mt-4 space-y-2">
                  {selectedDocumentTemplates.map((doc) => (
                    <div
                      key={doc.id}
                      className="rounded-[18px] border border-[rgb(12_91_65_/_0.08)] bg-[rgb(255_252_247_/_0.82)] px-3 py-3"
                    >
                      <p className="font-medium text-[var(--text)]">{doc.title}</p>
                      <p className="mt-1 text-sm text-[var(--muted)]">{doc.stage}</p>
                    </div>
                  ))}
                </div>
              </div>
            )}

            <div className="grid gap-3 sm:grid-cols-2">
              <button
                className="portal-btn portal-btn-primary"
                disabled={pendingId === selected.id}
                onClick={() =>
                  setActionState(createProposalActionState("approve", selected))
                }
                type="button"
              >
                {pendingId === selected.id ? <LoaderCircle className="h-4 w-4 animate-spin" /> : <CheckCircle2 className="h-4 w-4" />}
                {proposalsTabCopy.actionButtons.approve}
              </button>
              <button
                className="portal-btn border border-[rgb(194_65_12_/_0.14)] bg-[var(--danger-soft)] text-[var(--danger)]"
                disabled={pendingId === selected.id}
                onClick={() =>
                  setActionState(createProposalActionState("reject", selected))
                }
                type="button"
              >
                {pendingId === selected.id ? <LoaderCircle className="h-4 w-4 animate-spin" /> : <ShieldX className="h-4 w-4" />}
                {proposalsTabCopy.actionButtons.reject}
              </button>
            </div>
          </div>
        ) : (
          <p className="text-sm text-[var(--muted)]">{proposalsTabCopy.actionButtons.emptyLabel}</p>
        )}
      </Panel>

      <WorkspaceActionModal
        open={Boolean(actionState.mode && actionState.proposalId)}
        title={actionModal.title}
        description={actionModal.description}
        fields={actionModal.fields}
        values={{ note: actionState.note }}
        submitLabel={actionModal.submitLabel}
        closeLabel={proposalsTabCopy.modalButtons.close}
        cancelLabel={proposalsTabCopy.modalButtons.cancel}
        onChange={(key, value) => setActionState((current) => ({ ...current, [key]: value }))}
        onClose={() => setActionState(initialProposalActionState)}
        onSubmit={() => {
          if (actionState.proposalId && actionState.proposalNumber && actionState.mode) {
            void handleAction(actionState.proposalId, actionState.proposalNumber, actionState.mode);
          }
        }}
      />
    </div>
  );
}
