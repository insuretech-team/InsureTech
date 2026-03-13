"use client";

import * as React from "react";
import { LuLoader } from "react-icons/lu";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

// ─── Types ────────────────────────────────────────────────────────────────────

type DepartmentOption = { id: string; name: string };

type CatalogOption = {
  planId: string;
  productId: string;
  productName: string;
  planName: string;
  insuranceCategory: string;
  premiumAmount: string;
};

export type PurchaseOrderFormInput = {
  departmentId:       string;
  planId:             string;
  insuranceCategory:  string;  // auto-derived from selected plan
  employeeCount:      number;
  numberOfDependents: number;
  coverageAmount:     number;
  notes:              string;
};

type AddPurchaseOrderModalProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  departments: DepartmentOption[];
  catalog: CatalogOption[];
  submitting?: boolean;
  /** Return true on success (modal closes), false on error (modal stays open) */
  onSubmit: (payload: PurchaseOrderFormInput) => Promise<boolean>;
};

// ─── Helpers ──────────────────────────────────────────────────────────────────

const focusRing = "focus-visible:ring-primary focus-visible:border-primary focus-visible:ring-2";

const initialState: PurchaseOrderFormInput = {
  departmentId:       "",
  planId:             "",
  insuranceCategory:  "",
  employeeCount:      1,
  numberOfDependents: 0,
  coverageAmount:     0,
  notes:              "",
};

/** Visible field label above the control */
function FieldLabel({ htmlFor, children, required }: { htmlFor: string; children: React.ReactNode; required?: boolean }) {
  return (
    <label htmlFor={htmlFor} className="block text-sm font-medium text-foreground mb-1.5">
      {children}{required && <span className="text-destructive ml-0.5">*</span>}
    </label>
  );
}

// ─── Modal ────────────────────────────────────────────────────────────────────

const AddPurchaseOrderModal = ({
  open,
  onOpenChange,
  departments,
  catalog,
  submitting = false,
  onSubmit,
}: AddPurchaseOrderModalProps) => {
  const [form, setForm] = React.useState<PurchaseOrderFormInput>(initialState);
  const [errors, setErrors] = React.useState<Partial<Record<keyof PurchaseOrderFormInput, string>>>({});

  // Reset form when modal closes
  React.useEffect(() => {
    if (!open) { setForm(initialState); setErrors({}); }
  }, [open]);

  // Auto-derive insuranceCategory from the selected plan
  const selectedPlan = catalog.find((item) => item.planId === form.planId);
  React.useEffect(() => {
    if (selectedPlan) {
      setForm((cur) => ({ ...cur, insuranceCategory: selectedPlan.insuranceCategory }));
    }
  }, [selectedPlan]);

  function setField<K extends keyof PurchaseOrderFormInput>(key: K, value: PurchaseOrderFormInput[K]) {
    setForm((cur) => ({ ...cur, [key]: value }));
    setErrors((cur) => { const n = { ...cur }; delete n[key]; return n; });
  }

  function validate(): boolean {
    const e: typeof errors = {};
    if (!form.departmentId)     e.departmentId     = "Department is required";
    if (!form.planId)           e.planId           = "Product plan is required";
    if (form.employeeCount < 1) e.employeeCount    = "At least 1 employee is required";
    if (form.coverageAmount <= 0) e.coverageAmount = "Coverage amount must be greater than 0";
    setErrors(e);
    return Object.keys(e).length === 0;
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!validate()) return;
    const success = await onSubmit(form);
    if (success) onOpenChange(false);
  }

  const isLoading = departments.length === 0 || catalog.length === 0;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-2xl p-0 max-h-[90vh] overflow-y-auto">
        <DialogHeader className="px-6 py-4 border-b sticky top-0 bg-white z-10">
          <DialogTitle className="text-xl font-semibold">Create Purchase Order</DialogTitle>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="px-6 py-6 space-y-5">

          {/* Department */}
          <div>
            <FieldLabel htmlFor="po-dept" required>Department</FieldLabel>
            <Select
              value={form.departmentId}
              onValueChange={(v) => setField("departmentId", v)}
            >
              <SelectTrigger id="po-dept" className={`w-full h-11 ${focusRing} ${errors.departmentId ? "border-destructive" : ""}`}>
                <SelectValue placeholder={departments.length === 0 ? "Loading departments…" : "Select department"} />
              </SelectTrigger>
              <SelectContent>
                {departments.map((d) => (
                  <SelectItem key={d.id} value={d.id}>{d.name}</SelectItem>
                ))}
              </SelectContent>
            </Select>
            {errors.departmentId && <p className="mt-1 text-xs text-destructive">{errors.departmentId}</p>}
          </div>

          {/* Product / Plan */}
          <div>
            <FieldLabel htmlFor="po-plan" required>Product Plan</FieldLabel>
            <Select
              value={form.planId}
              onValueChange={(v) => setField("planId", v)}
            >
              <SelectTrigger id="po-plan" className={`w-full h-11 ${focusRing} ${errors.planId ? "border-destructive" : ""}`}>
                <SelectValue placeholder={catalog.length === 0 ? "Loading plans…" : "Select product plan"} />
              </SelectTrigger>
              <SelectContent>
                {catalog.map((item) => (
                  <SelectItem key={item.planId} value={item.planId}>
                    {item.productName} — {item.planName}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            {errors.planId && <p className="mt-1 text-xs text-destructive">{errors.planId}</p>}
          </div>

          {/* Plan summary card — shown when plan is selected */}
          {selectedPlan && (
            <div className="rounded-xl border border-border/70 bg-muted/20 px-4 py-3 space-y-1">
              <div className="text-sm font-semibold text-foreground">
                {selectedPlan.productName} / {selectedPlan.planName}
              </div>
              <div className="flex flex-wrap gap-x-4 gap-y-0.5 text-xs text-muted-foreground">
                <span>Category: <span className="font-medium text-foreground">{selectedPlan.insuranceCategory}</span></span>
                <span>Base premium: <span className="font-medium text-foreground">{selectedPlan.premiumAmount}</span></span>
              </div>
            </div>
          )}

          {/* Insurance Category — auto-filled from plan, shown as readonly info */}
          {form.insuranceCategory && (
            <div>
              <FieldLabel htmlFor="po-category">Insurance Category</FieldLabel>
              <Input
                id="po-category"
                value={form.insuranceCategory}
                readOnly
                className={`h-11 bg-muted/30 text-muted-foreground cursor-not-allowed ${focusRing}`}
              />
              <p className="mt-1 text-xs text-muted-foreground">Auto-assigned from selected plan.</p>
            </div>
          )}

          {/* Employee & Dependent counts */}
          <div className="grid gap-4 sm:grid-cols-2">
            <div>
              <FieldLabel htmlFor="po-emp-count" required>Number of Employees</FieldLabel>
              <Input
                id="po-emp-count"
                type="number"
                min={1}
                className={`h-11 ${focusRing} ${errors.employeeCount ? "border-destructive" : ""}`}
                value={String(form.employeeCount)}
                onChange={(e) => setField("employeeCount", Number(e.target.value))}
              />
              {errors.employeeCount && <p className="mt-1 text-xs text-destructive">{errors.employeeCount}</p>}
            </div>
            <div>
              <FieldLabel htmlFor="po-dep-count">Number of Dependents</FieldLabel>
              <Input
                id="po-dep-count"
                type="number"
                min={0}
                className={`h-11 ${focusRing}`}
                value={String(form.numberOfDependents)}
                onChange={(e) => setField("numberOfDependents", Number(e.target.value))}
              />
            </div>
          </div>

          {/* Coverage Amount */}
          <div>
            <FieldLabel htmlFor="po-coverage" required>Coverage Amount (BDT)</FieldLabel>
            <Input
              id="po-coverage"
              type="number"
              min={1}
              className={`h-11 ${focusRing} ${errors.coverageAmount ? "border-destructive" : ""}`}
              placeholder="e.g. 500000"
              value={form.coverageAmount === 0 ? "" : String(form.coverageAmount)}
              onChange={(e) => setField("coverageAmount", Number(e.target.value))}
            />
            {errors.coverageAmount && <p className="mt-1 text-xs text-destructive">{errors.coverageAmount}</p>}
          </div>

          {/* Notes */}
          <div>
            <FieldLabel htmlFor="po-notes">Notes / Justification</FieldLabel>
            <textarea
              id="po-notes"
              rows={3}
              className={`w-full rounded-md border px-3 py-2 text-sm resize-none ${focusRing}`}
              placeholder="Add any notes or justification for this purchase order…"
              value={form.notes}
              onChange={(e) => setField("notes", e.target.value)}
            />
          </div>

          <DialogFooter className="pt-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={submitting}
              className="mr-2"
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={submitting || isLoading}
              className="h-11 px-8 text-white bg-gradient-to-r from-primary to-accent hover:opacity-95"
            >
              {submitting ? (
                <span className="flex items-center gap-2">
                  <LuLoader className="animate-spin" /> Submitting…
                </span>
              ) : "Submit Purchase Order"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
};

export default AddPurchaseOrderModal;
