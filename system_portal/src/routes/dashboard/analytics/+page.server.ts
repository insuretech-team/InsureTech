import type { PageServerLoad } from './$types';
import {
	getClaimExposure,
	getOverviewData,
	getPartnerMix,
	getPolicyHealth
} from '$lib/server/system-data';

export const load: PageServerLoad = async (event) => {
	const overview = await getOverviewData(event);

	return {
		overview,
		partnerMix: getPartnerMix(overview.partners),
		policyHealth: getPolicyHealth(overview.policies),
		claimExposure: getClaimExposure(overview.claims)
	};
};
