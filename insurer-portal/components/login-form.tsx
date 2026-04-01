"use client";

import { LoaderCircle, LogIn, ShieldCheck } from "lucide-react";
import { useRouter } from "next/navigation";
import { useState } from "react";

import { api } from "@/lib/browser-client";
import { setStoredCurrentInsurerId } from "@/lib/current-insurer";
import { loginPageCopy } from "@/lib/tabs/login";

export function LoginForm() {
  const router = useRouter();
  const [mobileNumber, setMobileNumber] = useState("+88017");
  const [password, setPassword] = useState("");
  const [pending, setPending] = useState(false);
  const [message, setMessage] = useState("");

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setPending(true);
    setMessage("");

    try {
      const response = await api.auth.login({ mobileNumber, password });

      if (!response.ok) {
        setMessage(response.message ?? loginPageCopy.messages.invalidCredentials);
        return;
      }

      setStoredCurrentInsurerId("");
      router.replace("/");
      router.refresh();
    } catch {
      setMessage(loginPageCopy.messages.serviceDown);
    } finally {
      setPending(false);
    }
  }

  return (
    <form className="auth-card rounded-[32px] p-6 sm:p-8" onSubmit={handleSubmit}>
      <div className="mb-6 flex items-center gap-3">
        <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-[rgb(15_157_104_/_0.14)] text-[var(--brand-deep)]">
          <ShieldCheck className="h-6 w-6" />
        </div>
        <div>
          <p className="text-sm font-medium uppercase tracking-[0.18em] text-[var(--brand-deep)]">
            {loginPageCopy.header.eyebrow}
          </p>
          <h1 className="font-[family:var(--font-heading)] text-2xl font-semibold text-[var(--text)]">
            {loginPageCopy.header.title}
          </h1>
        </div>
      </div>

      <div className="space-y-4">
        <label className="block space-y-2">
          <span className="text-sm font-medium text-[var(--muted)]">{loginPageCopy.form.mobileLabel}</span>
          <input
            className="portal-input"
            value={mobileNumber}
            onChange={(event) => setMobileNumber(event.target.value)}
            placeholder={loginPageCopy.form.mobilePlaceholder}
            autoComplete="tel"
            required
          />
        </label>

        <label className="block space-y-2">
          <span className="text-sm font-medium text-[var(--muted)]">{loginPageCopy.form.passwordLabel}</span>
          <input
            className="portal-input"
            type="password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            placeholder={loginPageCopy.form.passwordPlaceholder}
            autoComplete="current-password"
            required
          />
        </label>
      </div>

      {message ? (
        <div className="mt-4 rounded-2xl border border-[rgb(194_65_12_/_0.16)] bg-[var(--danger-soft)] px-4 py-3 text-sm text-[var(--danger)]">
          {message}
        </div>
      ) : null}

      <button className="portal-btn portal-btn-primary mt-6 w-full" disabled={pending} type="submit">
        {pending ? <LoaderCircle className="h-4 w-4 animate-spin" /> : <LogIn className="h-4 w-4" />}
        {pending ? loginPageCopy.submit.busy : loginPageCopy.submit.idle}
      </button>

      <p className="mt-4 text-sm text-[var(--muted)]">
        {loginPageCopy.footer}
      </p>
    </form>
  );
}
