// Clean, well-structured analytics data for InsureTech business

// KPI Metrics
export interface KPICard {
	id: string;
	title: string;
	value: string;
	change: string;
	trend: 'up' | 'down';
	icon: string;
}

export const kpiMetrics: KPICard[] = [
	{
		id: 'revenue',
		title: 'Total Revenue',
		value: '৳45.2M',
		change: '+23.5%',
		trend: 'up',
		icon: 'dollar'
	},
	{
		id: 'policies',
		title: 'Active Policies',
		value: '2,156',
		change: '+12.3%',
		trend: 'up',
		icon: 'file'
	},
	{
		id: 'claims',
		title: 'Total Claims',
		value: '1,847',
		change: '+8.3%',
		trend: 'up',
		icon: 'activity'
	},
	{
		id: 'customers',
		title: 'Active Customers',
		value: '16,782',
		change: '+15.2%',
		trend: 'up',
		icon: 'users'
	}
];

// Chart Data - Monthly Trends
export interface ChartDataPoint {
	month: string;
	policies: number;
	claims: number;
	revenue: number;
}

export const monthlyChartData: ChartDataPoint[] = [
	{ month: 'Jan', policies: 1980, claims: 1420, revenue: 38.5 },
	{ month: 'Feb', policies: 2045, claims: 1498, revenue: 39.8 },
	{ month: 'Mar', policies: 2123, claims: 1556, revenue: 41.2 },
	{ month: 'Apr', policies: 2198, claims: 1612, revenue: 42.1 },
	{ month: 'May', policies: 2267, claims: 1678, revenue: 43.4 },
	{ month: 'Jun', policies: 2334, claims: 1734, revenue: 44.2 },
	{ month: 'Jul', policies: 2389, claims: 1789, revenue: 44.8 },
	{ month: 'Aug', policies: 2445, claims: 1823, revenue: 45.6 },
	{ month: 'Sep', policies: 2478, claims: 1847, revenue: 46.1 },
	{ month: 'Oct', policies: 2512, claims: 1812, revenue: 46.8 },
	{ month: 'Nov', policies: 2543, claims: 1847, revenue: 47.2 },
	{ month: 'Dec', policies: 2543, claims: 1847, revenue: 45.2 }
];

// Policy Analytics
export interface PolicyStats {
	total: number;
	active: number;
	pending: number;
	expired: number;
	life: number;
	nonLife: number;
	growthRate: number;
}

export const policyStats: PolicyStats = {
	total: 2543,
	active: 2156,
	pending: 234,
	expired: 153,
	life: 1689,
	nonLife: 854,
	growthRate: 12.3
};

// Claims Analytics
export interface ClaimStats {
	total: number;
	approved: number;
	pending: number;
	rejected: number;
	cashless: number;
	reimbursement: number;
	avgAmount: number;
	avgProcessingDays: number;
}

export const claimStats: ClaimStats = {
	total: 1847,
	approved: 1456,
	pending: 287,
	rejected: 104,
	cashless: 1234,
	reimbursement: 613,
	avgAmount: 45000,
	avgProcessingDays: 3.5
};

// Claim Distribution Table Data
export interface ClaimDistribution {
	type: string;
	count: number;
	amount: number;
	percentage: number;
}

export const claimDistributionData: ClaimDistribution[] = [
	{ type: 'Hospitalization', count: 687, amount: 412.0, percentage: 37.2 },
	{ type: 'Medication', count: 523, amount: 89.0, percentage: 28.3 },
	{ type: 'Consultation', count: 412, amount: 45.0, percentage: 22.3 },
	{ type: 'Auto Repair', count: 145, amount: 78.0, percentage: 7.8 },
	{ type: 'Device Repair', count: 80, amount: 18.0, percentage: 4.4 }
];

// Revenue Analytics
export interface RevenueStats {
	total: number;
	life: number;
	nonLife: number;
	recurring: number;
	discounts: number;
	net: number;
	growthRate: number;
}

export const revenueStats: RevenueStats = {
	total: 45.2,
	life: 29.8,
	nonLife: 15.4,
	recurring: 3.76,
	discounts: 2.28,
	net: 42.92,
	growthRate: 23.5
};

// Partner Performance Table Data
export interface PartnerPerformance {
	id: string;
	name: string;
	type: string;
	claims: number;
	revenue: number;
	rating: number;
	cashlessRate: number;
}

export const topPartnersData: PartnerPerformance[] = [
	{ id: 'H002', name: 'United Hospital', type: 'Hospital', claims: 523, revenue: 31.2, rating: 4.9, cashlessRate: 99 },
	{ id: 'H001', name: 'Square Hospital Ltd.', type: 'Hospital', claims: 487, revenue: 24.5, rating: 4.8, cashlessRate: 98 },
	{ id: 'H003', name: 'Apollo Hospital', type: 'Hospital', claims: 412, revenue: 22.3, rating: 4.7, cashlessRate: 97 },
	{ id: 'AR001', name: 'Auto Excellence', type: 'Auto Repair', claims: 234, revenue: 11.7, rating: 4.5, cashlessRate: 95 },
	{ id: 'P001', name: 'Lazz Pharma', type: 'Pharmacy', claims: 892, revenue: 8.9, rating: 4.5, cashlessRate: 100 }
];

// Partner Stats
export interface PartnerStats {
	total: number;
	active: number;
	life: number;
	nonLife: number;
	cashlessEnabled: number;
	discountEnabled: number;
	avgDiscount: number;
	avgSatisfaction: number;
}

export const partnerStats: PartnerStats = {
	total: 436,
	active: 412,
	life: 281,
	nonLife: 155,
	cashlessEnabled: 347,
	discountEnabled: 398,
	avgDiscount: 14.2,
	avgSatisfaction: 4.6
};

// Discount Impact Table Data
export interface DiscountImpact {
	partnerType: string;
	avgDiscount: number;
	claimsCount: number;
	totalSavings: number;
	satisfaction: number;
}

export const discountImpactData: DiscountImpact[] = [
	{ partnerType: 'Hospitals', avgDiscount: 15.2, claimsCount: 645, totalSavings: 9.78, satisfaction: 4.7 },
	{ partnerType: 'Pharmacies', avgDiscount: 7.5, claimsCount: 892, totalSavings: 6.69, satisfaction: 4.5 },
	{ partnerType: 'Doctors', avgDiscount: 12.8, claimsCount: 378, totalSavings: 4.84, satisfaction: 4.6 },
	{ partnerType: 'Auto Repair', avgDiscount: 18.5, claimsCount: 134, totalSavings: 2.48, satisfaction: 4.4 },
	{ partnerType: 'Device Repair', avgDiscount: 13.2, claimsCount: 76, totalSavings: 1.00, satisfaction: 4.3 }
];

// Age Group Distribution
export interface AgeDistribution {
	ageGroup: string;
	count: number;
	percentage: number;
}

export const ageDistributionData: AgeDistribution[] = [
	{ ageGroup: '18-25', count: 234, percentage: 9.2 },
	{ ageGroup: '26-35', count: 789, percentage: 31.0 },
	{ ageGroup: '36-45', count: 892, percentage: 35.1 },
	{ ageGroup: '46-55', count: 456, percentage: 17.9 },
	{ ageGroup: '56-65', count: 172, percentage: 6.8 }
];

// Customer Stats
export interface CustomerStats {
	total: number;
	active: number;
	newThisMonth: number;
	retentionRate: number;
	avgAge: number;
	satisfaction: number;
}

export const customerStats: CustomerStats = {
	total: 18945,
	active: 16782,
	newThisMonth: 1456,
	retentionRate: 88.6,
	avgAge: 38,
	satisfaction: 4.5
};

// Utility Functions
export function formatCurrency(amount: number): string {
	if (amount >= 1000000) {
		return `৳${(amount / 1000000).toFixed(1)}M`;
	}
	return `৳${amount.toFixed(2)}M`;
}

export function formatNumber(num: number): string {
	return new Intl.NumberFormat('en-BD').format(num);
}

export function formatPercent(num: number): string {
	return `${num.toFixed(1)}%`;
}
