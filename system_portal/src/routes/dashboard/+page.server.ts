import type { PageServerLoad } from './$types';
import {
	getClaimExposure,
	getOperationalBreakdown,
	getOverviewData,
	getPartnerMix,
	getPolicyHealth,
	getTenantSummary
} from '$lib/server/system-data';

export const load: PageServerLoad = async (event) => {
	const overview = await getOverviewData(event);

	return {
		overview,
		breakdown: getOperationalBreakdown(overview),
		partnerMix: getPartnerMix(overview.partners),
		policyHealth: getPolicyHealth(overview.policies),
		claimExposure: getClaimExposure(overview.claims),
		tenantSummary: getTenantSummary(overview.tenants)
	};
};
