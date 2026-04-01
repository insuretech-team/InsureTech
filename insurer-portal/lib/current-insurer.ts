const CURRENT_INSURER_KEY = "insurer-portal.current-insurer-id";

export function getStoredCurrentInsurerId() {
  if (typeof window === "undefined") return "";
  return window.localStorage.getItem(CURRENT_INSURER_KEY) ?? "";
}

export function setStoredCurrentInsurerId(value: string) {
  if (typeof window === "undefined") return;
  if (value) {
    window.localStorage.setItem(CURRENT_INSURER_KEY, value);
  } else {
    window.localStorage.removeItem(CURRENT_INSURER_KEY);
  }
  window.dispatchEvent(new CustomEvent("insurer:changed", { detail: value }));
}
