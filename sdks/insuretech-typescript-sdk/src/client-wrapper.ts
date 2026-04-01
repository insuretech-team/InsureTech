// Custom Client Wrapper for InsureTech SDK
// Provides a configured client instance for use with generated services

import { createClient, createConfig } from './client';

export interface InsureTechClientConfig {
  /** API key for authentication */
  apiKey: string;
  /** Base URL for the API (optional, defaults to production) */
  baseUrl?: string;
  /** Additional headers to include in all requests */
  headers?: Record<string, string>;
}

/**
 * Create a configured client for the InsureTech API
 * 
 * @example
 * ```typescript
 * import { createInsureTechClient, AiService } from '@lifeplus/insuretech-sdk';
 * 
 * const client = createInsureTechClient({
 *   apiKey: 'your-api-key',
 *   baseUrl: 'https://api.insuretech.com'
 * });
 * 
 * // Use with any service method
 * const response = await AiService.aiServiceChat({
 *   client,
 *   body: { message: 'Hello' }
 * });
 * ```
 */
export function createInsureTechClient(config: InsureTechClientConfig) {
  const c = createClient(createConfig({
    baseUrl: config.baseUrl || 'https://api.insuretech.com',
    headers: {
      'Authorization': `Bearer ${config.apiKey}`,
      ...config.headers,
    },
  }));

  // ── Unwrap ApiResponse envelope ─────────────────────────────────────────
  // The gateway wraps every response as { success, data, error, meta }.
  // hey-api puts the parsed JSON into result.data, so without this
  // interceptor consumers would need result.data.data to reach the payload.
  // By replacing the Response body with just the inner "data" field we make
  // result.data === T directly — no double-wrap.
  c.interceptors.response.use(async (response) => {
    const ct = response.headers.get('content-type') ?? '';
    if (!ct.includes('application/json')) return response;
    // Clone so we can read the body without consuming the original.
    const text = await response.clone().text();
    if (!text) return response;
    try {
      const envelope = JSON.parse(text);
      // Only unwrap if it looks like our standard ApiResponse envelope.
      if (
        typeof envelope === 'object' &&
        envelope !== null &&
        'success' in envelope &&
        'data' in envelope
      ) {
        // Success: unwrap envelope.data so result.data === T
        // Error: unwrap envelope.error so result.error has gateway error details
        const inner = envelope.success ? envelope.data : envelope.error;

        // Preserve Set-Cookie and X-CSRF-Token across the body rewrite.
        // Set-Cookie is a forbidden header in the Fetch API — constructing a
        // new Response(..., { headers }) silently drops it in both browser and
        // Node.js (undici). We copy it to the readable header x-set-cookie so
        // that server-side Next.js API route handlers (e.g. the login route)
        // can still forward the session cookie to the browser.
        const newHeaders = new Headers(response.headers);
        const setCookie = response.headers.get('set-cookie');
        if (setCookie) newHeaders.set('x-set-cookie', setCookie);
        const csrfToken = response.headers.get('x-csrf-token');
        if (csrfToken) newHeaders.set('x-csrf-token', csrfToken);

        return new Response(JSON.stringify(inner ?? {}), {
          status: response.status,
          statusText: response.statusText,
          headers: newHeaders,
        });
      }
    } catch { /* not JSON — pass through */ }
    return response;
  });

  return c;
}

// Re-export for convenience
export { createClient, createConfig } from './client';
