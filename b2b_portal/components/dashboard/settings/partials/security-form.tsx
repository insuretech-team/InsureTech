"use client";

import React, { useEffect, useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Field, FieldGroup } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { authClient } from "@lib/sdk";

const focusPurple =
  "focus-visible:ring-primary focus-visible:border-primary focus-visible:ring-2";

// ─── Change Password ──────────────────────────────────────────────────────────

const ChangePasswordSection = () => {
  const [form, setForm] = useState({ old_password: "", new_password: "", confirm: "" });
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState<{ text: string; ok: boolean } | null>(null);

  const set = (key: keyof typeof form) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setForm((f) => ({ ...f, [key]: e.target.value }));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setMessage(null);
    if (form.new_password !== form.confirm) {
      setMessage({ text: "New passwords do not match.", ok: false });
      return;
    }
    if (form.new_password.length < 8) {
      setMessage({ text: "New password must be at least 8 characters.", ok: false });
      return;
    }
    setSaving(true);
    try {
      const res = await authClient.changePassword({
        old_password: form.old_password,
        new_password: form.new_password,
      });
      setMessage({ text: res.message ?? (res.ok ? "Password changed successfully." : "Failed to change password."), ok: res.ok });
      if (res.ok) setForm({ old_password: "", new_password: "", confirm: "" });
    } catch {
      setMessage({ text: "An error occurred. Please try again.", ok: false });
    } finally {
      setSaving(false);
    }
  };

  return (
    <div>
      <h3 className="text-base font-semibold text-foreground mb-4">Change Password</h3>
      <form onSubmit={handleSubmit}>
        <FieldGroup className="space-y-4 gap-0">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Field>
              <Label htmlFor="old_password">Current Password</Label>
              <Input
                id="old_password"
                type="password"
                placeholder="Current password"
                className={focusPurple}
                value={form.old_password}
                onChange={set("old_password")}
                required
              />
            </Field>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Field>
              <Label htmlFor="new_password">New Password</Label>
              <Input
                id="new_password"
                type="password"
                placeholder="New password (min 8 chars)"
                className={focusPurple}
                value={form.new_password}
                onChange={set("new_password")}
                required
              />
            </Field>
            <Field>
              <Label htmlFor="confirm_password">Confirm New Password</Label>
              <Input
                id="confirm_password"
                type="password"
                placeholder="Confirm new password"
                className={focusPurple}
                value={form.confirm}
                onChange={set("confirm")}
                required
              />
            </Field>
          </div>
          {message && (
            <p className={`text-sm ${message.ok ? "text-green-600" : "text-red-500"}`}>
              {message.text}
            </p>
          )}
          <div className="flex items-center justify-end mt-2">
            <Button
              type="submit"
              variant="default"
              className="bg-primary hover:bg-accent"
              disabled={saving}
            >
              {saving ? "Saving…" : "Change Password"}
            </Button>
          </div>
        </FieldGroup>
      </form>
    </div>
  );
};

// ─── TOTP / 2FA ───────────────────────────────────────────────────────────────

const TotpSection = () => {
  const [totpData, setTotpData] = useState<Record<string, unknown> | null>(null);
  const [enabling, setEnabling] = useState(false);
  const [disabling, setDisabling] = useState(false);
  const [totpCode, setTotpCode] = useState("");
  const [message, setMessage] = useState<{ text: string; ok: boolean } | null>(null);

  const handleEnable = async () => {
    setEnabling(true);
    setMessage(null);
    try {
      const res = await authClient.enableTotp();
      if (res.ok && res.totp) {
        setTotpData(res.totp);
        setMessage({ text: "Scan the QR code with your authenticator app.", ok: true });
      } else {
        setMessage({ text: res.message ?? "Failed to enable 2FA.", ok: false });
      }
    } catch {
      setMessage({ text: "An error occurred.", ok: false });
    } finally {
      setEnabling(false);
    }
  };

  const handleDisable = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!totpCode.trim()) return;
    setDisabling(true);
    setMessage(null);
    try {
      const res = await authClient.disableTotp(totpCode);
      setMessage({ text: res.message ?? (res.ok ? "2FA disabled." : "Failed to disable 2FA."), ok: res.ok });
      if (res.ok) { setTotpData(null); setTotpCode(""); }
    } catch {
      setMessage({ text: "An error occurred.", ok: false });
    } finally {
      setDisabling(false);
    }
  };

  return (
    <div>
      <h3 className="text-base font-semibold text-foreground mb-1">Two-Factor Authentication (TOTP)</h3>
      <p className="text-sm text-muted-foreground mb-4">
        Add an extra layer of security using an authenticator app (Google Authenticator, Authy, etc).
      </p>

      {!totpData ? (
        <div className="flex flex-col gap-3">
          <Button
            variant="default"
            className="bg-primary hover:bg-accent w-fit"
            onClick={handleEnable}
            disabled={enabling}
          >
            {enabling ? "Setting up…" : "Enable 2FA"}
          </Button>

          {/* Disable 2FA form — shown if user already has TOTP enabled */}
          <form onSubmit={handleDisable} className="flex items-center gap-3 mt-2">
            <Input
              placeholder="Enter current TOTP code to disable"
              className={`${focusPurple} max-w-xs`}
              value={totpCode}
              onChange={(e) => setTotpCode(e.target.value)}
            />
            <Button
              type="submit"
              variant="default"
              className="bg-destructive hover:bg-destructive/80 text-white"
              disabled={disabling || !totpCode.trim()}
            >
              {disabling ? "Disabling…" : "Disable 2FA"}
            </Button>
          </form>
        </div>
      ) : (
        <div className="space-y-3">
          {typeof totpData.qr_code_url === "string" && totpData.qr_code_url && (
            // eslint-disable-next-line @next/next/no-img-element
            <img src={totpData.qr_code_url} alt="TOTP QR Code" className="w-40 h-40 border rounded" />
          )}
          {typeof totpData.secret === "string" && totpData.secret && (
            <p className="text-sm font-mono bg-muted px-3 py-2 rounded">
              Manual key: <strong>{totpData.secret}</strong>
            </p>
          )}
        </div>
      )}

      {message && (
        <p className={`text-sm mt-2 ${message.ok ? "text-green-600" : "text-red-500"}`}>
          {message.text}
        </p>
      )}
    </div>
  );
};

// ─── Active Sessions ──────────────────────────────────────────────────────────

type SessionRecord = {
  session_id?: string;
  device_name?: string;
  device_type?: string;
  ip_address?: string;
  user_agent?: string;
  created_at?: string;
  last_activity_at?: string; // SDK field name — was incorrectly last_active_at
  is_active?: boolean;
};

const SessionsSection = () => {
  const [sessions, setSessions] = useState<SessionRecord[]>([]);
  const [loading, setLoading] = useState(true);
  const [revoking, setRevoking] = useState<string | null>(null);
  const [revokingAll, setRevokingAll] = useState(false);
  const [message, setMessage] = useState<{ text: string; ok: boolean } | null>(null);

  const fetchSessions = (delay = 0) => {
    setLoading(true);
    const run = () => authClient.listSessions().then((res) => {
      if (res.ok && res.sessions) {
        // BFF returns { ok, sessions: SessionsListingResponse }
        // SessionsListingResponse = { sessions?: Array<Session>, total_count?, ... }
        // BFF passes active_only=true so revoked sessions are excluded at DB level.
        // Client-side guard: also filter out any is_active=false that slips through.
        const raw = res.sessions as Record<string, unknown>;
        const list = (Array.isArray(raw.sessions) ? raw.sessions : []) as SessionRecord[];
        setSessions(list.filter((s) => s.is_active !== false));
      }
      setLoading(false);
    }).catch(() => setLoading(false));
    if (delay > 0) setTimeout(run, delay); else run();
  };

  useEffect(() => { fetchSessions(); }, []);

  const handleRevoke = async (sessionId: string) => {
    setRevoking(sessionId);
    setMessage(null);
    try {
      const res = await authClient.revokeSession(sessionId);
      setMessage({ text: res.message ?? (res.ok ? "Session revoked." : "Failed."), ok: res.ok });
      // Small delay so the DB write propagates before we re-fetch.
      if (res.ok) fetchSessions(500);
    } catch {
      setMessage({ text: "An error occurred.", ok: false });
    } finally {
      setRevoking(null);
    }
  };

  const handleRevokeAll = async () => {
    setRevokingAll(true);
    setMessage(null);
    try {
      const res = await authClient.revokeAllSessions();
      setMessage({ text: res.message ?? (res.ok ? "All sessions revoked." : "Failed."), ok: res.ok });
      if (res.ok) {
        setSessions([]); // optimistic clear
        fetchSessions(800); // re-fetch after write propagates
      }
    } catch {
      setMessage({ text: "An error occurred.", ok: false });
    } finally {
      setRevokingAll(false);
    }
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-base font-semibold text-foreground">Active Sessions</h3>
        {sessions.length > 0 && (
          <Button
            variant="default"
            className="bg-destructive hover:bg-destructive/80 text-white text-xs"
            onClick={handleRevokeAll}
            disabled={revokingAll}
          >
            {revokingAll ? "Signing out…" : "Sign Out Everywhere"}
          </Button>
        )}
      </div>

      {loading ? (
        <p className="text-sm text-muted-foreground">Loading sessions…</p>
      ) : sessions.length === 0 ? (
        <p className="text-sm text-muted-foreground">No active sessions found.</p>
      ) : (
        <div className="space-y-2">
          {sessions.map((s) => (
            <div
              key={s.session_id}
              className="flex items-center justify-between bg-muted px-4 py-3 rounded text-sm"
            >
              <div>
                <p className="font-medium text-foreground">
                  {s.device_name ?? s.device_type ?? "Unknown device"}
                </p>
                <p className="text-muted-foreground text-xs">
                  {s.ip_address ?? "Unknown IP"} · Last active:{" "}
                  {s.last_activity_at ? new Date(s.last_activity_at).toLocaleString() : "—"}
                </p>
              </div>
              <Button
                variant="default"
                className="bg-destructive hover:bg-destructive/80 text-white text-xs h-7 px-3"
                onClick={() => s.session_id && handleRevoke(s.session_id)}
                disabled={revoking === s.session_id}
              >
                {revoking === s.session_id ? "…" : "Revoke"}
              </Button>
            </div>
          ))}
        </div>
      )}

      {message && (
        <p className={`text-sm mt-3 ${message.ok ? "text-green-600" : "text-red-500"}`}>
          {message.text}
        </p>
      )}
    </div>
  );
};

// ─── Security Form (combined) ─────────────────────────────────────────────────

const SecurityForm = () => {
  return (
    <Card>
      <CardContent className="py-6 space-y-8">
        <ChangePasswordSection />
        <Separator />
        <TotpSection />
        <Separator />
        <SessionsSection />
      </CardContent>
    </Card>
  );
};

export default SecurityForm;
