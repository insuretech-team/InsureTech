import type { PageServerLoad } from './$types';
import { getPartners } from '$lib/server/system-data';

export const load: PageServerLoad = async (event) => {
	return {
		partners: await getPartners(event, 'non-life')
	};
};
