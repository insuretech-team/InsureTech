# Employee Bulk Upload & Template Download — Implementation Plan

## Context & Current State

### What Already Exists ✅
- **Gateway handler**: `POST /v1/b2b/employees/bulk-upload` — `b2b_bulk_upload_handler.go`
- **BFF route**: `POST /api/employees/bulk-upload` — proxies raw multipart to gateway
- **DataTable**: `Upload (Excel, CSV)` button exists, opens `BulkUploadEmployeeModal`
- **Modal import**: `BulkUploadEmployeeModal` is imported in `data-table.tsx` but the **file does not exist yet**
- **Add Employee Modal**: `components/modals/add-employee-modal.tsx` — reference for UI patterns
- **Date format**: HTML `<input type="date">` → `YYYY-MM-DD` (ISO 8601) — gateway `normalizeDate()` accepts multiple formats
- **Currency**: All monetary values in **BDT (Bangladeshi Taka ৳)**, stored as `int64` paisa (×100)

### What is Missing ❌
1. `components/modals/bulk-upload-employee-modal.tsx` — the modal component
2. `app/api/employees/template/route.ts` — template download BFF route
3. Bengali names in example template rows
4. Download Template button wired up in the UI

---

## Gateway Bulk Upload — Exact Contract

**Endpoint**: `POST /v1/b2b/employees/bulk-upload`
**Content-Type**: `multipart/form-data`

**Form fields**:
| Field | Type | Required | Notes |
|---|---|---|---|
| `file` | File | ✅ | `.csv` or `.xlsx` |
| `business_id` | string | ✅ | Falls back to `X-Business-ID` header |

**Expected column headers** (case-insensitive, any order — handler uses alias map):

| Canonical Field | Accepted Aliases |
|---|---|
| `name` | name, full name, employee name, fullname |
| `employee_id` | employee_id, employee id, emp id, emp_id, staff id |
| `department_id` | department_id, department id, dept id, dept_id |
| `email` | email, work email |
| `mobile_number` | mobile_number, mobile number, mobile, phone, phone_number |
| `date_of_birth` | date_of_birth, date of birth, dob, birth_date |
| `date_of_joining` | date_of_joining, date of joining, joining_date, doj, start_date |
| `gender` | gender, sex |
| `insurance_category` | insurance_category, insurance category, insurance_type |
| `coverage_amount` | coverage_amount, coverage amount, coverage |
| `number_of_dependent` | number_of_dependent, number of dependents, dependents |

**Gender values accepted**: M, MALE, F, FEMALE, O, OTHER

**Insurance category values accepted**: LIFE (1), HEALTH (2), AUTO/MOTOR (3), TRAVEL (4), FIRE (5), MARINE (6), PROPERTY (7), LIABILITY (8)

**Date formats accepted**: YYYY-MM-DD, DD/MM/YYYY, MM/DD/YYYY, YYYY/MM/DD, DD-MM-YYYY, DD Mon YYYY

**Coverage amount**: Numeric, commas stripped. Stored as BDT paisa (×100). Example: `50000` = ৳500.00

**Response**:
```json
{
  "ok": true,
  "message": "12 created, 0 failed out of 12 rows",
  "result": {
    "created": 12,
    "failed": 0,
    "total": 12,
    "errors": []
  }
}
```

**Error response** (partial failure, HTTP 200):
```json
{
  "ok": false,
  "message": "10 created, 2 failed out of 12 rows",
  "result": {
    "created": 10,
    "failed": 2,
    "total": 12,
    "errors": [
      { "row": 3, "name": "করিম মিয়া", "message": "employee_id is required" },
      { "row": 7, "name": "সালমা বেগম", "message": "duplicate employee_id" }
    ]
  }
}
```

---

## Template Content

### CSV Template (`employees_template.csv`)

The template must use **Bengali names**, **BDT currency**, and **DD/MM/YYYY date format** (which `normalizeDate()` accepts).

```
name,employee_id,department_id,email,mobile_number,date_of_birth,date_of_joining,gender,insurance_category,coverage_amount,number_of_dependent
মোহাম্মদ রহিম উদ্দিন,EMP001,REPLACE_WITH_DEPT_UUID,rahim@company.com,+8801712345678,15/06/1990,01/01/2023,MALE,HEALTH,500000,2
ফাতেমা বেগম,EMP002,REPLACE_WITH_DEPT_UUID,fatema@company.com,+8801812345678,22/03/1988,15/03/2022,FEMALE,LIFE,300000,1
করিম হোসেন,EMP003,REPLACE_WITH_DEPT_UUID,karim@company.com,+8801912345678,10/11/1992,01/06/2021,MALE,HEALTH,400000,3
সালমা আক্তার,EMP004,REPLACE_WITH_DEPT_UUID,salma@company.com,+8801612345678,05/08/1995,01/09/2023,FEMALE,LIFE,250000,0
```

**Notes for template**:
- `department_id`: Must be a real UUID from the org's departments. The modal should help users find it.
- `coverage_amount`: Plain number in BDT (৳). No currency symbol. Example: `500000` = ৳5,00,000
- `date_of_birth` / `date_of_joining`: Use DD/MM/YYYY format
- `gender`: MALE / FEMALE / OTHER
- `insurance_category`: HEALTH / LIFE / AUTO / TRAVEL / FIRE / MARINE / PROPERTY / LIABILITY

---

## Implementation Steps

### Step 1 — BFF Template Download Route

**File to create**: `app/api/employees/template/route.ts`

```typescript
// GET /api/employees/template?format=csv (default) | format=xlsx
// Returns downloadable employee upload template file.
// No auth required beyond session check (anyone logged in can download the template).

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url);
  const format = searchParams.get("format") ?? "csv";

  const headers = [
    "name", "employee_id", "department_id", "email", "mobile_number",
    "date_of_birth", "date_of_joining", "gender", "insurance_category",
    "coverage_amount", "number_of_dependent"
  ];

  const exampleRows = [
    [
      "মোহাম্মদ রহিম উদ্দিন", "EMP001", "REPLACE_WITH_DEPT_UUID",
      "rahim@company.com", "+8801712345678", "15/06/1990", "01/01/2023",
      "MALE", "HEALTH", "500000", "2"
    ],
    [
      "ফাতেমা বেগম", "EMP002", "REPLACE_WITH_DEPT_UUID",
      "fatema@company.com", "+8801812345678", "22/03/1988", "15/03/2022",
      "FEMALE", "LIFE", "300000", "1"
    ],
    [
      "করিম হোসেন", "EMP003", "REPLACE_WITH_DEPT_UUID",
      "karim@company.com", "+8801912345678", "10/11/1992", "01/06/2021",
      "MALE", "HEALTH", "400000", "3"
    ],
  ];

  // Build CSV (UTF-8 with BOM for Excel Bengali rendering)
  const BOM = "\uFEFF";
  const csvRows = [headers, ...exampleRows]
    .map(row => row.map(cell => `"${cell.replace(/"/g, '""')}"`).join(","))
    .join("\r\n");
  const csv = BOM + csvRows;

  return new Response(csv, {
    status: 200,
    headers: {
      "Content-Type": "text/csv; charset=utf-8",
      "Content-Disposition": 'attachment; filename="employees_template.csv"',
      "Cache-Control": "no-store",
    },
  });
}
```

**CRITICAL**: Add `\uFEFF` BOM prefix so Excel renders Bengali UTF-8 correctly.

---

### Step 2 — BulkUploadEmployeeModal Component

**File to create**: `components/modals/bulk-upload-employee-modal.tsx`

**Props interface**:
```typescript
type BulkUploadEmployeeModalProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  organisationId?: string;
  onSaved?: () => void;
};
```

**State machine**:
```
"idle"      → user selects file
"uploading" → POST /api/employees/bulk-upload
"success"   → all rows created (result.failed === 0)
"partial"   → some created, some failed (result.created > 0 && result.failed > 0)
"error"     → all failed or upload error
```

**UI Sections**:

1. **Instructions panel** (always visible at top):
   - "Download the template, fill in employee data, then upload the file."
   - Button: `📥 Download Template (CSV)` → `window.open("/api/employees/template?format=csv")`
   - Note: "department_id must be a valid UUID from your organisation's departments"
   - Note: "Dates format: DD/MM/YYYY | Currency: BDT (৳) — enter plain numbers e.g. 500000"

2. **File dropzone**:
   - `<input type="file" accept=".csv,.xlsx">` styled as drag-drop zone
   - Show selected filename + size once chosen
   - Supported: Excel (.xlsx), CSV (.csv)

3. **Upload button**: disabled if no file or no organisationId

4. **Progress indicator**: spinner during upload

5. **Result panel** (after upload):
   - Green: `✅ X employees created successfully`
   - Amber: `⚠️ X created, Y failed`
   - Red: `❌ Upload failed`
   - Error table: columns Row | Name | Error — shown when `result.errors.length > 0`
   - "Upload Another File" button to reset state

**Implementation**:
```typescript
const handleUpload = async () => {
  if (!file || !organisationId) return;
  setState("uploading");

  const form = new FormData();
  form.append("file", file);
  form.append("business_id", organisationId);

  try {
    const res = await fetch("/api/employees/bulk-upload", {
      method: "POST",
      body: form,
      // No Content-Type header — browser sets multipart boundary automatically
    });
    const data = await res.json();
    setResult(data);
    setState(data.result?.failed === 0 ? "success" : "partial");
    if (data.result?.created > 0) onSaved?.(); // refresh table
  } catch {
    setState("error");
    setResult({ message: "Network error — could not upload file" });
  }
};
```

**File selection handler**:
```typescript
const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
  const f = e.target.files?.[0];
  if (!f) return;
  // Validate extension
  const name = f.name.toLowerCase();
  if (!name.endsWith(".csv") && !name.endsWith(".xlsx")) {
    setFileError("Only .csv and .xlsx files are supported");
    return;
  }
  // Validate size (32MB max — gateway limit)
  if (f.size > 32 * 1024 * 1024) {
    setFileError("File too large — maximum 32MB");
    return;
  }
  setFile(f);
  setFileError("");
  setState("idle");
  setResult(null);
};
```

---

### Step 3 — Wire Download Template in data-table.tsx

The `Export (Excel, Pdf, CSV)` button currently does nothing. Wire it as:

**File**: `components/dashboard/employees/data-table/data-table.tsx`

Replace the Export button:
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

---

### Step 4 — Result Display Details

**After successful upload** — show in modal:
```
✅ সফলভাবে আপলোড হয়েছে (Successfully Uploaded)
12 employees created, 0 failed out of 12 rows
[Close] [Upload Another]
```

**After partial failure**:
```
⚠️ আংশিক সফল (Partially Successful)
10 created, 2 failed out of 12 rows

| Row | Name          | Error                    |
|-----|---------------|--------------------------|
|  3  | করিম মিয়া   | employee_id is required  |
|  7  | সালমা বেগম   | duplicate employee_id    |

[Upload Another] [Close]
```

---

## File Structure

```
b2b_portal/
├── app/
│   └── api/
│       └── employees/
│           └── template/
│               └── route.ts          ← NEW (Step 1)
└── components/
    └── modals/
        └── bulk-upload-employee-modal.tsx  ← NEW (Step 2)
```

**Already exists (no changes needed)**:
- `app/api/employees/bulk-upload/route.ts` ✅
- `components/dashboard/employees/data-table/data-table.tsx` ✅ (only Export button needs wiring — Step 3)

---

## Key Constraints for Agent

1. **Bengali names in template** — use Unicode Bengali script. Add UTF-8 BOM (`\uFEFF`) to CSV so Excel renders correctly.

2. **Currency** — always BDT (৳). Template comment: "enter plain number e.g. 500000 = ৳5,00,000". No ৳ symbol in the cell itself — gateway parses plain float.

3. **Date format** — use `DD/MM/YYYY` in template examples (Bengali/BD standard). Gateway `normalizeDate()` accepts: `YYYY-MM-DD`, `DD/MM/YYYY`, `MM/DD/YYYY`, `DD-MM-YYYY`, `DD Mon YYYY`.

4. **department_id** — must be a real UUID. Template shows `REPLACE_WITH_DEPT_UUID`. Consider adding a helper in the modal: a collapsible "Available Departments" table fetched from `/api/departments?business_id=...`.

5. **No Content-Type header** on fetch — when posting `FormData`, let the browser set the multipart boundary automatically. Adding `Content-Type: multipart/form-data` manually breaks the boundary.

6. **File form field name** — must be exactly `file` (gateway: `r.FormFile("file")`).

7. **business_id form field** — must be exactly `business_id` (gateway: `r.FormValue("business_id")`).

8. **Max file size** — 32MB (gateway: `r.ParseMultipartForm(32 << 20)`).

9. **Modal pattern** — follow `add-employee-modal.tsx` for Dialog, DialogContent, DialogHeader, DialogFooter, Button, ToastBanner patterns.

10. **onSaved callback** — call `onSaved?.()` when `result.created > 0` to refresh the employee table.

11. **`bulk-upload-employee-modal.tsx` import** — already imported in `data-table.tsx` as:
    ```typescript
    import BulkUploadEmployeeModal from "../../../modals/bulk-upload-employee-modal";
    ```
    The component receives: `open`, `onOpenChange`, `organisationId`, `onSaved`.

---

## Departments Helper (Optional but Recommended)

Inside the modal, show a collapsible "Available Departments" section:

```typescript
// Fetch departments for the org so user can copy UUIDs into their template
useEffect(() => {
  if (!organisationId) return;
  fetch(`/api/departments?business_id=${organisationId}`)
    .then(r => r.json())
    .then(data => setDepartments(data.departments ?? []));
}, [organisationId]);
```

Display as a small table:
```
| Department Name | UUID (copy to paste in department_id column) |
|-----------------|----------------------------------------------|
| Operations      | 3fa85f64-5717-4562-b3fc-2c963f66afa6        |
| HR              | 8d3a1c22-ff1d-4e5a-9e3c-7a5b1d2c8f01        |
```

This prevents the most common error: incorrect `department_id`.

---

## Testing Checklist

- [ ] Download template CSV — opens/downloads correctly in browser
- [ ] Bengali characters render correctly in Excel (BOM required)
- [ ] Upload CSV with Bengali names — employees created
- [ ] Upload XLSX — employees created
- [ ] Upload file with missing `name` — row shows error, others succeed
- [ ] Upload file with invalid `department_id` UUID — row error displayed
- [ ] Upload file with `coverage_amount` = `500000` — stored as ৳5,00,000 (BDT)
- [ ] Upload file with date `15/06/1990` — accepted by `normalizeDate()`
- [ ] After successful upload, employee table refreshes
- [ ] Error table shows row number, name, error message
- [ ] File > 32MB shows client-side error before upload
- [ ] Non-.csv/.xlsx file shows client-side validation error
