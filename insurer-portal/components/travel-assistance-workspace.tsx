"use client";

import { Mail, Phone, Route, Siren } from "lucide-react";
import Link from "next/link";
import { useMemo, useState } from "react";

import { Panel } from "@/components/panel";
import { WorkspaceIconCardGrid } from "@/components/workspace-icon-card-grid";
import { WorkspaceNoteList } from "@/components/workspace-note-list";
import { getTravelAssistanceTabData } from "@/lib/tabs/travel-assistance";

export function TravelAssistanceWorkspace() {
  const workspace = useMemo(() => getTravelAssistanceTabData(), []);
  const [selectedId, setSelectedId] = useState(workspace.providers[0]?.id ?? "");
  const selectedProvider = workspace.providers.find((item) => item.id === selectedId) ?? workspace.providers[0];

  return (
    <div className="page-travel-assistance space-y-6">
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
                <p className="text-sm uppercase tracking-[0.16em] text-white/72">{workspace.copy.hero.providersLabel}</p>
                <p className="mt-3 text-4xl font-semibold">{workspace.providers.length}</p>
                <p className="mt-2 text-sm text-white/78">{workspace.copy.hero.providersDescription}</p>
              </div>
              <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/76 p-5">
                <p className="text-sm uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.hero.runbooksLabel}</p>
                <p className="mt-3 text-4xl font-semibold text-[var(--text)]">{workspace.runbooks.length}</p>
                <p className="mt-2 text-sm text-[var(--muted)]">{workspace.copy.hero.runbooksDescription}</p>
              </div>
            </div>
          </div>
        </Panel>

        <Panel title={workspace.copy.focusPanel.title} description={workspace.copy.focusPanel.description}>
          <WorkspaceNoteList items={workspace.routingFocus} />
        </Panel>
      </section>

      <div className="grid gap-6 xl:grid-cols-[360px_minmax(0,1fr)]">
        <Panel title={workspace.copy.providersPanel.title} description={workspace.copy.providersPanel.description}>
          <div className="space-y-3">
            {workspace.providers.map((provider) => (
              <button
                key={provider.id}
                className={`w-full rounded-[24px] border p-4 text-left transition ${
                  selectedProvider.id === provider.id
                    ? "border-[rgb(12_91_65_/_0.2)] bg-[rgb(12_91_65_/_0.08)]"
                    : "border-[rgb(12_91_65_/_0.08)] bg-white/72"
                }`}
                onClick={() => setSelectedId(provider.id)}
                type="button"
              >
                <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{provider.role}</p>
                <h3 className="mt-2 font-[family:var(--font-heading)] text-xl font-semibold text-[var(--text)]">
                  {provider.name}
                </h3>
                <p className="mt-2 text-sm text-[var(--muted)]">{provider.supportWindow}</p>
              </button>
            ))}
          </div>
        </Panel>

        <div className="space-y-6">
          <Panel title={selectedProvider.name} description={workspace.copy.detailPanel.description}>
            <div className="grid gap-4 lg:grid-cols-[minmax(0,1fr)_320px]">
              <div className="space-y-4">
                <div className="grid gap-3 md:grid-cols-2">
                  <div className="rounded-[22px] bg-[rgb(12_91_65_/_0.05)] p-4">
                    <div className="flex items-center gap-2 text-[var(--brand-deep)]">
                      <Phone className="h-4 w-4" />
                      <span className="text-xs font-semibold uppercase tracking-[0.16em]">{workspace.copy.detailPanel.phonesLabel}</span>
                    </div>
                    <div className="mt-3 space-y-2">
                      {selectedProvider.phones.map((phone) => (
                        <p key={phone} className="text-sm text-[var(--text)]">{phone}</p>
                      ))}
                    </div>
                  </div>
                  <div className="rounded-[22px] bg-[rgb(245_158_11_/_0.1)] p-4">
                    <div className="flex items-center gap-2 text-[#8a5200]">
                      <Mail className="h-4 w-4" />
                      <span className="text-xs font-semibold uppercase tracking-[0.16em]">{workspace.copy.detailPanel.emailsLabel}</span>
                    </div>
                    <div className="mt-3 space-y-2">
                      {selectedProvider.emails.length ? (
                        selectedProvider.emails.map((email) => (
                          <p key={email} className="text-sm text-[var(--text)]">{email}</p>
                        ))
                      ) : (
                        <p className="text-sm text-[var(--muted)]">{workspace.copy.detailPanel.emptyEmail}</p>
                      )}
                    </div>
                  </div>
                </div>

                <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/78 p-4">
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.detailPanel.useCasesLabel}</p>
                  <div className="mt-3 grid gap-3 md:grid-cols-3">
                    {selectedProvider.useCases.map((item) => (
                      <div key={item} className="rounded-[18px] bg-[rgb(255_252_247_/_0.86)] px-4 py-3 text-sm text-[var(--muted)]">
                        {item}
                      </div>
                    ))}
                  </div>
                </div>

                <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/78 p-4">
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.detailPanel.addressLabel}</p>
                  <p className="mt-3 text-sm leading-6 text-[var(--text)]">{selectedProvider.address}</p>
                </div>
              </div>

              <div className="space-y-4">
                <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/78 p-4">
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.detailPanel.handoffLabel}</p>
                  <div className="mt-3 space-y-2">
                    {selectedProvider.handoffPacket.map((item) => (
                      <div key={item} className="rounded-[18px] bg-[rgb(12_91_65_/_0.05)] px-4 py-3 text-sm text-[var(--muted)]">
                        {item}
                      </div>
                    ))}
                  </div>
                </div>

                <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/78 p-4">
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.detailPanel.sourceFilesLabel}</p>
                  <div className="mt-3 space-y-2">
                    {selectedProvider.sourceFiles.map((file) => (
                      <div key={file} className="rounded-[18px] bg-[rgb(12_91_65_/_0.05)] px-4 py-3 text-sm text-[var(--muted)]">
                        {file}
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          </Panel>

          <Panel title={workspace.copy.runbooksPanel.title} description={workspace.copy.runbooksPanel.description}>
            <div className="grid gap-4 lg:grid-cols-2">
              {workspace.runbooks.map((item) => (
                <div key={item.title} className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/76 p-4">
                  <div className="flex items-center gap-2 text-[var(--brand-deep)]">
                    <Siren className="h-4 w-4" />
                    <p className="text-xs font-semibold uppercase tracking-[0.16em]">{item.title}</p>
                  </div>
                  <p className="mt-3 text-sm leading-6 text-[var(--muted)]">{item.description}</p>
                  <div className="mt-4 space-y-3">
                    {item.steps.map((step) => (
                      <div key={step} className="flex gap-3 rounded-[18px] bg-[rgb(12_91_65_/_0.05)] px-4 py-3">
                        <Route className="mt-0.5 h-4 w-4 shrink-0 text-[var(--brand-deep)]" />
                        <p className="text-sm text-[var(--muted)]">{step}</p>
                      </div>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          </Panel>

          <Panel title={workspace.copy.rulesPanel.title} description={workspace.copy.rulesPanel.description}>
            <WorkspaceIconCardGrid items={workspace.timeSensitiveRules} />
          </Panel>
        </div>
      </div>
    </div>
  );
}
