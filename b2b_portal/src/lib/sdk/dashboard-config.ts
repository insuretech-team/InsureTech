/**
 * dashboard-config.ts
 * ────────────────────
 * Static configuration for the B2B portal dashboard.
 * Contains only nav items and settings rules — NO mock data.
 * All real data is fetched via bffClient or API routes.
 */
import type { NavItem, NotificationRule, WorkflowRule } from "@lib/types/ui";

// ─── Role-based navigation ────────────────────────────────────────────────────

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

const beneficiaryNav: NavItem[] = [
  { name: "My Coverage", href: "/employee", icon: "./navbar-icons/plans.svg" },
  { name: "My Profile", href: "/profile", icon: "./navbar-icons/employee.svg" },
];

// ─── Settings rules (static, not from API) ───────────────────────────────────

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

// ─── Dashboard config client ──────────────────────────────────────────────────

export const b2bDashboardClient = {
  /** Returns the correct nav items for the given portal role. */
  getNavigation: (role?: string): NavItem[] => {
    if (role === "SYSTEM_ADMIN") return superAdminNav;
    if (role === "B2B_ORG_ADMIN" || role === "BUSINESS_ADMIN") return b2bOrgAdminNav;
    if (role === "B2B_BENEFICIARY") return beneficiaryNav;
    return partnerUserNav;
  },
  getNotificationRules: (): NotificationRule[] => notificationRules,
  getWorkflowRules: (): WorkflowRule[] => workflowRules,
};
