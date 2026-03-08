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
import { LuCirclePlus, LuUpload, LuDownload } from "react-icons/lu";
import AddEmployeeModal from "../../../modals/add-employee-modal";

type DataTableProps<TData, TValue> = {
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
  loading?: boolean;
  /** height of scroll area (enables sticky header) */
  maxHeightClassName?: string; // e.g. "max-h-[420px]"
};

export function DataTable<TData, TValue>({
  columns,
  data,
  loading = false,
  maxHeightClassName = "max-h-screen",
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

  return (
    <>
      <div className="rounded-lg border bg-white overflow-hidden">
        <div className="border-b px-4 py-3 flex items-center justify-between gap-3">
          <div className="text-lg font-semibold text-[#242424]">
            Employee Table
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              className="text-[#FFFFFF] bg-gradient-to-r from-[#8C34C7] to-[#702A9F]"
            >
              <LuUpload />
              <span>Upload (Excel, CSV)</span>
            </Button>
            <Button
              variant="outline"
              className="text-[#2b2b2b] hover:text-[#8C34C7]"
            >
              <LuDownload />
              <span>Export (Excel, Pdf, CSV)</span>
            </Button>
            <Button
              variant="outline"
              className="text-[#FFFFFF] bg-gradient-to-r from-[#8C34C7] to-[#702A9F]"
              onClick={() => setAddEmployeeModalOpen(true)}
            >
              <LuCirclePlus />
              <span>Add Employee</span>
            </Button>
          </div>
        </div>

        <div className="px-4 py-3 border-b">
          <DataTableToolbar table={table} />
        </div>

        <div className={["overflow-auto", maxHeightClassName].join(" ")}>
          <table className="w-full text-sm">
            <thead className="sticky top-0 z-10 bg-gray-50">
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
                    className={[
                      "border-b last:border-b-0 hover:bg-gray-50 transition-colors",
                      row.getIsSelected() ? "bg-purple-50/40" : "",
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
