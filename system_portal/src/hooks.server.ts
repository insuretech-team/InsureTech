import type { Handle } from '@sveltejs/kit';
import type { CurrentSessionRetrievalResponse } from '@lifeplus/insuretech-sdk';
import { makeSdkClient, resolvePortalHeaders } from '$lib/server/sdk-client';
import { extractGatewayError } from '$lib/server/api-helpers';

const SESSION_COOKIE_NAME = 'session_token';

function mapUserTypeToRole(userType: string | undefined): string {
	const type = String(userType ?? '').toUpperCase();

	if (type === '4' || type === 'USER_TYPE_SYSTEM_USER' || type === 'SYSTEM_USER') {
		return 'SYSTEM_ADMIN';
	}

	if (type === '8' || type === 'USER_TYPE_B2B_ORG_ADMIN' || type === 'B2B_ORG_ADMIN') {
		return 'B2B_ORG_ADMIN';
	}

	if (type === '7' || type === 'USER_TYPE_BUSINESS_ADMIN' || type === 'BUSINESS_ADMIN') {
		return 'BUSINESS_ADMIN';
	}

	return 'SYSTEM_ADMIN';
}

function clearPortalCookies(event: Parameters<Handle>[0]['event']) {
	const cookieNames = [
		'session_token',
		'session',
		'csrf_token',
		'portal_role',
		'portal_user_id',
		'portal_biz_id',
		'portal_email',
		'portal_mobile',
		'portal_tenant_id',
		'portal_password_change_required'
	];

	for (const name of cookieNames) {
		event.cookies.delete(name, { path: '/' });
	}
}

function toPortalUser(
	event: Parameters<Handle>[0]['event'],
	data: CurrentSessionRetrievalResponse,
	sessionToken: string
): NonNullable<App.Locals['user']> | null {
	const session = data.session;
	if (!session?.session_id || !session.user_id) {
		return null;
	}

	const expiresAtRaw =
		session.refresh_token_expires_at ??
		session.access_token_expires_at ??
		new Date(Date.now() + 12 * 60 * 60 * 1000).toISOString();

	return {
		id: session.user_id,
		sessionId: session.session_id,
		token: sessionToken,
		role: mapUserTypeToRole(data.user_type),
		email: event.cookies.get('portal_email') ?? '',
		mobileNumber: event.cookies.get('portal_mobile') ?? '',
		expiresAt: Date.parse(expiresAtRaw)
	};
}

export const handle: Handle = async ({ event, resolve }) => {
	event.locals.user = null;
	event.locals.portalContext = null;

	const sessionToken = event.cookies.get(SESSION_COOKIE_NAME) ?? event.cookies.get('session');
	if (!sessionToken) {
		return resolve(event);
	}

	try {
		const sdk = makeSdkClient(event);
		const result = await sdk.authServiceGetCurrentSession({});

		if (!result.data?.session?.session_id) {
			clearPortalCookies(event);
			return resolve(event);
		}

		const user = toPortalUser(event, result.data, sessionToken);
		if (!user) {
			clearPortalCookies(event);
			return resolve(event);
		}

		event.locals.user = user;
		event.locals.portalContext =
			resolvePortalHeaders(event) ??
			({
				portal: 'PORTAL_SYSTEM',
				userId: user.id,
				businessId: '',
				tenantId:
					process.env.DEFAULT_TENANT_ID?.trim() ||
					event.cookies.get('portal_tenant_id') ||
					'00000000-0000-0000-0000-000000000001'
			} satisfies NonNullable<App.Locals['portalContext']>);
	} catch (error) {
		console.error('System portal session validation failed:', error);
		console.warn('Gateway session error:', extractGatewayError(error));
		clearPortalCookies(event);
	}

	return resolve(event);
};
