import type { PageServerLoad } from './$types';
import { getAuditEvents } from '$lib/server/system-data';

export const load: PageServerLoad = async (event) => {
	return {
		events: await getAuditEvents(event)
	};
};
