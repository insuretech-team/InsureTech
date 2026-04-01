import type { PageServerLoad } from './$types';
import { getPolicies } from '$lib/server/system-data';

export const load: PageServerLoad = async (event) => {
	return {
		policies: await getPolicies(event)
	};
};
