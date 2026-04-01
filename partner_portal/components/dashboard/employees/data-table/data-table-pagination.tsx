"use client";

import { Table } from "@tanstack/react-table";

export function DataTablePagination<TData>({ table }: { table: Table<TData> }) {
  return (
    <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div className="text-xs text-gray-600">
        Page {table.getState().pagination.pageIndex + 1} of{" "}
        {table.getPageCount()}
        {" • "}
        {table.getFilteredRowModel().rows.length} rows
      </div>

      <div className="flex items-center gap-2">
        <button
          className="h-9 rounded-md border px-3 text-sm disabled:opacity-50"
          onClick={() => table.setPageIndex(0)}
          disabled={!table.getCanPreviousPage()}
        >
          {"<<"}
        </button>

        <button
          className="h-9 rounded-md border px-3 text-sm disabled:opacity-50"
          onClick={() => table.previousPage()}
          disabled={!table.getCanPreviousPage()}
        >
          Prev
        </button>

        <button
          className="h-9 rounded-md border px-3 text-sm disabled:opacity-50"
          onClick={() => table.nextPage()}
          disabled={!table.getCanNextPage()}
        >
          Next
        </button>

        <button
          className="h-9 rounded-md border px-3 text-sm disabled:opacity-50"
          onClick={() => table.setPageIndex(table.getPageCount() - 1)}
          disabled={!table.getCanNextPage()}
        >
          {">>"}
        </button>
      </div>
    </div>
  );
}
