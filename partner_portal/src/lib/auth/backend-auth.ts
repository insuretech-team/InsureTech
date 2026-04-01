import { create } from "@bufbuild/protobuf";
import {
  authServiceGetCurrentSession,
  authServiceLogin,
  authServiceLogout,
  createInsureTechClient,
  type CurrentSessionRetrievalResponse,
  type LoginResponse,
  type User as SdkUser,
} from "@lifeplus/insuretech-sdk";

import {
  SessionSchema,
  SessionType,
  UserSchema,
  UserType,
  type Session,
  type User,
} from "@lib/proto";
import type { PortalPrincipal, PortalSession } from "@lib/types/auth";

const DEFAULT_API_BASE_URL = "http://localhost:8080";
const DEFAULT_SESSION_TTL_MS = 1000 * 60 * 60 * 12;
const DEFAULT_ROLE: PortalPrincipal["role"] = "PARTNER_ADMIN";

function getApiBaseUrl(): string {
  const url =
    process.env.INSURETECH_API_BASE_URL?.trim() ||
    process.env.NEXT_PUBLIC_INSURETECH_API_BASE_URL?.trim() ||
    DEFAULT_API_BASE_URL;

  // Guard: SDK will silently call wrong endpoint if URL is garbage
  try {
    new URL(url);
    return url;
  } catch {
    console.warn(
      `[backend-auth] INSURETECH_API_BASE_URL="${url}" is not a valid URL. ` +
      `Falling back to ${DEFAULT_API_BASE_URL}. ` +
      `Check your .env.local has INSURETECH_API_BASE_URL=http://localhost:8080`
    );
    return DEFAULT_API_BASE_URL;
  }
}

/**
 * Returns a per-request SDK client instance bound to the correct backend base URL.
 */
function buildRequestOptions(
  cookieHeader?: string,
  extraHeaders?: Record<string, string | undefined>
) {
  const headers: Record<string, string> = {};
  if (cookieHeader) {
    headers.cookie = cookieHeader;
  }
  if (extraHeaders) {
    for (const [key, value] of Object.entries(extraHeaders)) {
      if (value && value.trim()) {
        headers[key] = value;
      }
    }
  }
  const client = createInsureTechClient({ baseUrl: getApiBaseUrl(), apiKey: "" });
  return {
    client,
    headers,
  };
}

function inferDisplayName(email: string | undefined, fallback = "Partner User") {
  if (!email) {
    return fallback;
  }
  const value = email.split("@")[0]?.trim();
  return value ? value.replace(/[._-]+/g, " ") : fallback;
}

async function resolvePartnerContext(cookieHeader?: string): Promise<{ id: string; name: string }> {
  if (!cookieHeader?.trim()) {
    return { id: "", name: "" };
  }

  try {
    const { makeDirectHttp } = await import("@lib/sdk/partner-sdk-client");
    const fakeReq = new Request("http://localhost", {
      headers: { cookie: cookieHeader },
    });
    const result = await makeDirectHttp(fakeReq).get("/v1/partners/me");
    if (!result.ok) {
      return { id: "", name: "" };
    }
    const payload = result.data as Record<string, unknown>;
    return {
      id: typeof payload.partner_id === "string" ? payload.partner_id : "",
      name: typeof payload.partner_name === "string" ? payload.partner_name :
            typeof payload.name === "string" ? payload.name : "",
    };
  } catch {
    return { id: "", name: "" };
  }
}

function toPortalUser(user: SdkUser | undefined, fallbackUserId: string | undefined): User {
  return create(UserSchema, {
    userId: user?.user_id ?? fallbackUserId ?? "",
    email: user?.email ?? "",
    mobileNumber: user?.mobile_number ?? "",
  });
}

function toPortalSessionEntity(
  input:
    | CurrentSessionRetrievalResponse["session"]
    | {
      session_id?: string;
      user_id?: string;
      expires_at?: string;
    }
    | undefined
): Session {
  return create(SessionSchema, {
    sessionId: input?.session_id ?? "",
    userId: input?.user_id ?? "",
    sessionType: SessionType.SERVER_SIDE,
  });
}

function parseUserType(rawType: unknown): UserType {
  if (rawType === UserType.SYSTEM_USER || rawType === 4 || rawType === "USER_TYPE_SYSTEM_USER" || rawType === "SYSTEM_USER") {
    return UserType.SYSTEM_USER;
  }
  if (rawType === UserType.PARTNER || rawType === 5 || rawType === "USER_TYPE_PARTNER" || rawType === "PARTNER") {
    return UserType.PARTNER;
  }
  return UserType.UNSPECIFIED;
}

function mapUserTypeToRole(userType: UserType | undefined): PortalPrincipal["role"] {
  if (!userType) return DEFAULT_ROLE;

  if (userType === UserType.SYSTEM_USER) {
    return "SYSTEM_ADMIN";
  } else if (userType === UserType.PARTNER) {
    // Default to PARTNER_ADMIN for partner users
    // The actual role will be resolved from the partner context or user metadata
    return "PARTNER_ADMIN";
  }
  return DEFAULT_ROLE;
}

export async function loginWithMobile(input: {
  mobileNumber: string;
  password: string;
  deviceId?: string;
}) {
  return authServiceLogin({
    ...buildRequestOptions(),
    body: {
      mobile_number: input.mobileNumber,
      password: input.password,
      device_id: input.deviceId ?? "partner-portal-web",
      device_type: "WEB",
      device_name: "Partner Portal Web",
    },
  });
}

export async function getCurrentSession(cookieHeader: string) {
  return authServiceGetCurrentSession({
    ...buildRequestOptions(cookieHeader),
  });
}

export async function logoutCurrentSession(
  cookieHeader: string,
  csrfToken?: string,
  sessionId?: string
) {
  return authServiceLogout({
    ...buildRequestOptions(cookieHeader, { "X-CSRF-Token": csrfToken }),
    body: {
      session_id: sessionId ?? "",
      logout_reason: "user_initiated",
    },
  });
}

export function getSetCookieHeaders(headers: Headers): string[] {
  const value = headers as Headers & { getSetCookie?: () => string[] };
  if (typeof value.getSetCookie === "function") {
    const cookies = value.getSetCookie();
    if (cookies.length > 0) return cookies;
  }
  const single = headers.get("set-cookie");
  if (single) return [single];
  const xSetCookie = headers.get("x-set-cookie");
  return xSetCookie ? [xSetCookie] : [];
}

export function getErrorMessage(error: unknown, fallback = "Request failed") {
  if (typeof error === "string" && error.trim()) {
    return error;
  }
  if (error && typeof error === "object") {
    const candidates = ["message", "error", "detail", "description"] as const;
    for (const key of candidates) {
      const value = (error as Record<string, unknown>)[key];
      if (typeof value === "string" && value.trim()) {
        return value;
      }
    }
  }
  return fallback;
}

export async function toPortalSessionFromLogin(payload: LoginResponse | Record<string, unknown>, cookieHeader?: string): Promise<PortalSession> {
  const p = payload as unknown as LoginResponse;

  const user = toPortalUser(p.user, p.user_id);
  const session = toPortalSessionEntity({
    session_id: p.session_id,
    user_id: p.user_id ?? p.user?.user_id,
  });

  const rawUserType = p.user?.user_type;
  const userTypeEnum = parseUserType(rawUserType);
  const isSystem = userTypeEnum === UserType.SYSTEM_USER;
  const partnerCtx = isSystem ? { id: "", name: "" } : await resolvePartnerContext(cookieHeader);

  return {
    session,
    principal: {
      partnerId: partnerCtx.id,
      organisationName: partnerCtx.name,
      role: mapUserTypeToRole(userTypeEnum),
      displayName: inferDisplayName(user.email),
      user,
    },
    user,
    expiresAt: Date.now() + DEFAULT_SESSION_TTL_MS,
  };
}

export async function toPortalSessionFromCurrentSession(
  data: CurrentSessionRetrievalResponse | Record<string, unknown>,
  cookieHeader: string
): Promise<PortalSession | null> {
  const d = data as unknown as CurrentSessionRetrievalResponse;
  const currentSession = d.session;
  if (!currentSession) {
    return null;
  }

  const sessionUserId = currentSession.user_id;

  const user = toPortalUser(undefined, sessionUserId);
  const session = toPortalSessionEntity(currentSession);
  const parsedExpiry = currentSession.expires_at ? Date.parse(currentSession.expires_at) : Number.NaN;
  const expiresAt = Number.isNaN(parsedExpiry) ? Date.now() + DEFAULT_SESSION_TTL_MS : parsedExpiry;

  const rawUserType = d.user_type;
  const userTypeEnum = parseUserType(rawUserType);
  const role = mapUserTypeToRole(userTypeEnum);
  const isSystem = userTypeEnum === UserType.SYSTEM_USER;
  const partnerCtx = isSystem ? { id: "", name: "" } : await resolvePartnerContext(cookieHeader);

  return {
    session,
    principal: {
      partnerId: partnerCtx.id,
      organisationName: partnerCtx.name,
      role,
      displayName: inferDisplayName(user.email, "Partner User"),
      user,
    },
    user,
    expiresAt,
  };
}
