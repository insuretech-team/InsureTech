## 4. Sprint Planning & Timeline

### 4.1 Sprint Structure
- **Sprint Duration:** 2 weeks (12 working days)
- **Sprint Pattern:** Saturday to Thursday (6 days/week)
- **Sprint Ceremonies:**
  - Sprint Planning: Day 1 (4 hours)
  - Daily Standups: 15 minutes/day
  - Sprint Review: Last day (2 hours)
  - Sprint Retrospective: Last day (1 hour)
  - Backlog Refinement: Mid-sprint (2 hours)

---

### 4.2 Phase 1 Sprints (M1 Target: March 1, 2026)

**M1 SCOPE:** Business Admin Portal only, Policy purchase flow demo

#### Sprint 1: Dec 20, 2025 - Jan 2, 2026 (Foundation Sprint)
**Duration:** 12 working days (Includes Dec 25 holiday)
**Team:** 9 members (CTO, Mamoon, Sujon, Rumon, Nur, Sojol, QA, Sagor, React Dev)

| Team | Tasks | Hours | Deliverables |
|------|-------|-------|--------------|
| **CTO (40%)** | - API Gateway deployment (50% ready)<br>- User Service integration (100% ready) | 48 hrs | - Gateway running<br>- Auth APIs live |
| **Mamoon+Sujon** | - Payment Service setup (70% ready)<br>- Bkash sandbox testing | 134 hrs | - Payment test env |
| **React Dev** | - Project setup (Next.js)<br>- Design system integration<br>- Auth UI (reuse existing) | 96 hrs | - Login/Signup ready |
| **Rumon (UI/UX)** | - Design system finalization<br>- Business Admin portal wireframes | 96 hrs | - Design system<br>- Wireframes |
| **Nur+Sojol** | - Help React dev with setup<br>- Component library research | 192 hrs | - Support frontend |
| **Sagor (50%)** | - Cloud infrastructure setup<br>- Docker setup<br>- Database setup | 48 hrs | - Dev environment |
| **QA** | - Test plan development<br>- Test environment setup | 96 hrs | - Test plan |

**Sprint Goal:** Foundation ready, Auth working, designs done
**Hours Used:** 710 hrs (Dec capacity)

---

#### Sprint 2: Jan 3 - Jan 16, 2026
**Duration:** 12 working days
**Team:** 14 members (Full team, Project Manager joins Jan 1)

| Team | Tasks | Hours | Deliverables |
|------|-------|-------|--------------|
| **Delowar+C# Dev** | - Policy Service: Models & schema (C#)<br>- Policy CRUD operations start | 192 hrs | - Domain models<br>- Database setup |
| **CTO (40%)** | - Document Service integration (100% ready)<br>- Notification Service: Email setup | 77 hrs | - Doc API ready<br>- Email working |
| **Mamoon+Sujon** | - Payment Service: Premium payment flow<br>- Invoice generation | 120 hrs | - Payment API |
| **React Dev** | - Business Admin dashboard<br>- Product management UI (basic) | 96 hrs | - Dashboard live |
| **Rumon** | - High-fidelity mockups<br>- Policy management screens | 96 hrs | - Final designs |
| **Nur+Sojol** | - Component development support<br>- Responsive design help | 192 hrs | - UI components |
| **Sagor (50%)** | - CI/CD pipeline<br>- Database backup setup | 48 hrs | - Auto deploy |
| **QA** | - Test cases (auth, policy)<br>- API testing | 96 hrs | - Test coverage |
| **Project Manager** | - Sprint management<br>- Backlog refinement | 96 hrs | - Project tracking |

**Sprint Goal:** Policy service foundation, portal taking shape
**Hours Used:** 1,013 hrs

---

#### Sprint 3: Jan 17 - Jan 30, 2026
**Duration:** 12 working days
**Team:** 14 members (Delowar joins Jan 15, mid-sprint)

| Team | Tasks | Hours | Deliverables |
|------|-------|-------|--------------|
| **Delowar+C# Dev** | - Policy Service: CRUD complete<br>- Premium calculation engine<br>- Beneficiary management | 240 hrs | - Policy API complete |
| **CTO (40%)** | - Notification Service: SMS integration<br>- Policy events triggers | 77 hrs | - SMS working<br>- Email+SMS ready |
| **Mamoon+Sujon** | - Payment Service: Receipt generation<br>- Failed payment handling | 120 hrs | - Payment complete |
| **React Dev** | - Policy management full CRUD<br>- Policy search & filters | 96 hrs | - Policy pages done |
| **Rumon** | - Payment flow designs<br>- Notification UI mockups | 96 hrs | - Payment designs |
| **Nur+Sojol** | - Form components<br>- Table components<br>- Search UI help | 192 hrs | - Reusable components |
| **Sagor (50%)** | - SSL setup<br>- Basic monitoring | 48 hrs | - HTTPS enabled |
| **QA** | - Policy API testing<br>- Payment flow testing | 96 hrs | - Test reports |
| **Project Manager** | - Sprint coordination<br>- Risk tracking | 96 hrs | - Status reports |

**Sprint Goal:** Policy & Payment services complete
**Hours Used:** 1,061 hrs

---

#### Sprint 4: Jan 31 - Feb 13, 2026 (Includes Feb 4 holiday)
**Duration:** 11 working days (Shab-e-Barat holiday)
**Team:** 14 members

| Team | Tasks | Hours | Deliverables |
|------|-------|-------|--------------|
| **Delowar+C# Dev** | - Policy Service: Testing & bug fixes<br>- Policy status management<br>- Integration with Payment | 220 hrs | - Policy service stable |
| **CTO (40%)** | - Notification Service: Template management<br>- Event triggers complete | 71 hrs | - Notification ready |
| **Mamoon+Sujon** | - Payment integration testing<br>- Bkash production setup<br>- Payment audit logs | 110 hrs | - Payment production ready |
| **React Dev** | - Payment UI integration<br>- Document upload UI<br>- Notifications center | 88 hrs | - Payment flow UI done |
| **Rumon** | - Final design QA<br>- Asset library<br>- Design documentation | 88 hrs | - All designs final |
| **Nur+Sojol** | - Dashboard widgets<br>- Analytics charts<br>- Responsive polish | 176 hrs | - UI polish complete |
| **Sagor (50%)** | - Production environment setup<br>- Database backup | 44 hrs | - Prod environment |
| **QA** | - End-to-end testing<br>- Payment flow testing<br>- Bug verification | 88 hrs | - E2E test cases |
| **Project Manager** | - M1 readiness review<br>- Sprint planning | 88 hrs | - Go/No-go decision |

**Sprint Goal:** Integration & testing, production setup
**Hours Used:** 973 hrs

---

#### Sprint 5: Feb 14 - Feb 27, 2026 (Includes Feb 21 holiday)
**Duration:** 11 working days (Intl. Mother Language Day)
**Team:** 14 members

| Team | Tasks | Hours | Deliverables |
|------|-------|-------|--------------|
| **Delowar+C# Dev** | - Policy Service: Performance optimization<br>- API documentation<br>- Final polish | 220 hrs | - Policy service production ready |
| **CTO (40%)** | - All services integration testing<br>- Gateway performance tuning | 71 hrs | - All APIs integrated |
| **Mamoon+Sujon** | - Payment reconciliation<br>- Final payment testing | 110 hrs | - Payment certified |
| **React Dev** | - Final UI polish<br>- Error handling<br>- Loading states | 88 hrs | - UI production ready |
| **Rumon** | - Marketing materials<br>- Demo preparation support | 88 hrs | - Launch materials |
| **Nur+Sojol** | - Final responsive testing<br>- Browser compatibility<br>- Accessibility fixes | 176 hrs | - Cross-browser ready |
| **Sagor (50%)** | - Production deployment<br>- Monitoring setup<br>- SSL certificates | 44 hrs | - Deployed to prod |
| **QA** | - Full regression testing<br>- UAT support<br>- Sign-off testing | 88 hrs | - UAT complete |
| **Project Manager** | - M1 launch preparation<br>- Stakeholder coordination | 88 hrs | - Launch ready |

**Sprint Goal:** M1 launch readiness, final testing
**Hours Used:** 973 hrs

---

#### Sprint 5.5: Feb 28 - Mar 1, 2026 (Buffer & Launch)
**Duration:** 2 working days

| Team | Activities | Hours |
|------|------------|-------|
| **All Teams** | - Critical bug fixes<br>- Demo preparation<br>- Launch day coordination<br>- M1 deployment | 164 hrs |

**Sprint Goal:** M1 LAUNCH - March 1, 2026 (National Insurance Day)
**Hours Used:** 164 hrs

---

### M1 TOTAL SPRINT HOURS SUMMARY:
- Sprint 1 (Dec): 710 hrs
- Sprint 2 (Jan 3-16): 1,013 hrs  
- Sprint 3 (Jan 17-30): 1,061 hrs
- Sprint 4 (Jan 31-Feb 13): 973 hrs
- Sprint 5 (Feb 14-27): 973 hrs
- Sprint 5.5 (Feb 28-Mar 1): 164 hrs
- **TOTAL M1 USED:** 4,894 hrs ✅ (Exactly at capacity)

**M1 Components delivered:**
✅ User Service (Auth/AuthZ)
✅ Policy Service (Core CRUD + Premium calc)
✅ Payment Service (Bkash integration)
✅ Document Service (S3 storage)
✅ Notification Service (Email + SMS)
✅ Business Admin Portal
✅ DevOps Infrastructure (Basic)
✅ Testing & QA (Critical paths)

---

### 4.3 Phase 2 Sprints (M2 Target: April 14, 2026)

**M2 SCOPE:** Partner Portal, Mobile Apps (Customer+Agent), Claims Service

#### Sprint 6: Mar 2 - Mar 15, 2026
**Duration:** 12 working days (Less 5 days Eid = 7 working days)
**Team:** 14 members
**Note:** Eid-ul-Fitr Mar 21-25 impacts this sprint period

| Team | Tasks | Hours | Deliverables |
|------|-------|-------|--------------|
| **Delowar+C# Dev** | - Claims Service: Models & schema<br>- Claim submission flow | 168 hrs | - Claims foundation |
| **CTO (40%)** | - M1 production support<br>- API Gateway optimization | 54 hrs | - M1 stable |
| **Mamoon+Sujon** | - Partner Portal: Authentication<br>- Partner dashboard layout | 84 hrs | - Partner portal start |
| **React Dev** | - Partner Portal: Dashboard<br>- Partner profile UI | 67 hrs | - Partner UI |
| **Rumon** | - Partner Portal designs<br>- Mobile app wireframes | 67 hrs | - Designs ready |
| **Nur+Sojol** | - Customer Mobile App: Setup<br>- Authentication screens | 134 hrs | - Mobile project setup |
| **Sagor (50%)** | - M1 monitoring<br>- Performance tuning | 34 hrs | - M1 optimized |
| **QA** | - M1 production monitoring<br>- M2 test planning | 67 hrs | - Test plan M2 |
| **Project Manager** | - M2 sprint planning<br>- M1 post-launch review | 67 hrs | - M2 roadmap |

**Sprint Goal:** M2 foundation, M1 stabilization
**Hours Used:** 742 hrs (Reduced due to Eid holidays)

---

#### Sprint 7: Mar 16 - Mar 29, 2026 (Post-Eid, includes Mar 26 holiday)
**Duration:** 11 working days (Independence Day Mar 26)
**Team:** 14 members

| Team | Tasks | Hours | Deliverables |
|------|-------|-------|--------------|
| **Delowar+C# Dev** | - Claims Service: Approval workflow<br>- Claim assessment logic | 264 hrs | - Claims workflow ready |
| **CTO (40%)** | - Notification: Push notifications for mobile<br>- Mobile API support | 106 hrs | - Push notifications |
| **Mamoon+Sujon** | - Partner Portal: Commission tracking<br>- Partner analytics | 106 hrs | - Partner features |
| **React Dev** | - Partner Portal: Policy management<br>- Performance analytics | 88 hrs | - Partner portal 80% |
| **Rumon** | - Mobile app high-fidelity designs<br>- Agent app designs | 88 hrs | - Mobile designs |
| **Nur+Sojol** | - Customer App: Dashboard<br>- Customer App: Policy view<br>- Payment integration | 176 hrs | - Customer app core |
| **Sagor (50%)** | - Mobile app infrastructure<br>- Push notification setup | 44 hrs | - FCM configured |
| **QA** | - Claims API testing<br>- Partner portal testing | 88 hrs | - Test coverage |
| **Project Manager** | - Sprint coordination<br>- M2 progress tracking | 88 hrs | - Status reports |

**Sprint Goal:** Claims service advancing, mobile apps taking shape
**Hours Used:** 1,048 hrs

---

#### Sprint 8: Mar 30 - Apr 12, 2026
**Duration:** 11 working days (Apr 14 launch day not counted)
**Team:** 14 members

| Team | Tasks | Hours | Deliverables |
|------|-------|-------|--------------|
| **Delowar+C# Dev** | - Claims Service: Complete & test<br>- Claims settlement integration | 220 hrs | - Claims service ready |
| **CTO (40%)** | - All services integration<br>- Performance optimization | 85 hrs | - M2 services integrated |
| **Mamoon+Sujon** | - Partner Portal: Final features<br>- Commission calculation | 84 hrs | - Partner portal complete |
| **React Dev** | - Partner Portal: Polish & testing<br>- Final UI refinements | 67 hrs | - Partner portal production ready |
| **Rumon** | - M2 launch materials<br>- App store assets | 67 hrs | - Marketing ready |
| **Nur+Sojol** | - Customer App: Real API integration<br>- Customer App: bKash SDK integration<br>- App store submission | 134 hrs | - Customer apps submitted (NO Agent App) |
| **Sagor (50%)** | - M2 deployment<br>- Mobile app backend | 34 hrs | - M2 deployed |
| **QA** | - Full M2 regression<br>- Mobile app testing<br>- UAT sign-off | 67 hrs | - M2 certified |
| **Project Manager** | - M2 launch coordination<br>- Grand opening prep | 67 hrs | - Launch ready |

**Sprint Goal:** M2 LAUNCH READY
**Hours Used:** 825 hrs

**Note:** NO Agent Mobile App - agents use web portal on tablets per SRS requirements.

---

#### Buffer: Apr 13 - Apr 14, 2026
**Duration:** 1 working day (Apr 14 is launch/holiday)

| Team | Activities | Hours |
|------|------------|-------|
| **All Teams** | - Final bug fixes<br>- Launch day support<br>- Grand opening event | 112 hrs |

**Key Milestone:** M2 RELEASE - April 14, 2026 (Pohela Boishakh - Grand Launch)
**Hours Used:** 112 hrs

---

### M2 TOTAL SPRINT HOURS SUMMARY:
- Sprint 6 (Mar 2-15): 742 hrs
- Sprint 7 (Mar 16-29): 1,048 hrs
- Sprint 8 (Mar 30-Apr 12): 825 hrs
- Buffer (Apr 13-14): 112 hrs
- **TOTAL M2 USED:** 2,727 hrs

**M2 Components delivered:**
✅ Claims Service (Complete workflow)
✅ Partner Portal (Full featured)
✅ Customer Mobile App (Android + iOS) - Real API integration
❌ Agent Mobile App - NOT IN SCOPE (agents use web portal)
✅ Push Notifications
✅ Partner commission tracking
✅ Mobile payment integration (bKash real SDK)

---

### 4.4 Phase 3 Sprints (M3 Target: August 1, 2026)

**Sprint 9-15:** April 15 - August 1, 2026
**Total Duration:** 94 working days (7 two-week sprints + buffer)

#### High-Level Sprint Breakdown:

**Sprints 9-10 (Apr 15 - May 10):**
- Advanced Analytics & AI/ML foundation
- Machine learning model development
- Predictive analytics

**Sprints 11-12 (May 11 - Jun 7):**
- Complete integration ecosystem
- Advanced security features
- Performance optimization

**Sprints 13-14 (Jun 8 - Jul 5):**
- Additional features
- Platform enhancements
- Scalability improvements

**Sprint 15 (Jul 6 - Jul 19):**
- Final features
- Complete testing
- Documentation

**Buffer (Jul 20 - Aug 1):**
- Final polish
- M3 deployment preparation
- Production readiness

**Key Milestone:** M3 RELEASE - August 1, 2026

---

### 4.5 Sprint Capacity Planning

#### Per Sprint Team Capacity (2 weeks = 12 working days)

| Role | Count | Hours/Sprint | Utilization | Effective Hours |
|------|-------|--------------|-------------|-----------------|
| Backend Developers | 3 | 288 hrs | 85% | 245 hrs |
| Frontend Developer | 1 | 96 hrs | 85% | 82 hrs |
| Mobile Developers | 2 | 192 hrs | 85% | 163 hrs |
| DevOps Engineer | 1 | 96 hrs | 75% | 72 hrs |
| QA/Testers | 2 | 192 hrs | 80% | 154 hrs |
| UI/UX Designer | 1 | 96 hrs | 70% | 67 hrs |
| **Total per Sprint** | **10** | **960 hrs** | **82%** | **783 hrs** |

---

### 4.6 Sprint Velocity Tracking

| Sprint | Planned Story Points | Actual Story Points | Velocity | Notes |
|--------|---------------------|---------------------|----------|-------|
| Sprint 1 | TBD | TBD | TBD | Foundation |
| Sprint 2 | TBD | TBD | TBD | Core services |
| Sprint 3 | TBD | TBD | TBD | Service completion |
| Sprint 4 | TBD | TBD | TBD | Integration |
| Sprint 5 | TBD | TBD | TBD | M1 finalization |
| Sprint 6-8 | TBD | TBD | TBD | M2 development |
| Sprint 9-15 | TBD | TBD | TBD | M3 development |

**Note:** Velocity will be tracked and adjusted after Sprint 1 completion

---

### 4.7 Critical Path Analysis

#### M1 Critical Path:
```
Week 1-2:  Infrastructure + User Service (Foundation)
Week 3-4:  Policy Service + Document Service
Week 5-6:  Claim Service + Payment Service (70%)
Week 7-8:  Payment Service Complete + Customer Service + Integration
Week 9-10: Frontend Completion (with Mobile support) + Final Testing
```

#### Dependencies:
- **Week 1-2:** Must complete for all other work to proceed
- **Week 3-4:** Blocks claim and payment development
- **Week 5-6:** Critical for M1 functionality
- **Week 7-8:** Integration period, high risk
- **Week 9-10:** Final push, mobile dev support critical

#### Risk Mitigation:
- **Frontend bottleneck:** Mobile devs support from Week 8
- **Backend overload:** Prioritize critical path services
- **Integration issues:** Daily integration testing from Week 6
- **Testing delays:** Continuous testing throughout sprints

---

### 4.8 Sprint Ceremonies Schedule

#### Daily Standups (15 minutes)
- **Time:** 10:00 AM (Daily)
- **Participants:** All team members
- **Format:** What did you do? What will you do? Any blockers?

#### Sprint Planning (4 hours)
- **Day:** Sprint Day 1 (Saturday)
- **Time:** 9:00 AM - 1:00 PM
- **Participants:** All team members
- **Output:** Sprint backlog, task assignments

#### Backlog Refinement (2 hours)
- **Day:** Sprint Day 6 (Thursday)
- **Time:** 2:00 PM - 4:00 PM
- **Participants:** Product Owner, Tech Leads
- **Output:** Refined backlog for next sprint

#### Sprint Review (2 hours)
- **Day:** Last Sprint Day (Thursday)
- **Time:** 2:00 PM - 4:00 PM
- **Participants:** Team + Stakeholders
- **Output:** Demo, feedback, acceptance

#### Sprint Retrospective (1 hour)
- **Day:** Last Sprint Day (Thursday)
- **Time:** 4:00 PM - 5:00 PM
- **Participants:** Team members only
- **Output:** Improvement actions

---

### 4.9 Sprint Success Criteria

#### Definition of Ready (DoR):
- User story is clear and testable
- Acceptance criteria defined
- Dependencies identified
- Design mockups available (if UI)
- Technical approach discussed
- Story points estimated

#### Definition of Done (DoD):
- Code complete and reviewed
- Unit tests written and passing
- Integration tests passing
- API documentation updated
- Code merged to main branch
- Feature deployed to dev/staging
- QA testing complete
- No critical/high bugs
- Product Owner acceptance

---
