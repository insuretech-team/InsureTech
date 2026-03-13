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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { LuCalendarDays, LuLoader } from "react-icons/lu";
import { useEmployeeForm } from "@/src/hooks/useEmployeeForm";
import { useToast } from "@/src/hooks/useToast";
import { ToastBanner } from "@/components/ui/toast-banner";
import type { EmployeeFormValues } from "@/src/lib/types/employee-form";

// ─── Types ───────────────────────────────────────────────────────────────────

type Department = { id: string; name: string };

type AddEmployeeModalProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  /** Pass to enable edit mode. If undefined → create mode. */
  employeeUuid?: string;
  organisationId?: string;
  initialValues?: Partial<EmployeeFormValues>;
  /** Called after successful create/update so parent can refresh the list */
  onSaved?: () => void;
};

// ─── Styles ──────────────────────────────────────────────────────────────────

const focusPurple =
  "focus-visible:ring-primary focus-visible:border-primary focus-visible:ring-2";

// ─── Section: Personal Info ───────────────────────────────────────────────────

function PersonalInfoSection({
  values,
  errors,
  setField,
}: {
  values: EmployeeFormValues;
  errors: Record<string, string | undefined>;
  setField: <K extends keyof EmployeeFormValues>(k: K, v: EmployeeFormValues[K]) => void;
}) {
  return (
    <div className="space-y-4">
      <p className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Personal Info</p>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Field>
          <Label htmlFor="emp-name" className="sr-only">Name</Label>
          <Input
            id="emp-name"
            placeholder="Full Name*"
            value={values.name}
            onChange={(e) => setField("name", e.target.value)}
            className={`${focusPurple} ${errors.name ? "border-red-500" : ""}`}
            required
          />
          {errors.name && <p className="text-xs text-red-500 mt-1">{errors.name}</p>}
        </Field>

        <Field>
          <Label htmlFor="emp-email" className="sr-only">Email</Label>
          <Input
            id="emp-email"
            type="email"
            placeholder="Work Email"
            value={values.email}
            onChange={(e) => setField("email", e.target.value)}
            className={focusPurple}
          />
        </Field>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Field>
          <Label htmlFor="emp-mobile" className="sr-only">Mobile</Label>
          <Input
            id="emp-mobile"
            placeholder="Mobile Number (+880...)"
            value={values.mobileNumber}
            onChange={(e) => setField("mobileNumber", e.target.value)}
            className={focusPurple}
          />
        </Field>

        <Field>
          <Label htmlFor="emp-gender" className="sr-only">Gender</Label>
          <Select
            value={values.gender || undefined}
            onValueChange={(v) => setField("gender", v as EmployeeFormValues["gender"])}
          >
            <SelectTrigger id="emp-gender" className={`w-full ${focusPurple}`}>
              <SelectValue placeholder="Gender" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="EMPLOYEE_GENDER_MALE">Male</SelectItem>
              <SelectItem value="EMPLOYEE_GENDER_FEMALE">Female</SelectItem>
              <SelectItem value="EMPLOYEE_GENDER_OTHER">Other</SelectItem>
            </SelectContent>
          </Select>
        </Field>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Field>
          <Label htmlFor="emp-dob" className="sr-only">Date of Birth</Label>
          <div className="relative">
            <Input
              id="emp-dob"
              type="date"
              placeholder="Date of Birth"
              value={values.dateOfBirth}
              onChange={(e) => setField("dateOfBirth", e.target.value)}
              className={`${focusPurple} pr-10`}
            />
            <LuCalendarDays className="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-primary" />
          </div>
        </Field>
      </div>
    </div>
  );
}

// ─── Section: Employment Info ─────────────────────────────────────────────────

function EmploymentSection({
  values,
  errors,
  setField,
  departments,
  loadingDepts,
}: {
  values: EmployeeFormValues;
  errors: Record<string, string | undefined>;
  setField: <K extends keyof EmployeeFormValues>(k: K, v: EmployeeFormValues[K]) => void;
  departments: Department[];
  loadingDepts: boolean;
}) {
  return (
    <div className="space-y-4">
      <p className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Employment</p>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Field>
          <Label htmlFor="emp-eid" className="sr-only">Employee ID</Label>
          <Input
            id="emp-eid"
            placeholder="Employee ID*"
            value={values.employeeId}
            onChange={(e) => setField("employeeId", e.target.value)}
            className={`${focusPurple} ${errors.employeeId ? "border-red-500" : ""}`}
            required
          />
          {errors.employeeId && <p className="text-xs text-red-500 mt-1">{errors.employeeId}</p>}
        </Field>

        <Field>
          <Label htmlFor="emp-dept" className="sr-only">Department</Label>
          <Select
            value={values.departmentId || undefined}
            onValueChange={(v) => setField("departmentId", v)}
            disabled={loadingDepts}
          >
            <SelectTrigger
              id="emp-dept"
              className={`w-full ${focusPurple} ${errors.departmentId ? "border-red-500" : ""}`}
            >
              <SelectValue placeholder={loadingDepts ? "Loading departments…" : "Department*"} />
            </SelectTrigger>
            <SelectContent>
              {values.departmentId && !departments.some((d) => d.id === values.departmentId) && (
                <SelectItem value={values.departmentId}>
                  Unassigned / Missing
                </SelectItem>
              )}
              {departments.map((d) => (
                <SelectItem key={d.id} value={d.id}>
                  {d.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {errors.departmentId && <p className="text-xs text-red-500 mt-1">{errors.departmentId}</p>}
        </Field>
      </div>

      <Field>
        <Label htmlFor="emp-doj" className="sr-only">Date of Joining</Label>
        <div className="relative">
          <Input
            id="emp-doj"
            type="date"
            placeholder="Date of Joining*"
            value={values.dateOfJoining}
            onChange={(e) => setField("dateOfJoining", e.target.value)}
            className={`${focusPurple} pr-10 ${errors.dateOfJoining ? "border-red-500" : ""}`}
            required
          />
          <LuCalendarDays className="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-primary" />
        </div>
        {errors.dateOfJoining && <p className="text-xs text-red-500 mt-1">{errors.dateOfJoining}</p>}
      </Field>
    </div>
  );
}

// ─── Section: Insurance Info ──────────────────────────────────────────────────

function InsuranceSection({
  values,
  setField,
}: {
  values: EmployeeFormValues;
  setField: <K extends keyof EmployeeFormValues>(k: K, v: EmployeeFormValues[K]) => void;
}) {
  return (
    <div className="space-y-4">
      <p className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Insurance</p>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Field>
          <Label htmlFor="emp-ins-cat" className="sr-only">Insurance Category</Label>
          <Select
            value={values.insuranceCategory ? String(values.insuranceCategory) : undefined}
            onValueChange={(v) => setField("insuranceCategory", Number(v))}
          >
            <SelectTrigger id="emp-ins-cat" className={`w-full ${focusPurple}`}>
              <SelectValue placeholder="Insurance Category" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="1">Health Insurance</SelectItem>
              <SelectItem value="2">Life Insurance</SelectItem>
              <SelectItem value="3">Motor Insurance</SelectItem>
              <SelectItem value="4">Travel Insurance</SelectItem>
            </SelectContent>
          </Select>
        </Field>

        <Field>
          <Label htmlFor="emp-coverage" className="sr-only">Coverage Amount</Label>
          <Input
            id="emp-coverage"
            type="number"
            placeholder="Coverage Amount (BDT)"
            value={values.coverageAmount}
            onChange={(e) => setField("coverageAmount", e.target.value)}
            className={focusPurple}
            min={0}
          />
        </Field>
      </div>

      <Field>
        <Label htmlFor="emp-dependents" className="sr-only">Number of Dependents</Label>
        <Input
          id="emp-dependents"
          type="number"
          placeholder="Number of Dependents"
          value={values.numberOfDependent === 0 ? "" : String(values.numberOfDependent)}
          onChange={(e) => setField("numberOfDependent", Number(e.target.value))}
          className={focusPurple}
          min={0}
        />
      </Field>
    </div>
  );
}

// ─── Main Modal ───────────────────────────────────────────────────────────────

const AddEmployeeModal = ({
  open,
  onOpenChange,
  employeeUuid,
  organisationId,
  initialValues,
  onSaved,
}: AddEmployeeModalProps) => {
  const isEdit = Boolean(employeeUuid);
  const { toast, showToast } = useToast();

  // Hook must be declared BEFORE the departments effect so `values` is in scope.
  const { values, errors, submitting, loadingRecord, setField, submit } = useEmployeeForm({
    mode: isEdit ? "edit" : "create",
    employeeUuid,
    initialValues: {
      ...initialValues,
      businessId: initialValues?.businessId ?? organisationId ?? "",
    },
    onSuccess: (message) => {
      showToast("success", message);
      setTimeout(() => { onSaved?.(); onOpenChange(false); }, 1200);
    },
    onError: (message) => showToast("error", message),
  });

  // Departments fetched from API.
  // organisationId prop is used for create mode.
  // In edit mode, once the hook fetches the full record, values.businessId is populated —
  // we depend on both so departments reload when the business context is resolved.
  const [departments, setDepartments] = React.useState<Department[]>([]);
  const [loadingDepts, setLoadingDepts] = React.useState(false);

  React.useEffect(() => {
    if (!open) return;
    // Resolve the org id: explicit prop > loaded record's businessId > nothing
    const orgId = organisationId || values.businessId;
    if (!orgId) return;

    let cancelled = false;
    setLoadingDepts(true);

    const params = new URLSearchParams();
    params.set("business_id", orgId);

    fetch(`/api/departments?${params.toString()}`, { method: "GET", cache: "no-store" })
      .then((r) => r.json())
      .then((payload: { ok: boolean; departments?: Array<{ id: string; name: string }> }) => {
        if (cancelled) return;
        setDepartments(Array.isArray(payload.departments) ? payload.departments : []);
      })
      .catch(() => { if (!cancelled) setDepartments([]); })
      .finally(() => { if (!cancelled) setLoadingDepts(false); });

    return () => { cancelled = true; };
    // values.businessId changes after the edit-mode fetch resolves — reload depts then
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, organisationId, values.businessId]);

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-2xl p-0 max-h-[90vh] overflow-y-auto">
        <DialogHeader className="px-6 py-4 border-b sticky top-0 bg-white z-10">
          <DialogTitle className="text-xl font-semibold flex items-center gap-2">
            {isEdit ? "Edit Employee" : "Add Employee"}
            {loadingRecord && <LuLoader className="size-4 animate-spin text-muted-foreground" />}
          </DialogTitle>
        </DialogHeader>

        <ToastBanner toast={toast} />

        <form onSubmit={submit} className="px-6 py-6 space-y-6">
          <FieldGroup className="space-y-6 gap-0">
            <PersonalInfoSection values={values} errors={errors} setField={setField} />
            <EmploymentSection
              values={values}
              errors={errors}
              setField={setField}
              departments={departments}
              loadingDepts={loadingDepts}
            />
            <InsuranceSection values={values} setField={setField} />
          </FieldGroup>

          <DialogFooter className="mt-6">
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
              disabled={submitting}
              className="h-11 px-8 text-white bg-gradient-to-r from-primary to-accent hover:opacity-95"
            >
              {submitting ? (
                <span className="flex items-center gap-2">
                  <LuLoader className="animate-spin" />
                  {isEdit ? "Saving…" : "Adding…"}
                </span>
              ) : isEdit ? (
                "Save Changes"
              ) : (
                "Add Employee"
              )}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
};

export default AddEmployeeModal;
