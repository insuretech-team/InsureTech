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
  matcher: ["/((?!.*\\..*).*)"],
};
