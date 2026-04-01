"use client";

import { FileSpreadsheet, ReceiptText } from "lucide-react";
import Link from "next/link";
import { useMemo, useState } from "react";

import { Panel } from "@/components/panel";
import { WorkspaceIconCardGrid } from "@/components/workspace-icon-card-grid";
import { WorkspaceNoteList } from "@/components/workspace-note-list";
import { StatusPill } from "@/components/status-pill";
import { getPricingCommercialsTabData } from "@/lib/tabs/pricing-commercials";

export function PricingCommercialsWorkspace() {
  const workspace = useMemo(() => getPricingCommercialsTabData(), []);
  const [selectedId, setSelectedId] = useState(workspace.scenarios[0]?.id ?? "");
  const selectedScenario = workspace.scenarios.find((item) => item.id === selectedId) ?? workspace.scenarios[0];

  return (
    <div className="page-pricing space-y-6">
      <section className="grid gap-4 xl:grid-cols-[minmax(0,1.2fr)_minmax(340px,0.8fr)]">
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
                <Link className="portal-btn portal-btn-secondary" href="/enrollment-census">
                  {workspace.copy.hero.secondaryAction}
                </Link>
              </div>
            </div>

            <div className="grid gap-3">
              <div className="rounded-[24px] bg-[rgb(245_158_11_/_0.12)] p-5">
                <p className="text-sm uppercase tracking-[0.16em] text-[#8a5200]">{workspace.copy.hero.leadQuoteLabel}</p>
                <p className="mt-3 font-[family:var(--font-heading)] text-3xl font-semibold text-[var(--text)]">
                  {workspace.scenarios[0]?.headlinePremiumText}
                </p>
                <p className="mt-2 text-sm text-[var(--muted)]">{workspace.copy.hero.leadQuoteDescription}</p>
              </div>
              <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/74 p-5">
                <p className="text-sm uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.hero.activeCasesLabel}</p>
                <p className="mt-3 text-4xl font-semibold text-[var(--text)]">{workspace.scenarios.length}</p>
                <p className="mt-2 text-sm text-[var(--muted)]">{workspace.copy.hero.activeCasesDescription}</p>
              </div>
            </div>
          </div>
        </Panel>

        <Panel title={workspace.copy.outputsPanel.title} description={workspace.copy.outputsPanel.description}>
          <WorkspaceNoteList items={workspace.commercialOutputs} />
        </Panel>
      </section>

      <div className="grid gap-6 xl:grid-cols-[360px_minmax(0,1fr)]">
        <Panel title={workspace.copy.casesPanel.title} description={workspace.copy.casesPanel.description}>
          <div className="space-y-3">
            {workspace.scenarios.map((scenario) => (
              <button
                key={scenario.id}
                className={`w-full rounded-[24px] border p-4 text-left transition ${
                  selectedScenario.id === scenario.id
                    ? "border-[rgb(12_91_65_/_0.2)] bg-[rgb(12_91_65_/_0.08)]"
                    : "border-[rgb(12_91_65_/_0.08)] bg-white/72"
                }`}
                onClick={() => setSelectedId(scenario.id)}
                type="button"
              >
                <div className="flex items-start justify-between gap-3">
                  <div>
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">
                      {scenario.segment}
                    </p>
                    <h3 className="mt-2 font-[family:var(--font-heading)] text-xl font-semibold text-[var(--text)]">
                      {scenario.title}
                    </h3>
                    <p className="mt-1 text-sm text-[var(--muted)]">{scenario.clientName}</p>
                  </div>
                  <StatusPill status={scenario.status} />
                </div>
                <div className="mt-4 rounded-[18px] bg-white/88 p-3">
                  <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--muted)]">{workspace.copy.casesPanel.premiumLabel}</p>
                  <p className="mt-2 text-sm font-medium text-[var(--text)]">{scenario.headlinePremiumText}</p>
                </div>
              </button>
            ))}
          </div>
        </Panel>

        <div className="space-y-6">
          <Panel
            title={selectedScenario.title}
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
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.detailPanel.clientLabel}</p>
                    <p className="mt-2 text-sm text-[var(--text)]">{selectedScenario.clientName}</p>
                  </div>
                  <div className="rounded-[22px] bg-[rgb(245_158_11_/_0.1)] p-4">
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.detailPanel.segmentLabel}</p>
                    <p className="mt-2 text-sm text-[var(--text)]">{selectedScenario.segment}</p>
                  </div>
                </div>

                <div className="table-wrap">
                  <div className="overflow-x-auto">
                    <table className="table-base min-w-[720px]">
                      <thead>
                        <tr>
                          {workspace.copy.detailPanel.headers.map((header) => (
                            <th key={header}>{header}</th>
                          ))}
                        </tr>
                      </thead>
                      <tbody>
                        {selectedScenario.lineItems.map((item) => (
                          <tr key={item.label}>
                            <td>{item.label}</td>
                            <td>{item.members}</td>
                            <td>{item.benefit}</td>
                            <td>{item.rateText}</td>
                            <td>{item.annualPremiumText}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>

              <div className="space-y-4">
                <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/76 p-4">
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.detailPanel.assumptionsLabel}</p>
                  <div className="mt-3 space-y-2">
                    {selectedScenario.assumptions.map((item) => (
                      <div key={item} className="rounded-[18px] bg-[rgb(12_91_65_/_0.05)] px-4 py-3 text-sm text-[var(--muted)]">
                        {item}
                      </div>
                    ))}
                  </div>
                </div>

                <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/76 p-4">
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.detailPanel.openItemsLabel}</p>
                  <div className="mt-3 space-y-2">
                    {selectedScenario.openItems.map((item) => (
                      <div key={item} className="rounded-[18px] bg-[var(--danger-soft)] px-4 py-3 text-sm text-[var(--danger)]">
                        {item}
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          </Panel>

          <div className="grid gap-6 lg:grid-cols-2">
            <Panel title={workspace.copy.deliverablesPanel.title} description={workspace.copy.deliverablesPanel.description}>
              <div className="space-y-3">
                {selectedScenario.deliverables.map((item) => (
                  <div key={item} className="flex items-start gap-3 rounded-[20px] border border-[rgb(12_91_65_/_0.08)] bg-white/76 px-4 py-3">
                    <ReceiptText className="mt-0.5 h-4 w-4 text-[var(--brand-deep)]" />
                    <p className="text-sm text-[var(--muted)]">{item}</p>
                  </div>
                ))}
              </div>
            </Panel>

            <Panel title={workspace.copy.sourceFilesPanel.title} description={workspace.copy.sourceFilesPanel.description}>
              <div className="space-y-3">
                {selectedScenario.sourceFiles.map((file) => (
                  <div key={file} className="flex items-start gap-3 rounded-[20px] border border-[rgb(12_91_65_/_0.08)] bg-white/76 px-4 py-3">
                    <FileSpreadsheet className="mt-0.5 h-4 w-4 text-[var(--brand-deep)]" />
                    <p className="text-sm text-[var(--muted)]">{file}</p>
                  </div>
                ))}
              </div>
            </Panel>
          </div>

          <Panel title={workspace.copy.valuePanel.title} description={workspace.copy.valuePanel.description}>
            <WorkspaceIconCardGrid items={workspace.valueCards} />
          </Panel>
        </div>
      </div>
    </div>
  );
}
