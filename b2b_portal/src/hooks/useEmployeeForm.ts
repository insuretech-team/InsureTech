/**
 * useEmployeeForm.ts
 * ──────────────────
 * Form state hook for create / edit employee.
 *
 * In EDIT mode, when `employeeUuid` is provided, the hook automatically
 * fetches the full employee record from GET /api/employees/{id} and
 * populates ALL form fields — including email, mobileNumber, gender,
 * dateOfBirth, dateOfJoining, departmentId, insuranceCategory, etc.
 * This ensures the edit form is never partially blank.
 */
"use client";

import { useState, useCallback, useEffect } from "react";
import type {
  EmployeeFormValues,
  EmployeeFormErrors,
  EmployeeFormMode,
} from "@lib/types/employee-form";
import { EMPTY_EMPLOYEE_FORM } from "@lib/types/employee-form";
import { employeeClient, type EmployeeFullRecord } from "@lib/sdk/employee-client";

interface UseEmployeeFormOptions {
  mode: EmployeeFormMode;
  employeeUuid?: string;
  initialValues?: Partial<EmployeeFormValues>;
  onSuccess?: (message: string) => void;
  onError?: (message: string) => void;
}

function validate(values: EmployeeFormValues, mode: EmployeeFormMode): EmployeeFormErrors {
  const errors: EmployeeFormErrors = {};
  if (!values.name.trim()) errors.name = "Name is required";
  if (!values.employeeId.trim()) errors.employeeId = "Employee ID is required";
  if (mode === "create" && !values.businessId.trim()) errors.businessId = "Organisation is required";
  if (!values.departmentId.trim()) errors.departmentId = "Department is required";
  if (!values.dateOfJoining.trim()) errors.dateOfJoining = "Date of joining is required";
  return errors;
}

/** Parse a Money-like value from the API into a plain decimal string */
function parseCoverage(raw: unknown): string {
  if (!raw) return "";
  if (typeof raw === "number") return String(raw);
  if (typeof raw === "object" && raw !== null) {
    const m = raw as Record<string, unknown>;
    if (typeof m.decimal_amount === "number") return String(m.decimal_amount);
    if (typeof m.amount === "number") return String(m.amount / 100);
  }
  return "";
}

/** Map a full employee record (from GET /api/employees/{id}) to form values */
function mapApiToForm(raw: EmployeeFullRecord, initialValues?: Partial<EmployeeFormValues>): EmployeeFormValues {
  return {
    ...EMPTY_EMPLOYEE_FORM,
    ...initialValues,
    name:              raw.name              || initialValues?.name        || "",
    employeeId:        raw.employeeID        || initialValues?.employeeId  || "",
    email:             raw.email             || initialValues?.email       || "",
    mobileNumber:      raw.mobileNumber      || initialValues?.mobileNumber|| "",
    gender:            (raw.gender as EmployeeFormValues["gender"]) || initialValues?.gender || "",
    dateOfBirth:       raw.dateOfBirth       || initialValues?.dateOfBirth || "",
    dateOfJoining:     raw.dateOfJoining     || initialValues?.dateOfJoining || "",
    departmentId:      raw.departmentId      || initialValues?.departmentId || "",
    businessId:        raw.businessId        || initialValues?.businessId  || "",
    insuranceCategory: raw.insuranceCategory ?? initialValues?.insuranceCategory ?? 0,
    assignedPlanId:    raw.assignedPlanId    || initialValues?.assignedPlanId || "",
    coverageAmount:    raw.coverageAmount    || initialValues?.coverageAmount  || "",
    numberOfDependent: raw.numberOfDependent ?? initialValues?.numberOfDependent ?? 0,
  };
}

export function useEmployeeForm({
  mode,
  employeeUuid,
  initialValues,
  onSuccess,
  onError,
}: UseEmployeeFormOptions) {
  const [values, setValues] = useState<EmployeeFormValues>({
    ...EMPTY_EMPLOYEE_FORM,
    ...initialValues,
  });
  const [errors, setErrors] = useState<EmployeeFormErrors>({});
  const [submitting, setSubmitting] = useState(false);
  const [loadingRecord, setLoadingRecord] = useState(false);

  // In edit mode, fetch the full employee record so ALL fields are populated.
  // This replaces the sparse `initialValues` (which only has name + employeeId
  // from the table row) with the complete record from the backend.
  useEffect(() => {
    if (mode !== "edit" || !employeeUuid) return;
    let cancelled = false;
    setLoadingRecord(true);

    employeeClient.get(employeeUuid)
      .then((result) => {
        if (cancelled) return;
        if (result.ok && result.employee) {
          setValues(mapApiToForm(result.employee, initialValues));
        } else {
          // Fallback to whatever initialValues we have
          setValues({ ...EMPTY_EMPLOYEE_FORM, ...initialValues });
        }
      })
      .catch(() => {
        if (!cancelled) setValues({ ...EMPTY_EMPLOYEE_FORM, ...initialValues });
      })
      .finally(() => { if (!cancelled) setLoadingRecord(false); });

    return () => { cancelled = true; };
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [mode, employeeUuid]);

  const setField = useCallback(
    <K extends keyof EmployeeFormValues>(field: K, value: EmployeeFormValues[K]) => {
      setValues((prev) => ({ ...prev, [field]: value }));
      setErrors((prev) => {
        if (!prev[field]) return prev;
        const next = { ...prev };
        delete next[field];
        return next;
      });
    },
    []
  );

  const reset = useCallback(() => {
    setValues({ ...EMPTY_EMPLOYEE_FORM, ...initialValues });
    setErrors({});
  }, [initialValues]);

  const submit = useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();
      const validationErrors = validate(values, mode);
      if (Object.keys(validationErrors).length > 0) {
        setErrors(validationErrors);
        return;
      }
      setSubmitting(true);
      try {
        const payload = {
          name:              values.name.trim(),
          employeeId:        values.employeeId.trim(),
          businessId:        values.businessId,
          departmentId:      values.departmentId,
          email:             values.email.trim() || undefined,
          mobileNumber:      values.mobileNumber.trim() || undefined,
          dateOfBirth:       values.dateOfBirth || undefined,
          dateOfJoining:     values.dateOfJoining,
          gender:            values.gender || undefined,
          insuranceCategory: values.insuranceCategory || undefined,
          assignedPlanId:    values.assignedPlanId || undefined,
          coverageAmount:    values.coverageAmount ? Number.parseFloat(values.coverageAmount) : undefined,
          numberOfDependent: values.numberOfDependent || 0,
        };

        const result =
          mode === "edit" && employeeUuid
            ? await employeeClient.update(employeeUuid, payload)
            : await employeeClient.create(payload);

        if (!result.ok) {
          onError?.(result.message ?? "Operation failed");
          return;
        }
        onSuccess?.(result.message ?? (mode === "create" ? "Employee created" : "Employee updated"));
        if (mode === "create") reset();
      } catch (err) {
        onError?.(err instanceof Error ? err.message : "Unexpected error");
      } finally {
        setSubmitting(false);
      }
    },
    [values, mode, employeeUuid, onSuccess, onError, reset]
  );

  return { values, errors, submitting, loadingRecord, setField, reset, submit };
}
