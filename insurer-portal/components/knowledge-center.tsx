"use client";

import { BookOpenText, Search, Shapes, Users2 } from "lucide-react";
import Link from "next/link";
import { useMemo, useState } from "react";

import { Panel } from "@/components/panel";
import { getKnowledgeCenterTabData } from "@/lib/tabs/knowledge-center";

export function KnowledgeCenter() {
  const workspace = useMemo(() => getKnowledgeCenterTabData(), []);
  const [query, setQuery] = useState("");
  const [selectedId, setSelectedId] = useState(workspace.assets[0]?.id ?? "");

  const filteredAssets = useMemo(() => {
    const lowered = query.trim().toLowerCase();
    if (!lowered) return workspace.assets;

    return workspace.assets.filter((asset) =>
      [asset.title, asset.category, asset.summary, asset.sourceFile, asset.audience, asset.linkedTabs.join(" ")]
        .join(" ")
        .toLowerCase()
        .includes(lowered),
    );
  }, [query, workspace.assets]);

  const selectedAsset = filteredAssets.find((item) => item.id === selectedId) ?? filteredAssets[0] ?? workspace.assets[0];

  return (
    <div className="page-knowledge-center space-y-6">
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
                <Link className="portal-btn portal-btn-primary" href="/documents">
                  {workspace.copy.hero.primaryAction}
                </Link>
                <Link className="portal-btn portal-btn-secondary" href="/claim-checklists">
                  {workspace.copy.hero.secondaryAction}
                </Link>
              </div>
            </div>

            <div className="grid gap-3">
              <div className="rounded-[24px] bg-[rgb(12_91_65_/_0.92)] p-5 text-white">
                <p className="text-sm uppercase tracking-[0.16em] text-white/72">{workspace.copy.hero.modulesLabel}</p>
                <p className="mt-3 text-4xl font-semibold">{workspace.assets.length}</p>
                <p className="mt-2 text-sm text-white/78">{workspace.copy.hero.modulesDescription}</p>
              </div>
              <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/76 p-5">
                <p className="text-sm uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.hero.portalImpactLabel}</p>
                <p className="mt-3 text-sm leading-6 text-[var(--muted)]">
                  {workspace.impactSummary}
                </p>
              </div>
            </div>
          </div>
        </Panel>

        <Panel title={workspace.copy.searchPanel.title} description={workspace.copy.searchPanel.description}>
          <div className="flex items-center gap-3 rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/78 px-4 py-3">
            <Search className="h-4 w-4 text-[var(--muted)]" />
            <input
              className="w-full bg-transparent text-sm outline-none placeholder:text-[var(--muted)]"
              placeholder={workspace.copy.searchPanel.placeholder}
              value={query}
              onChange={(event) => setQuery(event.target.value)}
            />
          </div>
        </Panel>
      </section>

      <div className="grid gap-6 xl:grid-cols-[360px_minmax(0,1fr)]">
        <Panel title={workspace.copy.listPanel.title} description={workspace.copy.listPanel.description}>
          <div className="space-y-3">
            {filteredAssets.map((asset) => (
              <button
                key={asset.id}
                className={`w-full rounded-[24px] border p-4 text-left transition ${
                  selectedAsset.id === asset.id
                    ? "border-[rgb(12_91_65_/_0.2)] bg-[rgb(12_91_65_/_0.08)]"
                    : "border-[rgb(12_91_65_/_0.08)] bg-white/72"
                }`}
                onClick={() => setSelectedId(asset.id)}
                type="button"
              >
                <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{asset.category}</p>
                <h3 className="mt-2 font-[family:var(--font-heading)] text-xl font-semibold text-[var(--text)]">
                  {asset.title}
                </h3>
                <p className="mt-2 text-sm leading-6 text-[var(--muted)]">{asset.summary}</p>
              </button>
            ))}
          </div>
        </Panel>

        <div className="space-y-6">
          <Panel title={selectedAsset.title} description={selectedAsset.summary}>
            <div className="grid gap-4 lg:grid-cols-[minmax(0,1fr)_320px]">
              <div className="space-y-4">
                <div className="grid gap-3 md:grid-cols-3">
                  {selectedAsset.keyPoints.map((point) => (
                    <div key={point} className="rounded-[22px] border border-[rgb(12_91_65_/_0.08)] bg-white/74 p-4">
                      <BookOpenText className="h-5 w-5 text-[var(--brand-deep)]" />
                      <p className="mt-3 text-sm leading-6 text-[var(--muted)]">{point}</p>
                    </div>
                  ))}
                </div>
              </div>

              <div className="space-y-4">
                <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/76 p-4">
                  <div className="flex items-center gap-2 text-[var(--brand-deep)]">
                    <Users2 className="h-4 w-4" />
                    <p className="text-xs font-semibold uppercase tracking-[0.16em]">{workspace.copy.detailPanel.audienceLabel}</p>
                  </div>
                  <p className="mt-3 text-sm leading-6 text-[var(--text)]">{selectedAsset.audience}</p>
                </div>

                <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/76 p-4">
                  <div className="flex items-center gap-2 text-[var(--brand-deep)]">
                    <Shapes className="h-4 w-4" />
                    <p className="text-xs font-semibold uppercase tracking-[0.16em]">{workspace.copy.detailPanel.linkedTabsLabel}</p>
                  </div>
                  <div className="mt-3 flex flex-wrap gap-2">
                    {selectedAsset.linkedTabs.map((item) => (
                      <span
                        key={item}
                        className="rounded-full border border-[rgb(12_91_65_/_0.08)] bg-[rgb(12_91_65_/_0.05)] px-3 py-1 text-xs font-medium text-[var(--muted)]"
                      >
                        {item}
                      </span>
                    ))}
                  </div>
                </div>

                <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/76 p-4">
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{workspace.copy.detailPanel.sourceFileLabel}</p>
                  <p className="mt-3 text-sm leading-6 text-[var(--text)]">{selectedAsset.sourceFile}</p>
                </div>
              </div>
            </div>
          </Panel>
        </div>
      </div>
    </div>
  );
}
