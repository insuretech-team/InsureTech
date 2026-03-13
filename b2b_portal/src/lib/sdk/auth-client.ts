/**
 * auth-client.ts
 * ──────────────
 * Browser-side client for /api/auth/* BFF routes.
 * Components call these — never the gateway directly.
 *
 * All auth operations go through Next.js API route handlers which:
 *   1. Validate the session cookie server-side
 *   2. Call the SDK (gateway) with proper session/CSRF headers
 *   3. Return safe, typed responses to the browser
 */
import type { PortalAuthResponse, PortalLoginRequest } from "@lib/types/auth";
import { parseJson } from "./shared";

// ─── Types ────────────────────────────────────────────────────────────────────

export type AuthOkResponse = { ok: boolean; message?: string };
export type ProfileResponse = { ok: boolean; message?: string; profile?: Record<string, unknown> };
export type SessionsResponse = { ok: boolean; message?: string; sessions?: Record<string, unknown> };
export type TotpResponse = { ok: boolean; message?: string; totp?: Record<string, unknown> };
export type OtpResponse = { ok: boolean; message?: string; data?: Record<string, unknown> };
export type ProfilePhotoUrlResponse = { ok: boolean; message?: string; uploadUrl?: Record<string, unknown> };

// ─── Client ───────────────────────────────────────────────────────────────────

export const authClient = {
  // ── Session ────────────────────────────────────────────────────────────────
  async login(payload: PortalLoginRequest): Promise<PortalAuthResponse> {
    const response = await fetch("/api/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    return parseJson<PortalAuthResponse>(response);
  },

  async logout(): Promise<AuthOkResponse> {
    const response = await fetch("/api/auth/logout", { method: "POST", keepalive: true });
    return parseJson<AuthOkResponse>(response);
  },

  async getSession(): Promise<PortalAuthResponse> {
    const response = await fetch("/api/auth/session", { method: "GET", cache: "no-store" });
    return parseJson<PortalAuthResponse>(response);
  },

  async refreshToken(): Promise<AuthOkResponse> {
    const response = await fetch("/api/auth/refresh", { method: "POST" });
    return parseJson<AuthOkResponse>(response);
  },

  // ── Profile ────────────────────────────────────────────────────────────────
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

  async getProfilePhotoUploadUrl(): Promise<ProfilePhotoUrlResponse> {
    const response = await fetch("/api/auth/profile-photo-url", { method: "GET", cache: "no-store" });
    return parseJson<ProfilePhotoUrlResponse>(response);
  },

  // ── Password ───────────────────────────────────────────────────────────────
  async changePassword(payload: { old_password: string; new_password: string }): Promise<AuthOkResponse> {
    const response = await fetch("/api/auth/change-password", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    return parseJson<AuthOkResponse>(response);
  },

  // ── Sessions management ────────────────────────────────────────────────────
  async listSessions(): Promise<SessionsResponse> {
    const response = await fetch("/api/auth/sessions", { method: "GET", cache: "no-store" });
    return parseJson<SessionsResponse>(response);
  },

  async revokeSession(sessionId: string): Promise<AuthOkResponse> {
    const response = await fetch(`/api/auth/sessions/${sessionId}`, { method: "DELETE" });
    return parseJson<AuthOkResponse>(response);
  },

  async revokeAllSessions(): Promise<AuthOkResponse> {
    const response = await fetch("/api/auth/sessions", { method: "DELETE" });
    return parseJson<AuthOkResponse>(response);
  },

  // ── TOTP / 2FA ─────────────────────────────────────────────────────────────
  async enableTotp(): Promise<TotpResponse> {
    const response = await fetch("/api/auth/totp", { method: "POST" });
    return parseJson<TotpResponse>(response);
  },

  async disableTotp(totpCode: string): Promise<AuthOkResponse> {
    const response = await fetch("/api/auth/totp", {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ totp_code: totpCode }),
    });
    return parseJson<AuthOkResponse>(response);
  },

  // ── OTP ────────────────────────────────────────────────────────────────────
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

  async sendEmailOtp(purpose?: string): Promise<AuthOkResponse> {
    const response = await fetch("/api/auth/send-email-otp", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ purpose }),
    });
    return parseJson<AuthOkResponse>(response);
  },

  async verifyEmail(payload: { token?: string; otp?: string }): Promise<AuthOkResponse> {
    const response = await fetch("/api/auth/verify-email", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    return parseJson<AuthOkResponse>(response);
  },
};
