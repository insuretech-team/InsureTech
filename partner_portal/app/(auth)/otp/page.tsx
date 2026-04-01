"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import { useRouter } from "next/navigation";

function pad2(n: number) {
  return String(n).padStart(2, "0");
}

export default function OtpPage() {
  const router = useRouter();

  const [digits, setDigits] = useState<string[]>(Array(6).fill(""));
  const inputsRef = useRef<Array<HTMLInputElement | null>>([]);
  const [secondsLeft, setSecondsLeft] = useState(25);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const otp = useMemo(() => digits.join(""), [digits]);

  // useEffect(() => {
  //   const tempToken = sessionStorage.getItem("tempToken");
  //   if (!tempToken) router.replace("/login");
  // }, [router]);

  useEffect(() => {
    if (secondsLeft <= 0) return;
    const t = setInterval(() => setSecondsLeft((s) => s - 1), 1000);
    return () => clearInterval(t);
  }, [secondsLeft]);

  function setAt(index: number, value: string) {
    setDigits((prev) => {
      const next = [...prev];
      next[index] = value;
      return next;
    });
  }

  function onChange(index: number, raw: string) {
    const v = raw.replace(/\D/g, "");
    if (!v) {
      setAt(index, "");
      return;
    }

    // handle paste / multi input
    const chars = v.slice(0, 6 - index).split("");
    chars.forEach((c, i) => setAt(index + i, c));

    const nextIndex = Math.min(index + chars.length, 5);
    inputsRef.current[nextIndex]?.focus();
  }

  function onKeyDown(index: number, e: React.KeyboardEvent<HTMLInputElement>) {
    if (e.key === "Backspace" && !digits[index] && index > 0) {
      inputsRef.current[index - 1]?.focus();
    }
    if (e.key === "ArrowLeft" && index > 0)
      inputsRef.current[index - 1]?.focus();
    if (e.key === "ArrowRight" && index < 5)
      inputsRef.current[index + 1]?.focus();
  }

  async function onVerify() {
    setError(null);
    if (otp.length !== 6) {
      setError("Please enter the 6-digit OTP.");
      return;
    }

    setLoading(true);
    try {
      const tempToken = sessionStorage.getItem("tempToken");

      // TODO: Replace with your real backend endpoint
      // const res = await fetch("/api/auth/verify-otp", {
      //   method: "POST",
      //   headers: { "content-type": "application/json" },
      //   body: JSON.stringify({ otp, tempToken }),
      // });

      // if (!res.ok) {
      //   const msg = await res.text();
      //   throw new Error(msg || "OTP verification failed");
      // }

      // backend should set httpOnly cookie OR return access token
      sessionStorage.removeItem("tempToken");
      router.push("/");
    } catch (err: any) {
      setError(err?.message ?? "OTP verification failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="mx-auto w-full">
      <div className="text-left text-sm font-medium text-gray-700 mb-3">
        OTP
      </div>

      <div className="flex items-center justify-between gap-3">
        <div className="flex gap-3">
          {digits.map((d, i) => (
            <input
              key={i}
              ref={(el) => {
                inputsRef.current[i] = el;
              }}
              value={d}
              onChange={(e) => onChange(i, e.target.value)}
              onKeyDown={(e) => onKeyDown(i, e)}
              inputMode="numeric"
              maxLength={6}
              className="w-14 h-14 rounded-md border border-gray-200 text-center text-lg outline-none focus:ring-2 focus:ring-green-700/30"
              aria-label={`OTP digit ${i + 1}`}
            />
          ))}
        </div>

        <div className="text-xs text-gray-500 whitespace-nowrap">
          00m : {pad2(secondsLeft)}s
        </div>
      </div>

      {error && (
        <div className="mt-4 text-sm text-red-600 bg-red-50 border border-red-100 rounded-md p-3">
          {error}
        </div>
      )}

      <button
        type="button"
        onClick={onVerify}
        disabled={loading}
        className="mt-6 w-full h-12 rounded-md bg-[var(--primary)] text-white text-sm font-medium hover:bg-[var(--primary-light)] disabled:opacity-60"
      >
        {loading ? "Verifying..." : "Verify"}
      </button>
    </div>
  );
}
