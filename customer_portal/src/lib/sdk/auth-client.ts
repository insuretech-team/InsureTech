/**
 * auth-client.ts
 * ──────────────
 * Browser-side client for /api/auth/* BFF routes (customer portal).
 * Components call these — never the gateway directly.
 */
import { parseJson } from "./shared";

export type AuthOkResponse = { ok: boolean; message?: string };
export type ProfileResponse = { ok: boolean; message?: string; profile?: Record<string, unknown> };
export type SessionsResponse = { ok: boolean; message?: string; sessions?: Record<string, unknown> };
export type OtpResponse = { ok: boolean; message?: string; data?: Record<string, unknown> };

export const authClient = {
  async login(payload: { phone_number?: string; email?: string; password?: string; otp?: string }): Promise<{ ok: boolean; message?: string; data?: Record<string, unknown> }> {
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

  async getSession(): Promise<{ ok: boolean; message?: string; data?: Record<string, unknown> }> {
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

  async updateProfile(payload: Record<string, unknown>): Promise<ProfileResponse> {
    const response = await fetch("/api/auth/profile", {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    return parseJson<ProfileResponse>(response);
  },

  async changePassword(payload: { old_password: string; new_password: string }): Promise<AuthOkResponse> {
    const response = await fetch("/api/auth/change-password", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    return parseJson<AuthOkResponse>(response);
  },

  async sendOtp(purpose?: string): Promise<AuthOkResponse> {
    const response = await fetch("/api/auth/send-otp", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ purpose }),
    });
    return parseJson<AuthOkResponse>(response);
  },

  async verifyOtp(otp: string, purpose?: string): Promise<OtpResponse> {
    const response = await fetch("/api/auth/verify-otp", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ otp, purpose }),
    });
    return parseJson<OtpResponse>(response);
  },
};
