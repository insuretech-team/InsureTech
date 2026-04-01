import type { PageServerLoad } from './$types';
import { getReports } from '$lib/server/system-data';

export const load: PageServerLoad = async (event) => {
	return {
		reports: await getReports(event)
	};
};
