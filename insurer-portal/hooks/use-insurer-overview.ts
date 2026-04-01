"use client";

import { useCallback, useEffect, useState } from "react";

import { api } from "@/lib/browser-client";
import type { PortalOverview } from "@/lib/types";

export function useInsurerOverview(insurerId?: string) {
  const [overview, setOverview] = useState<PortalOverview | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const refresh = useCallback(async () => {
    setLoading(true);
    setError("");

    try {
      const response = await api.insurer.getOverview(insurerId || undefined);
      const data = response.data;

      if (!response.ok || !data) {
        setError(response.message ?? "Unable to load insurer overview.");
        return null;
      }

      setOverview(data);
      return data;
    } catch {
      setError("The overview service is unavailable right now.");
      return null;
    } finally {
      setLoading(false);
    }
  }, [insurerId]);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  return { overview, loading, error, refresh, setOverview };
}
