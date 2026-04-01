import type { ApiEnvelope, GenerateDocumentPayload, GenerateDocumentResult, InsurerConfigForm, LibraryResponse, LiveDocument, LiveDocumentTemplate, PortalOverview } from "@/lib/types";

type JsonRecord = Record<string, unknown>;

async function parseJson<T>(response: Response): Promise<T> {
  const data = (await response.json().catch(() => ({}))) as T;
  return data;
}

async function request<T>(input: string, init?: RequestInit): Promise<T> {
  const response = await fetch(input, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers ?? {}),
    },
    cache: "no-store",
  });

  return parseJson<T>(response);
}

function withSearch(path: string, params: Record<string, string | number | undefined>) {
  const url = new URL(path, typeof window === "undefined" ? "http://localhost" : window.location.origin);
  Object.entries(params).forEach(([key, value]) => {
    if (value === undefined || value === "") return;
    url.searchParams.set(key, String(value));
  });
  return `${url.pathname}${url.search}`;
}

export const api = {
  auth: {
    login(payload: { mobileNumber: string; password: string }) {
      return request<ApiEnvelope<{ session: JsonRecord }>>("/api/auth/login", {
        method: "POST",
        body: JSON.stringify(payload),
      });
    },
    getSession() {
      return request<ApiEnvelope<{ session: JsonRecord }>>("/api/auth/session");
    },
    logout() {
      return request<ApiEnvelope<null>>("/api/auth/logout", { method: "POST" });
    },
  },
  insurer: {
    getContext(insurerId?: string) {
      return request<ApiEnvelope<PortalOverview["currentInsurer"] | JsonRecord>>(
        withSearch("/api/insurer/context", { insurerId }),
      );
    },
    getOverview(insurerId?: string) {
      return request<ApiEnvelope<PortalOverview>>(withSearch("/api/overview", { insurerId }));
    },
    updateConfig(payload: InsurerConfigForm) {
      return request<ApiEnvelope<{ saved: boolean }>>("/api/insurer/config", {
        method: "PATCH",
        body: JSON.stringify(payload),
      });
    },
  },
  proposals: {
    list(insurerId?: string, status?: string) {
      return request<ApiEnvelope<PortalOverview["proposals"]>>(
        withSearch("/api/insurance-proposals", { insurerId, status }),
      );
    },
    updateStatus(
      proposalId: string,
      payload: { action: "approve" | "reject"; reason?: string },
    ) {
      return request<ApiEnvelope<{ updated: boolean }>>(`/api/insurance-proposals/${proposalId}/status`, {
        method: "POST",
        body: JSON.stringify(payload),
      });
    },
  },
  claims: {
    list(insurerId?: string, status?: string) {
      return request<ApiEnvelope<PortalOverview["claims"]>>(
        withSearch("/api/claims", { insurerId, status }),
      );
    },
    updateStatus(
      claimId: string,
      payload: {
        action: "approve" | "reject" | "settle";
        amount?: number;
        reason?: string;
        paymentReference?: string;
      },
    ) {
      return request<ApiEnvelope<{ updated: boolean }>>(`/api/claims/${claimId}/status`, {
        method: "POST",
        body: JSON.stringify(payload),
      });
    },
  },
  documents: {
    list(params?: { insurerId?: string; entityType?: string; entityId?: string }) {
      return request<ApiEnvelope<LiveDocument[]>>(
        withSearch("/api/documents", { insurerId: params?.insurerId, entityType: params?.entityType, entityId: params?.entityId }),
      );
    },
    generate(payload: GenerateDocumentPayload, insurerId?: string) {
      return request<ApiEnvelope<GenerateDocumentResult>>(
        withSearch("/api/documents/generate", { insurerId }),
        { method: "POST", body: JSON.stringify(payload) },
      );
    },
    downloadUrl(documentId: string, insurerId?: string) {
      return withSearch(`/api/documents/${encodeURIComponent(documentId)}/download`, { insurerId });
    },
    listTemplates(params?: { insurerId?: string; type?: string }) {
      return request<ApiEnvelope<LiveDocumentTemplate[]>>(
        withSearch("/api/documents/templates", { insurerId: params?.insurerId, type: params?.type }),
      );
    },
    generateFromCard(payload: {
      documentId: string;
      category: string;
      title?: string;
      kind?: string;
      data: Record<string, unknown>;
      filename?: string;
    }) {
      return request<ApiEnvelope<{ filename: string; downloadUrl: string; message: string }>>(
        "/api/documents/generate-from-card",
        { method: "POST", body: JSON.stringify(payload) },
      );
    },
    fileUrl(filename: string) {
      return `/api/documents/file/${encodeURIComponent(filename)}`;
    },
    library() {
      return request<ApiEnvelope<LibraryResponse>>("/api/documents/library");
    },
  },
};
