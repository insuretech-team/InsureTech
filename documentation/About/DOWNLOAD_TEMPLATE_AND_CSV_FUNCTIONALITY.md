# Download Template & CSV Functionality - Complete Code Reference

## Summary
The b2b_portal has a complete employee bulk upload system with a downloadable CSV template. The "Download Template" button triggers a GET request to `/api/employees/template?format=csv` which generates a CSV file with example rows and optional reference plan columns.

---

## 1. DOWNLOAD TEMPLATE BUTTON - data-table.tsx

**Location:** `E:/Projects/InsureTech/b2b_portal/components/dashboard/employees/data-table/data-table.tsx`

**Lines 149-157:** Download Template Button Implementation

```tsx
<Button
  variant="outline"
  className="brand-btn-ghost"
  onClick={() => window.open("/api/employees/template?format=csv", "_blank")}
  type="button"
>
  <LuDownload />
  <span>Download Template</span>
</Button>
```

**Key Details:**
- Uses `window.open()` to open the template URL in a new tab
- Endpoint: `/api/employees/template?format=csv`
- Located in the Employee Table header toolbar
- Accompanied by "Upload (Excel, CSV)" and "Add Employee" buttons

---

## 2. TEMPLATE GENERATION API ROUTE - route.ts

**Location:** `E:/Projects/InsureTech/b2b_portal/app/api/employees/template/route.ts`

**Complete File Content:**

```typescript
/**
 * GET /api/employees/template?format=csv
 *
 * Single table CSV template for bulk employee upload.
 *
 * Columns:
 *   name, employee_id, department_name, email, mobile_number,
 *   date_of_birth, date_of_joining, gender, insurance_category,
 *   coverage_amount (BDT), number_of_dependent, assigned_plan_name,
 *   [gap], available_plan_name (ref), premium_amount_BDT (ref)
 *
 * The last two columns are read-only reference columns showing available
 * plans fetched live from the catalog — users copy a plan name into
 * assigned_plan_name. These columns are ignored by the upload parser.
 *
 * UTF-8 BOM prepended so Excel renders Bengali/Unicode correctly.
 * Dates: DD/MM/YYYY. coverage_amount: plain BDT integer, no currency symbol.
 */

import { NextRequest, NextResponse } from "next/server";

// ─── Employee column headers (uploaded / parsed by gateway) ───────────────────

// Column order matches the portal employee table:
// Name → Employee ID → Department → Email → Mobile → DOB → DOJ → Gender →
// Insurance Category → Assigned Plan → Coverage → Dependents
const EMP_HEADERS = [
  "name",
  "employee_id",
  "department_name",
  "email",
  "mobile_number",
  "date_of_birth",
  "date_of_joining",
  "gender",
  "insurance_category",
  "assigned_plan_name",
  "coverage_amount",
  "number_of_dependent",
];

// ─── Reference column headers removed ────────────────────────────────────────
// Previously we appended reference plan columns to the right of the template.
// This caused "premium_amount (BDT)" and numeric premium values to be parsed as
// department names by the bulk upload handler when the column map misaligned.
// Reference info is now shown only in the instructions, not as extra CSV columns.
const REF_HEADERS: string[] = [];

// ─── Example employee rows (4 rows, Bengali names) ────────────────────────────

// Rows match EMP_HEADERS column order:
// name, employee_id, department_name, email, mobile_number, date_of_birth,
// date_of_joining, gender, insurance_category, assigned_plan_name, coverage_amount, number_of_dependent
const EXAMPLE_ROWS: string[][] = [
  ["মোহাম্মদ রহিম উদ্দিন", "EMP001", "Engineering",     "rahim@company.com",  "+8801712345678", "15/06/1990", "01/01/2023", "MALE",   "HEALTH", "", "500000", "2"],
  ["ফাতেমা বেগম",           "EMP002", "Human Resources", "fatema@company.com", "+8801812345678", "22/03/1988", "15/03/2022", "FEMALE", "LIFE",   "", "300000", "1"],
  ["করিম হোসেন",            "EMP003", "Engineering",     "karim@company.com",  "+8801912345678", "10/11/1992", "01/06/2021", "MALE",   "HEALTH", "", "400000", "3"],
  ["সালমা আক্তার",          "EMP004", "Finance",         "salma@company.com",  "+8801612345678", "05/08/1995", "01/09/2023", "FEMALE", "LIFE",   "", "250000", "0"],
];

// ─── Helpers ──────────────────────────────────────────────────────────────────

function q(value: string): string {
  return `"${value.replace(/\"/g, '\"\"')}"`;
}

// ─── Types ────────────────────────────────────────────────────────────────────

interface CatalogPlan {
  planName?: string;
  insuranceCategory?: string;
  premiumAmount?: string;
}

interface CatalogResponse {
  items?: CatalogPlan[];
}

// ─── Route handler ────────────────────────────────────────────────────────────

export async function GET(req: NextRequest) {
  try {
    // ── Fetch plan catalog (best-effort — silently omit if it fails) ──────────
    let plans: CatalogPlan[] = [];
    try {
      const res = await fetch(`${req.nextUrl.origin}/api/purchase-orders/catalog`, {
        headers: { cookie: req.headers.get("cookie") ?? "" },
        cache: "no-store",
      });
      if (res.ok) {
        const data = (await res.json()) as CatalogResponse;
        plans = Array.isArray(data.items) ? data.items : [];
      }
    } catch {
      // silently skip
    }

    // ── Build rows ────────────────────────────────────────────────────────────
    const hasPlans = plans.length > 0;
    const numRows = Math.max(EXAMPLE_ROWS.length, plans.length);

    // Header row: employee cols + gap + ref cols (only if plans available)
    const headerRow = hasPlans
      ? [...EMP_HEADERS, ...REF_HEADERS].map(q).join(",")
      : EMP_HEADERS.map(q).join(",");

    const rows: string[] = [headerRow];

    // assigned_plan_name column index in EMP_HEADERS
    const planNameColIdx = EMP_HEADERS.indexOf("assigned_plan_name");

    for (let i = 0; i < numRows; i++) {
      // Copy example row (or blank row if beyond example rows)
      const emp = (EXAMPLE_ROWS[i] ?? Array(EMP_HEADERS.length).fill("")).map((v) => v);

      // Fill assigned_plan_name from catalog for example rows
      if (hasPlans && planNameColIdx >= 0 && plans[i % plans.length]) {
        emp[planNameColIdx] = plans[i % plans.length].planName ?? "";
      }

      const empCells = emp.map(q).join(",");

      if (hasPlans) {
        const plan = plans[i];
        // premiumAmount formatted as ৳X,XX,XXX — strip all non-numeric to get plain integer
        const premiumRaw = plan ? (plan.premiumAmount ?? "") : "";
        const premiumPlain = premiumRaw.replace(/\\D/g, "");
        const planName = plan ? (plan.planName ?? "") : "";
        rows.push(`${empCells},${q(planName)},${q(premiumPlain)}`);
      } else {
        rows.push(empCells);
      }
    }

    // UTF-8 BOM so Excel renders Bengali correctly
    const csv = "\\uFEFF" + rows.join("\\r\\n");

    return new Response(csv, {
      status: 200,
      headers: {
        "Content-Type": "text/csv; charset=utf-8",
        "Content-Disposition": 'attachment; filename="employees_template.csv"',
        "Cache-Control": "no-store",
      },
    });
  } catch (err) {
    return NextResponse.json(
      { ok: false, message: err instanceof Error ? err.message : "Failed to generate template" },
      { status: 500 }
    );
  }
}
```

**Key Features:**
- **CSV Format:** UTF-8 BOM prepended for Excel/Bengali character support
- **Headers:** 12 required columns (name, employee_id, department_name, etc.)
- **Example Rows:** 4 Bengali-named example employees
- **Plan Integration:** Fetches plans from `/api/purchase-orders/catalog` to populate assigned_plan_name
- **File Naming:** `employees_template.csv`
- **Date Format:** DD/MM/YYYY
- **Currency:** BDT amounts as plain integers (no currency symbols)

---

## 3. BULK UPLOAD MODAL - bulk-upload-employee-modal.tsx

**Location:** `E:/Projects/InsureTech/b2b_portal/components/modals/bulk-upload-employee-modal.tsx`

**Key Features:**

### 3.1 Modal Props & Types (Lines 54-86)

```typescript
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

interface BulkUploadEmployeeModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  /** Organisation UUID — required so the gateway knows which org to insert into */
  organisationId?: string;
  /** Called after at least one employee was successfully created */
  onSaved?: () => void;
}
```

### 3.2 Template Download Button in Modal (Lines 256-265)

```tsx
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
```

### 3.3 File Upload Handler (Lines 152-192)

```typescript
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
```

### 3.4 File Drop Zone (Lines 269-319)

```tsx
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
```

### 3.5 Result Display Panel (Lines 338-415)

```tsx
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
```

### 3.6 Upload Field Reference Table (Lines 229-254)

Field specifications table included in the modal showing:
- **name** (Required) - Full name of the employee
- **employee_id** (Required) - Unique employee ID (e.g. EMP001)
- **department_name** (Optional) - Plain text name matched or created automatically
- **email** (Optional) - Work email address
- **mobile_number** (Optional) - e.g. +8801712345678
- **date_of_birth** (Optional) - DD/MM/YYYY or YYYY-MM-DD
- **date_of_joining** (Optional) - DD/MM/YYYY or YYYY-MM-DD (defaults to today)
- **gender** (Optional) - MALE / FEMALE / OTHER
- **insurance_category** (Optional) - HEALTH / LIFE / AUTO / TRAVEL
- **assigned_plan_name** (Optional) - Plan name from catalog — matched automatically
- **coverage_amount** (Optional) - Plain number in BDT (e.g. 500000)
- **number_of_dependent** (Optional) - Integer (e.g. 2)

### 3.7 Key Configuration Values (Lines 51)

```typescript
const MAX_FILE_SIZE_BYTES = 32 * 1024 * 1024; // 32 MB — matches gateway limit
```

**Upload Flow:**
1. User clicks "Download Template (CSV)" button
2. File opens in new tab from `/api/employees/template?format=csv`
3. User fills in CSV with employee data
4. User uploads file to modal via drag-and-drop or file picker
5. Modal sends FormData POST to `/api/employees/bulk-upload` with:
   - `file` (the CSV/XLSX file)
   - `business_id` (organisation UUID)
6. Response contains results with created/failed counts and error details
7. Error rows are displayed in a table for user correction

---

## 4. DATA TABLE INTEGRATION - data-table.tsx

**Location:** `E:/Projects/InsureTech/b2b_portal/components/dashboard/employees/data-table/data-table.tsx`

### 4.1 Button Toolbar (Lines 138-168)

```tsx
<div className="flex items-center gap-2">
  <Button
    variant="outline"
    className="brand-btn-gradient"
    onClick={() => setBulkUploadModalOpen(true)}
    disabled={!organisationId}
    type="button"
  >
    <LuUpload />
    <span>Upload (Excel, CSV)</span>
  </Button>
  <Button
    variant="outline"
    className="brand-btn-ghost"
    onClick={() => window.open("/api/employees/template?format=csv", "_blank")}
    type="button"
  >
    <LuDownload />
    <span>Download Template</span>
  </Button>
  <Button
    variant="outline"
    className="brand-btn-gradient"
    onClick={() => setAddEmployeeModalOpen(true)}
    type="button"
    disabled={!organisationId}
  >
    <LuCirclePlus />
    <span>Add Employee</span>
  </Button>
</div>
```

### 4.2 Modal Instances (Lines 263-274)

```tsx
<BulkUploadEmployeeModal
  open={bulkUploadModalOpen}
  onOpenChange={setBulkUploadModalOpen}
  organisationId={organisationId}
  onSaved={() => { setBulkUploadModalOpen(false); onRefresh?.(); }}
/>
```

**Imports Used:**
```tsx
import { LuCirclePlus, LuUpload, LuDownload, LuTrash2, LuLoader } from "react-icons/lu";
import BulkUploadEmployeeModal from "../../../modals/bulk-upload-employee-modal";
```

---

## 5. COMPONENT HIERARCHY

```
EmployeesPage (app/employees/page.tsx)
  └─ EmployeesTable (components/dashboard/employees/employees-table.tsx)
      └─ DataTable (components/dashboard/employees/data-table/data-table.tsx)
          ├─ [Download Template Button] → window.open("/api/employees/template?format=csv")
          ├─ [Upload Button] → Opens BulkUploadEmployeeModal
          ├─ BulkUploadEmployeeModal
          │   ├─ [Download Template button inside modal]
          │   ├─ File drag-drop zone
          │   └─ Result display with error table
          └─ [Add Employee Button] → Opens AddEmployeeModal
```

---

## 6. API ENDPOINTS SUMMARY

| Endpoint | Method | Purpose | Response |
|----------|--------|---------|----------|
| `/api/employees/template` | GET | Generate CSV template with example rows | CSV file with UTF-8 BOM |
| `/api/employees/bulk-upload` | POST | Upload and process employee CSV/XLSX | JSON with created/failed counts |
| `/api/purchase-orders/catalog` | GET | Fetch available insurance plans | JSON with plan list |

---

## 7. KEY TECHNICAL DETAILS

### CSV Generation (template/route.ts)
- **UTF-8 BOM:** `\uFEFF` prepended for Excel compatibility with Bengali characters
- **Example Rows:** 4 Bengali-named employees matching schema
- **Plan Integration:** Fetches from catalog and populates assigned_plan_name column
- **Date Format:** DD/MM/YYYY
- **CSV Quoting:** RFC 4180 compliant with proper quote escaping

### Bulk Upload (bulk-upload-employee-modal.tsx)
- **Form Data:** Multipart/form-data POST (no manual Content-Type header)
- **File Validation:** .csv and .xlsx only, max 32 MB
- **Drag & Drop:** Full support with visual feedback
- **Result Display:** Created count, failed count, detailed error table with row numbers

### Column Mapping
- Case-insensitive column matching via gateway
- Supports any column order in uploaded file
- Automatic department matching/creation by name
- Plan name matching from catalog
- Required fields: name, employee_id
- Optional fields: everything else

