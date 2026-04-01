"use client";

import { AlertTriangle, ClipboardList, Eye, ShieldCheck, X } from "lucide-react";
import Link from "next/link";
import { useMemo, useState } from "react";

import { Panel } from "@/components/panel";
import { StatusPill } from "@/components/status-pill";
import { WorkspaceIconCardGrid } from "@/components/workspace-icon-card-grid";
import { WorkspaceNoteList } from "@/components/workspace-note-list";
import { getClaimChecklistsTabData } from "@/lib/tabs/claim-checklists";

export function ClaimsChecklistWorkspace() {
  const workspace = useMemo(() => getClaimChecklistsTabData(), []);
  const [activeChecklistId, setActiveChecklistId] = useState("");
  const activeChecklist = workspace.checklists.find((item) => item.id === activeChecklistId) ?? null;

  function openChecklist(checklistId: string) {
    setActiveChecklistId(checklistId);
  }

  function closeChecklist() {
    setActiveChecklistId("");
  }

  return (
    <div className="page-claim-checklists space-y-6">
      <section className="grid gap-4 xl:grid-cols-[minmax(0,1.2fr)_minmax(320px,0.8fr)]">
        <Panel className="overflow-hidden">
          <div className="grid gap-5 md:grid-cols-[minmax(0,1fr)_260px]">
            <div>
              <p className="text-xs font-semibold uppercase tracking-[0.18em] text-[var(--brand-deep)]">
                {workspace.copy.hero.eyebrow}
              </p>
              <h2 className="mt-3 font-[family:var(--font-heading)] text-3xl font-semibold text-[var(--text)]">
                {workspace.copy.hero.title}
              </h2>
              <p className="mt-3 max-w-2xl text-sm leading-7 text-[var(--muted)]">
                {workspace.copy.hero.description}
              </p>
              <div className="mt-5 flex flex-wrap gap-3">
                <Link className="portal-btn portal-btn-primary" href="/claim-settlement">
                  {workspace.copy.hero.primaryAction}
                </Link>
                <Link className="portal-btn portal-btn-secondary" href="/documents">
                  {workspace.copy.hero.secondaryAction}
                </Link>
              </div>
            </div>

            <div className="grid gap-3">
              <div className="rounded-[24px] bg-[rgb(12_91_65_/_0.92)] p-5 text-white">
                <p className="text-sm uppercase tracking-[0.16em] text-white/72">{workspace.copy.hero.checklistLabel}</p>
                <p className="mt-3 text-4xl font-semibold">{workspace.checklists.length}</p>
                <p className="mt-2 text-sm text-white/78">{workspace.copy.hero.checklistDescription}</p>
              </div>
              <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/76 p-5">
                <p className="text-sm uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.hero.surveyorLabel}</p>
                <p className="mt-3 text-4xl font-semibold text-[var(--text)]">
                  {workspace.checklists.filter((item) => item.surveyorRequired).length}
                </p>
                <p className="mt-2 text-sm text-[var(--muted)]">{workspace.copy.hero.surveyorDescription}</p>
              </div>
            </div>
          </div>
        </Panel>

        <Panel title={workspace.copy.overviewPanel.title} description={workspace.copy.overviewPanel.description}>
          <WorkspaceNoteList items={workspace.overviewItems} />
        </Panel>
      </section>

      <Panel title={workspace.copy.lanesPanel.title} description={workspace.copy.lanesPanel.description}>
        <div className="checklist-card-grid">
          {workspace.checklists.map((item) => (
            <button key={item.id} className="checklist-lane-card" onClick={() => openChecklist(item.id)} type="button">
              <div className="flex items-start justify-between gap-3">
                <div>
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{item.category}</p>
                  <h3 className="mt-2 font-[family:var(--font-heading)] text-xl font-semibold text-[var(--text)]">
                    {item.title}
                  </h3>
                </div>
                <StatusPill status={item.surveyorRequired ? workspace.copy.lanesPanel.surveyorStatus : workspace.copy.lanesPanel.deskStatus} />
              </div>

              <p className="mt-3 text-sm leading-6 text-[var(--muted)]">{item.responseWindow}</p>

              <div className="mt-4 grid gap-3 sm:grid-cols-2">
                <div className="rounded-[18px] bg-[rgb(12_91_65_/_0.05)] p-3">
                  <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--muted)]">{workspace.copy.lanesPanel.docsLabel}</p>
                  <p className="mt-2 text-sm font-medium text-[var(--text)]">{item.documents.length}</p>
                </div>
                <div className="rounded-[18px] bg-[rgb(245_158_11_/_0.1)] p-3">
                  <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--muted)]">{workspace.copy.lanesPanel.blockersLabel}</p>
                  <p className="mt-2 text-sm font-medium text-[var(--text)]">{item.deficiencySignals.length}</p>
                </div>
              </div>

              <div className="mt-4 flex items-center justify-between text-sm">
                <span className="text-[var(--muted)]">{item.sourceFiles.length} {workspace.copy.lanesPanel.sourceSuffix}</span>
                <span className="inline-flex items-center gap-2 font-medium text-[var(--brand-deep)]">
                  <Eye className="h-4 w-4" />
                  {workspace.copy.lanesPanel.openLabel}
                </span>
              </div>
            </button>
          ))}
        </div>
      </Panel>

      <Panel title={workspace.copy.stancePanel.title} description={workspace.copy.stancePanel.description}>
        <WorkspaceIconCardGrid items={workspace.stanceCards} />
      </Panel>

      {activeChecklist ? (
        <div className="workspace-modal-backdrop" data-workspace-modal="true">
          <div className="workspace-modal-shell">
            <div className="workspace-modal-header">
              <div>
                <p className="text-xs font-semibold uppercase tracking-[0.18em] text-[var(--muted)]">
                  {activeChecklist.category}
                </p>
                <h2 className="mt-2 font-[family:var(--font-heading)] text-3xl font-semibold text-[var(--text)]">
                  {activeChecklist.title}
                </h2>
                <p className="mt-2 text-sm leading-6 text-[var(--muted)]">{activeChecklist.responseWindow}</p>
              </div>
              <div className="flex flex-wrap gap-2">
                {activeChecklist.surveyorRequired ? (
                  <Link className="portal-btn portal-btn-secondary" href="/surveyor-desk">
                    {workspace.copy.modal.surveyorAction}
                  </Link>
                ) : (
                  <Link className="portal-btn portal-btn-secondary" href="/claim-settlement">
                    {workspace.copy.modal.settlementAction}
                  </Link>
                )}
                <button className="portal-btn portal-btn-secondary" onClick={closeChecklist} type="button">
                  <X className="h-4 w-4" />
                  {workspace.copy.modal.closeButton}
                </button>
              </div>
            </div>

            <div className="workspace-modal-body">
              <main className="space-y-5">
                <div className="grid gap-3 sm:grid-cols-3">
                  <div className="rounded-[22px] bg-[rgb(12_91_65_/_0.06)] p-4">
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.modal.docsLabel}</p>
                    <p className="mt-2 text-3xl font-semibold text-[var(--text)]">{activeChecklist.documents.length}</p>
                  </div>
                  <div className="rounded-[22px] bg-[rgb(245_158_11_/_0.1)] p-4">
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.modal.blockersLabel}</p>
                    <p className="mt-2 text-3xl font-semibold text-[var(--text)]">{activeChecklist.deficiencySignals.length}</p>
                  </div>
                  <div className="rounded-[22px] bg-[rgb(14_165_233_/_0.1)] p-4">
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.modal.reviewModeLabel}</p>
                    <p className="mt-2 text-sm font-semibold text-[var(--text)]">
                      {activeChecklist.surveyorRequired ? workspace.copy.modal.surveyorMode : workspace.copy.modal.deskMode}
                    </p>
                  </div>
                </div>

                <div className="table-wrap">
                  <div className="overflow-x-auto">
                    <table className="table-base min-w-[720px]">
                      <thead>
                        <tr>
                          {workspace.copy.modal.headers.map((header) => (
                            <th key={header}>{header}</th>
                          ))}
                        </tr>
                      </thead>
                      <tbody>
                        {activeChecklist.documents.map((item) => (
                          <tr key={item.name}>
                            <td>{item.name}</td>
                            <td>{item.owner}</td>
                            <td>{item.purpose}</td>
                            <td>{item.critical ? workspace.copy.modal.criticalYes : workspace.copy.modal.criticalNo}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>

                <div className="grid gap-3 md:grid-cols-3">
                  {activeChecklist.intakeSteps.map((step) => (
                    <div key={step} className="rounded-[22px] border border-[rgb(12_91_65_/_0.08)] bg-white/74 p-4">
                      <ClipboardList className="h-5 w-5 text-[var(--brand-deep)]" />
                      <p className="mt-3 text-sm leading-6 text-[var(--muted)]">{step}</p>
                    </div>
                  ))}
                </div>
              </main>

              <aside className="workspace-modal-sidebar">
                <div className="workspace-sidebar-card">
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.modal.blockersSidebarLabel}</p>
                  <div className="mt-3 space-y-2">
                    {activeChecklist.deficiencySignals.map((item) => (
                      <div key={item} className="flex gap-3 rounded-[18px] bg-[var(--danger-soft)] px-4 py-3 text-sm text-[var(--danger)]">
                        <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0" />
                        <span>{item}</span>
                      </div>
                    ))}
                  </div>
                </div>

                <div className="workspace-sidebar-card">
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.modal.sourceFilesLabel}</p>
                  <div className="mt-3 space-y-2">
                    {activeChecklist.sourceFiles.map((file) => (
                      <div key={file} className="rounded-[18px] bg-[rgb(12_91_65_/_0.05)] px-4 py-3 text-sm text-[var(--muted)]">
                        {file}
                      </div>
                    ))}
                  </div>
                </div>

                <div className="workspace-sidebar-card">
                  <div className="flex items-center gap-3">
                    <ShieldCheck className="h-5 w-5 text-[var(--brand-deep)]" />
                    <p className="font-medium text-[var(--text)]">{workspace.copy.modal.lanePostureTitle}</p>
                  </div>
                  <p className="mt-3 text-sm leading-6 text-[var(--muted)]">
                    {activeChecklist.surveyorRequired
                      ? workspace.copy.modal.surveyorLaneBody
                      : workspace.copy.modal.deskLaneBody}
                  </p>
                </div>
              </aside>
            </div>
          </div>
        </div>
      ) : null}
    </div>
  );
}
