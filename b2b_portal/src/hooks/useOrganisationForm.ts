/**
 * useOrganisationForm.ts
 * ───────────────────────
 * Form state hook for create / edit organisation.
 */
"use client";

import { useState, useCallback } from "react";
import { organisationClient, type OrgCreatePayload } from "@lib/sdk/organisation-client";

export interface OrgFormValues {
  name: string;
  code: string;
  industry: string;
  contactEmail: string;
  contactPhone: string;
  address: string;
  adminFullName: string;
  adminEmail: string;
  adminPassword: string;
  adminMobileNumber: string;
}

export interface OrgFormErrors {
  name?: string;
  contactEmail?: string;
  adminEmail?: string;
  adminPassword?: string;
  adminMobileNumber?: string;
  [key: string]: string | undefined;
}

export const EMPTY_ORG_FORM: OrgFormValues = {
  name: "",
  code: "",
  industry: "",
  contactEmail: "",
  contactPhone: "",
  address: "",
  adminFullName: "",
  adminEmail: "",
  adminPassword: "",
  adminMobileNumber: "",
};

function looksLikeBangladeshMobile(raw: string): boolean {
  const digits = raw.replace(/\D/g, "");
  return /^(?:880|0)?1[3-9]\d{8}$/.test(digits) || /^1[3-9]\d{8}$/.test(digits);
}

function validate(v: OrgFormValues, mode: "create" | "edit"): OrgFormErrors {
  const e: OrgFormErrors = {};
  if (!v.name.trim()) e.name = "Organisation name is required";
  if (mode === "create") {
    if (!v.adminEmail.trim()) e.adminEmail = "Admin email is required";
    if (!v.adminPassword.trim()) {
      e.adminPassword = "Admin password is required";
    } else {
      const hasUpper = /[A-Z]/.test(v.adminPassword);
      const hasLower = /[a-z]/.test(v.adminPassword);
      const hasDigit = /\d/.test(v.adminPassword);
      const hasSymbol = /[^A-Za-z0-9]/.test(v.adminPassword);
      if (v.adminPassword.length < 8 || !hasUpper || !hasLower || !hasDigit || !hasSymbol) {
        e.adminPassword = "Use 8+ chars with upper, lower, number, and symbol";
      }
    }
    if (!v.adminMobileNumber.trim()) {
      e.adminMobileNumber = "Admin mobile number is required";
    } else if (!looksLikeBangladeshMobile(v.adminMobileNumber.trim())) {
      e.adminMobileNumber = "Use a valid BD mobile number";
    }
  }
  return e;
}

interface UseOrgFormOptions {
  mode: "create" | "edit";
  orgId?: string;
  initialValues?: Partial<OrgFormValues>;
  onSuccess?: (message: string) => void;
  onError?: (message: string) => void;
}

export function useOrganisationForm({ mode, orgId, initialValues, onSuccess, onError }: UseOrgFormOptions) {
  const [values, setValues] = useState<OrgFormValues>({ ...EMPTY_ORG_FORM, ...initialValues });
  const [errors, setErrors] = useState<OrgFormErrors>({});
  const [submitting, setSubmitting] = useState(false);

  const setField = useCallback(<K extends keyof OrgFormValues>(field: K, value: OrgFormValues[K]) => {
    setValues((prev) => ({ ...prev, [field]: value }));
    setErrors((prev) => { const n = { ...prev }; delete n[field]; return n; });
  }, []);

  const reset = useCallback(() => {
    setValues({ ...EMPTY_ORG_FORM, ...initialValues });
    setErrors({});
  }, [initialValues]);

  const submit = useCallback(async (e: React.FormEvent) => {
    e.preventDefault();
    const ve = validate(values, mode);
    if (Object.keys(ve).length > 0) { setErrors(ve); return; }
    setSubmitting(true);
    try {
      const payload: OrgCreatePayload = {
        name: values.name.trim(),
        code: values.code.trim() || undefined,
        industry: values.industry.trim() || undefined,
        contactEmail: values.contactEmail.trim() || undefined,
        contactPhone: values.contactPhone.trim() || undefined,
        address: values.address.trim() || undefined,
        admin: mode === "create" ? {
          fullName: values.adminFullName.trim() || undefined,
          email: values.adminEmail.trim(),
          password: values.adminPassword,
          mobileNumber: values.adminMobileNumber.trim(),
        } : undefined,
      };
      const result = mode === "edit" && orgId
        ? await organisationClient.update(orgId, payload)
        : await organisationClient.create(payload);
      if (!result.ok) { onError?.(result.message ?? "Operation failed"); return; }
      onSuccess?.(result.message ?? (mode === "create" ? "Organisation created" : "Organisation updated"));
      if (mode === "create") reset();
    } catch (err) {
      onError?.(err instanceof Error ? err.message : "Unexpected error");
    } finally {
      setSubmitting(false);
    }
  }, [values, mode, orgId, onSuccess, onError, reset]);

  return { values, errors, submitting, setField, reset, submit };
}
