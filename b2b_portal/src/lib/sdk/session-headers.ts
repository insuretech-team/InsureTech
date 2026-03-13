/**
 * session-headers.ts
 * ──────────────────
 * Server-side helper that resolves the correct x-portal, x-user-id,
 * x-business-id, x-tenant-id headers to inject into backend SDK calls.
 *
 * The portal is STATELESS — session data lives in the backend cookie, not a
 * local store. We resolve the portal role from the lightweight `portal_role`
 * cookie set at login time, and the businessId from the `x-business-id` cookie
 * (also set at login for B2B users).
 *
 * Fix for super-admin 403 on /departments and /purchase-orders:
 * The backend authz interceptor needs x-portal=PORTAL_SYSTEM to route the
 * request to the system:root Casbin domain (no x-business-id required).
 */

export interface PortalHeaders {
  portal: string;
  userId: string;
  businessId: string;
  tenantId: string;
}

/**
 * Extracts a named cookie value from a raw Cookie header string.
 */
function extractCookie(cookieHeader: string, name: string): string {
  const escaped = name.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  const m = cookieHeader.match(new RegExp(`(?:^|;\\s*)${escaped}=([^;]*)`));
  return m ? decodeURIComponent(m[1]) : "";
}

/**
 * Maps the portal principal role to the correct x-portal header value.
 * SYSTEM_ADMIN → PORTAL_SYSTEM (super admin, no org context required)
 * Everything else → PORTAL_B2B (b2b admin / hr manager / viewer)
 */
function roleToPortal(role: string): string {
  if (role === "SYSTEM_ADMIN") return "PORTAL_SYSTEM";
  return "PORTAL_B2B";
}

/**
 * Resolves portal auth headers from the request cookies.
 *
 * Cookie sources (all set by the login route):
 *   portal_role     — e.g. "SYSTEM_ADMIN" | "B2B_ORG_ADMIN" | "BUSINESS_ADMIN"
 *   portal_user_id  — user ID for x-user-id (optional, may be empty)
 *   portal_biz_id   — business/org ID for x-business-id (B2B users only)
 *
 * If portal_role is missing, we fall back to PORTAL_B2B (safest non-elevated default).
 * Returns null only if there's no session_token at all (unauthenticated).
 *
 * Usage in an API route:
 *   const hdrs = await resolvePortalHeaders(request);
 *   const sdk = makeSdkClient(request, hdrs ?? undefined);
 */
export async function resolvePortalHeaders(request: Request): Promise<PortalHeaders | null> {
  const cookieHeader = request.headers.get("cookie") ?? "";

  // Require a backend session cookie — if absent the request is unauthenticated
  const sessionToken = extractCookie(cookieHeader, "session_token");
  if (!sessionToken) return null;

  // Read lightweight metadata cookies written at login time
  const role = extractCookie(cookieHeader, "portal_role") || "BUSINESS_ADMIN";
  const userId = extractCookie(cookieHeader, "portal_user_id");
  const businessId = extractCookie(cookieHeader, "portal_biz_id");

  const portal = roleToPortal(role);
  const tenantId =
    process.env.DEFAULT_TENANT_ID?.trim() ||
    process.env.NEXT_PUBLIC_DEFAULT_TENANT_ID?.trim() ||
    "00000000-0000-0000-0000-000000000001";

  return { portal, userId, businessId, tenantId };
}
