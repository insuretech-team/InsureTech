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
const DEFAULT_ROLE: PortalPrincipal["role"] = "BUSINESS_ADMIN";

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
 * The global SDK client defaults to "https://api.insuretech.com" (the production URL
 * baked into the SDK bundle). Passing `client` in options overrides it so that
 * server-side Next.js routes always call the local gateway (localhost:8080 in dev).
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
  // Create a per-request client bound to the correct local gateway base URL.
  // The SDK's global client is hardcoded to "https://api.insuretech.com" (production).
  // createInsureTechClient requires apiKey for the Authorization header, but for
  // server-side session auth the gateway authenticates via the session cookie, not
  // Bearer token — so we pass an empty apiKey and rely on the cookie header instead.
  const client = createInsureTechClient({ baseUrl: getApiBaseUrl(), apiKey: "" });
  return {
    client,
    headers,
  };
}

function inferDisplayName(email: string | undefined, fallback = "Business User") {
  if (!email) {
    return fallback;
  }
  const value = email.split("@")[0]?.trim();
  return value ? value.replace(/[._-]+/g, " ") : fallback;
}

async function resolveBusinessContext(cookieHeader?: string): Promise<{ id: string; name: string }> {
  if (!cookieHeader?.trim()) {
    return { id: "", name: "" };
  }

  try {
    const res = await fetch(`${getApiBaseUrl()}/v1/b2b/organisations/me`, {
      method: "GET",
      headers: { cookie: cookieHeader },
      cache: "no-store",
    });
    if (!res.ok) {
      return { id: "", name: "" };
    }
    const data = (await res.json()) as Record<string, unknown>;
    return {
      id: typeof data.organisation_id === "string" ? data.organisation_id : "",
      name: typeof data.organisation_name === "string" ? data.organisation_name : "",
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
  if (rawType === UserType.B2B_ORG_ADMIN || rawType === 8 || rawType === "USER_TYPE_B2B_ORG_ADMIN" || rawType === "B2B_ORG_ADMIN") {
    return UserType.B2B_ORG_ADMIN;
  }
  if (rawType === UserType.BUSINESS_ADMIN || rawType === 7 || rawType === "USER_TYPE_BUSINESS_ADMIN" || rawType === "BUSINESS_ADMIN") {
    return UserType.BUSINESS_ADMIN;
  }
  return UserType.UNSPECIFIED;
}

function mapUserTypeToRole(userType: UserType | undefined): PortalPrincipal["role"] {
  if (!userType) return DEFAULT_ROLE;

  if (userType === UserType.SYSTEM_USER) {
    return "SYSTEM_ADMIN";
  } else if (userType === UserType.B2B_ORG_ADMIN) {
    return "B2B_ORG_ADMIN";
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
      device_id: input.deviceId ?? "customer-portal-web",
      device_type: "WEB",
      device_name: "Customer Portal Web",
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
    return value.getSetCookie();
  }
  const single = headers.get("set-cookie");
  return single ? [single] : [];
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

export async function toPortalSessionFromLogin(payload: LoginResponse, cookieHeader?: string): Promise<PortalSession> {
  const user = toPortalUser(payload.user, payload.user_id);
  const session = toPortalSessionEntity({
    session_id: payload.session_id,
    user_id: payload.user_id ?? payload.user?.user_id,
  });

  const rawUserType = payload.user?.user_type;
  const userTypeEnum = parseUserType(rawUserType);
  const isSystem = userTypeEnum === UserType.SYSTEM_USER;
  const bizCtx = isSystem ? { id: "", name: "" } : await resolveBusinessContext(cookieHeader);

  return {
    session,
    principal: {
      businessId: bizCtx.id,
      organisationName: bizCtx.name,
      role: mapUserTypeToRole(userTypeEnum),
      displayName: inferDisplayName(user.email),
      user,
    },
    user,
    expiresAt: Date.now() + DEFAULT_SESSION_TTL_MS,
  };
}

export async function toPortalSessionFromCurrentSession(
  data: CurrentSessionRetrievalResponse,
  cookieHeader: string
): Promise<PortalSession | null> {
  const currentSession = data.session;
  if (!currentSession) {
    return null;
  }

  const sessionUserId = currentSession.user_id;

  const user = toPortalUser(undefined, sessionUserId);
  const session = toPortalSessionEntity(currentSession);
  const parsedExpiry = currentSession.expires_at ? Date.parse(currentSession.expires_at) : Number.NaN;
  const expiresAt = Number.isNaN(parsedExpiry) ? Date.now() + DEFAULT_SESSION_TTL_MS : parsedExpiry;

  const rawUserType = data.user_type;
  const userTypeEnum = parseUserType(rawUserType);
  const role = mapUserTypeToRole(userTypeEnum);
  const isSystem = userTypeEnum === UserType.SYSTEM_USER;
  const bizCtx = isSystem ? { id: "", name: "" } : await resolveBusinessContext(cookieHeader);

  return {
    session,
    principal: {
      businessId: bizCtx.id,
      organisationName: bizCtx.name,
      role,
      displayName: inferDisplayName(user.email, "Business User"),
      user,
    },
    user,
    expiresAt,
  };
}

