"use client";

import { useEffect, useMemo, useState } from "react";

import { usePersistedState } from "@/hooks/use-persisted-state";
import { buildPreviewBundle, getDocumentSourceUrl } from "@/lib/document-preview-layouts";
import {
  buildDigitalBlocks,
  insurerManagedDocuments,
  isDigitalTemplate,
  type ReferenceDocument,
} from "@/lib/pragati-documents";
import {
  createDocumentEditorState,
  createInitialDocumentDraft,
  documentCustomStorageKey,
  documentDraftStorageKey,
  documentOverrideStorageKey,
  documentStageFilterOptions,
  filterManagedDocuments,
  getDocumentCategoryFilterOptions,
  getManagedDocuments,
  toCustomDocument,
  type DocumentEditorMode,
  type DocumentEditorState,
  type DocumentOverrides,
  type DocumentViewMode,
  type SheetDraft,
} from "@/lib/tabs/documents";

export function useDocumentLibraryWorkspace() {
  const [query, setQuery] = useState("");
  const [stageFilter, setStageFilter] = useState<(typeof documentStageFilterOptions)[number]>("All");
  const [categoryFilter, setCategoryFilter] = useState("All");
  const [activeId, setActiveId] = useState("");
  const [activePage, setActivePage] = useState(0);
  const [viewMode, setViewMode] = useState<DocumentViewMode>("preview");
  const [editorMode, setEditorMode] = useState<DocumentEditorMode>("none");
  const [editorState, setEditorState] = useState<DocumentEditorState>(createDocumentEditorState());
  const { value: drafts, setValue: setDrafts, ready: draftsReady } = usePersistedState<Record<string, SheetDraft>>(
    documentDraftStorageKey,
    {},
  );
  const { value: overrides, setValue: setOverrides, ready: overridesReady } = usePersistedState<DocumentOverrides>(
    documentOverrideStorageKey,
    {},
  );
  const { value: customDocuments, setValue: setCustomDocuments, ready: customReady } = usePersistedState<
    ReferenceDocument[]
  >(documentCustomStorageKey, []);

  const stageFilters = documentStageFilterOptions;
  const categoryOptions = useMemo(() => getDocumentCategoryFilterOptions(), []);
  const isStorageReady = draftsReady && overridesReady && customReady;

  const allDocuments = useMemo(
    () => getManagedDocuments([...insurerManagedDocuments, ...customDocuments], overrides),
    [customDocuments, overrides],
  );

  const filteredDocuments = useMemo(
    () => filterManagedDocuments(allDocuments, query, stageFilter, categoryFilter),
    [allDocuments, categoryFilter, query, stageFilter],
  );

  const activeDocument = allDocuments.find((document) => document.id === activeId) ?? null;
  const previewBundle = useMemo(() => (activeDocument ? buildPreviewBundle(activeDocument) : null), [activeDocument]);
  const sourceUrl = activeDocument ? getDocumentSourceUrl(activeDocument) : undefined;
  const currentPage = previewBundle?.pages[activePage] ?? null;

  const activeBlocks = useMemo(
    () => (activeDocument && isDigitalTemplate(activeDocument) ? buildDigitalBlocks(activeDocument) : []),
    [activeDocument],
  );

  const activeDraft = useMemo(() => {
    if (!activeDocument || !isDigitalTemplate(activeDocument)) return null;
    return drafts[activeDocument.id] ?? createInitialDocumentDraft(activeBlocks);
  }, [activeBlocks, activeDocument, drafts]);

  useEffect(() => {
    if (!isStorageReady || !activeDocument || !isDigitalTemplate(activeDocument) || !activeDraft) return;
    if (drafts[activeDocument.id]) return;

    setDrafts((current) => ({ ...current, [activeDocument.id]: activeDraft }));
  }, [activeDocument, activeDraft, drafts, isStorageReady, setDrafts]);

  function openDocument(documentId: string) {
    setActiveId(documentId);
    setActivePage(0);
    setViewMode("preview");
    setEditorMode("none");
  }

  function closeDocument() {
    setActiveId("");
    setActivePage(0);
    setViewMode("preview");
    setEditorMode("none");
  }

  function updateField(fieldId: string, value: string) {
    if (!activeDocument || !isDigitalTemplate(activeDocument)) return;

    setDrafts((current) => {
      const currentDraft = current[activeDocument.id] ?? createInitialDocumentDraft(activeBlocks);
      return {
        ...current,
        [activeDocument.id]: {
          ...currentDraft,
          fields: { ...currentDraft.fields, [fieldId]: value },
        },
      };
    });
  }

  function updateTable(tableId: string, rowIndex: number, columnIndex: number, value: string) {
    if (!activeDocument || !isDigitalTemplate(activeDocument)) return;

    setDrafts((current) => {
      const currentDraft = current[activeDocument.id] ?? createInitialDocumentDraft(activeBlocks);
      const table = currentDraft.tables[tableId]?.map((row) => [...row]) ?? [];
      while (table.length <= rowIndex) table.push([]);
      while (table[rowIndex].length <= columnIndex) table[rowIndex].push("");
      table[rowIndex][columnIndex] = value;

      return {
        ...current,
        [activeDocument.id]: {
          ...currentDraft,
          tables: {
            ...currentDraft.tables,
            [tableId]: table,
          },
        },
      };
    });
  }

  function beginCreate() {
    setEditorState(createDocumentEditorState());
    setEditorMode("create");
  }

  function beginEdit() {
    if (!activeDocument) return;
    setEditorState(createDocumentEditorState(activeDocument));
    setEditorMode("edit");
  }

  function saveEditor() {
    if (editorMode === "create") {
      const newDocument = toCustomDocument(editorState);
      setCustomDocuments((current) => [...current, newDocument]);
      setActiveId(newDocument.id);
      setEditorMode("none");
      return;
    }

    if (!activeDocument) return;

    setOverrides((current) => ({
      ...current,
      [activeDocument.id]: {
        title: editorState.title,
        category: editorState.category,
        stage: editorState.stage,
        kind: editorState.kind,
        summary: editorState.summary,
        owner: editorState.owner,
        uploadStatus: editorState.uploadStatus,
        suggestedUse: editorState.suggestedUse,
        fileName: editorState.fileName,
      },
    }));
    setEditorMode("none");
  }

  return {
    query,
    setQuery,
    stageFilter,
    setStageFilter,
    categoryFilter,
    setCategoryFilter,
    activeId,
    activePage,
    setActivePage,
    viewMode,
    setViewMode,
    editorMode,
    setEditorMode,
    editorState,
    setEditorState,
    stageFilters,
    categoryOptions,
    isStorageReady,
    allDocuments,
    filteredDocuments,
    activeDocument,
    previewBundle,
    sourceUrl,
    currentPage,
    activeBlocks,
    activeDraft,
    openDocument,
    closeDocument,
    updateField,
    updateTable,
    beginCreate,
    beginEdit,
    saveEditor,
  };
}
