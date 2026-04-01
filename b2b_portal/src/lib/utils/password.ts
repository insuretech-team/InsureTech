const PASSWORD_SYMBOL_RE = /[^A-Za-z0-9]/;

function joinRequirements(parts: string[]): string {
  if (parts.length <= 1) {
    return parts[0] ?? "";
  }
  if (parts.length === 2) {
    return `${parts[0]} and ${parts[1]}`;
  }
  return `${parts.slice(0, -1).join(", ")}, and ${parts.at(-1)}`;
}

export const PASSWORD_REQUIREMENTS_HINT =
  "Use at least 8 characters with uppercase, lowercase, number, and symbol.";

export function getPasswordValidationMessage(value: string, label = "Password"): string | null {
  if (!value.trim()) {
    return `${label} is required.`;
  }

  const missing: string[] = [];
  if (value.length < 8) missing.push("at least 8 characters");
  if (!/[A-Z]/.test(value)) missing.push("one uppercase letter");
  if (!/[a-z]/.test(value)) missing.push("one lowercase letter");
  if (!/\d/.test(value)) missing.push("one number");
  if (!PASSWORD_SYMBOL_RE.test(value)) missing.push("one symbol");

  if (missing.length === 0) {
    return null;
  }

  return `${label} must include ${joinRequirements(missing)}.`;
}
