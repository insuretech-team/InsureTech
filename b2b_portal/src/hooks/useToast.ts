/**
 * useToast.ts
 * ───────────
 * Lightweight in-component toast state hook.
 * Replaces the copy-pasted toast state + setTimeout logic in every modal.
 *
 * Usage:
 *   const { toast, showToast, clearToast } = useToast();
 *   showToast("success", "Employee created");
 */
"use client";

import { useCallback, useState } from "react";

export type ToastType = "success" | "error";

export interface Toast {
  type: ToastType;
  message: string;
}

export function useToast(autoDismissMs = 4000) {
  const [toast, setToast] = useState<Toast | null>(null);

  const showToast = useCallback(
    (type: ToastType, message: string) => {
      setToast({ type, message });
      setTimeout(() => setToast(null), autoDismissMs);
    },
    [autoDismissMs]
  );

  const clearToast = useCallback(() => setToast(null), []);

  return { toast, showToast, clearToast };
}
