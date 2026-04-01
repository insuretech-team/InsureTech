"use client";

import { useMemo } from "react";

import { useInsurerOverview } from "@/hooks/use-insurer-overview";
import { getInsurerPlaybooks } from "@/lib/product-playbooks";
import { getDashboardWorkspace } from "@/lib/tabs/dashboard";

export function useDashboardWorkspace(insurerId?: string) {
  const { overview, loading, error, refresh } = useInsurerOverview(insurerId);

  const playbooks = useMemo(
    () => getInsurerPlaybooks(overview?.currentInsurer?.name),
    [overview?.currentInsurer?.name],
  );

  const workspace = useMemo(
    () => getDashboardWorkspace(overview, playbooks),
    [overview, playbooks],
  );

  return {
    overview,
    loading,
    error,
    refresh,
    workspace,
  };
}
