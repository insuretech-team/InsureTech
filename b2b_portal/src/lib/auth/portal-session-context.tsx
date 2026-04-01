"use client";

import { createContext, useContext, useEffect, useMemo, useState, type ReactNode } from "react";

import type { PortalSession } from "@lib/types/auth";

type PortalSessionContextValue = {
  session: PortalSession | null;
  setSession: (session: PortalSession | null) => void;
};

const PortalSessionContext = createContext<PortalSessionContextValue | null>(null);

export function PortalSessionProvider({
  children,
  initialSession,
}: {
  children: ReactNode;
  initialSession: PortalSession | null;
}) {
  const [session, setSession] = useState<PortalSession | null>(initialSession);

  useEffect(() => {
    setSession(initialSession);
  }, [initialSession]);

  const value = useMemo(
    () => ({
      session,
      setSession,
    }),
    [session]
  );

  return <PortalSessionContext.Provider value={value}>{children}</PortalSessionContext.Provider>;
}

export function usePortalSession() {
  const context = useContext(PortalSessionContext);
  if (!context) {
    return { session: null, setSession: () => {} };
  }
  return context;
}

export function usePortalPrincipal() {
  return usePortalSession().session?.principal ?? null;
}
