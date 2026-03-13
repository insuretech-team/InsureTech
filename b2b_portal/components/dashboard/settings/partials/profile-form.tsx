"use client";

import React, { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Field, FieldGroup } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { authClient } from "@lib/sdk";

const focusPurple =
  "focus-visible:ring-primary focus-visible:border-primary focus-visible:ring-2";

type ProfileData = {
  full_name?: string;
  date_of_birth?: string;
  // address_line1 is the proto field name. The BFF also exposes an "address" alias
  // for the single-field form input — both are accepted on save.
  address_line1?: string;
  address?: string;
  // Read-only identity fields from the User record (auth credentials).
  // Injected by the BFF from portal_mobile / portal_email cookies set at login.
  email?: string;
  mobile_number?: string;
};

const ProfileForm = () => {
  const [profile, setProfile] = useState<ProfileData>({});
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState<{ text: string; ok: boolean } | null>(null);

  useEffect(() => {
    authClient.getProfile().then((res) => {
      if (res.ok && res.profile) {
        setProfile(res.profile as ProfileData);
      }
      setLoading(false);
    }).catch(() => setLoading(false));
  }, []);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setSaving(true);
    setMessage(null);
    // Strip read-only identity fields before sending to backend — they are not
    // updateable via the UserProfile table and would be silently ignored anyway.
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const { email, mobile_number, ...editableFields } = profile;
    try {
      const res = await authClient.updateProfile(editableFields as Record<string, unknown>);
      setMessage({ text: res.message ?? (res.ok ? "Profile updated." : "Update failed."), ok: res.ok });
      if (res.ok) {
        // Notify the header to re-fetch the profile and refresh the avatar + display name.
        window.dispatchEvent(new CustomEvent("profile:updated", { detail: editableFields }));
      }
    } catch {
      setMessage({ text: "An error occurred. Please try again.", ok: false });
    } finally {
      setSaving(false);
    }
  };

  const set = (key: keyof ProfileData) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setProfile((p) => ({ ...p, [key]: e.target.value }));

  return (
    <Card>
      <form className="py-3" onSubmit={handleSubmit}>
        <CardHeader>
          <CardTitle>My Profile</CardTitle>
        </CardHeader>
        <CardContent className="text-muted-foreground text-sm">
          {loading ? (
            <p className="text-muted-foreground py-4">Loading profile…</p>
          ) : (
            <FieldGroup className="space-y-4 gap-0">

              {/* ── Identity (read-only) ─────────────────────────────────────────
                  mobile_number and email are auth credentials from the User record.
                  They cannot be changed from this form — contact support to update. */}
              {(profile.mobile_number || profile.email) && (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 rounded-md border border-dashed border-border bg-muted/40 px-4 py-3">
                  {profile.mobile_number && (
                    <Field>
                      <Label htmlFor="mobile_number_ro" className="text-muted-foreground text-xs uppercase tracking-wide">
                        Mobile Number{" "}
                        <span className="ml-1 text-xs font-normal normal-case text-muted-foreground">
                          (login credential — read only)
                        </span>
                      </Label>
                      <Input
                        id="mobile_number_ro"
                        value={profile.mobile_number}
                        readOnly
                        disabled
                        className="cursor-not-allowed bg-muted/60 text-muted-foreground"
                      />
                    </Field>
                  )}
                  {profile.email && (
                    <Field>
                      <Label htmlFor="email_ro" className="text-muted-foreground text-xs uppercase tracking-wide">
                        Email{" "}
                        <span className="ml-1 text-xs font-normal normal-case text-muted-foreground">
                          (login credential — read only)
                        </span>
                      </Label>
                      <Input
                        id="email_ro"
                        type="email"
                        value={profile.email}
                        readOnly
                        disabled
                        className="cursor-not-allowed bg-muted/60 text-muted-foreground"
                      />
                    </Field>
                  )}
                </div>
              )}

              {/* ── Editable profile fields ──────────────────────────────────── */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <Field>
                  <Label htmlFor="full_name">Full Name</Label>
                  <Input
                    id="full_name"
                    placeholder="Full name"
                    className={focusPurple}
                    value={profile.full_name ?? ""}
                    onChange={set("full_name")}
                  />
                </Field>
                <Field>
                  <Label htmlFor="date_of_birth">Date of Birth</Label>
                  <Input
                    id="date_of_birth"
                    type="date"
                    className={focusPurple}
                    value={profile.date_of_birth ?? ""}
                    onChange={set("date_of_birth")}
                  />
                </Field>
              </div>

              <div className="grid grid-cols-1 gap-4">
                <Field>
                  <Label htmlFor="address_line1">Address</Label>
                  <Input
                    id="address_line1"
                    placeholder="Address line 1"
                    className={focusPurple}
                    value={profile.address_line1 ?? profile.address ?? ""}
                    onChange={(e) => setProfile((p) => ({ ...p, address_line1: e.target.value, address: e.target.value }))}
                  />
                </Field>
              </div>

              {message && (
                <p className={`text-sm ${message.ok ? "text-green-600" : "text-red-500"}`}>
                  {message.text}
                </p>
              )}

              <div className="flex items-center justify-end mt-4">
                <Button
                  type="submit"
                  variant="default"
                  className="bg-primary hover:bg-accent"
                  disabled={saving}
                >
                  {saving ? "Saving…" : "Save Changes"}
                </Button>
              </div>
            </FieldGroup>
          )}
        </CardContent>
      </form>
    </Card>
  );
};

export default ProfileForm;
