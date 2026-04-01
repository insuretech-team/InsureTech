"use client";

import { ColumnDef } from "@tanstack/react-table";
import { LuPen, LuTrash2 } from "react-icons/lu";

export type Quotation = {
  id: string;
  quotationID: string;
  insurerName: string;
  plan: string;
  insuranceCategory: string;
  department: string;
  employeeNo: number;
  estimatedPremium: string;
  quotedAmount: string;
  status: string;
  submissionDate: string;
  validUntil: string;
};

function SortHeader({ title, column }: { title: string; column: any }) {
  const dir = column.getIsSorted(); // false | "asc" | "desc"
  return (
    <button
      className="inline-flex items-center gap-2 select-none"
      onClick={() => column.toggleSorting(dir === "asc")}
      title="Sort"
      type="button"
    >
      {title}
      <span className="text-gray-400 text-xs">
        {dir === "asc" ? "▲" : dir === "desc" ? "▼" : "↕"}
      </span>
    </button>
  );
}

export const quotationColumns: ColumnDef<Quotation>[] = [
  {
    id: "srNo",
    header: "Sr. No.",
    cell: ({ row }) => row.index + 1,
    enableSorting: false,
  },
  {
    accessorKey: "quotationID",
    header: ({ column }) => <SortHeader title="Quotation ID" column={column} />,
    cell: ({ row }) => (
      <div className="font-medium text-gray-900">
        {row.original.quotationID}
      </div>
    ),
  },
  {
    accessorKey: "insurerName",
    header: ({ column }) => <SortHeader title="Insurer" column={column} />,
  },
  {
    accessorKey: "plan",
    header: ({ column }) => <SortHeader title="Plan" column={column} />,
  },
  {
    accessorKey: "insuranceCategory",
    header: ({ column }) => <SortHeader title="Category" column={column} />,
  },
  {
    accessorKey: "department",
    header: ({ column }) => <SortHeader title="Department" column={column} />,
  },
  {
    accessorKey: "employeeNo",
    header: ({ column }) => (
      <SortHeader title="No. of Employee" column={column} />
    ),
  },
  {
    accessorKey: "estimatedPremium",
    header: ({ column }) => (
      <SortHeader title="Estimated Premium" column={column} />
    ),
  },
  {
    accessorKey: "quotedAmount",
    header: ({ column }) => (
      <SortHeader title="Quoted Amount" column={column} />
    ),
  },

  {
    accessorKey: "status",
    header: ({ column }) => <SortHeader title="Status" column={column} />,
    cell: ({ row }) => {
      const s = row.original.status;

      const statusStyles: Record<string, string> = {
        Approved: "bg-green-50 text-green-700",
        Submitted: "bg-blue-50 text-blue-700",
        Received: "bg-purple-50 text-purple-700",
        "In Draft": "bg-yellow-50 text-yellow-700",
        Rejected: "bg-red-50 text-red-700",
      };

      return (
        <span
          className={`inline-flex items-center rounded-full px-2 py-1 text-xs font-medium ${
            statusStyles[s] || "bg-gray-50 text-gray-700"
          }`}
        >
          {s}
        </span>
      );
    },
  },
  {
    accessorKey: "submissionDate",
    header: ({ column }) => (
      <SortHeader title="Submission Date" column={column} />
    ),
  },
  {
    accessorKey: "validUntil",
    header: ({ column }) => <SortHeader title="Valid Until" column={column} />,
  },

  // Actions
  {
    id: "actions",
    enableSorting: false,
    enableHiding: false,
    header: "Actions",
    cell: ({ row }) => {
      const quotation = row.original;

      return (
        <details className="relative">
          <summary className="list-none cursor-pointer rounded-md border px-2 py-1 text-xs hover:bg-gray-50 w-max">
            •••
          </summary>

          <div className="absolute right-0 mt-2 w-40 rounded-md border bg-white shadow-sm z-20 overflow-hidden">
            <button
              className="w-full px-3 py-2 text-left text-sm hover:bg-gray-50"
              type="button"
              onClick={() => alert(`Edit: ${quotation.quotationID}`)}
            >
              <LuPen className="inline mr-2" />
              Edit
            </button>

            <button
              className="w-full px-3 py-2 text-left text-sm hover:bg-gray-50"
              type="button"
              onClick={() => alert(`Delete: ${quotation.quotationID}`)}
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
