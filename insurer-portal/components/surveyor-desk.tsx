"use client";

import { CheckCheck, FilePlus2, LoaderCircle, PhoneCall, Search, UserRoundSearch, Video } from "lucide-react";
import { Panel } from "@/components/panel";
import { StatusPill } from "@/components/status-pill";
import { useCurrentInsurerId } from "@/hooks/use-current-insurer-id";
import { useSurveyorDeskWorkspace } from "@/hooks/use-surveyor-desk-workspace";
import {
  surveyorCallActions,
  surveyorDeskCopy,
  surveyorDeskTabs,
  surveyorDocRequestFields,
} from "@/lib/tabs/surveyor-desk";
import { formatDate, formatDateTime } from "@/lib/utils";

export function SurveyorDesk() {
  const { insurerId } = useCurrentInsurerId();
  const {
    selected,
    setSelectedId,
    query,
    setQuery,
    loading,
    message,
    queueItems,
    activeTab,
    openTab,
    chatDraft,
    setChatDraft,
    docForm,
    setDocForm,
    chatLog,
    selectedMatrix,
    surveyor,
    sendChat,
    submitDocRequest,
  } = useSurveyorDeskWorkspace(insurerId || undefined);

  return (
    <div className="page-surveyor grid gap-6 xl:grid-cols-[minmax(0,1.1fr)_420px]">
      <Panel title={surveyorDeskCopy.queuePanel.title} description={surveyorDeskCopy.queuePanel.description}>
        <div className="mb-5 flex items-center gap-3 rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 px-4 py-3">
          <Search className="h-4 w-4 text-[var(--muted)]" />
          <input
            className="w-full bg-transparent text-sm outline-none placeholder:text-[var(--muted)]"
            placeholder={surveyorDeskCopy.queuePanel.searchPlaceholder}
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
                {surveyorDeskCopy.queueTable.map((header) => (
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
                      {surveyorDeskCopy.queuePanel.loadingLabel}
                    </div>
                  </td>
                </tr>
              ) : queueItems.length ? (
                queueItems.map(({ claim, review }) => {
                  return (
                    <tr
                      key={claim.id}
                      className={selected?.id === claim.id ? "bg-[rgb(12_91_65_/_0.05)]" : ""}
                      onClick={() => setSelectedId(claim.id)}
                    >
                      <td>
                        <p className="font-medium text-[var(--text)]">{claim.insuredName}</p>
                        <p className="mt-1 text-sm text-[var(--muted)]">{claim.claimNumber}</p>
                      </td>
                      <td className="text-sm text-[var(--muted)]">
                        {claim.category}
                        <div className="mt-1">{claim.planName}</div>
                      </td>
                      <td className="text-sm text-[var(--muted)]">
                        {review.name}
                        <div className="mt-1">{review.title}</div>
                      </td>
                      <td>
                        <StatusPill status={claim.status} />
                      </td>
                    </tr>
                  );
                })
              ) : (
                <tr>
                  <td colSpan={4}>
                    <div className="py-8 text-center text-sm text-[var(--muted)]">{surveyorDeskCopy.queuePanel.emptyLabel}</div>
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </Panel>

      <Panel title={surveyorDeskCopy.workspacePanel.title} description={surveyorDeskCopy.workspacePanel.description}>
        {selected && surveyor && selectedMatrix ? (
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
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{surveyorDeskCopy.fields.category}</p>
                  <p className="mt-2 text-sm text-[var(--text)]">{selected.category}</p>
                </div>
                <div>
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{surveyorDeskCopy.fields.plan}</p>
                  <p className="mt-2 text-sm text-[var(--text)]">{selected.planName}</p>
                </div>
                <div>
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{surveyorDeskCopy.fields.incidentDate}</p>
                  <p className="mt-2 text-sm text-[var(--text)]">{formatDate(selected.incidentDate)}</p>
                </div>
                <div>
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{surveyorDeskCopy.fields.filed}</p>
                  <p className="mt-2 text-sm text-[var(--text)]">{formatDateTime(selected.submittedAt)}</p>
                </div>
              </div>
            </div>

            <div className="grid grid-cols-2 gap-2 sm:grid-cols-4">
              {surveyorDeskTabs.map((tab) => (
                <button
                  key={tab.id}
                  className={activeTab === tab.id ? "portal-btn portal-btn-primary" : "portal-btn portal-btn-secondary"}
                  onClick={() => openTab(tab.id)}
                  type="button"
                >
                  {tab.label}
                </button>
              ))}
            </div>

            {activeTab === "details" ? (
              <div className="space-y-4">
                <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4">
                  <div className="flex items-start gap-3">
                    <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-2xl bg-[rgb(15_157_104_/_0.12)] text-[var(--brand-deep)]">
                      <UserRoundSearch className="h-5 w-5" />
                    </div>
                    <div>
                      <p className="font-medium text-[var(--text)]">{surveyor.name}</p>
                      <p className="mt-1 text-sm text-[var(--muted)]">
                        {surveyor.title} • {surveyor.status}
                      </p>
                    </div>
                  </div>
                </div>

                <div className="grid gap-3">
                  {surveyor.notes.map((note) => (
                    <div key={note} className="rounded-[20px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 px-4 py-3 text-sm text-[var(--muted)]">
                      {note}
                    </div>
                  ))}
                </div>

                <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-[rgb(255_252_247_/_0.82)] p-4">
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{surveyorDeskCopy.fields.requiredDocuments}</p>
                  <div className="mt-3 flex flex-wrap gap-2">
                    {selectedMatrix.primaryDocuments.map((item) => (
                      <span
                        key={item}
                        className="rounded-full border border-[rgb(12_91_65_/_0.08)] bg-white px-3 py-1 text-xs font-medium text-[var(--muted)]"
                      >
                        {item}
                      </span>
                    ))}
                  </div>
                </div>
              </div>
            ) : null}

            {activeTab === "chat" ? (
              <div className="space-y-4">
                <div className="max-h-72 space-y-3 overflow-y-auto rounded-[24px] bg-[rgb(12_91_65_/_0.04)] p-3">
                  {chatLog.map((entry) => (
                    <div key={entry.id} className="rounded-[18px] bg-white px-3 py-3">
                      <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--muted)]">{entry.author}</p>
                      <p className="mt-2 text-sm leading-6 text-[var(--text)]">{entry.text}</p>
                    </div>
                  ))}
                </div>
                <textarea
                  className="portal-textarea"
                  placeholder={surveyorDeskCopy.chat.placeholder}
                  value={chatDraft}
                  onChange={(event) => setChatDraft(event.target.value)}
                />
                <button className="portal-btn portal-btn-primary" onClick={sendChat} type="button">
                  <CheckCheck className="h-4 w-4" />
                  {surveyorDeskCopy.chat.sendButton}
                </button>
              </div>
            ) : null}

            {activeTab === "call" ? (
              <div className="space-y-4">
                <div className="rounded-[24px] bg-[rgb(12_91_65_/_0.92)] p-5 text-white">
                  <p className="text-sm uppercase tracking-[0.16em] text-white/72">{surveyorDeskCopy.call.roomLabel}</p>
                  <p className="mt-3 font-[family:var(--font-heading)] text-2xl font-semibold">
                    {selected.claimNumber} {surveyorDeskCopy.call.roomTitleSuffix}
                  </p>
                  <p className="mt-2 text-sm text-white/78">{surveyorDeskCopy.call.roomDescription}</p>
                </div>
                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4">
                    <div className="flex items-center gap-3">
                      <Video className="h-5 w-5 text-[var(--brand-deep)]" />
                      <p className="font-medium text-[var(--text)]">{surveyorDeskCopy.call.videoPanelTitle}</p>
                    </div>
                    <div className="mt-4 flex min-h-[180px] items-center justify-center rounded-[20px] bg-[rgb(12_91_65_/_0.06)] text-sm text-[var(--muted)]">
                      {surveyorDeskCopy.call.videoPlaceholder}
                    </div>
                  </div>
                  <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4">
                    <div className="flex items-center gap-3">
                      <PhoneCall className="h-5 w-5 text-[var(--brand-deep)]" />
                      <p className="font-medium text-[var(--text)]">{surveyorDeskCopy.call.controlsTitle}</p>
                    </div>
                    <div className="mt-4 grid gap-3">
                      {surveyorCallActions.map((action, index) => (
                        <button
                          key={action}
                          className={index === 0 ? "portal-btn portal-btn-primary" : "portal-btn portal-btn-secondary"}
                          type="button"
                        >
                          {action}
                        </button>
                      ))}
                    </div>
                  </div>
                </div>
              </div>
            ) : null}

            {activeTab === "documents" ? (
              <div className="space-y-4">
                <div className="grid gap-3 sm:grid-cols-2">
                  {surveyorDocRequestFields.map((field) => (
                    <label
                      key={field.key}
                      className={`flex items-center gap-3 rounded-[20px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4 ${
                        field.fullWidth ? "sm:col-span-2" : ""
                      }`}
                    >
                      <input
                        className="h-4 w-4 accent-[var(--brand)]"
                        checked={docForm[field.key]}
                        onChange={(event) =>
                          setDocForm((current) => ({ ...current, [field.key]: event.target.checked }))
                        }
                        type="checkbox"
                      />
                      <span className="text-sm text-[var(--text)]">{field.label}</span>
                    </label>
                  ))}
                </div>

                <textarea
                  className="portal-textarea"
                  placeholder={surveyorDeskCopy.documents.placeholder}
                  value={docForm.note}
                  onChange={(event) => setDocForm((current) => ({ ...current, note: event.target.value }))}
                />

                <button className="portal-btn portal-btn-primary" onClick={submitDocRequest} type="button">
                  <FilePlus2 className="h-4 w-4" />
                  {surveyorDeskCopy.documents.sendButton}
                </button>
              </div>
            ) : null}
          </div>
        ) : (
          <p className="text-sm text-[var(--muted)]">{surveyorDeskCopy.workspacePanel.emptyLabel}</p>
        )}
      </Panel>
    </div>
  );
}
