export interface PolicyMetrics {
	totalPolicies: number;
	activePolicies: number;
	pendingPolicies: number;
	expiredPolicies: number;
	lifePolicies: number;
	nonLifePolicies: number;
	monthlyGrowth: number;
}

export interface ClaimsMetrics {
	totalClaims: number;
	approvedClaims: number;
	pendingClaims: number;
	rejectedClaims: number;
	cashlessClaims: number;
	reimbursementClaims: number;
	averageClaimAmount: number;
	averageProcessingTime: number; // in days
}

export interface RevenueMetrics {
	totalRevenue: number;
	lifePremiums: number;
	nonLifePremiums: number;
	monthlyRecurring: number;
	discountsGiven: number;
	netRevenue: number;
	growthRate: number;
}

export interface PartnerMetrics {
	totalPartners: number;
	activePartners: number;
	lifePartners: number;
	nonLifePartners: number;
	cashlessEnabled: number;
	discountEnabled: number;
	averageDiscount: number;
	partnerSatisfaction: number;
}

export interface CustomerMetrics {
	totalCustomers: number;
	activeCustomers: number;
	newCustomersThisMonth: number;
	retentionRate: number;
	averageAge: number;
	customerSatisfaction: number;
}

// Monthly trend data for charts
export interface MonthlyData {
	month: string;
	policies: number;
	claims: number;
	revenue: number;
	customers: number;
}

// Analytics Dashboard Data
export const policyMetrics: PolicyMetrics = {
	totalPolicies: 2543,
	activePolicies: 2156,
	pendingPolicies: 234,
	expiredPolicies: 153,
	lifePolicies: 1689,
	nonLifePolicies: 854,
	monthlyGrowth: 12.3
};

export const claimsMetrics: ClaimsMetrics = {
	totalClaims: 1847,
	approvedClaims: 1456,
	pendingClaims: 287,
	rejectedClaims: 104,
	cashlessClaims: 1234,
	reimbursementClaims: 613,
	averageClaimAmount: 45000, // BDT
	averageProcessingTime: 3.5 // days
};

export const revenueMetrics: RevenueMetrics = {
	totalRevenue: 452000000, // 45.2M BDT in paisa
	lifePremiums: 298000000, // 29.8M BDT
	nonLifePremiums: 154000000, // 15.4M BDT
	monthlyRecurring: 37600000, // 3.76M BDT
	discountsGiven: 22800000, // 2.28M BDT
	netRevenue: 429200000, // 42.92M BDT
	growthRate: 23.5
};

export const partnerMetrics: PartnerMetrics = {
	totalPartners: 436,
	activePartners: 412,
	lifePartners: 281,
	nonLifePartners: 155,
	cashlessEnabled: 347,
	discountEnabled: 398,
	averageDiscount: 14.2,
	partnerSatisfaction: 4.6
};

export const customerMetrics: CustomerMetrics = {
	totalCustomers: 18945,
	activeCustomers: 16782,
	newCustomersThisMonth: 1456,
	retentionRate: 88.6,
	averageAge: 38,
	customerSatisfaction: 4.5
};

// Monthly trend data for last 12 months
export const monthlyTrends: MonthlyData[] = [
	{ month: 'Jan', policies: 1980, claims: 1420, revenue: 38500000, customers: 16234 },
	{ month: 'Feb', policies: 2045, claims: 1498, revenue: 39800000, customers: 16598 },
	{ month: 'Mar', policies: 2123, claims: 1556, revenue: 41200000, customers: 16892 },
	{ month: 'Apr', policies: 2198, claims: 1612, revenue: 42100000, customers: 17145 },
	{ month: 'May', policies: 2267, claims: 1678, revenue: 43400000, customers: 17423 },
	{ month: 'Jun', policies: 2334, claims: 1734, revenue: 44200000, customers: 17689 },
	{ month: 'Jul', policies: 2389, claims: 1789, revenue: 44800000, customers: 17934 },
	{ month: 'Aug', policies: 2445, claims: 1823, revenue: 45600000, customers: 18156 },
	{ month: 'Sep', policies: 2478, claims: 1847, revenue: 46100000, customers: 18367 },
	{ month: 'Oct', policies: 2512, claims: 1812, revenue: 46800000, customers: 18589 },
	{ month: 'Nov', policies: 2543, claims: 1847, revenue: 47200000, customers: 18734 },
	{ month: 'Dec', policies: 2543, claims: 1847, revenue: 45200000, customers: 18945 }
];

// Top performing partners
export interface TopPartner {
	id: string;
	name: string;
	type: string;
	claimsProcessed: number;
	revenue: number;
	rating: number;
	cashlessRate: number;
}

export const topPartners: TopPartner[] = [
	{ id: 'H001', name: 'Square Hospital Ltd.', type: 'Hospital', claimsProcessed: 487, revenue: 24500000, rating: 4.8, cashlessRate: 98 },
	{ id: 'H002', name: 'United Hospital', type: 'Hospital', claimsProcessed: 523, revenue: 31200000, rating: 4.9, cashlessRate: 99 },
	{ id: 'P001', name: 'Lazz Pharma', type: 'Pharmacy', claimsProcessed: 892, revenue: 8900000, rating: 4.5, cashlessRate: 100 },
	{ id: 'AR001', name: 'Auto Excellence', type: 'Auto Repair', claimsProcessed: 234, revenue: 11700000, rating: 4.5, cashlessRate: 95 },
	{ id: 'H003', name: 'Apollo Hospital', type: 'Hospital', claimsProcessed: 412, revenue: 22300000, rating: 4.7, cashlessRate: 97 }
];

// Claim distribution by type
export interface ClaimDistribution {
	type: string;
	count: number;
	amount: number;
	percentage: number;
}

export const claimDistribution: ClaimDistribution[] = [
	{ type: 'Hospitalization', count: 687, amount: 412000000, percentage: 37.2 },
	{ type: 'Medication', count: 523, amount: 89000000, percentage: 28.3 },
	{ type: 'Consultation', count: 412, amount: 45000000, percentage: 22.3 },
	{ type: 'Auto Repair', count: 145, amount: 78000000, percentage: 7.8 },
	{ type: 'Device Repair', count: 80, amount: 18000000, percentage: 4.4 }
];

// Policy distribution by age group
export interface PolicyByAge {
	ageGroup: string;
	count: number;
	percentage: number;
}

export const policyByAge: PolicyByAge[] = [
	{ ageGroup: '18-25', count: 234, percentage: 9.2 },
	{ ageGroup: '26-35', count: 789, percentage: 31.0 },
	{ ageGroup: '36-45', count: 892, percentage: 35.1 },
	{ ageGroup: '46-55', count: 456, percentage: 17.9 },
	{ ageGroup: '56-65', count: 172, percentage: 6.8 }
];

// Discount impact analysis
export interface DiscountImpact {
	partnerType: string;
	averageDiscount: number;
	claimsWithDiscount: number;
	totalSavings: number;
	customerSatisfaction: number;
}

export const discountImpact: DiscountImpact[] = [
	{ partnerType: 'Hospitals', averageDiscount: 15.2, claimsWithDiscount: 645, totalSavings: 9780000, customerSatisfaction: 4.7 },
	{ partnerType: 'Pharmacies', averageDiscount: 7.5, claimsWithDiscount: 892, totalSavings: 6690000, customerSatisfaction: 4.5 },
	{ partnerType: 'Doctors', averageDiscount: 12.8, claimsWithDiscount: 378, totalSavings: 4838400, customerSatisfaction: 4.6 },
	{ partnerType: 'Auto Repair', averageDiscount: 18.5, claimsWithDiscount: 134, totalSavings: 2478000, customerSatisfaction: 4.4 },
	{ partnerType: 'Device Repair', averageDiscount: 13.2, claimsWithDiscount: 76, totalSavings: 1003200, customerSatisfaction: 4.3 }
];

// Cashless vs Reimbursement trends
export interface PaymentTrend {
	month: string;
	cashless: number;
	reimbursement: number;
}

export const paymentTrends: PaymentTrend[] = [
	{ month: 'Jan', cashless: 945, reimbursement: 475 },
	{ month: 'Feb', cashless: 998, reimbursement: 500 },
	{ month: 'Mar', cashless: 1034, reimbursement: 522 },
	{ month: 'Apr', cashless: 1071, reimbursement: 541 },
	{ month: 'May', cashless: 1115, reimbursement: 563 },
	{ month: 'Jun', cashless: 1153, reimbursement: 581 },
	{ month: 'Jul', cashless: 1189, reimbursement: 600 },
	{ month: 'Aug', cashless: 1211, reimbursement: 612 },
	{ month: 'Sep', cashless: 1228, reimbursement: 619 },
	{ month: 'Oct', cashless: 1204, reimbursement: 608 },
	{ month: 'Nov', cashless: 1226, reimbursement: 621 },
	{ month: 'Dec', cashless: 1234, reimbursement: 613 }
];

// Utility functions
export function formatCurrency(amount: number): string {
	return new Intl.NumberFormat('en-BD', {
		style: 'currency',
		currency: 'BDT',
		minimumFractionDigits: 0
	}).format(amount / 100);
}

export function formatNumber(num: number): string {
	return new Intl.NumberFormat('en-BD').format(num);
}

export function formatPercentage(num: number): string {
	return `${num.toFixed(1)}%`;
}
