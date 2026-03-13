import type { LucideIcon } from "lucide-react";

// ─── Enums / Unions ───────────────────────────────────────────────────────────

export type EmployeeStatus = "Active" | "Inactive";

export type QuotationStatus =
  | "Approved"
  | "Submitted"
  | "Received"
  | "In Draft"
  | "Rejected"
  | "Pending";

export type PurchaseOrderStatusLabel =
  | "In Draft"
  | "Submitted"
  | "Approved"
  | "Fulfilled"
  | "Rejected"
  | "Pending";

export type OrganisationStatusLabel = "Active" | "Inactive" | "Suspended" | "Pending";

// ─── Session / Auth ───────────────────────────────────────────────────────────

export interface SessionUser {
  userId: string;
  businessId: string;
  /** Organisation name shown in the header/sidebar for B2B admin users. Empty for super_admin. */
  organisationName: string;
  name: string;
  email: string;
  role: "BUSINESS_ADMIN" | "FINANCE_MANAGER" | "HR_MANAGER" | "SYSTEM_ADMIN" | "B2B_ORG_ADMIN";
}

// ─── Navigation ───────────────────────────────────────────────────────────────

export interface NavItem {
  name: string;
  href: string;
  icon: string;
}

// ─── Stats / Dashboard ────────────────────────────────────────────────────────

export interface StatsCardItem {
  id: number;
  title: string;
  value: string | number;
  icon: string;
  bgColor: string;
  bgIcon: string;
}

export interface QuotationSummaryItem {
  id: number;
  title: string;
  value: number;
  icon: string;
  bgColor: string;
}

// ─── Core B2B Entities ────────────────────────────────────────────────────────

export interface Organisation {
  id: string;
  name: string;
  code?: string;
  industry: string;
  contactEmail: string;
  contactPhone: string;
  address: string;
  status: OrganisationStatusLabel | string;
  totalEmployees?: number;
  createdAt?: string;
}

export interface Department {
  id: string;
  name: string;
  employeeNo: number;
  totalPremium: string;
}

export interface Employee {
  id: string;
  name: string;
  employeeID: string;
  department: string;
  insuranceCategory?: string;
  assignedPlan?: string;
  coverage: string;
  premiumAmount: string;
  status: EmployeeStatus;
  numberOfDependent: number;
}

export interface PurchaseOrder {
  id: string;
  purchaseOrderNumber: string;
  productName: string;
  planName: string;
  insuranceCategory: string;
  department: string;
  employeeCount: number;
  numberOfDependents: number;
  coverageAmount: string;
  estimatedPremium: string;
  status: PurchaseOrderStatusLabel | string;
  submittedAt: string;
  notes: string;
}

// ─── Quotation ────────────────────────────────────────────────────────────────

export interface Quotation {
  id: string;
  quotationID: string;
  insurerName: string;
  plan: string;
  insuranceCategory: string;
  department: string;
  employeeNo: number;
  estimatedPremium: string;
  quotedAmount: string;
  status: QuotationStatus;
  submissionDate: string;
  validUntil: string;
}

// ─── Finance ──────────────────────────────────────────────────────────────────

export interface Invoice {
  id: string;
  dueDate: string;
  status: "Pending" | "Approved" | "Rejected";
  amount: number;
}

export interface Payment {
  id: string;
  amount: number;
  method: string;
  month: string;
  status: "Paid" | "Pending" | "Failed";
}

// ─── Settings ─────────────────────────────────────────────────────────────────

export interface NotificationRule {
  id: number;
  value: number;
  title: string;
  description: string;
}

export interface WorkflowRule {
  id: number;
  value: number;
  title: string;
  description: string;
}

// ─── Claims / Policies ────────────────────────────────────────────────────────

export interface ClaimItem {
  id: string;
  policyType: string;
  claimAmount: string;
  filedDate: string;
  status: string;
  statusClassName: string;
}

export interface PolicyItem {
  name: string;
  type: string;
  icon: LucideIcon;
  enrollmentId: string;
  coverage: string;
  premium: string;
  nextDue: string;
  status: "Active" | "Expiring";
  provider: string;
}

export interface UpcomingPaymentItem {
  policyName: string;
  type: string;
  icon: LucideIcon;
  amount: string;
  dueDate: string;
  status: "Due Soon" | "Upcoming";
}

export interface PaymentHistoryItem {
  policyName: string;
  amount: string;
  date: string;
  status: "Paid";
}
