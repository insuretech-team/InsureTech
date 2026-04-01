"use client";

interface WorkspaceNoteListProps {
  items: string[];
  className?: string;
  itemClassName?: string;
}

export function WorkspaceNoteList({ items, className, itemClassName }: WorkspaceNoteListProps) {
  return (
    <div className={className ?? "space-y-3"}>
      {items.map((item) => (
        <div
          key={item}
          className={
            itemClassName ??
            "rounded-[20px] border border-[rgb(12_91_65_/_0.08)] bg-white/76 px-4 py-3 text-sm text-[var(--muted)]"
          }
        >
          {item}
        </div>
      ))}
    </div>
  );
}
