import type { Handle } from '@sveltejs/kit';
import { getApiClient } from '$lib/server/api';
import { authServiceGetCurrentSession } from '@lifeplus/insuretech-sdk';

export const handle: Handle = async ({ event, resolve }) => {
    // 1. Get the session token from cookies
    const sessionToken = event.cookies.get('session');

    if (!sessionToken) {
        event.locals.user = null;
        return resolve(event);
    }

    try {
        // 2. Initialize the client using the session token to verify authenticity
        const client = getApiClient(sessionToken);

        // 3. Optional: Verify session against the Auth service.
        // It relies on the token being passed as a Bearer in the API client.
        const response = await authServiceGetCurrentSession({ client });

        if (response.data && response.data.session?.session_id) {
            // Populate locals with the validated session
            event.locals.user = {
                id: response.data.session.user_id,
                sessionId: response.data.session.session_id,
                token: sessionToken
            };
        } else {
            // If the backend doesn't recognize the token, clear it
            event.locals.user = null;
            event.cookies.delete('session', { path: '/' });
        }
    } catch (error) {
        console.error('Session validation failed:', error);
        event.locals.user = null;
        event.cookies.delete('session', { path: '/' });
    }

    return resolve(event);
};
