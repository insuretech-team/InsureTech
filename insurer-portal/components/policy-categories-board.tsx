"use client";

import { Layers3, LoaderCircle, Save } from "lucide-react";
import { useMemo, useState } from "react";

import { Panel } from "@/components/panel";
import { StatusPill } from "@/components/status-pill";
import { useCurrentInsurerId } from "@/hooks/use-current-insurer-id";
import { useInsurerOverview } from "@/hooks/use-insurer-overview";
import { usePersistedState } from "@/hooks/use-persisted-state";
import {
  policyCategoriesStorageKey,
  policyCategoriesTabCopy,
  policyCategoryOptions as categoryOptions,
} from "@/lib/tabs/policy-categories";
import { findPlaybook } from "@/lib/product-playbooks";
import type { PortalProduct } from "@/lib/types";

type DraftRecord = Record<
  string,
  {
    category: string;
    tags: string;
    notes: string;
  }
>;

export function PolicyCategoriesBoard() {
  const { insurerId } = useCurrentInsurerId();
  const { overview, loading } = useInsurerOverview(insurerId || undefined);
  const [selectedId, setSelectedId] = useState("");
  const { value: drafts, setValue: setDrafts } = usePersistedState<DraftRecord>(policyCategoriesStorageKey, {});
  const [saved, setSaved] = useState("");

  const grouped = useMemo(() => {
    const products = overview?.products ?? [];
    return categoryOptions.map((category) => ({
      category,
      count: products.filter((item) => (drafts[item.id]?.category || item.category) === category).length,
    }));
  }, [drafts, overview?.products]);

  const selectedProduct = overview?.products.find((item) => item.id === selectedId) ?? overview?.products[0] ?? null;
  const selectedDraft = selectedProduct ? drafts[selectedProduct.id] : undefined;
  const selectedPlaybook = selectedProduct ? findPlaybook(selectedProduct.name, selectedProduct.category) : undefined;

  function updateDraft(product: PortalProduct, patch: Partial<DraftRecord[string]>) {
    setSaved("");
    setDrafts((current) => ({
      ...current,
      [product.id]: {
        category: current[product.id]?.category || product.category,
        tags: current[product.id]?.tags || "",
        notes: current[product.id]?.notes || "",
        ...patch,
      },
    }));
  }

  function saveLocalChanges() {
    setSaved(policyCategoriesTabCopy.editorPanel.saveMessage);
  }

  return (
    <div className="page-policy-categories grid gap-6 xl:grid-cols-[minmax(0,1fr)_390px]">
      <div className="space-y-6">
        <Panel title={policyCategoriesTabCopy.footprintPanel.title} description={policyCategoriesTabCopy.footprintPanel.description}>
          {loading ? (
            <div className="flex min-h-[200px] items-center justify-center">
              <LoaderCircle className="h-5 w-5 animate-spin text-[var(--brand-deep)]" />
            </div>
          ) : (
            <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
              {grouped.map((item) => (
                <div
                  key={item.category}
                  className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4"
                >
                  <div className="flex items-center justify-between gap-3">
                    <p className="font-medium text-[var(--text)]">{item.category}</p>
                    <div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-[rgb(15_157_104_/_0.12)] text-[var(--brand-deep)]">
                      <Layers3 className="h-4 w-4" />
                    </div>
                  </div>
                  <p className="mt-4 font-[family:var(--font-heading)] text-4xl font-semibold text-[var(--text)]">
                    {item.count}
                  </p>
                  <p className="mt-1 text-sm text-[var(--muted)]">{policyCategoriesTabCopy.footprintPanel.productsDescription}</p>
                </div>
              ))}
            </div>
          )}
        </Panel>

        <Panel title={policyCategoriesTabCopy.mapPanel.title} description={policyCategoriesTabCopy.mapPanel.description}>
          <div className="grid gap-4">
            {(overview?.products ?? []).map((product) => {
              const effectiveCategory = drafts[product.id]?.category || product.category;

              return (
                <button
                  key={product.id}
                  className={`rounded-[24px] border p-4 text-left transition ${
                    selectedProduct?.id === product.id
                      ? "border-[rgb(15_157_104_/_0.3)] bg-[rgb(15_157_104_/_0.08)]"
                      : "border-[rgb(12_91_65_/_0.08)] bg-white/72"
                  }`}
                  onClick={() => setSelectedId(product.id)}
                  type="button"
                >
                  <div className="flex flex-wrap items-center gap-3">
                    <h3 className="font-[family:var(--font-heading)] text-xl font-semibold text-[var(--text)]">
                      {product.name}
                    </h3>
                    <StatusPill status={product.status} />
                    <span className="rounded-full bg-[rgb(245_158_11_/_0.12)] px-3 py-1 text-xs font-semibold text-[#8a5200]">
                      {effectiveCategory}
                    </span>
                  </div>
                  <p className="mt-2 text-sm text-[var(--muted)]">
                    {product.code} • {product.premiumRangeText}
                  </p>
                </button>
              );
            })}
          </div>
        </Panel>
      </div>

      <Panel title={policyCategoriesTabCopy.editorPanel.title} description={policyCategoriesTabCopy.editorPanel.description}>
        {selectedProduct ? (
          <div className="space-y-5">
            <div className="rounded-[24px] bg-[rgb(12_91_65_/_0.05)] p-4">
              <div className="flex items-center justify-between gap-4">
                <div>
                  <p className="text-sm uppercase tracking-[0.14em] text-[var(--muted)]">{selectedProduct.code}</p>
                  <h3 className="mt-2 font-[family:var(--font-heading)] text-2xl font-semibold text-[var(--text)]">
                    {selectedProduct.name}
                  </h3>
                </div>
                <StatusPill status={selectedProduct.status} />
              </div>
              <p className="mt-3 text-sm text-[var(--muted)]">
                {policyCategoriesTabCopy.editorPanel.coveragePrefix} {selectedProduct.coverageRangeText} • {policyCategoriesTabCopy.editorPanel.premiumPrefix} {selectedProduct.premiumRangeText}
              </p>
            </div>

            <label className="block space-y-2">
              <span className="text-sm font-medium text-[var(--muted)]">{policyCategoriesTabCopy.editorPanel.categoryLabel}</span>
              <select
                className="portal-select"
                value={selectedDraft?.category || selectedProduct.category}
                onChange={(event) => updateDraft(selectedProduct, { category: event.target.value })}
              >
                {categoryOptions.map((option) => (
                  <option key={option} value={option}>
                    {option}
                  </option>
                ))}
              </select>
            </label>

            <label className="block space-y-2">
              <span className="text-sm font-medium text-[var(--muted)]">{policyCategoriesTabCopy.editorPanel.tagsLabel}</span>
              <input
                className="portal-input"
                placeholder={policyCategoriesTabCopy.editorPanel.tagsPlaceholder}
                value={selectedDraft?.tags || ""}
                onChange={(event) => updateDraft(selectedProduct, { tags: event.target.value })}
              />
            </label>

            <label className="block space-y-2">
              <span className="text-sm font-medium text-[var(--muted)]">{policyCategoriesTabCopy.editorPanel.notesLabel}</span>
              <textarea
                className="portal-textarea"
                placeholder={policyCategoriesTabCopy.editorPanel.notesPlaceholder}
                value={selectedDraft?.notes || ""}
                onChange={(event) => updateDraft(selectedProduct, { notes: event.target.value })}
              />
            </label>

            {selectedPlaybook ? (
              <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-[rgb(255_252_247_/_0.82)] p-4">
                <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{policyCategoriesTabCopy.editorPanel.playbookLabel}</p>
                <p className="mt-2 text-sm leading-6 text-[var(--muted)]">{selectedPlaybook.summary}</p>
                <div className="mt-4">
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{policyCategoriesTabCopy.editorPanel.docsLabel}</p>
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
                  <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{policyCategoriesTabCopy.editorPanel.flagsLabel}</p>
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

            {saved ? (
              <div className="rounded-[20px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 px-4 py-3 text-sm text-[var(--muted)]">
                {saved}
              </div>
            ) : null}

            <button className="portal-btn portal-btn-primary w-full" onClick={saveLocalChanges} type="button">
              <Save className="h-4 w-4" />
              {policyCategoriesTabCopy.editorPanel.saveButton}
            </button>
          </div>
        ) : (
          <p className="text-sm text-[var(--muted)]">{policyCategoriesTabCopy.editorPanel.emptyLabel}</p>
        )}
      </Panel>
    </div>
  );
}
