import { redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

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

export const POST: RequestHandler = async ({ cookies }) => {
	for (const name of cookieNames) {
		cookies.delete(name, { path: '/' });
	}

	throw redirect(302, '/login');
};
