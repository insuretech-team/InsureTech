"use client";

import { ColumnDef } from "@tanstack/react-table";
import { LuPen, LuTrash2, LuLoader } from "react-icons/lu";
import { useState } from "react";
import type { Department } from "@lib/types/b2b";
import AddDepartmentModal from "@/components/modals/add-department-modal";
import { SortHeader } from "@/components/ui/sort-header";
import { departmentClient } from "@lib/sdk/department-client";

export type { Department };

// ─── Actions Cell ─────────────────────────────────────────────────────────────

function DepartmentActionsCell({
  dept,
  onRefresh,
}: {
  dept: Department;
  onRefresh?: () => void;
}) {
  const [editOpen, setEditOpen] = useState(false);
  const [deleting, setDeleting] = useState(false);

  async function handleDelete() {
    if (
      !confirm(
        `Delete department "${dept.name}"?\n\nThis will fail if the department has active employees.`
      )
    )
      return;

    setDeleting(true);
    try {
      const result = await departmentClient.delete(dept.id);
      if (!result.ok) { alert(result.message ?? "Failed to delete department"); return; }
      onRefresh?.();
    } catch {
      alert("Network error — could not delete department");
    } finally {
      setDeleting(false);
    }
  }

  return (
    <>
      <AddDepartmentModal
        open={editOpen}
        onOpenChange={setEditOpen}
        departmentId={dept.id}
        initialName={dept.name}
        onSaved={() => {
          setEditOpen(false);
          onRefresh?.();
        }}
      />

      <div className="flex items-center gap-1">
        <button
          className="rounded-md p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground"
          title="Edit department"
          onClick={() => setEditOpen(true)}
        >
          <LuPen className="size-4" />
        </button>
        <button
          className="rounded-md p-1.5 text-destructive hover:bg-destructive/10 disabled:opacity-50"
          title="Delete department"
          onClick={handleDelete}
          disabled={deleting}
        >
          {deleting ? <LuLoader className="size-4 animate-spin" /> : <LuTrash2 className="size-4" />}
        </button>
      </div>
    </>
  );
}

// ─── Column factory ───────────────────────────────────────────────────────────

export function buildDepartmentColumns(onRefresh?: () => void): ColumnDef<Department>[] {
  return [
    {
      id: "srNo",
      header: "Sr. No.",
      cell: ({ row }) => row.index + 1,
      enableSorting: false,
    },
    {
      accessorKey: "name",
      id: "name",
      header: ({ column }) => <SortHeader title="Department" column={column} />,
      cell: ({ row }) => (
        <div className="font-medium text-gray-900">{row.original.name}</div>
      ),
    },
    {
      accessorKey: "employeeNo",
      id: "employeeNo",
      header: ({ column }) => <SortHeader title="No. of Employees" column={column} />,
    },
    {
      accessorKey: "totalPremium",
      id: "totalPremium",
      header: ({ column }) => <SortHeader title="Total Premium" column={column} />,
    },
    {
      id: "actions",
      enableSorting: false,
      enableHiding: false,
      header: "Actions",
      cell: ({ row }) => (
        <DepartmentActionsCell dept={row.original} onRefresh={onRefresh} />
      ),
    },
  ];
}

// Backward-compatible static export
export const departmentColumns = buildDepartmentColumns();
