import type { PortalAuthResponse, PortalLoginRequest } from "@lib/types/auth";

async function parseJson<T>(response: Response): Promise<T> {
  const contentType = response.headers.get("content-type") ?? "";
  if (!contentType.includes("application/json")) {
    throw new Error("Unexpected response type");
  }
  return (await response.json()) as T;
}

export const authClient = {
  async login(payload: PortalLoginRequest): Promise<PortalAuthResponse> {
    const response = await fetch("/api/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    return parseJson<PortalAuthResponse>(response);
  },

  async logout(): Promise<PortalAuthResponse> {
    const response = await fetch("/api/auth/logout", { method: "POST" });
    return parseJson<PortalAuthResponse>(response);
  },

  async getSession(): Promise<PortalAuthResponse> {
    const response = await fetch("/api/auth/session", { method: "GET", cache: "no-store" });
    return parseJson<PortalAuthResponse>(response);
  },
};
