"use client";

import * as React from "react";
import Image from "next/image";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { XCircle, Sparkles, Clock } from "lucide-react";

type PlanDetailModalProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

const PlanDetailModal = ({ open, onOpenChange }: PlanDetailModalProps) => {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent
        className="
    fixed
    left-1/2 top-1/2
    w-[100vw] sm:w-[95vw]
    max-w-[100vw] sm:max-w-[720px] lg:max-w-[920px]
    h-[90vh] sm:h-auto
    max-h-[90vh]
    -translate-x-1/2 -translate-y-1/2
    overflow-y-auto
    rounded-none sm:rounded-2xl
    p-0
  "
      >
        {/* Header */}
        <DialogHeader className="px-6 pt-6">
          <DialogTitle className="text-lg font-semibold text-gray-900">
            Plan Details
          </DialogTitle>
        </DialogHeader>

        {/* Banner */}
        <div className="px-6 pt-4">
          <div className="relative h-[260px] w-full overflow-hidden rounded-xl">
            <Image
              src="/insurance-plans/plan_detail_banner.png"
              alt="Plan banner"
              fill
              className="object-cover"
              priority
            />
          </div>
        </div>

        {/* Content */}
        <div className="space-y-6 px-6 py-6">
          {/* Plan Name + Price */}
          <div className="flex items-start justify-between">
            <div>
              <h3 className="text-xl font-semibold text-gray-900">Seba</h3>
              <p className="mt-1 text-sm text-gray-500">
                Health coverage up to
                <span className="ml-1 font-medium block text-[var(--brand-jungle)]">
                  ৳ 25,000
                </span>
              </p>
            </div>

            <div className="text-right">
              <div className="flex items-center justify-end gap-1 ">
                <Sparkles className="text-[var(--brand-jungle)]" size={16} />
                <span className="text-sm font-medium">Premium price:</span>
              </div>
              <p className="text-lg font-semibold text-[var(--brand-jungle)]">
                ৳ 800
              </p>

              <div className="mt-1 flex items-center justify-end gap-1 ">
                <Clock className="text-[var(--brand-jungle)]" size={14} />
                <span className="text-sm font-medium">Policy duration:</span>
              </div>
              <p className="text-lg font-semibold text-[var(--brand-jungle)]">
                1 year
              </p>
            </div>
          </div>

          {/* Policy Purpose */}
          <div>
            <h4 className="mb-1 text-sm font-semibold text-gray-900">
              Policy Purpose
            </h4>
            <p className="text-sm leading-relaxed text-gray-600">
              To provide financial protection against medical expenses arising
              from illness, accidents, hospitalization, and emergency medical
              care.
            </p>
          </div>

          {/* Hospitalization Coverages */}
          <div>
            <h4 className="mb-3 text-sm font-semibold text-gray-900">
              Hospitalization Coverages
            </h4>

            <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
              <CoverageItem
                icon="/insurance-plans/fi_11569746.svg"
                label="In-patient treatment expenses"
              />
              <CoverageItem
                icon="/insurance-plans/fi_3466901.svg"
                label="Cabin room rent"
              />
              <CoverageItem
                icon="/insurance-plans/fi_13214156.svg"
                label="ICU / CCU charges"
              />
              <CoverageItem
                icon="/insurance-plans/fi_7381103.svg"
                label="Doctor & specialist consultation fees"
              />
            </div>
          </div>

          {/* Excluded Coverages */}
          <div>
            <h4 className="mb-3 text-sm font-semibold text-gray-900">
              Excluded Coverages
            </h4>

            <ul className="space-y-2">
              <ExcludeItem label="Cosmetic / aesthetic treatment" />
              <ExcludeItem label="Dental & vision (unless due to accident)" />
              <ExcludeItem label="Infertility treatments" />
              <ExcludeItem label="Self-inflicted injuries" />
              <ExcludeItem label="Non-prescribed medicines" />
              <ExcludeItem label="Treatment outside network hospitals (unless emergency)" />
            </ul>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
};

export default PlanDetailModal;

/* ---------- Small Components ---------- */

const CoverageItem = ({ icon, label }: { icon: string; label: string }) => (
  <div className="flex items-center gap-3 rounded-md bg-[var(--brand-surface-5)] px-3 py-3">
    <Image
      src={icon}
      width={24}
      height={24}
      alt={label}
      className="text-[var(--brand-jungle)]"
    />
    <span className="text-sm text-[var(--brand-jungle)]">{label}</span>
  </div>
);

const ExcludeItem = ({ label }: { label: string }) => (
  <li className="flex items-start gap-2 text-sm text-gray-600">
    <XCircle className="mt-0.5 text-red-500" size={16} />
    <span>{label}</span>
  </li>
);
