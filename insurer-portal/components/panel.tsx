import type { HTMLAttributes, ReactNode } from "react";

import { cn } from "@/lib/utils";

interface PanelProps extends HTMLAttributes<HTMLElement> {
  title?: string;
  description?: string;
  action?: ReactNode;
}

export function Panel({ title, description, action, className, children, ...props }: PanelProps) {
  return (
    <section className={cn("portal-panel rounded-[28px] p-5 lg:p-6", className)} {...props}>
      {(title || description || action) && (
        <div className="mb-5 flex flex-col gap-3 border-b border-[rgb(12_91_65_/_0.1)] pb-4 sm:flex-row sm:items-start sm:justify-between">
          <div className="space-y-1">
            {title ? (
              <div className="flex items-center gap-3">
                <span
                  className="h-3 w-3 rounded-full shadow-[0_0_0_6px_rgb(245_158_11_/_0.12)]"
                  style={{ background: "linear-gradient(135deg, var(--accent), var(--brand))" }}
                />
                <h2 className="font-[family:var(--font-heading)] text-xl font-semibold text-[var(--text)]">
                  {title}
                </h2>
              </div>
            ) : null}
            {description ? <p className="text-sm text-[var(--muted)]">{description}</p> : null}
          </div>
          {action ? <div className="shrink-0">{action}</div> : null}
        </div>
      )}
      {children}
    </section>
  );
}
