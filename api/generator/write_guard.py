"""
Write guard utility — prevents file writes when content is unchanged.

Eliminates git drift from pipeline regeneration by comparing new content
against existing files and skipping writes when identical.
"""

import os
import io
import json
import shutil

import yaml


# ---------------------------------------------------------------------------
# Custom YAML Dumper: always emit None as explicit `null` (not a bare key).
# This ensures API response examples render correctly in Swagger UI / Redoc /
# Apidog — e.g. `error: null` and `data: null` rather than `error:` / `data:`.
# ---------------------------------------------------------------------------
class _NullSafeDumper(yaml.Dumper):
    pass

def _represent_none(dumper, _):
    return dumper.represent_scalar('tag:yaml.org,2002:null', 'null')

_NullSafeDumper.add_representer(type(None), _represent_none)


def write_if_changed(path: str, content: str, encoding: str = 'utf-8') -> bool:
    """Write *content* to *path* only when the file doesn't exist or differs.

    Returns True if the file was (re)written, False if skipped.

    Line-ending normalisation: always compares and writes LF-only (\n) so that
    Windows CRLF files on disk are never seen as "different" from freshly
    generated LF content — preventing spurious git churn on every pipeline run.
    """
    # Normalise the new content to LF
    content_lf = content.replace('\r\n', '\n').replace('\r', '\n')

    if os.path.isfile(path):
        try:
            # Read raw bytes then decode — avoids Python's automatic CRLF→LF
            # translation which would silently differ from binary-mode writes.
            with open(path, 'rb') as f:
                existing_bytes = f.read()
            existing_lf = existing_bytes.decode(encoding).replace('\r\n', '\n').replace('\r', '\n')
            if existing_lf == content_lf:
                return False
        except Exception:
            pass  # unreadable → overwrite

    os.makedirs(os.path.dirname(path) or '.', exist_ok=True)
    # Always write LF line endings (universal, avoids CRLF/LF flip-flop)
    with open(path, 'w', encoding=encoding, newline='\n') as f:
        f.write(content_lf)
    return True


def yaml_dump_if_changed(path: str, data, encoding: str = 'utf-8', **kwargs) -> bool:
    """Serialize *data* to YAML and write only if the result differs from the existing file.

    Uses _NullSafeDumper so Python None values are always written as explicit
    `null` rather than a bare key — critical for correct OpenAPI example rendering.
    """
    buf = io.StringIO()
    # Inject our null-safe dumper unless the caller has overridden Dumper explicitly
    kwargs.setdefault('Dumper', _NullSafeDumper)
    yaml.dump(data, buf, **kwargs)
    content = buf.getvalue()
    return write_if_changed(path, content, encoding=encoding)


def json_dump_if_changed(path: str, data, encoding: str = 'utf-8', **kwargs) -> bool:
    """Serialize *data* to JSON and write only if the result differs from the existing file."""
    content = json.dumps(data, **kwargs) + '\n'
    return write_if_changed(path, content, encoding=encoding)


def copy_if_changed(src: str, dst: str) -> bool:
    """Copy *src* to *dst* only when content differs (or *dst* doesn't exist).

    Uses binary comparison so it works for any file type.
    """
    if os.path.isfile(dst):
        try:
            with open(src, 'rb') as fs, open(dst, 'rb') as fd:
                if fs.read() == fd.read():
                    return False
        except Exception:
            pass

    os.makedirs(os.path.dirname(dst) or '.', exist_ok=True)
    shutil.copy2(src, dst)
    return True
