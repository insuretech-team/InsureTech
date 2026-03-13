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
import { LuCirclePlus, LuDownload, LuUpload } from "react-icons/lu";

import AddPurchaseOrderModal from "@/components/modals/add-purchase-order-modal";
import { Button } from "@/components/ui/button";

import { DataTablePagination } from "./data-table-pagination";
import { DataTableToolbar } from "./data-table-toolbar";

type DepartmentOption = {
  id: string;
  name: string;
};

type CatalogOption = {
  planId: string;
  productId: string;
  productName: string;
  planName: string;
  insuranceCategory: string;
  premiumAmount: string;
};

type PurchaseOrderFormInput = {
  departmentId:       string;
  planId:             string;
  insuranceCategory:  string;
  employeeCount:      number;
  numberOfDependents: number;
  coverageAmount:     number;
  notes:              string;
};

type DataTableProps<TData, TValue> = {
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
  loading?: boolean;
  departments: DepartmentOption[];
  catalog: CatalogOption[];
  submitting?: boolean;
  onCreatePurchaseOrder: (payload: PurchaseOrderFormInput) => Promise<boolean>;
  maxHeightClassName?: string;
};

export function DataTable<TData, TValue>({
  columns,
  data,
  loading = false,
  departments,
  catalog,
  submitting = false,
  onCreatePurchaseOrder,
  maxHeightClassName = "max-h-screen",
}: DataTableProps<TData, TValue>) {
  const [sorting, setSorting] = React.useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>([]);
  const [globalFilter, setGlobalFilter] = React.useState("");
  const [columnVisibility, setColumnVisibility] = React.useState<VisibilityState>({});
  const [rowSelection, setRowSelection] = React.useState<RowSelectionState>({});
  const [addPurchaseOrderModalOpen, setAddPurchaseOrderModalOpen] = React.useState(false);

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

  return (
    <>
      <div className="portal-panel">
        <div className="border-b px-4 py-3 flex items-center justify-between gap-3">
          <div className="text-lg font-semibold text-foreground">Purchase Order Table</div>
          <div className="flex items-center gap-2">
            <Button variant="outline" className="brand-btn-gradient">
              <LuUpload />
              <span>Upload (Excel, CSV)</span>
            </Button>
            <Button variant="outline" className="brand-btn-ghost">
              <LuDownload />
              <span>Export (Excel, Pdf, CSV)</span>
            </Button>
            <Button
              variant="outline"
              className="brand-btn-gradient"
              onClick={() => setAddPurchaseOrderModalOpen(true)}
              disabled={departments.length === 0 || catalog.length === 0}
            >
              <LuCirclePlus />
              <span>Create Purchase Order</span>
            </Button>
          </div>
        </div>

        <div className="px-4 py-3 border-b">
          <DataTableToolbar table={table} />
        </div>

        <div className={["overflow-auto", maxHeightClassName].join(" ")}>
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
                  <tr key={row.id} className={["border-b last:border-b-0 table-row-hover", row.getIsSelected() ? "bg-primary/10" : ""].join(" ")}>
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
                    No purchase orders found.
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
      <AddPurchaseOrderModal
        open={addPurchaseOrderModalOpen}
        onOpenChange={setAddPurchaseOrderModalOpen}
        departments={departments}
        catalog={catalog}
        submitting={submitting}
        onSubmit={onCreatePurchaseOrder}
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
