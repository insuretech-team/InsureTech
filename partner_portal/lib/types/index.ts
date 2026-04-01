// Core entity types for Partner Portal

export type PartnerType =
  | 'HOSPITAL'
  | 'CLINIC'
  | 'PHARMACY'
  | 'AUTO_REPAIR'
  | 'PET_CLINIC'
  | 'FIRE_INSPECTOR'
  | 'LAPTOP_REPAIR'
  | 'MOBILE_REPAIR'
  | 'AMBULANCE'
  | 'MFS'
  | 'ECOMMERCE'
  | 'AGENT_NETWORK'
  | 'CORPORATE';

export type PartnerStatus =
  | 'PENDING_VERIFICATION'
  | 'ACTIVE'
  | 'SUSPENDED'
  | 'TERMINATED';

export type PolicyType =
  | 'HEALTH'
  | 'MOTOR'
  | 'LIFE'
  | 'FIRE'
  | 'PET'
  | 'DEVICE'
  | 'TRAVEL';

export type ClaimStatus =
  | 'SUBMITTED'
  | 'UNDER_REVIEW'
  | 'APPROVED'
  | 'REJECTED'
  | 'SETTLED';

export type CommissionStatus =
  | 'CALCULATED'
  | 'APPROVED'
  | 'PAID'
  | 'REVERSED';

export type AgentStatus = 'ACTIVE' | 'INACTIVE' | 'SUSPENDED';

export type DocumentStatus = 'PENDING' | 'VERIFIED' | 'REJECTED' | 'EXPIRED';

export type ReferralStatus =
  | 'LEAD_CREATED'
  | 'CUSTOMER_CONTACTED'
  | 'QUOTE_GENERATED'
  | 'OTP_VERIFIED'
  | 'PAYMENT_COMPLETED'
  | 'POLICY_ISSUED';

export type PolicyStatus = 'ACTIVE' | 'EXPIRING' | 'LAPSED' | 'CANCELLED';

export interface Partner {
  id: string;
  organizationName: string;
  partnerType: PartnerType;
  status: PartnerStatus;
  tradeLicense: string;
  tin: string;
  contactEmail: string;
  contactPhone: string;
  bankAccount: string;
  bankName: string;
  bankBranch?: string;
  policyTypes: PolicyType[];
  cashlessEnabled: boolean;
  cashlessLimit?: number;
  discountEnabled: boolean;
  discountPercentage?: number;
  autoApprovalThreshold?: number;
  serviceLocations: string[];
  nationwideCoverage: boolean;
  acquisitionRate: number;
  renewalRate: number;
  claimsAssistanceRate?: number;
  agentCount: number;
  activePolicies: number;
  totalClaims: number;
  conversionRate: number;
  createdAt: string;
}

export interface Agent {
  id: string;
  partnerId: string;
  fullName: string;
  phone: string;
  email?: string;
  nid: string;
  status: AgentStatus;
  commissionRate: number;
  territory?: string;
  policiesSold: number;
  commissionEarned: number;
  referralCount: number;
  conversionRate: number;
  createdAt: string;
}

export interface Claim {
  id: string;
  claimNumber: string;
  partnerId: string;
  partnerName: string;
  policyId: string;
  policyNumber: string;
  policyType: PolicyType;
  customerName: string;
  customerNid: string;
  status: ClaimStatus;
  claimType: string;
  incidentDate: string;
  submittedDate: string;
  amount: number;
  approvedAmount?: number;
  settlementAmount?: number;
  tat: number; // Turnaround time in days
  documents: ClaimDocument[];
  itemizedBill: BillItem[];
  incidentDetails: string;
}

export interface ClaimDocument {
  id: string;
  name: string;
  type: string;
  url: string;
  uploadedAt: string;
}

export interface BillItem {
  category: string;
  description: string;
  amount: number;
}

export interface Policy {
  id: string;
  policyNumber: string;
  partnerId: string;
  agentId?: string;
  agentName?: string;
  customerName: string;
  customerNid: string;
  productType: PolicyType;
  productName: string;
  coverage: number;
  premium: number;
  status: PolicyStatus;
  issuedDate: string;
  expiryDate: string;
  nextRenewalDate?: string;
  claimCount: number;
}

export interface Commission {
  id: string;
  commissionId: string;
  partnerId: string;
  agentId?: string;
  agentName?: string;
  policyId: string;
  policyNumber: string;
  type: 'ACQUISITION' | 'RENEWAL' | 'CLAIMS_ASSISTANCE';
  premiumAmount: number;
  rate: number;
  commissionAmount: number;
  partnerShare: number;
  agentShare: number;
  status: CommissionStatus;
  calculatedDate: string;
  paidDate?: string;
}

export interface CommissionPayout {
  id: string;
  payoutNumber: string;
  period: string;
  periodStart: string;
  periodEnd: string;
  commissionCount: number;
  totalAmount: number;
  status: 'PENDING' | 'APPROVED' | 'PROCESSING' | 'PAID' | 'FAILED' | 'CANCELLED';
  paymentMethod: string;
  paidDate?: string;
}

export interface Referral {
  id: string;
  referralId: string;
  partnerId: string;
  agentId: string;
  agentName: string;
  customerName: string;
  customerPhone: string;
  customerEmail?: string;
  productInterest: PolicyType;
  status: ReferralStatus;
  quoteId?: string;
  policyId?: string;
  createdDate: string;
  lastUpdated: string;
}

export interface Document {
  id: string;
  partnerId: string;
  category: 'KYB' | 'MOU' | 'AGENT_KYC' | 'OPERATIONAL' | 'CLAIM_SUPPORT';
  name: string;
  type: string;
  status: DocumentStatus;
  uploadedDate: string;
  verifiedDate?: string;
  expiryDate?: string;
  version: number;
  url: string;
}

export interface DashboardKPI {
  activePolicies: number;
  claimsSubmittedMTD: number;
  claimsSettledMTD: number;
  averageClaimTAT: number;
  commissionEarnedMTD: number;
  commissionPending: number;
  activeAgents: number;
  conversionRate: number;
}

export interface ActivityFeedItem {
  id: string;
  type: 'CLAIM' | 'POLICY' | 'COMMISSION' | 'AGENT' | 'DOCUMENT';
  title: string;
  description: string;
  timestamp: string;
  icon: string;
}

export interface Announcement {
  id: string;
  title: string;
  message: string;
  type: 'INFO' | 'WARNING' | 'SUCCESS' | 'ERROR';
  publishedDate: string;
  isRead: boolean;
}
