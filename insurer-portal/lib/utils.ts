export function cn(...parts: Array<string | false | null | undefined>) {
  return parts.filter(Boolean).join(" ");
}

export function formatMoney(value: number, currency = "BDT") {
  if (!Number.isFinite(value)) return "BDT 0";
  return new Intl.NumberFormat("en-BD", {
    style: "currency",
    currency,
    maximumFractionDigits: 0,
  }).format(value);
}

export function formatDate(value?: string) {
  if (!value) return "Unavailable";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return new Intl.DateTimeFormat("en-BD", {
    dateStyle: "medium",
  }).format(date);
}

export function formatDateTime(value?: string) {
  if (!value) return "Unavailable";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return new Intl.DateTimeFormat("en-BD", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(date);
}

export function titleFromEnum(value?: string) {
  if (!value) return "Unknown";
  return value
    .replace(/^[A-Z]+_STATUS_/g, "")
    .replace(/^[A-Z]+_/g, "")
    .toLowerCase()
    .split("_")
    .filter(Boolean)
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(" ");
}

export function inferCategory(...candidates: Array<string | undefined>) {
  const joined = candidates.join(" ").toLowerCase();
  if (joined.includes("travel")) return "Travel";
  if (joined.includes("motor") || joined.includes("auto") || joined.includes("car")) return "Auto";
  if (joined.includes("fire") || joined.includes("property")) return "Fire";
  if (joined.includes("health") || joined.includes("hospital")) return "Health";
  if (joined.includes("life")) return "Life";
  if (joined.includes("device") || joined.includes("gadget")) return "Device";
  return "General";
}

export function initialLetters(value?: string) {
  if (!value) return "IP";
  const parts = value.trim().split(/\s+/).filter(Boolean);
  if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase();
  return `${parts[0][0] ?? ""}${parts[1][0] ?? ""}`.toUpperCase();
}
