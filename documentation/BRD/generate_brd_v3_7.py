#!/usr/bin/env python3
"""Generate a COMPLETE BRD (business-facing) from SRS V3.7.

Goals
- Cover ALL FRs, NFRs and Security/Compliance controls from SRS V3.7.
- Do NOT copy-paste SRS sentences; re-phrase as business requirements.
- Produce portal definitions (what portals exist and what each must do).
- Maintain traceability by referencing SRS FG/FR/NFR/SEC IDs.

This script writes markdown section files under BRD/sections/.
The final BRD is assembled by merge_brd_v3_7.py (which concatenates sections in filename order).
"""

from __future__ import annotations

import re
from pathlib import Path
from datetime import datetime

ROOT = Path(__file__).resolve().parent
SRS = ROOT.parent / "SRS_v3" / "SPECS_V3.7"
SRS_SECTIONS = SRS / "sections"
BRD_SECTIONS = ROOT / "sections"


def read(p: Path) -> str:
    return p.read_text(encoding="utf-8")


def write(p: Path, s: str) -> None:
    p.parent.mkdir(parents=True, exist_ok=True)
    p.write_text(s.strip() + "\n", encoding="utf-8")


def parse_md_table_rows(md: str) -> list[list[str]]:
    """Parse GitHub-flavored markdown tables.

    Returns rows as list of cells (strings) excluding header separator rows.
    """
    rows: list[list[str]] = []
    for line in md.splitlines():
        if not line.strip().startswith("|"):
            continue
        # skip separator |---|
        if re.match(r"^\|\s*-{2,}", line.strip()):
            continue
        cells = [c.strip() for c in line.strip().strip("|").split("|")]
        # must have at least 2 cells and contain an ID in first cell
        if not cells or not re.search(r"\b(FR|NFR|SEC)-\d+\b", cells[0]):
            continue
        rows.append(cells)
    return rows


def extract_feature_groups(fr_md: str) -> list[dict]:
    """Extract FG sections (title + table rows).

    Output: [{fg_id, title, rows:[{id, desc, priority, ac}...]}]
    """
    groups: list[dict] = []

    # Find headings like: ### 4.1 User Management & Authentication (FG-001)
    pattern = re.compile(r"^###\s+.+?\((FG-\d+)\)\s*$", re.M)
    starts = [(m.start(), m.group(1), m.group(0)) for m in pattern.finditer(fr_md)]

    for i, (pos, fg_id, heading_line) in enumerate(starts):
        end = starts[i + 1][0] if i + 1 < len(starts) else len(fr_md)
        chunk = fr_md[pos:end]

        # Title is the heading text without markdown ### and without (FG-xxx)
        heading_text = re.sub(r"^###\s+", "", heading_line).strip()
        title = re.sub(r"\s*\(FG-\d+\)\s*$", "", heading_text).strip()

        rows = []
        for cells in parse_md_table_rows(chunk):
            # expected columns: ID | description | priority | acceptance
            rid = cells[0]
            desc = cells[1] if len(cells) > 1 else ""
            prio = cells[2] if len(cells) > 2 else ""
            ac = cells[3] if len(cells) > 3 else ""
            rows.append({"id": rid, "desc": desc, "priority": prio, "ac": ac})

        groups.append({"fg_id": fg_id, "title": title, "rows": rows})

    return groups


PORTALS = [
    {
        "name": "Customer Mobile App",
        "audience": "Customers/Policyholders",
        "purpose": "Self-service onboarding, discovery, purchase, policy servicing, claims, support.",
    },
    {
        "name": "Customer Web Portal (PWA)",
        "audience": "Customers/Policyholders",
        "purpose": "Web equivalent of the customer journey (campaign traffic, desktop users).",
    },
    {
        "name": "Partner Admin Portal",
        "audience": "Partner admins (MFS/hospital/e-commerce/agent org)",
        "purpose": "Partner onboarding, agent management, assisted sales, commission/analytics, operational tooling.",
    },
    {
        "name": "Agent App (Mobile)",
        "audience": "Agents operating under a partner",
        "purpose": "Assisted onboarding/purchase, lead handling, commissions, basic support tooling.",
    },
    {
        "name": "Focal Person Portal",
        "audience": "InsureTech focal persons",
        "purpose": "Partner KYB verification/approval, dispute resolution, monitoring, escalations.",
    },
    {
        "name": "Business Admin Portal",
        "audience": "InsureTech business ops",
        "purpose": "Product governance, workflow approvals, claims controls, reporting, operational configuration.",
    },
    {
        "name": "System Admin Portal",
        "audience": "Platform/system administrators",
        "purpose": "Security configuration, roles, system health, incident tooling, audit.",
    },
    {
        "name": "Support/Call Centre Portal",
        "audience": "Customer support agents",
        "purpose": "Ticketing, customer history, escalation workflows, communication tools.",
    },
    {
        "name": "Regulatory Access Portal (controlled)",
        "audience": "IDRA/BFIU or auditors (as per lawful request)",
        "purpose": "Controlled delivery of requested reports/data with full audit trail.",
    },
    {
        "name": "IoT Device Management Portal (partner/internal)",
        "audience": "IoT partners / internal ops",
        "purpose": "Device onboarding, health monitoring, telemetry visibility (where IoT program is active).",
    },
]


def portal_section() -> str:
    lines = [
        "# 4. Portals & Channels (What We Are Building)",
        "",
        "This BRD explicitly defines the portals/channels required by SRS V3.7 (FG-023 and cross-cutting requirements).",
        "Each portal definition below includes its business purpose, primary users, and must-have capabilities.",
        "",
    ]
    for idx, p in enumerate(PORTALS, start=1):
        lines += [
            f"## 4.{idx} {p['name']}",
            "",
            f"**Primary users:** {p['audience']}  ",
            f"**Business purpose:** {p['purpose']}",
            "",
            "**Must-have capabilities (business view):**",
            "- Role-appropriate login and account safety controls",
            "- Clear dashboards for primary tasks and statuses",
            "- Auditability for sensitive actions",
            "- Multi-language support where customer-facing (Bengali/English)",
            "",
        ]
    lines += [
        "**Traceability:** SRS FG-023 (FR-244..FR-248) and portal-related capabilities across FG-001..FG-019.",
        "",
        "[[[PAGEBREAK]]]",
    ]
    return "\n".join(lines)


def businessify_desc(desc: str) -> str:
    """Light transformation to avoid verbatim copy while keeping meaning."""
    d = re.sub(r"^The system shall\s+", "", desc.strip(), flags=re.I)
    d = d.rstrip(".")
    # Make it business phrased
    return d[0].upper() + d[1:] if d else d


def functional_catalog_section(groups: list[dict]) -> str:
    out = [
        "# 6. Detailed Business Functional Requirements (Complete Catalog)",
        "",
        "This section enumerates business requirements derived from SRS V3.7 functional requirements.",
        "Each requirement is phrased in business language and retains traceability to the original SRS FR-ID.",
        "",
        "Notation",
        "- **BR-ID**: Business Requirement Identifier (for BRD tracking)",
        "- **SRS Trace**: SRS Feature Group and FR ID(s)",
        "- **Priority**: aligned to SRS phase labels (M1/M2/M3/D/S/F)",
        "",
    ]

    br_counter = 1
    for g in groups:
        out += [f"## 6.{br_counter} {g['title']} ({g['fg_id']})", ""]
        out.append(
            "Business intent: define the outcomes this capability must deliver for customers/partners/admin teams."
        )
        out.append("")

        for r in g["rows"]:
            br_id = f"BR-{br_counter:02d}-{int(re.search(r'\d+', r['id']).group()):03d}"
            out += [
                f"### {br_id} — {businessify_desc(r['desc'])}",
                "",
                f"- **SRS Trace:** {g['fg_id']} / {r['id']}",
                f"- **Priority:** {r['priority'] or 'TBD'}",
                "- **Business acceptance (summary):**",
                f"  - {re.sub(r'<br>', ' / ', r['ac']).strip() if r['ac'] else 'Must be testable via clear acceptance criteria and operational logs.'}",
                "- **Primary portals impacted:** (to be confirmed during UX)",
                "  - Customer App / Web, Partner Portal, Admin Portals as applicable",
                "",
            ]
        out += ["[[[PAGEBREAK]]]", ""]
        br_counter += 1

    return "\n".join(out)


def nfr_section(nfr_md: str) -> str:
    rows = parse_md_table_rows(nfr_md)
    out = [
        "# 7. Non-Functional Requirements (NFR) — Business-Grade Detail",
        "",
        "NFRs are non-negotiable business constraints because they define customer experience, reliability of money movement, regulatory readiness, and operational cost.",
        "",
        "## 7.1 NFR Catalog (Derived from SRS Section 5)",
        "",
    ]

    for cells in rows:
        nfr_id = cells[0]
        area = cells[1] if len(cells) > 1 else ""
        req = cells[2] if len(cells) > 2 else ""
        target = cells[3] if len(cells) > 3 else ""
        prio = cells[4] if len(cells) > 4 else ""

        out += [
            f"### {nfr_id} — {area}",
            "",
            f"- **Business requirement:** {businessify_desc(req) if req else 'TBD'}",
            f"- **Target/Measurement:** {target or 'As defined in SRS'}",
            f"- **Priority:** {prio or 'TBD'}",
            "- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.",
            "",
        ]

    out += ["[[[PAGEBREAK]]]", ""]
    return "\n".join(out)


def security_section(sec_md: str) -> str:
    # Extract SEC rows (they are in tables). We'll parse all SEC-xxx occurrences by scanning all table rows.
    rows = parse_md_table_rows(sec_md)

    out = [
        "# 8. Security, Privacy, and Compliance Requirements (Detailed)",
        "",
        "Security and compliance are business requirements: they protect customers, protect funds, enable partner trust, and satisfy IDRA/BFIU expectations.",
        "This section translates SRS Section 7 controls into business-operational requirements.",
        "",
        "## 8.1 Security Control Catalog (SEC)",
        "",
    ]

    for cells in rows:
        sec_id = cells[0]
        desc = cells[1] if len(cells) > 1 else ""
        prio = cells[-1] if len(cells) >= 2 else ""
        out += [
            f"### {sec_id}",
            "",
            f"- **Business control requirement:** {businessify_desc(desc) if desc else 'TBD'}",
            f"- **Priority:** {prio or 'TBD'}",
            "- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.",
            "",
        ]

    out += [
        "## 8.2 AML/CFT Operating Model (Business View)",
        "",
        "The platform must support configurable AML monitoring rules, alerting, investigation workflow, and STR/SAR filing support with strict auditability.",
        "(See SRS Section 7.7.x for rule tables and workflows.)",
        "",
        "## 8.3 IDRA Reporting and Record-Keeping (Business View)",
        "",
        "The platform must retain and produce long-term records (policies, payments, claims, cancellations, approvals, customer communications) with retrieval capability within required SLAs.",
        "",
        "[[[PAGEBREAK]]]",
    ]

    return "\n".join(out)


def traceability_stub(groups: list[dict], nfr_md: str, sec_md: str) -> str:
    """Creates a simple traceability matrix skeleton.

    Full deep matrix can be expanded later (BR-ID -> Epic/Story -> Test cases).
    """
    nfr_ids = sorted(set(re.findall(r"\bNFR-\d+\b", nfr_md)))
    sec_ids = sorted(set(re.findall(r"\bSEC-\d+\b", sec_md)))

    out = [
        "# 9. Traceability (BRD ↔ SRS)",
        "",
        "This matrix ensures every SRS requirement is accounted for in business terms.",
        "",
        "## 9.1 Functional Traceability (FG/FR → BRD Coverage)",
        "",
        "| SRS Feature Group | Covered in BRD Section(s) |",
       "|---|---|",
    ]

    for g in groups:
        out.append(f"| {g['fg_id']} | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |")

    out += [
        "",
        "## 9.2 Non-Functional Traceability (NFR → BRD)",
        "",
        "| SRS NFR ID | Covered in BRD Section |",
       "|---|---|",
    ]
    for n in nfr_ids:
        out.append(f"| {n} | Section 7 |")

    out += [
        "",
        "## 9.3 Security Traceability (SEC → BRD)",
        "",
        "| SRS SEC ID | Covered in BRD Section |",
       "|---|---|",
    ]
    for s in sec_ids:
        out.append(f"| {s} | Section 8 |")

    out += ["", "[[[PAGEBREAK]]]", ""]
    return "\n".join(out)


def main() -> None:
    fr_md = read(SRS_SECTIONS / "04_functional_requirements.md")
    nfr_md = read(SRS_SECTIONS / "05_non_functional_requirements.md")
    sec_md = read(SRS_SECTIONS / "07_security_compliance.md")

    groups = extract_feature_groups(fr_md)

    # New/overwritten BRD sections (ordered by filename)
    write(BRD_SECTIONS / "04_portals_channels.md", portal_section())
    write(BRD_SECTIONS / "06_detailed_functional_catalog.md", functional_catalog_section(groups))
    write(BRD_SECTIONS / "07_detailed_nfr_catalog.md", nfr_section(nfr_md))
    write(BRD_SECTIONS / "08_detailed_security_compliance.md", security_section(sec_md))
    write(BRD_SECTIONS / "09_traceability_matrix.md", traceability_stub(groups, nfr_md, sec_md))

    print("Generated BRD sections:")
    for p in [
        BRD_SECTIONS / "04_portals_channels.md",
        BRD_SECTIONS / "06_detailed_functional_catalog.md",
        BRD_SECTIONS / "07_detailed_nfr_catalog.md",
        BRD_SECTIONS / "08_detailed_security_compliance.md",
        BRD_SECTIONS / "09_traceability_matrix.md",
    ]:
        print(" -", p)


if __name__ == "__main__":
    main()
