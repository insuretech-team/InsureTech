import { cn } from "@/lib/utils";

function toneForStatus(status?: string) {
  const normalized = (status ?? "").toLowerCase();

  if (
    normalized.includes("approved") ||
    normalized.includes("active") ||
    normalized.includes("settled") ||
    normalized.includes("live")
  ) {
    return "pill-live";
  }

  if (
    normalized.includes("review") ||
    normalized.includes("pending") ||
    normalized.includes("pilot") ||
    normalized.includes("draft")
  ) {
    return "pill-warn";
  }

  if (
    normalized.includes("reject") ||
    normalized.includes("inactive") ||
    normalized.includes("cancel") ||
    normalized.includes("error")
  ) {
    return "pill-danger";
  }

  return "pill-neutral";
}

export function StatusPill({ status }: { status?: string }) {
  return <span className={cn("pill", toneForStatus(status))}>{status || "Unknown"}</span>;
}
