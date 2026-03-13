"use client";

import { ColumnDef } from "@tanstack/react-table";
import { LuPen, LuEye, LuTrash2, LuLoader } from "react-icons/lu";
import { useState } from "react";
import type { Employee } from "@lib/types/b2b";
import AddEmployeeModal from "@/components/modals/add-employee-modal";
import { SortHeader } from "@/components/ui/sort-header";
import { StatusBadge } from "@/components/ui/status-badge";
import type { EmployeeFormValues } from "@/src/lib/types/employee-form";

export type { Employee };

// ─── Actions Cell ──────────────────────────────────────────────────────────────
// Extracted as a component so it can use hooks (useState).

function EmployeeActionsCell({
  emp,
  onRefresh,
}: {
  emp: Employee;
  onRefresh?: () => void;
}) {
  const [editOpen, setEditOpen] = useState(false);
  const [deleting, setDeleting] = useState(false);
  // View opens the edit modal in read-only context (same modal, employeeUuid populates all fields)
  const [viewOpen, setViewOpen] = useState(false);

  const initialValues: Partial<EmployeeFormValues> = {
    name: emp.name,
    employeeId: emp.employeeID,
    departmentId: "", // not available in list row — modal fetches full record by uuid
  };

  async function handleDelete() {
    if (!confirm(`Delete employee "${emp.name}"?\nThis cannot be undone.`)) return;
    setDeleting(true);
    try {
      const res = await fetch(`/api/employees/${emp.id}`, { method: "DELETE" });
      const payload = (await res.json()) as { ok: boolean; message?: string };
      if (!res.ok || !payload.ok) {
        alert(payload.message ?? "Failed to delete employee");
        return;
      }
      onRefresh?.();
    } catch {
      alert("Network error — could not delete employee");
    } finally {
      setDeleting(false);
    }
  }

  return (
    <>
      {/* View modal — opens with employeeUuid so the modal fetches full record */}
      <AddEmployeeModal
        open={viewOpen}
        onOpenChange={setViewOpen}
        employeeUuid={emp.id}
        initialValues={initialValues}
        onSaved={() => { setViewOpen(false); onRefresh?.(); }}
      />

      {/* Edit modal */}
      <AddEmployeeModal
        open={editOpen}
        onOpenChange={setEditOpen}
        employeeUuid={emp.id}
        initialValues={initialValues}
        onSaved={() => { setEditOpen(false); onRefresh?.(); }}
      />

      <div className="flex items-center gap-1">
        <button
          className="rounded-md p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground"
          title="View employee"
          onClick={() => setViewOpen(true)}
        >
          <LuEye className="size-4" />
        </button>
        <button
          className="rounded-md p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground"
          title="Edit employee"
          onClick={() => setEditOpen(true)}
        >
          <LuPen className="size-4" />
        </button>
        <button
          className="rounded-md p-1.5 text-destructive hover:bg-destructive/10 disabled:opacity-50"
          title="Delete employee"
          onClick={handleDelete}
          disabled={deleting}
        >
          {deleting ? <LuLoader className="size-4 animate-spin" /> : <LuTrash2 className="size-4" />}
        </button>
      </div>
    </>
  );
}

// ─── Column factory — accepts onRefresh so table can pass it through ──────────

export function buildEmployeeColumns(onRefresh?: () => void): ColumnDef<Employee>[] {
  return [
    {
      id: "select",
      header: ({ table }) => (
        <input
          type="checkbox"
          checked={table.getIsAllPageRowsSelected()}
          onChange={table.getToggleAllPageRowsSelectedHandler()}
          className="h-4 w-4 rounded border-gray-300 cursor-pointer"
          title="Select all"
        />
      ),
      cell: ({ row }) => (
        <input
          type="checkbox"
          checked={row.getIsSelected()}
          onChange={row.getToggleSelectedHandler()}
          onClick={(e) => e.stopPropagation()}
          className="h-4 w-4 rounded border-gray-300 cursor-pointer"
          title="Select row"
        />
      ),
      enableSorting: false,
      enableHiding: false,
    },
    {
      id: "srNo",
      header: "Sr. No.",
      cell: ({ row }) => row.index + 1,
      enableSorting: false,
    },
    {
      accessorKey: "name",
      id: "name",
      header: ({ column }) => <SortHeader title="Name" column={column} />,
      cell: ({ row }) => (
        <div className="font-medium text-gray-900">{row.original.name}</div>
      ),
    },
    {
      accessorKey: "employeeID",
      id: "employeeID",
      header: ({ column }) => <SortHeader title="Employee ID" column={column} />,
    },
    {
      accessorKey: "department",
      id: "department",
      header: ({ column }) => <SortHeader title="Department" column={column} />,
    },
    {
      accessorKey: "insuranceCategory",
      id: "insuranceCategory",
      header: ({ column }) => <SortHeader title="Insurance Category" column={column} />,
    },
    {
      accessorKey: "assignedPlan",
      id: "assignedPlan",
      header: ({ column }) => <SortHeader title="Assigned Plan" column={column} />,
    },
    {
      accessorKey: "coverage",
      id: "coverage",
      header: ({ column }) => <SortHeader title="Coverage" column={column} />,
    },
    {
      accessorKey: "premiumAmount",
      id: "premiumAmount",
      header: ({ column }) => <SortHeader title="Premium Amount" column={column} />,
    },
    {
      accessorKey: "numberOfDependent",
      id: "numberOfDependent",
      header: ({ column }) => <SortHeader title="Dependents" column={column} />,
    },
    {
      accessorKey: "status",
      id: "status",
      header: ({ column }) => <SortHeader title="Status" column={column} />,
      cell: ({ row }) => <StatusBadge status={row.original.status} />,
    },
    {
      id: "actions",
      enableSorting: false,
      enableHiding: false,
      header: "Actions",
      cell: ({ row }) => (
        <EmployeeActionsCell emp={row.original} onRefresh={onRefresh} />
      ),
    },
  ];
}

// Keep backward-compatible static export (no refresh)
export const employeeColumns = buildEmployeeColumns();
