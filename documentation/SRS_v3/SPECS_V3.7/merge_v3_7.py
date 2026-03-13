#!/usr/bin/env python3
"""
Merge Script for SRS V3.7
Combines modular sections, proto schemas, and examples into final SRS_V3.7.md
- Enforces SRS to be plan- and team-agnostic (no timelines/owners)
- Auto-injects proto definitions and examples
- Appends final non-appendix sections: Sign-off & Approval, Document Status
"""

import os
from pathlib import Path
from datetime import datetime
import re

# Base paths
BASE_PATH = Path(r"G:\_0LifePlus\InsureTech\SRS_v3\SPECS_V3.7")
SECTIONS_PATH = BASE_PATH / "sections"
PROTO_PATH = BASE_PATH / "proto"
EXAMPLES_PATH = BASE_PATH / "examples"
OUTPUT_FILE = Path(r"G:\_0LifePlus\InsureTech\SRS_v3\SRS_V3.7.md")


def read_file(filepath: Path) -> str:
    with open(filepath, 'r', encoding='utf-8') as f:
        return f.read()


def write_file(filepath: Path, content: str) -> None:
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)


# -------------------------
# Header processing (sanitize planning + approvals)
# -------------------------

def load_and_sanitize_header() -> str:
    """Load 00_header.md, remove planning/timeline blocks and embedded approvals/status.
    We'll append approvals and status at the end as final sections.
    """
    header_path = SECTIONS_PATH / "00_header.md"
    if not header_path.exists():
        return ""
    content = read_file(header_path)

    # Remove 'Phased Delivery' block (planning) in Executive Summary if present
    content = re.sub(r"\*\*Phased Delivery:\*\*[\s\S]*?(?:\n\n|---|\[\[\[PAGEBREAK\]\]\])", "", content, flags=re.IGNORECASE)

    # Remove inline Document Status line in header exec summary if present
    content = re.sub(r"\*\*Document Status:\*\*.*\n", "", content)

    # Remove 'Approval Signatures' section from header (will be appended at end)
    content = re.sub(r"## Approval Signatures[\s\S]*?(?:\n---|\[\[\[PAGEBREAK\]\]\]|\n## )", "", content, flags=re.IGNORECASE)

    # Remove 'milestone' mentions in Key Changes to keep SRS plan-agnostic
    content = re.sub(r"^\s*-\s*✅?\s*\*\*?CRITICAL:?\s*Milestone[\s\S]*?$", "", content, flags=re.IGNORECASE | re.MULTILINE)

    # Remove Table of Contents entry for Project Planning section to keep SRS plan-agnostic
    content = re.sub(r"^\s*14\.\s*\[Project Planning.*\)$", "", content, flags=re.IGNORECASE | re.MULTILINE)

    return content.strip()


def extract_approval_table_from_header() -> str:
    header_path = SECTIONS_PATH / "00_header.md"
    if not header_path.exists():
        return ""
    raw = read_file(header_path)
    m = re.search(r"## Approval Signatures\s*\n([\s\S]*?)\n\s*(?:---|\[\[\[PAGEBREAK\]\]\]|\n## )", raw, flags=re.IGNORECASE)
    table = m.group(1).strip() if m else "| Role | Name | Signature | Date |\n|------|------|-----------|------|"
    return table


def extract_header_status() -> dict:
    header_path = SECTIONS_PATH / "00_header.md"
    if not header_path.exists():
        return {"Version": "", "Date": "", "Status": ""}
    raw = read_file(header_path)
    ver = re.search(r"\*\*Version:\*\*\s*(.*?)\s{2,}\n", raw)
    date = re.search(r"\*\*Date:\*\*\s*(.*?)\s{2,}\n", raw)
    status = re.search(r"\*\*Status:\*\*\s*(.*?)\n", raw)
    return {
        "Version": ver.group(1).strip() if ver else "",
        "Date": date.group(1).strip() if date else "",
        "Status": status.group(1).strip() if status else "",
    }


# -------------------------
# Sections merge (exclude planning/team content)
# -------------------------

def sanitize_planning_and_team(content: str) -> str:
    """Remove obvious planning and team ownership constructs from any section content,
    and drop any in-body document control/sign-off duplicates (final sections will be appended)."""
    # Remove tables with a 'Team Owner' column by dropping the column.
    def drop_team_owner_column(md: str) -> str:
        lines = md.splitlines()
        out = []
        i = 0
        while i < len(lines):
            line = lines[i]
            # Detect start of a markdown table block
            if line.strip().startswith('|') and '|' in line:
                # Attempt to capture a table: header, separator, data rows
                table_block = [line]
                j = i + 1
                while j < len(lines) and lines[j].strip().startswith('|'):
                    table_block.append(lines[j])
                    j += 1
                # Analyze header
                header = table_block[0]
                header_cells = [c.strip() for c in header.strip().strip('|').split('|')]
                # If this table contains 'Team Owner' header, drop that specific column index
                if any(h.lower() == 'team owner' for h in header_cells):
                    drop_idx = next(idx for idx,h in enumerate(header_cells) if h.lower() == 'team owner')
                    new_block = []
                    for row in table_block:
                        if not row.strip().startswith('|'):
                            new_block.append(row)
                            continue
                        cells = [c for c in row.strip().strip('|').split('|')]
                        # Pad cells to at least drop_idx+1
                        if len(cells) > drop_idx:
                            cells = cells[:drop_idx] + cells[drop_idx+1:]
                        # Rebuild row
                        new_row = '| ' + ' | '.join(c.strip() for c in cells) + ' |'
                        new_block.append(new_row)
                    out.extend(new_block)
                else:
                    # Not a team-owner table: keep as-is
                    out.extend(table_block)
                i = j
                continue
            else:
                out.append(line)
                i += 1
        return "\n".join(out)

    content = drop_team_owner_column(content)

    # Remove explicit milestone/timeline headers within sections (defensive)
    content = re.sub(r"###\s*\d+\.\d+\s*(Milestones|Delivery Timeline|Phased Delivery)[\s\S]*?(?=\n#|\n##|\n###|\Z)", "", content, flags=re.IGNORECASE)

    # Remove any 'Document Control' blocks (in-body)
    content = re.sub(r"^\s*\*\*Document Control:\*\*[\s\S]*?(?=\n---|\n#|\Z)", "", content, flags=re.IGNORECASE | re.MULTILINE)

    # Remove any in-body 'Document Approval & Sign-off' sections; final sign-off will be appended
    content = re.sub(r"^\s*##?\s*Document Approval\s*&\s*Sign-off[\s\S]*?(?=\n---|\n#|\Z)", "", content, flags=re.IGNORECASE | re.MULTILINE)

    return content


def normalize_lists(md: str) -> str:
    """Convert numbered markdown lists to bullets ('- ') outside of code fences and tables."""
    lines = md.splitlines()
    out = []
    in_code = False
    for line in lines:
        if line.strip().startswith("```"):
            in_code = not in_code
            out.append(line)
            continue
        if in_code or line.strip().startswith('|'):
            out.append(line)
            continue
        # Replace '  1. text' => '- text'
        m = re.match(r'^(\s*)\d+\.\s+(.*)$', line)
        if m:
            indent, rest = m.groups()
            out.append(f"{indent}- {rest}")
        else:
            out.append(line)
    return "\n".join(out)


def merge_sections() -> str:
    sections = []
    # Sorted section files, exclude header, proto appendix, and project planning
    section_files = sorted([
        f for f in SECTIONS_PATH.glob("*.md")
        if f.name not in ["00_header.md", "15_appendix_protos.md", "13_project_planning.md"]
    ])

    for section_file in section_files:
        print(f"Processing: {section_file.name}")
        content = read_file(section_file)
        content = sanitize_planning_and_team(content)
        content = normalize_lists(content)
        sections.append(content)
        sections.append("\n\n---\n[[[PAGEBREAK]]]\n\n")

    return "\n\n".join(sections)


# -------------------------
# Appendices with PROTO + EXAMPLES (letter continuation)
# -------------------------

def renumber_appendix_letters(appendix_md: str, start_from_letter: str = "G") -> str:
    """Rename 'Appendix A', 'Appendix B' headings to continue from the provided letter."""
    def next_letter(ch):
        return chr(ord(ch) + 1)

    first = next_letter(start_from_letter)
    second = next_letter(first)

    # Replace only the first two appendix headings found in the template
    appendix_md = re.sub(r"#\s*Appendix\s*A:\s*", f"# Appendix {first}: ", appendix_md, count=1)
    appendix_md = re.sub(r"#\s*Appendix\s*B:\s*", f"# Appendix {second}: ", appendix_md, count=1)
    return appendix_md


def process_appendix_with_protos(start_from_letter: str = "G") -> str:
    appendix_template = SECTIONS_PATH / "15_appendix_protos.md"
    if not appendix_template.exists():
        return ""

    print("\n3a. Processing appendix with proto/example injections...")
    content = read_file(appendix_template)

    # Continue letters from given starting letter (default 'G')
    content = renumber_appendix_letters(content, start_from_letter=start_from_letter)

    # Replace placeholders
    proto_pattern = r"\{\{PROTO:(.*?)\}\}"
    example_pattern = r"\{\{EXAMPLE:(.*?)\}\}"

    def replace_proto(match):
        proto_path = match.group(1)
        full_path = PROTO_PATH / proto_path
        if full_path.exists():
            proto_content = read_file(full_path)
            return f"```protobuf\n{proto_content}\n```"
        else:
            return f"*Proto file not found: {proto_path}*"

    def replace_example(match):
        example_file = match.group(1)
        full_path = EXAMPLES_PATH / example_file
        if full_path.exists():
            example_content = read_file(full_path)
            return example_content
        else:
            return f"*Example file not found: {example_file}*"

    content = re.sub(proto_pattern, replace_proto, content)
    content = re.sub(example_pattern, replace_example, content)

    return content


# -------------------------
# Final sections (not appendices)
# -------------------------

def render_signoff_section(approval_table_md: str) -> str:
    return f"""
---
[[[PAGEBREAK]]]

# Sign-off & Approval

{approval_table_md}
""".strip()


def render_status_section(meta: dict) -> str:
    return f"""
---
[[[PAGEBREAK]]]

# Document Status

- Version: {meta.get('Version','')}
- Date: {meta.get('Date','')}
- Status: {meta.get('Status','')}
""".strip()


def generate_footer() -> str:
    return f"""

---

## Document Generation Information

**Generated:** {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}  
**Generator:** SRS V3.7 Merge Script  
**Source:** Modular sections in SPECS_V3.7/  

---

**End of Document**
"""


def main():
    print("=" * 60)
    print("SRS V3.7 Merge Script")
    print("=" * 60)

    document_parts = []

    # 1. Header (sanitized)
    print("\n1. Loading header...")
    header = load_and_sanitize_header()
    if header:
        document_parts.append(header)
        document_parts.append("\n\n---\n[[[PAGEBREAK]]]\n\n")

    # 2. Merge sections (excluding planning)
    print("\n2. Merging sections...")
    document_parts.append(merge_sections())

    # 3. Process appendix with proto/example injections (continue letters from G)
    print("\n3. Processing appendix with proto/example injections...")
    appendix_content = process_appendix_with_protos(start_from_letter="G")
    if appendix_content:
        document_parts.append("\n\n---\n[[[PAGEBREAK]]]\n\n")
        document_parts.append(appendix_content)

    # 4. Append final sections (Sign-off & Approval, Document Status)
    print("\n4. Appending final sections (sign-off, status)...")
    approval_table = extract_approval_table_from_header()
    meta = extract_header_status()
    document_parts.append(render_signoff_section(approval_table))
    document_parts.append(render_status_section(meta))

    # 5. Footer
    print("\n5. Generating footer...")
    document_parts.append(generate_footer())

    # Combine and write
    final_document = "\n".join(document_parts)

    print(f"\n6. Writing to {OUTPUT_FILE}...")
    write_file(OUTPUT_FILE, final_document)

    # Statistics
    line_count = final_document.count('\n')
    word_count = len(final_document.split())

    print("\n" + "=" * 60)
    print("✅ SRS V3.7 Generated Successfully!")
    print("=" * 60)
    print(f"Output File: {OUTPUT_FILE}")
    print(f"Total Lines: {line_count:,}")
    print(f"Total Words: {word_count:,}")
    print(f"File Size: {len(final_document):,} bytes")
    print("=" * 60)


if __name__ == "__main__":
    main()
