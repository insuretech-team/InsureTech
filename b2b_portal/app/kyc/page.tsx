import Image from "next/image";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { requireServerSession, SESSION_COOKIE_NAME } from "@lib/auth/session";
import { makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { resolveUserIdFromSession } from "@lib/auth/resolve-user-id";
import { KYCPageClient } from "./kyc-page-client";

// Force dynamic rendering — this page reads cookies and must never be statically cached.
export const dynamic = "force-dynamic";

/**
 * /kyc — Identity verification gate for B2B admins.
 *
 * This page is shown when:
 *   - The admin is authenticated (has a session cookie)
 *   - Their user profile has kyc_verified = false
 *
 * After the eKYC flow completes, they are redirected to the dashboard
 * and their account shows PENDING_REVIEW until an InsureTech reviewer
 * calls ApproveKYC.
 *
 * NOTE: KYC record initiation (creating the kyc_verifications row) is done
 * client-side inside KYCPageClient to avoid a self-referential server fetch
 * loop in dev mode. The server only checks whether the user is already verified.
 */
export default async function KYCPage() {
  await requireServerSession();
  const cookieStore = await cookies();
  const sessionToken = cookieStore.get(SESSION_COOKIE_NAME)?.value;
  if (!sessionToken) redirect("/login");

  // Build a dummy Request with the real backend session token so auth-backed
  // helpers validate against the gateway correctly. Using the in-memory
  // session UUID here causes /kyc -> /login -> /kyc redirect loops.
  const cookieHeader = `${SESSION_COOKIE_NAME}=${sessionToken}`;
  const dummyReq = new Request("http://localhost/api/kyc", {
    headers: { cookie: cookieHeader },
  });

  const hdrs = await resolvePortalHeaders(dummyReq);
  if (!hdrs) redirect("/login");

  const userId = await resolveUserIdFromSession(dummyReq, hdrs);
  if (!userId) redirect("/login");

  // Check if already KYC verified — if so, sync cookie and go to dashboard.
  const sdk = makeSdkClient(dummyReq, hdrs);
  let kycVerified = false;
  try {
    const profileResult = await sdk.getUserProfile({ path: { user_id: userId } });
    if (profileResult.response.ok && profileResult.data) {
      const raw = (profileResult.data as Record<string, unknown>);
      const profile = (raw.profile ?? raw) as Record<string, unknown>;
      kycVerified = Boolean(profile.kyc_verified);
    }
  } catch {
    // If profile fetch fails, show KYC page (safer to verify than skip)
  }

  if (kycVerified) {
    redirect("/api/auth/kyc/sync?status=true&next=/");
  }

  // KYC record initiation is handled client-side in KYCPageClient.
  // Doing it here (server-side) caused a self-referential fetch loop because
  // the server component was calling back into the same Next.js server
  // (http://localhost:3000/api/auth/kyc/initiate) on every render in dev mode.

  return (
    <div className="relative min-h-screen overflow-hidden bg-[radial-gradient(1100px_circle_at_15%_20%,rgb(var(--brand-jungle-rgb)/0.18),transparent_55%),radial-gradient(900px_circle_at_85%_85%,rgb(var(--brand-cold-rgb)/0.18),transparent_50%)]">
      {/* Grid overlay — matches login page */}
      <div className="pointer-events-none absolute inset-0 opacity-30 [background-image:linear-gradient(to_right,rgb(var(--brand-cold-rgb)/0.16)_1px,transparent_1px),linear-gradient(to_bottom,rgb(var(--brand-cold-rgb)/0.16)_1px,transparent_1px)] [background-size:32px_32px]" />

      <main className="relative mx-auto flex min-h-screen w-full max-w-2xl flex-col items-center justify-center px-4 py-12 sm:px-8">
        {/* Logo */}
        <div className="mb-8 animate-in fade-in duration-700">
          <Image
            src="/logos/insuretech-brand.png"
            alt="InsureTech"
            width={220}
            height={72}
            style={{ width: "auto", height: "auto" }}
            priority
          />
        </div>

        {/* Heading */}
        <div className="mb-8 text-center animate-in fade-in slide-in-from-bottom-4 duration-700">
          <p className="text-xs font-semibold uppercase tracking-[0.2em] text-[rgb(var(--brand-cold-rgb))] mb-2">Identity Verification</p>
          <h1 className="text-2xl font-semibold tracking-tight text-foreground sm:text-3xl mb-2">Verify Your Identity</h1>
          <p className="text-sm text-muted-foreground max-w-sm">
            Complete a quick face liveness check to access the B2B admin portal. This is a one-time process.
          </p>
        </div>

        {/* KYC Flow */}
        <div className="w-full animate-in fade-in slide-in-from-bottom-4 duration-700 delay-100">
          <KYCPageClient userId={userId} />
        </div>
      </main>
    </div>
  );
}
