import type { Session, User } from "@lib/proto";

export interface PortalPrincipal {
  businessId: string;
  /** Human-readable organisation name for B2B admin users. Empty for super_admin. */
  organisationName: string;
  role: "BUSINESS_ADMIN" | "FINANCE_MANAGER" | "HR_MANAGER" | "SYSTEM_ADMIN" | "B2B_ORG_ADMIN";
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
