#!/usr/bin/env python3
"""
Lint SRS sections for plan- and team-agnostic rules.
- No milestone/timeline content
- No team ownership columns ('Team Owner')
- No in-body 'Document Control' or 'Document Approval & Sign-off' (final sections allowed at end of merged doc)
- Header Executive Summary must not contain phase/milestone text

Exit code 0 if clean, 1 if violations found.
"""
from pathlib import Path
import re
import sys

BASE = Path(r"G:\_0LifePlus\InsureTech\SRS_v3\SPECS_V3.7")
SECTIONS = BASE / "sections"

VIOLATIONS = []

# Patterns to flag
PATTERNS = [
    (re.compile(r"\b(M1|M2|M3|Phase\s*[DSF]|milestone|timeline|delivery|phase[d]?)(?![A-Za-z])", re.IGNORECASE), "Planning/milestone language"),
    (re.compile(r"\|\s*Team Owner\s*\|", re.IGNORECASE), "Team Owner column in table"),
    (re.compile(r"^\s*\*\*Document Control:\*\*", re.IGNORECASE|re.MULTILINE), "In-body Document Control block"),
    (re.compile(r"^\s*##?\s*Document Approval\s*&\s*Sign-off", re.IGNORECASE|re.MULTILINE), "In-body Document Approval & Sign-off"),
]

# Files to scan (exclude planning file if present)
files = sorted([p for p in SECTIONS.glob('*.md') if p.name != '13_project_planning.md'])

for fp in files:
    text = fp.read_text(encoding='utf-8')
    # Header-specific check: remove 'Phased Delivery' or milestone bullets in Executive Summary
    if fp.name == '00_header.md':
        if re.search(r"\*\*Phased Delivery:\*\*", text, re.IGNORECASE):
            VIOLATIONS.append((fp, "Executive Summary contains 'Phased Delivery' block"))
        if re.search(r"Milestone", text, re.IGNORECASE):
            VIOLATIONS.append((fp, "Executive Summary mentions Milestone; should be plan-agnostic"))
    # General patterns
    for pat, desc in PATTERNS:
        for m in pat.finditer(text):
            # Ignore occurrences inside code fences
            before = text[:m.start()]
            code_fences_open = before.count('```') % 2 == 1
            if code_fences_open:
                continue
            VIOLATIONS.append((fp, f"{desc}: '{m.group(0)[:60]}'"))

if VIOLATIONS:
    print("❌ Lint violations found:")
    for fp, msg in VIOLATIONS:
        print(f" - {fp.name}: {msg}")
    sys.exit(1)
else:
    print("✅ Lint passed: No planning/team content detected in SRS sections.")
    sys.exit(0)
