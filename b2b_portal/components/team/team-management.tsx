"use client";

import { useCallback, useEffect, useState } from "react";
import { LuUserPlus, LuTrash2, LuLoader, LuUsers } from "react-icons/lu";
import DashboardLayout from "@/components/dashboard/dashboard-layout";
import { organisationClient } from "@lib/sdk/organisation-client";
import { authClient } from "@lib/sdk/auth-client";
import type { OrgMember } from "@lifeplus/insuretech-sdk";

function roleLabel(role: string | undefined): string {
  if (!role) return "Member";
  if (role.includes("BUSINESS_ADMIN")) return "Admin";
  if (role.includes("HR_MANAGER")) return "HR Manager";
  if (role.includes("VIEWER")) return "Viewer";
  return role.replace("ORG_MEMBER_ROLE_", "").replace(/_/g, " ");
}

export default function TeamManagementPage() {
  const [orgId, setOrgId] = useState("");
  const [members, setMembers] = useState<OrgMember[]>([]);
  const [loading, setLoading] = useState(true);
  const [addUserId, setAddUserId] = useState("");
  const [addRole, setAddRole] = useState<"ORG_MEMBER_ROLE_HR_MANAGER" | "ORG_MEMBER_ROLE_VIEWER">("ORG_MEMBER_ROLE_HR_MANAGER");
  const [adding, setAdding] = useState(false);
  const [error, setError] = useState("");

  // Resolve org_id from session
  useEffect(() => {
    authClient.getSession().then((res) => {
      const bizId = res.session?.principal.businessId ?? "";
      setOrgId(bizId);
    }).catch(() => setOrgId(""));
  }, []);

  const loadMembers = useCallback(async () => {
    if (!orgId) return;
    setLoading(true);
    setError("");
    try {
      const result = await organisationClient.listMembers(orgId);
      // Show only non-admin members (HR_MANAGER, VIEWER) on the team page
      const nonAdmins = (result.members ?? []).filter(
        (m) => !m.role?.includes("BUSINESS_ADMIN")
      );
      setMembers(result.ok ? nonAdmins : []);
      if (!result.ok) setError(result.message ?? "Failed to load team");
    } finally {
      setLoading(false);
    }
  }, [orgId]);

  useEffect(() => { loadMembers(); }, [loadMembers]);

  async function handleAdd() {
    if (!addUserId.trim()) { setError("User ID is required"); return; }
    if (!orgId) { setError("Organisation context not resolved"); return; }
    setAdding(true);
    setError("");
    try {
      const result = await organisationClient.addMember(orgId, addUserId.trim(), addRole);
      if (!result.ok) { setError(result.message ?? "Failed to add member"); return; }
      setAddUserId("");
      await loadMembers();
    } finally {
      setAdding(false);
    }
  }

  async function handleRemove(memberId: string) {
    if (!orgId) return;
    if (!confirm("Remove this team member?")) return;
    const result = await organisationClient.removeMember(orgId, memberId);
    if (!result.ok) { alert(result.message ?? "Remove failed"); return; }
    await loadMembers();
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        <div className="flex items-center gap-3">
          <LuUsers className="size-6 text-primary" />
          <div>
            <h1 className="text-xl font-semibold">Team Management</h1>
            <p className="text-sm text-muted-foreground">Manage HR Managers and Viewers for your organisation.</p>
          </div>
        </div>

        {/* Add member */}
        <div className="rounded-xl border bg-card p-5 space-y-3">
          <p className="text-sm font-semibold">Add Team Member</p>
          <div className="flex flex-wrap gap-2">
            <input
              className="flex-1 min-w-40 rounded-md border px-3 py-2 text-sm"
              placeholder="User ID"
              value={addUserId}
              onChange={(e) => setAddUserId(e.target.value)}
            />
            <select
              className="rounded-md border px-3 py-2 text-sm"
              value={addRole}
              onChange={(e) => setAddRole(e.target.value as typeof addRole)}
            >
              <option value="ORG_MEMBER_ROLE_HR_MANAGER">HR Manager</option>
              <option value="ORG_MEMBER_ROLE_VIEWER">Viewer</option>
            </select>
            <button
              onClick={handleAdd}
              disabled={adding}
              className="flex items-center gap-1.5 rounded-md bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
            >
              {adding ? <LuLoader className="size-4 animate-spin" /> : <LuUserPlus className="size-4" />}
              {adding ? "Adding…" : "Add Member"}
            </button>
          </div>
          {error && <p className="text-xs text-destructive">{error}</p>}
        </div>

        {/* Members table */}
        <div className="rounded-xl border bg-card overflow-hidden">
          <div className="border-b px-5 py-3 flex items-center justify-between">
            <p className="text-sm font-semibold">Team Members</p>
            <span className="text-xs text-muted-foreground">{members.length} member{members.length !== 1 ? "s" : ""}</span>
          </div>
          <table className="w-full text-sm">
            <thead className="bg-muted/50">
              <tr>
                <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground">User ID</th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground">Role</th>
                <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground">Status</th>
                <th className="px-4 py-3 text-right text-xs font-semibold text-muted-foreground">Actions</th>
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr>
                  <td colSpan={4} className="px-4 py-8 text-center text-muted-foreground">
                    <LuLoader className="inline-block animate-spin mr-2" />Loading…
                  </td>
                </tr>
              ) : members.length === 0 ? (
                <tr>
                  <td colSpan={4} className="px-4 py-8 text-center text-muted-foreground">No team members found. Add one above.</td>
                </tr>
              ) : (
                members.map((m) => (
                  <tr key={m.member_id} className="border-t hover:bg-muted/30">
                    <td className="px-4 py-3 font-medium">{m.user_id ?? "—"}</td>
                    <td className="px-4 py-3">
                      <span className="rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-semibold text-primary">
                        {roleLabel(m.role)}
                      </span>
                    </td>
                    <td className="px-4 py-3">
                      <span className={`rounded-full px-2.5 py-0.5 text-xs font-semibold ${
                        m.status === "ORG_MEMBER_STATUS_ACTIVE"
                          ? "bg-green-100 text-green-700"
                          : "bg-gray-100 text-gray-500"
                      }`}>
                        {m.status?.replace("ORG_MEMBER_STATUS_", "") ?? "Unknown"}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-right">
                      <button
                        onClick={() => handleRemove(m.member_id ?? "")}
                        className="rounded-md p-1.5 text-destructive hover:bg-destructive/10"
                        title="Remove"
                      >
                        <LuTrash2 className="size-4" />
                      </button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </DashboardLayout>
  );
}
