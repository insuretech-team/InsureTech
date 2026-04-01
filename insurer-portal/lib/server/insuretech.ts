import { NextResponse } from "next/server";
import {
  authServiceGetCurrentSession,
  authServiceLogin,
  authServiceLogout,
  createInsureTechClient,
  type LoginResponse,
} from "@lifeplus/insuretech-sdk";

import type { PortalRole, PortalSessionSnapshot } from "@/lib/types";

export const SESSION_COOKIE_NAME = "session_token";
export const CSRF_COOKIE_NAME = "csrf_token";

type DirectHttpResult = {
  ok: boolean;
  status: number;
  data: Record<string, unknown>;
  message?: string;
};

function getBaseUrl() {
  return (
    process.env.INSURETECH_API_BASE_URL ??
    process.env.NEXT_PUBLIC_INSURETECH_API_BASE_URL ??
    "http://localhost:8080"
  );
}

export function extractCookie(cookieHeader: string, name: string) {
  const escaped = name.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  const match = cookieHeader.match(new RegExp(`(?:^|;\\s*)${escaped}=([^;]*)`));
  return match ? decodeURIComponent(match[1]) : "";
}

function mapUserTypeToRole(userType?: string): PortalRole {
  const value = (userType ?? "").toUpperCase();
  if (value.includes("SYSTEM")) return "SYSTEM_ADMIN";
  if (value.includes("B2B_ORG_ADMIN")) return "B2B_ORG_ADMIN";
  if (value.includes("BUSINESS_ADMIN")) return "BUSINESS_ADMIN";
  if (value.includes("HR")) return "HR_MANAGER";
  if (value.includes("VIEWER")) return "VIEWER";
  if (value.includes("PARTNER_ADMIN")) return "PARTNER_ADMIN";
  return "PARTNER_USER";
}

export function roleToPortal(role: string) {
  return role === "SYSTEM_ADMIN" ? "PORTAL_SYSTEM" : "PORTAL_B2B";
}

function buildHeaders(cookieHeader: string, session?: Partial<PortalSessionSnapshot>) {
  const headers: Record<string, string> = {};
  if (cookieHeader) headers.cookie = cookieHeader;

  const csrf = extractCookie(cookieHeader, CSRF_COOKIE_NAME);
  if (csrf) headers["X-CSRF-Token"] = csrf;

  const role = session?.role ?? extractCookie(cookieHeader, "portal_role") ?? "PARTNER_USER";
  const portal = session?.portal ?? roleToPortal(role);
  const userId = session?.userId ?? extractCookie(cookieHeader, "portal_user_id");
  const businessId = session?.businessId ?? extractCookie(cookieHeader, "portal_biz_id");
  const tenantId =
    process.env.DEFAULT_TENANT_ID?.trim() ||
    process.env.NEXT_PUBLIC_DEFAULT_TENANT_ID?.trim() ||
    "00000000-0000-0000-0000-000000000001";

  if (portal) headers["x-portal"] = portal;
  if (userId) headers["x-user-id"] = userId;
  if (businessId) headers["x-business-id"] = businessId;
  if (tenantId) headers["x-tenant-id"] = tenantId;

  return headers;
}

function buildClient(cookieHeader: string, session?: Partial<PortalSessionSnapshot>) {
  return createInsureTechClient({
    apiKey: process.env.INSURETECH_API_KEY ?? "",
    baseUrl: getBaseUrl(),
    headers: buildHeaders(cookieHeader, session),
  });
}

export async function fetchCurrentSession(
  request: Request,
  cookieOverride?: string,
): Promise<PortalSessionSnapshot | null> {
  const cookieHeader = cookieOverride ?? request.headers.get("cookie") ?? "";
  if (!extractCookie(cookieHeader, SESSION_COOKIE_NAME)) return null;

  const client = buildClient(cookieHeader);
  const result = await authServiceGetCurrentSession({
    client,
    throwOnError: false,
  });

  if (!result.response.ok || !result.data) {
    return null;
  }

  const data = result.data as Record<string, unknown>;
  const session = (data.session as Record<string, unknown> | undefined) ?? {};
  const role = mapUserTypeToRole((data.user_type as string | undefined) ?? "");

  return {
    userId:
      (session.user_id as string | undefined) ??
      (data.user_id as string | undefined) ??
      extractCookie(cookieHeader, "portal_user_id"),
    sessionId: (session.session_id as string | undefined) ?? "",
    role,
    portal: roleToPortal(role),
    businessId:
      (data.business_id as string | undefined) ??
      extractCookie(cookieHeader, "portal_biz_id"),
    email: extractCookie(cookieHeader, "portal_email"),
    mobile: extractCookie(cookieHeader, "portal_mobile"),
    expiresAt: (session.expires_at as string | undefined) ?? "",
  };
}

export async function directHttp(
  request: Request,
  path: string,
  init?: {
    method?: "GET" | "POST" | "PATCH" | "PUT" | "DELETE";
    body?: unknown;
    session?: Partial<PortalSessionSnapshot>;
  },
): Promise<DirectHttpResult> {
  const cookieHeader = request.headers.get("cookie") ?? "";
  const response = await fetch(`${getBaseUrl()}${path}`, {
    method: init?.method ?? "GET",
    headers: {
      "Content-Type": "application/json",
      ...buildHeaders(cookieHeader, init?.session),
    },
    body: init?.body === undefined ? undefined : JSON.stringify(init.body),
    cache: "no-store",
  });

  const raw = await response.text();
  let parsed: Record<string, unknown> = {};
  try {
    parsed = raw ? (JSON.parse(raw) as Record<string, unknown>) : {};
  } catch {
    parsed = {};
  }

  if (typeof parsed.success === "boolean") {
    return {
      ok: Boolean(parsed.success),
      status: response.status,
      data: (parsed.success ? parsed.data : parsed.error) as Record<string, unknown> ?? {},
      message:
        !parsed.success && parsed.error && typeof parsed.error === "object"
          ? ((parsed.error as Record<string, unknown>).message as string | undefined)
          : undefined,
    };
  }

  return {
    ok: response.ok,
    status: response.status,
    data: parsed,
    message: !response.ok ? raw : undefined,
  };
}

export function normalizeMobileNumber(value: string): string | null {
  const stripped = value.trim().replace(/[^\d+]/g, "");
  const digits = stripped.startsWith("+") ? stripped.slice(1) : stripped;

  let normalized = "";
  if (digits.startsWith("00880")) normalized = digits.slice(2);
  else if (digits.startsWith("880")) normalized = digits;
  else if (digits.startsWith("0088")) normalized = `880${digits.slice(4)}`;
  else if (digits.startsWith("88") && digits.length === 13) normalized = `880${digits.slice(2)}`;
  else if (digits.startsWith("0")) normalized = `880${digits.slice(1)}`;
  else if (digits.length === 10) normalized = `880${digits}`;

  if (!/^880(13|14|15|16|17|18|19)\d{8}$/.test(normalized)) return null;
  return `+${normalized}`;
}

function getSetCookieHeaders(headers: Headers) {
  const extended = headers as Headers & { getSetCookie?: () => string[] };
  if (typeof extended.getSetCookie === "function") {
    const values = extended.getSetCookie();
    if (values.length) return values;
  }
  const raw = headers.get("x-set-cookie") ?? headers.get("set-cookie");
  return raw ? [raw] : [];
}

export async function loginWithMobile(request: Request, payload: { mobileNumber: string; password: string }) {
  const client = buildClient("");
  return authServiceLogin({
    client,
    throwOnError: false,
    body: {
      mobile_number: payload.mobileNumber,
      password: payload.password,
      device_id: "insurer-portal-web",
      device_type: "WEB",
      device_name: "Insurer Portal Web",
    },
  });
}

export async function buildSessionFromLogin(
  request: Request,
  loginData: LoginResponse | Record<string, unknown>,
  responseHeaders: Headers,
): Promise<{ session: PortalSessionSnapshot | null; sessionToken: string; csrfToken: string }> {
  const data = loginData as LoginResponse;
  const raw = loginData as Record<string, unknown>;
  const user = ((raw.user as Record<string, unknown> | undefined) ?? {}) as Record<string, unknown>;
  const setCookies = getSetCookieHeaders(responseHeaders);
  const sessionToken =
    data.session_token ??
    setCookies.find((value) => value.startsWith(`${SESSION_COOKIE_NAME}=`))?.split("=")[1]?.split(";")[0] ??
    "";
  const csrfToken =
    data.csrf_token ??
    responseHeaders.get("x-csrf-token") ??
    "";

  let session = await fetchCurrentSession(
    request,
    `${SESSION_COOKIE_NAME}=${sessionToken}${csrfToken ? `; ${CSRF_COOKIE_NAME}=${csrfToken}` : ""}`,
  );

  if (!session) {
    const role = mapUserTypeToRole((user.user_type as string | undefined) ?? "");
    session = {
      userId: data.user_id ?? ((user.user_id as string | undefined) ?? ""),
      sessionId: data.session_id ?? "",
      role,
      portal: roleToPortal(role),
      businessId: (user.business_id as string | undefined) ?? "",
      email: (user.email as string | undefined) ?? "",
      mobile: (user.mobile_number as string | undefined) ?? "",
      expiresAt: "",
    };
  }

  return { session, sessionToken, csrfToken };
}

export async function logoutCurrentSession(request: Request, session: PortalSessionSnapshot | null) {
  if (!session) return null;
  const client = buildClient(request.headers.get("cookie") ?? "", session);
  return authServiceLogout({
    client,
    throwOnError: false,
    body: {
      session_id: session.sessionId,
      logout_reason: "user_initiated",
    },
  });
}

export function applyPortalCookies(
  response: NextResponse,
  session: PortalSessionSnapshot | null,
  extras?: { sessionToken?: string; csrfToken?: string; email?: string; mobile?: string },
) {
  const secure = process.env.NODE_ENV === "production";

  if (extras?.sessionToken) {
    response.cookies.set({
      name: SESSION_COOKIE_NAME,
      value: extras.sessionToken,
      httpOnly: true,
      sameSite: "lax",
      secure,
      path: "/",
      maxAge: 60 * 60 * 12,
    });
  }

  if (extras?.csrfToken) {
    response.cookies.set({
      name: CSRF_COOKIE_NAME,
      value: extras.csrfToken,
      httpOnly: true,
      sameSite: "lax",
      secure,
      path: "/",
      maxAge: 60 * 60 * 12,
    });
  }

  if (!session) return;

  const cookieOpts = {
    httpOnly: false,
    sameSite: "lax" as const,
    secure,
    path: "/",
    maxAge: 60 * 60 * 12,
  };

  response.cookies.set({ name: "portal_role", value: session.role, ...cookieOpts });
  response.cookies.set({ name: "portal_user_id", value: session.userId, ...cookieOpts });
  response.cookies.set({ name: "portal_biz_id", value: session.businessId, ...cookieOpts });
  response.cookies.set({ name: "portal_email", value: extras?.email ?? session.email, ...cookieOpts });
  response.cookies.set({ name: "portal_mobile", value: extras?.mobile ?? session.mobile, ...cookieOpts });
}

export function clearPortalCookies(response: NextResponse) {
  [
    SESSION_COOKIE_NAME,
    CSRF_COOKIE_NAME,
    "portal_role",
    "portal_user_id",
    "portal_biz_id",
    "portal_email",
    "portal_mobile",
  ].forEach((name) => {
    response.cookies.set({
      name,
      value: "",
      path: "/",
      maxAge: 0,
    });
  });
}
