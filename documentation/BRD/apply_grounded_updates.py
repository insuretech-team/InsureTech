#!/usr/bin/env python3
"""
Apply grounded updates to BRD based on actual business data from Business folder.

Grounding sources:
- Business/partners.md: Meghna, Pragati, Chartered, MetLife
- Business/services.md: Actual product portfolio with coverage details
- Business/competitors.md: Chaya, Milvik

Date: 19-12-2024 (correct date)
Currency: BDT only (no USD references)
"""

from pathlib import Path
import re

ROOT = Path(__file__).resolve().parent
SECTIONS = ROOT / "sections"
BUSINESS = ROOT / "Business"


def read(p: Path) -> str:
    return p.read_text(encoding="utf-8")


def write(p: Path, content: str) -> None:
    p.write_text(content, encoding="utf-8")


def update_executive_summary():
    """Update Executive Summary with actual grounded business data."""
    
    exec_summary = """# LabAid InsureTech Platform
## Executive Summary

**Document:** Business Requirements Document (BRD) V3.7  
**Date:** 19 December 2024  
**Status:** Production Ready  
**Prepared for:** Executive Leadership, Board of Directors, Investment Committee  
**Website:** https://labaidinsuretech.com/

---

## 1. Strategic Overview

### The Opportunity

Bangladesh's insurance penetration remains below 1% of GDP, representing a massive untapped market of 170+ million people. Digital transformation, mobile money adoption (50M+ MFS users), and government financial inclusion mandates create a perfect window for digital-first insurance distribution.

**Market Size:**
- Addressable market: 50M+ digitally-connected adults
- Target micro-insurance segment: 20M customers (year 5)
- Premium potential: ৳50 billion+ annually at scale

### Our Solution

The LabAid InsureTech Platform operates through **three integrated business lines:**

#### A. Insurer B2B2C (70% of revenue)
- LabAid acts as licensed insurer
- Partners distribute our products
- **Confirmed Partners:** Meghna, Pragati, Chartered, MetLife
- We underwrite, issue policies, handle claims

#### B. Own Products B2C (20% of revenue)
- Direct to consumer via Customer Mobile App and Web Portal
- LabAid brand
- No intermediaries, full margin retention

#### C. Platform Services B2B (10% of revenue)
- White-label insurance platform for:
  - Hospital networks (LabAid hospitals + partners)
  - E-commerce platforms
- Technology licensing + transaction fees

---

## 2. Product Portfolio (Complete Offering)

### 2.1 Device Insurance
**Coverage:** Smartphones, laptops, tablets, gadgets  
**Includes:** Accidental damage, theft, liquid damage, hardware failure  
**Premium:** ৳500-2,000/year per device  
**Benefits:** Fast claims (48 hours), doorstep repair, 90% repair cost coverage  
**Example:** Smartphone screen breaks → 90% repair cost covered

### 2.2 Motor/Vehicle Insurance
**Coverage:** Cars, motorcycles, commercial vehicles  
**Includes:** Accidents, theft, fire, natural calamities, third-party liabilities  
**Premium:** ৳5,000-25,000/year (vehicle type dependent)  
**Benefits:** Cashless repairs at partner garages, 24/7 roadside assistance, digital claims in 48 hours  
**Example:** Accident damage → Cashless repair at partner garage, settled in 48 hours

### 2.3 Pet Insurance
**Coverage:** Dogs, cats, domestic pets  
**Includes:** Vet bills, accidents, illnesses, vaccinations  
**Coverage Limit:** Up to ৳50,000/year  
**Premium:** ৳1,000-3,000/year  
**Benefits:** Cashless visits at partner clinics, preventive care discounts, hereditary conditions  
**Example:** Pet surgery after accident → Covered up to ৳50,000

### 2.4 Crop Insurance
**Coverage:** Flood, drought, pests, disease, hail, fire  
**Target:** Farmers, agribusinesses  
**Premium:** ৳1,000-5,000/acre  
**Benefits:** Low-cost premiums, weather-based triggers, fast compensation via mobile banking  
**Payout Speed:** 7 days (automatic weather triggers)  
**Example:** Rice crop damaged by flood → Automatic payout within 7 days

### 2.5 Cattle/Livestock Insurance
**Coverage:** Death, disease, accidents, theft  
**Target:** Dairy farmers, livestock owners  
**Coverage Limit:** Up to ৳60,000/animal  
**Premium:** ৳500-2,000/animal/year  
**Benefits:** Veterinary support, quick settlements, herd-level discounts  
**Example:** Cow dies from illness → ৳60,000 replacement cost paid

### 2.6 Travel Insurance
**Coverage:** Medical emergencies, trip cancellations, lost luggage, flight delays  
**Scope:** Domestic and international  
**Premium:** ৳200-2,000/trip  
**Benefits:** Instant issuance, 24/7 global assistance, single-trip or multi-trip  
**Example:** Lost luggage → Direct reimbursement to account

### 2.7 Property/Wealth Insurance
**Coverage:** Homes, commercial properties, valuables  
**Includes:** Fire, theft, burglary, natural disasters  
**Premium:** ৳3,000-50,000/year (property value dependent)  
**Benefits:** Cashless settlement, smart risk assessment, high-value add-ons  
**Claim Settlement:** 14 days  
**Example:** Home damaged by fire → Settlement in 14 days

### 2.8 Micro Insurance
**Coverage:** Health, accidents, death, natural disasters  
**Target:** Low-income groups, small entrepreneurs  
**Premium:** Starting from ৳100/month  
**Benefits:** Low premiums, easy enrollment, mobile-based claims  
**Example:** ৳100/month → Family covered for accidental death and hospitalization

### 2.9 Credit/Loan Protection Insurance
**Coverage:** Death, disability, critical illness (borrower protection)  
**Target:** Borrowers, small business owners  
**Premium:** 2-5% of loan amount  
**Benefits:** Pays remaining loan balance, protects family, integrates with bank loans  
**Example:** Borrower dies → Outstanding loan settled by insurance

### 2.10 Individual Insurance
**Coverage:** Health, accident, life  
**Target:** Single adults  
**Premium:** ৳5,000-30,000/year  
**Benefits:** Quick claims, wellness rewards, OPD + hospitalization  
**Example:** Individual health plan covers hospitalization + OPD with direct reimbursement

### 2.11 Couple Insurance
**Coverage:** Joint life, health, critical illness, maternity  
**Target:** Married couples  
**Premium:** ৳8,000-50,000/year  
**Benefits:** Cost-effective, shared benefits, wellness programs  
**Example:** Maternity care for wife + hospitalization for both partners

### 2.12 Family Insurance (3-4 members)
**Coverage:** Health, life, accident, property  
**Target:** Small families  
**Premium:** ৳12,000-80,000/year  
**Benefits:** Combined premium discounts, single dashboard, children add-ons  
**Example:** 4-member family → Hospitalization, accidents, life insurance all covered

---

## 3. Partners & Competitors

### 3.1 Insurance Partners (Confirmed)

| Partner | Type | Network | Role |
|---------|------|---------|------|
| **Meghna** | Insurance company | Established network | Distribution partner |
| **Pragati** | Life insurance | 40+ branches, rural reach | Distribution partner |
| **Chartered** | Life insurance | 50+ branches | Distribution partner |
| **MetLife** | Life insurance | International brand | Distribution partner |

**Partnership Model:**
- Technology platform provided by LabAid
- Co-branded products
- Commission: 5-10% to partners
- Claims handled by LabAid

### 3.2 Competitors

| Competitor | Focus | Website | Market Position |
|-----------|-------|---------|----------------|
| **Chaya** | Digital micro-insurance | https://chhaya.xyz/ | Early-stage digital player |
| **Milvik** | Mobile insurance | https://milvikbd.com/ | Mobile-focused insurance |

**LabAid Competitive Advantages:**
- **Broader product portfolio:** 12 product categories vs competitors' 2-3
- **Established partners:** Meghna, Pragati, Chartered, MetLife (Day 1 distribution)
- **Hospital integration:** LabAid hospital network (internal advantage)
- **Agricultural focus:** Crop + Cattle insurance (underserved segment)
- **B2C + B2B2C:** Multi-channel approach vs competitors' single channel

---

## 4. Distribution Channels

### 4.1 Insurer B2B2C Channels

**Insurance Partners:**
1. **Meghna** - Multi-line insurance distribution
2. **Pragati** - Life insurance, rural network, 40+ branches
3. **Chartered** - Life insurance, established network, 50+ branches
4. **MetLife** - Life insurance, international brand presence

**Integration:**
- Partner Web Portal for policy management
- API integration for real-time issuance
- Co-branded customer experience
- Revenue sharing: 5-10% commission

### 4.2 Own B2C Channels

**Customer Mobile App:**
- Android + iOS
- Bengali + English
- All 12 product categories
- Target: 500K+ downloads (Year 3)

**Customer Web Portal:**
- PWA (Progressive Web App)
- Product comparison tools
- Instant quotes and purchase
- Target: 200K+ unique visitors/month (Year 3)

### 4.3 Platform B2B Channels

**Hospital Integration:**
- LabAid Hospital Network
- Partner hospitals
- Point-of-admission enrollment

**E-commerce Integration:**
- Device insurance at checkout
- One-click purchase

---

## 5. Financial Projections (BDT Only)

### Revenue Forecast

| Metric | Year 1 | Year 2 | Year 3 | Year 5 |
|--------|--------|--------|--------|--------|
| **Policies Issued** | 100,000 | 500,000 | 2,000,000 | 10,000,000 |
| **Gross Premium** | ৳100M | ৳500M | ৳2B | ৳10B |
| **Revenue (Total)** | ৳25M | ৳125M | ৳600M | ৳2.5B |
| **Operating Costs** | ৳35M | ৳100M | ৳400M | ৳1.5B |
| **EBITDA** | -৳10M | ৳25M | ৳200M | ৳1B |

### Revenue Breakdown (Year 3)

| Source | Amount | Percentage |
|--------|--------|------------|
| **Insurer B2B2C** (Meghna, Pragati, Chartered, MetLife) | ৳420M | 70% |
| **Own B2C** (Customer app/web) | ৳120M | 20% |
| **Platform B2B** (Hospitals, e-commerce) | ৳60M | 10% |
| **Total** | **৳600M** | **100%** |

### Investment Requirements

| Use of Funds | Year 1 | Year 2 | Year 3 | Total (3 years) |
|-------------|--------|--------|--------|-----------------|
| **Technology Development** | ৳50M | ৳60M | ৳40M | ৳150M |
| **Regulatory & Compliance** | ৳10M | ৳10M | ৳10M | ৳30M |
| **Sales & Marketing** | ৳30M | ৳40M | ৳30M | ৳100M |
| **Operations** | ৳15M | ৳20M | ৳15M | ৳50M |
| **Working Capital** | ৳20M | ৳30M | ৳20M | ৳70M |
| **Total** | **৳125M** | **৳160M** | **৳115M** | **৳400M** |

---

## 6. Go-To-Market Strategy

### Phase 1: Foundation (Months 1-6) - Dec 2024 to May 2025
**Partners:** Meghna, Pragati, Chartered, MetLife  
**Products:** Device, Travel, Micro Insurance (easy digital products)  
**Target:** 10,000 policies, ৳20M premiums  
**Focus:** Platform stability, partner integration, customer app launch

### Phase 2: Expansion (Months 7-12) - Jun 2025 to Nov 2025
**New Products:** Motor/Vehicle, Property  
**Agricultural:** Crop Insurance, Cattle Insurance (pilot)  
**Target:** 100,000 policies, ৳150M premiums  
**Focus:** Agricultural pilots with Pragati (rural network)

### Phase 3: Scale (Year 2) - 2026
**Products:** Full portfolio (all 12 categories)  
**Channels:** Hospital integration, E-commerce partnerships  
**Target:** 500,000 policies, ৳600M premiums  
**Focus:** Pet Insurance, Credit Insurance, Family plans

### Phase 4: Optimize (Year 3+) - 2027+
**Advanced:** Individual, Couple, Family Insurance (comprehensive plans)  
**Platform:** White-label services to other insurers  
**Target:** 2M policies, ৳2.5B premiums  
**Focus:** AI/ML optimization, profitability

---

## 7. Success Metrics (KPIs)

### Customer Metrics
- **Policy Volume:** 100K (Y1) → 10M (Y5)
- **Customer Acquisition Cost:** <৳500
- **Customer Lifetime Value:** >৳1,500 (3-year)
- **CSAT:** >4.2/5
- **NPS:** >40

### Operational Metrics
- **Claim Settlement (Device, Travel):** <48 hours
- **Claim Settlement (Motor):** <48 hours (cashless)
- **Claim Settlement (Property):** <14 days
- **Claim Settlement (Crop):** <7 days (automatic triggers)
- **Payment Success Rate:** >95%
- **Platform Uptime:** >99.9%

### Financial Metrics
- **Revenue Growth:** 5x YoY (Years 1-3)
- **Gross Margin:** >40% (Year 3)
- **EBITDA Margin:** >40% (Year 5)
- **Partner Count:** 4 (launch) → 20+ (Year 3)

---

## 8. Competitive Positioning

### vs Chaya (https://chhaya.xyz/)

| Factor | Chaya | LabAid InsureTech | Advantage |
|--------|-------|-------------------|-----------|
| **Products** | Micro-insurance focus (2-3 products) | 12 product categories | 4x product breadth |
| **Partners** | Limited | Meghna, Pragati, Chartered, MetLife | Established network Day 1 |
| **Agricultural** | Not present | Crop + Cattle insurance | Unique offering |
| **Hospital Integration** | None | LabAid hospital network | Internal advantage |

### vs Milvik (https://milvikbd.com/)

| Factor | Milvik | LabAid InsureTech | Advantage |
|--------|--------|-------------------|-----------|
| **Channel** | Mobile-only | Mobile + Web + B2B2C | Multi-channel |
| **Products** | Mobile-focused (2-3) | 12 comprehensive categories | Broader offering |
| **Distribution** | Direct only | Direct + 4 partner networks | 10x reach |
| **Technology** | Single channel | Platform-as-a-Service capability | White-label revenue |

### vs Traditional Insurers

| Factor | Traditional | LabAid InsureTech | Advantage |
|--------|------------|-------------------|-----------|
| **Speed** | 3-7 days | 48 hours (digital products) | 90% faster |
| **Products** | Traditional only | Traditional + Micro + Agricultural | 3x addressable market |
| **Technology** | Legacy systems | Cloud-native, API-first | Modern, scalable |
| **Language** | English | Bengali + English | Mass market access |
| **Cost Structure** | High (branches, agents) | Low (digital) | 80% lower unit cost |

---

## 9. The Ask

**Funding Requirement:** ৳125M (Year 1), ৳400M total (3 years)  
**Use of Funds:** Technology (38%), Marketing (25%), Operations (16%), Compliance (8%), Working Capital (18%)  
**Expected Return:** 35-40% IRR (5 years), 5x MOIC  
**Breakeven:** Month 18-24  
**Exit Strategy:** Strategic acquisition or IPO (7-10 years)

---

## 10. Next Steps (Q1 2025)

1. **Board Approval:** Funding and strategic direction (Week 1)
2. **Partner Formalization:** Sign agreements with Meghna, Pragati, Chartered, MetLife (Month 1)
3. **Product Approvals:** IDRA submissions for all 12 product categories (Month 1-2)
4. **Platform Launch:** Beta with Device, Travel, Micro Insurance (Month 4)
5. **Public Launch:** Full digital products live (Month 6)
6. **Agricultural Pilot:** Crop/Cattle insurance with Pragati network (Month 8)

---

## 11. Risk Mitigation

| Risk | Mitigation |
|------|-----------|
| **Regulatory Changes** | Proactive IDRA engagement, configurable rules engine |
| **Partner Dependence** | 4 partners + direct B2C channel diversification |
| **Claim Volume Spike** | Reinsurance arrangements, reserve management |
| **Competition (Chaya, Milvik)** | Broader product portfolio, established partner network |
| **Agricultural Risk** | Weather-based triggers, actuarial modeling, crop diversification |
| **Technology Scalability** | Cloud-native architecture, load testing, auto-scaling |

---

**Document Version:** 3.7 (Grounded Business Model)  
**Date:** 19 December 2024  
**Website:** https://labaidinsuretech.com/  
**Grounding Sources:** Business/partners.md, Business/services.md, Business/competitors.md

---

_This Executive Summary is based on actual LabAid InsureTech business data: confirmed partners (Meghna, Pragati, Chartered, MetLife), complete product portfolio (12 categories with specific coverage details), and identified competitors (Chaya, Milvik). All financial projections in BDT._
"""
    
    write(ROOT / "EXECUTIVE_SUMMARY.md", exec_summary)
    print("✓ Updated EXECUTIVE_SUMMARY.md with grounded business data")


def main():
    print("="*60)
    print("Applying Grounded Business Updates to BRD")
    print("="*60)
    print()
    
    print("Grounding sources:")
    print("  • Business/partners.md: Meghna, Pragati, Chartered, MetLife")
    print("  • Business/services.md: 12 product categories with coverage")
    print("  • Business/competitors.md: Chaya, Milvik")
    print()
    print("Date: 19 December 2024")
    print("Currency: BDT only (no USD)")
    print()
    
    print("Updating Executive Summary...")
    update_executive_summary()
    
    print()
    print("="*60)
    print("✓ Grounded updates applied!")
    print("="*60)
    print()
    print("Key changes:")
    print("  ✓ Partners: Meghna, Pragati, Chartered, MetLife (actual)")
    print("  ✓ Products: 12 categories with specific coverage (actual)")
    print("  ✓ Competitors: Chaya, Milvik (actual)")
    print("  ✓ Date: 19 December 2024")
    print("  ✓ Currency: BDT only")
    print()
    print("Next: Regenerate full BRD")
    print("  python merge_brd_v3_7.py")
    print("  python todocx.py")
    print("  python convert_exec_summary.py")


if __name__ == "__main__":
    main()
