/**
 * employee-client.ts
 * ──────────────────
 * Browser-side client for /api/employees. Used by client components and hooks.
 * All raw fetch calls go through this file — components never construct URLs.
 */
import { parseJson, type ApiResult } from "./api-client";
import type { Employee } from "@lib/types/b2b";

// ─── Payload types (mirrors what the API route expects) ───────────────────────

export type EmployeeCreatePayload = {
  name: string;
  employeeId: string;
  businessId: string;
  departmentId: string;
  email?: string;
  mobileNumber?: string;
  dateOfBirth?: string;
  dateOfJoining?: string;
  gender?: string;
  insuranceCategory?: number;
  assignedPlanId?: string;
  coverageAmount?: number;
  numberOfDependent?: number;
};

export type EmployeeUpdatePayload = Partial<EmployeeCreatePayload> & {
  status?: number;
};

export type EmployeeListResult = ApiResult<{ employees: Employee[]; total?: number }>;
/** Full employee record shape — returned by GET /api/employees/[id] with all form fields */
export type EmployeeFullRecord = {
  id: string; name: string; employeeID: string; department: string;
  insuranceCategory: number; assignedPlan: string; coverage: string;
  premiumAmount: string; status: "Active" | "Inactive"; numberOfDependent: number;
  // All form fields:
  email: string; mobileNumber: string; gender: string;
  dateOfBirth: string; dateOfJoining: string;
  departmentId: string; businessId: string;
  assignedPlanId: string; coverageAmount: string;
};
export type EmployeeSingleResult = ApiResult<{ employee?: EmployeeFullRecord }>;

// ─── Client ───────────────────────────────────────────────────────────────────

export const employeeClient = {
  async list(options?: {
    pageSize?: number;
    offset?: number;
    businessId?: string;
    departmentId?: string;
    status?: number;
  }): Promise<EmployeeListResult> {
    const params = new URLSearchParams({
      page_size: String(options?.pageSize ?? 50),
      offset: String(options?.offset ?? 0),
    });
    if (options?.businessId) params.set("business_id", options.businessId);
    if (options?.departmentId) params.set("department_id", options.departmentId);
    if (options?.status != null) params.set("status", String(options.status));
    const res = await fetch(`/api/employees?${params}`, { method: "GET", cache: "no-store" });
    return parseJson<EmployeeListResult>(res);
  },

  async get(id: string): Promise<EmployeeSingleResult> {
    const res = await fetch(`/api/employees/${id}`, { method: "GET", cache: "no-store" });
    return parseJson<EmployeeSingleResult>(res);
  },

  async create(payload: EmployeeCreatePayload): Promise<EmployeeSingleResult> {
    const res = await fetch("/api/employees", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    return parseJson<EmployeeSingleResult>(res);
  },

  async update(id: string, payload: EmployeeUpdatePayload): Promise<EmployeeSingleResult> {
    const res = await fetch(`/api/employees/${id}`, {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    return parseJson<EmployeeSingleResult>(res);
  },

  async delete(id: string): Promise<ApiResult> {
    const res = await fetch(`/api/employees/${id}`, { method: "DELETE" });
    return parseJson<ApiResult>(res);
  },
};
