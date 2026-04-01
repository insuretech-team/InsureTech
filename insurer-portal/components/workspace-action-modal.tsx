"use client";

import { X } from "lucide-react";

type FieldType = "text" | "number" | "textarea";

export interface WorkspaceActionField {
  key: string;
  label: string;
  type: FieldType;
  placeholder?: string;
  required?: boolean;
  min?: number;
}

interface WorkspaceActionModalProps {
  open: boolean;
  title: string;
  description: string;
  fields: WorkspaceActionField[];
  values: Record<string, string>;
  submitLabel: string;
  closeLabel?: string;
  cancelLabel?: string;
  onChange: (key: string, value: string) => void;
  onClose: () => void;
  onSubmit: () => void;
}

export function WorkspaceActionModal({
  open,
  title,
  description,
  fields,
  values,
  submitLabel,
  closeLabel = "Close",
  cancelLabel = "Cancel",
  onChange,
  onClose,
  onSubmit,
}: WorkspaceActionModalProps) {
  if (!open) return null;

  return (
    <div className="workspace-modal-backdrop" data-workspace-modal="true">
      <div className="workspace-modal-shell max-w-[640px]">
        <div className="workspace-modal-header">
          <div>
            <h2 className="font-[family:var(--font-heading)] text-3xl font-semibold text-[var(--text)]">{title}</h2>
            <p className="mt-2 text-sm leading-6 text-[var(--muted)]">{description}</p>
          </div>
          <button className="portal-btn portal-btn-secondary" onClick={onClose} type="button">
            <X className="h-4 w-4" />
            {closeLabel}
          </button>
        </div>

        <div className="space-y-4 p-6">
          {fields.map((field) => (
            <label key={field.key} className="block space-y-2">
              <span className="text-sm font-medium text-[var(--muted)]">{field.label}</span>
              {field.type === "textarea" ? (
                <textarea
                  className="portal-textarea"
                  placeholder={field.placeholder}
                  value={values[field.key] ?? ""}
                  onChange={(event) => onChange(field.key, event.target.value)}
                />
              ) : (
                <input
                  className="portal-input"
                  min={field.min}
                  placeholder={field.placeholder}
                  type={field.type}
                  value={values[field.key] ?? ""}
                  onChange={(event) => onChange(field.key, event.target.value)}
                />
              )}
            </label>
          ))}

          <div className="flex flex-wrap justify-end gap-2 pt-2">
            <button className="portal-btn portal-btn-secondary" onClick={onClose} type="button">
              {cancelLabel}
            </button>
            <button className="portal-btn portal-btn-primary" onClick={onSubmit} type="button">
              {submitLabel}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
