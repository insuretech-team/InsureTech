import { NextResponse } from "next/server";

import { directHttp } from "@/lib/server/insuretech";
import { mapLiveDocument } from "@/lib/server/mappers";
import { loadContext } from "@/lib/server/portal-data";

export async function POST(request: Request) {
  const searchParams = new URL(request.url).searchParams;
  const insurerId = searchParams.get("insurerId") ?? "";

  const context = await loadContext(request, insurerId);
  if (!context) {
    return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  }

  let body: Record<string, unknown> = {};
  try {
    body = (await request.json()) as Record<string, unknown>;
  } catch {
    return NextResponse.json({ ok: false, message: "Invalid request body" }, { status: 400 });
  }

  const { templateId, entityType, entityId, data, outputFormat, includeQrCode } = body as {
    templateId?: string;
    entityType?: string;
    entityId?: string;
    data?: Record<string, unknown>;
    outputFormat?: string;
    includeQrCode?: boolean;
  };

  if (!templateId || !entityType || !entityId) {
    return NextResponse.json(
      { ok: false, message: "templateId, entityType and entityId are required" },
      { status: 400 },
    );
  }

  const result = await directHttp(request, "/v1/documents/generate", {
    method: "POST",
    session: context.session,
    body: {
      template_id: templateId,
      entity_type: entityType,
      entity_id: entityId,
      data: data ?? {},
      output_format: outputFormat ?? "",
      include_qr_code: includeQrCode ?? false,
    },
  });

  if (!result.ok) {
    return NextResponse.json({
      ok: false,
      message: result.message ?? "Document generation failed",
    });
  }

  // The generation response includes document_id and file_url at minimum.
  const generated = result.data;

  return NextResponse.json({
    ok: true,
    data: {
      documentId: generated.document_id ?? generated.documentId ?? "",
      fileUrl: generated.file_url ?? generated.fileUrl ?? "",
      message: generated.message ?? "Document generated successfully",
    },
  });
}
