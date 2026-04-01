/**
 * docgen-sdk-client.ts
 * ────────────────────
 * Server-side client factory for Document Service (docgen) API route handlers.
 *
 * The generated SDK does not yet include docgen service functions, so this
 * module uses direct HTTP calls to the gateway REST endpoints.
 *
 * Endpoints:
 *   GET    /v1/document-templates              → listTemplates
 *   POST   /v1/document-templates              → createTemplate
 *   GET    /v1/document-templates/{template_id} → getTemplate
 *   PATCH  /v1/document-templates/{template_id} → updateTemplate
 *   DELETE /v1/document-templates/{template_id} → deleteTemplate
 */

// ─── Types ────────────────────────────────────────────────────────────────────

export interface CreateDocumentTemplatePayload {
  name: string;
  type?: string;
  content?: string;
  description?: string;
  [key: string]: unknown;
}

export interface UpdateDocumentTemplatePayload {
  name?: string;
  type?: string;
  content?: string;
  description?: string;
  is_active?: boolean;
  [key: string]: unknown;
}

interface ListTemplatesParams {
  type?: string;
  activeOnly?: boolean;
  pageSize?: number;
  pageToken?: string;
}

interface DocgenResult<T = unknown> {
  status: number;
  data: T;
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

function getBaseUrl(): string {
  return (
    process.env.INSURETECH_API_BASE_URL ??
    process.env.NEXT_PUBLIC_INSURETECH_API_BASE_URL ??
    "http://localhost:8080"
  );
}

function extractCsrf(cookieHeader: string): string {
  const m = cookieHeader.match(/(?:^|;\s*)csrf_token=([^;]*)/);
  return m ? decodeURIComponent(m[1]) : "";
}

// ─── Factory ──────────────────────────────────────────────────────────────────

export function makeDocgenClient(
  request: Request,
  sessionOverrides?: Record<string, string>
) {
  const cookieHeader = request.headers.get("cookie") ?? "";
  const csrf = extractCsrf(cookieHeader);

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };
  if (cookieHeader) headers["cookie"] = cookieHeader;
  if (csrf) headers["X-CSRF-Token"] = csrf;

  // Forward portal / business context headers
  if (sessionOverrides) {
    for (const [k, v] of Object.entries(sessionOverrides)) {
      if (v) headers[k] = v;
    }
  } else {
    for (const h of ["x-portal", "x-business-id", "x-user-id", "x-tenant-id"]) {
      const v = request.headers.get(h);
      if (v) headers[h] = v;
    }
  }

  const base = getBaseUrl();

  async function gw<T>(
    method: string,
    path: string,
    body?: unknown
  ): Promise<DocgenResult<T>> {
    const res = await fetch(`${base}${path}`, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    });
    const json = res.status === 204 ? null : await res.json();
    return { status: res.status, data: json as T };
  }

  return {
    listTemplates(params?: ListTemplatesParams) {
      const q = new URLSearchParams();
      if (params?.type) q.set("type", params.type);
      if (params?.activeOnly) q.set("active_only", "true");
      if (params?.pageSize) q.set("page_size", String(params.pageSize));
      if (params?.pageToken) q.set("page_token", params.pageToken);
      const qs = q.toString();
      return gw("GET", `/v1/document-templates${qs ? `?${qs}` : ""}`);
    },

    getTemplate(templateId: string) {
      return gw("GET", `/v1/document-templates/${encodeURIComponent(templateId)}`);
    },

    createTemplate(payload: CreateDocumentTemplatePayload) {
      return gw("POST", "/v1/document-templates", payload);
    },

    updateTemplate(templateId: string, payload: UpdateDocumentTemplatePayload) {
      return gw(
        "PATCH",
        `/v1/document-templates/${encodeURIComponent(templateId)}`,
        payload
      );
    },

    deleteTemplate(templateId: string) {
      return gw("DELETE", `/v1/document-templates/${encodeURIComponent(templateId)}`);
    },
  };
}
