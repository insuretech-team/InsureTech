"""
Phase 3: Add missing Traceability Matrix section and prepare for Word regeneration
"""

import re

# Read the current file
with open("SRS_V3_FINAL_DRAFT.md", "r", encoding="utf-8") as f:
    content = f.read()

# Add the missing Traceability Matrix & Change Control section
traceability_section = """
## 13. Traceability Matrix & Change Control

### 13.1 Requirements Traceability

| Business Objective | Related Functional Requirements | Success Metrics |
|--------------------|----------------------------------|-----------------|
| **Digital Onboarding: 40,000 policies by 2026** | FR-001 to FR-016 (User Management), FR-033 to FR-040 (Policy Purchase) | Monthly policy acquisition rate |
| **API Performance Optimization** | FR-107 to FR-118 (API Design), NFR-008 (gRPC <100ms) | API response time metrics |
| **Financial Transaction Integrity** | DomainServices (TigerBeetle integration), SEC-003 (PCI-DSS) | Transaction accuracy and speed |
| **Regulatory Compliance** | SEC-011 to SEC-020 (IDRA/AML/CFT) | Audit compliance score |
| **Partner Management Excellence** | FR-011 (Focal Person), FR-086 to FR-092 | Number of active partners |
| **Claims Efficiency** | FR-041 to FR-058, FR-133 to FR-137 | Average claim TAT |
| **VSA Architecture Implementation** | Section 4 (Architecture), CQRS/MediatR in C# services | Code maintainability score |
| **Proto-First Data Models** | Section 6 (Proto Buffers), proto/insuretech/v1/ structure | Cross-language compatibility |

### 13.2 Change Control Process

**Approval Hierarchy:**

1. **Dev** submits change request (Pull Request)
2. **Repository Admin** reviews code changes and architecture impact
3. **Database Admin** reviews data model impact (Proto changes, migrations)
4. **System Admin** reviews infrastructure impact (deployment, scaling)
5. **Business Admin** approves business impact and compliance
6. **Focal Person** approves partner-related changes
7. **CTO** final approval for major architectural changes

**Change Categories:**

| Category | Examples | Approval Level | Testing Required |
|----------|----------|----------------|------------------|
| **P1 - Critical** | Security patches, production bugs | System Admin + CTO | Hotfix testing |
| **P2 - High** | New features, API changes | Business Admin + Repository Admin | Full regression |
| **P3 - Medium** | UI/UX improvements, refactoring | Repository Admin | Feature testing |
| **P4 - Low** | Documentation, code comments | Dev | Peer review |

**Proto File Change Control:**
- All `.proto` file changes require Database Admin review
- Breaking changes to proto (field removal, type changes) require major version bump
- Adding optional fields is backward-compatible (minor version)
- Proto changes trigger code regeneration for all 4 languages (Go, C#, Node.js, Python)

---

"""

# Insert before "## 14. Appendices"
content = content.replace("## 12. Appendices", traceability_section + "## 14. Appendices")

# Also need to renumber the sections after insertion
content = content.replace("## 11. Acceptance Criteria", "## 12. Acceptance Criteria")
content = content.replace("## 10. Operational Requirements", "## 11. Operational Requirements")
content = content.replace("## 9. Security & Compliance", "## 10. Security & Compliance")
content = content.replace("## 8. Non-Functional Requirements", "## 9. Non-Functional Requirements")
content = content.replace("## 7. External Interface Requirements", "## 8. External Interface Requirements")
content = content.replace("## 6. Data Model", "## 7. Data Model")
content = content.replace("## 5. System Features", "## 6. System Features")

print("✅ Phase 3 complete!")
print("- Added Traceability Matrix & Change Control section (13)")
print("- Renumbered all sections properly")

# Write back
with open("SRS_V3_FINAL_DRAFT.md", "w", encoding="utf-8") as f:
    f.write(content)

print("\nNow regenerating Word document...")
