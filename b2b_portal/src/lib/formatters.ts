import type { Timestamp } from "@bufbuild/protobuf/wkt";

import {
  ClaimStatus,
  ClaimType,
  EmployeeStatus,
  InsuranceType,
  InvoiceStatus,
  PaymentMethod,
  PaymentStatus,
  PolicyStatus,
  PurchaseOrderStatus,
  QuotationStatus,
  type Money,
} from "@lib/proto";

export function formatMoney(money: Money | undefined): string {
  if (!money) {
    return "-";
  }
  const amount = Number(money.amount) / 100;
  return new Intl.NumberFormat("bn-BD", {
    style: "currency",
    currency: money.currency || "BDT",
    maximumFractionDigits: 0,
  }).format(amount);
}

export function formatDate(timestamp: Timestamp | undefined): string {
  if (!timestamp) {
    return "-";
  }
  const seconds = Number(timestamp.seconds ?? BigInt(0));
  if (!seconds) {
    return "-";
  }
  return new Date(seconds * 1000).toLocaleDateString("en-GB");
}

export function insuranceTypeLabel(value: InsuranceType): string {
  switch (value) {
    case InsuranceType.HEALTH:
      return "Health";
    case InsuranceType.LIFE:
      return "Life";
    case InsuranceType.AUTO:
      return "Auto";
    case InsuranceType.TRAVEL:
      return "Travel";
    case InsuranceType.DEVICE:
      return "Device";
    default:
      return "Unspecified";
  }
}

export function employeeStatusLabel(value: EmployeeStatus): "Active" | "Inactive" {
  return value === EmployeeStatus.ACTIVE ? "Active" : "Inactive";
}

export function quotationStatusLabel(value: QuotationStatus): string {
  switch (value) {
    case QuotationStatus.DRAFT:
      return "In Draft";
    case QuotationStatus.SUBMITTED:
      return "Submitted";
    case QuotationStatus.RECEIVED:
      return "Received";
    case QuotationStatus.APPROVED:
      return "Approved";
    case QuotationStatus.REJECTED:
      return "Rejected";
    default:
      return "Pending";
  }
}

export function purchaseOrderStatusLabel(value: PurchaseOrderStatus): string {
  switch (value) {
    case PurchaseOrderStatus.DRAFT:
      return "In Draft";
    case PurchaseOrderStatus.SUBMITTED:
      return "Submitted";
    case PurchaseOrderStatus.APPROVED:
      return "Approved";
    case PurchaseOrderStatus.FULFILLED:
      return "Fulfilled";
    case PurchaseOrderStatus.REJECTED:
      return "Rejected";
    default:
      return "Pending";
  }
}

export function invoiceStatusLabel(value: InvoiceStatus): string {
  switch (value) {
    case InvoiceStatus.PENDING:
      return "Pending";
    case InvoiceStatus.APPROVED:
      return "Approved";
    case InvoiceStatus.PAID:
      return "Paid";
    case InvoiceStatus.OVERDUE:
      return "Overdue";
    case InvoiceStatus.CANCELLED:
      return "Cancelled";
    default:
      return "Pending";
  }
}

export function paymentMethodLabel(value: PaymentMethod): string {
  switch (value) {
    case PaymentMethod.BANK_TRANSFER:
      return "Bank Transfer";
    case PaymentMethod.BKASH:
      return "bKash";
    case PaymentMethod.NAGAD:
      return "Nagad";
    case PaymentMethod.CARD:
      return "Card";
    default:
      return "Other";
  }
}

export function paymentStatusLabel(value: PaymentStatus): string {
  switch (value) {
    case PaymentStatus.SUCCESS:
      return "Paid";
    case PaymentStatus.PENDING:
      return "Pending";
    case PaymentStatus.FAILED:
      return "Failed";
    default:
      return "Pending";
  }
}

export function policyStatusLabel(value: PolicyStatus): "Active" | "Expiring" {
  switch (value) {
    case PolicyStatus.ACTIVE:
      return "Active";
    case PolicyStatus.GRACE_PERIOD:
      return "Expiring";
    default:
      return "Active";
  }
}

export function claimStatusLabel(value: ClaimStatus): string {
  switch (value) {
    case ClaimStatus.UNDER_REVIEW:
      return "Under Review";
    case ClaimStatus.PENDING_DOCUMENTS:
      return "Pending Documents";
    case ClaimStatus.APPROVED:
      return "Approved";
    case ClaimStatus.SETTLED:
      return "Paid";
    case ClaimStatus.REJECTED:
      return "Rejected";
    default:
      return "Processing";
  }
}

export function claimTypeLabel(value: ClaimType): string {
  switch (value) {
    case ClaimType.HEALTH_HOSPITALIZATION:
    case ClaimType.HEALTH_SURGERY:
      return "Health";
    case ClaimType.MOTOR_ACCIDENT:
    case ClaimType.MOTOR_THEFT:
      return "Auto";
    case ClaimType.DEATH:
      return "Life";
    default:
      return "Claim";
  }
}
