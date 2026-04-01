import type { RequestEvent } from '@sveltejs/kit';
import {
	getNumberField,
	getRecord,
	getStringField,
	parseMoneyDecimal
} from '$lib/server/api-helpers';
import { makeSdkClient } from '$lib/server/sdk-client';
import type {
	SystemAuditEvent,
	SystemClaim,
	SystemOverviewData,
	SystemPartner,
	SystemPolicy,
	SystemProduct,
	SystemReport,
	SystemTenant
} from '$lib/system/types';
import { formatCompactNumber, formatCurrency } from '$lib/system/format';

type JsonMap = Record<string, unknown>;

function asArray(value: unknown): JsonMap[] {
	return Array.isArray(value) ? value.map((item) => getRecord(item)) : [];
}

function pickArray(source: JsonMap, ...keys: string[]) {
	for (const key of keys) {
		if (Array.isArray(source[key])) {
			return asArray(source[key]);
		}
	}

	return [] as JsonMap[];
}

function classifyPartnerCategory(type: string): 'life' | 'non-life' | 'other' {
	const normalized = type.toUpperCase();

	if (
		normalized.includes('HOSPITAL') ||
		normalized.includes('PHARMACY') ||
		normalized.includes('DOCTOR') ||
		normalized.includes('AMBULANCE')
	) {
		return 'life';
	}

	if (
		normalized.includes('AUTO') ||
		normalized.includes('REPAIR') ||
		normalized.includes('LAPTOP') ||
		normalized.includes('MOBILE')
	) {
		return 'non-life';
	}

	return 'other';
}

function normalizeProduct(item: JsonMap): SystemProduct {
	return {
		id: getStringField(item, 'product_id', 'id'),
		name: getStringField(item, 'name', 'product_name'),
		code: getStringField(item, 'code', 'product_code', 'slug'),
		category: getStringField(item, 'category', 'product_category') || 'General',
		status: getStringField(item, 'status') || 'UNKNOWN',
		basePremium: parseMoneyDecimal(item.base_premium),
		minSumInsured: parseMoneyDecimal(item.min_sum_insured),
		maxSumInsured: parseMoneyDecimal(item.max_sum_insured),
		createdAt: getStringField(item, 'created_at'),
		updatedAt: getStringField(item, 'updated_at'),
		description: getStringField(item, 'description', 'short_description')
	};
}

function normalizePolicy(item: JsonMap): SystemPolicy {
	return {
		id: getStringField(item, 'policy_id', 'id'),
		policyNumber: getStringField(item, 'policy_number', 'number'),
		customerName:
			getStringField(item, 'customer_name', 'holder_name', 'customer_id') || 'Unassigned',
		productName:
			getStringField(item, 'product_name', 'plan_name', 'product_id') || 'Unknown product',
		status: getStringField(item, 'status') || 'UNKNOWN',
		premium: parseMoneyDecimal(item.premium_amount),
		sumInsured: parseMoneyDecimal(item.sum_insured),
		startDate: getStringField(item, 'start_date'),
		endDate: getStringField(item, 'end_date')
	};
}

function normalizeClaim(item: JsonMap): SystemClaim {
	return {
		id: getStringField(item, 'claim_id', 'id'),
		claimNumber: getStringField(item, 'claim_number', 'reference_number', 'id'),
		claimantName:
			getStringField(item, 'claimant_name', 'customer_name', 'submitted_by') || 'Unknown',
		policyNumber: getStringField(item, 'policy_number', 'policy_id') || 'Policy pending',
		status: getStringField(item, 'status') || 'UNKNOWN',
		amount: parseMoneyDecimal(item.claim_amount ?? item.amount),
		incidentDate: getStringField(item, 'incident_date', 'loss_date'),
		submittedAt: getStringField(item, 'submitted_at', 'created_at')
	};
}

function normalizePartner(item: JsonMap): SystemPartner {
	const type = getStringField(item, 'type', 'partner_type', 'category') || 'PARTNER';

	return {
		id: getStringField(item, 'partner_id', 'id'),
		name: getStringField(item, 'name', 'legal_name'),
		type,
		category: classifyPartnerCategory(type),
		status: getStringField(item, 'status') || 'UNKNOWN',
		email: getStringField(item, 'email', 'contact_email'),
		phone: getStringField(item, 'phone', 'mobile_number', 'contact_phone'),
		address: getStringField(item, 'address', 'address_line1', 'location'),
		joinedAt: getStringField(item, 'created_at', 'joined_at')
	};
}

function normalizeTenant(item: JsonMap): SystemTenant {
	return {
		id: getStringField(item, 'tenant_id', 'id'),
		name: getStringField(item, 'name', 'display_name'),
		code: getStringField(item, 'code', 'tenant_code'),
		status: getStringField(item, 'status') || 'UNKNOWN',
		domain: getStringField(item, 'domain', 'portal_domain', 'slug'),
		createdAt: getStringField(item, 'created_at')
	};
}

function normalizeReport(item: JsonMap): SystemReport {
	return {
		id: getStringField(item, 'report_definition_id', 'id'),
		name: getStringField(item, 'name', 'title'),
		code: getStringField(item, 'code', 'slug'),
		description: getStringField(item, 'description'),
		status: getStringField(item, 'status') || 'READY'
	};
}

function normalizeAudit(item: JsonMap): SystemAuditEvent {
	return {
		id: getStringField(item, 'audit_log_id', 'event_id', 'id'),
		action: getStringField(item, 'action', 'event_name', 'operation') || 'Unknown action',
		resource: getStringField(item, 'resource', 'entity_type', 'target') || 'System',
		actor: getStringField(item, 'actor', 'user_id', 'performed_by') || 'System',
		status: getStringField(item, 'status', 'outcome') || 'RECORDED',
		timestamp: getStringField(item, 'timestamp', 'created_at')
	};
}

async function attempt<T>(task: () => Promise<T>, fallback: T): Promise<T> {
	try {
		return await task();
	} catch (error) {
		console.error('System portal data fetch failed:', error);
		return fallback;
	}
}

export async function getOverviewData(event: RequestEvent): Promise<SystemOverviewData> {
	const sdk = makeSdkClient(event);

	const [products, partners, tenants, claims, policies] = await Promise.all([
		attempt(async () => {
			const result = await sdk.productServiceListProducts({});
			return pickArray(getRecord(result.data), 'products').map(normalizeProduct);
		}, [] as SystemProduct[]),
		attempt(async () => {
			const result = await sdk.partnerServiceListPartners({});
			return pickArray(getRecord(result.data), 'partners').map(normalizePartner);
		}, [] as SystemPartner[]),
		attempt(async () => {
			const result = await sdk.tenantServiceListTenants({});
			return pickArray(getRecord(result.data), 'tenants').map(normalizeTenant);
		}, [] as SystemTenant[]),
		attempt(async () => {
			const result = await sdk._directHttp.get('/v1/claims?page=1&page_size=8');
			return pickArray(getRecord(result.data), 'claims', 'items').map(normalizeClaim);
		}, [] as SystemClaim[]),
		attempt(async () => {
			const result = await sdk._directHttp.get('/v1/policies?page=1&page_size=8');
			return pickArray(getRecord(result.data), 'policies', 'items').map(normalizePolicy);
		}, [] as SystemPolicy[])
	]);

	return {
		metrics: [
			{
				label: 'Products',
				value: formatCompactNumber(products.length),
				description: 'Live catalog entries from the generated SDK',
				href: '/dashboard/products',
				tone: 'default'
			},
			{
				label: 'Policies',
				value: formatCompactNumber(policies.length),
				description: 'Recent policy records from the policy service',
				href: '/dashboard/policies',
				tone: 'secondary'
			},
			{
				label: 'Claims',
				value: formatCompactNumber(claims.length),
				description: 'Claim cases currently returned by the claims API',
				href: '/dashboard/claims',
				tone: 'outline'
			},
			{
				label: 'Partners',
				value: formatCompactNumber(partners.length),
				description: 'Network partners synced from partner management',
				href: '/dashboard/partners/life',
				tone: 'default'
			},
			{
				label: 'Tenants',
				value: formatCompactNumber(tenants.length),
				description: 'Active tenancy records visible to the system portal',
				href: '/dashboard/tenants',
				tone: 'secondary'
			},
			{
				label: 'Gross Premium',
				value: formatCurrency(
					policies.reduce((total, policy) => total + policy.premium, 0)
				),
				description: 'Total of visible policy premiums in this view',
				href: '/dashboard/policies',
				tone: 'outline'
			}
		],
		products,
		policies,
		claims,
		partners,
		tenants
	};
}

export async function getProducts(event: RequestEvent) {
	const sdk = makeSdkClient(event);
	return attempt(async () => {
		const result = await sdk.productServiceListProducts({});
		return pickArray(getRecord(result.data), 'products').map(normalizeProduct);
	}, [] as SystemProduct[]);
}

export async function getProductDetail(event: RequestEvent, productId: string) {
	const sdk = makeSdkClient(event);

	return attempt(async () => {
		const result = await sdk.productServiceGetProduct({ path: { product_id: productId } });
		return normalizeProduct(getRecord(result.data));
	}, null as SystemProduct | null);
}

export async function getPolicies(event: RequestEvent) {
	const sdk = makeSdkClient(event);
	return attempt(async () => {
		const result = await sdk._directHttp.get('/v1/policies?page=1&page_size=50');
		return pickArray(getRecord(result.data), 'policies', 'items').map(normalizePolicy);
	}, [] as SystemPolicy[]);
}

export async function getClaims(event: RequestEvent) {
	const sdk = makeSdkClient(event);
	return attempt(async () => {
		const result = await sdk._directHttp.get('/v1/claims?page=1&page_size=50');
		return pickArray(getRecord(result.data), 'claims', 'items').map(normalizeClaim);
	}, [] as SystemClaim[]);
}

export async function getPartners(
	event: RequestEvent,
	category?: 'life' | 'non-life'
) {
	const sdk = makeSdkClient(event);
	const partners = await attempt(async () => {
		const result = await sdk.partnerServiceListPartners({});
		return pickArray(getRecord(result.data), 'partners').map(normalizePartner);
	}, [] as SystemPartner[]);

	if (!category) return partners;
	return partners.filter((partner) => partner.category === category);
}

export async function getPartnerDetail(event: RequestEvent, partnerId: string) {
	const sdk = makeSdkClient(event);

	return attempt(async () => {
		const result = await sdk.partnerServiceGetPartner({ path: { partner_id: partnerId } });
		return normalizePartner(getRecord(result.data));
	}, null as SystemPartner | null);
}

export async function getTenants(event: RequestEvent) {
	const sdk = makeSdkClient(event);
	return attempt(async () => {
		const result = await sdk.tenantServiceListTenants({});
		return pickArray(getRecord(result.data), 'tenants').map(normalizeTenant);
	}, [] as SystemTenant[]);
}

export async function getReports(event: RequestEvent) {
	const sdk = makeSdkClient(event);
	return attempt(async () => {
		const result = await sdk.reportServiceListReportDefinitions({});
		return pickArray(getRecord(result.data), 'report_definitions').map(normalizeReport);
	}, [] as SystemReport[]);
}

export async function getAuditEvents(event: RequestEvent) {
	const sdk = makeSdkClient(event);
	return attempt(async () => {
		const result = await sdk._directHttp.get('/v1/audit-logs?page=1&page_size=25');
		return pickArray(getRecord(result.data), 'audit_logs', 'events', 'items').map(
			normalizeAudit
		);
	}, [] as SystemAuditEvent[]);
}

export function getPartnerMix(partners: SystemPartner[]) {
	return {
		life: partners.filter((partner) => partner.category === 'life').length,
		nonLife: partners.filter((partner) => partner.category === 'non-life').length,
		other: partners.filter((partner) => partner.category === 'other').length
	};
}

export function getPolicyHealth(policies: SystemPolicy[]) {
	return {
		active: policies.filter((policy) => policy.status.toUpperCase().includes('ACTIVE')).length,
		lapsed: policies.filter((policy) => policy.status.toUpperCase().includes('LAPSED')).length,
		cancelled: policies.filter((policy) => policy.status.toUpperCase().includes('CANCEL')).length
	};
}

export function getClaimExposure(claims: SystemClaim[]) {
	return {
		open: claims.filter((claim) => !claim.status.toUpperCase().includes('SETTLED')).length,
		totalAmount: claims.reduce((sum, claim) => sum + claim.amount, 0),
		latestIncident:
			claims
				.map((claim) => claim.incidentDate)
				.filter(Boolean)
				.sort()
				.at(-1) ?? ''
	};
}

export function getTenantSummary(tenants: SystemTenant[]) {
	return {
		total: tenants.length,
		active: tenants.filter((tenant) => tenant.status.toUpperCase().includes('ACTIVE')).length,
		lastCreated:
			tenants
				.map((tenant) => tenant.createdAt)
				.filter(Boolean)
				.sort()
				.at(-1) ?? ''
	};
}

export function getOperationalBreakdown(data: SystemOverviewData) {
	return {
		productCount: data.products.length,
		policyCount: data.policies.length,
		claimCount: data.claims.length,
		partnerCount: data.partners.length,
		livePremium: data.policies.reduce((sum, policy) => sum + policy.premium, 0),
		claimExposure: data.claims.reduce((sum, claim) => sum + claim.amount, 0),
		averagePolicyPremium:
			data.policies.length > 0
				? data.policies.reduce((sum, policy) => sum + policy.premium, 0) / data.policies.length
				: 0,
		totalPartnersByCategory: getNumberField(
			{
				count: data.partners.length
			},
			'count'
		)
	};
}
