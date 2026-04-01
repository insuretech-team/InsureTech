const BD_PHONE_RE = /^880(13|14|15|16|17|18|19)\d{8}$/;

export const BD_MOBILE_EXAMPLES = "01712345678, +8801712345678, or 008801712345678";

function stripPhoneInput(value: string): string {
  const trimmed = value.trim();
  if (!trimmed) return "";
  const compact = trimmed.replace(/[^\d+]/g, "");
  return compact.startsWith("+") ? compact.slice(1) : compact;
}

export function normalizeBangladeshMobile(value: string): string | null {
  const digits = stripPhoneInput(value);
  if (!digits) return null;

  let normalizedDigits: string;
  if (digits.startsWith("00880")) {
    normalizedDigits = digits.slice(2);
  } else if (digits.startsWith("0088")) {
    normalizedDigits = `88${digits.slice(4)}`;
  } else if (digits.startsWith("880")) {
    normalizedDigits = digits;
  } else if (digits.startsWith("88") && digits.length === 13) {
    normalizedDigits = `88${digits.slice(2)}`;
  } else if (digits.startsWith("0")) {
    normalizedDigits = `880${digits.slice(1)}`;
  } else if (digits.length === 10) {
    normalizedDigits = `880${digits}`;
  } else {
    return null;
  }

  if (!BD_PHONE_RE.test(normalizedDigits)) {
    return null;
  }

  return `+${normalizedDigits}`;
}

export function isValidBangladeshMobile(value: string): boolean {
  return normalizeBangladeshMobile(value) !== null;
}

export function normalizeBangladeshMobileOrRaw(value: string): string {
  return normalizeBangladeshMobile(value) ?? value.trim();
}

export function getBangladeshMobileValidationMessage(label = "Mobile number"): string {
  return `${label} must be a valid Bangladesh number. Use ${BD_MOBILE_EXAMPLES}.`;
}
