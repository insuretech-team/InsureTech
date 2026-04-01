"""
DocRender Sidecar — FastAPI service for rich document generation.

Endpoints:
  POST /render/docx   → Jinja2-templated DOCX via python-docx
  POST /render/pdf    → HTML → PDF via WeasyPrint (CSS-accurate, no browser needed)
  GET  /health        → liveness probe

Called exclusively by the Go docgen microservice over localhost.
"""

from __future__ import annotations

import io
import logging
from typing import Any

from fastapi import FastAPI, HTTPException
from fastapi.responses import Response
from pydantic import BaseModel

from docx_renderer import render_docx
from pdf_renderer import render_pdf_from_html

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
logger = logging.getLogger(__name__)

app = FastAPI(title="InsureTech DocRender Sidecar", version="1.0.0")


# ─── Request / Response models ────────────────────────────────────────────────

class DocxRequest(BaseModel):
    """Request body for DOCX generation."""
    # Jinja2 template string describing the document structure.
    # Supports sections: title, subtitle, paragraphs, tables, page_break.
    template_content: str
    # Key-value data to render into the template.
    data: dict[str, Any] = {}
    # Optional document metadata.
    title: str = ""
    author: str = "InsureTech"
    subject: str = ""


class PdfRequest(BaseModel):
    """Request body for WeasyPrint PDF generation."""
    html: str
    base_url: str = ""


# ─── Routes ──────────────────────────────────────────────────────────────────

@app.get("/health")
def health() -> dict:
    return {"status": "ok", "service": "docrender-sidecar"}


@app.post("/render/docx")
def docx_endpoint(req: DocxRequest) -> Response:
    """
    Render a DOCX file from a Jinja2 template + data payload.
    Returns raw .docx bytes with content-type application/vnd.openxmlformats-officedocument.wordprocessingml.document
    """
    try:
        buf = render_docx(
            template_content=req.template_content,
            data=req.data,
            title=req.title,
            author=req.author,
            subject=req.subject,
        )
        return Response(
            content=buf.getvalue(),
            media_type=(
                "application/vnd.openxmlformats-officedocument"
                ".wordprocessingml.document"
            ),
            headers={"Content-Disposition": 'attachment; filename="document.docx"'},
        )
    except ValueError as exc:
        logger.warning("DOCX render validation error: %s", exc)
        raise HTTPException(status_code=422, detail=str(exc)) from exc
    except Exception as exc:
        logger.exception("DOCX render failed")
        raise HTTPException(status_code=500, detail=f"DOCX render failed: {exc}") from exc


@app.post("/render/pdf")
def pdf_endpoint(req: PdfRequest) -> Response:
    """
    Render a PDF from an HTML string using WeasyPrint.
    Returns raw PDF bytes with content-type application/pdf.
    """
    try:
        pdf_bytes = render_pdf_from_html(html=req.html, base_url=req.base_url or None)
        return Response(
            content=pdf_bytes,
            media_type="application/pdf",
            headers={"Content-Disposition": 'attachment; filename="document.pdf"'},
        )
    except Exception as exc:
        logger.exception("PDF render failed")
        raise HTTPException(status_code=500, detail=f"PDF render failed: {exc}") from exc
