/**
 * add-organisation-modal.tsx
 * Create / Edit organisation modal dialog with admin management.
 */
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
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { LuLoader, LuRefreshCw, LuShield, LuTrash2, LuUserPlus } from "react-icons/lu";
import { useOrganisationForm } from "@/src/hooks/useOrganisationForm";
import { ToastBanner } from "@/components/ui/toast-banner";
import { useToast } from "@/src/hooks/useToast";
import { organisationClient } from "@lib/sdk/organisation-client";
import type { OrgFormValues } from "@/src/hooks/useOrganisationForm";
import type { OrgMember } from "@lifeplus/insuretech-sdk";

type Props = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  orgId?: string;
  initialValues?: Partial<OrgFormValues>;
  onSaved?: () => void;
};

type AdminDraft = {
  fullName: string;
  email: string;
  password: string;
  mobileNumber: string;
};

const focusPurple = "focus-visible:ring-primary focus-visible:border-primary focus-visible:ring-2";

const EMPTY_ADMIN_DRAFT: AdminDraft = {
  fullName: "",
  email: "",
  password: "",
  mobileNumber: "",
};

function looksLikeBangladeshMobile(raw: string): boolean {
  const digits = raw.replace(/\D/g, "");
  return /^(?:880|0)?1[3-9]\d{8}$/.test(digits) || /^1[3-9]\d{8}$/.test(digits);
}

function validateAdminDraft(values: AdminDraft): string | null {
  if (!values.email.trim() || !values.password.trim() || !values.mobileNumber.trim()) {
    return "Email, password, and mobile number are required";
  }

  const hasUpper = /[A-Z]/.test(values.password);
  const hasLower = /[a-z]/.test(values.password);
  const hasDigit = /\d/.test(values.password);
  const hasSymbol = /[^A-Za-z0-9]/.test(values.password);
  if (values.password.length < 8 || !hasUpper || !hasLower || !hasDigit || !hasSymbol) {
    return "Use 8+ chars with upper, lower, number, and symbol";
  }

  if (!looksLikeBangladeshMobile(values.mobileNumber.trim())) {
    return "Use a valid BD mobile number";
  }

  return null;
}

function readMemberString(member: OrgMember, ...keys: string[]) {
  const bag = member as unknown as Record<string, unknown>;
  for (const key of keys) {
    const value = bag[key];
    if (typeof value === "string" && value.trim()) return value;
  }
  return "";
}

function roleLabel(member: OrgMember) {
  const role = readMemberString(member, "role");
  if (role === "ORG_MEMBER_ROLE_BUSINESS_ADMIN" || role === "ORG_MEMBER_ROLE_ADMIN") return "B2B Admin";
  if (role === "ORG_MEMBER_ROLE_HR_STAFF") return "HR Staff";
  if (role === "ORG_MEMBER_ROLE_EMPLOYEE") return "Employee";
  return role || "Unknown";
}

function statusLabel(member: OrgMember) {
  const status = readMemberString(member, "status");
  if (status === "ORG_MEMBER_STATUS_ACTIVE") return "Active";
  if (status === "ORG_MEMBER_STATUS_INACTIVE") return "Inactive";
  return status || "Unknown";
}

function AdminFields({
  values,
  errors,
  setField,
}: {
  values: OrgFormValues;
  errors: Record<string, string | undefined>;
  setField: <K extends keyof OrgFormValues>(field: K, value: OrgFormValues[K]) => void;
}) {
  return (
    <div className="space-y-4 rounded-lg border border-dashed border-primary/30 bg-primary/5 p-4">
      <div>
        <div className="text-sm font-semibold text-foreground">Primary B2B Admin</div>
        <div className="text-xs text-muted-foreground">Super admin creates the organisation and its first portal admin together.</div>
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        <Field>
          <Label htmlFor="admin-name" className="sr-only">Admin Full Name</Label>
          <Input
            id="admin-name"
            placeholder="Admin Full Name"
            value={values.adminFullName}
            onChange={(e) => setField("adminFullName", e.target.value)}
            className={focusPurple}
          />
        </Field>

        <Field>
          <Label htmlFor="admin-email" className="sr-only">Admin Email</Label>
          <Input
            id="admin-email"
            type="email"
            placeholder="Admin Email*"
            value={values.adminEmail}
            onChange={(e) => setField("adminEmail", e.target.value)}
            className={`${focusPurple} ${errors.adminEmail ? "border-red-500" : ""}`}
            required
          />
          {errors.adminEmail && <p className="mt-1 text-xs text-red-500">{errors.adminEmail}</p>}
        </Field>
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        <Field>
          <Label htmlFor="admin-password" className="sr-only">Admin Password</Label>
          <Input
            id="admin-password"
            type="password"
            placeholder="Admin Password*"
            value={values.adminPassword}
            onChange={(e) => setField("adminPassword", e.target.value)}
            className={`${focusPurple} ${errors.adminPassword ? "border-red-500" : ""}`}
            required
          />
          {!errors.adminPassword && (
            <p className="mt-1 text-xs text-muted-foreground">Use 8+ characters with uppercase, lowercase, number, and symbol.</p>
          )}
          {errors.adminPassword && <p className="mt-1 text-xs text-red-500">{errors.adminPassword}</p>}
        </Field>

        <Field>
          <Label htmlFor="admin-mobile" className="sr-only">Admin Mobile Number</Label>
          <Input
            id="admin-mobile"
            placeholder="Admin Mobile Number*"
            value={values.adminMobileNumber}
            onChange={(e) => setField("adminMobileNumber", e.target.value)}
            className={`${focusPurple} ${errors.adminMobileNumber ? "border-red-500" : ""}`}
            required
          />
          {!errors.adminMobileNumber && (
            <p className="mt-1 text-xs text-muted-foreground">Accepted formats: `01712345678`, `8801712345678`, `+8801712345678`.</p>
          )}
          {errors.adminMobileNumber && <p className="mt-1 text-xs text-red-500">{errors.adminMobileNumber}</p>}
        </Field>
      </div>
    </div>
  );
}

function OrganisationFields({
  values,
  errors,
  setField,
}: {
  values: OrgFormValues;
  errors: Record<string, string | undefined>;
  setField: <K extends keyof OrgFormValues>(field: K, value: OrgFormValues[K]) => void;
}) {
  return (
    <FieldGroup className="space-y-4 gap-0">
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        <Field>
          <Label htmlFor="org-name" className="sr-only">Organisation Name</Label>
          <Input
            id="org-name"
            placeholder="Organisation Name*"
            value={values.name}
            onChange={(e) => setField("name", e.target.value)}
            className={`${focusPurple} ${errors.name ? "border-red-500" : ""}`}
            required
          />
          {errors.name && <p className="mt-1 text-xs text-red-500">{errors.name}</p>}
        </Field>

        <Field>
          <Label htmlFor="org-code" className="sr-only">Organisation Code</Label>
          <Input
            id="org-code"
            placeholder="Organisation Code"
            value={values.code}
            onChange={(e) => setField("code", e.target.value)}
            className={focusPurple}
          />
        </Field>
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        <Field>
          <Label htmlFor="org-industry" className="sr-only">Industry</Label>
          <Input
            id="org-industry"
            placeholder="Industry"
            value={values.industry}
            onChange={(e) => setField("industry", e.target.value)}
            className={focusPurple}
          />
        </Field>

        <Field>
          <Label htmlFor="org-email" className="sr-only">Contact Email</Label>
          <Input
            id="org-email"
            type="email"
            placeholder="Contact Email"
            value={values.contactEmail}
            onChange={(e) => setField("contactEmail", e.target.value)}
            className={focusPurple}
          />
        </Field>
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        <Field>
          <Label htmlFor="org-phone" className="sr-only">Contact Phone</Label>
          <Input
            id="org-phone"
            placeholder="Contact Phone"
            value={values.contactPhone}
            onChange={(e) => setField("contactPhone", e.target.value)}
            className={focusPurple}
          />
        </Field>

        <Field>
          <Label htmlFor="org-address" className="sr-only">Address</Label>
          <Input
            id="org-address"
            placeholder="Address"
            value={values.address}
            onChange={(e) => setField("address", e.target.value)}
            className={focusPurple}
          />
        </Field>
      </div>
    </FieldGroup>
  );
}

export default function AddOrganisationModal({ open, onOpenChange, orgId, initialValues, onSaved }: Props) {
  const isEdit = Boolean(orgId);
  const { toast, showToast } = useToast();
  const [activeTab, setActiveTab] = React.useState("organisation");
  const [members, setMembers] = React.useState<OrgMember[]>([]);
  const [membersLoading, setMembersLoading] = React.useState(false);
  const [memberActionId, setMemberActionId] = React.useState("");
  const [adminDraft, setAdminDraft] = React.useState<AdminDraft>(EMPTY_ADMIN_DRAFT);
  const [adminSubmitting, setAdminSubmitting] = React.useState(false);

  const { values, errors, submitting, setField, submit } = useOrganisationForm({
    mode: isEdit ? "edit" : "create",
    orgId,
    initialValues,
    onSuccess: (msg) => {
      showToast("success", msg);
      setTimeout(() => { onSaved?.(); onOpenChange(false); }, 1200);
    },
    onError: (msg) => showToast("error", msg),
  });

  const loadMembers = React.useCallback(async () => {
    if (!orgId) return;
    setMembersLoading(true);
    try {
      const result = await organisationClient.listMembers(orgId);
      if (!result.ok) {
        showToast("error", result.message ?? "Failed to load organisation members");
        return;
      }
      setMembers(result.members ?? []);
    } finally {
      setMembersLoading(false);
    }
  }, [orgId, showToast]);

  React.useEffect(() => {
    if (!open) return;
    setActiveTab("organisation");
    setAdminDraft(EMPTY_ADMIN_DRAFT);
    if (isEdit) {
      void loadMembers();
    }
  }, [open, isEdit, loadMembers]);

  const handleCreateAdmin = React.useCallback(async () => {
    if (!orgId) return;
    const validationMessage = validateAdminDraft(adminDraft);
    if (validationMessage) {
      showToast("error", validationMessage);
      return;
    }
    setAdminSubmitting(true);
    try {
      const result = await organisationClient.createAdmin(orgId, {
        fullName: adminDraft.fullName.trim() || undefined,
        email: adminDraft.email.trim(),
        password: adminDraft.password,
        mobileNumber: adminDraft.mobileNumber.trim(),
      });
      if (!result.ok) {
        showToast("error", result.message ?? "Failed to create B2B admin");
        return;
      }
      showToast("success", result.message ?? "B2B admin created");
      setAdminDraft(EMPTY_ADMIN_DRAFT);
      await loadMembers();
      onSaved?.();
    } finally {
      setAdminSubmitting(false);
    }
  }, [adminDraft, loadMembers, onSaved, orgId, showToast]);

  const handlePromote = React.useCallback(async (memberId: string) => {
    if (!orgId || !memberId) return;
    setMemberActionId(memberId);
    try {
      const result = await organisationClient.assignAdmin(orgId, memberId);
      if (!result.ok) {
        showToast("error", result.message ?? "Failed to assign admin");
        return;
      }
      showToast("success", result.message ?? "Admin assigned");
      await loadMembers();
      onSaved?.();
    } finally {
      setMemberActionId("");
    }
  }, [loadMembers, onSaved, orgId, showToast]);

  const handleRemove = React.useCallback(async (memberId: string) => {
    if (!orgId || !memberId) return;
    if (!confirm("Remove this organisation member?")) return;
    setMemberActionId(memberId);
    try {
      const result = await organisationClient.removeMember(orgId, memberId);
      if (!result.ok) {
        showToast("error", result.message ?? "Failed to remove member");
        return;
      }
      showToast("success", result.message ?? "Member removed");
      await loadMembers();
      onSaved?.();
    } finally {
      setMemberActionId("");
    }
  }, [loadMembers, onSaved, orgId, showToast]);

  const renderMembers = () => {
    if (membersLoading) {
      return <div className="flex items-center gap-2 text-sm text-muted-foreground"><LuLoader className="animate-spin" /> Loading members…</div>;
    }
    if (members.length === 0) {
      return <div className="rounded-md border border-dashed p-4 text-sm text-muted-foreground">No organisation members yet.</div>;
    }

    return (
      <div className="space-y-3">
        {members.map((member, index) => {
          const memberId = readMemberString(member, "member_id", "memberId");
          const userId = readMemberString(member, "user_id", "userId");
          const role = readMemberString(member, "role");
          const busy = memberActionId === memberId;
          return (
            <div key={memberId || `${userId}-${index}`} className="flex flex-col gap-3 rounded-lg border p-4 md:flex-row md:items-center md:justify-between">
              <div className="space-y-1">
                <div className="flex items-center gap-2">
                  <span className="text-sm font-medium text-foreground font-mono" title={userId}>
                    {userId ? `${userId.slice(0, 8)}…` : "Unknown user"}
                  </span>
                  {userId && (
                    <button
                      type="button"
                      className="rounded p-0.5 text-muted-foreground hover:text-foreground hover:bg-muted"
                      title="Copy user ID"
                      onClick={() => void navigator.clipboard.writeText(userId)}
                    >
                      <span className="text-xs">⎘</span>
                    </button>
                  )}
                </div>
                <div className="text-xs text-muted-foreground">Role: {roleLabel(member)} | Status: {statusLabel(member)}</div>
              </div>
              <div className="flex gap-2">
                {role !== "ORG_MEMBER_ROLE_BUSINESS_ADMIN" && role !== "ORG_MEMBER_ROLE_ADMIN" && (
                  <Button type="button" variant="outline" size="sm" disabled={busy} onClick={() => handlePromote(memberId)}>
                    {busy ? <LuLoader className="animate-spin" /> : <LuShield />}
                    Make Admin
                  </Button>
                )}
                <Button type="button" variant="outline" size="sm" disabled={busy} onClick={() => handleRemove(memberId)}>
                  {busy ? <LuLoader className="animate-spin" /> : <LuTrash2 />}
                  Remove
                </Button>
              </div>
            </div>
          );
        })}
      </div>
    );
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[90vh] overflow-y-auto p-0 sm:max-w-3xl">
        <DialogHeader className="sticky top-0 z-10 border-b bg-white px-6 py-4">
          <DialogTitle className="text-xl font-semibold">
            {isEdit ? "Edit Organisation" : "Add Organisation"}
          </DialogTitle>
        </DialogHeader>

        <ToastBanner toast={toast} />

        {!isEdit ? (
          <form onSubmit={submit} className="space-y-6 px-6 py-6">
            <OrganisationFields values={values} errors={errors} setField={setField} />
            <AdminFields values={values} errors={errors} setField={setField} />
            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => onOpenChange(false)} disabled={submitting}>
                Cancel
              </Button>
              <Button type="submit" disabled={submitting} className="h-11 px-8 text-white bg-gradient-to-r from-primary to-accent hover:opacity-95">
                {submitting ? (
                  <span className="flex items-center gap-2">
                    <LuLoader className="animate-spin" />
                    Creating…
                  </span>
                ) : "Create Organisation"}
              </Button>
            </DialogFooter>
          </form>
        ) : (
          <Tabs value={activeTab} onValueChange={setActiveTab} className="px-6 py-6">
            <TabsList className="mb-6">
              <TabsTrigger value="organisation">Organisation</TabsTrigger>
              <TabsTrigger value="admins">B2B Admins</TabsTrigger>
            </TabsList>

            <TabsContent value="organisation">
              <form onSubmit={submit} className="space-y-6">
                <OrganisationFields values={values} errors={errors} setField={setField} />
                <DialogFooter>
                  <Button type="button" variant="outline" onClick={() => onOpenChange(false)} disabled={submitting}>
                    Cancel
                  </Button>
                  <Button type="submit" disabled={submitting} className="h-11 px-8 text-white bg-gradient-to-r from-primary to-accent hover:opacity-95">
                    {submitting ? (
                      <span className="flex items-center gap-2">
                        <LuLoader className="animate-spin" />
                        Saving…
                      </span>
                    ) : "Save Changes"}
                  </Button>
                </DialogFooter>
              </form>
            </TabsContent>

            <TabsContent value="admins" className="space-y-6">
              <div className="flex items-center justify-between">
                <div>
                  <div className="text-sm font-semibold text-foreground">Organisation Admins</div>
                  <div className="text-xs text-muted-foreground">Create a new B2B admin or promote an existing member.</div>
                </div>
                <Button type="button" variant="outline" size="sm" onClick={() => void loadMembers()} disabled={membersLoading}>
                  {membersLoading ? <LuLoader className="animate-spin" /> : <LuRefreshCw />}
                  Refresh
                </Button>
              </div>

              <div className="space-y-4 rounded-lg border border-dashed border-primary/30 bg-primary/5 p-4">
                <div className="text-sm font-semibold text-foreground">Create New B2B Admin</div>
                <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                  <Input
                    placeholder="Full Name"
                    value={adminDraft.fullName}
                    onChange={(e) => setAdminDraft((prev) => ({ ...prev, fullName: e.target.value }))}
                    className={focusPurple}
                  />
                  <Input
                    placeholder="Email*"
                    type="email"
                    value={adminDraft.email}
                    onChange={(e) => setAdminDraft((prev) => ({ ...prev, email: e.target.value }))}
                    className={focusPurple}
                  />
                </div>
                <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                  <Input
                    placeholder="Password*"
                    type="password"
                    value={adminDraft.password}
                    onChange={(e) => setAdminDraft((prev) => ({ ...prev, password: e.target.value }))}
                    className={focusPurple}
                  />
                  <Input
                    placeholder="Mobile Number*"
                    value={adminDraft.mobileNumber}
                    onChange={(e) => setAdminDraft((prev) => ({ ...prev, mobileNumber: e.target.value }))}
                    className={focusPurple}
                  />
                </div>
                <div className="flex justify-end">
                  <Button type="button" onClick={() => void handleCreateAdmin()} disabled={adminSubmitting}>
                    {adminSubmitting ? <LuLoader className="animate-spin" /> : <LuUserPlus />}
                    Create B2B Admin
                  </Button>
                </div>
              </div>

              {renderMembers()}
            </TabsContent>
          </Tabs>
        )}
      </DialogContent>
    </Dialog>
  );
}
