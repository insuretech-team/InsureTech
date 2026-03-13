/**
 * docgen-sdk-client.ts
 * ─────────────────────
 * Server-side docgen helpers for Next.js API route handlers.
 *
 * The docgen service is NOT in the generated SDK (@lifeplus/insuretech-sdk).
 * It is exposed by the gateway under /v1/documents and /v1/document-templates.
 * We use makeDirectHttp (same cookie/CSRF auth as makeSdkClient) for all calls.
 *
 * Types are imported from src/lib/proto-generated — the canonical generated types.
 * The gateway serialises proto responses as JSON with snake_case field names
 * (protojson UseProtoNames=true), so field names on wire match the proto field names
 * (e.g. document_id, file_url, template_id, entity_type…).
 */

import { makeDirectHttp } from "./b2b-sdk-client";
import type {
  GenerateDocumentResponse,
  GetDocumentResponse,
  ListDocumentsResponse,
  DownloadDocumentResponse,
  DeleteDocumentResponse,
  CreateDocumentTemplateResponse,
  GetDocumentTemplateResponse,
  ListDocumentTemplatesResponse,
  UpdateDocumentTemplateResponse,
  DeactivateDocumentTemplateResponse,
  DeleteDocumentTemplateResponse,
} from "@lib/proto-generated/insuretech/document/services/v1/document_service_pb";

// Re-export proto types for convenience so callers only need this file.
export type {
  GenerateDocumentResponse,
  GetDocumentResponse,
  ListDocumentsResponse,
  DownloadDocumentResponse,
  DeleteDocumentResponse,
  CreateDocumentTemplateResponse,
  GetDocumentTemplateResponse,
  ListDocumentTemplatesResponse,
  UpdateDocumentTemplateResponse,
  DeactivateDocumentTemplateResponse,
  DeleteDocumentTemplateResponse,
};

export type {
  DocumentGeneration,
  GenerationStatus,
} from "@lib/proto-generated/insuretech/document/entity/v1/document_generation_pb";

export type {
  DocumentTemplate,
  DocumentType,
  OutputFormat,
} from "@lib/proto-generated/insuretech/document/entity/v1/document_template_pb";

// ─── Payload shapes (what we POST to the gateway — snake_case to match proto JSON) ──

export interface GenerateDocumentPayload {
  template_id: string;
  entity_type: string;
  entity_id: string;
  /** Key/value pairs merged into the template */
  data?: Record<string, unknown>;
  include_qr_code?: boolean;
}

export interface CreateDocumentTemplatePayload {
  name: string;
  type: string;
  description?: string;
  template_content: string;
  output_format: string;
  variables?: string[];
}

export interface UpdateDocumentTemplatePayload {
  template?: {
    name?: string;
    description?: string;
    template_content?: string;
    output_format?: string;
    is_active?: boolean;
  };
}

// ─── Factory ──────────────────────────────────────────────────────────────────

/**
 * makeDocgenClient — wraps makeDirectHttp with typed docgen helpers.
 * Call from API route handlers (server-side only).
 *
 * Usage:
 *   const docgen = makeDocgenClient(request, hdrs ?? undefined);
 *   const res = await docgen.generate({ template_id: "...", entity_type: "employee", entity_id: "..." });
 */
export function makeDocgenClient(
  request: Request,
  sessionOverrides?: { portal?: string; userId?: string; businessId?: string; tenantId?: string }
) {
  const http = makeDirectHttp(request, sessionOverrides);

  return {
    // ── Documents ────────────────────────────────────────────────────────────

    async generate(payload: GenerateDocumentPayload) {
      return http.post("/v1/documents", payload) as Promise<{
        ok: boolean; status: number;
        data: Partial<GenerateDocumentResponse> & Record<string, unknown>;
      }>;
    },

    async getDocument(documentId: string) {
      return http.get(`/v1/documents/${documentId}`) as Promise<{
        ok: boolean; status: number;
        data: Partial<GetDocumentResponse> & Record<string, unknown>;
      }>;
    },

    async listDocuments(options: {
      entityType: string;
      entityId: string;
      status?: string;
      page?: number;
      pageSize?: number;
    }) {
      const params = new URLSearchParams({
        entity_type: options.entityType,
        entity_id: options.entityId,
      });
      if (options.status) params.set("status", options.status);
      if (options.page) params.set("page", String(options.page));
      if (options.pageSize) params.set("page_size", String(options.pageSize));
      return http.get(`/v1/entities/${options.entityType}/${options.entityId}/documents?${params}`) as Promise<{
        ok: boolean; status: number;
        data: Partial<ListDocumentsResponse> & Record<string, unknown>;
      }>;
    },

    async downloadDocument(documentId: string) {
      return http.get(`/v1/documents/${documentId}/download`) as Promise<{
        ok: boolean; status: number;
        data: Partial<DownloadDocumentResponse> & Record<string, unknown>;
      }>;
    },

    async deleteDocument(documentId: string) {
      return http.delete(`/v1/documents/${documentId}`) as Promise<{
        ok: boolean; status: number;
        data: Partial<DeleteDocumentResponse> & Record<string, unknown>;
      }>;
    },

    // ── Templates ────────────────────────────────────────────────────────────

    async createTemplate(payload: CreateDocumentTemplatePayload) {
      return http.post("/v1/document-templates", payload) as Promise<{
        ok: boolean; status: number;
        data: Partial<CreateDocumentTemplateResponse> & Record<string, unknown>;
      }>;
    },

    async getTemplate(templateId: string) {
      return http.get(`/v1/document-templates/${templateId}`) as Promise<{
        ok: boolean; status: number;
        data: Partial<GetDocumentTemplateResponse> & Record<string, unknown>;
      }>;
    },

    async listTemplates(options?: {
      type?: string;
      activeOnly?: boolean;
      pageSize?: number;
      pageToken?: string;
    }) {
      const params = new URLSearchParams();
      if (options?.type) params.set("type", options.type);
      if (options?.activeOnly) params.set("active_only", "true");
      if (options?.pageSize) params.set("page_size", String(options.pageSize));
      if (options?.pageToken) params.set("page_token", options.pageToken);
      const qs = params.toString() ? `?${params}` : "";
      return http.get(`/v1/document-templates${qs}`) as Promise<{
        ok: boolean; status: number;
        data: Partial<ListDocumentTemplatesResponse> & Record<string, unknown>;
      }>;
    },

    async updateTemplate(templateId: string, payload: UpdateDocumentTemplatePayload) {
      return http.patch(`/v1/document-templates/${templateId}`, payload) as Promise<{
        ok: boolean; status: number;
        data: Partial<UpdateDocumentTemplateResponse> & Record<string, unknown>;
      }>;
    },

    async deactivateTemplate(templateId: string, reason?: string) {
      return http.post(`/v1/document-templates/${templateId}/deactivate`, { reason: reason ?? "" }) as Promise<{
        ok: boolean; status: number;
        data: Partial<DeactivateDocumentTemplateResponse> & Record<string, unknown>;
      }>;
    },

    async deleteTemplate(templateId: string) {
      return http.delete(`/v1/document-templates/${templateId}`) as Promise<{
        ok: boolean; status: number;
        data: Partial<DeleteDocumentTemplateResponse> & Record<string, unknown>;
      }>;
    },
  };
}

export type DocgenClient = ReturnType<typeof makeDocgenClient>;
