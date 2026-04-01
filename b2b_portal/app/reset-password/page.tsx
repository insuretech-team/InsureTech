"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { bffClient } from "@lib/sdk/b2b-sdk-client";
import {
  getPasswordValidationMessage,
  PASSWORD_REQUIREMENTS_HINT,
} from "@/src/lib/utils/password";

export default function ResetPasswordPage() {
  const router = useRouter();
  const [temporaryPassword, setTemporaryPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const newPasswordError = getPasswordValidationMessage(newPassword, "New password");

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setSuccess("");

    if (!temporaryPassword.trim()) {
      setError("Temporary password is required.");
      return;
    }
    if (newPasswordError) {
      setError(newPasswordError);
      return;
    }
    if (newPassword !== confirmPassword) {
      setError("Confirm password must match the new password.");
      return;
    }

    setSubmitting(true);
    try {
      const result = await bffClient.auth.changePassword({
        old_password: temporaryPassword,
        new_password: newPassword,
      });
      if (!result.ok) {
        setError(result.message ?? "Failed to update password.");
        return;
      }

      setSuccess("Password updated. Please sign in with your new password.");
      setTimeout(() => {
        router.replace("/login?passwordReset=1");
        router.refresh();
      }, 1000);
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <main className="min-h-screen bg-muted/20 px-4 py-10">
      <div className="mx-auto max-w-md rounded-2xl border bg-background p-6 shadow-sm">
        <div className="space-y-2">
          <h1 className="text-2xl font-semibold text-foreground">Set your new password</h1>
          <p className="text-sm text-muted-foreground">
            Your temporary password worked, but you need to replace it before continuing.
          </p>
        </div>

        <form onSubmit={handleSubmit} className="mt-6 space-y-4">
          <div className="space-y-2">
            <label htmlFor="temporary-password" className="text-sm font-medium text-foreground">
              Temporary password
            </label>
            <Input
              id="temporary-password"
              type="password"
              value={temporaryPassword}
              onChange={(event) => setTemporaryPassword(event.target.value)}
              placeholder="Enter temporary password"
              required
            />
          </div>

          <div className="space-y-2">
            <label htmlFor="new-password" className="text-sm font-medium text-foreground">
              New password
            </label>
            <Input
              id="new-password"
              type="password"
              value={newPassword}
              onChange={(event) => setNewPassword(event.target.value)}
              placeholder="Create new password"
              required
            />
            <p className={`text-xs ${newPasswordError ? "text-destructive" : "text-muted-foreground"}`}>
              {newPasswordError ?? PASSWORD_REQUIREMENTS_HINT}
            </p>
          </div>

          <div className="space-y-2">
            <label htmlFor="confirm-password" className="text-sm font-medium text-foreground">
              Confirm new password
            </label>
            <Input
              id="confirm-password"
              type="password"
              value={confirmPassword}
              onChange={(event) => setConfirmPassword(event.target.value)}
              placeholder="Repeat new password"
              required
            />
          </div>

          {error ? (
            <div className="rounded-md bg-destructive/10 px-3 py-2 text-sm text-destructive">{error}</div>
          ) : null}
          {success ? (
            <div className="rounded-md bg-emerald-50 px-3 py-2 text-sm text-emerald-700">{success}</div>
          ) : null}

          <Button type="submit" className="w-full" disabled={submitting}>
            {submitting ? "Updating password..." : "Update password"}
          </Button>
        </form>
      </div>
    </main>
  );
}
