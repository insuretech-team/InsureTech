"use client";

import { BadgeDollarSign, CheckCircle2, FileWarning, LoaderCircle, Search, ShieldX } from "lucide-react";
import Link from "next/link";
import { useState } from "react";

import { Panel } from "@/components/panel";
import { StatusPill } from "@/components/status-pill";
import { WorkspaceActionModal } from "@/components/workspace-action-modal";
import { useClaimsWorkspace } from "@/hooks/use-claims-workspace";
import { useCurrentInsurerId } from "@/hooks/use-current-insurer-id";
import { api } from "@/lib/browser-client";
import { findClaimCategoryMatrix, routeApprovalTier } from "@/lib/claims-intelligence";
import {
  claimStatusOptions as statusOptions,
  claimsTabCopy,
  createClaimActionState,
  extractClaimAmount,
  getClaimActionFailureMessage,
  getClaimActionModalConfig,
  getClaimActionSuccessMessage,
  initialClaimActionState,
} from "@/lib/tabs/claims";
import { findPlaybook } from "@/lib/product-playbooks";
import { formatDate, formatDateTime } from "@/lib/utils";

export function ClaimsBoard() {
  const [status, setStatus] = useState("All");
  const [actionState, setActionState] = useState(initialClaimActionState);
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
  } = useClaimsWorkspace(insurerId || undefined, status);

  const selectedPlaybook = selected ? findPlaybook(selected.planName, selected.category) : undefined;
  const selectedMatrix = selected ? findClaimCategoryMatrix(selected.category, selected.planName) : undefined;
  const selectedTier = selected ? routeApprovalTier(extractClaimAmount(selected.claimedAmountText)) : undefined;
  const actionModal = getClaimActionModalConfig(actionState.mode);

  async function handleAction(claimId: string, claimNumber: string, action: "approve" | "reject" | "settle") {
    const reason = actionState.reason.trim();
    let amount: number | undefined;
    const paymentReference = actionState.paymentReference.trim() || undefined;

    if (action === "reject" && !reason) {
      setMessage(claimsTabCopy.validation.rejectRequired);
      return;
    }

    if (action === "settle") {
      const parsed = Number(actionState.amount);
      if (!Number.isFinite(parsed) || parsed <= 0) {
        setMessage(claimsTabCopy.validation.settleAmount);
        return;
      }
      amount = parsed;
    }

    setPendingId(claimId);
    setMessage("");

    try {
      const response = await api.claims.updateStatus(claimId, {
        action,
        reason: reason || undefined,
        amount,
        paymentReference,
      });

      if (!response.ok) {
        setMessage(response.message ?? getClaimActionFailureMessage(action));
        return;
      }

      await refresh();
      setMessage(getClaimActionSuccessMessage(claimNumber));
      setActionState(initialClaimActionState);
    } catch {
      setMessage(claimsTabCopy.validation.saveFailed);
    } finally {
      setPendingId("");
    }
  }

  return (
    <div className="page-claims grid gap-6 xl:grid-cols-[minmax(0,1.15fr)_380px]">
      <Panel title={claimsTabCopy.queueTitle} description={claimsTabCopy.queueDescription}>
        <div className="mb-5 grid grid-cols-2 gap-2 md:grid-cols-3 xl:grid-cols-6">
          {statusOptions.map((option) => (
            <button
              key={option}
              className={`${status === option ? "portal-btn portal-btn-primary" : "portal-btn portal-btn-secondary"} w-full`}
              onClick={() => setStatus(option)}
              type="button"
            >
              {option}
            </button>
          ))}
        </div>

        <div className="mb-5 flex items-center gap-3 rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 px-4 py-3">
          <Search className="h-4 w-4 text-[var(--muted)]" />
          <input
            className="w-full bg-transparent text-sm outline-none placeholder:text-[var(--muted)]"
            placeholder={claimsTabCopy.searchPlaceholder}
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
                {claimsTabCopy.tableHeaders.map((header) => (
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
                      {claimsTabCopy.loadingLabel}
                    </div>
                  </td>
                </tr>
              ) : visibleItems.length ? (
                visibleItems.map((claim) => (
                  <tr
                    key={claim.id}
                    className={selected?.id === claim.id ? "bg-[rgb(12_91_65_/_0.05)]" : ""}
                    onClick={() => setSelectedId(claim.id)}
                  >
                    <td>
                      <p className="font-medium text-[var(--text)]">{claim.insuredName}</p>
                      <p className="mt-1 text-sm text-[var(--muted)]">{claim.claimNumber}</p>
                    </td>
                    <td>
                      <p className="font-medium text-[var(--text)]">{claim.planName}</p>
                      <p className="mt-1 text-sm text-[var(--muted)]">{claim.category}</p>
                    </td>
                    <td>
                      <StatusPill status={claim.status} />
                    </td>
                    <td className="text-sm text-[var(--muted)]">{claim.claimedAmountText}</td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={4}>
                    <div className="py-8 text-center text-sm text-[var(--muted)]">{claimsTabCopy.emptyLabel}</div>
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </Panel>

      <Panel title={claimsTabCopy.detailTitle} description={claimsTabCopy.detailDescription}>
        {selected ? (
          <div className="space-y-5">
            <div className="rounded-[24px] bg-[rgb(12_91_65_/_0.05)] p-4">
              <div className="flex items-start justify-between gap-4">
                <div>
                  <p className="text-sm uppercase tracking-[0.14em] text-[var(--muted)]">{selected.claimNumber}</p>
                  <h3 className="mt-2 font-[family:var(--font-heading)] text-2xl font-semibold text-[var(--text)]">
                    {selected.insuredName}
                  </h3>
                </div>
                <StatusPill status={selected.status} />
              </div>
              <div className="mt-4 grid gap-3 sm:grid-cols-2">
                <div>
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{claimsTabCopy.fields.filed}</p>
                  <p className="mt-2 text-sm text-[var(--text)]">{formatDateTime(selected.submittedAt)}</p>
                </div>
                <div>
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{claimsTabCopy.fields.incidentDate}</p>
                  <p className="mt-2 text-sm text-[var(--text)]">{formatDate(selected.incidentDate)}</p>
                </div>
                <div>
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{claimsTabCopy.fields.claimedAmount}</p>
                  <p className="mt-2 text-sm text-[var(--text)]">{selected.claimedAmountText}</p>
                </div>
                <div>
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{claimsTabCopy.fields.approvedAmount}</p>
                  <p className="mt-2 text-sm text-[var(--text)]">{selected.approvedAmountText}</p>
                </div>
                <div>
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{claimsTabCopy.fields.settledAmount}</p>
                  <p className="mt-2 text-sm text-[var(--text)]">{selected.settledAmountText}</p>
                </div>
                <div>
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{claimsTabCopy.fields.category}</p>
                  <p className="mt-2 text-sm text-[var(--text)]">{selected.category}</p>
                </div>
              </div>
            </div>

            <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4">
              <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{claimsTabCopy.fields.claimReason}</p>
              <p className="mt-3 text-sm leading-6 text-[var(--text)]">{selected.reason || claimsTabCopy.fields.emptyReason}</p>
            </div>

            {selectedPlaybook ? (
              <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-[rgb(255_252_247_/_0.82)] p-4">
                <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{claimsTabCopy.fields.guideTitle}</p>
                <h4 className="mt-2 font-[family:var(--font-heading)] text-xl font-semibold text-[var(--text)]">
                  {selectedPlaybook.name}
                </h4>
                <div className="mt-4">
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{claimsTabCopy.fields.requiredDocs}</p>
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
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{claimsTabCopy.fields.exclusions}</p>
                  <div className="mt-2 flex flex-wrap gap-2">
                    {selectedPlaybook.exclusions.map((item) => (
                      <span key={item} className="rounded-full bg-[var(--danger-soft)] px-3 py-1 text-xs font-semibold text-[var(--danger)]">
                        {item}
                      </span>
                    ))}
                  </div>
                </div>
                <div className="mt-4">
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{claimsTabCopy.fields.flow}</p>
                  <ol className="mt-2 grid gap-2 text-sm text-[var(--muted)]">
                    {selectedPlaybook.claimSteps.map((step) => (
                      <li key={step} className="rounded-[18px] bg-white px-3 py-2">
                        {step}
                      </li>
                    ))}
                  </ol>
                </div>
              </div>
            ) : null}

            {selectedMatrix && selectedTier ? (
              <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4">
                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="rounded-[20px] bg-[rgb(12_91_65_/_0.05)] p-3">
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{claimsTabCopy.fields.claimMode}</p>
                    <p className="mt-2 text-sm font-medium text-[var(--text)]">{selectedMatrix.claimMode}</p>
                    <p className="mt-1 text-xs text-[var(--muted)]">{selectedMatrix.typicalTat}</p>
                  </div>
                  <div className="rounded-[20px] bg-[rgb(245_158_11_/_0.1)] p-3">
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{claimsTabCopy.fields.approvalTier}</p>
                    <p className="mt-2 text-sm font-medium text-[var(--text)]">{selectedTier.approvalLevel}</p>
                    <p className="mt-1 text-xs text-[var(--muted)]">
                      {selectedTier.approvers} • {selectedTier.maxTat}
                    </p>
                  </div>
                </div>

                <div className="mt-4">
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{claimsTabCopy.fields.touchpoints}</p>
                  <div className="mt-2 flex flex-wrap gap-2">
                    {selectedMatrix.partnerTouchpoints.map((item) => (
                      <span
                        key={item}
                        className="rounded-full border border-[rgb(12_91_65_/_0.08)] bg-[rgb(12_91_65_/_0.05)] px-3 py-1 text-xs font-medium text-[var(--muted)]"
                      >
                        {item}
                      </span>
                    ))}
                  </div>
                </div>
              </div>
            ) : null}

            {selectedMatrix?.surveyorRequired ? (
              <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-[rgb(255_252_247_/_0.82)] p-4">
                <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{claimsTabCopy.fields.surveyorTitle}</p>
                <p className="mt-3 text-sm leading-6 text-[var(--muted)]">
                  {claimsTabCopy.fields.surveyorBody}
                </p>
                <div className="mt-4">
                  <Link className="portal-btn portal-btn-secondary" href="/surveyor-desk">
                    {claimsTabCopy.fields.surveyorAction}
                  </Link>
                </div>
              </div>
            ) : null}

            <div className="grid gap-3">
              <button
                className="portal-btn portal-btn-primary"
                disabled={pendingId === selected.id}
                onClick={() =>
                  setActionState(createClaimActionState("approve", selected))
                }
                type="button"
              >
                {pendingId === selected.id ? <LoaderCircle className="h-4 w-4 animate-spin" /> : <CheckCircle2 className="h-4 w-4" />}
                {claimsTabCopy.buttons.approve}
              </button>
              <button
                className="portal-btn border border-[rgb(245_158_11_/_0.2)] bg-[rgb(245_158_11_/_0.12)] text-[#8a5200]"
                disabled={pendingId === selected.id}
                onClick={() =>
                  setActionState(createClaimActionState("settle", selected))
                }
                type="button"
              >
                {pendingId === selected.id ? <LoaderCircle className="h-4 w-4 animate-spin" /> : <BadgeDollarSign className="h-4 w-4" />}
                {claimsTabCopy.buttons.settle}
              </button>
              <button
                className="portal-btn border border-[rgb(194_65_12_/_0.14)] bg-[var(--danger-soft)] text-[var(--danger)]"
                disabled={pendingId === selected.id}
                onClick={() =>
                  setActionState(createClaimActionState("reject", selected))
                }
                type="button"
              >
                {pendingId === selected.id ? <LoaderCircle className="h-4 w-4 animate-spin" /> : <ShieldX className="h-4 w-4" />}
                {claimsTabCopy.buttons.reject}
              </button>
            </div>

            <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-[rgb(255_252_247_/_0.82)] p-4">
              <div className="flex items-start gap-3">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-2xl bg-[rgb(245_158_11_/_0.14)] text-[#8a5200]">
                  <FileWarning className="h-5 w-5" />
                </div>
                <div>
                  <p className="font-medium text-[var(--text)]">{claimsTabCopy.fields.settlementNoteTitle}</p>
                  <p className="mt-1 text-sm leading-6 text-[var(--muted)]">
                    {claimsTabCopy.fields.settlementNoteBody}
                  </p>
                </div>
              </div>
            </div>
          </div>
        ) : (
          <p className="text-sm text-[var(--muted)]">{claimsTabCopy.fields.emptyLabel}</p>
        )}
      </Panel>

      <WorkspaceActionModal
        open={Boolean(actionState.mode && actionState.claimId)}
        title={actionModal.title}
        description={actionModal.description}
        fields={actionModal.fields}
        values={{
          reason: actionState.reason,
          amount: actionState.amount,
          paymentReference: actionState.paymentReference,
        }}
        submitLabel={actionModal.submitLabel}
        closeLabel={claimsTabCopy.modalButtons.close}
        cancelLabel={claimsTabCopy.modalButtons.cancel}
        onChange={(key, value) => setActionState((current) => ({ ...current, [key]: value }))}
        onClose={() => setActionState(initialClaimActionState)}
        onSubmit={() => {
          if (actionState.claimId && actionState.claimNumber && actionState.mode) {
            void handleAction(actionState.claimId, actionState.claimNumber, actionState.mode);
          }
        }}
      />
    </div>
  );
}
