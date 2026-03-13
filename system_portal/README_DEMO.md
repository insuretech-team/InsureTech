# InsureTech Admin Portal - Demo Documentation

## Overview

This is a comprehensive demo of the InsureTech Admin Portal built with **SvelteKit**, **TypeScript**, and **shadcn-svelte** components. The demo uses proto-generated TypeScript types and modular demo data that can be easily replaced with backend API calls.

## Features Implemented

### 1. Proto-Based Type Generation
- ✅ Proto files located in `/proto/insuretech/`
- ✅ TypeScript types auto-generated using `buf` and `@bufbuild/protoc-gen-es`
- ✅ Generated types in `src/lib/generated/`
- ✅ Run `npm run generate:proto` to regenerate types from proto files

### 2. Demo Data Files (API-Ready)
All demo data files are in `src/lib/data_detailed/`:

#### Products (`products_demo.ts`)
- 8 insurance products across all categories
- Motor, Health, Travel, Home, Device, Agricultural, Life
- Includes riders, exclusions, pricing
- Uses proto-generated `Product`, `Rider`, `ProductCategory`, `ProductStatus` types

#### Policies (`policies_demo.ts`)
- 5 sample policies with different statuses
- Active, Grace Period, Expired, Pending Payment
- Includes nominees and policy riders
- Uses proto-generated `Policy`, `Nominee`, `PolicyRider`, `PolicyStatus` types

#### Claims (`claims_demo.ts`)
- 5 claims with various statuses and types
- Approved, Under Review, Pending Documents, Rejected, Settled
- Includes fraud checks, documents, approvals
- Uses proto-generated `Claim`, `ClaimDocument`, `ClaimApproval`, `FraudCheckResult` types

### 3. Demo Pages

#### Dashboard (`/dashboard`)
- Overview stats with real data from demo files
- Active products, policies, premium totals
- Pending claims counter
- Coverage statistics
- Partner overview tabs

#### Products Page (`/dashboard/products`)
- **List View**: All products with search and filters
- **Category Filter**: Motor, Health, Travel, Home, Device, Agricultural, Life
- **Status Filter**: Active, Inactive, Draft
- **Stats Cards**: Total products, active products, category breakdown
- **Table View**: Product code, name, category, premium, coverage range, tenure, riders
- **Detail View** (`/dashboard/products/[id]`): Full product details with riders and exclusions

#### Policies Page (`/dashboard/policies`)
- **List View**: All policies with search and filters
- **Status Filter**: Active, Pending Payment, Grace Period, Expired, Cancelled
- **Stats Cards**: Total policies, active count, premium sum, coverage sum
- **Table View**: Policy number, customer, premium, sum insured, dates, nominees
- **Nominee Management**: View nominees per policy

#### Claims Page (`/dashboard/claims`)
- **List View**: All claims with search and filters
- **Status Filter**: All, Submitted, Under Review, Pending Documents, Approved, Rejected, Settled
- **Stats Cards**: Total claims, pending review, approved, fraud flagged, total claimed
- **Table View**: Claim number, type, customer, amounts, incident date, fraud score
- **Fraud Detection**: Visual indicators for flagged claims
- **Days Pending**: Auto-calculated for pending claims

### 4. API-Ready Architecture

All demo data files include utility functions that are ready to be replaced with API calls:

```typescript
// Current (Demo):
export function getProductById(id: string): Product | undefined {
  return productsDemo.find((p) => p.productId === id);
}

// Future (Backend):
export async function getProductById(id: string): Promise<Product | undefined> {
  const response = await fetch(`/api/products/${id}`);
  return await response.json();
}
```

Simply replace the function bodies with API calls - the function signatures remain the same!

### 5. Type Safety

All data uses proto-generated types:
- ✅ No manual type definitions
- ✅ Enums from proto files (ProductCategory, PolicyStatus, ClaimStatus, etc.)
- ✅ BigInt for monetary values (stored in paisa)
- ✅ ISO 8601 date strings
- ✅ Helper functions for formatting (formatBDT)

### 6. Component Library

Using **shadcn-svelte** components:
- Card, Table, Button, Badge, Input
- Select, Tabs, Dropdown Menu
- Sidebar, Sheet, Dialog
- All components customizable via Tailwind CSS

## Getting Started

### Prerequisites
- Node.js 18+
- pnpm (or npm)
- buf CLI (for proto generation)

### Installation

```bash
cd system_portal
pnpm install
```

### Generate Proto Types

```bash
# From project root
buf generate

# Or from admin_portal directory
npm run generate:proto
```

### Run Development Server

```bash
pnpm dev
```

Navigate to `http://localhost:5173/dashboard`

## Project Structure

```
admin_portal/
├── src/
│   ├── lib/
│   │   ├── components/ui/     # shadcn-svelte components
│   │   ├── data_detailed/     # Demo data files (API-ready)
│   │   │   ├── products_demo.ts
│   │   │   ├── policies_demo.ts
│   │   │   ├── claims_demo.ts
│   │   │   ├── partners.ts
│   │   │   └── analytics.ts
│   │   ├── generated/         # Proto-generated TypeScript types
│   │   │   └── insuretech/
│   │   └── static/
│   └── routes/
│       └── dashboard/
│           ├── +page.svelte           # Dashboard overview
│           ├── products/
│           │   ├── +page.svelte       # Products list
│           │   └── [id]/+page.svelte  # Product detail
│           ├── policies/
│           │   └── +page.svelte       # Policies list
│           ├── claims/
│           │   └── +page.svelte       # Claims list
│           ├── partners/              # Partner pages
│           └── analytics/             # Analytics pages
└── proto/                     # Proto definitions (shared with backend)
```

## Data Format Examples

### Product
```typescript
new Product({
  productId: 'prod_001',
  productCode: 'HLT-001',
  productName: 'LabAid Health Guard',
  category: ProductCategory.PRODUCT_CATEGORY_HEALTH,
  basePremium: BigInt(500000), // 5000 BDT in paisa
  minSumInsured: BigInt(10000000), // 100,000 BDT
  maxSumInsured: BigInt(100000000), // 1,000,000 BDT
  status: ProductStatus.PRODUCT_STATUS_ACTIVE,
  availableRiders: [...]
})
```

### Policy
```typescript
new Policy({
  policyId: 'pol_001',
  policyNumber: 'LBT-2024-HLTH-000001',
  productId: 'prod_001',
  customerId: 'cust_001',
  status: PolicyStatus.POLICY_STATUS_ACTIVE,
  premiumAmount: BigInt(650000),
  sumInsured: BigInt(50000000),
  nominees: [...],
  riders: [...]
})
```

### Claim
```typescript
new Claim({
  claimId: 'clm_001',
  claimNumber: 'CLM-2024-HLTH-000001',
  policyId: 'pol_001',
  status: ClaimStatus.CLAIM_STATUS_APPROVED,
  type: ClaimType.CLAIM_TYPE_HEALTH_HOSPITALIZATION,
  claimedAmount: BigInt(12000000),
  approvedAmount: BigInt(11500000),
  documents: [...],
  approvals: [...],
  fraudCheck: {...}
})
```

## Backend Integration Checklist

When integrating with the backend:

1. **API Endpoints**: Create corresponding REST/gRPC endpoints for each function
2. **Replace Functions**: Update utility functions to call APIs
3. **Add Loading States**: Implement loading indicators
4. **Error Handling**: Add error boundaries and toast notifications
5. **Authentication**: Implement JWT/session-based auth
6. **Pagination**: Add server-side pagination for large datasets
7. **Real-time Updates**: Consider WebSocket for live updates

## Currency Format

All monetary values are stored in **paisa** (smallest unit):
- 100 paisa = 1 BDT
- Use `formatBDT(bigint)` helper to display in BDT format
- Example: `BigInt(500000)` → ৳5,000

## Next Steps

1. ✅ Proto types generated
2. ✅ Demo data created
3. ✅ Pages implemented (Products, Policies, Claims)
4. 🔄 Test pages in browser
5. ⏳ Add detail pages for policies and claims
6. ⏳ Implement search and advanced filters
7. ⏳ Add create/edit forms
8. ⏳ Connect to backend API

## Scripts

```bash
# Development
pnpm dev

# Build
pnpm build

# Preview production build
pnpm preview

# Generate proto types
pnpm generate:proto

# Type checking
pnpm check

# Linting & formatting
pnpm lint
pnpm format
```

## Notes

- All data is **modular** and stored in separate files
- **No hardcoded values** in pages - everything uses imported data
- **Type-safe** with proto-generated types
- **Currency** stored as BigInt in paisa, formatted for display
- **Ready for backend** - just replace function implementations
- **Scalable** - easy to add more demo data or connect to real APIs

## Demo URLs

- Dashboard: `http://localhost:5173/dashboard`
- Products: `http://localhost:5173/dashboard/products`
- Policies: `http://localhost:5173/dashboard/policies`
- Claims: `http://localhost:5173/dashboard/claims`
- Partners: `http://localhost:5173/dashboard/partners/life`
- Analytics: `http://localhost:5173/dashboard/analytics`
