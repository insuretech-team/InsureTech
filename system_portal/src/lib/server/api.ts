import { createInsureTechClient, type InsureTechClientConfig } from '@lifeplus/insuretech-sdk';

/**
 * Creates a server-side SDK client authenticated with the current session token.
 * This should be used strictly in Server Load functions or Actions.
 * 
 * @param sessionToken - The HttpOnly session token from the user's cookies.
 */
export function getApiClient(sessionToken?: string) {
    const config: InsureTechClientConfig = {
        apiKey: sessionToken || '',
        baseUrl: process.env.VITE_API_URL || 'http://localhost:8080' // default to local backend in dev
    };

    return createInsureTechClient(config);
}
