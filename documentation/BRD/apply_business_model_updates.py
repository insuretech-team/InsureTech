#!/usr/bin/env python3
"""
Apply business model updates to BRD based on accurate LabAid InsureTech business model.

Changes:
1. Currency: USD → BDT
2. Business model: 3 lines (Insurer B2B2C, Own B2C, Platform B2B)
3. Partners: Chartered, Shanta, Pragati Life Insurance
4. Products: Traditional + Micro-insurance (Know Your Cattle, Know Your Crop)
5. Channels: Remove agent app, clarify web portals
"""

from pathlib import Path
import re

ROOT = Path(__file__).resolve().parent
SECTIONS = ROOT / "sections"


def read(p: Path) -> str:
    return p.read_text(encoding="utf-8")


def write(p: Path, content: str) -> None:
    p.write_text(content, encoding="utf-8")


def update_executive_summary():
    """Update Executive Summary with corrected business model."""
    
    exec_summary = """# LabAid InsureTech Platform
## Executive Summary

**Document:** Business Requirements Document (BRD) V3.7  
**Date:** January 2025  
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
- **Initial Partners:** Chartered Life Insurance, Shanta Life Insurance, Pragati Life Insurance
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

## 2. Business Value Proposition

### Revenue Model (All figures in BDT)

| Revenue Stream | Description | Year 1 | Year 3 | Year 5 |
|---------------|-------------|--------|--------|--------|
| **Insurer B2B2C** | Premium + Partner commission (70%) | ৳17.5M | ৳420M | ৳1.75B |
| **Own B2C** | Direct premium (20%) | ৳5M | ৳120M | ৳500M |
| **Platform B2B** | Tech fees + transactions (10%) | ৳2.5M | ৳60M | ৳250M |
| **Total Revenue** | | **৳25M** | **৳600M** | **৳2.5B** |

*Exchange reference: ৳110 = $1 USD*

### Cost Efficiency Gains

Compared to traditional insurance operations:
- **90% reduction** in policy issuance time (5 min vs 3 days)
- **80% reduction** in operational cost (digital vs paper)
- **70% reduction** in claims TAT (7 days vs 21 days)
- **50% reduction** in customer acquisition cost (via digital channels)

### Customer Impact

- **Instant coverage:** Policy issued within 5 minutes of payment
- **Affordable micro-insurance:** Starting from ৳200/year
- **Transparent pricing:** No hidden fees, instant quotes
- **Bengali language:** Mass-market accessibility
- **Mobile-first:** Works on low-end smartphones (2GB RAM)
- **Digital claims:** Submit from phone, track in real-time

---

## 3. Product Portfolio

### 3.1 Traditional Insurance Products (10,000-100,000 BDT range)

| Product | Premium Range | Target Market | Differentiation |
|---------|--------------|---------------|-----------------|
| **Health Insurance** | ৳10,000-50,000/year | Urban middle class | Cashless network, digital claims |
| **Motor Insurance** | ৳5,000-20,000/year | Vehicle owners | Instant issuance, photo claims |
| **Life Insurance** | ৳15,000-100,000/year | Families | Online purchase, no agent fees |
| **Travel Insurance** | ৳500-2,000/trip | Travelers | Real-time purchase, instant coverage |

### 3.2 Micro-Insurance Products (200-2,000 BDT range)

| Product | Premium Range | Target Market | Coverage |
|---------|--------------|---------------|----------|
| **Micro Health** | ৳200-2,000/year | Low-income families | Hospital cash, OPD coverage |
| **Micro Accident** | ৳300-1,000/year | Workers, students | Death, disability, medical |
| **Micro Life** | ৳500-2,000/year | Breadwinners | Term life, simple payout |
| **Know Your Cattle** | ৳500-2,000/animal | Farmers, rural | Livestock death, disease |
| **Know Your Crop** | ৳1,000-5,000/acre | Farmers | Crop damage, weather events |
| **Device Protection** | ৳500-1,500/device | E-commerce buyers | Accidental damage, theft |

**Micro-Insurance Advantages:**
- Low premiums (accessible to mass market)
- Simple coverage (easy to understand)
- Fast claims (<7 days)
- No medical exams (for low coverage amounts)
- Digital-first (mobile app enrollment)

---

## 4. Distribution Channels

### 4.1 Insurer B2B2C Channels

**Insurance Partners (Initial):**
1. **Chartered Life Insurance** - Established network, 50+ branches
2. **Shanta Life Insurance** - Digital-first approach, urban focus
3. **Pragati Life Insurance** - Rural reach, agent network

**Partner Integration:**
- Partner Web Portal for policy management
- API integration for real-time issuance
- Co-branded customer experience
- Revenue sharing: 5-10% commission to partners

### 4.2 Own B2C Channels

**Customer Mobile App:**
- Android + iOS
- Bengali + English
- Policy purchase, claims, renewals
- Target: 500K+ downloads (Year 3)

**Customer Web Portal:**
- PWA (Progressive Web App)
- Product comparison, quotes
- Full policy lifecycle management
- Target: 200K+ unique visitors/month (Year 3)

### 4.3 Platform B2B Channels

**Hospital Integration:**
- LabAid Hospital Network (5-10 hospitals)
- Partner hospitals (20+ hospitals by Year 2)
- Point-of-admission insurance enrollment
- Cashless claims integration

**E-commerce Integration:**
- Embed insurance at checkout
- Device protection, shipping insurance
- One-click purchase
- Target: 10+ e-commerce partners (Year 2)

---

## 5. Initial Partners (Confirmed)

### Insurance Distribution Partners

| Partner | Type | Network | Expected Volume (Year 1) |
|---------|------|---------|--------------------------|
| **Chartered Life Insurance** | Life insurer | 50+ branches, 500+ agents | 20,000 policies |
| **Shanta Life Insurance** | Life insurer | 30+ branches, digital channels | 15,000 policies |
| **Pragati Life Insurance** | Life insurer | 40+ branches, rural network | 15,000 policies |

**Partnership Model:**
- Technology platform provided by LabAid
- Co-branded products
- Shared underwriting (LabAid + partner)
- Commission: 5-10% to partners
- Claims handled by LabAid

### Hospital Partners

| Partner | Type | Integration |
|---------|------|-------------|
| **LabAid Hospital Network** | Internal | Full integration, cashless claims |
| **Partner Hospitals** | External | API integration, policy enrollment |

---

## 6. Technology Platform

### Architecture

- **Microservices:** Scalable, fault-tolerant
- **Multi-Tenant:** Partner data isolation
- **Cloud-Native:** Auto-scaling
- **API-First:** Easy partner integrations
- **Mobile-First:** Optimized for low-end devices

### Portal Ecosystem

| Portal | Users | Purpose |
|--------|-------|---------|
| **Customer Mobile App** | Customers | Purchase, manage, claim |
| **Customer Web Portal** | Customers | Product discovery, comparison |
| **Partner Web Portal** | Insurance partners, hospitals | Policy management, commissions |
| **Business Admin Portal** | LabAid ops team | Configuration, approvals, reports |
| **System Admin Portal** | IT team | Platform management, security |
| **Support Portal** | Call center | Ticketing, customer assistance |

**Note:** No dedicated agent mobile app. Partners use Partner Web Portal for all operations.

---

## 7. Regulatory & Compliance

### Licensing

- **Insurance Regulatory Authority (IDRA):** Licensed insurer status (required)
- **Bangladesh Financial Intelligence Unit (BFIU):** AML/CFT compliance
- **Data Protection:** Compliance with national regulations

### Compliance-First Design

| Requirement | Our Approach |
|------------|-------------|
| **IDRA Product Approval** | Configurable product engine (change rules without code) |
| **Policy Documentation** | Digital documents with QR verification, 7-year retention |
| **AML/CFT Monitoring** | Automated transaction monitoring, suspicious activity flagging |
| **Data Privacy** | Encryption at rest/transit, consent management, audit logs |
| **Audit Readiness** | Immutable audit trails, long-term retention, lawful access workflows |

---

## 8. Go-To-Market Strategy

### Phase 1: Foundation (Months 1-6)
**Goal:** Launch with anchor partners

- 3 insurance partners (Chartered, Shanta, Pragati)
- Focus: Micro health + traditional health
- Target: 10,000 policies, ৳20M premiums
- Channels: Partner portals + Customer app (B2C pilot)

### Phase 2: Expansion (Months 7-12)
**Goal:** Expand product portfolio

- Add motor insurance
- Launch agricultural insurance (Know Your Cattle, Know Your Crop)
- Onboard 5-10 hospital partners
- Target: 100,000 policies, ৳150M premiums

### Phase 3: Scale (Year 2)
**Goal:** Multi-channel distribution

- E-commerce partnerships (device insurance)
- Expand hospital network (20+ hospitals)
- Full B2C launch (customer app + web)
- Target: 500,000 policies, ৳600M premiums

### Phase 4: Optimize (Year 3+)
**Goal:** Platform services and profitability

- White-label platform for other insurers
- AI/ML features (fraud, underwriting)
- IoT/Telematics (future)
- Target: 2M policies, ৳2.5B premiums

---

## 9. Financial Projections (BDT)

### Revenue Forecast

| Metric | Year 1 | Year 2 | Year 3 | Year 5 |
|--------|--------|--------|--------|--------|
| **Policies Issued** | 100K | 500K | 2M | 10M |
| **Gross Premium** | ৳100M | ৳500M | ৳2B | ৳10B |
| **Commission Revenue (B2B2C)** | ৳17.5M | ৳87.5M | ৳420M | ৳1.75B |
| **Direct Revenue (B2C)** | ৳5M | ৳25M | ৳120M | ৳500M |
| **Platform Revenue (B2B)** | ৳2.5M | ৳12.5M | ৳60M | ৳250M |
| **Total Revenue** | ৳25M | ৳125M | ৳600M | ৳2.5B |
| **Operating Costs** | ৳35M | ৳100M | ৳400M | ৳1.5B |
| **EBITDA** | -৳10M | ৳25M | ৳200M | ৳1B |

### Investment Requirements (BDT)

| Use of Funds | Year 1 | Total (3 years) |
|-------------|--------|-----------------|
| **Technology Development** | ৳50M | ৳150M |
| **Regulatory & Compliance** | ৳10M | ৳30M |
| **Sales & Marketing** | ৳30M | ৳100M |
| **Operations** | ৳15M | ৳50M |
| **Working Capital** | ৳20M | ৳70M |
| **Total** | **৳125M** | **৳400M** |

*Equivalent to ~$1.14M USD (Year 1), ~$3.6M USD (3 years) at ৳110/$1*

### Return on Investment

- **Breakeven:** Month 18-24
- **IRR (5-year):** 35-40%
- **Exit Valuation (Year 5):** ৳15-20B (~$136-182M USD)

---

## 10. Success Metrics (KPIs)

### Customer Metrics
- **Policy Issuance Volume:** 100K (Y1) → 10M (Y5)
- **Customer Acquisition Cost (CAC):** <৳500
- **Customer Lifetime Value (LTV):** >৳1,500 (3-year)
- **CSAT:** >4.2/5
- **NPS:** >40

### Operational Metrics
- **Payment Success Rate:** >95%
- **Claims Settlement TAT:** <7 days
- **Platform Uptime:** >99.9%
- **Fraud Detection Rate:** >90%

### Financial Metrics
- **Revenue Growth:** 5x YoY (Years 1-3)
- **Gross Margin:** >40% (Year 3)
- **EBITDA Margin:** >40% (Year 5)

---

## 11. Competitive Advantage

### vs Traditional Insurers

| Factor | Traditional | LabAid Platform | Impact |
|--------|------------|----------------|--------|
| **Distribution** | Agents, branches | Digital + 3 business lines | 10x reach |
| **Products** | Traditional only | Traditional + micro-insurance | Addressable market 3x |
| **Customer Experience** | Paper, slow | Mobile-first, instant | NPS +30 points |
| **Cost Structure** | High fixed | Variable, scalable | 80% lower unit cost |

---

## 12. The Ask

**Funding Requirement:** ৳125M (Year 1), ৳400M total (3 years)  
**Use of Funds:** Technology, compliance, distribution, working capital  
**Expected Return:** 35-40% IRR, 5x MOIC (5 years)  
**Exit Strategy:** Strategic acquisition or IPO (7-10 years)

---

## 13. Next Steps

1. **Board Approval:** Funding and go-to-market strategy
2. **Partner Finalization:** Formalize agreements with Chartered, Shanta, Pragati
3. **Regulatory Submission:** IDRA product approvals
4. **Platform Development:** Complete integration with partner systems
5. **Pilot Launch:** Limited products, 3 partners (Month 4)
6. **Public Launch:** Full portfolio (Month 6)

---

**Website:** https://labaidinsuretech.com/  
**Document Version:** 3.7 (Updated for accurate business model)  
**Date:** January 2025  

---

_This Executive Summary reflects the actual LabAid InsureTech business model with three integrated business lines, confirmed partners (Chartered, Shanta, Pragati Life Insurance), and comprehensive product portfolio including micro-insurance and agricultural products._
"""
    
    write(ROOT / "EXECUTIVE_SUMMARY.md", exec_summary)
    print("✓ Updated EXECUTIVE_SUMMARY.md")


def update_business_context():
    """Update business context section with accurate model."""
    
    context = """# 2. Business Context

## 2.1 Market Overview

### Insurance Landscape in Bangladesh

- **Penetration:** <1% of GDP (among lowest in South Asia)
- **Total Market Size:** ~৳400 billion annual premiums
- **Growth Rate:** 15-20% annually (pre-COVID)
- **Challenges:** Limited distribution, low awareness, trust barriers, paper-heavy processes
- **Opportunities:** Mobile penetration >100%, MFS adoption (50M+ users), growing middle class, government financial inclusion push

### Target Segments

| Segment | Profile | Insurance Needs | Distribution Channel |
|---------|---------|-----------------|---------------------|
| **Urban Middle Class** | Salaried, smartphone users | Health, motor, life (traditional) | B2C (app/web) + Insurance partners |
| **Mass Market** | Lower income, MFS users | Micro health, accident, life | Insurance partners + Hospitals |
| **Rural/Agricultural** | Farmers, livestock owners | Cattle, crop, accident | Insurance partners + Field agents |
| **Digital Natives** | Online shoppers, young professionals | Device, travel, health | E-commerce integration + B2C |

## 2.2 Regulatory Environment

### Key Regulators

1. **Insurance Development and Regulatory Authority (IDRA)**
   - Regulates insurance companies and products
   - Requires product disclosure, policy document standards
   - Financial solvency reporting
   - Regular audits and inspections

2. **Bangladesh Financial Intelligence Unit (BFIU)**
   - Anti-Money Laundering (AML) and Countering the Financing of Terrorism (CFT)
   - Transaction monitoring requirements
   - Suspicious transaction reporting (STR/SAR)
   - Record retention obligations

### Licensing Requirements

**For LabAid InsureTech (Insurer License Required):**
- IDRA insurance company license
- Minimum paid-up capital: ৳400M (for life insurance)
- Board of Directors with insurance expertise
- Qualified actuary and compliance officer
- Adequate reserves and reinsurance arrangements

## 2.3 Business Model

### Three Integrated Business Lines

#### A. Insurer B2B2C (70% of Revenue)

**Model:** LabAid acts as licensed insurer, partners distribute

**Partners:**
- **Chartered Life Insurance** - Established network, 50+ branches
- **Shanta Life Insurance** - Digital-first, urban focus  
- **Pragati Life Insurance** - Rural reach, 40+ branches

**How It Works:**
1. Partners use LabAid technology platform
2. Partners sell LabAid-underwritten products
3. Co-branded customer experience
4. LabAid handles underwriting, policy issuance, claims
5. Partners receive 5-10% commission

**Target Products:**
- Traditional health, motor, life insurance
- Micro-insurance (health, accident, life)

#### B. Own Products B2C (20% of Revenue)

**Model:** Direct to consumer, LabAid brand

**Channels:**
- Customer Mobile App (Android, iOS)
- Customer Web Portal (PWA)

**How It Works:**
1. Customer discovers products on app/web
2. Instant quote and purchase
3. Digital policy document
4. Self-service claims and renewals
5. LabAid retains 100% of premium (higher margin)

**Target Products:**
- Micro-insurance (low complexity, digital-first)
- Travel insurance (instant purchase)
- Device protection (e-commerce checkout)

#### C. Platform Services B2B (10% of Revenue)

**Model:** White-label insurance technology platform

**Clients:**
- Hospital networks (LabAid hospitals + partners)
- E-commerce platforms
- Other insurers (future)

**How It Works:**
1. Client licenses LabAid platform
2. White-label or co-branded experience
3. LabAid provides technology + operations support
4. Revenue: licensing fees + transaction fees

**Target Products:**
- Hospital: Point-of-admission insurance, cashless claims
- E-commerce: Device protection, shipping insurance, warranty
- Future: Full platform-as-a-service for other insurers

## 2.4 Stakeholder Ecosystem

### Internal Stakeholders

| Stakeholder | Role | Primary Concerns |
|------------|------|------------------|
| **Board of Directors** | Governance, strategy | Regulatory compliance, profitability, growth |
| **CEO** | Overall strategy and execution | Revenue, partnerships, regulatory approvals |
| **CFO** | Financial management | Cash flow, reserves, unit economics |
| **CTO** | Technology platform | Uptime, scalability, security |
| **Head of Business** | Product and operations | Product-market fit, claims TAT, CSAT |
| **Compliance Officer** | Regulatory compliance | IDRA reporting, AML/CFT, audit readiness |

### External Stakeholders

| Stakeholder | Relationship | Expectations |
|------------|--------------|--------------|
| **Insurance Partners** | Distribution partners | Technology reliability, commission transparency, support |
| **Customers (B2C)** | End users | Affordable pricing, fast claims, Bengali support |
| **Hospital Partners** | Integration partners | Seamless enrollment, cashless claims, revenue share |
| **E-commerce Partners** | Integration partners | Easy integration, conversion boost, customer satisfaction |
| **Regulators (IDRA, BFIU)** | Oversight | Compliance, timely reporting, customer protection |
| **Reinsurers** | Risk partners | Accurate underwriting, fraud controls, financial stability |

## 2.5 Competitive Landscape

### Traditional Insurers (Life & General)

**Strengths:**
- Established brand, large agent networks
- Regulatory experience, financial reserves
- Reinsurance relationships

**Weaknesses:**
- Paper-heavy processes (3-7 days policy issuance)
- High operational costs (branches, agents)
- Limited digital presence
- Slow claims processing (30+ days)

**Examples:** Sadharan Bima, Jiban Bima, MetLife Bangladesh

### Emerging Digital Insurers

**Strengths:**
- Digital-first approach
- Mobile apps
- Faster issuance

**Weaknesses:**
- Limited product portfolio (1-2 products)
- Small distribution networks
- No micro-insurance focus
- English-only (limited mass market reach)

**Examples:** Local fintech/insurtech startups (early stage)

### LabAid Competitive Advantages

| Factor | Competitors | LabAid Platform | Differentiation |
|--------|------------|----------------|----------------|
| **Distribution** | Agents, branches | 3 business lines (B2B2C, B2C, B2B) | Multi-channel reach |
| **Products** | Traditional only | Traditional + micro-insurance + agricultural | Addressable market 3x |
| **Language** | English | Bengali + English | Mass market access |
| **Speed** | 3-7 days | 5 minutes | 90% faster |
| **Cost** | High (agents, branches) | Low (digital) | 80% lower |
| **Partners** | Limited | Chartered, Shanta, Pragati + hospitals | Established network Day 1 |
| **Technology** | Legacy systems | Cloud-native, microservices | Scalable, modern |

## 2.6 Market Entry Strategy

### Why Now?

✅ **Regulatory Window:** IDRA supportive of digital innovation  
✅ **Partnership Readiness:** Chartered, Shanta, Pragati eager to digitize  
✅ **Customer Adoption:** 50M+ MFS users, smartphone penetration growing  
✅ **Micro-Insurance Demand:** Government push for financial inclusion  
✅ **Agricultural Need:** 40% of population in agriculture, underinsured  
✅ **Competitive Gap:** No major player in digital micro-insurance  

### First-Mover Advantages

1. **Partner Lock-In:** Exclusive partnerships with top 3 life insurers
2. **Hospital Network:** LabAid hospital integration (internal advantage)
3. **Brand Credibility:** LabAid brand trusted in healthcare
4. **Regulatory Relationship:** Proactive engagement with IDRA
5. **Technology Lead:** 12-18 months ahead of traditional insurers

[[[PAGEBREAK]]]
"""
    
    write(SECTIONS / "02_business_context.md", context)
    print("✓ Updated 02_business_context.md")


def main():
    print("="*60)
    print("Applying Business Model Updates to BRD")
    print("="*60)
    print()
    
    print("Step 1: Updating Executive Summary...")
    update_executive_summary()
    
    print("Step 2: Updating Business Context...")
    update_business_context()
    
    print()
    print("="*60)
    print("✓ Core updates applied!")
    print("="*60)
    print()
    print("Next: Run full regeneration")
    print("  python generate_detailed_brd_v3_7.py")
    print("  python merge_brd_v3_7.py")
    print("  python todocx.py")
    print("  python convert_exec_summary.py")


if __name__ == "__main__":
    main()
