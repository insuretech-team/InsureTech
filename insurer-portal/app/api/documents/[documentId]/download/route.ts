import { NextResponse } from "next/server";

import { directHttp } from "@/lib/server/insuretech";
import { loadContext } from "@/lib/server/portal-data";

export async function GET(
  request: Request,
  { params }: { params: Promise<{ documentId: string }> },
) {
  const { documentId } = await params;
  const searchParams = new URL(request.url).searchParams;
  const insurerId = searchParams.get("insurerId") ?? "";

  const context = await loadContext(request, insurerId);
  if (!context) {
    return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
  }

  const result = await directHttp(
    request,
    `/v1/documents/${encodeURIComponent(documentId)}/download`,
    { session: context.session },
  );

  if (!result.ok) {
    return NextResponse.json(
      { ok: false, message: result.message ?? "Document not found" },
      { status: 404 },
    );
  }

  // The backend returns base64 content, content_type, and filename.
  const contentB64 = (result.data.content as string | undefined) ?? "";
  const contentType = (result.data.content_type as string | undefined) ?? "application/octet-stream";
  const filename = (result.data.filename as string | undefined) ?? "document";

  if (!contentB64) {
    // Redirect to file_url if no inline bytes
    const fileUrl = (result.data.file_url as string | undefined) ?? "";
    if (fileUrl) {
      return NextResponse.redirect(fileUrl);
    }
    return NextResponse.json({ ok: false, message: "No content available" }, { status: 404 });
  }

  const buffer = Buffer.from(contentB64, "base64");

  return new Response(buffer, {
    headers: {
      "Content-Type": contentType,
      "Content-Disposition": `attachment; filename="${filename}"`,
      "Content-Length": String(buffer.length),
    },
  });
}
