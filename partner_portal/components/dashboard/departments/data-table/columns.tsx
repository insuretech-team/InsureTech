"use client";

import { ColumnDef } from "@tanstack/react-table";
import { LuPen, LuTrash2 } from "react-icons/lu";

export type Department = {
  id: string;
  name: string;
  employeeNo: number;
  totalPremium: string;
};

function SortHeader({ title, column }: { title: string; column: any }) {
  const dir = column.getIsSorted(); // false | "asc" | "desc"
  return (
    <button
      className="inline-flex items-center gap-2 select-none"
      onClick={() => column.toggleSorting(dir === "asc")}
      title="Sort"
    >
      {title}
      <span className="text-gray-400 text-xs">
        {dir === "asc" ? "▲" : dir === "desc" ? "▼" : "↕"}
      </span>
    </button>
  );
}

export const departmentColumns: ColumnDef<Department>[] = [
  {
    accessorKey: "Sr. No.",
    id: "srNo",
    header: "Sr. No.",
    cell: ({ row }) => row.index + 1,
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
    header: ({ column }) => (
      <SortHeader title="No. of Employee" column={column} />
    ),
  },
  {
    accessorKey: "totalPremium",
    header: ({ column }) => (
      <SortHeader title="Total Premium" column={column} />
    ),
  },

  // Actions dropdown
  {
    id: "actions",
    enableSorting: false,
    enableHiding: false,
    header: "Actions",
    cell: ({ row }) => {
      const emp = row.original;

      return (
        <details className="relative">
          <summary className="list-none cursor-pointer rounded-md border px-2 py-1 text-xs hover:bg-gray-50 w-max">
            •••
          </summary>
          <div className="absolute right-0 mt-2 w-40 rounded-md border bg-white shadow-sm z-20 overflow-hidden">
            <button
              className="w-full px-3 py-2 text-left text-sm hover:bg-gray-50"
              onClick={() => alert(`Edit: ${emp.name}`)}
            >
              <LuPen className="inline mr-2" />
              Edit
            </button>
            <button
              className="w-full px-3 py-2 text-left text-sm hover:bg-gray-50"
              onClick={() => alert(`Delete: ${emp.name}`)}
            >
              <LuTrash2 className="inline mr-2" />
              Delete
            </button>
          </div>
        </details>
      );
    },
  },
];
