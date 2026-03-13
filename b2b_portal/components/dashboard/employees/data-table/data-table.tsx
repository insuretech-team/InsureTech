"use client";

import * as React from "react";
import {
  ColumnDef,
  SortingState,
  ColumnFiltersState,
  VisibilityState,
  RowSelectionState,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getSortedRowModel,
  getPaginationRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { DataTableToolbar } from "./data-table-toolbar";
import { DataTablePagination } from "./data-table-pagination";
import { Button } from "../../../ui/button";
import { LuCirclePlus, LuUpload, LuDownload, LuTrash2, LuLoader } from "react-icons/lu";
import AddEmployeeModal from "../../../modals/add-employee-modal";
import BulkUploadEmployeeModal from "../../../modals/bulk-upload-employee-modal";
import { EmployeeCard } from "../employee-card";
import type { Employee } from "@lib/types/b2b";

type DataTableProps<TData, TValue> = {
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
  loading?: boolean;
  /** height of scroll area (enables sticky header) */
  maxHeightClassName?: string; // e.g. "max-h-[420px]"
  /** Called by Add button after successful create to reload the list */
  onRefresh?: () => void;
  organisationId?: string;
};

export function DataTable<TData, TValue>({
  columns,
  data,
  loading = false,
  maxHeightClassName = "max-h-screen",
  onRefresh,
  organisationId,
}: DataTableProps<TData, TValue>) {
  const [sorting, setSorting] = React.useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
    [],
  );
  const [globalFilter, setGlobalFilter] = React.useState("");
  const [columnVisibility, setColumnVisibility] =
    React.useState<VisibilityState>({});
  const [rowSelection, setRowSelection] = React.useState<RowSelectionState>({});

  const table = useReactTable({
    data,
    columns,
    state: {
      sorting,
      columnFilters,
      globalFilter,
      columnVisibility,
      rowSelection,
    },
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    onGlobalFilterChange: setGlobalFilter,
    onColumnVisibilityChange: setColumnVisibility,
    onRowSelectionChange: setRowSelection,

    enableMultiSort: true,
    enableRowSelection: true,

    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
  });

  const [addEmployeeModalOpen, setAddEmployeeModalOpen] = React.useState(false);
  const [bulkUploadModalOpen, setBulkUploadModalOpen] = React.useState(false);
  const [bulkDeleting, setBulkDeleting] = React.useState(false);
  const [cardEmployee, setCardEmployee] = React.useState<Employee | null>(null);
  const [cardDeleting, setCardDeleting] = React.useState(false);
  const [cardEditOpen, setCardEditOpen] = React.useState(false);

  async function handleBulkDelete() {
    const selectedRows = table.getSelectedRowModel().rows;
    if (selectedRows.length === 0) return;
    if (!confirm(`Delete ${selectedRows.length} selected employee(s)?\nThis cannot be undone.`)) return;
    setBulkDeleting(true);
    try {
      const results = await Promise.allSettled(
        selectedRows.map((row) =>
          fetch(`/api/employees/${(row.original as import('@lib/types/b2b').Employee).id}`, { method: 'DELETE' })
        )
      );
      const failed = results.filter((r) => r.status === 'rejected').length;
      if (failed > 0) alert(`${failed} employee(s) could not be deleted.`);
      table.resetRowSelection();
      onRefresh?.();
    } finally {
      setBulkDeleting(false);
    }
  }

  async function handleCardDelete() {
    if (!cardEmployee) return;
    if (!confirm(`Delete employee "${cardEmployee.name}"?\nThis cannot be undone.`)) return;
    setCardDeleting(true);
    try {
      const res = await fetch(`/api/employees/${cardEmployee.id}`, { method: "DELETE" });
      const payload = (await res.json()) as { ok: boolean; message?: string };
      if (!res.ok || !payload.ok) {
        alert(payload.message ?? "Failed to delete employee");
        return;
      }
      setCardEmployee(null);
      onRefresh?.();
    } catch {
      alert("Network error — could not delete employee");
    } finally {
      setCardDeleting(false);
    }
  }

  function handleRowClick(e: React.MouseEvent<HTMLTableRowElement>, employee: Employee) {
    // Don't open card if clicking on action buttons
    if ((e.target as HTMLElement).closest("button")) return;
    setCardEmployee(employee);
  }

  return (
    <>
      <div className="portal-panel">
        <div className="border-b px-4 py-3 flex items-center justify-between gap-3">
          <div className="text-lg font-semibold text-foreground">
            Employee Table
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              className="brand-btn-gradient"
              onClick={() => setBulkUploadModalOpen(true)}
              disabled={!organisationId}
              type="button"
            >
              <LuUpload />
              <span>Upload (Excel, CSV)</span>
            </Button>
            <Button
              variant="outline"
              className="brand-btn-ghost"
              onClick={() => window.open("/api/employees/template?format=csv", "_blank")}
              type="button"
            >
              <LuDownload />
              <span>Download Template</span>
            </Button>
            <Button
              variant="outline"
              className="brand-btn-gradient"
              onClick={() => setAddEmployeeModalOpen(true)}
              type="button"
              disabled={!organisationId}
            >
              <LuCirclePlus />
              <span>Add Employee</span>
            </Button>
          </div>
        </div>

        <div className="px-4 py-3 border-b">
          <DataTableToolbar table={table} />
        </div>

        {table.getSelectedRowModel().rows.length > 0 && (
          <div className="px-4 py-2 bg-red-50 border-b flex items-center gap-3">
            <span className="text-sm text-red-700 font-medium">
              {table.getSelectedRowModel().rows.length} row(s) selected
            </span>
            <Button
              variant="outline"
              size="sm"
              onClick={handleBulkDelete}
              disabled={bulkDeleting}
              className="border-red-300 text-red-600 hover:bg-red-100 h-8 px-3 text-xs"
            >
              {bulkDeleting ? (
                <span className="flex items-center gap-1.5"><LuLoader className="animate-spin size-3.5" />Deleting…</span>
              ) : (
                <span className="flex items-center gap-1.5"><LuTrash2 className="size-3.5" />Delete Selected</span>
              )}
            </Button>
            <button
              className="ml-auto text-xs text-muted-foreground hover:text-foreground"
              onClick={() => table.resetRowSelection()}
            >Clear selection</button>
          </div>
        )}

        <div className={["overflow-auto", maxHeightClassName].join(" ")}>
          <table className="w-full text-sm">
            <thead className="table-head">
              {table.getHeaderGroups().map((hg) => (
                <tr key={hg.id} className="border-b">
                  {hg.headers.map((header) => (
                    <th
                      key={header.id}
                      className="px-4 py-3 text-left font-medium text-gray-700 whitespace-nowrap"
                    >
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext(),
                          )}
                    </th>
                  ))}
                </tr>
              ))}
            </thead>

            <tbody>
              {loading ? (
                <SkeletonRows colCount={table.getAllLeafColumns().length} />
              ) : table.getRowModel().rows.length ? (
                table.getRowModel().rows.map((row) => (
                  <tr
                    key={row.id}
                    onClick={(e) => handleRowClick(e, row.original as Employee)}
                    className={[
                      "border-b last:border-b-0 table-row-hover cursor-pointer",
                      row.getIsSelected() ? "bg-primary/10" : "",
                    ].join(" ")}
                  >
                    {row.getVisibleCells().map((cell) => (
                      <td key={cell.id} className="px-4 py-3 align-middle">
                        {flexRender(
                          cell.column.columnDef.cell,
                          cell.getContext(),
                        )}
                      </td>
                    ))}
                  </tr>
                ))
              ) : (
                <tr>
                  <td
                    className="px-4 py-10 text-center text-gray-500"
                    colSpan={table.getAllLeafColumns().length}
                  >
                    No results.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        <div className="border-t px-4 py-3">
          <DataTablePagination table={table} />
        </div>
      </div>
      <AddEmployeeModal
        open={addEmployeeModalOpen}
        onOpenChange={setAddEmployeeModalOpen}
        organisationId={organisationId}
        onSaved={() => { setAddEmployeeModalOpen(false); onRefresh?.(); }}
      />
      <BulkUploadEmployeeModal
        open={bulkUploadModalOpen}
        onOpenChange={setBulkUploadModalOpen}
        organisationId={organisationId}
        onSaved={() => { setBulkUploadModalOpen(false); onRefresh?.(); }}
      />
      <EmployeeCard
        employee={cardEmployee}
        open={cardEmployee !== null}
        onClose={() => setCardEmployee(null)}
        onEdit={() => { setCardEditOpen(true); }}
        onDelete={handleCardDelete}
        deleting={cardDeleting}
      />
      <AddEmployeeModal
        open={cardEditOpen}
        onOpenChange={setCardEditOpen}
        employeeUuid={cardEmployee?.id}
        onSaved={() => { setCardEditOpen(false); setCardEmployee(null); onRefresh?.(); }}
      />
    </>
  );
}

function SkeletonRows({ colCount }: { colCount: number }) {
  const rows = Array.from({ length: 6 });
  const cols = Array.from({ length: colCount });

  return (
    <>
      {rows.map((_, r) => (
        <tr key={r} className="border-b last:border-b-0">
          {cols.map((__, c) => (
            <td key={c} className="px-4 py-3">
              <div className="h-4 w-full rounded bg-gray-100 animate-pulse" />
            </td>
          ))}
        </tr>
      ))}
    </>
  );
}

