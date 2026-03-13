"use client";

/**
 * bulk-upload-employee-modal.tsx
 * ──────────────────────────────
 * Modal for bulk-uploading employees from an Excel (.xlsx) or CSV file.
 *
 * Flow:
 *   1. User downloads the CSV template (Bengali example rows, BOM-prefixed for Excel)
 *   2. User fills in data — use plain department names in the department_name column
 *   3. User picks a file (drag-and-drop or click)
 *   4. Client sends multipart/form-data POST to /api/employees/bulk-upload
 *      with fields: file + business_id
 *   5. Gateway handler (b2b_bulk_upload_handler.go) parses rows and calls
 *      CreateEmployee gRPC for each row
 *   6. We display a structured result: created / failed counts + error rows table
 *
 * Expected columns (case-insensitive, any order — gateway uses alias map):
 *   name, employee_id, department_name, email, mobile_number,
 *   date_of_birth, date_of_joining, gender, insurance_category,
 *   coverage_amount, number_of_dependent, assigned_plan_name
 *
 * Key constraints:
 *   - Do NOT set Content-Type on the fetch — browser sets multipart boundary automatically
 *   - form field "file" must match gateway: r.FormFile("file")
 *   - form field "business_id" must match gateway: r.FormValue("business_id")
 *   - Max file size: 32 MB (gateway: r.ParseMultipartForm(32 << 20))
 *   - Template served from GET /api/employees/template?format=csv (UTF-8 BOM included)
 */

import * as React from "react";
import {
  LuUpload,
  LuLoader,
  LuFileSpreadsheet,
  LuX,
  LuCircleCheck,
  LuCircleAlert,
  LuDownload,
} from "react-icons/lu";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";

// ─── Constants ─────────────────────────────────────────────────────────────────

const MAX_FILE_SIZE_BYTES = 32 * 1024 * 1024; // 32 MB — matches gateway limit

// ─── Types ────────────────────────────────────────────────────────────────────

interface BulkUploadError {
  row: number;
  name?: string;
  message: string;
}

interface BulkUploadResult {
  created: number;
  failed: number;
  total: number;
  errors?: BulkUploadError[];
}

interface BulkUploadResponse {
  ok: boolean;
  message?: string;
  result?: BulkUploadResult;
}

type UploadState = "idle" | "uploading" | "success" | "partial" | "error";


// ─── Props ────────────────────────────────────────────────────────────────────

interface BulkUploadEmployeeModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  /** Organisation UUID — required so the gateway knows which org to insert into */
  organisationId?: string;
  /** Called after at least one employee was successfully created */
  onSaved?: () => void;
}

// ─── Main Modal ───────────────────────────────────────────────────────────────

export default function BulkUploadEmployeeModal({
  open,
  onOpenChange,
  organisationId,
  onSaved,
}: BulkUploadEmployeeModalProps) {
  const [file, setFile] = React.useState<File | null>(null);
  const [fileError, setFileError] = React.useState<string>("");
  const [state, setState] = React.useState<UploadState>("idle");
  const [response, setResponse] = React.useState<BulkUploadResponse | null>(null);
  const [dragOver, setDragOver] = React.useState(false);
  const fileInputRef = React.useRef<HTMLInputElement>(null);

  // ── Reset when modal closes ─────────────────────────────────────────────────
  React.useEffect(() => {
    if (!open) {
      setFile(null);
      setFileError("");
      setState("idle");
      setResponse(null);
      setDragOver(false);
    }
  }, [open]);

  // ── File validation ─────────────────────────────────────────────────────────
  function validateAndSetFile(f: File | null) {
    if (!f) return;
    const ext = f.name.split(".").pop()?.toLowerCase();
    if (ext !== "xlsx" && ext !== "csv") {
      setFileError("Only .csv and .xlsx files are supported");
      return;
    }
    if (f.size > MAX_FILE_SIZE_BYTES) {
      setFileError("File too large — maximum 32 MB");
      return;
    }
    setFile(f);
    setFileError("");
    setState("idle");
    setResponse(null);
  }

  function handleFileInputChange(e: React.ChangeEvent<HTMLInputElement>) {
    validateAndSetFile(e.target.files?.[0] ?? null);
  }

  function handleDrop(e: React.DragEvent) {
    e.preventDefault();
    setDragOver(false);
    validateAndSetFile(e.dataTransfer.files[0] ?? null);
  }

  function clearFile(e: React.MouseEvent) {
    e.stopPropagation();
    setFile(null);
    setFileError("");
    setState("idle");
    setResponse(null);
    if (fileInputRef.current) fileInputRef.current.value = "";
  }

  // ── Upload handler ──────────────────────────────────────────────────────────
  async function handleUpload() {
    if (!file || !organisationId) return;
    setState("uploading");
    setResponse(null);

    try {
      const form = new FormData();
      form.append("file", file);
      form.append("business_id", organisationId);

      // Do NOT set Content-Type — browser sets multipart boundary automatically
      const res = await fetch("/api/employees/bulk-upload", {
        method: "POST",
        body: form,
      });

      const data = (await res.json()) as BulkUploadResponse;
      setResponse(data);

      const r = data.result;
      if (r) {
        if (r.failed === 0 && r.created > 0) {
          setState("success");
        } else if (r.created > 0 && r.failed > 0) {
          setState("partial");
        } else {
          setState("error");
        }
        if (r.created > 0) onSaved?.(); // refresh employee table
      } else {
        // No result object — treat as error
        setState("error");
      }
    } catch (err) {
      setResponse({
        ok: false,
        message: err instanceof Error ? err.message : "Network error — could not upload file",
      });
      setState("error");
    }
  }

  // ── Reset to try another file ───────────────────────────────────────────────
  function handleUploadAnother() {
    setFile(null);
    setFileError("");
    setState("idle");
    setResponse(null);
    if (fileInputRef.current) fileInputRef.current.value = "";
  }

  const result = response?.result;
  const isUploading = state === "uploading";
  const isDone = state === "success" || state === "partial" || state === "error";

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-xl p-0 max-h-[90vh] overflow-y-auto">
        {/* ── Header ── */}
        <DialogHeader className="px-6 py-4 border-b sticky top-0 bg-white z-10">
          <DialogTitle className="text-xl font-semibold flex items-center gap-2">
            <LuFileSpreadsheet className="size-5 text-primary" />
            Bulk Upload Employees
          </DialogTitle>
        </DialogHeader>

        <div className="px-6 py-5 space-y-5">
          {/* ── Instructions + Template Download ── */}
          <div className="rounded-lg border bg-muted/30 px-4 py-3 space-y-2">
            <p className="text-sm font-medium">How to bulk upload:</p>
            <ol className="text-xs text-muted-foreground space-y-1 list-decimal list-inside">
              <li>Download the template CSV below</li>
              <li>Fill in employee data — one employee per row</li>
              <li>Enter department names as plain text — they are auto-matched or created</li>
              <li>Save the file and upload it here</li>
            </ol>

            {/* Field reference table */}
            <div className="mt-2 overflow-x-auto">
              <table className="w-full text-xs border rounded">
                <thead className="bg-muted/60">
                  <tr>
                    <th className="px-2 py-1.5 text-left font-semibold text-muted-foreground">Column</th>
                    <th className="px-2 py-1.5 text-left font-semibold text-muted-foreground">Required</th>
                    <th className="px-2 py-1.5 text-left font-semibold text-muted-foreground">Format / Notes</th>
                  </tr>
                </thead>
                <tbody className="divide-y">
                  <tr><td className="px-2 py-1.5 font-mono text-primary/80">name</td><td className="px-2 py-1.5 text-green-700 font-medium">Yes</td><td className="px-2 py-1.5 text-muted-foreground">Full name of the employee</td></tr>
                  <tr><td className="px-2 py-1.5 font-mono text-primary/80">employee_id</td><td className="px-2 py-1.5 text-green-700 font-medium">Yes</td><td className="px-2 py-1.5 text-muted-foreground">Unique employee ID (e.g. EMP001)</td></tr>
                  <tr><td className="px-2 py-1.5 font-mono text-primary/80">department_name</td><td className="px-2 py-1.5 text-muted-foreground">No</td><td className="px-2 py-1.5 text-muted-foreground">Plain text name (e.g. Engineering) — matched or created automatically</td></tr>
                  <tr><td className="px-2 py-1.5 font-mono text-primary/80">email</td><td className="px-2 py-1.5 text-muted-foreground">No</td><td className="px-2 py-1.5 text-muted-foreground">Work email address</td></tr>
                  <tr><td className="px-2 py-1.5 font-mono text-primary/80">mobile_number</td><td className="px-2 py-1.5 text-muted-foreground">No</td><td className="px-2 py-1.5 text-muted-foreground">e.g. +8801712345678</td></tr>
                  <tr><td className="px-2 py-1.5 font-mono text-primary/80">date_of_birth</td><td className="px-2 py-1.5 text-muted-foreground">No</td><td className="px-2 py-1.5 text-muted-foreground">DD/MM/YYYY or YYYY-MM-DD</td></tr>
                  <tr><td className="px-2 py-1.5 font-mono text-primary/80">date_of_joining</td><td className="px-2 py-1.5 text-muted-foreground">No</td><td className="px-2 py-1.5 text-muted-foreground">DD/MM/YYYY or YYYY-MM-DD (defaults to today)</td></tr>
                  <tr><td className="px-2 py-1.5 font-mono text-primary/80">gender</td><td className="px-2 py-1.5 text-muted-foreground">No</td><td className="px-2 py-1.5 text-muted-foreground">MALE / FEMALE / OTHER</td></tr>
                  <tr><td className="px-2 py-1.5 font-mono text-primary/80">insurance_category</td><td className="px-2 py-1.5 text-muted-foreground">No</td><td className="px-2 py-1.5 text-muted-foreground">HEALTH / LIFE / AUTO / TRAVEL</td></tr>
                  <tr><td className="px-2 py-1.5 font-mono text-primary/80">assigned_plan_name</td><td className="px-2 py-1.5 text-muted-foreground">No</td><td className="px-2 py-1.5 text-muted-foreground">Plan name from catalog — matched automatically</td></tr>
                  <tr><td className="px-2 py-1.5 font-mono text-primary/80">coverage_amount</td><td className="px-2 py-1.5 text-muted-foreground">No</td><td className="px-2 py-1.5 text-muted-foreground">Plain number in BDT (e.g. 500000)</td></tr>
                  <tr><td className="px-2 py-1.5 font-mono text-primary/80">number_of_dependent</td><td className="px-2 py-1.5 text-muted-foreground">No</td><td className="px-2 py-1.5 text-muted-foreground">Integer (e.g. 2)</td></tr>
                </tbody>
              </table>
            </div>

            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => window.open("/api/employees/template?format=csv", "_blank")}
              className="mt-1 gap-1.5"
            >
              <LuDownload className="size-3.5" />
              Download Template (CSV)
            </Button>
          </div>

          {/* ── File Drop Zone ── */}
          {!isDone && (
            <div
              className={[
                "relative flex flex-col items-center justify-center gap-3 rounded-xl border-2 border-dashed p-8 text-center transition-colors cursor-pointer",
                dragOver
                  ? "border-primary bg-primary/5"
                  : "border-muted-foreground/30 hover:border-primary/50 hover:bg-muted/20",
              ].join(" ")}
              onDragOver={(e) => { e.preventDefault(); setDragOver(true); }}
              onDragLeave={() => setDragOver(false)}
              onDrop={handleDrop}
              onClick={() => fileInputRef.current?.click()}
            >
              <input
                ref={fileInputRef}
                type="file"
                accept=".xlsx,.csv"
                className="hidden"
                onChange={handleFileInputChange}
              />

              {file ? (
                <>
                  <LuFileSpreadsheet className="size-10 text-primary" />
                  <div>
                    <p className="font-medium text-sm">{file.name}</p>
                    <p className="text-xs text-muted-foreground">
                      {(file.size / 1024).toFixed(1)} KB
                    </p>
                  </div>
                  <button
                    type="button"
                    className="absolute top-3 right-3 rounded-full p-1 hover:bg-muted"
                    onClick={clearFile}
                  >
                    <LuX className="size-4 text-muted-foreground" />
                  </button>
                </>
              ) : (
                <>
                  <LuUpload className="size-10 text-muted-foreground/60" />
                  <div>
                    <p className="font-medium text-sm">Drop your file here</p>
                    <p className="text-xs text-muted-foreground">
                      or click to browse — .xlsx and .csv supported (max 32 MB)
                    </p>
                  </div>
                </>
              )}
            </div>
          )}

          {/* ── File validation error ── */}
          {fileError && (
            <p className="text-xs text-red-500 flex items-center gap-1.5">
              <LuCircleAlert className="size-3.5 shrink-0" />
              {fileError}
            </p>
          )}

          {/* ── Uploading indicator ── */}
          {isUploading && (
            <div className="flex items-center gap-2 text-sm text-muted-foreground py-2">
              <LuLoader className="size-4 animate-spin text-primary" />
              Uploading and processing file…
            </div>
          )}

          {/* ── Result Panel ── */}
          {isDone && result && (
            <div
              className={[
                "rounded-lg border p-4 space-y-3",
                state === "success"
                  ? "border-green-200 bg-green-50"
                  : state === "partial"
                  ? "border-yellow-200 bg-yellow-50"
                  : "border-red-200 bg-red-50",
              ].join(" ")}
            >
              {/* Status heading */}
              <div className="flex items-center gap-2">
                {state === "success" ? (
                  <LuCircleCheck className="size-5 text-green-500 shrink-0" />
                ) : state === "partial" ? (
                  <LuCircleAlert className="size-5 text-yellow-500 shrink-0" />
                ) : (
                  <LuCircleAlert className="size-5 text-red-500 shrink-0" />
                )}
                <p className="text-sm font-medium">
                  {state === "success"
                    ? "✅ All employees uploaded successfully"
                    : state === "partial"
                    ? "⚠️ Some rows were skipped — valid rows were saved"
                    : "❌ Upload failed — no employees were saved"}
                </p>
              </div>

              {/* Gateway message */}
              <p className="text-xs text-muted-foreground leading-relaxed">
                {response?.message}
              </p>

              <div className="flex gap-6 text-sm">
                <span className="text-green-700">
                  ✓ Saved: <strong>{result.created}</strong>
                </span>
                {result.failed > 0 && (
                  <span className="text-red-700">
                    ✗ Skipped: <strong>{result.failed}</strong>
                  </span>
                )}
                <span className="text-muted-foreground">
                  Total rows: {result.total}
                </span>
              </div>

              {/* Error rows table */}
              {result.errors && result.errors.length > 0 && (
                <div className="mt-2">
                  <p className="text-xs font-semibold text-muted-foreground uppercase tracking-wide mb-1">
                    Skipped Rows — fix these and re-upload
                  </p>
                  <div className="max-h-48 overflow-y-auto rounded border border-red-200">
                    <table className="w-full text-xs">
                      <thead className="bg-red-100/60">
                        <tr>
                          <th className="px-3 py-1.5 text-left font-medium text-red-800">Row</th>
                          <th className="px-3 py-1.5 text-left font-medium text-red-800">Name</th>
                          <th className="px-3 py-1.5 text-left font-medium text-red-800">Error</th>
                        </tr>
                      </thead>
                      <tbody>
                        {result.errors.map((e, i) => (
                          <tr key={i} className="border-t border-red-100 last:border-0">
                            <td className="px-3 py-1.5 font-mono text-red-700">{e.row}</td>
                            <td className="px-3 py-1.5 text-red-700">{e.name ?? "—"}</td>
                            <td className="px-3 py-1.5 text-red-700">{e.message}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              )}
            </div>
          )}

          {/* ── Network / generic error (no result object) ── */}
          {state === "error" && !result && (
            <div className="rounded-lg border border-red-200 bg-red-50 p-4">
              <div className="flex items-center gap-2">
                <LuCircleAlert className="size-4 text-red-500 shrink-0" />
                <p className="text-sm text-red-700">
                  {response?.message ?? "Upload failed. Please try again."}
                </p>
              </div>
            </div>
          )}
        </div>

        {/* ── Footer ── */}
        <DialogFooter className="px-6 pb-5 gap-2">
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={isUploading}
          >
            {isDone ? "Close" : "Cancel"}
          </Button>

          {isDone && (
            <Button
              type="button"
              variant="outline"
              onClick={handleUploadAnother}
            >
              <LuUpload className="size-4" />
              Upload Another File
            </Button>
          )}

          {!isDone && (
            <Button
              type="button"
              onClick={handleUpload}
              disabled={!file || isUploading || !organisationId || Boolean(fileError)}
              className="h-10 px-6 text-white bg-gradient-to-r from-primary to-accent hover:opacity-95"
            >
              {isUploading ? (
                <span className="flex items-center gap-2">
                  <LuLoader className="animate-spin size-4" />
                  Uploading…
                </span>
              ) : (
                <span className="flex items-center gap-2">
                  <LuUpload className="size-4" />
                  Upload
                </span>
              )}
            </Button>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
