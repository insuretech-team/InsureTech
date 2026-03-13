// Demo claims data using proto-generated types
// This file will be replaced with API calls when backend is ready

import {  Claim, ClaimStatus, ClaimType, ClaimDocument, ClaimApproval, FraudCheckResult, ApprovalDecision  } from '$lib/types';

export const claimsDemo: Claim[] = [
	<Claim> ({
		claim_id: 'clm_001',
		claim_number: 'CLM-2024-HLTH-000001',
		policy_id: 'pol_001',
		customer_id: 'cust_001',
		status: ClaimStatus.CLAIM_STATUS_APPROVED,
		type: ClaimType.CLAIM_TYPE_HEALTH_HOSPITALIZATION,
		claimed_amount: BigInt(12000000), // 120,000 BDT
		approved_amount: BigInt(11500000), // 115,000 BDT
		settled_amount: BigInt(11500000),
		incident_date: '2024-11-15T00:00:00Z',
		incident_description: 'Emergency appendectomy surgery at LabAid Cardiac Hospital',
		submitted_at: '2024-11-17T10:30:00Z',
		approved_at: '2024-11-20T15:45:00Z',
		settled_at: '2024-11-22T11:00:00Z',
		created_at: '2024-11-17T10:30:00Z',
		updated_at: '2024-11-22T11:00:00Z',
		documents: [
			<ClaimDocument> ({
				document_id: 'doc_001',
				claim_id: 'clm_001',
				document_type: 'medical_report',
				file_url: '/documents/claims/clm_001/medical_report.pdf',
				file_hash: 'sha256_hash_here',
				uploaded_at: '2024-11-17T10:35:00Z',
				verified: true,
				verified_by: 'admin_001',
				created_at: '2024-11-17T10:35:00Z',
				updated_at: '2024-11-18T09:00:00Z'
			}),
			<ClaimDocument> ({
				document_id: 'doc_002',
				claim_id: 'clm_001',
				document_type: 'bill',
				file_url: '/documents/claims/clm_001/hospital_bill.pdf',
				file_hash: 'sha256_hash_here_2',
				uploaded_at: '2024-11-17T10:36:00Z',
				verified: true,
				verified_by: 'admin_001',
				created_at: '2024-11-17T10:36:00Z',
				updated_at: '2024-11-18T09:00:00Z'
			})
		],
		approvals: [
			<ClaimApproval> ({
				approval_id: 'appr_001',
				claim_id: 'clm_001',
				approver_id: 'admin_002',
				approver_role: 'Claims Officer',
				approval_level: 1,
				decision: ApprovalDecision.APPROVAL_DECISION_APPROVED,
				approved_amount: BigInt(11500000),
				notes: 'All documents verified. Approved within policy limits.',
				approved_at: '2024-11-18T14:30:00Z',
				created_at: '2024-11-18T14:30:00Z'
			}),
			<ClaimApproval> ({
				approval_id: 'appr_002',
				claim_id: 'clm_001',
				approver_id: 'admin_003',
				approver_role: 'Claims Manager',
				approval_level: 2,
				decision: ApprovalDecision.APPROVAL_DECISION_APPROVED,
				approved_amount: BigInt(11500000),
				notes: 'Final approval granted.',
				approved_at: '2024-11-20T15:45:00Z',
				created_at: '2024-11-20T15:45:00Z'
			})
		],
		fraud_check: <FraudCheckResult> ({
			fraud_check_id: 'fraud_001',
			claim_id: 'clm_001',
			fraud_score: 15,
			risk_factors: [],
			flagged: false,
			reviewed_by: 'system',
			reviewed_at: '2024-11-17T10:40:00Z',
			created_at: '2024-11-17T10:40:00Z'
		})
	}),
	<Claim> ({
		claim_id: 'clm_002',
		claim_number: 'CLM-2024-MOTR-000001',
		policy_id: 'pol_003',
		customer_id: 'cust_003',
		status: ClaimStatus.CLAIM_STATUS_UNDER_REVIEW,
		type: ClaimType.CLAIM_TYPE_MOTOR_ACCIDENT,
		claimed_amount: BigInt(8500000), // 85,000 BDT
		incident_date: '2024-12-10T08:30:00Z',
		incident_description: 'Vehicle collision with another car at Mohakhali intersection. Front bumper and headlight damaged.',
		submitted_at: '2024-12-12T14:00:00Z',
		created_at: '2024-12-12T14:00:00Z',
		updated_at: '2024-12-18T10:00:00Z',
		documents: [
			<ClaimDocument> ({
				document_id: 'doc_003',
				claim_id: 'clm_002',
				document_type: 'police_report',
				file_url: '/documents/claims/clm_002/police_report.pdf',
				file_hash: 'sha256_hash_here_3',
				uploaded_at: '2024-12-12T14:05:00Z',
				verified: true,
				verified_by: 'admin_001',
				created_at: '2024-12-12T14:05:00Z',
				updated_at: '2024-12-13T09:00:00Z'
			}),
			<ClaimDocument> ({
				document_id: 'doc_004',
				claim_id: 'clm_002',
				document_type: 'photos',
				file_url: '/documents/claims/clm_002/damage_photos.zip',
				file_hash: 'sha256_hash_here_4',
				uploaded_at: '2024-12-12T14:10:00Z',
				verified: true,
				verified_by: 'admin_001',
				created_at: '2024-12-12T14:10:00Z',
				updated_at: '2024-12-13T09:00:00Z'
			})
		],
		approvals: [
			<ClaimApproval> ({
				approval_id: 'appr_003',
				claim_id: 'clm_002',
				approver_id: 'admin_004',
				approver_role: 'Claims Officer',
				approval_level: 1,
				decision: ApprovalDecision.APPROVAL_DECISION_PENDING,
				notes: 'Awaiting repair estimate from garage.',
				created_at: '2024-12-13T10:00:00Z'
			})
		],
		fraud_check: <FraudCheckResult> ({
			fraud_check_id: 'fraud_002',
			claim_id: 'clm_002',
			fraud_score: 25,
			risk_factors: ['Multiple claims in short period'],
			flagged: false,
			reviewed_by: 'system',
			reviewed_at: '2024-12-12T14:15:00Z',
			created_at: '2024-12-12T14:15:00Z'
		})
	}),
	<Claim> ({
		claim_id: 'clm_003',
		claim_number: 'CLM-2024-HLTH-000002',
		policy_id: 'pol_001',
		customer_id: 'cust_001',
		status: ClaimStatus.CLAIM_STATUS_PENDING_DOCUMENTS,
		type: ClaimType.CLAIM_TYPE_HEALTH_SURGERY,
		claimed_amount: BigInt(25000000), // 250,000 BDT
		incident_date: '2024-12-15T00:00:00Z',
		incident_description: 'Scheduled knee replacement surgery',
		submitted_at: '2024-12-20T09:00:00Z',
		created_at: '2024-12-20T09:00:00Z',
		updated_at: '2024-12-22T16:00:00Z',
		documents: [
			<ClaimDocument> ({
				document_id: 'doc_005',
				claim_id: 'clm_003',
				document_type: 'prescription',
				file_url: '/documents/claims/clm_003/prescription.pdf',
				file_hash: 'sha256_hash_here_5',
				uploaded_at: '2024-12-20T09:05:00Z',
				verified: false,
				created_at: '2024-12-20T09:05:00Z',
				updated_at: '2024-12-20T09:05:00Z'
			})
		],
		approvals: [],
		fraud_check: <FraudCheckResult> ({
			fraud_check_id: 'fraud_003',
			claim_id: 'clm_003',
			fraud_score: 10,
			risk_factors: [],
			flagged: false,
			reviewed_by: 'system',
			reviewed_at: '2024-12-20T09:10:00Z',
			created_at: '2024-12-20T09:10:00Z'
		})
	}),
	<Claim> ({
		claim_id: 'clm_004',
		claim_number: 'CLM-2024-MOTR-000002',
		policy_id: 'pol_003',
		customer_id: 'cust_003',
		status: ClaimStatus.CLAIM_STATUS_REJECTED,
		type: ClaimType.CLAIM_TYPE_MOTOR_THEFT,
		claimed_amount: BigInt(150000000), // 1,500,000 BDT
		incident_date: '2024-10-20T02:00:00Z',
		incident_description: 'Vehicle reported stolen from parking lot',
		submitted_at: '2024-10-22T11:00:00Z',
		created_at: '2024-10-22T11:00:00Z',
		updated_at: '2024-11-05T14:30:00Z',
		rejection_reason: 'Police investigation revealed inconsistencies in theft report. Vehicle tracking data shows movement after reported theft time.',
		documents: [
			<ClaimDocument> ({
				document_id: 'doc_006',
				claim_id: 'clm_004',
				document_type: 'police_report',
				file_url: '/documents/claims/clm_004/police_report.pdf',
				file_hash: 'sha256_hash_here_6',
				uploaded_at: '2024-10-22T11:10:00Z',
				verified: true,
				verified_by: 'admin_001',
				created_at: '2024-10-22T11:10:00Z',
				updated_at: '2024-10-23T09:00:00Z'
			})
		],
		approvals: [
			<ClaimApproval> ({
				approval_id: 'appr_004',
				claim_id: 'clm_004',
				approver_id: 'admin_005',
				approver_role: 'Claims Manager',
				approval_level: 2,
				decision: ApprovalDecision.APPROVAL_DECISION_REJECTED,
				notes: 'Claim rejected due to fraudulent activity detected.',
				created_at: '2024-11-05T14:30:00Z'
			})
		],
		fraud_check: <FraudCheckResult> ({
			fraud_check_id: 'fraud_004',
			claim_id: 'clm_004',
			fraud_score: 85,
			risk_factors: ['GPS data inconsistency', 'Multiple similar claims', 'Delayed reporting'],
			flagged: true,
			reviewed_by: 'admin_006',
			reviewed_at: '2024-10-25T10:00:00Z',
			created_at: '2024-10-22T11:15:00Z'
		})
	}),
	<Claim> ({
		claim_id: 'clm_005',
		claim_number: 'CLM-2024-TRVL-000001',
		policy_id: 'pol_005',
		customer_id: 'cust_005',
		status: ClaimStatus.CLAIM_STATUS_SETTLED,
		type: ClaimType.CLAIM_TYPE_TRAVEL_BAGGAGE_LOSS,
		claimed_amount: BigInt(5000000), // 50,000 BDT
		approved_amount: BigInt(5000000),
		settled_amount: BigInt(5000000),
		incident_date: '2024-05-15T00:00:00Z',
		incident_description: 'Checked baggage lost during international flight',
		submitted_at: '2024-05-17T10:00:00Z',
		approved_at: '2024-05-22T14:00:00Z',
		settled_at: '2024-05-25T10:00:00Z',
		created_at: '2024-05-17T10:00:00Z',
		updated_at: '2024-05-25T10:00:00Z',
		documents: [
			<ClaimDocument> ({
				document_id: 'doc_007',
				claim_id: 'clm_005',
				document_type: 'airline_report',
				file_url: '/documents/claims/clm_005/pir_report.pdf',
				file_hash: 'sha256_hash_here_7',
				uploaded_at: '2024-05-17T10:10:00Z',
				verified: true,
				verified_by: 'admin_001',
				created_at: '2024-05-17T10:10:00Z',
				updated_at: '2024-05-18T09:00:00Z'
			})
		],
		approvals: [
			<ClaimApproval> ({
				approval_id: 'appr_005',
				claim_id: 'clm_005',
				approver_id: 'admin_007',
				approver_role: 'Claims Officer',
				approval_level: 1,
				decision: ApprovalDecision.APPROVAL_DECISION_APPROVED,
				approved_amount: BigInt(5000000),
				notes: 'PIR report and receipts verified.',
				approved_at: '2024-05-22T14:00:00Z',
				created_at: '2024-05-22T14:00:00Z'
			})
		],
		fraud_check: <FraudCheckResult> ({
			fraud_check_id: 'fraud_005',
			claim_id: 'clm_005',
			fraud_score: 5,
			risk_factors: [],
			flagged: false,
			reviewed_by: 'system',
			reviewed_at: '2024-05-17T10:15:00Z',
			created_at: '2024-05-17T10:15:00Z'
		})
	})
];

// Utility functions for API-ready architecture
export function getClaimById(id: string): Claim | undefined {
	// TODO: Replace with API call: GET /api/claims/{id}
	return claimsDemo.find((c) => c.claim_id === id);
}

export function getClaimsByCustomer(customer_id: string): Claim[] {
	// TODO: Replace with API call: GET /api/claims?customer_id={customer_id}
	return claimsDemo.filter((c) => c.customer_id === customer_id);
}

export function getClaimsByStatus(status: ClaimStatus): Claim[] {
	// TODO: Replace with API call: GET /api/claims?status={status}
	return claimsDemo.filter((c) => c.status === status);
}

export function getClaimsByPolicy(policy_id: string): Claim[] {
	// TODO: Replace with API call: GET /api/claims?policy_id={policy_id}
	return claimsDemo.filter((c) => c.policy_id === policy_id);
}

export function searchClaims(query: string): Claim[] {
	// TODO: Replace with API call: GET /api/claims/search?q={query}
	const lowerQuery = query.toLowerCase();
	return claimsDemo.filter((c) => c.claim_number?.toLowerCase().includes(lowerQuery));
}

export function getPendingClaims(): Claim[] {
	// TODO: Replace with API call: GET /api/claims?status=UNDER_REVIEW,PENDING_DOCUMENTS
	return claimsDemo.filter(
		(c) =>
			c.status === ClaimStatus.CLAIM_STATUS_UNDER_REVIEW ||
			c.status === ClaimStatus.CLAIM_STATUS_PENDING_DOCUMENTS ||
			c.status === ClaimStatus.CLAIM_STATUS_SUBMITTED
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

// Helper to get status badge color
export function getClaimStatusColor(status: ClaimStatus): string {
	switch (status) {
		case ClaimStatus.CLAIM_STATUS_APPROVED:
		case ClaimStatus.CLAIM_STATUS_SETTLED:
			return 'success';
		case ClaimStatus.CLAIM_STATUS_UNDER_REVIEW:
		case ClaimStatus.CLAIM_STATUS_PENDING_DOCUMENTS:
			return 'warning';
		case ClaimStatus.CLAIM_STATUS_SUBMITTED:
			return 'info';
		case ClaimStatus.CLAIM_STATUS_REJECTED:
		case ClaimStatus.CLAIM_STATUS_DISPUTED:
			return 'destructive';
		default:
			return 'secondary';
	}
}
