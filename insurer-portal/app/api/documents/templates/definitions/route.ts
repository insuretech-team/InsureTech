import { NextResponse } from "next/server";
import path from "path";
import fs from "fs";

const DEFS_DIR = path.resolve(process.cwd(), "..", "backend", "inscore", "templates", "insurance", "definitions");

export async function GET() {
  fs.mkdirSync(DEFS_DIR, { recursive: true });
  const files = fs.readdirSync(DEFS_DIR).filter((f) => f.endsWith(".json"));
  const definitions = files.map((file) => {
    try {
      const raw = JSON.parse(fs.readFileSync(path.join(DEFS_DIR, file), "utf8")) as {
        id?: string;
        font?: string;
        company?: { name?: string };
        sections?: unknown[];
      };
      return {
        id: raw.id ?? file.replace(".json", ""),
        filename: file,
        name: toLabel(raw.id ?? file.replace(".json", "")),
        company: raw.company?.name ?? "Pragati Insurance PLC",
        sectionCount: (raw.sections ?? []).length,
        isBuiltIn: true,
      };
    } catch {
      return null;
    }
  }).filter(Boolean);

  return NextResponse.json({ ok: true, data: definitions });
}

export async function POST(request: Request) {
  let body: Record<string, unknown> = {};
  try {
    body = (await request.json()) as Record<string, unknown>;
  } catch {
    return NextResponse.json({ ok: false, message: "Invalid JSON" }, { status: 400 });
  }

  const { id, definition } = body as { id?: string; definition?: Record<string, unknown> };
  if (!id || !definition) {
    return NextResponse.json({ ok: false, message: "id and definition are required" }, { status: 400 });
  }

  const safeId = id.replace(/[^a-z0-9_]/gi, "_").toLowerCase();
  const filePath = path.join(DEFS_DIR, `${safeId}.json`);
  fs.mkdirSync(DEFS_DIR, { recursive: true });
  fs.writeFileSync(filePath, JSON.stringify({ id: safeId, ...definition }, null, 2), "utf8");

  return NextResponse.json({ ok: true, data: { id: safeId, filename: `${safeId}.json` } });
}

function toLabel(id: string): string {
  return id.replace(/_/g, " ").replace(/\b\w/g, (c) => c.toUpperCase());
}
