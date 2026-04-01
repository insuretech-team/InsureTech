import { create } from "@bufbuild/protobuf";
import {
  authServiceGetCurrentSession,
  authServiceGetUserProfile,
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

function extractCookieValue(cookieHeader: string | undefined, name: string): string {
  if (!cookieHeader?.trim()) {
    return "";
  }
  const escaped = name.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  const match = cookieHeader.match(new RegExp(`(?:^|;\\s*)${escaped}=([^;]*)`));
  return match ? decodeURIComponent(match[1]) : "";
}

async function resolveDisplayName(
  cookieHeader: string | undefined,
  userId: string | undefined,
  userType: UserType | undefined,
  fallbackEmail: string | undefined,
  fallback = "Business User"
): Promise<string> {
  const inferredName = inferDisplayName(fallbackEmail, fallback);
  if (!cookieHeader?.trim() || !userId?.trim()) {
    return inferredName;
  }

  try {
    const result = await authServiceGetUserProfile({
      ...buildRequestOptions(cookieHeader),
      path: { user_id: userId },
    });
    if (!result.response.ok) {
      return inferredName;
    }

    const responseData = (result.data ?? {}) as Record<string, unknown>;
    const rawProfile = (responseData.profile ?? responseData) as Record<string, unknown>;
    const fullName =
      typeof rawProfile.full_name === "string" ? rawProfile.full_name.trim() : "";

    if (fullName) {
      return fullName;
    }
  } catch {
    // fall through to other sources
  }

  if (userType === UserType.B2B_BENEFICIARY) {
    try {
      const { makeDirectHttp } = await import("@lib/sdk/b2b-sdk-client");
      const fakeReq = new Request("http://localhost", {
        headers: { cookie: cookieHeader },
      });
      const profileResult = await makeDirectHttp(fakeReq).get("/v1/b2b-self/profile");
      if (profileResult.ok) {
        const profilePayload = profileResult.data as Record<string, unknown>;
        const employeeView = (profilePayload.employee ?? profilePayload) as Record<string, unknown>;
        const employee = (employeeView.employee ?? employeeView) as Record<string, unknown>;
        const employeeName =
          typeof employee.name === "string" ? employee.name.trim() : "";
        if (employeeName) {
          return employeeName;
        }
      }
    } catch {
      // fall back to inferred name below
    }
  }

  try {
    const cookieEmail = extractCookieValue(cookieHeader, "portal_email");
    return inferDisplayName(cookieEmail || fallbackEmail, fallback);
  } catch {
    return inferredName;
  }
}

async function resolveBusinessContext(
  cookieHeader: string | undefined,
  userType: UserType
): Promise<{ id: string; name: string; role: string }> {
  if (!cookieHeader?.trim()) {
    return { id: "", name: "", role: "" };
  }

  try {
    // Use makeDirectHttp so there are no hardcoded fetch() calls in this file.
    // makeDirectHttp unwraps the unified { success, data, error, meta } envelope
    // and returns { ok, status, data } where data is already the inner payload.
    const { makeDirectHttp } = await import("@lib/sdk/b2b-sdk-client");
    // Build a minimal Request carrying only the session cookie.
    // Use a placeholder URL for the Request constructor — makeDirectHttp only
    // reads headers from it (cookie, x-portal, x-business-id etc.), not the URL.
    // We intentionally omit x-portal/x-business-id here: at login time we don't
    // know them yet, and the backend resolves /organisations/me from the session.
    const fakeReq = new Request("http://localhost", {
      headers: { cookie: cookieHeader },
    });
    if (userType === UserType.B2B_BENEFICIARY) {
      const coverageResult = await makeDirectHttp(fakeReq).get("/v1/b2b-self/coverage");
      if (coverageResult.ok) {
        const coveragePayload = coverageResult.data as Record<string, unknown>;
        const coverage = (coveragePayload.coverage ?? coveragePayload) as Record<string, unknown>;
        return {
          id: typeof coverage.organisation_id === "string" ? coverage.organisation_id : "",
          name: typeof coverage.organisation_name === "string" ? coverage.organisation_name : "",
          role: "B2B_BENEFICIARY",
        };
      }

      const profileResult = await makeDirectHttp(fakeReq).get("/v1/b2b-self/profile");
      if (profileResult.ok) {
        const profilePayload = profileResult.data as Record<string, unknown>;
        const employeeView = (profilePayload.employee ?? profilePayload) as Record<string, unknown>;
        const employee = (employeeView.employee ?? employeeView) as Record<string, unknown>;
        return {
          id: typeof employee.business_id === "string" ? employee.business_id : "",
          name: "",
          role: "B2B_BENEFICIARY",
        };
      }

      return { id: "", name: "", role: "B2B_BENEFICIARY" };
    }
    const result = await makeDirectHttp(fakeReq).get("/v1/b2b/organisations/me");
    if (!result.ok) {
      return { id: "", name: "", role: "" };
    }
    // result.data is the unwrapped inner payload from the unified envelope.
    const payload = result.data as Record<string, unknown>;
    return {
      id: typeof payload.organisation_id === "string" ? payload.organisation_id : "",
      name: typeof payload.organisation_name === "string" ? payload.organisation_name :
            typeof payload.name === "string" ? payload.name : "",
      role: typeof payload.role === "string" ? payload.role : "",
    };
  } catch {
    return { id: "", name: "", role: "" };
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
  if (rawType === UserType.B2B_BENEFICIARY || rawType === 9 || rawType === "USER_TYPE_B2B_BENEFICIARY" || rawType === "B2B_BENEFICIARY") {
    return UserType.B2B_BENEFICIARY;
  }
  return UserType.UNSPECIFIED;
}

function mapUserTypeToRole(userType: UserType | undefined): PortalPrincipal["role"] {
  if (!userType) return DEFAULT_ROLE;

  if (userType === UserType.SYSTEM_USER) {
    return "SYSTEM_ADMIN";
  } else if (userType === UserType.B2B_ORG_ADMIN) {
    return "B2B_ORG_ADMIN";
  } else if (userType === UserType.B2B_BENEFICIARY) {
    return "B2B_BENEFICIARY";
  }
  return DEFAULT_ROLE;
}

function mapOrgMemberRoleToPortalRole(
  rawRole: string | undefined,
  fallback: PortalPrincipal["role"]
): PortalPrincipal["role"] {
  const role = (rawRole ?? "").trim().toUpperCase();
  if (role === "ORG_MEMBER_ROLE_BUSINESS_ADMIN" || role === "ORG_MEMBER_ROLE_ADMIN") {
    return "B2B_ORG_ADMIN";
  }
  if (role === "ORG_MEMBER_ROLE_HR_MANAGER" || role === "ORG_MEMBER_ROLE_HR_STAFF") {
    return "HR_MANAGER";
  }
  return fallback;
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
  // Primary: use the standard getSetCookie() API (Node.js 18+)
  const value = headers as Headers & { getSetCookie?: () => string[] };
  if (typeof value.getSetCookie === "function") {
    const cookies = value.getSetCookie();
    if (cookies.length > 0) return cookies;
  }
  // Fallback 1: raw set-cookie header (may be present in some environments)
  const single = headers.get("set-cookie");
  if (single) return [single];
  // Fallback 2: x-set-cookie — the SDK interceptor copies Set-Cookie here because
  // the Fetch API forbids Set-Cookie in new Response() constructor headers, causing
  // it to be silently dropped. The interceptor preserves it as x-set-cookie so
  // server-side Next.js route handlers can still forward the session cookie.
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
  // The SDK interceptor (client-wrapper.ts) already unwraps the ApiResponse<T> envelope,
  // so payload is already the inner LoginResponse — NOT { success, data, error, meta }.
  // Do NOT double-unwrap here.
  const p = payload as unknown as LoginResponse;

  const user = toPortalUser(p.user, p.user_id);
  const session = toPortalSessionEntity({
    session_id: p.session_id,
    user_id: p.user_id ?? p.user?.user_id,
  });

  const rawUserType = p.user?.user_type;
  const userTypeEnum = parseUserType(rawUserType);
  const isSystem = userTypeEnum === UserType.SYSTEM_USER;
  const bizCtx = isSystem ? { id: "", name: "", role: "" } : await resolveBusinessContext(cookieHeader, userTypeEnum);
  const displayName = await resolveDisplayName(cookieHeader, user.userId, userTypeEnum, user.email);
  const role = isSystem
    ? mapUserTypeToRole(userTypeEnum)
    : userTypeEnum === UserType.B2B_BENEFICIARY
      ? "B2B_BENEFICIARY"
      : mapOrgMemberRoleToPortalRole(bizCtx.role, mapUserTypeToRole(userTypeEnum));

  return {
    session,
    principal: {
      businessId: bizCtx.id,
      organisationName: bizCtx.name,
      role,
      displayName,
      user,
    },
    user,
    passwordChangeRequired: Boolean((payload as Record<string, unknown>)?.password_change_required),
    expiresAt: Date.now() + DEFAULT_SESSION_TTL_MS,
  };
}

export async function toPortalSessionFromCurrentSession(
  data: CurrentSessionRetrievalResponse | Record<string, unknown>,
  cookieHeader: string
): Promise<PortalSession | null> {
  // The SDK interceptor (client-wrapper.ts) already unwraps the ApiResponse<T> envelope,
  // so data is already the inner CurrentSessionRetrievalResponse — NOT { success, data, error, meta }.
  // Do NOT double-unwrap here.
  const d = data as unknown as CurrentSessionRetrievalResponse;
  const currentSession = d.session;
  if (!currentSession) {
    return null;
  }

  const sessionUserId = currentSession.user_id;

  const user = create(UserSchema, {
    userId: sessionUserId,
    email: extractCookieValue(cookieHeader, "portal_email"),
    mobileNumber: extractCookieValue(cookieHeader, "portal_mobile"),
  });
  const session = toPortalSessionEntity(currentSession);
  const parsedExpiry = currentSession.expires_at ? Date.parse(currentSession.expires_at) : Number.NaN;
  const expiresAt = Number.isNaN(parsedExpiry) ? Date.now() + DEFAULT_SESSION_TTL_MS : parsedExpiry;

  const rawUserType = d.user_type;
  const userTypeEnum = parseUserType(rawUserType);
  const isSystem = userTypeEnum === UserType.SYSTEM_USER;
  const bizCtx = isSystem ? { id: "", name: "", role: "" } : await resolveBusinessContext(cookieHeader, userTypeEnum);
  const displayName = await resolveDisplayName(cookieHeader, user.userId, userTypeEnum, user.email);
  const role = isSystem
    ? mapUserTypeToRole(userTypeEnum)
    : userTypeEnum === UserType.B2B_BENEFICIARY
      ? "B2B_BENEFICIARY"
      : mapOrgMemberRoleToPortalRole(bizCtx.role, mapUserTypeToRole(userTypeEnum));

  return {
    session,
    principal: {
      businessId: bizCtx.id,
      organisationName: bizCtx.name,
      role,
      displayName,
      user,
    },
    user,
    passwordChangeRequired: Boolean((data as Record<string, unknown>)?.password_change_required),
    expiresAt,
  };
}

