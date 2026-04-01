export type BadgeTone = 'default' | 'secondary' | 'outline' | 'destructive';

export interface SystemMetric {
	label: string;
	value: string;
	description: string;
	href?: string;
	tone?: BadgeTone;
}

export interface SystemProduct {
	id: string;
	name: string;
	code: string;
	category: string;
	status: string;
	basePremium: number;
	minSumInsured: number;
	maxSumInsured: number;
	createdAt: string;
	updatedAt: string;
	description: string;
}

export interface SystemPolicy {
	id: string;
	policyNumber: string;
	customerName: string;
	productName: string;
	status: string;
	premium: number;
	sumInsured: number;
	startDate: string;
	endDate: string;
}

export interface SystemClaim {
	id: string;
	claimNumber: string;
	claimantName: string;
	policyNumber: string;
	status: string;
	amount: number;
	incidentDate: string;
	submittedAt: string;
}

export interface SystemPartner {
	id: string;
	name: string;
	type: string;
	category: 'life' | 'non-life' | 'other';
	status: string;
	email: string;
	phone: string;
	address: string;
	joinedAt: string;
}

export interface SystemTenant {
	id: string;
	name: string;
	code: string;
	status: string;
	domain: string;
	createdAt: string;
}

export interface SystemReport {
	id: string;
	name: string;
	code: string;
	description: string;
	status: string;
}

export interface SystemAuditEvent {
	id: string;
	action: string;
	resource: string;
	actor: string;
	status: string;
	timestamp: string;
}

export interface SystemOverviewData {
	metrics: SystemMetric[];
	products: SystemProduct[];
	policies: SystemPolicy[];
	claims: SystemClaim[];
	partners: SystemPartner[];
	tenants: SystemTenant[];
}
