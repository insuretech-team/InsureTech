import React from "react";
import { LuDownload } from "react-icons/lu";

const formatMoney = (n: number) =>
  new Intl.NumberFormat("bn-BD", {
    style: "currency",
    currency: "BDT",
    maximumFractionDigits: 0,
  }).format(n);

const statusStyles = (status: string) => {
  const s = (status || "").toLowerCase();
  if (s === "approved") {
    return "bg-green-100 text-green-700 ring-1 ring-green-200";
  }
  if (s === "pending") {
    return "bg-purple-100 text-purple-700 ring-1 ring-purple-200";
  }
  if (s === "rejected") {
    return "bg-red-100 text-red-700 ring-1 ring-red-200";
  }
  return "bg-gray-100 text-gray-700 ring-1 ring-gray-200";
};

export interface Invoice {
  id: string;
  dueDate: string;
  status: string;
  amount: number;
}

interface InvoiceCardProps {
  invoice: Invoice;
}

const InvoiceCard = ({ invoice }: InvoiceCardProps) => {
  const { id, dueDate, status, amount } = invoice || {};

  return (
    <div className="rounded-md border border-gray-200 bg-white p-3">
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0">
          <p className="text-sm font-semibold text-gray-800">
            Invoice ID: <span className="font-semibold">{id}</span>
          </p>
          <p className="mt-1 text-sm text-gray-500">Due Date: {dueDate}</p>

          <div className="mt-2">
            <span
              className={[
                "inline-flex items-center rounded-full px-2 py-0.5 text-sm font-medium",
                statusStyles(status),
              ].join(" ")}
            >
              {status}
            </span>
          </div>
        </div>

        <div className="flex shrink-0 flex-col items-end gap-2">
          <p className="text-sm font-semibold text-[#8C34C7]">
            {formatMoney(amount)}
          </p>

          <button
            type="button"
            className="inline-flex items-center gap-1.5 rounded-md bg-[#8C34C7] px-2.5 py-1.5 text-sm font-semibold text-white hover:bg-purple-700"
          >
            <LuDownload className="text-sm" />
            Download
          </button>
        </div>
      </div>
    </div>
  );
};

export default InvoiceCard;
