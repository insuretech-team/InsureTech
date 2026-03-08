# B2B Portal Frontend - Comprehensive Exploration Summary

## Project Overview
**Location:** `E:\Projects\InsureTech\b2b_portal`
**Type:** Next.js 16 with TypeScript, React 19
**Key Framework:** Next.js App Router, TailwindCSS, Radix UI, TanStack React Table

## 1. Key Files Identified

### 1.1 Add Employee Modal
**File:** `components/modals/add-employee-modal.tsx`
- **Purpose:** Modal for creating and editing employees
- **Features:**
  - Two modes: create and edit (controlled via `employeeUuid` prop)
  - Three sections: Personal Info, Employment, Insurance
  - Handles department dropdown fetching from `/api/departments`
  - Form submission via `useEmployeeForm` hook
  - Toast notifications on success/error
  - Full form validation
- **Key Props:**
  - `open`: boolean to control modal visibility
  - `employeeUuid`: optional, triggers edit mode when provided
  - `organisationId`: used to fetch departments in create mode
  - `initialValues`: partial form values (typically sparse from table row)
  - `onSaved`: callback after successful create/update

### 1.2 Employees Page & Table
**Files:**
- `app/employees/page.tsx` - Simple wrapper
- `components/dashboard/employees/employees-table.tsx` - Main component
- `components/dashboard/employees/data-table/columns.tsx` - Column definitions

**Features:**
- Role-based org selection (B2B admin locked to own org, super_admin sees dropdown)
- Employee data table with columns: name, ID, department, insurance category, plan, coverage, premium, dependents, status
- Action buttons: View (read-only), Edit, Delete
- Uses `employeeClient` for API calls
- Handles org switching and employee list refresh

### 1.3 API Routes - Employees
**Files:**
- `app/api/employees/route.ts` - GET (list) and POST (create)
- `app/api/employees/[id]/route.ts` - GET (single), PATCH (update), DELETE

**Key Features:**
- SDK-based backend communication via `makeSdkClient`
- Session header resolution (`resolvePortalHeaders`)
- Auto-resolves `business_id` from session or query parameter
- Maps protobuf EmployeeView to UI Employee type
- Full CRUD support
- Handles money conversion (amount in cents to decimal)

### 1.4 Client Library
**File:** `src/lib/clients/employee-client.ts`
- Browser-side fetch wrapper for `/api/employees` routes
- Exports:
  - `EmployeeCreatePayload` - shape for creating employees
  - `EmployeeUpdatePayload` - shape for updating (includes optional fields + status)
  - `EmployeeFullRecord` - complete record returned by GET /api/employees/[id]
  - `employeeClient` object with methods: `list()`, `get(id)`, `create()`, `update()`, `delete()`

### 1.5 Form Hook
**File:** `src/hooks/useEmployeeForm.ts`
- Manages employee form state for both create and edit modes
- Validation: name, employeeId, businessId (create only), departmentId, dateOfJoining
- In edit mode, automatically fetches full employee record via `employeeClient.get()`
- Parses Money objects from API into decimal strings
- Submit handler marshals form values to API payload
- Exports: `values`, `errors`, `submitting`, `loadingRecord`, `setField`, `reset`, `submit`

### 1.6 Type Definitions
**Files:**
- `src/lib/types/b2b.ts` - Core B2B entities (Employee, Organisation, Department, etc.)
- `src/lib/types/employee-form.ts` - Form-specific types

**Employee Type:**
```typescript
interface Employee {
  id: string;
  name: string;
  employeeID: string;
  department: string;
  insuranceCategory?: string;
  assignedPlan?: string;
  coverage: string;
  premiumAmount: string;
  status: EmployeeStatus;
  numberOfDependent: number;
}
```

**EmployeeFormValues:**
```typescript
interface EmployeeFormValues {
  name: string;
  employeeId: string;
  businessId: string;
  email: string;
  mobileNumber: string;
  departmentId: string;
  dateOfBirth: string;
  dateOfJoining: string;
  gender: EmployeeGender;
  insuranceCategory: number;
  assignedPlanId: string;
  coverageAmount: string;
  numberOfDependent: number;
}
```

## 2. Document Generation Infrastructure

### 2.1 Protobuf Definitions Located
**Path:** `src/lib/proto-generated/insuretech/document/`

**Entity Files:**
- `entity/v1/document_generation_pb.ts` - DocumentGeneration message
- `entity/v1/document_template_pb.ts` - DocumentTemplate message
- `services/v1/document_service_pb.ts` - Document service RPC definitions

### 2.2 DocumentGeneration Type
**Location:** `document_generation_pb.ts`

```typescript
type DocumentGeneration = {
  id: string;                          // generation_id (PK)
  documentTemplateId: string;          // reference to template
  entityType: string;                  // entity being documented
  entityId: string;                    // entity UUID
  data: string;                        // JSON data payload
  status: GenerationStatus;            // PENDING, PROCESSING, COMPLETED, FAILED
  fileUrl: string;                     // S3/storage URL
  fileSizeBytes: bigint;               // generated file size
  qrCodeData: string;                  // embedded QR code
  generatedBy: string;                 // user who triggered generation
  generatedAt?: Timestamp;             // generation time
  auditInfo?: AuditInfo;               // audit trail
}
```

**GenerationStatus Enum:**
- GENERATION_STATUS_UNSPECIFIED (0)
- GENERATION_STATUS_PENDING (1)
- GENERATION_STATUS_PROCESSING (2)
- GENERATION_STATUS_COMPLETED (3)
- GENERATION_STATUS_FAILED (4)

### 2.3 DocumentTemplate Type
**Location:** `document_template_pb.ts`

```typescript
type DocumentTemplate = {
  id: string;                    // template_id (PK)
  name: string;                  // template name
  type: DocumentType;            // document category
  description: string;           // template description
  templateContent: string;       // template markup/liquid/handlebars
  outputFormat: OutputFormat;    // PDF, HTML, DOCX
  variables: string;             // JSON list of template variables
  version: number;               // template version
  isActive: boolean;             // active/inactive flag
  auditInfo?: AuditInfo;         // audit trail
}
```

**DocumentType Enum:**
- DOCUMENT_TYPE_UNSPECIFIED (0)
- DOCUMENT_TYPE_POLICY_CERTIFICATE (1)
- DOCUMENT_TYPE_POLICY_SCHEDULE (2)
- DOCUMENT_TYPE_CLAIM_FORM (3)
- DOCUMENT_TYPE_ENDORSEMENT_NOTICE (4)
- DOCUMENT_TYPE_RENEWAL_NOTICE (5)
- DOCUMENT_TYPE_CANCELLATION_NOTICE (6)
- DOCUMENT_TYPE_RECEIPT (7)
- DOCUMENT_TYPE_INVOICE (8)

**OutputFormat Enum:**
- OUTPUT_FORMAT_UNSPECIFIED (0)
- OUTPUT_FORMAT_PDF (1)
- OUTPUT_FORMAT_HTML (2)
- OUTPUT_FORMAT_DOCX (3)

### 2.4 Document Service
**File:** `document_service_pb.ts`
- Contains RPC service definitions for document generation and management
- Service methods for:
  - Creating document generations
  - Listing document generations
  - Getting generation status
  - Querying document templates
  - Managing templates (CRUD)

## 3. SDK Client Architecture

### 3.1 B2B SDK Client
**File:** `src/lib/sdk/b2b-sdk-client.ts`

**Key Features:**
- Factory function `makeSdkClient(request, sessionOverrides?)` creates SDK client
- Auto-extracted CSRF token from cookies
- Session header forwarding: `x-portal`, `x-business-id`, `x-user-id`, `x-tenant-id`
- Authentication via cookie-based server sessions (not API keys)
- Wraps @lifeplus/insuretech-sdk auto-generated methods

**Available SDK Methods:**
- Employee: `listEmployees`, `createEmployee`, `getEmployee`, `updateEmployee`, `deleteEmployee`
- Department: `listDepartments`, `createDepartment`, `getDepartment`, `updateDepartment`, `deleteDepartment`
- Organisation: `listOrganisations`, `createOrganisation`, `deleteOrganisation`, `getOrganisation`, `updateOrganisation`, `listOrgMembers`, `addOrgMember`, `assignOrgAdmin`, `removeOrgMember`
- Purchase Order: `listPurchaseOrders`, `createPurchaseOrder`, `getPurchaseOrder`, `listPurchaseOrderCatalog`

**Direct HTTP Fallback:**
- `makeDirectHttp()` for endpoints not yet exposed as SDK methods
- Used for `/assign-admin`, `/members` endpoints
- Shares same CSRF/session auth as SDK client

### 3.2 Other Client Libraries
**Files:**
- `src/lib/clients/api-client.ts` - Base fetch wrapper with JSON parsing
- `src/lib/clients/auth-client.ts` - Authentication operations
- `src/lib/clients/department-client.ts` - Department API wrapper
- `src/lib/clients/organisation-client.ts` - Organisation API wrapper
- `src/lib/clients/purchase-order-client.ts` - Purchase order API wrapper
- `src/lib/clients/b2b-dashboard-client.ts` - Dashboard stats/activity

## 4. API Route Pattern

### Standard Pattern (Employees Example):

**GET /api/employees** (list)
- Query params: `page_size`, `business_id`, `department_id`
- Returns: `{ ok, employees, message }`
- Uses SDK: `listEmployees()`

**POST /api/employees** (create)
- Body: `EmployeeCreatePayload`
- Returns: `{ ok, message, employee }`
- Uses SDK: `createEmployee()`

**GET /api/employees/[id]** (get single)
- Returns: `{ ok, employee, message }` with all form fields
- Uses SDK: `getEmployee()`

**PATCH /api/employees/[id]** (update)
- Body: `EmployeeUpdatePayload`
- Returns: `{ ok, message, employee }`
- Uses SDK: `updateEmployee()`

**DELETE /api/employees/[id]**
- Returns: `{ ok, message }`
- Uses SDK: `deleteEmployee()`

## 5. Package Dependencies

**Key Dependencies:**
- `next`: 16.1.6
- `react`: 19.2.3
- `@lifeplus/insuretech-sdk`: file:../sdks/insuretech-typescript-sdk/
- `@bufbuild/protobuf`: 2.11.0
- `@tanstack/react-table`: 8.21.3
- `@radix-ui/*`: Various (avatar, dropdown-menu, slot)
- `tailwindcss`: 4.1.18
- `recharts`: 3.7.0
- `lucide-react`: 0.562.0

## 6. No Existing Bulk Upload or Document Generation UI

### Current Status:
- ✅ Document proto definitions exist (document service, generation, templates)
- ❌ No bulk-upload components found
- ❌ No bulk-upload API routes found
- ❌ No docgen/document-generation UI components
- ❌ No document client library

### Implications:
Document generation infrastructure exists at the backend/proto level but is **not yet integrated into the B2B portal UI**. This is a candidate area for new feature development.

## 7. Project Structure

```
b2b_portal/
├── app/
│   ├── api/
│   │   ├── employees/        [GET, POST, PATCH, DELETE]
│   │   ├── departments/      [GET, POST, PATCH, DELETE]
│   │   ├── organisations/    [GET, POST, PATCH operations]
│   │   ├── purchase-orders/  [GET, POST operations]
│   │   └── auth/             [login, logout, session]
│   ├── employees/            [page.tsx - wrapper]
│   └── [other pages]
├── components/
│   ├── modals/               [add-employee-modal, add-department-modal, etc.]
│   ├── dashboard/
│   │   ├── employees/        [employees-table, data-table, columns]
│   │   ├── departments/
│   │   ├── organisations/
│   │   └── [other sections]
│   └── ui/                   [generic components]
├── src/
│   ├── lib/
│   │   ├── sdk/              [b2b-sdk-client.ts, api-helpers.ts]
│   │   ├── clients/          [employee-client, department-client, etc.]
│   │   ├── types/            [b2b.ts, employee-form.ts]
│   │   ├── auth/             [session management]
│   │   └── proto-generated/  [protobuf definitions]
│   └── hooks/                [useEmployeeForm, useToast, useCrudList]
└── package.json
```

## 8. Key Integration Points for Future Document Generation

### Potential Integration Areas:

1. **Employee Module:**
   - Add "Generate Certificate" button in employee actions
   - Generate policy certificate/schedule on employee onboarding
   - Download generated documents via file URL

2. **Form Enhancement:**
   - Add document template selection to add-employee-modal
   - Pass documentTemplateId to backend on employee creation

3. **New API Routes Needed:**
   - `POST /api/documents/generate` - trigger generation
   - `GET /api/documents/generations` - list generations
   - `GET /api/documents/generations/[id]` - get generation status
   - `GET /api/documents/templates` - list available templates

4. **New Components Needed:**
   - Document generation modal
   - Document download component
   - Document list/history component
   - Document template selector

5. **Client Library Needed:**
   - `src/lib/clients/document-client.ts` with methods for:
     - `generate(entityType, entityId, templateId, data)`
     - `listGenerations(filters)`
     - `getGeneration(id)`
     - `getTemplates()`

## 9. Authentication & Authorization

- **Session-based:** Cookie authentication (no API keys for portal users)
- **Roles:** BUSINESS_ADMIN, FINANCE_MANAGER, HR_MANAGER, SYSTEM_ADMIN, B2B_ORG_ADMIN
- **Org Isolation:**
  - B2B admins locked to their organisation
  - Super admins can view all organisations
- **Header Forwarding:** `x-portal`, `x-business-id`, `x-user-id`, `x-tenant-id`
- **CSRF Protection:** Token extraction and forwarding from cookies

## 10. Hooks Available

- `useEmployeeForm` - Employee create/edit form state
- `useOrganisationForm` - Organisation create/edit form state
- `useCrudList` - Generic CRUD list operations
- `useToast` - Toast notification management
- `useToast()` returns `{ toast, showToast }` for showing success/error messages

---

**Document Generated:** 2024
**Last Explored:** Comprehensive frontend exploration
