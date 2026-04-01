"use client";

import { useCallback, useEffect, useMemo, useState } from "react";

import { api } from "@/lib/browser-client";
import type { PortalProposal } from "@/lib/types";

export function useProposalsWorkspace(insurerId?: string, status = "All") {
  const [items, setItems] = useState<PortalProposal[]>([]);
  const [selectedId, setSelectedId] = useState("");
  const [query, setQuery] = useState("");
  const [loading, setLoading] = useState(true);
  const [pendingId, setPendingId] = useState("");
  const [message, setMessage] = useState("");

  const refresh = useCallback(async () => {
    setLoading(true);
    setMessage("");

    try {
      const response = await api.proposals.list(insurerId || undefined, status === "All" ? undefined : status);
      const data = response.data;

      if (!response.ok || !data) {
        setMessage(response.message ?? "Unable to load proposals.");
        return null;
      }

      setItems(data);
      setSelectedId((current) => current || data[0]?.id || "");
      return data;
    } catch {
      setMessage("The proposals service could not be reached.");
      return null;
    } finally {
      setLoading(false);
    }
  }, [insurerId, status]);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  const visibleItems = useMemo(() => {
    const lowered = query.trim().toLowerCase();

    return items.filter((item) => {
      if (!lowered) return true;
      return [item.proposalNumber, item.customerName, item.planName, item.category].some((value) =>
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
