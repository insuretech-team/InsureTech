"use client";

import React, { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Field, FieldGroup } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { organisationClient } from "@lib/sdk";
import type { Organisation } from "@lib/types/b2b";

const focusPurple =
  "focus-visible:ring-primary focus-visible:border-primary focus-visible:ring-2";

const INDUSTRY_OPTIONS = [
  "Technology", "Healthcare", "Finance", "Education", "Manufacturing",
  "Retail", "Construction", "Transportation", "Government", "Other",
];

type OrgFormData = {
  name: string;
  industry: string;
  contactEmail: string;
  contactPhone: string;
  address: string;
};

const OrganizationForm = () => {
  const [orgId, setOrgId] = useState<string | null>(null);
  const [form, setForm] = useState<OrgFormData>({
    name: "", industry: "", contactEmail: "", contactPhone: "", address: "",
  });
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState<{ text: string; ok: boolean } | null>(null);
  const [isSystemAdmin, setIsSystemAdmin] = useState(false);

  useEffect(() => {
    organisationClient.getMe().then((res) => {
      if (res.ok && res.organisation) {
        const org: Organisation = res.organisation;
        setOrgId(org.id);
        setForm({
          name: org.name ?? "",
          industry: org.industry ?? "",
          contactEmail: org.contactEmail ?? "",
          contactPhone: org.contactPhone ?? "",
          address: org.address ?? "",
        });
      } else if (res.ok && !res.organisation) {
        // SYSTEM_ADMIN — no org context
        setIsSystemAdmin(true);
      }
      setLoading(false);
    }).catch(() => setLoading(false));
  }, []);

  const set = (key: keyof OrgFormData) => (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) =>
    setForm((f) => ({ ...f, [key]: e.target.value }));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!orgId) return;
    setSaving(true);
    setMessage(null);
    try {
      const res = await organisationClient.update(orgId, {
        name: form.name,
        industry: form.industry,
        contactEmail: form.contactEmail,
        contactPhone: form.contactPhone,
        address: form.address,
      });
      setMessage({
        text: res.message ?? (res.ok ? "Organisation updated successfully." : "Update failed."),
        ok: res.ok,
      });
    } catch {
      setMessage({ text: "An error occurred. Please try again.", ok: false });
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <Card>
        <CardContent className="py-8">
          <p className="text-sm text-muted-foreground">Loading organisation…</p>
        </CardContent>
      </Card>
    );
  }

  if (isSystemAdmin) {
    return (
      <Card>
        <CardContent className="py-8">
          <p className="text-sm text-muted-foreground">
            System administrators are not associated with a specific organisation.
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <form className="py-3" onSubmit={handleSubmit}>
        <CardHeader>
          <CardTitle>Organization Info.</CardTitle>
        </CardHeader>
        <CardContent className="text-muted-foreground text-sm">
          <FieldGroup className="space-y-4 gap-0">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Field>
                <Label htmlFor="organizationName">Organization Name *</Label>
                <Input
                  id="organizationName"
                  placeholder="Organization name"
                  className={focusPurple}
                  value={form.name}
                  onChange={set("name")}
                  required
                />
              </Field>
              <Field>
                <Label htmlFor="industry">Industry</Label>
                <select
                  id="industry"
                  name="industry"
                  className={`flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors placeholder:text-muted-foreground ${focusPurple} focus:outline-none disabled:cursor-not-allowed disabled:opacity-50`}
                  value={form.industry}
                  onChange={set("industry")}
                >
                  <option value="">Select industry</option>
                  {INDUSTRY_OPTIONS.map((opt) => (
                    <option key={opt} value={opt}>{opt}</option>
                  ))}
                </select>
              </Field>
            </div>

            <CardHeader className="px-0 pt-4 pb-2">
              <CardTitle>Contact Information</CardTitle>
            </CardHeader>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Field>
                <Label htmlFor="contactEmail">Contact Email</Label>
                <Input
                  id="contactEmail"
                  type="email"
                  placeholder="contact@company.com"
                  className={focusPurple}
                  value={form.contactEmail}
                  onChange={set("contactEmail")}
                />
              </Field>
              <Field>
                <Label htmlFor="contactPhone">Contact Phone</Label>
                <Input
                  id="contactPhone"
                  placeholder="e.g. 01712345678"
                  className={focusPurple}
                  value={form.contactPhone}
                  onChange={set("contactPhone")}
                />
              </Field>
            </div>

            <div className="grid grid-cols-1 gap-4">
              <Field>
                <Label htmlFor="address">Address</Label>
                <Input
                  id="address"
                  placeholder="Office address"
                  className={focusPurple}
                  value={form.address}
                  onChange={set("address")}
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
                disabled={saving || !orgId}
              >
                {saving ? "Saving…" : "Save Changes"}
              </Button>
            </div>
          </FieldGroup>
        </CardContent>
      </form>
    </Card>
  );
};

export default OrganizationForm;
