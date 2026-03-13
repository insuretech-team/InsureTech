# Person-Wise Responsibility Matrix

## Overview

This document provides a detailed breakdown of responsibilities, tasks, and workload for each team member across all three milestones (M1, M2, M3).

---

## Team Member Details & Availability

| # | Name | Role | Languages/Tech | Availability | Hours/Week | Total M1 Hours | Total M2 Hours | Total M3 Hours |
|---|------|------|---------------|--------------|------------|----------------|----------------|----------------|
| 1 | **CTO** | Technical Lead | Go, Architecture | 50% coding | 24 hrs | 240 hrs | 115 hrs | 374 hrs |
| 2 | **Mamoon** | Senior Full Stack | Node.js, React | 40% | 19 hrs | 190 hrs | 92 hrs | 299 hrs |
| 3 | **Sujon Ahmed** | Full Stack Dev | Node.js, React | 60% (realistic) | 48 hrs | 288 hrs | 144 hrs | 432 hrs |
| 4 | **Rumon** | UI/UX Designer | Figma, Design | 60% (realistic) | 48 hrs | 288 hrs | 144 hrs | 432 hrs |
| 5 | **Nur Hossain** | Android Developer | Kotlin, Java | 90% (realistic) | 48 hrs | 432 hrs | 216 hrs | 648 hrs |
| 6 | **Sojol Ahmed** | iOS Developer | Swift | 90% (realistic) | 48 hrs | 432 hrs | 216 hrs | 648 hrs |
| 7 | **QA** | QA/Test Engineer | Testing | 60% (realistic) | 48 hrs | 288 hrs | 144 hrs | 432 hrs |
| 8 | **Sagor** | DevOps Engineer | Docker, K8s, Go | 50% | 24 hrs | 240 hrs | 115 hrs | 374 hrs |
| 9 | **React Dev** | Frontend Developer | React, Next.js | 70% (realistic) | 48 hrs | 336 hrs | 168 hrs | 504 hrs |
| 10 | **Project Manager** | PM | Agile, Scrum | 90% (realistic) | 48 hrs | 346 hrs | 216 hrs | 648 hrs |
| 11 | **Mr. Delowar** | C# Lead | C# .NET, gRPC | 90% (from Jan 15) | 48 hrs | 221 hrs | 230 hrs | 749 hrs |
| 12 | **C# Developer** | C# Developer | C# .NET, gRPC | 90% (realistic) | 48 hrs | 346 hrs | 216 hrs | 648 hrs |
| 13 | **Python Dev 1** | Python Developer | Python, FastAPI | 50% (M3 AI focus) | 24 hrs | 192 hrs | 115 hrs | 374 hrs |
| 14 | **Python Dev 2** | Python Developer | Python, FastAPI | 50% (M3 AI focus) | 24 hrs | 192 hrs | 115 hrs | 374 hrs |

**Note on Utilization:** This table shows AVAILABLE hours per week. Actual utilization varies by role:
- CTO, Mamoon, Sagor, Python Devs: 40-50% (have other commitments)
- Sujon, Rumon, QA, React Dev: 60-70% (realistic with meetings, blockers)
- Mobile Devs, Backend Devs, PM: 90% (high utilization on focused work)
- See 03_TeamCapacity.md for REALISTIC capacity calculations used in planning.

---

## 1. CTO - Technical Lead & Gateway Developer

**Availability:** 50% coding time (50% management)
**Primary Tech:** Go, Architecture, API Gateway, Kafka, IoT Broker
**Total Project Hours:** 729 hours (240 M1 + 115 M2 + 374 M3)

### M1 Responsibilities (221 hours - Joins Jan 15, 2026)

| Service/Component | Hours | Description | Status |
|-------------------|-------|-------------|--------|
| **API Gateway** | 40 hrs | Deploy and configure 50% existing code | 50% Ready |
| **Authentication Service** | 8 hrs | Deploy existing service, insurance customization | ✅ 100% Ready |
| **Authorization Service** | 8 hrs | Deploy existing service, role setup | ✅ 100% Ready |
| **DBManager Service** | 8 hrs | Deploy existing service, connection setup | ✅ 100% Ready |
| **Storage Manager Service** | 8 hrs | Deploy existing service, S3 integration | ✅ 100% Ready |
| **Notification Service (Kafka)** | 120 hrs | Kafka setup, email/SMS integration, templates | New Development |
| **IoT Broker Setup (20%)** | 16 hrs | Basic MQTT setup for future | 80% Ready |
| **Technical Leadership** | 24 hrs | Architecture reviews, code reviews, mentoring | Ongoing |
| **Integration Support** | 8 hrs | Service integration troubleshooting | As needed |
| **TOTAL M1** | **240 hrs** | | |

**Primary Deliverables (M1):**
- ✅ API Gateway deployed and functional
- ✅ All existing Go services deployed (Auth, Authz, DBManager, Storage)
- 🆕 Kafka-based Notification Service operational
- 📋 Technical architecture documentation

---

### M2 Responsibilities (115 hours)

| Service/Component | Hours | Description |
|-------------------|-------|-------------|
| **API Gateway Enhancements** | 30 hrs | Rate limiting, advanced routing, webhooks |
| **Notification Service v2** | 35 hrs | Advanced templates, multi-language, scheduling |
| **Integration Service Support** | 24 hrs | Support C# team with gRPC integration |
| **Technical Leadership** | 18 hrs | M2 architecture, performance optimization |
| **Production Support** | 8 hrs | M1 production issues, incident response |
| **TOTAL M2** | **115 hrs** | |

---

### M3 Responsibilities (374 hours)

| Service/Component | Hours | Description |
|-------------------|-------|-------------|
| **IoT Broker Complete** | 120 hrs | Full MQTT setup, device management, 10K devices |
| **IoT Telemetry Processing** | 96 hrs | Kafka Streams, real-time & batch processing |
| **IoT Device Management Portal** | 80 hrs | Device registry, monitoring, health tracking |
| **API Gateway v3** | 24 hrs | Final optimizations, WebSocket support |
| **Technical Leadership** | 40 hrs | M3 architecture, IoT infrastructure |
| **Documentation** | 14 hrs | Technical docs, API documentation |
| **TOTAL M3** | **374 hrs** | |

**Primary Deliverables (M3):**
- 🚀 Complete IoT infrastructure (10,000 devices supported)
- 📡 Real-time telemetry processing with Kafka
- 🔧 IoT device management system
- 📚 Complete technical documentation

---

## 2. Mamoon - Senior Full Stack Developer

**Availability:** 40% (has other commitments)
**Primary Tech:** Node.js, React, MongoDB, Payment Integration
**Total Project Hours:** 613 hours

### M1 Responsibilities (190 hours)

| Service/Component | Hours | Description | Status |
|-------------------|-------|-------------|--------|
| **Payment Service** | 120 hrs | Complete bKash integration (70% done), add Nagad | 70% Ready |
| **Ticketing Service** | 50 hrs | Basic ticket management, Node.js backend | New Development |
| **Payment Testing & Integration** | 12 hrs | End-to-end payment testing | Testing |
| **Code Reviews** | 8 hrs | Review Sujon's work | Ongoing |
| **TOTAL M1** | **190 hrs** | | |

**Primary Deliverables (M1):**
- ✅ bKash payment gateway 100% operational
- 🆕 Basic ticketing system for customer support
- 💳 Payment reconciliation and audit trail

---

### M2 Responsibilities (92 hours)

| Service/Component | Hours | Description |
|-------------------|-------|-------------|
| **Payment Service v2** | 36 hrs | Add Rocket, credit/debit cards, installment plans |
| **Ticketing Service v2** | 30 hrs | Advanced features, SLA tracking, escalation |
| **Commission Service Support** | 18 hrs | Support C# team with payment integrations |
| **Code Reviews** | 8 hrs | Node.js code reviews |
| **TOTAL M2** | **92 hrs** | |

---

### M3 Responsibilities (299 hours)

| Service/Component | Hours | Description |
|-------------------|-------|-------------|
| **Payment Service v3** | 80 hrs | Blockchain integration, advanced reconciliation |
| **Ticketing Service v3** | 60 hrs | AI integration, auto-routing, chatbot integration |
| **Integration API Support** | 80 hrs | Partner payment integrations, webhooks |
| **Full-Stack Features** | 60 hrs | Various enhancements across platform |
| **Documentation** | 19 hrs | Payment integration guides, API docs |
| **TOTAL M3** | **299 hrs** | |

---

## 3. Sujon Ahmed - Mid-Level Full Stack Developer

**Availability:** 100% (exhaust)
**Primary Tech:** Node.js, React, PostgreSQL
**Total Project Hours:** 1,541 hours

### M1 Responsibilities (480 hours)

| Service/Component | Hours | Description |
|-------------------|-------|-------------|
| **Ticketing Service (Backend)** | 100 hrs | Work with Mamoon on Node.js backend |
| **Partner Portal Support** | 80 hrs | Backend APIs for Partner Portal |
| **System Admin Portal Support** | 80 hrs | Backend APIs for System Admin Portal |
| **Database Schema Design** | 40 hrs | PostgreSQL schema for all services |
| **API Integration Layer** | 60 hrs | Connect frontend to backend services |
| **Testing & Bug Fixes** | 80 hrs | Backend testing, integration testing |
| **Code Reviews** | 20 hrs | Peer reviews |
| **Documentation** | 20 hrs | API documentation |
| **TOTAL M1** | **480 hrs** | |

---

### M2 Responsibilities (230 hours)

| Service/Component | Hours | Description |
|-------------------|-------|-------------|
| **Business Admin Portal Backend** | 75 hrs | APIs for Business Admin Portal |
| **Agent Portal Backend** | 60 hrs | APIs for Agent Portal |
| **Analytics Backend Support** | 45 hrs | Data aggregation APIs |
| **Testing & Bug Fixes** | 38 hrs | M2 feature testing |
| **Documentation** | 12 hrs | API updates |
| **TOTAL M2** | **230 hrs** | |

---

### M3 Responsibilities (749 hours)

| Service/Component | Hours | Description |
|-------------------|-------|-------------|
| **Remaining Portals Backend** | 200 hrs | APIs for 6 remaining portals |
| **Performance Optimization** | 120 hrs | Database query optimization, caching |
| **Advanced Features** | 200 hrs | Various backend enhancements |
| **Integration Support** | 100 hrs | EHR, third-party integrations |
| **Testing & Bug Fixes** | 100 hrs | M3 comprehensive testing |
| **Documentation** | 29 hrs | Complete backend documentation |
| **TOTAL M3** | **749 hrs** | |

---
