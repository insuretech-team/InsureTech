# InsureTech Admin Portal - TODO List

## Project Requirements (From User)

### Source Files to Review
- ✅ `C:\_DEV\GO\InsureTech\logo-header.svg` - Logo colors: Purple (#322F70) and Green (#40D200)
- ⏳ `C:\_DEV\GO\InsureTech\EXECUTIVE_SUMMARY.docx` - Business requirements
- ⏳ `C:\_DEV\GO\InsureTech\proto\**\*.proto` - Data structures for insurer/partner/discount/cashless

### Core Data Structure
- [ ] **TWO MAIN ENTITIES**:
  1. **Insurer** - The insurance company
  2. **Partner** - Service providers with two types of relationships:
     - **Discount** - Percentage-based discounts
     - **Cashless** - Direct billing without upfront payment

### Life Insurance Partners (4 Types)
- [ ] **Pharmacy** - Medicine discount and cashless
- [ ] **Hospitals** - Treatment discount and cashless claims
- [ ] **Doctors** - Consultation fee discount
- [ ] **Ambulance** - Emergency transport cashless

### Non-Life Insurance Partners (3 Types)  
- [ ] **Auto Repairs** - Car/vehicle repair discount and cashless
- [ ] **Laptop Repairs** - Computer repair discount and cashless
- [ ] **Mobile Repairs** - Phone repair discount and cashless

---

## UI/UX Requirements

### Theme & Design
- ✅ Professional theme based on logo colors
  - Primary: Purple (#322F70)
  - Accent: Green (#40D200)
- ✅ All shadcn-svelte components downloaded
- ✅ Login page with blocks
- ⏳ Dashboard blocks need improvement

### Navigation & Layout
- ✅ Sidebar with navigation
- ✅ Tabs for different sections
- ✅ Avatar profile component
- ⏳ Full-featured dashboard (needs more work)
- ⏳ Colorful badges (partially done, needs more variety)

---

## Functional Requirements

### 1. Discount Configuration UI
- [ ] **Discount percentage input** for each partner type
- [ ] **Min/Max discount range** validation
- [ ] **Partner-specific discount rates**:
  - Life partners: 5%-25%
  - Non-life partners: 8%-30%
- [ ] **Discount category dropdown**:
  - Service discount
  - Product discount
  - Bulk discount
- [ ] **Save/Update discount settings** per partner
- [ ] **Discount history/audit log**

### 2. Cashless Configuration UI
- [ ] **Cashless toggle switch** (Enable/Disable) per partner
- [ ] **Cashless limit amount** input field
- [ ] **Approval workflow settings**:
  - Auto-approval threshold
  - Manual approval required above threshold
- [ ] **Pre-authorization setup**:
  - Required documents checklist
  - Authorization validity period
- [ ] **Cashless partner list** with filters:
  - Active/Inactive status
  - Service area/location
  - Available services

### 3. Partner Management Dashboard
- [ ] **Partner overview cards** showing:
  - Total partners by type
  - Active discount programs
  - Cashless-enabled partners
  - Pending approvals
- [ ] **Partner details page** with:
  - Basic information
  - Discount configuration section
  - Cashless settings section
  - Service offerings
  - Performance metrics
- [ ] **Add/Edit partner forms**
- [ ] **Partner search and filters**
- [ ] **Partner status management** (Active/Inactive/Pending)

### 4. Insurer Configuration
- [ ] **Insurer profile page**
- [ ] **Default discount policies** setup
- [ ] **Cashless program settings**
- [ ] **Commission structure** configuration
- [ ] **Partner onboarding workflow**

### 5. Colorful Badges System
- [ ] **Status badges**:
  - 🟢 Active (Green)
  - 🟡 Pending (Yellow/Orange)
  - 🔴 Inactive/Suspended (Red)
  - 🔵 Verified (Blue)
- [ ] **Type badges**:
  - 💊 Pharmacy (Green)
  - 🏥 Hospital (Blue)
  - 👨‍⚕️ Doctor (Purple)
  - 🚑 Ambulance (Red)
  - 🚗 Auto Repair (Indigo)
  - 💻 Laptop Repair (Cyan)
  - 📱 Mobile Repair (Pink)
- [ ] **Discount badges**:
  - Show percentage with color gradient
  - High discount: Green
  - Medium discount: Blue
  - Low discount: Gray
- [ ] **Cashless badges**:
  - ✅ Enabled (Green with checkmark)
  - ❌ Disabled (Gray)

### 6. Reports & Analytics
- [ ] **Discount utilization report**
- [ ] **Cashless claims summary**
- [ ] **Partner performance dashboard**
- [ ] **Financial reports** (discounts given, cashless payments)

---

## Technical Implementation

### Pages to Build/Improve
- ✅ `/login` - Login page (DONE)
- ⏳ `/dashboard` - Main dashboard (NEEDS MAJOR IMPROVEMENT)
- ⏳ `/dashboard/partners/life` - Life partners (EXISTS but needs discount/cashless UI)
- ⏳ `/dashboard/partners/non-life` - Non-life partners (EXISTS but needs discount/cashless UI)
- [ ] `/dashboard/partners/[id]/edit` - Partner edit page with discount/cashless forms
- [ ] `/dashboard/partners/new` - Add new partner
- [ ] `/dashboard/discounts` - Discount management page
- [ ] `/dashboard/cashless` - Cashless program management
- [ ] `/dashboard/insurer` - Insurer settings
- [ ] `/dashboard/reports` - Reports and analytics

### Components to Create
- [ ] `DiscountConfigCard.svelte` - Discount percentage input with validation
- [ ] `CashlessToggle.svelte` - Toggle with limit settings
- [ ] `PartnerTypeIcon.svelte` - Icon component for each partner type
- [ ] `StatusBadge.svelte` - Reusable status badge with colors
- [ ] `DiscountBadge.svelte` - Show discount percentage with color
- [ ] `CashlessBadge.svelte` - Show cashless status
- [ ] `PartnerStatsCard.svelte` - Statistics card for partners
- [ ] `DiscountForm.svelte` - Form for setting discount rates
- [ ] `CashlessForm.svelte` - Form for cashless configuration

---

## Current Status

### ✅ Completed
1. Theme configuration with logo colors
2. All shadcn-svelte components installed
3. Login page
4. Basic dashboard layout with sidebar
5. Basic partner pages (life and non-life)
6. Some colorful badges

### ⚠️ Needs Major Work
1. **DISCOUNT CONFIGURATION UI** - Not implemented at all!
2. **CASHLESS CONFIGURATION UI** - Not implemented at all!
3. **Partner edit forms** - Missing
4. **Insurer management** - Not started
5. **Full-featured dashboard** - Too basic, needs real widgets
6. **More colorful badges** - Need variety and consistency

### 🔴 Critical Missing Features
- **NO discount input fields anywhere**
- **NO cashless toggle switches**
- **NO partner configuration forms**
- **NO way to actually MANAGE discounts and cashless**

---

## Next Steps (Priority Order)

1. **CREATE discount configuration UI** with percentage inputs
2. **CREATE cashless toggle and limit UI**
3. **ADD edit forms to partner pages**
4. **IMPROVE main dashboard** with real widgets and data
5. **ADD more colorful badges** throughout
6. **CREATE partner detail/edit pages**
7. **ADD insurer configuration page**
8. **CREATE reports section**

---

## Notes
- User is RIGHT - the dashboard is too "naive" and basic
- Focus on DISCOUNT and CASHLESS features first - that's the core requirement
- Make it ACTUALLY FUNCTIONAL, not just displaying static data
- Use the proto files to understand exact data structure
