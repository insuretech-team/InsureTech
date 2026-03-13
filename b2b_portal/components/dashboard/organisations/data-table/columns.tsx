/**
 * Organisation table columns
 */
"use client";

import { ColumnDef } from "@tanstack/react-table";
import { LuPen, LuTrash2, LuLoader, LuEye, LuCheck } from "react-icons/lu";
import { useState } from "react";
import { SortHeader } from "@/components/ui/sort-header";
import { StatusBadge } from "@/components/ui/status-badge";
import AddOrganisationModal from "@/components/modals/add-organisation-modal";
import { organisationClient } from "@lib/sdk/organisation-client";
import type { Organisation } from "@lib/types/b2b";

export type { Organisation };

function OrgActionsCell({
  org,
  onRefresh,
  onView,
  onApprove,
}: {
  org: Organisation;
  onRefresh?: () => void;
  onView?: (org: Organisation) => void;
  onApprove?: (org: Organisation) => void;
}) {
  const [editOpen, setEditOpen] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [approving, setApproving] = useState(false);

  const isPending = (org.status ?? "").toUpperCase().includes("PENDING");

  async function handleDelete() {
    if (!confirm(`Delete organisation "${org.name}"?\nThis cannot be undone.`)) return;
    setDeleting(true);
    try {
      const result = await organisationClient.delete(org.id);
      if (!result.ok) { alert(result.message ?? "Delete failed"); return; }
      onRefresh?.();
    } finally {
      setDeleting(false);
    }
  }

  async function handleApprove() {
    if (!confirm(`Approve "${org.name}"?\nThis will activate their account and notify the admin.`)) return;
    setApproving(true);
    try {
      const result = await organisationClient.approve(org.id);
      if (!result.ok) { alert(result.message ?? "Approve failed"); return; }
      onApprove?.(org);
      onRefresh?.();
    } finally {
      setApproving(false);
    }
  }

  return (
    <>
      <AddOrganisationModal
        open={editOpen}
        onOpenChange={setEditOpen}
        orgId={org.id}
        initialValues={{
          name: org.name,
          code: org.code ?? "",
          industry: org.industry,
          contactEmail: org.contactEmail,
          contactPhone: org.contactPhone,
          address: org.address,
        }}
        onSaved={() => { setEditOpen(false); onRefresh?.(); }}
      />

      <div className="flex items-center gap-1">
        <button
          className="rounded-md p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground"
          title="View details"
          onClick={() => onView?.(org)}
        >
          <LuEye className="size-4" />
        </button>
        <button
          className="rounded-md p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground"
          title="Edit organisation"
          onClick={() => setEditOpen(true)}
        >
          <LuPen className="size-4" />
        </button>
        {/* Approve button — only shown for PENDING orgs */}
        {isPending && (
          <button
            className="rounded-md p-1.5 text-green-600 hover:bg-green-50 disabled:opacity-50"
            title="Approve organisation"
            onClick={handleApprove}
            disabled={approving}
          >
            {approving
              ? <LuLoader className="size-4 animate-spin" />
              : <LuCheck className="size-4" />}
          </button>
        )}
        <button
          className="rounded-md p-1.5 text-destructive hover:bg-destructive/10 disabled:opacity-50"
          title="Delete organisation"
          onClick={handleDelete}
          disabled={deleting}
        >
          {deleting ? <LuLoader className="size-4 animate-spin" /> : <LuTrash2 className="size-4" />}
        </button>
      </div>
    </>
  );
}

export function buildOrganisationColumns(
  onRefresh?: () => void,
  onRowClick?: (org: Organisation) => void,
  onApprove?: (org: Organisation) => void,
): ColumnDef<Organisation>[] {
  return [
    {
      id: "srNo",
      header: "Sr. No.",
      cell: ({ row }) => row.index + 1,
      enableSorting: false,
    },
    {
      accessorKey: "name",
      header: ({ column }) => <SortHeader title="Organisation Name" column={column} />,
      cell: ({ row }) => (
        <div
          className="font-medium text-gray-900 cursor-pointer hover:text-primary underline-offset-2 hover:underline"
          onClick={() => onRowClick?.(row.original)}
        >
          {row.original.name}
        </div>
      ),
    },
    {
      accessorKey: "code",
      header: ({ column }) => <SortHeader title="Code" column={column} />,
    },
    {
      accessorKey: "industry",
      header: ({ column }) => <SortHeader title="Industry" column={column} />,
    },
    {
      accessorKey: "contactEmail",
      header: ({ column }) => <SortHeader title="Email" column={column} />,
    },
    {
      accessorKey: "contactPhone",
      header: ({ column }) => <SortHeader title="Phone" column={column} />,
    },
    {
      accessorKey: "status",
      header: ({ column }) => <SortHeader title="Status" column={column} />,
      cell: ({ row }) => <StatusBadge status={row.original.status} />,
    },
    {
      id: "actions",
      header: "Actions",
      enableSorting: false,
      enableHiding: false,
      cell: ({ row }) => (
        <OrgActionsCell org={row.original} onRefresh={onRefresh} onView={onRowClick} onApprove={onApprove} />
      ),
    },
  ];
}

export const organisationColumns = buildOrganisationColumns();
