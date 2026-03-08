/**
 * organisation-client.ts
 * ───────────────────────
 * Browser-side client for /api/organisations.
 */
import { parseJson, type ApiResult } from "./api-client";
import type { Organisation } from "@lib/types/b2b";
import type { OrgMember } from "@lifeplus/insuretech-sdk";

export type OrgCreatePayload = {
  name: string;
  code?: string;
  industry?: string;
  contactEmail?: string;
  contactPhone?: string;
  address?: string;
  admin?: OrgAdminCreatePayload;
};

export type OrgUpdatePayload = Partial<OrgCreatePayload>;
export type OrgAdminCreatePayload = {
  email: string;
  password: string;
  fullName?: string;
  mobileNumber?: string;
};

export type OrgListResult = ApiResult<{ organisations: Organisation[] }>;
export type OrgSingleResult = ApiResult<{ organisation?: Organisation }>;
export type OrgMembersResult = ApiResult<{ members: OrgMember[] }>;
export type OrgMemberResult = ApiResult<{ member?: OrgMember }>;

export const organisationClient = {
  async list(): Promise<OrgListResult> {
    const res = await fetch("/api/organisations", { method: "GET", cache: "no-store" });
    return parseJson<OrgListResult>(res);
  },

  async get(id: string): Promise<OrgSingleResult> {
    const res = await fetch(`/api/organisations/${id}`, { method: "GET", cache: "no-store" });
    return parseJson<OrgSingleResult>(res);
  },

  async getMe(): Promise<OrgSingleResult> {
    const res = await fetch("/api/organisations/me", { method: "GET", cache: "no-store" });
    return parseJson<OrgSingleResult>(res);
  },

  async create(payload: OrgCreatePayload): Promise<OrgSingleResult> {
    const res = await fetch("/api/organisations", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    return parseJson<OrgSingleResult>(res);
  },

  async update(id: string, payload: OrgUpdatePayload): Promise<OrgSingleResult> {
    const res = await fetch(`/api/organisations/${id}`, {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    return parseJson<OrgSingleResult>(res);
  },

  async delete(id: string): Promise<ApiResult> {
    const res = await fetch(`/api/organisations/${id}`, {
      method: "DELETE",
    });
    return parseJson<ApiResult>(res);
  },

  async listMembers(id: string): Promise<OrgMembersResult> {
    const res = await fetch(`/api/organisations/${id}/members`, {
      method: "GET",
      cache: "no-store",
    });
    return parseJson<OrgMembersResult>(res);
  },

  /**
   * Promote an existing org member to B2B Admin role.
   * Uses /assign-admin (SDK assignOrgAdmin — takes member_id).
   * Do NOT use /admins for this — that endpoint creates a brand-new user.
   */
  async assignAdmin(id: string, memberId: string): Promise<OrgMemberResult> {
    const res = await fetch(`/api/organisations/${id}/assign-admin`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ memberId }),
    });
    return parseJson<OrgMemberResult>(res);
  },

  async createAdmin(id: string, payload: OrgAdminCreatePayload): Promise<OrgMemberResult> {
    const res = await fetch(`/api/organisations/${id}/admins`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    return parseJson<OrgMemberResult>(res);
  },

  async removeMember(id: string, memberId: string): Promise<ApiResult> {
    const res = await fetch(`/api/organisations/${id}/members/${memberId}`, {
      method: "DELETE",
    });
    return parseJson<ApiResult>(res);
  },

  /** Add a member (HR_MANAGER or VIEWER) to an org via /api/organisations/[id]/members POST */
  async addMember(id: string, userId: string, role: string): Promise<OrgMemberResult> {
    const res = await fetch(`/api/organisations/${id}/members`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ userId, role }),
    });
    return parseJson<OrgMemberResult>(res);
  },

  /** Assign an existing user as org admin via /api/organisations/[id]/assign-admin POST */
  async assignExistingAdmin(id: string, userId: string): Promise<ApiResult> {
    const res = await fetch(`/api/organisations/${id}/assign-admin`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ userId }),
    });
    return parseJson<ApiResult>(res);
  },

  /** Approve a pending org (sets status to ACTIVE) */
  async approve(id: string): Promise<OrgSingleResult> {
    const res = await fetch(`/api/organisations/${id}/approve`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
    });
    return parseJson<OrgSingleResult>(res);
  },
};
