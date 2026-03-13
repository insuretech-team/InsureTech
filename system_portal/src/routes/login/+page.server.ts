import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { getApiClient } from '$lib/server/api';
import { authServiceLogin } from '@lifeplus/insuretech-sdk';

export const load: PageServerLoad = async ({ locals }) => {
    // If already logged in, redirect away from login
    if (locals.user) {
        throw redirect(302, '/dashboard');
    }
    return {};
};

export const actions: Actions = {
    default: async ({ request, cookies }) => {
        const data = await request.formData();
        const email = data.get('email')?.toString();
        const password = data.get('password')?.toString();

        if (!email || !password) {
            return fail(400, {
                error: 'Email and password are required',
                email
            });
        }

        try {
            // Get an unauthenticated client instance
            const client = getApiClient();

            // Call the SDK login endpoint with device_type: 'WEB' for server-side sessions
            const res = await authServiceLogin({
                client,
                body: {
                    mobile_number: email, // Assuming email maps to mobile/username field in auth
                    password,
                    device_id: 'web-portal',
                    device_type: 'WEB', // Triggers server-side session instead of JWT
                    device_name: 'System Portal'
                }
            });

            if (res.data && res.data.session_token) {
                // Set the secure session cookie
                cookies.set('session', res.data.session_token, {
                    path: '/',
                    httpOnly: true,
                    sameSite: 'lax',
                    secure: process.env.NODE_ENV === 'production',
                    maxAge: 60 * 60 * 24 * 30 // 30 days
                });

                // Successfully authenticated, redirect to the dashboard
                throw redirect(302, '/dashboard');
            } else {
                // Should not happen if authn succeeds, but fallback
                return fail(401, {
                    error: 'Invalid response from authentication server',
                    email
                });
            }

        } catch (err: any) {
            if (err.status === 302) {
                throw err; // Re-throw SvelteKit redirects
            }
            console.error('Login error:', err);
            return fail(401, {
                error: err?.body?.message || 'Invalid credentials or unable to connect to authentication server',
                email
            });
        }
    }
};
