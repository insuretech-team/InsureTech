import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { getProductDetail } from '$lib/server/system-data';

export const load: PageServerLoad = async (event) => {
	const product = await getProductDetail(event, event.params.id);

	if (!product) {
		throw error(404, 'Product not found');
	}

	return { product };
};
