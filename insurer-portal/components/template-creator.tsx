"use client";

import {
  ChevronDown,
  ChevronUp,
  Download,
  Eye,
  FilePlus2,
  LoaderCircle,
  Plus,
  Save,
  Sparkles,
  Trash2,
  X,
} from "lucide-react";
import { useState } from "react";

// ─── Types ───────────────────────────────────────────────────────────────────

type SectionType =
  | "header" | "divider" | "title" | "subtitle" | "notice" | "paragraph"
  | "heading" | "key_value" | "table" | "signature" | "declaration" | "page_break" | "footer";

interface KVRow { number: string; label: string; value: string; }
interface ColDef { key: string; header: string; }
interface Signatory { label: string; name: string; }

interface Section {
  id: string;
  type: SectionType;
  // common
  text?: string;
  bold?: boolean;
  italic?: boolean;
  align?: string;
  size?: number;
  level?: number;
  // key_value
  label_width?: number;
  rows?: KVRow[];
  // table
  headers?: string[];
  columns?: ColDef[];
  rows_key?: string;
  // signature
  signatories?: Signatory[];
  // declaration
  items?: string[];
}

interface TemplateDef {
  id: string;
  font: string;
  nav_color: string;
  nav_light_color: string;
  company: { name: string; address: string; web: string };
  sections: Section[];
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

function uid() { return Math.random().toString(36).slice(2, 9); }

function emptySection(type: SectionType): Section {
  const base: Section = { id: uid(), type };
  switch (type) {
    case "key_value": return { ...base, label_width: 3.0, rows: [{ number: "1.", label: "Field Label", value: "{{ field_key }}" }] };
    case "table": return { ...base, headers: ["Column 1", "Column 2"], columns: [{ key: "col1", header: "Column 1" }, { key: "col2", header: "Column 2" }], rows_key: "table_data" };
    case "signature": return { ...base, signatories: [{ label: "Signature & Date", name: "{{ proposer_name }}" }, { label: "Authorised Signatory", name: "" }] };
    case "declaration": return { ...base, items: ["Declaration statement 1.", "Declaration statement 2."] };
    case "heading": return { ...base, text: "Section Heading", level: 1 };
    case "title": return { ...base, text: "Document Title" };
    case "subtitle": return { ...base, text: "Subtitle" };
    case "notice": return { ...base, text: "Notice or important information box." };
    case "paragraph": return { ...base, text: "Paragraph text." };
    case "footer": return { ...base, text: "Ref: {{ proposal_id }}   |   Generated: {{ generated_at }}" };
    default: return base;
  }
}

const SECTION_LABELS: Record<SectionType, string> = {
  header: "Company Header (Logo + Name)",
  divider: "Horizontal Divider",
  title: "Document Title",
  subtitle: "Subtitle",
  notice: "Notice / Info Box",
  paragraph: "Paragraph",
  heading: "Section Heading",
  key_value: "Key-Value Fields Table",
  table: "Data Table (from runtime data)",
  signature: "Signature Block",
  declaration: "Declaration List",
  page_break: "Page Break",
  footer: "Footer Reference Line",
};

const SECTION_TYPES = Object.keys(SECTION_LABELS) as SectionType[];

// ─── Sub-editors ─────────────────────────────────────────────────────────────

function KVEditor({ rows, onChange }: { rows: KVRow[]; onChange: (r: KVRow[]) => void }) {
  return (
    <div className="space-y-2">
      <p className="text-[10px] font-semibold uppercase tracking-wider text-[var(--muted)]">Fields</p>
      {rows.map((row, i) => (
        <div key={i} className="grid grid-cols-[56px_1fr_1fr_28px] gap-1 items-center">
          <input className="portal-input text-xs" placeholder="#" value={row.number}
            onChange={e => { const r = [...rows]; r[i] = { ...r[i], number: e.target.value }; onChange(r); }} />
          <input className="portal-input text-xs" placeholder="Label" value={row.label}
            onChange={e => { const r = [...rows]; r[i] = { ...r[i], label: e.target.value }; onChange(r); }} />
          <input className="portal-input text-xs" placeholder="{{ data_key }}" value={row.value}
            onChange={e => { const r = [...rows]; r[i] = { ...r[i], value: e.target.value }; onChange(r); }} />
          <button type="button" onClick={() => onChange(rows.filter((_, j) => j !== i))}
            className="flex h-7 w-7 items-center justify-center rounded-lg text-[var(--danger)] hover:bg-[var(--danger-soft)]">
            <Trash2 className="h-3.5 w-3.5" />
          </button>
        </div>
      ))}
      <button type="button" className="portal-btn portal-btn-secondary text-xs"
        onClick={() => onChange([...rows, { number: `${rows.length + 1}.`, label: "New Field", value: "{{ new_key }}" }])}>
        <Plus className="h-3.5 w-3.5" /> Add Field
      </button>
    </div>
  );
}

function TableEditor({ sec, onChange }: { sec: Section; onChange: (s: Section) => void }) {
  const cols = sec.columns ?? [];
  const headers = sec.headers ?? [];
  return (
    <div className="space-y-3">
      <div className="space-y-1">
        <p className="text-[10px] font-semibold uppercase tracking-wider text-[var(--muted)]">Data Key (from runtime data)</p>
        <input className="portal-input text-xs" placeholder="rows_key e.g. items" value={sec.rows_key ?? ""}
          onChange={e => onChange({ ...sec, rows_key: e.target.value })} />
        <p className="text-[10px] text-[var(--muted)]">Your API must return data[rows_key] as an array of objects.</p>
      </div>
      <div className="space-y-1">
        <p className="text-[10px] font-semibold uppercase tracking-wider text-[var(--muted)]">Columns</p>
        {cols.map((col, i) => (
          <div key={i} className="grid grid-cols-[1fr_1fr_28px] gap-1 items-center">
            <input className="portal-input text-xs" placeholder="Column header" value={headers[i] ?? col.header}
              onChange={e => {
                const h = [...headers]; h[i] = e.target.value;
                const c = [...cols]; c[i] = { ...c[i], header: e.target.value };
                onChange({ ...sec, headers: h, columns: c });
              }} />
            <input className="portal-input text-xs" placeholder="data key" value={col.key}
              onChange={e => {
                const c = [...cols]; c[i] = { ...c[i], key: e.target.value };
                onChange({ ...sec, columns: c });
              }} />
            <button type="button" onClick={() => {
              const c = cols.filter((_, j) => j !== i);
              const h = headers.filter((_, j) => j !== i);
              onChange({ ...sec, columns: c, headers: h });
            }} className="flex h-7 w-7 items-center justify-center rounded-lg text-[var(--danger)] hover:bg-[var(--danger-soft)]">
              <Trash2 className="h-3.5 w-3.5" />
            </button>
          </div>
        ))}
        <button type="button" className="portal-btn portal-btn-secondary text-xs"
          onClick={() => {
            const n = `col${cols.length + 1}`;
            onChange({ ...sec, columns: [...cols, { key: n, header: `Column ${cols.length + 1}` }], headers: [...headers, `Column ${cols.length + 1}`] });
          }}>
          <Plus className="h-3.5 w-3.5" /> Add Column
        </button>
      </div>
    </div>
  );
}

function SigEditor({ sigs, onChange }: { sigs: Signatory[]; onChange: (s: Signatory[]) => void }) {
  return (
    <div className="space-y-2">
      <p className="text-[10px] font-semibold uppercase tracking-wider text-[var(--muted)]">Signatories</p>
      {sigs.map((sig, i) => (
        <div key={i} className="grid grid-cols-[1fr_1fr_28px] gap-1 items-center">
          <input className="portal-input text-xs" placeholder="Label" value={sig.label}
            onChange={e => { const s = [...sigs]; s[i] = { ...s[i], label: e.target.value }; onChange(s); }} />
          <input className="portal-input text-xs" placeholder="{{ name_key }} or blank" value={sig.name}
            onChange={e => { const s = [...sigs]; s[i] = { ...s[i], name: e.target.value }; onChange(s); }} />
          <button type="button" onClick={() => onChange(sigs.filter((_, j) => j !== i))}
            className="flex h-7 w-7 items-center justify-center rounded-lg text-[var(--danger)] hover:bg-[var(--danger-soft)]">
            <Trash2 className="h-3.5 w-3.5" />
          </button>
        </div>
      ))}
      <button type="button" className="portal-btn portal-btn-secondary text-xs"
        onClick={() => onChange([...sigs, { label: "Signatory", name: "" }])}>
        <Plus className="h-3.5 w-3.5" /> Add Signatory
      </button>
    </div>
  );
}

function DeclEditor({ items, onChange }: { items: string[]; onChange: (items: string[]) => void }) {
  return (
    <div className="space-y-2">
      <p className="text-[10px] font-semibold uppercase tracking-wider text-[var(--muted)]">Declaration Items</p>
      {items.map((item, i) => (
        <div key={i} className="flex gap-1 items-center">
          <input className="portal-input flex-1 text-xs" value={item}
            onChange={e => { const a = [...items]; a[i] = e.target.value; onChange(a); }} />
          <button type="button" onClick={() => onChange(items.filter((_, j) => j !== i))}
            className="flex h-7 w-7 items-center justify-center rounded-lg text-[var(--danger)] hover:bg-[var(--danger-soft)]">
            <Trash2 className="h-3.5 w-3.5" />
          </button>
        </div>
      ))}
      <button type="button" className="portal-btn portal-btn-secondary text-xs"
        onClick={() => onChange([...items, "New declaration item."])}>
        <Plus className="h-3.5 w-3.5" /> Add Item
      </button>
    </div>
  );
}

// ─── Section Card ─────────────────────────────────────────────────────────────

function SectionCard({ sec, index, total, onChange, onMove, onRemove }: {
  sec: Section; index: number; total: number;
  onChange: (s: Section) => void;
  onMove: (dir: "up" | "down") => void;
  onRemove: () => void;
}) {
  const [open, setOpen] = useState(false);

  return (
    <div className="rounded-[18px] border border-[rgb(12_91_65_/_0.10)] bg-white/80">
      <div className="flex items-center gap-2 px-4 py-2">
        <span className="flex h-6 w-6 items-center justify-center rounded-full bg-[rgb(12_91_65_/_0.08)] text-[10px] font-bold text-[var(--accent)]">
          {index + 1}
        </span>
        <span className="flex-1 text-xs font-semibold text-[var(--text)]">{SECTION_LABELS[sec.type]}</span>
        <div className="flex items-center gap-1">
          <button type="button" disabled={index === 0} onClick={() => onMove("up")}
            className="flex h-6 w-6 items-center justify-center rounded text-[var(--muted)] hover:text-[var(--accent)] disabled:opacity-30">
            <ChevronUp className="h-3.5 w-3.5" />
          </button>
          <button type="button" disabled={index === total - 1} onClick={() => onMove("down")}
            className="flex h-6 w-6 items-center justify-center rounded text-[var(--muted)] hover:text-[var(--accent)] disabled:opacity-30">
            <ChevronDown className="h-3.5 w-3.5" />
          </button>
          <button type="button" onClick={() => setOpen(!open)}
            className="flex h-6 w-6 items-center justify-center rounded text-[var(--muted)] hover:text-[var(--accent)]">
            {open ? <ChevronUp className="h-3.5 w-3.5" /> : <ChevronDown className="h-3.5 w-3.5" />}
          </button>
          <button type="button" onClick={onRemove}
            className="flex h-6 w-6 items-center justify-center rounded text-[var(--danger)] hover:bg-[var(--danger-soft)]">
            <Trash2 className="h-3.5 w-3.5" />
          </button>
        </div>
      </div>

      {open && (
        <div className="border-t border-[rgb(12_91_65_/_0.08)] px-4 py-3 space-y-3">
          {/* Text sections */}
          {["title", "subtitle", "notice", "paragraph", "footer"].includes(sec.type) && (
            <div>
              <p className="text-[10px] font-semibold uppercase tracking-wider text-[var(--muted)]">Text</p>
              <textarea className="portal-input mt-1 w-full text-xs" rows={3} value={sec.text ?? ""}
                onChange={e => onChange({ ...sec, text: e.target.value })} />
              <div className="mt-2 flex gap-2">
                <label className="flex items-center gap-1 text-xs text-[var(--muted)]">
                  <input type="checkbox" checked={!!sec.bold} onChange={e => onChange({ ...sec, bold: e.target.checked })} /> Bold
                </label>
                <label className="flex items-center gap-1 text-xs text-[var(--muted)]">
                  <input type="checkbox" checked={!!sec.italic} onChange={e => onChange({ ...sec, italic: e.target.checked })} /> Italic
                </label>
                <select className="portal-select text-xs" value={sec.align ?? "left"} onChange={e => onChange({ ...sec, align: e.target.value })}>
                  <option value="left">Left</option>
                  <option value="center">Center</option>
                  <option value="right">Right</option>
                  <option value="justify">Justify</option>
                </select>
              </div>
            </div>
          )}

          {sec.type === "heading" && (
            <div className="grid grid-cols-[1fr_100px] gap-2">
              <div>
                <p className="text-[10px] font-semibold uppercase tracking-wider text-[var(--muted)]">Heading Text</p>
                <input className="portal-input mt-1 w-full text-xs" value={sec.text ?? ""}
                  onChange={e => onChange({ ...sec, text: e.target.value })} />
              </div>
              <div>
                <p className="text-[10px] font-semibold uppercase tracking-wider text-[var(--muted)]">Level</p>
                <select className="portal-select mt-1 text-xs" value={sec.level ?? 1} onChange={e => onChange({ ...sec, level: Number(e.target.value) })}>
                  <option value={1}>1 (Large)</option>
                  <option value={2}>2 (Medium)</option>
                </select>
              </div>
            </div>
          )}

          {sec.type === "key_value" && (
            <>
              <div>
                <p className="text-[10px] font-semibold uppercase tracking-wider text-[var(--muted)]">Label Column Width (inches)</p>
                <input type="number" className="portal-input mt-1 w-28 text-xs" step={0.25} min={1} max={5}
                  value={sec.label_width ?? 3.0} onChange={e => onChange({ ...sec, label_width: Number(e.target.value) })} />
              </div>
              <KVEditor rows={sec.rows ?? []} onChange={rows => onChange({ ...sec, rows })} />
            </>
          )}

          {sec.type === "table" && <TableEditor sec={sec} onChange={onChange} />}
          {sec.type === "signature" && <SigEditor sigs={sec.signatories ?? []} onChange={sigs => onChange({ ...sec, signatories: sigs })} />}
          {sec.type === "declaration" && <DeclEditor items={sec.items ?? []} onChange={items => onChange({ ...sec, items })} />}
        </div>
      )}
    </div>
  );
}

// ─── Main Component ───────────────────────────────────────────────────────────

export function TemplateCreator({ onClose }: { onClose?: () => void }) {
  const [templateId, setTemplateId] = useState("my_new_template");
  const [companyName, setCompanyName] = useState("PRAGATI INSURANCE PLC");
  const [companyAddress, setCompanyAddress] = useState("20-21 Kawran Bazar, Dhaka-1215");
  const [companyWeb, setCompanyWeb] = useState("info@pragatiinsurance.com");
  const [navColor, setNavColor] = useState("#1F3864");
  const [sections, setSections] = useState<Section[]>([
    emptySection("header"),
    emptySection("divider"),
    emptySection("title"),
    emptySection("key_value"),
    emptySection("signature"),
    emptySection("footer"),
  ]);
  const [addType, setAddType] = useState<SectionType>("key_value");
  const [saving, setSaving] = useState(false);
  const [saveMsg, setSaveMsg] = useState<string | null>(null);
  const [previewing, setPreviewing] = useState(false);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [previewDownload, setPreviewDownload] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  function buildDefinition(): TemplateDef {
    return {
      id: templateId.replace(/[^a-z0-9_]/gi, "_").toLowerCase(),
      font: "Times New Roman",
      nav_color: navColor.replace("#", ""),
      nav_light_color: "E9EFF7",
      company: { name: companyName, address: companyAddress, web: companyWeb },
      sections: sections.map(({ id: _id, ...rest }) => rest) as Section[],
    };
  }

  function moveSection(index: number, dir: "up" | "down") {
    const arr = [...sections];
    const target = dir === "up" ? index - 1 : index + 1;
    if (target < 0 || target >= arr.length) return;
    [arr[index], arr[target]] = [arr[target], arr[index]];
    setSections(arr);
  }

  async function handleSave() {
    setSaving(true); setError(null); setSaveMsg(null);
    try {
      const res = await fetch("/api/documents/templates/definitions", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ id: templateId, definition: buildDefinition() }),
      });
      const json = await res.json() as { ok: boolean; message?: string };
      if (json.ok) setSaveMsg(`Saved as ${templateId}.json — ready to use in any document card.`);
      else setError(json.message ?? "Save failed");
    } catch { setError("Network error"); }
    finally { setSaving(false); }
  }

  async function handlePreview() {
    setPreviewing(true); setError(null); setPreviewUrl(null);
    try {
      const res = await fetch("/api/documents/templates/preview", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ definition: buildDefinition(), sampleData: {} }),
      });
      const json = await res.json() as { ok: boolean; data?: { renderUrl: string; downloadUrl: string }; message?: string };
      if (json.ok && json.data) {
        setPreviewUrl(json.data.renderUrl);
        setPreviewDownload(json.data.downloadUrl);
      } else {
        setError(json.message ?? "Preview failed");
      }
    } catch { setError("Network error"); }
    finally { setPreviewing(false); }
  }

  return (
    <div className="document-modal-backdrop" data-document-modal="true">
      <div className="document-modal-shell" style={{ maxWidth: 960 }}>
        {/* Header */}
        <div className="document-modal-header" data-print-hide="true">
          <div>
            <p className="text-xs font-semibold uppercase tracking-[0.18em] text-[var(--muted)]">Document Template Builder</p>
            <h2 className="mt-2 font-[family:var(--font-heading)] text-3xl font-semibold text-[var(--text)]">Create New Template</h2>
          </div>
          <div className="flex gap-2">
            <button className="portal-btn portal-btn-secondary" disabled={previewing} onClick={() => void handlePreview()} type="button">
              {previewing ? <LoaderCircle className="h-4 w-4 animate-spin" /> : <Eye className="h-4 w-4" />}
              {previewing ? "Generating…" : "Preview DOCX"}
            </button>
            <button className="portal-btn portal-btn-primary" disabled={saving} onClick={() => void handleSave()} type="button">
              {saving ? <LoaderCircle className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
              {saving ? "Saving…" : "Save Template"}
            </button>
            {onClose && (
              <button className="portal-btn portal-btn-secondary" onClick={onClose} type="button">
                <X className="h-4 w-4" /> Close
              </button>
            )}
          </div>
        </div>

        <div className="document-modal-body">
          <div className="grid gap-4 lg:grid-cols-[320px_1fr]">
            {/* Left: Config */}
            <div className="space-y-4">
              <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/88 p-4 space-y-3">
                <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">Template Identity</p>
                <label className="block space-y-1">
                  <span className="text-[10px] font-semibold uppercase text-[var(--muted)]">Template ID (no spaces)</span>
                  <input className="portal-input w-full text-sm" placeholder="my_template_name" value={templateId}
                    onChange={e => setTemplateId(e.target.value)} />
                  <span className="text-[10px] text-[var(--muted)]">Saved as {templateId.replace(/[^a-z0-9_]/gi,"_").toLowerCase()}.json</span>
                </label>
                <label className="block space-y-1">
                  <span className="text-[10px] font-semibold uppercase text-[var(--muted)]">Header Color</span>
                  <div className="flex items-center gap-2">
                    <input type="color" className="h-8 w-12 rounded border" value={navColor} onChange={e => setNavColor(e.target.value)} />
                    <input className="portal-input flex-1 text-xs" value={navColor} onChange={e => setNavColor(e.target.value)} />
                  </div>
                </label>
              </div>

              <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/88 p-4 space-y-3">
                <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">Company Details</p>
                <label className="block space-y-1">
                  <span className="text-[10px] font-semibold uppercase text-[var(--muted)]">Company Name</span>
                  <input className="portal-input w-full text-xs" value={companyName} onChange={e => setCompanyName(e.target.value)} />
                </label>
                <label className="block space-y-1">
                  <span className="text-[10px] font-semibold uppercase text-[var(--muted)]">Address Line</span>
                  <input className="portal-input w-full text-xs" value={companyAddress} onChange={e => setCompanyAddress(e.target.value)} />
                </label>
                <label className="block space-y-1">
                  <span className="text-[10px] font-semibold uppercase text-[var(--muted)]">Website / Email</span>
                  <input className="portal-input w-full text-xs" value={companyWeb} onChange={e => setCompanyWeb(e.target.value)} />
                </label>
              </div>

              <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/88 p-4 space-y-3">
                <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">Add Section</p>
                <div className="flex gap-2">
                  <select className="portal-select flex-1 text-xs" value={addType} onChange={e => setAddType(e.target.value as SectionType)}>
                    {SECTION_TYPES.map(t => <option key={t} value={t}>{SECTION_LABELS[t]}</option>)}
                  </select>
                  <button type="button" className="portal-btn portal-btn-primary shrink-0"
                    onClick={() => setSections([...sections, emptySection(addType)])}>
                    <Plus className="h-4 w-4" />
                  </button>
                </div>
                <p className="text-[10px] text-[var(--muted)]">
                  Use {"{{ data_key }}"} placeholders — they are filled from runtime data when generating.
                </p>
              </div>

              {/* Status messages */}
              {saveMsg && (
                <div className="rounded-[18px] bg-[rgb(12_91_65_/_0.08)] px-4 py-3 text-sm text-[var(--accent)]">
                  <p className="font-semibold">✓ Template saved!</p>
                  <p className="mt-1 text-xs">{saveMsg}</p>
                </div>
              )}
              {error && (
                <div className="rounded-[18px] border border-[rgb(194_65_12_/_0.14)] bg-[var(--danger-soft)] px-4 py-3 text-sm text-[var(--danger)]">
                  {error}
                </div>
              )}
              {previewDownload && (
                <a className="portal-btn portal-btn-secondary w-full justify-center" href={previewDownload} download>
                  <Download className="h-4 w-4" /> Download Preview DOCX
                </a>
              )}
            </div>

            {/* Right: Sections + Preview */}
            <div className="space-y-3">
              {previewUrl && (
                <iframe
                  key={previewUrl}
                  src={previewUrl}
                  className="w-full rounded-[20px] border border-[rgb(12_91_65_/_0.10)] bg-white"
                  style={{ height: "45vh" }}
                  title="Template Preview"
                />
              )}

              <div className="space-y-2">
                {sections.map((sec, i) => (
                  <SectionCard
                    key={sec.id}
                    sec={sec}
                    index={i}
                    total={sections.length}
                    onChange={updated => setSections(sections.map((s, j) => j === i ? updated : s))}
                    onMove={dir => moveSection(i, dir)}
                    onRemove={() => setSections(sections.filter((_, j) => j !== i))}
                  />
                ))}
              </div>

              <button type="button" className="portal-btn portal-btn-secondary w-full"
                onClick={() => setSections([...sections, emptySection(addType)])}>
                <Plus className="h-4 w-4" /> Add {SECTION_LABELS[addType]}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
