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
  return createClient(createConfig({
    baseUrl: config.baseUrl || 'https://api.insuretech.com',
    headers: {
      'Authorization': `Bearer ${config.apiKey}`,
      ...config.headers,
    },
  }));
}

// Re-export for convenience
export { createClient, createConfig } from './client';
