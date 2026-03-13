"use client";

import * as React from "react";
import { ColumnDef } from "@tanstack/react-table";
import { LuEye, LuTrash2, LuLoader, LuX } from "react-icons/lu";
import { SortHeader } from "@/components/ui/sort-header";
import { StatusBadge } from "@/components/ui/status-badge";
import { purchaseOrderClient } from "@lib/sdk/purchase-order-client";

export type PurchaseOrder = {
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
  status: string;
  submittedAt: string;
  notes: string;
};

// ─── Detail Sheet ─────────────────────────────────────────────────────────────

function PODetailSheet({ po, onClose }: { po: PurchaseOrder; onClose: () => void }) {
  const rows: [string, React.ReactNode][] = [
    ["PO Number",         po.purchaseOrderNumber],
    ["Product",           po.productName],
    ["Plan",              po.planName],
    ["Insurance Category", po.insuranceCategory],
    ["Department",        po.department],
    ["Employees",         po.employeeCount],
    ["Dependents",        po.numberOfDependents],
    ["Coverage Amount",   po.coverageAmount],
    ["Estimated Premium", po.estimatedPremium],
    ["Status",            <StatusBadge key="s" status={po.status} />],
    ["Submitted On",      po.submittedAt],
    ["Notes",             po.notes || "—"],
  ];

  return (
    <div className="fixed inset-0 z-50 flex justify-end">
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/30" onClick={onClose} />
      {/* Panel */}
      <div className="relative z-10 w-full max-w-md bg-white shadow-xl flex flex-col h-full">
        <div className="flex items-center justify-between px-5 py-4 border-b">
          <h2 className="text-lg font-semibold">Purchase Order Details</h2>
          <button onClick={onClose} className="rounded-md p-1.5 hover:bg-muted">
            <LuX className="size-5" />
          </button>
        </div>
        <div className="flex-1 overflow-y-auto px-5 py-4 space-y-3">
          {rows.map(([label, value]) => (
            <div key={label} className="flex flex-col gap-0.5">
              <span className="text-xs font-semibold text-muted-foreground uppercase tracking-wide">{label}</span>
              <span className="text-sm text-foreground">{value}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

// ─── Actions Cell ─────────────────────────────────────────────────────────────

function POActionsCell({ po, onRefresh }: { po: PurchaseOrder; onRefresh?: () => void }) {
  const [detailOpen, setDetailOpen] = React.useState(false);
  const [deleting, setDeleting] = React.useState(false);

  async function handleDelete() {
    if (!confirm(`Delete purchase order "${po.purchaseOrderNumber}"?\nThis cannot be undone.`)) return;
    setDeleting(true);
    try {
      const result = await purchaseOrderClient.delete(po.id);
      if (!result.ok) {
        alert(result.message ?? "Failed to delete purchase order");
        return;
      }
      onRefresh?.();
    } finally {
      setDeleting(false);
    }
  }

  return (
    <>
      {detailOpen && <PODetailSheet po={po} onClose={() => setDetailOpen(false)} />}

      <div className="flex items-center gap-1">
        <button
          className="rounded-md p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground"
          title="View details"
          onClick={() => setDetailOpen(true)}
        >
          <LuEye className="size-4" />
        </button>
        <button
          className="rounded-md p-1.5 text-destructive hover:bg-destructive/10 disabled:opacity-50"
          title="Delete purchase order"
          onClick={handleDelete}
          disabled={deleting}
        >
          {deleting
            ? <LuLoader className="size-4 animate-spin" />
            : <LuTrash2 className="size-4" />}
        </button>
      </div>
    </>
  );
}

// ─── Column Definitions ───────────────────────────────────────────────────────

export function buildPurchaseOrderColumns(onRefresh?: () => void): ColumnDef<PurchaseOrder>[] {
  return [
    {
      id: "srNo",
      header: "Sr. No.",
      cell: ({ row }) => row.index + 1,
      enableSorting: false,
    },
    {
      accessorKey: "purchaseOrderNumber",
      header: ({ column }) => <SortHeader title="PO Number" column={column} />,
      cell: ({ row }) => <div className="font-medium text-gray-900">{row.original.purchaseOrderNumber}</div>,
    },
    {
      accessorKey: "productName",
      header: ({ column }) => <SortHeader title="Product" column={column} />,
    },
    {
      accessorKey: "planName",
      header: ({ column }) => <SortHeader title="Plan" column={column} />,
    },
    {
      accessorKey: "insuranceCategory",
      header: ({ column }) => <SortHeader title="Category" column={column} />,
    },
    {
      accessorKey: "department",
      header: ({ column }) => <SortHeader title="Department" column={column} />,
    },
    {
      accessorKey: "employeeCount",
      header: ({ column }) => <SortHeader title="Employees" column={column} />,
    },
    {
      accessorKey: "numberOfDependents",
      header: ({ column }) => <SortHeader title="Dependents" column={column} />,
    },
    {
      accessorKey: "coverageAmount",
      header: ({ column }) => <SortHeader title="Coverage" column={column} />,
    },
    {
      accessorKey: "estimatedPremium",
      header: ({ column }) => <SortHeader title="Estimated Premium" column={column} />,
    },
    {
      accessorKey: "status",
      header: ({ column }) => <SortHeader title="Status" column={column} />,
      cell: ({ row }) => <StatusBadge status={row.original.status} />,
    },
    {
      accessorKey: "submittedAt",
      header: ({ column }) => <SortHeader title="Submitted On" column={column} />,
    },
    {
      id: "actions",
      header: "Actions",
      enableSorting: false,
      enableHiding: false,
      cell: ({ row }) => <POActionsCell po={row.original} onRefresh={onRefresh} />,
    },
  ];
}

// Backward-compatible static export
export const purchaseOrderColumns = buildPurchaseOrderColumns();
