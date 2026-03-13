/**
 * status-badge.tsx
 * ─────────────────
 * Reusable coloured status pill used in all data tables.
 */

const STATUS_STYLES: Record<string, string> = {
  // Employee
  Active: "bg-green-50 text-green-700",
  Inactive: "bg-red-50 text-red-700",
  // Purchase Order / Quotation
  Approved: "bg-green-50 text-green-700",
  Submitted: "bg-blue-50 text-blue-700",
  Fulfilled: "bg-amber-50 text-amber-700",
  "In Draft": "bg-yellow-50 text-yellow-700",
  Rejected: "bg-red-50 text-red-700",
  Pending: "bg-gray-50 text-gray-700",
  // Organisation
  ACTIVE: "bg-green-50 text-green-700",
  INACTIVE: "bg-red-50 text-red-700",
  SUSPENDED: "bg-orange-50 text-orange-700",
};

interface StatusBadgeProps {
  status: string;
}

export function StatusBadge({ status }: StatusBadgeProps) {
  const style = STATUS_STYLES[status] ?? "bg-gray-50 text-gray-700";
  return (
    <span className={`inline-flex items-center rounded-full px-2 py-1 text-xs font-medium ${style}`}>
      {status}
    </span>
  );
}
