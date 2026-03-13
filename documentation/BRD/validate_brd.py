#!/usr/bin/env python3
"""
Comprehensive BRD validation and anomaly detection script.

Checks for:
- Structural issues
- Content inconsistencies
- Formatting problems
- Traceability gaps
- Image reference issues
"""

from __future__ import annotations

import re
from pathlib import Path
from collections import defaultdict, Counter

ROOT = Path(__file__).resolve().parent
BRD_FILE = ROOT / "BRDV3.7.md"
SRS_ROOT = ROOT.parent / "SRS_v3" / "SPECS_V3.7" / "sections"
IMAGES_DIR = ROOT / "images"


def read(p: Path) -> str:
    return p.read_text(encoding="utf-8")


class BRDValidator:
    def __init__(self):
        self.brd_content = read(BRD_FILE)
        self.errors = []
        self.warnings = []
        self.info = []

    def validate_all(self):
        """Run all validation checks."""
        print("=" * 60)
        print("BRD V3.7 COMPREHENSIVE VALIDATION")
        print("=" * 60)
        print()

        self.check_structure()
        self.check_feature_groups()
        self.check_user_stories()
        self.check_business_rules()
        self.check_traceability()
        self.check_images()
        self.check_tables()
        self.check_duplicate_content()
        self.check_srs_coverage()
        
        self.report()

    def check_structure(self):
        """Check BRD structural completeness."""
        print("Checking structure...")
        
        required_sections = [
            ("Executive Summary", r"# Executive Summary"),
            ("Business Context", r"# Business Context"),
            ("Portals & Channels", r"# 4\. Portals & Channels"),
            ("Detailed Functional Catalog", r"# 6\. Detailed Business Functional Requirements"),
            ("NFR Catalog", r"# 7\. Non-Functional Requirements"),
            ("Security & Compliance", r"# 8\. Security, Privacy, and Compliance"),
            ("Traceability", r"# 9\. Traceability"),
        ]
        
        for name, pattern in required_sections:
            if not re.search(pattern, self.brd_content):
                self.errors.append(f"Missing required section: {name}")
            else:
                self.info.append(f"✓ Section found: {name}")

    def check_feature_groups(self):
        """Check Feature Group consistency."""
        print("Checking Feature Groups...")
        
        # Find all FG sections
        fg_pattern = r"# Feature Group:.*?\((FG-\d+)\)"
        fgs = re.findall(fg_pattern, self.brd_content)
        
        # Expected FGs (FG-001 to FG-023, excluding FG-021 which is missing from SRS)
        # Note: FG-06 uses 2-digit format in SRS (inconsistency)
        expected_fgs = [f"FG-{i:03d}" for i in range(1, 24) if i not in [6, 21]]
        expected_fgs.append("FG-06")  # Add the 2-digit variant
        
        found_fgs = set(fgs)
        expected_set = set(expected_fgs)
        
        missing_fgs = expected_set - found_fgs
        extra_fgs = found_fgs - expected_set
        
        if missing_fgs:
            self.errors.append(f"Missing Feature Groups: {sorted(missing_fgs)}")
        
        if extra_fgs:
            self.warnings.append(f"Unexpected Feature Groups: {sorted(extra_fgs)}")
        
        # Check for duplicates
        fg_counts = Counter(fgs)
        duplicates = {fg: count for fg, count in fg_counts.items() if count > 1}
        if duplicates:
            self.errors.append(f"Duplicate Feature Groups: {duplicates}")
        
        self.info.append(f"✓ Found {len(found_fgs)} Feature Groups")

    def check_user_stories(self):
        """Check User Story consistency."""
        print("Checking User Stories...")
        
        # Find all user stories
        us_pattern = r"### (US-FG-\d+-\d+):"
        user_stories = re.findall(us_pattern, self.brd_content)
        
        # Check for duplicates
        us_counts = Counter(user_stories)
        duplicates = {us: count for us, count in us_counts.items() if count > 1}
        if duplicates:
            self.errors.append(f"Duplicate User Story IDs: {duplicates}")
        
        # Check US-FG numbering consistency
        fg_us_map = defaultdict(list)
        for us in user_stories:
            match = re.match(r"US-(FG-\d+)-(\d+)", us)
            if match:
                fg_id = match.group(1)
                us_num = int(match.group(2))
                fg_us_map[fg_id].append(us_num)
        
        # Check for gaps in numbering
        for fg_id, us_nums in fg_us_map.items():
            us_nums_sorted = sorted(us_nums)
            expected = list(range(1, len(us_nums_sorted) + 1))
            if us_nums_sorted != expected:
                self.warnings.append(f"{fg_id}: User Story numbering has gaps/jumps: {us_nums_sorted}")
        
        self.info.append(f"✓ Found {len(user_stories)} User Stories")

    def check_business_rules(self):
        """Check Business Rule consistency."""
        print("Checking Business Rules...")
        
        # Find all business rules
        br_pattern = r"\| (BR-[A-Z]+-\d+) \|"
        business_rules = re.findall(br_pattern, self.brd_content)
        
        # Check for duplicates
        br_counts = Counter(business_rules)
        duplicates = {br: count for br, count in br_counts.items() if count > 1}
        if duplicates:
            self.warnings.append(f"Duplicate Business Rule IDs: {duplicates}")
        
        self.info.append(f"✓ Found {len(business_rules)} Business Rules")

    def check_traceability(self):
        """Check FR/NFR/SEC traceability."""
        print("Checking traceability...")
        
        # Find all FR references
        fr_pattern = r"\bFR-\d+\b"
        frs = set(re.findall(fr_pattern, self.brd_content))
        
        # Find all NFR references
        nfr_pattern = r"\bNFR-\d+\b"
        nfrs = set(re.findall(nfr_pattern, self.brd_content))
        
        # Find all SEC references
        sec_pattern = r"\bSEC-\d+\b"
        secs = set(re.findall(sec_pattern, self.brd_content))
        
        self.info.append(f"✓ FR references: {len(frs)} unique IDs")
        self.info.append(f"✓ NFR references: {len(nfrs)} unique IDs")
        self.info.append(f"✓ SEC references: {len(secs)} unique IDs")
        
        # Check for suspiciously high FR numbers (potential typos)
        high_frs = [fr for fr in frs if int(fr.split('-')[1]) > 300]
        if high_frs:
            self.warnings.append(f"Suspiciously high FR numbers (check for typos): {sorted(high_frs)}")

    def check_images(self):
        """Check image references."""
        print("Checking image references...")
        
        # Find all image references
        img_pattern = r"!\[.*?\]\(images/(.*?\.png)\)"
        images = re.findall(img_pattern, self.brd_content)
        
        # Check for duplicates
        img_counts = Counter(images)
        
        # Check if image files exist
        missing_images = []
        for img in set(images):
            img_path = IMAGES_DIR / img
            if not img_path.exists():
                missing_images.append(img)
        
        if missing_images:
            self.warnings.append(f"Referenced images not found in images/ folder: {len(missing_images)} images")
        
        self.info.append(f"✓ Found {len(images)} image references ({len(set(images))} unique)")

    def check_tables(self):
        """Check markdown table formatting."""
        print("Checking tables...")
        
        lines = self.brd_content.split('\n')
        
        in_table = False
        table_start = 0
        issues = []
        
        for i, line in enumerate(lines):
            stripped = line.strip()
            
            # Detect table header separator
            if re.match(r'^\|[\s\-:|]+\|$', stripped):
                in_table = True
                table_start = i - 1
                continue
            
            # If in table, check for proper formatting
            if in_table:
                if stripped.startswith('|') and stripped.endswith('|'):
                    # Count pipes
                    pipe_count = stripped.count('|')
                    # Should have at least 2 pipes (start and end)
                    if pipe_count < 2:
                        issues.append(f"Line {i+1}: Malformed table row")
                elif stripped == '':
                    in_table = False
                elif not stripped.startswith('#') and not stripped.startswith('[[[PAGEBREAK]]]'):
                    # Table ended unexpectedly
                    in_table = False
        
        if issues:
            for issue in issues[:10]:  # Show first 10
                self.warnings.append(issue)
            if len(issues) > 10:
                self.warnings.append(f"... and {len(issues) - 10} more table issues")
        else:
            self.info.append("✓ All tables properly formatted")

    def check_duplicate_content(self):
        """Check for accidentally duplicated content."""
        print("Checking for duplicate content...")
        
        # Check for duplicate headings (might indicate copy-paste errors)
        heading_pattern = r"^(#{1,6})\s+(.+)$"
        headings = []
        for line in self.brd_content.split('\n'):
            match = re.match(heading_pattern, line.strip())
            if match:
                level = len(match.group(1))
                text = match.group(2).strip()
                headings.append((level, text))
        
        # Look for exact duplicate headings at same level (excluding common ones)
        exclude_headings = ["Business Objective", "Actors & Portals", "User Stories", 
                           "Business Rules", "Key Workflows", "Data Model Notes",
                           "Integration Touchpoints", "Security & Privacy", "NFR Constraints",
                           "Acceptance Criteria", "Traceability"]
        
        from collections import Counter
        heading_counts = Counter([h for h in headings if h[1] not in exclude_headings])
        duplicates = {h: count for h, count in heading_counts.items() if count > 1}
        
        if duplicates:
            self.warnings.append(f"Duplicate headings found: {len(duplicates)} (may be intentional)")
        else:
            self.info.append("✓ No unexpected duplicate headings")

    def check_srs_coverage(self):
        """Check coverage against SRS."""
        print("Checking SRS coverage...")
        
        # Read SRS FR section
        srs_fr_file = SRS_ROOT / "04_functional_requirements.md"
        if srs_fr_file.exists():
            srs_content = read(srs_fr_file)
            
            # Find all FR IDs in SRS
            srs_frs = set(re.findall(r"\bFR-\d+\b", srs_content))
            
            # Find all FR IDs in BRD
            brd_frs = set(re.findall(r"\bFR-\d+\b", self.brd_content))
            
            # Find missing FRs
            missing_frs = srs_frs - brd_frs
            
            if missing_frs:
                self.warnings.append(f"FRs from SRS not referenced in BRD: {len(missing_frs)} FRs")
                if len(missing_frs) <= 10:
                    self.warnings.append(f"Missing: {sorted(missing_frs)}")
            else:
                self.info.append("✓ All SRS FRs referenced in BRD")
            
            # Check coverage percentage
            coverage = (len(brd_frs) / len(srs_frs)) * 100 if srs_frs else 0
            self.info.append(f"✓ FR coverage: {coverage:.1f}% ({len(brd_frs)}/{len(srs_frs)})")

    def report(self):
        """Generate final report."""
        print()
        print("=" * 60)
        print("VALIDATION REPORT")
        print("=" * 60)
        print()
        
        if self.errors:
            print(f"❌ ERRORS ({len(self.errors)}):")
            for err in self.errors:
                print(f"  - {err}")
            print()
        
        if self.warnings:
            print(f"⚠️  WARNINGS ({len(self.warnings)}):")
            for warn in self.warnings[:20]:  # Show first 20
                print(f"  - {warn}")
            if len(self.warnings) > 20:
                print(f"  ... and {len(self.warnings) - 20} more warnings")
            print()
        
        if self.info:
            print(f"ℹ️  INFO ({len(self.info)}):")
            for inf in self.info:
                print(f"  {inf}")
            print()
        
        print("=" * 60)
        if not self.errors:
            print("✅ VALIDATION PASSED (no critical errors)")
        else:
            print("❌ VALIDATION FAILED (critical errors found)")
        print("=" * 60)


if __name__ == "__main__":
    validator = BRDValidator()
    validator.validate_all()
