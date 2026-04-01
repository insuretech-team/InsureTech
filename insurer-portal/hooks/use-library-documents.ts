"use client";

import { useCallback, useEffect, useState } from "react";

import { api } from "@/lib/browser-client";
import type { LibraryDocument, LibraryPack } from "@/lib/types";

export interface UseLibraryDocumentsResult {
  documents: LibraryDocument[];
  packs: LibraryPack[];
  loading: boolean;
  error: string | null;
  refresh: () => void;
  categoryOptions: string[];
  stageOptions: string[];
}

export function useLibraryDocuments(): UseLibraryDocumentsResult {
  const [documents, setDocuments] = useState<LibraryDocument[]>([]);
  const [packs, setPacks] = useState<LibraryPack[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchLibrary = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await api.documents.library();
      if (res.ok && res.data) {
        setDocuments(res.data.documents ?? []);
        setPacks(res.data.packs ?? []);
      } else {
        setError(res.message ?? "Failed to load document library.");
      }
    } catch {
      setError("Unable to load the document library. Please refresh.");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void fetchLibrary();
  }, [fetchLibrary]);

  const categoryOptions = [
    "All",
    ...Array.from(new Set(documents.map((d) => d.category))).sort(),
  ];

  const stageOptions = [
    "All",
    ...Array.from(new Set(documents.map((d) => d.stage))).sort(),
  ];

  return {
    documents,
    packs,
    loading,
    error,
    refresh: () => void fetchLibrary(),
    categoryOptions,
    stageOptions,
  };
}
