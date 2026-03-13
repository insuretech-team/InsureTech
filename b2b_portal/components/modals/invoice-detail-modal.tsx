"use client";

import * as React from "react";
import Image from "next/image";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Separator } from "@/components/ui/separator";

type InvoiceDetail = {
  invoiceNo: string;
  date: string;
  due: string;
  insuranceCategory: string;
  plan: string;
  employeesCovered: number | string;
  billingPeriod: string;
  grossPremium: string;
  vat: string;
  totalPayable: string;
  paymentMethod: string;
};

type InvoiceDetailModalProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  invoice?: Partial<InvoiceDetail>;
};

const defaultInvoice: InvoiceDetail = {
  invoiceNo: "INV-2026-001",
  date: "18 Feb 2026",
  due: "10 Mar 2026",
  insuranceCategory: "Health",
  plan: "Seba",
  employeesCovered: 425,
  billingPeriod: "Jan-Dec 2026",
  grossPremium: "৳14,20,000",
  vat: "৳2,13,000",
  totalPayable: "৳16,33,000",
  paymentMethod: "Bank Transfer",
};

const Row = ({ label, value }: { label: string; value: React.ReactNode }) => (
  <div className="flex items-center justify-between gap-4 px-6 py-3">
    <p className="text-[13px] text-gray-500">{label}</p>
    <p className="text-[13px] font-semibold text-gray-900">{value}</p>
  </div>
);

const InvoiceDetailModal = ({
  open,
  onOpenChange,
  invoice,
}: InvoiceDetailModalProps) => {
  const data = { ...defaultInvoice, ...(invoice || {}) };

  const rows: Array<{ label: string; value: React.ReactNode }> = [
    { label: "Invoice No", value: data.invoiceNo },
    { label: "Date", value: data.date },
    { label: "Due", value: data.due },
    { label: "Insurance Category", value: data.insuranceCategory },
    { label: "Plan", value: data.plan },
    { label: "Employees Covered", value: data.employeesCovered },
    { label: "Billing Period", value: data.billingPeriod },
    { label: "Gross Premium", value: data.grossPremium },
    { label: "VAT", value: data.vat },
    { label: "Total Payable", value: data.totalPayable },
    { label: "Payment Method", value: data.paymentMethod },
  ];

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="w-[92vw] max-w-[520px] overflow-hidden rounded-2xl p-0">
        {/* Header */}
        <DialogHeader className="px-6 pt-8">
          <DialogTitle className="flex items-center justify-center">
            <span className="flex h-24 w-24 items-center justify-center rounded-full bg-[var(--brand-surface-1)]">
              <Image
                src="/fi_9485672.svg"
                width={54}
                height={54}
                alt="Invoice"
                priority
              />
            </span>
          </DialogTitle>

          <h3 className="mt-5 text-left text-xl font-semibold text-gray-900">
            Invoice Details
          </h3>
        </DialogHeader>

        {/* Body */}
        <div className="pb-2">
          <Separator />
          {rows.map((r, i) => (
            <React.Fragment key={`${r.label}-${i}`}>
              <Row label={r.label} value={r.value} />
              {i !== rows.length - 1 ? <Separator /> : null}
            </React.Fragment>
          ))}
        </div>
      </DialogContent>
    </Dialog>
  );
};

export default InvoiceDetailModal;
