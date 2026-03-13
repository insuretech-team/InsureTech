# LabAid InsureTech Platform - Detailed Project Plan
## Version 1.1 December 16, 2025

---

## Document Information

| Field | Value |
|-------|-------|
| **Project Name** | LabAid InsureTech Platform |
| **Document Type** | Detailed Project Plan |
| **Version** | 1.1 |
| **Date** | December 2024 (Planning Phase) |
| **Project Start** | December 20, 2025 |
| **Status** | Final - Aligned |
| **Owner** | Project Management Team |
| **Reviewers** | CEO, CTO, Business Admin, Technical Leads |

---

## Executive Summary

### Project Overview
The LabAid InsureTech Platform is a comprehensive digital insurance ecosystem designed to streamline insurance operations in Bangladesh. The platform enables:
- End-to-end policy management
- Digital claims processing
- Multi-channel partner integration
- Mobile-first customer experience
- IoT-based usage-based insurance (UBI)
- AI-powered automation and fraud detection

### Key Milestones
- **Project Start:** December 20, 2025
- **M1 (March 1, 2026):** Beta Launch - National Insurance Day Demo
- **M2 (April 14, 2026):** Grand Launch - Pohela Boishakh Public Release
- **M3 (August 1, 2026):** Complete Platform with IoT & AI Features

**Note:** Project spans December 2025 - August 2026 (8.5 months)

### Project Timeline

```
Dec 2025          Jan 2026          Feb 2026          Mar 2026          Apr 2026          May-Jul 2026      Aug 2026
|-----------------|-----------------|-----------------|-----------------|-----------------|-----------------|---------|
|                        M1 DEVELOPMENT (10 weeks)                      |  M2 DEV (5 wks)  |   M3 DEV (16 wks)  |
|                                                                       |                 |                    |
Dec 20                                                              Mar 1            Apr 14              Aug 1
Start                                                            Beta Launch     Grand Launch      Complete Platform
  |                                                                 ↓                 ↓                    ↓
  └─► Sprint 1 (Foundation)                                    National      Pohela Boishakh        Full Features
      Sprint 2 (Policy Service)                             Insurance Day    Public Release      IoT + AI Ready
      Sprint 3 (Payment Complete)                              Demo                              
      Sprint 4 (Integration)                                                                      
      Sprint 5 (Testing)                                                                          
      Sprint 5.5 (Launch)                                                                         
                                                                   └─► Sprint 6-8 (Mobile Apps)                    
                                                                       Claims Service                               
                                                                       Partner Portal                              
                                                                                    └─► Sprint 9-15 (IoT/AI)
                                                                                        10 More Portals
                                                                                        Analytics
                                                                                        Performance
```

### Capacity Utilization Over Time

```
Capacity
(hrs)
8000 │
     │                                                                         █████████████████████████░░░░░░░░░░
7000 │                                                                         ██████████████████████████░░░░░░░░░░
     │                                                                         ██████████████████████████░░░░░░░░░░
6000 │                                                                         ██████████████████████████░░░░░░░░░░
     │                                                                         ██████████████████████████░░░░░░░░░░
5000 │ █████████████████████████████░░░░░░░░░░░                                ██████████████████████████░░░░░░░░░░
     │ █████████████████████████████░░░░░░░░░░░                                ██████████████████████████░░░░░░░░░░
4000 │ █████████████████████████████░░░░░░░░░░░                                ██████████████████████████░░░░░░░░░░
     │ █████████████████████████████░░░░░░░░░░░                                ██████████████████████████░░░░░░░░░░
3000 │ █████████████████████████████░░░░░░░░░░░   ████████████████████░░░░░    ██████████████████████████░░░░░░░░░░
     │ █████████████████████████████░░░░░░░░░░░   ████████████████████░░░░░    ██████████████████████████░░░░░░░░░░
2000 │ █████████████████████████████░░░░░░░░░░░   ████████████████████░░░░░    ██████████████████████████░░░░░░░░░░
     │ █████████████████████████████░░░░░░░░░░░   ████████████████████░░░░░    ██████████████████████████░░░░░░░░░░
1000 │ █████████████████████████████░░░░░░░░░░░   ████████████████████░░░░░    ██████████████████████████░░░░░░░░░░
     │ █████████████████████████████░░░░░░░░░░░   ████████████████████░░░░░    ██████████████████████████░░░░░░░░░░
   0 └─────────────────────────────────────────────────────────────────────────────────────────────────────────────
           M1 (10 weeks)                         M2 (5 weeks)                       M3 (15 weeks)
        Dec 20 - Mar 1                        Mar 2 - Apr 14                     Apr 15 - Aug 1
     2,740 hrs / 4,040 (68%) ✅           1,641 hrs / 2,256 (73%) ✅         4,863 hrs / 6,768 (72%) ✅

     █ = Hours Used       ░ = Available Buffer    ALL PHASES COMFORTABLE!
```

### Project Capacity Summary (FINAL - REALISTIC CAPACITY - CORRECTED)
| Phase | Available Hours | Required Hours | Utilization | Status |
|-------|----------------|----------------|-------------|--------|
| M1 (Beta) | 4,040 hrs | 2,740 hrs | 68% | ✅ **COMFORTABLE** |
| M2 (Launch) | 2,256 hrs | 1,641 hrs | 73% | ✅ **COMFORTABLE** |
| M3 (Complete) | 6,768 hrs | 4,863 hrs | 72% | ✅ Comfortable |
| **TOTAL** | **13,064 hrs** | **9,244 hrs** | **71%** | ✅ Achievable |

**Capacity Adjustments Made:**
- M1: +48 hrs (Delowar calculation corrected from 173 to 221 hrs)
- Total: +48 hrs (13,016 → 13,064 hrs)

### Critical Success Factors
1. ✅ **Existing Code Reuse:** ~1,420 hours saved from proven services (Auth, Storage, Payment, IoT Broker)
2. ✅ **Mock-Based Parallel Development:** All teams work simultaneously - Backend, Frontend, Mobile
3. ✅ **M1 Scope:** Customer Mobile App (Android + iOS) + 2 web portals (Business Admin + Partner) with mocks (69% utilization)
4. ✅ **M2 Integration:** Real API integration + Claims backend (73% utilization - comfortable)
5. ✅ **Team Allocation:** CTO 50%, Mamoon 40%, Sagor 50%, Python 50%, Sujon 60%, Rumon 60%, React 70%, QA 60%, Mobile/Backend 90%
6. ✅ **No Agent Mobile App:** Per SRS, agents use web portal on tablets - saves 580 hrs
7. ✅ **No Overtime Needed:** All phases comfortable (69-73% utilization)

---

## Table of Contents

### 1. Team & Technology Stack

- 1.1 Team Structure & Roles
- 1.2 Technology Stack Overview
- 1.3 Architecture & Services

### 2. Project Timeline & Capacity

- 2.1 Available Working Hours by Phase
- 2.2 Capacity Summary Across All Phases
- 2.3 Utilization Notes
- 2.4 Holiday Impact Analysis
- 2.5 Critical Path Resources

### 3. Effort Estimation by Component

- 3.1 Estimation Methodology
- 3.2 M1 Services - Core Platform
  - 3.2.1 User Service (Auth/AuthZ)
  - 3.2.2 Policy Service
  - 3.2.3 Claims Service (Moved to M2)
  - 3.2.4 Payment Service
  - 3.2.5 Document Service
  - 3.2.6 Notification Service
  - 3.2.7 Ticketing Service (Moved to M2)
  - 3.2.8 Web Admin Portals
  - 3.2.9 Mobile Apps (Moved to M2)
  - 3.2.10 Infrastructure & DevOps
  - 3.2.11 QA & Testing
  - 3.2.12 UI/UX Design
- 3.3 M1 Summary - Total Effort
- 3.4 M2 Services - Desirable Features
- 3.5 M3 Services - Should Have & Future
- 3.6 Overall Project Summary

### 4. Sprint Planning & Timeline

- 4.1 Sprint Structure
- 4.2 Phase 1 Sprints (M1 Target: March 1, 2026)
  - Sprint 1: Dec 20 - Jan 2 (Foundation)
  - Sprint 2: Jan 3 - 16 (Policy Service Start)
  - Sprint 3: Jan 17 - 30 (Policy & Payment Complete)
  - Sprint 4: Jan 31 - Feb 13 (Integration)
  - Sprint 5: Feb 14 - 27 (Testing & Polish)
  - Sprint 5.5: Feb 28 - Mar 1 (Launch)
- 4.3 Phase 2 Sprints (M2 Target: April 14, 2026)
  - Sprint 6: Mar 2 - 15 (Claims & Mobile Start)
  - Sprint 7: Mar 16 - 29 (Mobile Development)
  - Sprint 8: Mar 30 - Apr 12 (Final Integration)
  - Buffer: Apr 13 - 14 (Grand Launch)
- 4.4 Phase 3 Sprints (M3 Target: August 1, 2026)
- 4.5 Sprint Capacity Planning
- 4.6 Sprint Velocity Tracking
- 4.7 Critical Path Analysis
- 4.8 Sprint Ceremonies Schedule
- 4.9 Sprint Success Criteria

### 5. RACI Matrix

- 5.1 Project Governance
- 5.2 Service Development RACI
- 5.3 Sprint Activities RACI
- 5.4 Deployment & Operations RACI

### 6. Resource Reassignment Strategy
**File:** `06_ResourceReassignment.md`
- 6.1 Mobile Developer Reassignment (Week 8)
- 6.2 Cross-functional Support Strategy
- 6.3 Capacity Optimization

### 7. Risks & Mitigation

- 7.1 Technical Risks
- 7.2 Schedule Risks
- 7.3 Resource Risks
- 7.4 Business Risks
- 7.5 Mitigation Strategies

### 8. Requirements Distribution by Milestone

- 8.1 M1 Requirements (103 FRs)
- 8.2 M2 Requirements (92 FRs)
- 8.3 M3 Requirements (81 FRs)
- 8.4 Summary & Trade-offs

### 9. Person-Wise Responsibility

- 9.1 Individual Assignments
- 9.2 Accountability Matrix
- 9.3 Backup & Coverage

---

## Document Change Log

| Version | Date | Changes | Author |
|---------|------|---------|--------|
| 1.0 | Dec 16, 2025 | Initial aligned version - Fixed M1 effort estimation, aligned all planning docs |Farukhannan|
| 1.1 | Dec 2025 | Draft version with capacity issues ,grapgh, enhancement|AI Engine|

---

## Approval Sign-off

| Role | Name | Signature | Date |
|------|------|-----------|------|
| Project Sponsor | TBD | | |
| CTO | TBD | | |
| Business Admin | TBD | | |
| Project Manager | TBD | | |

---

[[[PAGEBREAK]]]
