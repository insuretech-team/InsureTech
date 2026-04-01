"use client";

import { LoaderCircle, Save, ShieldCheck } from "lucide-react";
import { useEffect, useState } from "react";

import { Panel } from "@/components/panel";
import { useCurrentInsurerId } from "@/hooks/use-current-insurer-id";
import { useInsurerOverview } from "@/hooks/use-insurer-overview";
import { api } from "@/lib/browser-client";
import {
  settingsAuthTypeOptions,
  settingsFeatureToggles,
  settingsSummaryMetrics,
  settingsTabCopy,
} from "@/lib/tabs/settings";
import type { InsurerConfigForm } from "@/lib/types";

export function SettingsPanel() {
  const { insurerId } = useCurrentInsurerId();
  const { overview, loading, error } = useInsurerOverview(insurerId || undefined);
  const [form, setForm] = useState<InsurerConfigForm | null>(null);
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState("");

  useEffect(() => {
    setMessage(error);
  }, [error]);

  useEffect(() => {
    if (overview?.config) {
      setForm(overview.config);
    }
  }, [overview?.config]);

  function updateForm<K extends keyof InsurerConfigForm>(key: K, value: InsurerConfigForm[K]) {
    setForm((current) => (current ? { ...current, [key]: value } : current));
  }

  async function handleSave() {
    if (!form) return;

    setSaving(true);
    setMessage("");

    try {
      const response = await api.insurer.updateConfig(form);
      setMessage(response.ok ? settingsTabCopy.messages.saved : response.message ?? settingsTabCopy.messages.saveFailed);
    } catch {
      setMessage(settingsTabCopy.messages.serviceDown);
    } finally {
      setSaving(false);
    }
  }

  return (
    <div className="page-settings grid gap-6 xl:grid-cols-[minmax(0,1fr)_360px]">
      <Panel title={settingsTabCopy.formPanel.title} description={settingsTabCopy.formPanel.description}>
        {loading || !form ? (
          <div className="flex min-h-[240px] items-center justify-center">
            <LoaderCircle className="h-5 w-5 animate-spin text-[var(--brand-deep)]" />
          </div>
        ) : (
          <div className="grid gap-5 md:grid-cols-2">
            <label className="block space-y-2 md:col-span-2">
              <span className="text-sm font-medium text-[var(--muted)]">{settingsTabCopy.formFields.apiBaseUrl}</span>
              <input
                className="portal-input"
                value={form.apiBaseUrl}
                onChange={(event) => updateForm("apiBaseUrl", event.target.value)}
              />
            </label>

            <label className="block space-y-2">
              <span className="text-sm font-medium text-[var(--muted)]">{settingsTabCopy.formFields.authType}</span>
              <select
                className="portal-select"
                value={form.authType}
                onChange={(event) => updateForm("authType", event.target.value)}
              >
                {settingsAuthTypeOptions.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </select>
            </label>

            <label className="block space-y-2">
              <span className="text-sm font-medium text-[var(--muted)]">{settingsTabCopy.formFields.businessModel}</span>
              <input
                className="portal-input"
                value={form.businessModel}
                onChange={(event) => updateForm("businessModel", event.target.value)}
              />
            </label>

            <label className="block space-y-2 md:col-span-2">
              <span className="text-sm font-medium text-[var(--muted)]">{settingsTabCopy.formFields.authCredentials}</span>
              <textarea
                className="portal-textarea"
                value={form.authCredentials}
                onChange={(event) => updateForm("authCredentials", event.target.value)}
              />
            </label>

            <label className="block space-y-2">
              <span className="text-sm font-medium text-[var(--muted)]">{settingsTabCopy.formFields.paymentTerms}</span>
              <input
                className="portal-input"
                value={form.paymentTerms}
                onChange={(event) => updateForm("paymentTerms", event.target.value)}
              />
            </label>

            <label className="block space-y-2">
              <span className="text-sm font-medium text-[var(--muted)]">{settingsTabCopy.formFields.claimSettlementDays}</span>
              <input
                className="portal-input"
                type="number"
                min={0}
                value={form.claimSettlementDays}
                onChange={(event) => updateForm("claimSettlementDays", Number(event.target.value))}
              />
            </label>

            {settingsFeatureToggles.map((toggle) => (
              <label
                key={toggle.key}
                className="flex items-start gap-3 rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4"
              >
                <input
                  className="mt-1 h-4 w-4 accent-[var(--brand)]"
                  checked={Boolean(form[toggle.key])}
                  onChange={(event) =>
                    updateForm(toggle.key, event.target.checked as InsurerConfigForm[typeof toggle.key])
                  }
                  type="checkbox"
                />
                <span>
                  <span className="block font-medium text-[var(--text)]">{toggle.title}</span>
                  <span className="mt-1 block text-sm text-[var(--muted)]">{toggle.description}</span>
                </span>
              </label>
            ))}

            {message ? (
              <div className="md:col-span-2 rounded-[20px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 px-4 py-3 text-sm text-[var(--muted)]">
                {message}
              </div>
            ) : null}

            <div className="md:col-span-2">
              <button className="portal-btn portal-btn-primary" disabled={saving} onClick={handleSave} type="button">
                {saving ? <LoaderCircle className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
                {saving ? settingsTabCopy.saveButton.saving : settingsTabCopy.saveButton.idle}
              </button>
            </div>
          </div>
        )}
      </Panel>

      <Panel title={settingsTabCopy.summaryPanel.title} description={settingsTabCopy.summaryPanel.description}>
        <div className="space-y-4">
          <div className="rounded-[24px] bg-[rgb(12_91_65_/_0.92)] p-5 text-white">
            <p className="text-sm uppercase tracking-[0.16em] text-white/72">{settingsTabCopy.summaryPanel.currentInsurerLabel}</p>
            <p className="mt-2 font-[family:var(--font-heading)] text-2xl font-semibold">
              {overview?.currentInsurer?.name ?? settingsTabCopy.summaryPanel.fallbackInsurer}
            </p>
            <p className="mt-2 text-sm text-white/76">{overview?.currentInsurer?.businessModel ?? settingsTabCopy.summaryPanel.fallbackModel}</p>
          </div>

          <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4">
            <div className="flex items-start gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-[rgb(15_157_104_/_0.12)] text-[var(--brand-deep)]">
                <ShieldCheck className="h-5 w-5" />
              </div>
              <div>
                <p className="font-medium text-[var(--text)]">{settingsTabCopy.summaryPanel.alignmentTitle}</p>
                <p className="mt-1 text-sm leading-6 text-[var(--muted)]">
                  {settingsTabCopy.summaryPanel.alignmentBody}
                </p>
              </div>
            </div>
          </div>

          <div className="grid gap-3">
            {settingsSummaryMetrics.map((metric) => (
              <div
                key={metric.key}
                className="rounded-[20px] border border-[rgb(12_91_65_/_0.08)] bg-[rgb(255_252_247_/_0.82)] p-4"
              >
                <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{metric.label}</p>
                <p className="mt-2 text-3xl font-semibold text-[var(--text)]">{overview?.metrics[metric.key] ?? 0}</p>
              </div>
            ))}
          </div>
        </div>
      </Panel>
    </div>
  );
}
