import React from "react";
import DashboardLayout from "../dashboard-layout";
import InvoiceCard from "./invoice-card";
import PaymentCard from "./payment-card";
import { invoices } from "@/lib/invoices";
import { payments } from "@/lib/payments";

const BillingInvoices = () => {
  return (
    <DashboardLayout>
      <div className="mt-3 grid grid-cols-1 gap-4 lg:grid-cols-2">
        {/* Recent Invoices */}
        <div className="rounded-md bg-white p-4 shadow-sm">
          <div className="flex items-center gap-2">
            <h5 className="text-lg font-semibold text-gray-800">
              Recent Invoices
            </h5>
          </div>

          <div className="mt-4 space-y-3">
            {invoices.map((invoice, idx) => (
              <InvoiceCard key={`${invoice.id}-${idx}`} invoice={invoice} />
            ))}
          </div>
        </div>

        {/* Payment History */}
        <div className="rounded-md bg-white p-4 shadow-sm">
          <div className="flex items-center gap-2">
            <h5 className="text-lg font-semibold text-gray-800">
              Payment History
            </h5>
          </div>

          <div className="mt-4 space-y-3">
            {payments.map((payment, idx) => (
              <PaymentCard key={`${payment.id}-${idx}`} payment={payment} />
            ))}
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
};

export default BillingInvoices;
