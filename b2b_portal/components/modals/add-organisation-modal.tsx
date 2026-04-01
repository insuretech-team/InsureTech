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
import { LuLoader, LuRefreshCw, LuSearch, LuShield, LuTrash2, LuUserPlus } from "react-icons/lu";
import { useOrganisationForm } from "@/src/hooks/useOrganisationForm";
import { ToastBanner } from "@/components/ui/toast-banner";
import { useToast } from "@/src/hooks/useToast";
import { bffClient } from "@lib/sdk/b2b-sdk-client";
import type { PortalUserLookup } from "@lib/sdk/b2b-sdk-client";
import type { OrgFormValues } from "@/src/hooks/useOrganisationForm";
import type { OrgMember } from "@lifeplus/insuretech-sdk";
import {
  BD_MOBILE_EXAMPLES,
  getBangladeshMobileValidationMessage,
  normalizeBangladeshMobile,
  normalizeBangladeshMobileOrRaw,
} from "@/src/lib/utils/bd-mobile";
import {
  getPasswordValidationMessage,
  PASSWORD_REQUIREMENTS_HINT,
} from "@/src/lib/utils/password";

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

type AdminDraftErrors = {
  email?: string;
  password?: string;
  mobileNumber?: string;
};

type ExistingAdminDraftErrors = {
  identifier?: string;
  temporaryPassword?: string;
};

type CreateAdminMode = "create" | "assign";

const focusPurple = "focus-visible:ring-primary focus-visible:border-primary focus-visible:ring-2";

const EMPTY_ADMIN_DRAFT: AdminDraft = {
  fullName: "",
  email: "",
  password: "",
  mobileNumber: "",
};

function validateAdminDraft(values: AdminDraft): AdminDraftErrors {
  const errors: AdminDraftErrors = {};
  if (!values.email.trim()) {
    errors.email = "Admin email is required.";
  }
  errors.password = getPasswordValidationMessage(values.password, "Admin password") ?? undefined;
  if (!values.mobileNumber.trim()) {
    errors.mobileNumber = "Admin mobile number is required.";
  } else if (!normalizeBangladeshMobile(values.mobileNumber.trim())) {
    errors.mobileNumber = getBangladeshMobileValidationMessage("Admin mobile number");
  }
  return errors;
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
            <p className="mt-1 text-xs text-muted-foreground">{PASSWORD_REQUIREMENTS_HINT}</p>
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
            onBlur={(e) => setField("adminMobileNumber", normalizeBangladeshMobileOrRaw(e.target.value))}
            className={`${focusPurple} ${errors.adminMobileNumber ? "border-red-500" : ""}`}
            required
          />
          {!errors.adminMobileNumber && (
            <p className="mt-1 text-xs text-muted-foreground">Accepted formats: {BD_MOBILE_EXAMPLES}.</p>
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
  isEdit,
}: {
  values: OrgFormValues;
  errors: Record<string, string | undefined>;
  setField: <K extends keyof OrgFormValues>(field: K, value: OrgFormValues[K]) => void;
  isEdit: boolean;
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
            placeholder={isEdit ? "Organisation Code" : "Generated Automatically"}
            value={values.code}
            className={focusPurple}
            readOnly
            disabled
          />
          <p className="mt-1 text-xs text-muted-foreground">
            Organisation code is generated automatically and cannot be edited here.
          </p>
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
            placeholder={`Contact Phone (${BD_MOBILE_EXAMPLES})`}
            value={values.contactPhone}
            onChange={(e) => setField("contactPhone", e.target.value)}
            onBlur={(e) => setField("contactPhone", normalizeBangladeshMobileOrRaw(e.target.value))}
            className={`${focusPurple} ${errors.contactPhone ? "border-red-500" : ""}`}
          />
          {!errors.contactPhone && (
            <p className="mt-1 text-xs text-muted-foreground">Accepted formats: {BD_MOBILE_EXAMPLES}.</p>
          )}
          {errors.contactPhone && <p className="mt-1 text-xs text-red-500">{errors.contactPhone}</p>}
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
  const [createAdminMode, setCreateAdminMode] = React.useState<CreateAdminMode>("create");
  const [members, setMembers] = React.useState<OrgMember[]>([]);
  const [membersLoading, setMembersLoading] = React.useState(false);
  const [memberActionId, setMemberActionId] = React.useState("");
  const [adminDraft, setAdminDraft] = React.useState<AdminDraft>(EMPTY_ADMIN_DRAFT);
  const [adminDraftErrors, setAdminDraftErrors] = React.useState<AdminDraftErrors>({});
  const [adminSubmitting, setAdminSubmitting] = React.useState(false);
  const [existingAdminIdentifier, setExistingAdminIdentifier] = React.useState("");
  const [existingAdminUser, setExistingAdminUser] = React.useState<PortalUserLookup | null>(null);
  const [existingAdminTemporaryPassword, setExistingAdminTemporaryPassword] = React.useState("");
  const [existingAdminErrors, setExistingAdminErrors] = React.useState<ExistingAdminDraftErrors>({});
  const [existingAdminSearching, setExistingAdminSearching] = React.useState(false);
  const [existingAdminSubmitting, setExistingAdminSubmitting] = React.useState(false);

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
      const result = await bffClient.organisations.listMembers(orgId);
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
    setCreateAdminMode("create");
    setAdminDraft(EMPTY_ADMIN_DRAFT);
    setAdminDraftErrors({});
    setExistingAdminIdentifier("");
    setExistingAdminUser(null);
    setExistingAdminTemporaryPassword("");
    setExistingAdminErrors({});
    if (isEdit) {
      void loadMembers();
    }
  }, [open, isEdit, loadMembers]);

  const handleCreateAdmin = React.useCallback(async () => {
    if (!orgId) return;
    const validationErrors = validateAdminDraft(adminDraft);
    setAdminDraftErrors(validationErrors);
    const firstValidationError = Object.values(validationErrors).find(Boolean);
    if (firstValidationError) {
      showToast("error", firstValidationError);
      return;
    }
    const normalizedMobileNumber = normalizeBangladeshMobile(adminDraft.mobileNumber.trim());
    setAdminSubmitting(true);
    try {
      const result = await bffClient.organisations.createAdmin(orgId, {
        fullName: adminDraft.fullName.trim() || undefined,
        email: adminDraft.email.trim(),
        password: adminDraft.password,
        mobileNumber: normalizedMobileNumber ?? adminDraft.mobileNumber.trim(),
      });
      if (!result.ok) {
        showToast("error", result.message ?? "Failed to create B2B admin");
        return;
      }
      showToast("success", result.message ?? "B2B admin created");
      setAdminDraft(EMPTY_ADMIN_DRAFT);
      setAdminDraftErrors({});
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
      const result = await bffClient.organisations.assignAdmin(orgId, memberId);
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
      const result = await bffClient.organisations.removeMember(orgId, memberId);
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

  const handleFindExistingUser = React.useCallback(async () => {
    const identifier = existingAdminIdentifier.trim();
    if (!identifier) {
      setExistingAdminErrors({ identifier: "Email or mobile number is required." });
      setExistingAdminUser(null);
      return;
    }

    setExistingAdminSearching(true);
    setExistingAdminErrors({});
    setExistingAdminUser(null);
    try {
      const lookupIdentifier = identifier.includes("@")
        ? identifier
        : normalizeBangladeshMobileOrRaw(identifier);
      const result = await bffClient.auth.findPortalUser(lookupIdentifier);
      if (!result.ok || !result.user?.userId) {
        setExistingAdminErrors({ identifier: result.message ?? "No matching user found." });
        return;
      }
      setExistingAdminIdentifier(lookupIdentifier);
      setExistingAdminUser(result.user);
    } finally {
      setExistingAdminSearching(false);
    }
  }, [existingAdminIdentifier]);

  const handleAssignExistingUser = React.useCallback(async () => {
    if (!orgId || !existingAdminUser?.userId) return;

    const nextErrors: ExistingAdminDraftErrors = {};
    if (!existingAdminIdentifier.trim()) {
      nextErrors.identifier = "Email or mobile number is required.";
    }
    nextErrors.temporaryPassword =
      getPasswordValidationMessage(existingAdminTemporaryPassword, "Temporary password") ?? undefined;
    setExistingAdminErrors(nextErrors);

    const firstError = Object.values(nextErrors).find(Boolean);
    if (firstError) {
      showToast("error", firstError);
      return;
    }

    setExistingAdminSubmitting(true);
    try {
      const result = await bffClient.organisations.assignExistingAdmin(
        orgId,
        existingAdminUser.userId,
        existingAdminTemporaryPassword
      );
      if (!result.ok) {
        showToast("error", result.message ?? "Failed to assign existing user as B2B admin");
        return;
      }
      showToast("success", result.message ?? "B2B admin assigned");
      setExistingAdminIdentifier("");
      setExistingAdminUser(null);
      setExistingAdminTemporaryPassword("");
      setExistingAdminErrors({});
      await loadMembers();
      onSaved?.();
    } finally {
      setExistingAdminSubmitting(false);
    }
  }, [
    existingAdminIdentifier,
    existingAdminTemporaryPassword,
    existingAdminUser,
    loadMembers,
    onSaved,
    orgId,
    showToast,
  ]);

  const handleCreateOrganisationSubmit = React.useCallback((e: React.FormEvent<HTMLFormElement>) => {
    if (createAdminMode === "create") {
      void submit(e);
      return;
    }

    e.preventDefault();
    if (!existingAdminUser?.userId) {
      setExistingAdminErrors({ identifier: "Find an existing user before creating the organisation." });
      showToast("error", "Find an existing user before creating the organisation.");
      return;
    }

    const passwordError =
      getPasswordValidationMessage(existingAdminTemporaryPassword, "Temporary password") ?? undefined;
    if (passwordError) {
      setExistingAdminErrors({ temporaryPassword: passwordError });
      showToast("error", passwordError);
      return;
    }

    void submit(e, {
      requireDefaultAdmin: false,
      createPayload: {
        adminAssignment: {
          userId: existingAdminUser.userId,
          temporaryPassword: existingAdminTemporaryPassword,
        },
      },
    });
  }, [
    createAdminMode,
    existingAdminTemporaryPassword,
    existingAdminUser,
    showToast,
    submit,
  ]);

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
          <form onSubmit={handleCreateOrganisationSubmit} className="space-y-6 px-6 py-6">
            <OrganisationFields values={values} errors={errors} setField={setField} isEdit={false} />
            <div className="space-y-4 rounded-lg border border-dashed border-primary/30 bg-primary/5 p-4">
              <div>
                <div className="text-sm font-semibold text-foreground">First B2B Admin</div>
                <div className="text-xs text-muted-foreground">
                  Choose whether to create a brand new admin or assign an existing portal user to this organisation.
                </div>
              </div>
              <Tabs value={createAdminMode} onValueChange={(value) => setCreateAdminMode(value as CreateAdminMode)}>
                <TabsList className="mb-4">
                  <TabsTrigger value="create">Create Admin</TabsTrigger>
                  <TabsTrigger value="assign">Assign Existing User</TabsTrigger>
                </TabsList>
                <TabsContent value="create" className="space-y-4">
                  <AdminFields values={values} errors={errors} setField={setField} />
                </TabsContent>
                <TabsContent value="assign" className="space-y-4">
                  <div className="space-y-4 rounded-lg border bg-background p-4">
                    <div>
                      <div className="text-sm font-semibold text-foreground">Assign Existing User</div>
                      <div className="text-xs text-muted-foreground">
                        Search by exact email or mobile number, then set the temporary password they will use for their first sign-in.
                      </div>
                    </div>
                    <div className="grid grid-cols-1 gap-4 md:grid-cols-[minmax(0,1fr)_auto]">
                      <div>
                        <Input
                          placeholder={`Email or mobile (${BD_MOBILE_EXAMPLES})`}
                          value={existingAdminIdentifier}
                          onChange={(e) => {
                            setExistingAdminIdentifier(e.target.value);
                            setExistingAdminErrors((prev) => ({ ...prev, identifier: undefined }));
                          }}
                          onBlur={(e) => {
                            const value = e.target.value;
                            if (!value.includes("@")) {
                              setExistingAdminIdentifier(normalizeBangladeshMobileOrRaw(value));
                            }
                          }}
                          className={`${focusPurple} ${existingAdminErrors.identifier ? "border-red-500" : ""}`}
                        />
                        <p className={`mt-1 text-xs ${existingAdminErrors.identifier ? "text-red-500" : "text-muted-foreground"}`}>
                          {existingAdminErrors.identifier ?? `Mobile accepts ${BD_MOBILE_EXAMPLES}, or use an email address.`}
                        </p>
                      </div>
                      <Button type="button" variant="outline" onClick={() => void handleFindExistingUser()} disabled={existingAdminSearching}>
                        {existingAdminSearching ? <LuLoader className="animate-spin" /> : <LuSearch />}
                        Find User
                      </Button>
                    </div>

                    {existingAdminUser ? (
                      <div className="space-y-3 rounded-lg border bg-background p-4">
                        <div className="space-y-1">
                          <div className="text-sm font-medium text-foreground">
                            {existingAdminUser.fullName || existingAdminUser.email || existingAdminUser.mobileNumber || existingAdminUser.userId}
                          </div>
                          <div className="text-xs text-muted-foreground">User ID: {existingAdminUser.userId}</div>
                          <div className="text-xs text-muted-foreground">
                            {existingAdminUser.email || "No email"} | {existingAdminUser.mobileNumber || "No mobile"}
                          </div>
                          <div className="text-xs text-muted-foreground">
                            Type: {existingAdminUser.userType || "UNKNOWN"} | KYC: {existingAdminUser.kycVerified ? "Verified" : "Required"}
                          </div>
                        </div>

                        <div>
                          <Input
                            placeholder="Temporary Password*"
                            type="password"
                            value={existingAdminTemporaryPassword}
                            onChange={(e) => {
                              setExistingAdminTemporaryPassword(e.target.value);
                              setExistingAdminErrors((prev) => ({ ...prev, temporaryPassword: undefined }));
                            }}
                            className={`${focusPurple} ${existingAdminErrors.temporaryPassword ? "border-red-500" : ""}`}
                          />
                          <p className={`mt-1 text-xs ${existingAdminErrors.temporaryPassword ? "text-red-500" : "text-muted-foreground"}`}>
                            {existingAdminErrors.temporaryPassword ?? "After sign-in, this user will be sent to eKYC first and then to the new-password page."}
                          </p>
                        </div>
                      </div>
                    ) : null}
                  </div>
                </TabsContent>
              </Tabs>
            </div>
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
                ) : createAdminMode === "assign" ? "Create Organisation & Assign Admin" : "Create Organisation"}
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
                <OrganisationFields values={values} errors={errors} setField={setField} isEdit />
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
                  <div className="text-xs text-muted-foreground">Assign an existing user with a temporary password, create a new admin, or promote an existing member.</div>
                </div>
                <Button type="button" variant="outline" size="sm" onClick={() => void loadMembers()} disabled={membersLoading}>
                  {membersLoading ? <LuLoader className="animate-spin" /> : <LuRefreshCw />}
                  Refresh
                </Button>
              </div>

              <div className="space-y-4 rounded-lg border border-dashed border-primary/30 bg-primary/5 p-4">
                <div>
                  <div className="text-sm font-semibold text-foreground">Assign Existing User</div>
                  <div className="text-xs text-muted-foreground">
                    Search by exact email or mobile number, then set the temporary password they will use for their first sign-in.
                  </div>
                </div>
                <div className="grid grid-cols-1 gap-4 md:grid-cols-[minmax(0,1fr)_auto]">
                  <div>
                    <Input
                      placeholder={`Email or mobile (${BD_MOBILE_EXAMPLES})`}
                      value={existingAdminIdentifier}
                      onChange={(e) => {
                        setExistingAdminIdentifier(e.target.value);
                        setExistingAdminErrors((prev) => ({ ...prev, identifier: undefined }));
                      }}
                      onBlur={(e) => {
                        const value = e.target.value;
                        if (!value.includes("@")) {
                          setExistingAdminIdentifier(normalizeBangladeshMobileOrRaw(value));
                        }
                      }}
                      className={`${focusPurple} ${existingAdminErrors.identifier ? "border-red-500" : ""}`}
                    />
                    <p className={`mt-1 text-xs ${existingAdminErrors.identifier ? "text-red-500" : "text-muted-foreground"}`}>
                      {existingAdminErrors.identifier ?? `Mobile accepts ${BD_MOBILE_EXAMPLES}, or use an email address.`}
                    </p>
                  </div>
                  <Button type="button" variant="outline" onClick={() => void handleFindExistingUser()} disabled={existingAdminSearching}>
                    {existingAdminSearching ? <LuLoader className="animate-spin" /> : <LuSearch />}
                    Find User
                  </Button>
                </div>

                {existingAdminUser ? (
                  <div className="space-y-3 rounded-lg border bg-background p-4">
                    <div className="space-y-1">
                      <div className="text-sm font-medium text-foreground">
                        {existingAdminUser.fullName || existingAdminUser.email || existingAdminUser.mobileNumber || existingAdminUser.userId}
                      </div>
                      <div className="text-xs text-muted-foreground">User ID: {existingAdminUser.userId}</div>
                      <div className="text-xs text-muted-foreground">
                        {existingAdminUser.email || "No email"} | {existingAdminUser.mobileNumber || "No mobile"}
                      </div>
                      <div className="text-xs text-muted-foreground">
                        Type: {existingAdminUser.userType || "UNKNOWN"} | KYC: {existingAdminUser.kycVerified ? "Verified" : "Required"}
                      </div>
                    </div>

                    <div>
                      <Input
                        placeholder="Temporary Password*"
                        type="password"
                        value={existingAdminTemporaryPassword}
                        onChange={(e) => {
                          setExistingAdminTemporaryPassword(e.target.value);
                          setExistingAdminErrors((prev) => ({ ...prev, temporaryPassword: undefined }));
                        }}
                        className={`${focusPurple} ${existingAdminErrors.temporaryPassword ? "border-red-500" : ""}`}
                      />
                      <p className={`mt-1 text-xs ${existingAdminErrors.temporaryPassword ? "text-red-500" : "text-muted-foreground"}`}>
                        {existingAdminErrors.temporaryPassword ?? "After sign-in, this user will be sent to eKYC first and then to the new-password page."}
                      </p>
                    </div>

                    <div className="flex justify-end">
                      <Button type="button" onClick={() => void handleAssignExistingUser()} disabled={existingAdminSubmitting}>
                        {existingAdminSubmitting ? <LuLoader className="animate-spin" /> : <LuShield />}
                        Assign Existing User
                      </Button>
                    </div>
                  </div>
                ) : null}
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
                    onChange={(e) => {
                      setAdminDraft((prev) => ({ ...prev, email: e.target.value }));
                      setAdminDraftErrors((prev) => ({ ...prev, email: undefined }));
                    }}
                    className={`${focusPurple} ${adminDraftErrors.email ? "border-red-500" : ""}`}
                  />
                </div>
                {adminDraftErrors.email && <p className="text-xs text-red-500">{adminDraftErrors.email}</p>}
                <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                  <div>
                    <Input
                      placeholder="Password*"
                      type="password"
                      value={adminDraft.password}
                      onChange={(e) => {
                        setAdminDraft((prev) => ({ ...prev, password: e.target.value }));
                        setAdminDraftErrors((prev) => ({ ...prev, password: undefined }));
                      }}
                      className={`${focusPurple} ${adminDraftErrors.password ? "border-red-500" : ""}`}
                    />
                    <p className={`mt-1 text-xs ${adminDraftErrors.password ? "text-red-500" : "text-muted-foreground"}`}>
                      {adminDraftErrors.password ?? PASSWORD_REQUIREMENTS_HINT}
                    </p>
                  </div>
                  <div>
                    <Input
                      placeholder="Mobile Number*"
                      value={adminDraft.mobileNumber}
                      onChange={(e) => {
                        setAdminDraft((prev) => ({ ...prev, mobileNumber: e.target.value }));
                        setAdminDraftErrors((prev) => ({ ...prev, mobileNumber: undefined }));
                      }}
                      onBlur={(e) => {
                        const nextValue = normalizeBangladeshMobileOrRaw(e.target.value);
                        setAdminDraft((prev) => ({ ...prev, mobileNumber: nextValue }));
                      }}
                      className={`${focusPurple} ${adminDraftErrors.mobileNumber ? "border-red-500" : ""}`}
                    />
                    <p className={`mt-1 text-xs ${adminDraftErrors.mobileNumber ? "text-red-500" : "text-muted-foreground"}`}>
                      {adminDraftErrors.mobileNumber ?? `Accepted formats: ${BD_MOBILE_EXAMPLES}.`}
                    </p>
                  </div>
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
