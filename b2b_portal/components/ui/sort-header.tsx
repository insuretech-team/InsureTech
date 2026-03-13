/**
 * sort-header.tsx
 * ────────────────
 * Reusable sortable column header for TanStack Table.
 * Replaces the identical SortHeader component duplicated in every columns.tsx.
 */
"use client";

interface SortHeaderProps {
  title: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  column: any;
}

export function SortHeader({ title, column }: SortHeaderProps) {
  const dir = column.getIsSorted();
  return (
    <button
      className="inline-flex items-center gap-2 select-none"
      onClick={() => column.toggleSorting(dir === "asc")}
      title="Sort"
      type="button"
    >
      {title}
      <span className="text-gray-400 text-xs">
        {dir === "asc" ? "▲" : dir === "desc" ? "▼" : "↕"}
      </span>
    </button>
  );
}
