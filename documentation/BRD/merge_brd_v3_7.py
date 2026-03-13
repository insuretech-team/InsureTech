#!/usr/bin/env python3

from pathlib import Path
from datetime import datetime

BASE_PATH = Path(__file__).resolve().parent
SECTIONS_PATH = BASE_PATH / "sections"
OUTPUT_FILE = BASE_PATH / "BRDV3.7.md"


def read_file(filepath: Path) -> str:
    return filepath.read_text(encoding="utf-8")


def write_file(filepath: Path, content: str) -> None:
    filepath.write_text(content, encoding="utf-8")


def merge_sections() -> str:
    parts = []

    section_files = sorted(SECTIONS_PATH.glob("*.md"))
    for section_file in section_files:
        content = read_file(section_file).strip()
        if not content:
            continue
        parts.append(content)
        parts.append("\n\n---\n[[[PAGEBREAK]]]\n\n")

    return "\n".join(parts).strip()


def generate_footer() -> str:
    return f"""

---

## Document Generation Information

**Generated:** {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}  
**Generator:** BRD V3.7 Merge Script  
**Source:** Modular sections in BRD/sections/  

---

**End of Document**
"""


def main() -> None:
    document = merge_sections()
    if document:
        document = document + "\n" + generate_footer()
    write_file(OUTPUT_FILE, document)


if __name__ == "__main__":
    main()
