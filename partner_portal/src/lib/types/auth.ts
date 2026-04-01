import type { Session, User } from "@lib/proto";

export interface PortalPrincipal {
  partnerId: string;
  /** Human-readable partner organisation name. Empty for SYSTEM_ADMIN. */
  organisationName: string;
  role: "PARTNER_ADMIN" | "PARTNER_FOCAL_PERSON" | "PARTNER_AGENT" | "VIEWER" | "SYSTEM_ADMIN";
  displayName: string;
  user: User;
}

export interface PortalSession {
  session: Session;
  principal: PortalPrincipal;
  user?: User;
  expiresAt: number;
}

export interface PortalLoginRequest {
  mobileNumber?: string;
  password: string;
  deviceId?: string;
}

// Backward-compatible alias used by existing route handlers.
export type LoginRequest = PortalLoginRequest;

export interface PortalAuthResponse {
  ok: boolean;
  message?: string;
  session?: PortalSession;
}

// Backward-compatible alias used by existing client code.
export type AuthResponse = PortalAuthResponse;
