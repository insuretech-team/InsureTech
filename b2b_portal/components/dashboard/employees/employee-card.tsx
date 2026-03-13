"use client";

import * as React from "react";
import { LuPen, LuTrash2, LuX, LuLoader } from "react-icons/lu";
import type { Employee } from "@lib/types/b2b";
import { StatusBadge } from "@/components/ui/status-badge";
import { Button } from "@/components/ui/button";

type EmployeeCardProps = {
  employee: Employee | null;
  open: boolean;
  onClose: () => void;
  onEdit: () => void;
  onDelete: () => void;
  deleting: boolean;
};

export function EmployeeCard({
  employee,
  open,
  onClose,
  onEdit,
  onDelete,
  deleting,
}: EmployeeCardProps) {
  if (!open || !employee) return null;

  return (
    <>
      {/* Backdrop overlay */}
      <div
        className="fixed inset-0 bg-black/50 z-40"
        onClick={onClose}
        aria-hidden="true"
      />

      {/* Slide-over panel - right side */}
      <div className="fixed right-0 top-0 bottom-0 w-full max-w-md bg-white shadow-lg z-50 overflow-y-auto animate-in slide-in-from-right">
        {/* Header with close button */}
        <div className="sticky top-0 bg-white border-b px-6 py-4 flex items-center justify-between">
          <h2 className="text-xl font-semibold text-foreground">
            Employee Details
          </h2>
          <button
            onClick={onClose}
            className="rounded-md p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground transition-colors"
            aria-label="Close panel"
          >
            <LuX className="size-5" />
          </button>
        </div>

        {/* Content */}
        <div className="p-6 space-y-6">
          {/* Personal Info Section */}
          <div>
            <h3 className="text-sm font-semibold uppercase tracking-wider text-muted-foreground mb-4">
              Personal Information
            </h3>
            <div className="space-y-3">
              <div>
                <p className="text-xs text-muted-foreground mb-1">Name</p>
                <p className="text-sm font-medium text-foreground">
                  {employee.name}
                </p>
              </div>
              <div>
                <p className="text-xs text-muted-foreground mb-1">
                  Employee ID
                </p>
                <p className="text-sm font-medium text-foreground">
                  {employee.employeeID}
                </p>
              </div>
            </div>
          </div>

          {/* Employment Info Section */}
          <div>
            <h3 className="text-sm font-semibold uppercase tracking-wider text-muted-foreground mb-4">
              Employment Information
            </h3>
            <div className="space-y-3">
              <div>
                <p className="text-xs text-muted-foreground mb-1">
                  Department
                </p>
                <p className="text-sm font-medium text-foreground">
                  {employee.department}
                </p>
              </div>
            </div>
          </div>

          {/* Insurance Info Section */}
          <div>
            <h3 className="text-sm font-semibold uppercase tracking-wider text-muted-foreground mb-4">
              Insurance Details
            </h3>
            <div className="space-y-3">
              {employee.insuranceCategory && (
                <div>
                  <p className="text-xs text-muted-foreground mb-1">
                    Insurance Category
                  </p>
                  <p className="text-sm font-medium text-foreground">
                    {employee.insuranceCategory}
                  </p>
                </div>
              )}
              {employee.assignedPlan && (
                <div>
                  <p className="text-xs text-muted-foreground mb-1">
                    Assigned Plan
                  </p>
                  <p className="text-sm font-medium text-foreground">
                    {employee.assignedPlan}
                  </p>
                </div>
              )}
              <div>
                <p className="text-xs text-muted-foreground mb-1">Coverage</p>
                <p className="text-sm font-medium text-foreground">
                  {employee.coverage}
                </p>
              </div>
              <div>
                <p className="text-xs text-muted-foreground mb-1">
                  Premium Amount
                </p>
                <p className="text-sm font-medium text-foreground">
                  {employee.premiumAmount}
                </p>
              </div>
              <div>
                <p className="text-xs text-muted-foreground mb-1">
                  Number of Dependents
                </p>
                <p className="text-sm font-medium text-foreground">
                  {employee.numberOfDependent}
                </p>
              </div>
            </div>
          </div>

          {/* Status Section */}
          <div>
            <h3 className="text-sm font-semibold uppercase tracking-wider text-muted-foreground mb-4">
              Status
            </h3>
            <StatusBadge status={employee.status} />
          </div>
        </div>

        {/* Footer with action buttons */}
        <div className="sticky bottom-0 bg-white border-t px-6 py-4 flex items-center gap-3">
          <Button
            variant="outline"
            onClick={onEdit}
            className="flex-1 flex items-center justify-center gap-2"
            disabled={deleting}
          >
            <LuPen className="size-4" />
            <span>Edit</span>
          </Button>
          <Button
            variant="destructive"
            onClick={onDelete}
            className="flex-1 flex items-center justify-center gap-2"
            disabled={deleting}
          >
            {deleting ? (
              <>
                <LuLoader className="size-4 animate-spin" />
                <span>Deleting…</span>
              </>
            ) : (
              <>
                <LuTrash2 className="size-4" />
                <span>Delete</span>
              </>
            )}
          </Button>
        </div>
      </div>
    </>
  );
}
