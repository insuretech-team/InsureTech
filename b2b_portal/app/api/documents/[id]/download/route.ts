/**
 * /api/documents/[id]/download  GET
 * Proxies GET /v1/documents/{id}/download from the docgen service.
 */
import { NextResponse } from "next/server";
import { makeDocgenClient } from "@lib/sdk/docgen-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";

type RouteContext = { params: Promise<{ id: string }> };

export async function GET(request: Request, { params }: RouteContext) {
  const { id } = await params;
  if (!id)
    return NextResponse.json(
      { ok: false, message: "document_id required" },
      { status: 400 }
    );
  try {
    const hdrs = await resolvePortalHeaders(request);
    const docgen = makeDocgenClient(request, hdrs ?? undefined);
    const result = await docgen.downloadDocument(id);
    return NextResponse.json(result.data, { status: result.status });
  } catch (err) {
    return NextResponse.json(
      { ok: false, message: err instanceof Error ? err.message : "Error" },
      { status: 502 }
    );
  }
}
