"use client";

import { FilePlus2, Grip, LoaderCircle, Plus, Save, Trash2 } from "lucide-react";
import { useEffect, useMemo, useState } from "react";

import { Panel } from "@/components/panel";
import { useCurrentInsurerId } from "@/hooks/use-current-insurer-id";
import { useInsurerOverview } from "@/hooks/use-insurer-overview";
import { usePersistedState } from "@/hooks/use-persisted-state";
import { productPlaybooks } from "@/lib/product-playbooks";
import {
  createEmptyPlanTemplate,
  createPlanTemplateFromPlaybook,
  createTemplateBlock,
  planTemplateKindOptions,
  planTemplatesTabCopy,
  planTemplatesStorageKey,
  type PlanTemplate,
  type TemplateBlock,
} from "@/lib/tabs/plan-templates";

export function PlanTemplateStudio() {
  const { insurerId } = useCurrentInsurerId();
  const { overview, loading } = useInsurerOverview(insurerId || undefined);
  const { value: templates, setValue: setTemplates, ready } = usePersistedState<PlanTemplate[]>(
    planTemplatesStorageKey,
    [],
  );
  const [selectedId, setSelectedId] = useState("");
  const [message, setMessage] = useState("");

  useEffect(() => {
    if (!ready || !overview || templates.length) return;

    const seeded = productPlaybooks.map((playbook) => {
      const matchingProduct = overview.products.find(
        (product) =>
          product.name.toLowerCase().includes(playbook.name.toLowerCase()) ||
          playbook.name.toLowerCase().includes(product.name.toLowerCase()),
      );
      return createPlanTemplateFromPlaybook(playbook, matchingProduct?.id || "");
    });
    setTemplates(seeded);
  }, [overview, ready, setTemplates, templates.length]);

  const selected = useMemo(
    () => templates.find((template) => template.id === selectedId) ?? templates[0] ?? null,
    [selectedId, templates],
  );

  function updateTemplate(id: string, patch: Partial<PlanTemplate>) {
    setMessage("");
    setTemplates((current) => current.map((template) => (template.id === id ? { ...template, ...patch } : template)));
  }

  function updateBlock(templateId: string, blockId: string, patch: Partial<TemplateBlock>) {
    setMessage("");
    setTemplates((current) =>
      current.map((template) =>
        template.id === templateId
          ? {
              ...template,
              blocks: template.blocks.map((block) => (block.id === blockId ? { ...block, ...patch } : block)),
            }
          : template,
      ),
    );
  }

  function addTemplate() {
    const next = createEmptyPlanTemplate(overview?.products[0]?.id || "");
    setTemplates((current) => [next, ...current]);
    setSelectedId(next.id);
  }

  function removeTemplate(id: string) {
    const next = templates.filter((template) => template.id !== id);
    setTemplates(next);
    setSelectedId(next[0]?.id || "");
  }

  function addBlock(templateId: string) {
    setTemplates((current) =>
      current.map((template) =>
        template.id === templateId
          ? {
              ...template,
              blocks: [
              ...template.blocks,
                createTemplateBlock("terms"),
              ],
            }
          : template,
      ),
    );
  }

  function removeBlock(templateId: string, blockId: string) {
    setTemplates((current) =>
      current.map((template) =>
        template.id === templateId
          ? { ...template, blocks: template.blocks.filter((block) => block.id !== blockId) }
          : template,
      ),
    );
  }

  function saveTemplates() {
    setMessage(planTemplatesTabCopy.studioPanel.saveMessage);
  }

  return (
    <div className="page-plan-templates grid gap-6 xl:grid-cols-[330px_minmax(0,1fr)]">
      <Panel
        title={planTemplatesTabCopy.libraryPanel.title}
        description={planTemplatesTabCopy.libraryPanel.description}
        action={
          <button className="portal-btn portal-btn-primary" onClick={addTemplate} type="button">
            <Plus className="h-4 w-4" />
            {planTemplatesTabCopy.libraryPanel.newButton}
          </button>
        }
      >
        {loading ? (
          <div className="flex min-h-[220px] items-center justify-center">
            <LoaderCircle className="h-5 w-5 animate-spin text-[var(--brand-deep)]" />
          </div>
        ) : (
          <div className="space-y-3">
            {templates.map((template) => {
              const product = overview?.products.find((item) => item.id === template.productId);

              return (
                <button
                  key={template.id}
                  className={`w-full rounded-[24px] border p-4 text-left transition ${
                    selected?.id === template.id
                      ? "border-[rgb(15_157_104_/_0.3)] bg-[rgb(15_157_104_/_0.08)]"
                      : "border-[rgb(12_91_65_/_0.08)] bg-white/72"
                  }`}
                  onClick={() => setSelectedId(template.id)}
                  type="button"
                >
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <p className="font-[family:var(--font-heading)] text-lg font-semibold text-[var(--text)]">
                        {template.title}
                      </p>
                      <p className="mt-1 text-sm text-[var(--muted)]">{product?.name ?? planTemplatesTabCopy.libraryPanel.emptyProduct}</p>
                    </div>
                    <FilePlus2 className="h-4 w-4 text-[var(--brand-deep)]" />
                  </div>
                  <p className="mt-3 text-sm leading-6 text-[var(--muted)]">{template.summary}</p>
                </button>
              );
            })}
          </div>
        )}
      </Panel>

      <Panel title={planTemplatesTabCopy.studioPanel.title} description={planTemplatesTabCopy.studioPanel.description}>
        {selected ? (
          <div className="space-y-5">
            <div className="grid gap-5 md:grid-cols-2">
              <label className="block space-y-2">
                <span className="text-sm font-medium text-[var(--muted)]">{planTemplatesTabCopy.studioPanel.titleLabel}</span>
                <input
                  className="portal-input"
                  value={selected.title}
                  onChange={(event) => updateTemplate(selected.id, { title: event.target.value })}
                />
              </label>

              <label className="block space-y-2">
                <span className="text-sm font-medium text-[var(--muted)]">{planTemplatesTabCopy.studioPanel.productLabel}</span>
                <select
                  className="portal-select"
                  value={selected.productId}
                  onChange={(event) => updateTemplate(selected.id, { productId: event.target.value })}
                >
                  <option value="">{planTemplatesTabCopy.studioPanel.emptyProductOption}</option>
                  {(overview?.products ?? []).map((product) => (
                    <option key={product.id} value={product.id}>
                      {product.name}
                    </option>
                  ))}
                </select>
              </label>

              <label className="block space-y-2 md:col-span-2">
                <span className="text-sm font-medium text-[var(--muted)]">{planTemplatesTabCopy.studioPanel.summaryLabel}</span>
                <textarea
                  className="portal-textarea"
                  value={selected.summary}
                  onChange={(event) => updateTemplate(selected.id, { summary: event.target.value })}
                />
              </label>
            </div>

            <div className="space-y-4">
              {selected.blocks.map((block) => (
                <div
                  key={block.id}
                  className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4"
                >
                  <div className="mb-4 flex items-center justify-between gap-3">
                    <div className="flex items-center gap-3">
                      <div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-[rgb(12_91_65_/_0.08)] text-[var(--brand-deep)]">
                        <Grip className="h-4 w-4" />
                      </div>
                      <div>
                        <p className="text-sm font-medium text-[var(--text)]">{block.title}</p>
                        <p className="text-xs uppercase tracking-[0.14em] text-[var(--muted)]">{block.kind}</p>
                      </div>
                    </div>
                    <button
                      className="portal-btn portal-btn-secondary"
                      onClick={() => removeBlock(selected.id, block.id)}
                      type="button"
                    >
                      <Trash2 className="h-4 w-4" />
                      {planTemplatesTabCopy.studioPanel.removeBlockButton}
                    </button>
                  </div>

                  <div className="grid gap-4 md:grid-cols-[180px_minmax(0,1fr)]">
                    <label className="block space-y-2">
                      <span className="text-sm font-medium text-[var(--muted)]">{planTemplatesTabCopy.studioPanel.sectionTypeLabel}</span>
                      <select
                        className="portal-select"
                        value={block.kind}
                        onChange={(event) =>
                          updateBlock(selected.id, block.id, {
                            kind: event.target.value as TemplateBlock["kind"],
                          })
                        }
                      >
                        {planTemplateKindOptions.map((option) => (
                          <option key={option.value} value={option.value}>
                            {option.label}
                          </option>
                        ))}
                      </select>
                    </label>

                    <label className="block space-y-2">
                      <span className="text-sm font-medium text-[var(--muted)]">{planTemplatesTabCopy.studioPanel.sectionTitleLabel}</span>
                      <input
                        className="portal-input"
                        value={block.title}
                        onChange={(event) => updateBlock(selected.id, block.id, { title: event.target.value })}
                      />
                    </label>
                  </div>

                  <label className="mt-4 block space-y-2">
                    <span className="text-sm font-medium text-[var(--muted)]">{planTemplatesTabCopy.studioPanel.sectionBodyLabel}</span>
                    <textarea
                      className="portal-textarea"
                      value={block.body}
                      onChange={(event) => updateBlock(selected.id, block.id, { body: event.target.value })}
                    />
                  </label>
                </div>
              ))}
            </div>

            {message ? (
              <div className="rounded-[20px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 px-4 py-3 text-sm text-[var(--muted)]">
                {message}
              </div>
            ) : null}

            <div className="flex flex-wrap gap-3">
              <button className="portal-btn portal-btn-secondary" onClick={() => addBlock(selected.id)} type="button">
                <FilePlus2 className="h-4 w-4" />
                {planTemplatesTabCopy.studioPanel.addSectionButton}
              </button>
              <button className="portal-btn portal-btn-primary" onClick={saveTemplates} type="button">
                <Save className="h-4 w-4" />
                {planTemplatesTabCopy.studioPanel.saveButton}
              </button>
              <button
                className="portal-btn border border-[rgb(194_65_12_/_0.14)] bg-[var(--danger-soft)] text-[var(--danger)]"
                onClick={() => removeTemplate(selected.id)}
                type="button"
              >
                <Trash2 className="h-4 w-4" />
                {planTemplatesTabCopy.studioPanel.deleteButton}
              </button>
            </div>
          </div>
        ) : (
          <p className="text-sm text-[var(--muted)]">{planTemplatesTabCopy.studioPanel.emptyLabel}</p>
        )}
      </Panel>
    </div>
  );
}
