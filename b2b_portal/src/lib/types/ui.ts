import type { LucideIcon } from "lucide-react";

import type {
  Claim,
  Invoice,
  Payment,
  Policy,
  Quotation,
} from "@lib/proto";

export interface NavItem {
  name: string;
  href: string;
  icon: string;
}

export interface StatsCardItem {
  id: number;
  title: string;
  value: string;
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

export interface AlertItem {
  id: number;
  count: number;
  label: string;
  barColor: string;
  bgColor: string;
}

export interface ActivityItem {
  id: number;
  title: string;
  subtitle: string;
  time: string;
  dotColor: string;
}

export interface PolicyCardItem {
  policy: Policy;
  icon: LucideIcon;
  type: string;
}

export interface UpcomingPaymentItem {
  payment: Payment;
  policyName: string;
  icon: LucideIcon;
  type: string;
  dueDate: string;
  status: "Due Soon" | "Upcoming";
}

export interface PaymentHistoryItem {
  payment: Payment;
  policyName: string;
  date: string;
}

export interface ClaimCardItem {
  claim: Claim;
}

export interface BillingInvoiceItem {
  invoice: Invoice;
}

export interface QuotationCardItem {
  quotation: Quotation;
}
