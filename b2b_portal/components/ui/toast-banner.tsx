/**
 * toast-banner.tsx
 * ─────────────────
 * Reusable inline toast/notification banner.
 * Used inside modals and pages to show success/error feedback.
 */
import type { Toast } from "@/src/hooks/useToast";

interface ToastBannerProps {
  toast: Toast | null;
}

export function ToastBanner({ toast }: ToastBannerProps) {
  if (!toast) return null;
  return (
    <div
      className={[
        "mx-6 mt-4 rounded-md px-4 py-2 text-sm font-medium",
        toast.type === "success"
          ? "bg-green-50 text-green-700 border border-green-200"
          : "bg-red-50 text-red-700 border border-red-200",
      ].join(" ")}
    >
      {toast.message}
    </div>
  );
}
