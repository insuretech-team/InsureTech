import { NextResponse } from "next/server";
import { execFile } from "child_process";
import { promisify } from "util";
import path from "path";
import fs from "fs";

const execFileAsync = promisify(execFile);

export async function POST(request: Request) {
  let body: Record<string, unknown> = {};
  try {
    body = (await request.json()) as Record<string, unknown>;
  } catch {
    return NextResponse.json({ ok: false, message: "Invalid JSON" }, { status: 400 });
  }

  const { definition, sampleData } = body as {
    definition?: Record<string, unknown>;
    sampleData?: Record<string, unknown>;
  };

  if (!definition) {
    return NextResponse.json({ ok: false, message: "definition is required" }, { status: 400 });
  }

  const projectRoot   = path.resolve(process.cwd(), "..");
  const generatedDir  = path.join(projectRoot, "backend", "inscore", "generated");
  const venvPython    = path.join(projectRoot, ".venv", "Scripts", "python.exe");
  const sidecarDir    = path.join(projectRoot, "backend", "inscore", "microservices", "docgen", "sidecar");
  const logoPath      = path.join(projectRoot, "web_shared", "insurers", "pragati_logo.png");

  fs.mkdirSync(generatedDir, { recursive: true });

  const filename   = `preview_${Date.now()}.docx`;
  const outPath    = path.join(generatedDir, filename);
  const tmpTpl     = path.join(generatedDir, `_tpl_${Date.now()}.json`);
  const tmpData    = path.join(generatedDir, `_dat_${Date.now()}.json`);
  const runnerPath = path.join(generatedDir, `_run_${Date.now()}.py`);

  // Build sample data with logo + placeholders for all {{ keys }} in definition
  const merged: Record<string, unknown> = {
    logo_path: logoPath,
    proposal_id: "PREVIEW-001",
    generated_at: new Date().toLocaleString("en-GB"),
    ...(sampleData ?? {}),
  };

  // Auto-fill any un-provided placeholder keys with "Sample Value"
  const defStr = JSON.stringify(definition);
  const keys = [...defStr.matchAll(/\{\{\s*(\w+)\s*\}\}/g)].map((m) => m[1]);
  for (const key of keys) {
    if (!(key in merged)) merged[key] = `[${key}]`;
  }

  // Auto-fill rows_key tables with 2 sample rows
  function fillTableRows(def: unknown): void {
    if (!def || typeof def !== "object") return;
    if (Array.isArray(def)) { def.forEach(fillTableRows); return; }
    const obj = def as Record<string, unknown>;
    if (obj.type === "table" && obj.rows_key) {
      const key = obj.rows_key as string;
      if (!(key in merged)) {
        const cols = (obj.columns as string[] | undefined) ?? [];
        merged[key] = [
          Object.fromEntries(cols.map((c) => [c, `[${c}]`])),
          Object.fromEntries(cols.map((c) => [c, `[${c}]`])),
        ];
      }
    }
    Object.values(obj).forEach(fillTableRows);
  }
  fillTableRows(definition);

  fs.writeFileSync(tmpTpl, JSON.stringify(definition), "utf8");
  fs.writeFileSync(tmpData, JSON.stringify(merged), "utf8");

  const runner = `
import sys, json
sys.path.insert(0, r"${sidecarDir.replace(/\\/g, "\\\\")}")
from docx_builder import DocxBuilder
with open(r"${tmpTpl.replace(/\\/g, "\\\\")}", encoding="utf-8") as f: tpl = json.load(f)
with open(r"${tmpData.replace(/\\/g, "\\\\")}", encoding="utf-8") as f: data = json.load(f)
buf = DocxBuilder(tpl, data).build()
with open(r"${outPath.replace(/\\/g, "\\\\")}", "wb") as f: f.write(buf.read())
print("OK")
`.trim();

  fs.writeFileSync(runnerPath, runner, "utf8");

  try {
    const { stdout, stderr } = await execFileAsync(venvPython, [runnerPath], { timeout: 30000 });
    fs.unlinkSync(tmpTpl); fs.unlinkSync(tmpData); fs.unlinkSync(runnerPath);
    if (!stdout.includes("OK")) {
      return NextResponse.json({ ok: false, message: stderr || stdout || "Preview failed" });
    }
    return NextResponse.json({
      ok: true,
      data: {
        filename,
        downloadUrl: `/api/documents/file/${encodeURIComponent(filename)}`,
        renderUrl: `/api/documents/render/${encodeURIComponent(filename)}`,
      },
    });
  } catch (err: unknown) {
    try { fs.unlinkSync(tmpTpl); } catch {}
    try { fs.unlinkSync(tmpData); } catch {}
    try { fs.unlinkSync(runnerPath); } catch {}
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : String(err) });
  }
}
