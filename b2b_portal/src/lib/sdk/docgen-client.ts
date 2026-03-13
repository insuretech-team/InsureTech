/**
 * docgen-client.ts
 * ────────────────
 * Browser-side client for /api/documents and /api/document-templates.
 * Used by client components. Mirrors the pattern of employee-client.ts.
 */
import { parseJson, type ApiResult } from "./shared";

// ─── Types ────────────────────────────────────────────────────────────────────

export type DocumentStatus =
  | "DOCUMENT_STATUS_PENDING"
  | "DOCUMENT_STATUS_PROCESSING"
  | "DOCUMENT_STATUS_COMPLETED"
  | "DOCUMENT_STATUS_FAILED";

export interface DocumentRecord {
  document_id: string;
  template_id: string;
  entity_type: string;
  entity_id: string;
  document_type: string;
  status: DocumentStatus;
  file_url?: string;
  download_url?: string;
  created_at?: string;
  updated_at?: string;
}

export interface GenerateDocumentPayload {
  template_id: string;
  entity_type: string;
  entity_id: string;
  data?: Record<string, unknown>;
  include_qr_code?: boolean;
}

export type DocumentListResult = ApiResult<{ documents?: DocumentRecord[]; total?: number }>;
export type DocumentSingleResult = ApiResult<{ document?: DocumentRecord }>;
export type DocumentDownloadResult = ApiResult<{
  content?: string;       // base64-encoded file content
  content_type?: string;  // MIME type
  file_name?: string;
}>;

// ─── Client ───────────────────────────────────────────────────────────────────

export const docgenClient = {
  /** Generate a new document */
  async generate(payload: GenerateDocumentPayload): Promise<DocumentSingleResult> {
    const res = await fetch("/api/documents", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    return parseJson<DocumentSingleResult>(res);
  },

  /** List documents for a given entity */
  async list(options: {
    entityType: string;
    entityId: string;
    status?: string;
    page?: number;
    pageSize?: number;
  }): Promise<DocumentListResult> {
    const params = new URLSearchParams({
      entity_type: options.entityType,
      entity_id: options.entityId,
    });
    if (options.status) params.set("status", options.status);
    if (options.page) params.set("page", String(options.page));
    if (options.pageSize) params.set("page_size", String(options.pageSize));
    const res = await fetch(`/api/documents?${params}`, { cache: "no-store" });
    return parseJson<DocumentListResult>(res);
  },

  /** Get a single document by ID */
  async get(documentId: string): Promise<DocumentSingleResult> {
    const res = await fetch(`/api/documents/${documentId}`, { cache: "no-store" });
    return parseJson<DocumentSingleResult>(res);
  },

  /** Download a document (returns base64 content) */
  async download(documentId: string): Promise<DocumentDownloadResult> {
    const res = await fetch(`/api/documents/${documentId}/download`, { cache: "no-store" });
    return parseJson<DocumentDownloadResult>(res);
  },

  /** Delete a document */
  async delete(documentId: string): Promise<ApiResult> {
    const res = await fetch(`/api/documents/${documentId}`, { method: "DELETE" });
    return parseJson<ApiResult>(res);
  },
};
