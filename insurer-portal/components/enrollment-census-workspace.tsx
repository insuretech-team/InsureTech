"use client";

import { ArrowRight, CheckCircle2, ShieldAlert } from "lucide-react";
import Link from "next/link";
import { useMemo, useState } from "react";

import { Panel } from "@/components/panel";
import { StatusPill } from "@/components/status-pill";
import { WorkspaceNoteList } from "@/components/workspace-note-list";
import { getEnrollmentCensusTabData } from "@/lib/tabs/enrollment-census";

export function EnrollmentCensusWorkspace() {
  const workspace = useMemo(() => getEnrollmentCensusTabData(), []);
  const [selectedId, setSelectedId] = useState(workspace.batches[0]?.id ?? "");

  const selectedBatch = workspace.batches.find((item) => item.id === selectedId) ?? workspace.batches[0];

  return (
    <div className="page-enrollment space-y-6">
      <section className="grid gap-4 lg:grid-cols-[minmax(0,1.25fr)_minmax(300px,0.75fr)]">
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
                <Link className="portal-btn portal-btn-primary" href="/documents">
                  {workspace.copy.hero.primaryAction}
                </Link>
                <Link className="portal-btn portal-btn-secondary" href="/proposals">
                  {workspace.copy.hero.secondaryAction}
                </Link>
              </div>
            </div>

            <div className="grid gap-3">
              <div className="rounded-[24px] bg-[rgb(12_91_65_/_0.92)] p-5 text-white">
                <p className="text-sm uppercase tracking-[0.16em] text-white/72">{workspace.copy.hero.membersLabel}</p>
                <p className="mt-3 text-4xl font-semibold">{workspace.metrics.totalMembers}</p>
                <p className="mt-2 text-sm text-white/78">{workspace.copy.hero.membersDescription}</p>
              </div>
              <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/76 p-5">
                <p className="text-sm uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.hero.dependentsLabel}</p>
                <p className="mt-3 text-4xl font-semibold text-[var(--text)]">{workspace.metrics.totalDependents}</p>
                <p className="mt-2 text-sm text-[var(--muted)]">{workspace.metrics.outstandingFlags} {workspace.copy.hero.dependentsDescriptionSuffix}</p>
              </div>
            </div>
          </div>
        </Panel>

        <div className="grid gap-4">
          <Panel title={workspace.copy.focusPanel.title} description={workspace.copy.focusPanel.description}>
            <WorkspaceNoteList items={workspace.operationalFocus} />
          </Panel>
        </div>
      </section>

      <div className="grid gap-6 xl:grid-cols-[380px_minmax(0,1fr)]">
        <Panel title={workspace.copy.batchesPanel.title} description={workspace.copy.batchesPanel.description}>
          <div className="space-y-3">
            {workspace.batches.map((batch) => (
              <button
                key={batch.id}
                className={`w-full rounded-[24px] border p-4 text-left transition ${
                  selectedBatch.id === batch.id
                    ? "border-[rgb(12_91_65_/_0.2)] bg-[rgb(12_91_65_/_0.08)]"
                    : "border-[rgb(12_91_65_/_0.08)] bg-white/72"
                }`}
                onClick={() => setSelectedId(batch.id)}
                type="button"
              >
                <div className="flex items-start justify-between gap-3">
                  <div>
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">
                      {batch.proposalNumber}
                    </p>
                    <h3 className="mt-2 font-[family:var(--font-heading)] text-xl font-semibold text-[var(--text)]">
                      {batch.title}
                    </h3>
                    <p className="mt-1 text-sm text-[var(--muted)]">{batch.clientName}</p>
                  </div>
                  <StatusPill status={batch.status} />
                </div>
                <div className="mt-4 grid gap-3 sm:grid-cols-2">
                  <div className="rounded-[18px] bg-white/88 p-3">
                    <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--muted)]">{workspace.copy.batchesPanel.membersLabel}</p>
                    <p className="mt-2 text-sm font-medium text-[var(--text)]">{batch.memberCount}</p>
                  </div>
                  <div className="rounded-[18px] bg-white/88 p-3">
                    <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--muted)]">{workspace.copy.batchesPanel.dependentsLabel}</p>
                    <p className="mt-2 text-sm font-medium text-[var(--text)]">{batch.dependentCount}</p>
                  </div>
                </div>
              </button>
            ))}
          </div>
        </Panel>

        <div className="space-y-6">
          <Panel
            title={selectedBatch.title}
            description={workspace.copy.detailPanel.description}
            action={
              <Link className="portal-btn portal-btn-secondary" href="/documents">
                {workspace.copy.detailPanel.actionLabel}
              </Link>
            }
          >
            <div className="grid gap-4 lg:grid-cols-[minmax(0,1fr)_320px]">
              <div className="space-y-4">
                <div className="grid gap-3 sm:grid-cols-2">
                  <div className="rounded-[22px] bg-[rgb(12_91_65_/_0.05)] p-4">
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.detailPanel.coverageLabel}</p>
                    <p className="mt-2 text-sm text-[var(--text)]">{selectedBatch.coverageWindow}</p>
                  </div>
                  <div className="rounded-[22px] bg-[rgb(245_158_11_/_0.1)] p-4">
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.detailPanel.postureLabel}</p>
                    <p className="mt-2 text-sm text-[var(--text)]">{selectedBatch.status}</p>
                  </div>
                </div>

                <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/74 p-4">
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.detailPanel.notesLabel}</p>
                  <div className="mt-3 space-y-3">
                    {selectedBatch.notes.map((note) => (
                      <div key={note} className="rounded-[18px] bg-[rgb(255_252_247_/_0.86)] px-4 py-3 text-sm leading-6 text-[var(--muted)]">
                        {note}
                      </div>
                    ))}
                  </div>
                </div>
              </div>

              <div className="space-y-4">
                <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/78 p-4">
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.detailPanel.flagsLabel}</p>
                  <div className="mt-3 space-y-2">
                    {selectedBatch.validationFlags.map((flag) => (
                      <div key={flag} className="flex gap-3 rounded-[18px] bg-[var(--danger-soft)] px-4 py-3 text-sm text-[var(--danger)]">
                        <ShieldAlert className="mt-0.5 h-4 w-4 shrink-0" />
                        <span>{flag}</span>
                      </div>
                    ))}
                  </div>
                </div>

                <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/78 p-4">
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.detailPanel.filesLabel}</p>
                  <div className="mt-3 space-y-2">
                    {selectedBatch.sourceFiles.map((file) => (
                      <div key={file} className="flex items-center justify-between gap-3 rounded-[18px] bg-[rgb(12_91_65_/_0.05)] px-4 py-3">
                        <span className="text-sm text-[var(--text)]">{file}</span>
                        <ArrowRight className="h-4 w-4 text-[var(--muted)]" />
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          </Panel>

          <Panel title={workspace.copy.rosterPanel.title} description={workspace.copy.rosterPanel.description}>
            <div className="table-wrap">
              <div className="overflow-x-auto">
                <table className="table-base min-w-[860px]">
                  <thead>
                    <tr>
                      {workspace.copy.rosterPanel.headers.map((header) => (
                        <th key={header}>{header}</th>
                      ))}
                    </tr>
                  </thead>
                  <tbody>
                    {selectedBatch.rows.map((row) => (
                      <tr key={row.employeeId}>
                        <td>{row.employeeId}</td>
                        <td>{row.memberName}</td>
                        <td>{row.relation}</td>
                        <td>{row.designation}</td>
                        <td>{row.sumAssuredText}</td>
                        <td>{row.coverageStart}</td>
                        <td>{row.nominee}</td>
                        <td>{row.phone}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </Panel>

          <Panel title={workspace.copy.checklistPanel.title} description={workspace.copy.checklistPanel.description}>
            <div className="grid gap-3 md:grid-cols-3">
              {workspace.dispatchChecklist.map((item) => (
                <div key={item} className="rounded-[22px] border border-[rgb(12_91_65_/_0.08)] bg-white/74 p-4">
                  <div className="flex items-start gap-3">
                    <CheckCircle2 className="mt-0.5 h-5 w-5 text-[var(--brand-deep)]" />
                    <p className="text-sm leading-6 text-[var(--muted)]">{item}</p>
                  </div>
                </div>
              ))}
            </div>
          </Panel>
        </div>
      </div>
    </div>
  );
}
