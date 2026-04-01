"use client";

import type { LucideIcon } from "lucide-react";

interface WorkspaceIconCardGridProps {
  items: Array<{
    icon: LucideIcon;
    text: string;
  }>;
  className?: string;
  cardClassName?: string;
}

export function WorkspaceIconCardGrid({
  items,
  className,
  cardClassName,
}: WorkspaceIconCardGridProps) {
  return (
    <div className={className ?? "grid gap-3 md:grid-cols-3"}>
      {items.map((item) => {
        const Icon = item.icon;
        return (
          <div
            key={item.text}
            className={
              cardClassName ??
              "rounded-[22px] border border-[rgb(12_91_65_/_0.08)] bg-white/74 p-4"
            }
          >
            <Icon className="h-5 w-5 text-[var(--brand-deep)]" />
            <p className="mt-3 text-sm leading-6 text-[var(--muted)]">{item.text}</p>
          </div>
        );
      })}
    </div>
  );
}
