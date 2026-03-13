/**
 * /api/document-templates/[id]  GET | PATCH | DELETE
 */
import { NextResponse } from "next/server";
import { makeDocgenClient } from "@lib/sdk/docgen-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import type { UpdateDocumentTemplatePayload } from "@lib/sdk/docgen-sdk-client";

type RouteContext = { params: Promise<{ id: string }> };

export async function GET(request: Request, { params }: RouteContext) {
  const { id } = await params;
  try {
    const hdrs = await resolvePortalHeaders(request);
    const docgen = makeDocgenClient(request, hdrs ?? undefined);
    const result = await docgen.getTemplate(id);
    return NextResponse.json(result.data, { status: result.status });
  } catch (err) {
    return NextResponse.json(
      { ok: false, message: err instanceof Error ? err.message : "Error" },
      { status: 502 }
    );
  }
}

export async function PATCH(request: Request, { params }: RouteContext) {
  const { id } = await params;
  try {
    const hdrs = await resolvePortalHeaders(request);
    const docgen = makeDocgenClient(request, hdrs ?? undefined);
    const body = (await request.json()) as UpdateDocumentTemplatePayload;
    const result = await docgen.updateTemplate(id, body);
    return NextResponse.json(result.data, { status: result.status });
  } catch (err) {
    return NextResponse.json(
      { ok: false, message: err instanceof Error ? err.message : "Error" },
      { status: 502 }
    );
  }
}

export async function DELETE(request: Request, { params }: RouteContext) {
  const { id } = await params;
  try {
    const hdrs = await resolvePortalHeaders(request);
    const docgen = makeDocgenClient(request, hdrs ?? undefined);
    const result = await docgen.deleteTemplate(id);
    return NextResponse.json(result.data, { status: result.status });
  } catch (err) {
    return NextResponse.json(
      { ok: false, message: err instanceof Error ? err.message : "Error" },
      { status: 502 }
    );
  }
}
