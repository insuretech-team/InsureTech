import type { PageServerLoad } from './$types';
import { getClaims } from '$lib/server/system-data';

export const load: PageServerLoad = async (event) => {
	return {
		claims: await getClaims(event)
	};
};
