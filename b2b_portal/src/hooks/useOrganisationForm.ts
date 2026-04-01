/**
 * useOrganisationForm.ts
 * ───────────────────────
 * Form state hook for create / edit organisation.
 */
"use client";

import { useState, useCallback } from "react";
import { bffClient, type OrgCreatePayload } from "@lib/sdk/b2b-sdk-client";
import {
  getBangladeshMobileValidationMessage,
  normalizeBangladeshMobile,
} from "@/src/lib/utils/bd-mobile";
import { getPasswordValidationMessage } from "@/src/lib/utils/password";

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
  contactPhone?: string;
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

function validateWithOptions(
  v: OrgFormValues,
  mode: "create" | "edit",
  requireDefaultAdmin: boolean,
): OrgFormErrors {
  const e: OrgFormErrors = {};
  if (!v.name.trim()) e.name = "Organisation name is required";
  if (v.contactPhone.trim() && !normalizeBangladeshMobile(v.contactPhone.trim())) {
    e.contactPhone = getBangladeshMobileValidationMessage("Contact phone");
  }
  if (mode === "create" && requireDefaultAdmin) {
    if (!v.adminEmail.trim()) e.adminEmail = "Admin email is required";
    e.adminPassword = getPasswordValidationMessage(v.adminPassword, "Admin password") ?? undefined;
    if (!v.adminMobileNumber.trim()) {
      e.adminMobileNumber = "Admin mobile number is required";
    } else if (!normalizeBangladeshMobile(v.adminMobileNumber.trim())) {
      e.adminMobileNumber = getBangladeshMobileValidationMessage("Admin mobile number");
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

interface SubmitOverrides {
  requireDefaultAdmin?: boolean;
  createPayload?: Partial<OrgCreatePayload>;
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

  const submit = useCallback(async (e: React.FormEvent, overrides?: SubmitOverrides) => {
    e.preventDefault();
    const requireDefaultAdmin = overrides?.requireDefaultAdmin ?? (mode === "create");
    const ve = validateWithOptions(values, mode, requireDefaultAdmin);
    if (Object.keys(ve).length > 0) { setErrors(ve); return; }
    setSubmitting(true);
    try {
      const normalizedContactPhone = values.contactPhone.trim()
        ? normalizeBangladeshMobile(values.contactPhone.trim())
        : null;
      const basePayload: OrgCreatePayload = {
        name: values.name.trim(),
        code: values.code.trim() || undefined,
        industry: values.industry.trim() || undefined,
        contactEmail: values.contactEmail.trim() || undefined,
        contactPhone: normalizedContactPhone ?? undefined,
        address: values.address.trim() || undefined,
      };
      const normalizedAdminMobile =
        mode === "create" && requireDefaultAdmin
          ? normalizeBangladeshMobile(values.adminMobileNumber.trim())
          : null;
      if (mode === "create" && requireDefaultAdmin) {
        basePayload.admin = {
          fullName: values.adminFullName.trim() || undefined,
          email: values.adminEmail.trim(),
          password: values.adminPassword,
          mobileNumber: normalizedAdminMobile ?? values.adminMobileNumber.trim(),
        };
      }
      const payload: OrgCreatePayload = {
        ...basePayload,
        ...(overrides?.createPayload ?? {}),
      };
      const result = mode === "edit" && orgId
        ? await bffClient.organisations.update(orgId, payload)
        : await bffClient.organisations.create(payload);
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
