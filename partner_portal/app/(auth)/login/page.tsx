"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [remember, setRemember] = useState(true);
  const [showPass, setShowPass] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);

    // try {

    //   const res = await fetch("/api/auth/login", {
    //     method: "POST",
    //     headers: { "content-type": "application/json" },
    //     body: JSON.stringify({ email, password, remember }),
    //   });

    //   if (!res.ok) {
    //     const msg = await res.text();
    //     throw new Error(msg || "Login failed");
    //   }

    //   const data = (await res.json()) as
    //     | { requires2fa: true; tempToken: string }
    //     | { requires2fa: false };

    //   if ("requires2fa" in data && data.requires2fa) {
    //     sessionStorage.setItem("tempToken", data.tempToken);
    //     router.push("/otp");
    //   } else {
    //     router.push("/");
    //   }
    // } catch (err: any) {
    //   setError(err?.message ?? "Login failed");
    // } finally {
    //   setLoading(false);
    // }
  }

  return (
    <div className="mx-auto w-full">
      <div className="text-sm font-medium text-gray-700 mb-4">Login As</div>

      <form onSubmit={onSubmit} className="space-y-4">
        <input
          className="w-full h-12 rounded-md border border-gray-200 px-4 text-sm outline-none focus:ring-2 focus:ring-[var(--primary-light)]/30"
          placeholder="Email*"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          autoComplete="email"
          required
        />

        <div className="relative">
          <input
            className="w-full h-12 rounded-md border border-gray-200 px-4 pr-12 text-sm outline-none focus:ring-2 focus:ring-[var(--primary-light)]/30"
            placeholder="Password"
            type={showPass ? "text" : "password"}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            autoComplete="current-password"
            required
          />
          <button
            type="button"
            onClick={() => setShowPass((s) => !s)}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
            aria-label="Toggle password visibility"
          >
            {/* simple eye icon */}
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path
                d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7-10-7-10-7Z"
                stroke="currentColor"
                strokeWidth="1.8"
              />
              <path
                d="M12 15a3 3 0 1 0 0-6 3 3 0 0 0 0 6Z"
                stroke="currentColor"
                strokeWidth="1.8"
              />
            </svg>
          </button>
        </div>

        <label className="flex items-center gap-2 text-xs text-gray-600 select-none">
          <input
            type="checkbox"
            checked={remember}
            onChange={(e) => setRemember(e.target.checked)}
          />
          Remember Me
        </label>

        {error && (
          <div className="text-sm text-red-600 bg-red-50 border border-red-100 rounded-md p-3">
            {error}
          </div>
        )}

        <button
          type="submit"
          disabled={loading}
          className="w-full h-12 rounded-md bg-gradient text-white text-sm font-medium  disabled:opacity-60"
        >
          {loading ? "Logging in..." : "Login"}
        </button>
      </form>
    </div>
  );
}
