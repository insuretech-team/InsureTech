"use client";

import * as React from "react";
import { useState } from "react";
import { LuX, LuBuilding2, LuUsers, LuLayoutDashboard, LuCheck, LuTrash2, LuLoader } from "react-icons/lu";
import { OrgMemberPanel } from "./org-member-panel";
import { organisationClient } from "@lib/sdk/organisation-client";
import type { Organisation } from "@lib/types/b2b";

interface OrgDetailPanelProps {
  org: Organisation | null;
  currentUserRole: string;
  open: boolean;
  onClose: () => void;
  onRefresh: () => void;
}

type Tab = "info" | "members" | "departments";

export function OrgDetailPanel({ org, currentUserRole, open, onClose, onRefresh }: OrgDetailPanelProps) {
  const [tab, setTab] = useState<Tab>("info");
  const [approving, setApproving] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [editName, setEditName] = useState(org?.name ?? "");
  const [editIndustry, setEditIndustry] = useState(org?.industry ?? "");
  const [editEmail, setEditEmail] = useState(org?.contactEmail ?? "");
  const [editPhone, setEditPhone] = useState(org?.contactPhone ?? "");
  const [editAddress, setEditAddress] = useState(org?.address ?? "");
  const [saving, setSaving] = useState(false);
  const [editError, setEditError] = useState("");

  // Sync edit fields whenever the selected org changes (panel re-opened for different org)
  // eslint-disable-next-line react-hooks/exhaustive-deps
  React.useEffect(() => {
    if (org) {
      setEditName(org.name ?? "");
      setEditIndustry(org.industry ?? "");
      setEditEmail(org.contactEmail ?? "");
      setEditPhone(org.contactPhone ?? "");
      setEditAddress(org.address ?? "");
      setEditError("");
    }
    setTab("info");
  }, [org?.id, open]);

  if (!open || !org) return null;

  const isPending = org.status === "Pending" || org.status === "ORGANISATION_STATUS_PENDING";
  const isSystemAdmin = currentUserRole === "SYSTEM_ADMIN";

  async function handleApprove() {
    if (!org) return;
    setApproving(true);
    try {
      const result = await organisationClient.approve(org.id);
      if (!result.ok) { alert(result.message ?? "Approve failed"); return; }
      onRefresh();
    } finally {
      setApproving(false);
    }
  }

  async function handleDelete() {
    if (!org) return;
    if (!confirm(`Delete organisation "${org.name}"? This cannot be undone.`)) return;
    setDeleting(true);
    try {
      const result = await organisationClient.delete(org.id);
      if (!result.ok) { alert(result.message ?? "Delete failed"); return; }
      onClose();
      onRefresh();
    } finally {
      setDeleting(false);
    }
  }

  async function handleSave() {
    if (!org) return;
    setSaving(true);
    setEditError("");
    try {
      const payload: Record<string, string> = {};
      if (editName.trim()) payload.name = editName.trim();
      if (editIndustry.trim()) payload.industry = editIndustry.trim();
      if (editEmail.trim()) payload.contactEmail = editEmail.trim();
      if (editPhone.trim()) payload.contactPhone = editPhone.trim();
      if (editAddress.trim()) payload.address = editAddress.trim();
      if (Object.keys(payload).length === 0) return;
      const result = await organisationClient.update(org.id, payload);
      if (!result.ok) { setEditError(result.message ?? "Update failed"); return; }
      onRefresh();
    } finally {
      setSaving(false);
    }
  }

  const tabs: { id: Tab; label: string; icon: React.ReactNode }[] = [
    { id: "info", label: "Info", icon: <LuBuilding2 className="size-4" /> },
    { id: "members", label: "Members", icon: <LuUsers className="size-4" /> },
    { id: "departments", label: "Departments", icon: <LuLayoutDashboard className="size-4" /> },
  ];

  return (
    <>
      {/* Backdrop */}
      <div className="fixed inset-0 z-40 bg-black/30 backdrop-blur-sm" onClick={onClose} />

      {/* Sheet */}
      <div className="fixed right-0 top-0 z-50 h-full w-full max-w-lg bg-background shadow-xl flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between border-b px-6 py-4">
          <div>
            <h2 className="text-lg font-semibold">{org.name}</h2>
            <p className="text-xs text-muted-foreground">{org.code} · {org.industry}</p>
          </div>
          <button onClick={onClose} className="rounded-md p-1.5 hover:bg-muted">
            <LuX className="size-5" />
          </button>
        </div>

        {/* Status bar + actions */}
        <div className="flex items-center gap-2 border-b px-6 py-2 bg-muted/30">
          <span className={`rounded-full px-2.5 py-0.5 text-xs font-semibold ${
            org.status === "Active" || org.status === "ORGANISATION_STATUS_ACTIVE"
              ? "bg-green-100 text-green-700"
              : org.status === "Pending" || org.status === "ORGANISATION_STATUS_PENDING"
              ? "bg-yellow-100 text-yellow-700"
              : "bg-gray-100 text-gray-500"
          }`}>
            {org.status?.replace("ORGANISATION_STATUS_", "") ?? "Unknown"}
          </span>
          <span className="text-xs text-muted-foreground">{org.totalEmployees ?? 0} employees</span>
          <div className="ml-auto flex gap-2">
            {isSystemAdmin && isPending && (
              <button
                onClick={handleApprove}
                disabled={approving}
                className="flex items-center gap-1 rounded-md bg-green-600 px-3 py-1.5 text-xs font-semibold text-white hover:bg-green-700 disabled:opacity-50"
              >
                {approving ? <LuLoader className="size-3 animate-spin" /> : <LuCheck className="size-3" />}
                Approve
              </button>
            )}
            {isSystemAdmin && (
              <button
                onClick={handleDelete}
                disabled={deleting}
                className="flex items-center gap-1 rounded-md bg-destructive px-3 py-1.5 text-xs font-semibold text-destructive-foreground hover:bg-destructive/90 disabled:opacity-50"
              >
                {deleting ? <LuLoader className="size-3 animate-spin" /> : <LuTrash2 className="size-3" />}
                Delete
              </button>
            )}
          </div>
        </div>

        {/* Tabs */}
        <div className="flex border-b px-6">
          {tabs.map((t) => (
            <button
              key={t.id}
              onClick={() => setTab(t.id)}
              className={`flex items-center gap-1.5 px-3 py-3 text-sm font-medium border-b-2 transition-colors ${
                tab === t.id
                  ? "border-primary text-primary"
                  : "border-transparent text-muted-foreground hover:text-foreground"
              }`}
            >
              {t.icon} {t.label}
            </button>
          ))}
        </div>

        {/* Tab content */}
        <div className="flex-1 overflow-y-auto px-6 py-4">
          {tab === "info" && (
            <div className="space-y-4">
              <div className="grid gap-3">
                {[
                  { label: "Name", setter: setEditName, stateVal: editName },
                  { label: "Industry", setter: setEditIndustry, stateVal: editIndustry },
                  { label: "Contact Email", setter: setEditEmail, stateVal: editEmail },
                  { label: "Contact Phone", setter: setEditPhone, stateVal: editPhone },
                  { label: "Address", setter: setEditAddress, stateVal: editAddress },
                ].map(({ label, setter, stateVal }) => (
                  <div key={label}>
                    <label className="text-xs font-semibold text-muted-foreground uppercase tracking-wide">{label}</label>
                    <input
                      className="mt-1 w-full rounded-md border px-3 py-1.5 text-sm"
                      value={stateVal}
                      placeholder={`Enter ${label.toLowerCase()}`}
                      onChange={(e) => setter(e.target.value)}
                    />
                  </div>
                ))}
              </div>
              {editError && <p className="text-xs text-destructive">{editError}</p>}
              <button
                onClick={handleSave}
                disabled={saving}
                className="flex items-center gap-1 rounded-md bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
              >
                {saving ? <LuLoader className="size-4 animate-spin" /> : null}
                {saving ? "Saving…" : "Save Changes"}
              </button>
            </div>
          )}

          {tab === "members" && (
            <OrgMemberPanel orgId={org.id} currentUserRole={currentUserRole} />
          )}

          {tab === "departments" && (
            <div className="space-y-3">
              <p className="text-sm text-muted-foreground">
                Departments for this organisation are managed on the Departments page.
              </p>
              <a
                href={`/departments?business_id=${org.id}`}
                className="inline-flex items-center gap-1 rounded-md bg-primary px-3 py-1.5 text-sm font-semibold text-primary-foreground hover:bg-primary/90"
              >
                <LuLayoutDashboard className="size-4" /> View Departments
              </a>
            </div>
          )}
        </div>
      </div>
    </>
  );
}
