import { fail, isRedirect, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import {
	authServiceEmailLogin,
	authServiceSendEmailOtp,
	createInsureTechClient
} from '@lifeplus/insuretech-sdk';
import { extractGatewayError } from '$lib/server/api-helpers';

const DEFAULT_TENANT_ID =
	process.env.DEFAULT_TENANT_ID?.trim() || '00000000-0000-0000-0000-000000000001';

function createAuthClient() {
	return createInsureTechClient({
		apiKey: process.env.INSURETECH_API_KEY ?? 'system-portal',
		baseUrl:
			process.env.INSURETECH_API_BASE_URL ??
			process.env.VITE_API_URL ??
			process.env.PUBLIC_API_URL ??
			'http://localhost:8080'
	});
}

function cookieOptions() {
	return {
		path: '/',
		httpOnly: true,
		sameSite: 'lax' as const,
		secure: process.env.NODE_ENV === 'production',
		maxAge: 60 * 60 * 12
	};
}

function metadataCookieOptions() {
	return {
		path: '/',
		httpOnly: false,
		sameSite: 'lax' as const,
		secure: process.env.NODE_ENV === 'production',
		maxAge: 60 * 60 * 12
	};
}

function maskEmail(email: string) {
	const [name, domain] = email.split('@');
	if (!name || !domain) return email;
	return `${name.slice(0, 2)}${'*'.repeat(Math.max(1, name.length - 2))}@${domain}`;
}

export const load: PageServerLoad = async ({ locals }) => {
	if (locals.user) {
		throw redirect(302, '/dashboard');
	}

	return {};
};

export const actions: Actions = {
	sendOtp: async ({ request }) => {
		const form = await request.formData();
		const email = form.get('email')?.toString().trim().toLowerCase() ?? '';

		if (!email) {
			return fail(400, {
				step: 'request',
				email,
				error: 'Enter the system user email address to continue.'
			});
		}

		try {
			const result = await authServiceSendEmailOtp({
				client: createAuthClient(),
				throwOnError: false,
				body: {
					email,
					type: 'email_login'
				}
			});

			const otpId = result.data?.otp_id;
			if (!otpId) {
				return fail(400, {
					step: 'request',
					email,
					error: extractGatewayError(result)
				});
			}

			return {
				step: 'verify',
				email,
				otpId,
				expiresIn: result.data?.expires_in_seconds ?? 300,
				maskedEmail: maskEmail(email),
				message: 'OTP sent. Check the inbox for the system user account.'
			};
		} catch (error) {
			return fail(400, {
				step: 'request',
				email,
				error: extractGatewayError(error)
			});
		}
	},

	verifyOtp: async ({ request, cookies }) => {
		const form = await request.formData();
		const email = form.get('email')?.toString().trim().toLowerCase() ?? '';
		const otpId = form.get('otpId')?.toString().trim() ?? '';
		const code = form.get('code')?.toString().trim() ?? '';

		if (!email || !otpId || !code) {
			return fail(400, {
				step: 'verify',
				email,
				otpId,
				error: 'Enter the 6-digit OTP to complete sign in.'
			});
		}

		try {
			const result = await authServiceEmailLogin({
				client: createAuthClient(),
				throwOnError: false,
				body: {
					email,
					otp_id: otpId,
					code,
					device_id: crypto.randomUUID(),
					device_name: 'System Portal Web'
				}
			});

			if (!result.data?.session_token) {
				return fail(401, {
					step: 'verify',
					email,
					otpId,
					error: extractGatewayError(result)
				});
			}

			cookies.set('session_token', result.data.session_token, cookieOptions());

			if (result.data.csrf_token) {
				cookies.set('csrf_token', result.data.csrf_token, metadataCookieOptions());
			}

			cookies.set('portal_role', 'SYSTEM_ADMIN', metadataCookieOptions());
			cookies.set(
				'portal_user_id',
				result.data.user_id ?? result.data.user?.user_id ?? '',
				metadataCookieOptions()
			);
			cookies.set('portal_biz_id', '', metadataCookieOptions());
			cookies.set('portal_email', result.data.user?.email ?? email, metadataCookieOptions());
			cookies.set(
				'portal_mobile',
				result.data.user?.mobile_number ?? '',
				metadataCookieOptions()
			);
			cookies.set('portal_tenant_id', DEFAULT_TENANT_ID, metadataCookieOptions());

			throw redirect(302, '/dashboard');
		} catch (error) {
			if (isRedirect(error)) throw error;

			return fail(401, {
				step: 'verify',
				email,
				otpId,
				error: extractGatewayError(error)
			});
		}
	}
};
