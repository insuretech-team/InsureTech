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
