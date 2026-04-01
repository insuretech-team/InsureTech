"use client";

import { useEffect, useState } from "react";

import { getStoredCurrentInsurerId } from "@/lib/current-insurer";

export function useCurrentInsurerId() {
  const [insurerId, setInsurerId] = useState(() => getStoredCurrentInsurerId());

  useEffect(() => {
    function handleInsurerChange() {
      setInsurerId(getStoredCurrentInsurerId());
    }

    window.addEventListener("insurer:changed", handleInsurerChange as EventListener);
    return () => {
      window.removeEventListener("insurer:changed", handleInsurerChange as EventListener);
    };
  }, []);

  return { insurerId, setInsurerId };
}
