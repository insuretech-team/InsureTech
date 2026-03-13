/**
 * useCrudList.ts
 * ──────────────
 * Generic hook that manages list + loading state for any CRUD entity.
 *
 * CRITICAL: Stores `fetcher` in a ref so it NEVER appears in the useEffect
 * dependency array — preventing the infinite-loop where an inline arrow
 * `() => client.list()` creates a new reference on every render.
 *
 * Usage (pass a stable reference or an inline arrow — both are safe):
 *   const { data, loading, reload } = useCrudList(employeeClient.list, "employees");
 *   // or:
 *   const { data, loading, reload } = useCrudList(
 *     () => employeeClient.list({ pageSize: 50 }), "employees"
 *   );
 */
"use client";

import { useCallback, useEffect, useRef, useState } from "react";

export function useCrudList<T>(
  fetcher: () => Promise<{ ok: boolean; [key: string]: unknown }>,
  dataKey: string
) {
  const [data, setData]       = useState<T[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError]     = useState<string | null>(null);

  // Keep fetcher + dataKey in refs so reload() never becomes stale
  // without needing to appear in dependency arrays.
  const fetcherRef = useRef(fetcher);
  const keyRef     = useRef(dataKey);
  fetcherRef.current = fetcher;
  keyRef.current     = dataKey;

  const reload = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const result = await fetcherRef.current();
      const rows   = Array.isArray(result[keyRef.current])
        ? (result[keyRef.current] as T[])
        : [];
      setData(rows);
      if (!result.ok) {
        setError((result.message as string | undefined) ?? "Failed to load data");
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unexpected error");
      setData([]);
    } finally {
      setLoading(false);
    }
  }, []); // stable — no deps needed thanks to refs

  // Fetch exactly once on mount
  // eslint-disable-next-line react-hooks/exhaustive-deps
  useEffect(() => { reload(); }, []);

  return { data, loading, error, reload };
}
