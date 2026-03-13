"use client";

import { useCallback, useEffect, useState } from "react";
import { LuTrash2, LuLoader, LuUserPlus, LuShield, LuCopy, LuCheck } from "react-icons/lu";
import { organisationClient } from "@lib/sdk/organisation-client";
import type { OrgMember } from "@lifeplus/insuretech-sdk";

interface OrgMemberPanelProps {
  orgId: string;
  currentUserRole: "SYSTEM_ADMIN" | "B2B_ORG_ADMIN" | string;
}

function roleLabel(role: string | undefined): string {
  if (!role) return "Member";
  if (role.includes("BUSINESS_ADMIN") || role.includes("ADMIN")) return "B2B Admin";
  if (role.includes("HR_MANAGER") || role.includes("HR_STAFF")) return "HR Manager";
  if (role.includes("VIEWER")) return "Viewer";
  return role.replace("ORG_MEMBER_ROLE_", "").replace(/_/g, " ");
}

function roleBadgeClass(role: string | undefined): string {
  if (!role) return "bg-gray-100 text-gray-600";
  if (role.includes("ADMIN")) return "bg-purple-100 text-purple-700";
  if (role.includes("HR")) return "bg-blue-100 text-blue-700";
  if (role.includes("VIEWER")) return "bg-gray-100 text-gray-600";
  return "bg-primary/10 text-primary";
}

/** Shows a truncated UUID with a one-click copy button — the backend doesn't return names. */
function UserIdCell({ userId }: { userId: string | undefined }) {
  const [copied, setCopied] = useState(false);
  const id = userId ?? "";
  const short = id ? `${id.slice(0, 8)}…` : "—";

  function handleCopy() {
    if (!id) return;
    void navigator.clipboard.writeText(id).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 1500);
    });
  }

  return (
    <div className="flex items-center gap-1.5">
      <span className="font-mono text-xs text-foreground" title={id}>{short}</span>
      {id && (
        <button
          onClick={handleCopy}
          className="rounded p-0.5 text-muted-foreground hover:text-foreground hover:bg-muted"
          title="Copy full user ID"
        >
          {copied ? <LuCheck className="size-3 text-green-600" /> : <LuCopy className="size-3" />}
        </button>
      )}
    </div>
  );
}

function MemberRow({
  member,
  isSuperAdmin,
  onRemove,
  onPromote,
}: {
  member: OrgMember;
  isSuperAdmin: boolean;
  onRemove: (memberId: string) => void;
  onPromote: (memberId: string) => void;
}) {
  const [removing, setRemoving] = useState(false);
  const [promoting, setPromoting] = useState(false);
  const memberId = member.member_id ?? "";
  const isAdmin = (member.role ?? "").includes("ADMIN");

  async function handleRemove() {
    if (!confirm("Remove this member from the organisation?\nThey will lose all access immediately.")) return;
    setRemoving(true);
    try { onRemove(memberId); } finally { setRemoving(false); }
  }

  async function handlePromote() {
    if (!confirm("Promote this member to B2B Admin?\nThey will gain full organisation admin access.")) return;
    setPromoting(true);
    try { onPromote(memberId); } finally { setPromoting(false); }
  }

  return (
    <tr className="border-b last:border-0 hover:bg-muted/30">
      <td className="px-3 py-2">
        <UserIdCell userId={member.user_id} />
      </td>
      <td className="px-3 py-2">
        <span className={`rounded-full px-2 py-0.5 text-xs font-semibold ${roleBadgeClass(member.role)}`}>
          {roleLabel(member.role)}
        </span>
      </td>
      <td className="px-3 py-2">
        <span className={`rounded-full px-2 py-0.5 text-xs font-semibold ${
          member.status === "ORG_MEMBER_STATUS_ACTIVE"
            ? "bg-green-100 text-green-700"
            : "bg-gray-100 text-gray-500"
        }`}>
          {member.status?.replace("ORG_MEMBER_STATUS_", "") ?? "Unknown"}
        </span>
      </td>
      <td className="px-3 py-2 text-right">
        <div className="flex items-center justify-end gap-1">
          {/* Only Super Admin can promote to B2B Admin */}
          {isSuperAdmin && !isAdmin && (
            <button
              onClick={handlePromote}
              disabled={promoting}
              className="rounded-md p-1 text-purple-600 hover:bg-purple-50 disabled:opacity-50"
              title="Promote to B2B Admin"
            >
              {promoting ? <LuLoader className="size-4 animate-spin" /> : <LuShield className="size-4" />}
            </button>
          )}
          <button
            onClick={handleRemove}
            disabled={removing}
            className="rounded-md p-1 text-destructive hover:bg-destructive/10 disabled:opacity-50"
            title="Remove member"
          >
            {removing ? <LuLoader className="size-4 animate-spin" /> : <LuTrash2 className="size-4" />}
          </button>
        </div>
      </td>
    </tr>
  );
}

export function OrgMemberPanel({ orgId, currentUserRole }: OrgMemberPanelProps) {
  const [members, setMembers] = useState<OrgMember[]>([]);
  const [loading, setLoading] = useState(true);
  const [addUserId, setAddUserId] = useState("");
  const [addRole, setAddRole] = useState<"ORG_MEMBER_ROLE_HR_MANAGER" | "ORG_MEMBER_ROLE_VIEWER">("ORG_MEMBER_ROLE_HR_MANAGER");
  const [adding, setAdding] = useState(false);
  const [error, setError] = useState("");
  const isSuperAdmin = currentUserRole === "SYSTEM_ADMIN";

  const loadMembers = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const result = await organisationClient.listMembers(orgId);
      setMembers(result.ok ? (result.members ?? []) : []);
      if (!result.ok) setError(result.message ?? "Failed to load members");
    } finally {
      setLoading(false);
    }
  }, [orgId]);

  useEffect(() => { loadMembers(); }, [loadMembers]);

  async function handleRemove(memberId: string) {
    const result = await organisationClient.removeMember(orgId, memberId);
    if (!result.ok) { setError(result.message ?? "Remove failed"); return; }
    await loadMembers();
  }

  async function handlePromote(memberId: string) {
    // Promote existing member to admin using member_id (not user_id)
    const result = await organisationClient.assignAdmin(orgId, memberId);
    if (!result.ok) { setError(result.message ?? "Failed to promote member"); return; }
    await loadMembers();
  }

  async function handleAddMember() {
    if (!addUserId.trim()) { setError("User ID is required"); return; }
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

  return (
    <div className="space-y-4">
      {/* Add member form */}
      <div className="rounded-lg border p-3 space-y-2">
        <p className="text-xs font-semibold text-muted-foreground uppercase tracking-wide">Add Team Member</p>
        <div className="flex gap-2 flex-wrap">
          <input
            className="flex-1 min-w-32 rounded-md border px-3 py-1.5 text-sm"
            placeholder="User ID (UUID)"
            value={addUserId}
            onChange={(e) => setAddUserId(e.target.value)}
          />
          <select
            className="rounded-md border px-2 py-1.5 text-sm"
            value={addRole}
            onChange={(e) => setAddRole(e.target.value as typeof addRole)}
          >
            <option value="ORG_MEMBER_ROLE_HR_MANAGER">HR Manager</option>
            <option value="ORG_MEMBER_ROLE_VIEWER">Viewer</option>
          </select>
          <button
            onClick={handleAddMember}
            disabled={adding}
            className="flex items-center gap-1 rounded-md bg-primary px-3 py-1.5 text-xs font-semibold text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
          >
            {adding ? <LuLoader className="size-3 animate-spin" /> : <LuUserPlus className="size-3" />}
            Add
          </button>
        </div>
        {isSuperAdmin && (
          <p className="text-xs text-muted-foreground">
            💡 To create a new B2B Admin, use the <strong>B2B Admins</strong> tab in the Edit Organisation modal.
            To promote an existing member, click the <LuShield className="inline size-3" /> icon in the table below.
          </p>
        )}
      </div>

      {error && (
        <p className="rounded-md bg-destructive/10 px-3 py-2 text-xs text-destructive">{error}</p>
      )}

      {/* Members table */}
      <div className="rounded-lg border overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-muted/50">
            <tr>
              <th className="px-3 py-2 text-left text-xs font-semibold text-muted-foreground">User ID</th>
              <th className="px-3 py-2 text-left text-xs font-semibold text-muted-foreground">Role</th>
              <th className="px-3 py-2 text-left text-xs font-semibold text-muted-foreground">Status</th>
              <th className="px-3 py-2 text-right text-xs font-semibold text-muted-foreground">Actions</th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr><td colSpan={4} className="px-3 py-8 text-center text-muted-foreground text-sm">
                <LuLoader className="inline-block animate-spin mr-2" />Loading members…
              </td></tr>
            ) : members.length === 0 ? (
              <tr><td colSpan={4} className="px-3 py-8 text-center text-muted-foreground text-sm">
                No members yet. Add one above.
              </td></tr>
            ) : (
              members.map((m) => (
                <MemberRow
                  key={m.member_id ?? m.user_id}
                  member={m}
                  isSuperAdmin={isSuperAdmin}
                  onRemove={handleRemove}
                  onPromote={handlePromote}
                />
              ))
            )}
          </tbody>
        </table>
        <div className="border-t px-3 py-2 text-xs text-muted-foreground">
          {members.length} member{members.length !== 1 ? "s" : ""}
        </div>
      </div>
    </div>
  );
}
