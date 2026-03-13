// Demo product data using proto-generated types
// This file will be replaced with API calls when backend is ready

import type { Product, Rider } from '$lib/types';
import { ProductCategory, ProductStatus } from '$lib/types';

export const productsDemo: Product[] = [
	<Product>({
		product_id: 'prod_001',
		product_code: 'HLT-001',
		product_name: 'LabAid Health Guard',
		category: ProductCategory.PRODUCT_CATEGORY_HEALTH,
		description: 'Comprehensive health insurance with cashless facility at LabAid network hospitals',
		base_premium: BigInt(500000), // 5000 BDT in paisa
		min_sum_insured: BigInt(10000000), // 100,000 BDT
		max_sum_insured: BigInt(100000000), // 1,000,000 BDT
		min_tenure_months: 12,
		max_tenure_months: 36,
		exclusions: [
			'Cosmetic surgery',
			'Dental treatment (unless accident)',
			'Pre-existing conditions (first 2 years)',
			'Self-inflicted injuries',
			'War and nuclear risks'
		],
		status: ProductStatus.PRODUCT_STATUS_ACTIVE,
		created_at: '2024-01-15T10:00:00Z',
		updated_at: '2024-12-01T14:30:00Z',
		created_by: 'admin_001',
		available_riders: [
			<Rider>{
				rider_id: 'rider_001',
				product_id: 'prod_001',
				rider_name: 'Pre-Post Hospitalization',
				description: 'Pre and post hospitalization expenses coverage',
				premium_amount: BigInt(100000),
				coverage_amount: BigInt(5000000),
				is_mandatory: false,
				created_at: '2024-01-15T10:00:00Z',
				updated_at: '2024-01-15T10:00:00Z'
			},
			<Rider>{
				rider_id: 'rider_002',
				product_id: 'prod_001',
				rider_name: 'Ambulance Service',
				description: 'Emergency ambulance transportation',
				premium_amount: BigInt(50000),
				coverage_amount: BigInt(2500000),
				is_mandatory: false,
				created_at: '2024-01-15T10:00:00Z',
				updated_at: '2024-01-15T10:00:00Z'
			}
		]
	}),
	<Product>({
		product_id: 'prod_002',
		product_code: 'HLT-002',
		product_name: 'Critical Care Shield',
		category: ProductCategory.PRODUCT_CATEGORY_HEALTH,
		description: 'Specialized coverage for critical illnesses with lump sum payout',
		base_premium: BigInt(1500000),
		min_sum_insured: BigInt(50000000),
		max_sum_insured: BigInt(500000000),
		min_tenure_months: 60,
		max_tenure_months: 240,
		exclusions: [
			'Pre-existing critical illnesses',
			'Self-inflicted injuries',
			'Alcohol or drug abuse related',
			'HIV/AIDS',
			'War and terrorism'
		],
		status: ProductStatus.PRODUCT_STATUS_ACTIVE,
		created_at: '2024-02-20T10:00:00Z',
		updated_at: '2024-11-15T16:20:00Z',
		created_by: 'admin_001'
	}),
	<Product>({
		product_id: 'prod_003',
		product_code: 'LIF-001',
		product_name: 'Life Protection Plus',
		category: ProductCategory.PRODUCT_CATEGORY_LIFE,
		description: 'Term life insurance with comprehensive coverage for your family',
		base_premium: BigInt(1200000),
		min_sum_insured: BigInt(100000000),
		max_sum_insured: BigInt(5000000000),
		min_tenure_months: 60,
		max_tenure_months: 360,
		exclusions: [
			'Suicide within first year',
			'Death due to illegal activities',
			'War and terrorism',
			'Self-inflicted injuries',
			'Death while under influence'
		],
		status: ProductStatus.PRODUCT_STATUS_ACTIVE,
		created_at: '2024-01-10T10:00:00Z',
		updated_at: '2024-12-10T10:15:00Z',
		created_by: 'admin_001',
		available_riders: [
			<Rider>{
				rider_id: 'rider_003',
				product_id: 'prod_003',
				rider_name: 'Accidental Death Benefit',
				description: 'Additional 100% sum assured on accidental death',
				premium_amount: BigInt(300000),
				coverage_amount: BigInt(10000000),
				is_mandatory: false,
				created_at: '2024-01-10T10:00:00Z',
				updated_at: '2024-01-10T10:00:00Z'
			}
		]
	}),
	<Product>({
		product_id: 'prod_004',
		product_code: 'MOT-001',
		product_name: 'Motor Comprehensive',
		category: ProductCategory.PRODUCT_CATEGORY_MOTOR,
		description: 'Comprehensive motor insurance covering own damage and third party liability',
		base_premium: BigInt(800000),
		min_sum_insured: BigInt(50000000),
		max_sum_insured: BigInt(1000000000),
		min_tenure_months: 12,
		max_tenure_months: 12,
		exclusions: [
			'Normal wear and tear',
			'Mechanical/electrical breakdown',
			'Driving without valid license',
			'Driving under influence',
			'War and terrorism'
		],
		status: ProductStatus.PRODUCT_STATUS_ACTIVE,
		created_at: '2024-03-01T10:00:00Z',
		updated_at: '2024-12-15T11:30:00Z',
		created_by: 'admin_001'
	}),
	<Product>({
		product_id: 'prod_005',
		product_code: 'TRV-001',
		product_name: 'Travel Shield International',
		category: ProductCategory.PRODUCT_CATEGORY_TRAVEL,
		description: 'Complete travel insurance for international trips',
		base_premium: BigInt(300000),
		min_sum_insured: BigInt(5000000),
		max_sum_insured: BigInt(500000000),
		min_tenure_months: 1,
		max_tenure_months: 12,
		exclusions: [
			'Pre-existing medical conditions',
			'Adventure sports without add-on',
			'War zones',
			'Pregnancy related (after 24 weeks)',
			'Alcohol/drug related incidents'
		],
		status: ProductStatus.PRODUCT_STATUS_ACTIVE,
		created_at: '2024-04-01T10:00:00Z',
		updated_at: '2024-12-01T09:45:00Z',
		created_by: 'admin_001',
		available_riders: [
			<Rider>{
				rider_id: 'rider_004',
				product_id: 'prod_005',
				rider_name: 'Adventure Sports Coverage',
				description: 'Coverage for adventure sports activities',
				premium_amount: BigInt(150000),
				coverage_amount: BigInt(10000000),
				is_mandatory: false,
				created_at: '2024-04-01T10:00:00Z',
				updated_at: '2024-04-01T10:00:00Z'
			}
		]
	}),
	<Product>({
		product_id: 'prod_006',
		product_code: 'DEV-001',
		product_name: 'Device Protection Plan',
		category: ProductCategory.PRODUCT_CATEGORY_DEVICE,
		description: 'Insurance coverage for mobile phones, laptops, and tablets',
		base_premium: BigInt(150000),
		min_sum_insured: BigInt(1000000),
		max_sum_insured: BigInt(20000000),
		min_tenure_months: 12,
		max_tenure_months: 24,
		exclusions: [
			'Wear and tear',
			'Cosmetic damage',
			'Software issues',
			'Lost or stolen without police report',
			'Damage during repair'
		],
		status: ProductStatus.PRODUCT_STATUS_ACTIVE,
		created_at: '2024-05-15T10:00:00Z',
		updated_at: '2024-11-20T13:15:00Z',
		created_by: 'admin_001'
	}),
	<Product>({
		product_id: 'prod_007',
		product_code: 'HOM-001',
		product_name: 'Home Shield',
		category: ProductCategory.PRODUCT_CATEGORY_HOME,
		description: 'Comprehensive home insurance covering structure and contents',
		base_premium: BigInt(1000000),
		min_sum_insured: BigInt(500000000),
		max_sum_insured: BigInt(10000000000),
		min_tenure_months: 12,
		max_tenure_months: 60,
		exclusions: [
			'Normal wear and tear',
			'Gradual deterioration',
			'War and terrorism',
			'Nuclear risks',
			'Intentional damage'
		],
		status: ProductStatus.PRODUCT_STATUS_ACTIVE,
		created_at: '2024-06-01T10:00:00Z',
		updated_at: '2024-12-05T15:45:00Z',
		created_by: 'admin_001'
	}),
	<Product>({
		product_id: 'prod_008',
		product_code: 'AGR-001',
		product_name: 'Crop Protection Plan',
		category: ProductCategory.PRODUCT_CATEGORY_AGRICULTURAL,
		description: 'Agricultural insurance for crop loss due to natural calamities',
		base_premium: BigInt(500000),
		min_sum_insured: BigInt(5000000),
		max_sum_insured: BigInt(500000000),
		min_tenure_months: 6,
		max_tenure_months: 12,
		exclusions: [
			'Poor farming practices',
			'Failure to report within 48 hours',
			'War and civil unrest',
			'Lack of proper documentation'
		],
		status: ProductStatus.PRODUCT_STATUS_ACTIVE,
		created_at: '2024-07-01T10:00:00Z',
		updated_at: '2024-11-30T10:20:00Z',
		created_by: 'admin_001'
	})
];

// Utility functions for API-ready architecture
export function getProductById(id: string): Product | undefined {
	// TODO: Replace with API call: GET /api/products/{id}
	return productsDemo.find((p) => p.product_id === id);
}

export function getProductsByCategory(category: ProductCategory): Product[] {
	// TODO: Replace with API call: GET /api/products?category={category}
	return productsDemo.filter((p) => p.category === category && p.status === ProductStatus.PRODUCT_STATUS_ACTIVE);
}

export function getProductsByStatus(status: ProductStatus): Product[] {
	// TODO: Replace with API call: GET /api/products?status={status}
	return productsDemo.filter((p) => p.status === status);
}

export function getActiveProducts(): Product[] {
	// TODO: Replace with API call: GET /api/products?status=ACTIVE
	return productsDemo.filter((p) => p.status === ProductStatus.PRODUCT_STATUS_ACTIVE);
}

export function searchProducts(query: string): Product[] {
	// TODO: Replace with API call: GET /api/products/search?q={query}
	const lowerQuery = query.toLowerCase();
	return productsDemo.filter(
		(p) =>
			p.product_name?.toLowerCase().includes(lowerQuery) ||
			p.product_code?.toLowerCase().includes(lowerQuery) ||
			p.description?.toLowerCase().includes(lowerQuery)
	);
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
