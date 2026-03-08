"use client";

import * as React from "react";
import { Table } from "@tanstack/react-table";

function cx(...s: Array<string | false | undefined | null>) {
  return s.filter(Boolean).join(" ");
}

export function DataTableToolbar<TData>({ table }: { table: Table<TData> }) {
  const [q, setQ] = React.useState(String(table.getState().globalFilter ?? ""));

  // debounce global search
  React.useEffect(() => {
    const t = setTimeout(() => table.setGlobalFilter(q), 250);
    return () => clearTimeout(t);
  }, [q, table]);

  return (
    <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div className="flex flex-1 flex-col gap-2 sm:flex-row sm:items-center">
        <span className="text-gray-500">List Show:</span>
        <select
          className="h-8 w-14 rounded-md border px-2 text-sm"
          value={table.getState().pagination.pageSize}
          onChange={(e) => table.setPageSize(Number(e.target.value))}
        >
          {[5, 10, 20, 50].map((n) => (
            <option key={n} value={n}>
              {n}
            </option>
          ))}
        </select>
      </div>

      {/* Column visibility */}
      <div className="flex items-center gap-2">
        {/* Global search */}
        <input
          value={q}
          onChange={(e) => setQ(e.target.value)}
          placeholder="Search..."
          className="h-9 w-full sm:w-72 rounded-md border px-3 text-sm outline-none focus:ring-2 focus:ring-purple-200"
        />
        <details className="relative">
          <summary className="h-9 cursor-pointer list-none rounded-md border px-3 text-sm flex items-center gap-2 hover:bg-gray-50">
            Columns
            <span className="text-gray-400">▾</span>
          </summary>

          <div className="absolute right-0 mt-2 w-56 rounded-md border bg-white shadow-sm p-2 z-20">
            <div className="text-xs font-semibold text-gray-700 px-2 py-1">
              Toggle columns
            </div>
            <div className="max-h-56 overflow-auto">
              {table
                .getAllLeafColumns()
                .filter((c) => c.getCanHide())
                .map((column) => (
                  <label
                    key={column.id}
                    className={cx(
                      "flex items-center gap-2 px-2 py-2 rounded hover:bg-gray-50 cursor-pointer",
                    )}
                  >
                    <input
                      type="checkbox"
                      checked={column.getIsVisible()}
                      onChange={column.getToggleVisibilityHandler()}
                    />
                    <span className="text-sm capitalize">{column.id}</span>
                  </label>
                ))}
            </div>
          </div>
        </details>
      </div>
    </div>
  );
}
