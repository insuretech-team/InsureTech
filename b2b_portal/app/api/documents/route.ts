/**
 * /api/documents  POST | GET
 *
 * POST → Generate a document  (POST /v1/documents)
 * GET  → List documents for an entity (GET /v1/entities/{entity_type}/{entity_id}/documents)
 */
import { NextResponse } from "next/server";
import { makeDocgenClient } from "@lib/sdk/docgen-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import type { GenerateDocumentPayload } from "@lib/sdk/docgen-sdk-client";

export async function POST(request: Request) {
  try {
    const hdrs = await resolvePortalHeaders(request);
    const docgen = makeDocgenClient(request, hdrs ?? undefined);
    const body = (await request.json()) as GenerateDocumentPayload;
    const result = await docgen.generate(body);
    return NextResponse.json(result.data, { status: result.status });
  } catch (err) {
    return NextResponse.json(
      { ok: false, message: err instanceof Error ? err.message : "Error" },
      { status: 502 }
    );
  }
}

export async function GET(request: Request) {
  try {
    const hdrs = await resolvePortalHeaders(request);
    const docgen = makeDocgenClient(request, hdrs ?? undefined);
    const url = new URL(request.url);
    const entityType = url.searchParams.get("entity_type") ?? "employee";
    const entityId = url.searchParams.get("entity_id") ?? "";
    const status = url.searchParams.get("status") ?? undefined;
    const page = url.searchParams.get("page");
    const pageSize = url.searchParams.get("page_size");
    const result = await docgen.listDocuments({
      entityType,
      entityId,
      status,
      page: page ? Number(page) : undefined,
      pageSize: pageSize ? Number(pageSize) : undefined,
    });
    return NextResponse.json(result.data, { status: result.status });
  } catch (err) {
    return NextResponse.json(
      { ok: false, message: err instanceof Error ? err.message : "Error" },
      { status: 502 }
    );
  }
}
