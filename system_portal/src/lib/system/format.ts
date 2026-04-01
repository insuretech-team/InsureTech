import type { BadgeTone } from './types';

export function formatCurrency(amount: number): string {
	return new Intl.NumberFormat('en-BD', {
		style: 'currency',
		currency: 'BDT',
		maximumFractionDigits: 0
	}).format(Number.isFinite(amount) ? amount : 0);
}

export function formatCompactNumber(value: number): string {
	return new Intl.NumberFormat('en', {
		notation: 'compact',
		maximumFractionDigits: 1
	}).format(Number.isFinite(value) ? value : 0);
}

export function formatDate(value: string): string {
	if (!value) return 'Not available';

	const date = new Date(value);
	if (Number.isNaN(date.getTime())) return value;

	return new Intl.DateTimeFormat('en-BD', {
		year: 'numeric',
		month: 'short',
		day: 'numeric'
	}).format(date);
}

export function formatRelativeDays(value: string): string {
	if (!value) return 'Unknown';
	const date = new Date(value);
	if (Number.isNaN(date.getTime())) return 'Unknown';

	const days = Math.max(0, Math.floor((Date.now() - date.getTime()) / 86400000));
	return `${days}d ago`;
}

export function humanizeStatus(value: string): string {
	if (!value) return 'Unknown';
	return value
		.replace(/^USER_TYPE_/, '')
		.replace(/^POLICY_STATUS_/, '')
		.replace(/^CLAIM_STATUS_/, '')
		.replace(/^PRODUCT_STATUS_/, '')
		.replace(/^PARTNER_STATUS_/, '')
		.replace(/^REPORT_STATUS_/, '')
		.replace(/_/g, ' ')
		.toLowerCase()
		.replace(/\b\w/g, (char) => char.toUpperCase());
}

export function statusTone(value: string): BadgeTone {
	const normalized = value.toUpperCase();

	if (
		normalized.includes('ACTIVE') ||
		normalized.includes('APPROVED') ||
		normalized.includes('COMPLETED') ||
		normalized.includes('SUCCESS')
	) {
		return 'default';
	}

	if (
		normalized.includes('PENDING') ||
		normalized.includes('REVIEW') ||
		normalized.includes('DRAFT')
	) {
		return 'secondary';
	}

	if (
		normalized.includes('REJECT') ||
		normalized.includes('FAILED') ||
		normalized.includes('CANCEL') ||
		normalized.includes('LAPSED') ||
		normalized.includes('SUSPEND')
	) {
		return 'destructive';
	}

	return 'outline';
}
