/**
 * /api/document-templates  GET | POST
 *
 * GET  → List templates   (GET /v1/document-templates)
 * POST → Create template  (POST /v1/document-templates)
 */
import { NextResponse } from "next/server";
import { makeDocgenClient } from "@lib/sdk/docgen-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import type { CreateDocumentTemplatePayload } from "@lib/sdk/docgen-sdk-client";

export async function GET(request: Request) {
  try {
    const hdrs = await resolvePortalHeaders(request);
    const docgen = makeDocgenClient(request, hdrs ?? undefined);
    const url = new URL(request.url);
    const result = await docgen.listTemplates({
      type: url.searchParams.get("type") ?? undefined,
      activeOnly: url.searchParams.get("active_only") === "true",
      pageSize: url.searchParams.get("page_size") ? Number(url.searchParams.get("page_size")) : undefined,
      pageToken: url.searchParams.get("page_token") ?? undefined,
    });
    return NextResponse.json(result.data, { status: result.status });
  } catch (err) {
    return NextResponse.json(
      { ok: false, message: err instanceof Error ? err.message : "Error" },
      { status: 502 }
    );
  }
}

export async function POST(request: Request) {
  try {
    const hdrs = await resolvePortalHeaders(request);
    const docgen = makeDocgenClient(request, hdrs ?? undefined);
    const body = (await request.json()) as CreateDocumentTemplatePayload;
    const result = await docgen.createTemplate(body);
    return NextResponse.json(result.data, { status: result.status });
  } catch (err) {
    return NextResponse.json(
      { ok: false, message: err instanceof Error ? err.message : "Error" },
      { status: 502 }
    );
  }
}
