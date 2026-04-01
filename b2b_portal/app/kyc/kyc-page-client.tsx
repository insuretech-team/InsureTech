"use client";

import { useRouter } from "next/navigation";
import { EKYCFlow } from "@/components/kyc/EKYCFlow";

interface KYCPageClientProps {
  userId: string;
}

function writeKycCookie(value: string) {
  const secure = typeof window !== "undefined" && window.location.protocol === "https:" ? "; Secure" : "";
  document.cookie = `portal_kyc_verified=${encodeURIComponent(value)}; Path=/; Max-Age=${60 * 60 * 12}; SameSite=Lax${secure}`;
}

function readCookie(name: string): string {
  if (typeof document === "undefined") return "";
  const escaped = name.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  const match = document.cookie.match(new RegExp(`(?:^|;\\s*)${escaped}=([^;]*)`));
  return match ? decodeURIComponent(match[1]) : "";
}

export function KYCPageClient({ userId }: KYCPageClientProps) {
  const router = useRouter();

  const handleComplete = (profileImageUrl: string) => {
    console.log("eKYC complete, profile image:", profileImageUrl);
    // Set cookie to pending_review — middleware lifts the /kyc gate.
    writeKycCookie("pending_review");
    const passwordChangeRequired = readCookie("portal_password_change_required") === "true";
    router.replace(passwordChangeRequired ? "/reset-password" : "/?kyc=pending");
  };

  const handleError = (message: string) => {
    console.error("eKYC error:", message);
  };

  return (
    <EKYCFlow
      userId={userId}
      onComplete={handleComplete}
      onError={handleError}
    />
  );
}
