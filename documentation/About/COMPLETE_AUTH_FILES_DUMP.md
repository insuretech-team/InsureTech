# Complete Auth System Code Dump

This document contains all authentication-related source files from the b2b_portal project.

## 1. src/lib/auth/session.ts

```typescript
import { cookies } from "next/headers";
import { redirect } from "next/navigation";

import { getSession } from "./session-store";

export const SESSION_COOKIE_NAME = "session_token";

export async function getServerSession() {
  const cookieStore = await cookies();
  return getSession(cookieStore.get(SESSION_COOKIE_NAME)?.value);
}

export async function requireServerSession() {
  const session = await getServerSession();
  if (!session) {
    redirect("/login");
  }
  return session;
}
```

## 2. src/lib/auth/session-store.ts

```typescript
import crypto from "node:crypto";

import { create } from "@bufbuild/protobuf";

import {
  DeviceType,
  MoneySchema,
  SessionSchema,
  SessionType,
  UserSchema,
  UserStatus,
  UserType,
  type User,
} from "@lib/proto";
import type { PortalPrincipal, PortalSession } from "@lib/types/auth";

const SESSION_TTL_MS = 1000 * 60 * 60 * 12;

const sessionStore = new Map<string, PortalSession>();

function toTimestamp(milliseconds: number) {
  return {
    seconds: BigInt(Math.floor(milliseconds / 1000)),
    nanos: (milliseconds % 1000) * 1_000_000,
  };
}

export function createSession(principal: PortalPrincipal): PortalSession {
  // Ensure organisationName is always present (backwards compat for callers that omit it)
  if (principal.organisationName === undefined) {
    (principal as PortalPrincipal).organisationName = "";
  }
  const now = Date.now();
  const sessionId = crypto.randomUUID();
  const expiresAt = now + SESSION_TTL_MS;
  const csrfToken = crypto.randomBytes(16).toString("hex");

  const session = {
    session: create(SessionSchema, {
      sessionId,
      userId: principal.user.userId,
      sessionType: SessionType.SERVER_SIDE,
      sessionTokenLookup: crypto.createHash("sha256").update(sessionId).digest("hex"),
      expiresAt: toTimestamp(expiresAt),
      ipAddress: "127.0.0.1",
      userAgent: "b2b_portal",
      deviceId: "web-browser",
      deviceName: "Web Browser",
      deviceType: DeviceType.WEB,
      createdAt: toTimestamp(now),
      lastActivityAt: toTimestamp(now),
      isActive: true,
      csrfToken,
    }),
    principal,
    user: principal.user,
    expiresAt,
  } satisfies PortalSession;

  sessionStore.set(sessionId, session);
  return session;
}

export function getSession(sessionId: string | undefined): PortalSession | null {
  if (!sessionId) {
    return null;
  }
  const session = sessionStore.get(sessionId);
  if (!session) {
    return null;
  }
  if (session.expiresAt <= Date.now()) {
    sessionStore.delete(sessionId);
    return null;
  }
  return session;
}

export function clearSession(sessionId: string | undefined): void {
  if (!sessionId) {
    return;
  }
  sessionStore.delete(sessionId);
}
```

## 3. src/lib/auth/resolve-user-id.ts

```typescript
/**
 * resolve-user-id.ts
 * ───────────────────
 * Utility to reliably resolve the current user's ID from session.
 *
 * Problem: The login response JSON body may not include user_id — it only sets
 * an HttpOnly session cookie. So portal_user_id cookie can be empty ("").
 * This helper falls back to calling getCurrentSession on the gateway to get it.
 */
import { authServiceGetCurrentSession, createInsureTechClient } from "@lifeplus/insuretech-sdk";
import { getApiBaseUrl } from "@lib/sdk/api-helpers";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import type { resolvePortalHeaders } from "@lib/sdk/session-headers";

type PortalHeaders = Awaited<ReturnType<typeof resolvePortalHeaders>>;

/**
 * Returns the real user_id string, or null if it cannot be determined.
 * 1. Prefers hdrs.userId if non-empty (already resolved from cookie)
 * 2. Falls back to calling getCurrentSession on the gateway
 */
export async function resolveUserIdFromSession(
  request: Request,
  hdrs: NonNullable<PortalHeaders>
): Promise<string | null> {
  // Fast path — already in cookie
  if (hdrs.userId) return hdrs.userId;

  // Fallback — call getCurrentSession via SDK with the request cookies
  try {
    const sdk = makeSdkClient(request, hdrs);
    const sessionResult = await sdk.getCurrentSession();
    if (sessionResult.response.ok && sessionResult.data) {
      const data = sessionResult.data as Record<string, unknown>;
      // Shape: { session: { user_id: string } } or { user_id: string }
      const userId = (((data.session as Record<string, unknown>)?.user_id) || data.user_id) as string | undefined;
      if (userId) return userId;
    }
  } catch { /* ignore */ }

  return null;
}
```

## 4. src/lib/auth/backend-auth.ts

```typescript
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
```

## 5. middleware.ts

```typescript
import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

const SESSION_COOKIE_NAME = "session_token";
const PUBLIC_PATHS = ["/login", "/api/auth/login"];

// Routes that require a specific role to access.
// Note: middleware runs in the edge runtime — we cannot call getSession() here
// (it's in-memory Node.js). We store the role in a separate lightweight cookie
// set at login time: "portal_role". If that cookie is absent but session exists,
// the API route will enforce via session store (defence-in-depth).
// The middleware provides UX-level redirects only.
const ROLE_GUARDS: Array<{ prefix: string; allowedRoles: string[] }> = [
  { prefix: "/organisations", allowedRoles: ["SYSTEM_ADMIN"] },
  { prefix: "/team",          allowedRoles: ["SYSTEM_ADMIN", "B2B_ORG_ADMIN", "BUSINESS_ADMIN"] },
  { prefix: "/departments",   allowedRoles: ["SYSTEM_ADMIN", "B2B_ORG_ADMIN", "BUSINESS_ADMIN", "HR_MANAGER", "VIEWER"] },
  { prefix: "/employees",     allowedRoles: ["SYSTEM_ADMIN", "B2B_ORG_ADMIN", "BUSINESS_ADMIN", "HR_MANAGER", "VIEWER"] },
  { prefix: "/purchase-orders", allowedRoles: ["SYSTEM_ADMIN", "B2B_ORG_ADMIN", "BUSINESS_ADMIN", "HR_MANAGER", "VIEWER"] },
];

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  if (
    pathname.startsWith("/_next") ||
    pathname.startsWith("/public") ||
    pathname.startsWith("/logos") ||
    pathname.startsWith("/navbar-icons") ||
    pathname.startsWith("/stats-cards") ||
    pathname.startsWith("/quotations/") ||
    pathname.startsWith("/insurance-plans") ||
    // All /api/* routes handle their own auth via session cookie forwarding.
    // The middleware must NOT redirect API routes — the SDK client forwards
    // the session cookie to the backend which validates it.
    pathname.startsWith("/api/") ||
    pathname === "/favicon.ico"
  ) {
    return NextResponse.next();
  }

  const isPublic = PUBLIC_PATHS.some((path) => pathname === path || pathname.startsWith(path + "/"));
  const hasSessionCookie = Boolean(request.cookies.get(SESSION_COOKIE_NAME)?.value);

  if (!hasSessionCookie && !isPublic) {
    const loginUrl = new URL("/login", request.url);
    loginUrl.searchParams.set("next", pathname);
    return NextResponse.redirect(loginUrl);
  }

  if (hasSessionCookie && pathname === "/login") {
    // Redirect to appropriate default page based on role
    const role = request.cookies.get("portal_role")?.value ?? "";
    const dest = role === "SYSTEM_ADMIN" ? "/organisations" : "/";
    return NextResponse.redirect(new URL(dest, request.url));
  }

  // Role-based route guard (UX-level, uses portal_role cookie set at login)
  if (hasSessionCookie) {
    const role = request.cookies.get("portal_role")?.value ?? "";
    if (role) {
      const guard = ROLE_GUARDS.find((g) => pathname.startsWith(g.prefix));
      if (guard && !guard.allowedRoles.includes(role)) {
        // Redirect to appropriate default page for their role
        const fallback = role === "SYSTEM_ADMIN" ? "/organisations" : "/";
        return NextResponse.redirect(new URL(fallback, request.url));
      }
    }
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/((?!.*\\..*).*)""],
};
```

## 6. next.config.ts

```typescript
import type { NextConfig } from "next";
import path from "path";

const nextConfig: NextConfig = {
  output: "standalone",

  // Transpile the local-tgz SDK so Next.js processes its ESM output correctly.
  transpilePackages: ["@lifeplus/insuretech-sdk"],

  // Webpack alias (used by `next build` and `next start`)
  webpack(config) {
    config.resolve.alias = {
      ...config.resolve.alias,
      "@lifeplus/insuretech-sdk": path.resolve(
        __dirname,
        "node_modules/@lifeplus/insuretech-sdk"
      ),
    };
    return config;
  },

  // Turbopack alias must use forward-slash relative path (Windows absolute paths
  // are not supported by Turbopack — "windows imports are not implemented yet")
  turbopack: {
    resolveAlias: {
      "@lifeplus/insuretech-sdk": "./node_modules/@lifeplus/insuretech-sdk/dist/index.mjs",
    },
  },
};

export default nextConfig;
```

## 7. app/api/auth/login/route.ts

```typescript
import { NextResponse } from "next/server";

import {
  getSetCookieHeaders,
  loginWithMobile,
  toPortalSessionFromLogin,
} from "@lib/auth/backend-auth";
import { SESSION_COOKIE_NAME } from "@lib/auth/session";
import type { LoginRequest } from "@lib/types/auth";

/**
 * Maps backend errors (gRPC codes, HTTP status codes, raw strings) to clean,
 * user-facing messages. Never leaks internal RPC details to the UI.
 */
function toUserFriendlyLoginError(error: unknown, httpStatus: number): string {
  // Extract the raw error string from whatever shape the SDK returns
  let raw = "";
  if (typeof error === "string") {
    raw = error;
  } else if (error && typeof error === "object") {
    for (const key of ["message", "error", "detail", "description"] as const) {
      const v = (error as Record<string, unknown>)[key];
      if (typeof v === "string" && v.trim()) { raw = v; break; }
    }
  }
  const lower = raw.toLowerCase();

  // ── gRPC / HTTP status → friendly message map ──────────────────────────────

  // Wrong password / user not found
  if (
    httpStatus === 401 ||
    lower.includes("unauthenticated") ||
    lower.includes("invalid password") ||
    lower.includes("invalid credentials") ||
    lower.includes("wrong password") ||
    lower.includes("password") ||
    lower.includes("not found") ||
    lower.includes("no user") ||
    lower.includes("user not found")
  ) {
    return "Mobile number or password is incorrect. Please try again.";
  }

  // Account locked / too many attempts
  if (
    lower.includes("locked") ||
    lower.includes("too many") ||
    lower.includes("rate limit") ||
    lower.includes("max attempt") ||
    lower.includes("blocked") ||
    httpStatus === 429
  ) {
    return "Your account has been temporarily locked due to too many failed attempts. Please try again later.";
  }

  // Account not active / disabled
  if (
    lower.includes("inactive") ||
    lower.includes("disabled") ||
    lower.includes("suspended") ||
    lower.includes("banned") ||
    lower.includes("not active")
  ) {
    return "Your account is not active. Please contact your administrator.";
  }

  // Invalid mobile number (should be caught client-side, but just in case authn rejects it)
  if (
    lower.includes("invalid mobile") ||
    lower.includes("invalid_argument") ||
    lower.includes("invalidargument") ||
    lower.includes("mobile_number") ||
    lower.includes("phone") ||
    httpStatus === 400
  ) {
    return "Invalid mobile number. Please enter a valid Bangladesh number (e.g. 01712345678).";
  }

  // Server / network errors
  if (httpStatus >= 500 || lower.includes("unavailable") || lower.includes("internal")) {
    return "The service is temporarily unavailable. Please try again in a moment.";
  }

  // Generic fallback — never show raw RPC text
  return "Login failed. Please check your mobile number and password and try again.";
}

const CSRF_COOKIE_NAME = "csrf_token";

function extractCookieValue(setCookieHeader: string, cookieName: string): string | undefined {
  const [nameValue] = setCookieHeader.split(";", 1);
  if (!nameValue) {
    return undefined;
  }
  const separatorIndex = nameValue.indexOf("=");
  if (separatorIndex <= 0) {
    return undefined;
  }
  const name = nameValue.slice(0, separatorIndex).trim();
  if (name !== cookieName) {
    return undefined;
  }
  return nameValue.slice(separatorIndex + 1);
}

// Valid Bangladesh operator prefixes: 013,014,015,016,017,018,019
const BD_PHONE_RE = /^880(13|14|15|16|17|18|19)\d{8}$/;

/**
 * Normalizes a Bangladesh mobile number to canonical E.164 form (+880XXXXXXXXXX).
 *
 * Accepted input variants (spaces, dashes, dots freely ignored):
 *   01712345678          → +8801712345678
 *   1712345678           → +8801712345678   (10 digits, no leading 0)
 *   8801712345678        → +8801712345678
 *   00 8801712345678     → +8801712345678
 *   +880 171-234-5678    → +8801712345678
 *   +88 01712345678      → +8801712345678   (typo with 88 instead of 880)
 *
 * Returns null when the number cannot be recognized as a valid BD number.
 */
function normalizeMobileNumber(value: string): string | null {
  // Strip everything except digits and a leading +
  const stripped = value.trim().replace(/[^\d+]/g, "");

  // Drop the leading + so we work purely with digits from here
  const digits = stripped.startsWith("+") ? stripped.slice(1) : stripped;

  let e164Digits: string; // will hold 880XXXXXXXXXX (13 digits)

  if (digits.startsWith("00880")) {
    // 008801712345678 → 8801712345678
    e164Digits = digits.slice(2);
  } else if (digits.startsWith("880")) {
    // 8801712345678
    e164Digits = digits;
  } else if (digits.startsWith("0088")) {
    // 00881712345678 — uncommon but handle gracefully
    e164Digits = "880" + digits.slice(4);
  } else if (digits.startsWith("88") && digits.length === 13) {
    // 88 followed by 01XXXXXXXXX — missing a zero: treat as typo
    e164Digits = "880" + digits.slice(2);
  } else if (digits.startsWith("0")) {
    // 01712345678 (11 digits local)
    e164Digits = "880" + digits.slice(1);
  } else if (digits.length === 10) {
    // 1712345678 — 10 digits without leading 0
    e164Digits = "880" + digits;
  } else {
    return null;
  }

  if (!BD_PHONE_RE.test(e164Digits)) {
    return null;
  }

  return `+${e164Digits}`;
}

export async function POST(request: Request) {
  let payload: LoginRequest;
  try {
    payload = (await request.json()) as LoginRequest;
  } catch {
    return NextResponse.json({ ok: false, message: "Invalid login payload" }, { status: 400 });
  }

  const mobileRaw = payload.mobileNumber?.trim();
  if (!mobileRaw) {
    return NextResponse.json(
      { ok: false, message: "Mobile number is required" },
      { status: 400 }
    );
  }
  if (!payload.password?.trim()) {
    return NextResponse.json(
      { ok: false, message: "Password is required" },
      { status: 400 }
    );
  }

  const normalizedMobile = normalizeMobileNumber(mobileRaw);
  if (!normalizedMobile) {
    return NextResponse.json(
      {
        ok: false,
        message:
          "Invalid mobile number. Please enter a valid Bangladesh number " +
          "(e.g. 01712345678, +8801712345678 or 008801712345678).",
      },
      { status: 400 }
    );
  }

  const result = await loginWithMobile({
    mobileNumber: normalizedMobile,
    password: payload.password,
    deviceId: payload.deviceId,
  });

  if (result.error) {
    const httpStatus = result.response?.status || 500;
    return NextResponse.json(
      { ok: false, message: toUserFriendlyLoginError(result.error, httpStatus) },
      { status: httpStatus }
    );
  }

  const response = NextResponse.json({ ok: true }, { status: result.response.status || 200 });
  const setCookieHeaders = getSetCookieHeaders(result.response.headers);
  const backendSessionCookie = setCookieHeaders.find((value) =>
    value.startsWith(`${SESSION_COOKIE_NAME}=`)
  );
  const sessionToken = backendSessionCookie
    ? extractCookieValue(backendSessionCookie, SESSION_COOKIE_NAME)
    : undefined;
  if (sessionToken) {
    response.cookies.set({
      name: SESSION_COOKIE_NAME,
      value: sessionToken,
      path: "/",
      httpOnly: true,
      sameSite: "strict",
      secure: process.env.NODE_ENV === "production",
      maxAge: 60 * 60 * 12,
    });
  }

  const csrfToken = result.response.headers.get("x-csrf-token") ?? result.data?.csrf_token;
  if (csrfToken) {
    response.cookies.set({
      name: CSRF_COOKIE_NAME,
      value: csrfToken,
      path: "/",
      httpOnly: true,
      sameSite: "lax",
      secure: process.env.NODE_ENV === "production",
      maxAge: 60 * 60 * 12,
    });
  }

  const sessionCookieHeader = sessionToken ? `${SESSION_COOKIE_NAME}=${sessionToken}` : undefined;
  const session = await toPortalSessionFromLogin(result.data ?? {}, sessionCookieHeader);
  const finalResponse = NextResponse.json({ ok: true, session }, { status: response.status });
  for (const cookie of response.cookies.getAll()) {
    finalResponse.cookies.set(cookie);
  }

  // Set lightweight metadata cookies used by:
  //   1. Edge middleware (portal_role) — for role-based page routing without hitting DB
  //   2. API route session-headers helper — to inject x-portal/x-user-id/x-business-id
  //      into backend SDK calls so the Casbin authz interceptor gets the right domain.
  //
  // These are NOT security boundaries — the backend session cookie enforces real authz.
  const portalRole = session.principal.role ?? "BUSINESS_ADMIN";
  const portalBizId = session.principal.businessId ?? "";

  // user_id may not be in the login JSON response body — it lives in the HttpOnly session cookie.
  // Prefer: from toPortalSessionFromLogin, then from result.data directly, then from getCurrentSession.
  let portalUserId = session.principal.user?.userId ?? (result.data as Record<string, unknown>)?.user_id as string ?? "";

  // If user_id is still empty, call getCurrentSession using the new session token to resolve it.
  // This handles gateways that don't return user_id in the login JSON body.
  if (!portalUserId && sessionToken) {
    try {
      const { authServiceGetCurrentSession, createInsureTechClient } = await import("@lifeplus/insuretech-sdk");
      const { getApiBaseUrl } = await import("@lib/sdk/api-helpers");
      const { toPortalSessionFromCurrentSession } = await import("@lib/auth/backend-auth");
      const tempClient = createInsureTechClient({ baseUrl: getApiBaseUrl(), apiKey: process.env.INSURETECH_API_KEY ?? "" });
      const cookieStr = `${SESSION_COOKIE_NAME}=${sessionToken}`;
      const sessionRes = await authServiceGetCurrentSession({
        client: tempClient,
        headers: { Cookie: cookieStr },
        throwOnError: false,
      });
      if (sessionRes.response.ok && sessionRes.data) {
        const portalSession = await toPortalSessionFromCurrentSession(sessionRes.data, cookieStr);
        portalUserId = portalSession?.principal?.user?.userId ?? "";
      }
    } catch { /* ignore — userId will remain empty */ }
  }
  const cookieOpts = {
    path: "/",
    httpOnly: false, // must be readable by edge middleware + session-headers helper
    sameSite: "strict" as const,
    secure: process.env.NODE_ENV === "production",
    maxAge: 60 * 60 * 12,
  };

  finalResponse.cookies.set({ name: "portal_role", value: portalRole, ...cookieOpts });
  finalResponse.cookies.set({ name: "portal_user_id", value: portalUserId, ...cookieOpts });
  finalResponse.cookies.set({ name: "portal_biz_id", value: portalBizId, ...cookieOpts });

  // Store user contact info cookies so the My Profile page can display
  // mobile_number and email — these live on the User record, not UserProfile.
  // They are read-only identity fields (auth credentials), not profile fields.
  const portalMobile = (result.data as Record<string, unknown>)?.user
    ? ((result.data as Record<string, unknown>).user as Record<string, unknown>)?.mobile_number as string ?? ""
    : "";
  const portalEmail = (result.data as Record<string, unknown>)?.user
    ? ((result.data as Record<string, unknown>).user as Record<string, unknown>)?.email as string ?? ""
    : "";
  finalResponse.cookies.set({ name: "portal_mobile", value: portalMobile, ...cookieOpts });
  finalResponse.cookies.set({ name: "portal_email", value: portalEmail, ...cookieOpts });

  return finalResponse;
}
```

## 8. app/api/auth/logout/route.ts

```typescript
import { NextResponse } from "next/server";

import {
  getCurrentSession,
  getErrorMessage,
  getSetCookieHeaders,
  logoutCurrentSession,
} from "@lib/auth/backend-auth";
import { SESSION_COOKIE_NAME } from "@lib/auth/session";

const CSRF_COOKIE_NAME = "csrf_token";

function getCookieValue(cookieHeader: string, cookieName: string): string | undefined {
  const target = `${cookieName}=`;
  for (const rawPart of cookieHeader.split(";")) {
    const part = rawPart.trim();
    if (part.startsWith(target)) {
      return decodeURIComponent(part.slice(target.length));
    }
  }
  return undefined;
}

function expireSessionCookie(response: NextResponse) {
  response.cookies.set({
    name: SESSION_COOKIE_NAME,
    value: "",
    path: "/",
    httpOnly: true,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    expires: new Date(0),
  });
}

function expireCsrfCookie(response: NextResponse) {
  response.cookies.set({
    name: CSRF_COOKIE_NAME,
    value: "",
    path: "/",
    httpOnly: true,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    expires: new Date(0),
  });
}

function expirePortalCookies(response: NextResponse) {
  for (const name of ["portal_role", "portal_user_id", "portal_biz_id"]) {
    response.cookies.set({
      name,
      value: "",
      path: "/",
      httpOnly: false,
      sameSite: "strict",
      secure: process.env.NODE_ENV === "production",
      expires: new Date(0),
    });
  }
}

export async function POST(request: Request) {
  const cookieHeader = request.headers.get("cookie") ?? "";
  const csrfToken = getCookieValue(cookieHeader, CSRF_COOKIE_NAME);
  let sessionId = "";

  try {
    const currentSessionResult = await getCurrentSession(cookieHeader);
    if (!currentSessionResult.error) {
      sessionId = currentSessionResult.data?.session?.session_id ?? "";
    }
  } catch {
    // Ignore session lookup failures; we can still clear local cookies.
  }

  if (!sessionId) {
    const response = NextResponse.json({ ok: true, message: "No active session" }, { status: 200 });
    expireSessionCookie(response);
    expireCsrfCookie(response);
    expirePortalCookies(response);
    return response;
  }

  let result: Awaited<ReturnType<typeof logoutCurrentSession>>;
  try {
    result = await logoutCurrentSession(cookieHeader, csrfToken, sessionId);
  } catch (error) {
    const response = NextResponse.json(
      { ok: false, message: getErrorMessage(error, "Logout failed") },
      { status: 502 }
    );
    expireSessionCookie(response);
    expireCsrfCookie(response);
    return response;
  }

  if (result.error) {
    const status = result.response?.status ?? 500;
    const response = NextResponse.json(
      { ok: false, message: getErrorMessage(result.error, "Logout failed") },
      { status }
    );
    for (const setCookie of getSetCookieHeaders(result.response.headers)) {
      response.headers.append("set-cookie", setCookie);
    }
    expireSessionCookie(response);
    expireCsrfCookie(response);
    return response;
  }

  const response = NextResponse.json({ ok: true }, { status: result.response.status || 200 });
  for (const setCookie of getSetCookieHeaders(result.response.headers)) {
    response.headers.append("set-cookie", setCookie);
  }
  expireSessionCookie(response);
  expireCsrfCookie(response);
  expirePortalCookies(response);
  return response;
}
```

## 9. app/api/auth/refresh/route.ts

```typescript
import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, badRequest } from "@lib/sdk/api-helpers";

/** POST /api/auth/refresh — refresh the current session token */
export async function POST(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  let body: { refresh_token?: string; device_id?: string; ip_address?: string };
  try { body = await request.json(); } catch { body = {}; }
  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.refreshToken({
    body: {
      refresh_token: body.refresh_token,
      device_id: body.device_id ?? "web",
      ip_address: body.ip_address,
    },
  });
  if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
  return NextResponse.json({ ok: true, data: result.data }, { status: 200 });
}
```

## 10. app/api/auth/session/route.ts

```typescript
import { NextResponse } from "next/server";

import {
  getCurrentSession,
  getErrorMessage,
  toPortalSessionFromCurrentSession,
} from "@lib/auth/backend-auth";
import { SESSION_COOKIE_NAME } from "@lib/auth/session";

export async function GET(request: Request) {
  try {
    const cookieHeader = request.headers.get("cookie") ?? "";
    const hasSessionCookie = cookieHeader.includes(`${SESSION_COOKIE_NAME}=`);
    if (!hasSessionCookie) {
      return NextResponse.json({ ok: false, message: "No active session" }, { status: 401 });
    }

    let result: Awaited<ReturnType<typeof getCurrentSession>>;
    try {
      result = await getCurrentSession(cookieHeader);
    } catch (error) {
      return NextResponse.json(
        { ok: false, message: getErrorMessage(error, "Session service unavailable") },
        { status: 502 }
      );
    }

    if (result.error) {
      const status = result.response?.status ?? 401;
      return NextResponse.json(
        { ok: false, message: getErrorMessage(result.error, "No active session") },
        { status }
      );
    }

    const session = await toPortalSessionFromCurrentSession(result.data ?? {}, cookieHeader);
    if (!session) {
      return NextResponse.json({ ok: false, message: "No active session" }, { status: 401 });
    }

    // Re-mint metadata cookies on every session refresh so they stay in sync
    // with the backend session. This is critical: if the lightweight cookies
    // (portal_role, portal_user_id, portal_biz_id) expire or are cleared while
    // the backend session_token is still valid, session-headers.ts would fall
    // back to PORTAL_B2B with no org context, causing 403 for superadmin on
    // org/dept tabs.
    const response = NextResponse.json({ ok: true, session }, { status: 200 });
    const portalRole = session.principal.role ?? "BUSINESS_ADMIN";
    const portalUserId = session.principal.user?.userId ?? "";
    const portalBizId = session.principal.businessId ?? "";
    const cookieOpts = {
      path: "/",
      httpOnly: false, // must be readable by edge middleware + session-headers helper
      sameSite: "strict" as const,
      secure: process.env.NODE_ENV === "production",
      maxAge: 60 * 60 * 12,
    };
    response.cookies.set({ name: "portal_role",    value: portalRole,    ...cookieOpts });
    response.cookies.set({ name: "portal_user_id", value: portalUserId,  ...cookieOpts });
    response.cookies.set({ name: "portal_biz_id",  value: portalBizId,   ...cookieOpts });

    // Re-mint contact info cookies on every session refresh so they stay in sync.
    // These are sourced from existing cookies (set at login) — if empty, preserve
    // whatever was already in the browser (no-op by writing empty string is safe).
    const existingCookieHeader = request.headers.get("cookie") ?? "";
    const extractCk = (name: string) => {
      const m = existingCookieHeader.match(new RegExp(`(?:^|;\\s*)${name}=([^;]*)`));
      return m ? decodeURIComponent(m[1]) : "";
    };
    const portalMobile = extractCk("portal_mobile");
    const portalEmail  = extractCk("portal_email");
    response.cookies.set({ name: "portal_mobile", value: portalMobile, ...cookieOpts });
    response.cookies.set({ name: "portal_email",  value: portalEmail,  ...cookieOpts });

    return response;
  } catch (error) {
    return NextResponse.json(
      { ok: false, message: getErrorMessage(error, "Session endpoint failed") },
      { status: 502 }
    );
  }
}
```

## 11. app/api/auth/profile/route.ts

```typescript
import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, badRequest } from "@lib/sdk/api-helpers";
import { resolveUserIdFromSession } from "@lib/auth/resolve-user-id";

// Reads a cookie value from a raw Cookie header string.
function extractCookieValue(cookieHeader: string, name: string): string {
  const m = cookieHeader.match(new RegExp(`(?:^|;\\s*)${name}=([^;]*)`));
  return m ? decodeURIComponent(m[1]) : "";
}

/** GET /api/auth/profile -- get current user profile */
export async function GET(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  const userId = await resolveUserIdFromSession(request, hdrs);
  if (!userId) return NextResponse.json({ ok: false, message: "Cannot resolve user identity" }, { status: 401 });
  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.getUserProfile({ path: { user_id: userId } });

  // Read user identity fields from lightweight cookies set at login.
  // mobile_number and email live on the User record, not UserProfile — they are
  // auth credentials and are read-only in the My Profile form.
  const cookieHeader = request.headers.get("cookie") ?? "";
  const mobile_number = extractCookieValue(cookieHeader, "portal_mobile");
  const email         = extractCookieValue(cookieHeader, "portal_email");

  // 404 means the user exists but has no profile row yet (new user).
  // Return identity fields so the form can still display mobile/email.
  if (result.response.status === 404) {
    return NextResponse.json({ ok: true, profile: { mobile_number, email } }, { status: 200 });
  }
  if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });

  // result.data is UserProfileRetrievalResponse = { profile?: UserProfile, error?: Error }
  // Unwrap the nested profile object — do NOT spread result.data directly or the
  // form state will contain a "profile" key which gets sent back on PATCH.
  const responseData = result.data as Record<string, unknown> ?? {};
  const raw = (responseData.profile as Record<string, unknown>) ?? {};

  // date_of_birth: the DB stores it as a DATE, the gateway serialises it via
  // protojson as an RFC3339 timestamp string (e.g. "1990-06-15T00:00:00Z").
  // Convert to YYYY-MM-DD for the HTML date input.
  // Treat zero-epoch / sentinel values (year ≤ 1970 or year === 1900) as empty
  // so new users see a blank date field rather than a confusing placeholder date.
  let dateOfBirth = "";
  const rawDOB = raw.date_of_birth as string | undefined;
  if (rawDOB) {
    const d = new Date(rawDOB);
    if (!isNaN(d.getTime()) && d.getFullYear() > 1970 && d.getFullYear() !== 1900) {
      dateOfBirth = d.toISOString().slice(0, 10); // "YYYY-MM-DD"
    }
  }

  // address_line1 is the proto field name — expose it as both address_line1 and
  // address so the form's single "Address" input maps correctly.
  const address_line1 = (raw.address_line1 as string) ?? "";

  const profile = {
    ...raw,
    date_of_birth: dateOfBirth,
    address_line1,
    // Convenience alias used by the form's "Address" field
    address: address_line1,
    mobile_number,
    email,
  };
  return NextResponse.json({ ok: true, profile }, { status: 200 });
}

/** PATCH /api/auth/profile -- update current user profile */
export async function PATCH(request: Request) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  let body: Record<string, unknown>;
  try { body = await request.json() as Record<string, unknown>; } catch { return badRequest("Invalid request body"); }
  const userId = await resolveUserIdFromSession(request, hdrs);
  if (!userId) return NextResponse.json({ ok: false, message: "Cannot resolve user identity" }, { status: 401 });

  // Build a clean payload for protojson.Unmarshal in the gateway:
  // 1. user_id is required by UpdateUserProfileRequest.
  // 2. date_of_birth must be RFC3339 (gateway uses protojson) — HTML date gives "YYYY-MM-DD".
  // 3. address (form convenience alias) maps to address_line1 (proto field name).
  // 4. Strip read-only identity fields (email, mobile_number) — User record only.
  // 5. Strip the "address" alias and any stale "profile" nesting from old form state.
  const transformed: Record<string, unknown> = {
    user_id: userId,
    full_name:      body.full_name      ?? "",
    occupation:     body.occupation     ?? "",
    employer:       body.employer       ?? "",
    address_line1:  body.address_line1  ?? body.address ?? "",
    address_line2:  body.address_line2  ?? "",
    city:           body.city           ?? "",
    district:       body.district       ?? "",
    division:       body.division       ?? "",
    country:        body.country        ?? "",
    postal_code:    body.postal_code    ?? "",
    nid_number:     body.nid_number     ?? "",
    marital_status: body.marital_status ?? "",
    gender:         body.gender         ?? "",
  };

  // Convert YYYY-MM-DD → RFC3339 for protojson Timestamp parsing.
  // Validate the 18-year age requirement (DB CHECK constraint) before hitting gateway.
  const dob = body.date_of_birth as string | undefined;
  if (dob && /^\d{4}-\d{2}-\d{2}$/.test(dob)) {
    const dobDate = new Date(dob);
    const minAge = new Date();
    minAge.setFullYear(minAge.getFullYear() - 18);
    if (dobDate > minAge) {
      return NextResponse.json({ ok: false, message: "Date of birth must be at least 18 years in the past." }, { status: 400 });
    }
    transformed.date_of_birth = `${dob}T00:00:00Z`;
  } else if (dob && dob.includes("T")) {
    transformed.date_of_birth = dob; // already RFC3339
  }
  // If no date_of_birth provided, omit it so the gateway skips the field.

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.updateUserProfile({
    path: { user_id: userId },
    body: transformed as Parameters<typeof sdk.updateUserProfile>[0]['body'],
  });
  if (!result.response.ok) return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
  return NextResponse.json({ ok: true, profile: result.data }, { status: 200 });
}
```

