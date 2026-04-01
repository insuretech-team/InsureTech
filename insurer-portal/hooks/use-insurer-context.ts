"use client";

import { useCallback, useEffect, useState } from "react";

import { api } from "@/lib/browser-client";
import type { InsurerConfigForm, PortalInsurer, PortalProduct, PortalSessionSnapshot } from "@/lib/types";

export type PortalContextResponse = {
  session?: PortalSessionSnapshot;
  insurers: PortalInsurer[];
  currentInsurer: PortalInsurer | null;
  config: InsurerConfigForm | null;
  products: PortalProduct[];
  source: "live" | "fallback" | "mixed";
};

const initialContext: PortalContextResponse = {
  insurers: [],
  currentInsurer: null,
  config: null,
  products: [],
  source: "fallback",
};

export function useInsurerContext(insurerId?: string) {
  const [context, setContext] = useState<PortalContextResponse>(initialContext);
  const [loading, setLoading] = useState(true);

  const refresh = useCallback(async () => {
    setLoading(true);

    try {
      const response = await api.insurer.getContext(insurerId || undefined);
      const data = response.data as PortalContextResponse | undefined;

      if (response.ok && data?.currentInsurer !== undefined) {
        setContext(data);
        return data;
      }

      return null;
    } finally {
      setLoading(false);
    }
  }, [insurerId]);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  return { context, loading, refresh, setContext };
}
