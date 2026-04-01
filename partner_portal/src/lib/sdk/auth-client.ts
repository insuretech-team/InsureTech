/**
 * auth-client.ts — Partner Portal
 * ────────────────────────────────
 * Browser-side client for /api/auth/* BFF routes.
 */
import { parseJson } from "./shared";

export type AuthOkResponse = { ok: boolean; message?: string };
export type ProfileResponse = { ok: boolean; message?: string; profile?: Record<string, unknown> };
export type OtpResponse = { ok: boolean; message?: string; data?: Record<string, unknown> };

export const authClient = {
  async login(payload: { phone_number?: string; email?: string; password?: string }): Promise<{ ok: boolean; message?: string; data?: Record<string, unknown> }> {
    const response = await fetch("/api/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    return parseJson(response);
  },

  async logout(): Promise<AuthOkResponse> {
    const response = await fetch("/api/auth/logout", { method: "POST", keepalive: true });
    return parseJson<AuthOkResponse>(response);
  },

  async getSession(): Promise<{ ok: boolean; data?: Record<string, unknown> }> {
    const response = await fetch("/api/auth/session", { method: "GET", cache: "no-store" });
    return parseJson(response);
  },

  async refreshToken(): Promise<AuthOkResponse> {
    const response = await fetch("/api/auth/refresh", { method: "POST" });
    return parseJson<AuthOkResponse>(response);
  },

  async getProfile(): Promise<ProfileResponse> {
    const response = await fetch("/api/auth/profile", { method: "GET", cache: "no-store" });
    return parseJson<ProfileResponse>(response);
  },

  async sendOtp(payload: { phone_number: string }): Promise<AuthOkResponse> {
    const response = await fetch("/api/auth/send-otp", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    return parseJson<AuthOkResponse>(response);
  },

  async verifyOtp(otp: string): Promise<OtpResponse> {
    const response = await fetch("/api/auth/verify-otp", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ otp }),
    });
    return parseJson<OtpResponse>(response);
  },
};
