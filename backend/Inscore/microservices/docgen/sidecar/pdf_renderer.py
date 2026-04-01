"""
PDF renderer using WeasyPrint.

WeasyPrint converts HTML+CSS to PDF with full CSS support including:
  - @page rules (size, margins, page numbers)
  - CSS Grid and Flexbox
  - Custom fonts via @font-face
  - Background colors and images
  - CSS counters for page numbering

Unlike Gotenberg (which needs a Chromium browser), WeasyPrint is a
pure Python library — lighter, faster for simple docs, and doesn't
require a browser runtime.

Gotenberg remains the primary PDF renderer in the Go service (for
complex JS-heavy pages). WeasyPrint is the sidecar alternative used
when Gotenberg is unavailable or when DOCX→PDF conversion is needed.
"""

from __future__ import annotations

import logging
from typing import Optional

import weasyprint

logger = logging.getLogger(__name__)

# Inject sensible @page defaults if none are present in the HTML.
_DEFAULT_CSS = weasyprint.CSS(string="""
@page {
    size: A4;
    margin: 20mm 18mm 25mm 18mm;

    @bottom-center {
        content: "Page " counter(page) " of " counter(pages);
        font-family: Arial, sans-serif;
        font-size: 9pt;
        color: #666;
    }
}

body {
    font-family: Arial, Helvetica, sans-serif;
    font-size: 11pt;
    line-height: 1.5;
    color: #1a1a2e;
}

table {
    border-collapse: collapse;
    width: 100%;
}

th, td {
    border: 1px solid #ccd;
    padding: 6px 8px;
    text-align: left;
}

th {
    background-color: #0D47A1;
    color: white;
    font-weight: bold;
}

tr:nth-child(even) td {
    background-color: #EEF2F7;
}

h1, h2, h3 {
    color: #0D47A1;
    page-break-after: avoid;
}

.page-break {
    page-break-before: always;
}

.signature-line {
    border-bottom: 1px solid #333;
    width: 200px;
    display: inline-block;
    margin-top: 40px;
}

.totals-table td:last-child {
    text-align: right;
    font-weight: bold;
}

.totals-total {
    background-color: #0D47A1;
    color: white;
}
""")


def render_pdf_from_html(
    html: str,
    base_url: Optional[str] = None,
) -> bytes:
    """
    Convert an HTML string to PDF bytes using WeasyPrint.

    Args:
        html:      Full HTML document string.
        base_url:  Optional base URL for resolving relative resources
                   (images, fonts, stylesheets).

    Returns:
        PDF content as bytes.
    """
    if not html or not html.strip():
        raise ValueError("HTML content is empty")

    wp_html = weasyprint.HTML(string=html, base_url=base_url)
    pdf_bytes = wp_html.write_pdf(stylesheets=[_DEFAULT_CSS])

    if not pdf_bytes:
        raise RuntimeError("WeasyPrint returned empty PDF")

    logger.info("WeasyPrint rendered PDF: %d bytes", len(pdf_bytes))
    return pdf_bytes
