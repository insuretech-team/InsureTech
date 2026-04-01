"use client";

import { useEffect, useState } from "react";
import DashboardLayout from "../dashboard-layout";
import { bffClient } from "@lib/sdk/b2b-sdk-client";
import { LuSearch, LuFileText, LuCreditCard, LuLoader } from "react-icons/lu";

type InvoiceRecord = Record<string, unknown>;
type PaymentRecord = Record<string, unknown>;

function StatusBadge({ status }: { status: string }) {
  const s = String(status ?? "").replace(/INVOICE_STATUS_|PAYMENT_STATUS_/g, "").replace(/_/g, " ");
  const color =
    s === "PENDING" ? "bg-yellow-100 text-yellow-700" :
    s === "PAID" || s === "SUCCESS" || s === "SETTLED" ? "bg-green-100 text-green-700" :
    s === "CANCELLED" || s === "FAILED" ? "bg-red-100 text-red-700" :
    "bg-gray-100 text-gray-600";
  return <span className={`rounded-full px-2 py-0.5 text-[10px] font-medium uppercase ${color}`}>{s || "—"}</span>;
}

const BillingInvoices = () => {
  const [invoices, setInvoices] = useState<InvoiceRecord[]>([]);
  const [payments, setPayments] = useState<PaymentRecord[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");

  useEffect(() => {
    setLoading(true);
    Promise.all([
      bffClient.billing.listInvoices({ pageSize: 20 }),
      bffClient.billing.listPayments({ pageSize: 20 }),
    ]).then(([invRes, payRes]) => {
      setInvoices((invRes.invoices ?? []) as InvoiceRecord[]);
      setPayments((payRes.payments ?? []) as PaymentRecord[]);
    }).catch(() => {}).finally(() => setLoading(false));
  }, []);

  const filteredPayments = search
    ? payments.filter((p) => String(p.invoice_id ?? p.payment_id ?? "").toLowerCase().includes(search.toLowerCase()))
    : payments;

  return (
    <DashboardLayout>
      <div className="mt-3 grid grid-cols-1 gap-4 lg:grid-cols-2">
        {/* Recent Invoices */}
        <div className="rounded-lg border border-gray-200 bg-white shadow-sm">
          <div className="flex items-center justify-between gap-3 border-b border-gray-100 px-4 py-3">
            <h5 className="text-sm font-semibold text-gray-900">Recent Invoices</h5>
          </div>
          <div className="divide-y divide-gray-100 px-4">
            {loading ? (
              <div className="flex items-center justify-center py-10 text-gray-400">
                <LuLoader className="animate-spin mr-2" /> Loading…
              </div>
            ) : invoices.length === 0 ? (
              <p className="py-8 text-center text-sm text-gray-400">No invoices found</p>
            ) : invoices.map((inv, idx) => (
              <div key={String(inv.invoice_id ?? idx)} className="flex items-center justify-between py-3 gap-2">
                <div className="flex items-center gap-3">
                  <div className="flex h-8 w-8 items-center justify-center rounded-full bg-purple-50">
                    <LuFileText className="text-purple-500 text-sm" />
                  </div>
                  <div>
                    <p className="text-xs font-medium text-gray-800">{String(inv.invoice_number ?? inv.invoice_id ?? "—")}</p>
                    <p className="text-[10px] text-gray-400">{String(inv.issued_at ?? inv.created_at ?? "")}</p>
                  </div>
                </div>
                <div className="flex flex-col items-end gap-1">
                  <p className="text-xs font-semibold text-gray-800">
                    {inv.amount ? `BDT ${Number((inv.amount as Record<string,unknown>)?.decimal_amount ?? inv.amount).toLocaleString()}` : "—"}
                  </p>
                  <StatusBadge status={String(inv.status ?? "")} />
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Payment History */}
        <div className="rounded-lg border border-gray-200 bg-white shadow-sm">
          <div className="flex items-center justify-between gap-3 border-b border-gray-100 px-4 py-3">
            <h5 className="text-sm font-semibold text-gray-900">Payment History</h5>
            <div className="relative w-[220px] max-w-full">
              <LuSearch className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-sm text-gray-400" />
              <input
                className="h-8 w-full rounded-md border border-gray-200 bg-white pl-9 pr-3 text-xs text-gray-700 placeholder:text-gray-400 focus:border-purple-400 focus:outline-none focus:ring-2 focus:ring-purple-100"
                placeholder="Search by invoice id"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
            </div>
          </div>
          <div className="divide-y divide-gray-100 px-4">
            {loading ? (
              <div className="flex items-center justify-center py-10 text-gray-400">
                <LuLoader className="animate-spin mr-2" /> Loading…
              </div>
            ) : filteredPayments.length === 0 ? (
              <p className="py-8 text-center text-sm text-gray-400">No payments found</p>
            ) : filteredPayments.map((pay, idx) => (
              <div key={String(pay.payment_id ?? idx)} className="flex items-center justify-between py-3 gap-2">
                <div className="flex items-center gap-3">
                  <div className="flex h-8 w-8 items-center justify-center rounded-full bg-green-50">
                    <LuCreditCard className="text-green-500 text-sm" />
                  </div>
                  <div>
                    <p className="text-xs font-medium text-gray-800">{String(pay.transaction_id ?? pay.payment_id ?? "—")}</p>
                    <p className="text-[10px] text-gray-400">{String(pay.completed_at ?? pay.initiated_at ?? pay.created_at ?? "")}</p>
                  </div>
                </div>
                <div className="flex flex-col items-end gap-1">
                  <p className="text-xs font-semibold text-gray-800">
                    {pay.amount ? `BDT ${Number((pay.amount as Record<string,unknown>)?.decimal_amount ?? pay.amount).toLocaleString()}` : "—"}
                  </p>
                  <StatusBadge status={String(pay.status ?? "")} />
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
};

export default BillingInvoices;
