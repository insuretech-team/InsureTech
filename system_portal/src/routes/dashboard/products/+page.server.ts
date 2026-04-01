import type { PageServerLoad } from './$types';
import { getProducts } from '$lib/server/system-data';

export const load: PageServerLoad = async (event) => {
	return {
		products: await getProducts(event)
	};
};
