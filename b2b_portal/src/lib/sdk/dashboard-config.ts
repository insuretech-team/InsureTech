import { create } from "@bufbuild/protobuf";
import { Car, Heart, Shield } from "lucide-react";

import {
  ClaimSchema,
  ClaimStatus,
  ClaimType,
  InsuranceType,
  InvoiceSchema,
  InvoiceStatus,
  MoneySchema,
  PaymentMethod,
  PaymentSchema,
  PaymentStatus,
  PaymentType,
  PolicySchema,
  PolicyStatus,
  QuotationSchema,
  QuotationStatus,
  type Claim,
  type Invoice,
  type Payment,
  type Policy,
  type Quotation,
} from "@lib/proto";
import type {
  ActivityItem,
  AlertItem,
  NavItem,
  NotificationRule,
  PaymentHistoryItem,
  PolicyCardItem,
  QuotationSummaryItem,
  StatsCardItem,
  UpcomingPaymentItem,
  WorkflowRule,
} from "@lib/types/ui";

function money(amountBdt: number) {
  return create(MoneySchema, {
    amount: BigInt(Math.round(amountBdt * 100)),
    currency: "BDT",
    decimalAmount: amountBdt,
  });
}

function timestamp(dateLike: string) {
  const milliseconds = Date.parse(dateLike);
  return {
    seconds: BigInt(Math.floor(milliseconds / 1000)),
    nanos: (milliseconds % 1000) * 1_000_000,
  };
}

const navigation: NavItem[] = [
  { name: "Dashboard", href: "/", icon: "./navbar-icons/dashboard.svg" },
  { name: "Organisations", href: "/organisations", icon: "./navbar-icons/department.svg" },
  { name: "Employees", href: "/employees", icon: "./navbar-icons/employee.svg" },
  { name: "Departments", href: "/departments", icon: "./navbar-icons/department.svg" },
  { name: "Insurance Plans", href: "/insurance-plans", icon: "./navbar-icons/plans.svg" },
  { name: "Purchase Orders", href: "/purchase-orders", icon: "./navbar-icons/quotation.svg" },
  { name: "Billing & Invoices", href: "/billing-invoices", icon: "./navbar-icons/billing.svg" },
  { name: "My Profile", href: "/profile", icon: "./navbar-icons/employee.svg" },
  { name: "Settings", href: "/settings", icon: "./navbar-icons/settings.svg" },
];

// const statsCards: StatsCardItem[] = [
//   {
//     id: 1,
//     title: "Total Employees Covered",
//     value: "425",
//     icon: "./stats-cards/employee-count-icon.svg",
//     bgColor: "var(--brand-surface-1)",
//     bgIcon: "./stats-cards/employee-count-lg.svg",
//   },
//   {
//     id: 2,
//     title: "Active Policy",
//     value: "6",
//     icon: "./stats-cards/policies-icon.svg",
//     bgColor: "var(--brand-surface-2)",
//     bgIcon: "./stats-cards/policies-lg.svg",
//   },
//   {
//     id: 3,
//     title: "Total Premium",
//     value: "৳ 1,236,598",
//     icon: "./stats-cards/premium-icon.svg",
//     bgColor: "var(--brand-surface-3)",
//     bgIcon: "./stats-cards/premium-lg.svg",
//   },
//   {
//     id: 4,
//     title: "Pending Actions",
//     value: "13",
//     icon: "./stats-cards/actions-icon.svg",
//     bgColor: "var(--brand-surface-4)",
//     bgIcon: "./stats-cards/actions-lg.svg",
//   },
// ];

// const purchaseOrderSummary: QuotationSummaryItem[] = [
//   { id: 1, title: "Total Purchase Orders", value: 64, icon: "./quotations/comment-quote.svg", bgColor: "var(--brand-surface-1)" },
//   { id: 2, title: "In Draft", value: 12, icon: "./quotations/form.svg", bgColor: "var(--brand-surface-2)" },
//   { id: 3, title: "Submitted", value: 9, icon: "./quotations/paper-plane.svg", bgColor: "var(--brand-surface-5)" },
//   { id: 4, title: "Fulfilled", value: 5, icon: "./quotations/inbox-in.svg", bgColor: "var(--brand-surface-neutral)" },
//   { id: 5, title: "Approved", value: 25, icon: "./quotations/check-circle.svg", bgColor: "var(--brand-surface-3)" },
//   { id: 6, title: "Rejected", value: 13, icon: "./quotations/cross-circle.svg", bgColor: "var(--brand-surface-4)" },
// ];

// const purchaseOrders: Quotation[] = [
//   create(QuotationSchema, {
//     quotationId: "quo-001",
//     quotationNumber: "QUO-2025-001",
//     businessId: "biz-labaid-001",
//     insurerName: "ABC Insurance",
//     planId: "plan-seba",
//     planName: "Seba",
//     insuranceCategory: InsuranceType.HEALTH,
//     departmentId: "dep-tech",
//     employeeNo: 32,
//     estimatedPremium: money(432650),
//     quotedAmount: money(432650),
//     status: QuotationStatus.SUBMITTED,
//     submissionDate: timestamp("2025-01-12T00:00:00Z"),
//     validUntil: timestamp("2025-02-12T00:00:00Z"),
//   }),
//   create(QuotationSchema, {
//     quotationId: "quo-002",
//     quotationNumber: "QUO-2025-002",
//     businessId: "biz-labaid-001",
//     insurerName: "Delta Assurance",
//     planId: "plan-premium-care",
//     planName: "Premium Care",
//     insuranceCategory: InsuranceType.HEALTH,
//     departmentId: "dep-accounts",
//     employeeNo: 18,
//     estimatedPremium: money(210000),
//     quotedAmount: money(205000),
//     status: QuotationStatus.APPROVED,
//     submissionDate: timestamp("2025-01-20T00:00:00Z"),
//     validUntil: timestamp("2025-02-20T00:00:00Z"),
//   }),
//   create(QuotationSchema, {
//     quotationId: "quo-003",
//     quotationNumber: "QUO-2025-003",
//     businessId: "biz-labaid-001",
//     insurerName: "Guardian Life",
//     planId: "plan-standard",
//     planName: "Standard",
//     insuranceCategory: InsuranceType.LIFE,
//     departmentId: "dep-business",
//     employeeNo: 45,
//     estimatedPremium: money(520000),
//     quotedAmount: money(515000),
//     status: QuotationStatus.REJECTED,
//     submissionDate: timestamp("2025-01-28T00:00:00Z"),
//     validUntil: timestamp("2025-02-28T00:00:00Z"),
//   }),
// ];

const invoices: Invoice[] = [
  create(InvoiceSchema, {
    invoiceId: "inv-001",
    invoiceNumber: "INV-2025-02",
    businessId: "biz-labaid-001",
    amount: money(124680),
    dueDate: timestamp("2025-04-05T00:00:00Z"),
    status: InvoiceStatus.PENDING,
    issuedAt: timestamp("2025-03-05T00:00:00Z"),
  }),
  create(InvoiceSchema, {
    invoiceId: "inv-002",
    invoiceNumber: "INV-2025-03",
    businessId: "biz-labaid-001",
    amount: money(223410),
    dueDate: timestamp("2025-05-05T00:00:00Z"),
    status: InvoiceStatus.APPROVED,
    issuedAt: timestamp("2025-04-05T00:00:00Z"),
  }),
  create(InvoiceSchema, {
    invoiceId: "inv-003",
    invoiceNumber: "INV-2025-04",
    businessId: "biz-labaid-001",
    amount: money(189320),
    dueDate: timestamp("2025-06-05T00:00:00Z"),
    status: InvoiceStatus.APPROVED,
    issuedAt: timestamp("2025-05-05T00:00:00Z"),
  }),
];

const payments: Payment[] = [
  create(PaymentSchema, {
    paymentId: "pay-001",
    transactionId: "TXN-20250404-6921",
    policyId: "policy-health-001",
    type: PaymentType.PREMIUM,
    method: PaymentMethod.BANK_TRANSFER,
    status: PaymentStatus.SUCCESS,
    amount: money(124680),
    currency: "BDT",
    payerId: "biz-admin-001",
    initiatedAt: timestamp("2025-04-04T09:00:00Z"),
    completedAt: timestamp("2025-04-04T09:03:00Z"),
  }),
  create(PaymentSchema, {
    paymentId: "pay-002",
    transactionId: "TXN-20250510-1982",
    policyId: "policy-auto-001",
    type: PaymentType.PREMIUM,
    method: PaymentMethod.CARD,
    status: PaymentStatus.SUCCESS,
    amount: money(223410),
    currency: "BDT",
    payerId: "biz-admin-001",
    initiatedAt: timestamp("2025-05-10T08:30:00Z"),
    completedAt: timestamp("2025-05-10T08:32:00Z"),
  }),
  create(PaymentSchema, {
    paymentId: "pay-003",
    transactionId: "TXN-20250612-3388",
    policyId: "policy-life-001",
    type: PaymentType.PREMIUM,
    method: PaymentMethod.BKASH,
    status: PaymentStatus.PENDING,
    amount: money(189320),
    currency: "BDT",
    payerId: "biz-admin-001",
    initiatedAt: timestamp("2025-06-12T10:05:00Z"),
  }),
];

const claims: Claim[] = [
  create(ClaimSchema, {
    claimId: "claim-001",
    claimNumber: "CLM-2024-HLT-000001",
    policyId: "policy-health-001",
    customerId: "biz-admin-001",
    status: ClaimStatus.UNDER_REVIEW,
    type: ClaimType.HEALTH_HOSPITALIZATION,
    claimedAmount: money(45000),
    submittedAt: timestamp("2024-12-15T09:30:00Z"),
    incidentDate: timestamp("2024-12-14T00:00:00Z"),
    incidentDescription: "Hospitalization due to acute condition.",
  }),
  create(ClaimSchema, {
    claimId: "claim-002",
    claimNumber: "CLM-2024-MTR-000002",
    policyId: "policy-auto-001",
    customerId: "biz-admin-001",
    status: ClaimStatus.SUBMITTED,
    type: ClaimType.MOTOR_ACCIDENT,
    claimedAmount: money(125000),
    submittedAt: timestamp("2024-12-10T12:30:00Z"),
    incidentDate: timestamp("2024-12-09T00:00:00Z"),
    incidentDescription: "Motor accident damage claim.",
  }),
  create(ClaimSchema, {
    claimId: "claim-003",
    claimNumber: "CLM-2024-HLT-000003",
    policyId: "policy-health-001",
    customerId: "biz-admin-001",
    status: ClaimStatus.APPROVED,
    type: ClaimType.HEALTH_SURGERY,
    claimedAmount: money(18500),
    approvedAmount: money(18500),
    submittedAt: timestamp("2024-11-28T13:30:00Z"),
    incidentDate: timestamp("2024-11-27T00:00:00Z"),
    incidentDescription: "Approved surgery reimbursement.",
  }),
  create(ClaimSchema, {
    claimId: "claim-004",
    claimNumber: "CLM-2024-LIF-000004",
    policyId: "policy-life-001",
    customerId: "biz-admin-001",
    status: ClaimStatus.SETTLED,
    type: ClaimType.DEATH,
    claimedAmount: money(75000),
    settledAmount: money(75000),
    submittedAt: timestamp("2024-11-15T10:00:00Z"),
    settledAt: timestamp("2024-11-25T10:00:00Z"),
    incidentDate: timestamp("2024-11-14T00:00:00Z"),
    incidentDescription: "Settled life claim.",
  }),
];

const policies: Policy[] = [
  create(PolicySchema, {
    policyId: "policy-health-001",
    policyNumber: "LBT-2025-HLT-000001",
    productId: "product-health-seba",
    customerId: "biz-admin-001",
    status: PolicyStatus.GRACE_PERIOD,
    premiumAmount: money(14500),
    sumInsured: money(225500),
    tenureMonths: 12,
    startDate: timestamp("2025-01-01T00:00:00Z"),
    endDate: timestamp("2025-12-28T00:00:00Z"),
    providerName: "Chartered Life",
  }),
  create(PolicySchema, {
    policyId: "policy-auto-001",
    policyNumber: "LBT-2025-AUT-000002",
    productId: "product-auto-standard",
    customerId: "biz-admin-001",
    status: PolicyStatus.ACTIVE,
    premiumAmount: money(8500),
    sumInsured: money(625000),
    tenureMonths: 12,
    startDate: timestamp("2025-01-12T00:00:00Z"),
    endDate: timestamp("2026-01-11T00:00:00Z"),
    providerName: "National Insurance",
  }),
  create(PolicySchema, {
    policyId: "policy-life-001",
    policyNumber: "LBT-2025-LIF-000003",
    productId: "product-life-plus",
    customerId: "biz-admin-001",
    status: PolicyStatus.ACTIVE,
    premiumAmount: money(11250),
    sumInsured: money(1200000),
    tenureMonths: 24,
    startDate: timestamp("2024-12-03T00:00:00Z"),
    endDate: timestamp("2026-12-02T00:00:00Z"),
    providerName: "MetLife Bangladesh",
  }),
];

const policyById = new Map(policies.map((policy) => [policy.policyId, policy]));

const notificationRules: NotificationRule[] = [
  { id: 1, value: 1, title: "Policy Expiry Alerts", description: "Get notified when policies are expiring soon." },
  { id: 2, value: 2, title: "Purchase Order Updates", description: "Receive alerts when purchase orders change status." },
  { id: 3, value: 3, title: "Invoice Reminders", description: "Get reminders for upcoming invoice due dates." },
  { id: 4, value: 4, title: "Employee Coverage Changes", description: "Notify when employee coverage is added or removed." },
  { id: 5, value: 5, title: "Weekly Summary Report", description: "Receive weekly summaries of insurance activities." },
];

const workflowRules: WorkflowRule[] = [
  { id: 1, value: 1, title: "Purchase Order Approval", description: "Require manager approval before submitting purchase orders." },
  { id: 2, value: 2, title: "Plan Changes", description: "Require approval for bulk employee plan changes." },
  { id: 3, value: 3, title: "Payment Authorization", description: "Require finance approval for payments over BDT 50,000." },
  { id: 4, value: 4, title: "Policy Renewals", description: "Require executive approval for policy renewals." },
];

const alerts: AlertItem[] = [
  { id: 1, count: 6, label: "Policies expiring in 30 days", barColor: "var(--brand-danger)", bgColor: "var(--brand-surface-danger)" },
  { id: 2, count: 4, label: "Employees without coverage", barColor: "var(--brand-warning)", bgColor: "var(--brand-surface-warning)" },
  { id: 3, count: 3, label: "Premium due", barColor: "var(--brand-alert)", bgColor: "var(--brand-surface-alert)" },
];

const activities: ActivityItem[] = [
  { id: 1, title: "Quotation requested", subtitle: "Health Insurance", time: "15 minutes ago", dotColor: "var(--brand-primary)" },
  { id: 2, title: "Plan updated", subtitle: "Health Insurance with premium BDT 6,200", time: "2 hours ago", dotColor: "var(--brand-secondary)" },
  { id: 3, title: "Policy expiring", subtitle: "1 health plan expiring in 16 days", time: "8 hours ago", dotColor: "var(--brand-danger)" },
  { id: 4, title: "Employee added", subtitle: "15 new employees added via CSV upload", time: "2 days ago", dotColor: "var(--brand-primary)" },
];

const policyCards: PolicyCardItem[] = [
  { policy: policies[0], icon: Heart, type: "Health" },
  { policy: policies[1], icon: Car, type: "Vehicle" },
  { policy: policies[2], icon: Shield, type: "Life" },
];

const upcomingPayments: UpcomingPaymentItem[] = [
  { payment: payments[0], policyName: "Health Insurance", type: "Health", icon: Heart, dueDate: "22-12-2025", status: "Due Soon" },
  { payment: payments[1], policyName: "Auto Insurance", type: "Vehicle", icon: Car, dueDate: "27-12-2025", status: "Due Soon" },
  { payment: payments[2], policyName: "Life Insurance", type: "Life", icon: Shield, dueDate: "02-12-2026", status: "Upcoming" },
];

const paymentHistory: PaymentHistoryItem[] = [
  { payment: payments[0], policyName: "Health Insurance", date: "22-11-2024" },
  { payment: payments[1], policyName: "Auto Insurance", date: "27-10-2024" },
  { payment: payments[2], policyName: "Life Insurance", date: "02-10-2024" },
];

// Role-based navigation definitions
// Super admin: sees everything including Organisations management
const superAdminNav: NavItem[] = [
  { name: "Dashboard", href: "/", icon: "./navbar-icons/dashboard.svg" },
  { name: "Organisations", href: "/organisations", icon: "./navbar-icons/department.svg" },
  { name: "Employees", href: "/employees", icon: "./navbar-icons/employee.svg" },
  { name: "Departments", href: "/departments", icon: "./navbar-icons/department.svg" },
  { name: "Insurance Plans", href: "/insurance-plans", icon: "./navbar-icons/plans.svg" },
  { name: "Purchase Orders", href: "/purchase-orders", icon: "./navbar-icons/quotation.svg" },
  { name: "Billing & Invoices", href: "/billing-invoices", icon: "./navbar-icons/billing.svg" },
  { name: "My Profile", href: "/profile", icon: "./navbar-icons/employee.svg" },
  { name: "Settings", href: "/settings", icon: "./navbar-icons/settings.svg" },
];

// B2B org admin: full access to their own org but no Organisations management
const b2bOrgAdminNav: NavItem[] = [
  { name: "Dashboard", href: "/", icon: "./navbar-icons/dashboard.svg" },
  { name: "Departments", href: "/departments", icon: "./navbar-icons/department.svg" },
  { name: "Employees", href: "/employees", icon: "./navbar-icons/employee.svg" },
  { name: "Insurance Plans", href: "/insurance-plans", icon: "./navbar-icons/plans.svg" },
  { name: "Purchase Orders", href: "/purchase-orders", icon: "./navbar-icons/quotation.svg" },
  { name: "Billing & Invoices", href: "/billing-invoices", icon: "./navbar-icons/billing.svg" },
  { name: "Team", href: "/team", icon: "./navbar-icons/employee.svg" },
  { name: "My Profile", href: "/profile", icon: "./navbar-icons/employee.svg" },
  { name: "Settings", href: "/settings", icon: "./navbar-icons/settings.svg" },
];

// HR Manager / Viewer / partner user: read + manage employees/departments, no team management
const partnerUserNav: NavItem[] = [
  { name: "Dashboard", href: "/", icon: "./navbar-icons/dashboard.svg" },
  { name: "Departments", href: "/departments", icon: "./navbar-icons/department.svg" },
  { name: "Employees", href: "/employees", icon: "./navbar-icons/employee.svg" },
  { name: "Insurance Plans", href: "/insurance-plans", icon: "./navbar-icons/plans.svg" },
  { name: "Purchase Orders", href: "/purchase-orders", icon: "./navbar-icons/quotation.svg" },
  { name: "Billing & Invoices", href: "/billing-invoices", icon: "./navbar-icons/billing.svg" },
  { name: "My Profile", href: "/profile", icon: "./navbar-icons/employee.svg" },
  { name: "Settings", href: "/settings", icon: "./navbar-icons/settings.svg" },
];

export const b2bDashboardClient = {
  getNavigation: (role?: string) => {
    if (role === "SYSTEM_ADMIN") return superAdminNav;
    if (role === "B2B_ORG_ADMIN" || role === "BUSINESS_ADMIN") return b2bOrgAdminNav;
    // HR_MANAGER, VIEWER, partner user, or not-yet-loaded (undefined → full nav as safe default)
    return partnerUserNav;
  },
  getStatsCards: () => [] as any,
  getPurchaseOrderSummary: () => [] as any,
  getPurchaseOrders: () => [] as any,
  getInvoices: () => invoices,
  getPayments: () => payments,
  getPolicyById: (policyId: string): Policy | undefined => policyById.get(policyId),
  getClaims: () => claims,
  getPolicies: () => policies,
  getPolicyCards: () => policyCards,
  getUpcomingPayments: () => upcomingPayments,
  getPaymentHistory: () => paymentHistory,
  getNotificationRules: () => notificationRules,
  getWorkflowRules: () => workflowRules,
  getAlerts: () => alerts,
  getActivities: () => activities,
};
