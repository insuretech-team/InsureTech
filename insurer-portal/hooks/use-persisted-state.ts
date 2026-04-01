"use client";

import { useEffect, useState } from "react";

export function usePersistedState<T>(storageKey: string, initialValue: T) {
  const [value, setValue] = useState<T>(initialValue);
  const [ready, setReady] = useState(false);

  useEffect(() => {
    try {
      const stored = window.localStorage.getItem(storageKey);
      if (stored) {
        setValue(JSON.parse(stored) as T);
      }
    } finally {
      setReady(true);
    }
  }, [storageKey]);

  useEffect(() => {
    if (!ready) return;
    window.localStorage.setItem(storageKey, JSON.stringify(value));
  }, [ready, storageKey, value]);

  return { value, setValue, ready };
}
