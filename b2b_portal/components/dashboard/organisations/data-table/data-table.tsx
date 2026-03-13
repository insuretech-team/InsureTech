/**
 * Organisation data table
 */
"use client";

import * as React from "react";
import {
  ColumnDef,
  SortingState,
  ColumnFiltersState,
  VisibilityState,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getSortedRowModel,
  getPaginationRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { LuCirclePlus } from "react-icons/lu";
import AddOrganisationModal from "@/components/modals/add-organisation-modal";

type Props<TData, TValue> = {
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
  loading?: boolean;
  onRefresh?: () => void;
};

function SkeletonRows({ colCount }: { colCount: number }) {
  return (
    <>
      {Array.from({ length: 5 }).map((_, r) => (
        <tr key={r} className="border-b last:border-b-0">
          {Array.from({ length: colCount }).map((__, c) => (
            <td key={c} className="px-4 py-3">
              <div className="h-4 w-full rounded bg-gray-100 animate-pulse" />
            </td>
          ))}
        </tr>
      ))}
    </>
  );
}

export function OrganisationDataTable<TData, TValue>({ columns, data, loading = false, onRefresh }: Props<TData, TValue>) {
  const [sorting, setSorting] = React.useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>([]);
  const [globalFilter, setGlobalFilter] = React.useState("");
  const [columnVisibility, setColumnVisibility] = React.useState<VisibilityState>({});
  const [addOpen, setAddOpen] = React.useState(false);

  const table = useReactTable({
    data,
    columns,
    state: { sorting, columnFilters, globalFilter, columnVisibility },
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    onGlobalFilterChange: setGlobalFilter,
    onColumnVisibilityChange: setColumnVisibility,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
  });

  return (
    <>
      <div className="portal-panel">
        <div className="border-b px-4 py-3 flex items-center justify-between gap-3">
          <div className="text-lg font-semibold text-foreground">Organisations</div>
          <div className="flex items-center gap-2">
            <input
              placeholder="Search organisations…"
              value={globalFilter}
              onChange={(e) => setGlobalFilter(e.target.value)}
              className="border rounded-md px-3 py-1.5 text-sm w-56 focus:outline-none focus:ring-2 focus:ring-primary"
            />
            <Button variant="outline" className="brand-btn-gradient" onClick={() => setAddOpen(true)} type="button">
              <LuCirclePlus />
              <span>Add Organisation</span>
            </Button>
          </div>
        </div>

        <div className="overflow-auto max-h-screen">
          <table className="w-full text-sm">
            <thead className="table-head">
              {table.getHeaderGroups().map((hg) => (
                <tr key={hg.id} className="border-b">
                  {hg.headers.map((header) => (
                    <th key={header.id} className="px-4 py-3 text-left font-medium text-gray-700 whitespace-nowrap">
                      {header.isPlaceholder ? null : flexRender(header.column.columnDef.header, header.getContext())}
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
                  <tr key={row.id} className="border-b last:border-b-0 table-row-hover">
                    {row.getVisibleCells().map((cell) => (
                      <td key={cell.id} className="px-4 py-3 align-middle">
                        {flexRender(cell.column.columnDef.cell, cell.getContext())}
                      </td>
                    ))}
                  </tr>
                ))
              ) : (
                <tr>
                  <td className="px-4 py-10 text-center text-gray-500" colSpan={table.getAllLeafColumns().length}>
                    No organisations found.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>

      <AddOrganisationModal
        open={addOpen}
        onOpenChange={setAddOpen}
        onSaved={() => { setAddOpen(false); onRefresh?.(); }}
      />
    </>
  );
}
