"use client";

import * as React from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Field, FieldGroup } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { LuLoader } from "react-icons/lu";
import { useToast } from "@/src/hooks/useToast";
import { ToastBanner } from "@/components/ui/toast-banner";

type AddDepartmentModalProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  /** Pass to enable edit mode */
  departmentId?: string;
  initialName?: string;
  /** Organisation to create the department under (required for create mode) */
  organisationId?: string;
  /** Called after successful save so parent can refresh */
  onSaved?: () => void;
};

const focusPurple =
  "focus-visible:ring-primary focus-visible:border-primary focus-visible:ring-2";

const AddDepartmentModal = ({
  open,
  onOpenChange,
  departmentId,
  initialName = "",
  organisationId = "",
  onSaved,
}: AddDepartmentModalProps) => {
  const isEdit = Boolean(departmentId);
  const [name, setName] = React.useState(initialName);
  const [error, setError] = React.useState("");
  const [submitting, setSubmitting] = React.useState(false);
  const { toast, showToast } = useToast();

  // Reset state when modal opens
  React.useEffect(() => {
    if (open) {
      setName(initialName);
      setError("");
    }
  }, [open, initialName]);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!name.trim()) {
      setError("Department name is required");
      return;
    }

    setSubmitting(true);
    setError("");

    try {
      let url = "/api/departments";
      let method = "POST";

      if (isEdit && departmentId) {
        url = `/api/departments/${departmentId}`;
        method = "PATCH";
      }

      const response = await fetch(url, {
        method,
        headers: { "Content-Type": "application/json" },
        // businessId required for create; ignored by PATCH (backend uses dept_id for scoping)
        body: JSON.stringify({ name: name.trim(), businessId: organisationId || undefined }),
      });

      const payload = (await response.json()) as { ok: boolean; message?: string };

      if (!response.ok || !payload.ok) {
        showToast("error", payload.message ?? "Operation failed");
        return;
      }

      showToast("success", payload.message ?? (isEdit ? "Department updated" : "Department created"));
      setTimeout(() => { onSaved?.(); onOpenChange(false); }, 1200);
    } catch (err) {
      showToast("error", err instanceof Error ? err.message : "Unexpected error");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-lg p-0">
        <DialogHeader className="px-6 py-4 border-b">
          <DialogTitle className="text-xl font-semibold">
            {isEdit ? "Edit Department" : "Add Department"}
          </DialogTitle>
        </DialogHeader>

        <ToastBanner toast={toast} />

        <form onSubmit={onSubmit} className="px-6 py-6">
          <FieldGroup className="space-y-4 gap-0">
            <Field>
              <Label htmlFor="dept-name" className="sr-only">Department Name</Label>
              <Input
                id="dept-name"
                name="name"
                placeholder="Department Name*"
                value={name}
                onChange={(e) => {
                  setName(e.target.value);
                  if (error) setError("");
                }}
                className={`${focusPurple} ${error ? "border-red-500" : ""}`}
                required
              />
              {error && <p className="text-xs text-red-500 mt-1">{error}</p>}
            </Field>
          </FieldGroup>

          <DialogFooter className="mt-8 flex gap-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={submitting}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={submitting}
              className="h-11 px-8 text-white bg-gradient-to-r from-primary to-accent hover:opacity-95"
            >
              {submitting ? (
                <span className="flex items-center gap-2">
                  <LuLoader className="animate-spin" />
                  {isEdit ? "Saving…" : "Creating…"}
                </span>
              ) : isEdit ? (
                "Save Changes"
              ) : (
                "Add Department"
              )}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
};

export default AddDepartmentModal;
