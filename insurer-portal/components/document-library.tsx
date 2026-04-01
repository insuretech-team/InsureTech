"use client";

import {
  Download,
  FileOutput,
  FilePlus2,
  LoaderCircle,
  Printer,
  RefreshCw,
  Search,
  Sparkles,
  X,
} from "lucide-react";
import { useEffect, useState } from "react";

import { api } from "@/lib/browser-client";
import { useLibraryDocuments } from "@/hooks/use-library-documents";
import { TemplateCreator } from "@/components/template-creator";
import type { LibraryDocument } from "@/lib/types";

export function DocumentLibrary() {
  const library = useLibraryDocuments();

  // Filters
  const [libQuery, setLibQuery] = useState("");
  const [libCategory, setLibCategory] = useState("All");
  const [libStage, setLibStage] = useState("All");

  const filteredLibDocs = library.documents.filter((doc) => {
    const q = libQuery.toLowerCase();
    const matchQ = !q || doc.title.toLowerCase().includes(q) || doc.category.toLowerCase().includes(q) || doc.summary.toLowerCase().includes(q);
    const matchCat = libCategory === "All" || doc.category === libCategory;
    const matchStage = libStage === "All" || doc.stage === libStage;
    return matchQ && matchCat && matchStage;
  });

  // Active modal document
  const [activeLibDoc, setActiveLibDoc] = useState<LibraryDocument | null>(null);

  // DOCX generation state
  const [cardGenLoading, setCardGenLoading] = useState(false);
  const [cardGenError, setCardGenError] = useState<string | null>(null);
  const [cardGenFilename, setCardGenFilename] = useState<string | null>(null);

  // Template creator
  const [showTemplateCreator, setShowTemplateCreator] = useState(false);

  // postMessage listener for editable iframe Save
  useEffect(() => {
    function onMessage(e: MessageEvent) {
      if (!e.data || e.data.type !== "DOC_SAVE") return;
      const doc = activeLibDoc;
      if (!doc) return;
      const editedFields = e.data.fields as Record<string, string>;
      function toSnakeKey(label: string): string {
        return label.toLowerCase().replace(/[^a-z0-9]+/g, "_").replace(/^_|_$/g, "");
      }
      const data: Record<string, unknown> = {
        proposal_id: doc.id,
        generated_at: new Date().toLocaleString("en-GB", { day: "2-digit", month: "short", year: "numeric", hour: "2-digit", minute: "2-digit" }),
      };
      for (const [label, value] of Object.entries(editedFields)) {
        data[toSnakeKey(label)] = value;
      }
      const safeTitle = doc.title.replace(/[^a-zA-Z0-9_]/g, "_");
      setCardGenLoading(true);
      setCardGenError(null);
      api.documents.generateFromCard({
        documentId: doc.id, category: doc.category, title: doc.title, kind: doc.kind,
        data, filename: `${safeTitle}_${Date.now()}`,
      }).then((res) => {
        if (res.ok && res.data) setCardGenFilename(res.data.filename);
        else setCardGenError(res.message ?? "Generation failed.");
      }).catch(() => {
        setCardGenError("Unable to reach the document service.");
      }).finally(() => setCardGenLoading(false));
    }
    window.addEventListener("message", onMessage);
    return () => window.removeEventListener("message", onMessage);
  }, [activeLibDoc]);

  async function handleGenerateFromCard(doc: LibraryDocument) {
    setCardGenLoading(true);
    setCardGenError(null);
    setCardGenFilename(null);
    const safeTitle = doc.title.replace(/[^a-zA-Z0-9_]/g, "_");
    try {
      const res = await api.documents.generateFromCard({
        documentId: doc.id, category: doc.category, title: doc.title, kind: doc.kind,
        data: {
          proposal_id: doc.id,
          generated_at: new Date().toLocaleString("en-GB", { day: "2-digit", month: "short", year: "numeric", hour: "2-digit", minute: "2-digit" }),
        },
        filename: `${safeTitle}_${Date.now()}`,
      });
      if (res.ok && res.data) setCardGenFilename(res.data.filename);
      else setCardGenError(res.message ?? "Generation failed.");
    } catch {
      setCardGenError("Unable to reach the document service.");
    } finally {
      setCardGenLoading(false);
    }
  }

  function openCard(doc: LibraryDocument) {
    setActiveLibDoc(doc);
    setCardGenLoading(false);
    setCardGenError(null);
    setCardGenFilename(null);
    void handleGenerateFromCard(doc);
  }

  function closeCard() {
    setActiveLibDoc(null);
    setCardGenFilename(null);
    setCardGenError(null);
  }

  return (
    <div className="page-documents space-y-6">

      {/* ── Hero ── */}
      <section className="document-library-hero" data-print-hide="true">
        <div>
          <p className="text-xs font-semibold uppercase tracking-[0.18em] text-[var(--muted)]">Document Centre</p>
          <h1 className="mt-3 font-[family:var(--font-heading)] text-4xl font-semibold text-[var(--text)]">Document Library</h1>
          <p className="mt-3 max-w-3xl text-sm leading-7 text-[var(--muted)]">
            Click any document to generate and view the real DOCX — edit inline and download. Create new templates with the Template Builder.
          </p>
        </div>
        <div className="flex flex-wrap gap-2">
          <button className="portal-btn portal-btn-primary" onClick={() => setShowTemplateCreator(true)} type="button">
            <FilePlus2 className="h-4 w-4" />
            Create Template
          </button>
          <button className="portal-btn portal-btn-secondary" onClick={library.refresh} type="button">
            <RefreshCw className={`h-4 w-4 ${library.loading ? "animate-spin" : ""}`} />
            Refresh
          </button>
        </div>
      </section>

      {/* ── Filters ── */}
      <section className="grid gap-3 md:grid-cols-[1fr_180px_180px]" data-print-hide="true">
        <div className="flex items-center gap-2 rounded-[20px] border border-[rgb(12_91_65_/_0.08)] bg-white/82 px-3 py-2.5">
          <Search className="h-4 w-4 shrink-0 text-[var(--muted)]" />
          <input
            className="w-full bg-transparent text-sm outline-none placeholder:text-[var(--muted)]"
            placeholder="Search documents…"
            value={libQuery}
            onChange={(e) => setLibQuery(e.target.value)}
          />
        </div>
        <select className="portal-select" value={libCategory} onChange={(e) => setLibCategory(e.target.value)}>
          {library.categoryOptions.map((c) => <option key={c} value={c}>{c}</option>)}
        </select>
        <select className="portal-select" value={libStage} onChange={(e) => setLibStage(e.target.value)}>
          {library.stageOptions.map((s) => <option key={s} value={s}>{s}</option>)}
        </select>
      </section>

      {/* ── Loading / Error ── */}
      {library.loading && (
        <div className="flex items-center justify-center gap-3 rounded-[28px] border border-[rgb(12_91_65_/_0.08)] bg-white/60 py-12 text-sm text-[var(--muted)]">
          <LoaderCircle className="h-5 w-5 animate-spin text-[var(--accent)]" /> Loading document library…
        </div>
      )}
      {library.error && !library.loading && (
        <div className="rounded-[20px] border border-[rgb(194_65_12_/_0.14)] bg-[var(--danger-soft)] px-4 py-3 text-sm text-[var(--danger)]">
          {library.error}
        </div>
      )}

      {/* ── Card Grid ── */}
      {!library.loading && filteredLibDocs.length > 0 && (
        <section className="document-library-grid" data-print-hide="true">
          {filteredLibDocs.map((doc) => (
            <button
              key={doc.id}
              className="document-card group text-left"
              onClick={() => openCard(doc)}
              type="button"
            >
              <div className="document-card-header">
                <div className="flex items-start justify-between gap-2">
                  <div>
                    <p className="text-[10px] font-semibold uppercase tracking-[0.14em] text-[var(--muted)]">{doc.category}</p>
                    <p className="mt-1 font-semibold leading-snug text-[var(--text)]">{doc.title}</p>
                  </div>
                  <div className="flex shrink-0 flex-col items-end gap-1">
                    <span className={`pill ${doc.stage === "Claims" ? "pill-danger" : doc.stage === "Pricing" || doc.stage === "Reference" ? "pill-warn" : "pill-live"}`}>
                      {doc.stage}
                    </span>
                    {doc.isGenerated && <span className="pill pill-neutral text-[9px]">DOCX</span>}
                  </div>
                </div>
              </div>
              <div className="document-card-body">
                <p className="line-clamp-3 text-xs leading-5 text-[var(--muted)]">{doc.summary}</p>
              </div>
              <div className="document-card-footer">
                <span className="text-[10px] text-[var(--muted)]">{doc.kind.replace(/-/g, " ")}</span>
                <Sparkles className="h-3.5 w-3.5 text-[var(--accent)] opacity-60 group-hover:opacity-100" />
              </div>
            </button>
          ))}
        </section>
      )}

      {!library.loading && filteredLibDocs.length === 0 && !library.error && (
        <div className="rounded-[28px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-8 text-sm text-[var(--muted)]">
          No documents found. Try adjusting the filters or click Create Template to add one.
        </div>
      )}

      {/* ── Template Creator Modal ── */}
      {showTemplateCreator && <TemplateCreator onClose={() => { setShowTemplateCreator(false); library.refresh(); }} />}

      {/* ── Document Modal ── */}
      {activeLibDoc && (
        <div className="document-modal-backdrop" data-document-modal="true">
          <div className="document-modal-shell">
            <div className="document-modal-header" data-print-hide="true">
              <div>
                <p className="text-xs font-semibold uppercase tracking-[0.18em] text-[var(--muted)]">{activeLibDoc.category} · {activeLibDoc.stage}</p>
                <h2 className="mt-2 font-[family:var(--font-heading)] text-3xl font-semibold text-[var(--text)]">{activeLibDoc.title}</h2>
                <p className="mt-1 text-sm text-[var(--muted)]">{activeLibDoc.summary}</p>
              </div>
              <div className="flex items-center gap-2" data-print-hide="true">
                {cardGenFilename && (
                  <a className="portal-btn portal-btn-primary" href={api.documents.fileUrl(cardGenFilename)} download>
                    <Download className="h-4 w-4" /> Download DOCX
                  </a>
                )}
                <button
                  className="portal-btn portal-btn-secondary"
                  disabled={cardGenLoading}
                  onClick={() => void handleGenerateFromCard(activeLibDoc)}
                  type="button"
                >
                  {cardGenLoading ? <LoaderCircle className="h-4 w-4 animate-spin" /> : <RefreshCw className="h-4 w-4" />}
                  {cardGenLoading ? "Generating…" : "Regenerate"}
                </button>
                <button className="portal-btn portal-btn-secondary" onClick={closeCard} type="button">
                  <X className="h-4 w-4" /> Close
                </button>
              </div>
            </div>

            <div className="document-modal-body">
              <div className="grid gap-4 md:grid-cols-[200px_1fr]">

                {/* Sidebar */}
                <aside className="space-y-3" data-print-hide="true">
                  <div className="rounded-[20px] bg-[rgb(12_91_65_/_0.06)] p-4">
                    <p className="text-[10px] font-semibold uppercase tracking-[0.14em] text-[var(--muted)]">Document Info</p>
                    <div className="mt-2 flex flex-wrap gap-1">
                      <span className="pill pill-live">{activeLibDoc.stage}</span>
                      <span className="pill pill-neutral">{activeLibDoc.kind.replace(/-/g, " ")}</span>
                    </div>
                    <p className="mt-2 text-xs leading-5 text-[var(--muted)]">{activeLibDoc.suggestedUse}</p>
                  </div>

                  <div className="rounded-[20px] border border-[rgb(12_91_65_/_0.08)] bg-white/88 p-4">
                    <p className="text-[10px] font-semibold uppercase tracking-[0.14em] text-[var(--muted)]">How to Edit</p>
                    <p className="mt-2 text-xs leading-5 text-[var(--muted)]">
                      Click any <span className="rounded bg-[#fffef0] px-1 font-semibold text-[#1f3864]">yellow cell</span> in the document to edit its value.
                      A <strong>Save &amp; Regenerate</strong> bar will appear — click it to update the DOCX.
                    </p>
                  </div>

                  {cardGenFilename && (
                    <div className="rounded-[20px] border border-[rgb(12_91_65_/_0.08)] bg-white/88 p-4 space-y-2">
                      <p className="text-[10px] font-semibold uppercase tracking-[0.14em] text-[var(--muted)]">Actions</p>
                      <a className="portal-btn portal-btn-secondary w-full justify-start" href={api.documents.fileUrl(cardGenFilename)} download>
                        <Download className="h-4 w-4" /> Download DOCX
                      </a>
                      <button className="portal-btn portal-btn-secondary w-full justify-start" onClick={() => window.print()} type="button">
                        <Printer className="h-4 w-4" /> Print
                      </button>
                    </div>
                  )}
                </aside>

                {/* Main: inline DOCX render */}
                <main className="min-w-0 space-y-3">
                  {cardGenError && (
                    <div className="flex items-start justify-between gap-4 rounded-[18px] border border-[rgb(194_65_12_/_0.14)] bg-[var(--danger-soft)] px-4 py-3">
                      <p className="text-sm text-[var(--danger)]">{cardGenError}</p>
                      <button
                        className="portal-btn portal-btn-secondary shrink-0"
                        disabled={cardGenLoading}
                        onClick={() => void handleGenerateFromCard(activeLibDoc)}
                        type="button"
                      >
                        <RefreshCw className="h-4 w-4" /> Retry
                      </button>
                    </div>
                  )}

                  {cardGenLoading && (
                    <div className="flex items-center justify-center gap-3 rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/60 py-16 text-sm text-[var(--muted)]">
                      <LoaderCircle className="h-5 w-5 animate-spin text-[var(--accent)]" /> Rendering document…
                    </div>
                  )}

                  {cardGenFilename && !cardGenLoading && (
                    <iframe
                      key={cardGenFilename}
                      src={`/api/documents/render/${encodeURIComponent(cardGenFilename)}`}
                      className="w-full rounded-[20px] border border-[rgb(12_91_65_/_0.10)] bg-white"
                      style={{ minHeight: "72vh", height: "72vh" }}
                      title="Document Preview"
                    />
                  )}

                  {!cardGenFilename && !cardGenLoading && !cardGenError && (
                    <div className="flex flex-col items-center justify-center gap-4 rounded-[28px] border border-dashed border-[rgb(12_91_65_/_0.15)] bg-white/40 py-20">
                      <FileOutput className="h-10 w-10 text-[var(--muted)]" />
                      <div className="text-center">
                        <p className="text-sm font-semibold text-[var(--text)]">Preparing document…</p>
                        <p className="mt-1 text-xs text-[var(--muted)]">Click Regenerate if it does not appear.</p>
                      </div>
                    </div>
                  )}
                </main>

              </div>
            </div>
          </div>
        </div>
      )}

    </div>
  );
}
