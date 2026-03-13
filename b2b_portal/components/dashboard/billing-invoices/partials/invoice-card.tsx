import Image from "next/image";
import { LuArrowRight } from "react-icons/lu";

const statusStyles = (status: string) => {
  const s = (status || "").toLowerCase();

  if (s === "paid" || s === "approved") {
    return "bg-green-100 text-green-700 ring-1 ring-green-200";
  }
  if (s === "pending") {
    return "bg-blue-100 text-blue-700 ring-1 ring-blue-200";
  }
  if (s === "failed" || s === "rejected") {
    return "bg-red-100 text-red-700 ring-1 ring-red-200";
  }
  if (s === "canceled" || s === "cancelled") {
    return "bg-orange-100 text-orange-700 ring-1 ring-orange-200";
  }
  return "bg-gray-100 text-gray-700 ring-1 ring-gray-200";
};

export interface Invoice {
  id: string;
  dueDate: string;
  status: string;
  amount: number;
  // optional: if your data has it
  paymentDate?: string;
}

interface InvoiceCardProps {
  invoice: Invoice;
}

const InvoiceCard = ({ invoice }: InvoiceCardProps) => {
  const { id, dueDate, status, amount, paymentDate } = invoice || {};
  const isPending = (status || "").toLowerCase() === "pending";

  return (
    <div className="flex items-start justify-between gap-4 my-2 p-2 border border-gray-200 rounded-md">
      {/* Left */}
      <div className="min-w-0">
        <p className="text-xs font-semibold text-gray-900">
          Invoice ID: <span className="font-semibold">{id}</span>
        </p>

        <p className="mt-1 text-[11px] text-gray-500">
          {paymentDate ? (
            <>
              Payment Date: <span className="text-gray-600">{paymentDate}</span>
            </>
          ) : (
            <>
              Due Date: <span className="text-gray-600">{dueDate}</span>
            </>
          )}
        </p>

        <div className="mt-2">
          <span className="text-sm me-2">Status:</span>
          <span
            className={[
              "inline-flex items-center rounded-sm px-2 py-0.5 text-[11px] font-medium",
              statusStyles(status),
            ].join(" ")}
          >
            {status}
          </span>
        </div>
      </div>

      {/* Right */}
      <div className="flex shrink-0 flex-col items-end gap-2">
        <p className="text-xs font-semibold text-[var(--brand-jungle)]">
          ৳{amount}
        </p>

        {isPending ? (
          <button
            type="button"
            className="inline-flex items-center gap-1.5 rounded-md bg-[var(--brand-jungle)] px-3 py-1.5 text-[11px] font-semibold text-white"
          >
            <Image
              src="./insurance-plans/credit-card.svg"
              width={16}
              height={16}
              alt="credit card"
            />
            Pay Now
          </button>
        ) : null}
      </div>
    </div>
  );
};

export default InvoiceCard;
