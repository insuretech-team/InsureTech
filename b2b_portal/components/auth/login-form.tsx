"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import { CheckCircle2, XCircle } from "lucide-react";

import { authClient } from "@lib/sdk";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

// Valid Bangladesh operator prefixes after country code 880: 13,14,15,16,17,18,19
const BD_PHONE_RE = /^880(13|14|15|16|17|18|19)\d{8}$/;

/**
 * Client-side mobile number normalizer — mirrors the server-side logic in
 * app/api/auth/login/route.ts so the UI can validate before hitting the API.
 * Returns the canonical +880XXXXXXXXXX string, or null if unrecognisable.
 */
function normalizeMobile(value: string): string | null {
  const stripped = value.trim().replace(/[^\d+]/g, "");
  const digits = stripped.startsWith("+") ? stripped.slice(1) : stripped;

  let e164: string;
  if (digits.startsWith("00880")) {
    e164 = digits.slice(2);
  } else if (digits.startsWith("880")) {
    e164 = digits;
  } else if (digits.startsWith("0088")) {
    e164 = "880" + digits.slice(4);
  } else if (digits.startsWith("88") && digits.length === 13) {
    e164 = "880" + digits.slice(2);
  } else if (digits.startsWith("0")) {
    e164 = "880" + digits.slice(1);
  } else if (digits.length === 10) {
    e164 = "880" + digits;
  } else {
    return null;
  }
  return BD_PHONE_RE.test(e164) ? `+${e164}` : null;
}

function getMobileHint(value: string): string | null {
  if (!value.trim()) return null;
  const normalized = normalizeMobile(value);
  if (normalized) return null; // valid — no hint needed
  // Give a hint only after the user has typed enough to be wrong
  const digits = value.replace(/\D/g, "");
  if (digits.length < 7) return null;
  return "Enter a valid Bangladesh number: 01712345678, +8801712345678, or 008801712345678";
}

type DialogState =
  | { open: false }
  | { open: true; kind: "success" }
  | { open: true; kind: "error"; message: string };

export default function LoginForm() {
  const router = useRouter();
  const [mobileNumber, setMobileNumber] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [mobileTouched, setMobileTouched] = useState(false);
  const [dialog, setDialog] = useState<DialogState>({ open: false });

  const mobileHint = mobileTouched ? getMobileHint(mobileNumber) : null;

  function closeDialog() {
    setDialog({ open: false });
  }

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setMobileTouched(true);

    const normalized = normalizeMobile(mobileNumber);
    if (!normalized) {
      setDialog({
        open: true,
        kind: "error",
        message:
          "Invalid mobile number. Use formats like 01712345678, +8801712345678 or 008801712345678.",
      });
      return;
    }

    setLoading(true);
    try {
      const response = await authClient.login({ mobileNumber: normalized, password });
      if (!response.ok) {
        setDialog({
          open: true,
          kind: "error",
          message: response.message ?? "Login failed. Please try again.",
        });
        return;
      }
      // Show brief success modal then redirect
      setDialog({ open: true, kind: "success" });
      setTimeout(() => {
        router.replace("/");
        router.refresh();
      }, 1200);
    } catch (submitError) {
      setDialog({
        open: true,
        kind: "error",
        message:
          submitError instanceof Error
            ? submitError.message
            : "An unexpected error occurred. Please try again.",
      });
    } finally {
      setLoading(false);
    }
  }

  const isError = dialog.open && dialog.kind === "error";
  const isSuccess = dialog.open && dialog.kind === "success";

  return (
    <>
      {/* ── Result Modal ────────────────────────────────────────────── */}
      <Dialog open={dialog.open} onOpenChange={(open) => !open && closeDialog()}>
        <DialogContent showCloseButton={isError} className="sm:max-w-sm">
          {isSuccess && (
            <>
              <DialogHeader className="items-center gap-3">
                <CheckCircle2 className="size-12 text-emerald-500" />
                <DialogTitle className="text-center text-lg">Signed in successfully</DialogTitle>
                <DialogDescription className="text-center">
                  Redirecting you to the dashboard…
                </DialogDescription>
              </DialogHeader>
            </>
          )}

          {isError && (
            <>
              <DialogHeader className="items-center gap-3">
                <XCircle className="size-12 text-destructive" />
                <DialogTitle className="text-center text-lg">Sign-in failed</DialogTitle>
                <DialogDescription className="text-center">
                  {dialog.message}
                </DialogDescription>
              </DialogHeader>
              <DialogFooter className="sm:justify-center">
                <Button variant="outline" onClick={closeDialog} className="w-full sm:w-auto">
                  Try again
                </Button>
              </DialogFooter>
            </>
          )}
        </DialogContent>
      </Dialog>

      {/* ── Login Form ──────────────────────────────────────────────── */}
      <form onSubmit={onSubmit} className="space-y-5">
        <div className="space-y-2">
          <label className="text-sm font-medium text-foreground" htmlFor="mobileNumber">
            Mobile Number
          </label>
          <input
            id="mobileNumber"
            type="tel"
            className={`auth-input${mobileHint ? " border-amber-500 focus:ring-amber-500" : ""}`}
            value={mobileNumber}
            onChange={(event) => setMobileNumber(event.target.value)}
            onBlur={() => setMobileTouched(true)}
            placeholder="01712345678 or +8801712345678"
            autoComplete="tel"
            required
          />
          {mobileHint ? (
            <p className="text-xs text-amber-600">{mobileHint}</p>
          ) : null}
        </div>

        <div className="space-y-2">
          <label className="text-sm font-medium text-foreground" htmlFor="password">
            Password
          </label>
          <input
            id="password"
            type="password"
            className="auth-input"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            placeholder="Enter password"
            required
          />
        </div>

        <button
          type="submit"
          disabled={loading}
          className="auth-submit"
        >
          {loading ? "Signing in..." : "Sign in"}
        </button>
      </form>
    </>
  );
}
