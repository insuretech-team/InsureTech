import DashboardLayout from "../dashboard-layout";
import InvoiceCard from "./partials/invoice-card";
import PaymentCard from "./partials/payment-card";
import { invoices } from "@/lib/invoices";
import { payments } from "@/lib/payments";
import { LuSearch } from "react-icons/lu";

const BillingInvoices = () => {
  return (
    <DashboardLayout>
      <div className="mt-3 grid grid-cols-1 gap-4 lg:grid-cols-2">
        {/* Recent Invoices */}
        <div className="rounded-lg border border-gray-200 bg-white shadow-sm">
          <div className="flex items-center justify-between gap-3 border-b border-gray-100 px-4 py-3">
            <h5 className="text-sm font-semibold text-gray-900">
              Recent Invoices
            </h5>
          </div>

          <div className="px-4">
            {invoices.map((invoice, idx) => (
              <InvoiceCard key={`${invoice.id}-${idx}`} invoice={invoice} />
            ))}
          </div>
        </div>

        {/* Payment History */}
        <div className="rounded-lg border border-gray-200 bg-white shadow-sm">
          <div className="flex items-center justify-between gap-3 border-b border-gray-100 px-4 py-3">
            <h5 className="text-sm font-semibold text-gray-900">
              Payment History
            </h5>

            <div className="relative w-[220px] max-w-full">
              <LuSearch className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-sm text-gray-400" />
              <input
                className="h-8 w-full rounded-md border border-gray-200 bg-white pl-9 pr-3 text-xs text-gray-700 placeholder:text-gray-400 focus:border-purple-400 focus:outline-none focus:ring-2 focus:ring-purple-100"
                placeholder="Search by invoice id"
              />
            </div>
          </div>

          <div className="divide-y divide-gray-100 px-4">
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
