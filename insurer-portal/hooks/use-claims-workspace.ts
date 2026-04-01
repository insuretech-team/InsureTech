"use client";

import { useCallback, useEffect, useMemo, useState } from "react";

import { api } from "@/lib/browser-client";
import { isSurveyorRequiredClaim } from "@/lib/claims-intelligence";
import type { PortalClaim } from "@/lib/types";

export function useClaimsWorkspace(insurerId?: string, status = "All", mode: "all" | "surveyor-only" = "all") {
  const [items, setItems] = useState<PortalClaim[]>([]);
  const [selectedId, setSelectedId] = useState("");
  const [query, setQuery] = useState("");
  const [loading, setLoading] = useState(true);
  const [pendingId, setPendingId] = useState("");
  const [message, setMessage] = useState("");

  const refresh = useCallback(
    async (filterFn?: (claim: PortalClaim) => boolean) => {
      setLoading(true);
      setMessage("");

      try {
        const response = await api.claims.list(insurerId || undefined, status === "All" ? undefined : status);
        const data = response.data ?? [];

        if (!response.ok && status !== "All") {
          setMessage(response.message ?? "Unable to load claims.");
          return null;
        }

        const scopedItems = mode === "surveyor-only" ? data.filter((claim) => isSurveyorRequiredClaim(claim)) : data;
        const nextItems = filterFn ? scopedItems.filter(filterFn) : scopedItems;
        setItems(nextItems);
        setSelectedId((current) => current || nextItems[0]?.id || "");
        return nextItems;
      } catch {
        setMessage("The claims service could not be reached.");
        return null;
      } finally {
        setLoading(false);
      }
    },
    [insurerId, mode, status],
  );

  useEffect(() => {
    void refresh();
  }, [refresh, mode]);

  const visibleItems = useMemo(() => {
    const lowered = query.trim().toLowerCase();

    return items.filter((item) => {
      if (!lowered) return true;
      return [item.claimNumber, item.insuredName, item.planName, item.category].some((value) =>
        value.toLowerCase().includes(lowered),
      );
    });
  }, [items, query]);

  const selected = visibleItems.find((item) => item.id === selectedId) ?? visibleItems[0] ?? null;

  return {
    items,
    setItems,
    selected,
    selectedId,
    setSelectedId,
    query,
    setQuery,
    loading,
    pendingId,
    setPendingId,
    message,
    setMessage,
    refresh,
    visibleItems,
  };
}
