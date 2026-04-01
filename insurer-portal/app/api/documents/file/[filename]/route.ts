import { NextResponse } from "next/server";
import path from "path";
import fs from "fs";

export async function GET(
  _request: Request,
  { params }: { params: Promise<{ filename: string }> },
) {
  const { filename } = await params;
  const safe = path.basename(decodeURIComponent(filename));
  const projectRoot = path.resolve(process.cwd(), "..");
  const filePath = path.join(projectRoot, "backend", "inscore", "generated", safe);

  if (!fs.existsSync(filePath)) {
    return NextResponse.json({ ok: false, message: "File not found" }, { status: 404 });
  }

  const buf = fs.readFileSync(filePath);
  const ext = path.extname(safe).toLowerCase();
  const contentType =
    ext === ".docx"
      ? "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
      : ext === ".pdf"
        ? "application/pdf"
        : "application/octet-stream";

  return new Response(buf, {
    headers: {
      "Content-Type": contentType,
      "Content-Disposition": `attachment; filename="${safe}"`,
      "Content-Length": String(buf.length),
    },
  });
}
