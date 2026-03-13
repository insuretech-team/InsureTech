import Image from "next/image";
import React, { useState } from "react";
import { LuDownload, LuEye } from "react-icons/lu";
import InvoiceDetailModal from "@/components/modals/invoice-detail-modal";
export interface Payment {
  id: string;
  amount: number;
  method: string;
  month: string; // used as "Payment Date" in UI
  status: string;
}

interface PaymentCardProps {
  payment: Payment;
}

const PaymentCard = ({ payment }: PaymentCardProps) => {
  const { amount, method, month, id } = payment || {};
  const [isInvoiceModalOpen, setIsInvoiceModalOpen] = useState(false);

  return (
    <>
      <div className="flex items-start justify-between gap-4 my-2 p-2 border border-gray-200 rounded-md">
        {/* Left */}
        <div className="min-w-0">
          <div className="flex items-center gap-2">
            <Image
              src="./quotations/check-circle.svg"
              width={16}
              height={16}
              alt="check"
            />
            <p className="text-xs font-semibold text-gray-900">৳{amount}</p>
          </div>

          <div className="mt-2 space-y-1 text-[11px] text-gray-500">
            <p>
              <span className="font-medium text-gray-600">Payment Date:</span>{" "}
              {month}
            </p>
            <p>
              <span className="font-medium text-gray-600">Method:</span>{" "}
              {method}
            </p>
            <p className="truncate max-w-[260px]">
              <span className="font-medium text-gray-600">Transaction ID:</span>{" "}
              {id}
            </p>
          </div>
        </div>

        {/* Right */}
        <div className="flex shrink-0 items-center gap-2">
          <button
            type="button"
            className="inline-flex items-center gap-1.5 rounded-md border border-purple-200 bg-white px-3 py-1.5 text-[11px] font-semibold text-[var(--primary-deep)] hover:bg-purple-50"
          >
            <LuDownload className="text-sm" />
            Download Invoice
          </button>

          <button
            type="button"
            onClick={() => setIsInvoiceModalOpen(true)}
            className="inline-flex items-center gap-1.5 rounded-md bg-[var(--brand-jungle)] px-3 py-1.5 text-[11px] font-semibold text-white"
          >
            <LuEye className="text-sm" />
            View
          </button>
        </div>
      </div>
      {isInvoiceModalOpen && (
        <InvoiceDetailModal
          open={isInvoiceModalOpen}
          onOpenChange={setIsInvoiceModalOpen}
        />
      )}
    </>
  );
};

export default PaymentCard;
