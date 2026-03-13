// Demo policy data using proto-generated types
// This file will be replaced with API calls when backend is ready

import type { Policy, Nominee, PolicyRider } from '$lib/types';
import { PolicyStatus } from '$lib/types';

export const policiesDemo: Policy[] = [
	<Policy>({
		policy_id: 'pol_001',
		policyNumber: 'LBT-2024-HLTH-000001',
		product_id: 'prod_001',
		customer_id: 'cust_001',
		partnerId: 'partner_001',
		status: PolicyStatus.POLICY_STATUS_ACTIVE,
		premium_amount: BigInt(650000),
		sumInsured: BigInt(50000000),
		tenureMonths: 12,
		startDate: '2024-01-15T00:00:00Z',
		endDate: '2025-01-14T23:59:59Z',
		issuedAt: '2024-01-15T10:30:00Z',
		created_at: '2024-01-15T09:00:00Z',
		updated_at: '2024-01-15T10:30:00Z',
		policyDocumentUrl: '/documents/policies/LBT-2024-HLTH-000001.pdf',
		nominees: [
			<Nominee>{
				nominee_id: 'nom_001',
				policy_id: 'pol_001',
				nominee_name: 'Fatima Rahman',
				nominee_relationship: 'Spouse',
				nominee_share_percent: 60,
				nominee_dob_text: '1992-05-20T00:00:00Z',
				created_at: '2024-01-15T09:15:00Z',
				updated_at: '2024-01-15T09:15:00Z'
			},
			<Nominee>{
				nominee_id: 'nom_002',
				policy_id: 'pol_001',
				nominee_name: 'Ayesha Rahman',
				nominee_relationship: 'Daughter',
				nominee_share_percent: 40,
				nominee_dob_text: '2015-08-10T00:00:00Z',
				created_at: '2024-01-15T09:15:00Z',
				updated_at: '2024-01-15T09:15:00Z'
			}
		],
		riders: [
			<PolicyRider>{
				rider_id: 'pol_rider_001',
				policy_id: 'pol_001',
				rider_name: 'Ambulance Service',
				premium_amount: BigInt(50000),
				coverage_amount: BigInt(2500000),
				created_at: '2024-01-15T09:15:00Z',
				updated_at: '2024-01-15T09:15:00Z'
			}
		]
	}),
	<Policy>({
		policy_id: 'pol_002',
		policyNumber: 'LBT-2024-LIFE-000001',
		product_id: 'prod_003',
		customer_id: 'cust_002',
		status: PolicyStatus.POLICY_STATUS_ACTIVE,
		premium_amount: BigInt(1500000),
		sumInsured: BigInt(200000000),
		tenureMonths: 120,
		startDate: '2024-02-01T00:00:00Z',
		endDate: '2034-01-31T23:59:59Z',
		issuedAt: '2024-02-01T11:00:00Z',
		created_at: '2024-01-28T10:00:00Z',
		updated_at: '2024-02-01T11:00:00Z',
		policyDocumentUrl: '/documents/policies/LBT-2024-LIFE-000001.pdf',
		nominees: [
			<Nominee>{
				nominee_id: 'nom_003',
				policy_id: 'pol_002',
				nominee_name: 'Rahim Ahmed',
				nominee_relationship: 'Son',
				nominee_share_percent: 100,
				nominee_dob_text: '1998-03-15T00:00:00Z',
				created_at: '2024-01-28T10:15:00Z',
				updated_at: '2024-01-28T10:15:00Z'
			}
		]
	}),
	<Policy>({
		policy_id: 'pol_003',
		policyNumber: 'LBT-2024-MOTR-000001',
		product_id: 'prod_004',
		customer_id: 'cust_003',
		status: PolicyStatus.POLICY_STATUS_ACTIVE,
		premium_amount: BigInt(1250000),
		sumInsured: BigInt(150000000),
		tenureMonths: 12,
		startDate: '2024-03-10T00:00:00Z',
		endDate: '2025-03-09T23:59:59Z',
		issuedAt: '2024-03-10T14:00:00Z',
		created_at: '2024-03-08T11:00:00Z',
		updated_at: '2024-03-10T14:00:00Z',
		policyDocumentUrl: '/documents/policies/LBT-2024-MOTR-000001.pdf'
	}),
	<Policy>({
		policy_id: 'pol_004',
		policyNumber: 'LBT-2024-HLTH-000002',
		product_id: 'prod_001',
		customer_id: 'cust_004',
		partnerId: 'partner_001',
		status: PolicyStatus.POLICY_STATUS_GRACE_PERIOD,
		premium_amount: BigInt(700000),
		sumInsured: BigInt(60000000),
		tenureMonths: 12,
		startDate: '2024-04-01T00:00:00Z',
		endDate: '2025-03-31T23:59:59Z',
		issuedAt: '2024-04-01T10:00:00Z',
		created_at: '2024-03-28T09:00:00Z',
		updated_at: '2024-11-15T10:00:00Z',
		policyDocumentUrl: '/documents/policies/LBT-2024-HLTH-000002.pdf'
	}),
	<Policy>({
		policy_id: 'pol_005',
		policyNumber: 'LBT-2024-TRVL-000001',
		product_id: 'prod_005',
		customer_id: 'cust_005',
		status: PolicyStatus.POLICY_STATUS_EXPIRED,
		premium_amount: BigInt(450000),
		sumInsured: BigInt(30000000),
		tenureMonths: 1,
		startDate: '2024-05-01T00:00:00Z',
		endDate: '2024-05-31T23:59:59Z',
		issuedAt: '2024-04-30T16:00:00Z',
		created_at: '2024-04-29T14:00:00Z',
		updated_at: '2024-06-01T00:00:00Z',
		policyDocumentUrl: '/documents/policies/LBT-2024-TRVL-000001.pdf'
	})
];

// Utility functions for API-ready architecture
export function getPolicyById(id: string): Policy | undefined {
	// TODO: Replace with API call: GET /api/policies/{id}
	return policiesDemo.find((p) => p.policy_id === id);
}

export function getPoliciesByCustomer(customer_id: string): Policy[] {
	// TODO: Replace with API call: GET /api/policies?customer_id={customer_id}
	return policiesDemo.filter((p) => p.customer_id === customer_id);
}

export function getPoliciesByStatus(status: PolicyStatus): Policy[] {
	// TODO: Replace with API call: GET /api/policies?status={status}
	return policiesDemo.filter((p) => p.status === status);
}

export function getPoliciesByProduct(product_id: string): Policy[] {
	// TODO: Replace with API call: GET /api/policies?product_id={product_id}
	return policiesDemo.filter((p) => p.product_id === product_id);
}

export function searchPolicies(query: string): Policy[] {
	// TODO: Replace with API call: GET /api/policies/search?q={query}
	const lowerQuery = query.toLowerCase();
	return policiesDemo.filter((p) => p.policyNumber?.toLowerCase().includes(lowerQuery));
}

export function getActivePolicies(): Policy[] {
	// TODO: Replace with API call: GET /api/policies?status=ACTIVE
	return policiesDemo.filter((p) => p.status === PolicyStatus.POLICY_STATUS_ACTIVE);
}

// Helper function to format bigint to BDT
export function formatBDT(paisa: bigint): string {
	const amount = Number(paisa) / 100;
	return new Intl.NumberFormat('en-BD', {
		style: 'currency',
		currency: 'BDT',
		minimumFractionDigits: 0
	}).format(amount);
}

// Helper to get status badge color
export function getStatusColor(status: PolicyStatus): string {
	switch (status) {
		case PolicyStatus.POLICY_STATUS_ACTIVE:
			return 'success';
		case PolicyStatus.POLICY_STATUS_GRACE_PERIOD:
			return 'warning';
		case PolicyStatus.POLICY_STATUS_PENDING_PAYMENT:
			return 'info';
		case PolicyStatus.POLICY_STATUS_EXPIRED:
		case PolicyStatus.POLICY_STATUS_CANCELLED:
		case PolicyStatus.POLICY_STATUS_LAPSED:
			return 'destructive';
		default:
			return 'secondary';
	}
}
