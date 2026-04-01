import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { getPartnerDetail } from '$lib/server/system-data';

export const load: PageServerLoad = async (event) => {
	const partner = await getPartnerDetail(event, event.params.id);

	if (!partner) {
		throw error(404, 'Partner not found');
	}

	return { partner };
};
