"use client";

import { ColumnDef } from "@tanstack/react-table";
import { LuPen, LuEye, LuTrash2 } from "react-icons/lu";

export type Employee = {
  id: string;
  name: string;
  employeeID: string;
  department: string;
  insuranceCategory?: string;
  assignedPlan?: string;
  coverage: string;
  premiumAmount: string;
  status: "Active" | "Inactive";
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

export const employeeColumns: ColumnDef<Employee>[] = [
  {
    accessorKey: "Sr. No.",
    id: "srNo",
    header: "Sr. No.",
    cell: ({ row }) => row.index + 1,
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
    header: ({ column }) => <SortHeader title="Department" column={column} />,
  },
  {
    accessorKey: "insuranceCategory",
    id: "insuranceCategory",
    header: ({ column }) => (
      <SortHeader title="Insurance Category" column={column} />
    ),
  },
  {
    accessorKey: "assignedPlan",
    id: "assignedPlan",
    header: ({ column }) => (
      <SortHeader title="Assigned Plan" column={column} />
    ),
  },
  {
    accessorKey: "coverage",
    id: "coverage",
    header: ({ column }) => <SortHeader title="Coverage" column={column} />,
  },
  {
    accessorKey: "premiumAmount",
    id: "premiumAmount",
    header: ({ column }) => (
      <SortHeader title="Premium Amount" column={column} />
    ),
  },
  {
    accessorKey: "status",
    header: ({ column }) => <SortHeader title="Status" column={column} />,
    cell: ({ row }) => {
      const s = row.original.status;
      return (
        <span
          className={[
            "inline-flex items-center rounded-full px-2 py-1 text-xs font-medium",
            s === "Active"
              ? "bg-green-50 text-green-700"
              : "bg-red-50 text-red-700",
          ].join(" ")}
        >
          {s}
        </span>
      );
    },
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
              onClick={() => alert(`View: ${emp.name}`)}
            >
              <LuEye className="inline mr-2" />
              View
            </button>
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
