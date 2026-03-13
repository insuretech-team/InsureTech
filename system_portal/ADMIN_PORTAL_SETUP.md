# InsureTech Admin Portal - Setup Complete

## Overview
Professional admin portal for LabAid InsureTech with a modern UI based on the company logo colors and comprehensive partner management features.

## Theme Colors (Based on Logo)
- **Primary (Purple)**: `#322F70` - Navigation, buttons, and primary actions
- **Accent (Green)**: `#40D200` - Success states, highlights, and active indicators
- **Additional Colors**:
  - Success: Green `#40D200`
  - Warning: Orange
  - Info: Blue
  - Destructive: Red

## Features Implemented

### ✅ 1. Dashboard Layout
- **Full-featured sidebar navigation** with collapsible sections
- **Mobile-responsive design** with slide-out menu
- **Header** with search, notifications, and user profile dropdown
- **Professional navigation structure** organized by:
  - Overview (Dashboard, Analytics)
  - Insurance (Products, Policies, Claims)
  - Partners (Life, Non-Life, Agents)
  - Finance (Payments, Commissions)
  - System (Users, Settings)

### ✅ 2. Authentication Pages
- **Login page** (`/login`) with:
  - Split-screen design (form + branding)
  - Email/password fields with show/hide password
  - Remember me checkbox
  - Forgot password link
  - Beautiful gradient background showcasing features

### ✅ 3. Main Dashboard (`/dashboard`)
- **Statistics cards** showing:
  - Total Policies
  - Active Partners
  - Revenue
  - Active Claims
- **Tabbed partner overview** (Life vs Non-Life)
- **Colorful partner type cards** with icons:
  - Life: Hospitals, Pharmacies, Doctors, Ambulances
  - Non-Life: Auto Repair, Laptop Repair, Mobile Repair
- **Cashless & Discount information panels**
- **Recent policies table** with status badges

### ✅ 4. Life Partners Management (`/dashboard/partners/life`)
Four tabs for different partner types:

#### Hospitals Tab
- Partner ID, Name, Location
- Services count badge
- Discount percentage (color-coded)
- Cashless status (green/gray badges)
- Status badge (active/pending)
- Edit actions

#### Pharmacies Tab
- Coverage area
- Number of outlets
- Discount rates (5-10%)
- Cashless enabled status

#### Doctors Tab
- Specialty badges
- Hospital affiliation
- Consultation discounts (10-18%)
- Individual doctor profiles

#### Ambulances Tab
- Service type (ALS/BLS)
- Vehicle count
- Cashless emergency services
- 24/7 availability

### ✅ 5. Non-Life Partners Management (`/dashboard/partners/non-life`)
Three tabs for repair services:

#### Auto Repair Tab
- Service capabilities (brands supported)
- Location with map pin icons
- Discount rates (15-25%)
- Star ratings
- Cashless repair status

#### Laptop Repair Tab
- Supported brands (Dell, HP, Lenovo, Apple)
- Technical service capabilities
- Discount rates (10-18%)
- Ratings and reviews

#### Mobile Repair Tab
- Brand support (Samsung, Apple, etc.)
- Screen, battery, software repairs
- Discount rates (8-15%)
- Quick turnaround badges

### ✅ 6. Discount & Cashless Configuration
Both Life and Non-Life partner pages include:
- **Cashless coverage summaries** with percentages
- **Average discount rate displays** with visual progress bars
- **Partner statistics** by category
- **Status indicators** for active/pending partners

### ✅ 7. Colorful Badges & Status Indicators
Implemented throughout the application:
- **Status Badges**: Active (primary), Pending (secondary), Suspended (warning)
- **Discount Badges**: Color-coded by partner type
  - Blue for hospitals
  - Green for pharmacies
  - Purple for doctors
  - Indigo for auto repair
  - Cyan for laptop repair
  - Pink for mobile repair
- **Cashless Badges**: Green for enabled, Gray for disabled
- **Rating Indicators**: Star ratings with numeric values
- **Service Badges**: Outlined badges for categories

## Shadcn-Svelte Components Installed
- ✅ Button
- ✅ Card (with Header, Content, Description, Title)
- ✅ Input
- ✅ Label
- ✅ Badge
- ✅ Avatar
- ✅ Dropdown Menu
- ✅ Separator
- ✅ Sidebar
- ✅ Tabs
- ✅ Sheet (mobile drawer)
- ✅ Scroll Area
- ✅ Table
- ✅ Select
- ✅ Checkbox
- ✅ Tooltip
- ✅ Skeleton

## Tech Stack
- **Framework**: SvelteKit 2.x with Svelte 5
- **Styling**: TailwindCSS 4.x
- **UI Components**: shadcn-svelte (New York style)
- **Icons**: Lucide Svelte
- **TypeScript**: Full type safety

## File Structure
```
admin_portal/
├── src/
│   ├── lib/
│   │   ├── components/
│   │   │   ├── ui/              # Shadcn components
│   │   │   └── dashboard-layout.svelte
│   │   └── utils.ts
│   ├── routes/
│   │   ├── +page.svelte         # Landing/home
│   │   ├── +layout.svelte       # Root layout
│   │   ├── login/
│   │   │   └── +page.svelte     # Login page
│   │   └── dashboard/
│   │       ├── +layout.svelte   # Dashboard wrapper
│   │       ├── +page.svelte     # Main dashboard
│   │       └── partners/
│   │           ├── life/
│   │           │   └── +page.svelte
│   │           └── non-life/
│   │               └── +page.svelte
│   └── app.css                  # Theme configuration
├── static/
│   └── logo.svg                 # Company logo
└── tailwind.config.js           # Tailwind configuration
```

## Running the Application

### Development Server
```bash
cd C:\_DEV\GO\InsureTech\admin_portal
pnpm run dev
```

Access at: `http://localhost:5173`

### Build for Production
```bash
pnpm run build
pnpm run preview
```

## Key Routes

| Route | Description |
|-------|-------------|
| `/` | Home/Landing page |
| `/login` | Authentication page |
| `/dashboard` | Main dashboard with overview |
| `/dashboard/partners/life` | Life insurance partners (Hospitals, Pharmacies, Doctors, Ambulances) |
| `/dashboard/partners/non-life` | Non-life partners (Auto, Laptop, Mobile repair) |

## Partner Types

### Life Insurance Partners (4 categories)
1. **Hospitals** - Healthcare facilities with 24 services, 15-20% discounts
2. **Pharmacies** - Medicine retailers with 5-10% discounts, nationwide coverage
3. **Doctors** - Individual practitioners with 10-18% consultation discounts
4. **Ambulances** - Emergency services with cashless availability

### Non-Life Insurance Partners (3 categories)
1. **Auto Repair** - Vehicle maintenance, 15-25% discounts, all brands
2. **Laptop Repair** - Computer services, 10-18% discounts, brand-specific
3. **Mobile Repair** - Smartphone services, 8-15% discounts, quick turnaround

## Discount & Cashless Features
- **Variable discount rates** by partner type and category
- **Cashless claim processing** for 79-84% of partners
- **Visual indicators** showing coverage percentages
- **Partner ratings** for quality assurance
- **Service badges** for quick identification

## Design Highlights
- **Professional color scheme** matching company branding
- **Responsive design** - works on desktop, tablet, and mobile
- **Accessible UI** with proper contrast ratios
- **Consistent spacing** using Tailwind utilities
- **Icon integration** with contextual meanings
- **Data visualization** with progress bars and statistics

## Next Steps (Future Enhancements)
1. Add partner edit/create forms with validation
2. Implement discount configuration module
3. Add cashless approval workflows
4. Create analytics/reporting pages
5. Integrate with backend API
6. Add user role management
7. Implement real-time notifications
8. Add partner performance metrics

## Notes
- Logo uses colors: Purple (#322F70) and Green (#40D200)
- Theme is configured for both light and dark modes
- All components follow shadcn-svelte conventions
- Ready for backend integration via SvelteKit load functions
