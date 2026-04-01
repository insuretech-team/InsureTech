import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

const SESSION_COOKIE_NAME = "session_token";
const PUBLIC_PATHS = ["/login", "/api/auth/login"];

// The KYC page is accessible to authenticated-but-unverified users.
// The KYC gate reads the lightweight "portal_kyc_verified" cookie set at login.
// If absent (existing sessions before this feature), we do NOT block — the
// /kyc page itself checks the DB and redirects if already verified.
const KYC_PATH = "/kyc";
const RESET_PASSWORD_PATH = "/reset-password";
const EMPLOYEE_HOME_PATH = "/employee";

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
  { prefix: "/billing-invoices", allowedRoles: ["SYSTEM_ADMIN", "B2B_ORG_ADMIN", "BUSINESS_ADMIN", "HR_MANAGER", "VIEWER"] },
  { prefix: "/insurance-plans", allowedRoles: ["SYSTEM_ADMIN", "B2B_ORG_ADMIN", "BUSINESS_ADMIN", "HR_MANAGER", "VIEWER"] },
  { prefix: "/settings", allowedRoles: ["SYSTEM_ADMIN", "B2B_ORG_ADMIN", "BUSINESS_ADMIN", "HR_MANAGER", "VIEWER"] },
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
    const role = request.cookies.get("portal_role")?.value ?? "";
    if (role === "B2B_BENEFICIARY") {
      return NextResponse.redirect(new URL(EMPLOYEE_HOME_PATH, request.url));
    }
    const kycVerified = request.cookies.get("portal_kyc_verified")?.value;
    const passwordChangeRequired = request.cookies.get("portal_password_change_required")?.value;
    if (kycVerified === "false" || (passwordChangeRequired === "true" && !kycVerified)) {
      return NextResponse.redirect(new URL(KYC_PATH, request.url));
    }
    if (passwordChangeRequired === "true") {
      return NextResponse.redirect(new URL(RESET_PASSWORD_PATH, request.url));
    }
    const dest = role === "SYSTEM_ADMIN" ? "/organisations" : "/";
    return NextResponse.redirect(new URL(dest, request.url));
  }

  const role = request.cookies.get("portal_role")?.value ?? "";
  const isBeneficiary = role === "B2B_BENEFICIARY";

  if (hasSessionCookie && isBeneficiary) {
    if (pathname === "/") {
      return NextResponse.redirect(new URL(EMPLOYEE_HOME_PATH, request.url));
    }
    if (pathname === KYC_PATH || pathname === RESET_PASSWORD_PATH) {
      return NextResponse.redirect(new URL(EMPLOYEE_HOME_PATH, request.url));
    }
  }

  // KYC gate: redirect unverified admins to /kyc on first login.
  // Uses the lightweight "portal_kyc_verified" cookie set at login time.
  // Cookie absent = existing session (pre-KYC feature) → do NOT block (backward compat).
  // Cookie = "false" → user not yet KYC verified → redirect to /kyc.
  if (hasSessionCookie && !isBeneficiary && pathname !== KYC_PATH && !pathname.startsWith("/api/")) {
    const kycVerified = request.cookies.get("portal_kyc_verified")?.value;
    const passwordChangeRequired = request.cookies.get("portal_password_change_required")?.value;
    if (kycVerified === "false" || (passwordChangeRequired === "true" && !kycVerified)) {
      return NextResponse.redirect(new URL(KYC_PATH, request.url));
    }
  }

  if (hasSessionCookie && !isBeneficiary && pathname !== RESET_PASSWORD_PATH && !pathname.startsWith("/api/")) {
    const kycVerified = request.cookies.get("portal_kyc_verified")?.value;
    const passwordChangeRequired = request.cookies.get("portal_password_change_required")?.value;
    if (passwordChangeRequired === "true" && kycVerified !== "false") {
      return NextResponse.redirect(new URL(RESET_PASSWORD_PATH, request.url));
    }
  }

  if (hasSessionCookie && pathname === RESET_PASSWORD_PATH) {
    if (isBeneficiary) {
      return NextResponse.redirect(new URL(EMPLOYEE_HOME_PATH, request.url));
    }
    const kycVerified = request.cookies.get("portal_kyc_verified")?.value;
    if (kycVerified === "false") {
      return NextResponse.redirect(new URL(KYC_PATH, request.url));
    }
  }

  // Role-based route guard (UX-level, uses portal_role cookie set at login)
  if (hasSessionCookie) {
    if (role) {
      const guard = ROLE_GUARDS.find((g) => pathname.startsWith(g.prefix));
      if (guard && !guard.allowedRoles.includes(role)) {
        // Redirect to appropriate default page for their role
        const fallback = role === "SYSTEM_ADMIN" ? "/organisations" : role === "B2B_BENEFICIARY" ? EMPLOYEE_HOME_PATH : "/";
        return NextResponse.redirect(new URL(fallback, request.url));
      }
    }
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/((?!.*\\..*).*)"],
};
