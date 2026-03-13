import Image from "next/image";
import React from "react";
import { LuEye } from "react-icons/lu";

const formatMoney = (n: number) =>
  new Intl.NumberFormat("bn-BD", {
    style: "currency",
    currency: "BDT",
    maximumFractionDigits: 0,
  }).format(n);

export interface Payment {
  id: string;
  amount: number;
  method: string;
  month: string;
  status: string;
}

interface PaymentCardProps {
  payment: Payment;
}

const PaymentCard = ({ payment }: PaymentCardProps) => {
  const { amount, method, month, id } = payment || {};

  return (
    <div className="rounded-md border border-gray-200 bg-white p-3">
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0">
          <div className="flex items-center gap-2">
            <Image
              src="./quotations/check-circle.svg"
              width={16}
              height={16}
              alt="check"
            />
            <p className="text-sm font-semibold text-gray-900">
              {formatMoney(amount)}
            </p>
          </div>

          <div className="mt-2 space-y-1 text-sm text-gray-500">
            <p>
              <span className="font-medium text-gray-600">Payment Date:</span>{" "}
              {month}
            </p>
            <p>
              <span className="font-medium text-gray-600">Method:</span>{" "}
              {method}
            </p>
            <p className="truncate">
              <span className="font-medium text-gray-600">Transaction ID:</span>{" "}
              {id}
            </p>
          </div>
        </div>

        <div className="shrink-0">
          <button
            type="button"
            className="inline-flex items-center gap-1.5 rounded-md bg-primary px-2.5 py-1.5 text-xs font-semibold text-white hover:bg-primary/90"
          >
            <LuEye className="text-sm" />
            View
          </button>
        </div>
      </div>
    </div>
  );
};

export default PaymentCard;

