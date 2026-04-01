import type { PageServerLoad } from './$types';
import { getTenants } from '$lib/server/system-data';

export const load: PageServerLoad = async (event) => {
	return {
		tenants: await getTenants(event)
	};
};
