import path from "path";
import fs from "fs";
import mammoth from "mammoth";

export async function GET(
  _request: Request,
  { params }: { params: Promise<{ filename: string }> },
) {
  const { filename } = await params;
  const safe = path.basename(decodeURIComponent(filename));
  const projectRoot = path.resolve(process.cwd(), "..");
  const filePath = path.join(projectRoot, "backend", "inscore", "generated", safe);

  if (!fs.existsSync(filePath)) {
    return new Response("<p>Document not found.</p>", { headers: { "Content-Type": "text/html" } });
  }

  const result = await mammoth.convertToHtml({ path: filePath }, {
    styleMap: [
      "p[style-name='Heading 1'] => h2:fresh",
      "p[style-name='Heading 2'] => h3:fresh",
    ],
  });

  const html = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <style>
    * { box-sizing: border-box; }
    html, body { margin: 0; padding: 0; background: #f8fafc; }
    body {
      font-family: "Times New Roman", Times, serif;
      font-size: 12pt;
      color: #0f172a;
      line-height: 1.55;
      padding: 28px 36px 48px;
    }
    h2 {
      color: #1f3864; font-size: 13pt; font-weight: 700;
      border-bottom: 2px solid #1f3864; padding-bottom: 3px;
      margin: 18px 0 8px;
    }
    h3 { color: #1f3864; font-size: 12pt; font-weight: 700; margin: 12px 0 6px; }
    p { margin: 3px 0; }
    img { max-width: 140px; }

    /* Tables */
    table { width: 100%; border-collapse: collapse; margin: 8px 0 14px; font-size: 11pt; }
    th {
      background: #1f3864; color: #fff;
      padding: 5px 8px; text-align: left; font-size: 10pt; font-weight: 700;
    }
    td {
      padding: 4px 8px; border: 1px solid #c8d6e5; vertical-align: top;
      transition: background 0.15s;
    }
    tr:nth-child(even) td { background: #e9eff7; }

    /* Editable value cells — last column of key-value tables */
    td[contenteditable="true"] {
      cursor: text;
      outline: none;
      border: 1.5px dashed #94a3b8;
      border-radius: 3px;
      background: #fffef0 !important;
      min-width: 80px;
    }
    td[contenteditable="true"]:hover { border-color: #1f3864; background: #f0f6ff !important; }
    td[contenteditable="true"]:focus {
      border: 2px solid #1f3864 !important;
      background: #fff !important;
      box-shadow: 0 0 0 3px rgba(31,56,100,0.12);
    }

    /* Edit mode indicator */
    .edit-hint {
      position: fixed; bottom: 16px; right: 16px;
      background: #1f3864; color: #fff;
      font-family: system-ui, sans-serif;
      font-size: 11px; padding: 6px 14px; border-radius: 20px;
      opacity: 0.85; pointer-events: none;
    }
    .save-bar {
      position: fixed; bottom: 0; left: 0; right: 0;
      background: #1f3864; color: #fff;
      display: flex; align-items: center; justify-content: space-between;
      padding: 10px 24px; gap: 12px;
      font-family: system-ui, sans-serif; font-size: 13px;
      transform: translateY(100%); transition: transform 0.25s ease;
      z-index: 100;
    }
    .save-bar.visible { transform: translateY(0); }
    .save-bar button {
      background: #fff; color: #1f3864;
      border: none; border-radius: 20px;
      padding: 6px 20px; font-size: 12px; font-weight: 700;
      cursor: pointer;
    }
    .save-bar button:hover { background: #e9eff7; }
    .save-bar .discard { background: transparent; color: rgba(255,255,255,0.7); }
    .save-bar .discard:hover { color: #fff; background: rgba(255,255,255,0.1); }
  </style>
</head>
<body>
${result.value}

<div class="edit-hint">✎ Click any yellow cell to edit</div>
<div class="save-bar" id="saveBar">
  <span>You have unsaved changes</span>
  <div style="display:flex;gap:8px">
    <button class="discard" onclick="discardChanges()">Discard</button>
    <button onclick="saveChanges()">Save &amp; Regenerate</button>
  </div>
</div>

<script>
  // Make the last cell of every 2-3 column table row editable (value column)
  const tables = document.querySelectorAll('table');
  tables.forEach(table => {
    const rows = table.querySelectorAll('tr');
    rows.forEach(row => {
      const cells = row.querySelectorAll('td');
      // Skip header rows and single-cell rows
      if (cells.length >= 2) {
        const lastCell = cells[cells.length - 1];
        // Don't make cells editable that look like headers (bold, dark bg)
        const style = window.getComputedStyle(lastCell);
        lastCell.setAttribute('contenteditable', 'true');
        lastCell.setAttribute('data-original', lastCell.innerText);
        lastCell.addEventListener('input', onEdit);
      }
    });
  });

  let hasChanges = false;
  const saveBar = document.getElementById('saveBar');

  function onEdit() {
    hasChanges = true;
    saveBar.classList.add('visible');
  }

  function collectChanges() {
    const fields = {};
    document.querySelectorAll('td[contenteditable]').forEach(cell => {
      const row = cell.closest('tr');
      // Use the first cell's text as the key (label), last cell as value
      const cells = row.querySelectorAll('td');
      let label = '';
      for (let i = 0; i < cells.length - 1; i++) {
        const t = cells[i].innerText.trim();
        if (t) label = (label ? label + ' ' : '') + t;
      }
      label = label.trim();
      const value = cell.innerText.trim();
      if (label && label.length < 120) {
        fields[label] = value;
      }
    });
    return fields;
  }

  function saveChanges() {
    const fields = collectChanges();
    window.parent.postMessage({ type: 'DOC_SAVE', fields }, '*');
  }

  function discardChanges() {
    document.querySelectorAll('td[contenteditable]').forEach(cell => {
      cell.innerText = cell.getAttribute('data-original') || '';
    });
    hasChanges = false;
    saveBar.classList.remove('visible');
  }

  // Listen for parent telling us to reset
  window.addEventListener('message', e => {
    if (e.data && e.data.type === 'DOC_RESET') discardChanges();
  });
</script>
</body>
</html>`;

  return new Response(html, {
    headers: { "Content-Type": "text/html; charset=utf-8" },
  });
}
